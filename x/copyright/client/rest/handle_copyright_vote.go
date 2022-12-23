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


func CopyrightVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyrightVote types.MsgVoteCopyright
	err := util.Json.Unmarshal(msgBytes, &copyrightVote)
	if err != nil {
		log.Error("Unmarshal")
		return err
	}
	log.Debug("do")

	/*flag := types2.JudgeLockedAccount(copyrightVote.Address)
	if flag{
		return sdkerrors.ErrLockedAccount
	}*/

	voteAccount, err := sdk.AccAddressFromBech32(copyrightVote.Address)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	copyrightData, _, err := grpcQueryCopyright(ctx, copyrightVote.DataHash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyright")
		return err
	}
	if copyrightData.DataHash == "" { 
		log.Warn("DataHashNotExist")
		return errors.New(DataHashNotExist)
	}
	if copyrightData.ApproveStatus != 0 { 
		log.Warn("ApproveNotVote")
		return errors.New(ApproveNotVote)
	}
	
	has, err := grpcQueryAccountVoteEnough(ctx, voteAccount, copyrightVote.Power)
	if err != nil {
		log.WithError(err).Error("grpcQueryAccountVoteEnough")
		return err
	}
	if !has {
		log.Warn("AccountHasNoVoteRight")
		return errors.New(AccountHasNoVoteRight)
	}

	return judgeFee(ctx, voteAccount, fee)
}
