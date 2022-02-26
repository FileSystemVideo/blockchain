<!--
order: 0
title: Capability Overview
parent:
  title: "capability"
-->

# `capability`

## Overview

`x/capability` is an implementation of a Cosmos SDK module, per [ADR 003](./../../../docs/architecture/adr-003-dynamic-capability-store.md),
that allows for provisioning, tracking, and authenticating multi-owner capabilities
at runtime.

守护者维持着两种状态：持久的和短暂的记忆。
持久存储维护一个全局唯一的自动递增索引，以及从功能索引到一组功能所有者（定义为模块和功能名称元组）的映射。
内存中的临时状态使用正向和反向索引跟踪实际功能，以本地内存中的地址表示。
正向索引将模块名称和功能元组映射到功能名称。
反向索引映射模块和功能名称以及功能本身。

keeper允许创建 `范围` 子keeper，这些子keeper通过名称绑定到特定模块。
必须在应用程序初始化时创建作用域保持器并将其传递给模块，然后模块可以使用它们声明它们接收的功能并检索它们按名称拥有的功能，此外还可以创建新功能并验证其他模块传递的功能。
作用域的keeper不能越出其作用域，因此模块不能干扰或检查其他模块拥有的功能。

keeper没有提供在其他模块（如queryer、REST和CLI处理程序以及genesis state）中可以找到的其他核心功能。

## Initialization

在应用程序初始化期间，必须使用持久存储密钥和内存存储密钥实例化keeper。


```go
type App struct {
  // ...

  capabilityKeeper *capability.Keeper
}

func NewApp(...) *App {
  // ...

  app.capabilityKeeper = capability.NewKeeper(codec, persistentStoreKey, memStoreKey)
}
```

创建keeper之后，可以使用它来创建作用域子keeper，这些子keeper被传递给可以创建、验证和声明功能的其他模块。
在创建了所有必要的作用域保持器并加载了状态之后，必须初始化和密封主功能保持器，以填充内存中的状态并防止创建更多的作用域保持器。

```go
func NewApp(...) *App {
  // ...

  // Initialize and seal the capability keeper so all persistent capabilities
  // are loaded in-memory and prevent any further modules from creating scoped
  // sub-keepers.
  ctx := app.BaseApp.NewContext(true, tmproto.Header{})
  app.capabilityKeeper.InitializeAndSeal(ctx)

  return app
}
```

## Contents

1. **[Concepts](01_concepts.md)**
1. **[State](02_state.md)**
