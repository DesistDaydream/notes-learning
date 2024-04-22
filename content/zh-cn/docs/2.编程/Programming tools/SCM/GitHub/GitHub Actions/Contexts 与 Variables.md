---
title: Contexts 与 Variables
---

# 概述

> 参考：
> 
> - [官方文档，学习 GitHub Actions-上下文](https://docs.github.com/en/actions/learn-github-actions/contexts)
> - [官方文档，学习 GitHub Actions-环境变量](https://docs.github.com/en/actions/learn-github-actions/environment-variables)

GitHub Actions 中可以通过 **Contexts(上下文)** 与 **Environment Variables(环境变量)** 来暴露工作流的信息或引用工作流的信息。就像下面的示例一样：

这是一个环境变量的示例：

> GitHub Action 中的环境变量本质上是 Shell 中的变量，引用方式也是一样的。

```yaml
name: Greeting on variable day

on: workflow_dispatch

env:
  DAY_OF_WEEK: Monday

jobs:
  greeting_job:
    runs-on: ubuntu-latest
    env:
      Greeting: Hello
    steps:
      - name: "Say Hello Mona it's Monday"
        run: echo "$Greeting $First_Name. Today is $DAY_OF_WEEK!"
        env:
          First_Name: Mona
```

这是一个上下文的示例：

```yaml
name: Greeting on variable day

on: workflow_dispatch

env:
  DAY_OF_WEEK: Monday

jobs:
  greeting_job:
    runs-on: ubuntu-latest
    env:
      Greeting: Hello
    steps:
      - name: "Say Hello Mona it's Monday"
        run: echo "${{ env.Greeting }} ${{ env.First_Name }}. Today is ${{ env.DAY_OF_WEEK }}!"
        env:
          First_Name: Mona
```

从示例中可以看到，想要使用 **Context(上下文)**，需要使用一种特殊的语法，这种语法称为 **Expressions(表达式)**。

## Expressions(表达式)

在 GitHub Actions 的 Workflow 文件中，我们可以使用 **Expressions(表达式)** 设置和访问** 环境变量 **或访问**上下文 **信息。表达式可以是 _字面量、上下文引用、函数_ 的任意组合。

以 `$` 开口，`{{ }}` 括起来的内容即为表达式的语法，当 GitHub Action 运行时，Workflow 中的 `${{ <EXPRESSION> }}` 内容会被解析为表达式进行处理，表达式就像模板一样，解析完成后，使用实际的值替换表达式。以实现以变成的方式设置 Workflow 文件。

> 注意：
>
> Workflow 中的 if 字段，会自动将其下的值解析为表达式，所以可以省略 `${{ }}` 符号。

表达式 Fiterals(字面量)

表达式 Operators(运算符)

表达式 Functions(函数)

- contains
- startsWith
- endsWith
- format
- join
- toJSON
- fromJSON
- hashFiles
- 状态检查函数
  - success
  - always
  - cancelled
  - failure
  - Evaluate Status Explicitly

函数使用示例

```yaml
name: print
on: push
env:
  continue: true
  time: 3
jobs:
  job1:
    runs-on: ubuntu-latest
    steps:
      - continue-on-error: ${{ fromJSON(env.continue) }}
        timeout-minutes: ${{ fromJSON(env.time) }}
        run: echo ...
```

## Context(上下文)

GitHub Action 中的上下文，是一种功能更丰富的环境变量，并且我们可以通过上下文的语法引用环境变量。

在 Expressions(表达式) 中使用 Context(上下文)，可以让我们在 Workflow 文件中访问工作流运行信息、运行器环境信息、Job 信息、每个 Job 下的 Step 信息。

每个上下文都是一个包含 **Properties(属性)** 的 **Object(对象**)，Properties 可以是字符串或其他 Objects

在表达式语中使用上下文的语法为：`${{ Object.Properties }}`。每个 Object 可以提供丰富的信息。

现阶段 Action 有多个上下文可用：

- GitHub 本身信息相关的上下文
  - **github** # Information about the workflow run.
  - **secrets** # Contains the names and values of secrets that are available to a workflow run.
- Job 信息相关的上下文
  - **needs** # Contains the outputs of all jobs that are defined as a dependency of the current job
  - **env** # Contains environment variables set in a workflow, job, or step.
  - **job** # Information about the currently running job.
  - **steps** # Information about the steps that have been run in the current job.
- **runner** # Information about the runner that is running the current job.
- **strategy** # Information about the matrix execution strategy for the current job.
- **matrix** # Contains the matrix properties defined in the workflow that apply to the current job.
- **inputs** # Contains the inputs of a reusable workflow. For more information, see inputscontext

# GitHub 本身信息相关的上下文

## github 上下文

github 上下文包含本次工作流的事件信息，还有很多 GitHub 信息，比如 `github.actor` 属性表示发起工作流的用户名，如果这个项目只有自己一个人，那就是这个仓库的拥有者名称~

## secrets 上下文

对于工作流程运行中的每个 Job，此上下文都是相同的。 您可以从 Job 中的任何步骤访问此上下文。 此对象包含下面列出的所有属性。

| 属性名称                  | 类型     | 描述                                                                                                                        |
| --------------------- | ------ | ------------------------------------------------------------------------------------------------------------------------- |
| secrets.GITHUB_TOKEN  | String | 为每个工作流程运行自动创建的令牌。 更多信息请参阅“[自动令牌身份验证](https://docs.github.com/cn/actions/security-guides/automatic-token-authentication)”。 |
| secrets.\<SecretName> | String | 特定 Secret 的值                                                                                                              |

SecretName 可以在在一个项目的设置中添加，如下图：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/qaf8qw/1652758616297-59e0ad52-f622-4d41-9235-b36a995ee87d.png)

我们为本仓库添加了一个 SecretName 为 DOCKER_REGISTRY_PASSWORD 的 Secret，那么，在 Actions 中，可以使用 `${{ secrets.DOCKER_REGISTRY_PASSWORD }}` 引用 DOCKER_REGISTRY_PASSWORD 的值。

这个上下文常用在登录行为，以防止将密码以明文形式呈现，比如下面的示例，我们将会通过 secrets 上下文登录容器镜像仓库，并推送构建的镜像到仓库中。

```yaml
jobs:
  # 推送到 docker.io
  push-docker:
    needs: [generate-tags]
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository code
        uses: actions/checkout@v2
      - name: 登录镜像仓库网站
        uses: docker/login-action@v1
        with:
          registry: docker.io
          username: ${{ secrets.DOCKER_REGISTRY_USERNAME }}
          password: ${{ secrets.DOCKER_REGISTRY_PASSWORD }}
      - name: 构建并推送容器镜像
        uses: docker/build-push-action@v2
        with:
          context: .
          file: Dockerfile
          push: true
          tags: docker.io/lchdzh/e37-exporter:${{needs.generate-tags.outputs.tag}}
```

# Job 信息相关的上下文

## needs 上下文

needs 上下文中包含了由 jobs.JOB_ID.needs 字段定义的依赖 Job 中的信息。包含如下属性：

| needs.\<JOB_ID>                        | OBJECT | A single job that the current job depends on.                                           |
| -------------------------------------- | ------ | --------------------------------------------------------------------------------------- |
| needs.\<JOB_ID>.outputs                | OBJECT | The set of outputs of a job that the current job depends on.                            |
| needs.\<JOB_ID>.outputs.\<OUTPUT_NAME> | STRING | 当前 Job 所依赖的 Job 的输出中，OUTPUT_NAME 输出的值。                                  |
| needs.\<JOB_ID>.result                 | STRING | 当前 Job 所依赖的 Job 的运行结果。可能的值有：success, failure, cancelled, or skipped。 |

## steps 上下文

此上下文针对作业中的每个步骤而改变。 您可以从作业中的任何步骤访问此上下文。 此对象包含下面列出的所有属性。
steps 上下文中包含了已指定 JOB_ID 且已运行的 Job 中的每个 step 的信息，包含如下属性

| 属性名称                                | 类型   | 描述                                                                                                                                                                                                                                                                                                            |
| --------------------------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| steps.\<step_id>.outputs                | 对象   | 为步骤定义的输出集。 更多信息请参阅“[GitHub Actions 的元数据语法](https://docs.github.com/cn/articles/metadata-syntax-for-github-actions#outputs-for-docker-container-and-javascript-actions)”。                                                                                                                |
| steps.\<step_id>.conclusion             | 字符串 | 在 [continue-on-error](https://docs.github.com/cn/actions/reference/workflow-syntax-for-github-actions#jobsjob_idstepscontinue-on-error) 应用之后完成的步骤的结果。 可能的值包括 success、failure、cancelled 或 skipped。 当 continue-on-error 步骤失败时，outcome 为 failure，但最终的 conclusion 为 success。 |
| steps.\<step_id>.outcome                | 字符串 | 在 [continue-on-error](https://docs.github.com/cn/actions/reference/workflow-syntax-for-github-actions#jobsjob_idstepscontinue-on-error) 应用之前完成的步骤的结果。 可能的值包括 success、failure、cancelled 或 skipped。 当 continue-on-error 步骤失败时，outcome 为 failure，但最终的 conclusion 为 success。 |
| steps.\<step_id>.outputs.\<output_name> | 字符串 | 指定 StepIP 步骤中指定的 OutputName 输出的值。                                                                                                                                                                                                                                                                  |

### 示例

```yaml
name: Generate random failure
on: push
jobs:
  randomly-failing-job:
    runs-on: ubuntu-latest
    steps:
      - id: checkout
        uses: actions/checkout@v3
      - name: Generate 0 or 1
        id: generate_number
        run: echo "::set-output name=random_number::$(($RANDOM % 2))"
      - name: Pass or fail
        run: |
          if [[ ${{ steps.generate_number.outputs.random_number }} == 0 ]]; then exit 0; else exit 1; fi
```

通过 `${{ steps.generate_number.outputs.random_number }}` 引用了当前 Job 中 ID 为 generate_number 这个步骤的输出中，random_number 的值。
