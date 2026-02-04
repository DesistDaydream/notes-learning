---
title: Linux DNS 管理
linkTitle: Linux DNS 管理
weight: 3
---

# 概述

> 参考：
>
> - [Manual(手册)，resolver(3)](https://man7.org/linux/man-pages/man3/resolver.3.html)
> - [Manual(手册)，resolv.conf(5)](https://man7.org/linux/man-pages/man5/resolv.conf.5.html)
>   - resolver(5) 手册也指向了 resolv.conf(5)

在 Linux 中，进行域名解析工作的是 Reslover(解析器)。

**Reslover(解析器)** 是 [Linux libc 库](/docs/1.操作系统/Linux%20源码解析/Linux%20libc%20库.md) 中用于提供 DNS 接口的程序集（其实就是对外暴露了 API 的 C 库），当某个进程调用这些程序时将同时读入 Reslover 的配置文件（resolv.conf），这个文件具有可读性并且包含大量可用的解析参数。

> Note：并不是所有程序都会使用 Linux Reslover。比如想要测试 resolv.conf 文件，不要使用 dig, host, nslook 这类工具，因为他们并没有调用 Resolver 的库(i.e. resolv.conf 文件中的 option 内的设置不会生效)。可以使用 getent 来测试。一般情况下正常的应用程序，都会调用 resolver，并使用 resolv.conf 文件(比如 ping 程序)。

Reslover 库定义了下面这些函数

```c
#include <netinet/in.h>
#include <arpa/nameser.h>
#include <resolv.h>

struct __res_state;
typedef struct __res_state *res_state;

int res_ninit(res_state statep);

void res_nclose(res_state statep);

int res_nquery(res_state statep,
    const char *dname, int class, int type,
    unsigned char answer[.anslen], int anslen);

int res_nsearch(res_state statep,
    const char *dname, int class, int type,
    unsigned char answer[.anslen], int anslen);

int res_nquerydomain(res_state statep,
    const char *name, const char *domain,
    int class, int type, unsigned char answer[.anslen],
    int anslen);

int res_nmkquery(res_state statep,
    int op, const char *dname, int class,
    int type, const unsigned char data[.datalen], int datalen,
    const unsigned char *newrr,
    unsigned char buf[.buflen], int buflen);

int res_nsend(res_state statep,
    const unsigned char msg[.msglen], int msglen,
    unsigned char answer[.anslen], int anslen);

int dn_comp(const char *exp_dn, unsigned char comp_dn[.length],
    int length, unsigned char **dnptrs,
    unsigned char **lastdnptr);

int dn_expand(const unsigned char *msg,
    const unsigned char *eomorig,
    const unsigned char *comp_dn, char exp_dn[.length],
    int length);
```

这些函数会暴露成 API 以供调用。比如 Go 编写的程序可以使用纯 Go 自己的解析器，也可以使用 cgo 来调用 glibc 的这个 Linux Resolver。

> tips: 这些函数名前面带个 n 表示这是新的函数，老的函数没有 n，已经被弃用了。

- res_ninit # 实例化并读取配置
  - 每次调用 res_ninit 都需要调用 res_nclose 以释放由 res_ninit 和后续 res_nquery 分配的内存
- res_nquery # 查询
- res_nsearch # 搜索
- res_nquerydomain # 串联查询

# Linux DNS 关联文件与配置

**/etc/resolv.conf** # Reslover 的配置文件

**/etc/hosts** # 更改本地主机名和 IP 的对应关系，用于解析指定域名

例：当 ping TEST-1 时，则 ping 192.168.2.3

```bash
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
192.168.2.3 TEST-1
```

**/PATH/TO/nsswitch.conf** # 名称服务切换配置。GUN C 库(glibc) 和 某些其他应用程序使用该配置文件来确定从哪些地方获取解析信息。比如是否要读取 /etc/hosts 文件

- 该文件属于 glibc 包中的一部分。但是由于 CentOS 与 Ubuntu 中 glibc 的巨大差异，该文件所在路径也不同：
  - CentOS 在 **/etc/nsswitch.conf**
  - Ubuntu 在 **/usr/share/libc-bin/nsswitch.conf**
