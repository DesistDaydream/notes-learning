---
title: "Job control"
linkTitle: "Job control"
weight: 20
---

# 概述

> 参考：
> 
> - [Manual(手册)，Bash(1) - JOB_CONTROL](https://www.man7.org/linux/man-pages/man1/bash.1.html#JOB_CONTROL)

Job control(作业控制) 是指有选择的**停止/挂起**进程的执行，并在稍后**继续/恢复**它们的执行能力。我们一般都是在 Shell 中使用此功能的，比如 Bash。

在前台执行的程序为前台 JOB，前台程序占用一个 shell，执行该程序后，shell 无法再进行别的操作

在后台执行的程序为后台 JOB，后台程序不占用 shell，可以在该 shell 下继续执行其余任务，不受 ctrl+c 的影响

常见操作：

- **ctrl+z** # 在正在运行的 porcess 中使用这个组合键，可以让前台进程暂停
- **fg %JobNumber** # 让后台的进程在前台工作
- **bg %JobNumber** # 让前台的工作在后台工作
- **nohup COMMAND** # 让命令触发的进程不随 shell 关闭而停止
- **COMMAND &** # 让命令触发的进程直接在后台运行

# jobs 命令

## Syntax(语法)

查看 jobs

**jobs [OPTIONS]**

OPTIONS：

- **-l** # 查看 jobs 的详细信息
