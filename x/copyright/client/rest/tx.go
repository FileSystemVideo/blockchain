package rest

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"fs.video/trerr"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"io/ioutil"
	"net/http"
	"strings"
)

//对签名后的tx进行广播，并返回结果
func BroadcastTxHandlerFn(clientCtx client.Context) http.HandlerFunc {
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
		tx, _ := clientCtx.TxConfig.TxDecoder()(txBytes) //字节编码成基础tx
		stdTx, err := txToStdTx(clientCtx, tx)           //基础tx转成  stdTx
		if err != nil {
			txResponse.Info = err.Error()
			SendReponse(w, clientCtx, txResponse)
			return
		}
		msgs := stdTx.GetMsgs() //tx包含的消息列表
		fee := stdTx.Fee        //tx支付的手续费
		memo := stdTx.Memo      //tx的备注
		//业务校验
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
		fmt.Println("广播返回:", res, ",code:", res.Code)
		if res.Code == 0 { //code=0 表示没有错误
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

//解析code、codeSpace、rowlog 返回错误信息
func parseErrorCode(code uint32, codeSpace string, rowlog string) string {
	if codeSpace == sdkErrors.RootCodespace {
		if code == sdkErrors.ErrInsufficientFee.ABCICode() { //手续费不足
			return FeeIsTooLess
		} else if code == sdkErrors.ErrOutOfGas.ABCICode() { //消耗的gas超出了客户端设定的上限
			return ErrorGasOut
		} else if code == sdkErrors.ErrUnauthorized.ABCICode() { //链id或者account number 错误
			return ErrUnauthorized
		} else if code == sdkErrors.ErrWrongSequence.ABCICode() { //账号序列号错误
			return ErrWrongSequence
		}
	}
	return rowlog
}

func broadcastMsgCheck(msgs []sdk.Msg, fee legacytx.StdFee, memo string) (err error) {
	for _, msg := range msgs {
		msgType := msg.Type()
		if txHandles.HaveRegistered(msgType) { //检查该消息类型是否已经注册
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

//解析TX里面包含的msg列表
func ParseTxMsgs(clientCtx client.Context, txhashBytes []byte) ([]sdk.Msg, string, error) {
	txhashHex := strings.ToUpper(hex.EncodeToString(txhashBytes))
	output, err := authclient.QueryTx(clientCtx, txhashHex)
	if err != nil {
		return nil, "", err
	}
	txBytes := output.Tx.Value
	txI, err := clientCtx.TxConfig.TxDecoder()(txBytes)
	if err != nil {
		return nil, "", err
	}
	tx, ok := txI.(signing.Tx)
	if !ok {
		return nil, "", err
	}

	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, tx)
	if !ok {
		return nil, "", err
	}
	return stdTx.GetMsgs(), txhashHex, nil
}

/**
根据tx的字节码解析为 StdTx 结构体
*/
func bytesToStdTx(clientCtx client.Context, txhashBytes []byte) (*legacytx.StdTx, string, error) {
	txhashHex := strings.ToUpper(hex.EncodeToString(txhashBytes))
	output, err := authclient.QueryTx(clientCtx, txhashHex)
	if err != nil {
		return nil, "", err
	}
	txBytes := output.Tx.Value
	tx, err := clientCtx.TxConfig.TxDecoder()(txBytes)
	if err != nil {
		return nil, "", err
	}
	stdTx, err := txToStdTx(clientCtx, tx)
	if err != nil {
		return nil, "", err
	}
	return stdTx, txhashHex, nil
}

func txToStdTx(clientCtx client.Context, tx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := tx.(signing.Tx)
	if !ok {
		return nil, errors.New("tx转stdtx失败")
	}
	stdTx, err := clienttx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}
