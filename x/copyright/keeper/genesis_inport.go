package keeper

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
)

//--------------------------------------------------------------------------------------------------------


func (k Keeper) InitInviteRelation(ctx sdk.Context, inviteRelation []types.InviteRelation) {
	store := k.KVHelper(ctx)
	for _, val := range inviteRelation {
		err := store.Set(val.InviteRelationKey, val.InviteAddress)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitInviteRecording(ctx sdk.Context, inviteRecords []types.InviteRecords) {
	store := k.KVHelper(ctx)
	for _, val := range inviteRecords {
		err := store.Set(val.InviteRecordingKey, val.AccountInviteRecord)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitInviteReward(ctx sdk.Context, inviteReward []types.InviteReward) {
	store := k.KVHelper(ctx)
	for _, val := range inviteReward {
		err := store.Set(val.InviteRewardKey, val.InviteSettlement)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitInviteStatistics(ctx sdk.Context, inviteStatistics []types.InvitesStatistics) {
	store := k.KVHelper(ctx)
	for _, val := range inviteStatistics {
		err := store.Set(val.InvitesStatisticsKey, val.InviteRewardStatistics)
		if err != nil {
			panic(err)
		}
	}
}

//----------------------------------------------------------------------------------------------------------


func (k Keeper) InitAccountSpace(ctx sdk.Context, accountSpace []types.AccountSpace) {
	store := k.KVHelper(ctx)
	for _, val := range accountSpace {
		err := store.Set(spaceMinerKey+val.Account, val)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitDeflationInfor(ctx sdk.Context, deflationInfor types.DeflationInfor) {
	store := k.KVHelper(ctx)
	deflationMinerInfor := DeflationMinerInfor{
		DeflationStatus:   int(deflationInfor.DeflationStatus),
		MinerTotalAmount:  decimal.RequireFromString(deflationInfor.MinerTotalAmount),
		HasMinerAmount:    decimal.RequireFromString(deflationInfor.HasMinerAmount),
		RemainMinerAmount: decimal.RequireFromString(deflationInfor.RemainMinerAmount),
		DayMinerAmount:    decimal.RequireFromString(deflationInfor.DayMinerAmount),
		DayMinerRemain:    deflationInfor.DayMinerRemain,
		MinerBlockNum:     core.MinerStartHeight, 
	}
	err := store.Set(deflationMinerKey, deflationMinerInfor)
	if err != nil {
		panic(err)
	}
	
	err = store.Set(spaceMinerAmountKey, deflationInfor.SpaceMinerAmount)
	if err != nil {
		panic(err)
	}
	
	err = store.Set(spaceMinerBonusKey, deflationInfor.SpaceMinerBonus)
	if err != nil {
		panic(err)
	}
	
	err = store.Set(deflationSpaceTotalKey, deflationInfor.DeflationSpaceTotal)
	if err != nil {
		panic(err)
	}
	//key
	err = store.Set(spaceMinerAccountKey, deflationInfor.SpaceMinerAccount)
	if err != nil {
		panic(err)
	}
	
	for _, val := range deflationInfor.SpaceTotalIndex {
		err := store.Set(val.SpaceTotalIndexKey, val.SpaceTotal)
		if err != nil {
			panic(err)
		}
	}
}

//----------------------------------------------------------------------------------------------------------

func (k Keeper) InitCopyrightParty(ctx sdk.Context, copyrightPart []types.GenesisCopyrightPart) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightPart {
		err := store.Set(val.CopyrightPartyKey, val.CopyrightParty)
		if err != nil {
			panic(err)
		}
	}
}

//id
func (k Keeper) InitCopyrightPublishId(ctx sdk.Context, copyrightId types.CopyrightPublishId) {
	store := k.KVHelper(ctx)
	err := store.Set(types.CopyrightPublishIdKey, copyrightId.PublishId)
	if err != nil {
		panic(err)
	}
}

//----------------------------------------------------------------------------------------------------------
//IP
func (k Keeper) InitCopyrightIp(ctx sdk.Context, copyrightIp []types.GenesisCopyrightIp) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightIp {
		err := store.Set(val.CopyrightIpKey, val.CopyrightIp)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightOriginHash(ctx sdk.Context, copyrightOriginDataHash []types.GenesisCopyrightOriginDataHash) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightOriginDataHash {
		err := store.Set(val.CopyrightOriginHashKey, val.OriginDataHash)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightBonusAddress(ctx sdk.Context, copyrightBonus []types.GenesisCopyrightBonus) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightBonus {
		err := store.Set(val.CopyrightBonusAddressKey, val.CopyrightBonus)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightNft(ctx sdk.Context, copyrightNft []types.GenesisNftInfo) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightNft {
		err := store.Set(val.NftTokenIdKey, val.NftInfo)
		if err != nil {
			panic(err)
		}
	}
}

//----------------------------------------------------------------------------------------------------------

func (k Keeper) InitCopyrightVote(ctx sdk.Context, copyrightVote []types.GenesisCopyrightVote) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightVote {
		err := store.Set(val.CopyrightVoteKey, val.VoteData)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightVoteList(ctx sdk.Context, copyrightVoteList []types.GenesisCopyrightVoteList) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightVoteList {
		err := store.Set(val.CopyrightVoteListKey, val.AccountVote)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightApproveResult(ctx sdk.Context, copyrightVoteList []types.CopyrightApproveResultData) {
	store := k.KVHelper(ctx)
	err := store.Set(copyrightApproveForResult, copyrightVoteList)
	if err != nil {
		panic(err)
	}
}


func (k Keeper) InitCopyrightVoteRedeem(ctx sdk.Context, voteRedeem string) {
	store := k.KVHelper(ctx)
	if voteRedeem != "" {
		err := store.Set(copyrightVoteRedeem, voteRedeem)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyright(ctx sdk.Context, copyright []types.GenesisCopyright) {
	store := k.KVHelper(ctx)
	for _, val := range copyright {
		err := store.Set(val.CopyrightKey, val.Copyright)
		if err != nil {
			panic(err)
		}
	}
}


func (k Keeper) InitCopyrightExtra(ctx sdk.Context, copyrightExtra []types.GenesisCopyrightExtra) {
	store := k.KVHelper(ctx)
	for _, val := range copyrightExtra {
		err := store.Set(val.CopyrightExtraKey, val.CopyrightExtra)
		if err != nil {
			panic(err)
		}
	}
}
