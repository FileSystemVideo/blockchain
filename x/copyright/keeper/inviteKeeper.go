package keeper

import (
	"encoding/json"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	InviteRelationKey = "inviteRelation_" //key

	InviteRewardKey = "inviteReward" //key

	InviteRecordingKey = "inviteRecording" //key

	InvitesStatisticsKey = "invitesStatistics" //key

	InviteFirstRateKey = "0.1" 

	InviteRateKey = "0.8" 

	InviteSpaceRateKey = "3" 
)


func (k Keeper) QueryRewardInfo(ctx sdk.Context, account string) (*types.Settlement, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	accountReward := new(types.Settlement)
	if !store.Has(InviteRewardKey + strings.ToLower(account)) {
		return accountReward, nil
	}
	
	err := store.GetUnmarshal(InviteRewardKey+strings.ToLower(account), &accountReward)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return accountReward, err
	}
	return accountReward, nil

}


func (k Keeper) IsInvite(ctx sdk.Context, address sdk.AccAddress) bool {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	log.WithField("Invitation relationship", string(store.Get(InviteRelationKey+address.String()))).Debug("IsInvite")
	return store.Has(InviteRelationKey + address.String())
}

// inviteAddress :  address:
func (k Keeper) InviteRelation(ctx sdk.Context, inviteAddress sdk.AccAddress, address sdk.AccAddress) error {
	store := k.KVHelper(ctx)
	key := InviteRelationKey + address.String()
	
	if inviteAddress.String() == "" {
		return nil
	}
	
	if store.Has(key) {
		return nil
	}
	
	return store.Set(InviteRelationKey+address.String(), inviteAddress.String())
}


func (k Keeper) InviteRecording(ctx sdk.Context, inviteAddress sdk.AccAddress, address sdk.AccAddress, inviteTime int64) error {
	store := k.KVHelper(ctx)
	var inviteRecording []types.InviteRecording
	inviteKey := InviteRecordingKey + inviteAddress.String()
	if store.Has(inviteKey) {
		bz := store.Get(inviteKey)
		err := json.Unmarshal(bz, &inviteRecording)
		if err != nil {
			return err
		}
	}
	var record = types.InviteRecording{
		Address:    address.String(),
		InviteTime: inviteTime,
	}
	inviteRecording = append(inviteRecording, record)
	recordByte, err := json.Marshal(inviteRecording)
	if err != nil {
		return err
	}
	return store.Set(inviteKey, recordByte)
}


func (k Keeper) GetInviteRecording(ctx sdk.Context, inviteAddress sdk.Address) ([]types.InviteRecording, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	key := InviteRecordingKey + inviteAddress.String()
	if !store.Has(key) {
		return nil, nil
	}
	var recording []types.InviteRecording
	err := store.GetUnmarshal(key, &recording)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return nil, err
	}
	
	for index, record := range recording {
		accountMiner := k.QueryAccountSpaceMinerInfor(ctx, record.Address)
		record.Space = accountMiner.BuySpace
		recording[index] = record
	}
	return recording, nil
}


func (k Keeper) GetInviteRewardStatistics(ctx sdk.Context, address sdk.AccAddress) (types.InviteRewardStatistics, error) {
	store := k.KVHelper(ctx)
	var statistics = types.InviteRewardStatistics{}
	if !store.Has(InvitesStatisticsKey + address.String()) {
		return statistics, nil
	}
	err := store.GetUnmarshal(InvitesStatisticsKey+address.String(), &statistics)
	if err != nil {
		return statistics, err
	}
	return statistics, nil
}


func (k Keeper) RewardSettlement(ctx sdk.Context, account string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	if !store.Has(InviteRewardKey + account) {
		log.WithField("account", account).Info("There are no rewards to settle")
		return nil
	}
	
	var accountReward types.Settlement
	err := store.GetUnmarshal(InviteRewardKey+account, &accountReward)
	if err != nil {
		log.WithError(err).Error("GetUnmarshal")
		return err
	}

	//2、 key->  val->
	var statistics = types.InviteRewardStatistics{}
	if store.Has(InvitesStatisticsKey + account) {
		err := store.GetUnmarshal(InvitesStatisticsKey+account, &statistics)
		if err != nil {
			log.WithError(err).Error("GetUnmarshal")
			return err
		}
	}
	

	
	err = k.updataSpaceMinerRewardSpace(ctx, account, accountReward, statistics)
	if err != nil {
		log.WithError(err).Error("updataSpaceMinerRewardSpace")
		return err
	}

	
	/*delete(rewardMap, account)
	rewardByte, err := json.Marshal(rewardMap)
	if err != nil {
		return err
	}*/
	//store.Delete(InviteRewardKey + account)
	return nil
}





func (k Keeper) updataSpaceMinerRewardSpace(ctx sdk.Context, addr string, settlement types.Settlement, statistics types.InviteRewardStatistics) error {
	
	accountMiner := k.QueryAccountSpaceMinerInfor(ctx, addr)
	accountMiner.Account = addr
	
	if accountMiner.BuySpace.Sign() <= 0 {
		return types.SpaceSettlementErr
	}
	
	/*var statistics types.InviteRewardStatistics
	if _, ok := rewardStatisticsMap[addr]; ok {
		statistics = rewardStatisticsMap[addr]
	}*/
	
	
	
	spaceTotalInv := settlement.ExpansionRewardSpace

	expansionReward := spaceTotalInv
	
	sp := accountMiner.BuySpace.Mul(decimal.RequireFromString(InviteSpaceRateKey)).Sub(accountMiner.RewardSpace)
	
	if expansionReward.Sign() > 0 {
		statistics.ExpansionRewardCounts = statistics.ExpansionRewardCounts + 1

		
		if sp.LessThan(expansionReward) {
			expansionReward = sp
			spaceTotalInv = sp
		}
		
		settlement.ExpansionRewardSpace = settlement.ExpansionRewardSpace.Sub(expansionReward)
		statistics.ExpansionRewardSpace = statistics.ExpansionRewardSpace.Add(expansionReward).Round(4)
	}
	
	if settlement.InviteRewardSpace.Sign() > 0 {
		inviteReward := settlement.InviteRewardSpace
		
		if sp.GreaterThan(expansionReward) {
			remainingSpace := sp.Sub(expansionReward)
			if remainingSpace.LessThan(settlement.InviteRewardSpace) {
				inviteReward = remainingSpace
			}
			
			settlement.InviteRewardSpace = settlement.InviteRewardSpace.Sub(inviteReward)

			spaceTotalInv = spaceTotalInv.Add(inviteReward)
			
			statistics.InviteRewardCounts = statistics.InviteRewardCounts + 1
			statistics.InviteRewardSpace = statistics.InviteRewardSpace.Add(inviteReward).Round(4)
		}
	}

	
	accountMiner.SpaceTotal = accountMiner.SpaceTotal.Add(spaceTotalInv).Round(4)
	
	accountMiner.RewardSpace = accountMiner.RewardSpace.Add(spaceTotalInv).Round(4)
	
	accountMiner = k.calSettlementMap(ctx, accountMiner)
	k.SetAccountSpaceMinerInfor(ctx, accountMiner)
	store := k.KVHelper(ctx)
	err := store.Set(InvitesStatisticsKey+addr, statistics)
	if err != nil {
		return err
	}
	
	err = store.Set(InviteRewardKey+addr, settlement)
	if err != nil {
		return err
	}
	
	k.SetDeflationSpaceTotal(ctx, spaceTotalInv)
	return nil
}

func (k Keeper) calSettlementMap(ctx sdk.Context, accountMiner AccountSpaceMiner) AccountSpaceMiner {
	height := ctx.BlockHeight()
	index := (height - core.MinerStartHeight) / core.SpaceMinerBonusBlockNum
	if index <= 0 {
		index = 1
	} else {
		index = index + 1
	}
	settlementMap := accountMiner.Settlement
	if settlementMap == nil {
		settlementMap = make(map[int64]Settlement)
	}
	set := Settlement{
		Index:      index,
		IndexSpace: accountMiner.SpaceTotal,
	}
	settlementMap[index] = set
	accountMiner.Settlement = settlementMap
	
	if accountMiner.SettlementEnd.IndexSpace.IsZero() {
		accountMiner.SettlementEnd = set
	}
	return accountMiner
}


func (k Keeper) InviteReward(ctx sdk.Context, space decimal.Decimal, address sdk.AccAddress, count int) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	log.WithFields(logrus.Fields{"address": address.String(), "": space}).Debug("")
	store := k.KVHelper(ctx)
	
	if !store.Has(InviteRelationKey + address.String()) {
		return nil
	}
	
	PreAddress := string(store.Get(InviteRelationKey + address.String()))
	
	k.RecursionInvite(ctx, PreAddress, space, count)
	return nil
}


/**
invite  
accountReward 
space 
counts  ,
*/
func (k Keeper) RecursionInvite(ctx sdk.Context, preAddress string, space decimal.Decimal, counts int) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	//1B 
	if space.Cmp(decimal.NewFromInt(1)) < 0 || counts > 1 { 
		return
	}
	
	var rewardSpace decimal.Decimal
	var accountReward types.Settlement
	//inviteRewardKey :=  + strings.ToLower()
	
	
	store := k.KVHelper(ctx)
	
	//,"".
	if preAddress == "" {
		return
	}
	
	if store.Has(InviteRewardKey + preAddress) {
		err := store.GetUnmarshal(InviteRewardKey+preAddress, &accountReward)
		if err != nil {
			log.WithError(err).WithField("store key", InviteRewardKey+preAddress).Error("GetUnmarshal")
			return
		}
	}
	
	if counts == 1 { 
		rewardSpace = space.Mul(decimal.RequireFromString(InviteFirstRateKey))
		accountReward.ExpansionRewardSpace = accountReward.ExpansionRewardSpace.Add(rewardSpace).Round(4)
	} else {
		
		rewardSpace = space.Mul(decimal.RequireFromString(InviteRateKey))
		
		accountReward.InviteRewardSpace = accountReward.InviteRewardSpace.Add(rewardSpace).Round(4)
	}
	
	inviteKey := InviteRewardKey + preAddress
	
	err := store.Set(inviteKey, accountReward)
	if err != nil {
		panic(err)
	}
	
	if !store.Has(InviteRelationKey + preAddress) {
		return
	}
	
	newPreAddress := string(store.Get(InviteRelationKey + preAddress))
	k.RecursionInvite(ctx, newPreAddress, rewardSpace.Round(4), 2)
	return
}
