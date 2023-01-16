---
title: Failed to get system container stats、failed to get cgroup stats
---

> Failed to get system container stats for "/system.slice/docker.service": failed to get cgroup stats for "/system.slice/docker.service": failed to get cgroup stats for "/system.slice/docker.service": failed to get container info for "/system.slice/docker.service": unknown container "/system.slice/docker.service"


参考：[Stackoverflow](https://stackoverflow.com/questions/46726216/kubelet-fails-to-get-cgroup-stats-for-docker-and-kubelet-services)

这个问题大概就是因为 kubelet 在 docker 之前就启动了。
