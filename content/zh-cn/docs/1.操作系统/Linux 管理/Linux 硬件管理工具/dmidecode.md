---
title: "dmidecode"
linkTitle: "dmidecode"
weight: 20
---

# 概述

> 参考：
>
> - [官网](https://www.nongnu.org/dmidecode/)
> - Manual(手册)，dmidecode
> - [Ubuntu 手册, dmidecode]()
> - [Wiki, dmidecode](https://en.wikipedia.org/wiki/Dmidecode)

dmidecode 命令可以让我们在 Linux 系统下获取有关硬件方面的信息。dmidecode 的作用是将数据库中的信息解码，以可读的文本方式显示。由于 DMI 信息可以人为修改，因此里面的信息不一定是系统准确的信息。dmidecode 遵循 [DMTF](/docs/Standard/IT/DMTF.md) 的 SMBIOS/DMI 标准，其输出的信息包括 BIOS、系统、主板、处理器、内存、缓存、etc. 。

dmidecode 附带三个额外的工具：

- biosdecode 打印它能找到的所有 BIOS 相关信息（参见 [示例输出](https://www.nongnu.org/dmidecode/sample/biosdecode.txt)）；
- ownership 检索可以在 Compaq 计算机上设置的“所有权标签”；
- vpddecode 打印可以在几乎所有 IBM 计算机中找到的“重要产品数据”信息（参见 [示例输出](https://www.nongnu.org/dmidecode/sample/vpddecode.txt)）。

## DMI 类型

截止到 XXXX-XX-XX 的 YY 版本的 SMBIOS 规范，定义了以下 DMI 类型:

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

此外：

 • 类型 `126` 用于表示已禁用（disabled）的条目；
 • 类型 `127` 用作表结束标记（end-of-table marker）；
 • 类型 `128` 到 `255` 保留用于 OEM（厂商）自定义数据。

`dmidecode` 默认会显示这些条目，但只有在硬件厂商提供了相关文档或解析代码的情况下，它才能正确解码这些内容。

在 `--type` 参数中，可以使用 Keywork(关键字) 代替类型编号。每个关键字都等价于一组类型编号:

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

关键字的匹配不区分大小写。以下命令行是等效的：

 • dmidecode --type 0 --type 13
 • dmidecode --type 0,13
 • dmidecode --type bios
 • dmidecode --type BIOS
```

# 关联文件与配置

- /dev/mem
- /sys/firmware/dmi/tables/smbios_entry_point (Linux only)
- /sys/firmware/dmi/tables/DMI (Linux only)

# Syntax(语法)

**dmidecode \[OPTIONS]**

**OPTIONS**

- **-t, --type**(STRING) # 只显示指定类型的条目。可用的值有: DMI 类型编号、以逗号分隔的类型编号列表、关键字(详见 [DMI 类型](#DMI%20类型))

# EXAMPLE

## 服务器信息

```bash
~]# dmidecode | grep "System Information" -A9|egrep  "Manufacturer|Product|Serial"
Manufacturer: VMware, Inc.
Product Name: VMware Virtual Platform
Serial Number: VMware-42 18 c8 32 77 c6 ec 16-3f 31 94 e9 d0 34 a6 ac
```

若是正规厂家出厂的服务器，甚至也能看到服务器的型号和序列号等信息

```bash
~]# dmidecode -t system
# dmidecode 3.4
Getting SMBIOS data from sysfs.
SMBIOS 3.0 present.

Handle 0x0001, DMI type 1, 27 bytes
System Information
        Manufacturer: XFUSION
        Product Name: 2288H V5
        Version: Purley
        Serial Number: 这里是真实的服务器序列号
        UUID: 8c2eb7d2-f46b-862f-ee11-062480781a3c
        Wake-up Type: Power Switch
        SKU Number: Purley
        Family: Purley
```

## 主板信息

```bash
dmidecode | grep -A16 "System Information$"
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

查看已经查了内存的插槽

```bash
dmidecode -t memory | grep "^[[:space:]]*Size: [0-9]"
```

Linux 查看内存的频率：

```bash
~]# dmidecode | grep -A16 "Memory Device" | grep 'Speed'
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: 667 MHz (1.5 ns)
        Speed: Unknown
```

在 linux 查看内存型号的命令：

内存槽及内存条：

```bash
dmidecode | grep -A16 "Memory Device$"
```

```bash
dmidecode | grep -P 'Maximum\s+Capacity'    # 最大支持几G内存
dmidecode | grep -P -A5 "Memory\s+Device"|grep Size|grep -v Range       # 总共几个插槽，已使用几个插槽
Size: 1024 MB       # 此插槽有1根1G内存
Size: 1024 MB       # 此插槽有1根1G内存
Size: 1024 MB       # 此插槽有1根1G内存
Size: 1024 MB       # 此插槽有1根1G内存
Size: No Module Installed       # 此插槽未使用
Size: No Module Installed       # 此插槽未使用
```


