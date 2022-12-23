package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)


var txHandles *TxHandles


func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

	
	txHandles = newTxHandles(clientCtx)

	/********  ********/

	
	txHandles.Add(types.TypeMsgCreateCopyright, CreateCopyrightHandlerFn)

	
	//txHandles.Add(types.TypeMsgRegisterCopyrightParty, RegisterCopyrightPartyHandlerFn)

	
	txHandles.Add(types.TypeMsgSpaceMiner, SpaceMinerHandlerFn)

	
	txHandles.Add(types.TypeMsgSpaceMinerReward, SpaceMinerBonusHandlerFn)

	
	//txHandles.Add(types.TypeMsgDeflationVote, DeflationVoteHandlerFn)

	//nft
	txHandles.Add(types.TypeMsgNftTransfer, NftTransferHandlerFn)

	
	txHandles.Add(types.TypeMsgMortgage, MortgageHandlerFn)

	
	txHandles.Add(types.TypeMsgEditorCopyright, EditorCopyrightHandlerFn)

	
	txHandles.Add(types.TypeMsgDeleteCopyright, DeleteCopyrightHandlerFn)

	
	txHandles.Add(types.TypeMsgCopyrightComplain, CopyrightComplainHandlerFn)

	
	txHandles.Add(types.TypeMsgComplainResponse, ComplainResponseHandlerFn)

	
	txHandles.Add(types.TypeMsgComplainVote, ComplainVoteHandlerFn)

	
	txHandles.Add(proto.MessageName((*types.MsgTransfer)(nil)), TransferHandlerFn)

	
	txHandles.Add(types.TypeMsgCopyrightVote, CopyrightVoteHandlerFn)

	
	txHandles.Add(types.TypeMsgCrossChainIn, CrossChainInHandlerFn)

	txHandles.Add(types.TypeMsgCrossChainOut, CrossChainOutHandlerFn)
}


func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	//, grpc ,
	//r.HandleFunc("/copyright/find/{hash}", CopyrightFindHandlerFn(clientCtx)).Methods("GET")
	
	r.HandleFunc("/copyright/accountNumberSeq/{address}", AccountNumberSeqHandlerFn(clientCtx)).Methods("GET")
	
	r.HandleFunc("/copyright/transfer/rate", QueryTransferRateHandlerFn(clientCtx)).Methods("GET")
	// Get the current parameter values
	r.HandleFunc("/copyright/parameters", paramsHandlerFn(clientCtx)).Methods("GET")
}


func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	//tx，
	r.HandleFunc("/copyright/tx/broadcast", BroadcastTxHandlerFn(clientCtx)).Methods("POST")

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

func paramsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		if !ok {
			return
		}

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams), nil)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
