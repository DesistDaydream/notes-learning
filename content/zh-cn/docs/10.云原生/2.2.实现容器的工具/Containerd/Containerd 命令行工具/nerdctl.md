---
title: nerdctl
---

# 概述

> 参考：
>
> - [GitHub 项目，containerd/nerdctl](https://github.com/containerd/nerdctl)
> - [官方文档,命令参考](https://github.com/containerd/nerdctl#command-reference)

## Network

nerdctl 本身没有像 docker 的 /etc/docker/daemon.json 这种配置文件，而是使用 CNI 的包 netutil 来执行网络相关的请求， CNI 默认有一个名为 nerdctl0 的 bridge 网络设备，都是常量：
[pkg/netutil/netutil_unix.go](https://github.com/containerd/nerdctl/blob/v0.14.0/pkg/netutil/netutil_unix.go)

```go
package netutil

const (
	DefaultNetworkName = "bridge"
	DefaultID          = 0
	DefaultCIDR        = "10.4.0.0/24"
)

// basicPlugins is used by ConfigListTemplate
var basicPlugins = []string{"bridge", "portmap", "firewall", "tuning"}
```

如果想要像 docker 一样配置网络，则需依赖于 CNI 默认的 /etc/cni/net.d/ 目录中创建配置文件，通常 Containerd 自带的 CNI 配置文件可以在其发布的 [Release](https://github.com/containerd/containerd/releases) 中带 cni 名称的包中找到。效果如下：

```bash
~]# tee /etc/cni/net.d/10-containerd-net.conflist <<-"EOF"
{
  "cniVersion": "0.4.0",
  "name": "containerd-net",
  "plugins": [
    {
      "type": "bridge",
      "bridge": "cni0",
      "isGateway": true,
      "ipMasq": true,
      "promiscMode": true,
      "ipam": {
        "type": "host-local",
        "ranges": [
          [{
            "subnet": "10.88.0.0/16"
          }],
          [{
            "subnet": "2001:4860:4860::/64"
          }]
        ],
        "routes": [
          { "dst": "0.0.0.0/0" },
          { "dst": "::/0" }
        ]
      }
    },
    {
      "type": "portmap",
      "capabilities": {"portMappings": true}
    }
  ]
}
EOF
```

此时，通过 `nerdctl network list` 命令即可看到一个名为 containerd-net 新的网络

```bash
]# nerdctl network ls
NETWORK ID    NAME              FILE
0             bridge
              containerd-net    /etc/cni/net.d/10-containerd-net.conflist
              host
              none
```

然后运行容器时，使用 `--net=containerd-net` 参数指定该网络，所有运行的容器，即可关联到指定的 containerd-net 网桥上。
注意：nerdctl 默认的网络是无法修改的

## 现存问题

v0.12.1 版本

通过 build 构建完镜像后，会产生一个相同 ID 无名的空镜像

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ws1t24/1631632584319-b6131f0b-6269-422d-a203-045ab0b2538f.png)

执行 nerdctl rmi 命令时，tab 无法补全，但是 `nerctl image rm` 可以补全，但是无法删除那两个空镜像

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ws1t24/1631632732271-b95cef43-e60b-4fd2-842a-1f0a80cb5dac.png)

# nerdctl 关联文件与配置

**/var/lib/nerdctl/\*** # nerdctl 的默认 DataRoot，即 nerdctl 运行容器时所产生的文件保存路径。

**/etc/nerdctl/nerdctl.toml** # nerdctl 的配置文件，配置文件中的内容通常与命令行标志可以对应。该文件可以通过 `${NERDCTL_TOML}` 改变位置。该配置文件与 /etc/containerd/config.toml 文件无关。

- **~/.config/nerdctl/nerdctl.toml** # Rootless 运行时的配置文件
- 参考：<https://github.com/containerd/nerdctl/blob/main/docs/config.md>

常见基本配置

```
mkdir -p /etc/nerdctl
tee /etc/nerdctl/nerdctl.toml > /dev/null <<EOF
address        = "unix:///run/k3s/containerd/containerd.sock"
namespace      = "k8s.io"
EOF
```

# Syntax(语法)

**nerdctl \[Global OPTIONS] COMMAND \[COMMAND OPTIONS] \[ARGUMENTS......]**

**Global OPTIONS**

nerdctl 除了可以通过全局选项改变运行行为，还可以通过环境变量改变。凡是可以通过环境变量指定的全局选项，都会有特殊说明

- **--aaddress, -a, --host, -H \<PATH>** # 指定容器地址。`默认值：/run/containerd/containerd.sock`。可以使用 `unix://` 前缀。
    - 环境变量：$CONTAINERD_ADDRESS
- **-n, --namespace \<STRING>** # 指定容器名称空间。`默认值：default`。通过 docker 运行的在 moby 名称空间中，通过 Kubernetes 运行的容器在 k8s.io 名称空间中。
    - 环境变量：$CONTAINERD_NAMESPACE
- **--cni-path \<PATH>** # 指定 CNI 插件所需的二进制文件所在目录。`默认值：/opt/cni/bin`
    - 环境变量：$CNI_PATH
- **--cni-netconfpath \<PATH>** # 指定 CNI 配置文件所在目录。`默认值：/etc/cni/net.d`
    - 环境变量：$NETCONFPATH
- **--data-root \<PATH>** # nerdctl 持久化数据所在目录。`默认值：/var/lib/nerdctl`。该目录由 nerdctl 管理，而不是 containerd。
- **--cgroup-manager \<STRING>** # 指定 nerdctl 要使用的 Cgroup 管理器。`默认值：cgroupfs`
    - 可用的值有： cgroupfs、systemd
- **--insecure-registry \<BOOLEAN>** # 是否跳过 HTTPS 证书的验证行为，并允许回退到纯 HTTP。`默认值：false`

# Management Commands (管理命令)

management command 在使用的时候，当后面还需要跟其子命令的时候，是可省的。直接使用子命令就表示对其执行，但是有的管理命令不行，比如 create，对于 container 可省，对于 network 不可省

# Other Commands

## inspect

显示 nerdctl 所能管理的所有对象的详细信息，包括 image、container、network、volume 等等

### Syntax(语法)

**nerdctl inspect \[OPTIONS] NAME|ID \[NAME|ID...]**

**OPTIONS**

- **--mode=\<STRING>** # 显示模式。`默认值：dockercompat`。可用的值有：dockercompat、native。
  - native # 显示更多信息
