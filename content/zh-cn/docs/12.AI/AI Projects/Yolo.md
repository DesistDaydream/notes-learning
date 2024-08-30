---
title: Yolo
linkTitle: Yolo
date: 2024-08-30T09:17
weight: 20
---

# 概述

> 参考：
>
> - [ultralytics 官网](https://www.ultralytics.com/)

[YOLO](https://arxiv.org/abs/1506.02640)(You Only Look Once）是一种流行的物体检测和图像分割模型，由华盛顿大学的约瑟夫-雷德蒙（Joseph Redmon）和阿里-法哈迪（Ali Farhadi）开发。YOLO 于 2015 年推出，因其高速度和高精确度而迅速受到欢迎。

- 2016 年发布的[YOLOv2](https://arxiv.org/abs/1612.08242) 通过纳入批量归一化、锚框和维度集群改进了原始模型。
- 2018 年推出的[YOLOv3](https://pjreddie.com/media/files/papers/YOLOv3.pdf) 使用更高效的骨干网络、多锚和空间金字塔池进一步增强了模型的性能。
- [YOLOv4](https://arxiv.org/abs/2004.10934)于 2020 年发布，引入了 Mosaic 数据增强、新的无锚检测头和新的损失函数等创新技术。
- [YOLOv5](https://github.com/ultralytics/yolov5)进一步提高了模型的性能，并增加了超参数优化、集成实验跟踪和自动导出为常用导出格式等新功能。
- [YOLOv6](https://github.com/meituan/YOLOv6)于 2022 年由[美团](https://about.meituan.com/)开源，目前已用于该公司的许多自主配送机器人。
- [YOLOv7](https://github.com/WongKinYiu/yolov7)增加了额外的任务，如 COCO 关键点数据集的姿势估计。
- [YOLOv8](https://github.com/ultralytics/ultralytics)是YOLO 的最新版本，由Ultralytics 提供。YOLOv8 YOLOv8 支持全方位的视觉 AI 任务，包括[检测](https://docs.ultralytics.com/tasks/detect/)、[分割](https://docs.ultralytics.com/tasks/segment/)、[姿态估计](https://docs.ultralytics.com/tasks/pose/)、[跟踪](https://docs.ultralytics.com/modes/track/)和[分类](https://docs.ultralytics.com/tasks/classify/)。这种多功能性使用户能够在各种应用和领域中利用YOLOv8 的功能。
- [YOLOv9](https://docs.ultralytics.com/models/yolov9/) 引入了可编程梯度信息 （PGI） 和广义高效层聚合网络 （GELAN） 等创新方法。
- [YOLOv10](https://docs.ultralytics.com/models/yolov10/)是由[清华大学](https://www.tsinghua.edu.cn/en/)的研究人员使用该软件包创建的。 [Ultralytics](https://ultralytics.com/)[Python 软件包](https://pypi.org/project/ultralytics/)创建的。该版本通过引入端到端头（End-to-End head），消除了非最大抑制（NMS）要求，实现了实时[目标检测](https://docs.ultralytics.com/tasks/detect/)的进步。

