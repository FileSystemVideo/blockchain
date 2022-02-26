<!--
order: 2
-->

# State

## Minter

minter是保存当前通胀信息的空间。

 - Minter: `0x00 -> amino(minter)`

```go
type Minter struct {
	Inflation        sdk.Dec   // 当前年通货膨胀率
	AnnualProvisions sdk.Dec   // 当前年度扣除准备金
}
```

## Params

铸币参数保存在全局参数存储中。

 - Params: `mint/params -> amino(params)`

```go
type Params struct {
	MintDenom           string  // 铸币厂的硬币类型
	InflationRateChange sdk.Dec // 通货膨胀率的最大年变化
	InflationMax        sdk.Dec // 最大通货膨胀率
	InflationMin        sdk.Dec // 最小通货膨胀率
	GoalBonded          sdk.Dec // 绑定的主币的百分比
	BlocksPerYear       uint64   // 每年预计区块数
}
```
