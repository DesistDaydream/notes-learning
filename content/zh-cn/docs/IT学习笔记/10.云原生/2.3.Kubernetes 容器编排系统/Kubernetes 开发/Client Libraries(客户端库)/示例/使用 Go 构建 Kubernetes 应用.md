---
title: 使用 Go 构建 Kubernetes 应用
---

[使用 Go 构建 Kubernetes 应用](https://mp.weixin.qq.com/s/849NTst8GGFMiZZQQgrgUg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)

Kubernetes 项目使用 Go 语言编写，对 Go api 原生支持非常便捷。本篇文章介绍了如何使用 kubernetes client-go 实践一个简单的与 K8s 交互过程。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)

## kubernetes 的 Go Client 项目（client-go）

go client 是 k8s client 中最古老的一个，具有很多特性。Client-go 没有使用 Swagger 生成器，它使用的是源于 k8s 项目中的源代码生成工具，这个工具的目的是要生成 k8s 风格的对象和序列化程序。

该项目是一组包的集合，该包能够满足从 REST 风格的原语到复杂 client 的不同的编程需求。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)

RESTClient 是一个基础包，它使用`api-machinery`库中的类型作为一组 REST 原语提供对 API 的访问。作为对`RESTClient`之上的抽象，\_clientset\_将是你创建 k8s client 工具的起点。它暴露了公开化的 API 资源及其对应的序列化。

注意：在 client-go 中还包含了如 discovery, dynamic, 和 scale 这样的包，虽然本次不介绍这些包，但是了解它们的能力还是很重要的。

## 一个简单的 k8s client 工具

让我们再次回顾我们将要构建的工具，来说明 go client 的用法。**pvcwatch**是一个简单的命令行工具，它可以监听集群中声明的 PVC 容量。当总数到达一个阈值的时候，他会采取一个 action（在这个例子中是在屏幕上通知显示）

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)

**你能在 github 上找到完整的例子**

这个例子是为了展示 k8s 的 go client 的以下几个方面：- 如何去连接 - 资源列表的检索和遍历 - 对象监听

## Setup

client-go 支持 Godep 和 dep 作为 vendor 的管理程序，我觉得 dep 便于使用所以继续使用 dep。例如，以下是 client-go v6.0 和 k8s API v1.9 所需最低限度的`Gopkg.toml`。

`[[constraint]] name = "k8s.io/api" version = "kubernetes-1.9.0" [[constraint]] name = "k8s.io/apimachinery" version = "kubernetes-1.9.0" [[constraint]] name = "k8s.io/client-go" version = "6.0.0"`

运行`dep ensure`确保剩下的工作。

## 连接 API Server

我们 Go client 的第一步就是建立一个与 API Server 的连接。为了做到这一点，我们要使用实体包中的`clientcmd`，如下代码所示：

`import ( ...     "k8s.io/client-go/tools/clientcmd" ) func main() {     kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config",     )     config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)     if err != nil {         log.Fatal(err)     } ... }`

\_Client-go\_通过提供实体功能来从不同的上下文中获取你的配置，从而使之成为一个不重要的任务。

### 从 config 文件

正如上面的例子所做的那样，你能从 kubeconfig 文件启动配置来连接 API server。当你的代码运行在集群之外的时候这是一个理想的方案。`clientcmd.BuildConfigFromFlags("", configFile)`

### 从集群

当你的代码运行在这个集群中的时候，你可以用上面的函数并且不使用任何参数，这个函数就会通过集群的信息去连接 api server。

`clientcmd.BuildConfigFromFlags("","")`

或者我们可以通过 rest 包来创建一个使用集群中的信息去配置启动的 (译者注：k8s 里所有的 Pod 都会以 Volume 的方式自动挂载 k8s 里面默认的 ServiceAccount, 所以会用默认的 ServiceAccount 的授权信息)，如下：

`import "k8s.io/client-go/rest" ... rest.InClusterConfig()`

### 创建一个 clientset

我们需要创建一个序列化的 client 为了让我们获取 API 对象。在`kubernetes`包中的 Clientset 类型定义，提供了去访问公开的 API 对象的序列化 client，如下：

`type Clientset struct {     *authenticationv1beta1.AuthenticationV1beta1Client     *authorizationv1.AuthorizationV1Client ...     *corev1.CoreV1Client }`

一旦我们有正确的配置连接，我们就能使用这个配置去初始化一个 clientset，如下：

`func main() {config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)     ...     clientset, err := kubernetes.NewForConfig(config)     if err != nil {         log.Fatal(err)     } }`

对于我们的例子，我们使用的是`v1`的 API 对象。下一步，我们要使用 clientset 通过`CoreV1()`去访问核心 api 资源，如下：

`func main() {     ...     clientset, err := kubernetes.NewForConfig(config)     if err != nil {         log.Fatal(err)     }     api := clientset.CoreV1() }`

你能在这里看到可以获得 clientsets。

## 获取集群的 PVC 列表

我们对 clientset 执行的最基本操作之一获取存储的 API 对象的列表。在我们的例子中，我们将要拿到一个 namespace 下面的 pvc 列表，如下：

`import ( ...     metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" ) func main() {     var ns, label, field string     flag.StringVar(&ns, "namespace", "","namespace")     flag.StringVar(&label, "l", "","Label selector")     flag.StringVar(&field, "f", "","Field selector") ...     api := clientset.CoreV1()     // setup list options     listOptions := metav1.ListOptions{         LabelSelector: label,          FieldSelector: field,     }     pvcs, err := api.PersistentVolumeClaims(ns).List(listOptions)     if err != nil {         log.Fatal(err)     }     printPVCs(pvcs) ... }`

在上面的代码中，我们使用`ListOptions`指定 label 和 field selectors （还有 namespace）来缩小 pvc 列表的范围，这个结果的返回类型是`v1.PeristentVolumeClaimList`。下面的这个代码展示了我们如何去遍历和打印从 api server 中获取的 pvc 列表。

`func printPVCs(pvcs *v1.PersistentVolumeClaimList) {     template := "%-32s%-8s%-8s\n"     fmt.Printf(template, "NAME", "STATUS", "CAPACITY")     for _, pvc := range pvcs.Items {         quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]         fmt.Printf(             template,              pvc.Name,              string(pvc.Status.Phase),              quant.String())     } }`

## 监听集群中 pvc

k8s 的 Go client 框架支持为指定的 API 对象在其生命周期事件中监听集群的能力，包括创建，更新，删除一个指定对象时候触发的`CREATED`,`MODIFIED`,`DELETED`事件。对于我们的命令行工具，我们将要监听在集群中已经声明的 PVC 的总量。

对于某一个 namespace，当 pvc 的容量到达了某一个阈值（比如说 200Gi），我们将会采取某个动作。为了简单起见，我们将要在屏幕上打印个通知。但是在更复杂的实现中，可以使用相同的办法触发一个自动操作。

### 启动监听功能

现在让我们为`PersistentVolumeClaim`这个资源通过`Watch`去创建一个监听器。然后这个监听器通过`ResultChan`从 go 的 channel 中访问事件通知。

`func main() { ...     api := clientset.CoreV1()     listOptions := metav1.ListOptions{         LabelSelector: label,          FieldSelector: field,     }     watcher, err :=api.PersistentVolumeClaims(ns).        Watch(listOptions)     if err != nil {       log.Fatal(err)     }     ch := watcher.ResultChan() ... }`

### 循环事件

接下来我们将要处理资源事件。但是在我们处理事件之前，我们先声明`resource.Quantity`类型的的两个变量为`maxClaimsQuant`和`totalClaimQuant`来分别表示我们的申请资源阈值（译者注：代表某个 ns 下集群中运行的 PVC 申请的上限）和运行总数。

`import(     "k8s.io/apimachinery/pkg/api/resource"     ... ) func main() {     var maxClaims string     flag.StringVar(&maxClaims, "max-claims", "200Gi",          "Maximum total claims to watch")     var totalClaimedQuant resource.Quantity     maxClaimedQuant := resource.MustParse(maxClaims) ...     ch := watcher.ResultChan()     for event := range ch {         pvc, ok := event.Object.(*v1.PersistentVolumeClaim)         if !ok {             log.Fatal("unexpected type")         }         ...     } }`

在上面的`for-range`循环中，watcher 的 channel 用于处理来自服务器传入的通知。每个事件赋值给变量 event，并且`event.Object`的类型被声明为`PersistentVolumeClaim`类型，所以我们能从中提取出来。

### 处理 ADDED 事件

当一个新的 PVC 创建的时候，`event.Type`的值被设置为`watch.Added`。然后我们用下面的代码去获取新增的声明的容量（`quant`），将其添加到正在运行的总容量中（`totalClaimedQuant`）。最后我们去检查是否当前的容量总值大于当初设定的最大值 (`maxClaimedQuant`)，如果大于的话我们就触发一个事件。

`import(     "k8s.io/apimachinery/pkg/watch"     ... ) func main() { ...     for event := range ch {pvc, ok := event.Object.(*v1.PersistentVolumeClaim)         if !ok {             log.Fatal("unexpected type")         }         quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]         switch event.Type {             case watch.Added:                 totalClaimedQuant.Add(quant)                 log.Printf("PVC %s added, claim size %s\n",                      pvc.Name, quant.String())                 if totalClaimedQuant.Cmp(maxClaimedQuant) == 1 {                     log.Printf(                         "\nClaim overage reached: max %s at %s",                         maxClaimedQuant.String(),                         totalClaimedQuant.String())                     // trigger action                     log.Println("*** Taking action ***")                 }             }         ...         }     } }`

### 处理 DELETED 事件

代码也会在 PVC 被删除的时候做出反应，它执行相反的逻辑以及把被删除的这个 PVC 申请的容量在正在运行的容量的总值里面减去。

`func main() { ...     for event := range ch {         ...         switch event.Type {         case watch.Deleted:             quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]             totalClaimedQuant.Sub(quant)             log.Printf("PVC %s removed, size %s\n",                 pvc.Name, quant.String())             if totalClaimedQuant.Cmp(maxClaimedQuant) <= 0 {                 log.Printf("Claim usage normal: max %s at %s",                     maxClaimedQuant.String(),                     totalClaimedQuant.String(),                 )                 // trigger action                 log.Println("*** Taking action ***")             }         }         ...     } }`

## 运行程序

当程序在一个运行中的集群被执行的时候，首先会列出 PVC 的列表。然后开始监听集群中新的`PersistentVolumeClaim`事件。

\`$> ./pvcwatch

Using kubeconfig:  /Users/vladimir/.kube/config

\--- PVCs ----

NAME                            STATUS  CAPACITY

my-redis-redis                  Bound   50Gi

my-redis2-redis                 Bound   100Gi

---

Total capacity claimed: 150Gi

---

\--- PVC Watch (max claims 200Gi) ----

2018/02/13 21:55:03 PVC my-redis2-redis added, claim size 100Gi

2018/02/13 21:55:03

At 50.0% claim capcity (100Gi/200Gi)

2018/02/13 21:55:03 PVC my-redis-redis added, claim size 50Gi

2018/02/13 21:55:03

At 75.0% claim capcity (150Gi/200Gi)

\`

下面让我们部署一个应用到集群中，这个应用会申请`75Gi`容量的存储。（例如，让我们通过 helm 去部署一个实例 influxdb）。

`helm install --name my-influx \ --set persistence.enabled=true,persistence.size=75Gi stable/influxdb`

正如下面你看到的，我们的工具立刻反应出来有个新的声明以及一个警告因为当前的运行的声明总量已经大于我们设定的阈值。

`--- PVC Watch (max claims 200Gi) ---- ... 2018/02/13 21:55:03 At 75.0% claim capcity (150Gi/200Gi) 2018/02/13 22:01:29 PVC my-influx-influxdb added, claim size 75Gi 2018/02/13 22:01:29 Claim overage reached: max 200Gi at 225Gi 2018/02/13 22:01:29 *** Taking action *** 2018/02/13 22:01:29 At 112.5% claim capcity (225Gi/200Gi)`

相反，从集群中删除一个 PVC 的时候，该工具会相应展示提示信息。

`... At 112.5% claim capcity (225Gi/200Gi) 2018/02/14 11:30:36 PVC my-redis2-redis removed, size 100Gi 2018/02/14 11:30:36 Claim usage normal: max 200Gi at 125Gi 2018/02/14 11:30:36 *** Taking action ***`

## 总结

这篇文章是进行的系列的一部分, 使用 Go 语言的官方 k8s 客户端与 API server 进行交互。和以前一样，这个代码会逐步的去实现一个命令行工具去监听指定 namespace 下面的 PVC 的大小。这个代码实现了一个简单的监听列表去触发从服务器返回的资源事件。

_文章转自：Linux 爱好者_

**![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)**

**Django Model 背后原理是什么？屡屡在面试中将你难倒的元编程究竟难在哪里？锁定 3 日（周三）晚 20:00，Python 大牛讲师 Wayne 带你搞懂 Django 背后的元编程！扫码报名立即参加 ↓↓**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/b2d2e933-67a9-43d6-8061-144b115d0818/640)
