---
title: Object detection
linkTitle: Object detection
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Object detection](https://en.wikipedia.org/wiki/Object_detection)
> - [Ultralytics 文档，术语 - 对象监测](https://www.ultralytics.com/glossary/object-detection)
> - Google 学术 https://scholar.google.com/scholar?q=Object+detection&hl=zh-CN&as_sdt=0&as_vis=1&oi=scholart

**Object detection(对象检测)** 是一种与计算机视觉和图像处理相关的计算机技术，用于检测数字图像和视频中特定类别（例如人类、建筑物或汽车）的语义对象的实例。

**Image Classification(图像分类)** # 确定图像中是否存在特定对象。

**Bounding Box(边界框)** # 用于突出显示图像中对象的位置的矩形。它由坐标 (x, y)、宽度和高度定义。

## 对象监测架构

**One-Stage Detectors(单级检测器)** 例如 YOLO、SSD：一步执行对象定位和分类，提供更快的推理速度，但有时会牺牲准确性。

**Two-Stage Detectors(两阶段检测器)** （例如 Faster R-CNN）：涉及区域提议阶段，然后进行对象分类，通常可以实现更高的准确度，但推理时间可能会更慢。
