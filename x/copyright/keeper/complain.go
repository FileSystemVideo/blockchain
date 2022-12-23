package keeper

import (
	"errors"
	"fmt"
	"fs.video/blockchain/core"
	"fs.video/blockchain/util"
	"fs.video/blockchain/x/copyright/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

//key
const (
	complainKeyForResult = "complain-for-result"
	voteName             = "_vote"
	voteShareName        = "_vote_share"
	complainHash         = "complain_hash" 
)

type CopyRightComplain struct {
	DataHash        string         `json:"datahash"`         //hash
	Author          string         `json:"author"`           
	Productor       string         `json:"productor"`        
	LegalNumber     string         `json:"legal_number"`     
	LegalTime       string         `json:"legal_time"`       
	ComplainInfor   string         `json:"complain_infor"`   
	ComplainAccount sdk.AccAddress `json:"complain_account"` 
	AccuseAccount   sdk.AccAddress `json:"accuse_account"`   
	ComplainId      string         `json:"complain_id"`      //id
	ComplainStatus  string         `json:"complain_status"`  //	 0  1 . 2   3 ,, 4 , 5 ,,.
	AccusedStatus   string         `json:"accused_status"`   
	AccusedInfo     string         `json:"accused_info"`     
	AccusedIp       string         `json:"accused_ip"`       //ip
	ComplainTime    int64          `json:"complain_time"`    
	ResponseTime    int64          `json:"response_time"`    
	//MortgAmount     sdk.Coins      `json:"mortg_amount"`     
}

type ComplainResult struct {
	ComplainIdList string `json:"complain_list"`
}


type ComplainVote struct {
	ComplainId   string        `json:"complain_id"`
	VoteInfor    string        `json:"voteinfor"`
	AccountsVote []AccountVote `json:"accounts_vote"`
}


type AccountVote struct {
	Account    string  `json:"account"`
	VoteStatus string  `json:"vote_status"`
	VoteShare  sdk.Dec `json:"vote_share"`
}


type ComplainShareInfor struct {
	ComplainId   string  `json:"complain_id"`
	FavorTotal   sdk.Dec `json:"favor_total"`
	AgainstTotal sdk.Dec `json:"against_total"`
}


func (keeper Keeper) ComplainVote(ctx sdk.Context, msg types.ComplainVoteData) error {
	
	complainInfor, err := keeper.GetCopyrightComplainInfor(ctx, msg.ComplainId)
	if err != nil {
		return err
	}
	
	if complainInfor.ComplainId == "" {
		return types.ErrComplainDoesNotExist
	}
	
	if complainInfor.ComplainStatus != "4" {
		return types.ErrComplainStatusInvalid
	}

	if msg.VoteStatus != "1" && msg.VoteStatus != "0" {
		return types.ErrComplainVoteStatusInvalid
	}

	
	dateTime := util.TimeStampToTime(complainInfor.ResponseTime)
	endTime := dateTime.Add(core.VoteResultTimePerioad)
	if endTime.Before(ctx.BlockTime()) { 
		return types.ErrComplainFinished
	}

	
	realShares, _ := keeper.GetAccountDelegatorShares(ctx, msg.VoteAccount)
	delegation := sdk.MustNewDecFromStr(realShares)
	if !delegation.IsPositive() {
		return types.ErrAccountHasNoVoteRight
	}

	
	
	err = keeper.SetComplainVoteInfor(ctx, msg.ComplainId, msg.VoteAccount.String(), msg.VoteStatus, msg.VoteShare)
	if err != nil {
		return err
	}
	err = keeper.SetComplainVoteShare(ctx, msg.ComplainId, msg.VoteStatus, msg.VoteShare)
	if err != nil {
		return err
	}
	return nil
}


func (keeper Keeper) ComplainResponse(ctx sdk.Context, msg types.ComplainResponseData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	//datahash 
	complainInfor, err := keeper.GetCopyrightComplainInfor(ctx, msg.ComplainId)
	if err != nil {
		return err
	}
	if complainInfor.ComplainId == "" {
		return types.ComplainIdInvalid
	}
	if complainInfor.AccusedStatus != "0" && complainInfor.ComplainStatus != "0" {
		return types.ErrComplainStatusInvalid
	}
	if msg.Status != "1" && msg.Status != "2" {
		return types.ErrResponseStatusInvalid
	}
	complainStatus := "0"
	if msg.Status == "1" { 
		complainInfor.ComplainStatus = "1"
		
		err = keeper.SetCopyrightAuthor(ctx, msg.DataHash, complainInfor.ComplainAccount)
		if err != nil {
			return err
		}
		//nft 
		err = keeper.NftTransfer(ctx, msg.DataHash, msg.AccuseAccount.String(), complainInfor.ComplainAccount.String())
		if err != nil {
			return err
		}
		
		copyright, err := keeper.GetCopyright(ctx, msg.DataHash)
		if err != nil {
			return err
		}
		var copyrightData types.CopyrightData
		err = util.Json.Unmarshal(copyright, &copyrightData)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
		if copyrightData.Size > 0 {
			
			flag := keeper.UpdateAccountSpace(ctx, msg.AccuseAccount, copyrightData.Size)
			if !flag {
				return types.AccountSpaceReturnError
			}
			
			flag = keeper.UpdateAccountSpaceUsed(ctx, complainInfor.ComplainAccount, copyrightData.Size)
			if !flag {
				return types.AccountSpaceNotEnough
			}
			
			err = keeper.LockAccountSpaceMinerReturn(ctx, complainInfor.ComplainAccount.String(), copyrightData.Size)
			if err != nil {
				return err
			}
		}
		complainStatus = "1"
		//,ids
		err = keeper.RemoveFromComplainIds(ctx, msg.ComplainId)
		if err != nil {
			return err
		}
		
		keeper.RemoveComplainHash(ctx, msg.DataHash)
	} else {
		complainStatus = "4" 
	}
	complainInfor.ComplainStatus = complainStatus
	complainInfor.AccusedStatus = msg.Status
	complainInfor.AccusedInfo = msg.AccuseInfor
	complainInfor.ResponseTime = msg.ResponseTime
	return keeper.SetCopyrightComplainInfor(ctx, complainInfor)
}


func (keeper Keeper) CopyrightComplain(ctx sdk.Context, msg types.CopyrightComplainData) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	//datahash 
	if !keeper.IsDataHashPresent(ctx, msg.DataHash) {
		return types.ErrDataHashDoesNotExist
	}

	
	copyrightInforByte, err := keeper.GetCopyright(ctx, msg.DataHash)
	if err != nil {
		return err
	}
	var copyrightData types.CopyrightData
	err = util.Json.Unmarshal(copyrightInforByte, &copyrightData)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}

	
	if !copyrightData.Creator.Equals(msg.AccuseAccount) {
		return types.AccuseAccountInvalid
	}

	
	spaceInfo := keeper.QueryAccountSpaceMinerInfor(ctx, msg.ComplainAccount.String())

	//  true  ,  false 
	var checkSpaceRemain = func(copyInfoSize int64, spaceInfo AccountSpaceMiner) bool {
		spaceRemain := spaceInfo.SpaceTotal.Sub(spaceInfo.UsedSpace).Sub(spaceInfo.LockedSpace)
		copySize := decimal.NewFromInt(copyInfoSize)
		
		return copySize.LessThan(spaceRemain)
	}

	
	if !checkSpaceRemain(copyrightData.Size, spaceInfo) {
		return types.AccountSpaceNotEnough
	}
	
	err = keeper.CreateCopyrightComplain(ctx, msg)
	if err != nil {
		log.WithError(err).Error("CreateCopyrightComplain")
	}
	keeper.SetComplainHash(ctx, msg.DataHash)
	
	return keeper.LockAccountSpaceMiner(ctx, msg.ComplainAccount.String(), copyrightData.Size)
}


func (k Keeper) SetComplainHash(ctx sdk.Context, dataHash string) {
	store := k.KVHelper(ctx)
	store.Set(complainHash+dataHash, []byte("exit"))
}


func (k Keeper) ComplainHashStatus(ctx sdk.Context, dataHash string) bool {
	store := k.KVHelper(ctx)
	return store.Has(complainHash + dataHash)
}


func (k Keeper) RemoveComplainHash(ctx sdk.Context, dataHash string) {
	store := k.KVHelper(ctx)
	store.Delete(complainHash + dataHash)
}


func (k Keeper) CreateCopyrightComplain(ctx sdk.Context, msg types.CopyrightComplainData) error {
	copyrightComplain := CopyRightComplain{}
	copyrightComplain.DataHash = msg.DataHash
	copyrightComplain.Author = msg.Author
	copyrightComplain.Productor = msg.Productor
	copyrightComplain.LegalNumber = msg.LegalNumber
	copyrightComplain.LegalTime = msg.LegalTime
	copyrightComplain.ComplainInfor = msg.ComplainInfor
	copyrightComplain.ComplainId = msg.ComplainId
	copyrightComplain.ComplainAccount = msg.ComplainAccount
	copyrightComplain.AccuseAccount = msg.AccuseAccount
	copyrightComplain.AccusedStatus = "0"
	copyrightComplain.ComplainStatus = "0"
	copyrightComplain.ComplainTime = msg.ComplainTime
	//copyrightComplain.MortgAmount = coins
	store := k.KVHelper(ctx)
	store.Set(msg.ComplainId, copyrightComplain)
	
	return k.SetComplainKeyForResult(store, msg.ComplainId)
}

//id
func (k Keeper) SetComplainKeyForResult(store storeHelper, complainId string) error {
	complainIdListByte := store.Get(complainKeyForResult)
	var complainIds ComplainResult
	if complainIdListByte != nil {
		err := util.Json.Unmarshal(complainIdListByte, &complainIds)
		if err != nil {
			return err
		}
	}
	if complainIds.ComplainIdList == "" {
		complainIds.ComplainIdList = complainId
	} else {
		complainIds.ComplainIdList = complainIds.ComplainIdList + "|" + complainId
	}
	return store.Set(complainKeyForResult, complainIds)
}


func (k Keeper) RemoveFromComplainIds(ctx sdk.Context, complainId string) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(complainKeyForResult)
	var complainids ComplainResult
	err := util.Json.Unmarshal(bz, &complainids)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
		return err
	}
	complainIdsString := complainids.ComplainIdList
	if complainIdsString != "" && strings.Contains(complainIdsString, complainId) {
		idsStringArray := strings.Split(complainIdsString, "|")
		index := 0
		for i := 0; i < len(idsStringArray); i++ {
			if idsStringArray[i] == complainId {
				index = i
				break
			}
		}
		idsStringArray = append(idsStringArray[:index], idsStringArray[index+1:]...)
		complainIdsString = strings.Join(idsStringArray, "|")
		if complainIdsString == "" {
			
			store.Delete(complainKeyForResult)
		} else {
			complainids.ComplainIdList = complainIdsString
			err = store.Set(complainKeyForResult, complainids)
			if err != nil {
				log.WithError(err).Error("store.Set")
				return err
			}
		}
	}
	return nil
}

//id
func (k Keeper) GetComplainKeyForResult(ctx sdk.Context) (ComplainResult, error) {
	store := k.KVHelper(ctx)
	bz := store.Get(complainKeyForResult)
	var complainids ComplainResult
	if bz != nil {
		err := util.Json.Unmarshal(bz, &complainids)
		if err != nil {
			return complainids, err
		}
	}
	return complainids, nil
}


func (k Keeper) GetCopyrightComplainInfor(ctx sdk.Context, complainId string) (CopyRightComplain, error) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	bz := store.Get(complainId)
	if bz == nil {
		return CopyRightComplain{}, nil
	}
	var copyrightComplainInfor CopyRightComplain
	err := util.Json.Unmarshal(bz, &copyrightComplainInfor)
	if err != nil {
		log.WithError(err).Error("Unmarshal")
	}
	return copyrightComplainInfor, err
}


func (k Keeper) SetCopyrightComplainInfor(ctx sdk.Context, copyrightComplain CopyRightComplain) error {
	store := k.KVHelper(ctx)
	return store.Set(copyrightComplain.ComplainId, copyrightComplain)
}

//  - 6+
func (k Keeper) GetVoteInfor(ctx sdk.Context, complain string) ComplainVote {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	voteKey := complain + voteName
	if !store.Has(voteKey) {
		return ComplainVote{}
	}
	bz := store.Get(voteKey)
	var voteInfor ComplainVote
	if bz != nil {
		err := util.Json.Unmarshal(bz, &voteInfor)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return voteInfor
		}
	}
	return voteInfor
}


func (k Keeper) SetComplainVoteInfor(ctx sdk.Context, complainId, account, status string, voteShare sdk.Dec) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	complainVoteKey := complainId + voteName
	complainVoteByte := store.Get(complainVoteKey)
	var complainVote ComplainVote
	accountVote := AccountVote{
		Account:    account,
		VoteStatus: status,
		VoteShare:  voteShare,
	}
	if complainVoteByte == nil { 
		accountVotes := []AccountVote{}
		accountVotes = append(accountVotes, accountVote)
		complainVote = ComplainVote{
			ComplainId:   complainId,
			AccountsVote: accountVotes,
		}
	} else {
		err := util.Json.Unmarshal(complainVoteByte, &complainVote)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
		if len(complainVote.AccountsVote) > 0 {
			for i := 0; i < len(complainVote.AccountsVote); i++ {
				if complainVote.AccountsVote[i].Account == account {
					return types.HasVoted
				}
			}
		}
		complainVote.AccountsVote = append(complainVote.AccountsVote, accountVote)
	}
	return store.Set(complainVoteKey, complainVote)
}


func (k Keeper) QueryComplainVoteShare(ctx sdk.Context, complainId string) ComplainShareInfor {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	complainVoteShareKey := complainId + voteShareName
	complainVoteByte := store.Get(complainVoteShareKey)
	var complainShareInfor ComplainShareInfor
	if complainVoteByte != nil {
		err := util.Json.Unmarshal(complainVoteByte, &complainShareInfor)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
		}
	}
	return complainShareInfor
}


func (k Keeper) SetComplainVoteShare(ctx sdk.Context, complainId, status string, voteShare sdk.Dec) error {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	store := k.KVHelper(ctx)
	complainVoteShareKey := complainId + voteShareName
	complainVoteByte := store.Get(complainVoteShareKey)
	var complainShareVote ComplainShareInfor
	if complainVoteByte == nil { 
		complainShareVote = ComplainShareInfor{
			ComplainId: complainId,
		}
		if status == "0" {
			complainShareVote.FavorTotal = voteShare
			complainShareVote.AgainstTotal = sdk.NewDec(int64(0))
		} else {
			complainShareVote.AgainstTotal = voteShare
			complainShareVote.FavorTotal = sdk.NewDec(int64(0))
		}
	} else {
		err := util.Json.Unmarshal(complainVoteByte, &complainShareVote)
		if err != nil {
			log.WithError(err).Error("Unmarshal")
			return err
		}
		if status == "0" {
			complainShareVote.FavorTotal = complainShareVote.FavorTotal.Add(voteShare)
		} else {
			complainShareVote.AgainstTotal = complainShareVote.AgainstTotal.Add(voteShare)
		}
	}
	return store.Set(complainVoteShareKey, complainShareVote)
}


func (k Keeper) ComplainResult(ctx sdk.Context) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	//ids
	complainResult, err := k.GetComplainKeyForResult(ctx)
	if err != nil {
		log.WithError(err).Error("GetComplainKeyForResult")
		return
	}
	if complainResult.ComplainIdList != "" { 
		complainIds := strings.Split(complainResult.ComplainIdList, "|")
		timeNow := ctx.BlockHeader().Time
		for i := 0; i < len(complainIds); i++ {
			complainId := complainIds[i]
			voteStatus, complainStatus, voteAmountInfor, complainHash, complainAddr, accuseAddr := k.CalculateComplainVote(ctx, complainId, timeNow)
			if voteStatus {
				log.WithFields(logrus.Fields{
					"voteStatus":      voteStatus,
					"complainStatus":  complainStatus,
					"voteAmountInfor": voteAmountInfor,
				}).Info("vote result")
				copyrightMove := types.CopyrightMove{}
				//timeoutStatus := "0"
				copyrightByte, err := k.GetCopyright(ctx, complainHash)
				var copyrightData types.CopyrightData
				err = util.Json.Unmarshal(copyrightByte, &copyrightData)
				if err != nil {
					return
				}
				
				if complainStatus == "2" || complainStatus == "5" { 
					
					err = k.NftTransfer(ctx, complainHash, accuseAddr, complainAddr)
					complainAddress, err := sdk.AccAddressFromBech32(complainAddr)
					if err != nil {
						log.WithError(err).WithField("Account", complainAddr).Error("AccAddressFromBech32")
					}
					copyrightData.Creator = complainAddress
					err = k.UpdateCopyrightCreator(ctx, copyrightData)
					if err != nil {
						return
					}
					accuseAddress, err := sdk.AccAddressFromBech32(accuseAddr)
					if copyrightData.Size > 0 {
						flag := k.UpdateAccountSpace(ctx, accuseAddress, copyrightData.Size)
						if !flag {
							errorString := fmt.Sprintf("The copyright claim was not replied, the copyright was transferred, and the storage space of the defendant failed to be returned account:%s copyright size:%d ", accuseAddress.String(), copyrightData.Size)
							log.WithError(errors.New(errorString)).Error("UpdateAccountSpace")
							panic(errorString)
						}
						
						flag = k.UpdateAccountSpaceUsed(ctx, complainAddress, copyrightData.Size)
						if !flag {
							errorString := fmt.Sprintf("The copyright complaint is not replied, the copyright is transferred, and the storage space of the copyright complainant is insufficient account:%s copyright size:%d ", complainAddress.String(), copyrightData.Size)
							log.WithError(errors.New(errorString)).Error("UpdateAccountSpaceUsed")
							panic(errorString)
						}
					}
					copyrightMove.FromAccount = accuseAddr
					copyrightMove.ToAccount = complainAddr
					copyrightMove.DataHash = complainHash
					copyrightMove.Status = 1
				}

				
				err = k.LockAccountSpaceMinerReturn(ctx, complainAddr, copyrightData.Size)
				if err != nil {
					log.WithError(err).Error("LockAccountSpaceMinerReturn")
					panic(err)
				}
				copyrightMove.ComplainStatus = complainStatus
				copyrightMove.ComplainId = complainId
				copyrightMove.TimeStemp = ctx.BlockTime().Unix()
				copyrightMove.BlockNum = ctx.BlockHeight()
				err = k.AddBlockRDS(ctx, types.NewBlockRD(copyrightMove))
				if err != nil {
					log.WithError(err).Error("AddBlockRDS-copyrightMove")
				}
				
				log.Info("Voting weight", voteAmountInfor)
				
				k.RemoveComplainHash(ctx, complainHash)
			}
		}
	}
}


func (k Keeper) CalculateComplainVote(ctx sdk.Context, complainId string, timeNow time.Time) (voteStatus bool, complainStatus, voteAmountInfor, complainHash, complainAccount, accuseAccount string) {
	log := core.BuildLog(core.GetFuncName(), core.LmChainKeeper)
	voteAmountInfor = "" 
	voteStatus = false   
	complainStatus = "3" 
	
	complainInfor, err := k.GetCopyrightComplainInfor(ctx, complainId)
	if err != nil {
		log.WithError(err).Error("GetCopyrightComplainInfor")
	}
	
	var updateTimeInt64 int64
	if complainInfor.AccusedStatus == "0" {
		updateTimeInt64 = complainInfor.ComplainTime
	} else if complainInfor.AccusedStatus == "2" {
		updateTimeInt64 = complainInfor.ResponseTime
	}

	dateTime := util.TimeStampToTime(updateTimeInt64)
	endTime := dateTime.Add(core.VoteResultTimePerioad)
	if endTime.After(timeNow) { 
		return voteStatus, complainStatus, voteAmountInfor, complainInfor.DataHash, complainInfor.ComplainAccount.String(), complainInfor.AccuseAccount.String()
	}

	if complainInfor.AccusedStatus == "0" { 
		
		complainStatus = "2" 
		err = k.SetCopyrightAuthor(ctx, complainInfor.DataHash, complainInfor.ComplainAccount)
		if err != nil {
			log.WithError(err).Error("SetCopyrightAuthor")
			return
		}
		voteStatus = true
		complainInfor.ComplainStatus = complainStatus
		err = k.SetCopyrightComplainInfor(ctx, complainInfor)
		if err != nil {
			log.WithError(err).Error("SetCopyrightComplainInfor")
			return
		}

	} else if complainInfor.AccusedStatus == "2" && complainInfor.ComplainStatus == "4" { 
		
		complainVoteInfor := k.GetVoteInfor(ctx, complainId)
		if complainVoteInfor.ComplainId != "" { 
			
			complainVoteShare := k.QueryComplainVoteShare(ctx, complainId)
			log.WithFields(logrus.Fields{
				"Support number":    complainVoteShare.FavorTotal,
				"Inverse logarithm": complainVoteShare.AgainstTotal,
				"Voting proportion": voteAmountInfor,
			}).Info("CalculateComplainVote")
			if complainVoteShare.ComplainId != "" && complainVoteShare.FavorTotal.GT(complainVoteShare.AgainstTotal) {
				complainInfor.ComplainStatus = "5"
				complainStatus = "5"
			} else {
				complainInfor.ComplainStatus = "3" 
				complainStatus = "3"
			}
			voteStatus = true
			
			err = k.SetCopyrightComplainInfor(ctx, complainInfor)
			if err != nil {
				log.WithError(err).Error("SetCopyrightComplainInfor")
				return
			}
		} else { 
			voteStatus = true
			complainInfor.ComplainStatus = "3"
			err = k.SetCopyrightComplainInfor(ctx, complainInfor)
			if err != nil {
				log.WithError(err).Error("SetCopyrightComplainInfor")
				return
			}
		}
	}
	err = k.RemoveFromComplainIds(ctx, complainId)
	if err != nil {
		return
	}
	return voteStatus, complainStatus, voteAmountInfor, complainInfor.DataHash, complainInfor.ComplainAccount.String(), complainInfor.AccuseAccount.String()
}
