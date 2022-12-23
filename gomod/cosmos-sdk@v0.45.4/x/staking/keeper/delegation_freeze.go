package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/json"
)

const DelegationFreezeKey = "delegation_freeze_"

func (k Keeper) DelegationFreeze(ctx sdk.Context, delAddr sdk.AccAddress, delegations sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	totalShares := sdk.NewDec(0)
	alldelegation := k.GetAllDelegatorDelegations(ctx, delAddr)
	for _, val := range alldelegation {
		totalShares = totalShares.Add(val.Shares)
	}
	if delegations.GT(totalShares) {
		return types.ErrNotEnoughDelegationShares
	}
	del := sdk.NewDec(0)
	if store.Has([]byte(DelegationFreezeKey + delAddr.String())) {
		bz := store.Get([]byte(DelegationFreezeKey + delAddr.String()))
		err := json.Unmarshal(bz, &del)
		if err != nil {
			return err
		}
		del = del.Add(delegations)
		if del.GT(totalShares) {
			return types.ErrNotEnoughDelegationShares
		}
		delByte, err := json.Marshal(del)
		if err != nil {
			return err
		}
		store.Set([]byte(DelegationFreezeKey+delAddr.String()), delByte)
	} else {
		delByte, err := json.Marshal(delegations)
		if err != nil {
			return err
		}
		store.Set([]byte(DelegationFreezeKey+delAddr.String()), delByte)
	}
	return nil
}

func (k Keeper) GetDelegationFreeze(ctx sdk.Context, delAddr sdk.AccAddress) (sdk.Dec, error) {
	store := ctx.KVStore(k.storeKey)
	del := sdk.NewDec(0)
	if store.Has([]byte(DelegationFreezeKey + delAddr.String())) {
		bz := store.Get([]byte(DelegationFreezeKey + delAddr.String()))
		err := json.Unmarshal(bz, &del)
		if err != nil {
			return del, err
		}
	}
	return del, nil
}

func (k Keeper) UnDelegationFreeze(ctx sdk.Context, delAddr sdk.AccAddress, delegations sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)
	del := sdk.NewDec(0)
	if store.Has([]byte(DelegationFreezeKey + delAddr.String())) {
		bz := store.Get([]byte(DelegationFreezeKey + delAddr.String()))
		err := json.Unmarshal(bz, &del)
		if err != nil {
			return err
		}
		if del.Equal(delegations) {
			store.Delete([]byte(DelegationFreezeKey + delAddr.String()))
		} else if del.GT(delegations) {
			del = del.Sub(delegations)
			delByte, err := json.Marshal(del)
			if err != nil {
				return err
			}
			store.Set([]byte(DelegationFreezeKey+delAddr.String()), delByte)
		}
	}
	return nil
}
