---
title: Go Client
---

# 概述

> 参考：
> - [GitHub 项目，kubernetes/client-go](https://github.com/kubernetes/client-go)
> - [Danielhu 博客](https://www.danielhu.cn/tags/client-go/)
> - [公众号-KubeSphere 云原生，Client-go 源码分析之 SharedInformer](https://mp.weixin.qq.com/s/13enj17ifaD-mSrjVzAQKw)

Go Client 就是 Kubernetes 针对 Go 编程语言 而言的 Client Library。Go Client 项目名称为 client-go，是用来与 Kubernetes 对话的 Go 编程语言 的第三方库。

# 安装 client-go

版本控制策略：k8s 版本 1.18.8 对应 client-go 版本 0.18.8，其他版本以此类推。

使用前注意事项：
使用 client-go 之前，需要手动获取对应版本的的 client-go 库。

如果使用的 Kubernetes 版本> = `v1.17.0`，请使用相应的 `v0.x.y`标签。例如，`k8s.io/client-go@v0.17.0`对应于 Kubernetes `v1.17.0`

根据版本控制策略，使用如下命令进行初始化:

```bash
# 初始化项目
go mod init github.com/DesistDaydream/kubernetes-development
# 为 go.mod 文件添加 require k8s.io/client-go v0.19.2 // indirect 信息
go get k8s.io/client-go@v0.19.2
# 整理 go.mod 与 go.sum 文件
go mod tidy
```

至此，才可以正常使用，否则会产生很多依赖问题

# Hello World 示例

参考：[阳明大佬](https://www.notion.so/K8S-b33520bf4f2c4005adb66f5ee1785502)

```go
package main

import (
	"context"
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InOrOut 判断当前环境是在集群内部，还是集群外部
func InOrOut() string {
	// 如果容器内具有环境变量 KUBERNETES_SERVICE_HOST 且不为空，则当前代码是在容器内运行，否则是在集群外部运行
	if h := os.Getenv("KUBERNETES_SERVICE_HOST"); h != "" {
		return "inCluster"
	}
	return "outCluster"
}

// Deployment 获取指定 namespace 下所有的 deployment 对象
func get(clientset *kubernetes.Clientset, namespace string) {
	// 获取指定 名称空间 下所有的 deployment 对象
	deployments, _ := clientset.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
	for i, deploy := range deployments.Items {
		fmt.Printf("%d -> %s\n", i+1, deploy.Name)
	}
}

func main() {
	var config *rest.Config
	// 根据代码所在环境，决定如何创建一个连接集群所需的配置。
	switch InOrOut() {
	case "inCluster":
		// 根据容器内的 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下的 token 与 ca.crt 文件创建一个用于连接集群的配置。
		config, _ = rest.InClusterConfig()
	case "outCluster":
		// 根据指定的 kubeconfig 文件创建一个用于连接集群的配置，/root/.kube/config 为 kubectl 命令所用的 config 文件
		config, _ = clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
		// 注意，clientcmd.BuildConfigFromFlags() 内部实际上也是有调用 rest.InClusterConfig() 的逻辑，只要满足条件即可。条件如下：
		// 若第二个参数为空的话，则会直接调用 rest.InClusterConfig()
	}

	// 根据 BuildConfigFromFlags 创建的配置，返回一个可以连接集群的指针
	clientset, _ := kubernetes.NewForConfig(config)

	// 获取指定 namespace 下所有的 deployment 对象
	get(clientset, "kube-system")
}
```

## 示例详解

从上面的示例中可以看到，通过 API 访问 Kubernetes 集群大体分为两类

- 从集群内部访问
- 从集群外部访问

所谓的从集群内部访问，其实就是代码运行在 pod 中。不过，不管是集群内还是集群外，至少需要一个配置信息来连接集群，而两种访问方式的唯一区别也在于此，不同的访问方式，获取到的配置信息是不同的。

- 但是，不管如何获得，这个配置信息实际上就是 k8s.io/client-go/rest 包中的 `*rest.Config` 结构体。这个结构体中包含了待连接到 Kubernetes 集群的各种信息。
- 然后，使用 `*rest.Config` 来实例化 `*kubernetes.Clientset`，这个 Clientset 就是用来通过 API 来操作 Kubernetes 的客户端集。这个 Clientset 实际上是 Kubernetes API 的集合。

```go
type Clientset struct {
	*discovery.DiscoveryClient
	admissionregistrationV1      *admissionregistrationv1.AdmissionregistrationV1Client
	admissionregistrationV1beta1 *admissionregistrationv1beta1.AdmissionregistrationV1beta1Client
	appsV1                       *appsv1.AppsV1Client
	appsV1beta1                  *appsv1beta1.AppsV1beta1Client
	appsV1beta2                  *appsv1beta2.AppsV1beta2Client
	authenticationV1             *authenticationv1.AuthenticationV1Client
	authenticationV1beta1        *authenticationv1beta1.AuthenticationV1beta1Client
	authorizationV1              *authorizationv1.AuthorizationV1Client
	authorizationV1beta1         *authorizationv1beta1.AuthorizationV1beta1Client
	autoscalingV1                *autoscalingv1.AutoscalingV1Client
	autoscalingV2beta1           *autoscalingv2beta1.AutoscalingV2beta1Client
	autoscalingV2beta2           *autoscalingv2beta2.AutoscalingV2beta2Client
	batchV1                      *batchv1.BatchV1Client
	batchV1beta1                 *batchv1beta1.BatchV1beta1Client
	batchV2alpha1                *batchv2alpha1.BatchV2alpha1Client
	certificatesV1               *certificatesv1.CertificatesV1Client
	certificatesV1beta1          *certificatesv1beta1.CertificatesV1beta1Client
	coordinationV1beta1          *coordinationv1beta1.CoordinationV1beta1Client
	coordinationV1               *coordinationv1.CoordinationV1Client
	coreV1                       *corev1.CoreV1Client
	discoveryV1alpha1            *discoveryv1alpha1.DiscoveryV1alpha1Client
	discoveryV1beta1             *discoveryv1beta1.DiscoveryV1beta1Client
	eventsV1                     *eventsv1.EventsV1Client
	eventsV1beta1                *eventsv1beta1.EventsV1beta1Client
	extensionsV1beta1            *extensionsv1beta1.ExtensionsV1beta1Client
	flowcontrolV1alpha1          *flowcontrolv1alpha1.FlowcontrolV1alpha1Client
	networkingV1                 *networkingv1.NetworkingV1Client
	networkingV1beta1            *networkingv1beta1.NetworkingV1beta1Client
	nodeV1alpha1                 *nodev1alpha1.NodeV1alpha1Client
	nodeV1beta1                  *nodev1beta1.NodeV1beta1Client
	policyV1beta1                *policyv1beta1.PolicyV1beta1Client
	rbacV1                       *rbacv1.RbacV1Client
	rbacV1beta1                  *rbacv1beta1.RbacV1beta1Client
	rbacV1alpha1                 *rbacv1alpha1.RbacV1alpha1Client
	schedulingV1alpha1           *schedulingv1alpha1.SchedulingV1alpha1Client
	schedulingV1beta1            *schedulingv1beta1.SchedulingV1beta1Client
	schedulingV1                 *schedulingv1.SchedulingV1Client
	settingsV1alpha1             *settingsv1alpha1.SettingsV1alpha1Client
	storageV1beta1               *storagev1beta1.StorageV1beta1Client
	storageV1                    *storagev1.StorageV1Client
	storageV1alpha1              *storagev1alpha1.StorageV1alpha1Client
}
```

### 在集群内部生成配置信息

使用 `rest.InClusterConfig()` 方法生成 \*rest.Config。使用该方法时，一般通过 /var/run/secrets/kubernetes.io/serviceaccount/ 目录下的 token 与 ca.crt 文件生成认证信息。通过容器内的环境变量生成待连接的 Server 端信息。

```go
func InClusterConfig() (*Config, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	token, err := ioutil.ReadFile(tokenFile)
	tlsClientConfig := TLSClientConfig{}
	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}
	return &Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}
```

### 在集群外部生成配置信息

使用 `clientcmd.BuildConfigFromFlags()` 方法生成 \*rest.Config。使用该方法时，直接将 kubectl 的配置文件，当作参数传递进去即可。

> 只不过有一点需要注意，就是当参数都为空时，该方法还是会调用的 `rest.InClusterConfig()` 方法。

```go
func BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error) {
	if kubeconfigPath == "" && masterUrl == "" {
		kubeconfig, err := restclient.InClusterConfig()
	}
	return NewNonInteractiveDeferredLoadingClientConfig(
		&ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterUrl}}).ClientConfig()
}
```
