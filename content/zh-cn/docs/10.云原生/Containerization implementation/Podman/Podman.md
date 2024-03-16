---
title: Podman
---

# 概述

alias docker=podman

podman 与 cri-o 同属 libpod 项目，<https://github.com/containers/podman>

项目地址：<https://github.com/containers/libpod>

# podman 命令行工具

podman \[OPTIONS] COMMAND

OPTIONS

COMMAND

attach Attach to a running container

build Build an image using instructions from Dockerfiles

commit Create new image based on the changed container

container Manage Containers

cp Copy files/folders between a container and the local filesystem

create Create but do not start a container

diff Inspect changes on container's file systems

events Show podman events

exec Run a process in a running container

export Export container's filesystem contents as a tar archive

## generate # 生成结构化数据

通过该命令，可以根据已经运行的容器生成 pod 类型的 yaml 文件或者 systemd 类型的 daemon 文件

Note：仅可对一个容器执行该命令，如果对俩容器执行命令则会报错：

```bash
~]# podman generate systemd generate_test test
Error: provide only one container name or ID`
```

pod 类型 yaml 生成效果如下：

```yaml
[root@lichenhao ~]# podman generate kube generate_test
# Generation of Kubernetes YAML is still under development!
#
# Save the output of this file and use kubectl create -f to import
# it into Kubernetes.
#
# Created with podman-1.4.4
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2020-01-21T03:02:52Z"
  labels:
    app: generatetest
  name: generatetest
spec:
  containers:
  - command:
    - ..... #内容省略
    image: docker.io/lchdzh/network-test:v2.0
    name: generatetest
    resources: {}
    securityContext:
      allowPrivilegeEscalation: true
      capabilities: {}
      privileged: false
      readOnlyRootFilesystem: false
    workingDir: /
status: {}
```

systemd 生成效果如下：

```ini
~]# podman generate systemd generate_test
[Unit]
Description=efa0ab2e648439a516372ecb907f5e506631d033e50978666f032ab5d9ecb788 Podman Container
[Service]
Restart=on-failure
ExecStart=/usr/bin/podman start efa0ab2e648439a516372ecb907f5e506631d033e50978666f032ab5d9ecb788
ExecStop=/usr/bin/podman stop -t 10 efa0ab2e648439a516372ecb907f5e506631d033e50978666f032ab5d9ecb788
KillMode=none
Type=forking
PIDFile=/var/lib/containers/storage/overlay-containers/efa0ab2e648439a516372ecb907f5e506631d033e50978666f032ab5d9ecb788/userdata/efa0ab2e648439a516372ecb907f5e506631d033e50978666f032ab5d9ecb788.pid
[Install]
WantedBy=multi-user.target
```

## healthcheck Manage Healthcheck

