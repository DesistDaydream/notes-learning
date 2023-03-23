---
title: TCP Wrappers
---

# 概述

> 参考：

注意：CentOS8 及 RHEL8 不再支持 TCP_Wrappers！！！！使用 firewalld 代替！！

**Transmission Control ProtocolWrappers(简称 TCP_Wrappers)** 是一个基于主机的网络访问控制表系统，用于过滤对类 Unix 系统（如 Linux 或 BSD）的网络访问。

其能将主机或子网 IP 地址、名称及 ident 查询回复作为筛选标记，实现访问控制。

## **Tcp_Wrappers 特点**

- 工作在第四层（传输层）的 TCP 协议
- 对有状态连接的特定服务进行安全检测并实现访问控制
- 以库文件形式实现
- 某进程是否接受 libwrap 的控制取决于发起此进程的程序在编译时是否针对 libwrap 进行编译的

## **判断程序是否支持 Tcp_Wrappers**

程序如果调用了 libwrap.so 库，表示支持。

    ldd 程序路径|grep libwrap.so
    strings 程序路径|grep libwrap.so
    #ldd /usr/sbin/sshd|grep libwrap.so
        libwrap.so.0 => /lib64/libwrap.so.0 (0x00007f9851678000)
    #ldd /usr/sbin/vsftpd |grep libwrap.so
        libwrap.so.0 => /lib64/libwrap.so.0 (0x00007f802ef50000)
    #strings `which sshd`|grep libwrap.so
    libwrap.so.0

## TCP_Wrappers 的执行处理机制了

TCP_Wrappers 只需要通过两个文件来处理，/etc/hosts.allow 和/etc/hosts.deny。匹配流程如下图
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cmiwuq/1624581688343-e7b5ceb8-1d56-4b6b-9872-f71462556d23.png)
Note：所以，如果想要仅允许一个 ip 可以通过，那么需要在 deny 中拒绝所有，否则只在 allow 添加该 ip，那么其他 ip 在 allow 没匹配到后，会去 deny 查找，如果还是无法匹配，则直接就通过了。

# TCP_Wrappers 使用

TCP Wrappers 是通过 /etc/hosts.allow 和 /etc/hosts.deny 这两个配置文件来实现一个类似防火墙的机制。

# TCP_Wrappers 关联文件

帮助参考：man 5 hosts_access，man 5 hosts_options
**/etc/hosts.allow** # 允许访问规则
**/etc/hosts.deny** # 拒绝访问规则

注意：说明文档中表示此文件也可以实现拒绝的规则，本着见名知义和管理清晰化的指引，这种写法不是讨论的重点。

**配置文件的写法: Daemon_List : Client_List \[:Shell_Command]**

**Daemon_List**# 单个应用程序的二进制文件名，而不是服务名，如果有多个，用逗号或空格分隔。如 sshd,vsftpd 或 sshd vsftpd

- 可以绑定服务地址,如，sshd@192.168.7.202:ALL(比如一台设备有俩网卡俩地址时)
- 支持通配符
- 内置关键字：
  - ALL 所有进程

**Client_List** # 客户端列表

- 基于单个 IP 地址：192.168.10.1
- 基于网段 IP 地址：192.168.1. 注意，192.168.1.0 这个写法是错误的。
- 基于主机名：www.hunk.tech .hunk.tech 较少用
- 基于网络/掩码：192.168.0.0/255.255.255.0
- 基于 net/prefixlen: 192.168.1.0/24（仅 CentOS7）
- 基于网络组（NIS 域）：@mynetwork
- 内置关键字：
  - ALL 所有主机
  - LOCAL 名称中不带点的主机
  - KNOWN 可以解析的主机名
  - UNKNOWN 无法解析的主机名
  - PARANOID 正、反向查询不匹配或无法解析
- 支持通配符

**Shell_Command # 执行指令**
如：sshd:all:spawn echo "`date +%%F-%%T` from %a pid=%p to %s" >> /app/sshd.log

- EXCEPT 是排除的意思，一行规则里可以有多个，后面的是对前面的结果集进行排除。
  - vsftpd:172.16. EXCEPT 172.16.100.0/24 EXCEPT 172.16.100.1
  - 匹配整个 172.16 网段，但是把 172.16.100 的网段排除，在排除 172.16.100 网段中又把 172.16.100.1 的 IP 给排除。
- spawn 启动一个外部程序完成执行的操作，可以支持内置变量。内置变量请 man ,找%的选项
  - %a (%A) 客户端 IP
  - %c 客户端信息，可以是 IP 或主机名(如果能解析)
  - %p 服务器进程信息 (PID)
  - %s 连接的服务端的信息
  - %% 当规则中包含%时，使用双%转义
- twist 特殊扩展
  - 以指定的命令执行，执行后立即结束该连接。需在 spawn 之后使用。

# 应用示例

对于 sshd 进程，仅允许 10.10.100.250 访问

    [root@lichenhao ~]# cat /etc/hosts.allow
    sshd:10.10.100.250
    [root@lichenhao ~]# cat /etc/hosts.deny
    sshd:ALL

或者

    [root@lichenhao ~]# cat /etc/hosts.deny
    sshd:ALL EXCEPT 10.10.100.250

使用的 2 台测试主机网络 IP 配置如下：

|      |                   |               |               |
| ---- | ----------------- | ------------- | ------------- |
| 简称 | 主机              | IP 1          | IP 2          |
| 6A   | 6-web-1.hunk.tech | 192.168.7.201 | 192.168.5.102 |
| 7B   | 7-web-2.hunk.tech | 192.168.7.202 | 192.168.5.103 |

默认 hosts.allow 和 hosts.deny 配置文件为空，表示全部允许。

拒绝某个 IP 访问：

    7B:
    vim /etc/hosts.deny
    sshd:192.168.7.201
    6A:
    #ssh 192.168.7.202
    ssh_exchange_identification: Connection closed by remote host
    配置规则保存后，立即生效
    7B：日志会明确记录
    #tail -n1 /var/log/secure
    Feb  8 11:18:29 7-web-2 sshd[1811]: refused connect from 192.168.7.201 (192.168.7.201)

把每个 ssh 登录日志记录到文件

    #vim /etc/hosts.allow
    sshd:all:spawn echo "`date +%%F-%%T` from %a pid=%p to %s" >> /app/sshd.log
    #cat /app/sshd.log
    2018-02-08-15:59:53 from 192.168.7.202 pid=2565 to sshd@192.168.7.202
    #ps aux |grep 2565
    root       2565  0.0  2.3 145696  5328 ?        Ss   15:59   0:00 sshd: root@pts/2

## 自动锁定登录失败的 IP

编写脚本 /root/bin/checkip.sh，每 5 分钟检查一次，如果发现通过 ssh 登录失败次数超过 10 次，自动将此远程 IP 放入 Tcp_Wrapper 的黑名单中予以禁止防问

```bash
#!/bin/bash
#定义 休眠时间
sleeptime=300
#定义 通过ssh登录失败次数
num=10
#定义 黑名单文件
file=/etc/hosts.deny
#无限循环
while true;do
  #将失败登录的记录逐行读入变量
  lastb | grep ssh|awk -F "[ ]+" '{print $3}'|uniq -c | while read conn ip;do
    #判断失败次数
    if  [ "$conn" -ge "$num" ];then
      #判断记录的IP是否存在
      egrep -q ^sshd.*$ip $file
      #如果不存在记录，将追加记录至指定黑名单文件
      [ $? -ne 0 ] &&  echo "sshd:$ip" >> $file
    fi
  done
  sleep $sleeptime
done
```

使用 `watch -n1 cat /etc/hosts.deny` 来观察动态文件
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/cmiwuq/1624581861186-486e5699-5dea-4062-83d7-f7bcb0d874e7.gif)
总结：TCP_Wrappers

适用于需求简单的应用场景，并且受到监控软件的是否支持 libwrap.so 库局限。
