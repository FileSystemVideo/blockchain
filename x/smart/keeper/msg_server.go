package keeper

import (
	"context"
	"encoding/hex"
	"fs.video/blockchain/x/smart/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	Keeper
	logPrefix string
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper, logPrefix: "smart | msgServer | "}
}


func (k msgServer) CreateSmartValidator(goCtx context.Context, msg *types.MsgCreateSmartValidator) (*types.MsgEmptyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	valAddress := sdk.ValAddress(addr)
	_, found := k.stakingKeeper.GetValidator(ctx, valAddress)
	if found {
		return nil, stakingTypes.ErrValidatorOwnerExists
	}
	err = k.createValidator(ctx, addr, valAddress, *msg, msg.Value)
	if err != nil {
		return &types.MsgEmptyResponse{}, err
	}
	return &types.MsgEmptyResponse{}, nil
}

//base64bech32
func ParseBech32ValConsPubkey(validatorInfoPubKeyBase64 string) (cryptotypes.PubKey, error) {
	//validatorInfoPubKeyBytes, err := base64.StdEncoding.DecodeString(validatorInfoPubKeyBase64)
	validatorInfoPubKeyBytes, err := hex.DecodeString(validatorInfoPubKeyBase64)
	if err != nil {
		return nil, err
	}
	pbk := ed25519.PubKey(validatorInfoPubKeyBytes) //ed25519
	pubkey, err := cryptocodec.FromTmPubKeyInterface(pbk)
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}
