---
title: kickstart 介绍
---

# KickStart 简介

官方文档：<https://docs.centos.org/en-US/centos/install-guide/Kickstart2>

# Kickstart 语法参考

官方文档：<https://docs.centos.org/en-US/centos/install-guide/Kickstart2/#sect-kickstart-syntax>

## 磁盘配置

bootloader \[OPTIONS] # 引导程序(boot loader)的相关配置。

\--location=VALUE # 指定引导程序的安装位置为 VALUE

1. mbr # 默认值。取决于磁盘格式是 MBR 还是 GUID
2. partition # 将引导程序安装在包含内核的分区的第一个扇区中
3. none # 不安装引导程序
4. boot # 未知，待更新

磁盘分区相关配置

part MntPoint \[OPTIONS] # 创建一个分区，挂载点为 MntPoint，

MntPoint # 可用的挂载点有如下几种。

- /PATH # 指定挂载点到具体路径下
- swap # 指定该分区为 swap
- raid.ID
- pv.ID # 指定该分区用于 lvm，即创建一个 pv(物理卷)
- biosboot # 指定该分区用于 BIOS 引导,建议大小为 2MiB
- /boot/efi # 指定该分区用于 UEFI 引导，建议大小为 200MiB.

OPTIONS

- --size=NUM # 指定该分区的大小，单位为 MiB。Note：NUM 为一个正整数(不包含单位)
- --grow # 将该分区大小设置为所有剩余可用的空间。如果指定了 --maximum 选项，则将该分区设置为该选项值的大小。
- --asprimary # 指定该分区为主分区。Note：对于 GUID 分区表(GPT),该选项没有任何意义。
- --fstype=TYPE # 指定该分区的文件系统类型。可用类型有 xfs、ext2、ext3、ext4、swap、vfat、efi、biosboot
- --ondisk= # 指定要使用的磁盘名称

volgroup Name PartName # 创建名为 NAME 卷组，使用名为 PartName 分区

logvol PATH --vgname=VGNAME --name=NAME \[OPTIONS] # 创建名为 NAME 逻辑卷，使用名为 VGNAME 的卷组

配置示例：

    # 配置引导程序的信息，将引导程序安装到 mbr 中，
    bootloader --location=mbr --append="crashkernel=auto"
    # 清理所有分区
    clearpart --all --initlabel
    # 创建主分区，挂载到/boot，文件系统类型为xfs，大小为500M，使用sda磁盘
    part /boot --fstype=xfs --asprimary --size=500 --ondisk=sda
    part biosboot --fstype=biosboot --asprimary --size=2 --ondisk=sda
    # 分区，创建一个名为pv.01的物理卷，大小为磁盘剩余的所有空间，使用sda磁盘
    part pv.01 --size=1 --grow --ondisk=sda
    # 使用pv.01来创建一个vg0的卷组
    volgroup vg0 pv.01
    # 使用vg0卷组来创建逻辑卷，逻辑卷为swap分区，文件系统类型为swap，大小为32G，名字为swap
    logvol swap --fstype swap --size=32768 --name=swap --vgname=vg0
    # 使用vg0卷组来创建逻辑卷，挂载到/目录下，文件系统类型为xfs，大小使用逻辑卷剩下的所有空间，名字为root
    logvol / --fstype xfs --size=1 --grow --name=root --vgname=vg0
    # 分区，创建一个名为pv.02的物理卷，大小为磁盘剩余的所有空间，使用sdb磁盘
    part pv.02 --size=1 --grow --ondisk=sdb
    volgroup vg1 pv.02
    logvol /app --fstype xfs --size=1 --grow --name=app --vgname=vg1

    # UEFI模式下需要增加下面一条分区，且不用创建biosboot分区
    part /boot/efi --fstype="efi" --ondisk=sda --size=500 --fsoptions="umask=0077,shortname=winnt"

Packages 安装配置

官方文档：<https://docs.centos.org/en-US/centos/install-guide/Kickstart2/#sect-kickstart-packages>

使用 ％packages %end 关键字指定 Kickstart 安装哪些软件包以及如何安装它们，可以指定 environment、group 或者单独的包名。这些可用的 packages 可以在 CentOS 的安装镜像的 .repodata/\* 目录中找到

比如当我把 iso 镜像文件挂载到 /mnt 目录下时，可以看到如下内容。(那些 xml 文件中可以看到 packages 列表及其依赖)

    [root@cobbler repodata]# pwd
    /mnt/BaseOS/repodata
    [root@cobbler repodata]# ls
    08a9add0907af002934460d81ea0edc8bb4154db679cdc113d4c51efcbddfce4-comps-BaseOS.x86_64.xml.xz  89376911bec34defd11535ee1f9e74237ec25e63c4dc5041ed519f1166d1cdca-other.xml.gz             repomd.xml
    190ce1d49a76eafb61b0e2738d7331a45f1efbdfb4c2c274176486f8c82f7f80-primary.xml.gz              a6e31179a63dffb12846a2f8ad848618e7669225deb4222acadcbabb04d522f0-filelists.sqlite.xz      TRANS.TBL
    2caf33eae61c01725123a875217fe6bb0754c2916e7336286fbc392d23f42b57-other.sqlite.xz             a9b7530c36c9681b97c159cd98693e7fffecf2d1018d508f990934a4fa71b447-primary.sqlite.xz
    53db8eac92f79abf479e202fb013aca86d5060fd948e1ed7ea8f9493e72fd4d1-filelists.xml.gz            fe7d9972481aaef922e8a2fafa4065c4eb9422125da783a058a9988cd9f3eb27-comps-BaseOS.x86_64.xml

.xml 文件中包含描述可用环境（用标记）和组（标记）的结构。每个条目都有一个 ID，用户可见性的值，名称，描述和包列表。如果选择了要安装的组，则始终会安装软件包列表中标记为必选的软件包，如果未在其他位置专门排除标记为默认的软件包，则即使已选择该组，也必须在其他位置特别包括标记为可选的软件包。

配置示例：

    %packages
    安装 minimal-enviroment 环境的所有软件包，一个 Kickstart 文件仅可指定一个 environment 包
    @^minimal-environment
    安装 core 组下的所有软件包
    @core
    安装 kexec-tools 软件包
    kexec-tools
    %end

kickstart 配置文件样例

    #platform=x86, AMD64, or Intel EM64T
    # System authorization information系统认证信息，使用加密，md5方式加密
    auth  --useshadow  --enablemd5
    # System bootloader configuration 系统引导配置
    # --location指定创建引导的位置，在mbr中创建引导；--append指定内核参数，crashkernel为开启kdump
    bootloader --location=mbr --append="crashkernel=auto"
    # Partition clearing information分区清除信息
    # 清除所有分区，并初始化磁盘标签
    clearpart --all --initlabel
    # Use text mode install
    text
    # Firewall configuration 关闭防火墙
    firewall --disabled
    # Run the Setup Agent on first boot
    firstboot --disable
    # System keyboard 系统键盘
    keyboard us
    # System language 系统语言
    lang en_US
    # Use network installation 指明安装系统的方式，这里使用网络方式安装，指明提供安装程序的服务器地址和路径
    url --url=$tree
    # If any cobbler repo definitions were referenced in the kickstart profile, include them here.
    $yum_repo_stanza
    # $SNIPPET 变量括号内的值是目录/var/lib/cobbler/snippets下的文件，该文件中可以写入linux命令，当做脚本文件来说明
    # Network information
    $SNIPPET('network_config')
    # Reboot after installation 安装完成后重启系统
    reboot
    #Root password设定root的密码
    rootpw --iscrypted $default_password_crypted
    # SELinux configuration关闭SELinux
    selinux --disabled
    # Do not configure the X Window System不要配置X window系统
    skipx
    # System timezone设定系统时区
    timezone  Asia/Shanghai
    # Install OS instead of upgrade 重新安装操作系统，而不是升级
    install
    # Clear the Master Boot Record清除主引导记录
    zerombr
    # Allow anaconda to partition the system as needed该选项用于自动分区
    #autopart
    # Disk partitioning information磁盘分区信息
    part /boot --fstype=xfs --asprimary --size=500
    part biosboot --fstype=biosboot --asprimary --size=2
    part pv.01 --size=1 --grow
    volgroup vg0 pv.01
    logvol / --fstype xfs --size=10240 --name=root --vgname=vg0
    logvol swap --fstype swap --size=1024 --name=swap --vgname=vg0
    logvol /var --fstype xfs --size=1 --grow --name=var --vgname=vg0
    %pre
    # %pre段落为安装前执行的任务
    %end
    %packages
    # %package段落中指定要安装的软件包
    @^minimal
    @core
    kexec-tools
    %end
    %post
    #`%post`段落为安装系统完成后执行的任务
    # 自定义系统配置
    # 安装工具，在/var/lib/cobbler/snippets目录下添加名为tools的文件，文件中可以写入想要执行的linux命令
    $SNIPPET('tools')
    # End final steps
    %end
