---
title: RsysLog 配置详解
---

# 概述

> 参考：
>
> - [Manual(手册),rsyslog.conf(5)](https://man7.org/linux/man-pages/man5/rsyslog.conf.5.html)

配置文件新老版本有两种格式，下面这两种写法，都可以表示让 Rsyslog 加载 imuxsock 模块

- **$ModLoad imuxsock # 老版本语法**
- **module(load="imuxsock") # 新版本语法**

# MODULES 模块配置

**module(load="imuxsock") **# 加载 imuxsock 模块，以便让 Rsyslog 可以监听 /dev/log 这个 Unix Socket 以接收日志消息<br />配置 TCP 协议的 syslog 接收，用于在日志服务器的时候配置

- **$ModLoad imtcp** # 使用 tcp 进行传输
- **$InputTCPServerRun 514** #监听在 514 端口上

# GLOBAL DIRECTIVES 全局指令配置

配置 rsyslogd 的全局属性，比如主信息队列大小等

# TEMPLATE 模板配置

该配置用于自定义日志的保存路径，日志格式，可以动态生成文件名等信息。定义后，可以在 RULE 中进行引用，该指令用法详见

    # 定义一个Location字段的模板，第一个是老版本的定义方法，第二个是新版的定义方法
    $template RemoteLogs,"/var/log/%HOSTNAME%/%PROGRAMNAME%.log" *
    template (name="RemoteLogs" type="string" string="/var/log/%HOSTNAME%/%PROGRAMNAME%.log")

## TEMPLATE 模板介绍

template 是 rsyslog 的一个关键功能，可以允许用户指定想要的任何日志格式。e.g.自定义日志的保存路径，日志格式。还可以动态生成文件名等信息

template 有两种表示方式：

- 7.0 之前的版本使用 $template
- 7.0 之后的版本使用 template()

### 语法结构：template(Parameters)

Parameters 必须包含 name 字段且 name 唯一，并指明类型，以及该类型的具体内容
template(NAME TYPE Descriptions)

1. TYPE
   1. list
   2. subtree
   3. string
   4. plugin
2. EXAMPLE
   1. template (name="RemoteLogs" type="string" string="/var/log/%HOSTNAME%/%PROGRAMNAME%.log")
   2. :programname, regex, "Keepalived.\*" -/var/log/keepalived/keepalived.log #根据程序名字，使用正则表达式，开头是 Keepalived 的日志，写入到/var/log/keepalived/keepalived.log 文件中

可用的属性详见：<https://www.rsyslog.com/doc/v8-stable/configuration/properties.html>

$template NAME,"PATH" #定义一个名为 NAME 的模板来作为 RULE 配置段中 Location 字段使用，在 Location 字段中通过?NAME 来引用该对应模板

- PATH 的可用变量
  - %HOSTNAME% #用来区分是哪台远程主机的。
  - %PROGRAMNAME% #通过日志标准格式中的 ProgramName 字段来进行分类保存日志。i.e.每个程序名是单独的一个文件
  - %$year%%$month%%$day% #用来以时间格式命名文件

# RULES 规则配置

Rules 配置段是 rsyslog 程序得以正常运行的最基础配置。规则的内容是告诉 rsyslog 处理日志的方式。i.e.每条规则用于定义以下几个内容：1.什么设施的。2.什么优先级。3.需要被记录在哪里

每条规则占用一行，规则的内容分为两个字段：Selector(选择器)和 Action(动作 i.e.即对选择器选择出来的设施和优先级进行什么操作)。这两个字段由一个或多个空格或者制表符分割。而 Selector 则是通过 Facility(设施)、分隔符、Priority(优先级，也可以用 Level(级别)来表示)来对整个系统的所有设施日志进行筛选

### Syntax(语法)

**Selector Action**

- **Selector(选择器) **# 根据匹配规则，选择要处理的日志。选择器由 Facility 和 Priority 组成，以 `.` 分隔
  - **Facility.Priority**
- **Action(动作)** # 描述了如何处理选择器选择出来的日志信息

#### Selectors(选择器)

多个选择器以 `;` 符号分隔

- **Facility(设施)** # Facility 定义了 rsyslog 可以选择的设施都有哪些(注：该字段用 \* 表示则表示所有支持的 Facility)。多个 Facility 以 `,` 分隔
  - 可选择设施见上文的对 Facility 的介绍
- **匹配符号** # 除了 `.` 还可以使用另外两个符号来进行更细致的匹配。
  - `.` # 选择包含且比 Prority 还要严重的优先级
  - `.=` # 仅选择包含 Prority 所定义的优先级
  - `.!` # 选择不包含 Prority 中所定义的优先级的其余优先级
- **Priority(日志的优先级，也可以叫日志的 Level 级别)** # Priority 定义了每条信息的严重程度，下面以严重程度从高到低进行排序。括号中的数字指级别(注：该字段用\*表示所有级别)
  - emerg(0)：错误信息。最严重日志等级，意味着系统将要宕机
  - alert(1)：错误信息。比 emerg 等级轻
  - crit(2)：错误信息。
  - err(3)：错误信息。err 就是 error
  - warn(4)：警告信息。可能有问题，但是还不至于影响到程序的运行。warn 就是 warnning
  - notice(5)：基本信息。
  - info(6)：基本信息。
  - debug(7)：特殊的等级，用来 troubleshooting 时产生的日志
  - none：特殊的等级。表示某个 Facility 不需要执行 Action。i.e.即不记录的级别

#### Action(动作)

匹配到的日志将要执行的动作。是保存本地文件、打印、保存到远程主机、转存到数据库中等等行为

- **RegularFile(常规文件)** # 把日志写入到某个文件，文件路径可以引用模板。
  - 如果在该字段前面记上 `-`，则表示先将日志保存在内存的 buffer 中，等数据量足够大时再一次性将数据写入磁盘文件中。
- **RemoteMachine(远程主机)** # `@HOST` 或者 `@@HOST`。用于把日志发送给远程主机。@使用 UDP 协议，@@使用 TCP 协议，默认使用 514 端口
- \*，表示把日志发送给目前在线的所有人，类似于 wall 命令
- | COMMAND：用于把日志信息通过管道符送给后面定义的 COMMAND 来进行处理
- 打印机或其他。
- 使用者名称(显示给用户)。
- **stop** # Discard(丢弃)，老版本使用 `~` 表示

### CentOS 系统中的默认配置示例

```bash
# 内核产生的日志全部送入到终端设施中，用于在系统出现严重问题，无法使用默认屏幕观察，可以使用笔记本连接到封闭服务器的RS232端口后，查看日志
# kern.*                                                 /dev/console
# 除了mail、authpriv、cron以外的优先级为info且以上的所有设施的信息写入/var/log/messages文件中，Action字段以绝对的路径表示
*.info;mail.none;authpriv.none;cron.none                /var/log/messages
# 认证方面的日志写入/var/log/secure文件中
authpriv.*                                              /var/log/secure
# 同上，用于邮件相关日志。
mail.*                                                  -/var/log/maillog
	# 同上，用于定时任务相关日志
cron.*                                                  /var/log/cron
# 任何优先级为emerg的设施日志以wall广播的方式给所有在系统登录的账号
*.emerg                                                 :omusrmsg:*
# 将uucp和news的等级为crit且以上的日志写入spooler文件中
uucp,news.crit                                          /var/log/spooler
# local7这个设施的所有优先级的日志写到/var/log/boot.log文件中
local7.*                                                /var/log/boot.log

# 下面是对Rule的应用实例：
# 除了sshd、keepalived、haproxy程序的日志以外，其余所有程序的所有等级的日志保存在/var/log/messages文件中
*.*;sshd,keepalived,haproxy.none	/var/log/messages
# 将local2设施的所有级别的日志写入名为RemoteLogs模板定义的文件中，Action字段引用模板RemoteLogs的路径
local2.*   ?RemoteLogs
# 以下符号用来告知 rsyslog 停止对日志消息的进一步处理，即只把日志写入指定的路径中，而不再重复写到默认的 /var/log/* 目录下
& stop # 注意，该指令仅对其上面一行的规则起作用，想让哪一条的规则生效，则在哪一行下面加上该组符号
# 把所有程序的所有级别的日志发送给192.168.10.10这台机器上，Action字段为远程主机
*.*	@@192.168.10.10
# 丢弃所有日志,Action字段为丢弃
*.* stop
# 来自于远程主机的日志，不包括本机的日志全部写入到RemoteLogs模板定义的路径中
:FROMHOST-IP, !isequal, "127.0.0.1"                     -?RemoteLogs
```

# outputs # 输出

# 配置实例：

注意：所有

## 实例一：设置一台主机为日志服务器，可以收集其余网络上的主机的日志信息

日志服务器与其余服务器形成服务端与客户端的关系，服务端的日志服务监听某个端口，来收集其余机器发送过来的日志信息

Server 端：

    $ModLoad imtcp
    $InputTCPServerRun 514
    $template RemoteLogs,"/var/log/%HOSTNAME%/%PROGRAMNAME%.log" *
    *.*  ?RemoteLogs
    & stop

:FROMHOST-IP, !isequal, "127.0.0.1" -?DynaFile

client 端：

在日志规则的 Location 配置段使用远程主机配置

    *.*	@@192.168.10.10	#把所有程序的所有级别的日志发送给192.168.10.10这台机器
    注意：如果使用UDP进行传输，则使用1个@

## 实例二：配置 keepalived 程序的日志到指定的目录

直接使用实例一 Server 端的后三条配置即可

## 实例三：根据正则匹配忽略某些日志信息写入到文件中

下面的配置表示：如果输出日志的程序为 kubelet，并且日志信息中包含 Setting volume ownership 这种内容，那么所有匹配到的日志全部丢弃，不写入到文件中。

    cat > /etc/rsyslog.d/ignore-kubelet-volume.conf << EOF
    if (\$programname == "kubelet") and (\$msg contains "Setting volume ownership") then {
      stop
    }
    EOF
