---
title: Bash 操作历史记录
---

# history 工具

> 参考：
> - [Manual(手册),history](https://www.man7.org/linux/man-pages/man3/history.3.html)
> - <https://blog.csdn.net/m0_38020436/article/details/78730631>
> - <https://blog.csdn.net/sz_bdqn/article/details/46527021>

history 工具可以通过如下几个 Bash 的环境变量来配置运行方式

- **HISTTIMEFORMAT** # 历史记录的格式
- **HISTSIZE** # 历史记录可以保留的最大命令数
- **HISTFILESIZE** # 历史记录可以保留的最大行数
- **HISTCONTROL** #

## 应用示例

- export HISTTIMEFORMAT="%Y-%m-%d:%H-%M-%S:`whoami`: "
- 持久化

```bash
    cat > /etc/profile.d/custom_ops.sh <<END
export HISTTIMEFORMAT="%Y-%m-%d %H:%M:%S `whoami` "
END
```

# 谁动了我的主机? 之活用 History 命令

> 参考：
> - <http://lab.xmirror.cn/2017/05/26/sdlwdzj/>

Linux 系统下可通过 history 命令查看用户所有的历史操作记录，在安全应急响应中起着非常重要的作用，但在未进行附加配置情况下，history 命令只能查看用户历史操作记录，并不能区分用户以及操作时间，不便于审计分析。

当然，一些不好的操作习惯也可能通过命令历史泄露敏感信息。

下面我们来介绍如何让 history 日志记录更细化，更便于我们审计分析。

### 1、命令历史记录中加时间

默认情况下如下图所示，没有命令执行时间，不利于审计分析。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148137-f9789cc5-18b7-4fef-9319-70bcf57795e7.jpeg)
通过设置 export HISTTIMEFORMAT='%F %T '，让历史记录中带上命令执行时间。

注意”%T”和后面的”’”之间有空格，不然查看历史记录的时候，时间和命令之间没有分割。

要一劳永逸，这个配置可以写在/etc/profile 中，当然如果要对指定用户做配置，这个配置可以写在/home/$USER/.bash_profile 中。

本文将以/etc/profile 为例进行演示。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148161-a02cf695-ebc3-4678-83bf-0dbb6a605411.jpeg)
要使配置立即生效请执行 source /etc/profile，我们再查看 history 记录，可以看到记录中带上了命令执行时间。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148125-d672b6e4-155b-43e0-90cf-4c34f4b15fd7.jpeg)
如果想要实现更细化的记录，比如登陆过系统的用户、IP 地址、操作命令以及操作时间一一对应，可以通过在/etc/profile 里面加入以下代码实现

export HISTTIMEFORMAT="%F %T`who -u am i 2>/dev/null| awk '{print $NF}'|sed -e 's/[()]//g'whoami` "，注意空格都是必须的。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148147-39797dc0-8f2b-4481-8a11-521ad0d3be06.jpeg)
修改/etc/profile 并加载后，history 记录如下，时间、IP、用户及执行的命令都一一对应。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148143-0955f121-58ad-4cc4-a350-c2103b8cb990.jpeg)

通过以上配置，我们基本上可以满足日常的审计工作了，但了解系统的朋友应该很容易看出来，这种方法只是设置了环境变量，攻击者 unset 掉这个环境变量，或者直接删除命令历史，对于安全应急来说，这无疑是一个灾难。

针对这样的问题，我们应该如何应对，下面才是我们今天的重点，通过修改 bash 源码，让 history 记录通过 syslog 发送到远程 logserver 中，大大增加了攻击者对 history 记录完整性破坏的难度。

### 2、修改 bash 源码，支持 syslog 记录

首先下载 bash 源码，可以从 gnu.org 下载，这里不做详细说明了，系统需要安装 gcc 等编译环境。我们用 bash4.4 版本做演示。

修改源码：bashhist.c
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148122-2898b5a5-b861-491b-ac1a-201aa6403243.jpeg)
修改源码 config-top.h，取消/#define SYSLOG_HISTORY/这行的注释
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148162-6b252c51-d8d8-4f98-9ad2-a67e574db7f6.jpeg)
编译安装，编译过程不做详细说明，本文中使用的编译参数为： ./configure --prefix=/usr/local/bash，安装成功后对应目录如下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148144-a4bcb15d-3f39-477d-bd7f-1167d58f615d.jpeg)
此时可以修改/etc/passwd 中用户 shell 环境，也可以用编译好的文件直接替换原有的 bash 二进制文件，但最好对原文件做好备份。

替换时要注意两点:

> 1、一定要给可执行权限，默认是有的，不过有时候下载到 windows 系统后，再上传就没有可执行权限了，这里一定要确定，不然你会后悔的；2、替换时原 bash 被占用，可以修改原用户的 bash 环境后再进行替换。

查看效果，我们发现 history 记录已经写到了/var/log/message 中。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148170-cb8d9280-1422-4701-95f3-2cdd4da96b59.jpeg)
如果要写到远程 logserver，需要配置 syslog 服务，具体配置这里不做详细讲解，大家自己研究，发送到远端 logserver 效果如下图所示。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1616165148152-c4072293-34eb-489d-ba93-70e05632e218.jpeg)
通过以上手段，可以有效保证 history 记录的完整性，避免攻击者登录系统后，通过取消环境变量、删除 history 记录等方式抹掉操作行为，为安全审计、应急响应等提供了完整的原始数据。

# linux 系统监控：记录用户操作轨迹，谁动过服务器

> **参考：**
>
> - <https://blog.51cto.com/ganbing/2053636>

**1、前言**

我们在实际工作当中，都碰到过误操作、误删除、误修改过配置文件等等事件。对于没有堡垒机的公司来说，要在 linux 系统上深究到底谁做过配置文件的修改、做过误删除是很头疼的事情，特别是遇到删库跑路的事件，更头大了。当然你可以通过 history 来查看历史命令记录，如果把 history 记录涂抹掉了，是不是啥也看不到了，如果你想查看在某个时间段到底是谁通过 vim 编辑过某个文件呢？

那么，有什么办法可以看见这些操作呢，答案是一定有的，具体怎么实现呢，linux script 命令正有如此强大的功能，可以满足我们的需求，script 可以记录终端会话，只要是 linux6.3 以上的系统，都会自带 script 命令，下面我用 centos 7 系统来测试一下。

**2、配置**

**2.1 验证 script 命令（我这里是有的）**

```bash
[root@localhost ~]# which script
/usr/bin/script
```

**2.2 配置 profile 文件，在末尾添加如下内容：**

```bash
[root@localhost ~]# vim /etc/profile

if [ $UID -ge 0 ]; then
        exec /usr/bin/script -t 2>/var/log/script/$USER-$UID-`date +%Y%m%d%H%M`.date -a -f -q /var/log/script/$USER-$UID-`date +%Y%m%d%H%M`.log
fi
```

```bash
  -t　　　 指明输出录制的时间数据
    -f     如果需要在输出到日志文件的同时，也可以查看日志文件的内容，可以使用 -f 参数。PS:可以用于教学,两个命令行接-f可以实时演示
    -a     输出录制的文件，在现有内容上追加新的内容
    -q     可以使script命令以静默模式运行
```

如下图所示：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/ieek6o/1647416324069-ae3a5928-f2d8-417c-a761-082d0c58a7fd.png)

说明：

用户登录执行的操作都会记录到/var/log/script/\*.log  里（保存日志的目录根据你自己定义），我们可以通过 more、vi 等命令查看目录里的日志。

注意：

- 我这里把用户 ID 大于 0 的都记录下来了，你可以重新登录用户，随便操作一些命令，查看生成的文件。
- root 用户的 ID 为 0，新建普通用户的 UID 是从 500 开始的(通过 cat /etc/password 可以查看用户的 UID)，如果你不想记录 root 用户的操作，你把 if 里面的值改成 500：  if \[ $UID - ge 500 ];

**2.3 创建目录、赋予权限**

你是不是以为写了这条 if 语句在/etc/profile 文件中就完事了，目录都没创建呢：

**2.4 使环境生效**

3、验证

好了，你可以退出 linux 终端，在重新登录一下，然后随便敲几个命令来看看。

```bash
[root@localhost ~]# cd /var/log/script/
[root@localhost script]# ll
total 16
-rw-r--r-- 1 root root   68 Dec 22 15:46 root-0-201712221545.date
-rw-r--r-- 1 root root  111 Dec 22 15:46 root-0-201712221545.log
-rw-r--r-- 1 root root    0 Dec 22 15:46 root-0-201712221546.date
-rw-r--r-- 1 root root 5693 Dec 22 15:46 root-0-201712221546.log
```

从上图可以看到，在/var/log/script 目录中，已经产生了 log 和 data 为后缀的文件，并且还看到了 root 用户和 UID 号 0。

.log：记录了操作

.data：可以回放操作

我们用 scriptreplay 来回放一下操作，看下效果如何：

```bash
[root@localhost script]# scriptreplay root-0-201712221545.date root-0-201712221545.log
```

**注意：** 先指定“时间文件 .data”，然后是“命令文件 .log”，不要颠倒了。

以上就完也了记录用户的所有操作，并且还可以随时查看，相当于有回放功能，像录像一样，以后定位是谁的问题就好找原因了。
