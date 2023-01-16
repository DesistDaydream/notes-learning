---
title: logcli 命令行工具
---

# 概述

> ## 参考：

export LOKI_ADDR=http://localhost:3100

在 <https://github.com/grafana/loki/releases> 该界面下载 logcli 工具的二进制文件，并放到系统 $PATH 下。

二进制文件有了之后，配置一下 logcli 所需要的环境变量 export LOKI_ADDR=http://localhost:3100 然后就可以使用了

下面是一些命令使用示例

```yaml
[root@master-1 bin]# logcli labels job
http://172.38.40.212:31000/loki/api/v1/label/job/values?end=1600402177427090944&start=1600398577427090944
varlogs
[root@master-1 bin]# logcli query '{job="varlogs"}' | more
http://172.38.40.212:31000/loki/api/v1/query_range?direction=BACKWARD&end=1600402187037107678&limit=30&query=%7Bjob%3D%22varlogs%22%7D&start=1600398587037107678
Common labels: {job="varlogs"}
2020-09-18T11:48:50+08:00 {filename="/var/log/host/messages"}          Sep 18 11:48:50 master-1 kubelet: W0918 11:48:50.468511   30889 container.go:526] Failed to update stats for container "/system
.slice/docker-f326688c0b9b38fb8190bba72eb12d55e2017a9624889948ac118e6b9eb1199b.scope": unable to determine device info for dir: /var/lib/docker/overlay2/51f0b901d76af9efd801abb473d0b7d5b27a193ccb990
d3db1cac1799a2a0432/diff: stat failed on /var/lib/docker/overlay2/51f0b901d76af9efd801abb473d0b7d5b27a193ccb990d3db1cac1799a2a0432/diff with error: no such file or directory, continuing to push stat
s
......
```
