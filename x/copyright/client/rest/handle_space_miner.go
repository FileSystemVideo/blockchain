package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/shopspring/decimal"
)

//fsv
func SpaceMinerHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var spaceMiner types.MsgSpaceMiner
	err := util.Json.Unmarshal(msgBytes, &spaceMiner)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	log.Debug("do")

	addr, err := sdk.AccAddressFromBech32(spaceMiner.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	var coin types.RealCoin
	err = util.Json.Unmarshal([]byte(spaceMiner.DeflationAmount), &coin)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return errors.New(ParseCoinError)
	}
	if coin.Denom != sdk.DefaultBondDenom {
		log.Warn("MainTokenOnly")
		return errors.New(MainTokenOnly)
	}
	flag := util.JudgeAmount(coin.Amount)
	if !flag {
		log.Warn("JudgeAmount fail")
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
		log.Warn("FeeCannotEmpty")
		return errors.New(FeeCannotEmpty)
	}
	if feeDec.IsZero() {
		log.Warn("FeeZero")
		return errors.New(FeeZero)
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(coin)
	ledgeCoin = ledgeCoin.Add(feeDec)
	
	balStatus, errStr := judgeBalance(ctx, addr, ledgeCoin.Amount.ToDec(), core.MainToken)
	if !balStatus {
		log.WithField("err", errStr).Error("judgeBalance fail")
		return errors.New(errStr)
	}
	return nil
}


func SpaceMinerBonusHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var spaceMiner types.MsgSpaceMinerReward
	err := util.Json.Unmarshal(msgBytes, &spaceMiner)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	log.Debug("do")

	_, err = sdk.AccAddressFromBech32(spaceMiner.Address)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	bonusAmount, err := grpcQueryMinerBonus(ctx, spaceMiner.Address)
	if err != nil {
		return err
	}
	bonusDecimal, err := decimal.NewFromString(bonusAmount)
	if bonusDecimal.LessThanOrEqual(decimal.Zero) {
		return errors.New(NoMinerBonusClain)
	}
	return nil
}


//func DeflationVoteHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
//	log := util.BuildLog(util.GetFuncName(), util.LmChainRest)
//	var spaceMiner types.MsgDeflationVote
//	err := util.Json.Unmarshal(msgBytes, &spaceMiner)
//	if err != nil {
//		log.WithError(err).Error("Unmarshal")
//		return err
//	}

//	log.Debug("do")

//	account, err := sdk.AccAddressFromBech32(spaceMiner.Creator)
//	if err != nil {
//		log.WithError(err).Error("AccAddressFromBech32")
//		return errors.New(ParseAccountError)
//	}

//	
//	has, err := grpcQueryDelegationVote(ctx, account)
//	if err != nil {
//		log.WithError(err).Error("grpcQueryDelegationVote")
//		return err
//	}
//	if !has {
//		log.Warn("AccountHasNoVoteRight")
//		return errors.New(AccountHasNoVoteRight)
//	}
//	has, err = grpcQueryDelegationVoteExist(ctx, account)
//	if err != nil {
//		log.WithError(err).Error("grpcQueryDelegationVoteExist")
//		return err
//	}
//	if has {
//		log.Warn("HasDelationVoteError")
//		return errors.New(HasDelationVoteError)
//	}
//	return judgeFee(ctx, account, fee)
//}
