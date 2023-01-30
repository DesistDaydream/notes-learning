---
title: "PowerShell"
linkTitle: "PowerShell"
weight: 1
---

# 概述
> 参考：
> - [官方文档，PowerShell](https://learn.microsoft.com/en-us/powershell)

PowerShell 是一种跨平台的任务自动化解决方案，由命令行 shell、脚本语言和配置管理框架组成。 PowerShell 在 Windows、Linux 和 macOS 上运行。

PowerShell 的独特之处在于，它接受并返回 .NET 对象，而非纯文本。这个特点让 PowerShell 可以更轻松地在一个管道中串联不通的命令。
> 这里面所说的对象，就是面向对象编程中常说的“对象”，就像 Go 语言中的 Struct 类似的东西，只不过是 .NET 语言中的对象。

这些对象在被接收后，再交给格式化函数处理，以人类可读的方式，输出出来。

## PowerShell 命令

PowerShell 中可以执行的命令分两类
- 系统上的可执行命令。
- cmdlet

PowerShell 内置了一组 **cmdlet(全称 command-lets)**，cmdlet 并不是一个独立的可执行文件，而是一种统称，cmdlet 被收集在 PowerShell 模块中，可以按需加载它们。可以用任何编译的 .NET 语言或 PowerShell 脚本语言来编写 cmdlet。

cmdlet 中每个命令的名称都是由 **Verb-Noun(动词-名词)** 组成，比如 Get-Command 命令用于获取在 CLI 中注册的所有 cmdlet。

我们可以通过如下几个命令来搜索可用的 cmdlet
- `Get-Verb` # 获取所有可用的动词
- `Get-Command` # 此命令会检索计算机上安装的所有命令的列表。
- `Get-Member` # 基于其他命令运行，可以获取 cmdlet 返回的对象信息，包括对象中的**属性、方法、等等**
- `Get-Help` # 以命令名称为参数调用此命令，将显示一个帮助页面，其中说明了命令的各个部分。

从本笔记的 [Windows 管理工具](/docs/IT学习笔记/1.操作系统/Y.Windows%20管理/Windows管理工具/_index.md) 目录查找所有可用的命令，以及查看命令的用法