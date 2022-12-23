package keeper

import (
	"context"
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/shopspring/decimal"
)

type msgServer struct {
	Keeper
	logPrefix string
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "chain | msgServer | "}
}

var (
	_ types.MsgServer = msgServer{}
)

// fsv->bsc
func (k msgServer) CrossChainOut(goCtx context.Context, msg *types.MsgCrossChainOut) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)

	saddress, err := sdk.AccAddressFromBech32(msg.SendAddress)
	if err != nil {
		log.WithError(err).WithField("SendAddress", msg.SendAddress).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	caddress, err := sdk.AccAddressFromBech32(core.CrossChainAccount)
	if err != nil {
		log.WithError(err).WithField("CrossChainAccount", core.CrossChainAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	faddress, err := sdk.AccAddressFromBech32(core.CrossChainFeeAccount)
	if err != nil {
		log.WithError(err).WithField("CrossChainFeeAccount", core.CrossChainFeeAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	amount, err := util.StringToCoinWithRate(msg.Coins)
	if err != nil {
		log.WithError(err).WithField("Coins", msg.Coins).Error("StringToCoinWithRate")
		return &types.MsgEmptyResponse{}, err
	}
	
	if amount.Amount.ToDec().LT(sdk.MustNewDecFromStr(core.CrossChainOutMinAmount)) {
		return &types.MsgEmptyResponse{}, errors.New("transfer amount is too low")
	}
	
	//lvdb kv,
	feeRatio, err := k.GetCrossChainOutFeeRatio(ctx)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}

	fee := amount.Amount.ToDec().Mul(sdk.MustNewDecFromStr(feeRatio))
	err = k.CoinKeeper.SendCoins(ctx, saddress, faddress, sdk.NewCoins(sdk.NewCoin(amount.Denom, fee.TruncateInt())))
	if err != nil {
		log.WithError(err).Error("SendCoins1")
		return &types.MsgEmptyResponse{}, err
	}
	
	//err = k.CoinKeeper.SendCoinsFromAccountToModule(ctx, saddress, types.ContractCrossChain, sdk.NewCoins(amount))
	freezeAmount := amount.Amount.ToDec().Sub(fee)
	err = k.CoinKeeper.SendCoins(ctx, saddress, caddress, sdk.NewCoins(sdk.NewCoin(amount.Denom, freezeAmount.TruncateInt())))
	if err != nil {
		log.WithError(err).Error("SendCoins2")
		return &types.MsgEmptyResponse{}, err
	}

	_, coinSymbol, err := util.StringDenom(msg.Coins)
	if err != nil {
		log.WithError(err).WithField("Coins", msg.Coins).Error("StringDenom")
		return &types.MsgEmptyResponse{}, err
	}

	//fmt.Println(":", freezeAmount.TruncateInt().String())

	//event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"cross_chain",
			sdk.NewAttribute("from_address", msg.SendAddress),
			sdk.NewAttribute("to_address", msg.ToAddress),
			sdk.NewAttribute("coin_num", freezeAmount.TruncateInt().String()),
			sdk.NewAttribute("coin_symbol", coinSymbol),
			sdk.NewAttribute("chain_type", msg.ChainType),
			sdk.NewAttribute("remark", msg.Remark),
		),
	)

	return &types.MsgEmptyResponse{}, nil
}

// bsc->fsv
func (k msgServer) CrossChainIn(goCtx context.Context, msg *types.MsgCrossChainIn) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	saddress, err := sdk.AccAddressFromBech32(msg.SendAddress)
	if err != nil {
		log.WithError(err).WithField("SendAddress", msg.SendAddress).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	caddress, err := sdk.AccAddressFromBech32(core.CrossChainAccount)
	if err != nil {
		log.WithError(err).WithField("CrossChainAccount", core.CrossChainAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	amount, err := util.StringToCoinWithRate(msg.Coins)
	if err != nil {
		log.WithError(err).WithField("Coins", msg.Coins).Error("StringToCoinWithRate")
		return &types.MsgEmptyResponse{}, err
	}
	//err = k.CoinKeeper.SendCoinsFromModuleToAccount(ctx, types.ContractCrossChain, saddress, sdk.NewCoins(amount))

	err = k.CoinKeeper.SendCoins(ctx, caddress, saddress, sdk.NewCoins(amount))

	if err != nil {
		log.WithError(err).Error("SendCoins")
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, nil
}


func (k msgServer) SpaceMinerReward(goCtx context.Context, msg *types.MsgSpaceMinerReward) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		log.WithError(err).WithField("Address", msg.Address).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	
	err = k.SpaceMinerRewardSettlement(ctx, address.String())
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) InviteReward(goCtx context.Context, msg *types.MsgInviteReward) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgEmptyResponse{}, nil
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		log.WithError(err).WithField("Address", msg.Address).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	
	err = k.RewardSettlement(ctx, address.String())
	if err != nil {
		log.WithError(err).WithField("address", address.String()).Error("RewardSettlement")
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) Transfer(goCtx context.Context, msg *types.MsgTransfer) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddress, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		log.WithError(err).WithField("FromAddress", msg.FromAddress).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	var toAddress sdk.AccAddress
	if msg.ToAddress == core.ContractAddressDestory.String() {
		toAddress = core.ContractAddressDestory
	} else {
		toAddress, err = sdk.AccAddressFromBech32(msg.ToAddress)
		if err != nil {
			log.WithError(err).WithField("ToAddress", msg.ToAddress).Error("AccAddressFromBech32")
			return &types.MsgEmptyResponse{}, err
		}
	}

	var realCoins types.RealCoins
	err = util.Json.Unmarshal([]byte(msg.Coins), &realCoins)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return &types.MsgEmptyResponse{}, err
	}
	coins := types.MustRealCoins2LedgerCoins(realCoins)
	
	if k.CoinKeeper.BlockedAddr(toAddress) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	err = k.Keeper.Transfer(ctx, fromAddress, toAddress, coins)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) ComplainVote(goCtx context.Context, msg *types.MsgComplainVote) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	voteAccount, err := sdk.AccAddressFromBech32(msg.VoteAccount)
	if err != nil {
		log.WithError(err).WithField("VoteAccount", msg.VoteAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	power, err := sdk.NewDecFromStr(msg.VotePower)
	if err != nil {
		log.WithError(err).WithField("VotePower", msg.VotePower).Error("NewDecFromStr")
		return &types.MsgEmptyResponse{}, err
	}
	data := types.ComplainVoteData{
		VoteAccount: voteAccount,
		VoteStatus:  msg.VoteStatus,
		ComplainId:  msg.ComplainId,
		VoteShare:   power,
	}
	err = k.Keeper.ComplainVote(ctx, data)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) ComplainResponse(goCtx context.Context, msg *types.MsgComplainResponse) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	/*complainAccount, err := sdk.AccAddressFromBech32(msg.)
	if err != nil {
		log.Error("", err)
		return &types.MsgEmptyResponse{}, err
	}*/
	accuseAccount, err := sdk.AccAddressFromBech32(msg.AccuseAccount)
	if err != nil {
		log.WithError(err).WithField("AccuseAccount", msg.AccuseAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	data := types.ComplainResponseData{
		DataHash:      msg.Datahash,
		AccuseInfor:   msg.AccuseInfor,
		RemoteIp:      msg.RemoteIp,
		Status:        msg.Status,
		AccuseAccount: accuseAccount,
		ComplainId:    msg.ComplainId,
		ResponseTime:  msg.ResponseTime,
	}
	err = k.Keeper.ComplainResponse(ctx, data)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) CopyrightComplain(goCtx context.Context, msg *types.MsgCopyrightComplain) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	complainAccount, err := sdk.AccAddressFromBech32(msg.ComplainAccount)
	if err != nil {
		log.WithError(err).WithField("ComplainAccount", msg.ComplainAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	
	accuseAccount, err := sdk.AccAddressFromBech32(msg.AccuseAccount)
	if err != nil {
		log.WithError(err).WithField("AccuseAccount", msg.AccuseAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	if k.ComplainHashStatus(ctx, msg.Datahash) {
		return &types.MsgEmptyResponse{}, types.CopyrightComplainErr
	}

	data := types.CopyrightComplainData{
		DataHash:        msg.Datahash,
		Author:          msg.Author,
		Productor:       msg.Productor,
		LegalTime:       msg.LegalTime,
		LegalNumber:     msg.LegalNumber,
		ComplainInfor:   msg.ComplainInfor,
		ComplainAccount: complainAccount,
		AccuseAccount:   accuseAccount,
		ComplainId:      msg.ComplainId,
		ComplainTime:    msg.ComplainTime,
	}
	err = k.Keeper.CopyrightComplain(ctx, data)
	return &types.MsgEmptyResponse{}, err
}

//V2
func (k msgServer) CopyrightBonusV2(goCtx context.Context, msg *types.MsgCopyrightBonusV2) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	dataHashAccount, err := sdk.AccAddressFromBech32(msg.DataHashAccount)
	if err != nil {
		log.WithError(err).WithField("DataHashAccount", msg.DataHashAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	data := types.CopyrightBonusData{
		Downer:            creator,
		DataHash:          msg.Datahash,
		OfferAccountShare: msg.OfferAccountShare,
		HashAccount:       dataHashAccount,
		BonusType:         msg.BonusType,
	}
	data.Fee = types.NewLedgerFee(core.ChainDefaultFee)
	err = k.Keeper.CopyrightBonus(ctx, data)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	
	err = k.SetBonusAddress(ctx, msg.BonusAddress, creator.String(), ctx.BlockHeight())
	if err != nil {
		log.WithError(err).Error("SetBonusAddress")
	}
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) CopyrightBonusRearV2(goCtx context.Context, msg *types.MsgCopyrightBonusRearV2) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	copyrightByte, err := k.GetCopyright(ctx, msg.Datahash)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightByte, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return &types.MsgEmptyResponse{}, err
	}
	txhash := formatTxHash(ctx.TxBytes())

	
	if k.IsExistBonusAddress(ctx, msg.BonusAddress) {
		err = dealBonusAuthrizeLogic(ctx, k.Keeper, msg.OfferAccountShare, txhash, "", core.KeyCopyrightBonus, ctx.BlockHeight(), copyrightData)
		if err != nil {
			log.WithError(err).Error("dealBonusAuthrizeLogic")
			return &types.MsgEmptyResponse{}, err
		}
		//key
		k.DeleteBonusAddress(ctx, msg.BonusAddress)
	}
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) CopyrightVote(goCtx context.Context, msg *types.MsgVoteCopyright) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	copyrightByte, err := k.GetCopyright(ctx, msg.DataHash)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightByte, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return &types.MsgEmptyResponse{}, err
	}

	if copyrightData.ApproveStatus != 0 {
		return &types.MsgEmptyResponse{}, types.CopyrightApproveHasFinished
	}
	txhash := formatTxHash(ctx.TxBytes())
	err = k.dealCopyrightVote(ctx, msg.Address, copyrightData.Name, copyrightData.DataHash, msg.Power, txhash, copyrightData.LinkMap)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) DeleteCopyright(goCtx context.Context, msg *types.MsgDeleteCopyright) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	data := types.DeleteCopyrightData{
		Creator:  creator,
		DataHash: msg.Datahash,
	}
	err = k.Keeper.DeleteCopyright(ctx, data)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) EditorCopyright(goCtx context.Context, msg *types.MsgEditorCopyright) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(msg.Price), &price)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return &types.MsgEmptyResponse{}, err
	}
	data := types.EditorCopyrightData{
		Creator:    creator,
		Name:       msg.Name,
		DataHash:   msg.Datahash,
		ChargeRate: msg.ChargeRate,
		Price:      price,
		Ip:         msg.Ip,
	}
	err = k.Keeper.EditorCopyright(ctx, data)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) Mortgage(goCtx context.Context, msg *types.MsgMortgage) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	mortgageAccount, err := sdk.AccAddressFromBech32(msg.MortageAccount)
	if err != nil {
		log.WithError(err).WithField("MortageAccount", msg.MortageAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	datahashAccount, err := sdk.AccAddressFromBech32(msg.DataHashAccount)
	if err != nil {
		log.WithError(err).WithField("DataHashAccount", msg.DataHashAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	/*copyrightPrice, err := strconv.ParseFloat(msg.CopyrightPrice, 64)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	price := types.NewLedgerDec(copyrightPrice)*/
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(msg.CopyrightPrice), &price)
	if err != nil {
		log.WithError(err).Error("Unmarshal1")
		return &types.MsgEmptyResponse{}, err
	}
	var mortgageAmount types.RealCoin
	err = util.Json.Unmarshal([]byte(msg.MortgageAmount), &mortgageAmount)
	if err != nil {
		log.WithError(err).Error("Unmarshal2")
		return &types.MsgEmptyResponse{}, err
	}
	mortgAmountDecimal, err := decimal.NewFromString(mortgageAmount.Amount)
	if err != nil {
		log.WithError(err).WithField("mortgageAmount", mortgageAmount.Amount).Error("NewFromString")
		return &types.MsgEmptyResponse{}, err
	}
	fee := mortgAmountDecimal.Mul(core.MortgageFee)
	data := types.MortgageData{
		MortgageAccount:   mortgageAccount,
		Creator:           creator,
		DataHashAccount:   datahashAccount,
		CopyrightPrice:    price,
		MortgageAmount:    mortgageAmount,
		DataHash:          msg.DataHash,
		CreateTime:        msg.CreateTime,
		OfferAccountShare: msg.OfferAccountShare,
		BonusType:         msg.BonusType,
	}
	data.Fee = types.NewLedgerFee(fee)
	err = k.Keeper.Mortgage(ctx, data)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) DistributeCommunityReward(goCtx context.Context, msg *types.MsgDistributeCommunityReward) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	coins, err := sdk.ParseCoinsNormalized(msg.Amount)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	if msg.Address != core.CommunityRewardAccount { 
		return &types.MsgEmptyResponse{}, types.UnAuthorizedAccountError
	}
	accountAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	err = k.distributionKeeper.DistributeFromFeePool(ctx, coins, accountAddress)
	return &types.MsgEmptyResponse{}, err
}


func (k msgServer) SpaceMiner(goCtx context.Context, msg *types.MsgSpaceMiner) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	
	awardAccount, err := sdk.AccAddressFromBech32(msg.AwardAccount)
	if err != nil {
		log.WithError(err).WithField("AwardAccount", msg.AwardAccount).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	var realCoin types.RealCoin
	
	err = util.Json.Unmarshal([]byte(msg.DeflationAmount), &realCoin)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return &types.MsgEmptyResponse{}, err
	}
	
	if realCoin.Denom != sdk.DefaultBondDenom {
		return &types.MsgEmptyResponse{}, types.OnlyMainTokenErr
	}
	
	data := types.SpaceMinerData{
		Creator:         creator,
		DeflationAmount: realCoin,
		AwardAccount:    awardAccount,
	}
	return &types.MsgEmptyResponse{}, k.Keeper.AddSpaceMiner(ctx, data)
}

func (k msgServer) NftTransfer(goCtx context.Context, msg *types.MsgNftTransfer) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.From).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.To).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}
	if msg.TokenId == "" {
		return &types.MsgEmptyResponse{}, types.TokenIdEmpty
	}
	if k.ComplainHashStatus(ctx, msg.TokenId) {
		return &types.MsgEmptyResponse{}, types.CopyrightComplainErr
	}
	data := types.NftTransferData{
		From:    from,
		To:      to,
		TokenId: msg.TokenId,
	}
	return &types.MsgEmptyResponse{}, k.Keeper.HandleNftTransfer(ctx, data)
}


func (k msgServer) RegisterCopyrightParty(goCtx context.Context, msg *types.MsgRegisterCopyrightParty) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		log.WithError(err).WithField("Creator", msg.Creator).Error("AccAddressFromBech32")
		return &types.MsgEmptyResponse{}, err
	}

	data := types.CopyrightPartyData{
		Creator: creator,
		Id:      msg.Id,
		Intro:   msg.Intro,
		Author:  msg.Author,
	}
	return &types.MsgEmptyResponse{}, k.Keeper.SetCopyrightParty(ctx, data)
}


func (k msgServer) CreateCopyright(goCtx context.Context, msg *types.MsgCreateCopyright) (*types.MsgEmptyResponse, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainMsgServer)
	ctx := sdk.UnwrapSDKContext(goCtx)
	data, err := types.NewCopyrightData(*msg)
	if err != nil {
		log.WithError(err).Error("NewCopyrightData")
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, k.Keeper.SetCopyright(ctx, *data)
}
