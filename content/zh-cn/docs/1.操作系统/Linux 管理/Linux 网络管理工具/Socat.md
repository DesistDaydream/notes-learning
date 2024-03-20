---
title: Socat
linkTitle: Socat
date: 2024-03-20T08:58
weight: 20
---

# 概述

> 参考：
>
> - 官网：<http://www.dest-unreach.org/socat/>
> - 官方文档：<http://www.dest-unreach.org/socat/doc/socat.html>

Socat 是一个多功能的网络工具，名字来由是” Socket CAT”，可以看作是 netcat 的 N 倍加强版。

Socat 是一个两个独立数据通道之间的双向数据传输的继电器。这些数据通道包含文件、管道、设备（终端或调制解调器等）、socket（Unix，IP4，IP6 - raw，UDP，TCP）、SSL、SOCKS4 客户端或代理 CONNECT。

Socat 支持广播和多播、抽象 Unix sockets、Linux tun/tap、GNU readline 和 PTY。它提供了分叉、记录和进程间通信的不同模式。多个选项可用于调整 socat 和其渠道，Socat 可以作为 TCP 中继（一次性或守护进程），作为一个守护进程基于 socksifier，作为一个 shell Unix 套接字接口，作为 IP6 的继电器，或面向 TCP 的程序重定向到一个串行线。

socat 的主要特点就是在两个数据流之间建立通道；且支持众多协议和链接方式：ip, tcp, udp, ipv6, pipe,exec,system,open,proxy,openssl,socket 等。

工作原理

socat 的运行有 4 个阶段:

- 初始化 解析命令行以及初始化日志系统。
- 打开连接 先打开第一个连接，再打开第二个连接。这个单步执行的。 如果第一个连接失败，则会直接退出。
- 数据转发 谁有数据就转发到另外一个连接上, read/write 互换。
- 关闭 其中一个连接掉开，执行处理另外一个连接。
- 地址类型 参数由 2 部分组成，第一个连接和第二个连接，最简单的用法就是 socat - - 其效果就是输入什么，回显什么其用法主要在于地址如何描述, 下面介绍几个常用的。
- TCP

TCP:<host>:<port> 目标机器 IP 对应端口 portTCP-LISTEN:<port> 本机监听端口。

- UDP

UDP:<host>:<port> 目标机器 host 对应端口 portUDP-LISTEN:<port> 本机监听端口。

- OPENSSL

需要一个证书,否则会失败提示: 2012/04/06 11:29:11 socat\[1614] E SSL_connect(): error:14077410:SSL routines:SSL23_GET_SERVER_HELLO:sslv3 alert handshake failure

OPENSSL:<host>:<port> 目标机器 host 对应端口 portOPENSSL-LISTEN:<port> 本机监听端口。

- TUN

Syntax(语法)

**socat \[OPTIONS] <ADDRESS> <ADDRESS>**

## OPTIONS

ADDRESS

ADDRESS 类似于一个文件描述符，socat 所做的工作就是在 2 个 ADDRESS 指定的描述符间建立一个 pipe 用于发送和接收数据。

格式为：`AddressTYPE[:ARGS][,OPTIONS]`

### AddressTYPE

- STDIN, STDOUT # 表示标准输入输出，可以使用 `-` 符号代替
- /PATH/TO/FILE # 任意路径，也可以使用相对路径，打开一个文件作为数据流
- tcp # 建立一个 TCP 连接作为数据流
- tcp-listen:PORT # 建立 TCP 监听端口
- EXEC # 指定一个程序作为数据流
- unix #
-

### AddressOPTIONS

# EXAMPLE

- 测试 172.19.42.243 的 161/UDP 端口
  - socat - udp:172.19.42.243:161
- 类似于 cat 命令，将 /var/log/messages 中的内容输出到标准输出
  - socat - /var/log/messages
- 监听本机 500 端口 并与 标准输入输出 建立连接。当其他设备 telnet 到本端口时，输入输出都会在两端互相显示。
  - socat tcp-listen:500 -
- 监听本地的 18080 端口，所有到 18080 的数据包，都会转发给 172.19.42.248 的 8080 端口
  - socat TCP-LISTEN:18080,fork,reuseaddr  TCP:172.19.42.248:8080
- 在本地 8080 端口与 docker 的 sock 建立连接。直接访问 <http://172.38.40.250:8080/debug/pprof/> 可以 debug docker(最后的 / 不能少)。
  - socat -d -d tcp-listen:8080,fork,bind:172.38.40.250 UNIX:/var/run/docker.sock

连接目标

`socat - tcp:192.168.1.18:80`

这个命令等同于 nc 192.168.1.18 80。

socat 里面，必须有两个流，所以第一个参数-代表标准的输入输出，第二个流连接到 192.168.1.18 的 80 端口。

`socat -d -d READLINE,history=$HOME/.http_history TCP4:www.qq.com:80`

这个例子支持历史记录查询，类似于 bash 的历史记录。

### 反向连接

再看一个反向 telnet 的例子： on server:

`socat tcp-listen:23 exec:cmd,pty,stderr`

这个命名把 cmd 绑定到端口 23，同时把 cmd 的 Stderr 复位向到 stdout。

on client:

`socat readline tcp:server:23`

连接到服务器的 23 端口，即可获得一个 cmd shell。readline 是 gnu 的命令行编辑器，具有历史功能。

### 向远程端口发数据

`echo “test” | socat – tcp-connect:127.0.0.1:12345`

### 本地开启端口

`socat tcp-l:7777,reuseaddr,fork system:bash`

同 nc -l -p 7777 -e bash。

### 执行 bash 的完美用法

|        |                                           |
| ------ | ----------------------------------------- |
| 服务端 | `socat tcp-l:8888 system:bash,pty,stderr` |
| 本地   | `socat readline tcp:$target:8888`         |

用 readline 替代-，就能支持历史功能了。在这个模式下的客户端有本地一样的效果

### 文件传递

再看文件传递的例子。nc 也经常用来传递文件，但是 nc 有一个缺点，就是不知道文件什么时候传完了，一般要用 Ctrl+c 来终止，或者估计一个时间，用-w 参数来让他自动终止。用 socat 就不用这么麻烦了：

|           |                                                      |
| --------- | ---------------------------------------------------- |
| on host 1 | socat -u open:myfile.exe,binary tcp-listen:999       |
| on host 2 | socat -u tcp:host1:999 open:myfile.exe,create,binary |

这个命令把文件 myfile.exe 用二进制的方式，从 host 1 传到 host 2。-u 表示数据单向流动，从第一个参数到第二个参数，-U 表示从第二个到第一个。文件传完了，自动退出。

### 转发

|                      |                                           |
| -------------------- | ----------------------------------------- |
| 本地端口转向远程主机 | socat TCP4-LISTEN:8888 TCP4:www.qq.com:80 |

如果需要使用并发连接，则加一个 fork,如下:

`socat TCP4-LISTEN:8888,fork TCP4:www.qq.com:80`

本地监听 8888 端口，来自 8888 的连接重定向到目标www.qq.com:80

### 端口映射

再来一个大家喜欢用的例子。在一个 NAT 环境，如何从外部连接到内部的一个端口呢？只要能够在内部运行 socat 就可以了。

|      |                                                 |
| ---- | ----------------------------------------------- |
| 外部 | socat tcp-listen:1234 tcp-listen:3389           |
| 内部 | socat tcp:outerhost:1234 tcp:192.168.12.34:3389 |

这样，你外部机器上的 3389 就映射在内部网 192.168.12.34 的 3389 端口上。

### VPN

| 服务端 | socat -d -d TCP-LISTEN:11443,reuseaddr TUN:192.168.255.1/24,up | | 客户端 | socat TCP:1.2.3.4:11443 TUN:192.168.255.2/24,up |

### 重定向

`socat TCP4-LISTEN:80,reuseaddr,fork TCP4:192.168.123.12:8080`

|             |                                                                        |
| ----------- | ---------------------------------------------------------------------- |
| TCP4-LISTEN | 在本地建立的是一个 TCP ipv4 协议的监听端口                             |
| reuseaddr   | 绑定本地一个端口；                                                     |
| fork        | 设定多链接模式，即当一个链接被建立后，自动复制一个同样的端口再进行监听 |

socat 启动监听模式会在前端占用一个 shell，因此需使其在后台执行。

`socat -d -d tcp4-listen:8900,reuseaddr,fork tcp4:10.5.5.10:3389`

或者

`socat -d -d -lf /var/log/socat.log TCP4-LISTEN:15000,reuseaddr,fork,su=nobody TCP4:static.5iops.com:15000-d -d -lf /var/log/socat.log`是参数，前面两个连续的`-d -d`代表调试信息的输出级别，`-lf`则指定输出信息的保存文件。`TCP4-LISTEN:15000,reuseaddr,fork,su=nobody`是一号地址，代表在 15000 端口上进行 TCP4 协议的监听，复用绑定的 IP，每次有连接到来就 fork 复制一个进程进行处理，同时将执行用户设置为 nobody 用户。`TCP4:static.5iops.com:15000`是二号地址，代表将 socat 监听到的任何请求，转发到`static.5iops.com:15000`上去。

### 读写分流

socat 还具有一个独特的读写分流功能，比如：

    socat open:read.txt!!open:write.txt,create,append tcp-listen:80,reuseaddr,fork

1
Plain Text

这个命令实现一个假的 web server，客户端连过来之后，就把 read.txt 里面的内容发过去，同时把客户的数据保存到 write.txt 里面。”！！”符号用户合并读写流，前面的用于读，后面的用于写。

### 通过 openssl 来加密传输过程

证书生成

    FILENAME=60.*.*.*
    openssl genrsa -out $FILENAME.key 1024
    openssl req -new -key $FILENAME.key -x509 -days 3653 -out $FILENAME.crtcat $FILENAME.key $FILENAME.crt >$FILENAME.pem

1
2
3
Plain Text

在当前目录下生成 `server.pem 、server.crt`

使用

|        |                                                                                        |
| ------ | -------------------------------------------------------------------------------------- |
| 服务端 | socat openssl-listen:4433,reuseaddr,cert=srv.pem,cafile=srv.crt system:bash,pty,stderr |
| 本地   | socat readline openssl:localhost:4433,cert=srv.pem,cafile=srv.crt                      |
