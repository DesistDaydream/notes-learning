---
title: Redhat 优化
---

# 概述
Redhat 官方文档：<https://access.redhat.com/documentation/zh-CN/Red_Hat_Enterprise_Linux/7/html/Performance_Tuning_Guide/>

参考文章：<https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/monitoring_and_managing_system_status_and_performance/index>

# Tuned 概述

Tuned 是一项服务，可监视您的系统并在某些工作负载下优化性能。Tuned 的核心是配置文件，它们可以针对不同的用例调整您的系统。

针对许多用例，已分发 Tuned 并附带了许多预定义的配置文件：

- High throughput

- Low latency

- Saving power

可以修改为每个配置文件定义的规则，并自定义如何调整特定设备。当您切换到另一个配置文件或停用 Tuned 时，以前的配置文件对系统设置所做的所有更改都将恢复为原始状态。

您还可以配置“Tuned”以对设备使用情况的变化做出反应，并调整设置以提高活动设备的性能并降低非活动设备的功耗。

## 配置文件

/usr/lib/tuned/_ # 特定于发行版的概要文件存储在目录中。每个配置文件都有其自己的目录。该配置文件由名为 tuned.conf 的主要配置文件以及其他文件（例如帮助程序脚本）组成。
/etc/tuned/_ # 如果需要定制概要文件，请将概要文件目录复制到用于定制概要文件的目录中。如果有两个同名的配置文件，则使用 /etc/tuned/ 中的自定义配置文件。

tuned-adm 命令行工具

usage: tuned-adm \[-h] \[--version] \[--debug] \[--async] \[--timeout TIMEOUT]

             [--loglevel LOGLEVEL]

             {list,active,off,profile,profile_info,recommend,verify,auto_profile,profile_mode}

             ...

positional arguments:

1. {list,active,off,profile,profile_info,recommend,verify,auto_profile,profile_mode}

2. list list available profiles or plugins (by default profiles)

3. active show active profile

4. off switch off all tunings

5. profile switch to a given profile, or list available profiles if no profile is given

6. profile_info show information/description of given profile or current profile if no profile is specified

7. recommend recommend profile

8. verify verify profile

9. auto_profile enable automatic profile selection mode, switch to the recommended profile

10. profile_mode show current profile selection mode

optional arguments:

1. -h, --help show this help message and exit

2. \--version, -v show program's version number and exit

3. \--debug, -d show debug messages

4. \--async, -a with dbus do not wait on commands completion and return immediately

5. \--timeout TIMEOUT, -t TIMEOUT with sync operation use specific timeout instead of the default 600 second(s)

6. \--loglevel LOGLEVEL, -l LOGLEVEL level of log messages to capture (one of debug, info,warn, error, console, none). Default: console
