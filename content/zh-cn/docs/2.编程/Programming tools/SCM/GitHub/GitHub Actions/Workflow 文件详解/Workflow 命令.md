---
title: Workflow 命令
---

# 概述

> 参考：
> - [官方文档，使用 Workflow-Workflow 命令](https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions)

我们可以在 Workflow 中执行 Shell 命令时，使用 GitHub Action Workflow 的特殊命令。这些 Workflow 命令可以与运行机器通信以 _设置环境变量、输出值、添加 debug 消息到 输出日志_ 等等。

大多数 Workflow 命令以特定格式通过 `echo` 命令使用，语法为：

`echo "::WorkflowCommand Parameter1={DATA},Parameter2={DATA},...::{COMMAND|VALUE}"`

下面是一个设置步骤的输出内容的示例：

```yaml
jobs:
  test-workflow-command:
    runs-on: ubuntu-latest
    steps:
      - name: 设置 ARCH 变量的值为 uname -m 命令执行的结果
        run: echo '::set-output name=ARCH::$(uname -m)'
        id: arch-generator
      - name: 从 id 为 arch-generator 步骤的输出中，获取 ARCH 的值
        run: echo "The selected color is ${{ steps.arch-generator.outputs.ARCH }}"
```
