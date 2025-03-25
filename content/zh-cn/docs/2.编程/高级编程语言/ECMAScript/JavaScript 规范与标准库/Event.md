---
title: Event
linkTitle: Event
weight: 20
---

# 概述

> 参考：
>
> - [网道，Javascript 教程-事件](https://wangdoc.com/javascript/events/index.html)

**Event(事件)** 的本质是程序各个组成部分之间的一种通信方式，也是异步编程的一种实现。DOM 支持大量的事件，

事件语法：

```javascript
elementObject.onEventType=Function(){}
```

- **elementObject** # 元素对象，也就是事件源，通常是通过类似 `getElementByID()`之类的方法获取到的 HTML 中的元素
- **onEventType** # 以`on`开头，后面跟一个事件名称
- **Function(){}** # elementObject 触发 onEventType 时要执行的代码，也就是 HTML 中某一元素触发事件时要执行的操作

比如：

```javascript
var divElement = document.getElementById("event")
// click 点击事件
divElement.onclick = function () {
  console.log(divElement, "元素，被点击了一下")
}
```

# 事件类型

## 鼠标事件

- click # 鼠标单击
- dblclick # 鼠标双击
- contextmenu # 左键单击
- mousedown # 鼠标按下
- mouseup # 鼠标抬起
- mousemove # 鼠标移动
- mouseenter # 鼠标移入
- mouseleave # 鼠标移出
- ...... 等等

## 键盘事件

- keydown # 键盘按下
- keyup # 键盘抬起
- keypress # 键盘输入
- ...... 等等

## 浏览器事件

- load # 加载完毕
- scroll # 滚动
- resize # 尺寸改变
- ...... 等等

## 触摸事件

- touchstart # 触摸开始
- touchmove # 触摸移动
- touchend # 触摸结束
- ...... 等等

## 表单事件

- focus # 聚焦
- blue # 失焦
- change # 改变
- input # 输入
- submit # 提交
- reset # 重置
- ...... 等等

# 事件对象

每个事件触发时，都会记录一组数据，这组数据是事件类型对象，事件对象数据中的数据包括该时间一系列属性信息，比如：

- **type** # 什么事件
- **target** # 谁触发的
- 如果是鼠标事件，那么还会记录
  - **x** # 光标 x 坐标
  - **y** # 光标 y 坐标
- **等等......**

比如我们模拟一下鼠标点击事件：

```javascript
var divElement = document.getElementById("event")

// click 触发一个点击事件
divElement.onclick = function (prop) {
  // 事件的属性
  console.log(prop)
}
```

这个鼠标点击事件具有如下属性：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nzci2k/1641961554049-e88a8221-0593-4b6e-be48-6a093e0d1b0a.png)
这是一个 PointerEvent 类型的对象

## 鼠标事件对象

坐标信息

- **offsetX 和 offsetY** # 相对于触发事件元素的坐标
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nzci2k/1641961872393-217cf16a-78e5-449c-9e33-d91f52c6dd96.png)
- **client 和 clientY**# 相对于浏览器可视窗口的坐标
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nzci2k/1641961917878-40b8b649-3846-4f31-81f1-9e98e407f320.png)
- **pageX 和 pageY** # 相对于页面文档流的坐标
  - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/nzci2k/1641961945785-c13cc660-5702-444f-b6a3-2d21a341341f.png)

## 键盘事件对象

-
