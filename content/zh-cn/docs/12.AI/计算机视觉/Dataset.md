---
title: Dataset
linkTitle: Dataset
date: 2024-10-17T10:14
weight: 20
---

# 概述

> 参考：
>
> - https://opencv.org/blog/data-annotation/

![800](https://opencv.org/wp-content/uploads/2024/02/Data-annotation-types-1536x864.png)

# Data annotations

> 参考：
>
> - [OpenCV 博客，Data Annotation – A Beginner’s Guide](https://opencv.org/blog/data-annotation/)

[计算机视觉](/docs/12.AI/计算机视觉/计算机视觉.md) 中的数据注释常见如下几种：

- **Bounding box**
- Segmentation
- KeyPoint
- Lines and Splines
- etc. TODO

下图中各种颜色的方框都是注释（数字代码是在模型训练时由数据集配置文件中定义其含义），通过 [Object detection](/docs/12.AI/计算机视觉/Object%20detection.md) 识别出对象并添加对象注释

![https://github.com/ultralytics/docs/releases/download/0/mosaiced-coco-dataset-sample.avif|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/mosaiced-coco-dataset-sample.jpg)

下图的注释除了有检测到的对象外，还有对象的姿态。识别出对象后，再通过 [Pose estimation](/docs/12.AI/计算机视觉/Pose%20estimation.md) 进行姿态估计并添加姿态注释

![https://github.com/ultralytics/docs/releases/download/0/mosaiced-training-batch-6.avif|500](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/mosaiced-coco-dataset-sample_pose-estimation.jpg)

上面只是用两种类型举例，还有很多其他的注释类型，这些都属于图像的 **Annotations(注释)**

## Bounding box

> 参考：
>
> - https://www.ultralytics.com/glossary/bounding-box
> - https://encord.com/glossary/bounding-box-definition/

**Bounding box(边界框)** 也称为 Bounding volume 或 Bounding region，是一种 *几何形状* 的 Label，用于在数字图像中包围或环绕一个或一组对象。Bounding box 的目的是在2D 或 3D 空间中定义对象的位置和大小，以执行各种 CV 任务，e.g. 对象检测、分割、分类。这是 CV 领域的基本概念，特别是在涉及图像和视频分析的应用中。

> 在 2D 图像中，Bounding box 通常用矩形表示，其长边与图像的 x 轴和 y 轴平行。矩形的大小由 x 轴和 y 轴的最小值和最大值决定，这些值由矩形的角坐标指定。矩形的大小和中心点也可用于创建 enclosing box(封闭框)。
>
> 在 3D 图像中，边界框通常用平行六面体（3D 矩形）表示，其各个面与图像的 x、y 和 z 轴平行。平行六面体的尺寸由其角的坐标决定，这些坐标表示 x、y 和 z 轴的最小值和最大值。平行六面体的大小和中心也可用于确定 Bounding box。

Bounding box 对于 对象检测 任务及衍生任务至关重要，尤其是对于让模型可以识别和分类图像中的物体。这些 Bounding box 作为真实标注，提供了模型训练时所需的信息（i.e. 确定物体的位置以及如何区分不同的物体）。在像 Ultralytics YOLO 这样的框架中，Bounding box 不仅用于标注，还用于在推理过程中预测物体的位置（也可说是为对象添加 Annotations(注释)）。

Bounding box 效果如下图所示，各种矩形框配上数字，以表示出图像中的对象及该对象的分类或名称。

![500](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/mosaiced-coco-dataset-sample.jpg)

> Tip: Bounding box 就像 [OpenCV](/docs/12.AI/计算机视觉/OpenCV/OpenCV.md) 的 Region Of Interest(感兴趣的区域，简称 ROI)，然后会生成图片对应的 Label，每个 ROI 都有一个数字表示的类别、以及用来定位 ROI 的坐标。

# 常见公开的数据集

### COCO

> 参考：
>
> - [官网](https://cocodataset.org/)

**Common Objects in Context(COCO)**  是一个大规模的对象检测、分割和字幕数据集。 COCO 有几个特点：

- Object segmentation
- Recognition in context
- Superpixel stuff segmentation
- 330K images (>200K labeled)
- 1.5 million object instances
- 80 object categories
- 91 stuff categories
- 5 captions per image
- 250,000 people with keypoints

COCO 数据集常见任务类型：

- **[Object detection](/docs/12.AI/计算机视觉/Object%20detection.md)(对象监测)** # 简称 Detection
  - https://cocodataset.org/#detection-2020
- **DensePose** # COCO DensePose 任务旨在推动在具有挑战性、不受控制的条件下对人体姿势进行密集估计的最先进技术。
  - https://cocodataset.org/#densepose-2020
- **Keypoint detection(关键点监测)** # 简称 Keypoints。COCO 关键点检测任务需要在具有挑战性、不受控制的条件下定位人物关键点。关键点任务涉及同时检测人员并定位其关键点（测试时未给出人员位置）。有关此任务的完整详细信息，请参阅[关键点评估](https://cocodataset.org/#keypoints-eval)页面。
  - https://cocodataset.org/#keypoints-2020
-  **Stuff Segmentation(东西分割)** # 简称 Stuff。COCO Stuff Segmentation Task 旨在推动 stuff 类语义分割的最新技术。对象检测任务针对的是事物类别（人、汽车、大象），而该任务则侧重于事物类别（草、墙壁、天空）。有关内容分割任务的完整详细信息，请参阅内容评估页面。注意：新引入的全景分割任务同时解决了事物和事物类别的识别。
  - https://cocodataset.org/#stuff-2019
- **Panoptic Segmentation(全景分割)** # 简称 Panoptic。全景分割解决了事物和事物类别，统一了通常不同的语义和实例分割任务。其目的是生成丰富且完整的连贯场景分割，这是迈向自动驾驶或增强现实等现实世界视觉系统的重要一步。有关全景分割任务的完整详细信息，请参阅[全景评估](https://cocodataset.org/#panoptic-eval)页面。
  - https://cocodataset.org/#panoptic-2020
- **Image Captioning** # 简称 Captions。用文字描述图像
  - https://cocodataset.org/#captions-2015

### ImageNet

https://www.image-net.org/
