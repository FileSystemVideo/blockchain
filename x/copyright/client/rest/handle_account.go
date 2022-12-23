package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)


func AccountNumberSeqHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		res := types.AccountNumberSeqResponse{
			AccountNumber: 0,
			Sequence:      0,
			NotFound:      false,
		}
		res.Status = 1
		res.Info = ""
		bech32addr := vars["address"]

		addr, err := sdk.AccAddressFromBech32(bech32addr)

		if err != nil {
			res.Info = err.Error()
			SendReponse(w, clientCtx, res)
			return
		}

		accountNumber, sequence, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, addr)
		if err != nil {
			
			if strings.Contains(err.Error(), "not found: key not found") {
				res.Status = 1
				res.Sequence = 0
				res.AccountNumber = 0
				res.NotFound = true
			} else {
				res.Status = 0
				res.Info = err.Error()
			}
			SendReponse(w, clientCtx, res)
			return
		}

		res.Status = 1
		res.AccountNumber = accountNumber
		res.Sequence = sequence
		SendReponse(w, clientCtx, res)
	}
}


func judgeBalance(cliCtx *client.Context, address sdk.AccAddress, amount sdk.Dec, denom string) (bool, string) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("addr", address.String())
	coin, err := grpcQueryBalance(cliCtx, address, denom)
	if err != nil {
		log.WithError(err).Error("grpcQueryBalance")
		return false, err.Error()
	}

	//log.Debug("", coin.String())

	
	if sdk.NewDecFromInt(coin.Amount).GTE(amount) {
		return true, ""
	}

	return false, AccountInsufficient
}

func judgeFee(ctx *client.Context, account sdk.AccAddress, fee legacytx.StdFee) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("addr", account.String())
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
		log.Warn("FeeCannotEmpty")
		return errors.New(FeeCannotEmpty)
	}
	if feeDec.IsZero() {
		log.Warn("FeeZero")
		return errors.New(FeeZero)
	}
	
	payAmount := types.NewLedgerDec(core.ChainDefaultFee)
	if feeDec.LT(payAmount) {
		log.Warn("FeeIsTooLess")
		return errors.New(FeeIsTooLess) 
	}

	balStatus, errStr := judgeBalance(ctx, account, payAmount, core.MainToken)
	if !balStatus {
		return errors.New(errStr)
	}
	return nil
}
