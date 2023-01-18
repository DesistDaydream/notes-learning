---
title: Storage Classes(存储类)
---

# 概述

> 参考：
> - [官方文档,概念-存储-存储类](https://kubernetes.io/docs/concepts/storage/storage-classes/)
> - [官方文档,任务-管理集群-改变默认 StorageClass](https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/)

在介绍的 [PV](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/9.Kubernetes%20 存储/Persistent%20Volume(持久卷).md Volume(持久卷).md) 时有个问题就是管理员需要先创建 pv 固定好容量，再让用户或者开发创建的 PVC 从中挑选，有时候 PVC 申请的时候未必会有满足容量要求的 PV 可以提供，甚至管理员维护大量的 PV 的工作也是非常繁重的。为了实现在创建完 PVC 后，K8S 可以自动创建 PV 的功能，则可以使用 **Storage Class(存储类)** 这个资源对象来满足这类需求。

**Storage Class(存储类)**，就像这个名字一样，Storage Class 是一个抽象的概念，用来抽象存储资源。一般情况都是把同类型的存储归为一类，比如 ssd 类型、hdd 类型等等，也可以按照功能划分，给订单组用的存储，给数据组用的存储等等。说白了，Storage Class 就是一块存储空间。

创建完 StorageClass 后，直接创建 PVC 并指定 storageClassName 参数的值为该 StorageClass 的名字，即可自动生成 PV，而不用手动创建。然后在 pod 中直接使用 PVC 作为 volume 进行挂载即可。

## Storage Class 的实现方式

### Storage Class Name(名字)

**Storage Class Name(存储类的名字)** 是 PV、PVC 选择的标准，PV 与 PVC 总是会选择具有相同名字的 StorageClass 来进行配对。

> StorageClass 可以看做是一个具有无数多个 PV 的存储池，每一个 StorageClass 都可以由管理员自定定义一个 name，管理员可以根据存储性能的高低、备份策略、服务质量等等类别来对 StorageClass 进行命名。k8s 本身并不能理解“类别”到底意味着什么，仅仅是将这些当做 PV、PVC 的特性描述。

### Provisioner(供应器)

**Provisioner(供应器)** 相当于 **StorageClass 的控制器**，也称为 **volume plugin(卷插件)**。**是 Storage Class 实现的主要工具**，用于动态创建 pv。为了让 k8s 可以操作后端 storage，该存储需要支持 restful 风格的 api 接口，否则 k8s 的 StorageClass 资源无法正常工作。但是一般的存储程序，是没有 restful 风格的 api 接口的，这样这种存储就无法与 k8s 进行交互。为了解决这个问题，就需要使用 Provisioner，Provisioner 作为存储与 k8s 之间的翻译者，让存储与 k8s 集群进行交互。

[这里面](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner)有可以使用的 Provisioner 列表，其中分为两类

1. internal，内部 Provisioners。k8s 内建的 Provisioner，名字都以 kubernetes.io 作为前缀
2. external，外部 Provisioners。第三方提供遵循 CSI 标准的 Provisioner。可用的 external provisioners 可以在<https://github.com/kubernetes-incubator/external-storage>这里找到，还有一些在第三方供给者在自己的网站上保存 provisioners 程序
   1. nfs 就是典型的 external provisioners

**parameters(参数)**

- 使用参数来描述该存储中数据的使用方式、如何关联到存储卷、如何关联到真实存储等等。不同的 Provisioner 可用的参数各不相同
- 内部 Provisioners 的参数详见：<https://kubernetes.io/docs/concepts/storage/storage-classes/#parameters>
- 外部 Provisioners 的参数详见：<https://github.com/kubernetes-incubator/external-storage>

### Reclaim Policy(回收策略)

在删除 Storage Class 创建的 PV 时， PV 中的数据应该如何处理，就是由 Reclaim Policy 决定的。一共由两种策略

- Delete # PV 删除时，数据删除
- Retain # PV 删除时，数据保留

## 简单示例

现在使用 NFS Provisioner 作为 Storage Class 的控制器，要根据[此项目](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner)来创建一个 nfs-client 的 Pod，并与后端存储关联。这个 Pod 就是 Provisioner。

```bash
[root@lichenhao rabbitmq-stack-allinone]# kubectl get pod -n sc-provisioner
NAME                                           READY   STATUS    RESTARTS   AGE
test-nfs-client-provisioner-58675b55fb-nprpl   1/1     Running   0          5h2m
```

这个 Provisioner 具有 3 各关键参数：

```yaml
- name: PROVISIONER_NAME # 该 Provisioner 的名字
  value: nfs-storage
- name: NFS_SERVER # 后端存储的地址
  value: 172.19.42.215
- name: NFS_PATH # 后端存储的数据储存路径
  value: /data/test
```

创建一个 Storage Class 资源，并关联 Provisioner

```bash
[root@lichenhao rabbitmq-stack-allinone]# kubectl get storageclasses.storage.k8s.io managed-nfs-storage -o yaml | neat
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: managed-nfs-storage
provisioner: nfs-storage # 指定要关联的 Profisioner，使用 名字 进行关联
parameters:
  archiveOnDelete: "false"
```

这时候，如果创建了一个 PVC，并指定要使用的 Storage Class 的名字，则 PVC 控制器会将请求交给对应的 Storage Class 并创建一个属于该 Storage Class 的 PV，此时由于 Provisioner 一直监控着 该 Storage Class 的状态，所以会触发工作，再根据自身配置规则，在本地目录中，创建对应名称的目录(比如这个 NFS Provisioner 就会创建名为 `名称空间-pvc名-pv名` 这样的目录)

```bash
[root@lichenhao rabbitmq-stack-allinone]# kubectl get pvc -n rabbitmq
NAME                        STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS          AGE
persistence-test-server-0   Bound    pvc-ed8e801f-2659-4829-912b-669145c8396b   10Gi       RWO            managed-nfs-storage   42m
persistence-test-server-1   Bound    pvc-72c85e1c-8c06-45b2-ba46-e223fafd24d5   10Gi       RWO            managed-nfs-storage   42m
persistence-test-server-2   Bound    pvc-9157e421-7150-45ab-8432-2be935dd69ef   10Gi       RWO            managed-nfs-storage   42m
[root@lichenhao rabbitmq-stack-allinone]# kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                                STORAGECLASS          REASON   AGE
pvc-72c85e1c-8c06-45b2-ba46-e223fafd24d5   10Gi       RWO            Delete           Bound    rabbitmq/persistence-test-server-1   managed-nfs-storage            42m
pvc-9157e421-7150-45ab-8432-2be935dd69ef   10Gi       RWO            Delete           Bound    rabbitmq/persistence-test-server-2   managed-nfs-storage            42m
pvc-ed8e801f-2659-4829-912b-669145c8396b   10Gi       RWO            Delete           Bound    rabbitmq/persistence-test-server-0   managed-nfs-storage            42m
[root@nfs-1 test]# ll -h
total 0
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-0-pvc-ed8e801f-2659-4829-912b-669145c8396b
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-1-pvc-72c85e1c-8c06-45b2-ba46-e223fafd24d5
drwxrwxrwx 4 root input 231 Dec  2 21:53 rabbitmq-persistence-test-server-2-pvc-9157e421-7150-45ab-8432-2be935dd69ef
```

## 总结

**Storage Class 与 Provisioner 的关系，就好像 ingress 与 ingress controller 的关系一样。**
PV 与 StorageClass 是集群管理员使用的，PVC 是用户或者开发者使用的。PV 与 StorageClass 就是管理员所管理的存储资源，一个是手动分配，一个是自动分配；PVC 就是用户需要找管理员索要的存储资源。

# 默认 Storage Class

默认的 Storage Class 用于为 [不指定 StorageClassName 的 PVC ](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5f9a47bc12d5ba00014970c7?scroll-to-block=5fc83d08eebc9352f7e0f5c7)动态提供 PV。

列出你的集群中的 StorageClasses：

    kubectl get storageclass

输出类似这样：

    NAME                 PROVISIONER               AGE
    standard (default)   kubernetes.io/gce-pd      1d
    gold                 kubernetes.io/gce-pd      1d

默认 StorageClass 以 `(default)` 标记。

请注意，最多只能有一个 StorageClass 能够被标记为默认。 如果它们中有两个或多个被标记为默认，Kubernetes 将忽略这个注解， 也就是它将表现为没有默认 StorageClass。

**将指定的 StorageClass 标记为非默认**
默认 StorageClass 的注解 `storageclass.kubernetes.io/is-default-class` 设置为 `true`。 注解的其它任意值或者缺省值将被解释为 `false`。要标记一个 StorageClass 为非默认的，你需要改变它的值为 `false`：

    kubectl patch storageclass StorageClassName -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}'

**将指定的 StorageClass 标记为默认**
和前面的步骤类似，需要添加/设置注解 `storageclass.kubernetes.io/is-default-class=true`。

    kubectl patch storageclass StorageClassName -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

# 特殊的 local 类型的 StorageClass

一般情况 StorageClass 都是动态创建 PV 的(i.e.创建一个 PVC，则会自动创建一个 PV 与之绑定，不用管理员手动创建 PV 了)，还有一种静态模式需要手动指定 PV，静态的就是指 local 模式的 storageclass

静态模式的 storageClass 同样需要创建 pv、pvc，不过与普通的 pv、pvc 不同，pv 与 pvc 并不会立刻绑定到一起，而是当 pod 消费 pvc 的时候，才会寻找一个合适的 pv 进行绑定

文末的简单示例，就是一个静态的 StorageClass,需要手动创建 PV 与 PVC，如果没有 StorageClass 的话，那么创建完 PV 与 PVC 后，这俩会自动绑定到一起

示例：

一般情况不会使用这种特殊的 StorageClass，而是使用一个具有正常后端存储的 StorageClass，这样创建 pvc 后可以自动创建 PV，而不再需要手动创建 PV 了

local 表示这个 PV 是使用宿主机的目录来当存储使用的，使用 local 的话，必须使用 nodeAffinity，以便让使用与该 PV 绑定的 PVC 的 pod 可以调度到 local 所在的宿主机上。nodeAffinity 里的 key 与 value 字段，就是某个宿主机的 label，该 PV 会在有相同 label 的宿主机上挂载指定的目录。详见<https://kubernetes.io/docs/concepts/storage/persistent-volumes/#node-affinity>，<https://kubernetes.io/docs/concepts/storage/volumes/#local>

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: test-pv
spec:
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 10Gi
  local:
    path: /mnt/testPV
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - master1
  persistentVolumeReclaimPolicy: Delete
  storageClassName: local-storage
  volumeMode: Filesystem
```

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: test-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 8Gi
  storageClassName: local-storage
  volumeMode: Filesystem
```

<https://kubernetes.io/docs/concepts/storage/storage-classes/#local>

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: local-storage
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```
