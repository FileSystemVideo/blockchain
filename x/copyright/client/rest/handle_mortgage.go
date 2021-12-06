package rest

import (
	"errors"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/shopspring/decimal"
	"strings"
)

func MortgageHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var mortgage types.MsgMortgage
	err := util.Json.Unmarshal(msgBytes, &mortgage)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(mortgage.MortageAccount)
	if err != nil {
		return errors.New(ParseAccountError)
	}

	err = judgeOfferAccount(mortgage.OfferAccountShare)
	if err != nil {
		return err
	}
	flag, err := grpcQueryMortgageAmount(ctx)
	if err != nil || !flag {
		return err
	}

	copyrightInfor, height, err := grpcQueryCopyright(ctx, mortgage.DataHash)
	if err != nil {
		return err
	}
	if copyrightInfor.Creator.Equals(addr) {
		return errors.New(HasPayedForDatahash)
	}
	if height < config.MortgageStartHeight {
		return errors.New(MortgageStartHeight)
	}

	data, err := grpcQueryCopyrightAndAccount(ctx, mortgage.MortageAccount, mortgage.DataHash)
	if err != nil {
		return errors.New(QueryChainInforError)
	}

	if data.Species == "buy" || (!data.Downer.Empty() && height < data.Height) {
		return errors.New(HasPayedForDatahash)
	}
	priceDecimal := copyrightInfor.Price.AmountDec()
	mortgAmountDecimal := priceDecimal.Mul(decimal.NewFromInt(config.MortgageRate))

	if len(fee.Amount) > 0 {
		feeDecimal, err := decimal.NewFromString(types.MustParseLedgerCoin(fee.Amount[0]))
		if err != nil {
			return errors.New(DecimalFromStringError)
		}
		if feeDecimal.LessThan(mortgAmountDecimal.Mul(decimal.RequireFromString(config.MortgageFee))) {
			return errors.New(FeeIsTooLess)
		}
		mortgAmountDecimal.Add(feeDecimal)
	} else {
		return errors.New(FeeCannotEmpty)
	}

	floatAmount, _ := mortgAmountDecimal.Float64()
	ledgeAmount := types.NewLedgerDec(floatAmount)
	balStatus, errStr := judgeBalance(ctx, addr, ledgeAmount, config.MainToken)
	if !balStatus {
		return errors.New(errStr)
	}

	return nil
}

func judgeOfferAccount(offerAccountShare string) error {
	accountArray := strings.Split(offerAccountShare, "-")
	size := len(accountArray)
	decimalTotal := decimal.RequireFromString("0")
	for i := 0; i < size; i++ {
		accountShare := strings.Split(accountArray[i], ":")
		if len(accountShare) != 2 {
			return errors.New(AcountRateError)
		}
		accountDecimal, err := decimal.NewFromString(accountShare[1])
		if err != nil {
			return errors.New(DecimalFromStringError)
		}
		decimalTotal = decimalTotal.Add(accountDecimal)
	}
	if !decimalTotal.Equal(decimal.RequireFromString("1000")) {
		return errors.New(AcountRateError)
	}
	return nil
}
