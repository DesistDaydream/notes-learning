---
title: Dockerfile 指令详解
---

# 概述

> 参考：
>
> - [官方文档](https://docs.docker.com/engine/reference/builder/)
> - [官方文档-构建镜像的最佳实践](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)

# FROM # 指定 base 镜像

`FROM` 指令将会初始化一个新的构建阶段，并设置用于后续指令所使用的 BaseImage(基础镜像)。**所以一个有效的 Dockerfile 必须从 FROM 指令开始**。`ARG` 是唯一一个可以在 `FROM` 指令前面的指令，除此以外，`FROM` 可以说是**必须存在的基础字段且为 DokerFile 文件的第一个非注释行。**

一个 Dockerfile 中可以有多个 `FROM` 指令，每出现一个 `FROM` 指令，即表示一个老阶段的结束，一个新阶段的开始。

## Syntax(语法)

**from \[--platform=\<PLATFORM>] \<image>\[:\<TAG> | @\[DIGEST] ] \[AS \<NAME>]**
附加指令：

1. **AS \<NAME>**# 为当前构建阶段起一个名字。该附加指令有如下几种用法：
   1. 在开始构建之前，可是使用 --target 指令指定指定要从 STRING 这个阶段开始构建镜像。
   2. 在构建中，`COPY` 指令可以使用 --from=\<NAME> 参数来指定数据源是来自某个构建阶段内的数据，也就是说，`COPY` 指令不止可以从宿主机拷贝文件到容器中，还可以从上一个构建阶段的容器中，拷贝其内容到当前容器中。这也为多阶段构建模式中，减少镜像体积打下来坚实基础。
2. **TAG 和 DIGEST** # 该附加指令是可选的，若不指定镜像的 TAG，则默认使用 latest。

## 用法

Docker 还存在一个特殊的镜像，名为 scratch。这个镜像是虚拟的概念，并不实际存在，它表示一个空白的镜像。

1. 如果以 scratch 为基础镜像的话，意味着本次构建阶段不以任何镜像为基础，接下来所写的指令将作为镜像第一层开始存在。
2. 不以任何系统为基础，直接将可执行文件复制进镜像的做法并不罕见，比如 coreos/etcd。对于 Linux 下静态编译的程序来说，并不需要有操作系统提供运行时支持，所需的一切库都已经在可执行文件里了，因此直接 FROM scratch 会让镜像体积更加小巧。使用 Go 语言 开发的应用很多会使用这种方式来制作镜像，这也是为什么有人认为 Go 是特别适合容器微服务架构的语言的原因之一。

# LABEL # 为镜像添加标签

## Syntax(语法)

`label <key>=<value> <key>=<value> <key>=<value> ...`

## 用法

代替曾经的 MAINTANER 指令。可以通过这种方式来添加维护者信息：`LABEL maintainer="SvenDowideit@home.org.au"`

# ENV # 设置环境变量，环境变量可被后面的指令使用

调用格式为：`$VariableName` 或 `${VariableName}`
定义格式 # Key 是变量名，Value 是变量的值，这是一个键值对的格式

- ENV Key Value # Key 之后的所有内容均被视为 Value 的一部分(包括各种特殊符号和空格等)，因此，一次只能定义一个变量
- ENV Key=Value Value ... # 一次给变量定义多个值，每个 Value 以空格分割，如果 Value 值中有空格，需要加\进行转义或者给 Value 加引号；另外反斜线也可用于续行

注意：在 run 的时候如果指定了变量变量的值，则会顶替调做 Image 时候用 ENV 指定的变量的值

# WORKDIR # 设置当前构建阶段的工作目录

为该命令后面的 RUN, CMD, ENTRYPOINT, ADD 或 COPY 指令设置镜像中的当前工作目录。

# RUN # 在构建 Image 时运行指定的命令

## Syntax(语法)

**run \<COMMAND>**
COMMAND 通常是一个 shell 命令，且 Docker 会以 `/bin/sh -c` 来运行这个命令，这意味着此进程在容器中的 PID 不为 1，不能接收 Unix 信号，因此当使用 docker stop 命令停止 Container 时，此进程接收不到 SIGTERM 信号

**run \["Executable","Param1","Param2",.....]**
Executable 是可执行的命令，参数是一个 JSON 格式的数组，不过这种格式指定的命令不会以“/bin/sh -c”来发起，因此常见的 shell 操作(如通配符，管道符等等)不会进行；如果要运行的命令想用 shell 特性，则可以写成如下格式 RUN \["/bin/sh","-c","Executable","Param1","Param2",.....]

## 用法

`RUN` 指令通常用来在本次构建阶段运行系统命令，以安装某些包或配置某些文件，比如使用 yum、apt、apk 等包管理工具安装，执行 go build 等命令构建代码，等等等等。

# COPY # 从指定的文件拷贝到镜像中

## Syntax(语法)

SRC 指源文件，即需要复制的源文件或目录，支持用通配符；DEST 指目标路径，即即将创建的 IMAGE 中的系统路径(若不适用绝对路径，则默认使用 WORKDIR 指令中指定的目录为起始路径)

1. copy \[] SRC1 SRC2 DEST

注意：src 的来源可以有两个地方

1. 只能指定 build context 中的文件或目录；如果指定了多个 SRC 或在 SRC 使用了通配符，则 DEST 必须是一个目录且以/结尾
2. 为 `COPY` 指令添加 --from 参数，可以让 src 从指定的构建阶段中获取源文件或目录。

## 用法

# ADD # 与 COPY 类似

与 COPT 指令的区别

1. 如果 SRC 为 URL 且 DEST 不以/结尾，则 SRC 指定的文件将被下载并直接被创建为 DEST；如果 DEST 以/结尾，则文件名 URL 指定的文件将被直接下载并保存为 DEST/FileName
2. 如果 SRC 是一个 tar 文件，则该文件会被展开为一个目录，类似于"tar -x"命令但是通过 URL 获取的 tar 文件不会自动展开.
3. 如果 SRC 有有多个，或使用了通配符，则 DEST 必须是一个以/结尾的目录路径，如果 DEST 不以/结尾，则被视作一个普通文件，SRC 的内容将被直接写入到 DEST

# VOLUME # 用于在 Image 中创建一个挂载点目录

以便挂载 Docker host 上的 Volume 或者其他 Container 上的 Volume
通过 VOLUME 命令创建的 Image 在启动成 Container 后，会在 host 上生成一个目录，以便让 Container 中 VOLUME 定义的目录与 host 目录关联

VOLUME MountPoint

# USER # 用于指定运行 Image 时的或运行 Dockerfile 中任何 RUN、CMD、ENTRYPOINT 指令的程序时的用户名或 UID

默认情况使用 root，如果想指定特殊用户，则在/etc/passwd 文件中有该用户才可以

# HEALTHCHECK # 健康检查

# EXPOSE # 指定容器中的进程会监听某个端口，Docker 可以将该端口暴露出来

指定完成后就算运行成 Container 也不会暴露端口，需要在 run 的时候指定-P 选项

# ENTRYPOINT # 用于为容器指定默认运行程序，从而使得容器像是一个单独的可执行程序

Dockerfile 中可以有多个 ENTRYPOINT 指令，但只有最后一个生效。

注意：**ENTRYPOINT 不会被 docker run 之后的参数替换**，CMD 指令的内容 或 docker run 命令最后手动添加的参数会被当做 ENTRYPOINT 指令设定的命令的参数。

## Syntax(语法)

**ENTRYPOINT \["executable", "param-1", "param-2",..."param-n"]**

注意：语法中 `[]` 符号不代表其内的内容是可选的，而是表示 `[]` 这个符号是语法中的一部分。

- executable # 将要执行的具体二进制程序名
- param-x # executable 的命令行参数

## 用法

`ENTRYPOINT` 指令的最佳用途是设置镜像的主命令，就好像这个镜像在运行时就是在执行这条命令似的，然后使用 `CMD` 指令为该命令添加标志。
假如现在有这么一个 Dockerfile，其中一部分是这样的：

```dockerfile
ENTRYPOINT ["s3cmd"]
CMD ["--help"]
```

这两条指令，其实就相当于这个镜像在运行时，执行了这么一条命令：`s3cmd --help`。

举个例子吧：

```bash
# 这是用 CMD 的情况
[root@ansible exporter]# docker run -p=8081:8081 lchdzh/xsky-exporter:v0.1 --web.listen-address=":8081"
docker: Error response from daemon: OCI runtime create failed: container_linux.go:349: starting container process caused "exec: \"--web.listen-address=:8081\": executable file not found in $PATH": unknown.
# 这是用 ENTRYPOINT 的情况
[root@ansible exporter]# docker run -p=8081:8081 lchdzh/xsky-exporter:v0.1 --web.listen-address=":8081"
time="2020-12-30 05:41:52" level=info msg="Scraper enabled cluster_info"
time="2020-12-30 05:41:52" level=info msg="Listening on address :8081"
```

ENTRYPOINT 与 CMD 的比较

1. ENTRYPOINT 一般用于设置 Container 启动后的第一个命令，这对一个 Container 来说是固定的
2. CMD 一般用于设置 Container 启动的第一个命令的默认参数，这对一个 Container 来说是可以变化的
3. 这俩个设置用于设定云原生应用的配置文件具体思路如下,现在以 Nginx 为例
   1. 使用一个写入配置文件的脚本做 ENTRYPOINT 的指令
   2. 使用软件的运行命令作为 CMD 的指令
   3. 在 a 中的脚本最后加上 exec $@以引用命令的所有参数，这样当 CMD 作为 ENTRYPOINT 的参数时，ENTRYPOINT 的指令执行完成之后会自动退出 shell 并运行 CMD 的指令，且 CMD 的命令也成为了该 Container 的 PID 为 1 的进程，并且注意让 nginx 在前台运行以保证容器的长时间运行。这样的话 ENTRYPOINT 的配置文件也带入进去了。然后只需要在 run 的时候指明变量的值，就可以实现不同环境使用不同配置的功能

# CMD # 启动 Container 时运行指定的命令

Dockerfile 中可以有多个 CMD 指令，但只有最后一个生效。

> 注意：CMD 会被 docker run 命令最后的参数替换掉。

CMD Command # 与 RUN 相同

CMD \["Executable","Param1","Param2",.....] # 与 RUN 相同

CMD \["Param1","Param2",.....] # 用于为 ENTRYPOINT 指令提供默认参数

# 其他

## ONBUILD

## STOPSIGNAL

## SHELL
