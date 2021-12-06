package rest

import (
	errors "errors"
	"fmt"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"strconv"
	"time"
)

func CreateCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyright types.MsgCreateCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(copyright.Price), &price)
	if err != nil {
		return errors.New(ParseCoinError)
	}
	flag := util.JudgeAmount(price.Amount)
	if !flag {
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
		return err
	}
	feeString := util.FilePrice(pubCount + 1)
	feeFloat, err := strconv.ParseFloat(feeString, 64)
	if err != nil {
		return errors.New(FormatStringToIntError)
	}
	pubFeeDec := types.NewLedgerDec(feeFloat)
	if !pubFeeDec.Equal(feeDec) {
		return errors.New(FeeInvalidError)
	}

	if CheckMinAmount(ledgeCoin.Amount.ToDec()) {
		return errors.New(LowerMinCoin)
	}

	balStatus, errStr := judgeBalance(ctx, addr, feeDec, config.MainToken)
	if !balStatus {
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
		return err
	}
	if exist {
		return errors.New(OriginDataHashExist)
	}
	exist, err = grpcQueryCopyrightOrigin(ctx, copyright.OriginDataHash)
	if err != nil {
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
		return errors.New(DecimalFromStringError)
	}
	if chargeRate.LT(config.ChargeRateLow) {
		return errors.New(ChargeRateTooSmall)
	}
	if chargeRate.GT(config.ChargeRateHigh) {
		return errors.New(ChargeRateTooHigh)
	}

	enough, err := grpcQueryAccountSpace(ctx, copyright.Size_, copyright.Creator)
	if err != nil {
		return err
	}
	if !enough {
		return errors.New(SpaceNotEnough)
	}
	return nil
}

func CheckMinAmount(coinAmount sdk.Dec) bool {
	minAmount := sdk.NewDec(1)
	return coinAmount.LT(minAmount)
}

func EditorCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyright types.MsgEditorCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		return err
	}

	account, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	var price types.RealCoin
	err = util.Json.Unmarshal([]byte(copyright.Price), &price)
	if err != nil {
		return errors.New(ParseCoinError)
	}
	flag := util.JudgeAmount(price.Amount)
	if !flag {
		return errors.New(ParseCoinError)
	}
	ledgeCoin := types.MustRealCoin2LedgerCoin(price)

	if CheckMinAmount(ledgeCoin.Amount.ToDec()) {
		return errors.New(LowerMinCoin)
	}

	if copyright.Datahash == "" {
		return errors.New(DataHashEmpty)
	}

	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyright.Datahash)
	if err != nil {
		return err
	}
	if copyrightInfor.DataHash == "" {
		return errors.New(DataHashNotExist)
	}
	if copyrightInfor.Creator.String() != copyright.Creator {
		return errors.New(DataHashNoRight)
	}

	if copyright.ChargeRate == "" {
		return errors.New(ChargeRateEmpty)
	}
	chargeRate, err := sdk.NewDecFromStr(copyright.ChargeRate)
	if err != nil {
		return errors.New(DecimalFromStringError)
	}
	if chargeRate.LT(config.ChargeRateLow) {
		return errors.New(ChargeRateTooSmall)
	}
	if chargeRate.GT(config.ChargeRateHigh) {
		return errors.New(ChargeRateTooHigh)
	}
	return judgeFee(ctx, account, fee)
}

func DeleteCopyrightHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyright types.MsgDeleteCopyright
	err := util.Json.Unmarshal(msgBytes, &copyright)
	if err != nil {
		return err
	}

	account, err := sdk.AccAddressFromBech32(copyright.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	copyrightInfor, _, err := grpcQueryCopyright(ctx, copyright.Datahash)
	if err != nil {
		return err
	}
	if copyrightInfor.DataHash == "" {
		return errors.New(DataHashNotExist)
	}
	if copyrightInfor.Creator.String() != copyright.Creator {
		return errors.New(DataHashNoRight)
	}
	if copyrightInfor.ApproveStatus == 0 {
		return errors.New(HasNotApprove)
	}

	return judgeFee(ctx, account, fee)
}

func CopyrightBonusHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var copyrightBonus types.MsgCopyrightBonus
	err := util.Json.Unmarshal(msgBytes, &copyrightBonus)
	if err != nil {
		return err
	}

	downer, err := sdk.AccAddressFromBech32(copyrightBonus.Creator)
	if err != nil {
		return errors.New(ParseAccountError)
	}
	copyrightInfor, height, err := grpcQueryCopyright(ctx, copyrightBonus.Datahash)
	if err != nil {
		return err
	}
	if copyrightInfor.DataHash == "" {
		return errors.New(DataHashNotExist)
	}

	data, err := grpcQueryCopyrightAndAccount(ctx, copyrightBonus.Creator, copyrightBonus.Datahash)
	if err != nil {
		return errors.New(QueryChainInforError)
	}

	if data.Species == "buy" || (!data.Downer.Empty() && height < data.Height) {
		return errors.New(HasPayedForDatahash)
	}

	err = judgeOfferAccount(copyrightBonus.OfferAccountShare)
	if err != nil {
		return err
	}

	ledgeCoin := types.MustRealCoin2LedgerCoin(copyrightInfor.Price)
	balStatus, errStr := judgeBalance(ctx, downer, ledgeCoin.Amount.ToDec(), config.MainToken)
	if !balStatus {
		return errors.New(errStr)
	}
	return nil
}
