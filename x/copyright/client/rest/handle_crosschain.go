package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	"fs.video/log"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

var logPrefix = "Rest"

func CrossChainOutHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	logPrefix := logPrefix + " | " + util.GetFuncName() + " | "
	log.Info(logPrefix)
	var crosschainout types.MsgCrossChainOut
	err := util.Json.Unmarshal(msgBytes, &crosschainout)
	if err != nil {
		log.Error(logPrefix, "Unmarshal error | "+err.Error())
		return err
	}
	amountString, demon, err := util.StringDenom(crosschainout.Coins)
	if err != nil {
		log.Error(logPrefix, "StringDenom error | "+err.Error())
		return err
	}
	amountDec, err := sdk.NewDecFromStr(amountString)
	if err != nil {
		log.Error(logPrefix, "NewDecFromStr error | "+err.Error())
		return err
	}
	minAmount := types.MustParseLedgerDec2(sdk.MustNewDecFromStr(config.CrossChainOutMinAmount))
	log.Info(amountDec, " < ", minAmount)
	if amountDec.LT(minAmount) {
		log.Error(logPrefix, "CrossChainOut Amount too low")
		return errors.New(CrossChainMinAmountErr)
	}
	addr, err := sdk.AccAddressFromBech32(crosschainout.SendAddress)
	if err != nil {
		log.Error(logPrefix, "AccAddressFromBech32 error | "+err.Error())
		return err
	}
	balStatus, errStr := judgeBalance(ctx, addr, amountDec, demon)
	if !balStatus {
		log.Error(logPrefix, "judgeBalance error | "+errStr)
		return errors.New(errStr)
	}

	return nil
}

func CrossChainInHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {

	return nil
}
