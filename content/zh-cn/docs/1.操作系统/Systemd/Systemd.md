---
title: Systemd
linkTitle: Systemd
date: 2023-08-01T12:40
weight: 1
---

# 概述

> 参考：
> 
> - [GitHub 项目，systemd/systemd](https://github.com/systemd/systemd)
> - [官网](https://systemd.io/)
> - [Systemd 中文手册, 金步国](http://www.jinbuguo.com/systemd/systemd.index.html)
> - [Manual(手册)，systemd](https://man.cx/systemd)

**System daemon(系统守护进程，简称 systemd)** 实质上：启动一个服务，就是启动一个程序，可以给该程序添加一些参数，也可以不添加，该程序的可执行文件一般是放在 /usr/lib/systemd/system/ 目录下的

历史上，Linux 的启动一直采用 init 进程。这种命令 `/etc/init.d/apache2 start 或者 service apache2 start`，就是用来启动服务。

这种方法有两个缺点。

1. 启动时间长。init 进程是串行启动，只有前一个进程启动完，才会启动下一个进程。
2. 启动脚本复杂。init 进程只是执行启动脚本，不管其他事情。脚本需要自己处理各种情况，这往往使得脚本变得很长。

Systemd 就是为了解决这些问题而诞生的。它的设计目标是，为系统的启动和管理提供一套完整的解决方案。

根据 Linux 惯例，字母 d 是 daemon(守护进程) 的缩写。 Systemd 这个名字的含义，就是它要守护整个系统。

使用了 Systemd，就不需要再用 init 了。Systemd 取代了 initd，成为系统的第一个进程(PID 等于 1)，其他进程都是它的子进程。

Systemd 的优点是功能强大，使用方便，缺点是体系庞大，非常复杂。事实上，现在还有很多人反对使用 Systemd，理由就是它过于复杂，与操作系统的其他部分强耦合，违反"keep simple, keep stupid"的 Unix 哲学。

注意：Systemd 启动的程序无法获取 shell 中的变量，需要通过在 Unit 的配置文件中设置环境变量。

## Unit(单元)

Systemd 将各种操作系统启动和运行的相关对象，抽象多种类型的 **Units(单元)**，并且提供了 Units 之间的依赖关系。**大多数 Units 是通过 Unit File(单元文件) 创建的**，没有 Unit File，也就不会存在所谓的 Units。**可以这么说，在特定目录创建了一个符合 Unit File 格式的文件，也就创建了一个 Unit**。

> 单元：比如以前上学总说：第一单元、第二单元，这种理解

现阶段有如下几种 Units：

1. **Automoount unit** # 自动挂载点
2. **Device unit** # 硬件设备
3. **Mount unit** # 文件系统挂载点
4. **Path unit** # 文件或路径
5. **Scope unit** # 与 Service unit 类似，但是由 systemd 根据 D-bus 接口接收到的信息自动创建， 可用于管理外部创建的进程。
6. **Service unit** # 用于启动和控制守护进程以及他们所包含的进程
7. **Slice unit** # 用于控制特定 CGroup 内(例如一组 service 与 scope 单元)所有进程的总体资源占用。
8. **Socket nuit** # 进程间通信的 socket
9. **Swap unit** # 关于 swap 文件
10. **Target nuit** # 是一群 Unit 的集合
11. **Timer unit** # 定时器

**Unit 的名称**。Unit 的名称由 Unit File 的名称决定。比如一个 crond.service 文件，将会创建出来一个类型为 Service，名为 crond.service 的 Unit。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/glcff3/1616167393721-79076d3b-2c04-48e9-a4a0-8ee4bfb69284.jpeg)

**Unit 的依赖** 。Systemd 能够处理 Units 之间的依赖关系，通过依赖关系，我们可以确定 Unit 之间启动的先后顺序、以及 Unit 之间是否可以同时运行。

**Unit 的状态**。 Unit 既可以处于活动(active)状态也可以处于停止(inactive)状态， 当然也可以处于启动中(activating)或停止中(deactivating)的状态。 还有一个特殊的失败(failed)状态， 意思是单元以某种方式失败了 (退出时返回了错误代码、进程崩溃、操作超时、触碰启动频率限制)。 当进入失败(failed)状态时， 导致故障的原因 将被记录到日志中以方便日后排查。 需要注意的是， 不同的单元可能还会有各自不同的"子状态"， 但它们都被映射到上述五种状态之一。通过 `systemctl list-units --all` 命令可以查看每个 Unit 的状态。

# Systemd 关联文件与配置

> 参考：
> 
> - [Manual(手册)，systemd-system.conf(5)](<https://man7.org/linux/man-pages/man5/systemd-system.conf.5.html>)

**/etc/systemd/**

- **./system.conf** # Systemd 程序运行时配置文件
- **./system.conf.d/\*.conf** # Systemd 程序运行时配置文件
- **./user.conf** # Systemd 以普通用户身份运行时的配置文件
- **./user.conf.d/\*.conf** # Systemd 以普通用户身份运行时的配置文件

**/run/systemd/** #

- **./system.conf.d/\*.conf** # Systemd 程序运行时配置文件
- **./user.conf.d/\*.conf** # Systemd 以普通用户身份运行时的配置文件

**/usr/lib/systemd/**

- **./system.conf.d/\*.conf** # Systemd 程序运行时配置文件
- **./user.conf.d/\*.conf** # Systemd 以普通用户身份运行时的配置文件

/etc、/run、/usr/lib 这三个目录的优先级从左至右由高到低。Systemd 会从最低优先级的目录 /usr/lib/ 下开始加载配置，注意加载其中的文件，直到最高优先级的目录 /etc/systemd/ 为止。

## Units 配置

Units 配置就是指 Unit File。Systemd 会从多个目录中加载 Unit File，以生成 Unit。下面列出的路径，优先级从上往下越来越低。也就是说，高优先级目录中的文件，将会覆盖低优先级目录中的同名文件。不同的 Systemd 运行方式，加载 Unit File 的路径不同。

### 使用 --system 参数，以系统实例运行 systemd

通过 `pkg-config systemd --variable=systemdsystemunitdir` 命令可以查看包管理器安装完程序后，生成 Unit File 的目录

通过 `pkg-config systemd --variable=systemdsystemconfdir` 命令可以查看优先级最高的存放 Unit File 的目录

- **/etc/systemd/system.control** # 通过 dbus API 创建的永久系统单元
- **/run/systemd/system.control** # 通过 dbus API 创建的临时系统单元
- **/run/systemd/transient** # 动态配置的临时单元(系统与全局用户共用)
- **/run/systemd/generator.early** # 生成的高优先级单元(系统与全局用户共用)(参见 systemd.generator(7) 手册中对 early-dir 的说明)
- **/etc/systemd/system/** # 人类根据需求，手动创建的 Unit File 所在路径。且当使用 systemctl enable UNIT 命令的时候，会自动在该目录中创建软连接到 /usr/lib/systemd/system/ 目录中的 Unit File
  - ./UnitFileName.d/\*.conf # 嵌入式单元文件 存放路径
- **/run/systemd/system/** # 程序运行时自动生成的 Unit File 所在路径。
  - ./UnitFileName.d/\*.conf # 嵌入式单元文件 存放路径
- **/run/systemd/generator** # 生成的中优先级系统单元(参见 systemd.generator(7) 手册中对 normal-dir 的说明)
- **/usr/local/lib/systemd/system** # 本地软件包安装的系统单元
- **/usr/lib/systemd/system/** # 通过系统的包管理器安装程序时，生成的 Unit File 所在路径。
  - ./UnitFileName.d/\*.conf # 嵌入式单元文件 存放路径
- **/run/systemd/generator.late** # 生成的低优先级系统单元(参见 systemd.generator(7) 手册中对 late-dir 的说明)

### 使用 --user 参数，以用户实例运行 systemd

- **$XDG_CONFIG_HOME/systemd/user.control 或 ~/.config/systemd/user.control** # 通过 dbus API 创建的永久私有用户单元(仅在未设置 $XDG_CONFIG_HOME 时才使用 ~/.config 来替代)
- $XDG_RUNTIME_DIR/systemd/user.control # 通过 dbus API 创建的临时私有用户单元
- /run/systemd/transient 动态配置的临时单元(系统与全局用户共用)
- /run/systemd/generator.early 生成的高优先级单元(系统与全局用户共用)(参见 systemd.generator(7) 手册中对 early-dir 的说明)
- $XDG_CONFIG_HOME/systemd/user 或 $HOME/.config/systemd/user 用户配置的私有用户单元(仅在未设置 $XDG_CONFIG_HOME 时才使用 ~/.config 来替代)
- /etc/systemd/user 本地配置的全局用户单元
- $XDG_RUNTIME_DIR/systemd/user 运行时配置的私有用户单元(仅当 $XDG_RUNTIME_DIR 已被设置时有效)
- /run/systemd/user 运行时配置的全局用户单元
- $XDG_RUNTIME_DIR/systemd/generator 生成的中优先级私有用户单元(参见 systemd.generator(7) 手册中对 normal-dir 的说明)
- $XDG_DATA_HOME/systemd/user 或 $HOME/.local/share/systemd/user 软件包安装在用户家目录中的私有用户单元(仅在未设置 $XDG_DATA_HOME 时才使用 ~/.local/share 来替代)
- `$dir/systemd/user`(对应 `$XDG_DATA_DIRS` 中的每一个目录($dir)) 额外安装的全局用户单元，对应 $XDG_DATA_DIRS(默认值="/usr/local/share/:/usr/share/") 中的每一个目录。
- /usr/local/lib/systemd/user 本地软件包安装的全局用户单元
- /usr/lib/systemd/user 发行版软件包安装的全局用户单元
- $XDG_RUNTIME_DIR/systemd/generator.late 生成的低优先级私有用户单元(参见 systemd.generator(7) 手册中对 late-dir 的说明)

可以使用环境变量来 扩充或更改 systemd 用户实例(`--user`)的单元文件加载路径。 环境变量可以通过环境变量生成器(详见 [systemd.environment-generator(7)](http://www.jinbuguo.com/systemd/systemd.environment-generator.html#) 手册)来设置。特别地， `$XDG_DATA_HOME` 与 `$XDG_DATA_DIRS` 可以方便的通过 [systemd-environment-d-generator(8)](http://www.jinbuguo.com/systemd/systemd-environment-d-generator.html#) 来设置。这样，上表中列出的单元目录正好就是默认值。 要查看实际使用的、基于编译选项与当前环境变量的单元目录列表，可以使用 `systemd-analyze --user unit-paths`

此外，还可以通过 [systemctl(1)](http://www.jinbuguo.com/systemd/systemctl.html#) 的 **link** 命令 向上述单元目录中添加额外的单元(不在上述常规单元目录中的单元)。

# 弃用的 System V

> 参考：
> 
> - [GitHub 项目，fedora-sysv/initscripts](https://github.com/fedora-sysv/initscripts)

service、chkconfig、etc. initscripts 相关命令已过时不再推荐使用。

initscripts-service、chkconfig 包也相对的不再推荐使用。