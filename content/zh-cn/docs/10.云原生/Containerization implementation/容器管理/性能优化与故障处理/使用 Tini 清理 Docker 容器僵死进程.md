---
title: 使用 Tini 清理 Docker 容器僵死进程
---

[使用 Tini 清理 Docker 容器僵死进程](https://mp.weixin.qq.com/s/Ktd56YQsU8pP6kUs3_uU4Q)

> 最近在Tini的仓库下看到作者对Tini优势的精彩回复，搬运过来，粗糙翻译，献给拥有同样疑惑的你。

## 写在前面

我们在查看一些大项目的Dockerfile时经常发现，它们的ENTRYPOINT中往往都有tini的身影：

![图片](https://mmbiz.qpic.cn/mmbiz_jpg/z9BgVMEm7YuABibMzousFl4z8FjppPLpJbBrOH8PLtyGDoDLbAzSzvlOfpNEnopz1dlwiav4N1zQKTIibJKrVP1TA/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1)

Rancher官方镜像

![图片](https://mmbiz.qpic.cn/mmbiz_jpg/z9BgVMEm7YuABibMzousFl4z8FjppPLpJUBBeEsMjKFZu4WJdxxHN6iar81Fx4NotscJh3VkyKPUg5N6ia0Eu2qJw/640?wx_fmt=jpeg&tp=webp&wxfrom=5&wx_lazy=1)

Jenkins官方镜像

那Tini到底是什么？为什么大家都喜欢在镜像中使用它呢？

## 开发者的疑问

我注意到Jenkins的官方镜像中使用了Tini，所以我很好奇它是什么。它看起来一定很有用，可能解决了一些我不知道的问题。你能用“说人话”的方式简单解释一下Tini相对于直接以CMD运行shell脚本的优势吗？

我的几个容器的ENTRYPOINT都设置了一个 docker-entrypoint.sh 脚本，里面基本上都是以“exec "$@"”的方式在运行，我应该使用Tini来代替吗？

## 来自作者的回复

问得好！但这解释可能有点长，所以请耐心听我说（我知道你要求简短，但我真的做不到，捂脸~）。

首先，我们先简单聊聊Jenkins。当您运行Docker容器时，Docker会将它与系统的其他部分隔离开来。这种隔离发生在不同的级别（例如网络、文件系统、进程）。

但Tini并不真正关注网络或文件系统，所以让我们把注意力放在Tini的一个重要概念上：进程。

每个Docker容器都是一个PID命名空间，这意味着容器中的进程与主机上的其他进程是隔离的。PID命名空间是一棵树，从PID 1开始，通常称为init。

注意：当你运行一个Docker容器时，镜像的ENTRYPOINT就是你的根进程，即PID 1（如果你没有ENTRYPOINT，那么CMD就会作为根进程，你可能配置了一个shell脚本，或其他的可执行程序，容器的根进程具体是什么，完全取决于你的配置）。

与其他进程不同的是，PID 1有一个独特的职责，那就是收割“僵尸进程”。

那何为“僵尸进程”呢？

“僵尸进程”是指：

- 已经退出。
- 没有被其父进程wait（wait是指syscall父进程用于检索其子进程的退出代码）。
- 父进程已丢失（也就是说，它们的父进程已经不存在了)，这意味着他们永远不会被其父进程处理。

当“僵尸进程”被创建时（也就是说，一旦它的父进程非正常退出了，它也就跟着无法正常退出了），它会继承成为PID 1的子级，最后PID 1会负责关闭它。

换句话说，有人必须在“不负责任”的父进程离开后，对这些“孤儿”进行清理，这是PID 1的作用。

请注意，创建“僵尸进程”通常是不被允许的（也就是说，理想情况下，您应该修复代码，这样就不会创建“僵尸进程”），但是对于像Jenkins这种应用来说，它们是不可避免的：因为Jenkins通常运行的代码不是由Jenkins维护者编写的（也就是您的Jenkins构建脚本），所以他们也无法“修复代码”。

这就是Jenkins使用Tini的原因：在构建了创建“僵尸进程”的脚本后进行清理。

---

但其实Bash实际上也做同样的事情（收割“僵尸进程”），所以你可能会想：为什么不把Bash当作PID 1呢？

第一个问题是，如果您将Bash作为PID 1运行，那么您发送到Docker容器的所有信号（例如，使用docker stop或docker kill）最终都会发送到Bash，Bash默认不会将它们转发到任何地方（除非您自己编写代码实现）。换句话说，如果你使用Bash来运行Jenkins，那么当你运行docker stop的时候，Jenkins将永远收不到停止信号！

而Tini通过“信号转发”解决了这个问题：如果你向Tini发送信号，那么它也会向你的子进程发送同样的信号（在你的例子中是Jenkins）。

第二个问题是，一旦您的进程退出，Bash也会继续退出。如果您不小心，Bash可能会退出，退出代码为0，而您的进程实际上崩溃了（但0表示“一切正常”；这将导致Docker重启策略不符合您的预期）。因为您真正想要的可能是Bash返回与您的进程相同的退出代码。

请注意，您可以通过在Bash中创建信号处理程序来实际执行转发，并返回适当的退出代码来解决这个问题。另一方面，这需要做更多的工作，而添加Tini只是文档中的几行。

---

其实还有另一个解决方案可以将Jenkins作为PID 1运行，即在Jenkins中添加另一个线程来负责收割“僵尸进程”。

但这也不理想，原因有二：

首先，如果将Jenkins以PID 1的身份运行，那么很难区分继承给Jenkins的进程（应该被收割）和Jenkins自己产生的进程（不应该被收割，因为还有其他代码已经在等待它们执行）。我相信你可以用代码来解决这个问题，但还是要问一遍：当你可以把Tini放进去的时候，为什么还要写呢？

其次，如果Jenkins以PID 1运行，那么它可能不会接收到您发送的信号！

这是PID 1进程中的微妙之处。与其他进程不同的是，PID 1没有默认的信号处理程序，这意味着如果Jenkins没有明确地为SIGTERM安装信号处理程序，那么该信号在发送时将被丢弃（而默认行为是终止该过程）。

Tini确实安装了显式信号处理程序（顺便说一下，是为了转发信号），所以这些信号不再被丢弃。相反，它们被发送到Jenkins，Jenkins并不像PID 1（Tini ）那样运行，因此有默认的信号处理程序（注意：这不是Jenkins使用Tini的原因，Jenkins使用它来获取信号，但在RabbitMQ的镜像中是这个作用）。

---

请注意，Tini中还有一些额外的功能，在Bash或Java中很难实现（例如，Tini可以注册为“子收割者”，因此它实际上不需要作为PID 1运行来完成“僵尸进程”收割工作），但是这些功能对于一些高级应用场景来说非常有用。

希望以上内容对你有所帮助！

如果您有兴趣了解更多，以下是一些可供参考的资料：

- 僵尸进程详解：https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/
- 更简洁的解释：https://github.com/docker-library/official-images#init

最后，请注意Tini还有更多的选择（比如Phusion的基础镜像）。

Tini的主要特性是：

- 做PID 1需要做的一切，而不做其他任何事情。像读取环境文件、改变用户、过程监控等事情不在Tini的范围内（还有其他更好的工具）；
- 零配置就能上手（如果运行不正常，Tini >= 0.6也会警告您）；
- 它有丰富的测试。

## 原文链接

- `What is advantage of Tini?`：https://github.com/krallin/tini/issues/8#issuecomment-146135930
- `译文来源` ：https://zhuanlan.zhihu.com/p/59796137
