---
title: Job(任务)控制
---

## Job Control # 在同一个 Shell 下，执行多个任务的控制

在前台执行的程序为前台 JOB，前台程序占用一个 shell，执行该程序后，shell 无法再进行别的操作

在后台执行的程序为后台 JOB，后台程序不占用 shell，可以在该 shell 下继续执行其余任务，不受 ctrl+c 的影响

语法格式：

- ctrl+z # 在正在运行的 porcess 中使用这个组合键，可以让前台进程暂停
- fg %JobNumber # 让后台的进程在前台工作
- bg %JobNumber # 让前台的工作在后台工作
- nohup COMMAND # 让命令触发的进程不随 shell 关闭而停止
- COMMAND & # 让命令触发的进程直接在后台运行

**jobs \[OPTIONS] **# 查看 jobs
OPTIONS：

- -l # 查看 jobs 的详细信息
