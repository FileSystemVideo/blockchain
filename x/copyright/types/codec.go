package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

/**
，
*/
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateCopyright{}, MSG_TYPE_CREATE_COPYRIGHT, nil)
	cdc.RegisterConcrete(&MsgRegisterCopyrightParty{}, MSG_TYPE_REGISTER_COPYRIGHT_PARTY, nil)
	cdc.RegisterConcrete(&MsgSpaceMiner{}, MSG_TYPE_SPACE_MINER, nil)
	cdc.RegisterConcrete(&MsgNftTransfer{}, MSG_TYPE_NFT_TRANSFER, nil)
	cdc.RegisterConcrete(&MsgDistributeCommunityReward{}, MSG_TYPE_DISTRIBUTE_COMMUNITY_REWARD, nil)
	cdc.RegisterConcrete(&MsgMortgage{}, MSG_TYPE_MORTGAGE, nil)
	cdc.RegisterConcrete(&MsgEditorCopyright{}, MSG_TYPE_EDITOR_COPYRIGHT, nil)
	cdc.RegisterConcrete(&MsgDeleteCopyright{}, MSG_TYPE_DELETE_COPYRIGHT, nil)
	cdc.RegisterConcrete(&MsgCopyrightBonusV2{}, MSG_TYPE_BONUS_COPYRIGHTV2, nil)
	cdc.RegisterConcrete(&MsgCopyrightComplain{}, MSG_TYPE_COPYRIGHT_COMPLAIN, nil)
	cdc.RegisterConcrete(&MsgComplainResponse{}, MSG_TYPE_COMPLAIN_RESPONSE, nil)
	cdc.RegisterConcrete(&MsgComplainVote{}, MSG_TYPE_COMPLAIN_VOTE, nil)
	cdc.RegisterConcrete(&MsgTransfer{}, MSG_TYPE_TRANSFER, nil)
	cdc.RegisterConcrete(&MsgInviteReward{}, MSG_TYPE_INVATE_REWARD, nil)
	cdc.RegisterConcrete(&MsgSpaceMinerReward{}, MSG_TYPE_SPACE_MINER_REWARD, nil)
	cdc.RegisterConcrete(&MsgCopyrightBonusRearV2{}, MSG_TYPE_BONUS_COPYRIGHT_REARV2, nil)
	cdc.RegisterConcrete(&MsgVoteCopyright{}, MSG_TYPE_COPYRIGHT_VOTE, nil)
	cdc.RegisterConcrete(&MsgCrossChainOut{}, MSG_TYPE_CROSSCHAIN_OUT, nil)
	cdc.RegisterConcrete(&MsgCrossChainIn{}, MSG_TYPE_CROSSCHAIN_IN, nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCopyright{},
		&MsgRegisterCopyrightParty{},
		&MsgSpaceMiner{},
		&MsgNftTransfer{},
		&MsgDistributeCommunityReward{},
		&MsgMortgage{},
		&MsgEditorCopyright{},
		&MsgDeleteCopyright{},
		&MsgCopyrightBonusV2{},
		&MsgCopyrightComplain{},
		&MsgComplainResponse{},
		&MsgComplainVote{},
		&MsgTransfer{},
		&MsgInviteReward{},
		&MsgSpaceMinerReward{},
		&MsgCopyrightBonusRearV2{},
		&MsgVoteCopyright{},
		&MsgCrossChainOut{},
		&MsgCrossChainIn{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
