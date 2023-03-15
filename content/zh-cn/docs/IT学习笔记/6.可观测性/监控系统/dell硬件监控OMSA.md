---
title: dell硬件监控OMSA
weight: 50
---

# OMSA（全称 Openmanage Server Administrator),是戴尔公司自主研发的 IT 系统管理解决方案

<http://linux.dell.com/>
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gknv4x/1616067497715-90aa8503-f2c3-4436-ac53-560ef9f66d5f.jpeg)

## OMSA 的安装

### 自动安装

<https://linux.dell.com/repo/hardware/dsu/>

1. 配置存储库
   1. curl -O <https://linux.dell.com/repo/hardware/dsu/bootstrap.cgi>
   2. bash bootstrap.cgi
2. yum install srvadmin-all.x86_64

Note：

1. srvadmin-base # 代理程序，可以生成 snmp 信息
2. srvadmin-server-cli # 命令行客户端，可以通过命令行查看硬件信息
3. srvadmin-storage-cli # 存储资源的命令行客户端，不安装这个则无法获取 raid 和硬盘的数据

下面是使用 yum 安装 OMSA 的 repo 文件

```bash
~]# cat /etc/yum.repos.d/DELL-OMSA.repo
[dell-system-update_independent]
name=dell-system-update_independent
baseurl=https://linux.dell.com/repo/hardware/dsu/os_independent/
gpgcheck=1
gpgkey=https://linux.dell.com/repo/hardware/dsu/public.key
      https://linux.dell.com/repo/hardware/dsu/public_gpg3.key
enabled=1
exclude=dell-system-update*.i386

[dell-system-update_dependent]
name=dell-system-update_dependent
mirrorlist=https://linux.dell.com/repo/hardware/dsu/mirrors.cgi?osname=el$releasever&basearch=$basearch&native=1
gpgcheck=1
gpgkey=https://linux.dell.com/repo/hardware/dsu/public.key
      https://linux.dell.com/repo/hardware/dsu/public_gpg3.key
enabled=1
```

参考文章：<http://www.madown.com/2017/05/23/81/>

### 手动安装

<https://www.dell.com/support/home/>去该网站输入主机号查询，然后根据关键字搜索 OMSA 并下载
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gknv4x/1616067497724-a84a9901-d1a0-4c18-9198-303c071c9960.jpeg)

1. 解压已下载的安装包
   1. mkdir dell-omsa
   2. tar -zxvf OM-SrvAdmin-Dell-Web-LX-9.3.0-3465_A00.tar -C dell-omsa #
2. 安装 rpm 包
   1. cd dell-omsa/linux/RPMS/supportRPMS/srvadmin/RHEL7/x86_64
   2. yum localinstall \*.rpm

将 dell-r740.tar.gz 拷贝到/root/Download 目录下并执行以下脚本

## OMSA 的配置与使用

安装 dell 监控 openManager 相关组件完成后，会在/opt/dell/\*目录下生成配置文件与可执行文件

- 配置 openManager

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gknv4x/1616067497704-782c72fb-2dcf-4422-9654-b71a7c89d1eb.jpeg)

- /opt/dell/srvadmin/sbin/srvadmin-services.sh start # 启动 openManger
  - 该脚本会通过 systemd 启动三个服务
  - instsvcdrv #
  - dataeng #
  - dsm_om_connsvc # web 控制台
- /opt/dell/srvadmin/sbin/srvadmin-services.sh enable # 设置开机自启 openManager
- systemctl stop dsm_om_connsvc # 关闭 openManager 的 web 服务
- systemctl disable dsm_om_connsvc # 将 openManager 的 web 服务开机自启关闭
- systemctl restart snmpd # 重启 snmp 服务。由于安装 openManager 会在 snmpd 的配置文件中写入内容，所以需要重启 snmpd 服务使得该配置生效
- 在 wiseman 上添加相关的 dell 硬件监控。效果如图，在主机的模板里添加 dell server 模板

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gknv4x/1616067497780-3afd659e-d460-4c2b-8d66-1f9a3c67890a.jpeg)

# 命令行工具使用说明

## omreport COMMAND

在任何时候都可以使用 omreport -?来获取命令帮助以查看都有哪些可用的 COMMAND，例如：omreport -?、omreport chassis -?等
COMMAND

1. about Product and version properties.
2. licenses Displays the digital licenses of the installed hardware devices.
3. preferences Report system preferences.
4. system System component properties.
5. chassis 机架组件的属性。i.e.基本硬件的信息。Chassis component properties.
6. storage # 显示存储组件的属性

## omreport storage

EXAMPLE

1. omreport storage pdisk controller=0 #
