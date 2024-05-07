---
title: Internationalization
linkTitle: Internationalization
date: 2024-05-07T15:27
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，Internationalization_and_localization](https://en.wikipedia.org/wiki/Internationalization_and_localization)

**Internationalization(国际化，简称 i18n)** 和 **Localization(本地化，简称 L10n)** 是使计算机软件适应不同语言、区域特性和目标语言环境的技术要求的方法。

操作系统中的 `locale` 功能和相关的环境设置是由底层的系统库和支持多语言的软件包共同控制和提供的。在 Linux 系统中，主要是通过 `glibc`（GNU C Library）这一核心库来实现对多语言和地区的支持。`glibc` 提供了对 i18n 和 L10n 的支持，包括对不同语言环境的文本处理、字符编码转换等功能。

Debian 系

- 由 libc-bin 和 systemd 管理

```bash
~]# dpkg-query -S /usr/bin/locale
libc-bin: /usr/bin/locale
~]# dpkg-query -S /usr/bin/localectl 
systemd: /usr/bin/localectl
```

RedHat 系

- 由 glibc-common 和 systemd 管理

```bash
~]# rpm -qf /usr/bin/locale
glibc-common-2.28-36.oe1.x86_64
~]# rpm -qf /usr/bin/localectl 
systemd-243-18.oe1.x86_64
```

此外，还有一些辅助工具和软件包负责生成和管理具体的`locale`数据，比如：

- **localedata**: 这是包含特定于区域信息的数据文件，如日期和时间格式、数字和货币表示方式等，通常随`glibc`一起提供。
- **localedef**: 一个用于从`.po`（gettext翻译项目文件）文件创建或更新`locale`定义的工具，也是`glibc`的一部分。
- **gettext**: 这是一个广泛使用的国际化的库，帮助程序实现多语言支持。它允许程序使用翻译过的字符串，并且维护了一个翻译消息的数据库。

在系统层面，管理`locale`设置的还包括系统配置文件（如 `/etc/default/locale`或`/etc/locale.conf`）以及一些系统管理命令，如`locale-gen`用于生成指定的`locale`，`update-locale`用于更新系统的`locale`设置。

因此，`locale`功能不是由单一的包控制，而是涉及系统库、命令行工具、系统配置文件及一系列相互协作的组件共同作用的结果。

# 关联文件与配置

Debian 系

- **/etc/default/locale**

RedHat 系

- **/etc/locale.conf**

# 命令行工具


## locale


## localectl


# 最佳实践

## 修改 Linux 系统语言

要将Linux系统的语言环境从其他语言（例如中文）修改为英语，你可以按照以下步骤操作：

### 方法 1: 使用 `locale` 命令临时修改

对于临时性的修改，可以在当前终端会话中设置`LANG`环境变量为英语环境。例如，设置为美国英语（UTF-8编码）：

`export LANG=en_US.UTF-8`

这只会改变当前终端会话的语言环境，一旦关闭终端或打开新的终端窗口，设置就会失效。

### 方法 2: 永久修改系统默认语言环境

要永久性地修改系统语言环境，你需要编辑系统配置文件。通常情况下，这个文件是`/etc/default/locale`（针对Debian系，如Ubuntu）或`/etc/locale.conf`（针对Red Hat系，如Fedora、CentOS）。

#### 对于使用 `/etc/default/locale` 的系统：

```bash
tee > /etc/default/locale > /dev/null <<EOF
LANG=en_US.UTF-8
LC_ALL=en_US.UTF-8
EOF
```

为了使改动生效，可以重新登录或者运行以下命令：

```bash
sudo locale-gen en_US.UTF-8
sudo update-locale LANG=en_US.UTF-8
```

#### 对于使用 `/etc/locale.conf` 的系统：

```bash
tee > /etc/locale.conf > /dev/null <<EOF
LANG=en_US.UTF-8
EOF
```

### 注意事项：

- 在某些系统上，你可能还需要使用`localectl`命令来设置语言，例如：
  - `sudo localectl set-locale LANG=en_US.UTF-8`
- 修改系统语言环境后，确保系统已经安装了对应的语言包，如果没有，可能需要使用包管理器（如`apt`或`yum`）来安装。
- 以上步骤适用于大多数Linux发行版，但具体命令和文件路径可能因发行版而异，请根据实际情况调整。