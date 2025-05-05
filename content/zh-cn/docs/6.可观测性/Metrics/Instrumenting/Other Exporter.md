---
title: Other Exporter
---

# 概述

官方推荐的一些第三方 exporter: https://prometheus.io/docs/instrumenting/exporters/

process-exporter: 采集进程指标

- https://github.com/ncabatoff/process-exporter

# K8S 集群中常用的 exporter

> - [Vermouth 博客，高可用 Prometheus 问题集锦](http://www.xuyasong.com/?p=1921)


可以在[这里](https://prometheus.io/docs/instrumenting/exporters/)看到官方、非官方的 exporter。如果还是没满足你的需求，你还可以自己编写 exporter，简单方便、自由开放，这是优点。

但是过于开放就会带来选型、试错成本。之前只需要在 zabbix agent 里面几行配置就能完成的事，现在你会需要很多 exporter 搭配才能完成。还要对所有 exporter 维护、监控。尤其是升级 exporter 版本时，很痛苦。非官方 exporter 还会有不少 bug。这是使用上的不足，当然也是 Prometheus 的设计原则。

K8S 生态的组件都会提供 / metric 接口以提供自监控，这里列下我们正在使用的：

- [cadvisor](http://www.xuyasong.com/?p=1483): 集成在 Kubelet 中。
- kubelet: 10255 为非认证端口，10250 为认证端口。
- apiserver: 6443 端口，关心请求数、延迟等。
- scheduler: 10251 端口。
- controller-manager: 10252 端口。
- etcd: 如 etcd 写入读取延迟、存储容量等。
- docker: 需要开启 experimental 实验特性，配置 metrics-addr，如容器创建耗时等指标。
- kube-proxy: 默认 127 暴露，10249 端口。外部采集时可以修改为 0.0.0.0 监听，会暴露：写入 iptables 规则的耗时等指标。
- [kube-state-metrics](http://www.xuyasong.com/?p=1525): K8S 官方项目，采集 pod、deployment 等资源的元信息。
- [node-exporter](http://www.xuyasong.com/?p=1539): Prometheus 官方项目，采集机器指标如 CPU、内存、磁盘。
- blackbox_exporter: Prometheus 官方项目，网络探测，dns、ping、http 监控
- process-exporter: 采集进程指标
- nvidia exporter: 我们有 gpu 任务，需要 gpu 数据监控
- node-problem-detector: 即 npd，准确的说不是 exporter，但也会监测机器状态，上报节点异常打 taint
- 应用层 exporter: mysql、nginx、mq 等，看业务需求。

还有各种场景下的[自定义 exporter](http://www.xuyasong.com/?p=1942)，如日志提取后面会再做介绍。

## Kubernetes Event Exporter

> 参考：
> 
> - [GitHub 项目，opsgenie/kubernetes-event-exporter](https://github.com/opsgenie/kubernetes-event-exporter)

将 Kubernetes 中 Event 资源的内容导出到多个目的地

# Process Exporter

> 参考：
>
> - [GitHub 项目，ncabatoff/process-exporter](https://github.com/ncabatoff/process-exporter)
> - https://mp.weixin.qq.com/s/sbnTByKJYFKQrvnU_iPnZA

process_names 下的数组定义进程组名称及该进程组的匹配条件，一共 3 个匹配方式

- **comm** # 与 /proc/${pid}/stat 中第二个字段进行匹配
- **exe** # 
- **cmdline** # 与进程的所有参数进行匹配

```yaml
process_names:
  # comm is the second field of /proc/<pid>/stat minus parens.
  # It is the base executable name, truncated at 15 chars.
  # It cannot be modified by the program, unlike exe.
  - comm:
    - bash

  # exe is argv[0]. If no slashes, only basename of argv[0] need match.
  # If exe contains slashes, argv[0] must match exactly.
  - exe:
    - postgres
    - /usr/local/bin/prometheus

  # cmdline is a list of regexps applied to argv.
  # Each must match, and any captures are added to the .Matches map.
  - name: "{{.ExeFull}}:{{.Matches.Cfgfile}}"
    exe:
    - /usr/local/bin/process-exporter
    cmdline:
    - -config.path\s+(?P<Cfgfile>\S+)
```

若进程组不指定 name 字段，默认为 `{{.ExeBase}}`(可执行文件的 basename)

## 关键指标

在实际监控进程时，主要使用的指标就是cpu和内存。process-exporter中进程的指标以 `namedprocess_namegroup` 开头：

- namedprocess\_namegroup\_cpu\_seconds\_total：cpu使用时间，通过mode区分是user还是system
- namedprocess\_namegroup\_memory\_bytes：内存占用，通过 memtype 区分不同的占用类型
- namedprocess\_namegroup\_num\_threads：线程数
- namedprocess\_namegroup\_open\_filedesc：打开的文件句柄数
- namedprocess\_namegroup\_read\_bytes\_total：进程读取的字节数
- namedprocess\_namegroup\_thread\_context\_switches\_total：线程上下文切换统计
- namedprocess\_namegroup\_thread\_count：线程数量统计
- namedprocess\_namegroup\_thread\_cpu\_seconds\_total：线程的cpu使用时间
- namedprocess\_namegroup\_thread\_io\_bytes\_total：线程的io

### cpu

cpu是我们最经常关注的指标，如果使用node-exporter采集节点的指标数据，可以得到机器的cpu占比。

而使用process-exporter采集的是进程的指标，具体来说就是采集/proc/pid/stat中与cpu时间有关的数据：

- 第14个字段：utime，进程在用户态运行的时间，单位为jiffies
- 第15个字段：stime，进程在内核态运行的时间，单位为jiffies
- 第16个字段：cutime，子进程在用户态运行的时间，单位为jiffies
- 第17个字段：cstime，子进程在内核态运行的时间，单位为jiffies

那么通过上述值就可以得到进程的单核CPU占比：

- 进程的单核CPU占比=(utime+stime+cutime+cstime)/时间差
- 进程的单核内核态CPU占比=(stime+cstime)/时间差

因此，进程的单核CPU占比的promsql语句为increase(namedprocess\_namegroup\_cpu\_seconds\_total{mode="user",groupname="procname"}\[30s\])\*100/30，单核内核态CPU占比的promsql语句为increase(namedprocess\_namegroup\_cpu\_seconds\_total{mode="system",groupname="procname"}\[30s\])\*100/30。

注意：实测发现，process-exporter获取的数据与/proc/pid/stat中的有一定差异，需要进一步看下。

### memory

process-exporter采集内存的指标时将内存分成5种类型：

- resident：进程实际占用的内存大小，包括共享库的内存空间，可以从/proc/pid/status中的VmRSS获取
- proportionalResident：与resident相比，共享库的内存空间会根据进程数量平均分配
- swapped：交换空间，系统物理内存不足时，会将不常用的内存页放到硬盘的交换空间，可以从/proc/pid/status中的VmSwap获取
- proportionalSwapped：将可能被交换的内存页按照可能性进行加权平均
- virtual：虚拟内存，描述了进程运行时所需要的总内存大小，包括哪些还没有实际加载到内存中的代码和数据，可以从/proc/pid/status中的VmSize获取

对于一般的程序来说，重点关注的肯定是实际内存，也就是resident和virtual，分别表示实际在内存中占用的空间和应该占用的总空间

# Nginx Exporter

参考：[GitHub 项目](https://github.com/nginxinc/nginx-prometheus-exporter)

端口默认 9113，Grafana 面板 ID：[GitHub](https://raw.githubusercontent.com/nginxinc/nginx-prometheus-exporter/master/grafana/dashboard.json)

```bash
docker run -d --name nginx-exporter\
  --network host \
  nginx/nginx-prometheus-exporter:0.8.0 \
  -nginx.scrape-uri http://localhost:9123/metrics
```

# MySQL Exporter

> 参考：
> - [GitHub 项目](https://github.com/prometheus/mysqld_exporter)

端口默认 9104

Grafana 面板：

- 7362
- 官方提供了一个面板在[这里](https://github.com/prometheus/mysqld_exporter/blob/main/mysqld-mixin/dashboards/mysql-overview.json)

```bash
docker run -d --name mysql-exporter \
  --network host  \
  -e DATA_SOURCE_NAME="root:oc@2020@(localhost:3306)/" \
  prom/mysqld-exporter
```

# Redis Exporter

参考：[GitHub 项目](https://github.com/oliver006/redis_exporter)

端口默认 9121，Grafana 面板 ID：11835

```bash
docker run -d --name redis-exporter \
  --network host \
  oliver006/redis_exporter \
  --redis.password='DesistDaydream' \
  --redis.addr='redis://127.0.0.1:6379'
```

# DELL OMSA 的 Exporter

https://github.com/galexrt/dellhw_exporter

运行方式: docker run -d --name dellhw_exporter --privileged -p 9137:9137 galexrt/dellhw_exporter

如果无法启动，一般都是因为 OMSA 的资源指标采集服务未安装导致。

该 Exporter 会读取 OMSA 的 SNMP 信息转换成 Metrics 并在 9137 端口输出

# IPMI Exporter

https://github.com/prometheus-community/ipmi_exporter

依赖于 FreeIPMI 套件中的工具来实现

# HAProxy 的 Exporter

项目地址：<https://github.com/prometheus/haproxy_exporter>
