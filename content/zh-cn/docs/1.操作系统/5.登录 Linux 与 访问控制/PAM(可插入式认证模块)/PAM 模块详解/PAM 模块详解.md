---
title: "PAM 模块详解"
linkTitle: "PAM 模块详解"
weight: 20
---

# 概述

> 参考：

# PAM 的各模块说明

全局参数

- file=/PATH/TO/FILE # 用于指定统计次数存放的位置，默认保存在/var/log/tallylog 文件中；
- onerr # 当意外发生时，返加 PAM_SUCCESS 或 pam 错误代码，一般该项不进行配置；
- audit # 如果登录的用户不存在，则将访问信息写入系统日志；
- silent # 静默模式，不输出任何日志信息；
- no_log_info # 不打印日志信息通过 syslog
- 上面的五项全局参数，一般在使用中都不需要单独配置。

# pam_faillock # 在指定的时间间隔内计算身份验证失败

> 参考：
>
> - [Manual(手册),pam_faillock(8)](https://man.cx/pam_faillock)
> - [Manual(手册),faillock.conf(5)](<https://man.cx/faillock.conf(5)>)
> - <https://github.com/dev-sec/ansible-collection-hardening/issues/377>
> - 红帽官方文档,安全指南-账户锁
>   - <https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/security_guide/chap-hardening_your_system_with_tools_and_services#sect-Security_Guide-Workstation_Security-Account_Locking>

提供 auth、account 管理类型的模块

pam_faillock 模块在指定的时间间隔内维护每个账户在尝试进行身份验证时的失败事件，并且在连续失败时锁定账户。

pam_faillock 与大部分模块有一点不同，不建议在 PAM 规则中配置参数，而是推荐使用默认的 /etc/security/faillock.conf 配置文件中配置参数

## 关联文件

**/etc/security/faillock.conf** # 运行时配置文件。除了在 /etc/pam.d/\* 文件中配置模块的参数，还可以通过这个文件配置模块的参数。

**/var/run/faillock/** # 记录用户身份验证失败的事件。目录中的文件名以用户名命名

## 模块参数

**preauth | authfail | authsucc** #

**conf=\</PATH/TO/FILE>** # 指定要使用的配置文件路径。

## 应用示例

登录失败 3 次会锁定用户 60 秒，账户登录失败 3 次 锁定 30 秒

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

## 命令行工具

### faillock

**faillock \[OPTIONS]**

管理登录失败锁定记录的工具

```bash
[root@LNDL-PSC-SCORE-PM-OS04-EBRS-HA02 pam.d]# faillock
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

# pam_nologin

这个模块可以限制一般用户是否能够登入主机之用。当 /etc/nologin 这个文件存在时，则所有一般使用者均无法再登入系统了！若 /etc/nologin 存在，则一般使用者在登入时， 在他们的终端机上会将该文件的内容显示出来！所以，正常的情况下，这个文件应该是不能存在系统中的。 但这个模块对 root 以及已经登入系统中的一般账号并没有影响。

# pam_pwhistory # 记住最后的密码

> 参考：
>
> - [Manual(手册),pam_pwhistory(8)](https://man7.org/linux/man-pages/man8/pam_pwhistory.8.html)

该模块用于记住用户设置过的密码，以防止用户在修改密码时频繁交替得使用相同的密码

## 关联文件

**/etc/security/opasswd** # 用户设置过的历史密码将会以加密方式保存在该文件中。

## 模块参数

- **remember=INT** # 用户设置过的 remember 个密码将会保存在 /etc/security/opasswd 文件中。`默认值：10`。值为 0 时，模块将会保持 opasswd 文件的现有内容不变

# pam_pwquality # 密码质量检查

> 参考：
>
> - [GitHub 项目，libpwquality/libpwquality](https://github.com/libpwquality/libpwquality/)
> - [Manual(手册)，pam_pwquality(8)](https://man.cx/pam_pwquality)

pam_pwquality 模块属于 libpwquality 库，最初基于 pam_cracklib 模块，用以执行密码质量检查。仅提供 password 模块类型。

> 注意：在红帽企业 Linux 7 中，pam_pwquality PAM 模块取代 pam_cracklib，该模块在红帽企业 Linux 6 中用作密码质量检查的默认模块。它使用与 pam_cracklib 相同的后端。详见[红帽官网](https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/security_guide/chap-hardening_your_system_with_tools_and_services)

该模块的作用是提示用户输入密码，并根据系统字典和一组识别不良选择的规则检查其强度。第一个操作是提示输入一个密码，检查其强度，然后，如果认为强度高，则第二次提示输入密码（以验证第一次输入的密码是否正确）。一切顺利，密码将传递给后续模块以作为新的身份验证令牌安装。

模块可以提供如下几种类型的检查

- **Palindrome(回文)**#
- **Case Change Only(仅大小写更改)** # 新密码是否与旧密码相同，只是大小写不同？
- **Similar(相似)** # 新密码是不是太像旧密码了？这主要由一个参数 difok 控制，它是旧密码和新密码之间足以接受新密码的一系列字符更改（插入、删除或替换）。
- **Rotated(旋转的)** # Is the new password a rotated version of the old password?新密码是旧密码的轮换版本吗？
- **Same consecutive characters(相同的连续字符)** # Optional check for same consecutive characters.可选检查相同的连续字符。
- **Too long monotonic character sequence(太长的单调字符序列)** # 可选检查太长的单调字符序列。
- **Contains user name(包含用户名)** # 检查密码是否包含某种形式的用户名。
- **Dictionary check(字典检查)** # 调用 Cracklib 例程来检查密码是否是字典的一部分。

上述这些检查可以通过使用模块参数或通过修改 `/etc/security/pwquality.conf` 配置文件来配置。模块参数 j 覆盖配置文件中的设置。

## 关联文件

**/etc/security/pwquality.conf** # 模块运行时配置文件
**/usr/lib64/security/pam_pwquality.so** # 模块文件

## 模块参数

通用参数

- **retry=\<N>** # 重试次数。`默认值：1`。允许用户输入密码错误的最大次数、对于输入新密码时的情况，则是允许用户输入不符合要求的新密码的最大次数。

设置新密码时执行检查的参数，某些参数的值为 -1 时，表示新密码至少需要有 1 位数字、大写字母、特殊字符 等等。值为 0 表示禁用检查。

- **minlen=\<N>** # 新密码的最小字符数。`默认值：8`
- **dcredit=\<N>** # 新密码中包含的 **digit(数字)**的字符数。
- **ucredit=\<N>** # 新密码中包含的 **uppercase(大写字母)**的字符数。
- **lcredit=\<N>** # 新密码中包含的 **lowercase(小写字母)** 的字符数。
- **ocredit=\<N>** # 新密码中包含的 **other(其他字符)** 的字符数。其他字符就是特殊字符
- **minclass=\<N>** # 新密码中包含的字符类型的数量。`默认值：0` 共有 4 中字符类型可用：
  - digits(数字)
  - uppercase(大写字母)
  - lettercase(小写字母)
  - other(其他字符)

## 应用示例

输入错误最多 3 次，至少 14 个字符，其中 大写字母、小写字母、数字、特殊符号 这四类字符每类至少有一个。

```bash
pam_pwquality.so try_first_pass retry=3 minlen=14 dcredit=-1 ucredit=-1 ocredit=-1 lcredit=-1
```

### pwquality.conf 文件示例

```bash
minlen = 14
lcredit = -1
ucredit = -1
dcredit = -1
ocredit = -1
```

# pam_succeed_if # 测试账户特性

pam_succeed_if 模块旨在根据 **账户的特征**或 其他

说白了，从编程角度看，这就是一个典型的 if...else... 控制结构，在一组规则栈中，通过使用该模块，对一些条件进行判断。

比如

- `pam_succeed_if.so uid >= 1000` 表示调用该模块的程序所使用的账号的 uid 要是大于等于 1000，则模块返回成功
- `pam_succeed_if.so service in crond` 表示调用该模块的程序是 crond 的话，则模块返回成功。

该模块的常见用途就是根据该模块的条件测试结果，决定是否加载其他模块。

## 模块参数

- **quiet_success** # 若模块返回成功，则不要将验证事件记录到系统日志中。
- **use_uid** # 使用运行应用程序的 UID 的用户的帐户而不是正在验证的用户来评估条件。
- **Conditions(条件)** # 条件参数由 3 部分组成：`Field Test Value`。
  - 可用的 Field 有：user、uid、gid、shell、home、ruser、rhost、tty、service
  - 可用的 Test 有：`<` `<=` `eq` `>=` `>` `ne` `=` `!=` `=~` `!~` `in` `not in`
  - 比如：
    - **uid >= 1000** # 表示运行程序所使用的账号的 uid 要是大于 1000，则模块返回成功
    - 更多语法说明，详见 Manual。

## 应用示例

该模块非常容易造成 PAM 认证时的结果与想要的结果产生偏差，比如在 CentOS7 下，/etc/pam.d/password-auth-ac 的配置如下：

```bash
auth        required      pam_env.so
auth        required      pam_faildelay.so delay=2000000
auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 1000 quiet_success
auth        required      pam_deny.so
```

此时，我们想要使用 pam_faillock 模块配置一下账户认证失败后锁定的功能，如果将配置修改为如下样子：

```bash
auth        required      pam_env.so
auth        required      pam_faillock.so preauth  audit deny=3 even_deny_root unlock_time=60
auth        required      pam_faildelay.so delay=2000000
auth        sufficient    pam_unix.so nullok try_first_pass
auth        requisite     pam_succeed_if.so uid >= 1000 quiet_success
auth        [default=die] pam_faillock.so authfail audit deny=3 even_deny_root unlock_time=60
auth        required      pam_deny.so
```

此时，所有 uid 小于 1000 的，包括 root 账户，都是无法享受到第 6 行规则的效果的，因为当执行到第 5 行时，发现此次认证行为的账户是 root(uid=0) 则直接返回失败，不在执行第 6 行的规则了~~

# pam_unix # 传统密码认证

> 参考：
>
> - [Manual(手册),pam_unix(8)](https://man7.org/linux/man-pages/man8/pam_unix.8.html)

注意：推荐使用 pam_pwquality 模块与 pam_unix 模块配合使用

若是不满足密码强度要求，将会出现类似如下的提示：

```bash
[root@common-centos-test pam.d]# passwd developer
Changing password for user developer.
New password:
BAD PASSWORD: The password contains less than 1 uppercase letters
Retype new password:
BAD PASSWORD: The password fails the dictionary check - it is too simplistic/systematic
```

注意：root 用户修改任何用户的密码不受此模块限制，只有普通用户修改自己的密码时才有效。

## 模块参数

- **nullok** # 此模块的默认操作是，如果用户的官方密码为空，则不允许用户访问服务。nullok 参数覆盖此默认值。
- **try_first_pass** # 这个选项指示本模块首先尝试使用已有的密码，即从第一个向用户提示输入密码的模块那里取得密码，并对该密码进行认证。如果密码认证失败，则再提示用户输入密码。
- **use_authtok** # “use_authtok”参数确保 pam_unix 模块不会提示输入密码，而是使用 pam_pwquality 提供的密码。
- **密码要求参数**
  - **minlen=INT** # 密码长度最少 minlen 位
  - **difok=INT** # 新旧密码最少 difok 个字符不同
  - **icredit=-1** # 最少 1 个数字.
  - **lcredit=-1** # 最少 1 个小写字符
  - **ucredit=-1** # 最少 1 个大写字符
  - **ocredit=-1** # 最少 1 个特殊字符

#