---
title: PAM
linkTitle: PAM
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，linux-pam](https://github.com/linux-pam/linux-pam)
> - [~~官方文档~~](http://www.linux-pam.org/Linux-PAM-html/)~~已停止更新不再维护~~
> - [Manual(手册)，PAM(8)](https://man7.org/linux/man-pages/man8/pam.8.html)
> - [Wiki, PAM](https://en.wikipedia.org/wiki/Pluggable_authentication_module)
> - [Wiki, Linux PAM](https://en.wikipedia.org/wiki/Linux_PAM)
> - [博客园，Linux 下 PAM 模块学习总结](https://www.cnblogs.com/kevingrace/p/8671964.html)
> - [博客园，PAM(Pluggable Authentication Modules)认证机制详情](https://www.cnblogs.com/yinzhengjie/p/8395279.html)
> - <https://www.redhat.com/sysadmin/pluggable-authentication-modules-pam>
> - [金步国，Linux PAM 学习笔记](http://www.jinbuguo.com/linux/pam.html)

**Pluggable Authentication Modules(可插入式认证模块，简称 PAM)** 是由 Sun 提出的一种认证机制。它通过提供一些动态链接库和一套统一的 API，将系统提供的服务和该服务的认证方式分开，使得系统管理员可以灵活地根据需要给不同的服务配置不同的认证方式而无需更改服务程序，同时也便于向系统中添加新的认证手段。

在过去，我们想要对一个使用者进行 Authentication(认证)，得要要求用户输入账号密码，然后通过自行撰写的程序来判断该账号密码是否正确。也因为如此，我们常常使用不同的机制来判断账号密码，所以搞的一部主机上面拥有多个各别的认证系统，也造成账号密码可能不同步的验证问题！为了解决这个问题因此有了 PAM 的机制！

以常见的 su 命令来说，它可以实现用户切换，从 root 切换至其他用户不需要密码、从非 root 用户切换至其他用户则需要验证目标用户的密码，一旦认证成功就以目标用户身份启动 shell 以供使用。本质上，su 只做两件事：(1)认证；(2)启动 shell 。按照传统思路，两件事都很容易实现，例如认证逻辑可以用伪代码这样简单的描述：

```bash
if ( uid == 0 ) 认证成功
elseif ( 输入的密码 == 目标用户的密码 ) 认证成功
else 认证失败
```

但是，认证需求不是一成不变的。例如：

- (1)为了方便运维团队成员(也就是 wheel 组)，希望 wheel 组中的用户无需输入密码也能直接进行用户切换；
- (2)为了加强安全性，希望额外验证手机短信；
- (3)为了避免频繁输入难记的用户密码，希望可以选用指纹方式进行验证；
- (4)为了方便某个特定的用户测试，希望仅凭手机短信也能完成验证
- ......

这样一来， su 的开发者将会被迫不断更改 su 的源代码，然后再重新调试、编译、分发，非常辛苦。这种情况下，PAM 就可以对 su 开发者说："认证的事交给我，能不能通过认证由我说了算，你只需做好其他事情(启动 shell)即可"；同时又对用户(系统管理员)说："只要学会了 PAM 配置语法，就可以利用各种 PAM 模块，编写出千变万化的认证策略。无需打扰 su 开发者，就能立即得到想要的效果"。通过把与认证相关的脏活累活都交给 PAM 来干， su 的开发者与用户之间实现了解耦，彼此皆大欢喜。

推而广之，如果一个应用程序想要使用 PAM 进行认证，只需在源代码中嵌入 PAM 支持即可(也就是引入相应的头文件)。然后开发者无需再为认证部分操心(是否通过认证交给 PAM 决定)，只需专注程序的其他部分即可。也就是说，无需修改应用程序就可以切换、修改、升级应用程序使用的认证机制。当然，如果只有一个不支持 PAM 的二进制文件，那就没有办法改造了。

# Linux-PAM

Linux-PAM 就是一种 Linux 平台上的 PAM 实现

我们可以通过 `ldd` 命令查出来进程是否支持 PAM

```bash
~]# ldd $(which login) | grep pam
 libpam.so.0 => /lib/x86_64-linux-gnu/libpam.so.0 (0x00007f53a3720000)
 libpam_misc.so.0 => /lib/x86_64-linux-gnu/libpam_misc.so.0 (0x00007f53a371b000)
```

Linux-PAM 可以说是一套 API 与 Library，为应用程序提供完整的 **Autherntication(认证)** 机制，只要使用者将认证阶段的需求告知 PAM 后，PAM 就能够回报使用者认证的结果(成功或失败)。由于 PAM 仅是一套认证的机制，又可以提供给其他程序所呼叫引用，因此不论你使用什么程序，都可以使用 PAM 来进行认证，如此一来，就能够让账号密码或者是其他方式的认证具有一致的结果！也让程序设计师方不在着重处理认证的问题。

## Linux-PAM 管理组(认证功能的分组)

Linux-PAM 将认证功能分为 4 个管理组(也可以称为：管理类型)，注意：相同的 PAM 模块关联到不同的管理类型时，产生的效果是不同的。

- **Authentication Management(身份认证管理)** # 认证用户并设置用户凭据。通常，这是通过用户必须满足的一些质询响应请求: 如果你是你声称的那个人，请输入你的密码。
  - 对应配置文件中的 **auth** 关键字
- **Account Management(账户管理)** # 提供账户认证服务类型。比如用户的密码是否已过期;该用户是否允许访问所请求的服务
  - 对应配置文件中的 **account** 关键字
- **Password Management(密码管理)** # 更新身份认证机制。通常，此类服务与 Authentication Management 紧密耦合。一些身份认证机制很适合使用这种功能进行更新。类 Unix 系统基于密码的访问是一个明显的例子
  - 对应配置文件中的 **password** 关键字
- **Session Management(会话管理)** # 会话管理任务涵盖应在提供服务前以及撤回后的服务。这些任务包括维护审计跟踪和用户主目录的安装。会话管理组很重要，因为它提供了用于模块的开放和关闭钩子，以影响用户可用的服务。
  - 对应配置文件中的 **session** 关键字

注意：上面说的"凭据"不仅仅指密码，而是泛指一切认证方式，例如：一次性密码、指纹、短信、IP 地址、二维码。除验证凭据之外，还可以进一步执行更多的关联操作，例如：修改账户所属的组、显示特定的提示信息、赋予账户某些权限。

也可以这么说，PAM 将模块分为 4 大类，每个模块都可以属于一个或多个类型，当指定一个模块的类别的时，当通过 API 调用该模块时，其行为也会变为对应模块类型的行为。

另外，从某种角度看，这 4 个管理类型，也可以称为 4 个 stack(栈)，每个**栈**都是由一组规则组成。**规则栈**的概念在下文有详细描述。

上述 4 种功能，都可以通过位于 /lib/security/ 或 /lib64/security/ 中的各种 PAM 模块实现，模块名称一般都符合 `pam_*.so` 格式

## Linux-PAM 的运行方式

```bash
  +----------------+
  | application: X |
  +----------------+       /  +----------+     +================+
  | authentication-[---->----] Linux-   |--<--| PAM config file|
  |       +        [----<--/--]   PAM    |     |================|
  |[conversation()][--+    \  |          |     | X auth .. a.so |
  +----------------+  |    /  +-n--n-----+     | X auth .. b.so |
  |                |  |       __|  |           |           _____/
  |  service user  |  A      |     |           |____,-----'
  |                |  |      V     A
  +----------------+  +------|-----|---------+ -----+------+
                         +---u-----u----+    |      |      |
                         |   auth....   |--[ a ]--[ b ]--[ c ]
                         +--------------+
                         |   acct....   |--[ b ]--[ d ]
                         +--------------+
                         |   password   |--[ b ]--[ c ]
                         +--------------+
                         |   session    |--[ e ]--[ c ]
```

PAM 大体分为三个部分

- 应用层 # 包括 ftpd、login、passwd 等应用程序
- 接口层 # 连接 应用层 与 身份认证协议层 的接口，包括一些配置文件，通过配置文件可以选择应用程序所采用的认证方案(即所使用的模块及其参数)
- 身份认证协议层 # 由各种模块实现的 Unix、Kerberos、S/Key 等身份认证方案

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ey0dgw/1635432839615-8af1927e-94eb-4119-8444-7dd023925c06.png)

- 应用程序的开发者，通过在代码中导入 PAM API 库，从而实现对认证方法的调用；
- PAM 模块的开发者，则利用 PAM SPI(Service Module API) 来编写模块，以便将不同的认证机制加入到系统中；
- PAM 核心库(libpam) 则会读取配置文件，以此为根据，将应用程序和对应的认证方法联系起来

PAM 通过一个与程序相同文件名的配置文件来进行一连串的认证分析需求(也可以通过程序的配置指定具体的 PAM 配置文件路径)。如上图所示，应用程序 X 调用 Linux-PAM 库，这是一个核心模块，它并不负责认证动作，仅用来暴露 PAM API 以及调用 PAM 模块，真正的认证动作是 PAM 库根据**自定义的 PAM 配置文件** 逐一调用模块进行认证的。PAM 读取配置文件之后会根据读取到的内容进行相应的认证操作，最后会将认证的结果返回给应用程序。

这个 PAM 配置文件，分为 4 个部分，分别是 4 个管理类型，这四个部分都有一个或多个规则来执行认证操作，多个规则组合成一个规则栈。就像上图描述的

- a、b、c 模块进行 auth 认证
- b、d 模块进行 account 认证
- b、c 模块进行 password 认证
- e、c 模块进行 session 认证

一个或多个模块对同一类管理类型进行认证操作，并统一处理结果，以一个整体响应给调用 PAM 的应用程序。所以，总结一下，调用 Linux-PAM 时，其实最多只会返回 4 个结果，因为同一个服务的相同管理类型下的规则，是作为一个整体的。

让我们看看当本地的用户登录到基于文本的控制台时所采取的步骤：

- auth 认证 # 登录应用程序提示输入用户名和密码，然后进行 libpam 身份认证调用以询问：“这个用户是他们所说的吗？” 该 pam_unix 模块负责检查本地账户认证。也可能会检查其他模块，最终将结果传递回登录过程。
- account 认证 # 登录过程接下来会询问“此用户是否允许连接？”，然后对 进行帐户调用 libpam。该 pam_unix 模块会检查诸如密码是否已过期之类的内容。其他模块可能会检查主机或基于时间的访问控制列表。总体响应将返回给流程。
  - 如果密码已过期，应用程序会做出响应。某些应用程序根本无法登录用户。登录过程会提示用户输入新密码。
- password 认证 # 为了验证密码并将其写入正确的位置，登录过程会对 进行密码调用 libpam。该 pam_unix 模块写入到本地 **shadow** 文件。也可以调用其他模块来验证密码强度。
  - 如果此时登录过程仍在继续，则已准备好创建会话。会话调用 libpam 导致 pam_unix 模块将登录时间戳写入**wtmp**文件。其他模块启用 X11 身份认证或 SELinux 用户上下文。
- session 认证 # 在注销时，当会话关闭时，可以对 进行另一个会话调用 libpam。这是 pam_unix 模块将注销时间戳写入**wtmp**文件的时间。

## Rules stack(规则栈)

Linux-PAM 可以通过一组规则栈，对一个程序进行多重验证，假如现在有下面一个服务所使用的 PAM 配置：

```bash
auth   requisite    pam_authtok_get.so.1
auth   required     pam_dhkeys.so.1
auth   required     pam_unix_cred.so.1
auth   required     pam_unix_auth.so.1
auth   required     pam_dial_auth.so.1
```

这就是一组规则栈，用来执行 auth 认证行为，每条规则的执行结果进行集合处理，按照从上到下的顺序逐一执行，每条规则的返回码都会根据 Control 的配置集成到整体结果中。所以，有可能这一组规则只执行了 1 条就停止了，也有可能所有都执行了。

# Linux-PAM 关联文件与配置

**/etc/pam.conf** # PAM 的默认配置文件，当存在 pam.d 文件夹时，自动接管配置，不再读取 pam.conf 下的配置

**/etc/pam.d/** # PAM 的配置文件，当存在该目录时，不再读取 pam.conf 文件中的配置。该目录下的文件，通常都是由其他应用程序在安装时自动创建的。

- **/etc/pam.d/sshd** # 使用 ssh 方式登录时候的配置文件
- **/etc/pam.d/login** # 使用 tty 方式登录时候的配置文件。i.e.通过设备直接登录或者 su 方式切换用户
- **/etc/pam.d/remote** # 使用 telnet 方式登录时候的配置文件
- **/etc/pam.d/kde** # 使用"图形界面"方式登录时候的配置文件
- **/etc/pam.d/system-auth** # 配置凡是调用 system-auth 文件的服务，都会生效
- **/etc/pam.d/common-auth** # 此文件中的安全策略可以限制用户不能更改为之前使用的历史密码
- ...... 等等

`/usr/lib64/security/` # CentOS 发行版的 PAM 模块存放目录

`/usr/lib/x86_64-linux-gnu/` # Ubuntu 发行版的 PAM 模块存放目录

# Linux-PAM 认证机制示例

## system-auth

```bash
~]# grep -v ^# /etc/pam.d/system-auth
auth        required      pam_env.so
auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 500 quiet
auth        required      pam_deny.so

account     required      pam_unix.so
account     sufficient    pam_localuser.so
account     sufficient    pam_succeed_if.so uid < 500 quiet
account     required      pam_permit.so

password    requisite     pam_cracklib.so try_first_pass retry=3 type=
password    sufficient    pam_unix.so sha512 shadow nullok try_first_pass use_authtok
password    required      pam_deny.so

session     optional      pam_keyinit.so revoke
session     required      pam_limits.so
session     [success=1 default=ignore] pam_succeed_if.so service in crond quiet use_uid
session     required      pam_unix.so
```

第一部分表示，当用户登录的时候，首先会通过 auth 类接口对用户身份进行识别和密码认证。所以在该过程中验证会经过几个带 auth 的配置项。

其中的第一步是通过 pam_env.so 模块来定义用户登录之后的环境变量， pam_env.so 允许设置和更改用户登录时候的环境变量，默认情况下，若没有特别指定配置文件，将依据/etc/security/pam_env.conf 进行用户登录之后环境变量的设置。

然后通过 pam_unix.so 模块来提示用户输入密码，并将用户密码与/etc/shadow 中记录的密码信息进行对比，如果密码比对结果正确则允许用户登录，而且该配置项的使用的是“sufficient”控制位，即表示只要该配置项的验证通过，用户即可完全通过认证而不用再去走下面的认证项。不过在特殊情况下，用户允许使用空密码登录系统，例如当将某个用户在/etc/shadow 中的密码字段删除之后，该用户可以只输入用户名直接登录系统。

下面的配置项中，通过 pam_succeed_if.so 对用户的登录条件做一些限制，表示允许 uid 大于 500 的用户在通过密码验证的情况下登录，在 Linux 系统中，一般系统用户的 uid 都在 500 之内，所以该项即表示允许使用 useradd 命令以及默认选项建立的普通用户直接由本地控制台登录系统。

最后通过 pam_deny.so 模块对所有不满足上述任意条件的登录请求直接拒绝，pam_deny.so 是一个特殊的模块，该模块返回值永远为否，类似于大多数安全机制的配置准则，在所有认证规则走完之后，对不匹配任何规则的请求直接拒绝。

第二部分的三个配置项主要表示通过 account 账户类接口来识别账户的合法性以及登录权限。

第一行仍然使用 pam_unix.so 模块来声明用户需要通过密码认证。第二行承认了系统中 uid 小于 500 的系统用户的合法性。之后对所有类型的用户登录请求都开放控制台。

第三部分会通过 password 类接口来确认用户使用的密码或者口令的合法性。第一行配置项表示需要的情况下将调用 pam_cracklib 来验证用户密码复杂度。如果用户输入密码不满足复杂度要求或者密码错，最多将在三次这种错误之后直接返回密码错误的提示，否则期间任何一次正确的密码验证都允许登录。需要指出的是，pam_cracklib.so 是一个常用的控制密码复杂度的 pam 模块，关于其用法举例我们会在之后详细介绍。之后带 pam_unix.so 和 pam_deny.so 的两行配置项的意思与之前类似。都表示需要通过密码认证并对不符合上述任何配置项要求的登录请求直接予以拒绝。不过用户如果执行的操作是单纯的登录，则这部分配置是不起作用的。

第四部分主要将通过 session 会话类接口为用户初始化会话连接。其中几个比较重要的地方包括，使用 pam_keyinit.so 表示当用户登录的时候为其建立相应的密钥环，并在用户登出的时候予以撤销。不过该行配置的控制位使用的是 optional，表示这并非必要条件。之后通过 pam_limits.so 限制用户登录时的会话连接资源，相关 pam_limit.so 配置文件是 /etc/security/limits.conf，默认情况下对每个登录用户都没有限制。关于该模块的配置方法在后面也会详细介绍。
