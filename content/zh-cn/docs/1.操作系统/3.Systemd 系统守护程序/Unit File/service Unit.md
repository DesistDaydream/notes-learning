---
title: service Unit
---

# 概述

> 参考：
> 
> - [Manual(手册)，systemd.service(5)](https://man7.org/linux/man-pages/man5/systemd.service.5.html)
> - [金步国 systemd.unit 中文手册，systemd.service 中文手册](https://jinbuguo.com/systemd/systemd.service.html)

所有名称以 `.service` 结尾的 Unit 都是由 Systemd 控制和监督的进程。说白了，就是一个一个的“服务”，这些“服务”就是一个一个的进程。

service Unit 是 systemd 使用数量最多，使用频率最高的单元。

# service 指令

**ExecStart=STRING** # 启动 Unit 所使用的命令。可以使用 `${VarName}` 引用环境变量

**ExecStartPre=STRING** # 启动 Unit 之前执行的命令

**ExecStartPost=STRING** # 启动 Unit 之后执行的命令

**ExecReload=STRING** # 重启 Unit 时执行的命令

**ExecStop=STRING** # 停止 Unit 时执行的命令

**ExecStopPost=STRING** # 停止 Unit 之后执行的命令

**RemainAfterExit=** # 即是 Service 启动的所有进程都退出了，该 Service 是否应该被视为活动状态。`默认值：no`

**RestartSec=INT** # 自动重启当前服务间隔的时间，单位为秒。`默认值：100ms`

**Restart=STRING** # 定义何种情况 Systemd 会自动重启当前服务。`默认值：no`

可用的值有：no(永不重启)、always(无条件重启)、on-success(仅在服务正常退出时重启)、on-failure()、on-abnormal、on-abort、on-watchdog。

下表描述了当由于何种原因退出时，将会执行重启操作的配置。表中有 X 的表示第一行 Restart 的值在第一列列出的退出原因时，将会重启

| 退出原因(↓) \| Restart= (→) | `no` | `always` | `on-success` | `on-failure` | `on-abnormal` | `on-abort` | `on-watchdog` |
| --------------------------- | ---- | -------- | ------------ | ------------ | ------------- | ---------- | ------------- |
| 正常退出                    | X    | X        |              |              |               |            |               |
| 退出码不为"0"               | X    |          | X            |              |               |            |               |
| 进程被强制杀死              | X    |          | X            | X            | X             |            |               |
| systemd 操作超时            | X    |          | X            | X            |               |            |               |
| 看门狗超时                  | X    |          | X            | X            |               | X          |               |

**TimeoutSec=STRING** # 定义 Systemd 停止当前服务之前等待的秒数

**Type=STRING** # 定义启动时的进程行为

它有以下几种值：

- **simple** # 默认值，执行 ExecStart 指定的命令，启动主进程
- **forking** # 以 fork 方式从父进程创建子进程，创建后父进程会立即退出
- **oneshot** # 一次性进程，Systemd 会等当前服务退出，再继续往下执行
- **dbus** # 当前服务通过 D-Bus 启动
- **notify** # 当前服务启动完毕，会通知 Systemd，再继续往下执行
- **idle** # 若有其他任务执行完毕，当前服务才会运行

# 分类

#systemd #unit-file #service
