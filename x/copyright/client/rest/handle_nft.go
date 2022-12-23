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

//nft 
func NftTransferHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var nftTransfer types.MsgNftTransfer
	err := util.Json.Unmarshal(msgBytes, &nftTransfer)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	log.Debug("do")

	account, err := sdk.AccAddressFromBech32(nftTransfer.From)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return errors.New(ParseAccountError)
	}
	_, err = sdk.AccAddressFromBech32(nftTransfer.To)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return errors.New(ParseAccountError)
	}

	/*feeCoin, err := sdk.NewDecFromStr(spaceMiner.Fee)
	if err != nil {
		return errors.New(ParseCoinError)
	}
	if feeCoin.IsPositive() {
		coin = coin.Add(feeCoin)
	}*/
	exists, err := grpcQueryCopyrightExist(ctx, nftTransfer.TokenId)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightExist")
		return err
	}
	if !exists { 
		log.Warn("DataHashNotExist")
		return errors.New(DataHashNotExist)
	}

	
	exists, err = grpcQueryCopyrightPartyExist(ctx, nftTransfer.To)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightPartyExist")
		return err
	}
	if !exists { 
		log.Warn("BindIdNotExist")
		return errors.New(BindIdNotExist)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, nftTransfer.TokenId)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyright")
		return err
	}
	
	enough, err := grpcQueryAccountSpace(ctx, copyrightInfor.Size, nftTransfer.To)
	if err != nil {
		log.WithError(err).Error("grpcQueryAccountSpace")
		return err
	}
	if !enough {
		log.Warn("SpaceNotEnough")
		return errors.New(SpaceNotEnough)
	}

	
	/*balStatus, errStr := judgeBalance(ctx, addr, coin, config.MainToken)
	if !balStatus {
		log.Error("", err)
		return errors.New(errStr)
	}*/

	return judgeFee(ctx, account, fee)
}
