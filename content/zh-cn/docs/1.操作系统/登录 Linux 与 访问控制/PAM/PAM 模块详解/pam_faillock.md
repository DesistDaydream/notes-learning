---
title: pam_faillock
linkTitle: pam_faillock
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，pam_faillock(8)](https://man7.org/linux/man-pages/man8/pam_faillock.8.html)
>   - [Ubuntu 22.04 TLS Manual](https://manpages.ubuntu.com/manpages/jammy/en/man8/pam_faillock.8.html)
> - [Manual(手册)，faillock.conf(5)](https://man7.org/linux/man-pages/man5/faillock.conf.5.html)
> - <https://github.com/dev-sec/ansible-collection-hardening/issues/377>
> - 红帽官方文档,安全指南-账户锁
>   - <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/security_guide/chap-hardening_your_system_with_tools_and_services#sect-Security_Guide-Workstation_Security-Account_Locking>

提供 auth、account 管理类型的模块

pam_faillock 模块在指定的时间间隔内维护每个账户在尝试进行身份验证时的失败事件，并且在连续失败时锁定账户。

pam_faillock 与大部分模块有一点不同，不建议在 PAM 规则中配置参数，而是推荐使用默认的 /etc/security/faillock.conf 配置文件中配置参数

## 关联文件

**/etc/security/faillock.conf** # 运行时配置文件。除了在 /etc/pam.d/ 目录下的文件中配置模块的参数，还可以通过这个文件配置模块的参数。

**/var/run/faillock/** # 记录用户身份验证失败的事件。目录中的文件名以用户名命名

## 模块参数

**preauth | authfail | authsucc** # 这 3 个参数必须根据该模块实例在 PAM 堆栈中的位置进行设置。

**conf=\</PATH/TO/FILE>** # 指定要使用的配置文件路径。

除了上面的参数外，模块的其他参数都可以在 /etc/security/faillock.conf 文件中进行配置

## 命令行工具

### faillock

**faillock \[OPTIONS]**

管理登录失败锁定记录的工具

```bash
~]# faillock
developer:
When                Type  Source                                           Valid
2021-10-21 21:42:50 RHOST 172.16.10.11                                         V
root:
When                Type  Source                                           Valid
2021-10-21 21:42:41 RHOST 172.16.10.11                                         V
```

OPTIONS

- **--user \<USERNAME>** # 指定要处理的用户名称
- **--reset** # 清除失败记录，解除锁定

# 最佳实践

普通登录失败 3 次会锁定用户 60 秒，root 账户登录失败 3 次 锁定 30 秒

```bash
    sudo tee /etc/pam.d/password-auth-local > /dev/null <<EOF
auth        required       pam_faillock.so preauth  audit deny=3 even_deny_root unlock_time=60 root_unlock_time=30
auth        include        password-auth-ac
auth        [default=die]  pam_faillock.so authfail audit deny=3 even_deny_root unlock_time=60 root_unlock_time=30

account     required       pam_faillock.so
account     include        password-auth-ac

password    include        password-auth-ac

session     include        password-auth-ac
EOF

    ln -sf /etc/pam.d/password-auth-local /etc/pam.d/password-auth
```

注意：由于 password-auth-ac 中有 pam_succeed_if.so uid >= 1000 quiet_success 这样一条规则，所以上述配置对 root 账户不起作用。

## Ubuntu 示例

参考:

https://manpages.ubuntu.com/manpages/jammy/en/man8/pam_faillock.8.html#examples

https://askubuntu.com/questions/1403438/how-do-i-set-up-pam-faillock

```bash
sudo tee /etc/security/faillock.conf > /dev/null <<EOF
deny=3
unlock_time=30
silent
even_deny_root
EOF
```

下面这个官方示例待研究，配置完没效果

```bash
sudo tee /etc/pam.d/config > /dev/null <<EOF
auth     required       pam_securetty.so
auth     required       pam_env.so
auth     required       pam_nologin.so
# optionally call: auth requisite pam_faillock.so preauth
# to display the message about account being locked
auth     [success=1 default=bad] pam_unix.so
auth     [default=die]  pam_faillock.so authfail
auth     sufficient     pam_faillock.so authsucc
auth     required       pam_deny.so
account  required       pam_unix.so
password required       pam_unix.so shadow
session  required       pam_selinux.so close
session  required       pam_loginuid.so
session  required       pam_unix.so
session  required       pam_selinux.so open
EOF
```

StackExchange 的回答中的配置可以生效，但是无法显示 `There were 4 failed login attempts since the last successful login` 这种信息

在 common-auth 文件中，pam_unix.so 模块前后添加共 3 条 pam_faillock.so 模块。效果如下:

```bash
auth    required                        pam_faillock.so preauth audit
auth    [success=1 default=ignore]      pam_unix.so nullok
auth    [default=die]                   pam_faillock.so authfail audit
auth    sufficient                      pam_faillock.so authsucc audit
auth    requisite                       pam_deny.so
auth    required                        pam_permit.so
auth    optional                        pam_cap.so
```

在 common-account 文件中，结尾添加 1 条 pam_faillock.so 模块。效果如下:

```bash
account [success=1 new_authtok_reqd=done default=ignore]        pam_unix.so
account requisite                       pam_deny.so
account required                        pam_permit.so
account required                        pam_faillock.so
```
