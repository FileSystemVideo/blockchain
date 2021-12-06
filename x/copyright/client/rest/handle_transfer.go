package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
)


func TransferHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	logPrefix := logPrefix + " | " + util.GetFuncName() + " | "
	logs.Info(logPrefix)
	var transfer types.MsgTransfer
	err := util.Json.Unmarshal(msgBytes, &transfer)
	if err != nil {
		logs.Error(logPrefix, "Unmarshal error1 | ", err.Error())
		return err
	}
	fromAccount, err := sdk.AccAddressFromBech32(transfer.FromAddress)
	if err != nil {
		logs.Error(logPrefix, "AccAddressFromBech32 error1 | ", err.Error())
		return errors.New(ParseAccountError)
	}

	if transfer.ToAddress != config.ContractAddressDestory.String() {
		_, err = sdk.AccAddressFromBech32(transfer.ToAddress)
		if err != nil {
			logs.Error(logPrefix, "AccAddressFromBech32 error2 | ", err.Error())
			return errors.New(ParseAccountError)
		}
	}
	lockFlag := types2.JudgeLockedAccount(transfer.FromAddress)
	if lockFlag {
		return sdkerrors.ErrLockedAccount
	}
	var realCoins types.RealCoins
	err = util.Json.Unmarshal([]byte(transfer.Coins), &realCoins)
	if err != nil {
		logs.Error(logPrefix, "Unmarshal error2 | ", err.Error())
		return err
	}
	for i := 0; i < len(realCoins); i++ {
		flag := util.JudgeAmount(realCoins[i].Amount)
		if !flag {
			return errors.New(ParseCoinError)
		}
	}
	ledgeCoins := types.MustRealCoins2LedgerCoins(realCoins)
	feeCoins := types.NewLedgerCoins(config.CopyrightInviteFee)
	isFsv := false
	isTip := false
	for _, coin := range ledgeCoins {
		if coin.Denom == config.MainToken {
			isFsv = true
		}
		if coin.Denom == config.InviteToken {
			isTip = true
		}
	}
	for _, coin := range ledgeCoins {
		if coin.Denom == config.MainToken {
			minTransfer := types.NewLedgerCoin(config.MinFsvTransfer)
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
			logs.Error(logPrefix, "judgeBalance error | ", err.Error())
			return errors.New(errStr)
		}
	}

	if !isFsv {
		return errors.New(InvalidAmountErr)
	}
	return nil
}
