---
title: Kind：一个容器创建K8S开发集群
---

## **什么是 Kind**

`kind`：是一种使用 Docker 容器`节点`运行本地 Kubernetes 集群的工具。该类型主要用于测试 Kubernetes，但可用于本地开发或 CI。

> 注意：kind 仍在开发中### **Mac & Linux** $ curl -Lo ./kind "https://kind.sigs.k8s.io/dl/v0.9.0/kind-$(uname)-amd64"

    chmod +x ./kind
    mv ./kind /some-dir-in-your-PATH/kind1

2
3
4
Plain Text

### **Mac 上使用 brew 安装** $ brew install kind1

2
Plain Text

### **Windows** $ curl.exe -Lo kind-windows-amd64.exe https://kind.sigs.k8s.io/dl/v0.9.0/kind-windows-amd64

Move-Item .\kind-windows-amd64.exe c:\some-dir-in-your-PATH\kind.exe # OR via Chocolatey (https://chocolatey.org/packages/kind)
$ choco install kind1
2
3
4
5
6
Plain Text

## **K8S 集群创建与删除** # 创建集群，默认集群名称为 kind

$ kind create cluster1
2
3
Plain Text![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sog8v3/1616120539030-1a544e75-7263-4229-8daa-01fc2037dd7c.png) # 定义集群名称
$ kind create cluster --name kind-2 # 查询集群
$ kind get clusters # 删除集群
$ kind delete cluster1
2
3
4
5
6
7
8
9
Plain Text

## **其它操作** # 列出集群镜像

$ docker exec -it my-node-name crictl images1
2
3
Plain Text![](https://notes-learning.oss-cn-beijing.aliyuncs.com/sog8v3/1616120538989-43863854-bbb8-4205-b2ac-8f73b48d3aac.jpeg)

## **参考链接**> <https://github.com/kubernetes-sigs/kind>

> <https://kind.sigs.k8s.io/docs/user/quick-start/#installation>
