---
title: "Fyne"
linkTitle: "Fyne"
date: "2023-06-28T22:42"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，fyne-io/fyne](https://github.com/fyne-io/fyne)
> - [官网](https://fyne.io/)
> - [简书，go fyne 开发桌面应用](https://www.jianshu.com/p/be97c0668252)
> - [稀土掘进，专栏-Fyne ( go跨平台GUI )中文文档](https://juejin.cn/column/7087843642252984351)

Fyne 是一个易于学习、免费、开源的工具包，用于构建适用于桌面、移动设备及其他设备的图形应用程序。

注意：使用 Fyne 需要安装 [MinGW-w64](https://sourceforge.net/projects/mingw-w64/)

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

# Fyne 的组成

Fyne 创建的 GUI 应用程序由一个或多个 **CanvasOjbects(画布对象)**

Fyne 用起来有点像写前端代码

- 使用时先创建一个 App（代码中是 App 接口），这个 App 就像创建了一个浏览器；
- 然后创建一个 Window（代码中是 Window 接口），这就像就为浏览器打开了一个标签；
- 在这个 Window 中创建 **Canvas(画布)**（代码中是 Canvas 接口），这就像标签页中打开了一个 HTML 页面；
- 在 Canvas 中我们可以填充各种内容，Fyne 将内容抽象为 **CanvasObject(画布对象)**（代码中是 CanvasObject 接口），这就像 HTML 中的元素；
  - 而 **Widget(小组件)** 是一种特殊类型的 CanvasObject，它具有与之关联的交互元素。HTML 中常见的元素（比如 表单、输入框、按钮、等等），我们都可以找到对应的 Widget
    - 各种 Widget 都是一个实现了 CanvasObject 接口的结构体，比如 Form、Entry、Button 等等，我们可以在[小部件列表](#小部件列表)找看到所有的小部件，这些小部件相关的代码都在 [widget 包](https://pkg.go.dev/fyne.io/fyne/v2/widget)中
  - 在一个前端页面中，除了元素以外还应该包含样式，对于 Fyne 来说，我们通过 **Container(容器)** 和 **Layouts(布局**) 实现，这就像 CSS
- **Container(容器)** 也是一种 CanvasObject，只不过 Container 中可以使用 Layouts，并包含其他 CanvasObject（甚至可以包含其他 Container），这个概念有点像 Vue 的组件。
  - 我们可以在 Container 中设计 **Layout(布局)** 以安放各种 CanvasObject 的位置，在[布局列表](#布局列表)中可以看到 Fyne 自带的所有 Layout 样式。

## 小部件列表

> 参考：
> 
> - [官方文档，探索 Fyne-小部件列表](https://developer.fyne.io/explore/widgets)

## 布局列表

> 参考：
> 
> - [官方文档，探索 Fyne-布局列表](https://developer.fyne.io/explore/layouts)

# 数据绑定

Fyne v2.0.0 中引入了数据绑定，使将许多小部件连接到将随时间更新的数据源变得更加容易。该`data/binding`包有许多有用的绑定，可以管理将在应用程序中使用的大多数标准类型。可以使用绑定 API（例如`NewString`）来管理数据绑定，也可以将其连接到外部数据项，例如 (`BindInt(*int))。

支持绑定的小部件通常有一个`...WithData`构造函数来在创建小部件时设置绑定。您还可以调用`Bind()`和 `Unbind()`管理现有小部件的数据。以下示例显示了如何管理`String`绑定到简单`Label`小部件的数据项。

# 应用案例

[公众号-Go语言中文网，使用 Go GUI 库 fyne 编写一个计算器程序](https://mp.weixin.qq.com/s/VrTFMhpYvzr78ULqsQ15Sw)