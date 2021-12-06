package rest

import (
	"context"
	"fs.video/blockchain/x/copyright/types"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)


func BlockMsgsHandlerDebugFn(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blockNum := vars["block"]
		res := types.TxPropsResponse{}
		res.Info = ""
		res.Status = 0
		txProps := []types.TxPropValue{}
		node, err := clientCtx.GetNode()
		if err != nil {
			res.Info = err.Error()
			SendReponse(w, clientCtx, res)
			return
		}


		height, err := strconv.ParseInt(blockNum, 10, 64)
		blockInfo, err := node.Block(context.Background(), &height)
		if err != nil {
			res.Info = err.Error()
			SendReponse(w, clientCtx, res)
			return
		}

		for i := 0; i < len(blockInfo.Block.Txs); i++ {
			txhashBytes := blockInfo.Block.Txs[i].Hash()
			txhashHex := strings.ToUpper(hex.EncodeToString(txhashBytes))
			resultTx, err := node.Tx(context.Background(), txhashBytes, true)
			if err != nil {
				continue
			}
			if resultTx.TxResult.Code != 0 {
				continue
			}
			tx, err := clientCtx.TxConfig.TxDecoder()(resultTx.Tx)
			if err != nil {
				continue
			}
			stdTx, err := txToStdTx(clientCtx, tx)
			if err != nil {
				continue
			}


			txFee := stdTx.Fee
			txMemo := stdTx.Memo
			blockTime := blockInfo.Block.Time.Unix()
			for j := 0; j < len(stdTx.Msgs); j++ {

				seq := int64(i + 1)
				txProps = append(txProps, types.TxPropValue{
					BlockTime: blockTime,
					TxHash:    txhashHex,
					Height:    height,
					Seq:       seq,
					Fee:       txFee,
					Memo:      txMemo,
					Events:    resultTx.TxResult.Events,
				})
			}
		}
		res.Props = txProps
		res.Status = 1
		rest.PostProcessResponseBare(w, clientCtx, res)
	}
}


func BlockMsgsHandlerFn(clientCtx client.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		blockNum := vars["block"]
		res := types.MsgsResponse{}
		res.Info = ""
		res.Status = 0
		messageResp := &types.MessageResp{}
		node, err := clientCtx.GetNode()
		if err != nil {
			res.Info = err.Error()
			SendReponse(w, clientCtx, res)
			return
		}


		height, err := strconv.ParseInt(blockNum, 10, 64)
		blockInfo, err := node.Block(context.Background(), &height)
		if err != nil {
			res.Info = err.Error()
			SendReponse(w, clientCtx, res)
			return
		}

		for i := 0; i < len(blockInfo.Block.Txs); i++ {
			txhashBytes := blockInfo.Block.Txs[i].Hash()
			txhashHex := strings.ToUpper(hex.EncodeToString(txhashBytes))
			resultTx, err := node.Tx(context.Background(), txhashBytes, true)
			if err != nil {
				continue
			}
			if resultTx.TxResult.Code != 0 {
				continue
			}
			tx, err := clientCtx.TxConfig.TxDecoder()(resultTx.Tx)
			if err != nil {
				continue
			}
			stdTx, err := txToStdTx(clientCtx, tx)
			if err != nil {
				continue
			}

			txFee := stdTx.Fee
			txMemo := stdTx.Memo
			blockTime := blockInfo.Block.Time.Unix()
			seq := int64(i + 1)
			txPropValue := &types.TxPropValue{
				BlockTime: blockTime,
				TxHash:    txhashHex,
				Height:    height,
				Seq:       seq,
				Fee:       txFee,
				Memo:      txMemo,
				Events:    resultTx.TxResult.Events,
			}
			for j := 0; j < len(stdTx.Msgs); j++ {
				err = parsingMsg(messageResp, stdTx.Msgs[j], txPropValue)
				if err != nil {
					continue
				}
			}

			err = parsingEvent(messageResp, txPropValue)
			if err != nil {
				continue
			}
		}
		res.Message = *messageResp
		res.Status = 1
		rest.PostProcessResponseBare(w, clientCtx, res)
	}
}
