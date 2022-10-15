package keeper

import (
	"alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) Delegate(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin) (*types.MsgDelegateResponse, error) {
	asset := k.GetAssetByDenom(ctx, coin.Denom)
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, types.ModuleName, sdk.NewCoins(coin))
	tokensToMint := asset.RewardWeight.MulInt(coin.Amount).TruncateInt()
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.Coin{
		Denom:  k.stakingKeeper.BondDenom(ctx),
		Amount: tokensToMint,
	}))
	if err != nil {
		return nil, err
	}
	_, err = k.stakingKeeper.Delegate(ctx, moduleAddr, tokensToMint, stakingtypes.Unbonded, validator, true)
	if err != nil {
		return nil, err
	}
	asset.TotalTokens = asset.TotalTokens.Add(coin.Amount)
	k.setAsset(ctx, asset)
	k.upsertDelegationWithNewShares(ctx, delAddr, validator, coin, asset)
	return &types.MsgDelegateResponse{}, nil
}

func (k Keeper) upsertDelegationWithNewShares(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, coin sdk.Coin, asset types.AllianceAsset) {
	newShares := convertNewTokenToShares(asset.TotalTokens, asset.TotalShares, coin.Amount)
	delegation, ok := k.GetDelegation(ctx, delAddr, validator, coin.Denom)
	if !ok {
		delegation = types.Delegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: validator.GetOperator().String(),
			Denom:            coin.Denom,
			Shares:           newShares,
		}
	} else {
		delegation.Shares = delegation.Shares.Add(newShares)
	}
	k.SetDelegation(ctx, delAddr, validator, coin.Denom, delegation)
}

func (k Keeper) Undelegate() {
	panic("implement me")
}

func (k Keeper) Redelegate() {
	panic("implement me")
}

func (k Keeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, denom string) (d types.Delegation, found bool) {
	key := types.GetDelegationWithDenomKey(delAddr, validator.GetOperator(), denom)
	b := ctx.KVStore(k.storeKey).Get(key)
	if b == nil {
		return d, false
	}
	k.cdc.MustUnmarshal(b, &d)
	return d, true
}

func (k Keeper) SetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, validator stakingtypes.Validator, denom string, del types.Delegation) {
	key := types.GetDelegationWithDenomKey(delAddr, validator.GetOperator(), denom)
	b := k.cdc.MustMarshal(&del)
	ctx.KVStore(k.storeKey).Set(key, b)
	return
}

func convertNewTokenToShares(totalTokens sdk.Int, totalShares sdk.Dec, newTokens sdk.Int) (shares sdk.Dec) {
	if totalShares.IsZero() {
		return sdk.NewDecFromInt(newTokens)
	}
	return totalShares.MulInt(newTokens).QuoInt(totalTokens)
}
