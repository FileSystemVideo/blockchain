package keeper

import (
	"context"
	"errors"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cid2 "github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	"github.com/ipfs/go-merkledag/dagutils"
	"github.com/ipfs/go-unixfs"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
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
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	copyrightVoteRedeemHash := k.GetCopyrightVoteRedeem(ctx)
	if copyrightVoteRedeemHash != "" {
		voteHashArray := strings.Split(copyrightVoteRedeemHash, "|")
		for i := 0; i < len(voteHashArray); i++ {
			voteHash := voteHashArray[i]
			accountVoteListInfor, err := k.QueryAccountCopyrightVoteList(ctx, voteHash)
			if err != nil {
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
				endTime := accountVoteInfor.MortgTime.Add(core.CopyrightVoteRedeemTimePerioad)
				if endTime.After(ctx.BlockTime()) { 
					continue
				}
				
				accountAddress, err := sdk.AccAddressFromBech32(accountVoteInfor.Account)
				if err != nil {
					log.WithError(err).WithField("Account", accountVoteInfor.Account).Error("AccAddressFromBech32")
					return
				}
				powerDec := types.NewLedgerDec(accountVoteInfor.Power)
				err = k.stakingKeeper.UnDelegationFreeze(ctx, accountAddress, powerDec)
				if err != nil {
					log.WithError(err).Error("UnDelegationFreeze")
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
				//rd ,
				err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVoteRedeem))
				if err != nil {
					log.WithError(err).Error("AddBlockRDS-copyrightVoteRedeem")
					panic(err)
				}
			}
			if editorFlag { 
				err = k.UpdateAccountCopyrightVoteList(ctx, voteHash, accountVoteListInfor)
				if err != nil {
					panic(err)
				}
			}
			if redeemNum == len(accountVoteListInfor) { 
				err = k.RemoveFromCopyrightVoteRedeem(ctx, voteHash)
				if err != nil {
					panic(err)
				}
			}

		}
	}
}


func (k Keeper) CalculateCopyrightVoteResult(ctx sdk.Context) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	copyrightResult, err := k.GetCopyrightVoteForResult(ctx)
	if err != nil {
		panic(err)
	}
	
	for i := 0; i < len(copyrightResult); i++ {
		copyrightResult := copyrightResult[i]
		endTime := copyrightResult.CreateTime.Add(core.CopyrightVoteTimePerioad)
		if endTime.After(ctx.BlockTime()) { 
			return
		}
		
		store := k.KVHelper(ctx)
		copyrightVoteInfor := make(map[string]CopyrightVoteShareInfor, 0)
		copyrightVoteByte := store.Get(copyrightVoteFor + copyrightResult.DataHash)
		if copyrightVoteByte != nil {
			err := util.Json.Unmarshal(copyrightVoteByte, &copyrightVoteInfor)
			if err != nil {
				log.WithError(err).Error("Unmarshal")
				panic(err)
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
				panic(err)
			}
			var copyrightData types.CopyrightData
			err = util.Json.Unmarshal(copyrighDataByte, &copyrightData)
			if err != nil {
				log.WithError(err).Error("Unmarshal")
				panic(err)
			}
			
			flag := k.UpdateAccountSpace(ctx, copyrightData.Creator, copyrightData.Size)
			if !flag {
				log.WithError(err).WithFields(logrus.Fields{
					"address": copyrightData.Creator.String(),
					"size":    copyrightData.Size,
				}).Error("UpdateAccountSpace")
				panic(err)
			}
			err = k.DeleteCopyrightInfor(ctx, copyrightResult.DataHash)
			if err != nil {
				panic(err)
			}
			copyrightApprove.Creator = copyrightData.Creator.String()
			//rd ,
			err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightApprove))
			if err != nil {
				log.WithError(err).Error("AddBlockRDS-copyrightApprove")
				panic(err)
			}
			
			err = k.RemoveFromCopyrightDataHash(ctx, copyrightResult.DataHash)
			if err != nil {
				panic(err)
			}
			continue
		}
		var linkMapKeySlice []string
		//map key
		for key, _ := range copyrightVoteInfor {
			linkMapKeySlice = append(linkMapKeySlice, key)
		}
		sort.Strings(linkMapKeySlice)
		var abandonHash []string
		for i := 0; i < len(linkMapKeySlice); i++ {
			key := linkMapKeySlice[i] //cid hash
			voteInfor := copyrightVoteInfor[key]
			if voteInfor.FavorTotal.LessThanOrEqual(voteInfor.AgainstTotal) {
				
				abandonHash = append(abandonHash, key)
			}
		}
		
		err = k.copyrightApprovedDeal(ctx, copyrightResult.DataHash, abandonHash)
		if err != nil {
			panic(err)
		}
	}
}

/*
	
 	datahash hash
 	abandonHash hash
*/
func (k Keeper) copyrightApprovedDeal(ctx sdk.Context, datahash string, abandonHash []string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
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
		log.WithError(err).Error("Unmarshal")
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
			//map key
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
					log.WithError(err).Error("cid2.Parse")
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
				log.WithError(err).Error("dagService.Add")
				panic(err)
			}
			getDagHash := func(editor *dagutils.Editor) string {
				nnode, err := editor.Finalize(context.Background(), dagService)
				if err != nil {
					log.WithError(err).Error("editor.Finalize")
					return ""
				}
				newHash := ipath.IpfsPath(nnode.Cid())
				return newHash.String()
			}
			e := dagutils.NewDagEditor(node1, dagService)
			newHash := getDagHash(e)
			if newHash == "" {
				log.WithError(err).Error("getDagHash")
				panic(errors.New("getDagHash error"))
			}
			nodeBytes, err := node1.MarshalJSON()
			if err != nil {
				log.WithError(err).Error("MarshalJSON")
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
			//,hash
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
	//rd ,
	err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightApprove))
	if err != nil {
		log.WithError(err).Error("AddBlockRDS-copyrightApprove")
	}
	
	store.Delete(copyrightVoteFor + datahash)

	
	return k.RemoveFromCopyrightDataHash(ctx, datahash)
}

/**
dataPower ,
*/
func (k Keeper) dealCopyrightVote(ctx sdk.Context, account, sourceName, copyrightHash, dataPower, txhash string, linkMap map[string]types.Link) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	var dataVoteMap map[string]DataHashVote
	err := util.Json.Unmarshal([]byte(dataPower), &dataVoteMap)
	if err != nil {
		log.WithError(err).Error("Unmarshal1")
		return err
	}
	
	store := k.KVHelper(ctx)
	voteKey := copyrightVoteFor + copyrightHash
	copyrightVoteInfor := make(map[string]CopyrightVoteShareInfor, 0)
	copyrightVoteByte := store.Get(voteKey)
	if copyrightVoteByte != nil {
		err := util.Json.Unmarshal(copyrightVoteByte, &copyrightVoteInfor)
		if err != nil {
			log.WithError(err).Error("Unmarshal2")
			return err
		}
	}
	var linkMapKeySlice []string
	//map key
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
				award := getVoteAward(ctx, dataVote.Power)
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
					Award:      util.DecimalStringFixed(award.String(), core.CoinPlaces),
					Power:      dataVote.Power.String(),
				}
				err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVote))
				if err != nil {
					return err
				}
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
				award := getVoteAward(ctx, dataVote.Power)
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
					Award:      util.DecimalStringFixed(award.String(), core.CoinPlaces),
					Power:      dataVote.Power.String(),
				}
				err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightVote))
				if err != nil {
					return err
				}
			}
		}
	}
	if totalVote.GreaterThan(decimal.Zero) {
		err = store.Set(voteKey, copyrightVoteInfor)
		if err != nil {
			log.WithError(err).Error("store.Set")
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
			log.WithError(err).Error("Marshal1")
			return err
		}
		err = store.Set(copyrightVoteListKey, copyrightVoteListByte)
		if err != nil {
			log.WithError(err).Error("Marshal2")
			return err
		}
		
		k.SetCopyrightVoteRedeem(ctx, copyrightHash)
		//staking 
		accountAddress, err := sdk.AccAddressFromBech32(accountVoteInfor.Account)
		if err != nil {
			log.WithError(err).WithField("Account", accountVoteInfor.Account).Error("AccAddressFromBech32")
			return err
		}
		//floatPower, _ := accountVoteInfor.Power.Float64()
		powerDec := types.NewLedgerDec(accountVoteInfor.Power)
		totalShares, _ := k.GetAccountDelegatorShares(ctx, accountAddress) //dpos
		totalSharesDecimal := decimal.RequireFromString(totalShares)

		totalSharesLedge := types.NewLedgerDec(totalSharesDecimal)
		freezeShares, _ := k.stakingKeeper.GetDelegationFreeze(ctx, accountAddress) 

		if totalSharesLedge.Sub(freezeShares).LT(powerDec) { // -  < 
			return types.CopyrightVoteNotEnoughErr
		}
		err = k.stakingKeeper.DelegationFreeze(ctx, accountAddress, powerDec)
		if err != nil {
			log.WithError(err).Error("DelegationFreeze")
			return err
		}
		
		//award := totalVote.Mul(config.CopyrightVoteAwardRate)
		award := getVoteAward(ctx, totalVote)
		awardString := util.DecimalStringFixed(award.String(), core.CoinPlaces)
		realCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, awardString)
		ledgerCoin := types.MustRealCoin2LedgerCoin(realCoin)

		copntractAddress := authtypes.NewModuleAddress(core.KeyCopyrighDeflation)
		
		if !k.CoinKeeper.HasBalance(ctx, copntractAddress, ledgerCoin) {
			ledgerCoin = k.CoinKeeper.GetBalance(ctx, copntractAddress, sdk.DefaultBondDenom)
		}
		
		if !ledgerCoin.IsPositive() {
			return nil
		}
		err = k.CoinKeeper.SendCoinsFromModuleToAccount(ctx, core.KeyCopyrighDeflation, accountAddress, sdk.NewCoins(ledgerCoin))
		if err != nil {
			log.WithError(err).WithFields(logrus.Fields{
				"fromAddr": core.KeyCopyrighDeflation,
				"toAddr":   accountAddress.String(),
				"amt":      ledgerCoin.String(),
			}).Error("SendCoins")
			return err
		}
		
		awardDecimal := decimal.RequireFromString(awardString)
		k.UpdateDeflationMinerInforByVote(ctx, awardDecimal)
		
		err = k.SetSpaceMinerBonusAmount(ctx, awardDecimal)
		if err != nil {
			return err
		}
		
		feeCoin := types.NewRealCoinFromStr(sdk.DefaultBondDenom, "0")
		voteAwardTradeInfor := types.TradeInfor{
			From:           copntractAddress.String(),
			To:             account,
			Txhash:         txhash,
			TradeType:      core.TradeTypeCopyrightVoteReward,
			Amount:         realCoin,
			Fee:            feeCoin,
			BlockNum:       ctx.BlockHeight(),
			TimeStemp:      ctx.BlockTime().Unix(),
			FromFsvBalance: k.GetBalance(ctx, core.MainToken, copntractAddress),
			FromTipBalance: k.GetBalance(ctx, core.InviteToken, copntractAddress),
			ToFsvBalance:   k.GetBalance(ctx, core.MainToken, accountAddress),
			ToTipBalance:   k.GetBalance(ctx, core.InviteToken, accountAddress),
		}
		err = k.AddBlockRDS(ctx, types.NewBlockRD(voteAwardTradeInfor))
		if err != nil {
			log.WithError(err).Error("AddBlockRDS-voteAwardTradeInfor")
			return err
		}
		return nil
	} else {
		return types.CopyrightVoteInvalidErr
	}

}


func (k Keeper) SetCopyrightVoteRedeem(ctx sdk.Context, dataHash string) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	var err error
	copyrightVoteRedeemByte := store.Get(copyrightVoteRedeem)
	if copyrightVoteRedeemByte == nil {
		err = store.Set(copyrightVoteRedeem, dataHash)
		if err != nil {
			log.WithError(err).Error("store.Set1")
		}
	} else {
		copyrightVoteRedeemString := string(copyrightVoteRedeemByte)
		if !strings.Contains(copyrightVoteRedeemString, dataHash) {
			copyrightVoteRedeemString += "|" + dataHash
			err = store.Set(copyrightVoteRedeem, copyrightVoteRedeemString)
			if err != nil {
				log.WithError(err).Error("store.Set2")
			}
		}
	}
	if err != nil {
		panic(err)
	}
}


func (k Keeper) RemoveFromCopyrightVoteRedeem(ctx sdk.Context, dataHash string) (err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
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
		copyrightVoteRedeemString = strings.Join(idsStringArray, "|")
		if copyrightVoteRedeemString == "" {
			
			store.Delete(copyrightVoteRedeem)
		} else {
			err = store.Set(copyrightVoteRedeem, copyrightVoteRedeemString)
			if err != nil {
				log.WithError(err).Error("store.Set")
			}
		}
	}
	return err
}


func (k Keeper) GetCopyrightVoteRedeem(ctx sdk.Context) string {
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightVoteRedeem)
	if bz != nil {
		return string(bz)
	}
	return ""
}


func (k Keeper) SetCopyrightVoteForResult(ctx sdk.Context, dataHash string) (err error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	complainIdListByte := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if complainIdListByte != nil {
		err = util.Json.Unmarshal(complainIdListByte, &copyrightArray)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
	}
	copyrightResult := CopyrightResult{
		DataHash:   dataHash,
		CreateTime: ctx.BlockTime(),
	}
	copyrightArray = append(copyrightArray, copyrightResult)
	log.WithField("data type", copyrightArray).Debug("SetCopyrightVoteForResult")
	err = store.Set(copyrightApproveForResult, copyrightArray)
	if err != nil {
		log.WithError(err).Error("store.Set")
	}
	return err
}


func (k Keeper) RemoveFromCopyrightDataHash(ctx sdk.Context, datahash string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if bz != nil {
		err := util.Json.Unmarshal(bz, &copyrightArray)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
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
				log.WithError(err).Error("Marshal")
				return err
			}
			err = store.Set(copyrightApproveForResult, copyrightArrayBytes)
			if err != nil {
				log.WithError(err).Error("store.Set")
				return err
			}
		}
	}
	return nil
}


func (k Keeper) GetCopyrightVoteForResult(ctx sdk.Context) ([]CopyrightResult, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(copyrightApproveForResult)
	var copyrightArray []CopyrightResult
	if bz != nil {
		err := util.Json.Unmarshal(bz, &copyrightArray)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return copyrightArray, err
		}
	}
	return copyrightArray, nil
}

func (k Keeper) QueryAccountCopyrightVoteList(ctx sdk.Context, copyrightHash string) ([]AccountVoteInfor, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	copyrightVoteListInfor := make([]AccountVoteInfor, 0)
	copyrightVoteListByte := store.Get(copyrightVoteListFor + copyrightHash)
	if copyrightVoteListByte != nil {
		err := util.Json.Unmarshal(copyrightVoteListByte, &copyrightVoteListInfor)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return copyrightVoteListInfor, err
		}
	}
	return copyrightVoteListInfor, nil
}


func (k Keeper) UpdateAccountCopyrightVoteList(ctx sdk.Context, copyrightHash string, accountVoteList []AccountVoteInfor) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	err := store.Set(copyrightVoteListFor+copyrightHash, accountVoteList)
	if err != nil {
		log.WithError(err).Error("store.Set")
	}
	return err
}

func buildNewPicHash(ctx sdk.Context, picMap map[string]types.Link) string {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	node1 := unixfs.EmptyDirNode()

	var linkMapKeySlice []string
	//map key
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

			log.WithError(err).Error("cid2.Parse")
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
		log.WithError(err).Error("MarshalJSON")
		panic(err)
	}
	dagService := dagutils.NewMemoryDagService()
	err = dagService.Add(context.Background(), node1)
	if err != nil {
		log.WithError(err).Error("dagService.Add")
		panic(err)
	}
	getDagHash := func(editor *dagutils.Editor) string {
		nnode, err := editor.Finalize(context.Background(), dagService)
		if err != nil {
			log.WithError(err).Error("editor.Finalize")
			return ""
		}
		newHash := ipath.IpfsPath(nnode.Cid())
		return newHash.String()
	}
	e := dagutils.NewDagEditor(node1, dagService)
	newHash := getDagHash(e)
	log.WithField("hash", newHash).Debug("buildNewPicHash")
	return string(picByte)
}

//,..
func fileNameSuffix(fileName string) string {
	var fileSuffix string
	fileSuffix = path.Ext(fileName)
	filenameOnly := strings.TrimSuffix(fileName, fileSuffix)
	if filenameOnly == "" {
		filenameOnly = strings.Replace(fileName, ".", "", -1)
	}
	return filenameOnly
}

func getVoteAward(ctx sdk.Context, power decimal.Decimal) decimal.Decimal {
	return power.Mul(core.CopyrightVoteAwardRateV2)
}
