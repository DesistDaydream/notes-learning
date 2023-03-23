---
title: Containerd 配置详解
weight: 1
---

# 概述

> 参考：
>
> - [Manual(手册),containerd-config.toml(5)](https://github.com/containerd/containerd/blob/main/docs/man/containerd-config.toml.5.md)
> - [Debian Manual](https://manpages.debian.org/bullseye/containerd/containerd-config.toml.5.en.html)

Containerd 使用 [TOML](/docs/2.编程/无法分类的语言/TOML.md) 作为配置文件的格式，默认配置文件为 /etc/containerd/config.toml，我们可以通过命令来生成一个包含所有配置字段的默认配置文件

```bash
mkdir -p /etc/containerd
containerd config default > /etc/containerd/config.toml
```

# 配置文件详解

# \[通用]配置

**version = 2** #
**root = \<STRING>** # Containerd 持久化数据路径。`默认值：/var/lib/containerd`。
**state = \<STRING>** # Containerd 临时数据路径。`默认值：/run/containerd`。
**oom_score = 0** # 设置 Containerd 的 OOM 权重。`默认值：0`。
Containerd 是容器的守护者，一旦发生内存不足的情况，理想的情况应该是先杀死容器，而不是杀死 Containerd。所以需要调整 Containerd 的 `OOM` 权重，减少其被 **OOM Kill** 的几率。最好是将 `oom_score` 的值调整为比其他守护进程略低的值。这里的 oom_socre 其实对应的是 `/proc/<pid>/oom_socre_adj`，在早期的 Linux 内核版本里使用 `oom_adj` 来调整权重, 后来改用 `oom_socre_adj` 了。该文件描述如下：
在计算最终的 `badness score` 时，会在计算结果是中加上 `oom_score_adj` ,这样用户就可以通过该在值来保护某个进程不被杀死或者每次都杀某个进程。其取值范围为 `-1000` 到 `1000`。如果将该值设置为 `-1000`，则进程永远不会被杀死，因为此时 `badness score` 永远返回 0。建议 Containerd 将该值设置为 `-999` 到 `0` 之间。如果作为 Kubernetes 的 Worker 节点，可以考虑设置为 `-999`。

# \[cgroup] 配置

# \[debug] 配置

# \[grpc] 配置表

**address = \<STRING>** # Containerd 监听的 GRPC 路径。`默认值：/run/containerd/containerd.sock`

# \[metrics] 配置

# \[plugins] 配置

详见 [《Plugin 配置》](/docs/IT学习笔记/10.云原生/2.2.实现容器的工具/Containerd/Containerd%20 配置详解/Plugin%20 配置.md 配置详解/Plugin 配置.md) 章节

# \[timeouts] 配置

# \[ttrpc] 配置表

# 配置文件示例

## 镜像加速配置示例

Containerd 的镜像仓库 mirror 与 Docker 相比有两个区别：

- Containerd 只支持通过 `CRI` 拉取镜像的 mirror，也就是说，只有通过 `crictl` 或者 Kubernetes 调用时 mirror 才会生效，通过 `ctr` 拉取是不会生效的。
- `Docker` 只支持为 `Docker Hub` 配置 mirror，而 `Containerd` 支持为任意镜像仓库配置 mirror。

所以需要修改的部分如下：

    [plugins."io.containerd.grpc.v1.cri".registry]
          [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
            [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
              endpoint = ["https://dockerhub.mirrors.nwafu.edu.cn"]
            [plugins."io.containerd.grpc.v1.cri".registry.mirrors."k8s.gcr.io"]
              endpoint = ["https://registry.aliyuncs.com/k8sxio"]
            [plugins."io.containerd.grpc.v1.cri".registry.mirrors."gcr.io"]
              endpoint = ["xxx"]

# Systemd 配置

建议通过 systemd 配置 Containerd 作为守护进程运行，配置文件在上文已经被解压出来了：

```bash
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target local-fs.target

[Service]
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/containerd

Type=notify
Delegate=yes
KillMode=process
Restart=always
RestartSec=5
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNPROC=infinity
LimitCORE=infinity
LimitNOFILE=1048576
# Comment TasksMax if your systemd version does not supports it.
# Only systemd 226 and above support this version.
TasksMax=infinity
OOMScoreAdjust=-999

[Install]
WantedBy=multi-user.target
```

这里有两个重要的参数：

- **Delegate** : 这个选项允许 Containerd 以及运行时自己管理自己创建的容器的 `cgroups`。如果不设置这个选项，systemd 就会将进程移到自己的 `cgroups` 中，从而导致 Containerd 无法正确获取容器的资源使用情况。
- **KillMode** : 这个选项用来处理 Containerd 进程被杀死的方式。默认情况下，systemd 会在进程的 cgroup 中查找并杀死 Containerd 的所有子进程，这肯定不是我们想要的。`KillMode`字段可以设置的值如下。我们需要将 KillMode 的值设置为 `process`，这样可以确保升级或重启 Containerd 时不杀死现有的容器。
  - **control-group**（默认值）：当前控制组里面的所有子进程，都会被杀掉
  - **process**：只杀主进程
  - **mixed**：主进程将收到 SIGTERM 信号，子进程收到 SIGKILL 信号
  - **none**：没有进程会被杀掉，只是执行服务的 stop 命令。
