---
title: OCR
---

# 概述

> 参考：
> - [Wiki，Optical_character_recognition](https://en.wikipedia.org/wiki/Optical_character_recognition)
> - https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.6/doc/doc_ch/models.md

**Optical character recognition(光学字符识别，简称 OCR)** 是将图像以电子或机械方式转换为机器编码文本，无论是来自扫描文档、文档照片、场景照片、叠加在图像上的字母文字等。目前是文字识别的统称，已不限于文档或书本文字识别，更包括识别自然场景下的文字，又可以称为 **Scene Text Recognition(场景文字识别，简称 STR)**。

OCR文字识别一般包括两个部分，文本检测和文本识别；文本检测首先利用检测算法检测到图像中的文本行；然后检测到的文本行用识别算法去识别到具体文字。

## Detection Model(检测模型)

文本检测就是要定位图像中的文字区域，然后通常以边界框的形式将单词或文本行标记出来。传统的文字检测算法多是通过手工提取特征的方式，特点是速度快，简单场景效果好，但是面对自然场景，效果会大打折扣。当前多是采用深度学习方法来做。

基于深度学习的文本检测算法可以大致分为以下几类：

1.  基于目标检测的方法；一般是预测得到文本框后，通过NMS筛选得到最终文本框，多是四点文本框，对弯曲文本场景效果不理想。典型算法为EAST、Text Box等方法。
2.  基于分割的方法；将文本行当成分割目标，然后通过分割结果构建外接文本框，可以处理弯曲文本，对于文本交叉场景问题效果不理想。典型算法为DB、PSENet等方法。
3.  混合目标检测和分割的方法；

## Recognition Model(识别模型)

OCR识别算法的输入数据一般是文本行，背景信息不多，文字占据主要部分，识别算法目前可以分为两类算法：

1.  基于CTC的方法；即识别算法的文字预测模块是基于CTC的，常用的算法组合为CNN+RNN+CTC。目前也有一些算法尝试在网络中加入transformer模块等等。
2.  基于Attention的方法；即识别算法的文字预测模块是基于Attention的，常用算法组合是CNN+RNN+Attention

# PaddleOCR

> 参考：
> - [GitHub 项目，PaddlePaddle/PaddleOCR](https://github.com/PaddlePaddle/PaddleOCR)
>     - <https://www.bilibili.com/video/BV1iY4y1s7fx>

PaddleOCR 是百度开源的 OCR 工具。

有 Python 库。

## 模型说明

PaddleOCR 中集成了很多OCR算法，文本检测算法有DB、EAST、SAST等等，文本识别算法有CRNN、RARE、StarNet、Rosetta、SRN等算法。

其中PaddleOCR针对中英文自然场景通用OCR，推出了PP-OCR系列模型，PP-OCR模型由DB+CRNN算法组成，利用海量中文数据训练加上模型调优方法，在中文场景上具备较高的文本检测识别能力。并且PaddleOCR推出了高精度超轻量PP-OCRv2模型，检测模型仅3M，识别模型仅8.5M，利用[PaddleSlim](https://github.com/PaddlePaddle/PaddleSlim)的模型量化方法，可以在保持精度不降低的情况下，将检测模型压缩到0.8M，识别压缩到3M，更加适用于移动端部署场景。


## 关联文件与配置

**~/.paddleocr/whl/** # 模型保存路径。注意：第一次运行调用 PaddleOCR 包运行代码时，会自动下载最新的模型。详见 [paddleocr.py 文件的 MODEL_URLS 变量](https://github.com/PaddlePaddle/PaddleOCR/blob/release/2.6/paddleocr.py#L58)

- **./cls/** # Direction Classification 方向分类器
- **./det/** # Detection 检测模型
- **./rec/** # Recognition 识别模型

# Umi-OCR

> 参考：
> - [GitHub 项目，hiroi-sora/Umi-OCR](https://github.com/hiroi-sora/Umi-OCR)
> - [公众号-差评，完全免费，不用联网，这套 OCR 工具比微信的还好用！](https://mp.weixin.qq.com/s/lkoBOAYCdIY8F2Y6FCR-7w)

OCR 图片转文字识别软件，完全离线。截屏/批量导入图片，支持多国语言、合并段落、竖排文字。可排除水印区域，提取干净的文本。基于 PaddleOCR 。
