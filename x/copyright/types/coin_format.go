package types

import (
	"fs.video/blockchain/x/copyright/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"strings"
)

func NewLedgerInt(realAmount float64) int64 {
	if realAmount < config.MinRealAmountFloat64 {
		return int64(config.MinRealAmountFloat64 * config.RealToLedgerRate)
	}
	return int64(realAmount * config.RealToLedgerRate)
}

func NewLedgerDec(realAmount float64) sdk.Dec {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewDec(ledgerInt)
}

func NewLedgerCoin(realAmount float64) sdk.Coin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoin(config.MainToken, sdk.NewInt(ledgerInt))
}

func NewLedgerDecCoin(realAmount float64) sdk.DecCoin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewDecCoin(config.MainToken, sdk.NewInt(ledgerInt))
}

func NewLedgerCoins(realAmount float64) sdk.Coins {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoins(sdk.NewCoin(config.MainToken, sdk.NewInt(ledgerInt)))
}

func NewLedgerFee(amount float64) legacytx.StdFee {
	ledgerInt := NewLedgerInt(amount)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewInt64Coin(config.MainToken, ledgerInt)))
	return fee
}

func NewLedgerFeeZero() legacytx.StdFee {
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewInt64Coin(config.MainToken, 0)))
	return fee
}

func MustParseLedgerCoins(ledgerCoins sdk.Coins) (realAmount string) {
	return MustParseLedgerCoin(ledgerCoins[0])
}

func MustParseLedgerFee(ledgerFee legacytx.StdFee) (realAmount string) {
	return MustParseLedgerCoins(ledgerFee.Amount)
}

func MustParseLedgerInt(ledgerInt sdk.Int) (realAmount string) {
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(sdk.NewDecFromInt(ledgerInt).Mul(rate).String())
}

func MustParseLedgerDec(ledgerDec sdk.Dec) (realAmount string) {
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerDec.Mul(rate).String())
}

func MustParseLedgerDec2(ledgerDec sdk.Dec) (realAmount sdk.Dec) {
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return ledgerDec.Mul(rate)
}

func MustParseLedgerCoinFromStr(ledgerCoinStr string) (realAmount string) {
	ledgerCoin, err := sdk.ParseCoinNormalized(ledgerCoinStr)
	if err != nil {
		panic(err)
	}
	ledgerAmount := ledgerCoin.Amount.ToDec()
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerCoin(ledgerCoin sdk.Coin) (realAmount string) {
	ledgerAmount := ledgerCoin.Amount.ToDec()
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerDecCoin(ledgerDecCoin sdk.DecCoin) (realAmount string) {
	ledgerAmount := ledgerDecCoin.Amount
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RemoveStringLastZero(ledgerAmount.Mul(rate).String())
}

func MustParseLedgerDecCoins(ledgerDecCoins sdk.DecCoins) (realAmount string) {
	return MustParseLedgerDecCoin(ledgerDecCoins[0])
}

func MustLedgerCoin2RealCoin(ledgerCoin sdk.Coin) (realCoin RealCoin) {
	ledgerAmount := ledgerCoin.Amount.ToDec()
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RealCoin{
		Denom:  ledgerCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(rate).String()),
	}
}

func MustLedgerDecCoin2RealCoin(ledgerDecCoin sdk.DecCoin) (realCoin RealCoin) {
	ledgerAmount := ledgerDecCoin.Amount
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return RealCoin{
		Denom:  ledgerDecCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(rate).String()),
	}
}

func MustLedgerDecCoins2RealCoins(ledgerDecCoins sdk.DecCoins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerDecCoins); i++ {
		realCoins = append(realCoins, MustLedgerDecCoin2RealCoin(ledgerDecCoins[0]))
	}
	return
}

func MustLedgerCoins2RealCoins(ledgerCoins sdk.Coins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerCoins); i++ {
		realCoins = append(realCoins, MustLedgerCoin2RealCoin(ledgerCoins[i]))
	}
	return
}

func MustRealCoin2LedgerDecCoin(realCoin RealCoin) (ledgerDecCoin sdk.DecCoin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	rate := sdk.NewDec(config.RealToLedgerRateInt64)
	return sdk.NewDecCoinFromDec(realCoin.Denom, realCoinAmount.Mul(rate))
}

func MustRealCoins2LedgerDecCoins(realCoins RealCoins) (ledgerDecCoins sdk.DecCoins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerDecCoins = append(ledgerDecCoins, MustRealCoin2LedgerDecCoin(realCoins[0]))
	}
	return
}

func MustRealInt2LedgerInt(ledgerInt sdk.Int) (realAmount sdk.Int) {
	realCoinAmount := sdk.NewDecFromInt(ledgerInt)
	rate := sdk.NewDec(config.RealToLedgerRateInt64)
	return realCoinAmount.Mul(rate).TruncateInt()
}

func MustLedgerInt2RealInt(ledgerInt sdk.Int) (realAmount sdk.Int) {
	realCoinAmount := sdk.NewDecFromInt(ledgerInt)
	rate, err := sdk.NewDecFromStr(config.LedgerToRealRate)
	if err != nil {
		panic(err)
	}
	return realCoinAmount.Mul(rate).TruncateInt()
}

/*
RealCoin 2 Coin
RealCoins 2 Coins
RealCoins 2 DecCoins
RealCoins 2 DecCoin
DecCoin 2 RealCoin
DecCoins 2 RealCoins
Coin 2 RealCoin
Coins 2 RealCoins
*/
func MustRealCoin2LedgerCoin(realCoin RealCoin) (ledgerCoin sdk.Coin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	rate := sdk.NewDec(config.RealToLedgerRateInt64)
	return sdk.NewCoin(realCoin.Denom, realCoinAmount.Mul(rate).TruncateInt())
}

func MustRealCoins2LedgerCoins(realCoins RealCoins) (ledgerCoins sdk.Coins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerCoins = append(ledgerCoins, MustRealCoin2LedgerCoin(realCoins[i])).Sort()
	}
	return
}

func NewRealCoinFromStr(denom string, amount string) RealCoin {
	return RealCoin{Denom: denom, Amount: amount}
}

func NewRealCoinsFromStr(denom string, amount string) RealCoins {
	return RealCoins{NewRealCoinFromStr(denom, amount)}
}

func NewRealCoins(realCoin RealCoin) RealCoins {
	return RealCoins{realCoin}
}

func RemoveStringLastZero(balance string) string {
	if !strings.Contains(balance, ".") {
		return balance
	}
	dataList := strings.Split(balance, ".")
	zhengshu := dataList[0]
	xiaoshu := dataList[1]
	if len(dataList[1]) > 6 {
		xiaoshu = xiaoshu[:6]
	}
	xiaoshu2 := ""
	for i := len(xiaoshu) - 1; i >= 0; i-- {
		if xiaoshu[i] != '0' {
			xiaoshu2 = xiaoshu[:i+1]
			break
		}
	}
	if xiaoshu2 == "" {
		return zhengshu
	} else {
		return zhengshu + "." + xiaoshu2
	}
}

func RemoveDecLastZero(amount sdk.Dec) string {
	balance := amount.String()
	if !strings.Contains(balance, ".") {
		return balance
	}
	dataList := strings.Split(balance, ".")
	zhengshu := dataList[0]
	xiaoshu := dataList[1]
	if len(dataList[1]) > 6 {
		xiaoshu = xiaoshu[:6]
	}
	xiaoshu2 := ""
	//fmt.Println("xiaoshu:",xiaoshu)
	for i := len(xiaoshu) - 1; i >= 0; i-- {
		if xiaoshu[i] != '0' {
			xiaoshu2 = xiaoshu[:i+1]
			break
		}
	}
	if xiaoshu2 == "" {
		return zhengshu
	} else {
		return zhengshu + "." + xiaoshu2
	}
}
