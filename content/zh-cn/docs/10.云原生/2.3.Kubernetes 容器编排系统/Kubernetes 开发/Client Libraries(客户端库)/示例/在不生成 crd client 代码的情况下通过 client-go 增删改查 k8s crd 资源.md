---
title: 在不生成 crd client 代码的情况下通过 client-go 增删改查 k8s crd 资源
---

原文链接：<https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html>
2020-07-19

[k8s](https://mozillazg.com/category/k8s.html) [k8s](https://mozillazg.com/tag/k8s.html) / [kubernetes](https://mozillazg.com/tag/kubernetes.html) / [crd](https://mozillazg.com/tag/crd.html) / [client-go](https://mozillazg.com/tag/client-go.html)

- [前言¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidid1)
- [示例 CRD¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidcrd)
- [list 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidlist)
- [get 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidget)
- [create 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidcreate)
- [update 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidupdate)
- [patch 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidpatch)
- [delete 资源¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hiddelete)
- [总结¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidid2)
- [参考资料¶](https://mozillazg.com/2020/07/k8s-kubernetes-client-go-list-get-create-update-patch-delete-crd-resource-without-generate-client-code-update-or-create-via-yaml.html#hidid3)

## 前言

一般情况下管理 crd 资源都是通过由 [code-generator](https://github.com/kubernetes/code-generator) 生成的 crd client 来操作，但是有时也会有只想简单的操作一下资源不想去导入或生成 crd client 相关代码的需求，这里简单的记录一下在不生成 crd client 代码的情况下通过 client-go 增删改查 k8s crd 资源的方法。

## 示例 CRD

先来定义一个测试用的 CRD （其实已有的 Pod 之类的也是可以的，没啥特别的不一定要自定义 CRD，这里只是展示这个能力，因为一般如果是内置的资源的话，直接用内置的 client 和内置的资源 struct 就可以了）（这个 crd 来自 [官方文档](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) ）：

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: crontabs.stable.example.com
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: stable.example.com
  # list of versions supported by this CustomResourceDefinition
  versions:
    - name: v1
      # Each version can be enabled/disabled by Served flag.
      served: true
      # One and only one version must be marked as the storage version.
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                cronSpec:
                  type: string
                image:
                  type: string
                replicas:
                  type: integer
  # either Namespaced or Cluster
  scope: Namespaced
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: crontabs
    # singular name to be used as an alias on the CLI and for display
    singular: crontab
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: CronTab
    # shortNames allow shorter string to match your resource on the CLI
    shortNames:
      - ct
```

然后通过 kubectl 创建一下这个 crd ，然后再创建几个 crd 对象

```shell
$ kubectl apply -f crd.yaml
customresourcedefinition.apiextensions.k8s.io/crontabs.stable.example.com created

$ cat c.yaml
---
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-1
spec:
  cronSpec: "* * * * */5"
  image: my-awesome-cron-image-1

---
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-2
spec:
  cronSpec: "* * * * */8"
  image: my-awesome-cron-image-2

---
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-3
spec:
  cronSpec: "* * * * */10"
  image: my-awesome-cron-image-3

$ kubectl apply -f c.yaml
crontab.stable.example.com/cron-1 created
crontab.stable.example.com/cron-2 created
crontab.stable.example.com/cron-3 created

$ kubectl get crontab.stable.example.com
NAME     AGE
cron-1   9s
cron-2   9s
cron-3   9s
```

## list 资源

首先是如何 list 前面创建的 3 个资源，类似 kubectl get [crontab.stable.example.com](http://crontab.stable.example.com/) 的效果。

简单来说就是通过 [k8s.io/client-go/dynamic](http://k8s.io/client-go/dynamic) 里的 Interface 提供的方法来操作 crd 资源。 关键是怎么拿到 NamespaceableResourceInterface 实例以及把结果转换为自定义的结构体。

完整的 list 资源的代码如下:

```go
package main

import (
        "encoding/json"
        "fmt"
        "os"
        "path/filepath"

        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        "k8s.io/apimachinery/pkg/runtime/schema"
        "k8s.io/client-go/dynamic"
        "k8s.io/client-go/tools/clientcmd"
)

var gvr = schema.GroupVersionResource{
        Group:    "stable.example.com",
        Version:  "v1",
        Resource: "crontabs",
}

type CrontabSpec struct {
        CronSpec string `json:"cronSpec"`
        Image    string `json:"image"`
}

type Crontab struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec CrontabSpec `json:"spec,omitempty"`
}

type CrontabList struct {
        metav1.TypeMeta `json:",inline"`
        metav1.ListMeta `json:"metadata,omitempty"`

        Items []Crontab `json:"items"`
}

func listCrontabs(client dynamic.Interface, namespace string) (*CrontabList, error) {
        list, err := client.Resource(gvr).Namespace(namespace).List(metav1.ListOptions{})
        if err != nil {
                return nil, err
        }
        data, err := list.MarshalJSON()
        if err != nil {
                return nil, err
        }
        var ctList CrontabList
        if err := json.Unmarshal(data, &ctList); err != nil {
                return nil, err
        }
        return &ctList, nil
}

func main() {
        kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
        config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
                panic(err)
        }

        client, err := dynamic.NewForConfig(config)
        if err != nil {
                panic(err)
        }
        list, err := listCrontabs(client, "default")
        if err != nil {
                panic(err)
        }
        for _, t := range list.Items {
                fmt.Printf("%s %s %s %s\n", t.Namespace, t.Name, t.Spec.CronSpec, t.Spec.Image)
        }
}
```

执行结果如下:

```shell
$ go run main.go
default cron-1 * * * * */5 my-awesome-cron-image-1
default cron-2 * * * * */8 my-awesome-cron-image-2
default cron-3 * * * * */10 my-awesome-cron-image-3
```

代码相对来说比较简单，有一个要注意的地方就是 gvr 里各个字段的值来自 crd 定义的 yaml 文件:

```yaml
spec:
  # group name to use for REST API: /apis/<group>/<version>
  # 对应 Group 字段的值
  group: stable.example.com
  # list of versions supported by this CustomResourceDefinition
  versions:
    - name: v1 # 对应 Version 字段的可选值
  # ...
names:
  # plural name to be used in the URL: /apis/<group>/<version>/<plural>
  # 对应 Resource 字段的值
  plural: crontabs
```

注意：因为这个 crd 定义的是 namespace 资源，如果是非 namespace 资源的话，应当改为使用不指定 namespace 的方法:

```go
client.Resource(gvr).List(metav1.ListOptions{})
```

## get 资源

get 资源的方法也是通过 dynamic.Interface 来实现，关键是怎么把结果转换为上面定义的结构体， 关键代码示例如下:

```go
func getCrontab(client dynamic.Interface, namespace string, name string) (*Crontab, error) {
        utd, err := client.Resource(gvr).Namespace(namespace).Get(name, metav1.GetOptions{})
        if err != nil {
                return nil, err
        }
        data, err := utd.MarshalJSON()
        if err != nil {
                return nil, err
        }
        var ct Crontab
        if err := json.Unmarshal(data, &ct); err != nil {
                return nil, err
        }
        return &ct, nil
}

func main() {
        // ...
        ct, err := getCrontab(client, "default", "cron-1")
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)
}
```

执行效果:

```latex
$ go run main.go
default cron-1 * * * * */5 my-awesome-cron-image-1
```

## create 资源

create 资源的方法也是通过 dynamic.Interface 来实现 ，这里主要记录一下怎么基于 yaml 文本的内容来创建资源。

关键代码示例如下:

```latex
func createCrontabWithYaml(client dynamic.Interface, namespace string, yamlData string) (*Crontab, error) {
        decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
        obj := &unstructured.Unstructured{}
        if _, _, err := decoder.Decode([]byte(yamlData), &gvk, obj); err != nil {
                return nil, err
        }

        utd, err := client.Resource(gvr).Namespace(namespace).Create(obj, metav1.CreateOptions{})
        if err != nil {
                return nil, err
        }
        data, err := utd.MarshalJSON()
        if err != nil {
                return nil, err
        }
        var ct Crontab
        if err := json.Unmarshal(data, &ct); err != nil {
                return nil, err
        }
        return &ct, nil
}

func main() {
        // ...
        createData := `
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-4
spec:
  cronSpec: "* * * * */15"
  image: my-awesome-cron-image-4
`
        ct, err := createCrontabWithYaml(client, "default", createData)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)
}
```

执行效果:

```latex
$ go run main.go
default cron-4 * * * * */15 my-awesome-cron-image-4

$ kubectl get crontab.stable.example.com cron-4
NAME     AGE
cron-4   5m33s
```

## update 资源

update 资源的方法也是通过 dynamic.Interface 来实现 ，这里主要记录一下怎么基于 yaml 文本的内容来更新资源。

关键代码示例如下:

```latex
func updateCrontabWithYaml(client dynamic.Interface, namespace string, yamlData string) (*Crontab, error) {
        decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
        obj := &unstructured.Unstructured{}
        if _, _, err := decoder.Decode([]byte(yamlData), &gvk, obj); err != nil {
                return nil, err
        }

        utd, err := client.Resource(gvr).Namespace(namespace).Get(obj.GetName(), metav1.GetOptions{})
        if err != nil {
                return nil, err
        }
        obj.SetResourceVersion(utd.GetResourceVersion())
        utd, err = client.Resource(gvr).Namespace(namespace).Update(obj, metav1.UpdateOptions{})
        if err != nil {
                return nil, err
        }

        data, err := utd.MarshalJSON()
        if err != nil {
                return nil, err
        }
        var ct Crontab
        if err := json.Unmarshal(data, &ct); err != nil {
                return nil, err
        }
        return &ct, nil
}

func main() {
        // ...
        updateData := `
apiVersion: "stable.example.com/v1"
kind: CronTab
metadata:
  name: cron-2
spec:
  cronSpec: "* * * * */15"
  image: my-awesome-cron-image-2-update
`
        ct, err := updateCrontabWithYaml(client, "default", updateData)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%s %s %s %s\n", ct.Namespace, ct.Name, ct.Spec.CronSpec, ct.Spec.Image)

}
```

执行效果:

```latex
 $ kubectl get crontab.stable.example.com cron-2 -o jsonpath='{.spec}'
map[cronSpec:* * * * */8 image:my-awesome-cron-image-2]

$ go run main.go
default cron-2 * * * * */15 my-awesome-cron-image-2-update

$ kubectl get crontab.stable.example.com cron-2 -o jsonpath='{.spec}'
map[cronSpec:* * * * */15 image:my-awesome-cron-image-2-update]
```

## patch 资源

patch 资源的方法跟 patch pod 之类的代码类似，关键代码示例如下:

```latex
func patchCrontab(client dynamic.Interface, namespace, name string, pt types.PatchType, data []byte) error {
        _, err := client.Resource(gvr).Namespace(namespace).Patch(name, pt, data, metav1.PatchOptions{})
        return err
}

func main() {
// ...
patchData := []byte(`{"spec": {"image": "my-awesome-cron-image-1-patch"}}`)
if err := patchCrontab(client, "default", "cron-1", types.MergePatchType, patchData); err != nil {
        panic(err)
}
}
```

执行效果:

```latex
$ kubectl get crontab.stable.example.com cron-1 -o jsonpath='{.spec.image}'
my-awesome-cron-image-1

$ go run main.go

$ kubectl get crontab.stable.example.com cron-1 -o jsonpath='{.spec.image}'
my-awesome-cron-image-1-patch
```

## delete 资源

delete 资源相对来说简单很多，关键代码示例如下:

```latex
func deleteCrontab(client dynamic.Interface, namespace string, name string) error {
        return client.Resource(gvr).Namespace(namespace).Delete(name, nil)
}

func main() {
        // ...
        if err := deleteCrontab(client, "default", "cron-3"); err != nil {
                panic(err)
        }
}
```

结果:

```latex
$ go run main.go
$ kubectl get crontab.stable.example.com
NAME     AGE
cron-1   4h5m
cron-2   4h5m
```

## 总结

简单记录了一下 list、get、create、update、patch、delete crd 资源的方法，其他方法大同小异就没记录了。 简单来说就是可以通过 dynamic.Interface 在不生成特定的 client 代码的情况下操作 crd 资源。

BTW, 另外一个非常规的操作 crd 资源的办法就是直接请求 api server 的 rest api 而不是借助封装好的方法，后面有时间的时候再记录一下这个方法。

## 参考资料

- [kubernetes/client-go: Go client for Kubernetes.](https://github.com/kubernetes/client-go/)
- [kubernetes/code-generator: Generators for kube-like API types](https://github.com/kubernetes/code-generator)
- [Extend the Kubernetes API with CustomResourceDefinitions | Kubernetes](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/)
- [Access Clusters Using the Kubernetes API | Kubernetes](https://kubernetes.io/docs/tasks/administer-cluster/access-cluster-api/)
