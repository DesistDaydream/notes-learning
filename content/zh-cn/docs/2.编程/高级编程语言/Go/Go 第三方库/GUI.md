---
title: GUI
---

# 概述

> 参考：

# Go OpenCV

> 参考：
> 
> - <http://www.codebaoku.com/it-go/it-go-146781.html>

<https://github.com/go-opencv/go-opencv>
<https://github.com/hybridgroup/gocv>

# Fyne

> 参考：
> 
> - [GitHub 项目，fyne-io/fyne](https://github.com/fyne-io/fyne)
> - [官网](https://fyne.io/)
> - [简书，go fyne 开发桌面应用](https://www.jianshu.com/p/be97c0668252)
> - [稀土掘进，专栏-Fyne ( go跨平台GUI )中文文档](https://juejin.cn/column/7087843642252984351)

Fyne 是一个易于学习、免费、开源的工具包，用于构建适用于桌面、移动设备及其他设备的图形应用程序。

注意：使用 Fyne 需要安装 [MinGW-w64](https://sourceforge.net/projects/mingw-w64/)

Fyne 用起来有点像写前端代码。使用时先创建一个 APP，就相当于创建了一个 HTML 页面，然后围绕这个 APP 编写代码。我们可以向这个 APP 中添加各种符合某种 **Layout(布局)** 的 **Widget(小组件)**，Widget 有点类似于 HTML 中的各种元素（比如表 单、输入框、按钮、等等），Layout 就是类似 CSS 一样的样式了。

这其中的 APP 概念，又有点像 Vue 的 APP 概念。不会单指一个 HTML 页面。

## Hello World

```go
package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 实例化一个应用
	a := app.New()
	// 为应用创建一个窗口
	w := a.NewWindow("Hello")

	// ######## 创建一些应该在窗口中显示的内容，以及设计窗口中的布局 ########
	// 创建一个 Label 小部件
	labelWidget := widget.NewLabel("Hello Fyne!")
	// 创建一个按钮小部件
	buttonWidget := widget.NewButton("Hi!", func() {
		labelWidget.SetText("Welcome :)")
	})

	// 创建布局并将指定的对象（小部件、等等）放到这个布局中。
	// NewVBox 中会使用 fyne 中自带的 VBox 布局，这种布局会将对象从上到下堆叠。
	layout := container.NewVBox(
		labelWidget,
		buttonWidget,
	)
	// #################################################################

	// 为 w 窗口设置应该在其中的内容
	w.SetContent(layout)

	// 显示窗口并运行程序。必须要在 main() 函数的末尾，因为该方法将会阻塞。
	w.ShowAndRun()
}
```

除了 Hello World 示例，官方还提供了一个演示程序，我们可以直接执行：

```bash
go run fyne.io/fyne/v2/cmd/fyne_demo@latest
```

在这个窗口中，我们可以看到 Fyne 的所有功能，就像下面这样：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/go/202306271643980.png)

# 应用案例

[公众号-Go语言中文网，使用 Go GUI 库 fyne 编写一个计算器程序](https://mp.weixin.qq.com/s/VrTFMhpYvzr78ULqsQ15Sw)

