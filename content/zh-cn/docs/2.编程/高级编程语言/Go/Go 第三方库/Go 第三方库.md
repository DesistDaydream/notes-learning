---
title: Go 第三方库
---

# 概述

> 参考：

第三方库一般属于由个人开发，实现更多丰富功能的库。在 [Go.dev ](https://pkg.go.dev/)可以搜索自己想要使用的所有库。
# 日志

## logrus

> 参考：
> 
> - [GitHub 项目，sirupsen/logrus](https://github.com/sirupsen/logrus)
> - <https://pkg.go.dev/github.com/sirupsen/logrus>

Logrus 是一种结构化得用于 Go 语言的日志处理器，完全与 Go 标准库中的 log 库。这名字来源于吉祥物 Walrus(海象)，所以在官方文档中，所有示例都与 Walrus 相关。

```go
package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	// Logrus 共有七个日志级别，由高到底分别为：Trace、Debug、Info、Warning、Error、Fatal、Panic
	// 默认情况下，只有 Info 及以下级别可以正常输出。如果想要输出高级别日志，通过 SetLevel() 函数设置日志级别即可
	// SetLevel() 函数的实参可以通过 ParseLevel() 函数将字符串解析为对应级别
	// logrus.SetLevel(logrus.InfoLevel)

	// 输出 Info 级别的日志内容
	logrus.Info("Hello World")
}

// 输出内容如下：
// time="2021-09-20T11:58:36+08:00" level=info msg="Hello World"
```

# 文件处理

## Excel 文件处理

### Excelize

> 参考：
> 
> - [GitHub 项目，xuri/excelize](https://github.com/xuri/excelize)
> - [官方文档](https://xuri.me/excelize/zh-hans/)
