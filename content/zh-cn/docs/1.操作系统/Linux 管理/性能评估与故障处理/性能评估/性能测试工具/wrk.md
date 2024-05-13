---
title: "wrk"
---

# 概述

> 参考：
>
> - [GitHub 项目,wrk](https://github.com/wg/wrk)

wrk 是一种现代 HTTP 基准测试小型工具，当在单个多核 CPU 上运行时，能够产生大量负载。它结合了多线程设计和可扩展的事件通知系统，例如 epoll 和 kqueue。

可选的 LuaJIT 脚本可以执行 HTTP 请求生成，响应处理和自定义报告。详细信息可在 SCRIPTING 中找到，几个示例位于 scripts /中

## Wrk2

> 参考
>
> - [GitHub 项目,wrk2](https://github.com/giltene/wrk2)
> - <https://www.wangbo.im/posts/usage-of-benchmarking-tool-wrk-and-wrk2/>

# Wrk 的安装

安装 wrk 需要从项目上 clone 项目然后编译获取二进制文件

1. yum groupinstall 'Development Tools' -y
2. yum install openssl-devel git -y
3. git clone <https://github.com/wg/wrk>
4. cd wrk
5. make
6. cp wrk /usr/local/bin

clone 和编译时间较长，这里给一个编译好的文件

# Wrk 命令行工具使用方法

wrk \[OPTIONS] URL
OPTIONS

- **-c,--connections NUM** # 指定总的 http 并发数。默认 10 个并发连接
- **-d,--duration NUM** # 指定压测的持续时间。默认 10s
- **-H,--header STRING** # 使用指定的头信息作为请求 header
- **-t,--threads NUM** # 指定总线程数。默认 2 个线程
- **--latency** # 输出延迟统计情况

EXAMPLE

1. wrk -t 12 -c 400 -d 30s http://127.0.0.1:8080/index.html # 使用 12 个线程运行 30 秒，模拟 400 个并发请求本地 8080 端口的 index.html
2. wrk -t 80 -d 60s -c 16000 -T 3s --latency http://10.10.9.60:80/ #

# WRK 结果解读

下面是一个最基本的测试结果

```bash
Running 1m test @ http://10.10.100.107:30000/
  12 threads and 10000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    线程信息  平均值   标准差    最大值   正负一个标准差所占比例
    Latency     1.28s   509.20ms   2.00s    76.80% # 延时、执行时间
    Req/Sec    70.18     74.26   717.00     88.88% # 每个线程每秒钟执行的连接数
  Latency Distribution #延迟分布。如果使用--latency参数，则会出现该字段信息
     50%    1.40s
     75%    1.60s
     90%    1.79s
     99%    1.99s
  31631 requests in 1.00m, 8.87MB read
  Socket errors: connect 0, read 2745465, write 0, timeout 17902
Requests/sec:    526.31 #平均每秒处理完成请求的个数。每秒请求数(QPS)，等于总请求数/测试总耗时
Transfer/sec:    151.11KB #平均每秒读取数据的值
```

# 并发测试案例

```bash
root@desistdaydream:/usr/local/bin# wrk -t 8 -c 1000 -d 30s --latency -H "Host: desistdaydream.ltd" http://172.19.42.217/robots.txt
Running 30s test @ http://172.19.42.217/robots.txt
  8 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    94.66ms   75.27ms   1.21s    47.54%
    Req/Sec     1.36k   235.57     1.98k    68.50%
  Latency Distribution
     50%  102.70ms
     75%  161.72ms
     90%  189.77ms
     99%  230.97ms
  324797 requests in 30.08s, 97.56MB read
Requests/sec:  10797.85
Transfer/sec:      3.24MB
root@desistdaydream:/usr/local/bin#
root@desistdaydream:/usr/local/bin#
root@desistdaydream:/usr/local/bin#
root@desistdaydream:/usr/local/bin#
root@desistdaydream:/usr/local/bin# wrk -t 8 -c 1000 -d 30s --latency -H "Host: desistdaydream.ltd" http://172.19.42.217
Running 30s test @ http://172.19.42.217
  8 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.00s     0.00us   2.00s   100.00%
    Req/Sec    23.36     20.77   121.00     77.18%
  Latency Distribution
     50%    2.00s
     75%    2.00s
     90%    2.00s
     99%    2.00s
  2629 requests in 30.09s, 118.62MB read
  Socket errors: connect 0, read 0, write 0, timeout 2628
Requests/sec:     87.37
Transfer/sec:      3.94MB

```
