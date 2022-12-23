package core

import (
	"fs.video/trerr"
)



/*
'transfer': '',
'down': '',
'delete': '',
'bind-pay': '',
'bonded': 'POS',
'unbonded': 'POS',
'delegate-reward': 'POS',
'bonus': '',
'copyright-share': '',
'mining-mortg': '',
'mining-redeem': '',
'slash':  'POS',
'node-auth'： ''
'complain':''
'complain-redeem':''
publish  
*/

var (
	TradeTypeTransfer               = RegisterTranserType("transfer", "", "transfer accounts")
	TradeTypeCopyrightBuy           = RegisterTranserType("down", "", "Copyright purchase")
	TradeTypeCopyrightPublish       = RegisterTranserType("publish", "", "Copyright foundry fee")
	TradeTypeCopyrightEdit          = RegisterTranserType("edit", "", "Modify copyright")
	TradeTypeCopyrightSharesReward  = RegisterTranserType("bonus", "", "Hard disk bandwidth sharing reward")
	TradeTypeCopyrightDelete        = RegisterTranserType("delete", "", "Delete copyright")
	TradeTypeCopyrightSell          = RegisterTranserType("copyright-share", "", "Copyright sales revenue")
	TradeTypeCopyrightBuyMortgage   = RegisterTranserType("mining-mortg", "", "Copyright mortgage purchase")
	TradeTypeCopyrightBuyRedeem     = RegisterTranserType("mining-redeem", "", "Copyright mortgage purchase and redemption")
	TradeTypeCopyrightPartyRegister = RegisterTranserType("bind-pay", "", "Copyright registration")
	TradeTypeDelegation             = RegisterTranserType("bonded", "POS", "POS mortgage")
	TradeTypeDelegationFee          = RegisterTranserType("bonded-fee", "POS", "POS mortgage service charge")
	TradeTypeUnbondDelegation       = RegisterTranserType("unbonded", "POS", "POS unbonded")
	TradeTypeUnbondDelegationFee    = RegisterTranserType("unbonded-fee", "POS", "POS unbonded fee")
	TradeTypeReDelegation           = RegisterTranserType("rebonded", "POS", "POS rebonded")
	TradeTypeReDelegationFee        = RegisterTranserType("rebonded-fee", "POS", "POS rebonded fee")

	TradeTypeFee              = RegisterTranserType("fee", "GAS", "GAS fee")
	TradeTypeFeeReturn        = RegisterTranserType("fee-return", "GAS", "GAS fee return")
	TradeTypeDelegationReward = RegisterTranserType("delegate-reward", "POS", "POS mortgage reward")
	TradeTypeCommissionReward = RegisterTranserType("commission-reward", "POS", "POS Commission reward")
	TradeTypeCommunityReward  = RegisterTranserType("community-reward", "", "Community reward")
	TradeTypeValidatorUnjail  = RegisterTranserType("unjail", "POS", "POS release")
	//TradeTypeNodeAuth               = RegisterTranserType("node-auth", "")
	TradeTypeFsvDestroyMiner = RegisterTranserType("space-miner", ""+MainToken+"", "Destruction for space")
	//TradeTypeFsvDestroyMinerVote    = RegisterTranserType("space-miner-vote", ""+config.MainToken+"")
	TradeTypeNftTransfer            = RegisterTranserType("nft-transfer", "nft", "NFT transfer")
	TradeTypeCopyrightComplain      = RegisterTranserType("cp-complain", "", "Copyright complaint")
	TradeTypeSpaceMinerBonus        = RegisterTranserType("space-miner-bonus", "", "Blockchain spatial incentive")
	TradeTypeCopyrightComplainReply = RegisterTranserType("cp-complain-reply", "", "Copyright response")
	TradeTypeCopyrightComplainVotes = RegisterTranserType("cp-complain-votes", "", "Copyright appeal vote")
	TradeTypeValidatorMinerBonus    = RegisterTranserType("validator-miner-bonus", "POS", "POS incentive")
	TradeTypeValidatorCreate        = RegisterTranserType("validator-create", "", "Create validator")
	TradeTypeValidatorEditor        = RegisterTranserType("validator-editor", "", "Editor validator")
	TradeTypeInviteReward           = RegisterTranserType("invite-reward", "", "Reward expansion")
	TradeTypeGovProposal            = RegisterTranserType("gov-prop", "", "Gov Proposal")
	TradeTypeGovVote                = RegisterTranserType("gov-vote", "", "Gov voting")
	TradeTypeGovDeposit             = RegisterTranserType("gov-deposit", "", "Gov deposit")
	TradeTypeCopyrightVotes         = RegisterTranserType("cp-votes", "", "Community governance voting")
	TradeTypeCopyrightVoteReward    = RegisterTranserType("copyright-vote-reward", "", "Community governance incentives")
	TradeTypeUsdtForFsv             = RegisterTranserType("usdt-for-fsv", "usdtfsv", "Usdt versus FSV")
	TradeTypeCrossChainOut          = RegisterTranserType("cross-chain-out", "", "Cross Chain Out")
	TradeTypeCrossChainFee          = RegisterTranserType("cross-chain-fee", "", "Cross Chain Fee")
	TradeTypeCrossChainIn           = RegisterTranserType("cross-chain-in", "", "Cross Chain In")
	TradeTypeIbcTransferOut         = RegisterTranserType("ibc-transfer-out", "IBC", "IBC Transfer Out")
	TradeTypeIbcTransferIn          = RegisterTranserType("ibc-transfer-in", "IBC", "IBC Transfer In")
	TradeTypeEvmContractCall        = RegisterTranserType("evm-contract-call", "EVM", "Evm Contract Call")
	TradeTypeEvmTokenTransfer       = RegisterTranserType("evm-token-transfer", "EVM", "Evm Token Transfer")
	TradeTypeEvmContractDeploy      = RegisterTranserType("evm-contract-deploy", "EVM", "Evm Contract Deploy")
)

var tradeTypeText = make(map[string]string)
var tradeTypeTextEn = make(map[string]string)


func RegisterTranserType(key, value, enValue string) TranserType {
	tradeTypeTextEn[key] = enValue
	tradeTypeText[key] = value
	return TranserType(key)
}


func GetTranserTypeConfig() map[string]string {
	if trerr.Language == "EN" {
		return tradeTypeTextEn
	} else {
		return tradeTypeText
	}
}

type TranserType string

func (this TranserType) GetValue() string {
	if text, ok := tradeTypeText[string(this[:])]; ok {
		return text
	} else {
		return ""
	}
}

func (this TranserType) GetKey() string {
	return string(this[:])
}
