---
title: Multi-Stage Builds(多阶段构建)
linkTitle: Multi-Stage Builds(多阶段构建)
weight: 20
---

# 概述

> 参考：
>
> - [官方文档](https://docs.docker.com/develop/develop-images/multistage-build/)
>     - https://docs.docker.com/build/building/multi-stage/

# 两个奇技淫巧，将 Docker 镜像体积减小 99%

https://mp.weixin.qq.com/s/6bgtD0Aer6-3u4u9jWBBhw

[^1]:Ephemeral Containers: https://kubernetes.io/docs/concepts/workloads/pods/ephemeral-containers/

[^2]: Dynamic-link library: https://en.wikipedia.org/wiki/Dynamic-link\_library

[^3]: VDSO: https://en.wikipedia.org/wiki/VDSO

对于刚接触容器的人来说，他们很容易被自己构建的 Docker 镜像体积吓到，我只需要一个几 MB 的可执行文件而已，为何镜像的体积会达到 `1 GB` 以上？本文将会介绍几个奇技淫巧来帮助你精简镜像，同时又不牺牲开发人员和运维人员的操作便利性。本系列文章将分为三个部分：

第一部分着重介绍多阶段构建（multi-stage builds），因为这是镜像精简之路至关重要的一环。在这部分内容中，我会解释静态链接和动态链接的区别，它们对镜像带来的影响，以及如何避免那些不好的影响。中间会穿插一部分对 `Alpine` 镜像的介绍。

第二部分将会针对不同的语言来选择适当的精简策略，其中主要讨论 `Go` ，同时也涉及到了 `Java` ， `Node` ， `Python` ， `Ruby` 和 `Rust` 。这一部分也会详细介绍 Alpine 镜像的避坑指南。什么？你不知道 `Alpine` 镜像有哪些坑？我来告诉你。

第三部分将会探讨适用于大多数语言和框架的通用精简策略，例如使用常见的基础镜像、提取可执行文件和减小每一层的体积。同时还会介绍一些更加奇特或激进的工具，例如 `Bazel` ， `Distroless` ， `DockerSlim` 和 `UPX` ，虽然这些工具在某些特定场景下能带来奇效，但大多情况下会起到反作用。

本文介绍第一部分。

## 01 万恶之源

我敢打赌，每一个初次使用自己写好的代码构建 Docker 镜像的人都会被镜像的体积吓到，来看一个例子。

让我们搬出那个屡试不爽的 `hello world` C 程序：

```
/* hello.c */
int main () {
  puts("Hello, world!");
  return 0;
}
```

并通过下面的 Dockerfile 构建镜像：

```
FROM gcc
COPY hello.c .
RUN gcc -o hello hello.c
CMD ["./hello"]
```

然后你会发现构建成功的镜像体积远远超过了 `1 GB` 。。。因为该镜像包含了整个 `gcc` 镜像的内容。

如果使用 `Ubuntu` 镜像，安装 C 编译器，最后编译程序，你会得到一个大概 `300 MB` 大小的镜像，比上面的镜像小多了。但还是不够小，因为编译好的可执行文件还不到 `20 KB` ：

```
$ ls -l hello
-rwxr-xr-x   1 root root 16384 Nov 18 14:36 hello
```

类似地，Go 语言版本的 `hello world` 会得到相同的结果：

```
package main

import "fmt"

func main () {
  fmt.Println("Hello, world!")
}
```

使用基础镜像 `golang` 构建的镜像大小是 `800 MB` ，而编译后的可执行文件只有 `2 MB` 大小：

```
$ ls -l hello
-rwxr-xr-x 1 root root 2008801 Jan 15 16:41 hello
```

还是不太理想，有没有办法大幅度减少镜像的体积呢？往下看。

> 为了更直观地对比不同镜像的大小，所有镜像都使用相同的镜像名，不同的标签。例如： `hello:gcc` ， `hello:ubuntu` ， `hello:thisweirdtrick` 等等，这样就可以直接使用命令 `docker images hello` 列出所有镜像名为 hello 的镜像，不会被其他镜像所干扰。

## 02 多阶段构建

要想大幅度减少镜像的体积，多阶段构建是必不可少的。多阶段构建的想法很简单：“我不想在最终的镜像中包含一堆 C 或 Go 编译器和整个编译工具链，我只要一个编译好的可执行文件！”

多阶段构建可以由多个 `FROM` 指令识别，每一个 `FROM` 语句表示一个新的构建阶段，阶段名称可以用 `AS` 参数指定，例如：

```
FROM gcc AS mybuildstage
COPY hello.c .
RUN gcc -o hello hello.c
FROM ubuntu
COPY --from=mybuildstage hello .
CMD ["./hello"]
```

本例使用基础镜像 `gcc` 来编译程序 `hello.c` ，然后启动一个新的构建阶段，它以 `ubuntu` 作为基础镜像，将可执行文件 `hello` 从上一阶段拷贝到最终的镜像中。最终的镜像大小是 `64 MB` ，比之前的 `1.1 GB` 减少了 `95%` ：

```
🐳 → docker images minimage
REPOSITORY          TAG                    ...         SIZE
minimage            hello-c.gcc            ...         1.14GB
minimage            hello-c.gcc.ubuntu     ...         64.2MB
```

还能不能继续优化？当然能。在继续优化之前，先提醒一下：

在声明构建阶段时，可以不必使用关键词 `AS` ，最终阶段拷贝文件时可以直接使用序号表示之前的构建阶段（从零开始）。也就是说，下面两行是等效的：

```
COPY --from=mybuildstage hello .
COPY --from=0 hello .
```

如果 `Dockerfile` 内容不是很复杂，构建阶段也不是很多，可以直接使用序号表示构建阶段。一旦 Dockerfile 变复杂了，构建阶段增多了，最好还是通过关键词 `AS` 为每个阶段命名，这样也便于后期维护。

### 使用经典的基础镜像

我强烈建议在构建的第一阶段使用经典的基础镜像，这里经典的镜像指的是 `CentOS` ， `Debian` ， `Fedora` 和 `Ubuntu` 之类的镜像。你可能还听说过 Alpine 镜像，不要用它！至少暂时不要用，后面我会告诉你有哪些坑。

### COPY --from 使用绝对路径

从上一个构建阶段拷贝文件时，使用的路径是相对于上一阶段的 **根目录** 的。如果你使用 `golang` 镜像作为构建阶段的基础镜像，就会遇到类似的问题。假设使用下面的 Dockerfile 来构建镜像：

```
FROM golang
COPY hello.go .
RUN go build hello.go
FROM ubuntu
COPY --from=0 hello .
CMD ["./hello"]
```

你会看到这样的报错：

```
COPY failed: stat /var/lib/docker/overlay2/1be...868/merged/hello: no such file or directory
```

这是因为 `COPY` 命令想要拷贝的是 `/hello` ，而 `golang` 镜像的 `WORKDIR` 是 `/go` ，所以可执行文件的真正路径是 `/go/hello` 。

当然你可以使用绝对路径来解决这个问题，但如果后面基础镜像改变了 `WORKDIR` 怎么办？你还得不断地修改绝对路径，所以这个方案还是不太优雅。最好的方法是在第一阶段指定 `WORKDIR` ，在第二阶段使用绝对路径拷贝文件，这样即使基础镜像修改了 `WORKDIR` ，也不会影响到镜像的构建。例如：

```
FROM golang
WORKDIR /src
COPY hello.go .
RUN go build hello.go
FROM ubuntu
COPY --from=0 /src/hello .
CMD ["./hello"]
```

最后的效果还是很惊人的，将镜像的体积直接从 `800 MB` 降低到了 `66 MB` ：

```
🐳 → docker images minimage
REPOSITORY     TAG                              ...    SIZE
minimage       hello-go.golang                  ...    805MB
minimage       hello-go.golang.ubuntu-workdir   ...    66.2MB
```

## 03 FROM scratch 的魔力

回到我们的 `hello world` ，C 语言版本的程序大小为 `16 kB` ，Go 语言版本的程序大小为 `2 MB` ，那么我们到底能不能将镜像缩减到这么小？能否构建一个只包含我需要的程序，没有任何多余文件的镜像？

答案是肯定的，你只需要将多阶段构建的第二阶段的基础镜像改为 `scratch` 就好了。 `scratch` 是一个虚拟镜像，不能被 pull，也不能运行，因为它表示空、nothing！这就意味着新镜像的构建是从零开始，不存在其他的镜像层。例如：

```
FROM golang
COPY hello.go .
RUN go build hello.go
FROM scratch
COPY --from=0 /go/hello .
CMD ["./hello"]
```

这一次构建的镜像大小正好就是 `2 MB` ，堪称完美！

然而，但是，使用 `scratch` 作为基础镜像时会带来很多的不便，且听我一一道来。

### 缺少 shell

`scratch` 镜像的第一个不便是没有 `shell` ，这就意味着 `CMD/RUN` 语句中不能使用字符串，例如：

```
...
FROM scratch
COPY --from=0 /go/hello .
CMD ./hello
```

如果你使用构建好的镜像创建并运行容器，就会遇到下面的报错：

```
docker: Error response from daemon: OCI runtime create failed: container_linux.go:345: starting container process caused "exec: \"/bin/sh\": stat /bin/sh: no such file or directory": unknown.
```

从报错信息可以看出，镜像中并不包含 `/bin/sh` ，所以无法运行程序。这是因为当你在   `CMD/RUN` 语句中使用字符串作为参数时，这些参数会被放到 `/bin/sh` 中执行，也就是说，下面这两条语句是等效的：

```
CMD ./hello
CMD /bin/sh -c "./hello"
```

解决办法其实也很简单：\*\*使用 JSON 语法取代字符串语法。\*\*例如，将 `CMD ./hello` 替换为 `CMD ["./hello"]` ，这样 Docker 就会直接运行程序，不会把它放到 shell 中运行。

### 缺少调试工具

`scratch` 镜像不包含任何调试工具， `ls` 、 `ps` 、 `ping` 这些统统没有，当然了，shell 也没有（上文提过了），你无法使用 `docker exec` 进入容器，也无法查看网络堆栈信息等等。

如果想查看容器中的文件，可以使用 `docker cp` ；如果想查看或调试网络堆栈，可以使用 `docker run --net container:`，或者使用 `nsenter` ；为了更好地调试容器，Kubernetes 也引入了一个新概念叫 Ephemeral Containers [^1] ，但现在还是 Alpha 特性。

虽然有这么多杂七杂八的方法可以帮助我们调试容器，但它们会将事情变得更加复杂，我们追求的是简单，越简单越好。

折中一下可以选择 `busybox` 或 `alpine` 镜像来替代 `scratch` ，虽然它们多了那么几 MB，但从整体来看，这只是牺牲了少量的空间来换取调试的便利性，还是很值得的。

### 缺少 libc

这是最难解决的问题。使用 `scratch` 作为基础镜像时，Go 语言版本的 `hello world` 跑得很欢快，C 语言版本就不行了，或者换个更复杂的 Go 程序也是跑不起来的（例如用到了网络相关的工具包），你会遇到类似于下面的错误：

```
standard_init_linux.go:211: exec user process caused "no such file or directory"
```

从报错信息可以看出缺少文件，但没有告诉我们到底缺少哪些文件，其实这些文件就是程序运行所必需的动态库（dynamic library）。

那么，什么是动态库？为什么需要动态库？

所谓动态库、静态库，指的是程序编译的 **链接阶段** ，链接成可执行文件的方式。 **静态库** 指的是在链接阶段将汇编生成的目标文件.o 与引用到的库一起链接打包到可执行文件中，因此对应的链接方式称为 **静态链接** （static linking）。而动态库在程序编译时并不会被连接到目标代码中，而是在程序运行是才被载入，因此对应的链接方式称为 **动态链接** （dynamic linking）。

90 年代的程序大多使用的是静态链接，因为当时的程序大多数都运行在软盘或者盒式磁带上，而且当时根本不存在标准库。这样程序在运行时与函数库再无瓜葛，移植方便。但对于 Linux 这样的分时系统，会在在同一块硬盘上并发运行多个程序，这些程序基本上都会用到标准的 C 库，这时使用动态链接的优点就体现出来了。使用动态链接时，可执行文件不包含标准库文件，只包含到这些库文件的索引。例如，某程序依赖于库文件 `libtrigonometry.so` 中的 `cos` 和 `sin` 函数，该程序运行时就会根据索引找到并加载 `libtrigonometry.so` ，然后程序就可以调用这个库文件中的函数。

使用动态链接的好处显而易见：

1. 节省磁盘空间，不同的程序可以共享常见的库。
2. 节省内存，共享的库只需从磁盘中加载到内存一次，然后在不同的程序之间共享。
3. 更便于维护，库文件更新后，不需要重新编译使用该库的所有程序。

严格来说，动态库与共享库（shared libraries）相结合才能达到节省内存的功效。Linux 中动态库的扩展名是 `.so` （ `shared object` ），而 Windows 中动态库的扩展名是 `.DLL` （ Dynamic-link library [^2]）。

回到最初的问题，默认情况下，C 程序使用的是动态链接，Go 程序也是。上面的 `hello world` 程序使用了标准库文件 `libc.so.6` ，所以只有镜像中包含该文件，程序才能正常运行。使用 `scratch` 作为基础镜像肯定是不行的，使用 `busybox` 和 `alpine` 也不行，因为 `busybox` 不包含标准库，而 alpine 使用的标准库是 `musl libc` ，与大家常用的标准库 `glibc` 不兼容，后续的文章会详细解读，这里就不赘述了。

那么该如何解决标准库的问题呢？有三种方案。

#### 1、使用静态库

我们可以让编译器使用静态库编译程序，办法有很多，如果使用 gcc 作为编译器，只需加上一个参数 `-static` ：

```
$ gcc -o hello hello.c -static
```

编译完的可执行文件大小为 `760 kB` ，相比于之前的 `16kB` 是大了好多，这是因为可执行文件中包含了其运行所需要的库文件。编译完的程序就可以跑在 `scratch` 镜像中了。

如果使用 alpine 镜像作为基础镜像来编译，得到的可执行文件会更小（< 100kB），下篇文章会详述。

#### 2、拷贝库文件到镜像中

为了找出程序运行需要哪些库文件，可以使用 `ldd` 工具：

```
$ ldd hello
    linux-vdso.so.1 (0x00007ffdf8acb000)
    libc.so.6 => /usr/lib/libc.so.6 (0x00007ff897ef6000)
    /lib64/ld-linux-x86-64.so.2 => /usr/lib64/ld-linux-x86-64.so.2 (0x00007ff8980f7000)
```

从输出结果可知，该程序只需要 `libc.so.6` 这一个库文件。 `linux-vdso.so.1` 与一种叫做   VDSO [^3] 的机制有关，用来加速某些系统调用，可有可无。 `ld-linux-x86-64.so.2` 表示动态链接器本身，包含了所有依赖的库文件的信息。

你可以选择将 `ldd` 列出的所有库文件拷贝到镜像中，但这会很难维护，特别是当程序有大量依赖库时。对于 `hello world` 程序来说，拷贝库文件完全没有问题，但对于更复杂的程序（例如使用到 DNS 的程序），就会遇到令人费解的问题： `glibc` （GNU C library）通过一种相当复杂的机制来实现 DNS，这种机制叫 `NSS` （Name Service Switch, 名称服务开关）。它需要一个配置文件 `/etc/nsswitch.conf` 和额外的函数库，但使用 `ldd` 时不会显示这些函数库，因为这些库在程序运行后才会加载。如果想让 DNS 解析正确工作，必须要拷贝这些额外的库文件（ `/lib64/libnss_*` ）。

我个人不建议直接拷贝库文件，因为它非常难以维护，后期需要不断地更改，而且还有很多未知的隐患。

#### 3、使用 busybox:glibc 作为基础镜像

有一个镜像可以完美解决所有的这些问题，那就是 `busybox:glibc` 。它只有 `5 MB` 大小，并且包含了 `glibc` 和各种调试工具。如果你想选择一个合适的镜像来运行使用动态链接的程序， `busybox:glibc` 是最好的选择。

注意：如果你的程序使用到了除标准库之外的库，仍然需要将这些库文件拷贝到镜像中。

## 04 总结

最后来对比一下不同构建方法构建的镜像大小：

- 原始的构建方法：1.14 GB
- 使用 `ubuntu` 镜像的多阶段构建：64.2 MB
- 使用 `alpine` 镜像和静态 `glibc` ：6.5 MB
- 使用 `alpine` 镜像和动态库：5.6 MB
- 使用 `scratch` 镜像和静态 `glibc` ：940 kB
- 使用 `scratch` 镜像和静态 `musl libc` ：94 kB

最终我们将镜像的体积减少了 `99.99%` 。

但我不建议使用 sratch 作为基础镜像，因为调试起来非常麻烦，但如果你喜欢，那我也不会拦着你。

下篇文章将会着重介绍 Go 语言的镜像精简策略，其中会花很大的篇幅来讨论 alpine 镜像，因为它实在是太酷了，在使用它之前必须得摸清它的底细。

