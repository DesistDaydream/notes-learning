---
title: Plugins
linkTitle: Plugins
weight: 20
---

# 概述

> 参考：
>
> -

## 关联文件与配置

在每个项目的根目录下有这么一个目录： `${Project}/.vscode/`，所有适用于本项目的插件配置通常都会保存在该目录中。

# Debug 插件

> 参考：
>
> - [官方文档，用户指南 - Debugging](https://code.visualstudio.com/docs/editor/debugging)
> - <https://www.barretlee.com/blog/2019/03/18/debugging-in-vscode-tutorial/>

Debug 插件的默认配置文件名为 `launch.json`

## launch.json

https://code.visualstudio.com/docs/editor/debugging#_launchjson-attributes

**cwd** # 运行程序的工作路径

**program** # 启动调试器时要运行的可执行文件或文件

**args** # 运行程序的参数

## 简单示例

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

# SFTP 插件

https://github.com/liximomo/vscode-sftp#connection-hopping

- 新的 SFTP 插件，上面的已不更新。 https://github.com/Natizyskunk/vscode-sftp

```json
{
  "name": "ansible",
  "host": "172.38.40.250",
  "protocol": "sftp",
  "port": 22,
  "username": "root",
  "password": "XXXXXX",
  "remotePath": "/root/projects",
  "uploadOnSave": true,
  "ignore": [".vscode", ".git", ".DS_Store"]
}
```

跳转服务器配置

```json
{
    "name": "DesistDaydream",
    "remotePath": "/root/projects",
 // 用于作为代理的服务器信息
    "host": "192.168.1.10",
    "protocol": "sftp",
    "port": 42203,
    "username": "root",
    "password": "XXXXX",
 // 最终目标服务器信息
    "hop": {
        "host": "172.19.42.248",
        "port": 22,
        "username": "root",
        "password": "XXXXX"
    },
    "uploadOnSave": true,
 "ignore": [
        ".vscode",
        ".git",
        ".DS_Store"
    ]
}
```
