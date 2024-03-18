---
title: 使用 Tini 清理 Docker 容器僵死进程
---

[使用 Tini 清理 Docker 容器僵死进程](https://mp.weixin.qq.com/s/Ktd56YQsU8pP6kUs3_uU4Q)

> 最近在 Tini 的仓库下看到作者对 Tini 优势的精彩回复，搬运过来，粗糙翻译，献给拥有同样疑惑的你。

## 写在前面

我们在查看一些大项目的 Dockerfile 时经常发现，它们的 ENTRYPOINT 中往往都有 tini 的身影：



Rancher 官方镜像



Jenkins 官方镜像

那 Tini 到底是什么？为什么大家都喜欢在镜像中使用它呢？

## 开发者的疑问

我注意到 Jenkins 的官方镜像中使用了 Tini，所以我很好奇它是什么。它看起来一定很有用，可能解决了一些我不知道的问题。你能用 “说人话” 的方式简单解释一下 Tini 相对于直接以 CMD 运行 shell 脚本的优势吗？

我的几个容器的 ENTRYPOINT 都设置了一个 docker-entrypoint.sh 脚本，里面基本上都是以 “exec "$@"” 的方式在运行，我应该使用 Tini 来代替吗？

## 来自作者的回复

问得好！但这解释可能有点长，所以请耐心听我说（我知道你要求简短，但我真的做不到，捂脸~）。

首先，我们先简单聊聊 Jenkins。当您运行 Docker 容器时，Docker 会将它与系统的其他部分隔离开来。这种隔离发生在不同的级别（例如网络、文件系统、进程）。

但 Tini 并不真正关注网络或文件系统，所以让我们把注意力放在 Tini 的一个重要概念上：进程。

每个 Docker 容器都是一个 PID 命名空间，这意味着容器中的进程与主机上的其他进程是隔离的。PID 命名空间是一棵树，从 PID 1 开始，通常称为 init。

注意：当你运行一个 Docker 容器时，镜像的 ENTRYPOINT 就是你的根进程，即 PID 1（如果你没有 ENTRYPOINT，那么 CMD 就会作为根进程，你可能配置了一个 shell 脚本，或其他的可执行程序，容器的根进程具体是什么，完全取决于你的配置）。

与其他进程不同的是，PID 1 有一个独特的职责，那就是收割 “僵尸进程”。

那何为 “僵尸进程” 呢？

“僵尸进程” 是指：

- 已经退出。
- 没有被其父进程 wait（wait 是指 syscall 父进程用于检索其子进程的退出代码）。
- 父进程已丢失（也就是说，它们的父进程已经不存在了)，这意味着他们永远不会被其父进程处理。

当 “僵尸进程” 被创建时（也就是说，一旦它的父进程非正常退出了，它也就跟着无法正常退出了），它会继承成为 PID 1 的子级，最后 PID 1 会负责关闭它。

换句话说，有人必须在 “不负责任” 的父进程离开后，对这些 “孤儿” 进行清理，这是 PID 1 的作用。

请注意，创建 “僵尸进程” 通常是不被允许的（也就是说，理想情况下，您应该修复代码，这样就不会创建“僵尸进程”），但是对于像 Jenkins 这种应用来说，它们是不可避免的：因为 Jenkins 通常运行的代码不是由 Jenkins 维护者编写的（也就是您的 Jenkins 构建脚本），所以他们也无法“修复代码”。

这就是 Jenkins 使用 Tini 的原因：在构建了创建 “僵尸进程” 的脚本后进行清理。

---

但其实 Bash 实际上也做同样的事情（收割 “僵尸进程”），所以你可能会想：为什么不把 Bash 当作 PID 1 呢？

第一个问题是，如果您将 Bash 作为 PID 1 运行，那么您发送到 Docker 容器的所有信号（例如，使用 docker stop 或 docker kill）最终都会发送到 Bash，Bash 默认不会将它们转发到任何地方（除非您自己编写代码实现）。换句话说，如果你使用 Bash 来运行 Jenkins，那么当你运行 docker stop 的时候，Jenkins 将永远收不到停止信号！

而 Tini 通过 “信号转发” 解决了这个问题：如果你向 Tini 发送信号，那么它也会向你的子进程发送同样的信号（在你的例子中是 Jenkins）。

第二个问题是，一旦您的进程退出，Bash 也会继续退出。如果您不小心，Bash 可能会退出，退出代码为 0，而您的进程实际上崩溃了（但 0 表示 “一切正常”；这将导致 Docker 重启策略不符合您的预期）。因为您真正想要的可能是 Bash 返回与您的进程相同的退出代码。

请注意，您可以通过在 Bash 中创建信号处理程序来实际执行转发，并返回适当的退出代码来解决这个问题。另一方面，这需要做更多的工作，而添加 Tini 只是文档中的几行。

---

其实还有另一个解决方案可以将 Jenkins 作为 PID 1 运行，即在 Jenkins 中添加另一个线程来负责收割 “僵尸进程”。

但这也不理想，原因有二：

首先，如果将 Jenkins 以 PID 1 的身份运行，那么很难区分继承给 Jenkins 的进程（应该被收割）和 Jenkins 自己产生的进程（不应该被收割，因为还有其他代码已经在等待它们执行）。我相信你可以用代码来解决这个问题，但还是要问一遍：当你可以把 Tini 放进去的时候，为什么还要写呢？

其次，如果 Jenkins 以 PID 1 运行，那么它可能不会接收到您发送的信号！

这是 PID 1 进程中的微妙之处。与其他进程不同的是，PID 1 没有默认的信号处理程序，这意味着如果 Jenkins 没有明确地为 SIGTERM 安装信号处理程序，那么该信号在发送时将被丢弃（而默认行为是终止该过程）。

Tini 确实安装了显式信号处理程序（顺便说一下，是为了转发信号），所以这些信号不再被丢弃。相反，它们被发送到 Jenkins，Jenkins 并不像 PID 1（Tini ）那样运行，因此有默认的信号处理程序（注意：这不是 Jenkins 使用 Tini 的原因，Jenkins 使用它来获取信号，但在 RabbitMQ 的镜像中是这个作用）。

---

请注意，Tini 中还有一些额外的功能，在 Bash 或 Java 中很难实现（例如，Tini 可以注册为 “子收割者”，因此它实际上不需要作为 PID 1 运行来完成“僵尸进程” 收割工作），但是这些功能对于一些高级应用场景来说非常有用。

希望以上内容对你有所帮助！

如果您有兴趣了解更多，以下是一些可供参考的资料：

- 僵尸进程详解：<https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/>
- 更简洁的解释：<https://github.com/docker-library/official-images#init>

最后，请注意 Tini 还有更多的选择（比如 Phusion 的基础镜像）。

Tini 的主要特性是：

- 做 PID 1 需要做的一切，而不做其他任何事情。像读取环境文件、改变用户、过程监控等事情不在 Tini 的范围内（还有其他更好的工具）；
- 零配置就能上手（如果运行不正常，Tini >= 0.6 也会警告您）；
- 它有丰富的测试。

## 原文链接

- `What is advantage of Tini?`：<https://github.com/krallin/tini/issues/8#issuecomment-146135930>
- `译文来源`：<https://zhuanlan.zhihu.com/p/59796137>

# Docker 容器中进程管理工具

为了防止容器中直接使用 ENTRYPOINT 或 CMD 指令启动命令或应用程序产生 PID 为 1 的进程无法处理传递信号给子进程或者无法接管孤儿进程，进而导致产生大量的僵尸进程。对于没有能力处理以上两个进程问题的 PID 进程，建议使用 dumb-int 或 tini 这种第三方工具来充当 1 号进程。

Linux 系统中，PID 为 1 的进程需要担任两个重要的使命：

1. **传递信号给子进程**
   如果 pid 为 1 的进程，无法向其子进程传递信号，可能导致容器发送 SIGTERM 信号之后，父进程等待子进程退出。此时，如果父进程不能将信号传递到子进程，则整个容器就将无法正常退出，除非向父进程发送 SIGKILL 信号，使其强行退出，这就会导致一些退出前的操作无法正常执行，例如关闭数据库连接、关闭输入输出流等。
2. **接管孤儿进程，防止出现僵尸进程**
   如果一个进程中 A 运行了一个子进程 B，而这个子进程 B 又创建了一个子进程 C，若子进程 B 非正常退出（通过 SIGKILL 信号，并不会传递 SIGKILL 信号给进程 C），那么子进程 C 就会由进程 A 接管，一般情况下，我们在进程 A 中并不会处理对进程 C 的托管操作（进程 A 不会传递 SIGTERM 和 SIGKILL 信号给进程 C），结果就导致了进程 B 结束了，倒是并没有回收其子进程 C，子进程 C 就变成了僵尸进程。

在 docker 中，`docker stop`命令会发送`SIGTERM`信号给容器的主进程来处理。如果主进程没有处理这个信号，docker 会在等待一段优雅 grace 的时间后，发送`SIGKILL`信号来强制终止

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/635f36ef-00c6-4161-bed8-02cbbf17e561/images)

详情参考：

1、<https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/>

2、<https://medium.com/@gchudnov/trapping-signals-in-docker-containers-7a57fdda7d86>

Github：<https://github.com/Yelp/dumb-init>

dumb-int 是一个用 C 写的轻量级进程管理工具。类似于一个初始化系统，

它充当 PID 1，并立即以子进程的形式允许您的命令，注意在接收到信号时正确处理和转发它们

dumb-init 解决上述两个问题：向子进程代理发送信号和接管子进程。

默认情况下，dumb-init 会向子进程的进程组发送其收到的信号。原因也很简单，前面已经提到过，像 bash 这样的应用，自己接收到信号之后，不会向子进程发送信号。当然，dumb-init 也可以通过设置环境变量`DUMB_INIT_SETSID=0`来控制只向它的直接子进程发送信号。

另外 dumb-init 也会接管失去父进程的进程，确保其能正常退出。

## 安装

- Alpine 镜像的 APK 可以直接安装
  `FROM alpine:3.11.5 RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \ && apk add --no-cache dumb-init ENTRYPOINT ["dumb-init", "--"] CMD ["/usr/local/bin/docker-entrypoint.sh"]`
- 二进制安装
  `RUN version=v1.2.2 && \ wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/$version/dumb-init_$version_amd64 && \ chmod +x /usr/local/bin/dumb-init`
- DEB/RPM 安装
  `RUN version=v1.2.2 && \<br />wget [https://github.com/Yelp/dumb-init/releases/download/](https://github.com/Yelp/dumb-init/releases/download/)`
- pip 安装
  `pip install dumb-init`

Github：<https://github.com/krallin/tini>

## 安装

Alpine 镜像的 APK 可以直接安装

\`FROM alpine:3.11.5
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \<br />&& apk add --no-cache tini

ENTRYPOINT \["tini", "--"]
CMD \["/your/program", "-and", "-its", "arguments"]

\`

## 1、php-fpm 进程的接管

针对 php 应用，通常采用`nginx+php-fpm`的架构来处理请求。为了保证 php-fpm 进程出现意外故障能够自动恢复，通常使用 supervisor 进程管理工具进行守护。php-fpm 的进程管理类也类似于 nginx，由 master，worker 进程组成。master 进程不处理请求，而是由 worker 进程处理！master 进程只负责管理 worker 进程。

master 进程负责监听子进程的状态，子进程挂掉之后，会发信号给 master 进程，然后 master 进程重新启一个新的 worker 进程。

`进程号 父进程号 进程 21 		10 			master 22		21				|----worker1 23		21				|----worker2`

使用 Supervisor 启动、守护 php-fpm 进程时的进程树

\`进程号 父进程号 进程
10    9       supervisor
21 10 |---master
22 21 |----worker1
23 21 |----worker2

# 使用 supervisor 启动、守护的是 php-fpm 的 master 进程，然后 master 进程再根据配置启动对应数量的 worker 进程。

\`

当 php-fpm 的 master 进程意外退出后的进程树

\`进程号 父进程号 进程
10    9       supervisor
22 1  worker1
23 1  worker2

# 此时 worker 进程成为僵尸进程，被 1 号进程接管

\`

此时 supervisor 检测到 php-fpm master 进程不存在就会在重新创建一个新的 php-fpm master 进程。但是会因为原先的 php-fpm worker 没有被杀掉，成为僵尸进程、依旧占用着端口而失败。本以为 php-fpm 会

1. <https://www.infoq.cn/article/2016/01/dumb-init-Docker>
2. <https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/>
3. <https://medium.com/@gchudnov/trapping-signals-in-docker-containers-7a57fdda7d86>

> 原文出处：<https://gitbook.curiouser.top/origin/docker-process-manager.html>
