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

# 华为

**Ascend(昇腾)**

昇思 是框架

**Compute Architecture for Neural Networks(神经网络异构计算架构，简称 CANN)**

## 学习资料

固件与驱动 社区版资源下载: https://www.hiascend.com/hardware/firmware-drivers/community

CANN 社区版资源下载: https://www.hiascend.com/developer/download/community

## npu-smi CLI

## 最佳实践

加载环境 `source /usr/local/Ascend/cann-8.5.2/set_env.sh`

