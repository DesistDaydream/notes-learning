---
title: Service Account 详解
---

# 概述

> 参考：
> - [官方文档,任务-配置 Pod 和 容器-为 Pods 配置服务账户](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/)
> - [官方文档,参考-API 访问控制-管理服务账户](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/)

详解 Service Account 之前，需要了解这么一个 k8s 的运行逻辑：**每个 Pod 在创建成功后，都会声明并使用一个 ServiceAccount 作为自己与集群通信的认证，默认使用 Pod 所在 Namepace 的名为 default 的 ServiceAccount**

上面这个默认的 default 具有很高的权限，如果想对该 default 进行安全保护，可以修改绑定到 default 的 role 的权限

每个 ServiceAccount 对象在被创建出来之后，都会自动生成一个对应的 Secrets 对象，认证信息即在该 Secret 中。

```bash
# 与 ServiceAccount 关联的 secret 会以 SAName-token-STRING 的方式命名。
# 效果如下，在monitoring 名称空间中每个 sa 都有一个 secret 对应。(SA 是 ServiceAccount 的简称)
[root@master alertmanager]# kubectl get serviceaccount -n monitoring
NAME                  SECRETS   AGE
alertmanager-main     1         2m18s
default               1         13d
....
[root@master alertmanager]# kubectl get secrets -n monitoring
NAME                              TYPE                                  DATA   AGE
alertmanager-main-token-h4vfx     kubernetes.io/service-account-token   3      2m22s
default-token-8hww4               kubernetes.io/service-account-token   3      13d
.......
```

该 Secret 中包含 3 个数据

1. ca.crt # 集群的 CA 的证书

2. namespace # 该 Pod 属于哪个 namespace

3. token # 该 Pod 所用的 Service Account 的令牌信息

这几个数据文件在 Pod 启动后会被自动挂载到 容器内的 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下(该目录下的证书和 token 和 namespace 就类似于 kubeclt 命令的配置文件里的证书和 token 和 namespace)。效果如下：

```bash
[root@master-1 ~]# kubectl exec -it myapp -- /bin/sh
/ # ls -l /var/run/secrets/kubernetes.io/serviceaccount/
total 0
lrwxrwxrwx    1 root     root            13 Sep 16 13:24 ca.crt -> ..data/ca.crt
lrwxrwxrwx    1 root     root            16 Sep 16 13:24 namespace -> ..data/namespace
lrwxrwxrwx    1 root     root            12 Sep 16 13:24 token -> ..data/token
```

Service Account 相关流程

从上面的介绍中可以看到，Pod 启动后，会自动加载 token ，SA 生成后会自动创建 token ，那么这些工作，都是由谁来完成的呢？~

在 Kubernetes 中有三个组件来协作完成 SA 的自动化工作

1. 服务账户准入控制器（Service account admission controller）

2. Token 控制器（Token controller）

3. 服务账户控制器（Service account controller）

### Service account admission controller 服务账户准入控制器

对 pod 的改动通过一个被称为 [Admission Controller](https://kubernetes.io/docs/admin/admission-controllers) 的插件来实现。它是 apiserver 的一部分。 当 pod 被创建或更新时，它会同步地修改 pod。 当该插件处于激活状态 ( 在大多数 k8s 发行版中都是默认的 )，当 pod 被创建或更新时它会进行以下动作：

1. 如果该 pod 没有 ServiceAccount 设置，将其 ServiceAccount 设为 default。

2. 保证 pod 所关联的 ServiceAccount 存在，否则拒绝该 pod。

3. 如果 pod 不包含 ImagePullSecrets 设置，那么 将 ServiceAccount 中的 ImagePullSecrets 信息添加到 pod 中。

4. 为 pod 添加一个用于访问 API 的卷(卷中包含访问 API 时所需的认证 token)

5. 将上面的卷挂载到 pod 中每个容器的 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下

反映在 pod 的 yaml 中，就是下面这个样子

        volumeMounts:
        - mountPath: /var/run/secrets/kubernetes.io/serviceaccount # 将卷挂载到 pod 中每个容器的 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下
          name: default-token-4w977
          readOnly: true
      volumes: # 为 pod 添加一个用于访问 API 的卷(卷中包含访问 API 时所需的认证 token)
      - name: default-token-4w977
        secret:
          defaultMode: 420
          secretName: default-token-4w977

### Token controller 令牌控制器&#xA;

Token 控制器 是 controller-manager 的一部分。 以异步的形式工作：

- 检测服务账户的创建，并且创建相应的 Secret 以支持 API 访问。

- 检测服务账户的删除，并且删除所有相应的服务账户 Token Secret。

- 检测 Secret 的增加，保证相应的服务账户存在，如有需要，为 Secret 增加 token。

- 检测 Secret 的删除，如有需要，从相应的服务账户中移除引用。

你需要通过 --service-account-private-key-file 参数项传入一个服务账户私钥文件至 Token 管理器。 私钥用于为生成的服务账户 token 签名。 同样地，你需要通过 --service-account-key-file 参数将对应的公钥传入 kube-apiserver。 公钥用于认证过程中的 token 校验。

创建额外的 API tokens

控制器中有专门的循环来保证每个服务账户中都存在 API token 对应的 Secret。 当需要为服务账户创建额外的 API token 时，创建一个类型为 ServiceAccountToken 的 Secret，并在 annotation 中引用服务账户，控制器会生成 token 并更新 :

### Service account controller 服务账户控制器&#xA;

服务账户管理器管理各命名空间下的服务账户，并且保证每个活跃的命名空间下存在一个名为 "default" 的服务账户

总结：

Secret 的作用与 kubeconfig 的作用有异曲同工之妙，ServiceAccount 生成后自动创建的 Secret 里的信息是在 pod 与 APIServer 交互时使用(可以把 pod 当做 kubectl )，比如当一个 pod 中的进程想要 get 或者 delete 资源时，相当于是对集群请求执行该指令，而集群是通过 ApiServer 来接收这些指令的。那么首先就要确认使用 pod 进行这操作的 User 是谁，这个 User 是否有证书来对 ApiServer 发起这些请求。如果 ApiServer 都不认可这个 User ,那么都不会接受这些指令请求。

ServiceAccount 的使用

官方文档：<https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/>

1. 常见用法：用来在 clusterrolebinding(或者 rolebinding) 中，将指定的 ServiceAccount 与指定的 ClusterRole(或者 role) 进行绑定，绑定之后，所有使用该 ServiceAccount 的 pod 就有了该 ClusterRole(或者 role) 中定义的操作授权

   1. 这是为了在部署服务(i.e.pod)时，可以不让 pod 获取太高的授权，一般都会自定义一个 sa，然后创建一个 role 仅有几种操作授权，再将 sa 与这个 role 绑定之后，就可以限制这个服务对集群可操作权限的多少

   2. 一般像 dashboard、prometheus、ingress 等服务都需要这么做，因为这些服务所创建的 pod 一般都会用来对 k8s 集群内的资源进行一系列操作来完成本身的功能(比如获取集群信息(至少需要 get pod)等等)。

2. 比如 Prometheus，当在集群中部署 Prometheus 时，一般是在 monitoring 名称空间中，如果 Prometheus 想要获取其他名称空间中 pod 的监控指标，则需要在对应的名称空间中创建 role 和 rolebing，并将 Prometheus 所用的 ServiceAccount 与该名称空间中的 role 绑定，这样，Prometheus 才可以获取指定名称空间中 pod 的相关指标。

3. 比如 dashboard 这个用于管理 kubernetes 集群的 web 界面的 Pod，该 Pod 可以实现对集群的增删改查的操作(通过 Pod 中 container 的进程实现)，与 kubectl 一样，这个 Pod 也会调用 API 对集群进行操作。但是在此之前需要对该 Pod 进行认证(只不过与 kubectl 不同，这个认证不是通过 UserAccount(i.e.kubeconfig)来进行的，而是通过 ServiceAccount 来进行)创建完 ServiceAccount 资源后，会自动生成一个 secrets 的资源，secrets 就是该 ServiceAccount 的认证，使用这个 Secrets 里的 token 来让 ServiceAccount 与 API 进行登录认证，然后再给 ServiceAccount 进行授权操作以决定 API 会接收哪些指令

   1. 创建 ServiceAccount 并获取此 ServiceAccount 的 secret，使用 kubectl describe 查看该 secret 的详细信息，使用其中的 token 来进行登录 web 界面

   2. 根据其管理目标，使用 rolebinding 或 clusterrolebinding 将 ServiceAccount 绑定至合理的 role 或 clusterrole 以实现权限管理

   3. secret 的密文都是通过把内容进行 base64 加密后得出的结果，base64 -w 0 为加密命令 base64 -d 为解密命令，可以使用`echo “要加/解密的内容” | base64  “加/解密选项”`命令
