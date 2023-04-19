---
title: Pod Manifest 详解
---

# 概述

> 参考：
>
> - [API 文档，单页](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.21/#pod-v1-core)
> - [官方文档，参考-Kubernetes API-工作负载资源-Pod](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/)

Pod 是可以在主机上运行的容器的集合。此资源由客户端创建并调度到主机上。

## Manifest 中的顶层字段

- **apiVersion**: v1
- **kind**: Pod
- **metadata**([metadata](#metadata))
- **spec**([spec](#spec))
- **status**([status](#status))

# metadata

**metadata** 字段用来描述一个 Pod 的元数据信息。该字段内容详见通用定义的 [ObjectMeta](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/ObjectMeta.md)

- **annotations**(STRING) # Pod 注释，不同于 label，仅用于为资源提供元数据
- **labels**(map[STRING]STRING) # Pod 的标签，通过“键值对”的方式定义，可以添加多个标签
  - KEY: VAL # 比如该键值可以是 run: httpd，标签名是 run，run 的值是 httpd，key 与 val 使用字母，数字，\_，-，.这几个字符且以字母或数字开头；val 可以为空。
- **name**(STRING) # Pod 对象名称。必须名称空间唯一。
- **ownerReferences**([]Object) # 该对象所依赖的对象列表，一般由控制器自动生成。也可以手动指定。

# spec

**spec** 字段用来描述一个 Pod 应该具有的属性。Pod 中的 spec 字段大体分为如下几类

- Containers # 用来描述 Pod 中容器的属性
- Volumes # 用来描述 Pod 所用卷，以及容器如何使用这些卷
- Scheduling # Pod 如何被调度到 node
- Lifecycle # Pod 的生命周期
- Hostname and Name resolution # 容器的主机名和域名解析
- Hosts namespaces # Pod 使用主机上名称空间的方法
- Service account # Pod 的服务账户
- Security context # Pod 安全相关

## Containers(容器) 相关字段

**containers**(\[][containers](#containers)) # 属于该 Pod 的 Containers 列表。[containers](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#container-v1-core) 字段**只会在 Pod 环境**中创建

**initContainers**(\[][containers](#containers)) # 属于该 Pod 的初始化容器的列表。该字段中所定义的容器都会比普通 containers 字段定义的容器先启动，并且 initContainer 会按顺序逐一启动，直到它们都启动并且退出之后，普通容器才会启动

**imagePullSecrets: <\[]Object>** # 拉取镜像时使用的私有仓库的信息。

拉取镜像时，如果是私有仓库，则使用该字段指定的 secret 中的信息。实际上就是代替 docker login 命令。 更多信息见官网：[Specifying imagePullSecrets on a Pod()章节](https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod)

## Volumes(卷) 相关字段

**volumes: \<[]Object>** # 给 pod 创建一个 Volume

- **TYPE:** # 选择要创建的 volume 的类型，具体支持的类型可以使用 kubectl explain pod.spec.volumes 命令查看
  - ...... # 定义该类型的 volume 相关参数
- **name: \<STRING>** # **必须的**。自定义该存储卷的名称

## Scheduling(调度) 相关字段

**nodeSelector: \<map\[string]string>** # 指明 Node 标签选择器，该 Pod 会运行在具有相同标签的 Node 上

**nodeName: \<STRING>** # Pod 运行在指定 Node 上

**affinity: \<Object>** # Pod 亲和性调度规则。用法详见[调度器章节](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/5.Scheduling(调度)/让%20Pod%20运行在指定%20Node%20上.md)

- **nodeAffinity**([nodeAffinity](#nodeAffinity)) # 为 Pod 定义节点亲和性的调度规则
- **podAffinity**([podAffinity](#podAffinity)) # 为 Pod 定义 Pod 亲和性的调度规则(e.g.将此 Pod 与其他一些 Pod 共同定位在同一节点、区域等中)。
- **podAntiAffinity**([PodAntiAffinity](podAntiAffinity)) # 描述 pod 反亲和性的调度规则(e.g. 避免将此 Pod 放在与其他某些 Pod 相同的节点、区域等中)

**tolerations: <\[]Object>** # 定义 Pod 容忍污点的容忍度。用法详见[调度器章节](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/5.Scheduling(调度)/让%20Pod%20运行在指定%20Node%20上.md)

## Lifecycle(生命周期) 相关字段

**restartPolicy: \<STRING>** # Pod 中容器失败后的重启策略，`默认值：Always`

- STRING 可用的值有：Always、OnFailure、Never

## Hostname and Name resolution(主机名和域名解析) 相关字段

**dnsConfig: \<OBJECT>** #

**dnsPolicy: \<STRING>** # pod 的域名解析策略。`默认值：ClusterFirst`

- 可用值有：
  - ClusterFirstWithHostNet
  - ClusterFirst
  - None

## Hosts namespaces(容器如何使用宿主机中的名称空间) 相关配置

**hostNetwork: \<BOOLEAN>** # 是否让 Pod 中的容器使用主机的网络名称空间。`默认值：false`

## Service account(服务账户) 相关字段

**serviceAccountName: \<STRING>** # 容器所使用 ServiceAccount。

## Security context(容器安全环境) 相关字段

**securityContext: \<Object> 和 .spec.container.securityContext: \<Object>**

`securityContext` 字段用于配置如何安全得运行 Pod，比如以 非特权用户运行容器、SELinux 等等。Pod 安全配置内容，在 [Security Context(安全环境) 文章](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/3.Pod%20集群最小的工作单元/Pod%20管理/Security%20Context(安全环境).md)中有更详细的描述。

注意：`.spec.securityContext` 字段属于 Pod 级别的安全配置，在 `.spec.containers` 字段下，还有一个 securityContext 字段，`.spec.containers.securityContext` 字段属于 Container 级别的安全配置，该字段的优先级要高于 `.spec.securityContext`。

- **fsGroup: \<INTEGER>**# 功能待确认。
- **runAsNonRoot: \<BOOLEAN>** # 容器内的进程是否不以 root 身份运行。`默认值：false`。
  - 若为 true，则必须指定 runAsUser 字段，除非构建镜像时已经指定了进程运行的 UID。
- **runAsUser: \<INTEGER>** # 指定运行容器内进程的 UID 为 INTEGER

## 其他类别的字段

# status

status 字段表示 Pod 的状态信息。状态可能会落后于系统的实际状态，尤其是当承载 Pod 的节点无法联系控制平面时。

**phase: STRING** # phase(阶段) 字段是对 Pod 在其生命周期中所处位置的简单、高级的总结。条件数组、原因和消息字段以及各个容器状态数组包含有关 pod 状态的更多详细信息。有五个可能的相位值：

- `Pending` # Pod 已被 Kubernetes 系统接受，但尚未创建容器镜像。 这包括 Pod 被调度之前的时间以及通过网络下载镜像所花费的时间。
- `Running` # Pod 已经被绑定到某个节点，并且所有的容器都已经创建完毕。至少有一个容器仍在运行，或者正在启动或重新启动过程中。
- `Succeeded` # Pod 中的所有容器都已成功终止，不会重新启动。
- `Failed` # Pod 中的所有容器都已终止，并且至少有一个容器因故障而终止。 容器要么以非零状态退出，要么被系统终止。
- `Unknown` # 由于某种原因无法获取 Pod 的状态，通常是由于与 Pod 的主机通信时出错。
- 更多信息： [https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle#pod-phase](https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle#pod-phase)

# 通用字段

## containers

**args: <\[]STRING>** # 定义容器运行的命令和参数。用于替换容器镜像中 CMD 指令。详见[为容器定义命令和参数章节](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/)

```yaml
# 注意，在使用 args 为容器传递 flags 时，不要使用空格。否则，会自动为 flags 和 参数 加上单引号，导致 flags 失效
# 比如现在有如下 args 配置
        args:
        - deletecr
        - --ns default
# 当容器运行后，会被转换成 "deletecr '--ns default'"，这时，deletecr 的 flag 变成了 --ns default，而不是 --ns。
# 这时，就会报错，并提示如下内容
flag provided but not defined: '--ns default'
# 可以看到，在容器中，将 --ns default 这个整体当作了一个 flags。
# 所以，如果想要使用 args，可以这样写：
        args:
        - deletecr
        - --ns=default
# 或者
        args:
        - deletecr
        - --ns
        - default
```

**command: <\[]STRING>** # 定义容器运行的命令和参数。用于替换容器镜像中的 ENTRYPOINT 指令。详见为[容器定义命令和参数章节](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/)

**env: <\[]Object>** # 要在容器中设置的环境变量列表。详见[为容器定义命令和参数章节](https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/)

**name: \<STRING>**# 容器的名称

**ports: <\[]Object>**# 容器对外暴露的端口，主要作为参考信息，就算不指定，通过 Service 也可以关联到容器中的指定端口，并通过 Service 访问容器内部。

- **containerPort: \<INTEGER>** # 容器内端口号
- **name: \<STRING>** # 端口号的名称，必须在 pod 中唯一。service 可以通过 name 引用端口。

**resources: \<OJBECT>**# 容器所需的资源。即.所需的 CPU、Memory 等等

- **limits: \<map\[STRING]STRING>** # 容器可以使用的最大资源
  - **cpu: \<STRING>** # 定义容器的 CPU 限额
  - **memory: \<STRING>** # 定义容器的 Memory 限额
  - ...... 其他资源限额
- **requests: \<map\[STRING]STRING>** # 容器所需的最小资源。如果 Requests 省略，则默认与 limits 下定义的值保持一直。
  - **cpu: \<STRING>** # 定义容器的 CPU 需求
  - **memory: \<STRING>** # 定义容器的 Memory 需求
  - ...... 其他资源需求

**volumeMounts: <\[]Object>** # 给 Container 挂载在 Pod 中创建的 Volume。Volume 通过下文的 [Volumes 字段](#Volumes(卷)%20相关字段)指定

- **mountPath: \<STRING>** # **必须的**。把 Volume 挂载到容器中的目录上
- **name: \<STRING>** # **必须的**。要挂载的 Volume 的名称。必须与 `spec.volumes.TYPE.volumeName` 字段的值相同，才可以引用到卷。

### Image(镜像) 相关字段

**image: \<STRING>** # 容器使用的镜像

**imagePullPolicy: \<STRING>** # 指明镜像拉取策略，公有三种 Always、IfNotPresent、Never。`默认值：IfNotPresent`

### Lifecycle(生命周期) 相关字段

Pod 中容器的生命周期功能详见[《Pod 的生命周期》](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/3.Pod%20集群最小的工作单元/Pod%20生命周期/Pod%20生命周期与探针.md)章节

**livenessProbe**([XXXProbe](#XXXProbe)) # 存活性探针，定期检测容器是否存活，容器**正常启动后**开始检查。若探针失败则容器将会重启

- 注意：如果 liveness 探测失败，kubelet 会杀死容器，并且容器会受到其重启策略的约束。如果不提供活性探测，则默认状态为成功。

**readinessProbe**([XXXProbe](#XXXProbe)) # 就绪状态探针，定期检测容器服务的准备状态，容器**正常启动前**开始检查。若探针失败则容器不会变为 Running 状态。

- 注意：如果就绪探测失败，端点控制器会从与 Pod 匹配的所有服务的端点中删除 Pod 的 IP 地址。初始延迟之前的默认就绪状态为失败。如果不提供就绪探测，则默认状态为成功。

### Debugging(调试) 相关字段

**stdin: \<BOOLEAN>** # 

**stdinOnce: \<BOOLEAN>** # 

**tty: \<BOOLEAN>** # 

## nodeAffinity

https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#NodeAffinity

**preferredDuringSchedulingIgnoredDuringExecution: <\[]Object>** # 调度程序将倾向于将 Pod 调度到满足此字段指定的反亲和行要求的节点，但是也可能会选择违反一个或多个该字段指定的调度规则。

- **preference: \<OBJECT> # 必须的**。
  - **matchExpressions: <\[]OBJECT>** # 该字段下的内容就是 [通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
  - **matchFields: <\[OBJECT]>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **weight: \<INTEGER> # 必须的**。

**requiredDuringSchedulingIgnoredDuringExecution: \<Object>** # 如果在调度时未满足该字段指定的反亲和性要求，则不会将 pod 调度到该节点上。

- **nodeSelectorTerms: <\[]OBJECT> # 必须的**。节点选择器列表。列表中元素之间是“或”的关系
  - **matchExpressions: <\[]OBJECT>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
  - **matchFields: <\[OBJECT]>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)

## podAffinity

https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#PodAffinity

## podAntiAffinity

https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#PodAntiAffinity

**preferredDuringSchedulingIgnoredDuringExecution: <\[]Object>** # 调度程序将倾向于将 Pod 调度到满足此字段指定的反亲和行要求的节点，但是也可能会选择违反一个或多个该字段指定的调度规则。

- **preference: \<OBJECT> # 必须的**。
  - **matchExpressions: <\[]OBJECT>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
  - **matchFields: <\[OBJECT]>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **weight: \<INTEGER> # 必须的**。

**requiredDuringSchedulingIgnoredDuringExecution: <\[]Object>** # 如果在调度时未满足该字段指定的反亲和性要求，则不会将 pod 调度到该节点上。

- **labelSelector: \<OBJECT>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **namespaceSelector: \<OBJECT>** # 该字段下的内容就是[通用的标签选择器字段](/docs/10.云原生/2.3.Kubernetes%20容器编排系统/1.API%20Resource%20与%20Object/API%20参考/Common%20Definitions(通用定义)/LabelSelector.md)
- **namespaces: <\[]STING>** # 名称空间。`默认值：该 Pod 所在的名称空间`
- **topologyKey: \<STRING>** # **必须的**。想要匹配的 Node 标签的键

## XXXProbe

XXXProbe 是 Probe(探针) 相关字段，比如 livenessProbe、readinessProbe 等字段的值都可以使用这部分内容。

**exec: \<Object>** # 通过在容器中执行命令作为探针检测方法

  - **exec.command: <\[]STRING>**

**httpGET: \<Object>** # 使用 HTTP 的 GET 的请求作为探针检测方法。

**tcpSocket: \<Object>** # 通过检测 TCP 的端口作为探针检测方法。

**grpc: \<Object>** #




# Pod Manifest 样例

以下是最简单最基本的 Pod 模板，具体 Pod 中可以实现的功能以及这些功能应该使用哪些 yaml 里的字段详见以下几处

- [Configure Pods and Containers(配置 Pod 和 Container)](https://kubernetes.io/docs/tasks/configure-pod-container/)章节下面的所有内容；每种可以在 Pod 中配置的功能，都是一小章节
- [Inject Data Into Applications(将数据注入应用程序)](https://kubernetes.io/docs/tasks/inject-data-application/) 章节

## 简单示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
  labels:
    name: myapp
spec:
  containers:
    - name: myapp
      image: lchdzh/network-test
```
