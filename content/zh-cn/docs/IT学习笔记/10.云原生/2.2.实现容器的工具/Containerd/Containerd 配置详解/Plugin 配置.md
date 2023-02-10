---
title: Plugin 配置
weight: 20
---

# 概述

> 参考：
> - [GitHub 项目文档，containerd/docs/PLUGINS.md](https://github.com/containerd/containerd/blob/main/docs/PLUGINS.md)

## 本篇笔记的记录格式

Containerd 在 TOML 配置文件中，通过 TOML 表的方式来描述一个插件及其具有的功能，效果如下：

- `[plugins."PLUGIN"]`
  - PLUGIN = TYPE.ID
  - TYPE = io.containerd.NAME.VERSION

所以，一个完整描述插件功能的的 TOML 表应该是这样的：

- `[plugins."io.containerd.NAME.VERSION.NAME".NAME....]`

这篇笔记在记录时，则省略前面的通用字符串(`plugins."io.containerd.`)，只以最后的关键字来描述，以获得更好的阅读效果。

比如下文中标题一的 `[grpc.v1.cri]` 下的标题二的 `[registry]` 下的标题三的 `[mirrors]`下的 `docker.io` 镜像仓库的镜像配置，反应到配置文件中，就是这样的：

```toml
[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://ac1rmo5p.mirror.aliyuncs.com"]
```

带 `[]` 的都是一个一个的表，表只是用来进行分组，表中的每一个 `键值对` 才是真实的配置。

# \[gc.v1.scheduler] # 调度器插件

# \[grpc.v1.cri] # CRI 插件

> 参考：
> - [GitHub 项目文档，containerd/docs/cri](https://github.com/containerd/containerd/tree/main/docs/cri)
> - [GitHub 项目文档，containerd/docs/cri/config.md-CRI 插件配置指南](https://github.com/containerd/containerd/blob/main/docs/cri/config.md)

注意：

> - CRI 插件是当 Containerd 作为 CRI 时所使用的配置，所以 ctr、nerdctl 工具在执行某些命令时，有可能不会调用这些配置，就比如其中的 registry 配置，就算配置了，ctr pull 和 nerdctl pull 命令也无法享受到效果。但是使用 crictl 命令是没问题的。

**sanbox_image = <STRING>** # 启动 Pod 时要使用的 Infra 容器。`默认值：k8s.gcr.io/pause:X.X`。这个默认值会根据当前 Containerd 的版本而改变。

## \[cni] # CNI 配置

**bin_dir = <STRING>** # CNI 二进制文件的目录 `默认值：/opt/cni/bin`
**conf_dir = <STRING> **# CNI 配置文件的目录`默认值：/etc/cni/net.d`

## \[containerd] # Containerd 运行时配置

**defautl_runtime_name = <STRING>** # containerd 进程工作时所调用的 runtime。`默认值：runc`

### \[runtimes.runc] # 当 Containerd 使用 runc 作为运行时生效的配置

**cni_conf_dir = <STRING>** # 特定于 runc 作为 runtime 时，所使用的 CNI 配置文件目录
**runtime_type = <STRING>** # 在 containerd 中要使用的 runtime 类型 `默认值：io.containerd.runc.v2`

#### \[options]

**SystemdCgroup = <BOOLEAN>** # 是否使用 systemd cgroup。`默认值：false`

## \[image_decryption]

## \[registry] # 访问镜像注册中心时的配置

> 参考：
> - [GitHub 项目文档，containerd/docs/hosts.md](https://github.com/containerd/containerd/blob/main/docs/hosts.md)

注意：从 Containerd 1.4 版本开始出现的 `registry.configs` 与 `registry.mirrors` 现在(2021 年 4 月)已弃用，只有在未指定 `config_path` 时才会生效

**config_path = <STRING>** # 指定一个目录来引用镜像注册中心的配置`默认值：空`
该目录的格式应该为：`STRING/REGISTRY/hosts.toml`，也就是说，以镜像注册中心的域名作为目录的名称，且目录下的文件名为 `hosts.toml`

假如现在有如下配置：`config_path = "/etc/containerd/registry.d"`，那么 registry.d 目录下的结构应该是下面这样的：

```bash
$ tree /etc/containerd/registry.d
/etc/containerd/registry.d
└── docker.io
    └── hosts.toml

$ cat /etc/containerd/registry.d/docker.io/hosts.toml
server = "https://docker.io"

[host."https://registry-1.docker.io"]
  capabilities = ["pull", "resolve"]
```

### \[configs] # 镜像注册中心的通用配置

**\[REGISTRY]** # 访问 REGISTRY 镜像仓库时的配置。说白了就是发起 HTTP 请求时要设置的那些东西。

- **\[tls]** # TLS 配置
  - **insecure_skip_verify = <BOOLEAN>** # 访问镜像仓库时是否跳过证书验证。`默认值：false`
- **\[auth]** # 发起 HTTP 请求时要使用的认证方式
  - **username = <STRING>** # 访问镜像仓库的用户名
  - **password = <STRING>** # 访问镜像仓库的密码

### \[mirrors] # 镜像注册中心的 mirrors 配置

**\[REGISTRY]** # 为指定的 REGISTRY 镜像仓库配置 mirrors。例如，`[略.registry.mirrors."docker.io"]` 表示配置 docker.io 的 mirror。

- **endpoint = <\[]STRING>** # 表示为 REGISTRY 提供 mirror 的镜像加速服务，是一个数组，可以使用多个镜像加速配置

### 注册中心配置

> 参考：
> - [GitHub 项目文档，containerd/containerd/docs/hosts.md](https://github.com/containerd/containerd/blob/main/docs/hosts.md)

### 配置示例

配置镜像加速
原始

```toml
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://registry-1.docker.io"]
```

hosts.toml 这种配置好像不生效，带验证

```bash
~]# cat registry.d/docker.io/hosts.toml
server = "https://docker.io"

[host."https://ac1rmo5p.mirror.aliyuncs.com"]
  capabilities = ["pull", "resolve"]
```

配置私有镜像仓库
原始

```toml
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."192.168.0.250"]
          endpoint = ["https://192.168.0.250"]
      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        [plugins."io.containerd.grpc.v1.cri".registry.configs."192.168.0.250".tls]
          insecure_skip_verify = true
        [plugins."io.containerd.grpc.v1.cri".registry.configs."192.168.0.250".auth]
          username = "admin"
          password = "Harbor12345"
```

hosts.toml

```bash
~]# cat registry.d/reg.superstor.com/hosts.toml
server = "https://reg.superstor.com"

[host."https://reg.superstor.com"]
  capabilities = ["pull", "resolve", "push"]
  skip_verify = true

```

## \[x509_key_pair_streaming]

# \[internal.v1.opt]

**path = <STRING>** # `默认值：/opt/containerd`

# \[internal.v1.restart]

**interval = <DURATION>** #

# \[metadata.v1.bolt]

# \[monitor.v1.cgroups]

# \[runtime.v2.task] # 运行时 v2 版本插件

**platforms = <\[]STRING>** # `默认值：linux/amd64`
