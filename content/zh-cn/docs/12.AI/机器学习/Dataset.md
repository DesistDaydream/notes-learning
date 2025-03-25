---
title: Dataset
linkTitle: Dataset
weight: 20
---

# 概述

> 参考：
>
> - TODO: Wiki, 没有 Wiki？
> - https://docs.lanyingim.com/quest/40_20240615_1_73_1718389635.html
>
> [^1]:[Wiki, Labeled data](https://en.wikipedia.org/wiki/Labeled_data)

**Datasete(数据集)** 是用于 训练模型、测试模型 的一组特定数据的集合。Dataset 中通常包含两部分内容

- **Unlabeled data(未标记数据)** i.e. 原始数据
- **Labeled data(已标记数据)** 也称为 Annotated data(已标注的数据)

通过数据标注，将 Unlabeled data 变为 Labeled data 后，将这些 Labeled data 打包成 Dataset 供模型训练。

有时候，训练模型并需要已标记的数据，比如 无监督学习、etc. 。

# Data annotation(数据注释)

> 参考：
>
> - [OpenCV 博客，Data Annotation – A Beginner’s Guide](https://opencv.org/blog/data-annotation/)
> - [Wiki, Labeled data](https://en.wikipedia.org/wiki/Labeled_data)
> - https://www.amantyatech.com/data-annotation-and-labeling-everything-you-need-to-know
> - [AWS，What is Data Labeling?](https://aws.amazon.com/what-is/data-labeling/)
> - https://toloka.ai/blog/annotation-vs-labeling/

**Data annotation(数据注释)** 是一组数据中的 信息 **Annotation(注释)/Label(标签)**。这些 **信息注释/标签** 在不同的场景下有不同的表示含义，e.g. 一张照片中哪部分是牛，哪部分是马；录音中说了哪些词；视频中正在执行什么类型的动作；新闻文章的主题是什么；推文想要表达的是一种什么情绪；X 光片中的一个点是否是肿瘤；etc. 。

> [!Note]
> 上面例子中标记照片中的牛、马通常用于 [Object detection](/docs/12.AI/计算机视觉/Object%20detection.md)(对象检测) 任务、识别音频中的问题通常用于 语音识别 任务、识别动作通常用于 [Pose estimation](/docs/12.AI/计算机视觉/Pose%20estimation.md)(姿态评估) 任务、etc.

人们通过 **Data annotating/Data Labeling(数据标注)** 行为，处理一组未标记的数据，为每一条数据添加数据注释。这些被添加了注释的数据称为 **Labeled data/Annotated data(已标记的数据)**

> [!Attention]
> Data annotation 是非常常见的词，在口语化或各种文章中，通常有多种含义。既可以表示动作，也可以当作形容词描述一个实体，也可以作名词表示一个实体。e.g. 我执行了 Data annotation(数据标注) 行为，添加了 Data annotation(数据注释)，让这些数据变成了 Data annotation(具有注释的数据)。( ╯□╰ )
>
> 个人感觉 Annotation 是 Label 的超集，Label 应该是 Annotation 的一种。

由于训练模型使用数据集进行训练，这些已标记数据的质量直接影响[机器学习](/docs/12.AI/机器学习/机器学习.md)的效果。

Data annotation 通常可以简单得做如下分类：

- Text annotation(文字注释)
- Image annotation(图像标注)
- Video annotation(视频注释)
- Audio Annotation(音频注释)
- Key-point Annotation(关键点标注)
- etc.
- TODO

基于不同的 AI 任务场景及实际需求，Data annotation 的样式也多种多样：

- CV 任务
  - **Bounding box(边界框)** # CV 的对象检测任务所需数据的标注方式
- NLP 任务
  - **Text Annotation(文字注释)** #
  - **Semantic Annotation(语义注释)**
- etc.
- TODO

## 数据标注工具

Roboflow

- https://github.com/roboflow
- https://roboflow.com/
- 在线
- https://www.youtube.com/watch?v=7YRJIAIhMpw 讲了一下基础

labelme

- https://github.com/wkentaro/labelme
- 本地

# CV 数据集

> 参考：
>
> -

[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md) 相关 [Dataset](/docs/12.AI/计算机视觉/Dataset.md) 中最常见、最基础的是用于 [Object detection](/docs/12.AI/计算机视觉/Object%20detection.md)(对象检测) 训练任务的数据集，需要进行人工标注，以确定让模型识别图片中的哪些对象。

2006年，斯坦福大学以人为中心的人工智能研究所联席主任李飞飞发起研究，通过大幅扩大训练数据来改进图像识别的人工智能模型和算法。研究人员从万维网上下载了数百万张图像，一组本科生开始为每张图像应用对象标签。 2007 年，李将亚马逊 Mechanical Turk（数字计件工作在线市场）上的数据标记工作外包。由 49,000 多名工作人员标记的 320 万张图像构成了 ImageNet 的基础，ImageNet 是最大的物体识别轮廓手工标记数据库之一。[^1]

# NLP 数据集
