package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channelutils "github.com/cosmos/ibc-go/v3/modules/core/04-channel/client/utils"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type AccountClient struct {
	TxClient  *TxClient
	key       *SecretKey
	ServerUrl string
	logPrefix string
}

type Account struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Address  string `json:"address"`
	Pubkey   string `json:"pubkey"`
	Mnemonic string `'json:"mnemonic"`
}

func (this *Account) Print() {
	fmt.Printf(":\t %s \n", this.Name)
	fmt.Printf(":\t %s \n", this.Address)
	fmt.Printf(":\t\t %s \n", this.Type)
	fmt.Printf(":\t\t %s \n", this.Pubkey)
	fmt.Printf(":\t\t %s \n", this.Mnemonic)
}

type AccountList struct {
	Accounts []Account
}

//SequenceAccountNumber
func (this *AccountClient) FindAccountNumberSeq(accountAddr string) (detail types.AccountNumberSeqResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithField("acc", accountAddr)
	reponse, err := GetRequest(this.ServerUrl, "/copyright/accountNumberSeq/"+accountAddr)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(reponse), &detail)
	if err != nil {
		log.Error("json.Unmarshal")
		return
	}
	return
}


func (this *AccountClient) GetAllAccounts() (accounts []string, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	reponseStr, _, err := clientCtx.QueryWithData("custom/auth/accounts", []byte{})
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &accounts)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON2")
		return
	}
	return
}

//token
func (this *AccountClient) FindAccountBalances(accountAddr string, height string) (coins types.RealCoins, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "height": height})
	url := "/bank/balances/" + accountAddr
	if height != "" {
		url += "?height=" + height
	}
	reponseStr, err := GetRequest(this.ServerUrl, url)
	if err != nil {
		log.Error("GetRequest")
		return
	}
	//fmt.Println(reponseStr)
	var resp = rest.ResponseWithHeight{}
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponseStr), &resp)
	if err != nil {
		log.Error("UnmarshalJSON1")
		return
	}
	var ledgerCoins sdk.Coins
	err = clientCtx.LegacyAmino.UnmarshalJSON(resp.Result, &ledgerCoins)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	coins = types.MustLedgerCoins2RealCoins(ledgerCoins)
	return
}


func (this *AccountClient) Transfer(data types.TransferData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	msg, err := types.NewMsgTransfer(data)
	if err != nil {
		log.Error("NewMsgTransfer")
		return
	}
	
	resp, err = this.TxClient.SignAndSendMsg(data.FromAddress, privateKey, data.Fee, data.Memo, msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}

//fsv
func (this *AccountClient) CrossChainOut(data types.CrossChainOutData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	
	msg, err := types.NewMsgCrossChainOut(data)
	if err != nil {
		log.Error("NewMsgCrossChainOut")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(data.SendAddress, privateKey, data.Fee, data.Memo, msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}

//fsv
func (this *AccountClient) CrossChainIn(data types.CrossChainInData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	
	msg, err := types.NewMsgCrossChainIn(data)
	if err != nil {
		log.Error("NewMsgCrossChainIn")
		return
	}

	
	resp, err = this.TxClient.SignAndSendMsg(data.SendAddress, privateKey, data.Fee, data.Memo, msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}

//IBC
func (this *AccountClient) IbcTransfer(data types.IbcTransferData, privateKey string) (resp *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)

	token := types.MustRealCoin2LedgerCoin(data.Token)
	if !strings.HasPrefix(token.Denom, "ibc/") {
		denomTrace := ibcTransferTypes.ParseDenomTrace(token.Denom)
		token.Denom = denomTrace.IBCDenom()
	}
	timeoutHeight, err := clienttypes.ParseHeight(ibcTransferTypes.DefaultRelativePacketTimeoutHeight)
	timeoutTimestamp := ibcTransferTypes.DefaultRelativePacketTimeoutTimestamp
	consensusState, height, _, err := channelutils.QueryLatestConsensusState(clientCtx, data.SourcePort, data.SourceChannel)
	if err != nil {
		log.WithError(err).Error("QueryLatestConsensusState")
		return
	}
	if !timeoutHeight.IsZero() {
		absoluteHeight := height
		absoluteHeight.RevisionNumber += timeoutHeight.RevisionNumber
		absoluteHeight.RevisionHeight += timeoutHeight.RevisionHeight
		timeoutHeight = absoluteHeight
	}

	if timeoutTimestamp != 0 {
		// use local clock time as reference time if it is later than the
		// consensus state timestamp of the counter party chain, otherwise
		// still use consensus state timestamp as reference
		now := time.Now().UnixNano()
		consensusStateTimestamp := consensusState.GetTimestamp()
		if now > 0 {
			now := uint64(now)
			if now > consensusStateTimestamp {
				timeoutTimestamp = now + timeoutTimestamp
			} else {
				timeoutTimestamp = consensusStateTimestamp + timeoutTimestamp
			}
		} else {
			log.Error("local clock time is not greater than Jan 1st, 1970 12:00 AM")
			return
		}
	}
	//IBC
	msg := ibcTransferTypes.NewMsgTransfer(data.SourcePort, data.SourceChannel, token, data.SendAddress, data.ToAddress, timeoutHeight, timeoutTimestamp)

	
	resp, err = this.TxClient.SignAndSendMsg(data.SendAddress, privateKey, data.Fee, data.Memo, msg)
	if err != nil {
		return
	}
	
	if resp.Status == 1 {
		return resp, nil
	} else {
		
		return resp, errors.New(resp.Info)
	}
}

func (this *AccountClient) FindBalanceByRpc(accountAddr string, denom string) (realCoins types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom})

	req := banktypes.QueryBalanceRequest{Address: accountAddr, Denom: denom}

	reqBytes, _ := clientCtx.LegacyAmino.MarshalJSON(req)

	reponseStr, _, err := clientCtx.QueryWithData("custom/bank/balance", reqBytes)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	var coin sdk.Coin
	err = clientCtx.LegacyAmino.UnmarshalJSON(reponseStr, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	realCoins = types.MustLedgerCoin2RealCoin(coin)
	return
}

//token
func (this *AccountClient) FindAccountBalance(accountAddr string, denom, height string) (realCoins types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom, "height": height})
	url := "/bank/balances/" + accountAddr + "?denom=" + denom
	if height != "" {
		url += "&height=" + height
	}
	reponseStr, err := GetRequest(this.ServerUrl, url)
	if err != nil {
		log.Error("GetRequest")
		return
	}

	var resp = rest.ResponseWithHeight{}
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponseStr), &resp)
	if err != nil {
		log.Error("UnmarshalJSON1")
		return
	}
	var coin sdk.Coin
	err = clientCtx.LegacyAmino.UnmarshalJSON(resp.Result, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	realCoins = types.MustLedgerCoin2RealCoin(coin)
	return
}

//token
func (this *AccountClient) FindAccountBalanceChain(accountAddr, denom string, height int64) (realCoins types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient).WithFields(logrus.Fields{"acc": accountAddr, "denom": denom, "height": height})
	/*url := "/bank/balances/" + accountAddr + "?denom=" + denom
	if height != "" {
		url += "&height=" + height
	}
	reponseStr, err := GetRequest(this.ServerUrl, url)
	if err != nil {
		log.Error("GetRequest")
		return
	}
	var resp = rest.ResponseWithHeight{}
	err = clientCtx.LegacyAmino.UnmarshalJSON([]byte(reponseStr), &resp)
	if err != nil {
		log.Error("UnmarshalJSON1")
		return
	}*/
	var coin sdk.Coin
	/*err = clientCtx.LegacyAmino.UnmarshalJSON(resp.Result, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}*/
	var addr sdk.AccAddress
	if accountAddr == sdk.AccAddress(sdk.BlackHoleAddress).String() {
		addr = sdk.AccAddress(sdk.BlackHoleAddress)
	} else {
		addr, err = sdk.AccAddressFromBech32(accountAddr)
	}
	if err != nil {
		return realCoins, err
	}

	params := banktypes.NewQueryBalanceRequest(addr, denom)
	route := fmt.Sprintf("custom/%s/%s", banktypes.QuerierRoute, banktypes.QueryBalance)
	clientCtx = clientCtx.WithHeight(height)
	/*resBytes, _, err := clientCtx.QueryWithData("custom/copyright/"+types.QueryHasMinerBonus, bz)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return "", err
	}*/
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		return
	}

	res, height, err := clientCtx.QueryWithData(route, bz)
	if err != nil {
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(res, &coin)
	if err != nil {
		log.Error("UnmarshalJSON2")
		return
	}
	realCoins = types.MustLedgerCoin2RealCoin(coin)
	return
}

//  hex.EncodeToString  
func (this *AccountClient) CreateAccountFromPrivEth(priv string) (*EthWallet, error) {
	return this.key.CreateAccountFromPrivEth(priv)
}

//  hex.EncodeToString   cosmos
func (this *AccountClient) CreateAccountFromPrivCos(priv string) (*CosmosWallet, error) {
	return this.key.CreateAccountFromPrivCos(priv)
}


func (this *AccountClient) CreateAccountFromSeedEth(seed string) (acc *EthWallet, err error) {
	return this.key.CreateAccountFromSeedEth(seed)
}

//() cosmos
func (this *AccountClient) CreateAccountFromSeedCos(seed string) (acc *CosmosWallet, err error) {
	return this.key.CreateAccountFromSeedCos(seed)
}


func (this *AccountClient) CreateSeedWord() (mnemonic string, err error) {
	return this.key.CreateSeedWord()
}
