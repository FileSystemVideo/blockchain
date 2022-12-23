package types

const (
	MSG_SMART_CREATE_VALIDATOR = "smart/MsgCreateSmartValidator"
	TypeContractProposal       = "contract/ContractProposal"
)


type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"` 
	ValidatorStatus   string `json:"validator_status"`   //0  Unbonded 1 Unbonding 2 Bonded 3  4 
	ValidatorPubAddr  string `json:"validator_pubaddr"`  
	ValidatorOperAddr string `json:"validator_operaddr"` 
	AccAddr           string `json:"acc_addr"`           
}

type MsgToByte struct {
	MsgType string `json:"msg_type"`
	Msg     string `json:"msg"`
}
