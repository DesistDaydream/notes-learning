---
title: "Httperf"
---

# 概述

Web 压力测试 - Httperf

Httperf 比 ab 更强大，能测试出 web 服务能承载的最大服务量及发现潜在问题；比如：内存使用、稳定性。最大优势：可以指定规律进行压力测试，模拟真实环境。

下载：<http://code.google.com/p/httperf/downloads/list>

1. \[root@localhost ~]# tar zxvf httperf-0.9.0.tar.gz
2. \[root@localhost ~]# cd httperf-0.9.0
3. \[root@localhost httperf-0.9.0]# ./configure
4. \[root@localhost httperf-0.9.0]# make && make install
5. \[root@localhost ~]# httperf --hog --server=192.168.0.202 --uri=/index.html --num-conns=10000 --wsess=10,10,0.1

参数说明：

- \--hog：让 httperf 尽可能多产生连接，httperf 会根据硬件配置，有规律的产生访问连接
- \--num-conns：连接数量，总发起 10000 请求
- \--wsess：用户打开网页时间规律模拟，第一个 10 表示产生 10 个会话连接，第二个 10 表示每个会话连接进行 10 次请求，0.1 表示每个会话连接请求之间的间隔时间 / s
