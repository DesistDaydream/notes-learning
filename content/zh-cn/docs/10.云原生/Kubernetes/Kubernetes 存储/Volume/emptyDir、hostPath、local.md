---
title: emptyDir、hostPath、local
---

# emptyDir # 在 Pod 的生命周期中的空目录

该类型的 Volume 会随着 Pod 的摧毁重建而自动还原成初始状态。Pod 会创建一个逻辑上的 volume，把该 volum 挂载到一个 Pod 中每个 container 所定义的目录，所有 container 对自己挂载 volume 的目录进行的操作都会在其余 container 挂载该 volume 的目录中看到(每个 container 用于挂载的目录可以不一样，但是用到的 volume 都是同一个)。

- medium：指定 volum 的存储媒介(即 volum 使用的存储资源)，默认使用 memory，这样两个容器的数据交互速度会非常快

- sizeLimit：容量大小限制，限制 volume 的最大存储空间，如果不做限制，那么对于设备来说，用户数据交互的 volume 会非常浪费资源

- 获取 emptyDir 类型的 volume 在宿主机的路径的方式：

  - 首先通过 kubectl get pod PODNAME -o yaml | grep uid 来获取 pod 的标识符

  - 然后在目录/var/lib/kubelet/pods/PodID/volumes/kubernets.io~empty-dir/下面找到所有该 pod 所挂载的 empty 类型的 volume，该目录与 pod 中的目录是同步的，在该目录增删改的信息同样也会影响到 pod 对应的目录中。

hostPath # Node 上的文件或目录

hostPath 类型的 Volume 会将节点上的文件或目录挂载到 Pod 中。

hostPath 常用于将节点上的目录挂载到容器中，以便让容器可以读取目录中的内容执行一些操作，比如：

- 运行一个需要访问节点日志的容器(比如 Loki)，使用 hostPath 类型的卷挂载节点上的 /var/log/pods 目录到容器中。以便 Loki 可以采集宿主机的日志。

- 在容器中运行 cAdvisor 时，以 hostPaht 的方式挂载节点的 /sys 目录。

- 使用 hostPath 作为 volume 时，有两个参数需要指定：

  - path：指明 Node 上的哪个文件或者目录作为 volume 提供给 container

  - type：Node 上目录或文件的使用行为，比如

    - DirectoryOrCreate # 如果指定 path 不存在，则自动创建，并设置权限为 0755，具有有 kubelet 相同的组和所有权

    - Directory # 指定 path 必须存在，不存在则报报错

    - ....等等，详见官网

注意：

- hostPath 类型的 volume 关联的目录，并不会自动变更权限，自动创建出来后，如果 容器 不以 root 身份运行，则有可能会没有权限对该目录进行读写操作

- hostPath 类型的 volume 不会随着 Pod 的摧毁重建而自动还原成初始状态。这种类型会把 Node 上的本地存储资源当作 volume 给 container 使用，保证了数据的持久性，但是当 pod 重建时，如果不被分配到原 Node，则该数据会变成另一个 Node 上的数据，没法使用原数据了，不够灵活。

# local # 指定 Node 上的存储设备(如 磁盘、分区、目录)

> 注意：local 类型的卷暂时只能在静态的 Persistent Volume 持久卷 中使用。

1. local 会根据容器所需要的权限，自动更改关联目录的权限

2. 相比于 hostPath 卷，local 卷更像是存储数据所用，而 hostPath 更偏向于容器共享宿主机目录。只不过 local 卷只能在 PV 中使用，并必须指定节点亲和性，以便将 PV 绑定到指定的节点上。这样，pod 在使用对应 PVC 时，则会自动调度到 PV 所在的 节点上，而不用给 pod 手动指定 nodeselector 了。

# Manifests 示例

## hostPath 类型

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pd
spec:
  containers:
    - image: k8s.gcr.io/test-webserver
      name: test-container
      volumeMounts:
        - mountPath: /test-pd # 指定要改volume要挂载到容器的哪个目录下
          name: test-volume # 指定要使用的volume，该字段的值为.volumes.name中的值
  volumes:
    - name: test-volume # 指定volume的名字，用于pod挂载volume时所用
      hostPath: # 指定volume的类型，该字段下面的字段为该类型的参数，不同的类型有不通的参数
        path: /data
        type: Directory
```
