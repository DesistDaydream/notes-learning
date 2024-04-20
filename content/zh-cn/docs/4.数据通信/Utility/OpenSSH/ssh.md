---
title: ssh
linkTitle: ssh
date: 2024-04-19T14:59
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册), ssh(1)](https://man.openbsd.org/ssh)
> - [Manual(手册), ssh_config(5)](https://man.openbsd.org/ssh_config)


# Syntax(语法)

**ssh \[OPTIONS] \[User@] HostIP \[COMMAND]**

## OPTIONS

- **-C** # 请求压缩所有数据(包括 stdin、stdout、stderr、X11 转发的数据、TCP 和 UNIX-domain 连接)
- **-i \</PATH/FILE>** # 使用指定的私钥文件来进行登录认证
- **-J \<DEST\[,DEST2,...]>** # 首先与 DEST 建立 ssh 连接，并通过 DEST 跳转到最终目标主机。如果需要多次跳转。可以指定多个 DEST 并以逗号分割。
  - DEST 格式为：`[USER@]HOST[:PORT]`
- **-o \<KEY=VALUE>** # 以命令行的方式配置本应该在 [OpenSSH 配置](/docs/4.数据通信/Utility/OpenSSH/OpenSSH%20配置.md) - ssh_config 文件中的内容。KEY 是 ssh_config 配置文件中的关键字。
- **-p \<PORT>** # 指定 HostIP 所在远程服务器监听的端口
- **-T** # 不要分配一个伪终端
- **-W** # 请求将客户端上的标准输入和输出通过安全通道转发到端口上的主机。 表示 -N，-T，ExitOnForwardFailure 和 ClearAllForwardings，尽管可以在配置文件中或使用 -o 命令行选项覆盖它们。
- **-X** # 启用 X11 转发
- **-Y** # 启用受信任的 X11 转发。 受信任的 X11 转发不受 X11 SECURITY 扩展控件的约束。

### 隧道选项

**-w LOCAL_TUN\[:REMOTE_TUN]** # 在客户端与服务端建立互相连接的 tun 设备. tun 设备的默认模式为 **point-to-point**.

- LOCAL_TUN 和 REMOTE_TUN 可以是 数字 或者 any 关键字. e.g. `-w 0:0` 则会创建名为 tun0 的设备
- Notes:
  - 必须使用具有创建网络设备权限的用户执行 ssh 命令, 比如 root 用户.
  - 必须保证 [OpenSSH 配置](/docs/4.数据通信/Utility/OpenSSH/OpenSSH%20配置.md) 中 sshd_config 文件中的 PermitTunnel 的值不为 no, 才能保证在两端创建 tun 设备, 并让 tun 设备互联.
  - 可以通过 ssh_config 中的 Tunnel 指令改变 tun 设备的类型. e.g. `ssh -p 20022 -o "Tunnel=ethernet" -w 0:0 root@1.1.1.1` 将会创建 tap 设备

### 端口转发选项

**-D \<\[Bind_Address]:PORT>** # Dynamic(动态) 转发。启用动态转发的 ssh 程序相当于一个代理服务，通过监听的端口，可以将流量送到指定的目标主机。

**-L \<XXX>** # Local(本地) 转发。发往本地的 TCP 端口 或 Unix Socket 上的流量转发到远端 TCP 端口 或 Unix Socket 上。

- XXX 有多种语法格式：
- **-L \[Bind_Address:]LocalPort:RemoteHost:RemoteHostPort**
- **-L \[Bind_Address:]LocalPort:RemoteSocket**
- **-L LocalSocket:RemoteHost:RemotePort **
- **-L LocalSocket:RemoteSocket**
- 假如现在有 A、B、C 三台主机，A 与 C 不通，A 与 B 通，C 与 B 通；也就是说 B 是中转站(运行 sshd 程序)。如果想要通过 A 访问 C，则需要在 A 上执行 `ssh -L XXX B-IP` 命令
- Local 表示 A 主机，Remote 表示 C 主机
- 访问 A Port 就是访问 C Port

**-R \<XXX>** # Remote(远程) 转发。发往指定 `RemotePort` 或 `RemoteSocket` 上的流量转发到 `本地端口` 或 `Unix Socket` 上。远程转发其实更像将 Local 服务通过 ssh 程序以类似 nat 的方式暴露到 Remote 上。

- XXX 有多种语法格式：
- **-R \[bind_address:]RemotePort:LocalHost:LocalPort**
- **-R \[bind_address:]RemotePort:LocalSocket**
- **-R RemoteSocket:LocalHost:LocalPort**
- **-R RemoteSocket:LocalSocket**
- **-R \[bind_address:]port**
- 假如现在有 A、B、C 三台主机，A 与 C 不通，A 与 B 通，C 与 B 通；也就是说 B 是中转站(运行 sshd 程序)。如果想要通过 A 访问 C，则需要在 C 上执行 `ssh -R XXX B-IP` 命令
- Local 表示 C 主机，Remote 表示 B 主机
- 从 A 访问 B Port 就是访问 C Port

**端口转发常用附加选项**

- **-f** # 在后台运行 ssh 程序。常与 -N 连用。
  - 如果单独使用，则会报错：`Cannot fork into background without a command to execute.`
- **-g** # 表示 ssh 隧道对应的转发端口将监听在主机的所有 IP 中，不使用 -g 时，转发端口默认只监听在主机的本地回环地址中，-g 表示开启网关模式，远程端口转发中，无法开启网关功能。
- **-N** # 不要执行远程命令。这对于仅让 ssh 用来端口转发时非常有用。

# EXAMPLE

- 常用在脚本中，对远程服务器执行本地脚本，由于在后台执行如法判断脚本退出状态导致 set -e 失效，所以加上了 echo $? > /dev/null 来处理后台脚本的执行状态。
  - ssh -T root@${IP} < ${WorkDir}/install/install-docker.sh; echo $? > /dev/null &
- 连接时，不要严格检查 HostKey。即跳过确实添加 HostKey 的过程，不用输入 yes 或 no。
  - ssh -o StrictHostKeyChecking=no 192.168.0.10
- 通过 IPv6 连接到目标主机。主意添加后面的 `%DEV`，DEV 就是指定的网络设备。
  - ssh -6 fe80::2c75:df14:7422:36a%ens3
- 依次通过 202.43.145.163、172.19.42.247 这两台主机的跳转后，连接到 172.19.42.248
  - ssh -J root@202.43.145.163:42203,172.19.42.247 root@172.19.42.248
- 建立隧道，并禁用远程执行命令，以通过公网远程连接 windows
  - 在内网的 Win10 上执行如下命令，将会在 122.9.154.106 设备上开启 13389 端口监听，其他客户端通过访问 122.9.154.106 的 13389 端口，即可远程登录连接本地 windows
  - ssh -N -R 13389:localhost:3389 root@122.9.154.106 -p 10022
- 建立隧道，并禁用远程执行命令，通过本地的 13389 端口，即可远程连接到 172.19.42.240 的桌面
  - ssh -N -L 13389:172.19.42.240:3389 root@202.43.145.163 -p 42201
- ssh 远程调用函数
  - <https://stackoverflow.com/questions/22107610/shell-script-run-function-from-script-over-ssh> #
- 连接时启用 X11 转发，常用于在 Linux 系统上连接后启动
  - ssh -X -C 192.168.0.1

## 端口转发语法示例

假如现在环境如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzs2hg/1639031352911-319c0d47-4ef2-4aa2-ae0f-b0e3f77615d0.jpeg)

**A 与 C 直接互通**，动态转发

- 在 A 主机上执行命令：
  - **ssh -D localhost:10022 root@B-HOST**
  - 这个命令的意思是：以 root 用户连接到 B 主机，并在本地 28080 端口上启用动态转发功能。相当于在 B 主机上与本地的 10022 端口建立了一个隧道。
- 此时在 A 主机上执行如下命令即可直接远程登录到 C 主机
  - **ssh -o ProxyCommand="nc -x 127.0.0.1:10022 %h %p" root@C-HOST**
  - 也可以使用 socat 工具
    - ssh -o ProxyCommand='socat - socks:127.0.0.1:%h:%p,socksport=10022' root@C-HOST
  - Windows 上可以使用 [Netcat](/docs/4.数据通信/Utility/Netcat.md) 工具
    - ssh -o ProxyCommand="ncat --proxy-type socks5 --proxy 127.0.0.1:10022 %h %p" root@C-HOST

**通过 B 中转，将 C 的端口映射到 A 的端口上**，本地转发

- **通过 A 访问 C 上的 mysql**
  - 在 A 主机上执行
    - **ssh -L 13306:C-IP:3306 root@B-HOST**
    - 这个命令的意思是：A 主机上以 root 用户连接到 B 主机，并在 A 主机本地 13306 端口上启用本地转发，所有到 A 主机的 13306 端口的流量都会通过 B 主机转发到 C 主机的 3306 端口上
  - 此时，即可通过 A 主机的 13306 端口连接 C 主机的 mysql 数据库。
- **通过 A 访问 C 上的 Prometheus**
  - 与通过 A 访问 C 的 mysql 类似
  - 在 A 主机上执行
    - **ssh -L 19090:C-IP:9090 root@B-HOST**
  - 此时通过浏览器访问 http://A-IP:19090 即可打开 C 在 9090 端口上的 Web 应用

**将 C 的端口映射到 B 的端口上**，远程转发

- 在 C 主机上执行
  - **ssh -R 19090:C-IP:9090 root@B-HOST**
- 此时访问 http://B-IP:19090 即可直接打开 C 主机在 9090 端口上的 Web 应用。
- 22 端口同理，只需要将 9090 改为 22，那么在任意一台机器上执行 **ssh root@B-IP -p 19090** 命令就是通过 B 的 19090 端口 ssh 连接到 C 主机。
- 注意：如果远程端口转发时遇到问题，需要在 sshd_config 配置文件中将 `GatewayPorts` 设置为 yes。因为不开启该配置，在 B 上监听的 IP 将会是 127.0.0.1，开启后是 \*

