package util

import (
	"fs.video/blockchain/x/copyright/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)


func StringToCoinWithRate(amStr string) (sdk.Coin, error) {
	decCoin, err := sdk.ParseDecCoin(amStr)
	if err != nil {
		return sdk.Coin{}, err
	}
	dec := decCoin.Amount.MulInt(sdk.NewInt(config.RealToLedgerRateInt64))
	decInt := dec.TruncateInt()
	coin := sdk.NewCoin(decCoin.Denom, decInt)
	return coin, nil
}



func StringDenom(amStr string) (string, string, error) {
	coin, err := StringToCoinWithRate(amStr)
	if err != nil {
		return "", "", err
	}
	amount := CoinToStringWithRateNoDenom(coin)
	return amount, coin.Denom, nil
}


func CoinToStringWithRateNoDenom(coin sdk.Coin) string {
	decCoin := sdk.NewDecCoin(coin.Denom, coin.Amount)
	dec := decCoin.Amount.QuoInt(sdk.NewInt(config.RealToLedgerRateInt64))
	decCoin = sdk.NewDecCoinFromDec(coin.Denom, dec)
	return stringBalance(decCoin.String(), false)
}


func stringBalance(balance string, MIP bool) string {
	if !strings.Contains(balance, ".") {
		return balance
	}
	if strings.Contains(balance, config.MainToken) {
		balance = balance[:len(balance)-len(config.MainToken)]
	}
	for i := len(balance) - 1; i > 0; i-- {
		if balance[i] == '.' {
			if MIP {
				return balance[:i] + config.MainToken
			} else {
				return balance[:i]
			}
		} else if balance[i] != '0' {
			if MIP {
				return balance[:i+1] + config.MainToken
			} else {
				return balance[:i+1]
			}
		}
	}
	return ""
}
