---
title: Driver
linkTitle: Driver
weight: 20
tags:
  - PCI
---

# 概述

> 参考：
>
> -

Linux 中的 **Driver(驱动)** 管理。

部分 BUS_TYPE 可以在 `/sys/bus/${BUS_TYPE}/drivers/` 可以找到各类型总线下的驱动。e.g. [PCI](/docs/1.操作系统/Kernel/Hardware/PCI.md) 设备、USB 设备、etc. 。

> Tip: 并不是所有的 BUS_TYPE 都有驱动，e.g. `/sys/bus/memory/drivers/` 目录是空的

# PCI

`/sys/bus/pci/drivers/${DRIVER}/` 目录是某个具体驱动下关联的设备的 PCI Addr 以及驱动相关操作

- **./${PCI_ADDR}** # 使用了本驱动的设备的 PCI Addr 软链接，指向 `/sys/devices/pciXXX/XXX/...` 某个目录
- **./bind** # 向该文件写入 PCI Addr 将会让设备与内核绑定
- **./unbind** # 向该文件写入 PCI Addr 将会让设备从内核解绑

# 最佳实践

## 网卡的绑定与解绑

我这里用网卡演示绑定和解绑的过程

可以通过 [lspci](/docs/1.操作系统/Linux%20管理/Linux%20硬件管理工具/lspci.md) 命令查看该设备在内核中可以使用的驱动。

```bash
~]# lspci -s 0000:21:00.1 -v
21:00.1 Ethernet controller: Intel Corporation Ethernet Controller X710 for 10GbE SFP+ (rev 02)
        Subsystem: Intel Corporation Ethernet Converged Network Adapter X710
        ......略
        Kernel driver in use: i40e
        Kernel modules: i40e
```

Kernel driver in use 表示当前设备使用的驱动；Kernel modules 表示当前设备使用的内核模块在。我们想要为某个 PCI 设备加载驱动时，可以参考 Kernel modules 的值。

比如现在我们可以在 i40e 驱动下看到 0000:21:00.1 设备

```bash
~]# ls /sys/bus/pci/drivers/i40e/ | grep 0000:21:00.1
```

利用 unbind 可以将该 PCI 设备（i.e. 网卡）从内核解绑

```bash
echo -n "0000:21:00.1" | sudo tee /sys/bus/pci/drivers/i40e/unbind
```

> 解绑后，`lspci -s 0000:21:00.1 -v` 命令不会显示 Kernel driver in use 这行内容

利用 bind 可以将该 PCI 设备（i.e. 网卡）绑定到内核的指定驱动上

```bash
echo -n "0000:21:00.1" | sudo tee /sys/bus/pci/drivers/i40e/unbind
```

> [!Note]
> 在解绑时，若该 PCI 设备使用的不是这个驱动，则会报错
>
> ```bash
> ~]# echo -n "0000:21:00.1" | sudo tee /sys/bus/pci/drivers/igb/unbind
> 0000:21:00.1tee: /sys/bus/pci/drivers/igb/unbind: No such device
> ~]# echo -n "0000:21:00.1" | sudo tee /sys/bus/pci/drivers/i40e/unbind
> 0000:21:00.1
> ```
>
> 可以看到，0000:21:00.1 设备是 i40e 驱动，如果对 igb 解绑，则会提示 No such device
>
> 同样的，若该 PCI 设备不支持某个驱动，在绑定时也会提示 No such device

通过 dmesg 过滤 PCI 地址，也可以看到在启动时，系统是如何为网卡分配 IP，以及如何为网卡命名的

```bash
~]# dmesg | grep 0000:0c:00.0
[    1.352192] pci 0000:0c:00.0: [8086:1533] type 00 class 0x020000
[    1.352219] pci 0000:0c:00.0: reg 0x10: [mem 0xc7200000-0xc72fffff]
[    1.352247] pci 0000:0c:00.0: reg 0x18: [io  0x2000-0x201f]
[    1.352262] pci 0000:0c:00.0: reg 0x1c: [mem 0xc7300000-0xc7303fff]
[    1.352305] pci 0000:0c:00.0: reg 0x30: [mem 0xc7100000-0xc71fffff pref]
[    1.352419] pci 0000:0c:00.0: PME# supported from D0 D3hot D3cold
[    1.927697] pci 0000:0c:00.0: Adding to iommu group 86
# 从这里开始就不是 pci 了，而是 igb 了，说明该 PCI 设备的驱动使用了 igb
[    2.864926] igb 0000:0c:00.0: added PHC on eth1
[    2.864928] igb 0000:0c:00.0: Intel(R) Gigabit Ethernet Network Connection
[    2.864930] igb 0000:0c:00.0: eth1: (PCIe:2.5Gb/s:Width x1) d8:cb:8a:fa:dc:7f
[    2.864975] igb 0000:0c:00.0: eth1: PBA No: 000200-000
[    2.864976] igb 0000:0c:00.0: Using MSI-X interrupts. 4 rx queue(s), 4 tx queue(s)
[    3.381286] igb 0000:0c:00.0 enp12s0: renamed from eth1 # 将网络设备重命名为 enp12s0
[   32.790007] igb 0000:0c:00.0 enp12s0: igb: enp12s0 NIC Link is Up 1000 Mbps Full Duplex, Flow Control: RX
```
