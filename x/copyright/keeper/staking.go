package keeper

import (
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k Keeper) GetAccountDelegatorShares(ctx sdk.Context, account sdk.AccAddress) (string, string) {
	delegations := k.stakingKeeper.GetAllDelegatorDelegations(ctx, account)
	totalShares := sdk.NewDec(0)
	totalBalance := sdk.NewDec(0)
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			continue
		}
		v, found := k.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			continue
		}
		totalBalance = totalBalance.Add(v.TokensFromShares(del.Shares))
		totalShares = totalShares.Add(del.Shares)
	}
	totalSharesStr := types.MustParseLedgerDec(totalShares)
	totalBalanceStr := types.MustParseLedgerDec(totalBalance)
	return totalSharesStr, totalBalanceStr
}


func (k Keeper) GetAllDelegatorShares(ctx sdk.Context) string {
	validators := k.stakingKeeper.GetAllValidators(ctx)
	total := sdk.NewDec(0)
	for _, val := range validators {
		total = total.Add(val.DelegatorShares)
	}
	totalString := types.MustParseLedgerDec(total)
	return totalString
}
