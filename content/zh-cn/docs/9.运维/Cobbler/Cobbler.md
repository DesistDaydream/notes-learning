---
title: Cobbler
---

# 前言

网络安装服务器套件 Cobbler(补鞋匠) 出现以前，我们一直在做装机民工这份很有前途的职业。自打若干年前 Red Hat 推出了 Kickstart，此后我们顿觉身价倍增。不再需要刻了光盘一台一台地安装 Linux，只要搞定 PXE、DHCP、 TFTP，还有那满屏眼花缭乱不知所云的 Kickstart 脚本，我们就可以像哈里波特一样，轻点魔棒，瞬间安装上百台服务器。这一堆花里胡哨的东西可不是一般人都能整明白的，没有大专以上学历，通不过英语四级， 根本别想玩转。

总而言之，这是一份多么有前途，多么有技术含量的工作啊。

很不幸，Red Hat 最新（Cobbler 项目最初在 2008 年左右发布）发布了网络安装服务器套件 Cobbler(补鞋匠)，它已将 Linux 网络安装的技术门槛，从大专以上文化水平，成功降低到初中以下，连补鞋匠都能学会。

对于我们这些在装机领域浸淫多年，经验丰富，老骥伏枥，志在千里的民工兄弟们来说，不啻为一个晴天霹雳。

# 概述

> 参考：
> - [Cobbler 官网](http://cobbler.github.io/)

Cobbler 是一个 Linux 服务器快速网络安装的服务，而且在经过调整也可以支持网络安装 windows。

该工具使用 python 开发，小巧轻便（才 15k 行 python 代码），可以通过网络启动(PXE)的方式来快速安装、重装物理服务器和虚拟机，同时还可以管理 DHCP，DNS，TFTP、RSYNC 以及 yum 仓库、构造系统 ISO 镜像。

Cobbler 可以使用命令行方式管理，也提供了基于 Web 的界面管理工具(cobbler-web)，还提供了 API 接口，可以方便二次开发使用。

Cobbler 是较早前的 kickstart 的升级版，优点是比较容易配置，还自带 web 界面比较易于管理。

Cobbler 内置了一个轻量级配置管理系统，但它也支持和其它配置管理系统集成，如 Puppet，暂时不支持 SaltStack。

Cobbler 客户端 Koan 支持虚拟机安装和操作系统重新安装，使重装系统更便捷。

## Cobbler 可以干啥

使用 Cobbler，您无需进行人工干预即可安装机器。Cobbler 设置一个 PXE 引导环境（它还可以使用 yaboot 支持 PowerPC），并 控制与安装相关的所有方面，比如网络引导服务（DHCP 和 TFTP）与存储库镜像。当希望安装一台新机器时，Cobbler 可以：

1）使用一个以前定义的模板来配置 DHCP 服务（如果启用了管理 DHCP）

2）将一个存储库（yum 或 rsync）建立镜像或解压缩一个媒介，以注册一个新操作系统

3）在 DHCP 配置文件中为需要安装的机器创建一个条目，并使用指定的参数（IP 和 MAC）

4）在 TFTP 服务目录下创建适当的 PXE 文件

5）重新启动 DHCP 服务来反应新的更改

6）重新启动机器以开始安装（如果电源管理已启动）

## Cobbler 支持的系统和功能

Cobbler 支持众多的发行版：RedHat、Fedora、CentOS、Debian、Ubuntu 和 SUSE。当添加一个操作系统（通常通过使用 ISO 文件）时，Cobbler 知道如何解压缩合适的文件并调整网络服务，以正确引导机器。

Cobbler 可以使用 kickstart 模板。基于 Red Hat 或 Fedora 的系统使用 kickstart 文件来自动化安装流程。通过使用模板，就会拥有基本的 kickstart 模板，然后定义如何针对一种配置文件或 机器配置而替换其中的变量。例如，一个模板可能包含两个变量 `$domain` 和 `$machine_name` 在 Cobbler 配置中，一个配置文件指定 domain=mydomain.com，并且每台使用该配置文件的机器在 machine_name 变量中指定其名称。该配置文件的所有机器都使用相同的 kickstart 安装且针对 domain=mydomain.com 进行配置，但每台机器拥有其自己的机器名称。您仍然可以使用 kickstart 模板 在不同的域中安装其他机器并使用不同的机器名称。

为了协助管理系统，Cobbler 可通过 fence scripts 连接到各个电源管理环境。Cobbler 支持 apc_snmp、bladecenter、bullpap、drac、 ether_wake、ilo、integrity、ipmilan、ipmitool、lpar、rsa、virsh 和 wti。要重新安装一台机器，可 运行 reboot system foo 命令，而且 Cobbler 会使用必要的 和信息来为您运行恰当的 fence scripts（比如机器插槽数）。

除了这些特性，还可以使用一个配置管理系统（CMS）。你有两种选择：该工具内的一个内部系统，或者现成的外部 CMS，比如 Chef 或 Puppet。借助内部系统，你可以指定文件模板，这些模板会依据配置参数进行处理（与 kickstart 模板的处理方式一样），然后复制到你指定的位 置。如果必须自动将配置文件部署到特定机器，那么此功能很有用。

使用 koan 客户端，Cobbler 可从客户端配置虚拟机并重新安装系统。

## Cobbler 各个组件之间关系

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tegfcw/1616125382993-9ef824a6-a456-4fc4-9eeb-8ca756ca7705.jpeg)

主要目的配置网络接口：

Cobbler 的配置结构基于一组注册的对象。每个对象表示一个与另一个实体相关联的实体（该对象指向另一个对象，或者另一个对象指向该对象）。当一个对象指向另一个对象时，它就继承了被指向对象的数据，并可覆盖或添加更多特定信息。以下对象类型的定义

1. Distros（发行版）：表示一个操作系统，它承载了内核和 initrd 的信息，以及内核参数等其他数据。使用 cobbler import 命令后即可生成该对象
2. Profile（配置文件）：包含一个 Distros、一个 kickstart 文件以及可能的存储库，还包含更多特定的内核参数等其他数据。使用 cobbler import 命令后，会默认使用名为/var/lib/cobbler/kickstarts/sample_end.ks 的 kickstart 文件。
3. Systems（系统）：表示将要安装的新机器。它包含一个配置文件或一个镜像，还包含该机器的 IP 和 MAC 地址、电源管理（地址、凭据、类型）、（网卡绑定、设置 valn 等）
4. Repository（镜像）：保存一个 yum 或 rsync 存储库的镜像信息
5. Image（存储库）：可替换一个包含不属于此类比的额文件的发行版对象（例如，无法分为内核和 initrd 的对象）。

基于注册的对象以及各个对象之间的关联，Cobbler 知道如何更改文件系统以反应具体配置。因为系统配置的内部是抽象的，所以可以仅关注想要执行的操作。
