package rest

import (
	"errors"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

func CopyrightVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyrightVote types.MsgVoteCopyright
	err := util.Json.Unmarshal(msgBytes, &copyrightVote)
	if err != nil {
		return err
	}

	voteAccount, err := sdk.AccAddressFromBech32(copyrightVote.Address)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	copyrightData, _, err := grpcQueryCopyright(ctx, copyrightVote.DataHash)
	if err != nil {
		return err
	}
	if copyrightData.DataHash == "" {
		return errors.New(DataHashNotExist)
	}
	if copyrightData.ApproveStatus != 0 {
		return errors.New(ApproveNotVote)
	}

	has, err := grpcQueryAccountVoteEnough(ctx, voteAccount, copyrightVote.Power)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(AccountHasNoVoteRight)
	}

	return judgeFee(ctx, voteAccount, fee)
}
