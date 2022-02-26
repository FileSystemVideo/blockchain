package types

// bank module event types
const (
	EventTypeTransfer = "transfer"

	AttributeKeyRecipient        = "recipient"
	AttributeKeySender           = "sender"
	AttributeKeyRecipientBalance = "recipientBalance"
	AttributeKeySenderBalance    = "senderBalance"

	AttributeValueCategory = ModuleName
)
