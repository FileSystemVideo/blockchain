package util

import (
	"fs.video/blockchain/core"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

//cion
func StringToCoinWithRate(amStr string) (sdk.Coin, error) {
	decCoin, err := sdk.ParseDecCoin(amStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	dec := decCoin.Amount.Mul(core.RealToLedgerRate)
	decInt := dec.TruncateInt()
	coin := sdk.NewCoin(decCoin.Denom, decInt)
	return coin, nil
}


//eg:1000MIP->1000,MIP
func StringDenom(amStr string) (string, string, error) {
	coin, err := StringToCoinWithRate(amStr)
	if err != nil {
		return "", "", err
	}
	amount := CoinToStringWithRateNoDenom(coin)
	return amount, coin.Denom, nil
}

//coin
func CoinToStringWithRateNoDenom(coin sdk.Coin) string {
	decCoin := sdk.NewDecCoin(coin.Denom, coin.Amount)
	dec := decCoin.Amount.Quo(core.RealToLedgerRate)
	decCoin = sdk.NewDecCoinFromDec(coin.Denom, dec)
	return stringBalance(decCoin.String(), false)
}


func stringBalance(balance string, MIP bool) string {
	if !strings.Contains(balance, ".") {
		return balance
	}
	if strings.Contains(balance, core.MainToken) {
		balance = balance[:len(balance)-len(core.MainToken)]
	}
	for i := len(balance) - 1; i > 0; i-- {
		if balance[i] == '.' {
			if MIP {
				return balance[:i] + core.MainToken
			} else {
				return balance[:i]
			}
		} else if balance[i] != '0' {
			if MIP {
				return balance[:i+1] + core.MainToken
			} else {
				return balance[:i+1]
			}
		}
	}
	return ""
}
