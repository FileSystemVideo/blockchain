package copyright

import (
	"fs.video/blockchain/x/copyright/keeper"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	
	k.InitAccountSpace(ctx, genState.AccountSpace)

	
	k.InitDeflationInfor(ctx, genState.DeflationInfor)

	
	k.InitInviteRelation(ctx, genState.InviteRelation)

	
	k.InitInviteRecording(ctx, genState.InviteRecords)

	
	k.InitInviteReward(ctx, genState.InviteReward)

	
	k.InitInviteStatistics(ctx, genState.InvitesStatistics)

	
	k.InitCopyrightParty(ctx, genState.CopyrightPart)

	//ID
	k.InitCopyrightPublishId(ctx, genState.CpyrightPublishId)

	
	k.InitCopyright(ctx, genState.Copyright)

	
	k.InitCopyrightExtra(ctx, genState.CopyrightExtra)

	//IP
	k.InitCopyrightIp(ctx, genState.CopyrightIp)

	
	k.InitCopyrightOriginHash(ctx, genState.CopyrightOriginHash)

	
	k.InitCopyrightBonusAddress(ctx, genState.CopyrightBonus)

	//NFT
	k.InitCopyrightNft(ctx, genState.NftInfo)

	
	k.InitCopyrightVote(ctx, genState.CopyrightVote)

	
	k.InitCopyrightVoteList(ctx, genState.CopyrightVoteList)

	
	k.InitCopyrightApproveResult(ctx, genState.ApproveResult)

	
	k.InitCopyrightVoteRedeem(ctx, genState.CopyrightVoteRedeem)

	
	k.SetParams(ctx, genState.Params)

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.ExportAccountSpace(ctx),
		k.ExportDeflationInfor(ctx),
		k.ExportInviteRelation(ctx),
		k.ExportInviteRecording(ctx),
		k.ExportInviteReward(ctx),
		k.ExportInviteStatistics(ctx),
		k.ExportCopyrightParty(ctx),
		k.ExportCopyrightPublishId(ctx),
		k.ExportCopyright(ctx),
		k.ExportCopyrightExtra(ctx),
		k.ExportCopyrightIp(ctx),
		k.ExportCopyrightOriginHash(ctx),
		k.ExportCopyrightBonusAddress(ctx),
		k.ExportCopyrightNft(ctx),
		k.ExportCopyrightVote(ctx),
		k.ExportCopyrightVoteList(ctx),
		k.ExportCopyrightApproveResult(ctx),
		k.ExportCopyrightVoteRedeem(ctx),
	)
}
