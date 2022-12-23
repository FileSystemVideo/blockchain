package keeper

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
)

const (
	deflation_rate     = "1"        // 20fsv：1M
	deflation_rate_min = "0.000001" // 0.000001
)

/**

*/
func (k Keeper) QueryDeflationRate(ctx sdk.Context) string {
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	if spaceTotal != "0" {
		return k.calculatDeflationRate(ctx)
	} else {
		return deflation_rate
	}
}


func (k Keeper) calculatDeflationRate(ctx sdk.Context) string {
	bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
	spaceTotal := k.QueryDeflatinSpaceTotal(ctx)
	spaceTotalDecimal := decimal.RequireFromString(spaceTotal)
	totalSpaceDecimalM := spaceTotalDecimal.Div(ByteToMb)
	//M 
	preMbonus := bonusDecimal.Div(totalSpaceDecimalM)
	currentRate := preMbonus.Mul(returnDays)
	if currentRate.GreaterThan(decimal.RequireFromString(deflation_rate_min)) {
		return util.DecimalStringFixed(currentRate.String(), core.CoinPlaces)
	} else {
		return deflation_rate_min
	}

}
