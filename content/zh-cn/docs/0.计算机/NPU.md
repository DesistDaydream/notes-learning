---
title: "NPU"
linkTitle: "NPU"
created: "2026-04-27T11:31"
weight: 100
---

# 概述

> 参考：
>
> - [Wiki, Neural processing unit](https://en.wikipedia.org/wiki/Neural_processing_unit)

**Neural processing unit(神经处理单元，简称 NPU)** 是一类专门的硬件加速器，旨在加速 AI 的机器学习相关应用的效率。

NPU 在 Linux 内核管理的 [PCI](/docs/1.操作系统/Kernel/Hardware/PCI.md) 上被划分为 Processing accelerators 类别，ID 是 1200。

# Ascend

**Ascend(昇腾)** 生态的 NPU

## 学习资料

驱动与固件 社区版资源下载: https://www.hiascend.com/hardware/firmware-drivers/community

## 安装 NPU 驱动

> [!Tip] [CANN](/docs/12.AI/Computing%20platform/CANN.md) 的安装文档中，包含了安装 “NPU 驱动与固件” 的内容。甚至提供了操作系统中的命令，可以识别出当前服务器使用的 NPU 是什么型号。非常便捷

安装 NPU 驱动指安装 **驱动** 和 **固件**。

软件包的命名格式：

- 驱动 # `{product name}-npu-driver_x.x.x_linux-{arch}._run`
- 固件 # `{product name}-npu-firmware_x.x.x._run`

### 最佳实践

使用其中一个场景举例

**一、确认 NPU 型号**

```bash
[root@localhost ~]# cat /sys/class/dmi/id/product_name
KunLun G2280
[root@localhost ~]# cat /sys/class/dmi/id/product_serial
210619FFNXXHR3000001
```

去河南昆仑的维保查询，查询该服务器的 NPU 卡型号，得到 `KunLun G2280-(1*2*S920S00-5220,2*2000W,16*64GB,4*3840GB-SATA,2*960GB-SATA,1*9460-8i,1*4*GE,2*XP382,4*Atlas 300I Duo)`，是 **Atlas 300I Duo**

**二、找到安装文档**

在[官方文档](https://www.hiascend.com/document)找到 “[硬件产品 - 加速卡 - Atlas 300I Duo 推理卡](https://support.huawei.com/enterprise/zh/ascend-computing/atlas-300i-duo-pid-252823107)”，进入后，从 “软件部署指南” 相关文字中，找到相关文档，e.g. [驱动安装](https://support.huawei.com/enterprise/zh/doc/EDOC1100245756/2645a51f?idPath=23710424|251366513|22892968|252309139|252823107) 与 [固件安装](https://support.huawei.com/enterprise/zh/doc/EDOC1100245756/f4a62e83?idPath=23710424|251366513|22892968|252309139|252823107)

**三、下载驱动与固件**

进入 [学习资料](#学习资料) 中的 “固件与驱动资源” 网站，选择 “加速卡” 产品，产品型号为 “Atlas 300I Duo”。CANN 版本 和 固件与驱动 默认。选择服务器的 CPU 架构，我们使用 .run 进行安装。下载两个文件

- 驱动: Ascend-hdk-310p-npu-driver_25.5.2_linux-aarch64.run
- 固件: Ascend-hdk-310p-npu-firmware_7.8.0.7.220.run

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/20260429103422812.png)

通过 [DevTools](/docs/Web/Browser/DevTools.md) 工具获取下载链接，在服务器中使用 wget, curl, etc. 工具下载

**四、安装驱动与固件**

创建用户

```bash
groupadd HwHiAiUser
useradd -g HwHiAiUser -d /home/HwHiAiUser -m HwHiAiUser -s /bin/bash
```

赋予 `XXX.run` 文件可执行权限 `chmod 755 Ascend-hdk-301p*`

设置环境变量

```bash
export chip_type="310p"
export arch=$(uname -m)
export driver_version="25.5.2"
export fimware_version="7.8.0.7.220"
```

**安装驱动**

```bash
./Ascend-hdk-${chip_type}-npu-driver_${driver_version}_linux-${arch}.run --full
```

根据提示决定是否需要重启。之后使用 `npu-smi` 程序检查 NPU 信息

```bash
~]# npu-smi info
+--------------------------------------------------------------------------------------------------------+
| npu-smi 25.5.2                                   Version: 25.5.2                                       |
+-------------------------------+-----------------+------------------------------------------------------+
| NPU     Name                  | Health          | Power(W)     Temp(C)           Hugepages-Usage(page) |
| Chip    Device                | Bus-Id          | AICore(%)    Memory-Usage(MB)                        |
+===============================+=================+======================================================+
| 1       310P3                 | OK              | NA           51                0     / 0             |
| 0       0                     | 0000:01:00.0    | 0            1605 / 44278                            |
+-------------------------------+-----------------+------------------------------------------------------+
......略
+-------------------------------+-----------------+------------------------------------------------------+
| NPU     Chip                  | Process id      | Process name             | Process memory(MB)        |
+===============================+=================+======================================================+
| No running processes found in NPU 1                                                                    |
+===============================+=================+======================================================+
......略
```

**安装固件**

```bash
./Ascend-hdk-${chip_type}-npu-firmware_${fimware_version}.run --full
```

若出现 `Firmware package installed successfully! Reboot now or after driver installation for the installation/upgrade to take effect.` 则**立刻重新启动服务器**

**五、安装完成**

使用 `/usr/local/Ascend/driver/tools/upgrade-tool --device_index -1 --component -1 --version` 命令检查芯片固件版本号，若与目标版本一致，则说明安装成功

所有组件默认安装在了 `/usr/local/Ascend/` 目录中。安装路径在安装驱动时指定，固件跟随驱动安装路径而变。

驱动安装完成后，需要安装 [CANN](/docs/12.AI/Computing%20platform/CANN.md) 才能开始使用模型。

## 关联文件与配置

**/usr/local/Ascend/** # 昇腾生态相关产品安装目录、运行时目录

- **./driver/** # 驱动
- **./firmware/** # 固件

## npu-smi CLI

## 最佳实践

加载环境 `source /usr/local/Ascend/cann-8.5.2/set_env.sh`
