---
title: Drone
---

# 概述

> 参考：
>
> - [GitHub 项目，drone/drone](https://github.com/drone/drone)

Drone 是一个比 Jenkins 更简单易用的现代化 CI/CD 平台。使用 Go 语言编写，基于 Docker 构建，使用简单的 yaml 配置文件来定义和执行 Docker 容器中定义的 Pipeline。

Drone 由两个部分组成：

1. Server # 用于对接 SCM，负责身份认证，SCM 配置，用户、Secrets 以及 webhook 相关的配置。当 Server 收到 SCM 的 webhook 消息后，会通知 Runner 执行 Pipeline。
   1. 可以启动多个 Server 来对接不同的 SCM。
2. Runners # 用于接收任务和运行 Pipeline。如果没有 Runner，那么在触发 Webhook，Drone 在开始 Pipeline 后，会处于 pending 状态并卡住。并且在 Drone Server 的日志中会看到如下报错：
   1. "error": "Cannot transition status via :enqueue from :pending (Reason(s): 状态 cannot transition via "enqueue")",
   2. "msg": "manager: cannot publish status",

Note：

1. Runner 的类型对应 Pipelines 类型，如果 Runner 类型与 Pipelines 类型不对应，则该 Pipelines 则无法执行，并报错。
2. 可以安装多个不同类型的 Runner 来处理不同类型的 Pipelines。

# Drone 部署

注意：Drone 的 Server 端会对接多种 SCM，不同的 SCM 使用不同的运行参数来运行 Server 程序。

由于 Server 与 Runner 的部署截图过多(主要是有对接各种 SCM 的参数，SCM 的配置都在 web 页面上完成的)，篇幅太长，所以另开一篇文档单独记录，详见：[Server 与 Runner 部署](https://www.yuque.com/go/doc/33153134)

# Drone 关联文件

drone.yml # Drone 通过该文件来定义流水线任务(Jenkins 使用的是 Groovy 语言)。参考 Drone 任务示例 文章中的第一个基本演示可以看到如何使用 yaml 来定义一个简单的 pipeline。参考 Drone Pipelines 详解 来学习如何通过一个 yaml 文件来开始一个流水线任务。
