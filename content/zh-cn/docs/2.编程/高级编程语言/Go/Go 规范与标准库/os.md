---
title: os
linkTitle: os
date: 2024-03-12T15:29
weight: 20
---

# 概述

> 参考：
> 
> - [Go 标准库，os](https://pkg.go.dev/os)

os 包提供了 [Operating system](docs/1.操作系统/Operating%20system/Operating%20system.md) 功能的接口，不受不同平台的影响。

这是一个简单的示例，打开一个文件并读取其中的一些内容

```go
file, err := os.Open("file.go") // For read access.
if err != nil {
	log.Fatal(err)
}
```