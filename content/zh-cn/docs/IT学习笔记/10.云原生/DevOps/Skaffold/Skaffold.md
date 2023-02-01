---
title: Skaffold
---

# 概述

> 参考：
> - [官方文档](https://skaffold.dev/docs/)

Skaffold(脚手架) 是一个命令行工具，可以使得 kubernetes 原生应用的持续开发(CI/CD)变得更加简单。skaffold 可以处理构建(代码与 docker 镜像)、推送到仓库、部署到本地或者远程 kubernetes 集群中这一工作流程。使用一个简单的 yaml 配置文件来定义和执行 Pipeline。

Skaffold 特性

1. Skaffold 更轻量。与 Drone 和 Jenkins 等工具不同，它仅是一个命令行工具，并没有任何附加的组件。所以 Skaffold 也就无法与 SCM 对接来接收 webhook 消息。
   1. 不过也正是因为这种轻量的特性，可以将 Skaffold 集成在 Jenkins 或 Drone 中，作为 shell 命令来执行 CI/CD 的过程，使得整个流水线更加简洁明了。
2. 在本地进行 Kubernetes 的快速开发。在 skaffold 运行之后，当在本地修改完代码之后，skaffold 可以自动触发流水线将应用程序部署到本地或者远程的 kubernetes 集群。

# Skaffold 部署

下载 skaffold 命令行工具的二进制文件，并将其放在$PATH 目录中

1. curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
2. chmod +x skaffold
3. mv skaffold /usr/local/bin

完成后，直接使用 skaffold 命令即可

# Skaffold 配置

skaffold.yaml #Skaffold 通过该文件定义来定义 Pipeline 任务。

/root/.skaffold/config #一些

# Skaffold 的基本使用示例

示例文件获取方式：

1. git clone <https://github.com/GoogleContainerTools/skaffold>
2. cd examples/getting-started

下面是想让 skaffold 正常工作所需的基本文件

    [root@master-1 getting-started]# ls
    Dockerfile  k8s-pod.yaml  main.go  README.md  skaffold.yaml

skaffold.yaml #skaffold 遵循该文件内容执行流水线过程

main.go #需要自动构建部署的代码文件

Dockerfile #用于构建 image

k8s-pod.yaml #用于部署到 k8s 的 yaml 文件

下面是这几个文件的内容

    cat > skaffold.yaml << EOF
    apiVersion: skaffold/v2alpha1
    kind: Config
    build:
      artifacts:
      - image: lchdzh/skaffold-example #指定构建完成后要上传的image信息
    deploy:
      kubectl:
        manifests:
          - k8s-*.yaml #指定上传完images后，通过哪些yaml文件来部署到kubernetes集群中
    EOF


    cat > main.go <<EOF
    package main
    import (
    	"fmt"
    	"time"
    )
    func main() {
    	for {
    		fmt.Println("Hello world!")
    		time.Sleep(time.Second * 1)
    	}
    }
    EOF


    cat > Dockerfile << EOF
    FROM golang:1.14.7 as builder
    COPY main.go .
    RUN go build -o /app main.go
    FROM alpine:3.10
    COPY --from=builder /app .
    CMD ["./app"]
    EOF


    cat > k8s-pod.yaml << EOF
    apiVersion: v1
    kind: Pod
    metadata:
      name: getting-started
    spec:
      containers:
      - name: getting-started
        image: lchdzh/skaffold-example
    EOF

## skaffold 的基本工作流程

1. 使用 Dockerfile 从源代码构建 image
2. 用 sha256 哈希值标记 image
3. 使用先前构建的 image 更新 kubernetes 清单
4. 使用 kubectl apply -f 部署
5. 从已部署的应用输出日志

## skaffold 运行效果

在该目录下执行 skaffold dev 命令，即可开始构建并部署应用，且会持续监控代码变化并以新的代码重新构建并部署

    [root@master-1 getting-started]# skaffold dev
    Listing files to watch...
     - lchdzh/skaffold-example
    Generating tags...
     - lchdzh/skaffold-example -> lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty
    Checking cache...
     - lchdzh/skaffold-example: Not found. Building
    Building [lchdzh/skaffold-example]...
    Sending build context to Docker daemon  3.072kB
    Step 1/6 : FROM golang:1.12.9-alpine3.10 as builder
     ---> e0d646523991
    Step 2/6 : COPY main.go .
     ---> Using cache
     ---> 1ba02146234e
    Step 3/6 : RUN go build -o /app main.go
     ---> Using cache
     ---> f10a56b4422c
    Step 4/6 : FROM alpine:3.10
     ---> 965ea09ff2eb
    Step 5/6 : CMD ["./app"]
     ---> Using cache
     ---> 4171bc10f8b7
    Step 6/6 : COPY --from=builder /app .
     ---> Using cache
     ---> 175e0527448f
    Successfully built 175e0527448f
    Successfully tagged lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty
    The push refers to repository [docker.io/lchdzh/skaffold-example]
    aa7e6c327313: Preparing
    77cae8ab23bf: Preparing
    77cae8ab23bf: Mounted from library/alpine
    aa7e6c327313: Pushed
    v1.1.0-29-gd88d5d0-dirty: digest: sha256:d5556b91934c1c48e3d67b28f0ef431985f63e6a7d37d121281ebf7dc844b803 size: 739
    Tags used in deployment:
     - lchdzh/skaffold-example -> lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty@sha256:d5556b91934c1c48e3d67b28f0ef431985f63e6a7d37d121281ebf7dc844b803
    Starting deploy...
     - pod/getting-started created
    Watching for changes...
    [getting-started] Hello world!
    [getting-started] Hello world!
    [getting-started] Hello world!

后续会持续输出 hello world，并监控代码变化，当 main.go 文件变化之后，将会触发流水线，并构建部署新的应用，效果如下

    [getting-started] Hello world!
    [getting-started] Hello world!
    Generating tags...
     - lchdzh/skaffold-example -> lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty
    Checking cache...
     - lchdzh/skaffold-example: Not found. Building
    Building [lchdzh/skaffold-example]...
    Sending build context to Docker daemon  3.072kB
    Step 1/6 : FROM golang:1.12.9-alpine3.10 as builder
     ---> e0d646523991
    Step 2/6 : COPY main.go .
     ---> 4e1d88698dab
    Step 3/6 : RUN go build -o /app main.go
     ---> Running in 6dc07a2fe012
     ---> 35baee78677b
    Step 4/6 : FROM alpine:3.10
     ---> 965ea09ff2eb
    Step 5/6 : CMD ["./app"]
     ---> Using cache
     ---> 4171bc10f8b7
    Step 6/6 : COPY --from=builder /app .
     ---> f7cb2f7531a0
    Successfully built f7cb2f7531a0
    Successfully tagged lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty
    The push refers to repository [docker.io/lchdzh/skaffold-example]
    abee0c770898: Preparing
    77cae8ab23bf: Preparing
    77cae8ab23bf: Layer already exists
    abee0c770898: Pushed
    v1.1.0-29-gd88d5d0-dirty: digest: sha256:62029ef04493af6c37d8034f1a592285e9e92c120e47cec439657a84ba9985ba size: 739
    Tags used in deployment:
     - lchdzh/skaffold-example -> lchdzh/skaffold-example:v1.1.0-29-gd88d5d0-dirty@sha256:62029ef04493af6c37d8034f1a592285e9e92c120e47cec439657a84ba9985ba
    Starting deploy...
     - pod/getting-started configured
    Watching for changes...
    [getting-started] Hello Skaffold!
    [getting-started] Hello Skaffold!
    [getting-started] Hello Skaffold!

当我们使用 kubectl 查看当前集群 pod 时，会发现该 pod 已经被部署到集群中了

    [root@master-1 ~]# kubectl get pod
    NAME              READY   STATUS    RESTARTS   AGE
    getting-started   1/1     Running   1          4h59m
