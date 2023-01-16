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

    [root@lichenhao ~]# podman generate systemd generate_test test
    Error: provide only one container name or ID

pod 类型 yaml 生成效果如下：

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

systemd 生成效果如下：

    [root@lichenhao ~]# podman generate systemd generate_test
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

## healthcheck Manage Healthcheck

help Help about any command

history Show history of a specified image

image Manage images

images List images in local storage

import Import a tarball to create a filesystem image

info Display podman system information

init Initialize one or more containers

inspect Display the configuration of a container or image

kill Kill one or more running containers with a specific signal

load Load an image from container archive

login Login to a container registry

logout Logout of a container registry

logs Fetch the logs of a container

mount Mount a working container's root filesystem

pause Pause all the processes in one or more containers

play Play a pod

pod Manage pods

port List port mappings or a specific mapping for the container

ps List containers

pull Pull an image from a registry

push Push an image to a specified destination

restart Restart one or more containers

rm Remove one or more containers

rmi Removes one or more images from local storage

run Run a command in a new container

save Save image to an archive

search Search registry for image

start Start one or more containers

stats Display a live stream of container resource usage statistics

stop Stop one or more containers

system Manage podman

tag Add an additional name to a local image

top Display the running processes of a container

umount Unmounts working container's root filesystem

unpause Unpause the processes in one or more containers

unshare Run a command in a modified user namespace

version Display the Podman Version Information

volume Manage volumes

wait Block on one or more containers
