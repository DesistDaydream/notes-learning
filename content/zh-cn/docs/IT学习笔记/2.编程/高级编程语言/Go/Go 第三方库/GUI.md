---
title: GUI
---

# 概述

> 参考：

# Go OpenCV

> 参考：
> - <http://www.codebaoku.com/it-go/it-go-146781.html>

<https://github.com/go-opencv/go-opencv>
<https://github.com/hybridgroup/gocv>

# Fyne

> 参考：
> - [GitHub 项目，fyne-io/fyne](https://github.com/fyne-io/fyne)
> - [官网](https://fyne.io/)
> - [简书，go fyne 开发桌面应用](https://www.jianshu.com/p/be97c0668252)
> - [稀土掘金，Fyne（go 跨平台 GUI）中文文档-小部件（五）](https://juejin.cn/post/7091103604492206087)

Fyne 是一个易于学习、免费、开源的工具包，用于构建适用于桌面、移动设备及其他设备的图形应用程序。

## Hello World

```go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
}
```
