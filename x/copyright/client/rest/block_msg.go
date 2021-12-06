package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

type AnalysisMethod int

var AnalysisMethodExcludeFee = AnalysisMethod(1)

var AnalysisMethodOnlyFee = AnalysisMethod(2)

var AnalysisMethodDelegationReward = AnalysisMethod(3)

func AnalysisEvents(resp *types.MessageResp, prop *types.TxPropValue, transerType types.TranserType, analysisMethod AnalysisMethod) error {
	for _, event := range prop.Events {
		if event.Type == "transfer" {
			AnalysisTransfer(resp, prop, event, transerType, analysisMethod)
		}
	}
	return nil
}

func AnalysisTransfer(resp *types.MessageResp, prop *types.TxPropValue, event abciTypes.Event, transType types.TranserType, analysisMethod AnalysisMethod) {
	transData := types.TransferData{}
	transData.TradeType = transType
	transData.Height = prop.Height
	transData.Memo = prop.Memo
	transData.TxHash = prop.TxHash
	transData.BlockTime = prop.BlockTime
	transCoins := types.RealCoins{}

	for _, attr := range event.Attributes {
		if string(attr.GetKey()) == "amount" {

			coin, err := sdk.ParseCoinNormalized(string(attr.GetValue()))
			if err != nil {
				return
			}

			transCoins = append(transCoins, types.MustLedgerCoin2RealCoin(coin))
		} else if string(attr.GetKey()) == "sender" {
			transData.FromAddress = string(attr.GetValue())
		} else if string(attr.GetKey()) == "recipient" {
			transData.ToAddress = string(attr.GetValue())
		} else if string(attr.GetKey()) == bankTypes.AttributeKeyRecipientBalance {
			coins, err := sdk.ParseCoinsNormalized(string(attr.GetValue()))
			realCoins := types.MustLedgerCoins2RealCoins(coins)
			if err != nil {
				return
			}
			for _, coin := range realCoins {
				if coin.Denom == config.MainToken {
					transData.ToFsvBalance = coin.Amount
				} else {
					transData.ToTipBalance = coin.Amount
				}
			}
		} else if string(attr.GetKey()) == bankTypes.AttributeKeySenderBalance {
			coins, err := sdk.ParseCoinsNormalized(string(attr.GetValue()))
			realCoins := types.MustLedgerCoins2RealCoins(coins)
			if err != nil {
				return
			}
			for _, coin := range realCoins {
				if coin.Denom == config.MainToken {
					transData.FromFsvBalance = coin.Amount
				} else {
					transData.FromTipBalance = coin.Amount
				}
			}
		}
	}

	isFeeTransfer := false
	isDelegationReward := false
	isUnDelegationTransfer := false
	isCrossChainFee := false

	if transData.ToAddress == config.ContractAddressFee.String() {
		isFeeTransfer = true
	}

	if transData.FromAddress == config.ContractAddressStakingBonded.String() && transData.ToAddress == config.ContractAddressStakingNotBonded.String() {
		isUnDelegationTransfer = true
	}

	if transData.FromAddress == config.ContractAddressDistribution.String() {
		isDelegationReward = true
	}

	if transData.ToAddress == config.CrossChainFeeAccount {
		isCrossChainFee = true
	}

	if transData.FromAddress == config.ContractAddressDistribution.String() {
		isDelegationReward = true
	}

	if isDelegationReward {
		transData.TradeType = types.TradeTypeDelegationReward
	}

	if isUnDelegationTransfer {
		transData.TradeType = types.TradeTypeUnbondDelegation
	}

	if isCrossChainFee {
		transData.TradeType = types.TradeTypeCrossChainFee
	}

	if analysisMethod == AnalysisMethodExcludeFee && isFeeTransfer {
		return
	}

	if analysisMethod == AnalysisMethodOnlyFee && !isFeeTransfer {
		return
	}

	if analysisMethod == AnalysisMethodDelegationReward && !isDelegationReward {
		return
	}

	switch analysisMethod {
	case AnalysisMethodOnlyFee:
		transData.Fee = prop.Fee
		transData.Coins = types.NewRealCoinsFromStr(config.MainToken, "0")
		break
	case AnalysisMethodExcludeFee:
		transData.Fee = types.NewLedgerFeeZero()
		transData.Coins = transCoins
		break
	case AnalysisMethodDelegationReward:
		transData.Fee = prop.Fee
		transData.Coins = transCoins
	}
	resp.Transfers = append(resp.Transfers, transData)
}

var analysisHandle MessageAnalysisHandle

type MessageAnalysisHandle struct {
	clientCtx    *client.Context
	handles      map[string]func(sdk.Msg, *types.MessageResp, *types.TxPropValue) error
	eventHandles map[string]func(*types.MessageResp, *types.TxPropValue) error
}

func (this *MessageAnalysisHandle) register(msgType string, callback func(sdk.Msg, *types.MessageResp, *types.TxPropValue) error) {
	this.handles[msgType] = callback
}

func (this *MessageAnalysisHandle) updateMsgCreateValidatorToData(msg sdk.Msg, obj *stakingTypes.MsgCreateValidator, callbak func() error) error {
	msgByte := this.clientCtx.LegacyAmino.MustMarshalJSON(msg)
	err := this.clientCtx.LegacyAmino.UnmarshalJSON(msgByte, &obj)
	if err != nil {
		return err
	}
	return callbak()
}

func updateMsgToData(msg sdk.Msg, obj interface{}, callbak func() error) error {
	msgByte, err := util.Json.Marshal(msg)
	if err != nil {
		return err
	}
	err = util.Json.Unmarshal(msgByte, &obj)
	if err != nil {
		return err
	}
	return callbak()
}

func GetAmountFromEvent(event abciTypes.Event) []byte {
	for _, v := range event.Attributes {
		if string(v.Key) == sdk.AttributeKeyAmount {
			return v.Value
		}
	}
	return nil
}

func GetSharesFromEvent(event abciTypes.Event) []byte {
	for _, v := range event.Attributes {
		if string(v.Key) == sdk.AttributeKeyShares {
			return v.Value
		}
	}
	return nil
}

func GetBalanceFromEvent(event abciTypes.Event) (string, string, error) {
	fsvBalance := "0"
	tipBalance := "0"
	for _, v := range event.Attributes {
		if string(v.Key) == stakingTypes.AttributeKeyDelegatorBalance {
			coins, err := sdk.ParseCoinsNormalized(string(v.Value))
			if err != nil {
				return fsvBalance, tipBalance, err
			}
			realCoins := types.MustLedgerCoins2RealCoins(coins)
			for _, coin := range realCoins {
				if coin.Denom == config.MainToken {
					fsvBalance = coin.Amount
				} else {
					tipBalance = coin.Amount
				}
			}
			return fsvBalance, tipBalance, nil
		}
	}
	return fsvBalance, tipBalance, nil
}

func CreateAnalysisHandle(ctx *client.Context) {
	initAnalysisHandle(ctx)
	initEventHandle()
}

func initEventHandle() {

	analysisHandle.eventHandles[distributionTypes.EventTypeWithdrawCommission] = func(resp *types.MessageResp, prop *types.TxPropValue) error {
		feeAdded := false
		for _, event := range prop.Events {
			if event.Type == distributionTypes.EventTypeWithdrawCommission {
				fee := prop.Fee
				if feeAdded {
					fee = types.NewLedgerFeeZero()
				} else {
					feeAdded = true
				}
				accAddre := string(event.Attributes[1].Value)

				newTransfer := types.TransferData{ToAddress: accAddre, FromAddress: config.ContractAddressDistribution.String()}
				err := newTransfer.UpdateBalance(prop, newTransfer.FromAddress, newTransfer.ToAddress)
				if err != nil {
					return err
				}
				newTransfer.TradeType = types.TradeTypeCommissionReward
				newTransfer.Fee = fee
				newTransfer.BlockTime = prop.BlockTime
				newTransfer.TxHash = prop.TxHash
				newTransfer.Height = prop.Height
				amount := types.MustParseLedgerCoinFromStr(string(event.Attributes[0].Value))
				newTransfer.Coins = types.NewRealCoinsFromStr(config.MainToken, amount)
				resp.Transfers = append(resp.Transfers, newTransfer)
			}
		}
		return nil
	}

	analysisHandle.eventHandles[distributionTypes.EventTypeWithdrawRewards] = func(resp *types.MessageResp, prop *types.TxPropValue) error {
		feeAdded := false
		for _, event := range prop.Events {

			if event.Type == distributionTypes.EventTypeWithdrawRewards {
				fee := prop.Fee
				if feeAdded {
					fee = types.NewLedgerFeeZero()
				} else {
					feeAdded = true
				}
				newTransfer := types.TransferData{ToAddress: string(event.Attributes[2].Value), FromAddress: config.ContractAddressDistribution.String()}
				err := newTransfer.UpdateBalance(prop, newTransfer.FromAddress, newTransfer.ToAddress)
				if err != nil {
					return err
				}
				newTransfer.TradeType = types.TradeTypeDelegationReward
				newTransfer.Fee = fee
				newTransfer.BlockTime = prop.BlockTime
				newTransfer.TxHash = prop.TxHash
				newTransfer.Height = prop.Height
				amount := types.MustParseLedgerCoinFromStr(string(event.Attributes[0].Value))
				newTransfer.Coins = types.NewRealCoinsFromStr(config.MainToken, amount)
				resp.Transfers = append(resp.Transfers, newTransfer)
			}
		}
		return nil
	}

	analysisHandle.eventHandles[stakingTypes.EventTypeUnbond] = func(resp *types.MessageResp, prop *types.TxPropValue) error {
		isUnbond := false
		for _, event := range prop.Events {

			if event.Type == stakingTypes.EventTypeUnbond {
				if len(event.Attributes) == 5 {
					isUnbond = true
					data := types.UndelegationData{ValidatorAddress: string(event.Attributes[0].Value), DelegatorAddress: string(event.Attributes[4].Value)}
					data.UpdateTxBase(prop)
					ledgerDec := sdk.MustNewDecFromStr(string(event.Attributes[2].Value))
					data.Shares = types.MustParseLedgerDec(ledgerDec)
					data.CompletionTime = string(event.Attributes[3].Value)
					coins, err := strconv.ParseInt(string(event.Attributes[1].Value), 10, 64)
					if err != nil {
						continue
					}
					data.Amount = types.MustLedgerCoin2RealCoin(sdk.NewCoin(config.MainToken, sdk.NewInt(coins)))
					resp.Undelegations = append(resp.Undelegations, data)
				}
			}
		}
		if isUnbond {
			AnalysisEvents(resp, prop, types.TradeTypeUnbondDelegationFee, AnalysisMethodOnlyFee)
			AnalysisEvents(resp, prop, types.TradeTypeUnbondDelegation, AnalysisMethodExcludeFee)
		}
		return nil
	}

}

func initAnalysisHandle(ctx *client.Context) {

	analysisHandle = MessageAnalysisHandle{
		clientCtx:    ctx,
		handles:      make(map[string]func(sdk.Msg, *types.MessageResp, *types.TxPropValue) error),
		eventHandles: make(map[string]func(*types.MessageResp, *types.TxPropValue) error),
	}

	analysisHandle.register(stakingTypes.TypeMsgDelegate, func(msg sdk.Msg, resp *types.MessageResp, prop *types.TxPropValue) error {
		msgObj := stakingTypes.MsgDelegate{}
		return updateMsgToData(msg, &msgObj, func() error {
			data := types.NewDelegateData(msgObj)
			data.UpdateTxBase(prop)
			AnalysisEvents(resp, prop, types.TradeTypeDelegationReward, AnalysisMethodExcludeFee)
			for _, event := range prop.Events {
				if event.Type == stakingTypes.EventTypeDelegate {

					newTransfer := types.TransferData{Coins: types.NewRealCoins(data.Coin), FromAddress: msgObj.DelegatorAddress, ToAddress: config.ContractAddressDistribution.String()}
					newTransfer.TradeType = types.TradeTypeDelegation
					newTransfer.Fee = prop.Fee
					newTransfer.BlockTime = prop.BlockTime
					newTransfer.TxHash = prop.TxHash
					newTransfer.Height = prop.Height
					fsv, tip, err := GetBalanceFromEvent(event)
					if err != nil {
						return err
					}
					newTransfer.FromFsvBalance = fsv
					newTransfer.FromTipBalance = tip

					if shares := GetSharesFromEvent(event); shares != nil {
						data.Shares = types.MustParseLedgerDec(sdk.MustNewDecFromStr(string(shares)))
					}
					resp.Transfers = append(resp.Transfers, newTransfer)
				}
			}

			resp.Delegations = append(resp.Delegations, *data)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgRegisterCopyrightParty, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgRegisterCopyrightParty{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewCopyrightPartyData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			resp.CopyrightPartys = append(resp.CopyrightPartys, *data)

			AnalysisEvents(resp, params, types.TradeTypeCopyrightPartyRegister, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgInviteReward, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgInviteReward{}
		return updateMsgToData(msg, &msgObj, func() error {
			AnalysisEvents(resp, params, types.TradeTypeInviteReward, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCreateCopyright, func(msg sdk.Msg, resp *types.MessageResp, prop *types.TxPropValue) error {
		msgObj := types.MsgCreateCopyright{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewCopyrightData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(prop)
			resp.Copyrights = append(resp.Copyrights, *data)

			AnalysisEvents(resp, prop, types.TradeTypeCopyrightPublish, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(stakingTypes.TypeMsgCreateValidator, func(msg sdk.Msg, resp *types.MessageResp, prop *types.TxPropValue) error {
		msgObj := stakingTypes.MsgCreateValidator{}
		return analysisHandle.updateMsgCreateValidatorToData(msg, &msgObj, func() error {

			msgDelegate := stakingTypes.MsgDelegate{
				DelegatorAddress: msgObj.DelegatorAddress,
				ValidatorAddress: msgObj.ValidatorAddress,
				Amount:           msgObj.Value,
			}
			data := types.NewDelegateData(msgDelegate)
			data.UpdateTxBase(prop)
			data.Shares = data.Coin.Amount
			resp.Delegations = append(resp.Delegations, *data)

			newTransfer := types.TransferData{Coins: types.NewRealCoins(types.MustLedgerCoin2RealCoin(msgObj.Value)), FromAddress: msgObj.DelegatorAddress, ToAddress: config.ContractAddressStakingBonded.String()}
			err := newTransfer.UpdateBalance(prop, newTransfer.FromAddress, config.ContractAddressFee.String())
			if err != nil {
				return err
			}
			newTransfer.UpdateTxBase(prop)
			newTransfer.TradeType = types.TradeTypeValidatorCreate
			resp.Transfers = append(resp.Transfers, newTransfer)
			return nil
		})
	})

	analysisHandle.register(slashingTypes.TypeMsgUnjail, func(msg sdk.Msg, resp *types.MessageResp, prop *types.TxPropValue) error {
		msgObj := slashingTypes.MsgUnjail{}
		return updateMsgToData(msg, &msgObj, func() error {
			AnalysisEvents(resp, prop, types.TradeTypeValidatorUnjail, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgTransfer, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgTransfer{}
		return updateMsgToData(msg, &msgObj, func() error {
			data := types.NewTransferData(msgObj)
			err := data.UpdateBalance(params, data.FromAddress, data.ToAddress)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeTransfer, config.ContractAddressFee.String())

			err = data.UpdateFee(params)
			if err != nil {
				return err
			}
			resp.Transfers = append(resp.Transfers, *data)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgSpaceMiner, func(msg sdk.Msg, resp *types.MessageResp, prop *types.TxPropValue) error {
		msgObj := types.MsgSpaceMiner{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewSpaceMinerData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(prop)
			data.UpdateTradeBase(types.TradeTypeFsvDestroyMiner, config.ContractAddressDestory.String())
			resp.SpaceMiners = append(resp.SpaceMiners, *data)

			newTransfer := types.TransferData{FromAddress: msgObj.Creator, ToAddress: config.ContractAddressDestory.String()}
			err = newTransfer.UpdateBalance(prop, newTransfer.FromAddress, newTransfer.ToAddress)
			if err != nil {
				return err
			}
			newTransfer.TradeType = types.TradeTypeFsvDestroyMiner
			coin := types.RealCoin{}
			err = util.Json.UnmarshalFromString(msgObj.DeflationAmount, &coin)
			if err != nil {
				return err
			}
			newTransfer.Coins = types.NewRealCoinsFromStr(coin.Denom, coin.Amount)
			newTransfer.Fee = prop.Fee
			newTransfer.BlockTime = prop.BlockTime
			newTransfer.TxHash = prop.TxHash
			newTransfer.Height = prop.Height
			resp.Transfers = append(resp.Transfers, newTransfer)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgNftTransfer, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgNftTransfer{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewNftTransferData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeNftTransfer, config.ContractAddressFee.String())
			resp.NftTransfers = append(resp.NftTransfers, *data)

			AnalysisEvents(resp, params, types.TradeTypeNftTransfer, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgMortgage, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgMortgage{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewMortgageData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightBuyMortgage, config.ContractAddressFee.String())
			resp.Mortgages = append(resp.Mortgages, *data)

			return nil
		})
	})

	analysisHandle.register(types.TypeMsgEditorCopyright, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgEditorCopyright{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewEditorCopyrightData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightEdit, config.ContractAddressFee.String())
			resp.EditorCopyrights = append(resp.EditorCopyrights, *data)
			AnalysisEvents(resp, params, types.TradeTypeCopyrightEdit, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgDeleteCopyright, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgDeleteCopyright{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewDeleteCopyrightData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightDelete, config.ContractAddressFee.String())
			resp.DeleteCopyrights = append(resp.DeleteCopyrights, *data)
			AnalysisEvents(resp, params, types.TradeTypeCopyrightDelete, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCopyrightBonus, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgCopyrightBonus{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewBonusCopyrightData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightBuy, config.ContractAddressFee.String())
			resp.BonusCopyrights = append(resp.BonusCopyrights, *data)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCopyrightComplain, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgCopyrightComplain{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewCopyrightComplainData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightComplain, config.ContractAddressFee.String())
			resp.CopyrightComplains = append(resp.CopyrightComplains, *data)
			AnalysisEvents(resp, params, types.TradeTypeCopyrightComplain, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgComplainResponse, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgComplainResponse{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewComplainResponseData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			data.UpdateTradeBase(types.TradeTypeCopyrightComplainReply, config.ContractAddressFee.String())
			resp.ComplainResponses = append(resp.ComplainResponses, *data)
			AnalysisEvents(resp, params, types.TradeTypeCopyrightComplainReply, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgComplainVote, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgComplainVote{}
		return updateMsgToData(msg, &msgObj, func() error {
			data, err := types.NewComplainVoteData(msgObj)
			if err != nil {
				return err
			}
			data.UpdateTxBase(params)
			resp.ComplainVotes = append(resp.ComplainVotes, *data)
			AnalysisEvents(resp, params, types.TradeTypeCopyrightComplainVotes, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCopyrightVote, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		msgObj := types.MsgVoteCopyright{}
		return updateMsgToData(msg, &msgObj, func() error {
			AnalysisEvents(resp, params, types.TradeTypeCopyrightVotes, AnalysisMethodOnlyFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCrossChainIn, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		logs.Info("analysisHandle TypeMsgCrossChainIn")
		msgObj := types.MsgCrossChainIn{}
		return updateMsgToData(msg, &msgObj, func() error {
			AnalysisEvents(resp, params, types.TradeTypeCrossChainIn, AnalysisMethodExcludeFee)
			return nil
		})
	})

	analysisHandle.register(types.TypeMsgCrossChainOut, func(msg sdk.Msg, resp *types.MessageResp, params *types.TxPropValue) error {
		logs.Info("analysisHandle TypeMsgCrossChainOut")
		msgObj := types.MsgCrossChainOut{}
		return updateMsgToData(msg, &msgObj, func() error {
			AnalysisEvents(resp, params, types.TradeTypeCrossChainOut, AnalysisMethodExcludeFee)
			return nil
		})
	})
}

func parsingMsg(resp *types.MessageResp, msg sdk.Msg, params *types.TxPropValue) error {

	if callback, ok := analysisHandle.handles[msg.Type()]; ok {

		err := callback(msg, resp, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func parsingEvent(resp *types.MessageResp, params *types.TxPropValue) error {

	for _, callback := range analysisHandle.eventHandles {
		err := callback(resp, params)
		if err != nil {
			continue
		}
	}
	return nil
}
