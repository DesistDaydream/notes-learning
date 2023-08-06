---
title: "Driver"
linkTitle: "Driver"
date: "2023-07-06T08:36"
weight: 1
---

# 概述

> 参考：
>
> - [官方文档，驱动程序](https://libvirt.org/drivers.html)

市面上有多种虚拟化平台，比如 [KVM/QEMU](/docs/10.云原生/1.2.实现虚拟化的工具/KVM_QEMU/KVM_QEMU.md)、[Hyper-V](/docs/10.云原生/1.2.实现虚拟化的工具/Hyper-V/Hyper-V.md)、等等，Libvirt 想要调用这些虚拟化平台的能力，需要对应平台的 **Driver(驱动程序)**，这个 Driver 可以对接虚拟化平台的的 **Hypervisor(虚拟机监视器)** 以控制整个虚拟化环境。这就好比 Windows 系统想要使用显卡的能力，就需要对应的显卡驱动程序一样。

想要连接到 Driver，我们需要使用 Libvirt API 开发的客户端应用程序（e.g. [virsh](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20命令行工具/virsh%20命令行工具.md)、[virt-manager](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/virt-manager.md)、等等）。Drivers 通常作为服务端都暴露了 [Libvirt API](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/Libvirt%20API.md)，这些客户端通过 [**URI**](#URI) 找到并连接到 Driver，这就好像 mysql 客户端连接 mysql 也需要 IP 端口、etcdctl 连接 etcd 同理，很多客户端都是同样的逻辑。

Libvirt 有如下几类 Drivers

- **Hypervisor drivers(Hypervisor 驱动)**
- **Storage drivers(存储驱动)**
- **Node device driver**
- **Secret driver**

~~一般情况应该始终有一个活跃状态的 Hypervisor driver，如果 libvirtd 进程可用的话，通常还会有一个活动状态的网络驱动和存储驱动~~

> 除了 Hypervisor 驱动有用以外，其他的几种驱动暂时没找到用途 —— 2023.7.5

使用 Libvirt 时，我们最常见的 virsh 和 libvirtd 就是一个客户端和服务端结构，virsh 是用 LIbvirt API 实现的客户端程序、libvirtd 则是暴露 LIbvirt API 的驱动程序。virsh 使用 URI 连接到 libvirtd 的指定驱动后，可以像 libvirtd 发出命令以管理虚拟化平台。

# URI

> 参考：
>
> - [官方文档，连接 URI](https://libvirt.org/uri.html)

Libvirt 的客户端程序（e.g. virsh、virt-manager、等等）**通过 [URI](https://datatracker.ietf.org/doc/html/rfc2396) 连接 Driver**。这就好像 mysql 客户端连接 MySQL 需要一个 URL、redis 客户端连接 Redis 需要、etcdctl 连接 ETCD 也需要、等等。其实就是 `virsh -c URI` 这种，类似 `mysql -h X.X.X.X`、`etcdctl --endpoint X.X.X.X` 之类的。只不过 Libvirt 的应用程序还可以使用 `qemu:///system` 这种方式连接本地 Socket，而不止是常见的 TCP/IP。

驱动程序通常都暴露了 Libvirt API，连接 Driver 的代码逻辑如下：

```c
virConnectPtr conn = virConnectOpen ("test:///default");
```

实际上，就是将 URI 作为 virConnectOpen() 函数的参数，以生成一个连接到 Driver 的实例。

## URI 语法

```text
Hypervisor[+Transport]://[UserName@][HostName][:PORT]/PATH[?Extraparameters]
```

Hypervisor # 虚拟化类型，可用的值有：

- qemu
- xen
- test # 专门测试的，libvirt 自带

Transport # 连接方式。默认为 unix

- unix
- ssh
- tcp
- libssh、libssh2
- auto
- netcat
- native
- ext
- tls

## 默认 URI

如果我们传递给 `virConnectOpen` 函数的参数为空，即 URI 为空，那么 Driver 通常使用如下逻辑来确定要使用的 URI

- 环境变量 LIBVIRT_DEFAULT_URI
- 客户端配置文件 uri_default 参数
- 依次探查每个虚拟机监控程序，直到找到有效的虚拟机监控程序

对于 virsh 这个客户端程序来说，可以使用配置文件（/etc/libvirt/libvirt.conf）中的 uri_default 参数、还可以使用 -c 选项覆盖这个配置。

## URI 种类

共分为两种 URI

- 本地 URI
- 远程 URI

### 本地 URI

其中本地 URI 语法通常都会省略所有 `[]` 部分，这样就会看到 3 个连续的 `/` 符号，就像 `qemu:///system`，这里的 Transport 就是默认的 unix，这就是本地 URI 了。

URI 示例

连接到系统模式守护程序

- `qemu:///system`

连接到会话模式守护程序

- `qemu:///session`

如果这样做 `libvirtd --help`，守护程序将打印出以各种不同方式监听的 Unix 域套接字的路径。

本地使用测试 URI

连接到驱动程序内置的一组默认主机定义

- `test:///default`

连接到指定文件中保存的一组主机定义

- `test:///PATH/TO/HOST/DEFINITIONS`

### 远程 URI

远程 URI 的语法有一些是不能省略的，其中 Transport、HostName 必须指定

```
Hypervisor+Transport://[UserName@]HostName[:PORT]/PATH[?Extraparameters]
```

URI 示例

通过 ssh 以 root 用户身份连接

- `qemu+ssh://root@172.38.180.96/system`

通过 tcp 连接。默认连接到 16509 端口，前提参考 [通过 TCP 连接](#通过%20TCP%20连接) 部分的内容为 libvirtd 开启 TCP 监听。

`qemu+tcp://172.38.180.95/system`

## URI 中额外的参数详解


# Hypervisor 驱动程序

Hypervisor Driver 有几种表现形式

- 内嵌在 Libvirt API 客户端应用程序中
- 内置到 libvirt.so 库
- [Libvirt 守护进程](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Driver/Libvirt%20守护进程.md)，比如 libvirtd
- 等等

## QEMU/KVM 驱动

> 参考：
>
> - [官方文档，QEMU/KVM/HVF hypervisor 驱动程序](https://libvirt.org/drvqemu.html)

Libvirt 的 KVM/QEMU 驱动程序将会探测 `/usr/bin/` 目录是否存在 `qemu`, `qemu-system-x86_64`, `qemu-system-microblaze`, `qemu-system-microblazeel`, `qemu-system-mips`,`qemu-system-mipsel`, `qemu-system-sparc`,`qemu-system-ppc`。来决定如何连接 QEMU emulator。

Libvirt 的 KVM/QEMU 驱动程序将会探测 `/usr/bin` 目录是否存在 `qemu-kvm`，以及 /dev/kvm 驱动是否存在。来决定如何连接 KVM hypervisor。

libvirt QEMU 驱动程序是一个**多实例**驱动程序，提供单个系统范围的特权驱动程序（“system”实例）和每用户非特权驱动程序（“session”实例）。 URI 驱动程序协议是“qemu”。 libvirt 驱动程序的一些连接 URI 示例如下：

```bash
qemu:///session                      (local access to per-user instance)
qemu+unix:///session                 (local access to per-user instance)

qemu:///system                       (local access to system instance)
qemu+unix:///system                  (local access to system instance)
qemu://example.com/system            (remote access, TLS/x509)
qemu+tcp://example.com/system        (remote access, SASl/Kerberos)
qemu+ssh://root@example.com/system   (remote access, SSH tunnelled)
```

通常，对于本地连接 QEMU/KVM 驱动来说，默认都是 `qemu:///system`，即连接到 system 实例以管理虚拟机。

### 嵌入式驱动

https://libvirt.org/drvqemu.html#embedded-driver

从 6.1.0 版本开始，可以使用 virsh 的嵌入式驱动。

用 `qemu:///embed?root=/some/dir` 这种 URI 可以使用 virsh 的内置 QEMU/KVM 驱动，并将运行时数据保存到 /some/dir/ 目录中，改目录的结构如下：

```bash
/some/dir
  |
  +- log
  |   |
  |   +- qemu
  |   +- swtpm
  +- etc
  |   |
  |   +- qemu
  |   +- pki
  |       |
  |       +- qemu
  +- run
  |   |
  |   +- qemu
  |   +- swtpm
  +- cache
  |   |
  |   +- qemu
  +- lib
      |
      +- qemu
      +- swtpm
```

# 连接驱动的方式示例

### 使用 virsh 连接 Driver

这里以 virsh 命令行工具作为示例，其他基于 Libvirt API 的第三方工具，都是同样的道理。

virsh 可以使用 -c 或者 --connect 选项连接到指定的 libvirtd。比如：`virsh -c qemu+tcp://172.38.180.96/system`

我们可以在客户端配置文件(/etc/libvirt/libvirt.conf)中，设定 `uri_default` 关键字的值以改变 virsh 默认链接到 libvirtd，通常默认值为：`qemu:///system`，即连接本地的 libvirtd 并管理 QEMU 虚拟机。

### 以代码方式连接 Driver

Python 客户端库

```python
conn = libvirt.open("qemu+tcp://172.38.180.95/system")
```

Go 客户端库

```go
conn, err := libvirt.NewConnect("qemu:///system")
```

