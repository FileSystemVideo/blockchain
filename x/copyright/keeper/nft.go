package keeper

import (
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	nftTokenIdKey = "nft_token_" //nft
	//nftAccountKey = "nft_account_"                  //nft
	nftApproveKey = "nft_approve_tfajg334eqvw1ga3_" //nft
)

type NftInfor struct {
	TokenId      string `json:"token_id"`
	EnghlishName string `json:"english_name"`
	ChineseName  string `json:"chinese_name"`
	Owner        string `json:"owner"`
	MetaData     string `json:"meta_data"`
}

func (k Keeper) HandleNftTransfer(ctx sdk.Context, nftTransfer types.NftTransferData) error {
	err := k.NftTransfer(ctx, nftTransfer.TokenId, nftTransfer.From.String(), nftTransfer.To.String())
	if err != nil {
		return err
	}
	err = k.SetCopyrightAuthor(ctx, nftTransfer.TokenId, nftTransfer.To)
	if err != nil {
		return err
	}
	copyrightByte, err := k.GetCopyright(ctx, nftTransfer.TokenId)
	if err != nil {
		return err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightByte, &copyrightData)
	if err != nil {
		return err
	}
	if copyrightData.Size > 0 {
		flag := k.UpdateAccountSpace(ctx, nftTransfer.From, copyrightData.Size)
		if !flag {
			return types.AccountSpaceReturnError
		}
		flag = k.UpdateAccountSpaceUsed(ctx, nftTransfer.To, copyrightData.Size)
		if !flag {
			return types.AccountSpaceNotEnough
		}
	}
	return nil
}

func (k Keeper) NftTransfer(ctx sdk.Context, tokenId, owner, account string) error {
	if tokenId == "" {
		return sdkerrors.Wrap(types.TokenIdEmpty, "")
	}
	store := k.KVHelper(ctx)
	var nftInfor NftInfor
	key := nftTokenIdKey + tokenId
	bz := store.Get(key)
	if bz != nil {
		util.Json.Unmarshal(bz, &nftInfor)
	}
	if nftInfor.Owner == "" {
		return sdkerrors.Wrap(types.TokenidNotExist, "")
	}
	if nftInfor.Owner != owner {
		return sdkerrors.Wrap(types.TokenidHasNoRight, "")
	}

	//oldOwer := nftInfor.Owner
	nftInfor.Owner = account
	nftInforByte, err := util.Json.Marshal(nftInfor)
	if err != nil {
		return sdkerrors.Wrap(types.TokenidFormatErr, "")
	}
	store.Set(key, nftInforByte)
	return nil
}

func (k Keeper) QueryNftInfor(ctx sdk.Context, tokenId string) NftInfor {
	store := k.KVHelper(ctx)
	var nftInfor NftInfor
	bz := store.Get(nftTokenIdKey + tokenId)
	if bz != nil {
		util.Json.Unmarshal(bz, &nftInfor)
	}
	return nftInfor
}

func (k Keeper) BuildNft(ctx sdk.Context, tokenId, account, englishName, chineseName string) error {
	store := k.KVHelper(ctx)
	var nftInfor NftInfor
	key := nftTokenIdKey + tokenId
	bz := store.Get(key)
	if bz != nil {
		return sdkerrors.Wrap(types.TokenidHasExist, "")
	}
	nftInfor.TokenId = tokenId
	nftInfor.Owner = account
	nftInfor.ChineseName = chineseName
	nftInfor.EnghlishName = englishName
	nftInfor.MetaData = tokenId
	return store.Set(key, nftInfor)
}

func (k Keeper) EditNft(ctx sdk.Context, tokenId, englishName, chineseName string) error {
	store := k.KVHelper(ctx)
	var nftInfor NftInfor
	key := nftTokenIdKey + tokenId
	bz := store.Get(key)
	if bz != nil {
		util.Json.Unmarshal(bz, &nftInfor)
	}
	nftInfor.TokenId = tokenId
	nftInfor.ChineseName = chineseName
	nftInfor.EnghlishName = englishName
	nftInfor.MetaData = tokenId
	return store.Set(key, nftInfor)
}

func (k Keeper) DeleteNft(ctx sdk.Context, tokenId string) {
	store := k.KVHelper(ctx)
	key := nftTokenIdKey + tokenId
	store.Delete(key)
}

func (k Keeper) UpdateCopyrightInfor(ctx sdk.Context, newdatahash, datahash string, newSize uint64) error {
	store := k.KVHelper(ctx)
	copyrightKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(copyrightKey)
	var copyrightData types.CopyrightData
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		return err
	}
	copyrightData.DataHash = newdatahash
	copyrightKeyNew := types.CopyrightDetailKey + newdatahash
	copyrightBytes, err = util.Json.Marshal(copyrightData)
	if err != nil {
		return err
	}
	oldSize := copyrightData.Size
	copyrightData.Size = int64(newSize)
	returnSize := oldSize - int64(newSize)
	store.Set(copyrightKeyNew, copyrightBytes)
	store.Delete(copyrightKey)
	flag := k.UpdateAccountSpace(ctx, copyrightData.Creator, returnSize)
	if !flag {
		return types.AccountSpaceReturnError
	}
	key := types.CopyrightIpKey + datahash
	var copyrightIpInfor CopyrightIpinfor
	err = store.GetUnmarshal(key, &copyrightIpInfor)
	if err != nil {
		return err
	}
	newKey := types.CopyrightIpKey + newdatahash
	store.Set(newKey, copyrightIpInfor)
	store.Delete(key)
	var nftInfor NftInfor
	nftKey := nftTokenIdKey + datahash
	bz := store.Get(nftKey)
	if bz != nil {
		util.Json.Unmarshal(bz, &nftInfor)
	}
	nftInfor.TokenId = newdatahash
	nftkeyNew := nftTokenIdKey + newdatahash
	err = store.Set(nftkeyNew, nftInfor)
	if err != nil {
		return err
	}
	store.Delete(nftkeyNew)
	return nil
}
