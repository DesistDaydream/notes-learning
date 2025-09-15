---
title: GPU
linkTitle: GPU
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Graphics processing unit](https://en.wikipedia.org/wiki/Graphics_processing_unit)
> - [Wiki, Graphics card](https://en.wikipedia.org/wiki/Graphics_card)

**Graphics card** (也称为 **video card**, **display card**, **graphics accelerator**, **graphics adapter**, **VGA card/VGA**, **video adapter**, **display adapter**, 或者通俗的说 **GPU**)

**Graphics Processing Unit(图形处理单元，简称 GPU)** 是 Graphics card 的核心计算组件，执行计算功能

**Neural Processing Unit(神经处理单元，简称 NPU)**

# 学习资料

[B 站 - 飞天闪客，你管这破玩意叫GPU？](https://www.bilibili.com/video/BV1fxpEzoEnw)

# 显卡常见配置

## 垂直同步

在不开启垂直同步的情况下，某些游戏可能会无线拉高游戏帧率。开启后，会让游戏的帧率限制到显示器的刷新率，某种程度上，关闭垂直同步会降低显卡的利用率

AI：

**开启垂直同步时：**

- 显卡帧率被限制在显示器刷新率（通常60Hz/144Hz等）
- 如果显卡性能过剩，使用率会**降低**，因为它不需要全力渲染超过刷新率的帧数
- 如果显卡性能刚好或不足，使用率可能保持较高水平

**关闭垂直同步时：**

- 显卡会尽全力渲染，追求最高帧率
- 使用率通常会**提高**，特别是在性能要求高的游戏中
- 可能出现画面撕裂，但帧率更高

**实际影响：**

- 高端显卡玩轻量游戏：开启V-Sync会明显降低使用率
- 中低端显卡玩大型游戏：开启V-Sync对使用率影响较小
- 竞技游戏玩家通常关闭V-Sync以获得更高帧率和更低延迟

现在还有**自适应同步技术**（如G-Sync、FreeSync），可以在防止撕裂的同时避免传统V-Sync的缺点，对显卡使用率的影响更加智能化。
