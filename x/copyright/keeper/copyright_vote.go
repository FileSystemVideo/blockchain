package keeper

import (
	"context"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/config"
	"fs.video/blockchain/x/copyright/types"
	logs "fs.video/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cid2 "github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/go-unixfs"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/shopspring/decimal"
	"path"
	"sort"
	"strings"
	"time"
)

type CopyrightVoteShareInfor struct {
	FavorTotal   decimal.Decimal `json:"favor_total"`
	AgainstTotal decimal.Decimal `json:"against_total"`
}

type DataHashVote struct {
	Status int             `json:"status"`
	Power  decimal.Decimal `json:"power"`
}

type AccountVoteInfor struct {
	Account   string          `json:"account"`
	Txhash    string          `json:"txhash"`
	Power     decimal.Decimal `json:"power"`
	MortgTime time.Time       `json:"mortg_time"`
	Status    int             `json:"status"`
}

const (
	copyrightApproveForResult = "copyright-approve-infor-"
	copyrightVoteFor          = "copyright-vote-infor-"
	copyrightVoteListFor      = "copyright-vote-list-infor-"
	copyrightVoteRedeem       = "copyright-vote-redeem-"
)

type CopyrightResult struct {
	DataHash   string    `json:"data_hash"`
	CreateTime time.Time `json:"create_time"`
}

func (k Keeper) CopyrightVoteRedeem(ctx sdk.Context) {
	copyrightVoteRedeemHash := k.GetCopyrightVoteRedeem(ctx)
	if copyrightVoteRedeemHash != "" {
		voteHashArray := strings.Split(copyrightVoteRedeemHash, "|")
		for i := 0; i < len(voteHashArray); i++ {
			voteHash := voteHashArray[i]
			accountVoteListInfor, err := k.QueryAccountCopyrightVoteList(ctx, voteHash)
			if err != nil {
				logs.Error("query vote infor error", err)
				return
			}
			var editorFlag bool
			var redeemNum int
			for j := 0; j < len(accountVoteListInfor); j++ {
				accountVoteInfor := accountVoteListInfor[j]
				if accountVoteInfor.Status == 1 {
					redeemNum += 1
					continue
				}
				endTime := accountVoteInfor.MortgTime.Add(config.CopyrightVoteRedeemTimePerioad)
				if endTime.After(ctx.BlockTime()) {
					continue
				}
				accountAddress, err := sdk.AccAddressFromBech32(accountVoteInfor.Account)
				if err != nil {
					logs.Error("format account error", err)
					return
				}
				floatPower, _ := accountVoteInfor.Power.Float64()
				powerDec := types.NewLedgerDec(floatPower)
				err = k.stakingKeeper.UnDelegationFreeze(ctx, accountAddress, powerDec)
				if err != nil {
					logs.Error("format account error", err)
					return
				}
				accountVoteInfor.Status = 1
				accountVoteListInfor[j] = accountVoteInfor
				redeemNum += 1
				editorFlag = true
				copyrightVoteRedeem := types.CopyrightVoteRedeem{}
				copyrightVoteRedeem.BlockNum = ctx.BlockHeight()
				copyrightVoteRedeem.TimeStemp = ctx.BlockTime().Unix()
				copyrightVoteRedeem.DataHash = voteHash
				copyrightVoteRedeem.TxHash = accountVoteInfor.Txhash
				k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVoteRedeem))
			}
			if editorFlag {
				k.UpdateAccountCopyrightVoteList(ctx, voteHash, accountVoteListInfor)
			}
			if redeemNum == len(accountVoteListInfor) {
				k.RemoveFromCopyrightVoteRedeem(ctx, voteHash)
			}

		}
	}
}

func (k Keeper) CalculateCopyrightVoteResult(ctx sdk.Context) {
	copyrightResult, err := k.GetCopyrightVoteForResult(ctx)
	if err != nil {
		logs.Error("query copyright infor error", err)
		panic(err)
	}
	for i := 0; i < len(copyrightResult); i++ {
		copyrightResult := copyrightResult[i]
		endTime := copyrightResult.CreateTime.Add(config.CopyrightVoteTimePerioad)
		if endTime.After(ctx.BlockTime()) { //版权审核还未到结束时间
			return
		}
		store := k.KVHelper(ctx)
		copyrightVoteInfor := make(map[string]CopyrightVoteShareInfor, 0)
		copyrightVoteByte := store.Get(copyrightVoteFor + copyrightResult.DataHash)
		if copyrightVoteByte != nil {
			err := util.Json.Unmarshal(copyrightVoteByte, &copyrightVoteInfor)
			if err != nil {
				return
			}
		} else {
			copyrightApprove := types.CopyrightApprove{}
			copyrightApprove.BlockNum = ctx.BlockHeight()
			copyrightApprove.TimeStemp = ctx.BlockTime().Unix()
			copyrightApprove.Status = 2
			copyrightApprove.DataHash = copyrightResult.DataHash
			copyrightApprove.PrimaryHash = copyrightResult.DataHash

			copyrighDataByte, err := k.GetCopyright(ctx, copyrightResult.DataHash)
			if err != nil {
				logs.Error("query copyright infor error", err)
				panic(err)
			}
			var copyrightData types.CopyrightData
			err = util.Json.Unmarshal(copyrighDataByte, &copyrightData)
			if err != nil {
				logs.Error("format data error", err)
				panic(err)
			}
			flag := k.UpdateAccountSpace(ctx, copyrightData.Creator, copyrightData.Size)
			if !flag {
				logs.Error("format data error", types.AccountSpaceReturnError)
				panic(err)
			}
			err = k.DeleteCopyrightInfor(ctx, copyrightResult.DataHash)
			if err != nil {
				logs.Error("delete data error", err)
				panic(err)
			}
			copyrightApprove.Creator = copyrightData.Creator.String()
			k.AddBlockRDS(ctx, types.NewBlockRD(copyrightApprove))
			k.RemoveFromCopyrightDataHash(ctx, copyrightResult.DataHash)
			continue
		}
		var linkMapKeySlice []string
		for key, _ := range copyrightVoteInfor {
			linkMapKeySlice = append(linkMapKeySlice, key)
		}
		sort.Strings(linkMapKeySlice)
		var abandonHash []string
		for i := 0; i < len(linkMapKeySlice); i++ {
			key := linkMapKeySlice[i] //cid
			voteInfor := copyrightVoteInfor[key]
			if voteInfor.FavorTotal.LessThanOrEqual(voteInfor.AgainstTotal) {
				abandonHash = append(abandonHash, key)
			}
		}
		err = k.copyrightApprovedDeal(ctx, copyrightResult.DataHash, abandonHash)
		if err != nil {
			logs.Error("deal copyright error", err)
			panic(err)
		}
	}
}


func (k Keeper) copyrightApprovedDeal(ctx sdk.Context, datahash string, abandonHash []string) error {
	copyrightApprove := types.CopyrightApprove{}
	copyrightApprove.BlockNum = ctx.BlockHeight()
	copyrightApprove.TimeStemp = ctx.BlockTime().Unix()
	copyrightApprove.PrimaryHash = datahash
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	storeKey := types.CopyrightDetailKey + datahash
	copyrightBytes := store.Get(storeKey)
	var copyrightData types.CopyrightData
	err := util.Json.Unmarshal(copyrightBytes, &copyrightData)
	if err != nil {
		return err
	}
	copyrightApprove.Creator = copyrightData.Creator.String()
	if len(abandonHash) > 0 {
		var removeName []string
		linkMap := copyrightData.LinkMap
		for i := 0; i < len(abandonHash); i++ {
			for key, link := range linkMap {
				if link.Cid == abandonHash[i] {
					delete(linkMap, key)
					removeName = append(removeName, link.Name)
					break
				}
			}
		}
		picLinkMap := copyrightData.PicLinkMap
		for i := 0; i < len(removeName); i++ {
			//nameArray := strings.Split(removeName[i], ".")
			name := fileNameSuffix(removeName[i])
			for key, _ := range picLinkMap {
				if strings.Contains(key, name) {
					delete(picLinkMap, key)
					continue
				}
			}
		}

		if len(linkMap) > 0 {

			var linkMapKeySlice []string
			for key, _ := range linkMap {
				linkMapKeySlice = append(linkMapKeySlice, key)
			}
			sort.Strings(linkMapKeySlice)

			node1 := unixfs.EmptyDirNode()
			dagService := dagutils.NewMemoryDagService()
			var newSize uint64
			for i := 0; i < len(linkMapKeySlice); i++ {
				key := linkMapKeySlice[i]
				link := linkMap[key]
				cid, err := cid2.Parse(link.Cid)
				if err != nil {
					logs.Error("query cid error", err)
					panic(err)
				}
				node1.AddRawLink(link.Name, &format.Link{
					Name: link.Name,
					Cid:  cid,
					Size: link.Size,
				})
				newSize += link.Size
			}
			err = dagService.Add(context.Background(), node1)
			if err != nil {
				logs.Error("6", err)
				panic(err)
			}
			getDagHash := func(editor *dagutils.Editor) string {
				nnode, err := editor.Finalize(context.Background(), dagService)
				if err != nil {
					logs.Error("serilize data error", err)
					return ""
				}
				newHash := ipath.IpfsPath(nnode.Cid())
				return newHash.String()
			}
			e := dagutils.NewDagEditor(node1, dagService)
			newHash := getDagHash(e)
			if newHash == "" {
				logs.Error("build data error", newHash)
				panic(err)
			}
			nodeBytes, err := node1.MarshalJSON()
			if err != nil {
				logs.Error("build data error", newHash)
				panic(err)
			}
			newHash = strings.Replace(newHash, "/ipfs/", "", 1)
			picNodeByte := buildNewPicHash(ctx, picLinkMap)
			copyrightApprove.DataHash = newHash
			copyrightApprove.PrimaryHash = datahash
			copyrightApprove.Status = 1
			copyrightApprove.RemoveName = removeName
			copyrightApprove.NodeByte = string(nodeBytes)
			copyrightApprove.NewSize = newSize
			copyrightApprove.PicNodeByte = picNodeByte
			copyrightApprove.EditorStatus = 1
			err = k.UpdateCopyrightStatus(ctx, datahash, 1)
			if err != nil {
				return err
			}
			err = k.UpdateCopyrightInfor(ctx, newHash, datahash, newSize)
			if err != nil {
				return err
			}
		} else {
			copyrightApprove.Status = 2
			copyrightApprove.DataHash = datahash
			err = k.DeleteCopyrightInfor(ctx, datahash)
			if err != nil {
				return err
			}
			flag := k.UpdateAccountSpace(ctx, copyrightData.Creator, copyrightData.Size)
			if !flag {
				return types.AccountSpaceReturnError
			}
		}
	} else {
		copyrightApprove.Status = 1
		copyrightApprove.DataHash = datahash
		err := k.UpdateCopyrightStatus(ctx, datahash, 1)
		if err != nil {
			return err
		}
	}
	k.AddBlockRDS(ctx, types.NewBlockRD(copyrightApprove))
	k.RemoveFromCopyrightDataHash(ctx, datahash)
	return nil
}


func (k Keeper) dealCopyrightVote(ctx sdk.Context, account, sourceName, copyrightHash, dataPower, txhash string, linkMap map[string]types.Link) error {

	var dataVoteMap map[string]DataHashVote
	err := util.Json.Unmarshal([]byte(dataPower), &dataVoteMap)
	if err != nil {
		return err
	}
	store := k.KVHelper(ctx)
	voteKey := copyrightVoteFor + copyrightHash
	copyrightVoteInfor := make(map[string]CopyrightVoteShareInfor, 0)
	copyrightVoteByte := store.Get(voteKey)
	if copyrightVoteByte != nil {
		err := util.Json.Unmarshal(copyrightVoteByte, &copyrightVoteInfor)
		if err != nil {
			return err
		}
	}
	var linkMapKeySlice []string
	for key, _ := range linkMap {
		linkMapKeySlice = append(linkMapKeySlice, key)
	}
	sort.Strings(linkMapKeySlice)
	var totalVote decimal.Decimal
	for _, linkKey := range linkMapKeySlice {
		link := linkMap[linkKey]
		key := link.Cid
		if vote, ok := copyrightVoteInfor[key]; ok {
			if dataVote, ok := dataVoteMap[key]; ok {
				if dataVote.Status == 1 {
					vote.FavorTotal = vote.FavorTotal.Add(dataVote.Power)
				} else {
					vote.AgainstTotal = vote.AgainstTotal.Add(dataVote.Power)
				}
				copyrightVoteInfor[key] = vote
				totalVote = totalVote.Add(dataVote.Power)
				award := getVoteAward(ctx,dataVote.Power)
				copyrightVote := types.CopyrightVote{
					DataHash:   copyrightHash,
					VideoHash:  key,
					SourceName: sourceName,
					VideoName:  link.Name,
					BlockNum:   ctx.BlockHeight(),
					TimeStemp:  ctx.BlockTime().Unix(),
					TxHash:     txhash,
					Account:    account,
					VoteStatus: dataVote.Status,
					Award:      util.DecimalStringFixed(award.String(), config.CoinPlaces),
					Power:      dataVote.Power.String(),
				}
				k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVote))
			}
		} else {
			if dataVote, ok := dataVoteMap[key]; ok {
				defaultVote := CopyrightVoteShareInfor{
					FavorTotal:   decimal.Zero,
					AgainstTotal: decimal.Zero,
				}
				if dataVote.Status == 1 {
					defaultVote.FavorTotal = dataVote.Power
				} else {
					defaultVote.AgainstTotal = dataVote.Power
				}
				copyrightVoteInfor[key] = defaultVote
				totalVote = totalVote.Add(dataVote.Power)
				award := getVoteAward(ctx,dataVote.Power)
				copyrightVote := types.CopyrightVote{
					DataHash:   copyrightHash,
					VideoHash:  key,
					SourceName: sourceName,
					VideoName:  link.Name,
					BlockNum:   ctx.BlockHeight(),
					TimeStemp:  ctx.BlockTime().Unix(),
					TxHash:     txhash,
					Account:    account,
					VoteStatus: dataVote.Status,
					Award:      util.DecimalStringFixed(award.String(), config.CoinPlaces),
					Power:      dataVote.Power.String(),
				}
				k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVote))
			}
		}
	}
	if totalVote.GreaterThan(decimal.Zero) {
		err = store.Set(voteKey, copyrightVoteInfor)
		if err != nil {
			return err
		}
		copyrightVoteListInfor, err := k.QueryAccountCopyrightVoteList(ctx, copyrightHash)
		if err != nil {
			return err
		}
		copyrightVoteListKey := copyrightVoteListFor + copyrightHash
		accountVoteInfor := AccountVoteInfor{
			Account:   account,
			Power:     totalVote,
			MortgTime: ctx.BlockTime(),
			Txhash:    txhash,
		}
		copyrightVoteListInfor = append(copyrightVoteListInfor, accountVoteInfor)
		copyrightVoteListByte, err := util.Json.Marshal(copyrightVoteListInfor)
		if err != nil {
			return err
		}
		store.Set(copyrightVoteListKey, copyrightVoteListByte)
		k.SetCopyrightVoteRedeem(ctx, copyrightHash)
		accountAddress, err := sdk.AccAddressFromBech32(accountVoteInfor.Account)
		if err != nil {
			logs.Error("format account error", err)
			return err
		}
		floatPower, _ := accountVoteInfor.Power.Float64()
		powerDec := types.NewLedgerDec(floatPower)
		totalShares, _ := k.GetAccountDelegatorShares(ctx, accountAddress)
		totalSharesFloat, _ := decimal.RequireFromString(totalShares).Float64()
		totalSharesLedge := types.NewLedgerDec(totalSharesFloat)
		freezeShares, err := k.stakingKeeper.GetDelegationFreeze(ctx, accountAddress)
		if totalSharesLedge.Sub(freezeShares).LT(powerDec) {
			return types.CopyrightVoteNotEnoughErr
		}
		err = k.stakingKeeper.DelegationFreeze(ctx, accountAddress, powerDec)
		if err != nil {
			logs.Error("format account error", err)
			return err
		}
		award := getVoteAward(ctx,totalVote)
		awardString := util.DecimalStringFixed(award.String(), config.CoinPlaces)
		realCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, awardString)
		ledgerCoin := types.MustRealCoin2LedgerCoin(realCoin)

		copntractAddress := authtypes.NewModuleAddress(config.KeyCopyrighDeflation)
		if !k.CoinKeeper.HasBalance(ctx,copntractAddress,ledgerCoin){
			ledgerCoin = k.CoinKeeper.GetBalance(ctx,copntractAddress,sdk.DefaultBondDenom)
		}
		if !ledgerCoin.IsPositive(){
			return nil
		}
		err = k.CoinKeeper.SendCoinsFromModuleToAccount(ctx, config.KeyCopyrighDeflation, accountAddress, sdk.NewCoins(ledgerCoin))
		if err != nil {
			logs.Error("vote award error", err)
			return err
		}
		awardDecimal := decimal.RequireFromString(awardString)
		k.UpdateDeflationMinerInforByVote(ctx,awardDecimal)
		err = k.SetSpaceMinerBonusAmount(ctx, awardDecimal)
		if err != nil {
			return err
		}
		feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
		voteAwardTradeInfor := types.TradeInfor{
			From:           copntractAddress.String(),
			To:             account,
			Txhash:         txhash,
			TradeType:      types.TradeTypeCopyrightVoteReward,
			Amount:         realCoin,
			Fee:            feeCoin,
			BlockNum:       ctx.BlockHeight(),
			TimeStemp:      ctx.BlockTime().Unix(),
			FromFsvBalance: k.GetBalance(ctx, config.MainToken, copntractAddress),
			FromTipBalance: k.GetBalance(ctx, config.InviteToken, copntractAddress),
			ToFsvBalance:   k.GetBalance(ctx, config.MainToken, accountAddress),
			ToTipBalance:   k.GetBalance(ctx, config.InviteToken, accountAddress),
		}
		k.AddBlockRDS(ctx, types.NewBlockRD(voteAwardTradeInfor))
		return nil
	} else {
		return types.CopyrightVoteInvalidErr
	}

}

func (k Keeper) SetCopyrightVoteRedeem(ctx sdk.Context, dataHash string) {
	store := k.KVHelper(ctx)
	copyrightVoteRedeemByte := store.Get(copyrightVoteRedeem)
	if copyrightVoteRedeemByte == nil {
		store.Set(copyrightVoteRedeem, dataHash)
	} else {
		copyrightVoteRedeemString := string(copyrightVoteRedeemByte)
		if !strings.Contains(copyrightVoteRedeemString, dataHash) {
			copyrightVoteRedeemString += "|" + dataHash
			store.Set(copyrightVoteRedeem, copyrightVoteRedeemString)
		}
	}

}

func (k Keeper) RemoveFromCopyrightVoteRedeem(ctx sdk.Context, dataHash string) error {
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightVoteRedeem)
	var copyrightVoteRedeemString string
	if bz != nil {
		copyrightVoteRedeemString = string(bz)
	}
	if copyrightVoteRedeemString != "" && strings.Contains(copyrightVoteRedeemString, dataHash) {
		idsStringArray := strings.Split(copyrightVoteRedeemString, "|")
		index := 0
		for i := 0; i < len(idsStringArray); i++ {
			if idsStringArray[i] == dataHash {
				index = i
				break
			}
		}
		idsStringArray = append(idsStringArray[:index], idsStringArray[index+1:]...)
		//idsStringArray = idsStringArray[:index] + idsStringArray[index+1:]
		copyrightVoteRedeemString = strings.Join(idsStringArray, "|")
		if copyrightVoteRedeemString == "" {
			store.Delete(copyrightVoteRedeem)
		} else {
			store.Set(copyrightVoteRedeem, copyrightVoteRedeemString)
		}
	}
	return nil
}

func (k Keeper) GetCopyrightVoteRedeem(ctx sdk.Context) string {
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightVoteRedeem)
	if bz != nil {
		return string(bz)
	}
	return ""
}

func (k Keeper) SetCopyrightVoteForResult(ctx sdk.Context, dataHash string) error {
	//store := ctx.KVStore(k.storeKey)
	store := k.KVHelper(ctx)
	complainIdListByte := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if complainIdListByte != nil {
		err := util.Json.Unmarshal(complainIdListByte, &copyrightArray)
		if err != nil {
			return err
		}
	}
	copyrightResult := CopyrightResult{
		DataHash:   dataHash,
		CreateTime: ctx.BlockTime(),
	}
	copyrightArray = append(copyrightArray, copyrightResult)
	return store.Set(copyrightApproveForResult, copyrightArray)
}

func (k Keeper) RemoveFromCopyrightDataHash(ctx sdk.Context, datahash string) error {
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if bz != nil {
		err := util.Json.Unmarshal(bz, &copyrightArray)
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	if len(copyrightArray) > 0 {
		for i := 0; i < len(copyrightArray); i++ {
			if copyrightArray[i].DataHash == datahash {
				copyrightArray = append(copyrightArray[:i], copyrightArray[i+1:]...)
				break
			}
		}
		if len(copyrightArray) == 0 {
			store.Delete(copyrightApproveForResult)
		} else {
			copyrightArrayBytes, err := util.Json.Marshal(copyrightArray)
			if err != nil {
				return err
			}
			store.Set(copyrightApproveForResult, copyrightArrayBytes)
		}
	}
	return nil
}

func (k Keeper) GetCopyrightVoteForResult(ctx sdk.Context) ([]CopyrightResult, error) {
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if bz != nil {
		err := util.Json.Unmarshal(bz, &copyrightArray)
		if err != nil {
			return copyrightArray, err
		}
	}
	return copyrightArray, nil
}

func (k Keeper) QueryAccountCopyrightVoteList(ctx sdk.Context, copyrightHash string) ([]AccountVoteInfor, error) {
	store := k.KVHelper(ctx)
	//copyrightAccountVoteList := make(map[string]AccountVoteInfor, 0)
	copyrightVoteListInfor := make([]AccountVoteInfor, 0)
	copyrightVoteListByte := store.Get(copyrightVoteListFor + copyrightHash)
	if copyrightVoteListByte != nil {
		err := util.Json.Unmarshal(copyrightVoteListByte, &copyrightVoteListInfor)
		if err != nil {
			return copyrightVoteListInfor, err
		}
	}
	return copyrightVoteListInfor, nil
}

func (k Keeper) UpdateAccountCopyrightVoteList(ctx sdk.Context, copyrightHash string, accountVoteList []AccountVoteInfor) error {
	store := k.KVHelper(ctx)
	store.Set(copyrightVoteListFor+copyrightHash, accountVoteList)
	return nil
}

func buildNewPicHash(ctx sdk.Context, picMap map[string]types.Link) string {
	node1 := unixfs.EmptyDirNode()

	var linkMapKeySlice []string
	for key, _ := range picMap {
		linkMapKeySlice = append(linkMapKeySlice, key)
	}
	sort.Strings(linkMapKeySlice)
	var newSize uint64
	for i := 0; i < len(linkMapKeySlice); i++ {
		key := linkMapKeySlice[i]
		link := picMap[key]
		cid, err := cid2.Parse(link.Cid)
		if err != nil {
			logs.Error("query cid error", err)
			panic(err)
		}
		node1.AddRawLink(link.Name, &format.Link{
			Name: link.Name,
			Cid:  cid,
			Size: link.Size,
		})
		newSize += link.Size
	}
	picByte, err := node1.MarshalJSON()
	if err != nil {
		logs.Error("query cid error", err)
		panic(err)
	}
	dagService := dagutils.NewMemoryDagService()
	err = dagService.Add(context.Background(), node1)
	if err != nil {
		logs.Error("6", err)
		panic(err)
	}
	getDagHash := func(editor *dagutils.Editor) string {
		nnode, err := editor.Finalize(context.Background(), dagService)
		if err != nil {
			logs.Error("query cid error", err)
			return ""
		}
		newHash := ipath.IpfsPath(nnode.Cid())
		return newHash.String()
	}
	e := dagutils.NewDagEditor(node1, dagService)
	newHash := getDagHash(e)
	logs.Info("newhash**********************************************", newHash)
	return string(picByte)
}

func fileNameSuffix(fileName string) string {
	var fileSuffix string
	fileSuffix = path.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	if filenameOnly == "" {
		filenameOnly = strings.Replace(fileName, ".", "", -1)
	}
	return filenameOnly
}

func getVoteAward(ctx sdk.Context,power decimal.Decimal)decimal.Decimal{
	var award decimal.Decimal
	award = power.Mul(config.CopyrightVoteAwardRate)
	return award
}
