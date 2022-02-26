<!--
order: 3
-->

# Begin-Block

铸币参数重新计算和通货膨胀支付在每个块的开始。

## NextInflationRate

目标年通货膨胀率在每个区块重新计算。
通货膨胀率也会发生变化（正或负），这取决于与理想比率（67%）的距离。
可能的最大利率变化被定义为每年13%，但是年通货膨胀率被限制在7%到20%之间。

```
NextInflationRate(params Params, bondedRatio sdk.Dec) (inflation sdk.Dec) {
	inflationRateChangePerYear = (1 - bondedRatio/params.GoalBonded) * params.InflationRateChange
	inflationRateChange = inflationRateChangePerYear/blocksPerYr

	// increase the new annual inflation for this next cycle
	inflation += inflationRateChange
	if inflation > params.InflationMax {
		inflation = params.InflationMax
	}
	if inflation < params.InflationMin {
		inflation = params.InflationMin
	}

	return inflation
}
```

## NextAnnualProvisions


根据当前总供给和通货膨胀率计算年度准备金。
每个块计算一次此参数。

```
NextAnnualProvisions(params Params, totalSupply sdk.Dec) (provisions sdk.Dec) {
	return Inflation * totalSupply
```

## BlockProvision

根据当前年度准备金计算每个区块产生的准备金。
然后由 `mint` 模块的 `ModuleMinterAccount` 生成年度准备金，然后将其传输到 `auth` 的 `FeeCollector` 模块的 `ModuleAccount`。

```
BlockProvision(params Params) sdk.Coin {
	provisionAmt = AnnualProvisions/ params.BlocksPerYear
	return sdk.NewCoin(params.MintDenom, provisionAmt.Truncate())
```
