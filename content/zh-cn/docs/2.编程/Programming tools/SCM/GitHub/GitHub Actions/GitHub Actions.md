---
title: GitHub Actions
weight: 1
---

# 概述

> 参考：
>
> - [官方文档](https://docs.github.com/cn/actions)
> - [官方文档，学习 GitHub Actions - GitHub Actions 简介](https://docs.github.com/en/actions/learn-github-actions/introduction-to-github-actions)
>     - 域名中的路径改成 [理解 GitHub Actions](https://docs.github.com/en/actions/learn-github-actions/understanding-github-actions)
> - GitHub Actions 官方市场：[Actions Marketplace](https://github.com/marketplace?type=actions)
> - 阮一峰老师的一篇文章：[GitHub Actions 入门教程](http://www.ruanyifeng.com/blog/2019/09/getting-started-with-github-actions.html)
> - <https://blog.csdn.net/sculpta/article/details/104142607>

GitHub Actions 是在 GitHub Universe 大会上发布的，被 Github 主管 Sam Lambert 称为 “再次改变软件开发” 的一款重磅功能（“_we believe we will once again revolutionize software development._”）。于 2018 年 10 月推出，内测了一段时间后，于 2019 年 11 月 13 日正式上线

GitHub 会提供一个以下配置的服务器做为 runner：

- 2-core CPU
- 7 GB of RAM memory
- 14 GB of SSD disk space

（免费额度最多可以同时运行 20 个作业，心动了有木有 💘）

GitHub Actions 是一个 `CI/CD（持续集成/持续部署）`工具，持续集成由很多操作组成，比如 **抓取代码**、**运行测试**、**登录远程服务器**、**发布到第三方服务** 等等。GitHub 把这些操作统称为 **`Actions(操作、行为)`**。

Actions 是 GitHub Actions 的核心，简单来说，它其实就是一段可以执行的代码，可以用来做很多事情。

> 比如，你在 python 3.7 环境下写了一个 python 项目放到了 GitHub 上，但是考虑到其他用户的生产环境各异，可能在不同的环境中运行结果都不一样，甚至无法安装，这时你总不能在自己电脑上把所有的 python 环境都测试一遍吧
>
> 但是如果有了 GitHub Actions，你可以在 runner 服务器上部署一段 actions 代码来自动完成这项任务。你不仅可以指定它的操作系统（支持 Windows Server 2019、Ubuntu 18.04、Ubuntu 16.04 和 macOS Catalina 10.15），还可以指定触发时机、指定 python 版本、安装其他库等等
>
> 此外，它还可以用来做很多有趣的事，比如当有人向仓库里提交 issue 时，给你的微信发一条消息；爬取课程表，每天早上准时发到你的邮箱；当向 master 分支提交代码时，自动构建 Docker 镜像并打上标签发布到 Docker Hub 上 ……

慢慢的，你会发现很多操作在不同项目里面是类似的，完全可以共享。GitHub 也注意到了这一点，于是它允许开发者把每个操作写成独立的脚本文件，存放到代码仓库，使得其他开发者可以引用。如果我们需要某个 action，不必自己写复杂的脚本，直接引用他人写好的 action 即可，整个 CI/CD 过程，就变成了一个个 action 的组合。这就是 GitHub Actions 最特别的地方。

> 总而言之，GitHub Actions 就是为我们提供了一个高效易用的 CI/CD 工作流，帮助我们自动构建、测试、部署我们的代码

GitHub 做了一个官方市场(暂且称为 Actions Hub)，在这里可以搜索到其他人提交的 Actions。另外，还有一个名为 [awesome-actions](https://github.com/shink/actions-bot) 的仓库，搜罗了不少好用的 actions。

既然 actions 是代码仓库，就有版本的概念，用户可以引用某个具体版本的 action。比如下面的例子，用的就是 Git 的指针的概念。

```bash
actions/setup-node@74bc508 # 指向一个 commit
actions/setup-node@v1.0    # 指向一个标签
actions/setup-node@master  # 指向一个分支
```

## Actions 基本概念

- **Workflow(工作流程)** # 持续集成一次运行的过程，就是一个 workflow。定义了由 Event 触发的执行流程。
    - **Job(任务)** # 一个 Workflow 由一个或多个 Jobs 构成，每个 Job 都是在 runner 中运行的程序，并且该程序具有一个或多个 Setp(步骤) 以完成 Job。
        - **Runner(运行器)** # 运行 JOB 的环境，可以是 [虚拟机](/docs/10.云原生/Virtualization/Virtualization.md) 或 [容器](/docs/10.云原生/Containerization/Containerization.md)。通常由 Workflow 文件中的 `.jobs.JOB_ID.runs-on` 字段指定。
        - **Step(步骤)** # 每个 job 由多个 Step 构成，一步步完成。
            - **Action(动作)** # 每个 step 可以依次执行一个或多个命令（action）。
- **Event(事件)** # 触发 Workflow 的特定活动。比如，推送新的提交到仓库或者创建 PR，甚至可以配置 cron 定时触发 Workflow，也可以设定只由手动惦记某个按钮触发。

### Workflow

与 Jenkins、Drone 这类 CI/CD 工具一样，GitHub Actions 也有一个配置文件，用来定义要执行的操作，这个配置文件叫做 **Workflow 文件**，需要默认存放在代码仓库的 **.github/workflows** 目录中。

Workflow 文件用来定义 GitHub Actions 要执行的操作，需要存放在代码仓库的 `.github/workflows/*.yml` 目录中。

Workflow 文件是 YAML 格式，后缀名必须统一为 `.yml`。一个代码库可以有多个 workflow 文件。GitHub 只要发现 .github/workflows 目录中有 .yml 文件，就会自动根据该文件的配置运行工作流 。

# Actions 关联文件与配置

**.github/workflows/** # 工作流文件保存目录

详见 [Actions 配置](/docs/2.编程/Programming%20tools/SCM/GitHub/GitHub%20Actions/Actions%20配置.md)

# 简单示例

- 从 GitHub 上的仓库，在 .github/workflows/ 目录中创建一个名为 github-actions-demo.yml 的新文件。 更多信息请参阅“[创建新文件](https://docs.github.com/cn/github/managing-files-in-a-repository/creating-new-files)”。
- 将以下 YAML 内容复制到 github-actions-demo.yml 文件中：

```yaml
# 定义了 Workflow 名称
name: GitHub Actions Demo
# 定义了触发 Workflow 的 Event
on: [push]
# 定义了 Workflow 中包含的 Job
jobs:
  Explore-GitHub-Actions:
    # 定义运行这个 job 的 Runner。类似 Dockerfile 中的 FORM 指令。
    runs-on: ubuntu-latest
    # 定义运行该 Job 应该执行的 Step
    steps:
      # 定义了该 Step 的 Action
      - run: echo "🎉 The job was automatically triggered by a ${{ github.event_name }} event."
      - run: echo "🐧 This job is now running on a ${{ runner.os }} server hosted by GitHub!"
      - run: echo "🔎 The name of your branch is ${{ github.ref }} and your repository is ${{ github.repository }}."
      # 该步骤使用一个actions 官方发布的名为 checkout 的 Action。
      # 这个 Action 用来将指定仓库的代码同步到工作流的 runner 中
      # 只要 runner 中有代码了，后续如何操作，就可以自己随便搞了~
      - name: Check out repository code
        uses: actions/checkout@v2
      - run: echo "💡 The ${{ github.repository }} repository has been cloned to the runner."
      - run: echo "🖥️ The workflow is now ready to test your code on the runner."
      - name: List files in the repository
        run: |
          ls ${{ github.workspace }}
      - run: echo "🍏 This job's status is ${{ job.status }}."
```

- 滚动到页面底部，然后选择 Create a new branch for this commit and start a pull request（为此提交创建一个新分支并开始拉取请求）。 然后，若要创建拉取请求，请单击 Propose new file（提议新文件）。![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537717320-0a2fe106-9eda-4c6f-a81b-6a5837803589.png)

向仓库的分支提交工作流程文件会触发 push 事件并运行工作流程。

## 查看工作流程结果

- 在 GitHub 上，导航到仓库的主页面。
- 在仓库名称下，单击 Actions（操作）。
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537717252-5a465a80-ace7-4a19-b689-c8a145ed90ee.png)
- 在左侧边栏中，单击您想要查看的工作流程。
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537717301-b7808d18-7c4f-40cc-85d4-83ef97121511.png)
- 从工作流程运行列表中，单击要查看的运行的名称。
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537717306-2e079ccf-8130-47fd-9642-f989e7b5fa74.png)
- 在 Jobs（作业）下，单击 Explore-GitHub-Actions 作业。
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537717287-fecb853f-8ee7-4868-81e3-7c843f665bcd.png)
- 日志显示每个步骤的处理方式。 展开任何步骤以查看其细节。
- ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/github_action/1627537718475-e6315bfa-71e1-48e5-9514-16a822265b81.png)

## 更多工作流程模板

GitHub 提供预配置的工作流程模板，您可以自定义以创建自己的持续集成工作流程。 GitHub 分析代码并显示可能适用于您的仓库的 CI 模板。 例如，如果仓库包含 Node.js 代码，您就会看到 Node.js 项目的建议。 您可以使用工作流程模板作为基础来构建自定义工作流程，或按原样使用模板。

您可以在 [actions/starter-workflows](https://github.com/actions/starter-workflows) 仓库中浏览工作流程模板的完整列表。

## 后续步骤

每次将代码推送到分支时，您刚刚添加的示例工作流程都会运行，并显示 GitHub Actions 如何处理仓库的内容。 但是，这只是您可以对 GitHub Actions 执行操作的开始：

- 您的仓库可以包含多个基于不同事件触发不同任务的工作流程。
- 您可以使用工作流程安装软件测试应用程序，并让它们自动在 GitHub 的运行器上测试您的代码。

GitHub Actions 可以帮助您自动执行应用程序开发过程的几乎每个方面。 准备好开始了吗？ 以下是一些帮助您对 GitHub Actions 执行后续操作的有用资源：
