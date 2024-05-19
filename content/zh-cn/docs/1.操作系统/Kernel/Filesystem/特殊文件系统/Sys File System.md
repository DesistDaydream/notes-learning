---
title: Sys File System
---

# Sys File System

> 参考：
>
> - [Manual(手册)，sysfs(5)](https://man7.org/linux/man-pages/man5/sysfs.5.html)
> - [Kernel 文档，用户与管理员指南-Linux ABI 描述](https://www.kernel.org/doc/html/latest/admin-guide/abi.html)
>   - [Linux Kernel 官方文档,Block-Queue sysfs files](https://www.kernel.org/doc/html/latest/block/queue-sysfs.html)（这个 404 了。。。o(╯□╰)o。。。官方目录结构改了。。）
> - [Wiki，Sysfs](https://en.wikipedia.org/wiki/Sysfs)

用于导出内核对象的文件系统

**Sys File System(sys 文件系统，简称 sysfs)** 是一个 pseudo-filesystem(伪文件系统)，提供内核数据结构的接口(更准确地说，sysfs 中的文件和目录提供了内核内部定义的 kobject 结构的视图)。sysfs 下的文件提供关于设备、内核模块、文件系统和其他内核组件的信息。`sysfs 一般挂载到 /sys 目录`。通常情况下，系统会自动挂载它，但也可以使用 `mount -t sysfs sysfs /sys` 命令手动挂载

sysfs 文件系统中的许多文件都是只读的，但是某些文件是可写的，从而允许更改内核变量。 为了避免冗余，符号链接被大量用于连接整个文件系统树中的条目。

**其中 [/sys/devices/](#/sys/devices/) 目录是非常重要且关键的目录，/sys/ 下的其他目录中的内容，有很多都是指向 /sys/devices/ 目录中的软链接**。

## sysfs 背景

Sysfs 文件系统是一个类似于 proc 文件系统的特殊文件系统，用于将系统中的设备组织成层次结构，并向用户模式程序提供详细的内核数据结构信息。

在 2.5 开发周期中，引入了 Linux 驱动程序模型来修复版本 2.4 的以下缺陷：

- 不存在表示驱动程序与设备关系的统一方法。
- 没有通用的热插拔机制。
- procfs 充斥着非过程信息。

Sysfs 的设计目的是导出设备树中存在的信息，从而不再使过程变得混乱。它是由 Patrick Mochel 撰写的。Maneesh Soni 后来编写了 sysfs 后备存储修补程序，以减少大型系统上的内存使用量。

在 2.5 开发的第二年，驱动程序模型和 driverfs（以前称为 ddfs）的基础结构功能开始被证明对其他子系统有用。开发了 kobjects 以提供中央对象管理机制，并且将 driverfs 重命名为 sysfs 以表示其子系统不可知论。

参考：<https://unix.stackexchange.com/questions/4884/what-is-the-difference-between-procfs-and-sysfs>

从一开始（在 Unix 时代开始），程序就了解系统上正在运行的进程的方法是直接从内核内存中读取进程结构（打开 `/dev/mem`，并直接解释原始数据）。这就是最初的 ps 命令的工作方式。随着时间的流逝，一些信息可以通过系统调用获得。

但是，通过 /dev/mem 将系统数据直接公开给用户空间是一种不好的形式，并且每次您要导出一些新的过程数据时都不断地创建新的系统调用是令人讨厌的，因此创建了一种新的方法访问用户空间应用程序的结构化数据以查找有关流程属性的信息。这是 /proc 文件系统。使用 /proc，即使内核中的基础数据结构发生了变化，接口和结构（目录和文件）也可以保持不变。与以前的系统相比，它不那么脆弱，并且扩展性更好。

/proc 文件系统最初旨在发布过程信息和一些关键系统属性，这些属性是“ ps”，“ top”，“ free”和其他一些系统实用程序所必需的。但是，由于易于使用（从内核和用户空间两个方面来看），它成为了整个系统信息的垃圾场。而且，它开始获取读/写文件，用于调整设置并控制内核或其各个子系统的操作。但是，实现控制接口的方法是临时的，并且 /proc 很快陷入混乱。

sysfs（或 /sys 文件系统）旨在为这种混乱增加结构，并提供一种统一的方式来从内核向用户空间公开系统信息和控制点（可设置的系统和驱动程序属性）。现在，注册驱动程序时，内核中的驱动程序框架会根据驱动程序类型及其数据结构中的值自动在 /sys 下创建目录。这意味着特定类型的驱动程序都将具有通过 sysfs 公开的相同元素。

/proc 中仍然可以访问许多旧版系统信息和控制点，但是所有新的总线和驱动程序都应通过 sysfs 公开其信息和控制点。

# /sys/block/

> 参考：
> 
> - [Kernel 文档，Linux 内核用户与管理员指南-Linux ABI 描述-ABI 稳定 链接-/sys/block 链接](https://www.kernel.org/doc/html/latest/admin-guide/abi-stable.html#symbols-under-sys-block)
> - [Kernel 文档，内核子系统文档-Block](https://www.kernel.org/doc/html/latest/block/index.html)

该目录下的所有子目录代表着系统中当前被发现的所有块设备。按照功能来说放置在 /sys/class/ 下会更合适，但由于历史遗留因素而一直存在于 /sys/block，但从 linux2.6.22 内核开始这部分就已经标记为过去时，只有打开了 CONFIG_SYSFS_DEPRECATED 配置编译才会有 这个目录存在，并且其中的内容在从 linux2.6.26 版本开始已经正式移到了 /sys/class/block/，旧的接口 /sys/block/ 为了向后兼容而保留存在，但其中的内容已经变为了指向它们在 `/sys/devices/` 中真实设备的**符号链接**文件。

```bash
~]# ll /sys/block/
total 0
drwxr-xr-x  2 root root 0 Apr  1 14:36 ./
dr-xr-xr-x 13 root root 0 Apr  1 14:36 ../
lrwxrwxrwx  1 root root 0 Apr  1 14:36 dm-0 -> ../devices/virtual/block/dm-0/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 dm-1 -> ../devices/virtual/block/dm-1/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop0 -> ../devices/virtual/block/loop0/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop1 -> ../devices/virtual/block/loop1/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop2 -> ../devices/virtual/block/loop2/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop3 -> ../devices/virtual/block/loop3/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop4 -> ../devices/virtual/block/loop4/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop5 -> ../devices/virtual/block/loop5/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop6 -> ../devices/virtual/block/loop6/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 loop7 -> ../devices/virtual/block/loop7/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 sr0 -> ../devices/pci0000:00/0000:00:01.1/ata1/host0/target0:0:0/0:0:0:0/block/sr0/
lrwxrwxrwx  1 root root 0 Apr  1 14:36 vda -> ../devices/pci0000:00/0000:00:07.0/virtio2/block/vda/
```

**/sys/block/\<BLOCK>/queue/**
  - **./rotational** # 块设备旋转的类型，旋转就是 HHD，不旋转就是 SSD，非常形象生动得比喻磁盘使用的情况~哈哈。`0 表示 SSD`，`1 表示 HDD`
    - 注意：如果磁盘已经被做了 Raid，那么这个值将会一直都是 1。这个说法忘记了出处，找到后补充。

# /sys/bus/

该目录下的每个子目录都是 kernel 支持并且已经注册了的总线类型。这是内核设备按照总线类型分层放置的目录结构，/sys/devices/ 中的所有设备都是连接于某种总线之下的，bus 子目录下的每种具体总线之下可以找到每个具体设备的符号链接，一般来说每个子目录(总线类型)下包含两个子目录，一个是 devices，另一个是 drivers；其中 devices 下是这个总线类型下的所有设备，这些设备**都是符号链接，它们分别指向真正的设备(/sys/devices/)**；而 drivers 下是所有注册在这个总线上的驱动，每个 driver 子目录下 是一些可以观察和修改的 driver 参数。

```bash
~]# ls /sys/bus/
i2c mdio_bus platform sdio usb
iio mmc scsi spi
```

应用 1：msp700 中计算电池电压

PipeADC5 = popen("cat /sys/bus/iio/devices/iio\\:device0/in_voltage5_raw", "r");

应用 2：改变提醒等级

echo 6 > /proc/sys/kernel/printk；

应用 3：msp700 中设置背光

echo 20 > /sys/class/backlight/pwm-backlight/brightness;

等价于：

echo 20 > /sys/bus/platform/devices/pwm-backlight/backlight/pwm-backlight/brightness;

# /sys/class/

该目录下包含所有注册在 [Linux Kernel](docs/1.操作系统/Kernel/Linux%20Kernel/Linux%20Kernel.md) 里面的设备，子目录按照设备功能分类，每个设备类型表达具有某种功能的设备。每个设备类型子目录下都是这种设备类型的各种具体设备的 **[Symbolic link](docs/1.操作系统/Kernel/Filesystem/文件管理/Symbolic%20link.md)(符号链接)，这些链接指向 /sys/devices/ 下的具体设备**。 设备类型和设备并没有一一对应的关系，一个物理设备可能具备多种设备类型；一个设备类型只表达具有一种功能的设备，e.g. 系统所有输入设备都会出现在 /sys/class/input/ 目录中，而不论它们是以何种总线连接到系统的。

```bash
~]# ls /sys/class/
ata_device     dma             i2c-dev   pci_epc       rfkill        tpmrm
ata_link       dmi             input     phy           rtc           tty
ata_port       drm             iommu     powercap      scsi_device   vc
backlight      drm_dp_aux_dev  leds      power_supply  scsi_disk     vfio
bdi            extcon          mdio_bus  ppp           scsi_generic  virtio-ports
block          firmware        mem       pps           scsi_host     vtconsole
bsg            gpio            misc      ptp           sound         wakeup
dax            graphics        mmc_host  pwm           spi_master    watchdog
devcoredump    hidraw          nd        rapidio_port  spi_slave
devfreq        hwmon           net       regulator     thermal
devfreq-event  i2c-adapter     pci_bus   remoteproc    tpm

# 从 /sys/class/net 可以看到，所有网络设备信息，其中包括网络设备的各种状态，比如传输的总字节数等等。
~]# ls -l /sys/class/net
总用量 0
drwxr-xr-x  2 root root 0 9月   1 10:35 ./
drwxr-xr-x 73 root root 0 9月   1 10:35 ../
lrwxrwxrwx  1 root root 0 9月   1 10:35 ens3 -> ../../devices/pci0000:00/0000:00:03.0/virtio0/net/ens3/
lrwxrwxrwx  1 root root 0 9月   1 10:35 lo -> ../../devices/virtual/net/lo/
lrwxrwxrwx  1 root root 0 9月   1 10:50 wg0 -> ../../devices/virtual/net/wg0/
```

**/sys/class/block/DEVICE/** # 块设备信息，DEVICE 是块设备的名称，用来顶替 [/sys/block/](#/sys/block) 目录，软链接到 **/sys/device/** 中的某个目录。

**/sys/class/net/DEVICE** # 网络设备信息，DEVICE 是网络设备的名称。目录中的信息详见 [Linux 网络设备详解](docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Linux%20网络设备详解/Linux%20网络设备详解.md)

# /sys/dev/

该目录下存放主次设备号文件，其中分成字符设备、块设备的主次设备号码(major:minor)组成的文件名，该文件是链接文件并且链接到其真实的设备(/sys/devices)。

# /sys/devices/

该目录下是**全局设备结构体系**，包含所有被发现的注册在各种总线上的各种物理设备。**/sys/ 目录中，只要有关于设备信息的目录，都会是指向 /sys/devices/ 目录中的软链接**。由于 /sys/devices/ 目录结构对设备的分类是按照总线拓扑结构分的，那么对于设备类型来说，就缺乏分类了，所以至今还保留了 /sys/block/ 之类的目录，将设备以类型进行区分。这些区分设备类型的目录下存放的，实际上是指向 /sys/devices/ 目录的**软连接**。

就用下面的 /sys/block/ 目录举例的话，其中的软连接分别会指向 /sys/devices/ 下的不用目录，这就更加说明了 /sys/devices/ 目录对设备的划分并不是按照类型来的。

```bash
]# ll /sys/block/
total 0
drwxr-xr-x  2 root root 0 Mar  6 19:01 ./
dr-xr-xr-x 13 root root 0 Aug 31  2022 ../
lrwxrwxrwx  1 root root 0 Mar  6 19:01 dm-0 -> ../devices/virtual/block/dm-0/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop0 -> ../devices/virtual/block/loop0/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop1 -> ../devices/virtual/block/loop1/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop2 -> ../devices/virtual/block/loop2/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop3 -> ../devices/virtual/block/loop3/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop4 -> ../devices/virtual/block/loop4/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop5 -> ../devices/virtual/block/loop5/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop6 -> ../devices/virtual/block/loop6/
lrwxrwxrwx  1 root root 0 Mar  6 23:10 loop7 -> ../devices/virtual/block/loop7/
lrwxrwxrwx  1 root root 0 Mar  6 19:01 vda -> ../devices/pci0000:00/0000:00:05.0/0000:02:01.0/virtio3/block/vda/
lrwxrwxrwx  1 root root 0 Mar  6 19:01 vdb -> ../devices/pci0000:00/0000:00:05.0/0000:02:02.0/virtio4/block/vdb/
```

一般来说，所有的物理设备都按其在总线上的拓扑结构来显示，但有两个例外即 platform devices 和 system devices。

- platform devices 一般是挂在芯片内部的高速或者低速总线上的各种控制器和外设，它们能被 CPU 直接寻址；
- system devices 不是外设，而是芯片内部的核心结构，比如 CPU，timer 等，它们一般没有相关的驱动，但是会有一些体系结构相关的代码来配置它们。

/sys/devices/ 是内核对系统中所有设备的分层次表达模型，也是 **sysfs 管理设备的最重要的目录结构**。

```bash
~]# ls /sys/devices/
breakpoint  kprobe       msr         platform  software  tracepoint  virtual
isa         LNXSYSTM:00  pci0000:00  pnp0      system    uprobe

~]# ls /sys/devices/platform/
 eisa.0              intel_rapl_msr.0   power        uevent
'Fixed MDIO bus.0'   kgdboc             reg-dummy    vesa-framebuffer.0
 i8042               pcspkr             serial8250

~]# ls /sys/devices/system/
clockevents  clocksource  container  cpu  edac  machinecheck  memory  node
```

# /sys/firmware/

这里是系统加载固件机制的对用户空间的接口，关于固件有专用于固件加载的一套 API，在附录 LDD3 一书中有关于内核支持固件加载机制的更详细的介绍；

# /sys/fs/

此目录包含某些文件系统的子目录。仅当文件系统选择显式创建子目录时，才会在此处具有子目录。

# /sys/kernel/

root@RPM /sys$ ls kernel/

debug mm profiling uevent_seqnum

fscaps notes uevent_helper

这个目录下存放的是内核中所有可调整的参数。

# /sys/module/

该目录下有系统中所有的模块信息，不论这些模块是以内联(inlined)方式编译到内核映像文件中还是编译为外模块(.ko 文件)，都可能出现在/sys/module 中。即 module 目录下包含了所有的被载入 kernel 的模块。

root@RPM /sys$ ls module/

auth_rpcgss lockd rfkill tea5761

block loop scsi_mod tea5767

can mailbox snd tuner_simple

cpuidle mmc_core snd_pcm tuner_xc2028

cryptodev mmcblk snd_pcm_oss ubi

cryptosoft mousedev snd_timer ubifs

d_can mt20xx soundcore usb_storage

d_can_platform mtdoops spurious usbcore

dns_resolver nf_conntrack sunrpc vt

g_mass_storage nf_conntrack_ipv4 tcp_cubic xc4000

hostap nfs tda18271 xc5000

input_polldev ocf tda827x xz_dec

kernel oprofile tda8290

keyboard printk tda9887

# /sys/power/

该目录是系统中的电源选项，对正在使用的 power 子系统的描述。这个目录下有几个属性文件可以用于控制整个机器的电源状态，如可以向其中写入控制命令让机器关机/重启等等。

```bash
root@RPM /sys$ ls power/
pm_async      pm_test       state         wakeup_count
```

