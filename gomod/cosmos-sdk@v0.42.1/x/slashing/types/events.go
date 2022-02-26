//noalias
package types

// Slashing module event types
const (
	EventTypeSlash    = "slash"
	EventTypeLiveness = "liveness"

	AttributeKeyValAddress    = "val_address"
	AttributeKeySlashFraction = "slash_fraction"

	AttributeKeyAddress      = "address"
	AttributeKeyHeight       = "height"
	AttributeKeyPower        = "power"
	AttributeKeyReason       = "reason"
	AttributeKeyJailed       = "jailed"
	AttributeKeyMissedBlocks = "missed_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
