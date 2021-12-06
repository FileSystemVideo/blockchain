package types

import (
	"fs.video/blockchain/x/copyright/config"
	"fs.video/trerr"
)



var (
	TradeTypeTransfer               = RegisterTranserType("transfer", "转账", "transfer accounts")
	TradeTypeCopyrightBuy           = RegisterTranserType("down", "版权购买", "Copyright purchase")
	TradeTypeCopyrightPublish       = RegisterTranserType("publish", "版权铸造费", "Copyright foundry fee")
	TradeTypeCopyrightEdit          = RegisterTranserType("edit", "修改版权", "Modify copyright")
	TradeTypeCopyrightSharesReward  = RegisterTranserType("bonus", "硬盘带宽分享奖励", "Hard disk bandwidth sharing reward")
	TradeTypeCopyrightDelete        = RegisterTranserType("delete", "删除版权", "Delete copyright")
	TradeTypeCopyrightSell          = RegisterTranserType("copyright-share", "版权销售收入", "Copyright sales revenue")
	TradeTypeCopyrightBuyMortgage   = RegisterTranserType("mining-mortg", "版权抵押购买", "Copyright mortgage purchase")
	TradeTypeCopyrightBuyRedeem     = RegisterTranserType("mining-redeem", "版权抵押购买赎回", "Copyright mortgage purchase and redemption")
	TradeTypeCopyrightPartyRegister = RegisterTranserType("bind-pay", "版权方注册", "Copyright registration")
	TradeTypeDelegation             = RegisterTranserType("bonded", "POS抵押", "POS mortgage")
	TradeTypeDelegationFee          = RegisterTranserType("bonded-fee", "POS抵押手续费", "POS mortgage service charge")
	TradeTypeUnbondDelegation       = RegisterTranserType("unbonded", "POS赎回", "POS redemption")
	TradeTypeUnbondDelegationFee    = RegisterTranserType("unbonded-fee", "POS赎回手续费", "POS redemption fee")
	TradeTypeFee                    = RegisterTranserType("fee", "手续费支出", "Service charge expenditure")
	TradeTypeDelegationReward       = RegisterTranserType("delegate-reward", "POS抵押奖励", "POS mortgage reward")
	TradeTypeCommissionReward       = RegisterTranserType("commission-reward", "POS佣金奖励", "POS Commission reward")
	TradeTypeValidatorUnjail        = RegisterTranserType("unjail", "POS解除监禁", "POS release")
	//TradeTypeNodeAuth               = RegisterTranserType("node-auth", "节点认证")
	TradeTypeFsvDestroyMiner = RegisterTranserType("space-miner", "销毁"+config.MainToken+"换空间", "Destruction for space")
	//TradeTypeFsvDestroyMinerVote    = RegisterTranserType("space-miner-vote", "销毁"+config.MainToken+"换空间投票")
	TradeTypeNftTransfer            = RegisterTranserType("nft-transfer", "nft转账", "NFT transfer")
	TradeTypeCopyrightComplain      = RegisterTranserType("cp-complain", "发起版权申诉", "Copyright complaint")
	TradeTypeSpaceMinerBonus        = RegisterTranserType("space-miner-bonus", "通缩挖矿分红", "Deflation dividend")
	TradeTypeCopyrightComplainReply = RegisterTranserType("cp-complain-reply", "版权应诉", "Copyright response")
	TradeTypeCopyrightComplainVotes = RegisterTranserType("cp-complain-votes", "版权申诉投票", "Copyright appeal vote")
	TradeTypeValidatorMinerBonus    = RegisterTranserType("validator-miner-bonus", "通缩挖矿DPOS奖励", "Verifier deflation dividend")
	TradeTypeValidatorCreate        = RegisterTranserType("validator-create", "创建验证器", "Create validator")
	TradeTypeInviteReward           = RegisterTranserType("invite-reward", "空间结算手续费", "Space settlement fee")
	TradeTypeCopyrightVotes         = RegisterTranserType("cp-votes", "版权审核投票", "Copyright audit voting")
	TradeTypeCopyrightVoteReward    = RegisterTranserType("copyright-vote-reward", "版权审核投票奖励", "Copyright audit voting Award")
	TradeTypeUsdtForFsv             = RegisterTranserType("usdt-for-fsv", "usdt兑fsv", "Usdt versus FSV")
	TradeTypeChainBck               = RegisterTranserType("chain-bck", "区块链空间挖矿回归", "Blockchain spatial mining regression")
	TradeTypeCrossChainOut          = RegisterTranserType("cross-chain-out", "跨链转出", "Cross Chain Out")
	TradeTypeCrossChainFee          = RegisterTranserType("cross-chain-fee", "跨链手续费", "Cross Chain Fee")
	TradeTypeCrossChainIn           = RegisterTranserType("cross-chain-in", "跨链转入", "Cross Chain In")
	TradeTypeGenesis          		= RegisterTranserType("genesis-in", "创世转入", "Chuangshi transfer in")
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
