package rest

import (
	"fs.video/blockchain/x/copyright/types"
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/gorilla/mux"
	"net/http"

)

const (
	MethodGet = "GET"
)


var txHandles *TxHandles


func RegisterRoutes(clientCtx client.Context, r *mux.Router) {

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)


	txHandles = newTxHandles(clientCtx)




	txHandles.Add(types.TypeMsgCreateCopyright, CreateCopyrightHandlerFn)


	txHandles.Add(types.TypeMsgRegisterCopyrightParty, RegisterCopyrightPartyHandlerFn)


	txHandles.Add(types.TypeMsgSpaceMiner, SpaceMinerHandlerFn)


	txHandles.Add(types.TypeMsgDeflationVote, DeflationVoteHandlerFn)


	txHandles.Add(types.TypeMsgNftTransfer, NftTransferHandlerFn)


	txHandles.Add(types.TypeMsgMortgage, MortgageHandlerFn)


	txHandles.Add(types.TypeMsgEditorCopyright, EditorCopyrightHandlerFn)


	txHandles.Add(types.TypeMsgDeleteCopyright, DeleteCopyrightHandlerFn)


	txHandles.Add(types.TypeMsgCopyrightBonus, CopyrightBonusHandlerFn)


	txHandles.Add(types.TypeMsgCopyrightComplain, CopyrightComplainHandlerFn)


	txHandles.Add(types.TypeMsgComplainResponse, ComplainResponseHandlerFn)


	txHandles.Add(types.TypeMsgComplainVote, ComplainVoteHandlerFn)


	txHandles.Add(types.TypeMsgAuthorizeAccount, AuthorizeAccountHandlerFn)


	txHandles.Add(types.TypeMsgTransfer, TransferHandlerFn)


	txHandles.Add(types.TypeMsgCopyrightVote, CopyrightVoteHandlerFn)

	//
	txHandles.Add(types.TypeMsgCrossChainIn, CrossChainInHandlerFn)

	txHandles.Add(types.TypeMsgCrossChainOut, CrossChainOutHandlerFn)
}


func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {




	r.HandleFunc("/copyright/accountNumberSeq/{address}", AccountNumberSeqHandlerFn(clientCtx)).Methods("GET")


	r.HandleFunc("/copyright/list/{block}", BlockMsgsHandlerFn(clientCtx)).Methods("GET")


	r.HandleFunc("/copyright/list/{block}/debug", BlockMsgsHandlerDebugFn(clientCtx)).Methods("GET")

	CreateAnalysisHandle(&clientCtx)
}


func registerTxHandlers(clientCtx client.Context, r *mux.Router) {

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
	if !this.HaveRegistered(msgType) {
		return errors.New("msgType:" + msgType )
	}
	return this.funcMap[msgType](msgBytes, &this.ctx, fee, memo)
}
