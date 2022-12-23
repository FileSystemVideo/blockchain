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

/*
******************************************************************
**    ，client grpc  **
******************************************************************
 */


func RegisterCopyrightPartyHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyrightParty types.MsgRegisterCopyrightParty
	err := util.Json.Unmarshal(msgBytes, &copyrightParty)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	log.Debug("do")
	_, err = sdk.AccAddressFromBech32(copyrightParty.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	creator, err := sdk.AccAddressFromBech32(copyrightParty.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	
	enough, err := grpcQueryAccountSpace(ctx, 1, copyrightParty.Creator)
	if err != nil {
		log.WithError(err).Error("grpcQueryAccountSpace")
		return err
	}
	if !enough {
		log.Warn("space NotEnough")
		return errors.New(SpaceNotEnough)
	}

	if copyrightParty.Id == "" {
		log.Warn("BindId is empty")
		return errors.New(BindIdIsEmpty)
	}

	
	exists, err := grpcQueryCopyrightPartyExist(ctx, copyrightParty.Creator)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightPartyExist")
		return err
	}
	if !exists { 
		//return errors.New(NoRelationShip)
		//id
		has, err := grpcQueryPublisherExist(ctx, copyrightParty.Id)
		if err != nil {
			log.WithError(err).Error("grpcQueryPublisherExist")
			return err
		}
		if has {
			log.Warn("BindId Has Used")
			return errors.New(BindIdHasUsed)
		}
	}

	return judgeFee(ctx, creator, fee)
}
