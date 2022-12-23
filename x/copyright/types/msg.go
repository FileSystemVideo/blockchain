package types

import (
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	TypeMsgDistributeCommunityReward  = "distributeCommunityReward"
	TypeMsgDistributeDelegationReward = "distributeDelegationReward"
	TypeMsgDistributeCommissionReward = "distributeCommissionReward"
	TypeMsgCreateCopyright            = "createCopyright"
	TypeMsgRegisterCopyrightParty     = "registerCopyrightParty"
	TypeMsgSpaceMiner                 = "spaceMiner"
	TypeMsgDeflationVote              = "deflationVote"
	TypeMsgNftTransfer                = "nftTransfer"
	TypeMsgInviteCode                 = "inviteCode"
	TypeMsgMortgage                   = "mortgage"
	TypeMsgEditorCopyright            = "editorCopyright"
	TypeMsgDeleteCopyright            = "deleteCopyright"
	TypeMsgCopyrightBonus             = "copyrightBonus"
	TypeMsgCopyrightBonusV2           = "copyrightBonusV2"
	TypeMsgCopyrightComplain          = "copyrightComplain"
	TypeMsgComplainResponse           = "complainResponse"
	TypeMsgComplainVote               = "complainVote"
	TypeMsgAuthorizeAccount           = "authorizeAccount"
	TypeMsgTransfer                   = "copyrightTransfer"
	TypeMsgInviteReward               = "inviteReward"
	TypeMsgSpaceMinerReward           = "spaceMinerReward"
	TypeMsgCopyrightBonusRear         = "copyrightBonusRear"
	TypeMsgCopyrightBonusRearV2       = "copyrightBonusRearV2"
	TypeMsgCopyrightVote              = "copyrightVote"
	TypeMsgCrossChainIn               = "crossChainIn"
	TypeMsgCrossChainOut              = "crossChainOut"
)

var (
	_ sdk.Msg = &MsgCreateCopyright{}
	_ sdk.Msg = &MsgRegisterCopyrightParty{}
	_ sdk.Msg = &MsgSpaceMiner{}
	_ sdk.Msg = &MsgDistributeCommunityReward{}
	_ sdk.Msg = &MsgMortgage{}
	_ sdk.Msg = &MsgEditorCopyright{}
	_ sdk.Msg = &MsgDeleteCopyright{}
	_ sdk.Msg = &MsgCopyrightBonusV2{}
	_ sdk.Msg = &MsgCopyrightComplain{}
	_ sdk.Msg = &MsgComplainResponse{}
	_ sdk.Msg = &MsgComplainVote{}
	_ sdk.Msg = &MsgTransfer{}
	_ sdk.Msg = &MsgInviteReward{}
	_ sdk.Msg = &MsgSpaceMinerReward{}
	_ sdk.Msg = &MsgCopyrightBonusRearV2{}
	_ sdk.Msg = &MsgVoteCopyright{}
	_ sdk.Msg = &MsgCrossChainIn{}
	_ sdk.Msg = &MsgCrossChainOut{}
)

func NewMsgCrossChainOut(data CrossChainOutData) (*MsgCrossChainOut, error) {
	return &MsgCrossChainOut{
		SendAddress: data.SendAddress,
		ToAddress:   data.ToAddress,
		Coins:       data.Coins,
		ChainType:   data.ChainType,
		Remark:      data.Remark,
	}, nil
}

func (m MsgCrossChainOut) Route() string { return RouterKey }

func (m MsgCrossChainOut) Type() string { return TypeMsgCrossChainOut }

func (m MsgCrossChainOut) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.SendAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCrossChainOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCrossChainOut) ValidateBasic() error {
	return nil
}

func (m MsgCrossChainOut) XXX_MessageName() string {
	return TypeMsgCrossChainOut
}

func NewMsgCrossChainIn(data CrossChainInData) (*MsgCrossChainIn, error) {
	return &MsgCrossChainIn{
		SendAddress: data.SendAddress,
		Coins:       data.Coins,
		ChainType:   data.ChainType,
		Remark:      data.Remark,
	}, nil
}

func (m MsgCrossChainIn) Route() string { return RouterKey }

func (m MsgCrossChainIn) Type() string { return TypeMsgCrossChainIn }

func (m MsgCrossChainIn) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(core.CrossChainInManageAccount) //todo  
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCrossChainIn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCrossChainIn) ValidateBasic() error {
	return nil
}

func (m MsgCrossChainIn) XXX_MessageName() string {
	return TypeMsgCrossChainIn
}

func NewMsgCopyrightVote(data CopyrightVoteData) (*MsgVoteCopyright, error) {
	return &MsgVoteCopyright{
		Address:  data.Address,
		DataHash: data.DataHash,
		Power:    data.Power,
	}, nil
}

func (m MsgVoteCopyright) Route() string { return RouterKey }

func (m MsgVoteCopyright) Type() string { return TypeMsgCopyrightVote }

func (m MsgVoteCopyright) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgVoteCopyright) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgVoteCopyright) ValidateBasic() error {
	return nil
}

func (m MsgVoteCopyright) XXX_MessageName() string {
	return TypeMsgCopyrightVote
}
func NewMsgCopyrightBonusRearV2(data CopyrightBonusRearData) (*MsgCopyrightBonusRearV2, error) {
	return &MsgCopyrightBonusRearV2{
		Creator:           data.Downer.String(),
		Datahash:          data.DataHash,
		OfferAccountShare: data.OfferAccountShare,
		BonusAddress:      data.BonusAddress,
	}, nil
}

func (m MsgCopyrightBonusRearV2) Route() string { return RouterKey }

func (m MsgCopyrightBonusRearV2) Type() string { return TypeMsgCopyrightBonusRearV2 }

func (m MsgCopyrightBonusRearV2) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.BonusAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCopyrightBonusRearV2) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCopyrightBonusRearV2) ValidateBasic() error {
	return nil
}

func (m MsgCopyrightBonusRearV2) XXX_MessageName() string {
	return TypeMsgCopyrightBonusRearV2
}

func NewMsgSpaceMinerReward(data SpaceMinerRewardData) (*MsgSpaceMinerReward, error) {
	return &MsgSpaceMinerReward{
		Address: data.Address,
	}, nil
}

func (m MsgSpaceMinerReward) Route() string { return RouterKey }

func (m MsgSpaceMinerReward) Type() string { return TypeMsgSpaceMinerReward }

func (m MsgSpaceMinerReward) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgSpaceMinerReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgSpaceMinerReward) ValidateBasic() error {
	return nil
}

func (m MsgSpaceMinerReward) XXX_MessageName() string {
	return TypeMsgSpaceMinerReward
}

func NewMsgInviteReward(data InviteRewardData) (*MsgInviteReward, error) {
	return &MsgInviteReward{
		Address: data.Address,
	}, nil
}

func (m MsgInviteReward) Route() string { return RouterKey }

func (m MsgInviteReward) Type() string { return TypeMsgInviteReward }

func (m MsgInviteReward) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgInviteReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgInviteReward) ValidateBasic() error {
	return nil
}

func (m MsgInviteReward) XXX_MessageName() string {
	return TypeMsgInviteReward
}

func NewMsgTransfer(data TransferData) (*MsgTransfer, error) {
	coinsBytes, err := util.Json.Marshal(data.Coins)
	if err != nil {
		return nil, err
	}
	return &MsgTransfer{
		FromAddress: data.FromAddress,
		ToAddress:   data.ToAddress,
		Coins:       string(coinsBytes),
	}, nil
}

func (m MsgTransfer) Route() string { return RouterKey }

func (m MsgTransfer) Type() string { return TypeMsgTransfer }

func (m MsgTransfer) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgTransfer) ValidateBasic() error {
	var prices RealCoins
	err := util.Json.Unmarshal([]byte(m.Coins), &prices)
	if err != nil {
		return err
	}
	for _, price := range prices {
		//prices1fsv,2tip
		if price.Denom == core.MainToken {
			amountDec, err := decimal.NewFromString(price.Amount)
			if err != nil {
				return err
			}
			if len(prices) == 1 { //fsv
				if amountDec.LessThan(core.MinFsvTransfer) {
					return QuantityIllegalErr
				}
			} else if len(prices) == 2 { //tip   tip 、fsv
				if amountDec.LessThan(core.CopyrightInviteFee) {
					return QuantityIllegalErr
				}
			}
		}
		
		priceList := strings.Split(price.Amount, ".")
		if len(priceList) > 1 {
			if len(priceList[1]) > core.TransferAmountDecimalInputLength {
				return AmountDecimalMax18Err
			}
		}
	}

	return nil
}

func NewMsgComplainVote(data ComplainVoteData) (*MsgComplainVote, error) {
	return &MsgComplainVote{
		VoteAccount: data.VoteAccount.String(),
		VotePower:   data.VoteShare.String(),
		VoteStatus:  data.VoteStatus,
		ComplainId:  data.ComplainId,
	}, nil
}

func (m MsgComplainVote) Route() string { return RouterKey }

func (m MsgComplainVote) Type() string { return TypeMsgComplainVote }

func (m MsgComplainVote) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.VoteAccount)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgComplainVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgComplainVote) ValidateBasic() error {
	return nil
}

func (m MsgComplainVote) XXX_MessageName() string {
	return TypeMsgComplainVote
}

func NewMsgComplainResponse(data ComplainResponseData) (*MsgComplainResponse, error) {
	return &MsgComplainResponse{
		Datahash:      data.DataHash,
		AccuseInfor:   data.AccuseInfor,
		RemoteIp:      data.RemoteIp,
		Status:        data.Status,
		AccuseAccount: data.AccuseAccount.String(),
		ResponseTime:  data.ResponseTime,
		ComplainId:    data.ComplainId,
	}, nil
}

func (m MsgComplainResponse) Route() string { return RouterKey }

func (m MsgComplainResponse) Type() string { return TypeMsgComplainResponse }

func (m MsgComplainResponse) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.AccuseAccount)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgComplainResponse) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgComplainResponse) ValidateBasic() error {
	return nil
}

func (m MsgComplainResponse) XXX_MessageName() string {
	return TypeMsgComplainResponse
}

func NewMsgCopyrightComplain(data CopyrightComplainData) (*MsgCopyrightComplain, error) {
	return &MsgCopyrightComplain{
		Datahash:        data.DataHash,
		Author:          data.Author,
		Productor:       data.Productor,
		LegalNumber:     data.LegalNumber,
		LegalTime:       data.LegalTime,
		ComplainInfor:   data.ComplainInfor,
		ComplainAccount: data.ComplainAccount.String(),
		AccuseAccount:   data.AccuseAccount.String(),
		ComplainTime:    data.ComplainTime,
		ComplainId:      data.ComplainId,
	}, nil
}

func (m MsgCopyrightComplain) Route() string { return RouterKey }

func (m MsgCopyrightComplain) Type() string { return TypeMsgCopyrightComplain }

func (m MsgCopyrightComplain) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.ComplainAccount)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCopyrightComplain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCopyrightComplain) ValidateBasic() error {
	return nil
}

func (m MsgCopyrightComplain) XXX_MessageName() string {
	return TypeMsgCopyrightComplain
}

func NewMsgCopyrightBonusV2(data CopyrightBonusData) (*MsgCopyrightBonusV2, error) {
	return &MsgCopyrightBonusV2{
		Creator:           data.Downer.String(),
		Datahash:          data.DataHash,
		OfferAccountShare: data.OfferAccountShare,
		DataHashAccount:   data.HashAccount.String(),
		BonusType:         data.BonusType,
		BonusAddress:      data.BonusAddress,
	}, nil
}

func (m MsgCopyrightBonusV2) Route() string { return RouterKey }

func (m MsgCopyrightBonusV2) Type() string { return TypeMsgCopyrightBonusV2 }

func (m MsgCopyrightBonusV2) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCopyrightBonusV2) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCopyrightBonusV2) ValidateBasic() error {
	return nil
}

func (m MsgCopyrightBonusV2) XXX_MessageName() string {
	return TypeMsgCopyrightBonusV2
}

func NewMsgDeleteCopyright(data DeleteCopyrightData) (*MsgDeleteCopyright, error) {
	return &MsgDeleteCopyright{
		Creator:  data.Creator.String(),
		Datahash: data.DataHash,
	}, nil
}

func (m MsgDeleteCopyright) Route() string { return RouterKey }

func (m MsgDeleteCopyright) Type() string { return TypeMsgDeleteCopyright }

func (m MsgDeleteCopyright) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgDeleteCopyright) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgDeleteCopyright) ValidateBasic() error {
	return nil
}

func (m MsgDeleteCopyright) XXX_MessageName() string {
	return TypeMsgDeleteCopyright
}

func NewMsgEditorCopyright(data EditorCopyrightData) (*MsgEditorCopyright, error) {
	priceByte, err := util.Json.Marshal(data.Price)
	if err != nil {
		return nil, err
	}
	return &MsgEditorCopyright{
		Creator:    data.Creator.String(),
		Datahash:   data.DataHash,
		Name:       data.Name,
		ChargeRate: data.ChargeRate,
		Price:      string(priceByte),
		Ip:         data.Ip,
	}, nil
}

func (m MsgEditorCopyright) Route() string { return RouterKey }

func (m MsgEditorCopyright) Type() string { return TypeMsgEditorCopyright }

func (m MsgEditorCopyright) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgEditorCopyright) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgEditorCopyright) ValidateBasic() error {
	var price RealCoin
	err := util.Json.Unmarshal([]byte(m.Price), &price)
	if err != nil {
		return err
	}
	
	priceList := strings.Split(price.Amount, ".")
	if len(priceList) > 1 {
		if len(priceList[1]) > core.CopyrightPriceDecimalInputlLength {
			return CopyrightPirceErr
		}
	}
	
	chargeRateList := strings.Split(m.ChargeRate, ".")
	if len(chargeRateList) > 1 {
		if len(chargeRateList[1]) > core.CoryrightChargeRateDecimalInputlLength { 
			return CopyrightChargeRateErr
		}
	}

	return nil
}

func (m MsgEditorCopyright) XXX_MessageName() string {
	return TypeMsgEditorCopyright
}

func NewMsgMsgMortgage(data MortgageData) (*MsgMortgage, error) {
	priceByte, err := util.Json.Marshal(data.CopyrightPrice)
	if err != nil {
		return nil, err
	}
	mortgageByte, err := util.Json.Marshal(data.MortgageAmount)
	if err != nil {
		return nil, err
	}
	return &MsgMortgage{
		Creator:           data.Creator.String(),
		MortageAccount:    data.MortgageAccount.String(),
		DataHash:          data.DataHash,
		CopyrightPrice:    string(priceByte),
		CreateTime:        data.CreateTime,
		MortgageAmount:    string(mortgageByte),
		OfferAccountShare: data.OfferAccountShare,
		DataHashAccount:   data.DataHashAccount.String(),
		BonusType:         data.BonusType,
	}, nil
}

func (m MsgMortgage) Route() string { return RouterKey }

func (m MsgMortgage) Type() string { return TypeMsgMortgage }

func (m MsgMortgage) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.MortageAccount)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgMortgage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgMortgage) ValidateBasic() error {
	return nil
}

func (m MsgMortgage) XXX_MessageName() string {
	return TypeMsgMortgage
}

func NewMsgSpaceMiner(data SpaceMinerData) (*MsgSpaceMiner, error) {
	realCoinByte, err := util.Json.Marshal(data.DeflationAmount)
	if err != nil {
		return nil, err
	}
	return &MsgSpaceMiner{
		Creator:         data.Creator.String(),
		DeflationAmount: string(realCoinByte),
		AwardAccount:    data.AwardAccount.String(),
	}, nil
}

func (m MsgSpaceMiner) Route() string { return RouterKey }

func (m MsgSpaceMiner) Type() string { return TypeMsgSpaceMiner }

func (m MsgSpaceMiner) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgSpaceMiner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgSpaceMiner) ValidateBasic() error {
	var realCoin RealCoin
	
	err := util.Json.Unmarshal([]byte(m.DeflationAmount), &realCoin)
	if err != nil {
		return err
	}
	amountDec, err2 := sdk.NewDecFromStr(realCoin.Amount)
	if err2 != nil {
		return err2
	}
	if amountDec.LT(core.MinRealAmountDec) {
		return QuantityIllegalErr
	}
	return nil
}

func (m MsgSpaceMiner) XXX_MessageName() string {
	return TypeMsgSpaceMiner
}

func NewMsgNftTransfer(data NftTransferData) (*MsgNftTransfer, error) {
	return &MsgNftTransfer{
		From:    data.From.String(),
		To:      data.To.String(),
		TokenId: data.TokenId,
	}, nil
}

func (m MsgNftTransfer) Route() string { return RouterKey }

func (m MsgNftTransfer) Type() string { return TypeMsgNftTransfer }

func (m MsgNftTransfer) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.From)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgNftTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgNftTransfer) ValidateBasic() error {
	return nil
}

func (m MsgNftTransfer) XXX_MessageName() string {
	return TypeMsgNftTransfer
}

func NewMsgCreateCopyright(data CopyrightData) (*MsgCreateCopyright, error) {
	filesBytes, err := util.Json.Marshal(data.Files)
	if err != nil {
		return nil, err
	}
	priceByte, err := util.Json.Marshal(data.Price)
	if err != nil {
		return nil, err
	}
	linkMapByte, err := util.Json.Marshal(data.LinkMap)
	if err != nil {
		return nil, err
	}
	picLinkMapByte, err := util.Json.Marshal(data.PicLinkMap)
	if err != nil {
		return nil, err
	}

	return &MsgCreateCopyright{
		Datahash:       data.DataHash,
		Price:          string(priceByte),
		Creator:        data.Creator.String(),
		ResourceType:   data.ResourceType,
		PreHash:        data.PreHash,
		VideoHash:      data.VideoHash,
		Name:           data.Name,
		Files:          filesBytes,
		Size_:          data.Size,
		CreateTime:     int32(data.CreateTime),
		Password:       data.Password,
		ChargeRate:     data.ChargeRate,
		Ip:             data.Ip,
		OriginDataHash: data.OriginDataHash,
		ClassifyUid:    data.ClassifyUid,
		Ext:            data.Ext,
		LinkMap:        linkMapByte,
		PicLinkMap:     picLinkMapByte,
	}, nil
}

func (m MsgCreateCopyright) Route() string { return RouterKey }

func (m MsgCreateCopyright) Type() string { return TypeMsgCreateCopyright }

func (m MsgCreateCopyright) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgCreateCopyright) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgCreateCopyright) ValidateBasic() error {
	var price RealCoin
	err := util.Json.Unmarshal([]byte(m.Price), &price)
	if err != nil {
		return err
	}
	
	priceList := strings.Split(price.Amount, ".")
	if len(priceList) > 1 {
		if len(priceList[1]) > core.CopyrightPriceDecimalInputlLength {
			return CopyrightPirceErr
		}
	}
	
	chargeRateList := strings.Split(m.ChargeRate, ".")
	if len(chargeRateList) > 1 {
		if len(chargeRateList[1]) > core.CoryrightChargeRateDecimalInputlLength { 
			return CopyrightChargeRateErr
		}
	}
	return nil
}

func (m MsgCreateCopyright) XXX_MessageName() string {
	return TypeMsgCreateCopyright
}

func NewMsgRegisterCopyrightParty(data CopyrightPartyData) (*MsgRegisterCopyrightParty, error) {
	return &MsgRegisterCopyrightParty{
		Id:      data.Id,
		Intro:   data.Intro,
		Creator: data.Creator.String(),
		Author:  data.Author,
	}, nil
}

func (m MsgRegisterCopyrightParty) Route() string { return RouterKey }

func (m MsgRegisterCopyrightParty) Type() string { return TypeMsgRegisterCopyrightParty }

func (m MsgRegisterCopyrightParty) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Creator)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgRegisterCopyrightParty) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgRegisterCopyrightParty) ValidateBasic() error {
	return nil
}

func (m MsgRegisterCopyrightParty) XXX_MessageName() string {
	return TypeMsgRegisterCopyrightParty
}

func NewMsgDistributeCommunityReward(address sdk.AccAddress, coins string) (*MsgDistributeCommunityReward, error) {
	return &MsgDistributeCommunityReward{
		Address: address.String(),
		Amount:  coins,
	}, nil
}

func (m MsgDistributeCommunityReward) Route() string { return RouterKey }

func (m MsgDistributeCommunityReward) Type() string { return TypeMsgDistributeCommunityReward }

func (m MsgDistributeCommunityReward) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(m.Address)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	return addrs
}

func (m MsgDistributeCommunityReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&m)
	return sdk.MustSortJSON(bz)
}

func (m MsgDistributeCommunityReward) ValidateBasic() error {
	return nil
}

func (m MsgDistributeCommunityReward) XXX_MessageName() string {
	return TypeMsgDistributeCommunityReward
}
