---
title: 当master丢失一个节点后，如何恢复
---

# 故障现象

在日常维护中，如果三台 master 节点的其中一个节点故障，并不可恢复，我们如何重新建建立一个 master 节点并重新加入进去呢？

假设曾经有三个节点

1. master-1.tj-test

2. master-2.tj-test

3. master-3.tj-test

其中一个节点丢失后，想要新建一个节点并重新加入集群，但是失败了

# 故障排查

当 master-3 挂掉并不可恢复时，首先需要通过 kubectl delete node master-3.tj-test 命令来删除该节点。然后使用一台新的设备初始化环境，并通过 kubeadm join 命令来加入集群，但是这时候，加入集群是失败的。

因为虽然使用命令删除了 master-3 节点，但是 etcd 集群的 master-3 这个 member 还存在

```shell
[root@master-1 ~]# etcdv3 member list
13b7460f0eebe6ea, started, master-1.tj-test, https://172.38.40.212:2380, https://172.38.40.212:2379
fdddf32d7b4d4498, started, master-3.tj-test, https://172.38.40.214:2380, https://172.38.40.214:2379
fed9f57af62ba6a0, started, master-2.tj-test, https://172.38.40.213:2380, https://172.38.40.213:2379
```

# 故障处理

这时候需要通过 etcdctl 命令 etcdv3 member remove fdddf32d7b4d4498 将该 member 移除，再重新让 master-3 加入集群，就可以了。
