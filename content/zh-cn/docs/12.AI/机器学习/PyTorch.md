---
title: PyTorch
linkTitle: PyTorch
weight: 20
tags:
  - Python
---

# 概述

> 参考：
>
> - [GitHub 项目，pytorch/pytorch](https://github.com/pytorch/pytorch)
> - [官网](https://pytorch.org/)

PyTorch 是一个使用 [Torch](/docs/12.AI/科学计算/Torch.md) 构建的 Python 包，提供两个高级特性：

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
  - `pip3 install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121`

GPU 版的 PyTorch 依赖 CUDA

> Note: 如果我们想要使用 GPU 但是却安装的 CPU 版的 PyTorch，将会报错：`Torch not compiled with CUDA enabled`。说白了就是下载的 PyTorch 不是在 CUDA 环境下编译的，无法处理 CUDA 的请求。

> [!Tip]
> 若安装速度太慢，可以在 pip install 命令中看到 Downloading 的 URL，手动下载，比如 `https://download.pytorch.org/whl/cu121/torch-2.4.1+cu121-cp311-cp311-win_amd64.whl`，然后先执行 `pip intall torch-2.4.1+cu121-cp311-cp311-win_amd64.whl` 进行本地安装，再执行上面的命令安装其他包

安装完成后可以通过如下代码在 Python 解释器中验证 CUDA 是否可用，若可用，将输出 True

```python
import torch
print(torch.cuda.is_available())
```

# 学习

[B 站，10分钟入门神经网络 PyTorch 手写数字识别](https://www.bilibili.com/video/BV1GC4y15736)
