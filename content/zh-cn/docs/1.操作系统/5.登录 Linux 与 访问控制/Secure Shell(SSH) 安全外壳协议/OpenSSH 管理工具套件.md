---
title: OpenSSH 管理工具套件
---

# 概述

> - [官方文档,Manual(手册)](http://www.openssh.com/manual.html)
> - <https://www.myfreax.com/how-to-setup-ssh-tunneling/>
> - <https://hellolzc.github.io/2020/04/port-forwarding-with-ssh/>
> - <http://tuxgraphics.org/~guido/scripts/ssh-socks5-port-forwarding.html>

- ssh、scp、sftp # 客户端管理工具
- ssh-add、ssh-keysign、ssh-keyscan、ssh-keygen # 密钥管理工具
- sshd、sftp-server、ssh-agent # 服务端管理工具

# ssh # OpenSSH 的 ssh 客户端工具(远程登录程序)

## Syntax(语法)

**ssh \[OPTIONS] \[User@] HostIP \[COMMAND]**

OPTIONS

- **-C** # 请求压缩所有数据(包括 stdin、stdout、stderr、X11 转发的数据、TCP 和 UNIX-domain 连接)
- **-i \</PATH/FILE>** # 使用指定的私钥文件来进行登录认证
- **-J \<DEST\[,DEST2,...]>** # 首先与 DEST 建立 ssh 连接，并通过 DEST 跳转到最终目标主机。如果需要多次跳转。可以指定多个 DEST 并以逗号分割。
  - DEST 格式为：`[USER@]HOST[:PORT]`
- **-o \<Key=Value>** # 以命令行的方式配置本应该在 /etc/ssh/ssh_config 文件中配置的内容。可用于以配置文件中使用的格式提供选项。Key 是 /etc/ssh/ssh_config 配置文件中的关键字。可用的 OPTIONS 详见 /etc/ssh/ssh_config 文件详解。
- **-p \<PORT>** # 指定 HostIP 所在远程服务器监听的端口
- **-T** # 不要分配一个伪终端
- **-W** # 请求将客户端上的标准输入和输出通过安全通道转发到端口上的主机。 表示-N，-T，ExitOnForwardFailure 和 ClearAllForwardings，尽管可以在配置文件中或使用-o 命令行选项覆盖它们。
- **-X** # 启用 X11 转发
- **-Y** # 启用受信任的 X11 转发。 受信任的 X11 转发不受 X11 SECURITY 扩展控件的约束。
- **端口转发选项**
  - **-D \<[Bind_Address]:PORT>** # Dynamic(动态) 转发。启用动态转发的 ssh 程序相当于一个代理服务，通过监听的端口，可以将流量送到指定的目标主机。
  - **-L \<XXX>** # Local(本地) 转发。发往本地的 TCP 端口 或 Unix Socket 上的流量转发到远端 TCP 端口 或 Unix Socket 上。
    - XXX 有多种语法格式：
    - **-L \[Bind_Address:]LocalPort:RemoteHost:RemoteHostPort**
    - **-L \[Bind_Address:]LocalPort:RemoteSocket**
    - **-L LocalSocket:RemoteHost:RemotePort **
    - **-L LocalSocket:RemoteSocket**
    - 假如现在有 A、B、C 三台主机，A 与 C 不通，A 与 B 通，C 与 B 通；也就是说 B 是中转站(运行 sshd 程序)。如果想要通过 A 访问 C，则需要在 A 上执行 `ssh -L XXX B-IP` 命令
    - Local 表示 A 主机，Remote 表示 C 主机
    - 访问 A Port 就是访问 C Port
  - **-R \<XXX>** # Remote(远程) 转发。发往指定 `RemotePort` 或 `RemoteSocket` 上的流量转发到 `本地端口` 或 `Unix Socket` 上。远程转发其实更像将 Local 服务通过 ssh 程序以类似 nat 的方式暴露到 Remote 上。
    - XXX 有多种语法格式：
    - **-R \[bind_address:]RemotePort:LocalHost:LocalPort**
    - **-R \[bind_address:]RemotePort:LocalSocket**
    - **-R RemoteSocket:LocalHost:LocalPort**
    - **-R RemoteSocket:LocalSocket**
    - **-R \[bind_address:]port**
    - 假如现在有 A、B、C 三台主机，A 与 C 不通，A 与 B 通，C 与 B 通；也就是说 B 是中转站(运行 sshd 程序)。如果想要通过 A 访问 C，则需要在 C 上执行 `ssh -R XXX B-IP` 命令
    - Local 表示 C 主机，Remote 表示 B 主机
    - 从 A 访问 B Port 就是访问 C Port
- **端口转发常用附加选项**
  - **-f** # 在后台运行 ssh 程序。常与 -N 连用。
    - 如果单独使用，则会报错：`Cannot fork into background without a command to execute.`
  - **-g** # 表示 ssh 隧道对应的转发端口将监听在主机的所有 IP 中，不使用 -g 时，转发端口默认只监听在主机的本地回环地址中，-g 表示开启网关模式，远程端口转发中，无法开启网关功能。
  - **-N** # 不要执行远程命令。这对于仅让 ssh 用来端口转发时非常有用。

## 端口转发语法示例

假如现在环境如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/mzs2hg/1639031352911-319c0d47-4ef2-4aa2-ae0f-b0e3f77615d0.jpeg)

**A 与 C 直接互通**

- 在 A 主机上执行命令：
  - **ssh -D localhost:10022 root@B-HOST**
  - 这个命令的意思是：以 root 用户连接到 B 主机，并在本地 28080 端口上启用动态转发功能。相当于在 B 主机上与本地的 10022 端口建立了一个隧道。
- 此时在 A 主机上执行如下命令即可直接远程登录到 C 主机
  - **ssh -o ProxyCommand="nc -x 127.0.0.1:10022 %h %p" root@C-HOST**
  - 也可以使用 socat 工具
    - ssh -o ProxyCommand='socat - socks:127.0.0.1:%h:%p,socksport=10022' root@C-HOST
  - Windows 上可以使用 [ncat ](https://nmap.org/download.html)工具
    - ssh -o ProxyCommand="ncat --proxy-type socks5 --proxy 127.0.0.1:10022 %h %p" root@C-HOST

**通过 B 中转，将 C 的端口映射到 A 的端口上**

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

**将 C 的端口映射到 B 的端口上**

- 在 C 主机上执行
  - **ssh -R 19090:C-IP:9090 root@B-HOST**
- 此时访问 http://B-IP:19090 即可直接打开 C 主机在 9090 端口上的 Web 应用。
- 22 端口同理，只需要将 9090 改为 22，那么在任意一台机器上执行 **ssh root@B-IP -p 19090** 命令就是通过 B 的 19090 端口 ssh 连接到 C 主机。
- 注意：如果远程端口转发时遇到问题，需要在 sshd_config 配置文件中将 `GatewayPorts` 设置为 yes。因为不开启该配置，在 B 上监听的 IP 将会是 127.0.0.1，开启后是 \*

## EXAMPLE

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

# scp # 基于 ssh 协议的文件传输工具

## Syntax(语法)

**scp [OPTIONS] SourceFILE DestinationFILE**

Note：远程 FILE 的格式为：USER@IP:/PATH/FILE)

OPTIONS：

- -p 
- **-r** # 以递归方式复制，用于复制整个目录

## EXAMPLE

把本地 nginx 文件推上去复制到以 root 用户登录的 10.10.10.10 这台机器的/opt/soft/scptest 目录下

  - scp /opt/soft/nginx-0.5.38.tar.gz root@10.10.10.10:/opt/soft/scptest

把以 root 用户登录的 10.10.10.10 机器中的 nginx 文件拉下来复制到本地/opt/soft 目录下

  - scp root@10.10.10.10:/opt/soft/nginx-0.5.38.tar.gz /opt/soft/

基于密钥的认证,当对方主机 ssh 登录的用户的家目录存在公钥，并且公钥设置密码为空，那么以后 ssh 协议登录传输都可以直接登录而不用密码

# ssh-keygen # 在客户端生成密钥对

**ssh-keygen -t rsa \[-P ''] \[-f ~/.ssh/id_rsa]**

- EXAMPLE：
  - ssh-keygen -t rsa -P '' -f ~/.ssh/id_rsa

# ssh-copy-id # 把生成的公钥传输至远程服务器对应用户的家目录

**ssh-copy-id \[-i \[Identity_File]] \[User@]HostIP**
Identity_File(身份文件) # 一般为 /root/.ssh/id_rsa.pub
EXAMPLE：

- 将公钥拷贝到服务端
  - ssh-copy-id -i /root/.ssh/id_rsa.pub root@192.168.0.10
- 若没有 ssh-copy-id 命令，则可以这么这么弄
  - cat ~/.ssh/id_rsa.pub | ssh root@192.168.0.10 'cat >> .ssh/authorized_keys'
