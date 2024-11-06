---
title: SFTP Subsystem
linkTitle: SFTP Subsystem
date: 2024-11-06T09:53
weight: 20
---

# 概述

> 参考：
>
> - [Manual, sftp-server](https://man.openbsd.org/sftp-server)
> - https://serverfault.com/questions/660160/difference-between-openssh-internal-sftp-and-sftp-server

[OpenSSH](/docs/4.数据通信/Utility/OpenSSH/OpenSSH.md) 的 SFTP 子系统早期由独立的二进制程序 sftp-server 实现（默认路径 /usr/libexec/openssh/sftp-server），后面内嵌到 sshd 中，通过 internal-sftp 参数开启 SFTP 子系统。

OpenSSH 的 SFTP 子系统可以实现 S[FTP](/docs/4.数据通信/Protocol/FTP.md) 能力。

要启用 SFTP 子系统，在 [sshd_config](/docs/4.数据通信/Utility/OpenSSH/OpenSSH%20配置.md) 配置文件中使用 `Subsystem sftp XXX` 指令，XXX 用来指示是用外部二进制程序实现 SFTP 还是内嵌的 sftp 逻辑实现 SFTP。

下面两个指令都可以启动 SFTP 子系统（sftp-server 在新版已经内嵌到 sshd 程序中，更推荐使用 `Subsystem sftp internal-sftp`）

```bash
Subsystem sftp /usr/libexec/openssh/sftp-server
Subsystem sftp internal-sftp
```

# sftp-server

sftp-server 程序的命令参数同样可以在 sshd_config 配置文件中指定，比如：

```bash
Subsystem sftp internal-sftp -f LOCAL5 -l INFO
```

## OPTIONS

**-f LOG_FACILITY** # 指定记录日志时使用的 Facility 代码。`默认值: AUTH`。可用的值有: DAEMON, USER, AUTH, LOCAL0, LOCAL1, LOCAL2, LOCAL3, LOCAL4, LOCAL5, LOCAL6, LOCAL7

**-l LOG_LEVEL** # 指定需要记录哪些级别的日志。`默认值: ERROR`。可用的值有: QUIET, FATAL, ERROR, INFO, VERBOSE, DEBUG, DEBUG1, DEBUG2, DEBUG3。

- INFO 和 VERBOSE 等价，记录客户端的操作信息
- DEBUG 和 DEBUG1 等级；DEBUG2 和 DEBUG3 将会记录比 DEBUG1 更多的日志。

**-u UMASK** # 设置新创建的文件和目录的 umask。（不设置的话则使用用户默认的 umask）

# 最佳实践

## 利用 Match 指令为不同用户设置 ChrootDirectory 并在 rsyslog 记录日志

修改 sshd_config 文件

```bash
Subsystem sftp internal-sftp -f local3 -l INFO

Match User desistdaydream
  ChrootDirectory /sftp_test/desistdaydream/
  ForceCommand internal-sftp -f local4 -l INFO
```

修改 rsyslog.conf 文件

```bash
$AddUnixListenSocket /sftp_test/desistdaydream/dev/log
local3.*                                                /var/log/sftp/sftp.log
local4.*                                                /var/log/sftp/lichenhao.log
```

`mkdir /sftp_test/desistdaydream/dev/` 在 ChrootDirectory 指定的目录下创建 dev/ 目录，以便 rsyslog 在目录下创建名为 log 的 Socket 文件。

# 故障处理

## 配置 ChrootDirectory 后，日志无法被 rsyslog 采集

Google 搜索关键字：ssh sftp ChrootDirectory can't logging (rsyslog)

- https://serverfault.com/questions/710487/sftp-logging-doesnt-work
- https://ubuntuforums.org/showthread.php?t=2081637&p=12342713#post12342713

原因：chroot 后，用户所在的 `/` 目录缺少 [Rsyslog](/docs/6.可观测性/Logs/Rsyslog/Rsyslog.md) 所需的 Socket 类型的 /dev/log 文件

