---
title: iperf 网络性能测量和调整工具
---

# 概述

> 参考：
> - [官网](https://iperf.fr/)

在 server 端监听某个端口，然后 client 用同样的 iperf 访问服务端，来进行性能测试

所以该工具需要在两台设备之间一同使用，其中一台作为服务端，另外一台作为客户端，两端互相通信，才可测试网络性能。而命令行的 OPTIONS 也是分为全局、服务端特定、客户端特定 这三类

# Syntax(语法)

**iperf3 <-s | -c HOST> \[OPTIONS]**

**通用选项**

- **-p, --port** # 指定服务端监听的端口或者客户端要连接的端口
- **-f, --format \[kmgKMG]** # 指定输出格式。可以是：Kbits, Mbits, KBytes, MBytes
- **-i, --interval** # 指定每次带宽报告间隔的秒数。i.e.每隔几秒输出一次数据。默认每 1 秒报告一次
- -F, --file name # xmit/recv the specified file
- -A, --affinity n/n,m # set CPU affinity
- -B, --bind <host> # bind to a specific interface
- -V, --verbose # more detailed output
- -J, --json # output in JSON format
- \--logfile f # send output to a log file
- \--forceflush # force flushing output at every interval
- -d, --debug # emit debugging output

**服务端的特定选项**

- **-s, --server** # 在服务器模式下运行 iperf3，默认程序监听在 5201 端口上
- **-D, --daemon** # 以守护进程的形式运行服务端
- -I, --pidfile file # write PID file
- -1, --one-off # handle one client connection then exit

**客户端的特定选项**

- **-c, --client <HOST>** # 在客户端模式下运行 iperf3，并连接到指定的服务端主机 HOST
- **-u, --udp** # 使用 UPD 模式进行测试。默认为 TCP
- **-b, --bandwidth <NUM>** # 指定目标带宽上限，单位是 bits/s（0 表示无限制）（UDP 默认为 1 Mbit / sec，TCP 无限制）
  - 该选项为每个线程的带宽上限，比如我如果 -P 选项为 2，-b 为 100M ，那么当前测试每个线程的带宽上限为 100M，总上限 200M
- \--fq-rate #\[KMG]enable fair-queuing based socket pacing in bits/sec (Linux only)
- **-t, --time** # 指定传输数据的总时间。(默认为 10 秒)
- **-n, --bytes \[KMG]** # 要传输的字节数 (不可与 -t 选项同用)
- -k, --blockcount \[KMG] # number of blocks (packets) to transmit (instead of -t or -n)
- -l, --len \[KMG] # length of buffer to read or write (default 128 KB for TCP, dynamic or 1 for UDP)
- \--cport <port> # bind to a specific client port (TCP and UDP, default: ephemeral port)
- **-P, --parallel <NUM>** # 并发数
- -R, --reverse # run in reverse mode (server sends, client receives)
- -w, --window \[KMG] # set window size / socket buffer size
- -C, --congestion <algo> #set TCP congestion control algorithm (Linux and FreeBSD only)
- -M, --set-mss # set TCP/SCTP maximum segment size (MTU - 40 bytes)
- -N, --no-delay # set TCP/SCTP no delay, disabling Nagle's Algorithm
- -4, --version4 only use IPv4
- -6, --version6 only use IPv6
- -S, --tos N set the IP 'type of service'
- -L, --flowlabel N set the IPv6 flow label (only supported on Linux)
- -Z, --zerocopy use a 'zero copy' method of sending data
- -O, --omit N omit the first n seconds
- -T, --title str prefix every output line with this string
- \--get-server-output get results from server
- \--udp-counters-64bit use 64-bit counters in UDP test packets

# 应用实例

## 基础用法

服务端命令：iperf3 -s
客户端命令：iperf3 -c 10.10.100.250

这时候服务端的 iperf3 程序会监听在 5201 端口上，客户端会访问服务端(这里 ip 是 10.10.100.250)的 2501 端口进行网络测试，测试效果如图
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fdemaq/1616164267994-3d8e4e2d-0c26-4b52-8054-12aeac917398.png)
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fdemaq/1616164268013-f44eee21-25f2-48a5-acad-07b32cf7af7a.png)
客户端每秒会往服务端发送一次数据，Interval 表示时间间隔，Transfer 表示传输的数据量，Bandwidth 表示带宽的大小，Retr 表示重传次数

在客户端的最后两行表示 10 秒钟的传送的总数据量，以及平均带宽，第一行是发送的，第二行是接收的。这次测试结果就是两台服务器之间最大带宽是 20G

## 查看网络丢包率和延迟

服务端命令：iperf3 -s
客户端命令：iperf3 -c 10.10.100.250 -u -b 10M -t 10 -i 1 -P 100
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fdemaq/1616164268023-89190ba4-c7a7-4587-b4f7-389086f5d465.png)
Jitter 表示抖动 i.e.数据包延迟时间。Lost/Total Datagrams 表示“丢失/总数据包”的数量，扩内的百分比为丢包率。

## 吞吐测试

服务端命令：iperf3 -s -i 1 -p 10000
客户端命令：iperf3 -c 172.19.42.221 -p 10000 -b 1G -t 15 -P 2

    [ ID] Interval           Transfer     Bandwidth       Retr
    [  4]   0.00-15.00  sec   781 MBytes   437 Mbits/sec  4856             sender
    [  4]   0.00-15.00  sec   779 MBytes   436 Mbits/sec                  receiver
    [  6]   0.00-15.00  sec   876 MBytes   490 Mbits/sec  7074             sender
    [  6]   0.00-15.00  sec   874 MBytes   489 Mbits/sec                  receiver
    [SUM]   0.00-15.00  sec  1.62 GBytes   927 Mbits/sec  11930             sender
    [SUM]   0.00-15.00  sec  1.61 GBytes   925 Mbits/sec                  receiver

TCP 吞吐(带宽)大概为 900+M，也就是千兆基本 能跑慢
