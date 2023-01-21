---
title: Plugins(插件)
---

# 概述

## 关联文件与配置

在每个项目的根目录下有这么一个目录： `${Project}/.vscode/`，所有适用于本项目的插件配置通常都会保存在该目录中。

# Debug 插件

> 参考：
> - [官方文档，用户指南-Debugging](https://code.visualstudio.com/docs/editor/debugging)
> - <https://www.barretlee.com/blog/2019/03/18/debugging-in-vscode-tutorial/>

Debug 插件的默认配置文件名为 `launch.json`

cwd # 运行程序的工作路径
program # 要运行的程序
args # 运行程序的参数

简单示例

```json
{
  // 使用 IntelliSense 了解相关属性。
  // 悬停以查看现有属性的描述。
  // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Package",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "cwd": "${workspaceRoot}",
      "program": "cmd/statistics/main.go",
      "args": ["-s", "dp"]
    }
  ]
}
```
