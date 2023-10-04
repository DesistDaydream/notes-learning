---
title: Deno 工具
linkTitle: Deno 工具
date: 2024-01-15T22:01
weight: 20
---

# 概述

> 参考：
> 
> - [官方文档，运行时-手册-CLI 命令参考](https://docs.deno.com/runtime/manual/tools/)

## compile

https://docs.deno.com/runtime/manual/tools/compiler

把脚本编译成独立的可执行文件。


## info

https://docs.deno.com/runtime/manual/tools/dependency_inspector

可用于显示有关缓存位置的信息，若指定了具体模块则可检查 ES 模块及其所有依赖项。

可以显示下面这些缓存信息

```text
DENO_DIR location: C:\Users\DesistDaydream\AppData\Local\deno
Remote modules cache: C:\Users\DesistDaydream\AppData\Local\deno\deps
npm modules cache: C:\Users\DesistDaydream\AppData\Local\deno\npm
Emitted modules cache: C:\Users\DesistDaydream\AppData\Local\deno\gen
Language server registries cache: C:\Users\DesistDaydream\AppData\Local\deno\registries
Origin storage: C:\Users\DesistDaydream\AppData\Local\deno\location_data
```

# 其他工具
