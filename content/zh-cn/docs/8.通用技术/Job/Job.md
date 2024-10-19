---
title: Job
linkTitle: Job
date: 2023-11-23T10:59
weight: 1
---

# 概述

> 参考：
>
> -[Wiki, Job_(computing)](https://en.wikipedia.org/wiki/Job_(computing))

在计算中，Job 是一个工作单元或执行单元（执行所述工作）。Job 的组成部分（作为工作单元）称为任务或步骤（如果是连续的，称为 Job Stream）。作为执行单元，Job 可以具体地标识为单个进程，该进程又可以具有子进程，这些子进程执行构成作业的工作的任务或步骤。

# Job 调度

> 参考：
>
> - [Wiki, Job_scheduler](https://en.wikipedia.org/wiki/Job_scheduler)

大多数操作系统（例如 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 和 [Microsoft OS](/docs/1.操作系统/Operating%20system/Microsoft%20OS/Microsoft%20OS.md)）都提供基本的作业调度功能，特别是通过 批处理、cron 和 Windows 任务调度程序。 Web 托管服务通过控制面板或 webcron 解决方案提供作业调度功能。许多程序（例如 DBMS、备份、ERP 和 BPM）还包含相关的作业调度功能。操作系统（“OS”）或点程序提供的作业调度通常不会提供超出单个操作系统实例或超出特定程序范围的调度能力。需要自动化不相关 IT 工作负载的组织还可以利用作业调度程序的更多高级功能

# Job 的衍生

[Cron](/docs/1.操作系统/Linux%20管理/Linux%20系统管理工具/Cron.md) 是类 Unix 系统中的 Job 调度工具

GitHub 的 Action 使用 Workflow(工作流) 称呼多个多个 Job 组成的任务。
