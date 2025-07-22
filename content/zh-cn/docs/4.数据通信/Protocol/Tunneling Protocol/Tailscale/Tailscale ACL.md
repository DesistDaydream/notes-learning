---
title: Tailscale ACL
linkTitle: Tailscale ACL
weight: 20
---

# 概述

> 参考：
>
> - [官方文档，参考 - Tailnet 策略文件语法](https://tailscale.com/kb/1337/policy-syntax)
> - [Tailscale 博客，RBAC 的初衷](https://tailscale.com/blog/rbac-like-it-was-meant-to-be)

# groups

https://tailscale.com/kb/1337/policy-syntax#groups

# hosts

https://tailscale.com/kb/1337/policy-syntax#hosts

# acls

https://tailscale.com/kb/1337/policy-syntax#acls

## dst

`dst` 字段设置一个访问目标的列表，该列表是一组适用于某 acl 规则的目标。

列表中的每个元素格式为 `HOST:PORTS`。i.e. 1 个 host，1 个或多个 ports。

HOST 可以是一以下任意类型

| **Type**                                                                         | **Example**               | **Description**                                                                                                                                                                         |
| -------------------------------------------------------------------------------- | ------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Any                                                                              | `*`                       | Includes any destination (no restrictions).                                                                                                                                             |
| User                                                                             | `shreya@example.com`      | Includes any device currently signed in as the provided user.                                                                                                                           |
| [Group](https://tailscale.com/kb/1337/policy-syntax#groups)                      | `group:<group-name>`      | Includes all users in the provided group.                                                                                                                                               |
| Tailscale IP address                                                             | `100.101.102.103`         | Includes only the device that owns the provided Tailscale IP address.                                                                                                                   |
| [Hosts](https://tailscale.com/kb/1337/policy-syntax#hosts)                       | `example-host-name`       | Includes the Tailscale IP address in the [`hosts` section](https://tailscale.com/kb/1337/policy-syntax#hosts).                                                                          |
| [Subnet](https://tailscale.com/kb/1019/subnets) CIDR Range                       | `192.168.1.0/24`          | Includes any IP address within the given subnet.                                                                                                                                        |
| [Tags](https://tailscale.com/kb/1068/tags)                                       | `tag:<tag-name>`          | Includes any device with the provided tag.                                                                                                                                              |
| Internet access through an [exit node](https://tailscale.com/kb/1103/exit-nodes) | `autogroup:internet`      | Includes devices with access to the internet through [exit nodes](https://tailscale.com/kb/1103/exit-nodes).                                                                            |
| Own devices                                                                      | `autogroup:self`          | Includes devices where the same user is authenticated on both the `src` and the `dst`. This does not include devices for which the user has [tags](https://tailscale.com/kb/1068/tags). |
| Tailnet devices                                                                  | `autogroup:member`        | Includes devices in the tailnet where the user is a direct member (not a shared user) of the tailnet.                                                                                   |
| Admin devices                                                                    | `autogroup:admin`         | Includes devices where the user is an [Admin](https://tailscale.com/kb/1138/user-roles#admin).                                                                                          |
| Network admin devices                                                            | `autogroup:network-admin` | Includes devices where the user is a [Network admin](https://tailscale.com/kb/1138/user-roles#network-admin).                                                                           |
| IT admin devices                                                                 | `autogroup:it-admin`      | Includes to devices where the user is an [IT admin](https://tailscale.com/kb/1138/user-roles#it-admin).                                                                                 |
| Billing admin devices                                                            | `autogroup:billing-admin` | Includes devices where the user is a [Billing admin](https://tailscale.com/kb/1138/user-roles#billing-admin).                                                                           |
| Auditor devices                                                                  | `autogroup:auditor`       | Includes devices where the user is an [Auditor](https://tailscale.com/kb/1138/user-roles#auditor).                                                                                      |
| Owner devices                                                                    | `autogroup:owner`         | Includes devices where the user is the tailnet [Owner](https://tailscale.com/kb/1138/user-roles#owner).                                                                                 |
| [IP sets](https://tailscale.com/kb/1387/ipsets)                                  | `ipset:<ip-set-name>`     | Includes all targets in the IP set.                                                                                                                                                     |

PORTS 可以是以下任意类型

| **Type** | **Description**                                        | **Example** |
| -------- | ------------------------------------------------------ | ----------- |
| Any      | Includes any port number.                              | `*`         |
| Single   | Includes a single port number.                         | `22`        |
| Multiple | Includes two or more port numbers separated by commas. | `80,443`    |
| Range    | Includes a range of port numbers.                      | `1000-2000` |

# Tailscale ACL 访问控制策略完全指南！

原文: [公众号-云原生实验室，ailscale ACL 访问控制策略完全指南！](https://mp.weixin.qq.com/s/JIbKEWJBDzP3mjWzlZ1DIA)

前面几篇文章给大家给介绍了 Tailscale 和 Headscale，包括 [👉 Headscale 的安装部署和各个平台客户端的接入，以及如何打通各个节点所在的局域网](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247504037&idx=1&sn=b059e0ed24be4ae39a25e5724700ff54&scene=21#wechat_redirect) 。同时还介绍了 [👉 如何自建私有的 DERP 服务器，并让 Tailscale 使用我们自建的 DERP 服务器](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247504288&idx=1&sn=93d74eb52ac6d1bb176c1599b3c27962&scene=21#wechat_redirect) 。

今天我们来探索一下更复杂的场景。想象有这么一个场景，我系统通过 Tailscale 方便的连接一台不完全属于我的设备， 这台设备可能还有其他人也在使用。如果我仅仅是安装一个 Tailscale， 那么所有能登录这台设备的人都可以通过 Tailscale 连接我所有的设备。

我能不能实现这样一种需求： **我可以连接这台节点，但是这台节点不能连接我的其他节点？**

这就是 Tailscale ACL（Access Control List）干的事情。ACL 可以严格限制特定用户或设备在 Tailscale 网络上访问的内容。

> ❝
> 
> 虽然 Headscale 兼容 Tailscale 的 ACL，但还是有些许差异的。 **本文所讲的 ACL 只适用于 Headscale** ，如果你使用的是官方的控制服务器，有些地方可能跟预期不符，请自行参考 Tailscale 的官方文档。

Tailscale/Headscale 的默认访问规则是 `default deny` ，也就是黑名单模式，只有在访问规则明确允许的情况下设备之间才能通信。所以 Tailscale/Headscale 默认会使用 `allowall` 访问策略进行初始化，该策略允许加入到 Tailscale 网络的所有设备之间可以相互访问。

Tailscale/Headscale 通过使用 group 这种概念，可以 **只用非常少的规则就能表达大部分安全策略** 。除了 group 之外，还可以为设备打 tag 来进一步扩展访问策略。结合 group 和 tag 就可以构建出强大的基于角色的访问控制（RBAC）策略。

关于 Tailscale 访问控制系统的详情可以参考这篇文章： 基于角色的访问控制（RBAC）：演进历史、设计理念及简洁实现 <sup>[1]</sup> 。这篇文章深入探讨了访问控制系统的历史，从设计层面分析了 `DAC -> MAC -> RBAC -> ABAC` 的演进历程及各模型的优缺点、适用场景等， 然后从实际需求出发，一步步地设计出一个实用、简洁、真正符合 RBAC 理念的访问控制系统。

## Tailscale ACL 语法

Tailscale ACL 需要保存为 HuJSON 格式，也就是 human JSON <sup>[2]</sup> 。HuJSON 是 JSON 的超集，允许添加注释以及结尾处添加逗号。这种格式更易于维护，对人类和机器都很友好。

> ❝
> 
> Headscale 除了支持 HuJSON 之外，还支持使用 YAML 来编写 ACL。本文如不作特殊说明，默认都使用 YAML 格式。

Headscale 的 ACL 策略主要包含以下几个部分：

- `acls` ：ACL 策略定义。
- `groups` ：用户的集合。Tailscale 官方控制器的“用户”指的是登录名，必须是邮箱格式。而 **Headscale 的用户就是 namesapce** 。
- `hosts` ：定义 IP 地址或者 CIDR 的别名。
- `tagOwners` ：指定哪些用户有权限给设备打 tag。
- `autoApprovers` ：允许哪些用户不需要控制端确认就可以宣告 Subnet 路由和 Exit Node。

### ACL 规则

acls 部分是 ACL 规则主体，每个规则都是一个 HuJSON 对象，它授予从一组访问来源到一组访问目标的访问权限。

所有的 ACL 规则最终表示的都是 **允许从特定源 IP 地址到特定目标 IP 地址和端口的流量** 。虽然可以直接使用 IP 地址来编写 ACL 规则，但为了可读性以及方便维护，建议使用用户、Group 以及 tag 来编写规则，Tailscale 最终会将其转换为具体的 IP 地址和端口。

每一个 ACL 访问规则长这个样子：

```
- action: accept
    src:
      - xxx
      - xxx
      - ...
    dst:
      - xxx
      - xxx
      - ...
    proto: protocol # 可选参数
```

Tailscale/Headscale 的默认访问规则是 `default deny` ，也就是黑名单模式，只有在访问规则明确允许的情况下设备之间才能通信。所以 ACL 规则中的 `action` 值一般都写 `accept` ，毕竟默认是 deny 嘛。

`src` 字段表示访问来源列表，该字段可以填的值都在这个表格里：

| 类型 | 示例 | 含义 |
| --- | --- | --- |
| Any | \* | 无限制（即所有来源） |
| 用户(Namespace) | dev1 | Headscale namespace 中的所有设备 |
| Group (ref) <sup>[3]</sup> | group:example | Group 中的所有用户 |
| Tailscale IP | 100.101.102.103 | 拥有给定 Tailscale IP 的设备 |
| Subnet CIDR (ref) <sup>[4]</sup> | 192.168.1.0/24 | CIDR 中的任意 IP |
| Hosts (ref) <sup>[5]</sup> | my-host | `hosts` 字段中定义的任意 IP |
| Tags (ref) <sup>[6]</sup> | tag:production | 分配指定 tag 的所有设备 |
| Tailnet members | autogroup:members | Tailscale 网络中的任意成员（设备） |

`proto` 字段是可选的，指定允许访问的协议。如歌不指定，默认可以访问所有 TCP 和 UDP 流量。

`proto` 可以指定为 IANA IP 协议编号 <sup>[7]</sup> 1-255（例如 16）或以下命名别名之一（例如 sctp）：

只有 TCP、UDP 和 SCTP 流量支持指定端口，其他协议的端口必须指定为 `*` 。

dst 字段表示访问目标列表，列表中的每个元素都用 `hosts:ports` 来表示。hosts 的取值范围如下：

| 类型 | 示例 | 含义 |
| --- | --- | --- |
| Any | \* | 无限制（即所有访问目标） |
| 用户（Namespace） | dev1 | Headscale namespace 中的所有设备 |
| Group (ref) <sup>[8]</sup> | group:example | Group 中的所有用户 |
| Tailscale IP | 100.101.102.103 | 拥有给定 Tailscale IP 的设备 |
| Hosts (ref) <sup>[9]</sup> | my-host | `hosts` 字段中定义的任意 IP |
| Subnet CIDR (ref) <sup>[10]</sup> | 192.168.1.0/24 | CIDR 中的任意 IP |
| Tags (ref) <sup>[11]</sup> | tag:production | 分配指定 tag 的所有设备 |
| Internet access (ref) <sup>[12]</sup> | autogroup:internet | 通过 Exit Node 访问互联网 |
| Own devices | autogroup:self | 允许 src 中定义的来源访问自己（不包含分配了 tag 的设备） |
| Tailnet devices | autogroup:members | Tailscale 网络中的任意成员（设备） |

`ports` 的取值范围：

| 类型 | 示例 |
| --- | --- |
| Any | \* |
| Single | 22 |
| Multiple | 80,443 |
| Range | 1000-2000 |

### Groups

groups 定义了一组用户的集合，YAML 格式示例配置如下：

```
groups:
  group:admin:
    - "admin1"
  group:dev:
    - "dev1"
    - "dev2"
```

huJSON 格式：

```
"groups": {
  "group:admin": ["admin1"],
  "group:dev": ["dev1", "dev2"],
},
```

每个 Group 必须以 `group:` 开头，Group 之间也不能相互嵌套。

### Autogroups

autogroup 是一个特殊的 group，它自动包含具有相同属性的用户或者访问目标，可以在 ACL 规则中调用 autogroup。

| Autogroup | 允许在 ACL 的哪个字段调用 | 含义 |
| --- | --- | --- |
| autogroup:internet | dst | 用来允许任何用户通过任意 Exit Node 访问你的 Tailscale 网络 |
| autogroup:members | src 或者 dst | 用来允许 Tailscale 网络中的任意成员（设备）访问别人或者被访问 |
| autogroup:self | dst | 用来允许 src 中定义的来源访问自己 |

示例配置：

```
acls:
  # 允许所有员工访问自己的设备
  - action: accept
    src:
      - "autogroup:members"
    dst:
      - "autogroup:self:*"
  # 允许所有员工访问打了标签 tag:corp 的设备
  - action: accept
    src:
      - "autogroup:members"
    dst:
      - "tag:corp:*"
```

### Hosts

Hosts 用来定义 IP 地址或者 CIDR 的别名，使 ACL 可读性更强。示例配置：

```
hosts:
  example-host-1: "100.100.100.100"
  example-network-1: "100.100.101.100/24
```

### Tag Owners

`tagOwners` 定义了哪些用户有权限给设备分配指定的 tag。示例配置：

```
tagOwners:
  tag:webserver:
    - group:engineering
  tag:secure-server:
    - group:security-admins
    - dev1
  tag:corp:
    - autogroup:members
```

这里表示的是允许 Group `group:engineering` 给设备添加 tag `tag:webserver` ；允许 Group `group:security-admins` 和用户（也就是 namespace）dev1 给设备添加 tag `tag:secure-server` ；允许 Tailscale 网络中的任意成员（设备）给设备添加 tag `tag:corp` 。

每个 tag 名称必须以 `tag:` 开头，每个 tag 的所有者可以是用户、Group 或者 `autogroup:members` 。

### Auto Approvers

`autoApprovers` 定义了 **无需 Headscale 控制端批准即可执行某些操作** 的用户列表，包括宣告特定的子网路由或者 Exit Node。

当然了，即使可以通过 `autoApprovers` 自动批准，Headscale 控制端仍然可以禁用路由或者 Exit Node，但不推荐这种做法，因为控制端只能临时修改， `autoApprovers` 中定义的用户列表仍然可以继续宣告路由或 Exit Node，所以正确的做法应该是修改 `autoApprovers` 中的用户列表来控制宣告的路由或者 Exit Node。

autoApprovers 示例配置：

```
autoApprovers:
  exitNode:
    - "default"
    - "tag:bar"
  routes:
    "10.0.0.0/24":
      - "group:engineering"
      - "dev1"
      - "tag:foo"
```

这里表示允许 `default` namespace 中的设备（以及打上标签 `tag:bar` 的设备）将自己宣告为 Exit Node；允许 Group `group:engineering` 中的设备（以及 dev1 namespace 中的设备和打上标签 `tag:foo` 的设备）宣告子网 `10.0.0.0/24` 的路由。

## Headscale 配置 ACL 的方法

要想在 Headscale 中配置 ACL，只需使用 HuJSON 或者 YAML 编写相应的 ACL 规则（HuJSON 格式的文件名后缀为 hujson），然后在 Headscale 的配置文件中引用 ACL 规则文件即可。

```
# Path to a file containg ACL policies.
# ACLs can be defined as YAML or HUJSON.
# https://tailscale.com/kb/1018/acls/
acl_policy_path: "./acl.yaml"
```

## ACL 规则示例

### 允许所有流量

默认的 ACL 规则允许所有访问流量，规则内容如下：

```
# acl.yaml
acls:
  - action: accept
    src:
      - "*"
    dst:
      - "*:*"
```

### 允许特定 ns 访问所有流量

假设 Headscale 有两个 namesapce： `default` 和 `guest` 。管理员的设备都在 `default` namespace 中，访客的设备都在 `guest` namespace 中。

```
$ headscale ns ls
ID | Name    | Created
1  | default | 2022-08-20 06:15:17
2  | guest   | 2022-11-27 09:20:25

$ headscale -n default node ls
ID | Hostname               | Name                            | NodeKey | Namespace | IP addresses | Ephemeral | Last seen           | Online  | Expired
2  | OpenWrt                | openwrt-njprohi0                | [7LdVc] | default   | 10.1.0.2,    | false     | 2022-08-26 04:18:43 | offline | no
5  | tailscale              | tailscale-home                  | [pwlFE] | default   | 10.1.0.5,    | false     | 2022-11-27 10:02:35 | online  | no
10 | k3s-worker05           | share                           | [5Z38M] | default   | 10.1.0.9,    | false     | 2022-11-22 18:49:25 | offline | no
11 | Galaxy a52s            | galaxy-a52s-arg5owsh            | [U+0qY] | default   | 10.1.0.1,    | false     | 2022-11-27 10:02:34 | online  | no

$ headscale -n guest node ls
ID | Hostname  | Name      | NodeKey | Namespace | IP addresses | Ephemeral | Last seen           | Online | Expired
12 | guest-1 | guest-1 | [75qSK] | guest     | 10.1.0.10,   | false     | 2022-11-27 10:05:33 | online | no
13 | guest-2 | guest-2 | [8lONp] | guest     | 10.1.0.11,   | false     | 2022-11-27 10:05:31 | online | no
```

现在我想让 `default` namespace 中的设备可以访问所有设备，而 `guest` namespace 中的设备只能访问 `guest` namespace 中的设备，那么规则应该这么写：

```
# acl.yaml
acls:
  - action: accept
    src:
      - "default"
    dst:
      - "*:*"
      - "guest:*"
  - action: accept
    src:
      - "guest"
    dst:
      - "guest:*"
```

在 `guest-1` 上查看 Tailscale 状态：

```
$ tailscale status
10.1.0.10       ks-node-2            guest        linux   -
                desktop-aoulurh-j7dfnsul.default.example.com default      windows offline
                galaxy-a52s-arg5owsh.default.example.com default      android active; relay "hs", tx 12112 rx 11988
                guest-3            guest        linux   active; direct 172.31.73.176:41641, tx 2552 rx 2440
                openwrt-njprohi0.default.example.com default      linux   offline
                tailscale-home.default.example.com default      linux   active; direct 60.184.243.56:41641, tx 3416 rx 25576
```

看起来 `guest-1` 可以看到所有的设备，但事实上它只能 ping 通 `guest-2` ，我们来验证一下：

```
$ ping 10.1.0.1
PING 10.1.0.1 (10.1.0.1) 56(84) bytes of data.
^C
--- 10.1.0.1 ping statistics ---
9 packets transmitted, 0 received, 100% packet loss, time 8169ms
```

果然是 ping 不通的。但是 10.1.0.1 这个设备是 **可以反向 ping 通** guest-1 的：

```
# 在 10.1.0.1 所在的设备操作
$ ping 10.1.0.10
PING 10.1.0.10 (10.1.0.10) 56(84) bytes of data.
64 bytes from 10.1.0.10: icmp_seq=1 ttl=64 time=68.9 ms
64 bytes from 10.1.0.10: icmp_seq=2 ttl=64 time=91.5 ms
64 bytes from 10.1.0.10: icmp_seq=3 ttl=64 time=85.3 ms
64 bytes from 10.1.0.10: icmp_seq=4 ttl=64 time=79.7 ms
^C
--- 10.1.0.10 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 68.967/81.389/91.551/8.306 ms
```

ssh 测试一下：

```
$ ssh root@10.1.0.10
root@10.1.0.10's password:
```

完美。

下面再来看看 `guest-1` 能不能 ping 通 `guest-2` ：

```
# 在 guest-1 设备上操作
$ ping 10.1.0.11
PING 10.1.0.11 (10.1.0.11) 56(84) bytes of data.
64 bytes from 10.1.0.11: icmp_seq=1 ttl=64 time=2.93 ms
64 bytes from 10.1.0.11: icmp_seq=2 ttl=64 time=1.33 ms
^C
--- 10.1.0.11 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1001ms
rtt min/avg/max/mdev = 1.325/2.128/2.931/0.803 ms
```

和我在上面预测的效果一样，ACL 规则生效了。

### 神奇的 tag

tag 有一个非常神奇的功效：它可以让 `src` 和 `dst` 中的元素失效。具体什么意思呢？ **假设你的 src 或 dst 中指定了 namespace 或者 group，那么这个规则只对这个 namespace 或者 group 中（没有分配 tag 的设备）生效。**

举个例子你就明白了，现在我给 guest-2 打上一个 tag：

```
$ headscale node tag -i 13 -t tag:test
Machine updated

$ headscale -n guest node ls -t
ID | Hostname  | Name      | NodeKey | Namespace | IP addresses | Ephemeral | Last seen           | Online | Expired | ForcedTags | InvalidTags | ValidTags
12 | ks-node-2 | ks-node-2 | [75qSK] | guest     | 10.1.0.10,   | false     | 2022-11-27 10:18:35 | online | no      |            |             |
13 | ks-node-3 | ks-node-3 | [8lONp] | guest     | 10.1.0.11,   | false     | 2022-11-27 10:18:31 | online | no      | tag:test   |             |
```

此时 guest-1 就 ping 不通 guest-2 了：

```
# 在 guest-1 设备上操作
$ ping 10.1.0.11
PING 10.1.0.11 (10.1.0.11) 56(84) bytes of data.
^C
--- 10.1.0.11 ping statistics ---
4 packets transmitted, 0 received, 100% packet loss, time 3070ms
```

这就说明 guest-2 并不包含在 `guest:*` 这个访问目标中，也就是说打了 tag 的设备并不包含在 `guest:*` 这个访问目标中。

此时其他设备如果还想继续 guest-2，必须在 dst 中指定 `tag:test` ：

```
acls:
  - action: accept
    src:
      - "default"
    dst:
      - "*:*"
      - "guest:*"
      - "tag:test:*"
  - action: accept
    src:
      - "guest"
    dst:
      - "guest:*"
      - "tag:test:*"
```

再次测试访问：

```
# 在 guest-1 设备上操作
$ ping 10.1.0.11
PING 10.1.0.11 (10.1.0.11) 56(84) bytes of data.
64 bytes from 10.1.0.11: icmp_seq=1 ttl=64 time=1.31 ms
64 bytes from 10.1.0.11: icmp_seq=2 ttl=64 time=3.40 ms
^C
--- 10.1.0.11 ping statistics ---
2 packets transmitted, 2 received, 0% packet loss, time 1002ms
rtt min/avg/max/mdev = 1.314/2.355/3.397/1.041 ms
```

果然可以 ping 通了。

## 总结

Tailscale/Headscale 的 ACL 非常强大，你可以基于 ACL 实现各种各样的访问控制策略，本文只是给出了几个关键示例，帮助大家理解其用法，更多功能大家可以自行探索（比如 group 等）。下篇文章将会给大家介绍如何配置 Headscale 的 Exit Node，以及各个设备如何使用 Exit Node，届时会用到 ACL 里面的 `autoApprovers` ，敬请期待！

### 引用链接

\[1\]

基于角色的访问控制（RBAC）：演进历史、设计理念及简洁实现: *http://arthurchiao.art/blog/rbac-as-it-meant-to-be-zh/*

\[2\]

human JSON: *https://github.com/tailscale/hujson*

\[3\]

(ref): *https://tailscale.com/kb/1018/acls/#groups*

\[4\]

(ref): *https://tailscale.com/kb/1019/subnets*

\[5\]

(ref): *https://tailscale.com/kb/1018/acls/#hosts*

\[6\]

(ref): *https://tailscale.com/kb/1068/acl-tags*

\[7\]

IANA IP 协议编号: *https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml*

\[8\]

(ref): *https://tailscale.com/kb/1018/acls/#groups*

\[9\]

(ref): *https://tailscale.com/kb/1018/acls/#hosts*

\[10\]

(ref): *https://tailscale.com/kb/1019/subnets*

\[11\]

(ref): *https://tailscale.com/kb/1068/acl-tags*

\[12\]

(ref): *https://tailscale.com/kb/1103/exit-nodes*
