---
title: BUG
---

## orphaned pod "XX" found, but volume paths are still present on disk

问题跟踪：[issue #60987](https://github.com/kubernetes/kubernetes/issues/60987)

kubelet 执行逻辑：<https://github.com/kubernetes/kubernetes/blob/release-1.19/pkg/kubelet/kubelet_volumes.go#L173>

解决方式：

- 更新至 1.19.8 版本及以上，[ChangeLog 中提到，在 #95301 Merged](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md#changelog-since-v1198) 中已解决
- 未更新的话，通过 [ali 提供的脚本](https://raw.githubusercontent.com/AliyunContainerService/kubernetes-issues-solution/master/kubelet/kubelet.sh)，进行一些修改，该脚本会手动 umount 和 rm 目录

# 已修复

#

## aggregator_unavailable_apiservice

问题描述：聚合 API 删除之后，依然存在于 kube-apiserver 的 metrics 中，这会导致频繁告警

跟踪连接：<https://github.com/kubernetes/kubernetes/issues/92671>

解决方式：<https://github.com/kubernetes/kubernetes/pull/96421>

**将在 1.20 版本解决**

## Scope libcontainer-21733-systemd-test-default-dependencies.scope has no PIDs. Refusing.

问题跟踪：<https://github.com/kubernetes/kubernetes/issues/71887>

解决方式：（1.16 及以后的版本中，无该问题。主要是 18+版本 docker 无该问题）

忽略该告警：<https://www-01.ibm.com/support> ... .wss?uid=ibm10883724

## Error while processing event ("/sys/fs/cgroup/devices/libcontainer_2434_systemd_test_default.slice": 0x40000100 == IN_CREATE|IN_ISDIR): open /sys/fs/cgroup/devices/libcontainer_2434_systemd_test_default.slice: no such file or directory

解决方式：1.16 及以后版本解决该问题

<https://github.com/kubernetes/kubernetes/issues/76531#issuecomment-548230839>

## Setting volume ownership for XXX and fsGroup set. If the volume has a lot of files then setting volume ownership could be slow,

    Apr 20 11:03:37 lxkubenode01 kubelet[9103]: W0420 11:03:37.275020    9103 volume_linux.go:49] Setting volume ownership for /var/lib/kubelet/pods/63c4a49f-bf23-4b87-989e-102f5fcdb315/volumes/kubernetes.io~secret/seq-token-whpfk and fsGroup set. If the volume has a lot of files then setting volume ownership could be slow, see https://github.com/kubernetes/kubernetes/issues/69699
    Apr 20 11:03:46 lxkubenode01 kubelet[9103]: W0420 11:03:46.198559    9103 volume_linux.go:49] Setting volume ownership for /var/lib/kubelet/pods/cdb79fa0-4942-4211-9261-a4928f872bd6/volumes/kubernetes.io~secret/prometheus-operator-prometheus-node-exporter-token-snbtf and fsGroup set. If the volume has a lot of files then setting volume ownership could be slow, see https://github.com/kubernetes/kubernetes/issues/69699
    Apr 20 11:03:47 lxkubenode01 kubelet[9103]: W0420 11:03:47.201212    9103 volume_linux.go:49] Setting volume ownership for /var/lib/kubelet/pods/8cafe45f-2d14-4eb1-8c38-71a54b34f83c/volumes/kubernetes.io~secret/default-token-9hq84 and fsGroup set. If the volume has a lot of files then setting volume ownership could be slow, see https://github.com/kubernetes/kubernetes/issues/69699

该报警发生于 1.17 及其以后的版本，由于 kubelet 代码更改后，导致频繁得大量刷新类似的日志信息

问题跟踪：
<https://github.com/kubernetes/kubernetes/issues/90293>

解决方式：
治标不治本方法：配置 rayslog 忽略

    cat > /etc/rsyslog.d/ignore-kubelet-volume.conf << EOF
    if (\$programname == "kubelet") and (\$msg contains "Setting volume ownership") then {
      stop
    }
    EOF

根治方法：升级集群
<https://github.com/kubernetes/kubernetes/pull/92878>，pr 已被 merged，在 1.18.8 及 1.19.0 之后得版本已修复，参见：
<https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.18.md#other-cleanup-or-flake>
<https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.19.md#other-cleanup-or-flake-1>

## AggregatedAPIDown 报警

问题描述：添加聚合 API 再删除后，kube-apiserver 中的 metircs 并不会删除，导致一致产生报警。报警规则：

    aggregator_unavailable_apiservice

问题跟踪：<https://github.com/kubernetes/kubernetes/issues/92671>

解决方式：v1.20 有解决 PR

## kubectl get cs 获取信息 unknown

相关报警信息：

    watch chan error: etcdserver: mvcc: required revision has been compacted


    [root@master-1 manifests]# kubectl get cs
    NAME                 AGE
    controller-manager   <unknown>
    scheduler            <unknown>
    etcd-0               <unknown>

问题连接：<https://github.com/kubernetes/kubernetes/issues/83024#issuecomment-559538245>

大意是在 1.17 中得到解决，1.19 及其以后版本不再支持 cs

## 一次 kube-apiserver 响应超时的问题记录

首先是收到 apiserver 响应时间过长的告警，查看日志，发现频繁出现如下内容

    I0331 01:40:30.953289       1 trace.go:116] Trace[2133477734]: "Get" url:/api/v1/namespaces/kube-system (started: 2020-03-31 01:40:21.623714299 +0000 UTC m=+338766.344413381) (total time: 9.329480544s):
    Trace[2133477734]: [9.329404093s] [9.329362028s] About to write a response
    I0331 01:40:36.528652       1 trace.go:116] Trace[1431450424]: "Get" url:/api/v1/namespaces/default (started: 2020-03-31 01:40:28.063278623 +0000 UTC m=+338772.783977705) (total time: 8.465254793s):
    Trace[1431450424]: [8.465073901s] [8.465027207s] About to write a response
    I0331 01:40:37.333718       1 trace.go:116] Trace[1319947973]: "Get" url:/api/v1/namespaces/kube-public (started: 2020-03-31 01:40:30.954280125 +0000 UTC m=+338775.674979196) (total time: 6.379382999s):
    Trace[1319947973]: [6.37929667s] [6.379253238s] About to write a response

查看 etcd 日志，发现很多 took too long 的信息：

```bash
etcdserver: read-only range request ..... took too long  to execute
```

怀疑有可能是磁盘性能问题，使用 dd 测试，发现只有 10+M/s，查看 Raid 卡发现模式是直写，更换 Raid 卡，修改模式为强制回写。问题解决
