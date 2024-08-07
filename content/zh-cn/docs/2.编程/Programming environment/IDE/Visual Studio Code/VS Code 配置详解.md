---
title: VS Code 配置详解
---

# 概述

# 配置详解

```json
{
  // 通用配置
  // 资源管理器自动展开功能。即.追踪代码时，是否自动展开资源管理器中的目录。
  // 该配置在前端代码中，最好关闭，因为前端的依赖都在本地目录中，展开之后太乱了
  "explorer.autoReveal": true,
  // 启用后，将不会显示扩展建议的通知。
  "extensions.ignoreRecommendations": true,
  // 终端在 Windows 上使用的 shell 的路径(默认: C:\WINDOWS\System32\WindowsPowerShell\v1.0\powershell.exe)。
  // [详细了解如何配置 shell](https://code.visualstudio.com/docs/editor/integrated-terminal#_configuration)。
  // "terminal.integrated.shell.windows": "D:\\Tools\\Git\\bin\\bash.exe",
  // 上面的指令已弃用，使用下面的方式来配置 shell
  "terminal.integrated.profiles.windows": {
    "PowerShell": {
      "source": "PowerShell",
      "icon": "terminal-powershell"
    },
    "bash": {
      "path": "D:\\Tools\\Git\\bin\\bash.exe"
    }
  },
  "terminal.integrated.defaultProfile.windows": "Ubuntu-20.04 (WSL)",
  "terminal.integrated.automationShell.linux": "",
  // 控制侧边栏和活动栏的位置。它们可以显示在工作台的左侧或右侧。
  "workbench.sideBar.location": "left",
  // 默认行尾字符，可用的值有如下几个：
  // \n # 表示 LF
  // \r\n # 表示 CRLF
  // auto # 表示 使用具体操作系统规定的行末字符。
  "files.eol": "\n",
  // 配置语言的文件关联(如: `"*.extension": "html"`)。这些关联的优先级高于已安装语言的默认关联。
  "files.associations": {},
  "editor.columnSelection": false,
  // 控制编辑器是否显示控制字符。
  "editor.renderControlCharacters": false,
  // 扩展 Git
  // 若设置为 true，则自动从当前 Git 存储库的默认远程库提取提交。若设置为“全部”，则从所有远程库进行提取。
  "git.autofetch": true,
  "git.enableSmartCommit": true,
  // 扩展 Go
  "go.toolsManagement.autoUpdate": true,
  "go.useLanguageServer": true,
  // 扩展 YAML
  "yaml.schemas": {
    "file:///toc.schema.json": "/toc\\.yml/i"
  },
  "yaml.schemaStore.enable": false,
  "redhat.telemetry.enabled": false,
  // 扩展 Kubernetes
  "vs-kubernetes": {
    "vs-kubernetes.knownKubeconfigs": ["D:\\Tools\\vs-kubernetes\\config"],
    "vs-kubernetes.kubeconfig": "D:\\Tools\\vs-kubernetes\\config",
    "vs-kubernetes.kubectl-path.windows": "D:\\Tools\\vs-kubernetes\\tools\\kubectl\\kubectl.exe",
    "vs-kubernetes.minikube-path.windows": "D:\\Tools\\vs-kubernetes\\tools\\minikube\\windows-amd64\\minikube.exe",
    "vs-kubernetes.helm-path.windows": "D:\\Tools\\vs-kubernetes\\tools\\helm\\windows-amd64\\helm.exe",
    "vs-kubernetes.draft-path.windows": "D:\\Tools\\vs-kubernetes\\tools\\draft\\windows-amd64\\draft.exe"
  },
  // 扩展 Ansible
  // 由于插件依赖于 ansible-lint 程序，所以在 windows 中关闭该功能以防止持续报错 `spawnSync C:\WINDOWS\system32\cmd.exe ENOENT`
  "ansible.ansibleLint.enabled": false
}
```

# 调试配置

**configurations**

- **cwd**(STRING) # 指定调试时的工作目录。
  - 在程序中指定某个需要读取的文件路径时，若 env 为空则读取文件可能会失败。因为调试程序默认以运行文件为当前工作目录。
  - 可以设置为 `${workspaceFolder}` 则让调试程序以当前工作空间目录作为调试时的工作目录。

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch file",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${file}",
      "env": {
        "NAMESPACE": "default"
      },
      "cwd": "${workspaceRoot}"
      "args": [
        "alidns",
        "-d",
        "desistdaydream.ltd",
        "-o",
        "list",
        "-u",
        "断灬念梦",
        "-F",
        "owner.yaml"
      ]
    }
  ]
}
```
