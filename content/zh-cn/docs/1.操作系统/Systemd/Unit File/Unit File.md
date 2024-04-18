---
title: Unit File
---

# 概述

> 参考：
>
> - [Manual(手册)，systemd.unit(5)](https://man7.org/linux/man-pages/man5/systemd.unit.5.html) # Unit 的介绍
> - [Manual(手册)，systemd.syntax(7)](https://man7.org/linux/man-pages/man7/systemd.syntax.7.html) # Unit 的配置语法
> - [金步国 systemd.unit 中文手册](http://www.jinbuguo.com/systemd/systemd.unit.html#)

**Unit File**，是 [INI](/docs/2.编程/无法分类的语言/INI.md) 格式的纯文本文件。在这个文件中，由 **Directives(指令)** 和 **Sections(部分)** 组成，这里的 Directve 就是 INI 格式中的 `键/值对`。

- **Directives(指令)** # 指令由 名称 与 值 组成，以 `=` 分割
- **Sections(部分)** # 与 INI 中的 Sections 概念一样。是一组 Directives 的集合

Unit File 与 INI 格式文件不同的地方是关于注释，Unit File 使用 `#` 作为注释行的开头。

## Unit File 最简单示例

```bash
[Unit]
# Unit 的描述
Description=Foo

[Service]
# 如何启动该 Unit
ExecStart=/usr/sbin/foo-daemon

[Install]
# 当 enable 该 Unit 时，应该在 multi-user.tartet.wants/ 目录中创建软链接
WantedBy=multi-user.target
```

# Unit File 关联文件与配置

systemd 处理 Unit 时，默认从 3 个路径注意查找要读取的 Unit File。

**/etc/systemd/system/** # Unit File 的存放路径，具有最高优先级

- **./UnitName.d/\*.conf** # Unit File 的 include 功能，该路径下的的以 .conf 结尾的文件，将会附加到主 Unit File 中
- .**/UnitName.wants/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Wants 指令的值。
- .**/UnitName.requires/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Requires 指令的值。

**/run/systemd/system/** # Unit File 的存放路径，具有中等优先级

- **./UnitName.d/\*.conf** # Unit File 的 include 功能，该路径下的的以 .conf 结尾的文件，将会附加到主 Unit File 中
- .**/UnitName.wants/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Wants 指令的值。
- .**/UnitName.requires/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Requires 指令的值。

**/usr/lib/systemd/system/** # Unit File 的存放路径，具有最低优先级

- **./UnitName.d/\*.conf** # Unit File 的 include 功能，该路径下的的以 .conf 结尾的文件，将会附加到主 Unit File 中
- .**/UnitName.wants/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Wants 指令的值。
- .**/UnitName.requires/** # 与 include 功能类似，区别在于，该路径下的文件都是其他 Unit File 的软链接，这些 Unit File 的文件名将会作为主 Unit File 配置中 \[Unit] 部分中 Requires 指令的值。

Systemd 会从最低优先级的目录 `/usr/lib/`下开始加载配置，注意加载其中的文件，直到最高优先级的目录 /etc/systemd/_ 为止。

# Unit File 规范

## Unit File 名称

一个有效的 Unit File 名称由三部分组成

- **NAME** # Unit 的名称。
- **DOT** # `.` 符号。
- **TYPE**# Unit 的类型。
  - TYPE 必须是 ".service", ".socket", ".device", ".mount", ".automount", ".swap", ".target", ".path", ".timer", ".slice", or ".scope" 中的一个。

比如 `foo.server` 就是一个有效的 Unit File 名称。

除了手册中列出的选项之外，单元文件还可以包含更多其他选项。 无法识别的选项不会中断单元文件的加载，但是 systemd 会输出一条警告日志。 如果选项或者小节的名字以 `X-` 开头， 那么 systemd 将会完全忽略它。 以 `X-` 开头的小节中的选项没必要再以 `X-` 开头， 因为整个小节都已经被忽略。 应用程序可以利用这个特性在单元文件中包含额外的信息。

如果想要给一个单元赋予别名，那么可以按照需求，在系统单元目录或用户单元目录中， 创建一个软连接(以别名作为文件名)，并将其指向该单元的单元文件。 例如 `systemd-networkd.service` 在安装时就通过 `/usr/lib/systemd/system/dbus-org.freedesktop.network1.service` 软连接创建了 `dbus-org.freedesktop.network1.service` 别名。 此外，还可以直接在单元文件的 \[Install] 部分中使用 `Alias=` 创建别名。 注意，单元文件中设置的别名会随着单元的启用(enable)与禁用(disable)而生效和失效， 也就是别名软连接会随着单元的启用(enable)与禁用(disable)而创建与删除。 例如，因为 `reboot.target` 单元文件中含有 `Alias=ctrl-alt-del.target` 的设置，所以启用(enable)此单元之后，按下 CTRL+ALT+DEL 组合键将会导致启动该单元。单元的别名可以用于 **enable**, **disable**, **start**, **stop**, **status**, … 这些命令中，也可以用于 `Wants=`, `Requires=`, `Before=`, `After=`, … 这些依赖关系选项中。 但是务必注意，不可将单元的别名用于 **preset** 命令中。 再次提醒，通过 `Alias=` 设置的别名仅在单元被启用(enable)之后才会生效。

## Unit File 模板 && 多实例

Unit File 可以通过 `@` 符号声明一个**模板文件**，通过在 `@` 添加字符串，可以根据模板文件，生成一个**实例文件**。经过实例化后的实例文件，就是一个真实可用的 Unit File 了

就拿 Wireguard 的 Unit File 为例

```bash
~]$ systemctl cat wg-quick@.service
# /lib/systemd/system/wg-quick@.service
[Unit]
Description=WireGuard via wg-quick(8) for %I
After=network-online.target nss-lookup.target
Wants=network-online.target nss-lookup.target
PartOf=wg-quick.target

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/wg-quick up %i
ExecStop=/usr/bin/wg-quick down %i
Environment=WG_ENDPOINT_RESOLUTION_RETRIES=infinity

[Install]
WantedBy=multi-user.target

```

上面就是一个模板文件，我们需要实例化他，那么只需要在 `@` 符号后面添加字符串即可，`systemctl enable wg-quick@wg0.service`，此时，将会生成如下 Unit File

```bash
[lichenhao@hw-cloud-xngy-jump-server-linux-2 ~]$ systemctl cat wg-quick@wg0.service
# /lib/systemd/system/wg-quick@.service
[Unit]
Description=WireGuard via wg-quick(8) for %I
After=network-online.target nss-lookup.target
Wants=network-online.target nss-lookup.target
PartOf=wg-quick.target

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/usr/bin/wg-quick up %i
ExecStop=/usr/bin/wg-quick down %i
Environment=WG_ENDPOINT_RESOLUTION_RETRIES=infinity

[Install]
WantedBy=multi-user.target
```

当我们启动该服务时，其实就是执行了 ExecStart 指令中的命令，并且 `%i` 被替换成了 wg0，即 `/usr/bin/wg-quick up wg0`

这种模板化的 Unit File 特别适合这种服务来使用，还有那些配置 TTY、磁盘名称的，都可以将想要实例化的 TTY 或者 分区，作为 %i 配置。

## drop-in(嵌入式) Unit File(就是 include 功能)

**drop-in(嵌入式)** 单元文件就是一种类似配置文件的 **include**的功能(比如 Nginx 中的 include 指令)，可以让主配置文件包含其他子配置文件。**Systemd 设定了两种 include 的规范**

假如现在有一个名为 `foo.service` 的 Unit File，那么，Systemd 会从加载 Unit File 的目录中，加载与之相关联的一系列文件

- **UnitFileName.wants/ 与 UnitFileName.requires/** # 比如 foo.service.wants/  与  foo.service.requires/。该目录中可以放置许多指向其他 Unit Files 的软连接。 软连接所指向的 Unit 将会被当做  `foo.service`  的 Unit 文件中  `Wants=`  与  `Requires=`  指令的值
  - 注意：即使文件中不存在 Wants 和 Requires 指令。只要存在对应的 _.wants/ 和_.requires/ 目录，就相当于为 Unit File 中加上了这两个指令。
  - 这样就可以方便的为 Unit 添加依赖关系，而无需修改单元文件本身。 向  `*.wants/`  与  `*.requires/`  目录中添加软连接的首选方法是使用  [systemctl(1)](http://www.jinbuguo.com/systemd/systemctl.html#)  的  **enable**  命令， 它会读取 Unit File 的 \[Install] 部分。
- **UnitFileName.d/** # 比如 foo.service.d/。这就是配置文件的 include 功能。当解析完主 Unit File 之后，该目录中所有以 `.conf` 结尾的文件，都会被依次附加到主 Unit File 的末尾。
  - 这样就可以方便的修改 Unit 的设置，或者为 Unit 添加额外的设置，而无需修改 Unit File 本身。
  - 注意，include 功能中的文件遵守如下规则：
    - 必须包含明确的 Sections (例如 `[Service]` 之类)。
    - 对于从模板文件实例化而来的 Unit，会优先读取与此实例对应的 `UnitFileName.d/` 目录(例如 "`foo@bar.service.d/`")中的配置片段(`*.conf` 文件)， 然后才会读取与模板对应的 "`.d/`" 目录(例如 "`foo@.service.d/`")中的配置片段("`.conf`" 文件)。
    - 对于名称中包含连字符("`-`")的单元，将会按特定顺序依次在一组(而不是一个)目录中搜索单元配置片段。 例如对于  `foo-bar-baz.service`  单元来说，将会依次在  `foo-.service.d/`, `foo-bar-.service.d/`, `foo-bar-baz.service.d/`  目录下搜索单元配置片段。
      - 这个机制可以方便的为一组相关单元(单元名称的前缀都相同)定义共同的单元配置片段， 特别适合应用于 mount, automount, slice 类型的单元， 因为这些单元的命名规则就是基于连字符构建的。
      - 注意，在前缀层次结构的下层目录中的单元配置片段，会覆盖上层目录中的同名文件， 也就是  `foo-bar-.service.d/10-override.conf`  会覆盖(取代) `foo-.service.d/10-override.conf`  文件。

注意：通常情况下，drop-in Unit File 放在 `/etc/systemd/{system,user}`  目录中， 还可以放置在  `/usr/lib/systemd/{system,user}`  与  `/run/systemd/{system,user}`  目录中。 虽然在优先级上，`/etc`  中的配置片段优先级最高、`/run`  中的配置片段优先级居中、 `/usr/lib`  中的配置片段优先级最低。但是这仅对同名配置片段之间的覆盖关系有意义。 因为所有 `.d/` 目录中的配置片段，无论其位于哪个目录， 都会被按照文件名的字典顺序，依次覆盖单元文件中的设置(相当于依次附加到主单元文件的末尾)。

### drop-in 示例

- /usr/lib/systemd/system/ssh.service # ssh 这个 Unit 的主配置文件
- /etc/systemd/system/ssh.service.d/CUSTOM.conf # 在 /etc/systemd/system 下面创建与配置文件相同文件名的目录,但是要加上 .d 的扩展名。然后在该目录下创建 .conf 结尾的配置文件即可。
- /etc/systemd/system/vsftpd.service.wants/ # 此目录内的文件为链接文件,设置相依服务的链接。意思是启动了 vsftpd.service 之后,最好再加上这目录下面建议的服务。
- /etc/systemd/system/vsftpd.service.requires/ # 此目录内的文件为链接文件,设置相依服务的链接。意思是在启动 vsftpd.service 之前,需要事先启动哪些服务的意思。

配置文件分为 4 种状态，当启用(enable)该文件的时候，从配置文件目录建立一个软连接到 /etc/systemd/system/ 目录下，当禁用(disable)该文件的时候，会把该软连接删除。如果在/etc/systemd/system/目录下有 Unit 的配置文件，则开机则会自动加载并启动 Unit。可以使用命令 systemctl list-unit-files 命令查看所有的配置文件状态。结论：建立了连接则说明该 Unit 会开机启动，没建立连接则该 Unit 不会开机启动；还可以禁止该 Unit 建立连接，则说明该 Unit 永远不能开机启动。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/gvdc29/1620129636770-82d464a7-f1ef-4e5c-9749-61b258d1b0e6.png)

- enabled：启用(该文件已建立链接)
- disabled：禁用(该文件没建立链接)
- static：这个文件不可以被建立链接，但是可以被其他服务进行关联启动
- masked：该配置文件被禁止建立启动链接，可以使用 systemctl unmask 命令开启

## 总结以及注意事项

在使用 systemctl 命令时，会有两个子命令，systemctl list-units 和 systemctl list-unit-files，这俩命令的区别如下

- list-units 是列出已经载入的 Unit(加上-a 选项可以看 start 和 stop 所有的)，其中还包括各种.device 的 unit，这类设备是硬件，没有配置文件的，不会再 list-unit-files 命令下列出来
- list-unit-files 是列出所有已经安装的 Unit 的配置文件(有一部分 Unit 是没有配置文件的)，安装完的 Unit 会把配置文件存放在/lib/systemd/system/目录下，而且通过查看 list-unit-files 还能看出来是否是开机启动，具体原因请见下文的《systemd 的 unit 配置文件说明》
- 如果一个单元文件的大小为零字节或者是指向 `/dev/null` 的软连接， 那么它的所有相关配置都将被忽略。同时，该单元将被标记为 "`masked`" 状态，并且无法被启动。 这样就可以 彻底屏蔽一个单元(即使手动启动也不行)。

# 各种 Unit 额外说明

## Target Unit

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gvdc29/1616167393736-e7bb0f5d-83be-4d38-856d-a22c56e9fab1.jpeg)

启动计算机的时候，需要启动大量的 Unit。如果每一次启动，都要一一写明本次启动需要哪些 Unit，显然非常不方便。Systemd 的解决方案就是 Target。

简单说，Target 就是一个 Unit 组，包含许多相关的 Unit 。启动某个 Target 的时候，Systemd 就会启动里面所有的 Unit。从这个意义上说，Target 这个概念类似于"状态点"，启动某个 Target 就好比启动到某种状态。

注意

- 在启动系统的时候，也是首先通过查询 target 中的内容，以便启动相应的 Unit
- 传统的 init 启动模式里面，有 RunLevel 的概念，Target 就能起到同样的效果。不同的是，RunLevel 是互斥的，不可能多个 RunLevel 同时启动，但是多个 Target 可以同时启动，并且更加灵活，可以自己定义每个 Target 可以包含的 Unit。比如启动 graphical.target 则里面就包含运行图形界面的 Unit 和 multi-user.target 中的所有 Unit。
- 我们可以通过 systemctl set-default UNIT 命令来设定系统启动时，默认启动的一组 Unit

# Unit File 加载示例

可以通过 [系统启动流程](/docs/1.操作系统/Operating%20system/Unix-like%20OS/系统启动流程.md) 看出来 Systemd 是如何加载 Unit File 的

