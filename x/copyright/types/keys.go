package types

const (
	// ModuleName defines the module name
	ModuleName = "copyright"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"

	CopyrightDetailKey = "copyright_detail_"

	CopyrightIpKey = "copyright_ip_"

	CopyrightOriginHashKey = "copyright_origin_hash_"

	CopyrightCountKey = "copyright_count_"

	CopyrightPartyKey = "copyright_party_"

	CopyrightPublishIdKey = "current_copyright_publish_id"

	BlockRelationDataKey = "block_relation_"

	ContractCrossChain = "contract_cross_chain"

	CopyrightCrossChainOutFeeRatioKey = "copyright_cross_chain_out_fee_ratio"
	EventTypeTransferFee = "transfer_fee"

	AttributeKeyTransferFee = "fee"
)

var (
	ChatsPrefix         = []byte("chats")
	SupplyKey           = []byte{0x00}
	DenomMetadataPrefix = []byte{0x1}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
