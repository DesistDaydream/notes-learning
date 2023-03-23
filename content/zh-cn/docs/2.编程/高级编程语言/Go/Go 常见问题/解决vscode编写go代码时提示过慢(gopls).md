---
title: 解决vscode编写go代码时提示过慢(gopls)
---

前言

之前用的 vscode 的自动代码提示，发现太慢了，隔 3，4 秒才会出提示，所以换为 Google 推荐的 **gopls**来代替。

## 下载过程

方案一

打开 VS Code 的**setting**, 搜索 **go.useLanguageServe**, 并勾选上.

默认情况下, 会提示叫你 reload，重新打开之后，右下角会自动弹出下载的框框，点击 **install**即可。

如果下载时间过长，不成功，可以看**方案二**

方案二

直接上 github 下载，下载下来 之后`go install github.com/golang/tools/cmd/gopls`安装

方案三

`go get golang.org/x/tools/gopls@latest`，不需要加 u 可以去 github 上看看文档是怎么说的

github

## 配置过程

在 github 文档里有提示 `Use the VSCode-Go plugin, with the following configuration:`

```json
"go.useLanguageServer": true,
"[go]": {
    "editor.snippetSuggestions": "none",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true,
    }
},
"gopls": {
    "usePlaceholders": true, // add parameter placeholders when completing a function

    // Experimental settings
    "completeUnimported": true, // autocomplete unimported packages
    "watchFileChanges": true,  // watch file changes outside of the editor
    "deepCompletion": true,     // enable deep completion
},
"files.eol": "\n", // formatting only supports LF line endings
```

所以现在需要配置一下 `setting.json`配置文件。

在 VSCode 中按下`Ctrl + Shift + P`，在搜索框中输入`settings`,找到`Open Settings:JSON` 添加上面那段代码即可~
