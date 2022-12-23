package rest

import (
	"encoding/json"
	"errors"
	"fs.video/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gorilla/mux"
	"net/http"
)


var txHandles *TxHandles


func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)

	
	txHandles = newTxHandles(clientCtx)

	/********  ********/

	
}


func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	
	r.HandleFunc("/smart/msg", MsgHandlerFun(clientCtx)).Methods("POST")
}

func SendReponse(w http.ResponseWriter, clientCtx client.Context, body interface{}) {
	resBytes, err := json.Marshal(body)
	if err != nil {
		return
	}
	rest.PostProcessResponseBare(w, clientCtx, resBytes)
}

type TxHandles struct {
	ctx     client.Context
	funcMap map[string]func([]byte, *client.Context, legacytx.StdFee, string) error
}

func (this *TxHandles) Add(type1 string, func1 func([]byte, *client.Context, legacytx.StdFee, string) error) {
	this.funcMap[type1] = func1
}

func newTxHandles(ctx client.Context) *TxHandles {
	txHandles := &TxHandles{
		ctx:     ctx,
		funcMap: make(map[string]func([]byte, *client.Context, legacytx.StdFee, string) error),
	}
	return txHandles
}


func (this *TxHandles) HaveRegistered(msgType string) bool {
	_, ok := this.funcMap[msgType]
	return ok
}


func (this *TxHandles) Handle(msgType string, msgBytes []byte, fee legacytx.StdFee, memo string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainRest).WithField("msg", msgType)
	if !this.HaveRegistered(msgType) {
		log.Error("No handle registered!")
		return errors.New("msgType:" + msgType + " No handle registered!!")
	}
	log.Info("do") 
	return this.funcMap[msgType](msgBytes, &this.ctx, fee, memo)
}
