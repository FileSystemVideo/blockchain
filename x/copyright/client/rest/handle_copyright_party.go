package rest

import (
	"errors"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

func RegisterCopyrightPartyHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyrightParty types.MsgRegisterCopyrightParty
	err := util.Json.Unmarshal(msgBytes, &copyrightParty)
	if err != nil {
		return err
	}

	_, err = sdk.AccAddressFromBech32(copyrightParty.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}

	creator, err := sdk.AccAddressFromBech32(copyrightParty.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}

	enough, err := grpcQueryAccountSpace(ctx, 1, copyrightParty.Creator)
	if err != nil {
		return err
	}
	if !enough {
		return errors.New(SpaceNotEnough)
	}

	if copyrightParty.Id == "" {
		return errors.New(BindIdIsEmpty)
	}

	exists, err := grpcQueryCopyrightPartyExist(ctx, copyrightParty.Creator)
	if err != nil {
		return err
	}
	if !exists {

		has, err := grpcQueryPublisherExist(ctx, copyrightParty.Id)
		if err != nil {
			return err
		}
		if has {
			return errors.New(BindIdHasUsed)
		}
	}

	return judgeFee(ctx, creator, fee)
}
