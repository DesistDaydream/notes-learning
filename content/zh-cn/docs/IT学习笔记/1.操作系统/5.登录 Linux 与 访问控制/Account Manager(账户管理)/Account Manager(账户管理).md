---
title: Account Manager(账户管理)
weight: 1
---

# 概述

> 参考：
>
> - [红帽官方文档,RedHat7-管理用户账户的基础知识](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-getting_started#sec-Managing_User_Accounts)
> - [红帽官方文档,RedHat7-系统管理员指南-第四章-管理用户和组](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-managing_users_and_groups)
> - [红帽官方文档,RedHat7-安全指南](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/security_guide/index)

Linux 操作系统是一个多用户操作系统，所以除了 **Terminal(终端)** 以外，还需 **Account(账户)** 才可以登录上去，Linux 操作系统允许多个用户访问安装在一台机器上的单个系统。每个 User(用户) 都在自己的 Account(账户) 下操作。因此，Account Manager 代表了 Linux 系统管理的核心要素。

> User 与 Account 傻傻分不清楚，在 Linux 操作系统中，通常都会将 Account 称为 User，但是，这个称呼实际上并不准确。因为一个 User，比如 root，可以被多个现实世界中多个人使用，那么 root 这个 User 就会对应多个真实的 User~这种描述是非常矛盾的~~~~
> 只不过，随着时间的推移，人们慢慢叫习惯了，各种文档和源码也都一直使用 UID 这种名称，也就不再特别区分 Account 和 User 了。只需要知道，UID 更准确的描述应该是 AID。

同时，一个 Linux Account 也可以不代表一个真实的 User，这样的 Account 只被应用程序使用，一个应用程序使用某个 Account 运行，以便让系统更方便得对程序进行精细化控制。这种控制方式称为 **Access Control(访问控制)**，所以，从这种角度看，Account 也可以称为 **Role(角色)**，详见 [访问控制](/docs/IT学习笔记/1.操作系统/5.登录%20Linux%20 与%20 访问控制/Access%20Control(访问控制).md Control(访问控制).md) 章节。

为了方便得对多个 Account 管理，可以将多个 Account 组合起来，称为 **Group(组)**，一个 Group 就是一个或多个 Account 的集合。

通常，Linux 将账户分为两类

- Normal Accounts(普通账户)
- System Accounts(系统账户)

每个账户都有一个对应的 UID 作为其唯一标识符(纯数字)。同样，每个组也有一个对应的 GID 作为其唯一标识符(纯数字)。通常来说：

- 1000 以下是系统账户与保留账户 和 系统组与保留组
- 1000 以上是普通账户和组

每当我们使用 useradd 命令新建一个普通用户时，用户的 UID 都是 1000 之后的数字，这种行为可以通过修改 /etc/login.defs 文件中的 UID_MIN、GID_MIN 等参数来改变。

## Password(密码)

**Password(密码)** 是用来验证用户身份的最主要方法。当用户使用一个账户登录 Linux 操作系统时，密码是用来证明账户属于该用户的一种非常高效的方式。

Linux 系统使用 **Secure Hash Algorithm 512(SHA512)**和 **shadow passwords**。默认情况下，账户信息保存在 /etc/passwd 文件中，对应的密码信息经过哈希后保存在 /etc/shadow 文件中。

# 关联文件

**/etc/group** # 账户组信息
**/etc/passwd** # 账户信息
**/etc/shadow** # 安全账户信息
**/etc/login.defs** # login 工具包中的配置文件，部分账户管理工具会读取该文件中的参数
**/etc/pam.d/** #

- ./chfn
- ./chpasswd
- ./chsh
- ./newusers
- ./passwd

**/etc/skel/** # 该目录为账户目录模板。该目录下包含多个隐藏的文件，当创建用户时，会拷贝该目录下的所有文件到所创建用户的家目录中
**/home/AccountName/** # UserName 为该账户同名的家目录
**/var/spool/mail/AccountName**# 该文件为该账户的邮件池
注意：

- 若 /etc/shadow 被 selinux 所管理，有的时候密码修改将会失败，报错 `passwd: Authentication token manipulation error`

# 账户管理工具

Linux 系统的账户管理功能，通常由 **shadow-utils 包** 或 **passwd 包** 中的各种工具和库提供。

> 在有的发行版中(比如 CentOS)，只会将 passwd 包中的 passwd 程序保留，而将其余的程序，放在名为 **shadow-utils** 的包中。

不同的 Linux 发型，还会用到某些个别的包与主包配合提供完整的账户管理功能：

- **base-passwd** # 这是 Ubuntu 发型版中独有的包。这个包中包含一个 `update-passwd` 的程序，将会根据 /var/lib/dpkg/info/base-passwd.preinst 脚本生成 /etc/passwd 和 /etc/group 文件
  - 参考：<https://unix.stackexchange.com/questions/470126/how-is-the-etc-passwd-file-instantiated>
- **setup** # 这是 CentOS 发行版中独有的包。包含了一组重要的系统配置文件和安装文件，例如 /etc/passwd、/etc/group、/etc/shadow、/etc/profile 等等

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/oib9pm/1635476577308-bd0e79ab-ffc9-41f8-ba65-471f0a3e2918.png)

## useradd # 添加用户

> 参考：
>
> - [Manual(手册),useradd(8)](https://man7.org/linux/man-pages/man8/useradd.8.html)

useradd 程序在添加用户时，会读取 /etc/login.defs 文件中的 PASS_MAX_DAYS、PASS_MIN_DAYS、PASS_WARN_AGE 等参数，并将参数的值写入到 /etc/shadow 文件中的对应字段

### Syntax(语法)

**useradd \[OPTIONS] NAME**

OPTIONS：

- **-m** # 自动建立用户的登入目录。
- **-u** # 指定用户 ID 号。该值在系统中必须是唯一的。0~499 默认是保留给系统用户账号使用的，所以该值必须大于 499。
- **-g GID** # 指定该用户的基本组 ID
- **-c** # 加上备注文字，备注文字保存在 passwd 的备注栏中。
- -**d** # 指定该用户的家目录，默认为 /home 目录下的与用户名同名的文件夹
- **-e** # 指定账号的失效日期，日期格式为 MM/DD/YY，例如 06/30/12。缺省表示永久有效。
- -f # 指定在密码过期后多少天即关闭该账号。如果为 0 账号立即被停用；如果为-1 则账号一直可用。默认值为-1.
- **-G, --groups \<GROUPS>** # 指定用户所属的附加群组。
- -l # 不要把用户添加到 lastlog 和 failog 中, 这个用户的登录记录不需要记载
- -M # 不要自动建立用户的登入目录。
- -n # 取消建立以用户名称为名的群组。
- -p # PASSWORD 指定新用户的密码
- -r # 建立一个系统帐号
- -s # 指定用户登入后所使用的 shell。默认值为/bin/bash。

EXAMPLE：

- 创建一个名为 lichenhao 的用户，并设置该用户密码为 lch@password
  - useradd -m lichenhao && echo 'lch@password' | passwd lichenhao --stdin
- 新增一个用户 user，并指定用户组 ftp
  - useradd -g ftp user
- 新增一个用户：user 并指定家目录为/mnt/bak/。如果没有此目录，则报错，就需要手动创建目录
  - useradd -d /mnt/back/ user
- 新增一个用户：user 并设置密码为 123456
  - useradd -p 123456 user
- 新增一个 FTP 用户：ftp2018 （无需登录系统）
  - useradd -g ftp -s /sbin/nologin ftp2018

## userdel # 删除用户

### Syntax(语法)

**userdel \[OPTIONS] NAME**

OPTIONS：

- **-f, --force** # 强制删除
- **-r, --remove** # 删除用户的时候同时移除该用户的家目录和邮件池。不加参数的话，只会删除用户，但是不会删除在/home 目录中的用户家目录。如果你想要连此用户的家目录也一并删除，可以加上 –remove-home 这个参数
- -**R, --root CHROOT_DIR** # chroot 到的目录
- **-Z, --selinux-user** # 为用户删除所有的 SELinux 用户映射

EXAMPLE

- userdel -r user #删除名为 user 的用户，同时删除该用户的家目录和邮件池文件

## usermod # 修改用户

### Syntax(语法)

**usermod \[OPTIONS] NAME**

OPTIONS：

- **-a, --append GROUP** # 将用户追加至上边 -G 中指定附加组中，并不从其它组中删除此用户
- **-c, --comment \<STRING>** # GECOS 字段的新值
- **-d, --home \<HOME_DIR>** # 用户的新主目录
- **-e, --expiredate \<EXPIRE_DATE>** # 设定帐户过期的日期为 EXPIRE_DATE
- **-f, --inactive INACTIVE** # 过期 INACTIVE 天数后，设定密码为失效状态
- **-g, --gid GROUP** # 强制使用 GROUP 为新主组
- **-G, --groups GROUPS** # 新的附加组列表 GROUPS。将用户从该选项指定的组列表以外的组中删除。可以与 -a 选项配合，变更此行为。-a 会将用户附加到指定的组中而不删除。
- **-l, --login LOGIN** # 新的登录名称
- **-L, --lock** # 锁定用户帐号
- **-m, --move-home** # 将家目录内容移至新位置 (仅于 -d 一起使用)
- **-o, --non-unique** # 允许使用重复的(非唯一的) UID
- **-p, --password PASSWORD** # 将加密过的密码 (PASSWORD) 设为新密码
- -**R, --root CHROOT_DIR** # chroot 到的目录
- **-s, --shell SHELL** # 该用户帐号的新登录 shell
- **-u, --uid UID** # 用户帐号的新 UID
- **-U, --unlock** # 解锁用户帐号
- **-Z, --selinux-user SEUSER** # 用户账户的新 SELinux 用户映射

EXAMPLE

- 修改 newname 用户所在群组为 test
  - usermod -g test newname
- 一次将一个用户添加到多个群组
  - usermod -G friends,happy,funny newname
  - 注意：使用 usermod 时要小心，因为配合-g 或-G 参数时，它会把用户从原先的群组里剔除，加入到新的群组。如果你不想离开原先的群组，又想加入新的群组，可以在-G 参数的基础上加上-a 参数，a 是英语 append 的缩写，表示“追加”。
- -a 追加用户到新的用户组，保留原来的组
  - usermod -aG happy newname

## groupadd、groupdel、groupmod、gpasswd # 用户组管理相关命令

OPTIONS：

- -f, --force 如果组已经存在则成功退出，并且如果 GID 已经存在则取消 -g
- -g, --gid GID # 为新组使用 GID
- -K, --key KEY=VALUE # 不使用 /etc/login.defs 中的默认值
- -o, --non-unique # 允许创建有重复 GID 的组
- -p, --password PASSWORD # 为新组使用此加密过的密码
- -r, --system # 创建一个系统账户
- -R, --root CHROOT_DIR # chroot 到的目录

EXAMPLE

- 创建一个名为 newname 的组
  - groupadd newname
- 修改组
  - groupmod -n test2group testgroup
- 删除名为 test2group 的组
  - groupdel test2group
- 查看当前登陆用户所在的组
  - groups
- 查看 testnewuser 所在的组
  - groups testnewuser

## who # 显示当前登录用户的相关信息

### Syntax(语法)

**who \[OPTION]... \[ FILE | ARG1 ARG2 ]**

OPTIONS

- -a 打印能打印的全部
- -d 打印死掉的进程
- -m 同 am i,mom likes
- -q 打印当前登录用户数及用户名
- -u 打印当前登录用户登录信息
- -r 打印运行等级

EXAMPLE

- whoami # 要查看当前登录用户的用户名
- who am i # 表示打开当前伪终端的用户的用户名，可以简写为 who

who

# 密码管理工具

## passwd # 改变用户的密码

> 参考：
>
> - [Manual(手册),passwd(1)](https://man7.org/linux/man-pages/man1/passwd.1.html)

### Syntax(语法)

## chage # 控制用户的密码到期信息

> 参考：
>
> - [Manual(手册),chage(1)](https://man7.org/linux/man-pages/man1/chage.1.html)

passwd 软件包将会记录用户上次更改密码的时间、应该间隔多久更改一次密码 等等，chage 工具就可以对上述信息进行管理

`chage` 工具仅控制 /etc/shadow 文件中的信息，/etc/passwd 文件并不会影响到 `chage` 程序的实现。并且，`chage` 程序也不会报告 /etc/passwd 和 /etc/shaodw 文件的不一致情况，`pwck` 工具可用于检测两个文件的不一致处。

`chage` 工具仅限于 root 用户，但是 -l 选项除外，非特权用户可以使用 -l 选项来确定自身的密码或账户合适到期。chage 可以修改 /etc/shadow 文件中多个字段的配置。

### Syntax(语法)

**chage \[OPTIONS] LOGIN**

**OPTIONS**

- **-d, --lastday \<INT>** # 设置上次更改密码的日期。值是从 1970 年 1 月 1 日开始到某年某月某日的天数。
  - 若指定空值，则表示从没修改过密码，即 -l 选项查看的第一行的值为 never。
  - 若指定 0，则用户再次登录时，则会被强制要求立刻修改密码，否则无法登录
- **-l, --list** # 显示账户的老化信息

```bash
[root@hw-cloud-xngy-jump-server-linux-2 ~]# chage -l root
Last password change     : Oct 01, 2021 # 最后一次修改密码的时间
Password expires     : never
Password inactive     : never
Account expires      : never
Minimum number of days between password change  : 0
Maximum number of days between password change  : 99999
Number of days of warning before password expires : 7
```

- **-m, --mindays \<INT>** # 密码可以修改的最小间隔天数。如果 INT 为 0，则表示不用等待，任何时候都可以修改密码
  - 对应 shadow 文件中的第 4 个字段 minimum password age
- **-M, --maxdays \<INT>** #
  - 对应 shadow 文件中的第 5 个字段 maximum password age
- **-W, --warndays \<INT>** #
  - 对应 shadow 文件中的第 6 个字段 password warning period

## pwck
