---
title: 通过pprof监控docker
---

debug 模式启动 docker

    $ /usr/bin/docker daemon -D -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock

# 通过 socat 端口转发

    $ socat -d -d TCP-LISTEN:8080,fork,bind=192.168.1.137 UNIX:/var/run/docker.sock

## 测试

    [root@reg pprof]#  curl -s http://10.39.0.102:8080/debug/vars | jq .
    {
    "cmdline": [
    "/usr/bin/dockerd",
    "-D"
    ],
    "memstats": {
    "Alloc": 13847856,
    "TotalAlloc": 71577968,
    "Sys": 27052344,
    "Lookups": 7829,
    "Mallocs": 891300,
    "Frees": 772846,
    "HeapAlloc": 13847856,
    "HeapSys": 18743296,
    "HeapIdle": 1941504,
    "HeapInuse": 16801792,
    "HeapReleased": 1810432,
    "HeapObjects": 118454,
    "StackInuse": 1179648,
    "StackSys": 1179648,
    "MSpanInuse": 225280,
    "MSpanSys": 262144,
    "MCacheInuse": 4800,
    "MCacheSys": 16384,
    "BuckHashSys": 1460436,
    "GCSys": 1374208,
    "OtherSys": 4016228,
    "NextGC": 25872553,
    "LastGC": 1512984476111075800,
    "PauseTotalNs": 29246607,
    "PauseNs": [
    317474,
    1159328,
    271770,
      ...

## 获取命令行

    $ curl -s http://10.39.0.102:8080/debug/pprof/cmdline
    /usr/bin/dockerd-D

## 通过客户端获取

    $ go tool pprof http://10.39.0.102:8080/debug/pprof/profile
    Fetching profile from http://10.39.0.102:8080/debug/pprof/profile
    Please wait... (30s)
    Saved profile in /root/pprof/pprof.dockerd.10.39.0.102:8080.samples.cpu.001.pb.gz
    Entering interactive mode

## 生成文件转成 pdf

    $  go tool pprof --pdf pprof.dockerd.10.39.0.102\:8080.samples.cpu.001.pb >call.pdf

## get symbol

    $ curl -s http://10.39.0.102:8080/debug/pprof/symbol
      num_symbols: 1
    1
    2

## 如果你感兴趣，其它的信息都可以获取到

    $ culr -s http://10.39.0.102:8080/debug/pprof/block
    $ curl -s http://10.39.0.102:8080/debug/pprof/heap
    $ curl -s http://10.39.0.102:8080/debug/pprof/goroutine
    $ curl -s http://10.39.0.102:8080/debug/pprof/threadcreate
