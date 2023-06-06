package keeper

import (
	"time"

	"github.com/terra-money/alliance/x/alliance/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RewardDelayTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.RewardDelayTime, &res)
	return
}

func (k Keeper) RewardClaimInterval(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.TakeRateClaimInterval, &res)
	return
}

func (k Keeper) LastRewardClaimTime(ctx sdk.Context) (res time.Time) {
	k.paramstore.Get(ctx, types.LastTakeRateClaimTime, &res)
	return
}

func (k Keeper) SetLastRewardClaimTime(ctx sdk.Context, lastTime time.Time) {
	k.paramstore.Set(ctx, types.LastTakeRateClaimTime, &lastTime)
}

func (k Keeper) TakeRateReceiver(ctx sdk.Context) (res sdk.AccAddress) {
	var str string
	k.paramstore.Get(ctx, types.TakeRateReceiver, &str)
	res, err := sdk.AccAddressFromBech32(str)
	if err != nil {
		panic(err)
	}

	return res
}

func (k Keeper) SetTakeRateReceiver(ctx sdk.Context, addr sdk.AccAddress) {
	str := addr.String()
	k.paramstore.Set(ctx, types.TakeRateReceiver, &str)
}
