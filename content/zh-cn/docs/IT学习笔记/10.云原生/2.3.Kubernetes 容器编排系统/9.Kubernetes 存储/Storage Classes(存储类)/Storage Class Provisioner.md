---
title: Storage Class Provisioner
---

# NFS Provisioner

> 参考：
> - [GitHub](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner)
> - [GitHubOld](https://github.com/kubernetes-retired/external-storage/tree/master/nfs-client)

注意：NFS Provisioner 不支持容量限制功能
NFS subdir external provisioner 是一个自动 Provisioner，它使用您现有的和已配置的 NFS 服务器来通过持久卷声明来动态供应 Kubernetes 持久卷。PV 配置为 $ {namespace}-$ {pvcName}-$ {pvName}。

## 部署 NFS Provisioner

> 參考：
>
> - [arifacthub](https://artifacthub.io/packages/helm/ckotzbauer/nfs-client-provisioner)
> - [GitHub](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner#how-to-deploy-nfs-subdir-external-provisioner-to-your-cluster)
> - [GitHubOld](https://github.com/helm/charts/tree/master/stable/nfs-client-provisioner#installing-the-chart)

创建 namespace

    apiVersion: v1
    kind: Namespace
    metadata:
      name: storage

创建 rbac

```yaml
kind: ServiceAccount
apiVersion: v1
metadata:
  name: nfs-client-provisioner
  namespace: storage
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-client-provisioner-runner
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["create", "update", "patch"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-client-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-client-provisioner
    namespace: storage
roleRef:
  kind: ClusterRole
  name: nfs-client-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-client-provisioner
  namespace :storage
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-client-provisioner
  namespace :storage
subjects:
  - kind: ServiceAccount
    name: nfs-client-provisioner
    # replace with namespace where provisioner is deployed
    namespace: storage
roleRef:
  kind: Role
  name: leader-locking-nfs-client-provisioner
  apiGroup: rbac.authorization.k8s.io
```

部署 client

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-client-provisioner
  namespace: storage
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: nfs-client-provisioner
  namespace: storage
spec:
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      serviceAccountName: nfs-client-provisioner
      containers:
        - name: nfs-client-provisioner
          imagePullPolicy: Never #不让每次部署都拉取镜像，直接使用本地镜像即可
          image: quay.io/external_storage/nfs-client-provisioner:latest
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: nfs-storage
            - name: NFS_SERVER
              value: 填写nfs设备的IP或者主机名
            - name: NFS_PATH
              value: 填写nfs设备上设定的数据存储路径
      volumes:
        - name: nfs-client-root
          nfs:
            server: 填写nfs设备的IP或者主机名
            path: 填写nfs设备上设定的数据存储路径
```

创建 storageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: managed-nfs-storage
provisioner: nfs-storage # 必须与 manifests 文件中 PROVISIONER_NAME 环境变量的值相同
parameters:
  archiveOnDelete: "false" # 设置为 false 时，PV 将不会在删除 PVC 时存档
```

# OpenEBS

## 概述

> 参考：
> - [官网](https://openebs.io/)

OpenEBS 建立在 Kubernetes 之上，使有状态应用程序能够轻松访问动态 **Local PV** 或 **Replicated PV**。通过使用容器附加存储模式，用户报告说他们的团队成本更低，管理更容易，控制力更强。
OpenEBS 是![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/cbnzrv/1627391323750-5ed0ed9d-7753-4790-bbb5-e0bfdcc25523.png) 由 MayaData 和社区共同开发的 100% 开源 CNCF 项目。著名用户包括 Arista、Optoro、Orange、Comcast 和 CNCF 本身。
