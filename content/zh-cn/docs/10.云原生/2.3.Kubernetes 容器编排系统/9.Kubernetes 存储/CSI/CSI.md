---
title: CSI
weight: 1
---

# 概述

> 参考：
>
> - [公众号-阿里云云原生，一文读懂容器存储接口 CSI](https://mp.weixin.qq.com/s/A9xWKMmrxPyOEiCs_sicYQ)

**导读：** 在[《一文读懂 K8s 持久化存储流程》](https://mp.weixin.qq.com/s?__biz=MzUzNzYxNjAzMg==&mid=2247490043&idx=1&sn=c09ad4a9bc790f4b742abd8ca1301ffb&scene=21#wechat_redirect)一文我们重点介绍了 K8s 内部的存储流程，以及 PV、PVC、StorageClass、Kubelet 等之间的调用关系。接下来本文将将重点放在 CSI（Container Storage Interface）容器存储接口上，探究什么是 CSI 及其内部工作原理。

# 背景

K8s 原生支持一些存储类型的 PV，如 iSCSI、NFS、CephFS 等等（详见链接），这些 in-tree 类型的存储代码放在 Kubernetes 代码仓库中。这里带来的问题是 K8s 代码与三方存储厂商的代码**强耦合**：

- 更改 in-tree 类型的存储代码，用户必须更新 K8s 组件，成本较高
- in-tree 存储代码中的 bug 会引发 K8s 组件不稳定
- K8s 社区需要负责维护及测试 in-tree 类型的存储功能
- in-tree 存储插件享有与 K8s 核心组件同等的特权，存在安全隐患
- 三方存储开发者必须遵循 K8s 社区的规则开发 in-tree 类型存储代码

CSI 容器存储接口标准的出现解决了上述问题，将三方存储代码与 K8s 代码解耦，使得三方存储厂商研发人员只需实现 CSI 接口（无需关注容器平台是 K8s 还是 Swarm 等）。

# CSI 核心流程介绍

在详细介绍 CSI 组件及其接口之前，我们先对 K8s 中 CSI 存储流程进行一个介绍。[《一文读懂 K8s 持久化存储流程》](https://mp.weixin.qq.com/s?__biz=MzUzNzYxNjAzMg==&mid=2247490043&idx=1&sn=c09ad4a9bc790f4b742abd8ca1301ffb&scene=21#wechat_redirect)一文介绍了 K8s 中的 Pod 在挂载存储卷时需经历三个的阶段：Provision/Delete（创盘/删盘）、Attach/Detach（挂接/摘除）和 Mount/Unmount（挂载/卸载），下面以图文的方式讲解 K8s 在这三个阶段使用 CSI 的流程。

### 1. Provisioning Volumes

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrBeY2LVjZxMVenKvU5ZKicmSw8VHGuc0Oz4e5pf4p3Obx6I0IEpTRyPg/640?wx_fmt=jpeg)

1. **集群管理员**创建 StorageClass 资源，该 StorageClass 中包含 CSI 插件名称（provisioner:pangu.csi.alibabacloud.com）以及存储类必须的参数（parameters: type=cloud_ssd）。sc.yaml 文件如下：

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrANJr4ZpicWLs1OCnFLhicFAU62B2k6A0ziarheficmo68kCzlCWQM9HHicQ/640?wx_fmt=png)

2. **用户**创建 PersistentVolumeClaim 资源，PVC 指定存储大小及 StorageClass（如上）。pvc.yaml 文件如下：

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrutp4BzfXLmBicwlbZKic2icwyfdrONeCXriaL69aHV1dRGsV45xs8FwrAA/640?wx_fmt=png)

3. **卷控制器（PersistentVolumeController）**观察到集群中新创建的 PVC 没有与之匹配的 PV，且其使用的存储类型为 out-of-tree，于是为 PVC 打 annotation：volume.beta.kubernetes.io/storage-provisioner=\[out-of-tree CSI 插件名称\]（本例中即为 provisioner:pangu.csi.alibabacloud.com）。

4. **External Provisioner 组件**观察到 PVC 的 annotation 中包含 "volume.beta.kubernetes.io/storage-provisioner" 且其 value 是自己，于是开始创盘流程。

- 获取相关 StorageClass 资源并从中获取参数（本例中 parameters 为  type=cloud_ssd），用于后面 CSI 函数调用。

- 通过 unix domain socket 调用**外部 CSI 插件**的 **CreateVolume 函数**。

5. **外部 CSI 插件**返回成功后表示盘创建完成，此时 **External Provisioner 组件**会在集群创建一个 PersistentVolume 资源。

6. **卷控制器**会将 PV 与 PVC 进行绑定。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrJDTQYxaVvDX4VPpF2qCeXNPAuG4lYLYsO91gG51s5mrrQtoAkFK0gQ/640?wx_fmt=png)

### **2. Attaching Volumes**

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrDnVdvz0ibIzE0rqPvICXvWdHGosbgU2UCDGtImAVbCYNSQ6gEgoXfow/640?wx_fmt=jpeg)

1. **AD 控制器（AttachDetachController）**观察到使用 CSI 类型 PV 的 Pod 被调度到某一节点，此时 **AD 控制器**会调用**内部 in-tree CSI 插件（csiAttacher）**的 Attach 函数。

2. **内部 in-tree CSI 插件（csiAttacher）**会创建一个 VolumeAttachment 对象到集群中。

3. **External Attacher**观察到该 VolumeAttachment 对象，并调用**外部 CSI** **插件**的 **ControllerPublish 函数**以将卷挂接到对应节点上。**外部 CSI 插件**挂载成功后，**External Attacher** 会更新相关 VolumeAttachment 对象的 .Status.Attached 为 true。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPr3VQrSYwJScqw9NmzQWEb2EiavN62qibOtY8UqsHf7pJOgd0dGfH51J7A/640?wx_fmt=png)

4. **AD 控制器内部 in-tree CSI 插件（csiAttacher）**观察到 VolumeAttachment 对象的 .Status.Attached 设置为 true，于是更新 **AD 控制器**内部状态（ActualStateOfWorld），该状态会显示在 Node 资源的 .Status.VolumesAttached 上。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrCEstrgCCXic5WRlHKcqo9tsay8vlFtC16iaw7ibNgztibNUBF7T0Wa5GjA/640?wx_fmt=png)

### **3. Mounting Volumes**

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrLCMro57TU3k6Y4n5rJshOXjyuUDSxm07DFrkwFc0PrJnQWCIuNjJMQ/640?wx_fmt=jpeg)

1. **Volume Manager（Kubelet 组件）**观察到有新的使用 CSI 类型 PV 的 Pod 调度到本节点上，于是调用**内部 in-tree CSI 插件（csiAttacher）**的 WaitForAttach 函数。

2. **内部 in-tree CSI 插件（csiAttacher）**等待集群中 VolumeAttachment 对象状态 .Status.Attached 变为 true。

3. **in-tree CSI 插件（csiAttacher）**调用 MountDevice 函数，该函数内部通过 unix domain socket 调用**外部 CSI 插件**的 **NodeStageVolume 函数**；之后**插件（csiAttacher）**调用**内部 in-tree CSI 插件（csiMountMgr）**的 SetUp 函数，该函数内部会通过 unix domain socket 调用**外部 CSI 插件**的 **NodePublishVolume 函数**。

### **4. Unmounting Volumes**

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrQe0rqRuybVvWT35SIHKS5SicLNrkyy0NeEWssojJpqyKw2cScnET8cg/640?wx_fmt=jpeg)

1. **用户**删除相关 Pod。

2. **Volume Manager（Kubelet 组件）**观察到包含 CSI 存储卷的 Pod 被删除，于是调用**内部 in-tree CSI 插件（csiMountMgr）**的 TearDown 函数，该函数内部会通过 unix domain socket 调用**外部 CSI 插件**的 **NodeUnpublishVolume 函数**。

3. **Volume Manager（Kubelet 组件）**调用**内部 in-tree CSI 插件（csiAttacher）**的 UnmountDevice 函数，该函数内部会通过 unix domain socket 调用**外部 CSI 插件**的 **NodeUnpublishVolume 函数**。

### **5. Detaching Volumes**

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrXZxdCa6icByCacuAZetQqxVfVRNFx3EIt4g8DvFrNPmvyIjSHGSKRzg/640?wx_fmt=jpeg)

1. **AD 控制器**观察到包含 CSI 存储卷的 Pod 被删除，此时该控制器会调用**内部 in-tree CSI 插件（csiAttacher）**的 Detach 函数。

2. **csiAttache****r** 会删除集群中相关 VolumeAttachment 对象（但由于存在 finalizer，va 对象不会立即删除）。

3. **External Attacher** 观察到集群中 VolumeAttachment 对象的 DeletionTimestamp 非空，于是调用**外部 CSI 插件**的 **ControllerUnpublish 函数**以将卷从对应节点上摘除。**外部 CSI 插件**摘除成功后，**External Attacher** 会移除相关 VolumeAttachment 对象的 finalizer 字段，此时 VolumeAttachment 对象被彻底删除。

4. **AD 控制器**中**内部 in-tree CSI 插件（csiAttacher）**观察到 VolumeAttachment 对象已删除，于是更新 **AD 控制器**中的内部状态；同时 **AD 控制器**更新 Node 资源，此时 Node 资源的 .Status.VolumesAttached 上已没有相关挂接信息。

### **6. Deleting Volumes**

![](https://mmbiz.qpic.cn/mmbiz_jpg/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrUnvXdic5KwU7cRfGTUkOG17ol0Dnicj177dzzzRLiam3XLKZUvPQwcicjQ/640?wx_fmt=jpeg)

1. **用户**删除相关 PVC。

2. **External Provisioner 组件**观察到 PVC 删除事件，根据 PVC 的回收策略（Reclaim）执行不同操作：

- Delete：调用**外部 CSI 插件**的 **DeleteVolume 函数**以删除卷；一旦卷成功删除，**Provisioner** 会删除集群中对应 PV 对象。

- Retain：**Provisione** 不执行卷删除操作。

# CSI Sidecar 组件介绍

为使 K8s 适配 CSI 标准，社区将与 K8s 相关的存储流程逻辑放在了 CSI Sidecar 组件中。

## 1. Node Driver Registrar

### **1）功能**

**Node-Driver-Registrar 组件**会将**外部 CSI 插件**注册到 **Kubelet**，从而使 **Kubelet** 通过特定的 Unix Domain Socket 来调用**外部 CSI 插件函数**（Kubelet 会调用外部 CSI 插件的 NodeGetInfo、NodeStageVolume、NodePublishVolume、NodeGetVolumeStats 等函数）。

### **2）原理**

**Node-Driver-Registrar 组件**通过Kubelet 外部插件注册机制实现注册，注册成功后：

- **Kubelet** 为本节点 Node 资源打 annotation：**Kubelet** 调用**外部 CSI 插件**的 **NodeGetInfo 函数**，其返回值 \[nodeID\]、\[driverName\] 将作为值用于 "csi.volume.kubernetes.io/nodeid" 键。

- **Kubelet** 更新 Node Label：将 **NodeGetInfo 函数**返回的 \[AccessibleTopology\] 值用于节点的 Label。

- **Kubelet** 更新 Node Status：将 **NodeGetInfo 函数**返回的 maxAttachLimit（节点最大可挂载卷数量）更新到 Node 资源的 Status.Allocatable：attachable-volumes-csi-\[driverName\]=\[maxAttachLimit\]。

![](https://mmbiz.qpic.cn/mmbiz_png/yvBJb5IiafvnZrRu67Cf0RQ5ToaxdqVPrBmrHGPKfvIpb11zVlVHicPqfFjRCCbUNSgNsufaCca3U0mgWVp3qQGg/640?wx_fmt=png)

- **Kubelet** 更新 CSINode 资源（没有则创建）：将 \[driverName\]、\[nodeID\]、\[maxAttachLimit\]、\[AccessibleTopology\] 更新到 Spec 中（拓扑仅保留 Key 值）。

## 2. External Provisioner

### **1）功能**

创建/删除实际的存储卷，以及代表存储卷的 PV 资源。

### **2）原理**

**External-Provisioner** 在启动时需指定参数 \-\- provisioner，该参数指定 Provisioner 名称，与 StorageClass 中的 provisioner 字段对应。

**External-Provisioner** 启动后会 watch 集群中的 PVC 和 PV 资源。

对于集群中的 PVC 资源：

- 判断 PVC 是否需要动态创建存储卷，标准如下：

- PVC 的 annotation 中是否包含 "volume.beta.kubernetes.io/storage-provisioner" 键（由卷控制器创建），并且其值是否与 Provisioner 名称相等。

- PVC 对应 StorageClass 的 VolumeBindingMode 字段若为 WaitForFirstConsumer，则 PVC 的 annotation 中必须包含 "volume.kubernetes.io/selected-node" 键（详见调度器如何处理 WaitForFirstConsumer），且其值不为空；若为 Immediate 则表示需要 Provisioner 立即提供动态存储卷。

- 通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **CreateVolume 函数**。

- 创建 PV 资源，PV 名称为 \[Provisioner 指定的 PV 前缀\] - \[PVC uuid\]。

对于集群中的 PV 资源：

- 判断 PV 是否需要删除，标准如下：

- 判断其 .Status.Phase 是否为 Release。

- 判断其 .Spec.PersistentVolumeReclaimPolicy 是否为 Delete。

- 判断其是否包含 annotation（pv.kubernetes.io/provisioned-by），且其值是否为自己。

- 通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **DeleteVolume 接口**。

- 删除集群中的 PV 资源。

## 3. External Attacher

### **1）功能**

挂接/摘除存储卷。

### **2）原理**

**External-Attacher**内部会时刻 watch 集群中的 VolumeAttachment 资源和 PersistentVolume 资源。

对于 VolumeAttachment 资源：

- 从 VolumeAttachment 资源中获得 PV 的所有信息，如 volume ID、node ID、挂载 Secret 等。

- 判断 VolumeAttachment 的 DeletionTimestamp 字段是否为空来判断其为卷挂接或卷摘除：若为卷挂接则通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **ControllerPublishVolume 接口**；若为卷摘除则通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **ControllerUnpublishVolume 接口**。

对于 PersistentVolume 资源：

- 在挂接时为相关 PV 打上 Finalizer：external-attacher/\[driver 名称\]。

- 当 PV 处于删除状态时（DeletionTimestamp 非空），删除 Finalizer：external-attacher/\[driver 名称\]。

## 4. External Resizer

### **1）功能**

扩容存储卷。

### **2）原理**

**External-Resizer** 内部会 watch 集群中的 PersistentVolumeClaim 资源。

对于 PersistentVolumeClaim 资源：

- 判断 PersistentVolumeClaim 资源是否需要扩容：PVC 状态需要是 Bound 且 .Status.Capacity 与 .Spec.Resources.Requests 不等。

- 更新 PVC 的 .Status.Conditions，表明此时处于 Resizing 状态。

- 通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **ControllerExpandVolume 接口**。

- 更新 PV 的 .Spec.Capacity。

- 若 CSI 支持文件系统在线扩容，ControllerExpandVolume 接口返回值中 NodeExpansionRequired 字段为 true，**External-Resizer** 更新 PVC 的 .Status.Conditions 为 FileSystemResizePending 状态；若不支持则扩容成功，**External-Resizer** 更新 PVC 的 .Status.Conditions 为空，且更新 PVC 的 .Status.Capacity。

**Volume Manager（Kubelet 组件）**观察到存储卷需在线扩容，于是通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **NodeExpandVolume 接口**实现文件系统扩容。

## 5. livenessprobe

### **1）功能**

检查 CSI 插件是否正常。

### **2）原理**

通过对外暴露一个 / healthz HTTP 端口以服务 kubelet 的探针探测器，内部是通过特定的 Unix Domain Socket 调用**外部 CSI 插件**的 **Probe 接口**。

# CSI 接口介绍

三方存储厂商需实现 CSI 插件的三大接口：**IdentityServer、ControllerServer、NodeServer**。

## 1. IdentityServer

IdentityServer 主要用于认证 CSI 插件的身份信息。

```typescript

type IdentityServer interface {

  GetPluginInfo(context.Context, *GetPluginInfoRequest) (*GetPluginInfoResponse, error)

    GetPluginCapabilities(context.Context, *GetPluginCapabilitiesRequest) (*GetPluginCapabilitiesResponse, error)

    Probe(context.Context, *ProbeRequest) (*ProbeResponse, error)
}
```

## 2. ControllerServer

ControllerServer 主要负责存储卷及快照的创建/删除以及挂接/摘除操作。

```typescript

type ControllerServer interface {

  CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)

    DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)

    ControllerPublishVolume(context.Context, *ControllerPublishVolumeRequest) (*ControllerPublishVolumeResponse, error)

    ControllerUnpublishVolume(context.Context, *ControllerUnpublishVolumeRequest) (*ControllerUnpublishVolumeResponse, error)

    ValidateVolumeCapabilities(context.Context, *ValidateVolumeCapabilitiesRequest) (*ValidateVolumeCapabilitiesResponse, error)

    ListVolumes(context.Context, *ListVolumesRequest) (*ListVolumesResponse, error)

    GetCapacity(context.Context, *GetCapacityRequest) (*GetCapacityResponse, error)

    ControllerGetCapabilities(context.Context, *ControllerGetCapabilitiesRequest) (*ControllerGetCapabilitiesResponse, error)

    CreateSnapshot(context.Context, *CreateSnapshotRequest) (*CreateSnapshotResponse, error)

    DeleteSnapshot(context.Context, *DeleteSnapshotRequest) (*DeleteSnapshotResponse, error)

    ListSnapshots(context.Context, *ListSnapshotsRequest) (*ListSnapshotsResponse, error)

    ControllerExpandVolume(context.Context, *ControllerExpandVolumeRequest) (*ControllerExpandVolumeResponse, error)
}
```

## 3. NodeServer

NodeServer 主要负责存储卷挂载/卸载操作。

```typescript

type NodeServer interface {

  NodeStageVolume(context.Context, *NodeStageVolumeRequest) (*NodeStageVolumeResponse, error)

    NodeUnstageVolume(context.Context, *NodeUnstageVolumeRequest) (*NodeUnstageVolumeResponse, error)

    NodePublishVolume(context.Context, *NodePublishVolumeRequest) (*NodePublishVolumeResponse, error)

    NodeUnpublishVolume(context.Context, *NodeUnpublishVolumeRequest) (*NodeUnpublishVolumeResponse, error)

    NodeGetVolumeStats(context.Context, *NodeGetVolumeStatsRequest) (*NodeGetVolumeStatsResponse, error)

    NodeExpandVolume(context.Context, *NodeExpandVolumeRequest) (*NodeExpandVolumeResponse, error)

    NodeGetCapabilities(context.Context, *NodeGetCapabilitiesRequest) (*NodeGetCapabilitiesResponse, error)

    NodeGetInfo(context.Context, *NodeGetInfoRequest) (*NodeGetInfoResponse, error)
}
```

# K8s CSI API 对象

K8s 为支持 CSI 标准，包含如下 API 对象：

- CSINode

- CSIDriver

- VolumeAttachment

## 1. CSINode

```properties
apiVersion: storage.k8s.io/v1beta1
kind: CSINode
metadata:
  name: node-10.212.101.210
spec:
  drivers:
  - name: yodaplugin.csi.alibabacloud.com
    nodeID: node-10.212.101.210
    topologyKeys:
    - kubernetes.io/hostname
  - name: pangu.csi.alibabacloud.com
    nodeID: a5441fd9013042ee8104a674e4a9666a
    topologyKeys:
    - topology.pangu.csi.alibabacloud.com/zone
```

作用：

1. 判断**外部 CSI 插件**是否注册成功。在 Node Driver Registrar 组件向 Kubelet 注册完毕后，Kubelet 会创建该资源，故不需要显式创建 CSINode 资源。

2. 将 Kubernetes 中 Node 资源名称与三方存储系统中节点名称（nodeID）一一对应。此处 **Kubelet** 会调用**外部 CSI 插件** NodeServer 的 **GetNodeInfo 函数**获取 nodeID。

3. 显示卷拓扑信息。CSINode 中 topologyKeys 用来表示存储节点的拓扑信息，卷拓扑信息会使得 **Scheduler** 在 Pod 调度时选择合适的存储节点。

## 2. CSIDriver

```properties
apiVersion: storage.k8s.io/v1beta1
kind: CSIDriver
metadata:
  name: pangu.csi.alibabacloud.com
spec:

  attachRequired: true

  podInfoOnMount: true

  volumeLifecycleModes:
  - Persistent
```

作用：

1. 简化**外部 CSI 插件**的发现。由集群管理员创建，通过 kubectl get csidriver 即可得知环境上有哪些 CSI 插件。

2. 自定 义Kubernetes 行为，如一些外部 CSI 插件不需要执行卷挂接（VolumeAttach）操作，则可以设置 .spec.attachRequired 为 false。

## 3. VolumeAttachment

```properties
apiVersion: storage.k8s.io/v1
kind: VolumeAttachment
metadata:
  annotations:
    csi.alpha.kubernetes.io/node-id: 21481ae252a2457f9abcb86a3d02ba05
  finalizers:
  - external-attacher/pangu-csi-alibabacloud-com
  name: csi-0996e5e9459e1ccc1b3a7aba07df4ef7301c8e283d99eabc1b69626b119ce750
spec:
  attacher: pangu.csi.alibabacloud.com
  nodeName: node-10.212.101.241
  source:
    persistentVolumeName: pangu-39aa24e7-8877-11eb-b02f-021234350de1
status:
  attached: true
```

作用：VolumeAttachment 记录了存储卷的挂接/摘除信息以及节点信息。

# 支持特性

## 1. 拓扑支持

在 StorageClass 中有 AllowedTopologies 字段：

```properties
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: csi-pangu
provisioner: pangu.csi.alibabacloud.com
parameters:
  type: cloud_ssd
volumeBindingMode: Immediate
allowedTopologies:
- matchLabelExpressions:
  - key: topology.pangu.csi.alibabacloud.com/zone
    values:
    - zone-1
    - zone-2
```

**外部 CSI 插件**部署后会为每个节点打标，打标内容 **NodeGetInfo 函数**返回的 \[AccessibleTopology\] 值（详见 Node Driver Registrar 部分）。

**External Provisioner** 在调用 CSI 插件的 CreateVolume 接口之前，会在请求参数设置 AccessibilityRequirements：

- 对于 WaitForFirstConsumer

- 当 PVC 的 anno 中包含 "volume.kubernetes.io/selected-node" 且不为空，则先获取对应节点 CSINode 的 TopologyKeys，然后根据该 TopologyKeys 键从 Node 资源的 Label 获取 Values 值，最后拿该 Values 值与 StorageClass 的 AllowedTopologies 比对，判断其是否包含于其中；若不包含则报错。

- 对于 Immediately

- 将 StorageClass 的 AllowedTopologies 的值填进来，若 StorageClass 没有设置 AllowedTopologies 则将所有包含 TopologyKeys 键的节点 Value 添进来。

## Scheduler 如何处理使用存储卷调度

> 基于社区 1.18 版本调度器

调度器的调度过程主要有如下三步：

- **预选（Filter）**：筛选满足Pod调度要求的节点列表。

- **优选（Score）**：通过内部的优选算法为节点打分，获得最高分数的节点即为选中的节点。

- **绑定（Bind）**：调度器将调度结果通知给 kube-apiserver，更新 Pod 的 .spec.nodeName 字段。

调度器预选阶段：处理 Pod 的 PVC/PV 绑定关系以及动态供应 PV（Dynamic Provisioning），同时使调度器调度时考虑 Pod 所使用 PV 的节点亲和性。详细调度过程如下：

1. Pod 不包含 PVC 直接跳过。

2. FindPodVolumes

- 获取 Pod 的 boundClaims、claimsToBind 以及 unboundClaimsImmediate。

- boundClaims：已 Bound 的 PVC

- claimsToBind：PVC 对应 StorageClass 的 VolumeBindingMode 为 VolumeBindingWaitForFirstConsumer

- unboundClaimsImmediate：PVC 对应 StorageClass 的 VolumeBindingMode 为 VolumeBindingImmediate

- 若 len(unboundClaimsImmediate) 不为空，表示这种 PVC 需要立即绑定 PV（即存 PVC 创建后，立刻动态创建 PV 并将其绑定到 PVC，该过程不走调度），若 PVC 处于 unbound 阶段则报错。

- 若 len(boundClaims) 不为空，则检查 PVC 对应 PV 的节点亲和性与当前节点的 Label 是否冲突，若冲突则报错（可检查 Immediate 类型的 PV 拓扑）。

- 若 len(claimsToBind) 不为空

- 先检查环境中已有的 PV 能否与该 PVC 匹配（findMatchingVolumes），将能够匹配 PVC 的 PV 记录在调度器的 cache 中。

- 未匹配到 PV 的 PVC 走动态调度流程，动态调度主要通过 StorageClass 的 AllowedTopologies 字段判断当前调度节点是否满足拓扑要求（针对 WaitForFirstConsumer 类型的 PVC）。

调度器优选阶段不讨论。

调度器 Assume 阶段

> 调度器会先 Assume PV/PVC，再 Assume Pod。

1. 将当前待调度的 Pod 进行深拷贝。

2. AssumePodVolumes（针对 WaitForFirstConsumer 类型的 PVC）

- 更改调度器 cache 中已经 Match 的 PV 信息：设置 annotation：pv.kubernetes.io/bound-by-controller="yes"。

- 更改调度器 cache 中未匹配到 PV 的 PVC，设置 annotation：volume.kubernetes.io/selected-node=【所选节点】。

3. Assume Pod 完毕

- 更改调度器 cache 中 Pod 的 .Spec.NodeName 为【所选节点】。

调度器 Bind 阶段

BindPodVolumes：

- 调用 Kubernetes 的 API 更新集群中 PV/PVC 资源，使其与调度器 Cache 中的 PV/PVC 一致。

- 检查 PV/PVC 状态：

- 检查所有 PVC 是否已处于 Bound 状态。

- 检查所有 PV 的 NodeAffinity 是否与节点 Label 冲突。

- 调度器执行 Bind 操作：调用 Kubernetes 的 API 更新 Pod 的 .Spec.NodeName 字段。

### **2. 存储卷扩容**

存储卷扩容部分在 External Resizer 部分已提到，故不再赘述。用户只需要编辑 PVC 的 .Spec.Resources.Requests.Storage 字段即可，注意只可扩容不可缩容。

若 PV 扩容失败，此时 PVC 无法重新编辑 spec 字段的 storage 为原来的值（只可扩容不可缩容）。参考 K8s 官网提供的 PVC 还原方法：

_https://kubernetes.io/docs/concepts/storage/persistent-volumes/#recovering-from-failure-when-expanding-volumes_

## 3. 单节点卷数量限制

卷数量限制在 Node Driver Registrar 部分已提到，故不再赘述。

## 4. 存储卷监控

存储商需实现 CSI 插件的 NodeGetVolumeStats 接口，Kubelet 会调用该函数，并反映在其 metrics上：

- kubelet\_volume\_stats\_capacity\_bytes：存储卷容量

- kubelet\_volume\_stats\_used\_bytes：存储卷已使用容量

- kubelet\_volume\_stats\_available\_bytes：存储卷可使用容量

- kubelet\_volume\_stats_inodes：存储卷 inode 总量

- kubelet\_volume\_stats\_inodes\_used：存储卷 inode 使用量

- kubelet\_volume\_stats\_inodes\_free：存储卷 inode 剩余量

## 5. Secret

CSI 存储卷支持传入 Secret 来处理不同流程中所需要的私密数据，目前 StorageClass 支持如下 Parameter：

- csi.storage.k8s.io/provisioner-secret-name

- csi.storage.k8s.io/provisioner-secret-namespace

- csi.storage.k8s.io/controller-publish-secret-name

- csi.storage.k8s.io/controller-publish-secret-namespace

- csi.storage.k8s.io/node-stage-secret-name

- csi.storage.k8s.io/node-stage-secret-namespace

- csi.storage.k8s.io/node-publish-secret-name

- csi.storage.k8s.io/node-publish-secret-namespace

- csi.storage.k8s.io/controller-expand-secret-name

- csi.storage.k8s.io/controller-expand-secret-namespace

Secret 会包含在对应 CSI 接口的参数中，如对于 CreateVolume 接口而言则包含在 CreateVolumeRequest.Secrets 中。

## 6. 块设备

```properties
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: nginx-example
spec:
  selector:
    matchLabels:
      app: nginx
  serviceName: "nginx"
  volumeClaimTemplates:
  - metadata:
      name: html
    spec:
      accessModes:
        - ReadWriteOnce
      volumeMode: Block
      storageClassName: csi-pangu
      resources:
        requests:
          storage: 40Gi
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        volumeDevices:
        - devicePath: "/dev/vdb"
          name: html
```

三方存储厂商需实现 NodePublishVolume 接口。Kubernetes 提供了针对块设备的工具包（"k8s.io/kubernetes/pkg/util/mount"），在 NodePublishVolume 阶段可调用该工具的 EnsureBlock 和 MountBlock 函数。

## 7. 卷快照/卷克隆能力

鉴于本文篇幅，此处不做过多原理性介绍。读者感兴趣见官方介绍：卷快照、卷克隆。

# 总结

本文首先对 CSI 核心流程进行了大体介绍，并结合 CSI Sidecar 组件、CSI 接口、API 对象对 CSI 标准进行了深度解析。在 K8s 上，使用任何一种 CSI 存储卷都离不开上面的流程，环境上的容器存储问题也一定是其中某个环节出现了问题。
