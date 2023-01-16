---
title: Linux 代理配置
---

#

# Linux 代理服务相关变量：

- http_proxy |https_proxy | ftp_proxy |all_proxy # 此变量值用于所有 http、https、ftp 或者所有流量
- socks_proxy # 在大多数情况下，它用于 TCP 和 UDP 流量。其值通常采用 socks：// address：port 格式。
- rsync_proxy # 这用于 rsync 流量，尤其是在 Gentoo 和 Arch 等发行版中。
- no_proxy # 以逗号分隔的域名或 IP 列表，应绕过代理。该本地主机就是一个很好的例子。一个例子是 localhost，127.0.0.1。

语法格式

XXXX_proxy='http://\[USER:PASSWORD@]ServerIP:PORT/' #需要设置用户名，密码，代理服务器的 IP 和端口，用户名和密码可省

EXAMPLE

- http_proxy="http://tom:secret@10.23.42.11:8080/" #设置本机的 http 代理服务器为 10.23.42.11:8080，用户名是 tom，密码是 secret

- 同时设置 3 种类型代理，没有用户名和密码，代理服务器是 192.168.19.79:1080
  - export {https,ftp,http}\_proxy="127.0.0.1:8889"
- all_proxy="socks5://localhost:10808" #使用本地 10808 端口的 socks 协议代理所有流量(e.g.安装完 v2ray 客户端并配置好启动后，即可使用该变量来让设备使用代理进行翻墙)

- no_proxy="10._._._,192.168._._,_.local,localhost,127.0.0.1" #忽略指定 ip 的代理

**注意：通过 Systemd 启动的进程，无法识别这些环境变量，只能通过 Unit File 中的 \[Service] 部分的 Environment 指令指定代理信息。**

# Linux bash 终端设置代理（proxy）访问

Linux 是开源操作系统，有很多软件包需要从国外网站上下载，而其中大部分国外网站是被墙的，这时我们需要通过代理来访问这些网站。下面我们介绍 Linux bash shell 终端代理设置方法，包括 socks 代理，http 代理。

一、linux shell 终端代理设置方法：

linux 要在 shell 终端为 http、https、ftp 协议设置代理，值需要设置对应的环境变量即可。下面是一些关于代理的环境变量：

|            |                        |                 |
| ---------- | ---------------------- | --------------- |
| 环境变量   | 描述                   | 值示例          |
| http_proxy | 为 http 网站设置代理； | 10.0.0.51:8080; |

user:pass@10.0.0.10:8080
socks4://10.0.0.51:1080
socks5://192.168.1.1:1080 |
| https_proxy | 为 https 网站设置代理； | 同上 |
| ftp_proxy | 为 ftp 协议设置代理； | socks5://192.168.1.1:1080 |
| no_proxy | 无需代理的主机或域名；
可以使用通配符；
多个时使用“,”号分隔； | _.aiezu.com,10._._._,192.168._._,
\*.local,localhost,127.0.0.1 |

可以将上面 4 个环境变量设置项放于~/.bashrc 文件尾部，这样用户打开 bash shell 终端时会自动调用此脚本，读入它们。

二、linux bash 为 http 站点设置代理：

根据代理类型，将下面对应的设置项添加到~/.bashrc 文件末尾，然后运行". ~/.bashrc"（前面是一个“.”号）命令使用之在当前环境生效。

1、为 http 站点设置 http 代理（默认）：

|     |                                  |
| --- | -------------------------------- |
| 1   | export http_proxy=10.0.0.52:8080 |

2、为 http 站点设置 sock4、sock5 代理：

|     |     |
| --- | --- |

| 1
2
3
4
5
6 | # 设置 socks 代理，自动识别 socks 版本
export http_proxy=socks://10.0.0.52:1080
\# 设置 socks4 代理
export http_proxy=socks4://10.0.0.52:1080
\# 设置 socks5 代理
export http_proxy=socks5://10.0.0.52:1080 |

3、代理使用用户名密码认证：

|     |                                              |
| --- | -------------------------------------------- |
| 1   | export http_proxy=user:pass@192.158.8.8:8080 |

三、linux bash 为 https 站点设置代理：

如果需要为 https 网站设置代理，设置 https_proxy 环境变量即可；设置方法完全与 http_proxy 环境变量相同：

|     |     |
| --- | --- |

| 1
2
3
4
5
6 | # 任意使用一项
export https_proxy=10.0.0.52:8080
export https_proxy=user:pass@192.158.8.8:8080
export https_proxy=socks://10.0.0.52:1080
export https_proxy=socks4://10.0.0.52:1080
export https_proxy=socks5://10.0.0.52:1080 |

四、举例：

现在我们要设置 http、https 网站都使用 socks5 代理 10.0.0.52:1080，下面为完整设置方法：

1、vim ~/.bashrc，在文件尾部添加下面内容：

|     |     |
| --- | --- |

| 1
2
3 | export http_proxy=socks5://10.0.0.52:1080
export https_proxy=socks5://10.0.0.52:1080
export no_proxy="_.aiezu.com,10._._._,192.168._._,\*.local,localhost,127.0.0.1" |

2、加载设置:

|     |     |
| --- | --- |

| 1
2
3
4
5 | \[root@aiezu.com ~]# . ~/.bashrc
\[root@aiezu.com ~]# echo $http_proxy
socks5://10.0.0.52:1080
\[root@aiezu.com ~]# echo $https_proxy
socks5://10.0.0.52:1080 |

3、测试代理：

|     |     |
| --- | --- |

| 1
2
3
4
5
6
7
8
9
10 | \[root@aiezu.com ~]# curl -I <http://www.fackbook.com>
HTTP/1.1 200 OK
Content-Length: 2423
Content-Type: text/html
Last-Modified: Mon, 14 Nov 2016 22:03:32 GMT
Accept-Ranges: bytes
ETag: "0521af0c23ed21:0"
Server: Microsoft-IIS/7.5
X-Powered-By: ASP.NET
Date: Sun, 11 Dec 2016 13:21:33 GMT |

# proxychains

项目地址：<https://github.com/haad/proxychains>

凡是通过 proxychains 程序运行的程序都会通过 proxychains 配置文件中设置的代理配置来发送数据包。

apt install proxychains 即可

修改配置文件

sock5 127.0.0.1 10808 #指定本地代理服务所监听的地址

proxychains /opt/google/chrome/chrome #即可通过代理打开 chrome 浏览器

proxychains curl -I <https://www.google.com> 会成功
