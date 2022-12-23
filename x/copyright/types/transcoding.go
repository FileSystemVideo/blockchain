package types

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

/**

*/

//dpos msg  data
func NewWithdrawDelegatorRewardData(msg distributionTypes.MsgWithdrawDelegatorReward) *WithdrawDelegatorRewardData {
	data := WithdrawDelegatorRewardData{
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
	}
	return &data
}

//dpos msg  data
func NewUndelegateData(msg stakingTypes.MsgUndelegate) *UndelegationData {
	data := UndelegationData{
		Amount:           MustLedgerCoin2RealCoin(msg.Amount),
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
	}
	return &data
}

//dpos msg  data
func NewRedelegateData(msg stakingTypes.MsgBeginRedelegate) *RedelegationData {
	data := RedelegationData{
		Amount:              MustLedgerCoin2RealCoin(msg.Amount),
		DelegatorAddress:    msg.DelegatorAddress,
		ValidatorDstAddress: msg.ValidatorDstAddress,
		ValidatorSrcAddress: msg.ValidatorSrcAddress,
	}
	return &data
}

//dpos msg  data
func NewDelegateData(msg stakingTypes.MsgDelegate) *DelegationData {
	data := DelegationData{
		Coin:             MustLedgerCoin2RealCoin(msg.Amount),
		DelegatorAddress: msg.DelegatorAddress,
		ValidatorAddress: msg.ValidatorAddress,
	}
	return &data
}

// msg  data
func NewTransferData(msg MsgTransfer) *TransferData {
	var realCoins RealCoins
	util.Json.Unmarshal([]byte(msg.Coins), &realCoins)
	data := TransferData{
		FromAddress: msg.FromAddress,
		ToAddress:   msg.ToAddress,
		Coins:       realCoins,
	}
	return &data
}

// msg  data
func NewSpaceMinerData(msg MsgSpaceMiner) (*SpaceMinerData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return nil, err
	}
	awardAccount, err := sdk.AccAddressFromBech32(msg.AwardAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return nil, err
	}
	var realCoin RealCoin
	err = util.Json.Unmarshal([]byte(msg.DeflationAmount), &realCoin)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	data := SpaceMinerData{
		DeflationAmount: realCoin,
		Creator:         creator,
		AwardAccount:    awardAccount,
	}
	return &data, nil
}

// msg  data
func NewNftTransferData(msg MsgNftTransfer) (*NftTransferData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return nil, err
	}

	data := NftTransferData{
		From:    from,
		To:      to,
		TokenId: msg.TokenId,
	}
	return &data, nil
}

// msg  data
func NewCopyrightPartyData(msg MsgRegisterCopyrightParty) (*CopyrightPartyData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return nil, err
	}
	data := CopyrightPartyData{
		Id:      msg.Id,
		Intro:   msg.Intro,
		Creator: creator,
		Author:  msg.Author,
	}
	return &data, nil
}

// msg  data
func NewCopyrightData(msg MsgCreateCopyright) (*CopyrightData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return nil, err
	}
	var files Files
	err = util.Json.Unmarshal(msg.Files, &files)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal 1")
		return nil, err
	}
	var price RealCoin
	err = util.Json.Unmarshal([]byte(msg.Price), &price)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal 2")
		return nil, err
	}

	var linkMap map[string]Link
	err = util.Json.Unmarshal(msg.LinkMap, &linkMap)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal 3")
		return nil, err
	}
	var picLinkMap map[string]Link
	err = util.Json.Unmarshal(msg.PicLinkMap, &picLinkMap)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal 4")
		return nil, err
	}
	data := CopyrightData{
		DataHash:       msg.Datahash,
		Price:          price,
		Creator:        creator,
		ResourceType:   msg.ResourceType,
		PreHash:        msg.PreHash,
		VideoHash:      msg.VideoHash,
		Name:           msg.Name,
		Files:          files,
		Size:           msg.Size_,
		CreateTime:     int(msg.CreateTime),
		Password:       msg.Password,
		ChargeRate:     msg.ChargeRate,
		Ip:             msg.Ip,
		OriginDataHash: msg.OriginDataHash,
		ClassifyUid:    msg.ClassifyUid,
		Ext:            msg.Ext,
		LinkMap:        linkMap,
		ApproveStatus:  0,
		PicLinkMap:     picLinkMap,
	}
	return &data, nil
}

// msg  data
func NewEditorCopyrightData(msg MsgEditorCopyright) (*EditorCopyrightData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return nil, err
	}
	var price RealCoin
	err = util.Json.Unmarshal([]byte(msg.Price), &price)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return nil, err
	}
	data := EditorCopyrightData{
		DataHash:   msg.Datahash,
		Price:      price,
		Creator:    creator,
		Name:       msg.Name,
		ChargeRate: msg.ChargeRate,
		Ip:         msg.Ip,
	}
	return &data, nil
}

// msg  data
func NewDeleteCopyrightData(msg MsgDeleteCopyright) (*DeleteCopyrightData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return nil, err
	}
	data := DeleteCopyrightData{
		DataHash: msg.Datahash,
		Creator:  creator,
	}
	return &data, nil
}

// msg  data
func NewCopyrightComplainData(msg MsgCopyrightComplain) (*CopyrightComplainData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	complainAccount, err := sdk.AccAddressFromBech32(msg.ComplainAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return nil, err
	}
	accuseAccount, err := sdk.AccAddressFromBech32(msg.AccuseAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return nil, err
	}
	data := CopyrightComplainData{
		DataHash:        msg.Datahash,
		ComplainAccount: complainAccount,
		AccuseAccount:   accuseAccount,
		Author:          msg.Author,
		Productor:       msg.Productor,
		LegalTime:       msg.LegalTime,
		LegalNumber:     msg.LegalNumber,
		ComplainInfor:   msg.ComplainInfor,
		ComplainId:      msg.ComplainId,
		ComplainTime:    msg.ComplainTime,
	}
	return &data, nil
}

// msg  data
func NewComplainResponseData(msg MsgComplainResponse) (*ComplainResponseData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	accuseAccount, err := sdk.AccAddressFromBech32(msg.AccuseAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return nil, err
	}
	data := ComplainResponseData{
		DataHash:      msg.Datahash,
		AccuseAccount: accuseAccount,
		Status:        msg.Status,
		RemoteIp:      msg.RemoteIp,
		AccuseInfor:   msg.AccuseInfor,
		ComplainId:    msg.ComplainId,
		ResponseTime:  msg.ResponseTime,
	}
	return &data, nil
}

// msg  data
func NewComplainVoteData(msg MsgComplainVote) (*ComplainVoteData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	voteAccount, err := sdk.AccAddressFromBech32(msg.VoteAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return nil, err
	}
	decPower, err := sdk.NewDecFromStr(msg.VotePower)
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return nil, err
	}
	data := ComplainVoteData{
		VoteAccount: voteAccount,
		ComplainId:  msg.ComplainId,
		VoteShare:   decPower,
		VoteStatus:  msg.VoteStatus,
	}
	return &data, nil
}

// msg  data
func NewMortgageData(msg MsgMortgage) (*MortgageData, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainType)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 1")
		return nil, err
	}
	mortgageAccount, err := sdk.AccAddressFromBech32(msg.MortageAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 2")
		return nil, err
	}
	dataAccount, err := sdk.AccAddressFromBech32(msg.DataHashAccount)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32 3")
		return nil, err
	}

	var mortgageAmount RealCoin
	err = util.Json.Unmarshal([]byte(msg.MortgageAmount), &mortgageAmount)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return nil, err
	}
	var copyrightPrice RealCoin
	err = util.Json.Unmarshal([]byte(msg.CopyrightPrice), &copyrightPrice)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return nil, err
	}
	data := MortgageData{
		Creator:           creator,
		MortgageAmount:    mortgageAmount,
		MortgageAccount:   mortgageAccount,
		DataHash:          msg.DataHash,
		OfferAccountShare: msg.OfferAccountShare,
		DataHashAccount:   dataAccount,
		CopyrightPrice:    copyrightPrice,
		CreateTime:        msg.CreateTime,
	}
	return &data, nil
}

// msg  data
func NewCopyrightVoteData(msg MsgVoteCopyright) (*CopyrightVoteData, error) {
	data := CopyrightVoteData{
		Address:  msg.Address,
		DataHash: msg.DataHash,
		Power:    msg.Power,
	}

	return &data, nil
}
