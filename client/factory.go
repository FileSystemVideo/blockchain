package client

import (
	"bytes"
	"fs.video/blockchain/app"
	"fs.video/blockchain/core"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authType "github.com/cosmos/cosmos-sdk/x/auth/types"
	flag "github.com/spf13/pflag"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tharsis/ethermint/encoding"
)

var clientCtx client.Context

var clientFactory tx.Factory

var encodingConfig params.EncodingConfig

func SetClientChainID(chainID string) {
	clientCtx = clientCtx.WithChainID(chainID)
}

func NewEvmClient() EvmClient {
	return EvmClient{core.EvmRpcURL}
}

func NewBlockClient() BlockClient {
	return BlockClient{core.ServerURL, "BlockClient"}
}

func NewTxClient() TxClient {
	return TxClient{core.ServerURL, "TxClient"}
}

func NewBlockClientUrl(url string) BlockClient {
	return BlockClient{url, "BlockClient"}
}

func NewTxClientServerURl(url string) TxClient {
	return TxClient{url, "TxClient"}
}

func NewNodeClient(dposClient *DposClient) NodeClient {
	return NodeClient{dposClient, core.ServerURL}
}

func NewAccountClient(txClient *TxClient) AccountClient {
	return AccountClient{txClient, NewSecretKey(), core.ServerURL, "AccountClient"}
}

func NewCopyrightClient(txClient *TxClient, accClient *AccountClient) *CopyrightClient {
	return &CopyrightClient{txClient, accClient, core.ServerURL, "CopyrightClient"}
}

func NewCopyrightPartyClient(txClient *TxClient) *CopyrightPartyClient {
	return &CopyrightPartyClient{txClient, core.ServerURL, "CopyrightPartyClient"}
}

func NewDposClient(txClient *TxClient) *DposClient {
	return &DposClient{txClient, core.ServerURL, "DposClient"}
}

func NewSpaceMinerClient(txClient *TxClient) SpaceMinerClient {
	return SpaceMinerClient{txClient, core.ServerURL, "SpaceMinerClient"}
}

func NewNftClient(txClient *TxClient) NftClient {
	return NftClient{txClient, core.ServerURL}
}

func NewComplainClient(txClient *TxClient) ComplainClient {
	return ComplainClient{txClient, core.ServerURL, "ComplainClient"}
}
func NewMortgageClient(txClient *TxClient) MortgageClient {
	return MortgageClient{txClient, core.ServerURL}
}
func NewAuthorizeClient(txClient *TxClient) AuthorizeClient {
	return AuthorizeClient{txClient, core.ServerURL}
}

func NewInviteCodeClient(txClient *TxClient) ShareClient {
	return ShareClient{txClient, core.ServerURL, "InviteClient"}
}


func MsgToStruct(msg sdk.Msg, obj interface{}) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	msgByte, err := encodingConfig.Amino.Marshal(msg)
	if err != nil {
		log.WithError(err).Error("MarshalBinaryBare")
		return err
	}
	err = encodingConfig.Amino.Unmarshal(msgByte, obj)
	if err != nil {
		log.WithError(err).Error("UnmarshalBinaryBare")
		return err
	}
	return nil
}

func init() {
	encodingConfig = encoding.MakeConfig(app.ModuleBasics)

	rpcClient, err := rpchttp.New(core.RpcURL, "/websocket")
	if err != nil {
		panic("start ctx client error.")
	}

	clientCtx = client.Context{}.
		WithChainID(core.ChainID).
		WithCodec(encodingConfig.Marshaler).
		//WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithOffline(true).
		WithNodeURI(core.RpcURL).
		WithClient(rpcClient).
		WithAccountRetriever(authType.AccountRetriever{})

	//cfg := network.DefaultConfig()

	//clientCtx = clientCtx.WithLegacyAmino(cfg.LegacyAmino)

	flags := flag.NewFlagSet("chat", flag.ContinueOnError)

	flagErrorBuf := new(bytes.Buffer)

	flags.SetOutput(flagErrorBuf)

	//gas  gas 
	clientFactory = tx.NewFactoryCLI(clientCtx, flags)
	clientFactory.WithChainID(core.ChainID).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
}
