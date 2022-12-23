package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cmssecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	"github.com/shopspring/decimal"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	ttypes "github.com/tendermint/tendermint/types"
	evmhd "github.com/tharsis/ethermint/crypto/hd"
	"regexp"
	"strconv"
)

type TxInfo struct {
	Height             string   //id
	Txhash             string   
	Data               string   
	LastCommitHash     string   //hash
	Datahash           string   //hash
	ValidatorsHash     string   //hash
	NextValidatorsHash string   //hash
	ConsensusHash      string   //Hash
	Apphash            string   //hash
	LastResultsHash    string   //Hash
	EvidenceHash       string   //hash
	ProposerAddress    string   
	Txs                []string //hash
	Signatures         []Signature
}

type TxClient struct {
	ServerUrl string
	logPrefix string
}

func (this *TxClient) ConvertTxToStdTx(cosmosTx sdk.Tx) (*legacytx.StdTx, error) {
	signingTx, ok := cosmosTx.(xauthsigning.Tx)
	if !ok {
		return nil, errors.New("txstdtx")
	}
	stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, signingTx)
	if err != nil {
		return nil, err
	}
	return &stdTx, nil
}

//tx tendermint tx  cosmos  tx
func (this *TxClient) TermintTx2CosmosTx(signTxs ttypes.Tx) (sdk.Tx, error) {
	return clientCtx.TxConfig.TxDecoder()(signTxs)
}

//tx
func (this *TxClient) SignTx2Bytes(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

func (this *TxClient) SetFee(signTxs xauthsigning.Tx) ([]byte, error) {
	return clientCtx.TxConfig.TxEncoder()(signTxs)
}

/**
tx
*/
func (this *TxClient) SignTx(privateKey string, seqDetail types.AccountNumberSeqResponse, fee legacytx.StdFee, memo string, msgs ...sdk.Msg) (xauthsigning.Tx, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.WithError(err).Error("hex.DecodeString")
		return nil, err
	}
	keyringAlgos := keyring.SigningAlgoList{evmhd.EthSecp256k1}
	algo, err := keyring.NewSigningAlgoFromString("eth_secp256k1", keyringAlgos)
	if err != nil {
		return nil, err
	}
	privKey := algo.Generate()(privKeyBytes)
	//gas,gas
	if fee.Gas == flags.DefaultGasLimit {
		feeCoin, gas, err := this.GasInfo(seqDetail, msgs...)
		if err != nil {
			log.WithError(err).Error("CulGas")
			return nil, core.Errformat(err)
		}
		log.WithField("gas", gas).Info("CulGas:")
		
		ledgeCoin := types.MustRealCoin2LedgerCoin(feeCoin)
		if ledgeCoin.Amount.GT(fee.Amount.AmountOf(sdk.DefaultBondDenom)) {
			fee.Amount = sdk.NewCoins(ledgeCoin)
		}
		fee.Gas = gas
	}
	signMode := clientCtx.TxConfig.SignModeHandler().DefaultMode()
	signerData := xauthsigning.SignerData{
		ChainID:       clientCtx.ChainID,
		AccountNumber: seqDetail.AccountNumber,
		Sequence:      seqDetail.Sequence,
	}
	txBuild, err := tx.BuildUnsignedTx(clientFactory, msgs...)
	if err != nil {
		log.WithError(err).Error("tx.BuildUnsignedTx")
		return nil, err
	}
	txBuild.SetGasLimit(fee.Gas)     //gas
	txBuild.SetFeeAmount(fee.Amount) 
	txBuild.SetMemo(memo)            
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privKey.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	
	if err := txBuild.SetSignatures(sig); err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}
	signV2, err := tx.SignWithPrivKey(signMode, signerData, txBuild, privKey, clientCtx.TxConfig, seqDetail.Sequence)
	if err != nil {
		log.WithError(err).Error("SignWithPrivKey")
		return nil, err
	}
	err = txBuild.SetSignatures(signV2)
	if err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}

	signedTx := txBuild.GetTx()
	//fmt.Println("getSigners:",signedTx.GetSigners())
	return signedTx, nil
}

/**
tx
*/
func (this *TxClient) SignTxCos(privateKey string, seqDetail types.AccountNumberSeqResponse, fee legacytx.StdFee, memo string, msgs ...sdk.Msg) (xauthsigning.Tx, error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	var privkey secp256k1.PrivKey
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		log.WithError(err).Error("hex.DecodeString")
		return nil, err
	}
	privkey = privKeyBytes
	privkey1 := cmssecp256k1.PrivKey{Key: privkey}
	//gas,gas
	if fee.Gas == flags.DefaultGasLimit {
		feeCoin, gas, err := this.GasInfo(seqDetail, msgs...)
		if err != nil {
			log.WithError(err).Error("CulGas")
			return nil, core.Errformat(err)
		}
		log.WithField("gas", gas).Info("CulGas:")
		
		ledgeCoin := types.MustRealCoin2LedgerCoin(feeCoin)

		if ledgeCoin.Amount.GT(fee.Amount.AmountOf(sdk.DefaultBondDenom)) && proto.MessageName(msgs[0]) != "copyright.v1beta1.MsgTransfer" {
			fee.Amount = sdk.NewCoins(ledgeCoin)
		}
		fee.Gas = gas
	}
	signMode := clientCtx.TxConfig.SignModeHandler().DefaultMode()
	signerData := xauthsigning.SignerData{
		ChainID:       clientCtx.ChainID,
		AccountNumber: seqDetail.AccountNumber,
		Sequence:      seqDetail.Sequence,
	}
	txBuild, err := tx.BuildUnsignedTx(clientFactory, msgs...)
	if err != nil {
		log.WithError(err).Error("tx.BuildUnsignedTx")
		return nil, err
	}
	txBuild.SetGasLimit(fee.Gas)     //gas
	txBuild.SetFeeAmount(fee.Amount) 
	txBuild.SetMemo(memo)            
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privkey1.PubKey(),
		Data:     &sigData,
		Sequence: seqDetail.Sequence,
	}
	
	if err := txBuild.SetSignatures(sig); err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}
	signV2, err := tx.SignWithPrivKey(signMode, signerData, txBuild, &privkey1, clientCtx.TxConfig, seqDetail.Sequence)
	if err != nil {
		log.WithError(err).Error("SignWithPrivKey")
		return nil, err
	}
	err = txBuild.SetSignatures(signV2)
	if err != nil {
		log.WithError(err).Error("SetSignatures")
		return nil, err
	}

	signedTx := txBuild.GetTx()
	return signedTx, nil
}


func (this *TxClient) Send(req []byte) (txRes *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	response, err := PostRequest(this.ServerUrl, "/copyright/tx/broadcast", req)
	if err != nil {
		log.WithError(err).Error("PostRequest")
		return
	}
	txRes = &types.BroadcastTxResponse{}
	err = json.Unmarshal([]byte(response), txRes)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
		return
	}
	return
}

//tx
func (this *TxClient) FindByByte(txhash []byte) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	node, err := clientCtx.GetNode()
	if err != nil {
		log.WithError(err).Error("GetNode")
		return
	}
	resultTx, err = node.Tx(context.Background(), txhash, true)
	if err != nil {
		
		notFound = this.isTxNotFoundError(err.Error())
		if notFound {
			err = nil
		} else {
			log.WithError(err).WithField("txhash", hex.EncodeToString(txhash)).Error("node.Tx")
		}
		return
	}
	return
}

//tx,,err=nil notFound=true
//notFound 
//err    
func (this *TxClient) FindByHex(txhashStr string) (resultTx *ctypes.ResultTx, notFound bool, err error) {
	var txhash []byte
	notFound = false
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	txhash, err = hex.DecodeString(txhashStr)
	if err != nil {
		log.WithError(err).WithField("txhash", txhashStr).Error("hex.DecodeString")
		return
	}
	return this.FindByByte(txhash)
}

//tx
func (this *TxClient) isTxNotFoundError(errContent string) (ok bool) {
	errRegexp := `tx\ \([0-9A-Za-z]{64}\)\ not\ found`
	r, err := regexp.Compile(errRegexp)
	if err != nil {
		return false
	}
	if r.Match([]byte(errContent)) {
		return true
	} else {
		return false
	}
}

/**

*/
func (this *TxClient) SignAndSendMsg(address string, privateKey string, fee legacytx.StdFee, memo string, msg ...sdk.Msg) (txRes *types.BroadcastTxResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	
	seqDetail, err := this.FindAccountNumberSeq(address)
	if err != nil {
		return
	}
	var signedTx xauthsigning.Tx
	if address[:3] == "dex" {
		
		signedTx, err = this.SignTx(privateKey, seqDetail, fee, memo, msg...)
		if err != nil {
			return
		}
	} else {
		
		signedTx, err = this.SignTxCos(privateKey, seqDetail, fee, memo, msg...)
		if err != nil {
			return
		}
	}
	
	signPubkeyBytes, err := signedTx.GetPubKeys()
	if err != nil {
		return
	}
	senderPubkeyBytes := signedTx.GetSigners()[0].Bytes()
	if !bytes.Equal(signPubkeyBytes[0].Address().Bytes(), senderPubkeyBytes) {
		return nil, errors.New("sign error")
	}

	
	signedTxBytes, err := this.SignTx2Bytes(signedTx)
	if err != nil {
		log.WithError(err).Error("SignTx2Bytes")
		return
	}

	
	txRes, err = this.Send(signedTxBytes)
	if txRes != nil {
		txRes.SignedTxStr = hex.EncodeToString(signedTxBytes)
	}
	return
}

//SequenceAccountNumber
func (this *TxClient) FindAccountNumberSeq(accountAddr string) (detail types.AccountNumberSeqResponse, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	reponse, err := GetRequest(this.ServerUrl, "/copyright/accountNumberSeq/"+accountAddr)
	if err != nil {
		log.WithError(err).Error("GetRequest")
		return
	}
	err = json.Unmarshal([]byte(reponse), &detail)
	if err != nil {
		log.WithError(err).Error("json.Unmarshal")
	}
	return
}

/**

*/
func (this *TxClient) Balance(accountAddr, demo string) (realCoin *types.RealCoin, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	params := bankTypes.QueryBalanceRequest{Address: accountAddr, Denom: demo}
	bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
	if err != nil {
		log.WithError(err).Error("MarshalJSON")
		return nil, err
	}
	resBytes, _, err := clientCtx.QueryWithData("custom/bank/"+bankTypes.QueryBalance, bz)
	if err != nil {
		log.WithError(err).WithField("acc", accountAddr).Error("QueryWithData")
		return nil, err
	}
	var coin sdk.Coin
	err = util.Json.Unmarshal(resBytes, &coin)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return nil, err
	}
	realCoin1 := types.MustLedgerCoin2RealCoin(coin)
	return &realCoin1, nil
}


func (this *TxClient) ListenBlock() error {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	err := clientCtx.Client.Start()
	if err != nil {
		log.WithError(err).Error("Start")
		return err
	}
	eventCh, err := clientCtx.Client.Subscribe(context.Background(), "TestBlockEvents", ttypes.QueryForEvent(ttypes.EventNewBlock).String())
	if err != nil {
		log.WithError(err).Error("Subscribe")
		return err
	}
	for {
		event := <-eventCh
		blockEvent, ok := event.Data.(ttypes.EventDataNewBlock)
		if !ok {
			continue
		}
		block := blockEvent.Block
		fmt.Println(block)
	}
	return nil
}


func (this *TxClient) CulGas(seqDetail types.AccountNumberSeqResponse, msg ...sdk.Msg) (gas uint64, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	clientFactory = clientFactory.WithSequence(seqDetail.Sequence)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg...)
	if err != nil {
		log.WithError(err).Error("tx.CalculateGas")
		return
	}
	gas = gasInfo.GasInfo.GasUsed * 2
	return
}


func (this *TxClient) GasInfo(seqDetail types.AccountNumberSeqResponse, msg ...sdk.Msg) (coin types.RealCoin, gas uint64, err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainClient)
	clientFactory = clientFactory.WithSequence(seqDetail.Sequence)
	gasInfo, _, err := tx.CalculateGas(clientCtx, clientFactory, msg...)
	if err != nil {
		log.WithError(err).Error("tx.CalculateGas")
		return
	}
	gas = gasInfo.GasInfo.GasUsed * 2
	gasUsed := decimal.RequireFromString(strconv.Itoa(int(gas)))
	//gasUsed = gasUsed.Add(decimal.RequireFromString("2000000"))
	gasUsed = gasUsed.Mul(core.MinimumGasPrices).Add(decimal.NewFromFloat(1))
	gasDec, err := sdk.NewDecFromStr(gasUsed.StringFixed(6))
	if err != nil {
		log.WithError(err).Error("NewDecFromStr")
		return
	}
	gasDecCoin := sdk.NewDecCoinFromDec(core.MainToken, gasDec)
	amount := types.MustParseLedgerDecCoin(gasDecCoin)
	if decimal.RequireFromString(amount).LessThan(core.ChainDefaultFee) {
		
		amount = core.ChainDefaultFee.String()
	}
	coin = types.NewRealCoinFromStr(core.MainToken, amount)
	return
}

type Params struct {
	StakingParam      stakingTypes.Params      `json:"staking_param"`
	DistributionParam distributionTypes.Params `json:"distribution_param"`
	SlashingParam     slashingTypes.Params     `json:"slashing_param"`
}

func (this *TxClient) GetParams() (param Params, err error) {
	log := core.BuildLog(core.GetStructFuncName(this), core.LmChainClient)
	stakingParam := stakingTypes.Params{}
	resBytes, _, err := clientCtx.QueryWithData("custom/staking/"+stakingTypes.QueryParameters, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &stakingParam)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return
	}
	param.StakingParam = stakingParam
	distributionParam := distributionTypes.Params{}
	resBytes, _, err = clientCtx.QueryWithData("custom/distribution/"+distributionTypes.QueryParams, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &distributionParam)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return
	}
	param.DistributionParam = distributionParam
	slashingParam := slashingTypes.Params{}
	resBytes, _, err = clientCtx.QueryWithData("custom/slashing/"+slashingTypes.QueryParameters, nil)
	if err != nil {
		log.WithError(err).Error("QueryWithData")
		return
	}
	err = clientCtx.LegacyAmino.UnmarshalJSON(resBytes, &slashingParam)
	if err != nil {
		log.WithError(err).Error("UnmarshalJSON")
		return
	}
	param.SlashingParam = slashingParam
	return
}
