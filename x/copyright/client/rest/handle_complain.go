package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"time"
)


func ComplainVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var complainVote types.MsgComplainVote
	err := util.Json.Unmarshal(msgBytes, &complainVote)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	log.Debug("do")
	voteAccount, err := sdk.AccAddressFromBech32(complainVote.VoteAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	complainInfor, err := grpcQueryCopyrightComplain(ctx, complainVote.ComplainId)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightComplain")
		return err
	}
	if complainInfor.DataHash == "" {
		log.Warn("CompainIdNotExist")
		return errors.New(CompainIdNotExist)
	}
	
	if complainInfor.ComplainStatus == "" {
		log.Warn("CompainIdNotExist")
		return errors.New(CompainIdNotExist)
	}
	if complainInfor.ComplainStatus != "4" {
		log.Warn("ComplainStatusValid")
		return errors.New(ComplainStatusValid)
	}

	
	dateTime := util.TimeStampToTime(complainInfor.ResponseTime)
	endTime := dateTime.Add(core.VoteResultTimePerioad)
	if endTime.Before(time.Now()) { 
		log.Warn("ComplainFinished")
		return errors.New(ComplainFinished)
	}

	
	has, err := grpcQueryDelegationVote(ctx, voteAccount)
	if err != nil {
		log.WithError(err).Error("grpcQueryDelegationVote error | ", err.Error())
		return err
	}
	if !has {
		log.Warn("AccountHasNoVoteRight")
		return errors.New(AccountHasNoVoteRight)
	}
	return judgeFee(ctx, voteAccount, fee)
}


func ComplainResponseHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var complainResponse types.MsgComplainResponse
	err := util.Json.Unmarshal(msgBytes, &complainResponse)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	log.Debug("do")

	complainInfor, err := grpcQueryCopyrightComplain(ctx, complainResponse.ComplainId)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightComplain")
		return err
	}
	if complainInfor.DataHash == "" {
		log.Warn("CompainIdNotExist")
		return errors.New(CompainIdNotExist)
	}

	if complainInfor.AccusedStatus != "0" {
		log.Warn("ComplainHasResponse")
		return errors.New(ComplainHasResponse)
	}
	
	if complainInfor.ComplainStatus == "2" {
		log.Warn("ComplainFinished")
		return errors.New(ComplainFinished)
	}
	if complainInfor.AccuseAccount.String() != complainResponse.AccuseAccount {
		log.Warn("CurrentAccountHasNoRight")
		return errors.New(CurrentAccountHasNoRight)
	}
	account, err := sdk.AccAddressFromBech32(complainResponse.AccuseAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return err
	}
	return judgeFee(ctx, account, fee)
}


func CopyrightComplainHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyrightComplain types.MsgCopyrightComplain
	err := util.Json.Unmarshal(msgBytes, &copyrightComplain)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	log.Debug("do")

	account, err := sdk.AccAddressFromBech32(copyrightComplain.ComplainAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyrightComplain.Datahash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyright")
		return err
	}
	if copyrightInfor.DataHash == "" {
		log.Warn("DataHashNotExist")
		return errors.New(DataHashNotExist)
	}
	
	exists, err := grpcQueryCopyrightPartyExist(ctx, copyrightComplain.ComplainAccount)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightPartyExist")
		return err
	}
	if !exists { 
		log.Warn("NoRelationShip")
		return errors.New(NoRelationShip)
	}

	//grpc
	enough, err := grpcQueryAccountSpace(ctx, copyrightInfor.Size, copyrightComplain.ComplainAccount)
	if err != nil {
		log.WithError(err).Error("grpcQueryAccountSpace")
		return err
	}
	if !enough {
		log.Warn("SpaceNotEnough")
		return errors.New(SpaceNotEnough)
	}
	
	/*balStatus, errStr := judgeBalance(ctx, complainAccount, copyrightInfor.PublishPrice, config.MainToken)
	if !balStatus {
		log.Error("", err)
		return errors.New(errStr)
	}*/
	return judgeFee(ctx, account, fee)
}
