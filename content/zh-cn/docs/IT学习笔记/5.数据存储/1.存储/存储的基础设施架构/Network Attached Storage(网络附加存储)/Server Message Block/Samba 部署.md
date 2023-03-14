---
title: Samba 部署
---

# [Linux 下部署 Samba 服务环境的操作记录](https://www.cnblogs.com/kevingrace/p/8550810.html)

关于 Linux 和 Windows 系统之间的文件传输，很多人选择使用 FTP，相对较安全，但是有时还是会出现一些问题，比如上传文件时，文件名莫名出现乱码，文件大小改变等问题。相比较来说，使用 Samba 作为文件共享，就省事简洁多了。Samba 服务器通信协议（Server Messages Block）就是是为了解决局域网内的文件或打印机等资源的共享服务问题，让多个主机之间共享文件变成越来越简单。下面简单介绍下，在 Centos7 下部署 Samba 服务的操作记录（测试机 192.168.10.204）：

## 1）安装 Samba

    [root@samba-server ~]# cat /etc/redhat-release
    CentOS Linux release 7.4.1708 (Core)
    [root@samba-server ~]# rpm -qa|grep samba
    [root@samba-server ~]# yum install -y samba

## 2）安全角度考虑，需要设置防火墙策略（不要关闭防火墙）

    添加samba服务到防火墙策略中
    [root@samba-server ~]# firewall-cmd --add-service samba --permanent
    success
      
    重启防火墙
    [root@samba-server ~]# firewall-cmd --reload
    success
      
    查看samba服务是否添加到防火墙中：
    [root@samba-server ~]# firewall-cmd --list-all|grep samba
      services: ssh dhcpv6-client samba
       
    记住：一定要关闭selinux（否则会造成windows客户机连接Samba失败）
    [root@samba-server ~]# vim /etc/sysconfig/selinux
    .....
    SELINUX=disabled
     
    [root@samba-server kevin]# setenforce 0
    [root@samba-server kevin]# getenforce
    Permissive

## 3）配置 Samba 服务文件

    [root@samba-server ~]# cp /etc/samba/smb.conf /etc/samba/smb.conf.bak
    [root@samba-server ~]# vim /etc/samba/smb.conf
    # See smb.conf.example for a more detailed config file or
    # read the smb.conf manpage.
    # Run 'testparm' to verify the config is correct after
    # you modified it.
     
    [global]                                                   //全局配置
         workgroup = SAMBA
         security = user
     
         passdb backend = tdbsam
     
         printing = cups
         printcap name = cups
         load printers = yes
         cups options = raw
     
    [homes]
         comment = Home Directories
         valid users = %S, %D%w%S
         browseable = No
         read only = No
         inherit acls = Yes
     
    [printers]                                                 //共享打印机配置
         comment = All Printers
         path = /var/tmp
         printable = Yes
         create mask = 0600
         browseable = No
     
    [print$]
         comment = Printer Drivers
         path = /var/lib/samba/drivers
         write list = root
         create mask = 0664
         directory mask = 0775
     
    [kevin]                                                    //这个是共享文件夹标识，表示登录samba打开时显示的文件夹名称。配置了多少个共享文件夹标识，登录samba时就会显示多少文件夹。
           comment = please do not modify it all will          //comment是对该共享的描述，可以是任意字符串
           path= /home/kevin                                   //共享的路径
           writable = yes                                      //是否写入
           public = no                                         //是否公开

## 4）添加 kevin 账号（如上配置中添加的内容）

    设置为不予许登入系统,且用户的家目录为 /home/kevin（相当于虚拟账号）的kevin账号。
    [root@samba-server ~]# useradd -d /home/kevin -s /sbin/nologin kevin

## 5）pdbedit 命令说明

    pdbedit 命令用于管理Samba服务的帐户信息数据库，格式为："pdbedit [选项] 帐户"
    第一次把用户信息写入到数据库时需要使用-a参数，以后修改用户密码、删除用户等等操作就不再需要了。
     
    pdbedit -L ：查看samba用户
    pdbedit -a -u user：添加samba用户
    pdbedit -r -u user：修改samba用户信息
    pdbedit -x -u user：删除samba用户
     
    samba服务数据库的密码也可以用 smbpasswd 命令 操作
    smbpasswd -a user：添加一个samba用户
    smbpasswd -d user：禁用一个samba用户
    smbpasswd -e user：恢复一个samba用户
    smbpasswd -x user：删除一个samba用户

## 6）将 kevin 添加为 samba 用户

    [root@samba-server ~]# id kevin
    uid=1001(kevin) gid=1001(kevin) groups=1001(kevin)
     
    [root@samba-server ~]# pdbedit -a -u kevin
    new password:                              //设置kevin使用的samba账号密码，比如123456
    retype new password:                       //确认密码
    Unix username:        kevin
    NT username:         
    Account Flags:        [U          ]
    User SID:             S-1-5-21-33923925-2092173964-3757452328-1000
    Primary Group SID:    S-1-5-21-33923925-2092173964-3757452328-513
    Full Name:           
    Home Directory:       \\samba-server\kevin
    HomeDir Drive:       
    Logon Script:        
    Profile Path:         \\samba-server\kevin\profile
    Domain:               SAMBA-SERVER
    Account desc:        
    Workstations:        
    Munged dial:         
    Logon time:           0
    Logoff time:          Wed, 06 Feb 2036 23:06:39 CST
    Kickoff time:         Wed, 06 Feb 2036 23:06:39 CST
    Password last set:    Mon, 12 Mar 2018 18:07:58 CST
    Password can change:  Mon, 12 Mar 2018 18:07:58 CST
    Password must change: never
    Last bad password   : 0
    Bad password count  : 0
    Logon hours         : FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
     
    接着修改samba用户的家目录权限
    [root@samba-server ~]# chown -Rf kevin.kevin /home/kevin

## 7）启动 Samba 服务

    [root@samba-server ~]# systemctl start smb
    [root@samba-server ~]# systemctl enable smb
    Created symlink from /etc/systemd/system/multi-user.target.wants/smb.service to /usr/lib/systemd/system/smb.service.
    [root@samba-server ~]# systemctl restart smb
    [root@samba-server ~]# systemctl status smb
    ● smb.service - Samba SMB Daemon
       Loaded: loaded (/usr/lib/systemd/system/smb.service; enabled; vendor preset: disabled)
       Active: active (running) since Mon 2018-03-12 18:11:20 CST; 3s ago
     Main PID: 977 (smbd)
       Status: "smbd: ready to serve connections..."
       CGroup: /system.slice/smb.service
               ├─977 /usr/sbin/smbd
               ├─978 /usr/sbin/smbd
               ├─979 /usr/sbin/smbd
               └─980 /usr/sbin/smbd
     
    Mar 12 18:11:19 samba-server systemd[1]: Starting Samba SMB Daemon...
    Mar 12 18:11:19 samba-server systemd[1]: smb.service: Supervising process 977 which is not our child. We'll most likely not... exits.
    Mar 12 18:11:20 samba-server smbd[977]: [2018/03/12 18:11:20.065982,  0] ../lib/util/become_daemon.c:124(daemon_ready)
    Mar 12 18:11:20 samba-server systemd[1]: Started Samba SMB Daemon.
    Mar 12 18:11:20 samba-server smbd[977]:   STATUS=daemon 'smbd' finished starting up and ready to serve connections
    Hint: Some lines were ellipsized, use -l to show in full.

## 8）开始测试

先往共享路径 / home/kevin 里添加点内容

    [root@samba-server kevin]# touch test1 test2 test3
    [root@samba-server kevin]# mkdir a1 a2 a3
    [root@samba-server kevin]# ls
    a1  a2  a3  test1  test2  test3

接着再 windos 客户机本地测试。”Win+E 键 "打开，在最上面的" 网络 " 地址栏输入 “\192.168.10.204”，然后回车，输入上面设置的 samba 账号 kevin 及其密码，就能共享到 linux 上的 / home/kevin 下的文件了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221290-ef6820a7-3ebd-469e-bf21-e6d8e84b757d.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221300-59931478-7448-43c1-81a9-ebe72272bdf7.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221336-d1f24efb-36e2-4f25-901c-bcab699c2c72.png)

连接上后，就可以在 windows 和 linux 直接进行文件夹的共享操作了，可以让里面放点测试文件

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221396-c2e1f154-1f36-4b92-ba6e-1e7ef93b2d29.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221323-a2bdb63b-f19a-40e2-bced-80e6216090f5.png)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/fqfqm4/1619091221293-e705eb86-3b36-45a5-ac92-39daaa0d8677.png)

**如果在 windows 客户机上连接 Samba 出现网络错误：Windows 无法访问\192.168.10.204\kevin**，解决办法如下：

    查看上下文的安全关系
    [root@samba-server ~]# semanage kevin -a -t samba_share_t /home/kevin/
    -bash: semanage: command not found
     
    如果系统出现上面的报错 ，说明你系统里没有安装 semanage命令，下面开始安装semanage：
     
    [root@samba-server ~]# yum provides /usr/sbin/semanage
    Loaded plugins: fastestmirror
    Loading mirror speeds from cached hostfile
     * base: mirror.0x.sg
     * epel: mirror.dmmlabs.jp
     * extras: mirror.0x.sg
     * updates: mirror.0x.sg
    policycoreutils-python-2.5-17.1.el7.x86_64 : SELinux policy core python utilities        //这个是安装包
    Repo        : base
    Matched from:
    Filename    :
     
    [root@samba-server ~]# yum install -y policycoreutils-python
     
    然后再执行一次，执行完成后，不要忘了刷新上下文关系
    [root@samba-server ~]# semanage fcontext -a -t samba_share_t /home/kevin
    [root@samba-server ~]# restorecon -Rv /home/kevin
     
    允许SElinux对于SMB用户共享家目录的布尔值
     
    重启Samba
    [root@samba-server ~]# systemctl restart smb

### 如何在 windows 本机访问 samba 时用切换另一个用户登录

    方法如下：
    1）按键ctrl+r，打开"运行"，输入"cmd"
    2）输入命令"net use * /delete"，接着输入"Y",即先取消所有的net 连接
    3）输入切换账号的命令"net use \\192.168.10.204\IPC$ grace@123 /user:grace"，即表示切换到grace账号（密码为grace@123）

### 重置 samba 账号密码

    [root@samba-server ~]# smbpasswd -a kevin     //即重置kevin密码

**======================================================**

### 清理 windows 下连接 linux 的 samba 服务缓存

在安装配置 linux 服务器 samba 服务之初，samba 服务难免会经过多次修改配置 / 重启，在期间 win 的系统或许早已连接上了 samba。samba 修改配置修改后，特别是用户权限，再次从 win 登录就很容易出现因缓存的权限原因导致不允许访问或者操作。
这时一般要等很久会清理缓存，另外重启 win 也会清理。但这效率很低。用以下手动的方法可以实时清理。

清理步骤：
1）打开 win 的命令行（ctrl+R，输入 cmd）。
2）在命令行里输入**net use**，就会打印出当前缓存的连接上列表。
3\) 根据列表，一个个删除连接： **net use 远程连接名称 /del**；
或者一次性全部删除：**net use \* /del**。

这样再次命令行输入 samba 服务地址的时候，就会重新让你输入访问的账户和密码了。

**======================================================**

### 可以在一个 samba 环境下建立多个业务组的共享目录

    比如：
    创建一个运维部门的samba共享磁盘，可以看到所有的共享内容；
    创建一个产品风控组的samba共享磁盘，只能看到自己组的共享内容；
     
    [root@samba ~]# cd /etc/samba/
    [root@samba samba]# ls
    lmhosts  ops.smb.conf  smb.conf  smb.conf.bak  smbusers  chanpinfengkong.smb.conf
    [root@samba samba]# diff smb.conf smb.conf.bak
    103d102
    <         config file = /etc/samba/%U.smb.conf     #使用config file时，当用户访问Samba服务器，只能看到自己，其他在smb.conf中定义的共享资源都无法看到。
     
    [root@samba samba]# cat ops.smb.conf
    [信息科技部-运维小窝]                                                 
           comment = please do not modify it all will       
           path= /data/samba                                                                
           public = no        
           valid users = wangshibo,linan,@samba
           printable = no
           write list = @samba
     
    [root@samba samba]# cat chanpinfengkong.smb.conf
    [产品风控组共享目录]                                                 
           comment = please do not modify it all will       
           path= /data/samba/产品风控组                                                           
           public = no        
           valid users = xiaomin,haokun,@samba
           printable = no
           write list = @samba
     
     
    useradd创建以上的几个用户，并设置好用户家目录
    [root@samba ~]# useradd wangshibo -d /data/samba -s /sbin/nologin
    [root@samba ~]# useradd linan -d /data/samba -s /sbin/nologin
    [root@samba ~]# useradd xiaomin -d /data/samba/产品风控组 -s /sbin/nologin
    [root@samba ~]# useradd haokun -d /data/samba/产品风控组 -s /sbin/nologin
    [root@samba ~]# cat /etc/passwd
    ......
    wangshibo:x:507:508::/data/samba:/sbin/nologin
    lijinhe:x:508:509::/data/samba:/sbin/nologin
    ......
    xiaomin:x:1006:1006::/data/samba/产品风控组:/sbin/nologin
    haokun:x:1007:1007::/data/samba/产品风控组:/sbin/nologin
    chanpinfengkong:x:1010:1010::/home/chanpinfengkong:/bin/bash
     
    将这几个用户添加到samba里
    [root@samba ~]# pdbedit -a -u wangshibo
    [root@samba ~]# pdbedit -a -u linan
    [root@samba ~]# pdbedit -a -u xiaomin
    [root@samba ~]# pdbedit -a -u haokun
     
    [root@samba ~]# pdbedit -L
    wangshibo:507:
    linan:510:
    xiaomin:1006:
    haokun:1007:
     
    创建chanpinfengkong组，将xiaomin和haokun添加到这个组内
    [root@samba ~]# useradd chanpinfengkong
    [root@samba ~]# usermod -G chanpinfengkong xiaomin
    [root@samba ~]# usermod -G chanpinfengkong haokun
     
    创建samba共享目录
    [root@samba ~]# cd /data/
    [root@samba data]# mkdir samba
    [root@samba data]# mkdir samba/产品风控组
    [root@samba data]# chown -R samba.samba samba
    [root@samba data]# chmod -R 777 samba
    [root@samba data]# setfacl -R -m g:chanpinfengkong:rwx samba/产品风控组
     
    赋权脚本
    [root@samba ~]# cat /opt/samba.sh
    #!/bin/bash
     
    while [ "1" = "1" ]
    do
       /bin/chmod -R 777 /data/samba
       /usr/bin/setfacl -R -m g:chanpinfengkong:rwx /data/samba/产品风控组
    done
     
    [root@samba ~]# nohup sh -x /opt/samba.sh &
    [root@samba ~]# ps -ef|grep samba.sh
    root      62836      1 16 May09 ?        14-23:47:39 sh -x /opt/samba.sh
    root     185455 117471  0 15:41 pts/2    00:00:00 grep samba.sh
     
    如上配置后，登录samba：
    1）用wangshibo，linan账号登录samba，能看到"/data/samba"下面所有的共享内容。
    2）用xiaomin,haokun账号登录samba，只能看到"/data/samba/产品风控组" 下面的共享内容
    3）如果还需要分更多的组，就如上面的"产品风控组"一样进行配置即可！

**\*\***\*\*\***\*\*** 当你发现自己的才华撑不起野心时，就请安静下来学习吧！**\*\***\*\*\***\*\***
