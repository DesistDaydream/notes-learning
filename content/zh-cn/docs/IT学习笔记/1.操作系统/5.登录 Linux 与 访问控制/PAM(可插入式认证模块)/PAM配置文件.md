---
title: PAM配置文件
---

# 概述

> 参考：
> - [Manual(手册),pam.conf(5)](https://man7.org/linux/man-pages/man5/pam.conf.5.html)

这是一个最基本的配置文件示例：

```bash
#%PAM-1.0
auth     required pam_deny.so
account  required pam_deny.so
password required pam_deny.so
session  required pam_deny.so
```

# Syntax(语法)

PAM 配置文件由 **Rules(规则)** 列表组成，每条规则一行。规则是由空格分割的多个 **Tokens** 组成

> 不知道官方为什么把每个字段要称为 Token~~~o(╯□╰)o

**Service Type Control Module-Path Module-Arguments**

- **Service** # 需要调用 PAM 的应用程序的名称。比如 su、login、sshd 等等
  - 注意：/etc/pam.conf 和 /etc/pam.d/\* 配置文件有一点差别，在于 Service 字段。/etc/pam.d/ 目录下的所有配置文件，没有 Service 字段，取而代之的是文件名称，也就是说，Service 字段的值，就是 /etc/pam.d/ 目录下的文件名。
- **Type **# 管理类型，这个类型就是 《[PAM(可插入式认证模块)](docs/IT学习笔记/1.操作系统/5.登录%20Linux%20与%20访问控制/PAM(可插入式认证模块)/_index.md#Linux-PAM%20管理组(认证功能的分组))》 的简写。即.本条规则中使用的模块要与哪个管理组关联。
  - 可用的类型有 auth、account、password、session
  - 若在类型前面加上 `-`，则表示即使模块不存在，也不会影响认证结果，也不会将事件记录到日志中。
- **Control **# 规则执行完成后的行为。即调用 PAM API 完成后，会有返回值，根据返回值，决定如何进行后续认证。
- **Module-Path** # 规则调用的 PAM 模块名称，模块默认在 **/usr/lib64/security/** 目录(CentOS 系统)下。
  - 不同系统中，模块所在的默认路径可能不一样。
  - 若调用的 PAM 模块不在默认目录下，则该字段需要使用**模块的绝对路径**。
- **Module-Arguments** # 规则调用的 PAM 模块的参数。每个参数以空格分隔。

## Service

Service 除了以应用程序命名，还可以使用自定义的名称，这些名称通常通过 include 这种 Control 行为引用该 Service。

## Type

指定规则所属的管理组。用于定义规则调用的模块需要与哪个管理类型关联。

> 也就是指定这条规则指定的模块应该使用的模块类型。

- **account **# 对应账户管理。验证用户是否有权限访问。
  - 比如验证用户的密码是否过期、验证用户是否有权访问所请求的服务
- **auth **# 对应身份验证管理。验证用户身份，就是证明 root 是 root
  - 比如让应用程序提示用户输入密码来确定该用户就是其所声称的身份。
- **password **# 对应密码管理，用于更改用户密码以及强制使用强密码配置
  - 比如修改密码时，必须满足强度要求。
- **session **# 对应会话管理，用户管理和配置用户会话。会话在用户成功认证后启动生效

account 与 auth 的配合可以实现这么一个场景：

- 比如现在有这么一个场景，张三要去商场买酱油。当张三到达商场后，保安人员首先要对张三进行认证，确认张三这个人可以进入商场；然后张三到达货柜拿走酱油去结账，收银人员进行授权检查，核验张三是否有权力购买酱油。

## Control

Control 会根据当前规则的执行结果，执行后续操作，也就是控制。比如当一条规则失败时，是否继续执行后面的规则；当一条规则成功时，是否执行后面的规则；等等。

Control 有两种语法，简单与复杂。简单语法通过单一的指令，来定义规则执行后的行为；复杂指令通过 1 个或多个键值对来定义规则执行后的行为。

#### 简单语法

- **requisite** # 验证失败时，立即结束整个验证过程，返回 failure。
  - 就好比让你答题 100 道，如果在答题的过程中有一道做错了直接让你出去，不会进行下面的答题过程。拥有一票否决，此关不过，立即返回 failure。
- **required** # 验证失败时，最后会返回 failure，但仍需执行同一个规则栈中的其他规则。拥有参考其他模块意见基础之上的一票否决权。可以通过其它模块来检查为什么验证没有通过。
- **sufficient** # 验证成功且之前的 required 模块没有失败时，立即结束整个验证过程，返回 true。验证失败时，忽略失败结果并继续执行栈中的后续规则。
  - 换句话说，sufficient 的验证失败对整个验证没有任何影响。
- **optional** # 可选条件，无论验证结果如何，均不会影响。通常用于 session 类型。
  - 该模块返回的通过/失败结果被忽略。当没有其他模块被引用时，标记为 optional 模块并且成功验证时该模块才是必须的。该模块被调用来执行一些操作，并不影响模块堆栈的结果。
- **include** # 包含另外一个配置文件中**相同类型**的行。比如 `password  include  system-auth-ac` 则会从 system-auth-ac 文件中，将 Type 字段为 password 的行填充到本文件中。
  - 为当前规则中指定的 Type 引用 Module-Path 中定义的规则
- **substack **# 子栈。这与 include 的不同之处在于，对子规则栈中的 done 和 die 操作的评估不会导致跳过完整模块堆栈的其余部分

#### 复杂语法

**\[Value1=Acton1 Value2=Action2 ... ValueN=ActionN]**

- Value # 该规则调用的模块执行完成后的返回码。
  - 可用的返回码有：**success**; **open_err**; **symbol_err**; **service_err**; **system_err**; **buf_err**; **perm_denied**; **auth_err**; **cred_insufficient**; **authinfo_unavail**; **user_unknown**; **maxtries**; **new_authtok_reqd**; **acct_expired**; **session_err**; **cred_unavail**; **cred_expired**; **cred_err**; **no_module_data**; **conv_err**; **authtok_err**; **authtok_recover_err**; **authtok_lock_busy**; **authtok_disable_aging**; **try_again**; **ignore**; **abort**; **authtok_expired**; **module_unknown**; **bad_item**; and **default**
- Action # 表示当发现该返回码时，要执行的行为。
  - 可用的行为有：**ignore、bad、die、ok、done、reset、N**

#### 简单语法与复杂语法的对应关系

| 简单语法   | 复杂语法                                                    |
| ---------- | ----------------------------------------------------------- |
| required   | \[success=ok new_authtok_reqd=ok ignore=ignore default=bad] |
| requisite  | \[success=ok new_authtok_reqd=ok ignore=ignore default=die] |
| sufficient | \[success=done new_authtok_reqd=done default=ignore]        |
| optional   | \[success=ok new_authtok_reqd=ok default=ignore]            |

## Module-Path

在 CentOS 中，模块文件默认在 /usr/lib64/security/ 目录中，以 .so 结尾

```bash
~]# ls /usr/lib64/security/
pam_access.so   pam_cracklib.so  pam_env.so        pam_filter     pam_issue.so    pam_listfile.so   pam_mkhomedir.so  pam_permit.so      pam_rhosts.so          pam_selinux.so   pam_succeed_if.so  pam_timestamp.so  pam_unix_auth.so     pam_userdb.so
pam_cap.so      pam_debug.so     pam_exec.so       pam_filter.so  pam_keyinit.so  pam_localuser.so  pam_motd.so       pam_postgresok.so  pam_rootok.so          pam_sepermit.so  pam_systemd.so     pam_tty_audit.so  pam_unix_passwd.so   pam_warn.so
pam_chroot.so   pam_deny.so      pam_faildelay.so  pam_ftp.so     pam_lastlog.so  pam_loginuid.so   pam_namespace.so  pam_pwhistory.so   pam_securetty.so       pam_shells.so    pam_tally2.so      pam_umask.so      pam_unix_session.so  pam_wheel.so
pam_console.so  pam_echo.so      pam_faillock.so   pam_group.so   pam_limits.so   pam_mail.so       pam_nologin.so    pam_pwquality.so   pam_selinux_permit.so  pam_stress.so    pam_time.so        pam_unix_acct.so  pam_unix.so          pam_xauth.so

```

在 Ubuntu 中，模块文件默认在 /usr/lib/x86_64-linux-gnu/security/ 目录中，以 .so 结尾

```bash
~]$ ls /usr/lib/x86_64-linux-gnu/security/
pam_access.so  pam_echo.so        pam_faildelay.so  pam_ftp.so            pam_issue.so    pam_listfile.so   pam_mkhomedir.so  pam_permit.so     pam_securetty.so  pam_stress.so      pam_tally.so      pam_umask.so   pam_wheel.so
pam_cap.so     pam_env.so         pam_faillock.so   pam_gdm.so            pam_keyinit.so  pam_localuser.so  pam_motd.so       pam_pwhistory.so  pam_selinux.so    pam_succeed_if.so  pam_time.so       pam_unix.so    pam_xauth.so
pam_debug.so   pam_exec.so        pam_filter.so     pam_gnome_keyring.so  pam_lastlog.so  pam_loginuid.so   pam_namespace.so  pam_rhosts.so     pam_sepermit.so   pam_systemd.so     pam_timestamp.so  pam_userdb.so
pam_deny.so    pam_extrausers.so  pam_fprintd.so    pam_group.so          pam_limits.so   pam_mail.so       pam_nologin.so    pam_rootok.so     pam_shells.so     pam_tally2.so      pam_tty_audit.so  pam_warn.so

```

## Module-Arguments

详见 《[PAM 模块详解](✏IT 学习笔记/📄1.操作系统/5.登录%20Linux%20 与%20 访问控制/PAM(可插入式认证模块)/PAM%20 模块详解.md 模块详解.md)》

# /etc/pam.d/password-auth 与 /etc/pam.d/system-auth

- CentOS 发行版中，这个两个文件分别是 password-auth-ac 和 system-auth-ac 的软链接。两个 \*-ac 文件，则是由 `authconfig` 程序生成的
- Ubuntu 发行版中，这个两个文件分别是 common-password 和 common-auth。这两个文件，则是由 `pam-auth-update` 程序生成的。

通常情况下，如果想要添加更多的认证配置，推荐使用一个新的文件，并使用 include 指令包含这俩文件即可。

## CentOS 发行版配置

auth required pam_env.so
auth required pam_faildelay.so delay=2000000
auth sufficient pam_unix.so nullok try_first_pass
auth requisite pam_succeed_if.so uid >= 1000 quiet_success
auth required pam_deny.so

account required pam_unix.so
account sufficient pam_localuser.so
如果用户 ID 小于 1000，直接退出，不再进行验证
account sufficient pam_succeed_if.so uid < 1000 quiet
account required pam_permit.so

password requisite pam_pwquality.so try_first_pass local_users_only retry=3 authtok_type=
password sufficient pam_unix.so sha512 shadow nullok try_first_pass use_authtok

password required pam_deny.so

session optional pam_keyinit.so revoke
session required pam_limits.so
-session optional pam_systemd.so
session \[success=1 default=ignore] pam_succeed_if.so service in crond quiet use_uid
session required pam_unix.so

# /etc/pam.d/sshd

这是用于安全的 Shell 服务的 PAM 配置，比如 OpenSSH

Standard Un\*x authentication.
@include common-auth
====================

Disallow non-root logins when /etc/nologin exists.
account required pam_nologin.so
=======================================

Uncomment and edit /etc/security/access.conf if you need to set complex
\# access limits that are hard to express in sshd_config.
\# account required pam_access.so
=======================================

Standard Un\*x authorization.
@include common-account
=======================

SELinux needs to be the first session rule. This ensures that any
\# lingering context has been cleared. Without this it is possible that a
\# module could execute code in the wrong domain.
session \[success=ok ignore=ignore module_unknown=ignore default=bad] pam_selinux.so close
===================================================================================================

Set the loginuid process attribute.
session required pam_loginuid.so
========================================

Create a new session keyring.
session optional pam_keyinit.so force revoke
====================================================

Standard Un\*x session setup and teardown.
@include common-session
=======================

登录成功后打印当天消息。这包括来自 /run/motd.dynamic 的动态生成部分和来自 /etc/motd 的静态部分。
注释这两行，将会禁用登录后的消息提示功能
session optional pam_motd.so motd=/run/motd.dynamic
session optional pam_motd.so noupdate

Print the status of the user's mailbox upon successful login.
session optional pam_mail.so standard noenv # \[1]
==========================================================

Set up user limits from /etc/security/limits.conf.
session required pam_limits.so
======================================

Read environment variables from /etc/environment and
\# /etc/security/pam_env.conf.
session required pam_env.so # \[1]
\# In Debian 4.0 (etch), locale-related environment variables were moved to
\# /etc/default/locale, so read that as well.
session required pam_env.so user_readenv=1 envfile=/etc/default/locale
===============================================================================

SELinux needs to intervene at login time to ensure that the process starts
\# in the proper default security context. Only sessions which are intended
\# to run in the user's context should be run after this.
session \[success=ok ignore=ignore module_unknown=ignore default=bad] pam_selinux.so open
==================================================================================================

Standard Un\*x password updating.
@include common-password
========================

# /etc/pam.d/su

**auth sufficient pam_rootok.so**
当开始使用 pam_wheel.so 模块时，只有属于 wheel 组的用户，才可以使用 su 命令切换到 root 用户
**auth sufficient pam_wheel.so trust use_uid**

- 当用户在 wheel 组时，使用 su - root 命令不需要密码即可切换到 root 用户

**auth required pam_wheel.so use_uid**

- 当用户在 wheel 组时，使用 su - root 命令需要密码即可切换到 root 用户

**auth substack system-auth**
**auth include postlogin**
**account sufficient pam_succeed_if.so uid = 0 use_uid quiet**
**account include system-auth**
**password include system-auth**
**session include system-auth**
**session include postlogin**
**session optional pam_xauth.so**

# 配置示例

**/etc/pam.d/sshd 配置文件示例**
注意 sshd、login、remote、kde 这几个文件中的配置大部分都相同，

    # %PAM-1.0
    # 最多连续三认认证登录都出错时，60秒后解锁，root用户也可以被锁定，root用户15秒后解锁。
    auth required pam_tally2.so deny=3 unlock_time=60 even_deny_root root_unlock_time=15

**/etc/pam.d/common-password 文件配置示例**

    # 限制用户不能更改为之前使用的历史密码
    password required pam_pwhistory.so use_authtok remember=6 retry=3

问题实例：限制用户不能更改为之前使用的历史密码

- Linux 历史密码在 /etc/security/opasswd 中存放
- 解决方法
  - 临时更改 commen-password 文件修改密码修改策略，去除历史密码的限制，更改完密码后在恢复原来的策略
  - 删掉 /etc/security/opasswd 中关于被修改文件的内容，这样就检测不到之前的历史密码了
