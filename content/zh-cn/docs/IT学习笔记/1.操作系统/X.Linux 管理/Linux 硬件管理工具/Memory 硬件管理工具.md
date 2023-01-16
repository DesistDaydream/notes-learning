---
title: Memory 硬件管理工具
---

# 概述

> ## 参考：

# 看内存的品牌及型号

背景：想加个内存，Mysql 服务器上的内存不够了，就算可以插(插槽都有，现在是 4 条 4G 内存，还有空闲八个槽。)，不知道兼不兼容，否则不稳定的,不兼容的话，死的更惨,这个不一定了，停产了，就没有办法了。怎么办?得看内存的品牌及型号。

```bash
[root@jackxiang ~]# rpm -qa|grep dmidecode
dmidecode-2.11-2.el6.x86_64
[root@jackxiang ~]# dmidecode
```

查看服务器型号、序列号：

```bash
[root@jackxiang ~]# dmidecode|grep "System Information" -A9|egrep  "Manufacturer|Product|Serial"
        Manufacturer: VMware, Inc.
        Product Name: VMware Virtual Platform
        Serial Number: VMware-42 18 c8 32 77 c6 ec 16-3f 31 94 e9 d0 34 a6 ac
```

Linux 查看内存的插槽数,已经使用多少插槽.每条内存多大：

    [root@jackxiang ~]# dmidecode|grep -A5 "Memory Device"|grep Size|grep -v Range
            Size: 4096 MB
            Size: 2048 MB
            Size: No Module Installed
            Size: No Module Installed

Linux 查看内存的频率：

    [root@localhost htdocs]# dmidecode|grep -A16 "Memory Device"|grep 'Speed'
            Speed: 667 MHz (1.5 ns)
            Speed: 667 MHz (1.5 ns)
            Speed: 667 MHz (1.5 ns)
            Speed: 667 MHz (1.5 ns)
            Speed: Unknown

在 linux 查看内存型号的命令：

    dmidecode -t memory

查看主板型号：

    dmidecode |grep -A16 "System Information$"

内存槽及内存条：

    dmidecode |grep -A16 "Memory Device$"

硬盘：

    fdisk -l
    smartctl -a /dev/sda

网卡：

    mii-tool

<!---->

    dmidecode|grep -P 'Maximum\s+Capacity'    //最大支持几G内存
    # dmidecode  |grep  "Product Name"  //查看服务器品牌和型号
    # dmidecode|grep -P -A5 "Memory\s+Device"|grep Size|grep -v Range       //总共几个插槽，已使用几个插槽
    Size: 1024 MB       //此插槽有1根1G内存
    Size: 1024 MB       //此插槽有1根1G内存
    Size: 1024 MB       //此插槽有1根1G内存
    Size: 1024 MB       //此插槽有1根1G内存
    Size: No Module Installed       //此插槽未使用
    Size: No Module Installed       //此插槽未使用

<!---->

    # dmidecode -t 17        //数字17是dmidecode的参数，本文最后有其他数字参数

<!---->

    # dmidecode 2.7
    SMBIOS 2.4 present.
    Handle 0x0015, DMI type 17, 27 bytes.
    Memory Device
      Array Handle: 0x0013
      Error Information Handle: Not Provided
      Total Width: 72 bits
      Data Width: 64 bits
      Size: 2048 MB 【插槽1有1条2GB内存】
      Form Factor: DIMM
      Set: None
      Locator: DIMM00
      Bank Locator: BANK
      Type: Other
      Type Detail: Other
      Speed: 667 MHz (1.5 ns)
      Manufacturer:
      Serial Number: BZACSKZ001
      Asset Tag: RAM82
      Part Number: MT9HTF6472FY-53EA2
    Handle 0x0017, DMI type 17, 27 bytes.
    Memory Device
      Array Handle: 0x0013
      Error Information Handle: Not Provided
      Total Width: 72 bits
      Data Width: 64 bits
      Size: 2048 MB 【插槽2有1条2GB内存】
      Form Factor: DIMM
      Set: None
      Locator: DIMM10
      Bank Locator: BANK
      Type: Other
      Type Detail: Other
      Speed: 667 MHz (1.5 ns)
      Manufacturer:
      Serial Number: BZACSKZ001
      Asset Tag: RAM83
      Part Number: MT9HTF6472FY-53EA2
    Handle 0x0019, DMI type 17, 27 bytes.
    Memory Device
      Array Handle: 0x0013
      Error Information Handle: Not Provided
      Total Width: 72 bits
      Data Width: 64 bits
      Size: 2048 MB 【插槽3有1条2GB内存】
      Form Factor: DIMM
      Set: None
      Locator: DIMM20
      Bank Locator: BANK
      Type: Other
      Type Detail: Other
      Speed: 667 MHz (1.5 ns)
      Manufacturer:
      Serial Number: BZACSKZ001
      Asset Tag: RAM84
      Part Number: MT9HTF6472FY-53EA2
    Handle 0x001B, DMI type 17, 27 bytes.
    Memory Device
      Array Handle: 0x0013
      Error Information Handle: Not Provided
      Total Width: 72 bits
      Data Width: 64 bits
      Size: 2048 MB 【插槽4有1条2GB内存】
      Form Factor: DIMM
      Set: None
      Locator: DIMM30
      Bank Locator: BANK
      Type: Other
      Type Detail: Other
      Speed: 667 MHz (1.5 ns)
      Manufacturer:
      Serial Number: BZACSKZ001
      Asset Tag: RAM85
      Part Number: MT9HTF6472FY-53EA2

实践来源：

<http://www.jbxue.com/LINUXjishu/10053.html>

<http://www.linuxsir.org/bbs/thread309696.html>

<http://xclinux.diandian.com/post/2013-04-16/40049844350>作者：justwinit@向东博客 专注 WEB 应用 构架之美 --- 构架之美，在于尽态极妍 | 应用之美，在于药到病除

地址：<http://www.justwinit.cn/post/7496/>
