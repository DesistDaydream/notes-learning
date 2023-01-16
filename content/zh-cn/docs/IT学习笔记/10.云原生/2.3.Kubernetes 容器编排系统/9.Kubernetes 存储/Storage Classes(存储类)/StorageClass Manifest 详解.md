---
title: StorageClass Manifest 详解
---

# 概述

> 参考：[API 文档](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#storageclass-v1-storage-k8s-io)

其中带有 `-required-` 标志的为必须的字段

## apiVersion: storage.k8s.io/v1

## kind: StorageClass

## metadata:

- name: <STRING> # 该 StorageClass 的名字

## provisioner: <STRING> -required- # 指定要使用的 provisioner。

## parameters \<map\[string]string> # Provisioner 的配置参数，不同的 Provisioner 具有不同的参数

### NFS Parameters

- **archiveOnDelete: false** # PV 删除时，是否存档数据。如果存档数据，则会在 PV 关联的目录名前加上 archived 字符串，以便保存数据。效果如下：

```bash
[root@nfs-1 test]# ll -h
total 0
drwxrwxrwx 4 root input 179 Dec  2 17:02 archived-rabbitmq-persistence-test-server-0-pvc-9e1aabf2-a072-44d1-92df-9cd0864f9fda
drwxrwxrwx 4 root input 179 Dec  2 17:02 archived-rabbitmq-persistence-test-server-1-pvc-0409c70e-be04-43e1-9b4c-17b96930cb26
drwxrwxrwx 4 root input 179 Dec  2 17:02 archived-rabbitmq-persistence-test-server-2-pvc-01937680-cc6c-4118-86fc-ed420fdc275b
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-0-pvc-ed8e801f-2659-4829-912b-669145c8396b
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-1-pvc-72c85e1c-8c06-45b2-ba46-e223fafd24d5
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-2-pvc-9157e421-7150-45ab-8432-2be935dd69ef
# 可以看到，前面三个时之前层经删除过的 PV，archiveOnDelete 参数改为 true 后，数据并不会被删除
```

## reclaimPolicy: <STRING> # 回收策略

StorageClass 动态创建的 PV 的回收策略。

- Retain(保留) # 关联的 PVC 删除后，PV 保留，并变为 Released 状态。
  - 再次创建 PVC 后，并不会绑定到被保留的 PV 上。
- Delete(删除) # 关联的 PVC 删除后，PV 删除。`默认值`
  - PV 中的数据将会根据参数的 archiveOnDelete 的值来决定如何处理。

## allowVolumeExpansion: <BOOLEAN> # 存储类是否允许卷扩展

## allowedTopologies <\[]Object>

Restrict the node topologies where volumes can be dynamically provisioned. Each volume plugin defines its own supported topology specifications. An empty TopologySelectorTerm list means there is no topology restriction. This field is only honored by servers that enable the VolumeScheduling feature.

## mountOptions <\[]string>

Dynamically provisioned PersistentVolumes of this storage class are created with these mountOptions, e.g. \["ro", "soft"]. Not validated - mount of the PVs will simply fail if one is invalid.

## volumeBindingMode: <STRING>

VolumeBindingMode indicates how PersistentVolumeClaims should be provisioned and bound. When unset, VolumeBindingImmediate is used. This field is only honored by servers that enable the VolumeScheduling feature.

# StorageClass Manifest 样例

    apiVersion: storage.k8s.io/v1
    kind: StorageClass
    metadata:
      name: managed-nfs-storage
    provisioner: nfs-storage
    parameters:
      archiveOnDelete: "false"
