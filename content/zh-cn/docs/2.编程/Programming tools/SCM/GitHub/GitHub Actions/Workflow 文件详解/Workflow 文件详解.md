---
title: Workflow 文件详解
---

# 概述

> 参考：
>
> - [官方文档，使用工作流-触发工作流](https://docs.github.com/en/actions/using-workflows/triggering-a-workflow)
> - [官方文档，使用工作流-触发工作流的事件](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows)
> - [官方文档，使用工作流-Workflow 语法](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

GitHub 的 Actions 通过 [YAML](/docs/2.编程/无法分类的语言/YAML.md) 格式的文件来定义运行方式。工作流文件必须保存在项目根目录下的 `.github/workflows/` 目录下

# 顶层字段

- **name**(STRING) # Workflow 的名称。`默认值：当前 Workflow 的文件名`。
- **run-name** #
- **on**[(on)](#on) # 指定触发 Workflow 的条件
- **permissions**
- **env**
- **defaults**
- **concurrency**
- **jobs**[(jobs)](#jobs) #  Workflow 文件的主体，用于定义要执行的一项或多项任务

# on

https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#on

这个字段用来定义触发工作流的事件，在这里可以看到 GitHub 支持的所有事件，通常包含如下字段

- **push**([push](#push)) # 当上传代码时，触发 Workflow
- **pull_request**([pull_request](#pull_request)) # 当发生 PR 时，触发 orkflow
- **schedule**(\[][schedule](#schedule)) # 定时触发 Worlkflow
- **workflow_dispatch**([workflow_dispatch](#workflow_dispatch)) # 手动触发 Workflow
- ......

## push

**branches: <[]STRING>** # 指定出发条件，当上传代码到该字段指定的分支时，触发 Workflow

## pull_request

## schedule

使用 POSIX cron 语法让 Worlkflow 在指定时间运行

```yaml
on:
  schedule:
    # * is a special character in YAML so you have to quote this string
    - cron: "30 5,17 * * *"
```

## workflow_dispatch

```yaml
on:
  workflow_dispatch:
    inputs:
      file:
        description: "指定要使用的镜像同步文件的路径"
        type: string
        required: true
jobs:
  build:
    ......
    steps:
      - name: images sync
        # 这里可以调用 inputs 中定义的变量，这些变量通过 Web 页面传递进来，也可以通过 CLI 传递进来。
        run: |
          echo ${{ github.event.inputs.file }}
```

下面对话框中填写的值将传入 Action 中，作为 `file` 变量的值

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sytu80/1643186313475-dfed2719-28b6-4680-8a28-b6a6772763c8.png)

**inputs(OBJECT)** # 触发 Workflow 时，传入的信息

更多 GitHub 可用的传入信息，详见 [Contexts](https://docs.github.com/en/actions/learn-github-actions/contexts#github-context)

**NAME(OBJECT)** # 定义变量。这里的 NAME 可以任意字符串，然后在 workflow 文件中使用`${{ github.event.inputs.NAME }}`的方式调用

- **description(STRING)** # 对 NAME 的描述
- **type(STRING)** # 可用的类型有 string、number、boolean、choice、environment
- **required(BOLLEAN)** #
- **options([]TYPE)** # 为 choice 类型提供可用选择的列表

# jobs

https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobs

jobs 字段是 Workflow 文件的主体，用于定义要执行的一项或多项任务

使用 **`jobs.<JOB_ID>`** 为我们的工作提供唯一标识符，JOB_ID 是一个字符串，必须以字母或 `*` 开头，并且仅能包含字母、数字、下划线、中横线。一个最简单的不用执行任何具体行为的 jobs 配置如下：

```yaml
jobs:
  my_first_job:
    name: My first job
  my_second_job:
    name: My second job
```

示例中 `my_first_job` 就是 JOB_ID

通常包含如下字段

- **JOB_ID.needs**(\[][needs](#needs) | [needs](#needs)) # 此 Job 必须在指定的 JOB_ID 成功后才可以执行
- **JOB_ID.runs-on**(STRING) # 必须的。运行 JOB_ID 的运行器。GitHub 自带的运行器有：ubuntu-latest、windows-latest、macos-latest 等等
- **JOB_ID.steps**(\[][steps](#steps)) # Job 的运行步骤
- **JOB_ID.outputs**([outputs](#outputs)) # 创建输出映射

## needs

```yaml
jobs:
  job1:
  job2:
    needs: job1
  job3:
    needs: [job1, job2]
```

上面这个示例表示 job2 等待 job1 成功后开始执行；job3 等待 job1 和 job2 都成功后开始执行。

## steps

**env**(map\[STRING]STRING) # 设定前 Job 中可用的环境变量。

**name**(STRING) # 当前 Job 的名称。

**run**(STRING) # 运行命令。使用 runs-on 中指定的操作系统的 shell 运行。

```yaml
# 单行命令
- name: Install Dependencies
  run: npm install
# 多行命令
- name: Clean install dependencies and build
  run: |
    npm ci
    npm run build
# working-directory 字段与 run 关联使用，可以用来指定运行命令的工作目录
- name: Clean temp directory
  run: rm -rf *
  working-directory: ./temp
# shell 字段与 run 关联使用。可以用来指定运行命令的 shell
steps:
  - name: Display the path
    run: echo $PATH
    shell: bash
```

**uses**(STRING) # 当前步骤要使用的 Action。

在这里可以指定其他 Action 作为工作流的一部分来运行，本质上，Action 是可重用的代码。其实就类似于在代码中调用函数一样，`uses` 字段可以理解为调用某个函数，这个函数就是指其他的 Action。在[这篇文章](/docs/2.编程/SCM/GitHub/GitHub%20Actions/好用的%20Action.md Action.md)中，介绍了很多比较好用的 Action。

通过使用其他 Action，可以大大简化自身工作流的配置文件。比如 Git Action 官方提供的 [actions/checkout](https://github.com/actions/checkout) 这个 Action，可以用来将仓库中的代码，拷贝到运行 Action 的容器中，然后进行后续操作，如果不使用这个 Action，那么我们就要写很多命令来 pull 代码了~

## outputs

这是一个 `map[STRING]STRING` 类型的数据。

通过 outputs 行为可以为本 job 创建一个输出映射，本 job 的输出可以用于其他依赖本 job 的所有下游 job。job 之间的依赖关系通过 JOB_ID.needs 行为确定。

```yaml
jobs:
  job1:
    runs-on: ubuntu-latest
    # 将本 job 的 step 的输出映射到 job 的输出
    outputs:
      output1: ${{ steps.step1.outputs.test }}
      output2: ${{ steps.step2.outputs.test }}
    steps:
      - id: step1
        run: echo "::set-output name=test::hello"
      - id: step2
        run: echo "::set-output name=test::world"
  job2:
    runs-on: ubuntu-latest
    needs: job1
    steps:
      - run: echo ${{needs.job1.outputs.output1}} ${{needs.job1.outputs.output2}}
```

> 双冒号中的语法为 [Workflow 命令](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/Workflow%20文件详解/Workflow%20命令.md)

job1 创建了 `output1` 变量，值为 `hello`，同时创建了 `output2` 变量，值为 `world`。

job2 中首先通过 needs 创建依赖关系，然后通过 `${{ needs.job1.outputs.output1 }}` 与 `${{ needs.job1.outputs.output2 }}` 引用 job1 中输出的变量。

通过 needs 上下文引用值得表达式语法详见 [Context,needs](https://docs.github.com/en/actions/learn-github-actions/contexts#needs-context)

## strategy

**strategy(策略，尤指为获得某物制定长期的策略)** 可以帮我们创建一个 matrix strategy(矩阵策略)，这是一个类似循环的功能，比如我们创建了这么一个 Workflow

```yaml
jobs:
  example_matrix:
    strategy:
      matrix:
        version: [10, 12, 14]
        os: [ubuntu-latest, windows-latest]
```

将为每个可能的变量组合运行一个 job。在此示例中，Workflod 将运行六个 job，每个作业对应操作系统和版本变量的组合。

默认情况下，GitHub 将根据运行器的可用性最大化并行运行的 job 数量。矩阵中变量的顺序决定了作业的创建顺序。您定义的第一个变量将是在工作流程运行中创建的第一个作业。例如，上面的 strategy 将按以下顺序创建 job：

- `{version: 10, os: ubuntu-latest}`
- `{version: 10, os: windows-latest}`
- `{version: 12, os: ubuntu-latest}`
- `{version: 12, os: windows-latest}`
- `{version: 14, os: ubuntu-latest}`
- `{version: 14, os: windows-latest}`

> Notes: 每个 Workflow 通过 strategy 最多生成 256 个作业。

strategy.matrix 下定义的字段还可以作为 [Contexts 与 Variables](docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/Contexts%20与%20Variables.md) 变量使用，e.g. `${{ matrix.version }}` 将会获取当前 job 下对应的 version 的值，比如上面例子中，进行到第三个任务，那么 version 的值为 12