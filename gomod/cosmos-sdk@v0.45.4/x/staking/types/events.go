package types

// staking module event types
const (
	EventTypeCompleteUnbonding    = "complete_unbonding"
	EventTypeCompleteRedelegation = "complete_redelegation"
	EventTypeCreateValidator      = "create_validator"
	EventTypeEditValidator        = "edit_validator"
	EventTypeDelegate             = "delegate"
	EventTypeUnbond               = "unbond"
	EventTypeRedelegate           = "redelegate"

	AttributeKeyValidator         = "validator"
	AttributeKeyCommissionRate    = "commission_rate"
	AttributeKeyMinSelfDelegation = "min_self_delegation"
	AttributeKeySrcValidator      = "source_validator"
	AttributeKeyDstValidator      = "destination_validator"
	AttributeKeyDelegator         = "delegator"
	AttributeKeyCompletionTime    = "completion_time"
	AttributeKeyNewShares         = "new_shares"
	AttributeValueCategory        = ModuleName

	EventTypeSlashAmount         = "slash_amount"
	AttributeKeyDelegatorAddr    = "delegatorAddress"
	AttributeKeyDelegatorBalance = "delegator_balance"
)
