---
title: VS Code 快捷键
---

# 概述

> 参考：
> - [官方文档](https://code.visualstudio.com/shortcuts/keyboard-shortcuts-windows.pdf),pdf
> - [官方文档](https://docs.microsoft.com/zh-cn/visualstudio/ide/default-keyboard-shortcuts-in-visual-studio)
> - <https://docs.microsoft.com/zh-cn/visualstudio/ide/productivity-shortcuts>

多个组合按键表示需要连续按

# 全局快捷键

Ctrl+q 搜索 Visual Studio。主要用于快速切换 VS Code 本身的功能
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/izq5cy/1640249077536-e62dfef8-9671-42f9-a9be-cb61bf2ad504.png)
Ctrl+k,Ctrl+s 打开键盘快捷方式列表

# 视图

**Ctrl+Shift+e** # 显示资源管理器
**Ctrl+Shift+g** # 显示源代码管理

- **Ctrl+k,Ctrl+F5** # (自定义)Git 刷新

**Ctrl+Shift+f** # 显示搜索

**Ctrl+Alt+→ **将编辑器拆分到新组
**Ctrl+Alt+→ **将编辑器合并到前一个组
**Ctrl+PageUp** 选中上一个编辑器
**Ctrl+PageDown** 选中下一个编辑器
**Ctrl+k,Ctrl+→** 选中右侧编辑器组
**Ctrl+k,Ctrl+←** 选中左侧编辑器组

## 组内操作

**Alt+NUM** 选中组中第 NUM 号编辑器

# 编辑时操作

操作光标所在代码

- 折叠
  - **Ctrl+k,Ctrl+\[** 递归折叠光标所在代码块所有层
  - **Ctrl+Shift+\[** 折叠光标所在代码代码 1 层
- 展开
  - **Ctrl+k,Ctrl+]** 递归展开光标所在代码块所有层
  - **Ctrl+Shift+]** 展开光标所在代码代码 1 层

操作文件所有代码

- 折叠
  - **Ctrl+k,Ctrl+0** 折叠所有代码全部层
  - **Ctrl+k,Ctrl+NUM **折叠所有代码的 NUM 层
  - **Ctrl+k,Ctrl+/** 折叠所有代码的注释
- 展开
  - **Ctrl+k,Ctrl+j** 展开所有代码的全部层

**Alt+Shift+a** 添加代码块注释
**Ctrl+/** 添加代码行注释
**Alt + z** # 开启/关闭 自动换行功能。自动换行用于当一行内容过长时，可以自动换行显示。

## 光标移动

光标移动通常以 Ctrl 与 方向键 为主
**Ctrl+Alt+↑** 向上添加一个光标。多次使用，可以添加多个光标
**Ctrl+Alt+↓** 向下添加一个光标。多次使用，可以添加多个光标

## 光标选中

代码选择快捷键通常都是以 Shift 与 方向键 为主
**Shift+→** 选中光标右侧的字符，多次使用可选中多个
**Shift+←** 选中光标左侧的字符，多次使用可选中多个
**Shift+Alt+→** 选中光标所在的单词。将下划线连接的单词当做多个个体
**Shift+Alt+→** 取消选中光标所在的单词。将下划线连接的单词当做多个个体
**Shift+Ctrl+→ **选中光标所在位置到右侧下一个单词的结尾。将下划线连接的单词当做一个整体
**Shift+Ctrl+← **选中光标所在位置到左侧下一个单词的结尾。将下划线连接的单词当做一个整体

## 搜索内容

**Ctrl+d** 搜索选中的字符。多次使用的话，就是搜索下一个的效果。

> 若没有选中的字符，则搜索当前光标所在的单词

**Ctrl+Shift+l** 搜索所有与选中字符一样的字符

# 代码跳转

Ctrl + F12 跳转到当前代码的实现
Shift + F12 跳转到当前代码的引用
Alt + ← 跳回。在跳转代码后，跳回上一个位置

# 其他

Ctrl + Shift + p 打开命令面板
Ctrl+r 打开工作区
