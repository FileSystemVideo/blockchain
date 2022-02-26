<!--
order: 4
-->

# BeginBlock

## Liveness Tracking

在每个块的开头，我们更新每个验证器的“ValidatorSigningInfo”，并检查它们是否通过滑动窗口越过了活跃度阈值以下。
此滑动窗口由“SignedBlocksWindow”定义，此窗口中的索引由验证器的“ValidatorSigningInfo”中的“IndexOffset”确定。
对于处理的每个块，“IndexOffset”都会递增，无论验证程序是否签名。
一旦确定了索引，就会相应地更新“MissedBlocksBitArray”和“MissedBlocksCounter”。

最后，为了确定验证器是否超过活跃度阈值，我们获取丢失的最大块数“maxMissed”，
即“maxMissed = SignedBlocksWindow-（MinSignedPerWindow*SignedBlocksWindow）”，以及可以确定活跃度的最小高度“minHeight”。
如果当前块大于“minHeight”，而验证器的“MissedBlocksCounter”大于“maxMissed”，则它们将被“SlashFractionDowntime”删除，将因“DowntimeJailDuration”而入狱，并重置以下值：“MissedBlocksBitArray”、“MissedBlocksCounter”和“IndexOffset”。


每当验证器: 未参与区块签名, 则 未签名的区块数 - 1
每当验证器: 参与了区块签名, 则 未签名额区块数 + 1

最大允许的未签名区块数 = SignedBlocksWindow - (SignedBlocksWindow * MinSignedPerWindow) 大约 50个
未签名的区块数 > 最大允许的未签名区块数 = 执行监禁和削减

**Note**: Liveness slashes do **NOT** lead to a tombstombing.

```go
height := block.Height

for vote in block.LastCommitInfo.Votes {
  signInfo := GetValidatorSigningInfo(vote.Validator.Address)

  // This is a relative index, so we counts blocks the validator SHOULD have
  // signed. We use the 0-value default signing info if not present, except for
  // start height.
  index := signInfo.IndexOffset % SignedBlocksWindow()
  signInfo.IndexOffset++

  // Update MissedBlocksBitArray and MissedBlocksCounter. The MissedBlocksCounter
  // just tracks the sum of MissedBlocksBitArray. That way we avoid needing to
  // read/write the whole array each time.
  missedPrevious := GetValidatorMissedBlockBitArray(vote.Validator.Address, index)
  missed := !signed

  switch {
  case !missedPrevious && missed:
    // array index has changed from not missed to missed, increment counter
    SetValidatorMissedBlockBitArray(vote.Validator.Address, index, true)
    signInfo.MissedBlocksCounter++  //如果本次检测到没有签名,则未签名计数器加1

  case missedPrevious && !missed:
    // array index has changed from missed to not missed, decrement counter
    SetValidatorMissedBlockBitArray(vote.Validator.Address, index, false)
    signInfo.MissedBlocksCounter--  //如果本次检测到已经签名，则未签名计数器减1(也就是说，只要验证器不是连续50个区块不在线，那么她的未签名区块数会慢慢恢复)

  default:
    // array index at this index has not changed; no need to update counter
  }

  if missed {
    // emit events...
  }

  minHeight := signInfo.StartHeight + SignedBlocksWindow()
  maxMissed := SignedBlocksWindow() - MinSignedPerWindow()

  // If we are past the minimum height and the validator has missed too many
  // jail and slash them.
  if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
    validator := ValidatorByConsAddr(vote.Validator.Address)

    // emit events...

    // We need to retrieve the stake distribution which signed the block, so we
    // subtract ValidatorUpdateDelay from the block height, and subtract an
    // additional 1 since this is the LastCommit.
    //
    // Note, that this CAN result in a negative "distributionHeight" up to
    // -ValidatorUpdateDelay-1, i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
    // That's fine since this is just used to filter unbonding delegations & redelegations.
    distributionHeight := height - sdk.ValidatorUpdateDelay - 1

    Slash(vote.Validator.Address, distributionHeight, vote.Validator.Power, SlashFractionDowntime())
    Jail(vote.Validator.Address)

    signInfo.JailedUntil = block.Time.Add(DowntimeJailDuration())

    // We need to reset the counter & array so that the validator won't be
    // immediately slashed for downtime upon rebonding.
    signInfo.MissedBlocksCounter = 0
    signInfo.IndexOffset = 0
    ClearValidatorMissedBlockBitArray(vote.Validator.Address)
  }

  SetValidatorSigningInfo(vote.Validator.Address, signInfo)
}
```
