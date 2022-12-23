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


func TransferHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)

	var transfer types.MsgTransfer
	err := util.Json.Unmarshal(msgBytes, &transfer)
	if err != nil {
		log.WithError(err).Error("Unmarshal1")
		return err
	}

	log.Debug("do")

	fromAccount, err := sdk.AccAddressFromBech32(transfer.FromAddress)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return errors.New(ParseAccountError)
	}

	if transfer.ToAddress != core.ContractAddressDestory.String() {
		_, err = sdk.AccAddressFromBech32(transfer.ToAddress)
		if err != nil {
			log.WithError(err).Error("AccAddressFromBech32 2")
			return errors.New(ParseAccountError)
		}
	}

	var realCoins types.RealCoins
	err = util.Json.Unmarshal([]byte(transfer.Coins), &realCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal2")
		return err
	}
	for i := 0; i < len(realCoins); i++ {
		flag := util.JudgeAmount(realCoins[i].Amount)
		if !flag {
			return errors.New(ParseCoinError)
		}
	}
	ledgeCoins := types.MustRealCoins2LedgerCoins(realCoins)
	feeCoins := types.NewLedgerCoins(core.CopyrightInviteFee)
	isFsv := false
	isTip := false
	for _, coin := range ledgeCoins {
		if coin.Denom == core.MainToken {
			isFsv = true
		}
		if coin.Denom == core.InviteToken {
			isTip = true
		}
	}
	for _, coin := range ledgeCoins {
		if coin.Denom == core.MainToken {
			minTransfer := types.NewLedgerCoin(core.MinFsvTransfer)
			if coin.IsLT(minTransfer) && !coin.IsEqual(feeCoins[0]) {
				return errors.New(InvalidAmountErr)
			}
			if !isTip {
				if coin.IsLT(minTransfer) {
					return errors.New(InvalidAmountErr)
				}
			}
		}
		
		balStatus, errStr := judgeBalance(ctx, fromAccount, coin.Amount.ToDec(), coin.Denom)
		if !balStatus {
			log.WithField("err", errStr).Error("judgeBalance fail")
			return errors.New(errStr)
		}
	}
	//fsv
	if !isFsv {
		log.Warn("InvalidAmountErr")
		return errors.New(InvalidAmountErr)
	}
	return nil
}
