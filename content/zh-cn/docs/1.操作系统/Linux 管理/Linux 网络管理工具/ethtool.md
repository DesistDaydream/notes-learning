---
title: ethtool
linkTitle: ethtool
date: 2024-03-20T08:15
weight: 20
---

# 概述

> 参考：
>
> - [Manual(手册)，ethtool(8)](https://man7.org/linux/man-pages/man8/ethtool.8.html)

ethtool 是一个工具，用来查询或控制网络驱动程序和硬件设备。

> Notes: ethtool 无法查看到 [DPDK](/docs/4.数据通信/DPDK/DPDK.md) 等程序接管的网卡的网络设备。因为已经绕过了内核，ethtool 只能查询或控制直接由内核管理的网卡。

```bash
~]# ethtool em1
Settings for em1:
    Supported ports: [ TP ]
    Supported link modes:   10baseT/Half 10baseT/Full
                            100baseT/Half 100baseT/Full
                            1000baseT/Half 1000baseT/Full
    Supported pause frame use: No
    Supports auto-negotiation: Yes
    Supported FEC modes: Not reported
    Advertised link modes:  10baseT/Half 10baseT/Full
                            100baseT/Half 100baseT/Full
                            1000baseT/Half 1000baseT/Full
    Advertised pause frame use: Symmetric
    Advertised auto-negotiation: Yes
    Advertised FEC modes: Not reported
    Link partner advertised link modes:  10baseT/Full
                                         100baseT/Full
                                         1000baseT/Full
    Link partner advertised pause frame use: No
    Link partner advertised auto-negotiation: Yes
    Link partner advertised FEC modes: Not reported
    Speed: 1000Mb/s
    Duplex: Full
    Port: Twisted Pair
    PHYAD: 1
    Transceiver: internal
    Auto-negotiation: on
    MDI-X: on
    Supports Wake-on: g
    Wake-on: d
    Current message level: 0x000000ff (255)
                   drv probe link timer ifdown ifup rx_err tx_err
    Link detected: yes
```

# Syntax(语法)

**ethtool \[OPTIONS] DeviceName**

## OPTIONS

### 查询选项

- **-i, --driver** # 查询指定网络设备的关联驱动程序的信息。

```bash
]# ethtool -i eno3
driver: igb
version: 5.10.0-136.12.0.86.oe2203sp1.x8
firmware-version: 1.67, 0x80000d72, 17.0.12
expansion-rom-version:
bus-info: 0000:01:00.0
supports-statistics: yes
supports-test: yes
supports-eeprom-access: yes
supports-register-dump: yes
supports-priv-flags: yes
```

- **-l, --show-channels** # 显示指定网络设备的通道数。通道是一个 IRQ，是可以触发该 IRQ 的队列集。

```bash
~]# ethtool -l ens4f0
Channel parameters for ens4f0:
Pre-set maximums:
RX:  0
TX:  0
Other:  1
Combined: 63
Current hardware settings:
RX:  0
TX:  0
Other:  1
Combined: 32
```

- **-m, --dump-module-eeprom, --module-info** # 从 “EEPROM 插件模块”检索并解码硬件信息。比如 SFP+、QSFP。如果驱动程序和模块支持它，光学诊断信息也会被读取和解码。如果指定了 page、bank 或 i2c 参数之一，则仅允许转储单个页面或其部分。在这种情况下，偏移和长度参数相对于 EEPROM 页面边界进行处理。
  - “插件模块”，就是指网卡口插的模块。通常情况下来说，都是光模块。
    - 光模块能显示的信息详见下文“[光模块插件信息](#yVQY0)”
  - 这种信息一般只存在于具有模块网卡的物理机上。虚拟机或者无法插模块的网卡，是没有这种信息的，这时候获取信息，将会出现如下报错：

```bash
~]# ethtool -m enp25s0f0
Cannot get module EEPROM information: Invalid argument
```

- **-S, --statistics** # 获得特定于 NIC 和 驱动程序 的统计信息。

```bash
~]# ethtool -S ens4f0
NIC statistics:
     rx_packets: 100463573
     tx_packets: 59794837
     rx_bytes: 74687073845
     tx_bytes: 7381975533
     rx_pkts_nic: 100463561
     tx_pkts_nic: 59794820
     rx_bytes_nic: 75088927199
     tx_bytes_nic: 7672892990
   ......略
```

### 其他选项

- **-p, --identify** # 可以让指定网络设备关联的网卡上的 led 灯闪烁。**常用来在机房中识别机器网卡**。
- **-s,--change \<DEV>** # 允许更改指定网络设备的部分或全部设置。 以下所有选项仅在指定了-s 时适用。
  - **speed N** # 用 Mb/s 作为单位设置网卡速率。Set speed in Mb/s. ethtool with just the device name as an argument will show you the supported device speeds.
  - **duplex {half|full}** # 设置全双工或半双工模式

# EXAMPLE

显示 veth1.1 对的状态信息，可以在 container 或 namespace 中查看绑定的另一半 veth 在 Host 上的网卡号

- ethtool -S veth1.1

设置 eth0 网卡速率为 100Mb/s，全双工

- ethtool -s eth0 speed 100 duplex full

# Reference
## -i 选项信息

```bash
]# ethtool -i eno3
driver: igb
version: 5.10.0-136.12.0.86.oe2203sp1.x8
firmware-version: 1.67, 0x80000d72, 17.0.12
expansion-rom-version:
bus-info: 0000:01:00.0
supports-statistics: yes
supports-test: yes
supports-eeprom-access: yes
supports-register-dump: yes
supports-priv-flags: yes
```

- **driver** # 网络设备连接网卡所使用的驱动程序. 也可以当作网络设备的类型, 还有 bridge、tun、etc. 这些.
- **version** # 驱动程序的版本号
- **firmware-version** # 与驱动程序一起使用的固件的版本号
- **bus-info** # 网卡在系统总线上的位置信息，比如 PCI 总线上的位置
- supports-XXX # 网络设备的驱动是否支持某些特性
  - **supperts-statistics** # 驱动是否支持收集统计信息
  - **supports-test** # 驱动是否支持测试功能
  - **supports-eeprom-access** # 驱动是否支持 EEPROM 访问
  - **supports-register-dump** # 驱动是否支持寄存器转储功能
  - **supports-priv-flags** # 驱动是否支持私有标志

## -m 选项信息

### 光模块插件信息

- **Identifier** # 标识符。即该口上所插的模块信息。示例值：`0x03 (SFP)`
- **Extended identifier** # 扩展标识符。模块的扩展信息。示例值：`0x04 (GBIC/SFP defined by 2-wire interface ID)`
- **Connector** # 连接器信息。即. 示例值：`0x07 (LC)`
- **Transceiver codes : 0x10 0x00 0x00 0x00 0x00 0x00 0x00 0x00**
- **Transceiver type** # 收发器类型。示例值：`10G Ethernet: 10G Base-SR`
- **Encoding : 0x06 (64B/66B)**
- **BR, Nominal : 10300MBd**
- **Rate identifier : 0x00 (unspecified)**
- **Length (SMF,km) : 0km**
- **Length (SMF) : 0m**
- **Length (50um) : 80m**
- **Length (62.5um) : 20m**
- **Length (Copper) : 0m**
- **Length (OM3) : 300m**
- **Laser wavelength : 850nm**
- **Vendor name : Hisense**
- **Vendor OUI : 00:00:00**
- **Vendor PN : LTF8502-BC+-H3C**
- **Vendor rev : 1**
- **Option values : 0x00 0x1a**
- **Option : RX_LOS implemented**
- **Option : TX_FAULT implemented**
- **Option : TX_DISABLE implemented**
- **BR margin, max : 20%**
- **BR margin, min : 20%**
- **Vendor SN : U8AA7H05H97**
- **Date code : 200727\_\_**
- **Optical diagnostics support : Yes**
- **光模块特有信息**
  - **Laser bias current** # 激光偏置电流。示例值：`5.762 mA`
  - **Laser output power** # 激光输出功率(光模块发送功率)。示例值：`0.5240 mW / -2.81 dBm`
  - **Receiver signal average optical power** # 接收信号的平均光功率(光模块接收功率)。示例值：`0.4646 mW / -3.33 dBm`
- **Module temperature : 26.32 degrees C / 79.38 degrees F**
- **Module voltage : 3.3360 V**
- **模块告警开关**
  - **Alarm/warning flags implemented : Yes**
  - **Laser bias current high alarm : Off**
  - **Laser bias current low alarm : Off**
  - **Laser bias current high warning : Off**
  - **Laser bias current low warning : Off**
  - **Laser output power high alarm : Off**
  - **Laser output power low alarm : Off**
  - **Laser output power high warning : Off**
  - **Laser output power low warning : Off**
  - **Module temperature high alarm : Off**
  - **Module temperature low alarm : Off**
  - **Module temperature high warning : Off**
  - **Module temperature low warning : Off**
  - **Module voltage high alarm : Off**
  - **Module voltage low alarm : Off**
  - **Module voltage high warning : \`Off**
  - **Module voltage low warning : Off**
  - **Laser rx power high alarm : Off**
  - **Laser rx power low alarm : Off**
  - **Laser rx power high warning : Off**
  - **Laser rx power low warning : Off**
- **模块告警阈值**
  - **Laser bias current high alarm threshold : 16.500 mA**
  - **Laser bias current low alarm threshold : 1.000 mA**
  - **Laser bias current high warning threshold : 15.000 mA**
  - **Laser bias current low warning threshold : 1.000 mA**
  - **Laser output power high alarm threshold : 1.2589 mW / 1.00 dBm**
  - **Laser output power low alarm threshold : 0.0933 mW / -10.30 dBm**
  - **Laser output power high warning threshold : 0.7943 mW / -1.00 dBm**
  - **Laser output power low warning threshold : 0.1862 mW / -7.30 dBm**
  - **Module temperature high alarm threshold : 81.00 degrees C / 177.80 degrees F**
  - **Module temperature low alarm threshold : 0.00 degrees C / 32.00 degrees F**
  - **Module temperature high warning threshold : 78.00 degrees C / 172.40 degrees F**
  - **Module temperature low warning threshold : 3.00 degrees C / 37.40 degrees F**
  - **Module voltage high alarm threshold : 3.7950 V**
  - **Module voltage low alarm threshold : 2.8050 V**
  - **Module voltage high warning threshold : 3.5000 V**
  - **Module voltage low warning threshold : 3.1000 V**
  - **Laser rx power high alarm threshold : 1.2589 mW / 1.00 dBm**
  - **Laser rx power low alarm threshold : 0.0646 mW / -11.90 dBm**
  - **Laser rx power high warning threshold : 0.7943 mW / -1.00 dBm**
  - **Laser rx power low warning threshold : 0.1023 mW / -9.90 dBm**
