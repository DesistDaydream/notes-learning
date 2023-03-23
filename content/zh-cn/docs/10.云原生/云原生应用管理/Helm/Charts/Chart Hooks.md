---
title: Chart Hooks
---

# Chart Hooks 概述

参考：[**官方文档**](https://helm.sh/docs/topics/charts_hooks/)

Helm 提供了一种** Hook(钩子) **机制可以在一个 release 的生命周期内进行干预，比如：

1.

安装任何资源之前，提前先安装 ConfigMap 或者 Secret

2.

在安装一个新的 Chart 执行，执行一个 Job 以备份数据库，然后在升级后执行第二个 Job 以还原数据。

3.

在删除一个 Release 之前，运行一个 Job，以便在删除服务之前优雅得停止服务。

4.

等等等

说白了：就是让我们在操作 Chart 中的资源时，可以运行一个 Job 或某些资源(比如删除 operator 之前，运行一个 job 先删除所有 CRD 资源)。有点类似与就类似 [**crds 目录**](https://www.teambition.com/project/5f90e312755d8a00446050eb/app/5eba5fba6a92214d420a3219/workspaces/5f90e312c800160016ea22fb/docs/5fae78274cc5830001b9bbd6?scroll-to-block=5fae7842875a1b0329e90945) 的作用一样，但并不完全一样。

## Hooks 种类

1.

pre-install # 在渲染模板之后、创建资源之前，执行安装

2.

post-install # 在 Chart 中所有资源创建之后(并不用等待 running)，执行安装

3.

pre-delete # 在从 Kubernetes 删除任何资源之前，执行安装

4.

post-delete # 删除所有 releases 资源后，执行安装

5.

pre-upgrade # 在渲染模板之后，在任何资源更新之前，执行安装

6.

post-upgrade # 在所有资源都升级后，执行安装

7.

pre-rollback # Executes on a rollback request after templates are rendered, but before any resources are rolled back

8.

post-rollback # Executes on a rollback request after all resources have been modified

9.

test Executes when the Helm test subcommand is invoked ( view test docs)

# Release 的生命周期

1.

运行 helm install myapp

2.

The Helm library install API is called

3.

经过一些验证后，开始渲染 myapp 的模板

4.

渲染后生成的资源加载到 Kubernetes 中

5.

library 将 release 的数据返回给客户端

6.

客户端退出

等待钩子准备好意味着什么？ 这取决于挂钩中声明的资源。 如果资源是 Job 或 Pod 类型，Helm 将等待直到成功运行完成为止。 如果挂钩失败，释放将失败。 这是一项阻止操作，因此 Helm 客户端将在 Job 运行时暂停。

# 使用 Hooks

Helm 会读取 manifest 文件中的 `.annotations."helm.sh/hook"`、`annotations."helm.sh/hook-weight"`、`annotations."helm.sh/hook-delete-policy"`这三个字段，来为具有这三个字段的资源执行 Hooks 功能。
