---
title: "PyTorch"
linkTitle: "PyTorch"
weight: 20
---

# 概述

> 参考：
> 
> - [GitHub 项目，pytorch/pytorch](https://github.com/pytorch/pytorch)
> - [官网](https://pytorch.org/)

PyTorch 是一个使用 [Torch](/docs/12.人工智能/科学计算/Torch.md) 构建的 Python 包，提供两个高级特性：

- 带有强大 GPU 加速的张量计算（类似于 NumPy）
- 基于计算图的自动微分系统构建的深度神经网络

## 安装 PyTorch

> 参考：
> 
> - [官方文档，开始](https://pytorch.org/get-started/locally/)

安装 PyTorch 分为使用 GPU 和 CPU 两种，比如：

- CPU
  - `pip3 install torch torchvision torchaudio`
- GPU
  - `pip3 install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118`

GPU 版的 PyTorch 依赖 CUDA

> 如果我们想要使用 GPU 但是却安装的 CPU 版的 PyTorch，将会报错：`Torch not compiled with CUDA enabled`。说白了就是下载的 PyTorch 不是在 CUDA 环境下编译的，无法处理 CUDA 的请求。

# 分类

#人工智能 #机器学习 #Python