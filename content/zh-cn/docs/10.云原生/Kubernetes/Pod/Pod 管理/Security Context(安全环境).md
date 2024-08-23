---
title: Security Context(安全环境)
---

# 概述

> 参考：
>
> - [官方文档,任务-配置 Pod 和 Containers](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
> - [公众号-阳明](https://mp.weixin.qq.com/s/NFgQrvn_LyU0qQbhMZwDAQ)

**Security Context(安全环境)** 用来定义 Pod 或 Container 的特权与访问控制设置。 安全上下文包括但不限于：

- 自主访问控制（Discretionary Access Control）：基于 [用户 ID（UID）和组 ID（GID）](https://wiki.archlinux.org/index.php/users_and_groups). 来判定对对象（例如文件）的访问权限。
- [安全性增强的 Linux（SELinux）](https://zh.wikipedia.org/wiki/%E5%AE%89%E5%85%A8%E5%A2%9E%E5%BC%BA%E5%BC%8FLinux)： 为对象赋予安全性标签。
- 以特权模式或者非特权模式运行。
- [Linux 权能](https://linux-audit.com/linux-capabilities-hardening-linux-binaries-by-removing-setuid/): 为进程赋予 root 用户的部分特权而非全部特权。
- [AppArmor](https://kubernetes.io/zh/docs/tutorials/clusters/apparmor/)：使用程序框架来限制个别程序的权能。
- [Seccomp](https://en.wikipedia.org/wiki/Seccomp)：过滤进程的系统调用。
- AllowPrivilegeEscalation：控制进程是否可以获得超出其父进程的特权。 此布尔值直接控制是否为容器进程设置 `[no_new_privs](https://www.kernel.org/doc/Documentation/prctl/no_new_privs.txt)`标志。 当容器以特权模式运行或者具有 `CAP_SYS_ADMIN` 权能时，AllowPrivilegeEscalation 总是为 true。
- readOnlyRootFilesystem：以只读方式加载容器的根文件系统。

以上条目不是安全上下文设置的完整列表 -- 请参阅 [SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#securitycontext-v1-core) 了解其完整列表。
关于在 Linux 系统中的安全机制的更多信息，可参阅 [Linux 内核安全性能力概述](https://www.linux.com/learn/overview-linux-kernel-security-features)。

特别注意：限制自由，会产生很多问题，比如：

- 使用 hostPath 类型的 volume 时，如果容器不以 root 用户运行，则无法对 hostPath 所在目录执行操作，任何写操作将会提示权限不够。因为目录权限 755

# 应该了解的 10 个 Kubernetes 安全上下文配置

在 Kubernetes 中安全地运行工作负载是很困难的，有很多配置都可能会影响到整个 Kubernetes API 的安全性，这需要我们有大量的知识积累来正确的实施。Kubernetes 在安全方面提供了一个强大的工具 securityContext，每个 Pod 和容器清单都可以使用这个属性。在本文中我们将了解各种 securityContext 的配置，探讨它们的含义，以及我们应该如何使用它们。
securityContext 设置在 PodSpec 和 ContainerSpec 规范中都有定义，这里我们分别用\[P]和\[C]来表示。需要注意的是，如果一个设置在两个作用域中都可以使用和配置，那么我们应该优先考虑设置容器级别的。

## 1runAsNonRoot \[P/C]

我们知道容器是使用 namespaces 和 cgroups 来限制其进程，但只要在部署的时候做了一次错误的配置，就可以让这些进程访问主机上的资源。如果该进程以 root 身份运行，它对这些资源的访问权限与主机 root 账户是相同的。此外，如果其他 pod 或容器设置被用来减少约束（比如 procMount 或 capabilities），拥有一个 root UID 就会提高风险，除非你有一个非常好的原因，否则你不应该以 root 身份运行一个容器。
那么，如果你有一个使用 root 的镜像需要部署，那应该怎么办呢？

### 1.1 使用基础镜像中提供的用户

通常情况下，基础镜像已经创建并提供了一个用户，例如，官方的 Node.js 镜像带有一个 UID 为 1000 的名为 node 的用户，我们就可以使用该身份来运行容器，但他们并没有在 Dockerfile 中明确地设置当前用户。我们可以在运行时用 runAsUser 设置来配置它，或者用自定义的 Dockerfile 来更改镜像中的当前用户。这里我们来看看使用自定义的 Dockerfile 来构建我们自己的镜像的例子。
在不深入了解镜像构建的情况下，让我们假设我们有一个预先构建好的 npm 应用程序。这里是一个最小的 Dockerfile 文件，用来构建一个基于 node:slim 的镜像，并以提供的 node 用户身份运行。
FROM node:slim
COPY --chown=node . /home/node/app/   *# <--- Copy app into the home directory with right ownership*
USER 1000                             *# <--- Switch active user to “node” (by UID)*
WORKDIR /home/node/app                *# <--- Switch current directory to app*
ENTRYPOINT \["npm", "start"]           *# <--- This will now exec as the “node” user instead of root*

其中以 USER 开头的一行就是关键设置，这使得 node 成为从这个镜像启动的任何容器里面的默认用户。我们使用 UID 而不是用户的名字，因为 Kubernetes 无法在启动容器前将镜像的默认用户名映射到 UID 上，并且在部署时指定 runAsNotRoot: true，会返回有关错误。

### 1.2 基础镜像没有提供用户

如果我们使用的基础镜像没有提供一个可以使用的用户，那么我们又应该怎么做呢？对于大部分进程来说，我们只需在自定义的 Dockerfile 中创建一个用户并使用它即可。如下所示：
FROM node:slim
RUN useradd somebody -u 10001 --create-home --user-group  *# <--- Create a user*
COPY --chown=somebody . /home/somebody/app/
USER 10001
WORKDIR /home/somebody/app
ENTRYPOINT \["npm", "start"]

这里我们增加了一行创建用户的 RUN 命令即可。不过需要注意的是这对于 node.js 和 npm 来说，这很好用，但是其他工具可能需要文件系统的不同元素进行所有权变更。如果遇到任何问题，需要查阅对应工具的文档。

## 2runAsUser/runAsGroup \[P/C]

容器镜像可能有一个特定的用户或组，我们可以用 runAsUser 和 runAsGroup 来进行覆盖。通常，这些设置与包含具有相同所有权 ID 的文件的卷挂载结合在一起。
....
spec:
  containers:
    - name: web
      image: mycorp/webapp:1.2.3
  securityContext:
    runAsNonRoot: true
    runAsUser: 10001
....

不过使用这些配置也是有风险的，因为你为容器做出的运行时决定可能与原始镜像不兼容。例如，jenkins/jenkins 镜像以名为 jenkins:jenkins 的**组:用户**身份运行，其应用文件全部由该用户拥有。如果我们配置一个不同的用户，它将无法启动，因为该用户不存在于镜像的/etc/passwd 文件中。即使它以某种方式存在，它也很可能在读写 jenkins:jenkins 拥有的文件时出现问题。我们可以用一个简单的 docker 运行命令来验证这个问题。
$ docker run --rm -it -u eric:eric jenkins/jenkins
docker: Error response from daemon: unable to find user eric: no matching entries in passwd file.

上面我们提到确保容器进程不以 root 用户身份运行是一个非常好的主意，但不要依赖 runAsUser 或 runAsGroup 设置来保证这一点，未来有人可能会删除这些配置，请确保同时将 runAsNonRoot 设置为 true。

## 3seLinuxOptions \[P/C]

SELinux 是一个用于控制对 Linux 系统上的应用、进程和文件进行访问的策略驱动系统，它在 Linux 内核中实现了 Linux 安全模块框架。SELinux 是基于标签的策略，它将一些标签应用于系统中的所有元素，然后将元素进行分组。这些标签被称为**安全上下文**（不要和 Kubernetes 中的 securityContext 混淆了）- 由用户、角色、类型和可选的一些其他属性组成，格式为：user:role:type:level。
然后，SELinux 使用策略来定义特定上下文中的哪些进程可以访问系统中其他被标记的对象。SELinux 可以是严格执行 enforced 模式，在这种情况下，访问将被拒绝，如果被配置为允许的 permissive 模式，那么安全策略没有被强制执行，当安全策略规则应该拒绝访问时，访问仍然被允许，然而，此时会向日志文件发送一条消息，表示该访问应该被拒绝。在容器中，SELinux 通常给容器进程和容器镜像打上标签，以限制该进程只能访问镜像中的文件。
默认的 SELinux 策略将在实例化容器时由容器运行时应用，securityContext 中的 seLinuxOptions 允许配置自定义的 SELinux 策略标签，请注意，改变容器的 SELinux 策略标签有可能允许容器进程摆脱容器镜像并访问主机文件系统。
当然只有当宿主机操作系统支持 SELinux 时，这个功能才会起作用。

## 4seccompProfile \[P/C]

Seccomp 表示一种安全计算模式，是 Linux 内核的一项功能，它可以限制一个特定进程从用户空间到内核的调用。seccomp 配置文件是使用一个 JSON 文件进行定义的，通常由一组系统调用和发生这些系统调用时的默认动作组成。如下配置所示：
{
    "defaultAction": "SCMP_ACT_ERRNO",
    "architectures": \[
        "SCMP_ARCH_X86_64",
        "SCMP_ARCH_X86",
        "SCMP_ARCH_X32"
    ],
    "syscalls": \[
        {
            "name": "accept",
            "action": "SCMP_ACT_ALLOW",
            "args": \[]
        },
        {
            "name": "accept4",
            "action": "SCMP_ACT_ALLOW",
            "args": \[]
        },
        ...
    ]
}

Kubernetes 通过在 securityContext 中的 seccompProfile 属性来提供一个使用自定义配置文件的机制。
seccompProfile:
  type: Localhost
  localhostProfile: profiles/myprofile.json

这里配置的 type 字段有三个可选的值：

- Localhost：其中 localhostProfile 配置为容器内的 seccomp 配置文件路径。
- Unconfined：其中没有配置文件。
- RuntimeDefault：其中使用容器运行时的默认值--如果没有指定类型，就是默认值。

我们可以在 PodSecurityContext 或 securityContext 中使用这些配置，如果两者都配置了，就会使用容器级别中的配置。
此外与大多数安全相关的设置一样，**最小权限原则**在此同样适用。只给你的容器访问它所需要的权限即可。首先创建一个配置文件，简单地记录哪些系统调用正在发生，然后测试你的应用程序，建立一套允许的系统调用规则。我们可以在 Kubernetes 教程(https://kubernetes.io/docs/tutorials/clusters/seccomp)中找到关于Seccomp的更多信息。

## 5 避免使用特权容器 \[C]

给容器授予特权模式是非常危险的，一般会有一种更简单的方式来实现特定的权限，或者可以通过授予 Linux Capabilities 权限来控制。容器运行时控制器着特权模式的具体实现，但是它会授予容器所有的特权，并解除由 cgroup 控制器执行的限制，它还可以修改 Linux 安全模块的配置，并允许容器内的进程逃离容器。
容器在宿主机中提供了进程隔离，所以即使容器是使用 root 身份运行的，也有容器运行时不授予容器的 Capabilities。如果配置了特权模式，容器运行时就会授予系统 root 的所有能力，从安全角度来看，这是很危险的，因为它允许对底层宿主机系统的所有操作访问。
避免使用特权模式，如果你的容器确实需要额外的能力，只需通过添加 capabilities 来满足你的需求。除非你的容器需要控制主机内核中的系统级设置，如访问特定的硬件或重新配置网络，并且需要访问主机文件系统，那么它就不需要特权模式。

## 6Linux Capabilities \[C]

Capabilities 是一个内核级别的权限，它允许对内核调用权限进行更细粒度的控制，而不是简单地以 root 身份运行。Capabilities 包括更改文件权限、控制网络子系统和执行系统管理等功能。在 securityContext 中，Kubernetes 可以添加或删除 Capabilities，单个 Capabilities 或逗号分隔的列表可以作为一个字符串数组进行配置。另外，我们也可以使用 all 来添加或删除所有的配置。这种配置会被传递给容器运行时，在它创建容器的时候会配置上 Capabilities 集合，如果 securityContext 中没有配置，那么容器将会直接容器运行时提供的所有默认配置。
securityContext:
  capabilities:
    drop:
      - all
    add: \["MKNOD"]

一般推荐的做法是先删除所有的配置，然后只添加你的应用程序实际需要的，在大部分情况下，应用程序在正常运行中实际上不需要任何 Capabilities，通过删除所有配置来测试，并通过监控审计日志来调试问题，看看哪些功能被阻止了。
请注意，当在 securityContext 中列出要放弃或添加的 Capabilities 时，你要删除内核在命名 Capabilities 时使用的 CAP\_前缀。capsh 工具可以给我们一个比较友好的调试信息，可以来说明你的容器中到底启用了哪些 Capabilities，当然不要在生产容器中使用这个工具，因为这使得攻击者很容易弄清楚哪些 Capabilities 被启用了。

## 7 以只读文件系统运行 \[C]

如果你的容器被入侵，而且它有一个可读写的文件系统，那么攻击者就可以随意地改变它的配置、安装软件，并有可能启动其他的漏洞。拥有一个只读的文件系统有助于防止这些类型的安全问题，因为它限制了攻击者可以执行的操作。一般来说，容器不应该要求对容器文件系统进行写入，如果你的应用程序是有状态数据，那么你应该使用外部持久化方法，如数据库、volume 或其他一些服务。另外，确保所有的日志都写到 stdout 或日志转发器上。

## 8procMount \[C]

默认情况下，为了防止潜在的安全问题，容器运行时会屏蔽容器内/proc 文件系统的某些部分文件。然而有时需要访问/proc 的这些文件，特别是在使用嵌套容器时，因为它经常被用作集群内构建过程的一部分。该配置只有两个有效的选项：

- Default：保持标准的容器运行时行为
- Unmasked：它删除/proc 文件系统的所有屏蔽行为

显然只有当我们知道在做什么的时候才应该使用这个配置，如果你是为了构建镜像而使用它，请检查构建工具的最新版本，因为许多工具不再需要这个设置了，最好升级下工具并设置为 Default 默认的 procMount。

## 9fsGroup/fsGroupChangePolicy \[P]

fsGroup 设置定义了一个组，当卷被 pod 挂载时，Kubernetes 将把卷中所有文件的权限改为该组。这里的行为也由 fsGroupChangePolicy 控制，它可以被设置为 onRootMismatch 或 Always。如果设置为 onRootMismatch 则只有当权限与容器 root 的权限不匹配时才会被改变。
不过在使用 fsGroup 时也要慎重，改变整个 volume 卷的组所有权会导致**变慢**，如果是大型文件系统**启动也会延迟**。如果共享同一卷的其他进程没有对新的 GID 的访问权限，它也会对这些进程造成损害。由于这个原因，一些共享文件系统如 NFS，没有实现这个功能。这些设置也不影响临时的 ephemeral 卷。

## 10sysctls \[P]

Sysctls 是 Linux 内核的一个功能，它允许管理员修改内核配置。在一个完整的 Linux 操作系统中，这些是通过使用/etc/sysctl.conf 定义的，也可以使用 sysctl 工具进行修改。
securityContext 中的 sysctls 配置允许在容器中修改特定的 sysctls。只有一小部分的 sysctls 可以在每个容器的基础上进行修改，它们都在内核中被命名的。在这个可以配置的子集中，有些被认为是安全的，而更多的则被认为是不安全的，这取决于对其他 pod 的潜在影响。在集群中，不安全的 sysctls 通常是被禁用，需要由集群管理员专门开启。
鉴于有可能破坏底层操作系统的稳定，除非你有非常特殊的要求，否则应该避免通过 sysctls 修改内核参数。

## 11 总结

在用 securityContext 加固你的应用时，有很多事情需要注意。如果使用得当，它们是一种非常有效的工具，我们希望这个列表能帮助你的团队为你的工作负载和环境进行正确的安全配置。
