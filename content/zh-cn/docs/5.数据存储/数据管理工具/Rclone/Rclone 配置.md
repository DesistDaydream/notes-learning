---
title: "Rclone 配置"
linkTitle: "Rclone 配置"
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，用法](https://rclone.org/docs/)
> - [官方文档，全局命令行标志](https://rclone.org/flags/)

Rclone 的配置有点混乱，不过大体分为两部分，**Backend(与后端相关)** 和 **Non Backend(与后端无关)** 两种配置。

- **Non Backend** # 通常是指 Rclone 自身的运行时方式。比如使用 sync、copy 等命令时设置并发数、等等。
  - 这部分配置无法通过 rclone.conf 文件设置。
- **Backend** # 针对各种 Remote 所使用的 Backend 的配置。

这两类配置，可以使用**一种或多种**方式进行配置，可用的配置方式有：
- 命令行标志
- 环境变量
- 配置文件

## Backend 无关配置

Rclone 可以通过如下几种方式配置，优先级从上至下由高到低：

- 命令行标志
- 环境变量

rclone 命令的每个选项都可以通过环境变量设置其默认值。选项与环境变量的对应关系规则如下：

- 去掉长选项名称开头的 `--`
- 选项名称中的 `-` 改为 `_`
- 字母改为大写
- 前面加上 `RCLONE_` 前缀

比如：

- --stats 对应 RCLONE_STATS
- --drive-use-trash 对应 RCLONE_DRIVE_USE_TRASH

> 其中 --verbose 的对应关系比较特殊，-v 对应 RCLONE_VERBOSE=1；-vv 对应 RCLONE_VERBOSE=2

## Backend 相关配置

Backend 配置本质上就是配置 Remote，Rclone 通过如下几种方式配置 Backend，优先级从上至下由高到低

- 命令行选项
- 环境变量
- rclone.conf 配置文件

对于 Backend 相关配置来说，配置 rclone.conf 文件，**就像设置默认值**似的，**我们平时在使用 Rclone 时，一般都需要先创建配置文件**，并配置想要管理的 Remote 信息(就是使用 `rclone config` 命令或手动配置 rclone.conf 文件)，在这个文件中配置每个 Remote 的 Backend 信息。有了 Remote 之后，Rclone 就有了控制目标，可以对目标中的数据进行控制了。

比如我们配置了一个名为 test 的 Remtoe，使用 s3 类型的 Backend，然后我们还可以通过 Flags 为 test 这个 Remote 补充其他的配置信息。

```bash
~]# rclone config show
[test]
type = s3
provider = Alibaba
~]#  rclone lsd test: --s3-access-key-id XXXXXXXXXXXX --s3-secret-access-key YYYYYYYYYYYY --s3-endpoint oss-cn-beijing.aliyuncs.com
```

这些 Backend 相关配置的命令行标志，可以在[全局标志-后端相关标志](https://rclone.org/flags/#backend-flags)中找到

命令行标志、环境变量、rclone.conf 配置文件这三种配置的对应关系规则如下：

  - - 去掉长选项名称开头的 `--`
  - 选项名称中的 `-` 改为 `_`
  - 还有一些规则对于配置文件和环境变量不一样
    - 环境变量：
      - 字母改为大写
      - 前面加上 `RCLONE_` 前缀
    - 配置文件：
      - 去掉第一个 `-` 前的字符串
      - 第一个 `-` 之前的字符串作为配置文件中每个部分的名称，即 Backend 名称
      - 保持小写不变
      - 不加前缀

由于 Backend 相关配置都是与各种类型的存储相关联的，所以我们可以查看官方文档中每个存储页面，在其中有一个名为 `Standard options` 的段落，详细描写了该 Backend 的所有可用配置，比如用 S3 类型的 Backend 举例，效果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/rclone/20230506163807.png)

其中 Config 和 Env Var 分别对应 --s3-secret-access-key 这个选项在配置文件和环境变量中应该使用的名称。

## rclone.conf

**rclone.conf** 是 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的配文件，其中保存了各种 Remote 信息。可以直接编辑该文件，也可以使用 `rclone config` 命令在命令行的交互模式下进行配置。

rclone.conf 配置文件的每个字段的用途在官方文档中没有找到说明，但是命令行标志跟配置文件的字段有对应关系，以下是我自己总结的：

- 每个部分的名称就是 Remote 的名称
- 每个 Remote 都有一个 type 字段用来指定该 Remote 使用的 Backend
- 每个 Remote 的其他字段，根据其 Backend 的类型决定，不同的 Backend，有不同的字段。

### 读取逻辑

rclone 运行后，根据如下顺序从上到下依次查找 rclone.conf 文件

- **rclone 二进制可执行文件所在目录**中的 rclone.conf 文件
- **$APPDATA/rclone/rclone.conf**（该位置仅在 Windows 系统中有效）
- **$XDG_CONFIG_HOME/rclone/rclone.conf**
- **~/.config/rclone/rclone.conf**
- **~/.rclone.conf**

如果上述位置都没有找到 rclone.conf 文件，则会在以下位置创建一个新的配置文件：

- Windows 上
  - 在 **$APPDATA/rclone/rclone.conf**
- 类 Unix 上
  - 如果定义了 `$XDG_CONFIG_HOME`，则在 **$XDG_CONFIG_HOME/rclone/rclone.conf**
  - 如果未定义 `$XDG_CONFIG_HOME`，则在 **~/.config/rclone/rclone.conf**

# Backend 无关配置详解

> 参考：
>
> - [官方文档，命令行标志，与后端无关的标志](https://rclone.org/flags/#non-backend-flags)

这部分配置无法通过 rclone.conf 文件配置。只能使用环境变量和命令行选项。除了在这里记录的通用配置意外，各个 Rclone 功能(比如 mount、sync、copy 等命令)还有自身的配置用以配置运行时行为。

- **--cache-dir** # 用于缓存的目录。`默认值：~/.cache/rclone`
- **-n, --dry-run** # 试运行，不会真的执行
- **-i, --interactive** # 开启交互模式
- **-p, --progress** # 显示传输进度、传输速度
- **--transfers INT** # 并行运行的文件传输数。`默认值: 4`
- **-v** # 输出更多的内容，重复只用该选项会增加输出的内容，比如 -vv、-vvv。比如 -vv 就会输出 Debug 信息。

# Backend 相关配置详解

> 参考：
>
> - [官方文档，命令行标志，与后端有关的标志](https://rclone.org/flags/#backend-flags)

**这部分配置本质就是配置 Remote**，为每个 Remote 配置 Backend，可以用的配置内容非常多，并且每种类型的 Backend 通常都有其独立的配置内容。

## S3

https://rclone.org/s3/

## WebDav

https://rclone.org/webdav

# 配置示例

