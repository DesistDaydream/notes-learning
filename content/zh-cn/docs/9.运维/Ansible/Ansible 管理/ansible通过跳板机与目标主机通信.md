---
title: ansible通过跳板机与目标主机通信
---

在公司开发中，为了安全起见，生产环境跟开发环境是相互隔离开来的。也就是说在开发环境网络中无法直接 ssh 登录到生产环境的机器， 如果需要登录生产环境的机器，通常会需要借助跳板机，先登录到跳板机，然后通过跳板机登录到生产环境。

那么，使用 Ansible 时，如何配置，可以直接穿过跳板机呢？

大致的过程如下面的图示：

```bash
+-------------+       +----------+      +--------------+
| 开发环境机器A  | <---> |   跳板机B  | <--> | 生产环境机器B   |
+-------------+       +----------+      +--------------+
```

我们可以通过 ssh 命令的 ProxyCommand 选项来解决以上问题。

通过 ProxyCommand 选项，机器 A 能够灵活使用任意代理机制与机器 C 上的 SSH Server 端口建立连接，接着机器 A 上的 SSH Client 再与该连接进行数据交互，从而机器 A 上的 SSH Client 与机器 C 上的 SSH Server 之间建立了与一般直接 SSH 连接不太一样的间接 SSH 连接。

不过由于间接 SSH 连接的透明性，逻辑上可认为机器 A 上的 SSH Client 与机器 C 上的 SSH Server 建立了直接 SSH 连接。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wawlgb/1616124930604-fefb5961-13b2-48ca-ad59-0e623b3bdc35.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wawlgb/1616124930609-2ca4dea8-3c64-4b8f-9e78-b37dd522fdba.jpeg)

ssh 命令自提供的代理机制，在机器 A 上另外单独建立与 B 的 SSH 连接，该 SSH 连接的 B 端侧与机器 C 上的 SSH Server 端口建立连接，该 SSH 连接的 A 端侧与机器 A 上的 SSH Client 建立连接。

# 测试环境

A-本机：192.22.9.23

B-跳板机：192.22.9.21

C-目标机：192.22.4.46

条件：A->B 的互信

# 测试步骤

## 测试 1

A–>B：root 用户的互信
C：root 的登录信息

```bash
# ansible all -m ping --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q root@10.0.13.251"'
192.22.4.46 | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```

可以成功穿过跳板机。

## 测试 2

A–>B：一般用户（luke）的互信

C：root 的登录信息

```bash
# ansible all -m ping --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"'
192.22.4.46 | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```

跳板机普通用户，访问目标机器的 root 目录

```bash
# ansible all --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"' -m command -a 'ls /root/'
192.22.4.46 | CHANGED | rc=0 >>
1115.txt
anaconda-ks.cfg
deploytelegraf
deploytelegraf.tar.gz
epic-brook.18.11.20.01.tar.gz
epic-josh-threshold.0.1.26.tar.gz
images.tar.gz
metric.tar.gz
python-httplib2-0.9.2-1.el7.noarch.rpm
restates
sensu-1.6.1-1.el7.x86_64.rpm
sshpass-1.06-2.el7.x86_64.rpm
test1115.txt
tl.txt
```

结论

访问权限取决于目标机器的登录用户。

## 测试 3

A–>B：一般用户（luke）的互信

C：一般用户（luke）的登录信息

```bash
# ansible all -m ping --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"'
192.22.4.46 | SUCCESS => {
    "changed": false,
    "ping": "pong"
}
```

```bash
# ansible all --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"' -m command -a 'pwd'
192.22.4.46 | CHANGED | rc=0 >>
/home/desistdaydream
```

跳板机普通用户，目标机普通用户，无法访问 root

```bash
# ansible all --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"' -m command -a 'ls /root/'
192.22.4.46 | FAILED | rc=2 >>
ls: cannot open directory /root/: Permission deniednon-zero return code
```

目标机器的普通用户无法访问 root 权限的路径。

## 测试 4

A–>B：root 用户的互信
C：luke 的登录信息

```bash
# ansible all --ssh-common-args='-o ProxyCommand="ssh -W %h:%p -q desistdaydream@10.0.13.251"' -m command -a 'pwd'
192.22.4.46 | CHANGED | rc=0 >>
/home/luke
```

# 总结

1、若要使用跳板机功能，需要本机和跳板机的互信，任一用户的互信都可以。

2、目标机器的操作权限，取决于目标机器的登录用户信息，与跳板机的登录信息无关。

3、穿过跳板机，无法保证所有的 Ansible 功能都支持，特别是一些复杂的功能，如 synchronize 模块。后续需要用的模块需要一一测试。
