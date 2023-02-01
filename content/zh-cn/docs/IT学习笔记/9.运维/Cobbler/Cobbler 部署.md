---
title: Cobbler 部署
---

# 概述

## 基础环境准备

- 检测 selinux 是否关闭(必须关闭)
  - setenforce 0
  - sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
- 检测防火墙是否关闭(必须关闭)
  - systemctl disable firewalld && systemctl stop firewalld

注意：

- 虚拟机网卡采用桥接模式，不使用 NAT 模式，我们会搭建 DHCP 服务器，在同一局域网多个 DHCP 服务会有冲突，所以最好把路由器的 DHCP 服务关闭

## 安装 Cobbler 以及相关功能配件

- yum install cobbler cobbler-web pykickstart httpd dhcp tftp-server fence-agents -y
  - cobbler #cobbler 程序包
  - cobbler-web #cobbler 的 web 服务包
  - pykickstart #cobbler 检查 kickstart 语法错误
  - httpd #Apache web 服务
  - dhcp #dhcp 服务
  - tftp-server #tftp 服务
- systemctl enable cobblerd httpd tftp rsyncd dhcpd && systemctl start cobblerd httpd tftp rsyncd dhcpd

## 检测 Cobbler

cobbler 的运行依赖于 dhcp、tftp、rsync 及 dns 服务，其中 dhcp 可由 dhcpd（isc）提供，也可由 dnsmasq 提供；tftp 可由 tftp-server 程序包提供，也可由 cobbler 功能提供，rsync 有 rsync 程序包提供，dns 可由 bind 提供，也可由 dnsmasq 提供

cobbler 可自行管理这些服务中的部分甚至是全部，但需要配置文件/etc/cobbler/settings 中的“manange_dhcp”、“manager_tftpd”、“manager_rsync”、“manager_dns”分别来进行定义，另外，由于各种服务都有着不同的实现方式，如若需要进行自定义，需要通过修改/etc/cobbler/modules.conf 配置文件中各服务的模块参数的值来实现。

### 检查配置文件，并修改其中的错误提示。需要在 cobblerd 和 httpd 启动的情况下检查

检查 cobbler 配置存在的问题,逐一解决

cobbler check

使用该命令后，会出现若干报错，以下是对几个报错的总结

备份将要修改的文件文件

cp /etc/cobbler/settings{,.bak}

cp /etc/xinetd.d/tftp{,.bak}

1. The 'server' field in /etc/cobbler/settings must be set to something other than localhost, or kickstarting features will not work. This should be a resolvable hostname or IP for the boot server as reachable by all machines that will use it.
   1. 修改/etc/cobbler/settings 文件中的 server 参数的值为提供 cobbler 服务的主机相应的 IP 地址或主机名 server，Cobbler 服务器的 IP，修改 384 行 server: 127.0.0.1
      1. sed -i 's/server: 127.0.0.1/server: 192.168.10.11/' /etc/cobbler/settings
2. For PXE to be functional, the 'next_server' field in /etc/cobbler/settings must be set to something other than 127.0.0.1, and should match the IP of the boot server on the PXE network
   1. 修改/etc/cobbler/settings 文件中的 next_server 参数的值为提供 PXE 服务的主机相应的 IP 地址，如 next_server: 192.168.31.73；
      1. server，pxe 服务器的 IP，由于这里使用的是同一台机器，所以填 Cobbler 服务器的 IP，修改 272 行 next_server: 127.0.0.1
         1. sed -i 's/next_server: 127.0.0.1/next_server: 192.168.10.11/' /etc/cobbler/settings
3. change 'disable' to 'no' in /etc/xinetd.d/tftp
   1. 修改/etc/xinetd.d/tftp 文件中的 disable 参数修改为 disable = no
      1. sed -i 's/disable.\*= yes/disable = no/g' /etc/xinetd.d/tftp
4. some network boot-loaders are missing from /var/lib/cobbler/loaders, you may run 'cobbler get-loaders' to download them, or, if you only want to handle x86/x86*64 netbooting, you may ensure that you have installed a \_recent* version of the syslinux package installed and can ignore this message entirely. Files in this directory, should you want to support all architectures, should include pxelinux.0, menu.c32, elilo.efi, and yaboot. The 'cobbler get-loaders' command is the easiest way to resolve these requirements.
   1. 执行 cobbler get-loaders 命令即可；否则，需要安装 syslinux 程序包，而后复制/usr/share/syslinux/{pxelinux.0,memu.c32}等文件至/var/lib/cobbler/loaders/目录中
      1. cobbler get-loaders
5. debmirror package is not installed, it will be required to manage debian deployments and repositories
   1. 安装 debmirror
      1. yum install debmirror -y
      2. cp /etc/debmirror.conf{,.bak}
      3. sed -i '/@dists/s/^/#&/' /etc/debmirror.conf
      4. sed -i '/@arches/s/^/#&/' /etc/debmirror.conf
6. The default password used by the sample templates for newly installed machines (default_password_crypted in /etc/cobbler/settings) is still set to 'cobbler' and should be changed, try: "openssl passwd -1 -salt 'random-phrase-here' 'your-password-here'" to generate new one
   1. 生成密码来取代默认的密码，更安全。该密码用来给新安装的设备提供 root 密码。根据提示 your-password-here，这里是自己的密码，random-phrase-here，这里是随机的干扰码
      1. openssl passwd -1 -salt 'lichenhao' '123456'
         1. $1$lichenha$UkZ9KiUaiwS/0C324YtoP0
      2. sed -i 's/'default_password_crypted:.\*'/'default_password_crypted: "$1$lichenha$UkZ9KiUaiwS/0C324YtoP0"/g' /etc/cobbler/settings

其他一些没有提示报错的小修改

1. 用 cobbler 管理 DHCP，修改 242 行 manage_dhcp: 0
   1. sed -i 's/manage_dhcp: 0/manage_dhcp: 1/g' /etc/cobbler/settings
2. 防止循环装系统，适用于服务器第一启动项是 PXE 启动，修改 292 行 pxe_just_once: 0
   1. sed -i 's/pxe_just_once: 0/pxe_just_once: 1/' /etc/cobbler/settings

修改完后重启服务再进行一次检测,若提示 No configuration problems found. All systems go. 则可继续后续步骤

1. systemctl restart cobblerd.service
2. cobbler check

### 配置 DHCP

1. 修改 cobbler 的 dhcp 模版，这个模板会覆盖 dhcp 本身的配置文件。
2. cp /etc/cobbler/dhcp.template{,.bak}
3. vim /etc/cobbler/dhcp.template
4. 在文件末尾添加如下内容
   1. subnet 192.168.10.0 netmask 255.255.255.0 { #规划一段子网以便用于被安装系统的设备来获取这一段 IP
   2. option domain-name-servers 114.114.114.114; #指定该子网的 DNS
   3. option routers 192.168.10.2; #指定该段子网的网关
   4. range dynamic-bootp 192.168.10.100 192.168.10.250; #指定给被安装系统的设备可用的 IP 段
   5. option subnet-mask 255.255.255.0; #指定该子网的掩码
   6. next-server $next_server;
   7. default-lease-time 43200;
   8. max-lease-time 86400;
   9. }

## 导入镜像

1. 挂载 centos 光盘镜像
   1. mount /dev/cdrom /mnt/
2. 导入镜像,创建一个可以提供给其余被安装设备将要使用到的系统镜像
   1. cobbler import --path=/mnt/ --name=CentOS7-2003 --arch=x86_64
   2. 安装源的唯一标示就是根据 name 参数来定义，本例导入成功后，安装源的唯一标示就是：CentOS7-2003-x86_64，如果重复，系统会提示导入失败。
3. 同步 cobbler 的配置，可以看到同步干了哪些事
   1. cobbler sync
4. 查看镜像列表
   1. cobbler distro list
      1. 镜像存放目录，cobbler 会将镜像中的所有安装文件拷贝到本地一份，放在/var/www/cobbler/ks_mirror 下的 CentOS7-2003-x86_64 目录下。因此/var/www/cobbler 目录必须具有足够容纳安装文件的空间，查看一下目录内容
         1. cd /var/www/cobbler/ks_mirror/

### 配置 ks.cfg

1. cd /var/lib/cobbler/kickstarts/
2. vim centos7.ks #（sample_end.ks（默认使用的 ks 文件））修改成以下内容
3. 注意：文件中不能有中文，即使注释掉也不行，否则会导致安装失败


    #platform=x86, AMD64, or Intel EM64T
    # System authorization information系统认证信息，使用加密，md5方式加密
    auth  --useshadow  --enablemd5
    # System bootloader configuration系统引导配置
    # --location指定创建引导的位置，在mbr中创建引导；--append指定内核参数，crashkernel为开启kdump
    bootloader --location=mbr --append="crashkernel=auto"
    # Partition clearing information分区清除信息
    # 清除所有分区，并初始化磁盘标签
    clearpart --all --initlabel
    # Use text mode install
    text
    # Firewall configuration关闭防火墙
    firewall --disabled
    # Run the Setup Agent on first boot
    firstboot --disable
    # System keyboard系统键盘
    keyboard us
    # System language系统语言
    lang en_US
    # Use network installation指明安装系统的方式，这里使用网络方式安装，指明提供安装程序的服务器地址和路径
    url --url=$tree
    # If any cobbler repo definitions were referenced in the kickstart profile, include them here.
    $yum_repo_stanza
    # $SNIPPET变量括号内的值是目录/var/lib/cobbler/snippets下的文件，该文件中可以写入linux命令，当做脚本文件来说明
    # Network information
    $SNIPPET('network_config')
    # Reboot after installation安装完成后重启系统
    reboot

    #Root password设定root的密码
    rootpw --iscrypted $default_password_crypted
    # SELinux configuration关闭SELinux
    selinux --disabled
    # Do not configure the X Window System不要配置X window系统
    skipx
    # System timezone设定系统时区
    timezone  Asia/Shanghai
    # Install OS instead of upgrade重新安装操作系统，而不是升级
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

1. 在第一次导入系统镜像后，Cobbler 会给镜像指定一个默认的 kickstart 自动安装文件在/var/lib/cobbler/kickstarts 下的 sample_end.ks。
2. cobbler list
3. 查看所有的 profile 设置
   1. cobbler profile report
4. 编辑 profile，修改关联的 ks 文件
   1. cobbler profile edit --name=CentOS7-2003-x86_64 --kickstart=/var/lib/cobbler/kickstarts/CentOS7.ks
5. （可选）修改安装系统的内核参数，使得网卡名变为 eth
   1. cobbler profile edit --name=CentOS7-2003-x86_64 --kopts='net.ifnames=0 biosdevname=0'

可以看到下面 Kickstart 那里的配 置 cfg 文件地址被改变了\[root@cobbler kickstarts]# cobbler profile report --name=CentOS7-2003-x86_64 Name : CentOS7-2003-x86_64 TFTP Boot Files : {}Comment : DHCP Tag : defaultDistribution : CentOS7-2003-x86_64 Enable gPXE? : 0Enable PXE Menu? : 1Fetchable Files : {}Kernel Options : {}Kernel Options (Post Install) : {}Kickstart : /var/lib/cobbler/kickstarts/CentOS7-2003-x86_64.cfgKickstart Metadata : {}Management Classes : \[]Management Parameters : <>Name Servers : \[]Name Servers Search Path : \[]Owners : \['admin']Parent Profile : Internal proxy : Red Hat Management Key : <>Red Hat Management Server : <>Repos : \[]Server Override : <>Template Files : {}Virt Auto Boot : 1Virt Bridge : xenbr0Virt CPUs : 1Virt Disk Driver Type : rawVirt File Size(GB) : 5Virt Path : Virt RAM (MB) : 512Virt Type : kvm#同步下 cobbler 数据，每次修改完都要镜像同步

## 开始批量安装系统

cobbler 已经配置完成，开机后，使用 PXE 启动后，如果 DHCP 运行正常，则可看到如下界面

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125382981-7146a22b-241d-4c33-96f5-7167cabf7bf4.jpeg)

这里如果方向键不选择第二项就无法装机，不算自动化，我们需要进行手动指定才可以完全自动化

## 全自动化安装

上面可以看到在开机之后，还需要手动选择一下，还不够完全自动，如果想要开机之后，不用任何操作即可安装系统，那么可以使用 cobbler system 命令，通过 mac 来区分设备，当 cobbler 识别到 mac 时，就会自动选择 profiles 进行安装。

VMware 查看物理 MAC 地址

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383095-bb54ec6c-c5ff-49ea-a276-16b88900e2d7.jpeg)

1. cobbler system add --name=test --mac=52:54:00:be:29:97 --profile=CentOS8-2004-x86_64 --ip-address=10.10.100.230 --subnet=255.255.255.0 --gateway=10.10.100.254 --interface=eth0 --static=1 --hostname=cobbler-system-test --name-servers="223.5.5.5"

Note：

- 此功能需要在 ks 文件中添加如下内容，否则在自动安装中，读取 ks 文件的时候会报错：unable to open input kickstart file:error opening file:No such file or directory: '/tmp/pre_install_network_config'


    ...
    %pre
    $SNIPPET('pre_install_network_config')
    %end
    ... # 后面是%packages的配置

- 如果有多块网卡的话，只自定义第二块网卡为外网，第一块网卡不一起配置，就会还是出现上文。多网卡像这样指定 IP 和网卡


    [root@cobbler ~]# cobbler system add --name=test --mac=52:54:00:be:29:97  --profile=CentOS8-2004-x86_64 --ip-address=10.0.0.82 --subnet=255.255.255.0 --interface=eth0 --static=1
    [root@cobbler ~]# cobbler system edit --name=test --mac=52:54:00:04:89:ee  --profile=CentOS8-2004-x86_64 --ip-address=192.168.31.82 --subnet=255.255.255.0 --gateway=192.168.31.1 --interface=eth1 --static=1 --hostname=zabbix --name-servers="223.5.5.5"

- 其实 interface 的值可以随便写，因为 cobbler 是根据 mac 选择的，interface 的值没啥意义

# Cobbler 配置

/etc/cobbler/\* # 基础配置文件路径

1. ,/setting #cobbler 服务的配置文件

/var/lib/cobbler/\* # cobbler 数据保存路径

1. ./kickstarts/\* # 从该路径下读取 kickstarts 文件。
2. ./snippets/\* #kickstart 文件中$SNIPPET 变量内所用到的文件

/var/www/cobbler/\* # cobble 数据保存路径

1. ,/ks_mirror/\* #镜像存放目录，cobbler 会将镜像中的所有安装文件拷贝到本地一份，放在该目录下，该目录下的目录名是以 import 时指定的 name 来命名的。因此/var/www/cobbler 目录必须具有足够容纳安装文件的空间

# Cobbler 的 Web 管理界面的安装与配置

新版 Cobbler 的 Web 界面使用的是 https

如果 web 界面报错，详见该文章的解决方式<http://www.trgeek.com/linux/2019/01/94.html>，如果 google 浏览器打不开就用火狐，证书问题

登录 URL: <https://10.10.100.250/cobbler_web>

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383019-05bbc151-eeb0-4b67-b89b-e089ef3df4ef.jpeg)

默认用户名：cobbler

默认密码 ：cobbler

/etc/cobbler/users.conf # Web 服务授权配置文件

/etc/cobbler/users.digest # 用于 web 访问的用户名密码配置文件

\#在 Cobbler 组添加 ren 用户，提示输入 2 遍密码确认

\[root@cobbler kickstarts]# htdigest /etc/cobbler/users.digest "Cobbler" ren

Adding user ren in realm Cobbler

New password:

Re-type new password:

同步下

\[root@cobbler kickstarts]# cobbler sync

# Cobbler 的跨网段部署

如果想实现跨网段部署，需要注意以下几点

1. 修改/etc/cobbler/dhcp.template 中的配置信息，想给哪个网段部署系统，就加上该网段的配置。配置如下面的
2. 在没有 DHCP 服务的网段中的设备上开启 dhcp 中继服务，需要安装 DHCP 包。
3. i.e.使用命令`dhcrelay IP`(IP 为 DHCP 服务所在设备的 IP)来开启 DHCP 代理服务，这样在与 cobbler 不同网段的设备，也可以实现自动部署系统了。dhcrelay 命令还可以使用`-i`选项指定网卡名，来确定代理从哪个网卡进来的 dhcp 流量
4. 还可以直接修改 dhcrelay 服务的配置文件，e.g.`sed -i '/ExecStart=/s/$/& 192.168.20.12/g' /etc/systemd/system/dhcrelay.service`
5. e.g.Cobbler 所在网段为 192.168.1.0/24，现在想给 192.168.2.0/24 网段的其中一台设备安装系统，这时候就需要在 192.168.2.0/24 网段那种找一台系统，开启 dhcp 中继

e.g./etc/cobbler/dhcp.template

注意：

1. 需要注释掉 cobbler 与 dhcp 设备所在网段的网关配置，否则部署不同网段会失败，原因未知，待验证


    ddns-update-style interim;

    allow booting;
    allow bootp;

    ignore client-updates;
    set vendorclass = option vendor-class-identifier;

    option pxe-system-type code 93 = unsigned integer 16;

    subnet 192.168.20.0 netmask 255.255.255.0 {
    #     option routers             192.168.20.254;
         option domain-name-servers 114.114.114.114;
         option subnet-mask         255.255.255.0;
         range dynamic-bootp        192.168.20.100 192.168.20.254;
         default-lease-time         21600;
         max-lease-time             43200;
         next-server                $next_server;
         class "pxeclients" {
              match if substring (option vendor-class-identifier, 0, 9) = "PXEClient";
              if option pxe-system-type = 00:02 {
                      filename "ia64/elilo.efi";
              } else if option pxe-system-type = 00:06 {
                      filename "grub/grub-x86.efi";
              } else if option pxe-system-type = 00:07 {
                      filename "grub/grub-x86_64.efi";
              } else if option pxe-system-type = 00:09 {
                      filename "grub/grub-x86_64.efi";
              } else {
                      filename "pxelinux.0";
              }
         }
    }

    subnet 192.168.30.0 netmask 255.255.255.0 {
    option domain-name-servers 114.114.114.114;
    option routers 192.168.30.254;
    range dynamic-bootp 192.168.30.100 192.168.30.250;
    option subnet-mask 255.255.255.0;
    next-server $next_server;
    default-lease-time 21600;
    max-lease-time 43200;
    }

# 附录

添加 repos 源

    1.添加repo源
    #举个栗子，centos7.2版本的openstack的repo源
    [root@cobbler02 ~]# cobbler repo add --name=centos7.2-openstack-mitaka --mirror=http://mirrors.aliyun.com/centos/7.2.1511/cloud/x86_64/openstack-mitaka/ --arch=x86_64 --breed=yum
    #添加repo源，举个栗子，centos7版本的epel源
    [root@cobbler02 ~]# cobbler repo add --name=centos7-x86_64-epel --mirror=http://mirrors.aliyun.com/epel/7/x86_64/ --arch=x86_64 --breed=yum
    2.同步repo
    [root@cobbler02 ~]# cobbler reposync
    3.添加repo到对应的profile
    cobbler profile edit --name=Centos-7-x86_64  --repos=
    4.修改kickstart文件。添加。（些到%post %end中间）
    %post
    systemctl disable postfix.service

    $yum_config_stanza
    %end

    5.添加定时任务，定期同步repo

客户机重装系统教程：1

centos7 系列：

1\)不指定详细 system 模板，让 cobbler 自己装一台 centos7 的镜像

    #注意上面如果没有指定epel源的是无法装koan包的
    #执行这个命令：wget -O /etc/yum.repos.d/epel.repo http://mirrors.aliyun.com/repo/epel-7.repo
    [root@zabbix ~]# yum install -y koan
    #指定cobbler的服务器
    [root@zabbix ~]# koan --server=192.168.31.73 --list=profiles
    - looking for Cobbler at http://192.168.31.73:80/cobbler_api
    CentOS-7.2-x86_64
    #指定从哪个镜像进行安装
    [root@zabbix ~]# koan --replace-self --server=192.168.31.73 --profile=CentOS-7.2-x86_64
    - looking for Cobbler at http://192.168.31.73:80/cobbler_api
    - reading URL: http://192.168.31.73/cblr/svc/op/ks/profile/CentOS-7.2-x86_64
    install_tree: http://192.168.31.73/cblr/links/CentOS-7.2-x86_64
    downloading initrd initrd.img to /boot/initrd.img_koan
    url=http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/initrd.img
    - reading URL: http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/initrd.img
    downloading kernel vmlinuz to /boot/vmlinuz_koan
    url=http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/vmlinuz
    - reading URL: http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/vmlinuz
    - ['/sbin/grubby', '--add-kernel', '/boot/vmlinuz_koan', '--initrd', '/boot/initrd.img_koan', '--args', '"ksdevice=link lang= text net.ifnames=0 ks=http://192.168.31.73/cblr/svc/op/ks/profile/CentOS-7.2-x86_64 biosdevname=0 kssendmac "', '--copy-default', '--make-default', '--title=kick1464687061']
    - ['/sbin/grubby', '--update-kernel', '/boot/vmlinuz_koan', '--remove-args=root']
    - reboot to apply changes
    #重启系统
    [root@zabbix ~]# reboot

2\)指定设定好的系统配置，让系统生成一个指定的 mac 地址绑定的 ip 和其他你指定的东西

    [root@zabbix ~]# yum install -y koan
    #指定cobbler的服务器选择system模板
    [root@MiWiFi-R2D-srv ~]# koan --server=192.168.31.73 --list=system
    - looking for Cobbler at http://192.168.31.73:80/cobbler_api
    koan does not know how
    [root@MiWiFi-R2D-srv ~]# koan --replace-self --server=192.168.31.73 --system=zabbix02
    - looking for Cobbler at http://192.168.31.73:80/cobbler_api
    - reading URL: http://192.168.31.73/cblr/svc/op/ks/system/zabbix02
    install_tree: http://192.168.31.73/cblr/links/CentOS-7.2-x86_64
    downloading initrd initrd.img to /boot/initrd.img_koan
    url=http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/initrd.img
    - reading URL: http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/initrd.img
    downloading kernel vmlinuz to /boot/vmlinuz_koan
    url=http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/vmlinuz
    - reading URL: http://192.168.31.73/cobbler/p_w_picpaths/CentOS-7.2-x86_64/vmlinuz
    - ['/sbin/grubby', '--add-kernel', '/boot/vmlinuz_koan', '--initrd', '/boot/initrd.img_koan', '--args', '"ksdevice=link lang= text net.ifnames=0 ks=http://192.168.31.73/cblr/svc/op/ks/system/zabbix02 biosdevname=0 kssendmac "', '--copy-default', '--make-default', '--title=kick1464688081']
    - ['/sbin/grubby', '--update-kernel', '/boot/vmlinuz_koan', '--remove-args=root']
    - reboot to apply changes to list that
    [root@MiWiFi-R2D-srv ~]# reboot

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383019-fc08ace0-55f9-4c41-873e-43b8a4aeae90.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383004-27cf30e1-7d74-49f2-9840-c7f2fce73950.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383082-57129950-42f7-4ea1-8837-579b371c54e0.jpeg)

已经按照指定的 system 配置安装好了，但是这里出现了一个问题，双网卡的话，两个网段，现在的路由都是同一个了，所以上不了网了，求高手解决

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wf5w0g/1616125383011-006ecd1b-39f8-40de-8836-92a1b578d936.jpeg)
