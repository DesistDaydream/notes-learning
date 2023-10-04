---
title: OCR
linkTitle: OCR
date: 2023-10-31T14:39
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，Optical_character_recognition](https://en.wikipedia.org/wiki/Optical_character_recognition)
> - https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.6/doc/doc_ch/models.md

**Optical character recognition(光学字符识别，简称 OCR)** 是将图像以电子或机械方式转换为机器编码文本，无论是来自扫描文档、文档照片、场景照片、叠加在图像上的字母文字等。目前是文字识别的统称，已不限于文档或书本文字识别，更包括识别自然场景下的文字，又可以称为 **Scene Text Recognition(场景文字识别，简称 STR)**。

OCR 文字识别一般包括两个部分，**文本检测**和**文本识别**

- 文本检测首先利用检测算法检测到图像中的文本块
- 然后文本识别利用识别算法去识别文本块中的具体文字

## Detection Model(检测模型)

文本检测就是要定位图像中的文字区域，然后通常以边界框的形式将单词或文本行标记出来。传统的文字检测算法多是通过手工提取特征的方式，特点是速度快，简单场景效果好，但是面对自然场景，效果会大打折扣。当前多是采用深度学习方法来做。

基于深度学习的文本检测算法可以大致分为以下几类：

1.  基于目标检测的方法；一般是预测得到文本框后，通过NMS筛选得到最终文本框，多是四点文本框，对弯曲文本场景效果不理想。典型算法为EAST、Text Box等方法。
2.  基于分割的方法；将文本行当成分割目标，然后通过分割结果构建外接文本框，可以处理弯曲文本，对于文本交叉场景问题效果不理想。典型算法为DB、PSENet等方法。
3.  混合目标检测和分割的方法；

## Recognition Model(识别模型)

OCR 识别算法的输入数据一般是文本行，背景信息不多，文字占据主要部分，识别算法目前可以分为两类算法：

1.  基于 CTC 的方法；即识别算法的文字预测模块是基于 CTC 的，常用的算法组合为 CNN+RNN+CTC。目前也有一些算法尝试在网络中加入 transformer 模块等等。
2.  基于 Attention 的方法；即识别算法的文字预测模块是基于 Attention 的，常用算法组合是 CNN+RNN+Attention

## 预处理

为了可以让程序快速检测到字符块后精准识别字符，有的时候还需要对图片进行预处理

- 比如图片是斜的，我们可以把图片正过来
- 若是图片有干扰，可以去掉干扰
- 等等......

## 总结

用稍微简单一些的话说，检测模型用来检查一个图片中，哪些地方可以被识别模型识别；然后交给识别模型。若将图片直接交给识别模型，那么是无法获得任何结果的

用 PaddleOCR 的识别识别逻辑举例，至少需要用到两种模型：文本检测模型 和 文本识别模型。提供给 PaddleOCR 一张图后，首先先检测图片中包含的文字信息并定位为文本框，然后识别文本框中的文本。

> Tips: 若想要识别倒转的文字，还可以通过 方向分类器 模型进行预处理。有的 OCR 程序还有很多其他的**预处理**操作，比如去斑、二值化、线条去除、布局分析 等等。

如下图：

![image.png|800](https://notes-learning.oss-cn-beijing.aliyuncs.com/ocr/202310311306270.png)

用红框框起来的就是检测到的文本框，每个文本框都由 `[[196.0, 10.0], [237.0, 10.0], [237.0, 28.0], [196.0, 28.0]]` （这里用 驯兽师 三个字的文本块为例）这样的多维数组进行定位

- 外层数组共 4 个元素，分别表示文本框的 4 个顶点；0 号元素为 **左上角**，1 号元素为 **右上角**，2 号元素为 **右下角**，3 号元素为 **左下角**。
- 内层数组共 2 个元素，分别表示顶点的横/纵坐标；0 号元素为像素点的**横轴**坐标，1 号元素为像素点的**纵轴**坐标。

然后对每个文本框进行文字识别，以识别出其中的文字。

# PaddleOCR

> 参考：
> - [GitHub 项目，PaddlePaddle/PaddleOCR](https://github.com/PaddlePaddle/PaddleOCR)
>     - <https://www.bilibili.com/video/BV1iY4y1s7fx>

PaddleOCR 是百度开源的 OCR 工具。旨在打造一套丰富、领先、且实用的OCR工具库，助力开发者训练出更好的模型，并应用落地。

## 模型说明

PaddleOCR 中集成了很多OCR算法，文本检测算法有 DB、EAST、SAST 等等，文本识别算法有CRNN、RARE、StarNet、Rosetta、SRN等算法。

其中PaddleOCR针对中英文自然场景通用OCR，推出了PP-OCR系列模型，PP-OCR模型由DB+CRNN算法组成，利用海量中文数据训练加上模型调优方法，在中文场景上具备较高的文本检测识别能力。并且 PaddleOCR 推出了高精度超轻量 PP-OCRv2 模型，检测模型仅3M，识别模型仅8.5M，利用 [PaddleSlim](https://github.com/PaddlePaddle/PaddleSlim) 的模型量化方法，可以在保持精度不降低的情况下，将检测模型压缩到 0.8M，识别压缩到 3M，更加适用于移动端部署场景。

## 关联文件与配置

**~/.paddleocr/whl/** # 模型保存路径。注意：第一次运行调用 PaddleOCR 包运行代码时，会自动下载最新的模型。详见 [paddleocr.py 文件的 MODEL_URLS 变量](https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.6/paddleocr.py#L58)

- **./cls/** # Direction Classification(方向分类器) 模型保存路径
- **./det/** # Detection(检测) 模型保存路径
- **./rec/** # Recognition(识别) 模型保存路径

## 模型下载

在 [PP-OCR系列模型列表](https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.7/doc/doc_ch/models_list.md) 处可以找到三个基本模型以及一个超轻量模型的简介、配置文件、下载地址。

- 基本模型
  - 文本检测模型
  - 文本识别模型
  - 文本方向分类模型
- 轻量模型

模型都分为多个种类

- 推理模型 # 用于预测引擎推理。通常默认下载这种模型。
- 训练模型 与 预训练模型 # 训练过程中保存的模型的参数、优化器状态和训练中间信息，多用于模型指标评估和恢复训练
  - 训练模型 # 是基于预训练模型在真实数据与竖排合成文本数据上finetune得到的模型，在真实应用场景中有着更好的表现
  - 预训练模型 # 则是直接基于全量真实数据与合成数据训练得到，更适合用于在自己的数据集上finetune。
- nb模型 # 经过飞桨Paddle-Lite工具优化后的模型，适用于移动端/IoT端等端侧部署场景（需使用飞桨Paddle Lite部署）。

选择自己感兴趣的模型，下载即可。下载后，将对应的模型，解压到 `ch_PP-OCRv4_server_rec` 目录下对应模型的目录中。比如在 `2.1 中文识别模型` 章节中，找到 `ch_PP-OCRv4_server_rec` 推理模型，下载并解压到 `ch_PP-OCRv4_server_rec/rec/` 目录下即可。其他两个模型同理。这样就可以更换我们感兴趣的模型。

## Python 库

详见：Python 第三方库 [paddleocr](/docs/2.编程/高级编程语言/Python/Python%20第三方库/图像处理/paddleocr.md) 包 

# 实用 OCR工具
## Umi-OCR

> 参考：
> 
> - [GitHub 项目，hiroi-sora/Umi-OCR](https://github.com/hiroi-sora/Umi-OCR)
> - [公众号-差评，完全免费，不用联网，这套 OCR 工具比微信的还好用！](https://mp.weixin.qq.com/s/lkoBOAYCdIY8F2Y6FCR-7w)

OCR 图片转文字识别软件，完全离线。截屏/批量导入图片，支持多国语言、合并段落、竖排文字。可排除水印区域，提取干净的文本。基于 PaddleOCR 。
