package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

func CrossChainOutHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)

	var crosschainout types.MsgCrossChainOut
	err := util.Json.Unmarshal(msgBytes, &crosschainout)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	log.Debug("do")

	amountString, demon, err := util.StringDenom(crosschainout.Coins)
	if err != nil {
		log.WithError(err).Error("StringDenom")
		return err
	}
	amountDec, err := sdk.NewDecFromStr(amountString)
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return err
	}
	minAmount := types.MustParseLedgerDec2(sdk.MustNewDecFromStr(core.CrossChainOutMinAmount))
	log.Info(amountDec, " < ", minAmount)
	if amountDec.LT(minAmount) {
		log.Error("CrossChainOut Amount too low")
		return errors.New(CrossChainMinAmountErr)
	}
	addr, err := sdk.AccAddressFromBech32(crosschainout.SendAddress)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return err
	}
	balStatus, errStr := judgeBalance(ctx, addr, amountDec, demon)
	if !balStatus {
		log.Error("judgeBalance fail")
		return errors.New(errStr)
	}

	return nil
}

func CrossChainInHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}
