<!--
order: 0
title: "Auth Overview"
parent:
  title: "auth"
-->

# `auth`

## Abstract

本文档指定了cosmossdk的auth模块。

auth模块负责指定应用程序的基本事务和帐户类型，因为SDK本身对这些细节是不可知的。

它包含ante处理程序，在这里执行所有基本事务有效性检查（签名、nonce、辅助字段），并公开account keeper，它允许其他模块读取、写入和修改帐户。


## Contents

1. **[Concepts](01_concepts.md)**
   - [Gas & Fees](01_concepts.md#gas-&-fees)
2. **[State](02_state.md)**
   - [Accounts](02_state.md#accounts)
3. **[AnteHandlers](03_antehandlers.md)**
   - [Handlers](03_antehandlers.md#handlers)
4. **[Keepers](04_keepers.md)**
   - [Account Keeper](04_keepers.md#account-keeper)
5. **[Vesting](05_vesting.md)**
   - [Intro and Requirements](05_vesting.md#intro-and-requirements)
   - [Vesting Account Types](05_vesting.md#vesting-account-types)
   - [Vesting Account Specification](05_vesting.md#vesting-account-specification)
   - [Keepers & Handlers](05_vesting.md#keepers-&-handlers)
   - [Genesis Initialization](05_vesting.md#genesis-initialization)
   - [Examples](05_vesting.md#examples)
   - [Glossary](05_vesting.md#glossary)
6. **[Parameters](07_params.md)**
