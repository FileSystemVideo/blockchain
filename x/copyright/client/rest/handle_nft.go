package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)


func NftTransferHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var nftTransfer types.MsgNftTransfer
	err := util.Json.Unmarshal(msgBytes, &nftTransfer)
	if err != nil {
		return err
	}

	account, err := sdk.AccAddressFromBech32(nftTransfer.From)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	_, err = sdk.AccAddressFromBech32(nftTransfer.To)
	if err != nil {
		return errors.New(ParseAccountError)
	}


	exists, err := grpcQueryCopyrightExist(ctx, nftTransfer.TokenId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(DataHashNotExist)
	}


	exists, err = grpcQueryCopyrightPartyExist(ctx, nftTransfer.To)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(BindIdNotExist)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, nftTransfer.TokenId)
	if err != nil {
		return err
	}

	enough, err := grpcQueryAccountSpace(ctx, copyrightInfor.Size, nftTransfer.To)
	if err != nil {
		return err
	}
	if !enough {
		return errors.New(SpaceNotEnough)
	}




	return judgeFee(ctx, account, fee)
}
