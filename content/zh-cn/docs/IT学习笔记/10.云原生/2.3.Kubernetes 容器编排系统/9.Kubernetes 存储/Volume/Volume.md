---
title: Volume
---

# 概述

> 参考：
>
> - [官方文档，概念-存储-卷](https://kubernetes.io/docs/concepts/storage/volumes/)

Volume(卷) 的工作流程：可以把 volume 想象成一个中间人，数据流走向：Container—Volum—StorageResource

要使用 Volume，Pod 的 Manifests 中需要指定要为 Pod 提供卷的类型（`.spec.volumes` 字段）以及将这些卷挂载到容器的位置（`.spec.containers [*].volumeMounts` 字段）。

## Volume 中数据所在目录

/var/lib/kubelet/pods/PodUID/volumes/XXX # pod 中挂载的 volume，其数据都会保存在该目录下，就算是 nfs 这种远程存储，pod 也是会读取该目录下的某个目录，因为这个目录就是挂载 nfs 的目录。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ymmgaq/1616117653866-393d3afc-56b4-4ad5-b1f1-428a013719d3.jpeg)

# Volume 的类型

> 参考：
>
> - [官方文档，概念-存储-卷-卷类型](https://kubernetes.io/docs/concepts/storage/volumes/#types-of-volumes)(全部 Volume 类型的官方介绍)

注意：大量的第三方 In-Tree 类型的第三方卷插件(i.e.非 ConfigMap、Secret 等内置资源的卷插件)将会逐步被弃用，详情参考 [Kubernetes 博客，Kubernetes 1.23: Kubernetes In-Tree to CSI Volume Migration Status Update](https://kubernetes.io/blog/2021/12/10/storage-in-tree-to-csi-migration-status-update/)

现阶段(1.21 版本)，Kubernetes 支持以下类型的 Volume：

- [awsElasticBlockStore](https://kubernetes.io/zh/docs/concepts/storage/volumes/#awselasticblockstore)
- [azureDisk](https://kubernetes.io/zh/docs/concepts/storage/volumes/#azuredisk)
- [azureFile](https://kubernetes.io/zh/docs/concepts/storage/volumes/#azurefile)
- [cephfs](https://kubernetes.io/zh/docs/concepts/storage/volumes/#cephfs)
- [cinder](https://kubernetes.io/zh/docs/concepts/storage/volumes/#cinder)
- [configMap](https://kubernetes.io/zh/docs/concepts/storage/volumes/#configmap) # 一种 Kubernetes 资源，使用 ConfigMap 资源中定义的内容作为 Volume。比如 key 是文件名，value 是文件内容。
- [downwardAPI](https://kubernetes.io/zh/docs/concepts/storage/volumes/#downwardapi)
- [emptyDir](https://kubernetes.io/zh/docs/concepts/storage/volumes/#emptydir) # 把宿主机上的目录作为 Volume
- [fc (光纤通道)](https://kubernetes.io/zh/docs/concepts/storage/volumes/#fc)
- [gcePersistentDisk](https://kubernetes.io/zh/docs/concepts/storage/volumes/#gcepersistentdisk)
- [hostPath](https://kubernetes.io/zh/docs/concepts/storage/volumes/#hostpath) # 把宿主机上的目录作为 Volume
- [iscsi](https://kubernetes.io/zh/docs/concepts/storage/volumes/#iscsi)
- [local](https://kubernetes.io/zh/docs/concepts/storage/volumes/#local) # 把宿主机上的目录作为 Volume
- [nfs](https://kubernetes.io/zh/docs/concepts/storage/volumes/#nfs) # 将 NFS 服务提供的目录作为 Volume
- [persistentVolumeClaim](https://kubernetes.io/zh/docs/concepts/storage/volumes/#persistentvolumeclaim) # 一种 Kubernetes 资源。详见[Persistent Volume 持久卷](https://www.yuque.com/go/doc/33163956) 中关于 PVC 的说明
- [portworxVolume](https://kubernetes.io/zh/docs/concepts/storage/volumes/#portworxvolume)
- [projected](https://kubernetes.io/zh/docs/concepts/storage/volumes/#projected)
- [quobyte](https://kubernetes.io/zh/docs/concepts/storage/volumes/#quobyte)
- [rbd](https://kubernetes.io/zh/docs/concepts/storage/volumes/#rbd)
- [secret](https://kubernetes.io/zh/docs/concepts/storage/volumes/#secret) # 一种 Kubernetes 资源。使用 Secret 资源中定义的内容作为 Volume。比如 key 是文件名，value 是文件内容。
- [storageOS](https://kubernetes.io/zh/docs/concepts/storage/volumes/#storageos)
- [vsphereVolume](https://kubernetes.io/zh/docs/concepts/storage/volumes/#vspherevolume)

上述类型都有各自的用法，常用类型的卷详见本章节下的子章节。

已弃用：[glusterfs](https://kubernetes.io/zh/docs/concepts/storage/volumes/#glusterfs)、[scaleIO](https://kubernetes.io/zh/docs/concepts/storage/volumes/#scaleio)

\*\*
Projected 类型的卷可以将现在的某些卷映射到同一个目录上。
目前，有以下几种 Volume 可以统一映射

1. Secret
2. ConfigMap
3. Downward API
4. ServiceAccountToken

# subPath

参考：<https://kubernetes.io/docs/concepts/storage/volumes/#using-subpath>
注意：如果使用 subPath 挂载为 Container 的 Volume，Kubernetes 不会做自动热更新，因为子路径属于 Volume 下的内容，而不属于 k8s 本身。

有时候，在一个 Pod 中共享一个 Volume 以用于多种用途是很有用的。volumeMounts.subPath 字段可引用 Volume 内的子子路径，而不是该 Volume 的根。

比如

```yaml
    volumeMounts:
    - name: sub
      mountPath: /logs/sub1
      subPath: sub1
    - name: sub
      mountPath: /opt/sub2
      subPath: sub2
  volumes:
  - name: sub
    hostPath:
      path: /root
```

在这个示例中，宿主机的 /root/sub1 目录，将会被挂载到容器中的 /logs/sub1 目录上；宿主机得 /root/sub2 目录，将会被挂载到容器中的 /opt/sub2 目录上。如果宿主机不存在该目录，则会自动创建该目录。

subPath 的行为根据 Volume 类型的不同而不同

1. Volume 类型为 hostPath、nfs、pvc 等 # 挂载该类型 Volume 目录下的子目录
2. Volume 类型为 configmap # 挂载 configmap 中指定的文件

## Pod 中挂载单个文件的方法

在一般情况下 configmap 挂载文件时，会先覆盖掉挂载目录，然后再将 congfigmap 中的内容作为文件挂载进行。如果想不对原来的文件夹下的文件造成覆盖，只是将 configmap 中的每个 key，按照文件的方式挂载到目录下，可以使用 subpath 字段。

1. 创建一个 configmap，让其中 key 为文件名，val 为文件内容
2. 然后创建 configmap 类型的 volumes
3. 配置文件中，
4. 在 pod 的 mountPath 下面使用 subPath 键指定其值为所需要挂载的文件的文件名，即可把文件挂载到相应的目录中，注意 mountPath 需要写文件的绝对路径，不能只写一个目录的路径。这样，在 host 上的 grafana.ini 文件就被挂载到容器中的 /etc/frafana/ 目录下，并且该目录其余文件不受影响

```yaml
kubectl create configmap -n montiroing grafanaini --from-file=grafana.ini
cat XXXX.yaml
.......
        - mountPath: /etc/grafana/grafana.ini # 指定要挂载到容器的哪个路径下
          name: grafana-ini # 指定要使用的volumes名称
          subPath: grafana.ini
      volumes:
      - name: grafana-ini
        configMap: # 指定要使用的configmap名称
          name: grafanaini
.......
```
