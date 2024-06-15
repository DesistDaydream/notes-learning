---
title: "Linux 硬件管理工具"
linkTitle: "Linux 硬件管理工具"
weight: 1
---

# 概述

> 参考：
>
> -

[dmidecode](docs/1.操作系统/Linux%20管理/Linux%20硬件管理工具/dmidecode.md)

# lshw

> 参考：
>
> - [GitHub 项目，lyonel/lshw](https://github.com/lyonel/lshw)
> - [官网](http://lshw.ezix.org/)

lshw 是一个提供机器硬件配置详细信息的小工具。

lshw 输出的内容详解见：<https://ezix.org/project/wiki/HardwareLiSter#Howtointerpretlshwsoutput>

## Syntax(语法)

**lshw [FORMAT] [OPTIONS]**

FORMAT

- **-X** # 启用 GUI（如果可用的话）
- **-html** # 以 [HTML](/docs/2.编程/标记语言/HTML.md) 格式输出
- **-xml** # 以 [XML](/docs/2.编程/标记语言/XML.md) 格式输出
- **-json** # 以 [JSON](/docs/2.编程/无法分类的语言/JSON.md) 格式输出
- **-short** # 打印硬件路径。
```bash
H/W path        Device      Class          Description
======================================================
                            system         Standard PC (i440FX + PIIX, 1996)
/0                          bus            Motherboard
/0/0                        memory         96KiB BIOS
/0/400                      processor      AMD EPYC 7542 32-Core Processor
/0/401                      processor      AMD EPYC 7542 32-Core Processor
/0/1000                     memory         8GiB System Memory
/0/1000/0                   memory         8GiB DIMM RAM
/0/100                      bridge         440FX - 82441FX PMC [Natoma]
/0/100/1                    bridge         82371SB PIIX3 ISA [Natoma/Triton II]
/0/100/1.1                  storage        82371SB PIIX3 IDE [Natoma/Triton II]
/0/100/1.2                  bus            82371SB PIIX3 USB [Natoma/Triton II]
......略
```
- **-businfo** # 打印 [Bus(总线)](docs/0.计算机/Motherboard/Bus.md) 的信息
```bash
Bus info          Device      Class          Description
========================================================
                              system         Standard PC (i440FX + PIIX, 1996)
                              bus            Motherboard
                              memory         96KiB BIOS
cpu@0                         processor      AMD EPYC 7542 32-Core Processor
cpu@1                         processor      AMD EPYC 7542 32-Core Processor
                              memory         8GiB System Memory
                              memory         8GiB DIMM RAM
pci@0000:00:00.0              bridge         440FX - 82441FX PMC [Natoma]
pci@0000:00:01.0              bridge         82371SB PIIX3 ISA [Natoma/Triton II]
pci@0000:00:01.1              storage        82371SB PIIX3 IDE [Natoma/Triton II]
pci@0000:00:01.2              bus            82371SB PIIX3 USB [Natoma/Triton II]
......略
```

若是不适用 -businfo 或 -short 格式选项，输出效果类似 tree、pstree 等命令这种带缩进的样式

```bash
bj-test-desistdaydream-1
    description: Computer
    product: Standard PC (i440FX + PIIX, 1996)
    vendor: QEMU
    version: pc-i440fx-focal
    width: 64 bits
    capabilities: smbios-2.8 dmi-2.8 smp vsyscall32
    configuration: boot=normal uuid=9F023923-B44E-2C40-A62E-403E6A1F69D9
  *-core
       description: Motherboard
       physical id: 0
     *-firmware
          description: BIOS
          vendor: SeaBIOS
          physical id: 0
          version: 1.13.0-1ubuntu1.1
          date: 04/01/2014
          size: 96KiB
     *-cpu:0
          description: CPU
          product: AMD EPYC 7542 32-Core Processor
          ......略
     *-cpu:1
          description: CPU
          product: AMD EPYC 7542 32-Core Processor
          vendor: Advanced Micro Devices [AMD]
          physical id: 401
          bus info: cpu@1
          version: pc-i440fx-focal
          slot: CPU 1
          size: 2GHz
          capacity: 2GHz
          width: 64 bits
          capabilities: fpu fpu_exception wp vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 syscall nx mmxext fxsr_opt pdpe1gb rdtscp x86-64 rep_good nopl cpuid extd_apicid tsc_known_freq pni pclmulqdq ssse3 fma cx16 sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm cmp_legacy svm cr8_legacy abm sse4a misalignsse 3dnowprefetch osvw perfctr_core ssbd ibrs ibpb stibp vmmcall fsgsbase tsc_adjust bmi1 avx2 smep bmi2 rdseed adx smap clflushopt clwb sha_ni xsaveopt xsavec xgetbv1 xsaves clzero xsaveerptr wbnoinvd arat npt nrip_save umip rdpid arch_capabilities
          configuration: cores=1 enabledcores=1 threads=1
     *-memory
          description: System Memory
          physical id: 1000
          size: 8GiB
          capabilities: ecc
          configuration: errordetection=multi-bit-ecc
        *-bank
             description: DIMM RAM
             vendor: QEMU
             physical id: 0
             slot: DIMM 0
             size: 8GiB
     *-pci
          description: Host bridge
          product: 440FX - 82441FX PMC [Natoma]
          vendor: Intel Corporation
          physical id: 100
          bus info: pci@0000:00:00.0
          version: 02
          width: 32 bits
          clock: 33MHz
        *-isa
             description: ISA bridge
             product: 82371SB PIIX3 ISA [Natoma/Triton II]
             vendor: Intel Corporation
             physical id: 1
             bus info: pci@0000:00:01.0
             version: 00
             width: 32 bits
             clock: 33MHz
             capabilities: isa
             configuration: latency=0
        *-ide
             description: IDE interface
             product: 82371SB PIIX3 IDE [Natoma/Triton II]
             vendor: Intel Corporation
             physical id: 1.1
             bus info: pci@0000:00:01.1
             version: 00
             width: 32 bits
             clock: 33MHz
             capabilities: ide isa_compat_mode bus_master
             configuration: driver=ata_piix latency=0
             resources: irq:0 ioport:1f0(size=8) ioport:3f6 ioport:170(size=8) ioport:376 ioport:c200(size=16)
        *-usb
             description: USB controller
             product: 82371SB PIIX3 USB [Natoma/Triton II]
             vendor: Intel Corporation
             physical id: 1.2
             bus info: pci@0000:00:01.2
             version: 01
             width: 32 bits
             clock: 33MHz
             capabilities: uhci bus_master
             configuration: driver=uhci_hcd latency=0
             resources: irq:11 ioport:c1c0(size=32)
           *-usbhost
                product: UHCI Host Controller
                vendor: Linux 5.4.0-100-generic uhci_hcd
                physical id: 1
                bus info: usb@1
                logical name: usb1
                version: 5.04
                capabilities: usb-1.10
                configuration: driver=hub slots=2 speed=12Mbit/s
              *-usb
                   description: Mouse
                   product: QEMU USB Mouse
                   vendor: QEMU
                   physical id: 1
                   bus info: usb@1:1
                   version: 0.00
                   serial: 89126-0000:00:01.2-1
                   capabilities: usb-2.00
                   configuration: driver=usbhid maxpower=100mA speed=12Mbit/s
        *-bridge
             description: Bridge
             product: 82371AB/EB/MB PIIX4 ACPI
             vendor: Intel Corporation
             physical id: 1.3
             bus info: pci@0000:00:01.3
             version: 03
             width: 32 bits
             clock: 33MHz
             capabilities: bridge
             configuration: driver=piix4_smbus latency=0
             resources: irq:9
        *-display
             description: VGA compatible controller
             product: GD 5446
             vendor: Cirrus Logic
             physical id: 2
             bus info: pci@0000:00:02.0
             version: 00
             width: 32 bits
             clock: 33MHz
             capabilities: vga_controller rom
             configuration: driver=cirrus latency=0
             resources: irq:0 memory:fc000000-fdffffff memory:feb90000-feb90fff memory:c0000-dffff
        *-network
             description: Ethernet controller
             product: Virtio network device
             vendor: Red Hat, Inc.
             physical id: 3
             bus info: pci@0000:00:03.0
             version: 00
             width: 64 bits
             clock: 33MHz
             capabilities: msix bus_master cap_list rom
             configuration: driver=virtio-pci latency=0
             resources: irq:10 ioport:c1e0(size=32) memory:feb91000-feb91fff memory:fe000000-fe003fff memory:feb00000-feb7ffff
......略
```

OPTIONS

- `-enable` TEST to enable a test
- `-disable` TEST to disable a test
- **-C, -class CLASS** # 只输出指定 CLASS 的信息。
    - 所有可用的 CLASS 列表见：<https://ezix.org/project/wiki/HardwareLiSter#Deviceclasses>

# lspci

详见: [lspci](docs/1.操作系统/Linux%20管理/Linux%20硬件管理工具/lspci.md)

# smartctl

smartctl -a /dev/sda

# 网卡

## mii-tool



# USB 管理工具

## usbutils

> 参考：
> - [GitHub 项目，gregkh/usbutils](https://github.com/gregkh/usbutils)
> - [官网](http://www.linux-usb.org/)
> - [Manual(手册)，lsusb(8)](https://man7.org/linux/man-pages/man8/lsusb.8.html)

适用于 Linux 的 USB 实用程序，包括 lsusb。这是在 Linux 和 BSD 系统上使用的 USB 工具的集合，用于查询连接到系统的 USB 设备类型。这将在 USB 主机 (即您插入USB设备的机器) 上运行，而不是在 USB 设备 (即您插入USB主机的设备) 上运行。

包括如下几个工具

- lsusb
- usb-devices
- usbhid-dump
- usbreset

### lsusb Syntax(语法)

列出系统上的USB总线和USB设备的详细信息。在输出中，您将看到USB控制器的制造商、型号和当前的状态。

- `sudo lshw -class bus -class usb`

只查看有关USB设备的更详细信息

- `lshw -class usb`