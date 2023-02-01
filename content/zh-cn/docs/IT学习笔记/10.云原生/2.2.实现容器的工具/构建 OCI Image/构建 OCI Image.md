---
title: 构建 OCI Image
---

# 概述

> 参考：
> - [官方文档，使用 Docker 开发-构建镜像-编写 Dockerfile 的最佳实践](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

在容器刚刚流行的时候，想要构建一个容器镜像通常只有两种方式：

- 通过对 Container 执行 commit 命令来创建基于该 Container 的 Image
- 通过 Dockerfile 功能来构建 Image

Dockerfile 构建镜像的方式逐渐成为主流甚至标准，但是随着各个项目的去 Docker 化，大家都想消除自身对 Docker 的依赖，这其中包括 Docker 项目的起源 Moby，从[这里(moby/moby 的 issue #34227)](https://github.com/moby/moby/issues/34227)可以略窥 12。但是 Dockerfile 的影响已经深入人心，所以各家一时半会也无法完全舍弃，只能说基于 Dockerfile 形式进行优化。时至今日(2022 年 6 月 3 日)，Dockerfile 依然是最常见最通用的构建镜像的方式，不管构建程序是什么，总归是要通过 Dockerfile 文件的。

# Dockerfile

> 参考：
> - [官方文档](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
> - [Dockerfile 参考](https://docs.docker.com/engine/reference/builder/)

Docker 通过读取 **Dockerfile 文件**中的指令来构建符合 OCI 标准的容器镜像。Dockerfile 这个称呼有多个理解方式，可以是指一个**功能**，也可以指一个文件的**文件名**，也可以代指一类文件的**统称。**

## DockerFile 功能的工作逻辑：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ymn9n3/1616121926992-8640dd4a-029c-487c-9fa2-f20ee1b61b83.png)

- 找一个专用目录，在目录中放入默认的名为 Dockerfile 的文件，该文件名首字母必须大写
  - Dockerfile 文件，类似于一个脚本，使用 docker build 命令创建镜像的时候默认使用名为 Dockerfile 的文件，通过该文件中的各种指令来执行操作。（如果想使用其余名字的文件，则需要使用-f 参数来指明需要使用的 DockerFile 的文件，这时候可以使用名字不为 Dockerfile 的文件）
- 如果该 Image 中需要打包进去很多文件(比如 rpm 包、配置文件等等)，则这些文件必须做好后，放到 Dockerfile 所在的目录中(可以有子目录)。
- 使用 docker build 命令并用指定目录路径，则该命令会自动查找该目录下的名为 Dockerfile 文件并根据其中内容创建 Image，效果如上图所示

### 具体过程详解如下

- 首先在当前目录创建一个名为 Dockerfile 的文件，在该文件中写入需要执行的指令

```bash
~]# ll
total 4
-rw-r--r--. 1 root root 49 Nov 27 21:02 test
~]# cat test
FROM ubuntu
RUN apt update && apt install -y vim
RUN apt install -y iproute2
```

- 使用 docker build -t ubuntu-vi -f test /dockerfile/ 命令创建镜像
  - 当创建镜像的时候，会使用命令中定义的 PATH 中的默认名为 Dockerfile 文件中的指令来进行自动操作，可以通过-f 选项来选择指定路径下的 Dockerfile 文件（注：命令会从/dockerfile/目录中查找 Dockerfile 文件，然后把/root 目录中的所有文件发送给 Docker daemon 来使用，所以定义创建环境的时候最好使用一个空目录）
    - sending(发送)build context(创建环境)to(给)docker daemon(容器守护进程) 17.92KB(这个文件 17.92K)

```bash
~]# docker build -t ubuntu-vi -f test /dockerfile/
Sending build context to Docker daemon  2.048kB
```

- 执行文件中 FROM 指令,将 ubuntu 作为基础镜像，IMAGE ID 为 93fd78260bd1


    Step 1/3 : FROM ubuntu
     ---> 93fd78260bd1

- 基于 93fd78260bd1 这个 IMAGE 启动名为 607ce2e8553f 的临时容器，执行 RUN 后面的命令执行文件中的第二行 RUN 指令,安装 VIM


    Step 2/3 : RUN apt update && apt install -y vim
     ---> Running in 607ce2e8553f

- 开始执行安装程序，会有警告，具体过程忽略不截图了


    WARNING: apt does not have a stable CLI interface. Use with caution in scripts.
    Get:1 http://archive.ubuntu.com/ubuntu bionic InRelease [242 kB]
    Get:2 http://security.ubuntu.com/ubuntu bionic-security InRelease [83.2 kB]
    Get:3 http://security.ubuntu.com/ubuntu bionic-security/multiverse amd64 Packages [1364 B]
    Get:4 http://archive.ubuntu.com/ubuntu bionic-updates InRelease [88.7 kB]
    Get:5 http://security.ubuntu.com/ubuntu bionic-security/main amd64 Packages [264 kB]
    Get:6 http://archive.ubuntu.com/ubuntu bionic-backports InRelease [74.6 kB]
    Get:7 http://archive.ubuntu.com/ubuntu bionic/main amd64 Packages [1344 kB]
    。。。。。。。后面省略

- 生成临时 IMAGE，ID 为 f5d8205bae1b，移除临时容器 607ce2e8553f
- 注：这一步使用的就是 docker commit 命令类似的功能，提交一个运行中的容器生成镜像，只不过该容器是缓存状态，最后会彻底删除


    Removing intermediate container 607ce2e8553f
     ---> f5d8205bae1b

- 基于 f5d8205bae1b 这个临时 IMAGE 启动 ddceac75c0ef 容器，执行 Dockerfile 中第三行的 RUN 指令


    Step 3/3 : RUN apt install -y iproute2
     ---> Running in ddceac75c0ef .......命令执行结果省略

- 处理成功，生成最终 IMAGE，ID 为 9e0ddfd39bb1，并移除临时容器 ddceac75c0ef
- 注：这一步使用的就是 docker commit 命令类似的功能，提交一个运行中的容器生成镜像，只不过该容器是缓存状态，最后会彻底删除


    Processing triggers for libc-bin (2.27-3ubuntu1) ...
    Removing intermediate container ddceac75c0ef
     ---> 9e0ddfd39bb1

- 创建 IMAGE ID 为 9e0ddfd39bb1 的 IMAGE 成功，并成功给 IMAGE 打上标签为 ubuntu-vi:latest


    Successfully built 9e0ddfd39bb1
    Successfully tagged ubuntu-vi:latest

- 下图是创建过程中查询容器，会发现，中间创建的临时 Container 已经都被删除了，并且还会有一个没有名字只有 IMAGE ID 的 IMAGE


    [root@master0 ~]# docker ps -a
    CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS               NAMES
    607ce2e8553f        93fd78260bd1        "/bin/sh -c 'apt ins…"   2 seconds ago       Up 2 seconds                            jovial_almeida
    [root@master0 ~]# docker ps -a
    CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS               NAMES
    ddceac75c0ef        f5d8205bae1b        "/bin/sh -c 'apt ins…"   2 seconds ago       Up 2 seconds                            serene_rosalind
    [root@master0 ~]# docker ps -a
    CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
    [root@master0 ~]# docker ps -a
    [root@master0 ~]# docker images -a
    ubuntu-vi           latest              9e0ddfd39bb1        58 seconds ago      169M
    <none>              <none>              f5d8205bae1b        58 seconds ago      151MB
    ubuntu              latest              93fd78260bd1        7 days ago          86.2MB

- 该过程结束
- 总结
  - dockerfile 的每一个命令，就是给 base image 上新加一层 image

### Dockerfile 构建镜像的过程：

1. 从 base 镜像运行一个容器。
2. 执行一条指令，对容器做修改。
3. 执行类似 docker commit 的操作，生成一个新的镜像层。
4. Docker 再基于刚刚提交的镜像运行一个新容器。
5. 重复 2-4 步，直到 Dockerfile 中的所有指令执行完毕

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ymn9n3/1616121927006-3431594b-1d7d-4764-b6ce-6227bcfc9966.png)

# Dockerfile 格式

Dockerfile 遵循特定的格式和指令集，也就是说，Dockerfile 实际上是 **Instruction(指令)** 的合集，Dockerfile 文件中每行都是一个指令极其参数:

```dockerfile
INSTRUCTION-1 arguments
INSTRUCTION-2 arguments
.......
INSTRUCTION-n arguments
```

实际上指令是不区分大小写的，但是，一般情况，都将他们写成大写，以便在人类阅读时，可以一眼就与参数区分开。

## Dockerfile 指令

详见 [Dockerfile 指令详解](https://www.yuque.com/go/doc/33171452)
