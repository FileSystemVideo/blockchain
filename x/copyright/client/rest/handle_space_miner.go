package rest

import (
	"errors"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

func SpaceMinerHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var spaceMiner types.MsgSpaceMiner
	err := util.Json.Unmarshal(msgBytes, &spaceMiner)
	if err != nil {
		return err
	}
	addr, err := sdk.AccAddressFromBech32(spaceMiner.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	var coin types.RealCoin
	err = util.Json.Unmarshal([]byte(spaceMiner.DeflationAmount), &coin)
	if err != nil {
		return errors.New(ParseCoinError)
	}
	if coin.Denom != sdk.DefaultBondDenom {
		return errors.New(MainTokenOnly)
	}
	flag := util.JudgeAmount(coin.Amount)
	if !flag {
		return errors.New(ParseCoinError)
	}

	var feeDec sdk.Coin
	if fee.Amount.Len() > 0 {
		for i := 0; i < fee.Amount.Len(); i++ {
			coin := fee.Amount[i]
			if coin.Denom == sdk.DefaultBondDenom {
				feeDec = coin
				break
			}
		}
	} else {
		return errors.New(FeeCannotEmpty)
	}
	if feeDec.IsZero() {
		return errors.New(FeeZero)
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(coin)
	ledgeCoin = ledgeCoin.Add(feeDec)

	balStatus, errStr := judgeBalance(ctx, addr, ledgeCoin.Amount.ToDec(), config.MainToken)
	if !balStatus {
		return errors.New(errStr)
	}

	return nil
}

func DeflationVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var spaceMiner types.MsgDeflationVote
	err := util.Json.Unmarshal(msgBytes, &spaceMiner)
	if err != nil {
		return err
	}

	account, err := sdk.AccAddressFromBech32(spaceMiner.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}

	has, err := grpcQueryDelegationVote(ctx, account)
	if err != nil {
		return err
	}
	if !has {
		return errors.New(AccountHasNoVoteRight)
	}
	has, err = grpcQueryDelegationVoteExist(ctx, account)
	if err != nil {
		return err
	}
	if has {
		return errors.New(HasDelationVoteError)
	}
	return judgeFee(ctx, account, fee)
}
