package client

import (
	"context"
	"encoding/base64"
	"fs.video/blockchain/core"
	"fs.video/blockchain/x/copyright/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"strconv"
)

type NodeClient struct {
	dposClient *DposClient
	ServerUrl  string
}

//  ValidatorStatus:  0 Unbonded , 1 Unbonding , 2 Bonded , 3  , 4 
func (this *NodeClient) ValidatorInfo() (validatorInfo *types.ValidatorInfor, err error) {
	nodeStatus, err := this.StatusInfo()
	if err != nil {
		return nil, err
	}

	
	consAddress, err := sdk.ConsAddressFromHex(nodeStatus.ValidatorInfo.Address.String())
	if err != nil {
		return nil, err
	}
	validatorInfo = &types.ValidatorInfor{
		Status:            "1",
		ValidatorStatus:   "", //0  Unbonded 1 Unbonding 2 Bonded 3  4 
		ValidatorPubAddr:  nodeStatus.ValidatorInfo.PubKey.Address().String(),
		ValidatorConsAddr: consAddress.String(),
	}

	
	validator, notFound, err := this.dposClient.FindValidatorByConsAddress(consAddress.String())
	if notFound {
		validatorInfo.ValidatorStatus = "4" 
		return validatorInfo, nil
	}

	if err != nil {
		return nil, err
	}

	validatorInfo.ValidatorStatus = strconv.Itoa(this.dposClient.GetValidatorStatus(validator.Status, validator.Jailed))
	validatorInfo.ValidatorOperAddr = validator.GetOperator().String()
	accAddre := sdk.AccAddress(validator.GetOperator())
	validatorInfo.AccAddr = accAddre.String()
	return validatorInfo, nil
}


func (this *NodeClient) StatusInfo() (statusInfo *ctypes.ResultStatus, err error) {
	node, err := clientCtx.GetNode()
	return node.Status(context.Background())
}


func (this *NodeClient) NetInfo() (statusInfo *ctypes.ResultNetInfo, err error) {
	node, err := clientCtx.GetNode()
	return node.NetInfo(context.Background())
}


func (this *NodeClient) GetTranserTypeConfig() map[string]string {
	return core.GetTranserTypeConfig()
}

//base64bech32
func (this *NodeClient) ParseBech32ValConsPubkey(validatorInfoPubKeyBase64 string) (cryptotypes.PubKey, error) {
	validatorInfoPubKeyBytes, err := base64.StdEncoding.DecodeString(validatorInfoPubKeyBase64)
	if err != nil {
		return nil, err
	}
	pbk := ed25519.PubKey(validatorInfoPubKeyBytes) //ed25519
	pubkey, err := cryptocodec.FromTmPubKeyInterface(pbk)
	if err != nil {
		return nil, err
	}

	//pubkeyBech32, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey)
	//pubkeyBech32, err := legacybech32.MarshalPubKey(legacybech32.ConsPK, pubkey)
	//if err != nil {
	//	return nil, err
	//}
	return pubkey, nil
}
