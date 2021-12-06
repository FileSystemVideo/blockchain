package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)


func AuthorizeAccountHandlerFn(msgBytes []byte, ctx *client.Context, fee legacytx.StdFee, memo string) error {
	var authorize types.MsgAuthorizeAccount
	err := util.Json.Unmarshal(msgBytes, &authorize)
	if err != nil {
		return err
	}
	account, err := sdk.AccAddressFromBech32(authorize.Account)
	if err != nil {
		return errors.New(ParseAccountError)
	}

	flag, address := util.PricheckSign(authorize.Message, authorize.Sign)
	if !flag || address == "" {
		return errors.New(VerificationSignError)
	}
	_, err = sdk.ConsAddressFromBech32(authorize.ConsAddr)
	if err != nil {
		return errors.New(ConAddressError)
	}





	return judgeFee(ctx, account, fee)
}
