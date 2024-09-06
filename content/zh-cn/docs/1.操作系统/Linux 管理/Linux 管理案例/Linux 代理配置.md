---
title: Linux 代理配置
---


# 概述

在 [Unix-like OS](/docs/1.操作系统/Operating%20system/Unix-like%20OS/Unix-like%20OS.md) 中，很多程序都会读取 [Terminal 与 Shell](/docs/1.操作系统/Terminal%20与%20Shell/Terminal%20与%20Shell.md) 中的某些变量来读取代理信息

> TODO: 这些变量到底应该大写还是小写？wget 命令无法识别到大写的变量。

- http_proxy | https_proxy | ftp_proxy | all_proxy # 此变量值用于所有 http、https、ftp 或者所有流量
- socks_proxy # 在大多数情况下，它用于 TCP 和 UDP 流量。其值通常采用 socks：// address：port 格式。
- rsync_proxy # 这用于 rsync 流量，尤其是在 Gentoo 和 Arch 等发行版中。
- no_proxy # 以逗号分隔的域名或 IP 列表，应绕过代理。该本地主机就是一个很好的例子。一个例子是 localhost，127.0.0.1。

TODO: 但是这些变量却不是 Bash 的自带变量，但是这些程序却无一例外得统一使用这些变量，具体为什么暂时不知道

语法格式

XXXX_proxy="http://[USER:PASSWORD@]ServerIP:PORT/" # 需要设置用户名，密码，代理服务器的 IP 和端口，用户名和密码可省

EXAMPLE

- http_proxy="http://tom:secret@10.23.42.11:8080/" # 设置本机的 http 代理服务器为 10.23.42.11:8080，用户名是 tom，密码是 secret

- 同时设置 3 种类型代理，没有用户名和密码，代理服务器是 192.168.19.79:1080
  - `export {https,ftp,http}_proxy="127.0.0.1:8889"`
- `all_proxy="socks5://localhost:10808"` # 使用本地 10808 端口的 socks 协议代理所有流量(e.g.安装完 v2ray 客户端并配置好启动后，即可使用该变量来让设备使用代理进行翻墙)
- `no_proxy="10._._._,192.168._._,_.local,localhost,127.0.0.1`" # 忽略指定 ip 的代理

**注意：通过 Systemd 启动的进程，无法识别这些环境变量，只能通过 Unit File 中的 \[Service] 部分的 Environment 指令指定代理信息。**

# 每种变量详解

## http_proxy

为 http 网站设置代理

示例值：

- 10.0.0.51:8080
- user:pass@10.0.0.10:8080
- socks4://10.0.0.51:1080
- socks5://192.168.1.1:1080

## https_proxy

为 https 网站设置代理

示例值：

- 与 http_proxy 相同

## ftp_proxy

为 ftp 协议设置代理

示例值：

- socks5://192.168.1.1:1080

## no_proxy

无需代理的主机或域名，可以使用通配符；多个时使用 `,` 号分隔

示例值：

- `*.aiezu.com,10.*.*.*,192.168.*.*`
- `*.local,localhost,127.0.0.1`

# 最佳实践

## 为 WSL2 设置代理

设置为本地计算机的 Clash

```bash
#!/bin/bash
#
export hostip=$(cat /etc/resolv.conf |grep -oP '(?<=nameserver\ ).*')
export http_proxy="http://${hostip}:7890"
export https_proxy="http://${hostip}:7890"
export all_proxy="sock5://${hostip}:7890"
export ALL_PROXY="sock5://${hostip}:7890"
```

# 可用的代理程序

参考 [代理](/docs/Web/代理.md)

