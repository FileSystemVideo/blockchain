package keeper

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

//--------------------------------------------------------------------------------------------------------


func (k Keeper) ExportInviteRelation(ctx sdk.Context) []types.InviteRelation {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(InviteRelationKey)
	defer iterator.Close()
	inviteRelationArray := make([]types.InviteRelation, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("Invitation relationship information:", string(iterator.Key()), string(iterator.Value()))
		inviteRelation := types.InviteRelation{
			InviteRelationKey: string(iterator.Key()),
			InviteAddress:     string(iterator.Value()),
		}
		inviteRelationArray = append(inviteRelationArray, inviteRelation)
	}
	return inviteRelationArray
}


func (k Keeper) ExportInviteRecording(ctx sdk.Context) []types.InviteRecords {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(InviteRecordingKey)
	defer iterator.Close()
	inviteRecordsArray := make([]types.InviteRecords, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("Invitation record information:", string(iterator.Key()), string(iterator.Value()))
		var inviteRecording []types.AccountInviteRecord
		err := util.Json.Unmarshal(iterator.Value(), &inviteRecording)
		if err != nil {
			panic(err)
		}
		inviteRecords := types.InviteRecords{}
		inviteRecords.InviteRecordingKey = string(iterator.Key())
		inviteRecords.AccountInviteRecord = inviteRecording
		inviteRecordsArray = append(inviteRecordsArray, inviteRecords)
	}
	return inviteRecordsArray
}


func (k Keeper) ExportInviteReward(ctx sdk.Context) []types.InviteReward {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(InviteRewardKey)
	defer iterator.Close()
	inviteRewardArray := make([]types.InviteReward, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("InviteReward information:", string(iterator.Key()), string(iterator.Value()))
		var inviteSettlement types.InviteSettlement
		err := util.Json.Unmarshal(iterator.Value(), &inviteSettlement)
		if err != nil {
			panic(err)
		}
		inviteReward := types.InviteReward{}
		inviteReward.InviteRewardKey = string(iterator.Key())
		inviteReward.InviteSettlement = inviteSettlement
		inviteRewardArray = append(inviteRewardArray, inviteReward)
	}
	return inviteRewardArray
}


func (k Keeper) ExportInviteStatistics(ctx sdk.Context) []types.InvitesStatistics {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(InvitesStatisticsKey)
	defer iterator.Close()
	inviteStatisticsArray := make([]types.InvitesStatistics, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("Invitation record information:", string(iterator.Key()), string(iterator.Value()))
		var genInviteStatistics types.GenesisInviteRewardStatistics
		err := util.Json.Unmarshal(iterator.Value(), &genInviteStatistics)
		if err != nil {
			panic(err)
		}
		inviteStatistics := types.InvitesStatistics{}
		inviteStatistics.InvitesStatisticsKey = string(iterator.Key())
		inviteStatistics.InviteRewardStatistics = genInviteStatistics
		inviteStatisticsArray = append(inviteStatisticsArray, inviteStatistics)
	}
	return inviteStatisticsArray
}

//----------------------------------------------------------------------------------------------------------


func (k Keeper) ExportAccountSpace(ctx sdk.Context) []types.AccountSpace {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(spaceMinerKey)
	defer iterator.Close()
	accountSpaceArray := make([]types.AccountSpace, 0)
	for ; iterator.Valid(); iterator.Next() {
		var accountSpace types.AccountSpace
		ctx.Logger().Info("spatial information:", string(iterator.Key()), string(iterator.Value()))
		key := string(iterator.Key())
		//key
		if strings.Contains(key, spaceMinerAccountKey) || strings.Contains(key, spaceMinerAmountKey) || strings.Contains(key, spaceMinerBonusKey) {
			continue
		}
		err := util.Json.Unmarshal(iterator.Value(), &accountSpace)
		if err != nil {
			
			//var accountMap map[string]AccountSpaceMiner
			//err = util.Json.Unmarshal(iterator.Value(), &accountMap)
			//if err != nil {
			//	panic(err)
			//}
			panic(err)
		}
		accountSpaceArray = append(accountSpaceArray, accountSpace)
	}
	return accountSpaceArray
}


func (k Keeper) ExportDeflationInfor(ctx sdk.Context) types.DeflationInfor {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(deflationMinerKey)
	defer iterator.Close()
	deflationInfor := types.DeflationInfor{}
	for ; iterator.Valid(); iterator.Next() {
		err := util.Json.Unmarshal(iterator.Value(), &deflationInfor)
		if err != nil {
			log.Error("Failed to export deflationary mining basic information")
			panic(err)
		}
	}
	
	spaceMinerAmountByte := store.Get(spaceMinerAmountKey)
	if spaceMinerAmountByte != nil {
		deflationInfor.SpaceMinerAmount = string(spaceMinerAmountByte)
	}

	
	spaceMinerBonusByte := store.Get(spaceMinerBonusKey)
	if spaceMinerBonusByte != nil {
		deflationInfor.SpaceMinerBonus = string(spaceMinerBonusByte)
	}

	
	deflationSpaceTotalByte := store.Get(deflationSpaceTotalKey)
	if deflationSpaceTotalByte != nil {
		deflationInfor.DeflationSpaceTotal = string(deflationSpaceTotalByte)
	}

	//key
	spaceMinerAccountByte := store.Get(spaceMinerAccountKey)
	if spaceMinerAccountByte != nil {
		deflationInfor.SpaceMinerAccount = string(spaceMinerAccountByte)
	}

	
	iterator = store.KVStorePrefixIterator(spaceTotalIndexKey)
	spaceTotalArray := make([]types.SpaceTotalIndex, 0)
	for ; iterator.Valid(); iterator.Next() {
		var spaceTotal types.SpaceTotalIndex
		spaceTotal.SpaceTotalIndexKey = string(iterator.Key())
		spaceTotal.SpaceTotal = string(iterator.Value())
		spaceTotalArray = append(spaceTotalArray, spaceTotal)
	}
	deflationInfor.SpaceTotalIndex = spaceTotalArray

	return deflationInfor
}

//----------------------------------------------------------------------------------------------------------

func (k Keeper) ExportCopyrightParty(ctx sdk.Context) []types.GenesisCopyrightPart {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightPartyKey)
	defer iterator.Close()
	copyrightPartyArray := make([]types.GenesisCopyrightPart, 0)
	for ; iterator.Valid(); iterator.Next() {
		ctx.Logger().Info("ExportCopyrightParty:", string(iterator.Key()), string(iterator.Value()))
		var copyrightPart types.CopyrightParty
		err := util.Json.Unmarshal(iterator.Value(), &copyrightPart)
		if err != nil {
			panic(err)
		}
		var genCopyrightPart types.GenesisCopyrightPart
		genCopyrightPart.CopyrightPartyKey = string(iterator.Key())
		genCopyrightPart.CopyrightParty = copyrightPart
		copyrightPartyArray = append(copyrightPartyArray, genCopyrightPart)
	}
	return copyrightPartyArray
}

//ID
func (k Keeper) ExportCopyrightPublishId(ctx sdk.Context) types.CopyrightPublishId {
	store := k.KVHelper(ctx)
	var copyrightPublishId types.CopyrightPublishId
	copyrightPublishIdByte := store.Get(types.CopyrightPublishIdKey)
	if copyrightPublishIdByte == nil {
		return copyrightPublishId
	}
	data := make(map[string]string)
	err := util.Json.Unmarshal(copyrightPublishIdByte, &data)
	if err != nil {
		panic(err)
	}
	copyrightPublishId.PublishId = data
	return copyrightPublishId
}

//----------------------------------------------------------------------------------------------------------
//IP
func (k Keeper) ExportCopyrightIp(ctx sdk.Context) []types.GenesisCopyrightIp {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightIpKey)
	defer iterator.Close()
	copyrightIpArray := make([]types.GenesisCopyrightIp, 0)
	for ; iterator.Valid(); iterator.Next() {
		var copyrightIp types.CopyrightIp
		err := util.Json.Unmarshal(iterator.Value(), &copyrightIp)
		if err != nil {
			panic(err)
		}
		var genCopyrightIp types.GenesisCopyrightIp
		genCopyrightIp.CopyrightIpKey = string(iterator.Key())
		genCopyrightIp.CopyrightIp = copyrightIp
		copyrightIpArray = append(copyrightIpArray, genCopyrightIp)
	}
	return copyrightIpArray
}


func (k Keeper) ExportCopyrightOriginHash(ctx sdk.Context) []types.GenesisCopyrightOriginDataHash {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightOriginHashKey)
	defer iterator.Close()
	originDataHashArray := make([]types.GenesisCopyrightOriginDataHash, 0)
	for ; iterator.Valid(); iterator.Next() {
		var originDataHash types.CopyrightOriginDataHash
		err := util.Json.Unmarshal(iterator.Value(), &originDataHash)
		if err != nil {
			panic(err)
		}
		var genCopyrightOriginDataHash types.GenesisCopyrightOriginDataHash
		genCopyrightOriginDataHash.CopyrightOriginHashKey = string(iterator.Key())
		genCopyrightOriginDataHash.OriginDataHash = originDataHash
		originDataHashArray = append(originDataHashArray, genCopyrightOriginDataHash)
	}
	return originDataHashArray
}


func (k Keeper) ExportCopyrightBonusAddress(ctx sdk.Context) []types.GenesisCopyrightBonus {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightBonusAddressKey)
	defer iterator.Close()
	bonusArray := make([]types.GenesisCopyrightBonus, 0)
	for ; iterator.Valid(); iterator.Next() {
		var copyrightBonus types.CopyrightBonus
		err := util.Json.Unmarshal(iterator.Value(), &copyrightBonus)
		if err != nil {
			panic(err)
		}
		var genCopyrightBonus types.GenesisCopyrightBonus
		genCopyrightBonus.CopyrightBonusAddressKey = string(iterator.Key())
		genCopyrightBonus.CopyrightBonus = copyrightBonus
		bonusArray = append(bonusArray, genCopyrightBonus)
	}
	return bonusArray
}

//NFT
func (k Keeper) ExportCopyrightNft(ctx sdk.Context) []types.GenesisNftInfo {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(nftTokenIdKey)
	defer iterator.Close()
	nftArray := make([]types.GenesisNftInfo, 0)
	for ; iterator.Valid(); iterator.Next() {
		var copyrightNft types.NftInfoData
		err := util.Json.Unmarshal(iterator.Value(), &copyrightNft)
		if err != nil {
			panic(err)
		}
		var genesisNftInfo types.GenesisNftInfo
		genesisNftInfo.NftTokenIdKey = string(iterator.Key())
		genesisNftInfo.NftInfo = copyrightNft
		nftArray = append(nftArray, genesisNftInfo)
	}
	return nftArray
}

//----------------------------------------------------------------------------------------------------------

func (k Keeper) ExportCopyrightVote(ctx sdk.Context) []types.GenesisCopyrightVote {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(copyrightVoteFor)
	defer iterator.Close()
	voteArray := make([]types.GenesisCopyrightVote, 0)
	for ; iterator.Valid(); iterator.Next() {
		copyrightVote := make(map[string]types.CopyrightVoteShare)
		err := util.Json.Unmarshal(iterator.Value(), &copyrightVote)
		if err != nil {
			panic(err)
		}
		var genesisVote types.GenesisCopyrightVote
		genesisVote.CopyrightVoteKey = string(iterator.Key())
		genesisVote.VoteData = copyrightVote
		voteArray = append(voteArray, genesisVote)
	}
	return voteArray
}


func (k Keeper) ExportCopyrightVoteList(ctx sdk.Context) []types.GenesisCopyrightVoteList {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(copyrightVoteListFor)
	defer iterator.Close()
	voteListArray := make([]types.GenesisCopyrightVoteList, 0)
	for ; iterator.Valid(); iterator.Next() {
		copyrightVoteList := make([]types.AccountVote, 0)
		err := util.Json.Unmarshal(iterator.Value(), &copyrightVoteList)
		if err != nil {
			panic(err)
		}
		var genesisVoteList types.GenesisCopyrightVoteList
		genesisVoteList.CopyrightVoteListKey = string(iterator.Key())
		genesisVoteList.AccountVote = copyrightVoteList
		voteListArray = append(voteListArray, genesisVoteList)
	}
	return voteListArray
}


func (k Keeper) ExportCopyrightApproveResult(ctx sdk.Context) []types.CopyrightApproveResultData {
	store := k.KVHelper(ctx)
	copyrightApproveForResultByte := store.Get(copyrightApproveForResult)
	data := []types.CopyrightApproveResultData{}
	if copyrightApproveForResultByte != nil {
		err := util.Json.Unmarshal(copyrightApproveForResultByte, &data)
		if err != nil {
			panic(err)
		}
	}
	return data
}


func (k Keeper) ExportCopyrightVoteRedeem(ctx sdk.Context) string {
	store := k.KVHelper(ctx)
	copyrightVoteRedeemByte := store.Get(copyrightVoteRedeem)
	if copyrightVoteRedeemByte != nil {
		return string(copyrightVoteRedeemByte)
	}
	return ""
}


func (k Keeper) ExportCopyright(ctx sdk.Context) []types.GenesisCopyright {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightDetailKey)
	defer iterator.Close()
	copyrightArray := make([]types.GenesisCopyright, 0)
	for ; iterator.Valid(); iterator.Next() {
		var copyright types.Copyright
		err := util.Json.Unmarshal(iterator.Value(), &copyright)
		if err != nil {
			panic(err)
		}
		var genesisCopyright types.GenesisCopyright
		genesisCopyright.CopyrightKey = string(iterator.Key())
		genesisCopyright.Copyright = copyright
		copyrightArray = append(copyrightArray, genesisCopyright)
	}
	return copyrightArray
}


func (k Keeper) ExportCopyrightExtra(ctx sdk.Context) []types.GenesisCopyrightExtra {
	store := k.KVHelper(ctx)
	
	accounts := k.authKeeper.GetAllAccounts(ctx)
	accountList := []string{}
	for _, acc := range accounts {
		accountList = append(accountList, acc.GetAddress().String())
	}
	
	iterator := store.KVStorePrefixIterator(types.CopyrightDetailKey)
	defer iterator.Close()
	var copyrightArray []types.Copyright
	for ; iterator.Valid(); iterator.Next() {
		var copyright types.Copyright
		err := util.Json.Unmarshal(iterator.Value(), &copyright)
		if err != nil {
			panic(err)
		}
		copyrightArray = append(copyrightArray, copyright)
	}
	copyrightExtraArray := make([]types.GenesisCopyrightExtra, 0)
	//+key
	for _, acc := range accountList {
		for _, val := range copyrightArray {
			keys := acc + "_" + val.DataHash
			if store.Has(keys) {
				var copyrightExtra types.CopyrightExtra
				copyrightExtraBytes := store.Get(keys)
				err := util.Json.Unmarshal(copyrightExtraBytes, &copyrightExtra)
				if err != nil {
					panic(err)
				}
				var genesisCopyrightExtra types.GenesisCopyrightExtra
				//key
				genesisCopyrightExtra.CopyrightExtraKey = types.CopyrightRelationKey + keys
				genesisCopyrightExtra.CopyrightExtra = copyrightExtra
				copyrightExtraArray = append(copyrightExtraArray, genesisCopyrightExtra)
			}
		}
	}
	return copyrightExtraArray
}


func (k Keeper) ExportCopyrightExtraNew(ctx sdk.Context) []types.GenesisCopyrightExtra {
	store := k.KVHelper(ctx)
	iterator := store.KVStorePrefixIterator(types.CopyrightRelationKey)
	defer iterator.Close()
	copyrightExtraArray := make([]types.GenesisCopyrightExtra, 0)
	for ; iterator.Valid(); iterator.Next() {
		var copyrightExtra types.CopyrightExtra
		err := util.Json.Unmarshal(iterator.Value(), &copyrightExtra)
		if err != nil {
			panic(err)
		}
		var genesisCopyrightExtra types.GenesisCopyrightExtra
		genesisCopyrightExtra.CopyrightExtraKey = string(iterator.Key())
		genesisCopyrightExtra.CopyrightExtra = copyrightExtra
		copyrightExtraArray = append(copyrightExtraArray, genesisCopyrightExtra)
	}
	return copyrightExtraArray
}
