package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"time"
)

//版权申述投票的回调函数
func ComplainVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var complainVote types.MsgComplainVote
	err := util.Json.Unmarshal(msgBytes, &complainVote)
	if err != nil {
		return err
	}
	voteAccount, err := sdk.AccAddressFromBech32(complainVote.VoteAccount)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	complainInfor, err := grpcQueryCopyrightComplain(ctx, complainVote.ComplainId)
	if err != nil {
		return err
	}
	if complainInfor.DataHash == "" {
		return errors.New(CompainIdNotExist)
	}

	if complainInfor.ComplainStatus == "" {
		return errors.New(CompainIdNotExist)
	}
	if complainInfor.ComplainStatus != "4" {
		return errors.New(ComplainStatusValid)
	}


	dateTime := util.TimeStampToTime(complainInfor.ResponseTime)
	endTime := dateTime.Add(config.VoteResultTimePerioad)
	if endTime.Before(time.Now()) {
		return errors.New(ComplainFinished)
	}


	has, err := grpcQueryDelegationVote(ctx, voteAccount)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(AccountHasNoVoteRight)
	}
	return judgeFee(ctx, voteAccount, fee)
}


func ComplainResponseHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var complainResponse types.MsgComplainResponse
	err := util.Json.Unmarshal(msgBytes, &complainResponse)
	if err != nil {
		return err
	}

	complainInfor, err := grpcQueryCopyrightComplain(ctx, complainResponse.ComplainId)
	if err != nil {
		return err
	}
	if complainInfor.DataHash == "" {
		return errors.New(CompainIdNotExist)
	}

	if complainInfor.AccusedStatus != "0" {
		return errors.New(ComplainHasResponse)
	}

	if complainInfor.ComplainStatus == "2" {
		return errors.New(ComplainFinished)
	}
	if complainInfor.AccuseAccount.String() != complainResponse.AccuseAccount {
		return errors.New(CurrentAccountHasNoRight)
	}
	account, _ := sdk.AccAddressFromBech32(complainResponse.AccuseAccount)
	return judgeFee(ctx, account, fee)
}


func CopyrightComplainHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyrightComplain types.MsgCopyrightComplain
	err := util.Json.Unmarshal(msgBytes, &copyrightComplain)
	if err != nil {
		return err
	}

	account, err := sdk.AccAddressFromBech32(copyrightComplain.ComplainAccount)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyrightComplain.Datahash)
	if err != nil {
		return err
	}
	if copyrightInfor.DataHash == "" {
		return errors.New(DataHashNotExist)
	}

	exists, err := grpcQueryCopyrightPartyExist(ctx, copyrightComplain.ComplainAccount)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New(NoRelationShip)
	}


	enough, err := grpcQueryAccountSpace(ctx, copyrightInfor.Size, copyrightComplain.ComplainAccount)
	if err != nil {
		return err
	}
	if !enough {
		return errors.New(SpaceNotEnough)
	}


	return judgeFee(ctx, account, fee)
}
