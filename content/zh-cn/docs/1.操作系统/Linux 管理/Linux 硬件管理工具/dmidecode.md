---
title: "dmidecode"
linkTitle: "dmidecode"
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://www.nongnu.org/dmidecode/)
> - [Wiki，dmidecode](https://en.wikipedia.org/wiki/Dmidecode)
> - Manual(手册)，dmidecode

dmidecode 命令可以让我们在 Linux 系统下获取有关硬件方面的信息。dmidecode 的作用是将 [DMI](/docs/Standard/IT/DMTF.md#DMI) 数据库中的信息解码，以可读的文本方式显示。由于 DMI 信息可以人为修改，因此里面的信息不一定是系统准确的信息。dmidecode 遵循 SMBIOS/DMI 标准，其输出的信息包括BIOS、系统、主板、处理器、内存、缓存等等。

dmidecode 附带三个额外的工具：

- biosdecode 打印它能找到的所有 BIOS 相关信息（参见[示例输出](https://www.nongnu.org/dmidecode/sample/biosdecode.txt)）；
- ownership 检索可以在 Compaq 计算机上设置的“所有权标签”；
- vpddecode 打印可以在几乎所有 IBM 计算机中找到的“重要产品数据”信息（参见 [示例输出](https://www.nongnu.org/dmidecode/sample/vpddecode.txt)）。

## DMI 类型

SMBIOS规范定义了以下DMI类型:

```
Type   Information
 ────────────────────────────────────────────
    0   BIOS
    1   System
    2   Baseboard
    3   Chassis
    4   Processor
    5   Memory Controller
    6   Memory Module
    7   Cache
    8   Port Connector
    9   System Slots
   10   On Board Devices
   11   OEM Strings
   12   System Configuration Options
   13   BIOS Language
   14   Group Associations
   15   System Event Log
   16   Physical Memory Array
   17   Memory Device
   18   32-bit Memory Error
   19   Memory Array Mapped Address
   20   Memory Device Mapped Address
   21   Built-in Pointing Device
   22   Portable Battery
   23   System Reset
   24   Hardware Security
   25   System Power Controls
   26   Voltage Probe
   27   Cooling Device
   28   Temperature Probe
   29   Electrical Current Probe
   30   Out-of-band Remote Access
   31   Boot Integrity Services
   32   System Boot
   33   64-bit Memory Error
   34   Management Device
   35   Management Device Component
   36   Management Device Threshold Data
   37   Memory Channel
   38   IPMI Device
   39   Power Supply
   40   Additional Information
   41   Onboard Devices Extended Information
   42   Management Controller Host Interface

 Additionally,  type 126 is used for disabled entries and type 127 is an end-of-table marker. Types 128 to 255 are for OEM-specific data.  dmidecode will display these entries by
 default, but it can only decode them when the vendors have contributed documentation or code for them.

 Keywords can be used instead of type numbers with --type.  Each keyword is equivalent to a list of type numbers:

 Keyword     Types
 ──────────────────────────────
 bios        0, 13
 system      1, 12, 15, 23, 32
 baseboard   2, 10, 41
 chassis     3
 processor   4
 memory      5, 6, 16, 17
 cache       7
 connector   8
 slot        9

 Keywords are matched case-insensitively. The following command lines are equivalent:

 • dmidecode --type 0 --type 13
 • dmidecode --type 0,13
 • dmidecode --type bios
 • dmidecode --type BIOS
```

# Syntax(语法)

**dmidecode [OPTIONS]**

**OPTIONS**

- **-t, --type STRING** # 只显示指定类型的条目。可用的值可以使 DMI 类型编号、以逗号分隔的类型编号列表、关键字(详见 [DMI 类型](#DMI%20类型))

# EXAMPLE

## 服务器信息

```bash
~]# dmidecode | grep "System Information" -A9|egrep  "Manufacturer|Product|Serial"
Manufacturer: VMware, Inc.
Product Name: VMware Virtual Platform
Serial Number: VMware-42 18 c8 32 77 c6 ec 16-3f 31 94 e9 d0 34 a6 ac
```

## 主板信息

```bash
dmidecode |grep -A16 "System Information$"
```

## CPU 信息

```bash
dmidecode -t processor
```

## 内存信息

查看内存信息

```bash
dmidecode -t memory
```

查看内存的插槽数,已经使用多少插槽.每条内存多大：

```bash
~]# dmidecode -t memory | grep Size | grep -v Range
Size: 4096 MB
Size: 2048 MB
Size: No Module Installed
Size: No Module Installed
```

Linux 查看内存的频率：

```bash
~]# dmidecode|grep -A16 "Memory Device"|grep 'Speed'
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: Unknown
```

在 linux 查看内存型号的命令：

内存槽及内存条：

```bash
dmidecode |grep -A16 "Memory Device$"
```

```bash
dmidecode|grep -P 'Maximum\s+Capacity'    //最大支持几G内存
dmidecode|grep -P -A5 "Memory\s+Device"|grep Size|grep -v Range       //总共几个插槽，已使用几个插槽
Size: 1024 MB       //此插槽有1根1G内存
Size: 1024 MB       //此插槽有1根1G内存
Size: 1024 MB       //此插槽有1根1G内存
Size: 1024 MB       //此插槽有1根1G内存
Size: No Module Installed       //此插槽未使用
Size: No Module Installed       //此插槽未使用
```

```bash
  # dmidecode -t 17        //数字17是dmidecode的参数，本文最后有其他数字参数
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
```
