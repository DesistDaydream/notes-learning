---
title: 11.图形处理
---

# 概述

> 参考：
> - <https://www.jiqizhixin.com/articles/2019-03-22-10>

## 什么是计算机视觉

为了简化这个问题的答案， 让我们来试想一个场景。
假设你和你的朋友去度假，然后你上传了很多照片到 Facebook 上。但是现在在每张照片中找到你朋友的脸并标记它们要花费很多时间。实际上，Facebook 已经足够智能，它可以帮你标记人物。
那么，你认为自动的特征标记是如何工作的呢？ 简单来说，它通过计算机视觉来实现。
计算机视觉是一个跨学科领域，它解决如何使计算机从数字图像或视频中获得高层次的理解的问题。
这里的想法是将人类视觉系统可以完成的任务自动化。因此，计算机应该能够识别诸如人脸或者灯柱甚至雕像之类的物体。

### 计算机如何读取图像？

思考以下图片：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fhwfe4/1658568420855-b25fc9bb-0d76-4e6f-a1b7-f77be1aff1da.png)
我们可以认出它是纽约天际线的图片。 但是计算机可以自己发现这一切吗？答案是不！
计算机将任何图片都读取为一组 0 到 255 之间的值。
对于任何一张彩色图片，有三个主通道——红色(R)，绿色(G)和蓝色(B)。它的工作原理非常简单。
对每个原色创建一个矩阵，然后，组合这些矩阵以提供 R, G 和 B 各个颜色的像素值。
每一个矩阵的元素提供与像素的亮度强度有关的数据。
思考下图：
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fhwfe4/1658568420764-dd8cadc1-402a-49a6-bdff-ccb3e61ba2eb.png)
如图所示，图像的大小被计算为 B x A x 3。
注意：对于黑白图片，只有一个单一通道。

# GPU

# OpenCV

> 参考：
> - [GitHub 项目，opencv/opencv](https://github.com/opencv/opencv)
> - [官网](https://opencv.org/)
> - [官方文档](https://docs.opencv.org/)，从左侧 Nightly 中选择想要查看的版本
> - <https://zhuanlan.zhihu.com/p/115321759>

**Open Source Computer Vision Library(开源计算机视觉库，简称 OpenCV)** 是一个包含数百种计算机视觉算法的开源库。
官方提供了 Python 语言的 OpenCV 接口~~~在官方这没找到其他语言的

# Halcon

##
