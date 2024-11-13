---
title: client-go连接K8s集群进行pod的增删改查
---

# 背景

最近在看`client-go`源码最基础的部分，`client-go`的四类客户端，`RestClient、ClientSet、DynamicClient、DiscoveryClient`。其中`RestClient`是最基础的客户端，它对`Http`进行了封装，支持`JSON`和`protobuf`格式数据。其它三类客户端都是通过在`RestClient`基础上再次封装而得来。不过我对`ClientSet`和`DynamicClient`傻傻分不清，虽然很多资料上说它两最大区别是，`ClientSet`能够使用预先生成的`Api`和`ApiServer`进行通信；而`DynamicClient`更加强大，不仅仅能够调用预先生成的`Api`，还能够对一些`CRD`资源通过结构化嵌套类型跟`ApiServer`进行通信。意思大致明白前者能够调用`Kubernetes`本地资源类型，后者还可以调用一些自定资源，那么他们究竟是如何跟`ApiServer`进行交互、`Pod`的增删改查呢？本文通过分析`ClientSet`代码和`client-go`客户端调用`Kubernetes`集群的方式来演示下整个交互过程。

# 准备工作

已经安装`Kubernetes`集群和配置本地`IDE`环境

1. 根据`kubernetes`集群版本选择`clone client-go`到本地：`https://github.com/kubernetes/client-go/tree/release-14.0`。
2. 导入到`IDE`。
3. 运行 `examples/create-update-delete-deployment/main.go` 正常情况下会提示如下错误：

```
panic: CreateFile C:\Users\shj\.kube\config: The system cannot find the path spe
cified.
```

错误信息提示很清楚，没有找到本地文件夹下的`config`文件，处理方式也很简单，只需要把你`Kubernetes`集群中`$HOME/.kube/config`复制到本地即可；仔细阅读代码可以发现，也可以通过自行配置客户端连接信息（生产环境慎用）。

4、运行 main 函数即可进行 Pod 增删改查操作。

# client-go 连接 ApiServer 进行 Pod 的增删改查

1. 获取`APIserver`连接地址、认证配置等信息

```go
var kubeconfig *string
 //获取当前用户home文件夹，并获取kubeconfig配置
 if home := homedir.HomeDir(); home != "" {
  kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
 } else {//如果没有获取到，则需要自行配置kubeconfig
  kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
 }
 //把用户传递的命令行参数，解析为响应变量的值
 flag.Parse()
 //加载kubeconfig中的apiserver地址、证书配置等信息
 config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tbu177/311z1fa79ef320e339b2cf84c6924ecaf224)

    debug信息

2、创建`Clientset`客户端

```go
//NewForConfig为给定的配置创建一个新的Clientset（如下图所示包含所有的api-versions，这样做的目的是便于其它
//资源类型对这个Pod进行管理和控制？）。
clientset, err := kubernetes.NewForConfig(config)
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tbu177/311z40ea6463a53d2e1fbffae83e513a8852)

    debug信息

3、创建`deployment`客户端

```go
//这个过程中会把包含RESTClient配置信息、命名空间信息赋值到deploymentsClient中，具体如下图信息。
deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
 //构造deployment
 deployment := &appsv1.Deployment{
  ObjectMeta: metav1.ObjectMeta{
   Name: "demo-deployment",
  },
  Spec: appsv1.DeploymentSpec{
   Replicas: int32Ptr(2),
   Selector: &metav1.LabelSelector{
    MatchLabels: map[string]string{
     "app": "demo",
    },
   },
   Template: apiv1.PodTemplateSpec{
    ObjectMeta: metav1.ObjectMeta{
     Labels: map[string]string{
      "app": "demo",
     },
    },
    Spec: apiv1.PodSpec{
     Containers: []apiv1.Container{
      {
       Name:  "web",
       Image: "nginx:1.12",
       Ports: []apiv1.ContainerPort{
        {
         Name:          "http",
         Protocol:      apiv1.ProtocolTCP,
         ContainerPort: 80,
        },
       },
      },
     },
    },
   },
  },
 }
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tbu177/311zfe33d3d2814176bc1ef2cc4f72eeb2e7)

看到这里之后，`ClientSet`之所以只能处理预先声明的原生资源类型，是因为对象都是使用的内置元数据类型，不存在的自然没有办法使用了，这时我们在看下`DynamicClient`，部分代码如下所示，它使用`unstructured.Unstructured`表示来自 `API Server`的所有对象值。`Unstructured`类型是一个嵌套的`map[string]inferface{}` 值的集合来创建一个内部结构，通过这种方式，可以表示自定义资源`CRD`资源对象。具体示例，请参考：`examples/dynamic-create-update-delete-deployment/main.go`

```go
deployment := &unstructured.Unstructured{
  Object: map[string]interface{}{
   "apiVersion": "apps/v1",
   "kind":       "Deployment",
   "metadata": map[string]interface{}{
    "name": "demo-deployment",
   },
   "spec": map[string]interface{}{
    "replicas": 2,
    "selector": map[string]interface{}{
     "matchLabels": map[string]interface{}{
      "app": "demo",
     },
    },
    "template": map[string]interface{}{
     "metadata": map[string]interface{}{
      "labels": map[string]interface{}{
       "app": "demo",
      },
     },

     "spec": map[string]interface{}{
      "containers": []map[string]interface{}{
       {
        "name":  "web",
        "image": "nginx:1.12",
        "ports": []map[string]interface{}{
         {
          "name":          "http",
          "protocol":      "TCP",
          "containerPort": 80,
         },
        },
       },
      },
     },
    },
   },
  },
 }
```

4、发送`Post`请求

```go
//发送Post请求到Kubernetes集群
result, err := deploymentsClient.Create(deployment)
```

执行下`kubectl get pod`发现 Kubernetes 集群中 Pod 已经创建。5、更新`Pod`

```go
  //尝试更新资源
 retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
  //获取Get（）返回的“result”
  result, getErr := deploymentsClient.Get("demo-deployment", metav1.GetOptions{})
  if getErr != nil {
   panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
  }
  //replica数量降低到1
  result.Spec.Replicas = int32Ptr(1)
  //修改nginx镜像
  result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13"
  //更新（result）
  _, updateErr := deploymentsClient.Update(result)
  return updateErr
 })
```

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tbu177/311ze17817d585f6936fa2f18f17f2678de6)

    更新模板信息

`RetryOnConflict`用于需要考虑更新冲突的情况下对资源进行更新，出现这种场景，大多因为存在其它客户端但或者代码同一时间内操作该资源对象。如果`update`函数返回冲突错误，`RetryOnConflict`将按指定策略等待一段时间退后，再次尝试更新。

6、查询操作

```go
//发送http get请求获取pod列表
list, err := deploymentsClient.List(metav1.ListOptions{})
```

其内部查询接口如下图所示：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tbu177/311zdad11d29dbefd5810a31c0e3a02fa976)

其中`c.client`读取配置实例化`RESTClient`对象和`ns`,其中`deployments`这个对象是在这行`deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)`进行初始化; `Get()`通过`Get`请求，同样支持`POST PUT DELETE PATCH`;`Resource`设置请求的资源名称;`VersionedParams` 设置查询选项，例如：`TimeoutSeconds`;`Do()`执行请求，结果`Into`到`Result`。

7、删除操作

    //指定删除策略
    deletePolicy := metav1.DeletePropagationForeground
    //针对特定deployment进行删除操作
     if err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
      PropagationPolicy: &deletePolicy,
     });

`Kubernetes`控制器的删除有 3 种模式：

- `Foreground`: 删除控制器之前，先删除控制器所管理的资源对象删除。
- `Background`：删除控制器后，控制器所管理的资源对象由`GC`在后台进行删除。
- `Orphan`：只删除控制器，不删除控制器所管理的资源对象(举个例子，比如你删除了`deployment`，那么对应的`Pod`不会被删除)。

8、观察`Pod`变化

    [root@k8s-m1 ~]# kubectl get pod  --watch
    NAME                               READY   STATUS    RESTARTS   AGE
    demo-deployment-5fc8ffdb68-8xdcx   1/1     Running   0          26s
    demo-deployment-5fc8ffdb68-w555g   1/1     Running   0          26s
    demo-deployment-5fc8ffdb68-w555g   1/1     Terminating   0          42s
    demo-deployment-5cb6f65f77-tn5bn   0/1     Pending       0          0s
    demo-deployment-5cb6f65f77-tn5bn   0/1     Pending       0          0s
    demo-deployment-5cb6f65f77-tn5bn   0/1     ContainerCreating   0          0s
    demo-deployment-5fc8ffdb68-w555g   0/1     Terminating         0          43s
    demo-deployment-5cb6f65f77-tn5bn   1/1     Running             0          2s
    demo-deployment-5fc8ffdb68-8xdcx   1/1     Terminating         0          44s
    demo-deployment-5fc8ffdb68-8xdcx   0/1     Terminating         0          45s
    demo-deployment-5fc8ffdb68-8xdcx   0/1     Terminating         0          48s
    demo-deployment-5fc8ffdb68-8xdcx   0/1     Terminating         0          48s
    demo-deployment-5fc8ffdb68-w555g   0/1     Terminating         0          51s
    demo-deployment-5fc8ffdb68-w555g   0/1     Terminating         0          51s
    demo-deployment-5cb6f65f77-tn5bn   1/1     Terminating         0          52s
    demo-deployment-5cb6f65f77-tn5bn   0/1     Terminating         0          53s
    demo-deployment-5cb6f65f77-tn5bn   0/1     Terminating         0          56s
    demo-deployment-5cb6f65f77-tn5bn   0/1     Terminating         0          56s

# 总结

本文主要通过在本地运行`client-go/ClientSet`客户端对`Pod`的增删改查，并解释了代码的执行过程。同时加深了对`ClientSet`客户端的理解。
