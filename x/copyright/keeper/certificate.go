package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	authorizekey = "authorize_"
	pubkeys      = "pubkey_all"
	pubkey       = "pubkey_"
)


type AuthorizeAccount struct {
	ValidatorPubkey string `json:"validator_pubkey"`
	Account         string `json:"account"`
}

type AuthorizeMapObj struct {
	AuthorizeMap map[string]string `json:"authorize_map"`
}


func (k Keeper) GetBindValidatorPubkey(ctx sdk.Context, validatorConsPubAddr string) bool {
	store := ctx.KVStore(k.storeKey)
	validatorBindInfor := store.Get([]byte(authorizekey + validatorConsPubAddr))
	if validatorBindInfor == nil || len(validatorBindInfor) == 0 {
		return false
	} else {
		return true
	}
}



func (k Keeper) JudgeAuthorize(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(authorizekey + address)) {
		return true
	} else {
		return false
	}
}

func (k Keeper) Get(ctx sdk.Context, validatorConsPubAddr string) bool {
	store := ctx.KVStore(k.storeKey)
	validatorBindInfor := store.Get([]byte(pubkey + validatorConsPubAddr))
	if validatorBindInfor == nil || len(validatorBindInfor) == 0 {
		return false
	} else {
		return true
	}
}

func (k Keeper) SetAuthorize(ctx sdk.Context, address sdk.AccAddress, publickey string, consAddr sdk.ConsAddress) error {
	store := ctx.KVStore(k.storeKey)
	authorizeAccount := AuthorizeAccount{}
	if !store.Has([]byte(publickey)) {
		authorizeAccount.ValidatorPubkey = consAddr.String()
		authorizeAccount.Account = address.String()

		k.distributionKeeper.AuthenticationValidator(ctx, "", consAddr)
	} else {
		bz := store.Get([]byte(publickey))
		err := util.Json.Unmarshal(bz, &authorizeAccount)
		if err != nil {
			return err
		}
		if address.String() != authorizeAccount.Account {
			store.Delete([]byte(authorizekey + authorizeAccount.Account))
		}
		if authorizeAccount.ValidatorPubkey != consAddr.String() {
			k.distributionKeeper.AuthenticationValidator(ctx, authorizeAccount.ValidatorPubkey, consAddr)
		} else {
			authorizeAccount.ValidatorPubkey = consAddr.String()
		}
		authorizeAccount.Account = address.String()
	}
	authorizeAccountByte, err := util.Json.Marshal(authorizeAccount)
	if err != nil {
		return err
	}
	store.Set([]byte(authorizekey+address.String()), authorizeAccountByte)
	store.Set([]byte(pubkey+publickey), authorizeAccountByte)
	return nil
}

func queryPubkeyInfor(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAuthorizeAccountParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	store := ctx.KVStore(keeper.storeKey)
	return store.Get([]byte(pubkey + params.Pubkey)), nil
}

func queryAuthorizeValidator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAuthorizeValidatorParams
	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	consAddress, err := sdk.ConsAddressFromBech32(params.ConsAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	flag := keeper.distributionKeeper.GetAuthenticationValidator(ctx, consAddress)
	isExist := "0"
	if flag {
		isExist = "1"
	}
	return []byte(isExist), nil
}
