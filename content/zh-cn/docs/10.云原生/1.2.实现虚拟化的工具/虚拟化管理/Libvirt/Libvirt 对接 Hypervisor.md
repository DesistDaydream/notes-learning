---
title: Libvirt 对接 Hypervisor
---

# 概述

> 参考：
> 
> - [官方文档，驱动程序](https://libvirt.org/drivers.html)

Libvirt 支持多不同类型的虚拟化平台，想要对接这些虚拟化平台并管理它们，需要 **Drivers(驱动程序)**。通过 Drivers 可以连接各种虚拟化平台的 **Hypervisors(虚拟机监视器)**，可以是本地或网络上的远程 Hypervisors。

这些 Drivers 通常都暴露了了 [Libvirt API](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/Libvirt%20API.md)。使用 Libvirt API 的应用程序(e.g. [virsh](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20命令行工具/virsh%20命令行工具.md)、virt-manager、等等)使用 **RFC 2396 定义的 URI** 以连接指定的 Driver。

Libvirt 有如下几类 Drivers

- Hypervisor drivers
- Storage drivers
- Node device driver
- Secret driver

~~一般情况应该始终有一个活跃状态的 Hypervisor driver，如果 libvirtd 进程可用的话，通常还会有一个活动状态的网络驱动和存储驱动~~

除了 Hypervisor 驱动有用以外，其他的几种驱动暂时没找到用途 —— 2023.7.5

libvirtd 可以算是一种 Driver。

我们最常见的 virsh 和 libvirtd 就是一个标准的客户端和服务端组合程序，virsh 是用 LIbvirt API 实现的客户端程序、libvirtd 则是暴露 LIbvirt API 的驱动程序。virsh 使用 URI 连接到 libvirtd 的指定驱动后，可以像 libvirtd 发出命令以管理虚拟化平台。

# URI

> 参考：
> 
> - [官方文档，连接 URI](https://libvirt.org/uri.html)

我们使用客户端程序（e.g. virsh、virt-manager、等等）**通过 URI 连接 Driver**。驱动程序通常都暴露了 Libvirt API，连接 Driver 的代码逻辑如下：

```c
virConnectPtr conn = virConnectOpen ("test:///default");
```

实际上，就是将 URI 作为 virConnectOpen() 函数的参数，以生成一个连接到 Driver 的实例。从代码角度看，这就好像连接 MySQL 需要一个 URI、连接 Redis 需要一个 URI、连接 ETCD 也需要一个。

如果用不用代码的方式描述，其实就是 `virsh -c URI` 这种，就类似 `mysql -h X.X.X.X`、`etcdctl --endpoint X.X.X.X` 之类的。只不过 Libvirt 的应用程序还可以使用 `qemu:///system` 这种方式连接本地 Socket，而不止是常见的 TCP/IP。

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

这两种 URI 其实都符合如下语法：

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

### 本地 URI

其中本地 URI 通常都会省略所有 `[]` 部分，这样就会看到 3 个连续的 `/` 符号，就像 `qemu:///system`，这里的 Transport 就是默认的 unix，这就是本地 URI 了。

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

## QEMU/KVM 驱动

> 参考：
> 
> - [官方文档，QEMU/KVM/HVF hypervisor 驱动程序](https://libvirt.org/drvqemu.html)

Libvirt 的 KVM/QEMU 驱动程序将会探测 `/usr/bin/` 目录是否存在 `qemu`, `qemu-system-x86_64`, `qemu-system-microblaze`, `qemu-system-microblazeel`, `qemu-system-mips`,`qemu-system-mipsel`, `qemu-system-sparc`,`qemu-system-ppc`。来决定如何连接 QEMU emulator。

Libvirt 的 KVM/QEMU 驱动程序将会探测 /usr/bin 目录是否存在 `qemu-kvm`，以及 /dev/kvm 驱动是否存在。来决定如何连接 KVM hypervisor。

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

用 `qemu:///embed?root=/some/dir` 这种 URI 可以使用 virsh 的内置 QEMU/KVM 驱动，并将运行数据保存到 /some/dir/ 目录中，改目录的结构如下：

```bash
/some/dir
  |
  +- log
  |   |
  |   +- qemu
  |   +- swtpm
  |
  +- etc
  |   |
  |   +- qemu
  |   +- pki
  |       |
  |       +- qemu
  |
  +- run
  |   |
  |   +- qemu
  |   +- swtpm
  |
  +- cache
  |   |
  |   +- qemu
  |
  +- lib
      |
      +- qemu
      +- swtpm
```

# libvirtd

> 参考：
> 
> - [官方文档，手册-libvirtd](https://libvirt.org/manpages/libvirtd.html)

libvirtd 程序是 libvirt 虚拟化管理系统的服务器端守护进程组件。内置了部分 Hypervisor 驱动，并暴露了 [Libvirt API](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/Libvirt%20API.md)。

该守护进程在主机服务器上运行，并为虚拟来宾执行所需的管理任务。这包括启动、停止和在主机服务器之间迁移来宾、配置和操作网络以及管理供来宾使用的存储等活动。

libvirt 客户端库和实用程序连接到此守护进程以发出任务并收集有关主机系统和来宾的配置和资源的信息。

默认情况下，libvirtd 守护进程侦听本地 Unix 域套接字上的请求。使用 -l | --listen 命令行选项，可以指示 libvirtd 守护进程另外侦听 TCP/IP 套接字。要使用的 TCP/IP 套接字在 libvirtd 配置文件中定义。

> 这里官方说的其实有一些问题，参考 https://stackoverflow.com/questions/65663825/could-not-add-the-parameter-listen-to-open-tcp-socket ，并且在下面关于启动模式中也有说明。

重新启动 libvirtd 不会影响正在运行的 guest 虚拟机。如果定义了 XML 配置，来宾将继续操作并将被自动接听。任何尚未定义 XML 配置的来宾都将从配置中丢失。

## libvirtd 的启动模式

libvirtd 守护进程能够以两种模式启动

- **传统模式** # 它将自行创建并侦听 UNIX 套接字。
  - 如果给出了 --listen 参数，它还将根据 /etc/libvirt/libvirtd.conf 中的 listen_tcp 和 listen_tls 选项监听 TCP/IP 套接字
- **套接字激活模式** # 它将依靠 systemd 在 UNIX 和可选的 TCP/IP 套接字上创建和侦听，并将它们作为预打开的文件描述符传递。
  - <font color="#ff0000">注意</font>：在这种模式下，不允许传递 --listen 参数，并且 /etc/libvirt/libvirtd.conf 中大多数与套接字相关的配置选项将不再起作用。
  - 如果想要启用 TCP 或 TLS 套接字，可以开启 libvirtd-tcp.socket 或 libvirtd-tls.socket 这两个 Unit。

在使用 systemd 的主机操作系统上运行时，套接字激活模式通常是默认模式。要恢复到传统模式，必须使用如下命令屏蔽所有套接字单元文件：

```bash
$ systemctl mask libvirtd.socket libvirtd-ro.socket \
   libvirtd-admin.socket libvirtd-tls.socket libvirtd-tcp.socket
```

最重要的是，请确保 --timeout 参数不用于守护进程，因为它不会在以后的任何连接中启动。

如果使用 libvirt-guests 服务，则需要调整该服务的排序，以便将其排序在服务单元而不是套接字单元之后。由于依赖项和顺序无法通过直接覆盖来更改，因此需要更改整个 libvirt-guests 单元文件。为了保留此类更改，请将已安装的 /usr/lib/systemd/system/libvirt-guests.service 复制到 /etc/systemd/system/libvirt-guests.service 并在那里进行更改，特别确保 After= 排序提及libvirtd.service 而不是 libvirtd.socket：

```ini
[Unit]
After=libvirtd.service
```

# 连接驱动的方式实例

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

# 应用实例

## 通过 libvirt 远程管理虚拟机

### 通过 TCP 连接

> 参考：
> 
> - [StackOverflow，could-not-add-the-parameter-listen-to-open-tcp-socket](https://stackoverflow.com/questions/65663825/could-not-add-the-parameter-listen-to-open-tcp-socket)
> - [libvirtd 官方手册](https://libvirt.org/manpages/libvirtd.html)

`systemctl stop libvirtd.service`

在 `/etc/libvirt/libvirtd.conf` 文件中添加 `auth_tcp="none"`

让 libvirtd 监听本地 TCP 端口

- `systemctl enable libvirtd-tcp.socket --now`

`systemctl start libvirtd.service`

最后使用 `virsh -c qemu+tcp://192.168.1.66/system` 即可连接到远程 libvirtd

### 通过 SSH 连接

```
virsh -c qemu+ssh://root@192.168.1.166/system
```

## 配置URI别名 

为了简化管理员的工作，可以在 libvirt 客户端配置文件中设置 URI 别名。对于 root 用户，配置文件为 `/etc/libvirt/libvirt.conf`；对于任何非特权用户，配置文件为 `$XDG_CONFIG_HOME/libvirt/libvirt.conf`。在此文件中，可以使用以下语法来设置别名

```bash
uri_aliases = [
  "hail=qemu+ssh://root@hail.cloud.example.com/system",
  "sleet=qemu+ssh://root@sleet.cloud.example.com/system",
]
```

URI 别名应该是由字符 `a-Z`、`0-9`、`_`、`-` 组成的字符串。 `=` 后面可以是任何 libvirt URI 字符串，包括任意 URI 参数。 URI 别名将应用于任何打开 libvirt 连接的应用程序，除非它已显式地将 VIR_CONNECT_NO_ALIASES 参数传递给 virConnectOpenAuth。如果传入的 URI 包含允许的别名字符集之外的字符，则不会尝试别名查找。













