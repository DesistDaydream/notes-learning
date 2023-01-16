---
title: Pod 的生命周期，Probe(探针)，Hook(钩子)
---

# 概述

> 参考：
> - [官方文档,概念-工作负载-Pods-Pod 的生命周期](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
> - [官方文档,任务-配置 Pods 与 容器-配置 Liveness、Readiness、Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
> - [公众号,YP 小站-怎么使用 Pod 的 liveness 和 readiness 与 startupProbe](https://mp.weixin.qq.com/s/jPkAj2C0ZNHbaSZRwTOk9g)

## Pod 从开始到结束，会有以下几个 phase(阶段)

1. Pending：调度尚未完成。Pod 已经在 apiserver 中创建，但还没有调度到 Node 上面
2. Running：运行中。od 已经调度到 Node 上面，所有容器都已经创建，并且至少有一个容器还在运行或者正在启动
3. Failed：失败。 Pod 调度到 Node 上面后至少有一个容器运行失败（即退出码不为 0 或者被系统终止）
4. Succeeded：已经成功。Pod 调度到 Node 上面后成功运行结束，并且不会重启。使用 kubectl 命令看到的就是 Completed
5. Unkown：得不到该 Pod 的信息。 状态未知，通常是由于 apiserver 无法与 kubelet 通信导致
6. Completed：已完成。主要用于 Job 模式的 Pod，表示该 Job 正常执行结束

# 容器生命周期钩子(Container Lifecycle Hooks)

Pod 启动的时候，先运行多个 Container 的初始化程序，然后运行 Container 的主程序(主程序中可以在开始 postStart 和结尾 postStop 处执行一些用户自定义“钩子”，这个钩子类似于 awk 命令的 START 和 STOP 功能)。以下是两种钩子的描述。

1. postStart： 容器创建后立即执行，注意由于是异步执行，它无法保证一定在 ENTRYPOINT 之前运行。如果失败，容器会被杀死，并根据 RestartPolicy 决定是否重启
2. preStop：容器终止前执行，常用于资源清理。如果失败，容器同样也会被杀死

钩子的回调函数支持两种方式：

1. exec：在容器内执行命令，如果命令的退出状态码是 0 表示执行成功，否则表示失败
2. httpGet：向指定 URL 发起 GET 请求，如果返回的 HTTP 状态码在 \[200, 400) 之间表示请求成功，否则表示失败

postStart 和 preStop 钩子示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: lifecycle-demo
spec:
  containers:
    - name: lifecycle-demo-container
      image: nginx
      lifecycle:
        postStart:
          httpGet:
            path: /
            port: 80
        preStop:
          exec:
            command: ["/usr/sbin/nginx", "-s", "quit"]
```

在 Container 运行起来之后，Pod 会执行对容器的健康状态检查。i.e.探针，下文就会描述

## Pod 里 Container 的 probe(探针)，健康状态检查

对于 Pod 的健康检查状态有以下几点说明：

1. 对于 Pod 中的 Container 有两种检查方式
   1. 存活性 liveness 探测：周期性探测 Container 的活性。如果探测失败那么 Container 将重新启动。无法更新
   2. 就绪状态 readiness 检测：定期探测 Container 中服务的准备情况。如果探测失败的话 Container 将从服务的后端移除(即使用 kubectl get pod 命令中 READY 标签中左侧数字会减少，减少的就是该 Pod 中某个不在准备状态的 Container)。无法更新。
2. 检查方式的探针类型：
   1. exec，发送命令进行检查。在容器中执行一个命令，如果命令的退出状态码为 0，则探针成功，否则失败
   2. tcpSocket：通过 TCP 协议来检查。对指定容器 IP 和 PORT 执行一个 TCP 检查，如果端口是开发的则探针成功，否则失败
   3. httpGet：通过 HTTP 返回的响应报文来检查。对指定容器的 IP、Port、Path 执行一个 http 的 get 请求，如果返回的状态码在 200 到 400 之间则表示成功，否则表示失败
3. Pod 中容器失败时候(存活性)的重启策略，Always，OnFailure，Never：Always(一失败就重启)
4. Pod 删除的时候：先发送 terminal 信号，有一个宽限期，宽限期一过发送终止信号

## Pod 运行中的几种状态

1. CrashLoopBackOff： 容器退出，kubelet 正在将它重启
2. InvalidImageName： 无法解析镜像名称
3. ImageInspectError： 无法校验镜像
4. ErrImageNeverPull： 策略禁止拉取镜像
5. ImagePullBackOff： 正在重试拉取
6. RegistryUnavailable： 连接不到镜像中心
7. ErrImagePull： 通用的拉取镜像出错
8. CreateContainerConfigError： 不能创建 kubelet 使用的容器配置
9. CreateContainerError： 创建容器失败
10. m.internalLifecycle.PreStartContainer 执行 hook 报错
11. RunContainerError： 启动容器失败
12. PostStartHookError： 执行 hook 报错
13. ContainersNotInitialized： 容器没有初始化完毕
14. ContainersNotReady： 容器没有准备完毕
15. ContainerCreating：容器创建中
16. PodInitializing：pod 初始化中
17. DockerDaemonNotReady：docker 还没有完全启动
18. NetworkPluginNotReady： 网络插件还没有完全启动

# Probe 的 yaml 样例

    apiVersion: v1
    kind: Pod
    metadata:
      labels:
        app: nginx
      name: nginx
    spec:
        containers:
        - image: nginx
          imagePullPolicy: Always
          name: http
          livenessProbe:
            httpGet:
              path: /
              port: 80
              httpHeaders:
              - name: X-Custom-Header
                value: Awesome
            initialDelaySeconds: 15
            timeoutSeconds: 1
          readinessProbe:
            exec:
              command:
              - cat
              - /usr/share/nginx/html/index.html
            initialDelaySeconds: 5
            timeoutSeconds: 1
        - name: goproxy
          image: gcr.io/google_containers/goproxy:0.1
          ports:
          - containerPort: 8080
          readinessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 20
