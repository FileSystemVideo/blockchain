package client

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	acc "fs.video/blockchain/util/account"
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strconv"
)

type SpaceMinerClient struct {
	TxClient  *TxClient
	ServerUrl string
	logPrefix string
}

func (this *SpaceMinerClient) QuerySpaceMinerInfor(account string) (data *keeper.DeflationAccountMinerInfor, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	accountSpaceMiner, blockNum, err := this.QueryAccountSpaceMinerInfor(account)
	if err != nil {
		return nil, 0, err
	}
	spaceMinerReward, err := this.QueryAccountSpaceMinerReward(account)
	if err != nil {
		return nil, 0, err
	}
	deflationMinerBaseInfor, err := this.QueryDeflationMinerInfor()
	if err != nil {
		return nil, 0, err
	}
	spaceTotal, err := this.QuerySpaceTotal()

	if err != nil {
		return nil, 0, err
	}
	destoryCoins, err := this.QuerySpaceAmount()
	if err != nil {
		return nil, 0, err
	}
	
	var defaltionMinerInfor keeper.DeflationAccountMinerInfor
	bytes, err := util.Json.Marshal(deflationMinerBaseInfor)
	log.WithField("data", string(bytes)).Debug("")
	err = util.Json.Unmarshal(bytes, &defaltionMinerInfor)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return nil, 0, err
	}
	defaltionMinerInfor.AccountSpaceTotal = accountSpaceMiner.SpaceTotal
	//+-?
	leftBlockNum := defaltionMinerInfor.MinerBlockNum + int64(core.SpaceMinerBonusBlockNum) - blockNum
	defaltionMinerInfor.LeftBlockNum = leftBlockNum
	if core.MinerStartHeight > blockNum {
		defaltionMinerInfor.LeftReachHeight = core.MinerStartHeight - blockNum
	} else {
		defaltionMinerInfor.DeflationStatus = 1
		defaltionMinerInfor.LeftReachHeight = 0
	}
	defaltionMinerInfor.SpaceTotal = spaceTotal
	spaceTotalDecimal := decimal.RequireFromString(spaceTotal)
	//accountSpace := decimal.RequireFromString("0")
	//if accountSpaceMiner.SpaceTotal != "" {
	//	accountSpace = decimal.RequireFromString(accountSpaceMiner.SpaceTotal)
	//}
	spaceAward, err := this.QuerySpaceAward()
	if err != nil {
		return nil, 0, err
	}
	//spaceAwardDecimal := decimal.RequireFromString(spaceAward)
	spaceAwardDec, err := sdk.NewDecFromStr(spaceAward)
	if err != nil {
		return nil, 0, err
	}
	percent := decimal.Decimal{}
	if spaceTotalDecimal.Sign() != 0 {
		spaceTotalDec, err := sdk.NewDecFromStr(spaceTotalDecimal.String())
		if err != nil {
			return nil, 0, err
		}
		percent = decimal.RequireFromString(spaceAwardDec.Quo(spaceTotalDec).String())
		defaltionMinerInfor.SpacePercent = accountSpaceMiner.SpaceTotal.Div(spaceTotalDecimal).StringFixed(2)
	}
	//bonusDecimal := k.QuerySpaceMinerBonusAmount(ctx)
	//accountExpectBonus := config.SpaceMinerPerDayStandand.Mul(percent).StringFixed(6)
	accountExpectBonus := util.DecimalStringFixed(accountSpaceMiner.SpaceTotal.Mul(percent).String(), core.CoinPlaces)
	defaltionMinerInfor.ExpectBonus = accountExpectBonus
	defaltionMinerInfor.DestoryAddress = core.ContractAddressDestory.String() 
	defaltionMinerInfor.DestoryAmount = destoryCoins.String()                 
	if spaceMinerReward.Amount == "" {
		defaltionMinerInfor.SpaceMinerReward = "0"
	} else {
		defaltionMinerInfor.SpaceMinerReward = spaceMinerReward.Amount 
	}
	return &defaltionMinerInfor, blockNum, err
}


func (this *SpaceMinerClient) QueryAccountSpaceMinerReward(account string) (data *types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryAccountSpaceMinerParams{Account: account}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QuerySpaceMinerRewardInfo, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, err
	}
	data = &types.RealCoin{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *SpaceMinerClient) QueryAccountSpaceMinerInfor(account string) (data *keeper.AccountSpaceMiner, blockNum int64, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := types.QueryAccountSpaceMinerParams{Account: account}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, 0, err
	}
	resBytes, blockNum, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryAccountSpaceMiner, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return nil, 0, err
	}
	data = &keeper.AccountSpaceMiner{}
	err = util.Json.Unmarshal(resBytes, data)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return
}


func (this *SpaceMinerClient) QueryDeflationRate() (*keeper.RateAndVoteIndex, int64, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	rateVoteIndex := &keeper.RateAndVoteIndex{}
	bz := []byte("")
	resBytes, height, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDeflationRateInfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return rateVoteIndex, height, err
	}
	err = util.Json.Unmarshal(resBytes, rateVoteIndex)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return rateVoteIndex, height, err
}


func (this *SpaceMinerClient) QuerySpaceAmount() (*sdk.Dec, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	spaceAmount := &sdk.Dec{}
	bz := []byte("")
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QuerySpaceAmount, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return spaceAmount, err
	}
	err = util.Json.Unmarshal(resBytes, spaceAmount)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return spaceAmount, err
}


func (this *SpaceMinerClient) QuerySpaceTotal() (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	bz := []byte("")

	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryTotalSpaceInfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return "", err
	}
	return string(resBytes), err
}


func (this *SpaceMinerClient) QuerySpaceAward() (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	bz := []byte("")
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QuerySpaceAward, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return "", err
	}
	return string(resBytes), err
}


func (this *SpaceMinerClient) QueryHasMinerBonus(height int64) (string, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	bz := []byte("")
	clientCtx = clientCtx.WithHeight(height)
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryHasMinerBonus, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return "", err
	}
	return string(resBytes), err
}


func (this *SpaceMinerClient) QueryDeflationMinerInfor() (*keeper.DeflationMinerInfor, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	deflationMinerInfor := &keeper.DeflationMinerInfor{}
	bz := []byte("")
	resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryDeflationMinerInfor, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return deflationMinerInfor, err
	}
	err = util.Json.Unmarshal(resBytes, deflationMinerInfor)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return deflationMinerInfor, err
}


func (this *SpaceMinerClient) Create(data types.SpaceMinerData, privateKey string, gas uint64) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgSpaceMiner(data)
	if err != nil {
		log.WithError(err).Error("NewMsgSpaceMiner")
		return
	}
	//data
	data.Fee.Gas = gas
	resp, err = this.TxClient.SignAndSendMsg(acc.AccAddressToString(data.Creator), privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *SpaceMinerClient) SpaceMinerReward(data types.SpaceMinerRewardData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgSpaceMinerReward(data)
	if err != nil {
		log.WithError(err).Error("NewMsgSpaceMinerReward")
		return
	}
	
	resp, err = this.TxClient.SignAndSendMsg(data.Address, privateKey, data.Fee, "", msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}


func (this *SpaceMinerClient) GasInfo(data types.SpaceMinerData) (coin types.RealCoin, gas string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgSpaceMiner(data)
	if err != nil {
		log.WithError(err).Error("NewMsgSpaceMiner")
		return
	}
	seqDetail, err := this.TxClient.FindAccountNumberSeq(msg.Creator)
	if err != nil {
		return
	}
	clientFactory = clientFactory.WithSequence(seqDetail.Sequence)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg)
	if err != nil {
		log.WithError(err).Error("CalculateGas")
		return
	}
	gasUsed := decimal.RequireFromString(strconv.Itoa(int(gasInfo.GasInfo.GasUsed)))
	gasUsed = gasUsed.Add(decimal.RequireFromString("2000000"))
	gas = gasUsed.String()
	gasUsed = gasUsed.Mul(core.MinimumGasPrices).Add(decimal.NewFromFloat(1))
	gasDec, err := sdk.NewDecFromStr(gasUsed.StringFixed(6))
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return
	}
	gasDecCoin := sdk.NewDecCoinFromDec(core.MainToken, gasDec)
	amount := types.MustParseLedgerDecCoin(gasDecCoin)
	if decimal.RequireFromString(amount).LessThan(core.ChainDefaultFee) {
		
		amount = core.ChainDefaultFee.String()
	}
	coin = types.NewRealCoinFromStr(core.MainToken, amount)
	return
}
