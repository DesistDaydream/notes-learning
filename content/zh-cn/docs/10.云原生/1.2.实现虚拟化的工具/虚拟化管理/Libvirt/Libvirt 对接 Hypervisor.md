---
title: Libvirt 对接 Hypervisor
---

# 概述

> 参考：
> 
> - [官方文档，连接 URI](https://libvirt.org/uri.html)

Libvirt 支持多不同类型的虚拟化（通常称为 **Drviers(驱动程序)** 或 **Hypervisors(虚拟机监视器)**），因此我们需要一种方法来连接到指定的 Hypervisors。另外，我们可能希望引用网络上远程的 Hypervisors。

为此，**Libvirt 使用的 RFC 2396中 定义的 URI 来实现此功能**

注意：由于常用 QEMU-KVM 类型虚拟化，所以后文介绍的都以 KVM 虚拟化为主，Xen 类型会单独注明。


## URI 格式

```
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


### 本地 URI 格式示例

```
qemu:///system # 连接到系统模式守护程序。
qemu:///session # 连接到会话模式守护程序。
```

如果这样做 libvirtd -help，守护程序将打印出以各种不同方式监听的 Unix 域套接字的路径。


### 远程 URI 格式示例

```
# 通过非标准Unix套接字连接到本地qemu实例（在这种情况下，显式提供了Unix套接字的完整路径）。
qemu+unix:///system?socket=/opt/libvirt/run/libvirt/libvirt-sock
# 连接到提供本地主机端口5000上未加密的TCP / IP连接的libvirtd守护程序，并使用具有默认设置的测试驱动程序。
test+tcp://localhost:5000/default
# 使用与libssh2驱动程序的ssh连接连接到远程主机，并使用其他known_hosts文件。
qemu+libssh2://user@host/system?known_hosts=/home/user/.ssh/known_hosts
# 使用与libssh驱动程序的ssh连接连接到远程主机，并使用其他known_hosts文件。
qemu+libssh://user@host/system?known_hosts=/home/user/.ssh/known_hosts
```

### 测试 URI 格式示例

```
test:///default connects to a default set of host definitions built into the driver.
test:///path/to/host/definitions connects to a set of host definitions held in the named file.
```

## URI 中额外的参数详解

可以将额外的参数作为查询字符串的一部分（后面的部分?）添加到远程URI中。远程URI了解下面显示的其他参数。任何其他内容都未经修改地传递到后端。请注意，参数值必须是 URI转义的。

| 名称 | 运输工具 | 含义 |  |
| --- | --- | --- | --- |
| name | 任何运输 | 传递给远程virConnectOpen函数的名称。该名称通常是通过从远程URI中删除传输，主机名，端口号，用户名和其他参数形成的，但是在某些非常复杂的情况下，最好显式提供名称。 |  |
|  |  |  | 例： name=qemu:///system |
| tls_priority | tls | 无效的GNUTLS优先级字符串 |  |
|  |  |  | 例： tls_priority=NORMAL:-VERS-SSL3.0 |
| mode | Unix，ssh，libssh，libssh2 | auto
自动确定守护程序
direct
连接到每个驱动程序守护程序
legacy
连接到libvirtd
也可以设置libvirt.conf为remote_mode |  |
|  |  |  | 例： mode=direct |
| command | ssh，ext | 外部命令。对于外部运输，这是必需的。对于ssh，默认值为ssh。在PATH中搜索命令。 |  |
|  |  |  | 例： command=/opt/openssh/bin/ssh |
| socket | Unix，ssh，libssh2，libssh | Unix域套接字的路径，它将覆盖已编译的缺省值。对于ssh传输，这将传递到远程netcat命令（请参阅下一个）。 |  |
|  |  |  | 例： socket=/opt/libvirt/run/libvirt/libvirt-sock |
| netcat | ssh，libssh2，libssh | 远程计算机上的netcat命令的名称。默认值为nc。对于ssh传输，libvirt构造一个ssh命令，如下所示：
命令 -p 端口 [-l 用户名 ] 主机名 netcat -U 套接字
其中port，username，hostname可以指定为远程URI的一部分，而command，netcat 和socket来自额外的参数（或合理的默认值）。 |  |
|  |  |  | 例： netcat=/opt/netcat/bin/nc |
| keyfile | ssh，libssh2，libssh | 用于对远程计算机进行身份验证的私钥文件的名称。如果不使用此选项，则使用默认密钥。 |  |
|  |  |  | 例： keyfile=/root/.ssh/example_key |
| no_verify | ssh，tls | SSH：如果设置为非零值，则禁用客户端的严格主机密钥检查，使其自动接受新的主机密钥。现有主机密钥仍将被验证。
TLS：如果设置为非零值，则将禁用客户端对服务器证书的检查。请注意，要禁用服务器对客户机证书或IP地址的检查，必须 更改libvirtd配置。 |  |
|  |  |  | 例： no_verify=1 |
| no_tty | ssh | 如果设置为非零值，如果它无法自动登录到远程计算机（例如，使用ssh-agent等），它将阻止ssh询问密码。当您无权访问终端时（例如在使用libvirt的图形程序中），请使用此选项。 |  |
|  |  |  | 例： no_tty=1 |
| pkipath | tls | 指定客户端的x509证书路径。如果缺少任何CA证书，客户端证书或客户端密钥，则连接将失败并出现致命错误。 |  |
|  |  |  | 例： pkipath=/tmp/pki/client |
| known_hosts | libssh2，libssh | 验证主机密钥所依据的known_hosts文件的路径。尽管LibSSH2不支持所有密钥类型，但LibSSH2和libssh支持OpenSSH风格的known_hosts文件，因此，使用由OpenSSH二进制文件创建的文件可能会导致截断known_hosts文件。因此，对于LibSSH2，建议使用默认的known_hosts文件，该文件位于libvirt的客户端本地配置目录中，例如：〜/ .config / libvirt / known_hosts。注意：使用绝对路径。 |  |
|  |  |  | 例： known_hosts=/root/.ssh/known_hosts |
| sshauth | libssh2，libssh | 用逗号分隔的身份验证方法列表。默认值（是“ agent，privkey，password，keyboard-interactive”。保留方法的顺序。某些方法可能需要其他参数。 |  |
|  |  |  | 例： sshauth=privkey,agent |

# 连接 Hypervisor 方式

官方文档：[https://libvirt.org/drvqemu.html](https://libvirt.org/drvqemu.html)

要连接 QEMU Hypervisor，则必须运行 libvirtd 守护进程(systemctl start libvirtd)，该守护进程的目的是管理 qemu 实例

Libvirt 的 KVM/QEMU 驱动程序将会探测 /usr/bin 目录是否存在`qemu`, `qemu-system-x86_64`, `qemu-system-microblaze`, `qemu-system-microblazeel`, `qemu-system-mips`,`qemu-system-mipsel`, `qemu-system-sparc`,`qemu-system-ppc`。来决定如何连接 QEMU emulator。

Libvirt 的 KVM/QEMU 驱动程序~~将会~~探测 /usr/bin 目录是否存在`qemu-kvm`，以及 /dev/kvm 驱动是否存在。来绝对如何连接 KVM hypervisor。


## 以 libvirt 的 配置文件 或 命令行选项 连接 Hypervisor

这里以 virsh 命令行工具作为示例，其他基于 Libvirt API 的第三方工具，都是同样的道理

1. 使用 -c 或者 --connect 选项。比如：
   1. virsh -c test:///default list
2. 在客户端配置文件(/etc/libvirt/libvirt.conf)中，设定`uri_default`关键字的值


## 以代码方式通过 Libvirt API 连接 Hypervisor

URI 作为 `name` 参数传递给 virConnectOpen 或 virConnectOpenReadOnly 函数。例如：
```
virConnectPtr conn = virConnectOpenReadOnly ("test:///default");
```

如果传递给 virConnectOpen* 的 URI 为NULL，则 libvirt 将使用以下逻辑来确定要使用的URI。

1. `LIBVIRT_DEFAULT_URI` 环境变量
2. 在客户端配置文件(/etc/libvirt/libvirt.conf)中，`uri_default`关键字的值
3. 依次探查每个 hypervisor 程序，直到找到有效的虚拟机监控程序


# 应用实例

## 通过 libvirt 远程管理虚拟机

2台主机：

node4: 192.168.1.166

node5: 192.168.1.143

node4 作为远程libvirt的服务器，上面有已经创建的虚拟机，现在node5上通过以下2种方式管理远程服务器上的虚拟机：

1. 通过qemu+ssh方式
2. 通过qemu+tcp方式

node5上安装 libvirt 及相关工具包，我这里安装了这些，

```
yum groupinstall "Virtualization"
yum install libvirt libvirt-python python-virtinst virt-viewer
```

通过qemu+ssh连接方式比较简单，只需node5能用ssh远程访问node4即可,

命令如下：

```
virsh -c qemu+ssh://root@192.168.1.166/system
```

如果2个节点设置了互信，免密钥登录，可直接执行virsh相关命令，

```
~]# virsh -c qemu+ssh://root@192.168.1.166/system list
 Id    名称                         状态
 ---------------------------------------------------- 
3     vm01                           running
```

下面介绍通过qemu+tcp方式登录远程节点的virsh：

node4上

修改/etc/sysconfig/libvirtd,开启以下2个配置项：

```
~]# egrep -v "^#|^$" /etc/sysconfig/libvirtd
LIBVIRTD_CONFIG=/etc/libvirt/libvirtd.conf
LIBVIRTD_ARGS="--listen
```

修改配置文件，

```
vim /etc/libvirt/libvirtd.conf
listen_tls = 0
listen_tcp = 1
tcp_port ="16509"
listen_addr ="0.0.0.0"
auth_tcp ="none"
```

重启libvirtd并查看监听的端口，

```
# /etc/init.d/libvirtd restart
# netstat -anltp|grep 16509
tcp   0      0 0.0.0.0:16509    0.0.0.0:*      LISTEN      28843/libvirtd
```

node5上远程访问（需要确保可以访问node4的16509 tcp端口）：

```
[root@node5 ~]# virsh -c qemu+tcp://192.168.1.166/system list 
Id    名称                         状态
---------------------------------------------------- 
3     vm01                        running
```


## 配置URI别名 

为了简化管理员的工作，可以在libvirt客户端配置文件中设置URI别名。该配置文件/etc/libvirt/libvirt.conf 适用于root用户或 `${XDG_CONFIG_HOME}/libvirt/libvirt.conf` 任何非特权用户。在此文件中，以下语法可用于设置别名

uri_aliases = [  “ hail = qemu + ssh：//root@hail.cloud.example.com/system”，  “ sleet = qemu + ssh：//root@sleet.cloud.example.com/system”，]

URI别名应该是由字符组成的字符串 a-Z, 0-9, _, -。继= 可以是任何libvirt的URI字符串，包括任意URI参数。URI别名将适用于任何打开libvirt连接的应用程序，除非它已将VIR_CONNECT_NO_ALIASES 参数明确传递给virConnectOpenAuth。如果传入的URI包含允许的别名字符集之外的字符，则不会尝试别名查找。













