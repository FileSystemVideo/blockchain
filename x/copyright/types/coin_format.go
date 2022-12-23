package types

import (
	"fs.video/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/shopspring/decimal"
	"strings"
)

//int64 ()
func NewLedgerInt(realAmount decimal.Decimal) sdk.Int {
	realAmountDec, err := sdk.NewDecFromStr(realAmount.String())
	if err != nil {
		panic(err)
	}
	if realAmountDec.LT(core.MinRealAmountDec) {
		return sdk.NewInt(1)
	}
	return realAmountDec.Mul(core.RealToLedgerRate).TruncateInt()
}

//dec ()
func NewLedgerDec(realAmount decimal.Decimal) sdk.Dec {
	ledgerInt := NewLedgerInt(realAmount)
	return ledgerInt.ToDec()
}

//coin ()
func NewLedgerCoin(realAmount decimal.Decimal) sdk.Coin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoin(core.MainToken, ledgerInt)
}

//coin ()
func NewLedgerDecCoin(realAmount decimal.Decimal) sdk.DecCoin {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewDecCoin(core.MainToken, ledgerInt)
}

//coins ()
func NewLedgerCoins(realAmount decimal.Decimal) sdk.Coins {
	ledgerInt := NewLedgerInt(realAmount)
	return sdk.NewCoins(sdk.NewCoin(core.MainToken, ledgerInt))
}

//gasfee ()
func NewLedgerFeeFromGas(gas uint64, amount decimal.Decimal) legacytx.StdFee {
	ledgerInt := NewLedgerInt(amount)
	fee := legacytx.NewStdFee(gas, sdk.NewCoins(sdk.NewCoin(core.MainToken, ledgerInt)))
	return fee
}

//fee ()
func NewLedgerFee(amount decimal.Decimal) legacytx.StdFee {
	ledgerInt := NewLedgerInt(amount)
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(core.MainToken, ledgerInt)))
	return fee
}


func NewLedgerFeeZero() legacytx.StdFee {
	fee := legacytx.NewStdFee(flags.DefaultGasLimit, sdk.NewCoins(sdk.NewCoin(core.MainToken, sdk.ZeroInt())))
	return fee
}

//coins  ()
func MustParseLedgerCoins(ledgerCoins sdk.Coins) (realAmount string) {
	return MustParseLedgerCoin(ledgerCoins[0])
}

//fee  ()
func MustParseLedgerFee(ledgerFee legacytx.StdFee) (realAmount string) {
	return MustParseLedgerCoins(ledgerFee.Amount)
}

//Int  ()
func MustParseLedgerInt(ledgerInt sdk.Int) (realAmount string) {
	return RemoveStringLastZero(sdk.NewDecFromInt(ledgerInt).Mul(core.LedgerToRealRate).String())
}

//dec  ()
func MustParseLedgerDec(ledgerDec sdk.Dec) (realAmount string) {
	return RemoveStringLastZero(ledgerDec.Mul(core.LedgerToRealRate).String())
}

//dec  ()
func MustParseLedgerDec2(ledgerDec sdk.Dec) (realAmount sdk.Dec) {
	return ledgerDec.Mul(core.LedgerToRealRate)
}

// coin  ()
func MustParseLedgerCoinFromStr(ledgerCoinStr string) (realAmount string) {
	ledgerCoin, err := sdk.ParseCoinNormalized(ledgerCoinStr)
	if err != nil {
		panic(err)
	}
	ledgerAmount := ledgerCoin.Amount.ToDec()
	return RemoveStringLastZero(ledgerAmount.Mul(core.LedgerToRealRate).String())
}

//coin  ()
func MustParseLedgerCoin(ledgerCoin sdk.Coin) (realAmount string) {
	ledgerAmount := ledgerCoin.Amount.ToDec()
	return RemoveStringLastZero(ledgerAmount.Mul(core.LedgerToRealRate).String())
}

//deccoin  ()
func MustParseLedgerDecCoin(ledgerDecCoin sdk.DecCoin) (realAmount string) {
	ledgerAmount := ledgerDecCoin.Amount
	return RemoveStringLastZero(ledgerAmount.Mul(core.LedgerToRealRate).String())
}

//deccoins  ()
func MustParseLedgerDecCoins(ledgerDecCoins sdk.DecCoins) (realAmount string) {
	return MustParseLedgerDecCoin(ledgerDecCoins[0])
}

// LedgerCoin  RealCoin
func MustLedgerCoin2RealCoin(ledgerCoin sdk.Coin) (realCoin RealCoin) {
	ledgerAmount := ledgerCoin.Amount.ToDec()
	return RealCoin{
		Denom:  ledgerCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(core.LedgerToRealRate).String()),
	}
}

// LedgerDecCoin  RealCoin
func MustLedgerDecCoin2RealCoin(ledgerDecCoin sdk.DecCoin) (realCoin RealCoin) {
	ledgerAmount := ledgerDecCoin.Amount
	return RealCoin{
		Denom:  ledgerDecCoin.Denom,
		Amount: RemoveStringLastZero(ledgerAmount.Mul(core.LedgerToRealRate).String()),
	}
}

// LedgerCoins  RealCoins
func MustLedgerDecCoins2RealCoins(ledgerDecCoins sdk.DecCoins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerDecCoins); i++ {
		realCoins = append(realCoins, MustLedgerDecCoin2RealCoin(ledgerDecCoins[0]))
	}
	return
}

// LedgerCoins  RealCoins
func MustLedgerCoins2RealCoins(ledgerCoins sdk.Coins) (realCoins RealCoins) {
	for i := 0; i < len(ledgerCoins); i++ {
		realCoins = append(realCoins, MustLedgerCoin2RealCoin(ledgerCoins[i]))
	}
	return
}

// RealCoin  LedgerCoin
func MustRealCoin2LedgerDecCoin(realCoin RealCoin) (ledgerDecCoin sdk.DecCoin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	return sdk.NewDecCoinFromDec(realCoin.Denom, realCoinAmount.Mul(core.RealToLedgerRate))
}

// RealCoins  LedgerCoins
func MustRealCoins2LedgerDecCoins(realCoins RealCoins) (ledgerDecCoins sdk.DecCoins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerDecCoins = append(ledgerDecCoins, MustRealCoin2LedgerDecCoin(realCoins[0]))
	}
	return
}

//Int   ()
func MustRealInt2LedgerInt(ledgerInt sdk.Int) (realAmount sdk.Int) {
	realCoinAmount := sdk.NewDecFromInt(ledgerInt)
	return realCoinAmount.Mul(core.RealToLedgerRate).TruncateInt()
}

//Int   ()
func MustLedgerInt2RealInt(ledgerInt sdk.Int) string {
	realCoinAmount := sdk.NewDecFromInt(ledgerInt)
	realAmount := RemoveStringLastZero(realCoinAmount.Mul(core.LedgerToRealRate).String())
	return realAmount
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
// RealCoin  LedgerCo6in
func MustRealCoin2LedgerCoin(realCoin RealCoin) (ledgerCoin sdk.Coin) {
	realCoinAmount, err := sdk.NewDecFromStr(realCoin.Amount)
	if err != nil {
		panic(err)
	}
	return sdk.NewCoin(realCoin.Denom, realCoinAmount.Mul(core.RealToLedgerRate).TruncateInt())
}

// RealCoins  LedgerCoins
func MustRealCoins2LedgerCoins(realCoins RealCoins) (ledgerCoins sdk.Coins) {
	for i := 0; i < len(realCoins); i++ {
		ledgerCoins = append(ledgerCoins, MustRealCoin2LedgerCoin(realCoins[i])).Sort()
	}
	return
}

//coin
func NewRealCoinFromStr(denom string, amount string) RealCoin {
	return RealCoin{Denom: denom, Amount: amount}
}

//coins
func NewRealCoinsFromStr(denom string, amount string) RealCoins {
	return RealCoins{NewRealCoinFromStr(denom, amount)}
}

//coincoins
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
	if len(dataList[1]) > 18 {
		xiaoshu = xiaoshu[:18]
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
	if len(dataList[1]) > 18 {
		xiaoshu = xiaoshu[:18]
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
