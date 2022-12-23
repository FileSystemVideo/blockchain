package rest

import (
	errors "errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/shopspring/decimal"
	"time"
)

/*
******************************************************************
**    ，client grpc  **
******************************************************************
 */


func CreateCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyright types.MsgCreateCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	log.Debug("do")

	addr, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(copyright.Price), &price)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return errors.New(ParseCoinError)
	}
	flag := util.JudgeAmount(price.Amount)
	if !flag {
		log.Warn("JudgeAmount fail")
		return errors.New(ParseCoinError)
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(price)
	var feeDec sdk.Dec
	if fee.Amount.Len() > 0 {
		for i := 0; i < fee.Amount.Len(); i++ {
			coin := fee.Amount[i]
			if coin.Denom == sdk.DefaultBondDenom {
				feeDec = coin.Amount.ToDec()
				break
			}
		}
	} else {
		return errors.New(FeeCannotEmpty)
	}
	
	now := time.Now()
	dayString := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	pubCount, err := grpcQueryPubCount(ctx, dayString)
	if err != nil {
		log.WithError(err).Error("grpcQueryPubCount")
		return err
	}
	feeString := util.FilePrice(pubCount + 1)
	feeDecimal, err := decimal.NewFromString(feeString)
	if err != nil {
		log.WithError(err).Error("ParseFloat")
		return errors.New(FormatStringToIntError)
	}
	pubFeeDec := types.NewLedgerDec(feeDecimal)
	//fmt.Println("", pubFeeDec.String(), feeDec.String())
	if !pubFeeDec.Equal(feeDec) {
		log.Error("!pubFeeDec.Equal(feeDec) error")
		return errors.New(FeeInvalidError)
	}
	
	if CheckMinAmount(ledgeCoin.Amount.ToDec()) {
		log.Error("CheckMinAmount error")
		return errors.New(LowerMinCoin)
	}
	
	balStatus, errStr := judgeBalance(ctx, addr, feeDec, core.MainToken)
	if !balStatus {
		log.Error("judgeBalance fail")
		return errors.New(errStr)
	}

	
	if copyright.Password == "" || len(copyright.Password) != 64 {
		return errors.New(PasswordGroupError)
	}

	if copyright.Datahash == "" {
		return errors.New(DataHashEmpty)
	}

	if copyright.OriginDataHash == "" {
		return errors.New(OriginDataHashEmpty)
	}

	if copyright.ClassifyUid < 0 {
		return errors.New(ClassifyIdInvalid)
	}
	exist, err := grpcQueryCopyrightExist(ctx, copyright.Datahash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightExist")
		return err
	}
	if exist {
		return errors.New(OriginDataHashExist)
	}
	exist, err = grpcQueryCopyrightOrigin(ctx, copyright.OriginDataHash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyrightOrigin")
		return err
	}
	if exist {
		return errors.New(OriginDataHashExist)
	}

	
	if copyright.ChargeRate == "" {
		return errors.New(ChargeRateEmpty)
	}
	chargeRate, err := sdk.NewDecFromStr(copyright.ChargeRate)
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return errors.New(DecimalFromStringError)
	}
	if chargeRate.LT(core.ChargeRateLow) {
		return errors.New(ChargeRateTooSmall)
	}
	if chargeRate.GT(core.ChargeRateHigh) {
		return errors.New(ChargeRateTooHigh)
	}

	enough, err := grpcQueryAccountSpace(ctx, copyright.Size_, copyright.Creator)
	if err != nil {
		log.WithError(err).Error("grpcQueryAccountSpace")
		return err
	}
	if !enough {
		return errors.New(SpaceNotEnough)
	}
	return nil
}

//dec6,
func CheckMinAmount(coinAmount sdk.Dec) bool {
	minAmount := sdk.NewDec(1)
	return coinAmount.LT(minAmount)
}


func EditorCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyright types.MsgEditorCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return err
	}
	log.Debug("do")

	account, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(copyright.Price), &price)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return errors.New(ParseCoinError)
	}
	flag := util.JudgeAmount(price.Amount)
	if !flag {
		log.Error("JudgeAmountl fail")
		return errors.New(ParseCoinError)
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(price)

	
	if CheckMinAmount(ledgeCoin.Amount.ToDec()) {
		log.Error("CheckMinAmount fail")
		return errors.New(LowerMinCoin)
	}

	if copyright.Datahash == "" {
		log.Error("copyright.Datahash is empty")
		return errors.New(DataHashEmpty)
	}

	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyright.Datahash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyright")
		return err
	}
	if copyrightInfor.DataHash == "" {
		log.Error("DataHash is empty")
		return errors.New(DataHashNotExist)
	}
	if copyrightInfor.Creator.String() != copyright.Creator {
		log.Error("copyrightInfor.Creator != copyright.Creator")
		return errors.New(DataHashNoRight)
	}

	
	if copyright.ChargeRate == "" {
		return errors.New(ChargeRateEmpty)
	}
	chargeRate, err := sdk.NewDecFromStr(copyright.ChargeRate)
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return errors.New(DecimalFromStringError)
	}
	if chargeRate.LT(core.ChargeRateLow) {
		log.Error("chargeRate < ChargeRateLow")
		return errors.New(ChargeRateTooSmall)
	}
	if chargeRate.GT(core.ChargeRateHigh) {
		log.Error("chargeRate > ChargeRateHigh")
		return errors.New(ChargeRateTooHigh)
	}
	return judgeFee(ctx, account, fee)
}


func DeleteCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	var copyright types.MsgDeleteCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	log.Debug("do")

	account, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		log.WithError(err).Error("AccAddressFromBech32")
		return errors.New(ParseAccountError)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyright.Datahash)
	if err != nil {
		log.WithError(err).Error("grpcQueryCopyright")
		return err
	}
	if copyrightInfor.DataHash == "" {
		log.Error("DataHash is empty")
		return errors.New(DataHashNotExist)
	}
	if copyrightInfor.Creator.String() != copyright.Creator {
		log.Error("copyrightInfor.Creator != copyright.Creator")
		return errors.New(DataHashNoRight)
	}
	if copyrightInfor.ApproveStatus == 0 {
		log.Error("ApproveStatus is 0")
		return errors.New(HasNotApprove)
	}

	return judgeFee(ctx, account, fee)
}
