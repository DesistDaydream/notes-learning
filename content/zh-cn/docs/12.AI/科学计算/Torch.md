---
title: "Torch"
linkTitle: "Torch"
weight: 20
---
# 概述

> 参考：
>
> - [GitHub 项目，torch/torch7](https://github.com/torch/torch7)
> - [官网](http://torch.ch/)
> - [Wiki，Torch(机器学习)](https://en.wikipedia.org/wiki/Torch_(machine_learning))

**Torch** 是一个开源的机器学习库，一个科学计算框架，也是一种基于 Lua 的脚本语言。它为用 C 实现的深度学习算法提供 LuaJIT 接口。它是在 EPFL 的 IDIAP 创建的。 Torch 开发于 2017 年转移到 [PyTorch](/docs/12.AI/机器学习/PyTorch.md)，这是 Python 库的一个端口。

Torch 自称为神经网络界的 Numpy，因为他能将 Torch 产生的 **Tensor(张量)** 放在 GPU 中加速运算 (前提是你有合适的 GPU)，就像 Numpy 会把 array 放在 CPU 中加速运算。所以神经网络的话, 当然是用 Torch 的 tensor 形式数据最好。就像 Tensorflow 当中的 Tensor 一样。

# 安装

略，一般都安装 [PyTorch](/docs/12.AI/机器学习/PyTorch.md)
