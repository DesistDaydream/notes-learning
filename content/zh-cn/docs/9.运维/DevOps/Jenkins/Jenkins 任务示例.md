---
title: Jenkins 任务示例
---

#

通过 Jenkins 执行 shell 命令

通过 Jenkins 流水线功能 构建 go 项目 hello world

通过 Jenkins 执行 shell 命令

在 Jenkins 首页点击新建任务，选择构建一个自由风格的软件项目，并输入名称，点击确定

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848094-96bfa9b9-8514-46db-ac32-6f2155bc57a4.jpeg)

选择构建标签，增加执行 shell 步骤，

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848109-395b9055-45e1-4fe1-81a9-b37dcee6b2b9.jpeg)

在命令框中，输入想要执行的 shell 命令(类似于 shell 脚本)，点击保存

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848105-6278dba5-7768-4937-9112-bc73e552a7ad.jpeg)

点击立即构建后，Jenkins 就会根据配置的 shell 命令执行任务。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848120-db58946b-bd06-4bc8-ad3e-bda7771d493c.jpeg)

点击下面红框中内容可以查看该次任务的详细信息，点击控制台输出，可以观察任务执行过程以及结果

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848107-37f21311-066a-4ec3-a9e8-091186e22181.jpeg)

至此，一次简单的 Jenkins 任务就执行完成了。

Note：

1. 本示例是使用 docker 运行的 Jenkins。并且本机器为 k8s 集群的 master 节点

2. 由于 Jenkins 运行在 docker 环境中，是获取不到 docker 外部的 kubectl 命令的

3. 所以本示例是将 kubectl 命令以及配置文件拷贝到 Jenkins 的容器中，才让该任务成功执行，否则该任务执行失败，并提示无法找到 kubectl 命令。

通过 Jenkins 流水线功能 构建 go 项目 hello world

新建一个流水线任务

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zab1lx/1616077848144-9811b384-6531-4994-87ef-b9e8bc660e90.jpeg)

点击流水线标签，在脚本中输入 Groovy 语言来指明流水线(pipeline)的步骤。

    pipeline {
      # 使用 docker 来代理执行本次 pipeline,代理镜像为golang。后续所有阶段的步骤都会在，通过该镜像运行的容器内进行。
      agent { docker 'golang' }
      # 指明本次 pileline 的各个阶段
      stages {
        # 第一阶段名为 wget
        stage('wget') {
          # 该阶段第一步，执行 shell 命令，下载main.go文件
          steps {
            sh 'wget https://raw.githubusercontent.com/DesistDaydream/jenkins_pipeline/master/main.go -O main.go'
          }
        }
        # 第二阶段名为 run
        stage('run') {
          # 该阶段第一步，执行 shell 命令运行go程序
          steps {
            sh 'go run main.go'
          }
        }
      }
    }

点击保存后，开始立即构建，并观察控制台输出内容

    Console Output
    Started by user desistdaydream
    Running in Durability level: MAX_SURVIVABILITY
    [Pipeline] Start of Pipeline
    [Pipeline] node
    Running on Jenkins in /var/jenkins_home/workspace/go-hello-world
    [Pipeline] {
    [Pipeline] isUnix
    [Pipeline] sh
    + docker inspect -f . golang
    .
    # 判断本次流水线代理
    [Pipeline] withDockerContainer
    Jenkins seems to be running inside container 05695fdb9f0d2f41b9ba37335af3083539a4f50b0df119bbd8955b8f5c8a6fd9
    # 当使用 docker 作为流水线代理时，会根据指定的镜像(这里指定的是golang)来启动一个容器，后续所有步骤都会在启动的容器中进行。
    $ docker run -t -d -u 0:0 -w /var/jenkins_home/workspace/go-hello-world --volumes-from 05695fdb9f0d2f41b9ba37335af3083539a4f50b0df119bbd8955b8f5c8a6fd9 -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** -e ******** golang cat
    $ docker top 9491ee360563c343b9b2fc603b65d1dc0c1cfa06cbc4c1c3530cd9b7922042da -eo pid,comm
    # 流水线第一阶段，下载 main.go
    [Pipeline] {
    [Pipeline] stage
    [Pipeline] { (wget)
    [Pipeline] sh
    + wget https://raw.githubusercontent.com/DesistDaydream/jenkins_pipeline/master/main.go -O main.go
    --2020-07-31 16:10:42--  https://raw.githubusercontent.com/DesistDaydream/jenkins_pipeline/master/main.go
    Resolving raw.githubusercontent.com (raw.githubusercontent.com)... 151.101.108.133
    Connecting to raw.githubusercontent.com (raw.githubusercontent.com)|151.101.108.133|:443... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: 71 [text/plain]
    Saving to: 'main.go.1'
         0K                                                       100% 1.99M=0s
    2020-07-31 16:10:43 (1.99 MB/s) - 'main.go.1' saved [71/71]
    [Pipeline] }
    # 流水线第二阶段，运行 main.go
    [Pipeline] // stage
    [Pipeline] stage
    [Pipeline] { (run)
    [Pipeline] sh
    + go run main.go
    hi,jenkins
    [Pipeline] }
    [Pipeline] // stage
    [Pipeline] }
    # 所有阶段执行完成后，停止并删除流水线代理容器
    $ docker stop --time=1 9491ee360563c343b9b2fc603b65d1dc0c1cfa06cbc4c1c3530cd9b7922042da
    $ docker rm -f 9491ee360563c343b9b2fc603b65d1dc0c1cfa06cbc4c1c3530cd9b7922042da
    [Pipeline] // withDockerContainer
    [Pipeline] }
    [Pipeline] // node
    [Pipeline] End of Pipeline
    # 流水线结束，输出流水线工作结果。
    Finished: SUCCESS
