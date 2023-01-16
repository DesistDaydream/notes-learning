---
title: Drone 任务示例
---

说明：

1. Drone 无法通过手动来触发任务，需要在代码仓库开启 webhook，并配置该 webhook 的地址指向 Drone。

2. 这样在提交代码时，会自动通过 webhook 通知 Drone，有新的代码更改，可以触发任务。

# Drone 的基本示例演示

在 git 仓库中的根目录中，需要创建一个名为 .drone.yml 的文件(该文件也可以使用别的名字，如果要使用别的名字，那么在 Drone 的配置中也要进行相应修改)

    kind: pipeline
    type: docker
    name: default
    steps:
    - name: greeting
      image: alpine
      commands:
      - echo hello
      - echo world

项目的根目录下创建好该文件后，由于有 webhook 的存在，GitLab 会自动通知 Dreon 开始任务，如果没通知，那么可以通过下图的方式来手动触发 Webhook(在一个项目的 webhook 设置，手动测试 push)。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eluf5n/1616077689964-63e8cca0-f50d-42ec-8b95-4a3c4982ba2c.jpeg)

Drone 收到 Webhook 的通知后，开始任务，首先会 clone 该仓库，然后启动 alpine 容器，并在容器中，执行 commands 中执行的命令。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eluf5n/1616077689966-c28f6642-8925-48c5-98b8-6d06fdf5612a.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eluf5n/1616077689966-37666074-eae6-4135-9864-9d0264238779.jpeg)

至此，一套最简单的 CI/CD 流程就算完成了

# 通过 Drone 运行 go 语言的 hello world

.drone.yaml 文件内容如下

    kind: pipeline
    name: default
    steps:
    - name: test
      image: golang
      commands:
      - go test
      - go run main.go

在当前目录写一个名为 main.go 的 go 语言代码并提交，触发 Drone 任务。执行效果如下：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/eluf5n/1616077689965-fbbe7b4f-87e1-406f-afcc-3d2f94018ac8.jpeg)

# 构建代码后构建 docker 镜像并推送到镜像仓库

    kind: pipeline
    name: default
    steps:
    - name: build
      image: golang:latest
      commands:
      - go build -o hello-world
    - name: docker
      image: plugins/docker
      settings:
        repo: lchdzh/drone-docker-build-test
        use_cache: true
        username:
          from_secret: docker_username
        password:
          from_secret: docker_password
      tags:
      - latest
