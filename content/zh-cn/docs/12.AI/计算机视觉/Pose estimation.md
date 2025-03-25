---
title: Pose estimation
linkTitle: Pose estimation
weight: 20
---

# 概述

> 参考：
>
> - [Wiki, Pose](https://en.wikipedia.org/wiki/Pose_(computer_vision))
> - [Wiki, 3D pose estimation](https://en.wikipedia.org/wiki/3D_pose_estimation)

通过 [Object detection](/docs/12.AI/计算机视觉/Object%20detection.md) 检测到对象后，利用 **Pose estimation(姿态估计)** 确定对象的位置和方向，通常是在三维空间中。姿态通常以**变换矩阵**的形式在内部存储。术语 **pose(姿态)** 在很大程度上与 **transform(变换)** 同义，但 transform 常常包含 **缩放**，而姿态通常不包含。
