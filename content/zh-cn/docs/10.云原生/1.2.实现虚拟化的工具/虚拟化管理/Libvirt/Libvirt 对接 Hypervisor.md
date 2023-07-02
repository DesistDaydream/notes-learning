---
title: Libvirt 对接 Hypervisor
---

# 概述

> 参考：
> 
> - [官方文档，QEMU/KVM/HVF hypervisor 驱动程序](https://libvirt.org/drvqemu.html)

Libvirt 支持多不同类型的虚拟化（通常称为 **Drviers(驱动程序)** 或 **Hypervisors(虚拟机监视器)**），因此我们需要一种方法来连接到指定的 Hypervisors。另外，我们还可以通过网络，远程使用其他计算机上的 Hypervisors。

为此，**Libvirt 使用的 RFC 2396 中 定义的 URI 来实现此功能**。就像通过 URL 经由 Browser 访问 Web 资源似的，通过 URI 经由 libvirtd 访问 Hypervisor

> 注意：由于常用 QEMU-KVM 类型虚拟化，所以后文介绍的都以 KVM 虚拟化为主，Xen 类型会单独注明。

使用 Libvirt 客户端编写的代码、virsh 命令行工具、等等，想要连接到 Hypervisor 以管理虚拟机，需要通过一种唯一标识符来找到 Hypervisor，这就是 **URI**；使用 URI 首先需要连接到 **libvirtd 程序**，libvirtd 根据 URI 找到 Hypervisor，然后处理收到的指令以管理虚拟机。

## 连接 Hypervisor 方式

Libvirt 的 KVM/QEMU 驱动程序将会探测 /usr/bin 目录是否存在 `qemu`, `qemu-system-x86_64`, `qemu-system-microblaze`, `qemu-system-microblazeel`, `qemu-system-mips`,`qemu-system-mipsel`, `qemu-system-sparc`,`qemu-system-ppc`。来决定如何连接 QEMU emulator。

Libvirt 的 KVM/QEMU 驱动程序将会探测 /usr/bin 目录是否存在 `qemu-kvm`，以及 /dev/kvm 驱动是否存在。来决定如何连接 KVM hypervisor。

如果使用的 URI 为空，则 libvirt 将使用以下逻辑来确定要使用的 URI

- `LIBVIRT_DEFAULT_URI` 环境变量
- 在客户端配置文件(/etc/libvirt/libvirt.conf)中，`uri_default`关键字的值
- 依次探查每个 hypervisor 程序，直到找到有效的虚拟机监控程序

### 使用 virsh 连接 Hypervisor

这里以 virsh 命令行工具作为示例，其他基于 Libvirt API 的第三方工具，都是同样的道理。

virsh 可以使用 -c 或者 --connect 选项已连接到指定的 libvirtd。比如：`virsh -c qemu+tcp://172.38.180.96/system`

我们可以在客户端配置文件(/etc/libvirt/libvirt.conf)中，设定 `uri_default` 关键字的值以改变 virsh 默认链接到 libvirtd，通常默认值为：`qemu:///system`，即连接本地的 libvirtd 并管理 QEMU 虚拟机。

### 以代码方式通过 Libvirt API 连接 Hypervisor

以 Python 客户端库为例

```python
conn = libvirt.open("qemu+tcp://172.38.180.95/system")
```

# URI

> 参考：
> 
> - [官方文档，连接 URI](https://libvirt.org/uri.html)

这种 URI 格式通常只会由 Libvirt 相关应用程序才可以识别，比如 [virsh 命令行工具](/docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/virsh%20命令行工具/virsh%20命令行工具.md)、[Libvirt 客户端库](docs/10.云原生/1.2.实现虚拟化的工具/虚拟化管理/Libvirt/Libvirt%20API/Libvirt%20客户端库.md)、libvirtd 程序等等。共分为两种 URI

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

## 本地 URI

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

## 远程 URI

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

# libvirtd

> 参考：
> 
> - [官方文档，手册-libvirtd](https://libvirt.org/manpages/libvirtd.html)

libvirtd 程序是 libvirt 虚拟化管理系统的服务器端守护进程组件。

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













