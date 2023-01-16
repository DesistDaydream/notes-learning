---
title: get 子命令
---

# 概述

> 参考：
> - [5 个冷门但非常实用的 Kubectl 使用技巧](https://mp.weixin.qq.com/s/DlYcJNNCc9C_YUZlvADuMQ)

展示对象的信息，get 获得的是该对象的个性信息，describe 获得的是该对象的集群信息

# Syntax(语法)

**kubectl get (TYPE\[.VERSION]\[.GROUP] \[NAME | -l label] | TYPE\[.VERSION]\[.GROUP]/NAME ...) \[FLAGS]**

**FLAGS**
Note：在 kubectl 命令中的 全局 flags 中还有很多有用的 flags 可以用于 get 子命令。比如 -v 指定 debug 等级，-n 指定要操作的 namespace，等等

- **-A, --all-namespaces** # 列出在所有名称空间中的对象。
- \--allow-missing-template-keys=true: If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats.
- \--chunk-size=500: Return large lists in chunks rather than all at once. Pass 0 to disable. This flag is beta and may change in the future.
- **--field-selector=''** # 根据一个或多个资源字段的值[筛选 Kubernetes 对象](https://kubernetes.io/zh/docs/concepts/overview/working-with-objects/kubernetes-objects)。支持 '=', '==', and '!='.(比如 --field-selector key1=value1,key2=value2)。注意，仅支持部分字段筛选
- -f, --filename=\[]: Filename, directory, or URL to files identifying the resource to get from a server.
- \--ignore-not-found=false: If the requested object does not exist the command will return exit code 0.
- **-k, --kustomize=<DIR>** # 处理指定的 Kustomize 目录。这个标志不能与 -f 或 -R 同时使用。
- **-L, --label-columns=\[]** # 显示所有展示出的对象具有 KEY 这个键所对应的值(KEY=VAL，显示那个 VAL)Accepts a comma separated list of labels that are going to be presented as columns. Names are case-sensitive. You can also use multiple flag options like -L label1 -L label2...
- **--no-headers** # 当使用 默认的 或者 custom-column 格式输出信息时，不显示标题(标题就是 NAME 那一行)。
- **-o, --output=FORMAT **# 指定输出信息的输出格式
  - FORMAT 包括 json|yaml|wide|name|custom-columns=...|custom-columns-file=...|go-template=...|go-template-file=...|jsonpath=...|jsonpath-file=...
    - 官方说明：<https://kubernetes.io/docs/reference/kubectl/overview/#formatting-output>
  - yaml | json # 输出 yaml 或 json 格式的信息
  - wide # 多显示该对象的 IP 和所在 NODE 两个信息
  - name #仅打印对象名称，而不打印其他任何内容。Note：资源名称格式为：资源类型/名字
  - custom-columns=<HEADER>:<JSON-PATH-EXPR> # 自定义以一列一列的形式显示列表。参考：\[http://kubernetes.io/docs/user-guide/kubectl-overview/#custom-columns]
  - golang template\[http://golang.org/pkg/text/template/#pkg-overview]
  - jsonpath template # 使用 json 格式里的路径来查看某个字段的状态，样例如下，`.`符号是字符分隔符。用法详见[官方文档](https://kubernetes.io/docs/reference/kubectl/jsonpath/)
    - e.g.-o jsonpath="{.status.phase}"
- \--output-watch-events=false: Output watch event objects when --watch or --watch-only is used. Existing objects are output as initial ADDED events.
- **--raw <URL Path>** #从 API Server 请求原始 URI。显示指定 URL Path 路径下的原始 URI 信息，默认输出为 JSON 格式
- -R, --recursive=false: Process the directory used in -f, --filename recursively. Useful when you want to manage related manifests organized within the same directory.
- **-l, --selector=KEY\[=VAL,KEY2=VAL2,...]** # 根据标签对输出进行过滤。可以只指定标签中的 key，或者指定多个 key，或者指定 key 不匹配的 value
- e.g. #-l key1=value1,key2=value2,Note:k/v 中的 = 还可以使用 == 和 !=
- \--server-print=true: If true, have the server return the appropriate table output. Supports extension APIs and CRDs.
- \--show-kind=false # 列出所请求对象的资源类型。
- **--show-labels** # 输出信息时，在最后一列显示该对象的 label。(默认不显示)
- \--sort-by='': If non-empty, sort list types using this field specification. The field specification is expressed as a JSONPath expression (e.g. '{.metadata.name}'). The field in the API resource specified by this JSONPath expression must be an integer or a string.
- \--template='': Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates \[http://golang.org/pkg/text/template/#pkg-overview].
- -w, --watch # 实时监控。类似于在命令前加 wathch。只有当资源状态变化时，才会显示
- \--watch-only # Watch for changes to the requested object(s), without listing/getting first.

# EXAMPLE

- 获取原始 URL，显示 `/` 下的资源
  - kubectl get --raw /
- 查看所有名称空间下的所有资源
  - kubectl get all -A
- 显示 myapp-pod 这个 pod 的信息
  - kubectl get pod/myapp-pod
- 自定义显示内容，仅显示 node 名字 和 标签
  - kubectl get nodes -o custom-columns=NAME:.metadata.name,LABELS:.metadata.labels

```bash
[root@master-test-1 .kube]# kubectl get nodes -o custom-columns=NAME:.metadata.name,LABELS:.metadata.labels
NAME            LABELS
master-test-1   map[beta.kubernetes.io/arch:amd64 beta.kubernetes.io/os:linux kubernetes.io/arch:amd64 kubernetes.io/hostname:master-test-1 kubernetes.io/os:linux node-role.kubernetes.io/controlplane:true node-role.kubernetes.io/etcd:true]
node-test-1     map[beta.kubernetes.io/arch:amd64 beta.kubernetes.io/os:linux kubernetes.io/arch:amd64 kubernetes.io/hostname:node-test-1 kubernetes.io/os:linux node-role.kubernetes.io/worker:true]
node-test-2     map[beta.kubernetes.io/arch:amd64 beta.kubernetes.io/os:linux kubernetes.io/arch:amd64 kubernetes.io/hostname:node-test-2 kubernetes.io/os:linux node-role.kubernetes.io/worker:true]
```

- 显示 pod 的名字及其启动时间。
  - kubectl get pods -o=jsonpath='{range .items\[\*]}{.metadata.name}{"\t"}{.status.startTime}{"\n"}{end}'

```bash
[root@master-test-1 .kube]# kubectl get pods -A -o=jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.startTime}{"\n"}{end}'
myapp-pod	2020-08-21T15:57:35Z
default-http-backend-598b7d7dbd-xrp6s	2020-08-17T05:04:41Z
nginx-ingress-controller-7b9mp	2020-08-17T05:04:40Z
nginx-ingress-controller-lhbgl	2020-08-22T08:17:10Z
coredns-849545576b-7xt9p	2020-08-17T07:27:43Z
```

- 获取所有节点的污点
  - kubectl get nodes --template='{{range .items}}{{.metadata.name}}{{"\t"}}{{.spec.taints}}{{"\n"}}{{end}}'
- 列出事件（Events），按时间戳排序
  - kubectl get events -A --sort-by='{.metadata.creationTimestamp}'
- 删除所有名称空间中标签为 authz.cluster.cattle.io/rtb-owner-updated 的 rolebinding 对象
  - for k in $(kubectl get ns -o name | awk -F/ '{print $2}'); do for i in $(kubectl get -n $k rolebindings.rbac.authorization.k8s.io -l "authz.cluster.cattle.io/rtb-owner-updated" -o name); do kubectl delete -n $k $i; done;done
- 从 kubernetes-dashboard 这个 deployment 中获取 SA，并从 SA 中提取 Token。
  - APPNAME=$(kubectl get deployments.apps -n kubernetes-dashboard -o name)
  - SA=$(kubectl get ${APPNAME} -n kubernetes-dashboard -ojsonpath='{.spec.template.spec.serviceAccountName}')
  - kubectl get secrets -n kubernetes-dashboard -o jsonpath="{.items\[?(@.metadata.annotations\['kubernetes.io/service-account.name']=='${SA}')].data.token}" | base64 -d
- ## 获取 admin TOKEN

- 查看 replicaset 的历史版本号
  - kubectl get replicasets.apps -n bluestore-console -ojsonpath='{range .items\[\*]}{.metadata.annotations.deployment.kubernetes.io/revision}{"\t"}{.metadata.name}{"\n"}{end}'

## 过滤 Pod

获取 kube-system 名称空间下，标签 k8s-app 的值为 kube-dns 的所有 pod。

- kubectl get pod -n kube-system --selector="k8s-app=kube-dns"

获取 node-1 节点上的所有 Pods

- kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=node-1

获取指定状态的 pod。(Succeeded 就是 Completed)

- kubectl get pods -A --field-selector status.phase=Running
- kubectl get pod -A --field-selector status.phase=Succeeded
- kubectl get pod -A --field-selector status.reason=Evicted

获取 test 名称空间下所有资源

- kubectl api-resources -o name --verbs=list --namespaced | xargs -n 1 kubectl get --show-kind --ignore-not-found -n test

获取正在使用 pvc 的 pod

- kubectl get pods --all-namespaces -o=json | jq -c '.items\[] | {name: .metadata.name, namespace: .metadata.namespace, claimName:.spec.volumes\[] | select( has ("persistentVolumeClaim") ).persistentVolumeClaim.claimName }'

### 获取 Pod 在 Node 上的分布数量

```bash
kubectl get pods -A -o wide -l app="flannel" | awk '{print $8}'|\
 awk '{ count[$0]++  }
 END {
   printf("%-35s: %s\n","Word","Count");
   for(ind in count){
    printf("%-35s: %d\n",ind,count[ind]);
   }
 }'
```

若是指定了名称空间，awk 则应筛选 $7 列。

### 获取指定状态的 Pod，并删除

获取 Pod 状态为 Pending 的所有 Pod 并生成删除指令

```bash
export PodStatus="Pending"
kubectl get pods --all-namespaces --field-selector status.phase=${PodStatus} -o json | \
  jq '.items[] | "kubectl get pods \(.metadata.name) -o wide -n \(.metadata.namespace)"'
```

！！！注意：执行 `xargs -n 1 bash -c` 删除操作前，需要详细调试命令
最后加上 `xargs -n 1 bash -c` 以执行生成的删除 Pod 指令

```bash
kubectl get pods --all-namespaces --field-selector status.phase=Pending -o json | \
  jq '.items[] | "kubectl delete pods \(.metadata.name) -n \(.metadata.namespace)"' | \
  xargs -n 1 bash -c
```

如果只是获取单一 Namespace 下的 Pods，直接只用 grep 命令筛选更快~

- kubectl -n default get pods | grep Completed | awk '{print $1}' | xargs kubectl -n default delete pods

## finalizers 字段处理

快速清空对象中的 finalizers 字段

```bash
kubectl get namespace test -o json \
 | tr -d "\n" | sed "s/\"finalizers\": \[[^]]\+\]/\"finalizers\": []/" \
 | kubectl replace --raw /api/v1/namespaces/test/finalize -f -
```
