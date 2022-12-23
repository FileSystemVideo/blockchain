package rest

import (
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"fs.video/trerr"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/gogo/protobuf/proto"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//tx，
func BroadcastTxHandlerFn(clientCtx client.Context) http.HandlerFunc {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest)
	return func(w http.ResponseWriter, r *http.Request) {
		var txBytes []byte
		if r.Body != nil {
			txBytes, _ = ioutil.ReadAll(r.Body)
		}

		txResponse := types.BroadcastTxResponse{
			BaseResponse: types.BaseResponse{Info: "", Status: 0},
			TxHash:       "",
			Height:       0,
		}
		tx, _ := clientCtx.TxConfig.TxDecoder()(txBytes) //tx
		stdTx, err := txToStdTx(clientCtx, tx)           //tx  stdTx
		if err != nil {
			txResponse.Info = err.Error()
			SendReponse(w, clientCtx, txResponse)
			return
		}
		msgs := stdTx.GetMsgs() //tx
		fee := stdTx.Fee        //tx
		memo := stdTx.Memo      //tx
		
		err = broadcastMsgCheck(msgs, fee, memo)
		if err != nil {
			errmsg := trerr.TransError(err.Error())
			txResponse.Info = errmsg.Error()
			SendReponse(w, clientCtx, txResponse)
			return
		}
		res, err := clientCtx.BroadcastTx(txBytes)
		if err != nil {
			txResponse.Info = err.Error()
			SendReponse(w, clientCtx, txResponse)
			return
		}
		log.WithFields(logrus.Fields{
			"txhash": res.TxHash,
			"rawlog": res.RawLog,
			"code":   res.Code,
			"logs":   res.Logs,
		}).Info("result")

		if res.Code == 0 { //code=0 
			txResponse.Status = 1
		} else {
			txResponse.Status = 0
		}

		txResponse.Code = res.Code
		txResponse.Codespace = res.Codespace
		txResponse.TxHash = res.TxHash
		txResponse.Info = parseErrorCode(res.Code, res.Codespace, res.RawLog)
		txResponse.Height = res.Height
		SendReponse(w, clientCtx, txResponse)
	}
}

//code、codeSpace、rowlog 
func parseErrorCode(code uint32, codeSpace string, rowlog string) string {
	if codeSpace == sdkErrors.RootCodespace {
		if code == sdkErrors.ErrInsufficientFee.ABCICode() { 
			return FeeIsTooLess
		} else if code == sdkErrors.ErrOutOfGas.ABCICode() { //gas
			return ErrorGasOut
		} else if code == sdkErrors.ErrUnauthorized.ABCICode() { //idaccount number 
			return ErrUnauthorized
		} else if code == sdkErrors.ErrWrongSequence.ABCICode() { 
			return ErrWrongSequence
		}
	}
	return rowlog
}

func broadcastMsgCheck(msgs []sdk.Msg, fee legacytx.StdFee, memo string) (err error) {
	for _, msg := range msgs {
		msgType := proto.MessageName(msg)
		if txHandles.HaveRegistered(msgType) { 
			msgByte, err := util.Json.Marshal(msg)
			if err != nil {
				return err
			}
			err = txHandles.Handle(msgType, msgByte, fee, memo)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func txToStdTx(clientCtx client.Context, tx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := tx.(signing.Tx)
	if !ok {
		return nil, errors.New("txstdtx")
	}
	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}
