package types

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		AccountSpace: []AccountSpace{},
		DeflationInfor: DeflationInfor{
			
			MinerTotalAmount: "1039500000",
			
			HasMinerAmount: "396000",
			
			RemainMinerAmount: "1039140000",
			
			DayMinerAmount: "36000",
			
			DayMinerRemain: 28865,
			
			DeflationStatus: 1,
			
			SpaceMinerAmount: "7000000",
			
			SpaceMinerBonus: "396000",
			
			DeflationSpaceTotal: "7340032000000",
			
			SpaceTotalIndex: []SpaceTotalIndex{},
			
			SpaceMinerAccount: "",
		},
		InviteRecords:  []InviteRecords{},
		InviteRelation: []InviteRelation{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	return nil
}

func NewGenesisState(accountSpaces []AccountSpace,
	deflationInfor DeflationInfor,
	inviteRelation []InviteRelation,
	inviteRecord []InviteRecords,
	inviteReward []InviteReward,
	statistics []InvitesStatistics,
	party []GenesisCopyrightPart,
	publishId CopyrightPublishId,
	copyright []GenesisCopyright,
	extra []GenesisCopyrightExtra,
	copyrightIp []GenesisCopyrightIp,
	originDataHash []GenesisCopyrightOriginDataHash,
	copyrightBonus []GenesisCopyrightBonus,
	nftInfo []GenesisNftInfo,
	copyrightVote []GenesisCopyrightVote,
	copyrightVoteList []GenesisCopyrightVoteList,
	approveResult []CopyrightApproveResultData,
	copyrightVoteRedeem string) *GenesisState {
	return &GenesisState{
		AccountSpace:        accountSpaces,
		DeflationInfor:      deflationInfor,
		InviteRelation:      inviteRelation,
		InviteRecords:       inviteRecord,
		InviteReward:        inviteReward,
		InvitesStatistics:   statistics,
		CopyrightPart:       party,
		CpyrightPublishId:   publishId,
		Copyright:           copyright,
		CopyrightExtra:      extra,
		CopyrightIp:         copyrightIp,
		CopyrightOriginHash: originDataHash,
		CopyrightBonus:      copyrightBonus,
		NftInfo:             nftInfo,
		CopyrightVote:       copyrightVote,
		CopyrightVoteList:   copyrightVoteList,
		ApproveResult:       approveResult,
		CopyrightVoteRedeem: copyrightVoteRedeem,
	}
}
