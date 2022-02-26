<!--
order: 1
-->

# State

## ConstantFee

由于预计需要大量的天然气成本来验证一个不变量（并且可能超过最大允许的区块天然气限制），因此使用固定费用代替标准天然气消耗方法。
固定费用旨在大于使用标准气体消耗法运行不变量的预期气体成本。

ConstantFee参数保存在全局参数存储中。

 - Params: `mint/params -> legacy_amino(sdk.Coin)`

