package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BaseResponse struct {
	Info   string `json:"info"`
	Status int    `json:"status"`
}

type TxResponse struct {
	Height string `json:"height"`
	TxHash string `json:"txhash"`
	Info   string `json:"info"`
}

type BroadcastTxResponse struct {
	BaseResponse
	Height      int64  `json:"height"`
	TxHash      string `json:"txhash"`
	Codespace   string `json:"codespace"`
	Code        uint32 `json:"code"`
	SignedTxStr string `json:"signed_tx_str"`
}

func (this *BaseResponse) IsSuccess() bool {
	return this.Status == 0
}

type AccountNumberSeqResponse struct {
	BaseResponse
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}

type MsgsResponse struct {
	BaseResponse
	Message MessageResp `json:"message"`
}

type TxPropsResponse struct {
	BaseResponse
	Props []TxPropValue `json:"props"`
}

type MessageResp struct {
	CopyrightPartys []CopyrightPartyData `json:"copyright_partys"`
	Copyrights      []CopyrightData      `json:"copyrights"`
	Transfers       []TransferData       `json:"transfers"`
	SpaceMiners     []SpaceMinerData     `json:"space_miners"`
	DeflationVotes  []DeflationVoteData  `json:"deflation_votes"`
	NftTransfers    []NftTransferData    `json:"nft_transfers"`
	Delegations     []DelegationData     `json:"delegations"`
	Undelegations   []UndelegationData   `json:"undelegations"`
	Mortgages          []MortgageData          `json:"mortgages"`
	EditorCopyrights   []EditorCopyrightData   `json:"editor_copyrights"`
	DeleteCopyrights   []DeleteCopyrightData   `json:"delete_copyrights"`
	BonusCopyrights    []CopyrightBonusData    `json:"bonus_copyrights"`
	CopyrightComplains []CopyrightComplainData `json:"copyright_complains"`
	ComplainResponses  []ComplainResponseData  `json:"complain_responses"`
	ComplainVotes      []ComplainVoteData      `json:"complain_votes"`
}

func (this *MessageResp) Count() int {
	return len(this.CopyrightPartys) +
		len(this.Copyrights) +
		len(this.Transfers) +
		len(this.SpaceMiners) +
		len(this.DeflationVotes) +
		len(this.NftTransfers) +
		len(this.Delegations) +
		len(this.Undelegations) +
		len(this.Mortgages) +
		len(this.EditorCopyrights) +
		len(this.DeleteCopyrights) +
		len(this.BonusCopyrights) +
		len(this.CopyrightComplains) +
		len(this.ComplainResponses) +
		len(this.ComplainVotes)
}

type BlockMessagesResponse struct {
	List []BlockMessageI `json:"list"`
}

type BalanceResponse struct {
	BaseResponse
	Height string  `json:"height"`
	Token  []Token `json:"result"`
}

type Token struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type DelegationPreviewResponse struct {
	Shares        string `json:"shares"`
	SourceAmount  string `json:"source_amount"`
	SourceShares  string `json:"source_shares"`
	BalanceAmount string `json:"balance_amount"`
	BalanceShares string `json:"balance_shares"`
}

type UnbondingDelegationPreviewResponse struct {
	Amount        string `json:"amount"`
	Shares        string `json:"shares"`
	BalanceAmount string `json:"balance_amount"`
	BalanceShares string `json:"balance_shares"`
	SourceAmount  string `json:"source_amount"`
	SourceShares  string `json:"source_shares"`
}

type DelegationDetailResponse struct {
	DelegationAddr string `json:"delegation_addr"`
	ValidatorAddr  string `json:"validator_addr"`
	Amount         string `json:"amount"`
	Shares         string `json:"shares"`
}

type PosReportFormResponse struct {
	UnbondAmount        string `json:"unbond_amount"`
	PosRewardReceived   string `json:"pos_reward_received"`
	PosRewardUnreceived string `json:"pos_reward_unreceived"`
	MortgAmount         string `json:"mortg_amount"`
	Shares              string `json:"shares"`
	ValidatorShares     string `json:"validator_shares"`
	TotalShares         string `json:"total_shares"`
	Account             string `json:"account"`
	AccountTotalShares  string `json:"account_total_shares"`
}

type ValidatorRegisterLimit struct {
	MortgAmount string `json:"mortg_amount"`
	Status      string `json:"status"`
}

type ValidatorInfor struct {
	ValidatorConsAddr string `json:"validator_consaddr"`
	ValidatorStatus   string `json:"validator_status"`
	Status            string `json:"status"`
	ValidatorPubAddr  string `json:"validator_pubaddr"`
	ValidatorOperAddr string `json:"validator_operaddr"`
	AccAddr           string `json:"acc_addr"`
}

type ValidatorsDelegationResp struct {
	Shares  string `json:"shares"`
	Balance string `json:"balance"`
}

type ValidatorCommission struct {
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Reward           sdk.DecCoin    `json:"reward"`
}
type ValidatorCommissionResp struct {
	ValidatorCommissions []ValidatorCommission `json:"validator_commissions"`
	Total                sdk.DecCoin           `json:"total"`
}
