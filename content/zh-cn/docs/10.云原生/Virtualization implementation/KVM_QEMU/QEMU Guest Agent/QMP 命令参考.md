---
title: QMP 命令参考
---

# 概述

> 参考：
>
> - [官方文档，系统模拟管理与交互-QEMU Guest Agent 协议参考](https://www.qemu.org/docs/master/interop/qemu-ga-ref.html)
> - [简书，qemu-agent-command 命令含义](https://www.jianshu.com/p/27d8491ed100)

通过指令

    virsh qemu-agent-command 虚拟机 --cmd '{"execute":"guest-info"}'

可以查看其所有支持的命令，返回的数据如下

    {"return":{"version":"2.8.0",
    "supported_commands":  [
    {"enabled":true,"name":"guest-sync-delimited","success-response":true},{"enabled":true,"name":"guest-sync","success-response":true},
    {"enabled":true,"name":"guest-suspend-ram","success-response":false},
    {"enabled":true,"name":"guest-suspend-hybrid","success-response":false},
    {"enabled":true,"name":"guest-suspend-disk","success-response":false},
    {"enabled":true,"name":"guest-shutdown","success-response":false},
    {"enabled":true,"name":"guest-set-vcpus","success-response":true},
    {"enabled":true,"name":"guest-set-user-password","success-response":true},
    {"enabled":true,"name":"guest-set-time","success-response":true},
    {"enabled":true,"name":"guest-set-memory-blocks","success-response":true},
    {"enabled":true,"name":"guest-ping","success-response":true},
    {"enabled":true,"name":"guest-network-get-interfaces","success-response":true},
    {"enabled":true,"name":"guest-info","success-response":true},
    {"enabled":true,"name":"guest-get-vcpus","success-response":true},
    {"enabled":true,"name":"guest-get-time","success-response":true},
    {"enabled":true,"name":"guest-get-memory-blocks","success-response":true},
    {"enabled":true,"name":"guest-get-memory-block-info","success-response":true},
    {"enabled":true,"name":"guest-get-fsinfo","success-response":true},
    {"enabled":true,"name":"guest-fstrim","success-response":true},
    {"enabled":true,"name":"guest-fsfreeze-thaw","success-response":true},
    {"enabled":true,"name":"guest-fsfreeze-status","success-response":true},
    {"enabled":true,"name":"guest-fsfreeze-freeze-list","success-response":true},
    {"enabled":true,"name":"guest-fsfreeze-freeze","success-response":true},
    {"enabled":false,"name":"guest-file-write","success-response":true},
    {"enabled":false,"name":"guest-file-seek","success-response":true},
    {"enabled":false,"name":"guest-file-read","success-response":true},
    {"enabled":false,"name":"guest-file-open","success-response":true},
    {"enabled":false,"name":"guest-file-flush","success-response":true},
    {"enabled":false,"name":"guest-file-close","success-response":true},
    {"enabled":false,"name":"guest-exec-status","success-response":true},
    {"enabled":false,"name":"guest-exec","success-response":true}
    ]}}

返回为数据，其中 supported_command 为所有命令的数组

其官方地址为：[QEMU Guest Agent Protocol Reference](https://links.jianshu.com/go?to=https%3A%2F%2Fqemu.weilnetz.de%2Fdoc%2Fqemu-ga-ref.html%23API-Reference)

各命令含义如下：
1\. guest-sync-delimited

宿主机发送一个 int 数字给 qga，qga 返回这个数字，并且在后续返回字符串响应中加入 ascii 码为 0xff 的字符。
2\. guest-sync

回文唯一的整数，这个命令进行测试。

这个命令用于确保 client 与 guest agent 是同步的，不包含之前 client 的陈旧数据。直到返回特定的数字之前的 guest agent 响应都应被忽略。当含有 client 接收到陈旧数据时，这个命令并不能可靠的执行。一个特定的场景是，如果 qemu-ga 响应被逐个字符地输入到 JSON 解析器中。在这些情况下，使用 guest-sync-delimited 可能是最佳选择。对于逐行获取响应并将其转换为 JSON 对象的客户机，guest-sync 应该足够了，但请注意，在通道不干净的情况下，一些解析响应的尝试可能会导致解析器错误。此类客户端还应该在此命令之前加上 0xFF 字节，以确保客户代理刷新前一个会话的部分读取的 JSON 数据。

**Arguments:**
_id: int,随机生成的 64-bit 整数_

**Returns:**
_客户端发出的特定整数_

**测试:**

    virsh qemu-agent-command centos --cmd '{"execute":"guest-sync", "arguments":{"id":1234567890}}'
    {"return":1234567890}

3. guest-ping
   Ping the guest agent，如果不返回错误信息，则成功

**测试:**

    virsh qemu-agent-command centos --cmd
    '{"execute":"guest-ping"}'
    {"return":{}}

4. guest-get-time
   获取虚拟机系统时间（相对于 1970-01-01 in UTC)；

**Returns:**
_纳秒格式的时间_

**测试:**

    virsh qemu-agent-command centos --cmd
    '{"execute":"guest-get-time"}'
    {"return":1534345952638400000}

5. guest-set-time
   设置虚拟机时间

Arguments:
time: int (optional)

时间格式为纳秒，相对于 1970-01-01 in UTC

**Returns:**
_成功则无返回值_

6. guest-info
   获取 guest agent 信息

**Returns:**
_GuestAgentInfo 对象_

7. guest-shutdown
   开启虚拟机关机任务，其为异步命令，不保证关机成功

Arguments:

_mode: string (optional)_

_"halt"， "powerdown" ， "reboot"三种状态可以选择，powerdown 为默认选项，命令成功执行无返回。成功的标志是，VM 以 0 的推出状态推出，或使用 QMP 命令查询时返回 VM 状态为 shutdown_

8. guest-file-open

打开虚拟机内文件并返回文件句柄

Arguments:

path: string，虚拟机所打开文件完整路径

mode: string (optional)，打开文件方式，与 fopen()函数相同，默认为"r"

**Returns:**

执行成功则返回文件句柄

9. guest-file-close

关闭虚拟机文件

Arguments:

_handle: int，guest-file-open 所返回的文件句柄_

**Returns:**

_成功无返回值_

10. guest-file-read

读取虚拟机中打开的文件（Data will be base64-encoded）

Arguments:

handle: int，guest-file-open 所返回的文件句柄

count：int，最少读取位数（默认为 64K）

**Returns:**

_成功则返回 GuestFileRead 类_

11. guest-file-write

写入虚拟机打开的文件

Arguments:

handle: int，guest-file-open 所返回的文件句柄

count：int，最少读取位数（默认为 64K）

buf-b64: string，表示要写入数据的 base64 编码字符串

count: int (optional)，写入的位数，默认是在 buffer 中的全部位数

**Returns:**

_成功返回 GuestFileWrite 类_

12. guest-file-seek

同 fseek()用法相同，seek 到文件的指定位置，

Arguments:

handle: int，guest-file-open 所返回的文件句柄

offset：int，文件位移量

whence: GuestFileWhence，描述 offset

**Returns:**

成功则返回 GuestFileSeek 类

**13. guest-file-flush**

将用户缓冲区数据写入磁盘或内核缓冲区

Arguments:

_handle: int，guest-file-open 所返回的文件句柄_

**Returns:**

成功则无返回值\*

14. guest-fsfreeze-status

获取虚拟机文件冻结状态

**Returns:**

_GuestFsfreezeStatus 枚举，包括 thawed，frozen 两种状态_

15. guest-fsfreeze-freeze

同步并冻结虚拟机文件系统

**Returns:**

返回目前冻结的文件个数，如果执行错误，则解冻当前所有文件。

16. guest-fsfreeze-freeze-list

同步和冻结指定的虚拟机文件，

Arguments:

_mountpoints: array of string (optional)，要冻结的文件系统挂载点数组。如果省略，每个挂载的文件系统都会被冻结。无效的挂载点被忽略。_

17. guest-fsfreeze-thaw

解冻所有冻结的文件

**Returns:**

_解冻的文件个数_

18. guest-fstrim

文件系统未使用的硬盘空间

Arguments：

minimum: int (optional)

最小可丢弃的连续自由范围，单位为字节。通过增加这个值，fstrim 操作将更快地完成具有严重碎片化的空闲空间的文件系统，尽管并非所有块都将被丢弃。默认值为零，意思是“丢弃所有空闲块”。

19. guest-suspend-disk

挂起虚拟机磁盘，如成功则不返回值

20. guest-suspend-ram

挂起虚拟机 ram

21. guest-suspend-hybrid

将虚拟机状态写入磁盘，并挂起 ram

## guest-network-get-interfaces # 获取虚拟机 IP 地址，MAC 地址，子网掩码

{"execute": "guest-network-get-interfaces"}

23. guest-get-vcpus
    检索客户的逻辑处理器列表。这是一个只读操作。

**Returns:**

_虚拟机的 VCPUs 列表，以 GuestLogicalProcessor 类形式返回_

24. guest-set-vcpus

尝试重新配置客户内部的逻辑处理器(当前:启用/禁用)。

Arguments:

_vcpus: array of GuestLogicalProcessor_

25. guest-get-fsinfo

获取在虚拟机中挂载的文件系统列表

26. guest-set-user-password

Arguments:

username: string，需要更改密码的用户名

password: string，新的密码（base64 encoded）

crypted: boolean，如果以被 crypt()加密则为真，否则为 false

**Returns:**

_如成功则无返回值_

27. guest-get-memory-block

获取虚拟机内存块信息，返回虚拟机所知的所有内存块，以 GuestMemoryBlock 对象展示

28. guest-set-memory-blocks

设置虚拟机中的内存块信息

## guest-exec-status # 获取虚拟机中的进程状态，如进程退出，则获取其相关元数据

Arguments：

- pid: int

    { "execute": "guest-exec-status", "arguments": { "pid": PID } }

**Returns:**
成功则返回 GuestExecStatus 类对象,GuestExecStatus 含有如下成员

> exited: boolean，如进程已经终止则为真
> exitcode: int (optional)，进程退出码
> signal: int (optional)，异常终止代码
> out-data: string (optional)，程序 stdout(base64-encoded)
> err-data: string (optional)，程序 stderr(base64-encoded)
> out-truncated: boolean (optional)，如果由于大小限制而未完全捕获 stdout，则为真。

err-truncated: boolean (optional)，如果由于大小限制而没有完全捕获 stderr，则为真。

## guest-exec # 在虚拟机中执行命令

**Arguments:**

- path: string，执行的路径或名称
- arg: array of string (optional)，执行命令所需参数
- env: array of string (optional)，执行所需的环境变量
- input-data: string (optional)，所需数据
- capture-output: boolean (optional)，获取进程的 stdout/stderr

    { "execute": "guest-exec", "arguments": { "path": "ip", "arg": [ "addr", "list" ], "capture-output": true } }

**Returns:**
如执行成功则返回其 PID

## guest-get-host-name # 返回机器名称

## guest-get-timezone # 获取虚拟机时区信息

## guest-get-osinfo # 获取操作系统信息

# 应用示例

以 virsh qemu-agent-command 命令为例，通过 socat 等工具与 VM 交互，只需要直接输入 QMP 指令即可

## 在 VM 中执行命令，并在宿主机接收执行结果

    # 在 VM 中执行命令，并返回该命令 PID
    [root@host-3 ~]# virsh qemu-agent-command desistdaydream.bj-net --pretty '{ "execute": "guest-exec", "arguments": { "path": "ip", "arg": [ "addr", "list" ], "capture-output": true } }'
    {
      "return": {
        "pid": 1826
      }
    }

    # 通过 PID 获取命令输出结果，这个结果是 base64 编码的。
    [root@host-3 ~]# virsh qemu-agent-command desistdaydream.bj-net --pretty '{ "execute": "guest-exec-status", "arguments": { "pid": 1826 } }'
    {
      "return": {
        "exitcode": 0,
        "out-data": "MTogbG86IDxMT09QQkFDSyxVUCxMT1dFUl9VUD4gbXR1IDY1NTM2IHFkaXNjIG5vcXVldWUgc3RhdGUgVU5LTk9XTiBncm91cCBkZWZhdWx0IHFsZW4gMTAwMAogICAgbGluay9sb29wYmFjayAwMDowMDowMDowMDowMDowMCBicmQgMDA6MDA6MDA6MDA6MDA6MDAKICAgIGluZXQgMTI3LjAuMC4xLzggc2NvcGUgaG9zdCBsbwogICAgICAgdmFsaWRfbGZ0IGZvcmV2ZXIgcHJlZmVycmVkX2xmdCBmb3JldmVyCiAgICBpbmV0NiA6OjEvMTI4IHNjb3BlIGhvc3QgCiAgICAgICB2YWxpZF9sZnQgZm9yZXZlciBwcmVmZXJyZWRfbGZ0IGZvcmV2ZXIKMjogZW5zMzogPEJST0FEQ0FTVCxNVUxUSUNBU1QsVVAsTE9XRVJfVVA+IG10dSAxNTAwIHFkaXNjIGZxX2NvZGVsIHN0YXRlIFVQIGdyb3VwIGRlZmF1bHQgcWxlbiAxMDAwCiAgICBsaW5rL2V0aGVyIDUyOjU0OjAwOjZkOmZhOmYwIGJyZCBmZjpmZjpmZjpmZjpmZjpmZgogICAgaW5ldCAxNzIuMTkuNDIuMjQ4LzI0IGJyZCAxNzIuMTkuNDIuMjU1IHNjb3BlIGdsb2JhbCBub3ByZWZpeHJvdXRlIGVuczMKICAgICAgIHZhbGlkX2xmdCBmb3JldmVyIHByZWZlcnJlZF9sZnQgZm9yZXZlcgo0OiBkb2NrZXIwOiA8Tk8tQ0FSUklFUixCUk9BRENBU1QsTVVMVElDQVNULFVQPiBtdHUgMTUwMCBxZGlzYyBub3F1ZXVlIHN0YXRlIERPV04gZ3JvdXAgZGVmYXVsdCAKICAgIGxpbmsvZXRoZXIgMDI6NDI6NWU6MjQ6Mjg6YmQgYnJkIGZmOmZmOmZmOmZmOmZmOmZmCiAgICBpbmV0IDEwLjM4LjAuMS8yNCBicmQgMTAuMzguMC4yNTUgc2NvcGUgZ2xvYmFsIGRvY2tlcjAKICAgICAgIHZhbGlkX2xmdCBmb3JldmVyIHByZWZlcnJlZF9sZnQgZm9yZXZlcgo=",
        "exited": true
      }
    }
    # 使用 base64 将数据解码
    [root@host-3 ~]# echo "MTogbG86IDxMT09QQkFDSyxVUCxMT1dFUl9VUD4gbXR1IDY1NTM2IHFkaXNjIG5vcXVldWUgc3RhdGUgVU5LTk9XTiBncm91cCBkZWZhdWx0IHFsZW4gMTAwMAogICAgbGluay9sb29wYmFjayAwMDowMDowMDowMDowMDowMCBicmQgMDA6MDA6MDA6MDA6MDA6MDAKICAgIGluZXQgMTI3LjAuMC4xLzggc2NvcGUgaG9zdCBsbwogICAgICAgdmFsaWRfbGZ0IGZvcmV2ZXIgcHJlZmVycmVkX2xmdCBmb3JldmVyCiAgICBpbmV0NiA6OjEvMTI4IHNjb3BlIGhvc3QgCiAgICAgICB2YWxpZF9sZnQgZm9yZXZlciBwcmVmZXJyZWRfbGZ0IGZvcmV2ZXIKMjogZW5zMzogPEJST0FEQ0FTVCxNVUxUSUNBU1QsVVAsTE9XRVJfVVA+IG10dSAxNTAwIHFkaXNjIGZxX2NvZGVsIHN0YXRlIFVQIGdyb3VwIGRlZmF1bHQgcWxlbiAxMDAwCiAgICBsaW5rL2V0aGVyIDUyOjU0OjAwOjZkOmZhOmYwIGJyZCBmZjpmZjpmZjpmZjpmZjpmZgogICAgaW5ldCAxNzIuMTkuNDIuMjQ4LzI0IGJyZCAxNzIuMTkuNDIuMjU1IHNjb3BlIGdsb2JhbCBub3ByZWZpeHJvdXRlIGVuczMKICAgICAgIHZhbGlkX2xmdCBmb3JldmVyIHByZWZlcnJlZF9sZnQgZm9yZXZlcgo0OiBkb2NrZXIwOiA8Tk8tQ0FSUklFUixCUk9BRENBU1QsTVVMVElDQVNULFVQPiBtdHUgMTUwMCBxZGlzYyBub3F1ZXVlIHN0YXRlIERPV04gZ3JvdXAgZGVmYXVsdCAKICAgIGxpbmsvZXRoZXIgMDI6NDI6NWU6MjQ6Mjg6YmQgYnJkIGZmOmZmOmZmOmZmOmZmOmZmCiAgICBpbmV0IDEwLjM4LjAuMS8yNCBicmQgMTAuMzguMC4yNTUgc2NvcGUgZ2xvYmFsIGRvY2tlcjAKICAgICAgIHZhbGlkX2xmdCBmb3JldmVyIHByZWZlcnJlZF9sZnQgZm9yZXZlcgo=" | base64 -d
    1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
        link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
        inet 127.0.0.1/8 scope host lo
           valid_lft forever preferred_lft forever
        inet6 ::1/128 scope host
           valid_lft forever preferred_lft forever
    2: ens3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
        link/ether 52:54:00:6d:fa:f0 brd ff:ff:ff:ff:ff:ff
        inet 172.19.42.248/24 brd 172.19.42.255 scope global noprefixroute ens3
           valid_lft forever preferred_lft forever
    4: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default
        link/ether 02:42:5e:24:28:bd brd ff:ff:ff:ff:ff:ff
        inet 10.38.0.1/24 brd 10.38.0.255 scope global docker0
           valid_lft forever preferred_lft forever
