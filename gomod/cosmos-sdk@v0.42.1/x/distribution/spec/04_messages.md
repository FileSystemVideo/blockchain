<!--
order: 4
-->

# Messages

## MsgSetWithdrawAddress

默认情况下，取款地址是委托人地址。
如果委托人要更改其取款地址，则必须发送`MsgSetWithdrawAddress`。

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/distribution/v1beta1/tx.proto#L29-L37

```go

func (k Keeper) SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) error 
	if k.blockedAddrs[withdrawAddr.String()] {
		fail with "`{withdrawAddr}` is not allowed to receive external funds"
	}

	if !k.GetWithdrawAddrEnabled(ctx) {
		fail with `ErrSetWithdrawAddrDisabled`
	}

	k.SetDelegatorWithdrawAddr(ctx, delegatorAddr, withdrawAddr)
```

## MsgWithdrawDelegatorReward

在特殊情况下，委托人可能希望仅从一个验证器中提取奖励。

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0/proto/cosmos/distribution/v1beta1/tx.proto#L42-L50

```go
// withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	val := k.stakingKeeper.Validator(ctx, valAddr)
	if val == nil {
		return nil, types.ErrNoValidatorDistInfo
	}

	del := k.stakingKeeper.Delegation(ctx, delAddr, valAddr)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, val, del)
	if err != nil {
		return nil, err
	}

	// reinitialize the delegation
	k.initializeDelegation(ctx, valAddr, delAddr)
	return rewards, nil
}
```

## 提取验证器所有奖励

当验证器想要取回他们的奖励时，它必须发送一个 `MsgWithdrawDelegatorReward` 数组。
请注意，此事务逻辑的各个部分也会随着单个委托的任何更改而触发，例如取消绑定、重新委托或将附加令牌委托给特定验证器。
此交易提取验证人的佣金，以及他们自授权获得的任何奖励。

```go

for _, valAddr := range validators {
    val, err := sdk.ValAddressFromBech32(valAddr)
    if err != nil {
        return err
    }

    msg := types.NewMsgWithdrawDelegatorReward(delAddr, val)
    if err := msg.ValidateBasic(); err != nil {
        return err
    }
    msgs = append(msgs, msg)
}
```

## Common calculations 

### 更新验证器累计总数

为了确定验证器在特定块上有权获得的token池的数量，必须计算验证器累计的总量。
累加器总是与现有累加器相加。
每次从系统中提取奖励时，都会更新此它。

```go
func (g FeePool) UpdateTotalValAccum(height int64, totalBondedTokens Dec) FeePool
    blocks = height - g.TotalValAccumUpdateHeight
    g.TotalValAccum += totalDelShares * blocks
    g.TotalValAccumUpdateHeight = height
    return g
```

### 更新验证器的累计值


必须更新委托人累计总数，以确定每个委托人相对于该验证器的其他委托人有权获得的池令牌的数量。

累加器总是与现有累加器相加。每次从验证器取款时都要更新它。

``` go
func (vi ValidatorDistInfo) UpdateTotalDelAccum(height int64, totalDelShares Dec) ValidatorDistInfo
    blocks = height - vi.TotalDelAccumUpdateHeight
    vi.TotalDelAccum += totalDelShares * blocks
    vi.TotalDelAccumUpdateHeight = height
    return vi
```

### 手续费池 到 验证器池

每当验证器或委托人执行取回,或验证器是提议者并接收新token时，相关验证器必须将token从被动全局池移动到自己的池中。
正是在这一点上，委托会被取回

```go
func (vi ValidatorDistInfo) TakeFeePoolRewards(g FeePool, height int64, totalBonded, vdTokens, commissionRate Dec) (
                                vi ValidatorDistInfo, g FeePool)

    g.UpdateTotalValAccum(height, totalBondedShares)
    
    // update the validators pool
    blocks = height - vi.FeePoolWithdrawalHeight
    vi.FeePoolWithdrawalHeight = height
    accum = blocks * vdTokens
    withdrawalTokens := g.Pool * accum / g.TotalValAccum 
    commission := withdrawalTokens * commissionRate
    
    g.TotalValAccum -= accumm
    vi.PoolCommission += commission
    vi.PoolCommissionFree += withdrawalTokens - commission
    g.Pool -= withdrawalTokens

    return vi, g
```


### 委托奖励提取

对于委托（包括验证者的自委托），所有奖励池中的奖励已经被验证者的佣金拿走。

```go
func (di DelegationDistInfo) WithdrawRewards(g FeePool, vi ValidatorDistInfo,
    height int64, totalBonded, vdTokens, totalDelShares, commissionRate Dec) (
    di DelegationDistInfo, g FeePool, withdrawn DecCoins)

    vi.UpdateTotalDelAccum(height, totalDelShares) 
    g = vi.TakeFeePoolRewards(g, height, totalBonded, vdTokens, commissionRate) 
    
    blocks = height - di.WithdrawalHeight
    di.WithdrawalHeight = height
    accum = delegatorShares * blocks 
     
    withdrawalTokens := vi.Pool * accum / vi.TotalDelAccum
    vi.TotalDelAccum -= accum

    vi.Pool -= withdrawalTokens
    vi.TotalDelAccum -= accum
    return di, g, withdrawalTokens

```

### Validator commission withdrawal

每次奖励进入验证器时都会计算佣金。

```go
func (vi ValidatorDistInfo) WithdrawCommission(g FeePool, height int64, 
          totalBonded, vdTokens, commissionRate Dec) (
          vi ValidatorDistInfo, g FeePool, withdrawn DecCoins)

    g = vi.TakeFeePoolRewards(g, height, totalBonded, vdTokens, commissionRate) 
    
    withdrawalTokens := vi.PoolCommission 
    vi.PoolCommission = 0

    return vi, g, withdrawalTokens
```
