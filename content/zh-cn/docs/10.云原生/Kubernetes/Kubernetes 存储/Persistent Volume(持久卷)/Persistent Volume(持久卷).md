---
title: Persistent Volume(持久卷)
weight: 1
---

# PersistentVolume(持久卷)

<https://kubernetes.io/docs/concepts/storage/persistent-volumes/>
kubernetes 的 volume 有一个问题就是不够灵活，且使用 volume 的用户必须要明确知道要使用的 volume 后端是用的什么类型的存储系统(例如 NFS 就需要配置 IP 和 PATH)。这与 kubernetes 的向用户和开发隐藏底层架构的目的有所背离，对存储资源最好的使用方式是能向计算资源一样，用户和开发人员无需了解 pod 资源究竟运行于哪个节点，也无需了解存储系统是什么类型的设备以及位于何处。他们只需要提出容量大小的需求(i.e.PVC)，k8s 管理员可以为其分配所需的空间。这就是 PV 与 PVC 的作用，抽象了底层的存储，使得存储系统的使用和管理两个职能互相解耦。这就好比创建一台虚拟机，并需求 20G 的存储空间，然后虚拟机管理系统就会自动创建出来，而不用去手动指定使用哪个存储空间了。

**PersistentVolume(持久卷,简称 PV)** 与 **PersistentVolumeClaims(持久卷申请,简称 PVC)** 是 kubernetes 中的一个 **Resource(资源)**，PVC 与 volume 不同，volume 的定义需要写进 Pod 的 manifest 中。而 PVC 对于 volume 来说就是一种 "volume 的类型"。

1. PV 就相当于虚拟机中的存储卷，虚拟了宿主机或者远程存储的存储资源。PVC 就相当于虚拟机里的一个物理磁盘。虚拟机里的磁盘其实就是通过存储卷来实现的，这与 PV 与 PVC 的关系基本一致。而下文讲的 StorageClass(存储类) 则相当于虚拟机中的存储池的概念了。
2. PV 与 PVC 为用户和管理员提供了一个 API 接口，抽象定义了存储的消费者-生产者模型(即从 k8s 系统来说，管理员使用 PV 生产一个 storage 资源，用户使用 PVC 消费一个 storage 资源)

## PV 与 PVC 的生命周期

详见：[PV 与 PVC 的状态变化示例](/docs/10.云原生/Kubernetes/Kubernetes%20存储/Persistent%20Volume(持久卷)/PV%20与%20PVC%20的状态变化示例.md)

PV 是集群中的资源。 PVC 是对这些资源的请求，并且还充当对资源的声明检查。 PV 和 PVC 之间的交互遵循以下生命周期：

1. Provisioning(供应)
2. Binding(绑定)
3. Using(使用)
4. Storage Object in User(保护使用中存储对象)
5. Protection(回收)
6. Reserving a PersistentVolume(预留 PV)
7. Expanding Persistent Volumes Claims(扩充 PVC 申领)

PV、PVC 在 Pod 中实现方式简述：

PV 与 PVC 的工作流程：Container—Volume—PVC—PV—StorageResource

1. 创建 PV 资源类型的对象，每个 PV 包括几个属性：名字、capacity(容量)、reclaim policy(回收策略) 等等
2. 创建 PVC 资源类型的对象，每个 PVC 中包括几个属性：会话模式、名字、所需求的容量的大小等等
3. 当 PVC 创建后，会自动选择一个满足 PVC 所需求容量的最小 PV 来与 PVC 进行绑定，此时 PV 的状态变为 bond
4. 创建 Pod，在 Pod 中定义 Volume 类型为 PVC，并指定 PVC 名字以选择一个 PVC 挂载上来，然后在 container 中指定要挂载的 volume 名以及在 container 中的挂载路径

### Provisioning(供应)

**Provisioning(供应)** 就是指供应 PV，也就是**创建 PV** 的意思

可以通过两种方式创建 PV：**statically(静态地)** 或 **dynamically(动态地)**

- **Static(静态)** # 手动编写 PV 的 manifest，并 apply 到集群中。
- **Dynamic(动态)** # 当管理员创建的所有 Statci PV(静态 PV) 都不匹配用户的 PVC 时，集群将会尝试为 PVC 动态地提供一个卷
  - 这种动态供应必须基于 StorageClasses：PVC 必须配置指定的 StorageClass，并且集群必须存在与 PVC 配置的同名的 StorageClass，才能进行动态供应。
  - 在声明一个 PVC 对象时，如果 manifest 中 storageClassName 字段为空，则该声明会禁用动态供应。除非集群中有一个默认的 StorageClass。

### Binding(绑定)

### Using(使用)

### Storage Object in User Protection(保护使用中存储对象)

### Reclaiming(回收)

当删除 PVC 后，与之关联的 PV 会变为 **Released(已释放)** 状态。此时，控制器根据一个**Reclaim Policy(回收策略)**来处理已释放的 PV。目前, 回收策略包括三种：

1. Retained(保留)
2. Recycled(回收)
3. Deleted(删除)。

也就是说，PV 可以进行这三种处理方式。

### Reserving a PersistentVolume(预留 PV)

### Expanding Persistent Volumes Claims(扩充 PVC 申领)

# 动态 PV 交互流程

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/okh44l/1616117503785-5b51e61b-c925-49f1-97e4-d97e9e020268.jpeg)

1. 用户创建一个包含 PVC 的 Pod
2. PV Controller 会观察 ApiServer，如果它发现一个 PVC 已经创建完毕但仍然是未绑定的状态，它就会试图把一个 PV 和 PVC 绑定
3. Provision 就是从远端上一个具体的存储介质创建一个 Volume，并且在集群中创建一个 PV 对象，然后将此 PV 和 PVC 进行绑定
4. Scheduler 进行多个维度考量完成后,把 Pod 调度到一个合适的 Node
5. Kubelet 不断 watch APIServer 是否有 Pod 要调度到当前所在节点
6. Pod 调度到某个节点之后，它所定义的 PV 还没有被挂载（Attach），此时 AD Controller 就会调用 VolumePlugin，把远端的 Volume 挂载到目标节点中的设备上（/dev/vdb）；当 Volum Manager 发现一个 Pod 调度到自己的节点上并且 Volume 已经完成了挂载，它就会执行 mount 操作，将本地设备（也就是刚才得到的/dev/vdb）挂载到 Pod 在节点上的一个子目录中
7. 启动容器,并将已经挂载到本地的 Volume 映射到容器中

# PV、PVC 的 manifest 样例

## 一个简单的创建 PV、PVC,然后 pod 使用 PVC 的实例

<https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes> 该链接下有 PV 的每个键值对的解释

<https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims> 该链接下有 PVC 的每个键值对的机制

注意：

1. storageClassName 相当于 PV 的一个标签，当 PVC 申请的时候，只会选择具有相同 storageClassName 的 PV 来绑定，如果集群中存在 StorageClass 资源的话，则会在对应的 StorageClass 中创建该 PV。该参数的值可以为空，如果不指定 StorageClass，则 PVC 也不能指定，才会绑定到该 PV 上。此时 PV 不使用任何 StorageClass。StorageClass 概念详见：Storage Classes 存储类

该示例先创建 PV，然后创建 PVC，最后让 pod 使用该 PVC，node-1.tj-test 节点上 /opt/myapp 目录会挂载到容器中的 /persistentVolume 目录上

    apiVersion: v1
    kind: PersistentVolume
    metadata:
      labels:
        app: myapp
      name: myapp-pv
    spec:
      accessModes:
      - ReadWriteOnce
      capacity:
        storage: 10Gi
      local:
        path: /opt/myapp
      nodeAffinity:
        required:
          nodeSelectorTerms:
          - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
              - node-1.tj-test
    ---
    kind: PersistentVolumeClaim
    apiVersion: v1
    metadata:
      name: myapp-pvc
      labels:
        app: myapp
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 10Gi
      selector:
        matchLabels:
          app: myapp
    ---
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: myapp
      labels:
        name: myapp
    spec:
      replicas: 1
      selector:
        matchLabels:
          name: myapp
      template:
        metadata:
          name: myapp
          labels:
            name: myapp
        spec:
          containers:
          - name: myapp
            image: lchdzh/network-test
            volumeMounts:
            - mountPath: "/persistentVolume"
              name: myapp-pvc
          volumes:
          - name: myapp-pvc
            persistentVolumeClaim:
              claimName: test-claim
