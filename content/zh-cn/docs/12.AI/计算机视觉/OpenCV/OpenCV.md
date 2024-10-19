---
title: OpenCV
linkTitle: OpenCV
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，opencv/opencv](https://github.com/opencv/opencv)
> - [官网](https://opencv.org/)
> - [官方文档](https://docs.opencv.org/)，从左侧 Nightly 中选择想要查看的版本
> - <https://zhuanlan.zhihu.com/p/115321759>
> - [手把手教你使用OpenCV库（附实例、Python代码解析）](https://www.jiqizhixin.com/articles/2019-03-22-10)

**Open Source Computer Vision Library(开源计算机视觉库，简称 OpenCV)** 是一个包含数百种计算机视觉算法的开源库。

## 各语言的库

官方提供了 Python、C++ 的 OpenCV 库

go https://github.com/hybridgroup/gocv

# Modules(模块)

> 参考：
>
> - [4.x 官网文档，主页](https://docs.opencv.org/4.x/index.html)
> - [4.x 官方文档，介绍]()
> - [4.x 官方文档，模块](https://docs.opencv.org/4.x/modules.html)
>   - [4.x 官方文档，Class 列表](https://docs.opencv.org/4.x/annotated.html)
> - https://zhuanlan.zhihu.com/p/19988205

OpenCV 具有模块化的结构，整个 OpenCV 的功能由一个个模块提供，每个模块具有自己的类、函数、方法，并且可以多个模块共享使用。这种模块化的结构可以让 OpenCV 像一门编程语言一样，具有自己的标准库和第三方库，标准库中的标准模块可以实现自身的核心功能，第三方库的模块可以基于核心功能扩展其他功能。就像 https://pkg.go.dev/ 中的各种包，可以看到类型、方法、函数等等的描述。

所有 OpenCV 的类和函数都放在 cv Namespace 中(Namespace 是 C++ 编程语言的基本概念)，如果我们要使用 C++ 代码调用 OpenCV 的模块，需要使用 `cv::` 或在头部添加 `using namespace cv;` 指令。

模块分为两类：

- Main Modules(主模块)
- Extra Modules(额外模块)

## Main Modules(主模块)

- **core** # 核心功能模块，全称 [Core functionality](https://docs.opencv.org/4.x/d0/de1/group__core.html) 。定义了基本的数据结构，包括最重要的 Mat 类、XML 读写、opengl三维渲染等。
- **imgproc** # [图像处理模块](/docs/12.AI/计算机视觉/OpenCV/图像处理模块.md)，全称 Image processing。包括图像滤波、集合图像变换、直方图计算、形状描述、物体检测、等等。图像处理是计算机视觉的重要工具。
- **imgcodecs** # 图像文件读写模块，全称 [Image file reading and writing](https://docs.opencv.org/4.x/d4/da8/group__imgcodecs.html)。
- **videoio** # 视频文件读写模块，全称 [Video I/O](https://docs.opencv.org/4.x/dd/de7/group__videoio.html)。也包括摄像头、Kinect 等的输入。
- **highgui** # 高级图形界面及与 QT 框架的整合。
  - [High-level GUI](https://docs.opencv.org/4.x/d7/dfc/group__highgui.html)
- **video** # 视频分析模块。包括背景提取、光流跟踪、卡尔曼滤波等，做视频监控的读者会经常使用这个模块。
  - [Video Analysis](https://docs.opencv.org/4.x/d7/de9/group__video.html)
- calib3d. [Camera Calibration and 3D Reconstruction](https://docs.opencv.org/4.x/d9/d0c/group__calib3d.html)
- features2d. [2D Features Framework](https://docs.opencv.org/4.x/da/d9b/group__features2d.html)
- objdetect. [Object Detection](https://docs.opencv.org/4.x/d5/d54/group__objdetect.html)
- dnn. [Deep Neural Network module](https://docs.opencv.org/4.x/d6/d0f/group__dnn.html)
- ml. [Machine Learning](https://docs.opencv.org/4.x/dd/ded/group__ml.html)
- flann. [Clustering and Search in Multi-Dimensional Spaces](https://docs.opencv.org/4.x/dc/de5/group__flann.html)
- photo. [Computational Photography](https://docs.opencv.org/4.x/d1/d0d/group__photo.html)
- stitching. [Images stitching](https://docs.opencv.org/4.x/d1/d46/group__stitching.html)
- gapi. [Graph API](https://docs.opencv.org/4.x/d0/d1e/gapi.html)

## Extra Modules(额外模块)

- alphamat. [Alpha Matting](https://docs.opencv.org/4.x/d4/d40/group__alphamat.html)
- aruco. [Aruco markers, module functionality was moved to objdetect module](https://docs.opencv.org/4.x/d9/d6a/group__aruco.html)
- barcode. [Barcode detecting and decoding methods](https://docs.opencv.org/4.x/d2/dea/group__barcode.html)
- bgsegm. [Improved Background-Foreground Segmentation Methods](https://docs.opencv.org/4.x/d2/d55/group__bgsegm.html)
- bioinspired. [Biologically inspired vision models and derivated tools](https://docs.opencv.org/4.x/dd/deb/group__bioinspired.html)
- ccalib. [Custom Calibration Pattern for 3D reconstruction](https://docs.opencv.org/4.x/d3/ddc/group__ccalib.html)
- cudaarithm. [Operations on Matrices](https://docs.opencv.org/4.x/d5/d8e/group__cudaarithm.html)
- cudabgsegm. [Background Segmentation](https://docs.opencv.org/4.x/d6/d17/group__cudabgsegm.html)
- cudacodec. [Video Encoding/Decoding](https://docs.opencv.org/4.x/d0/d61/group__cudacodec.html)
- cudafeatures2d. [Feature Detection and Description](https://docs.opencv.org/4.x/d6/d1d/group__cudafeatures2d.html)
- cudafilters. [Image Filtering](https://docs.opencv.org/4.x/dc/d66/group__cudafilters.html)
- cudaimgproc. [Image Processing](https://docs.opencv.org/4.x/d0/d05/group__cudaimgproc.html)
- cudalegacy. [Legacy support](https://docs.opencv.org/4.x/d5/dc3/group__cudalegacy.html)
- cudaobjdetect. [Object Detection](https://docs.opencv.org/4.x/d9/d3f/group__cudaobjdetect.html)
- cudaoptflow. [Optical Flow](https://docs.opencv.org/4.x/d7/d3f/group__cudaoptflow.html)
- cudastereo. [Stereo Correspondence](https://docs.opencv.org/4.x/dd/d47/group__cudastereo.html)
- cudawarping. [Image Warping](https://docs.opencv.org/4.x/db/d29/group__cudawarping.html)
- cudev. [Device layer](https://docs.opencv.org/4.x/df/dfc/group__cudev.html)
- cvv. [GUI for Interactive Visual Debugging of Computer Vision Programs](https://docs.opencv.org/4.x/df/dff/group__cvv.html)
- datasets. [Framework for working with different datasets](https://docs.opencv.org/4.x/d8/d00/group__datasets.html)
- dnn_objdetect. [DNN used for object detection](https://docs.opencv.org/4.x/d5/df6/group__dnn__objdetect.html)
- dnn_superres. [DNN used for super resolution](https://docs.opencv.org/4.x/d9/de0/group__dnn__superres.html)
- dpm. [Deformable Part-based Models](https://docs.opencv.org/4.x/d9/d12/group__dpm.html)
- face. [Face Analysis](https://docs.opencv.org/4.x/db/d7c/group__face.html)
- freetype. [Drawing UTF-8 strings with freetype/harfbuzz](https://docs.opencv.org/4.x/d4/dfc/group__freetype.html)
- fuzzy. [Image processing based on fuzzy mathematics](https://docs.opencv.org/4.x/df/d5b/group__fuzzy.html)
- hdf. [Hierarchical Data Format I/O routines](https://docs.opencv.org/4.x/db/d77/group__hdf.html)
- hfs. [Hierarchical Feature Selection for Efficient Image Segmentation](https://docs.opencv.org/4.x/dc/d29/group__hfs.html)
- img_hash. [The module brings implementations of different image hashing algorithms.](https://docs.opencv.org/4.x/d4/d93/group__img__hash.html)
- intensity_transform. [The module brings implementations of intensity transformation algorithms to adjust image contrast.](https://docs.opencv.org/4.x/dc/dfe/group__intensity__transform.html)
- julia. [Julia bindings for OpenCV](https://docs.opencv.org/4.x/d7/d44/group__julia.html)
- line_descriptor. [Binary descriptors for lines extracted from an image](https://docs.opencv.org/4.x/dc/ddd/group__line__descriptor.html)
- mcc. [Macbeth Chart module](https://docs.opencv.org/4.x/dd/d19/group__mcc.html)
- optflow. [Optical Flow Algorithms](https://docs.opencv.org/4.x/d2/d84/group__optflow.html)
- ovis. [OGRE 3D Visualiser](https://docs.opencv.org/4.x/d2/d17/group__ovis.html)
- phase_unwrapping. [Phase Unwrapping API](https://docs.opencv.org/4.x/df/d3a/group__phase__unwrapping.html)
- plot. [Plot function for Mat data](https://docs.opencv.org/4.x/db/dfe/group__plot.html)
- quality. [Image Quality Analysis (IQA) API](https://docs.opencv.org/4.x/dc/d20/group__quality.html)
- rapid. [silhouette based 3D object tracking](https://docs.opencv.org/4.x/d4/dc4/group__rapid.html)
- reg. [Image Registration](https://docs.opencv.org/4.x/db/d61/group__reg.html)
- rgbd. [RGB-Depth Processing](https://docs.opencv.org/4.x/d2/d3a/group__rgbd.html)
- saliency. [Saliency API](https://docs.opencv.org/4.x/d8/d65/group__saliency.html)
- sfm. [Structure From Motion](https://docs.opencv.org/4.x/d8/d8c/group__sfm.html)
- shape. [Shape Distance and Matching](https://docs.opencv.org/4.x/d1/d85/group__shape.html)
- stereo. [Stereo Correspondance Algorithms](https://docs.opencv.org/4.x/dd/d86/group__stereo.html)
- structured_light. [Structured Light API](https://docs.opencv.org/4.x/d1/d90/group__structured__light.html)
- superres. [Super Resolution](https://docs.opencv.org/4.x/d7/d0a/group__superres.html)
- surface_matching. [Surface Matching](https://docs.opencv.org/4.x/d9/d25/group__surface__matching.html)
- text. [Scene Text Detection and Recognition](https://docs.opencv.org/4.x/d4/d61/group__text.html)
- tracking. [Tracking API](https://docs.opencv.org/4.x/d9/df8/group__tracking.html)
- videostab. [Video Stabilization](https://docs.opencv.org/4.x/d5/d50/group__videostab.html)
- viz. [3D Visualizer](https://docs.opencv.org/4.x/d1/d19/group__viz.html)
- wechat_qrcode. [WeChat QR code detector for detecting and parsing QR code.](https://docs.opencv.org/4.x/dd/d63/group__wechat__qrcode.html)
- xfeatures2d. [Extra 2D Features Framework](https://docs.opencv.org/4.x/d1/db4/group__xfeatures2d.html)
- ximgproc. [Extended Image Processing](https://docs.opencv.org/4.x/df/d2d/group__ximgproc.html)
- xobjdetect. [Extended object detection](https://docs.opencv.org/4.x/d4/d54/group__xobjdetect.html)
- xphoto. [Additional photo processing algorithms](https://docs.opencv.org/4.x/de/daa/group__xphoto.html)

# 核心功能模块

> 参考

## Mat 类

Mat 类记录在 “核心功能模块 - 基本结构” 中。

当我们使用 OpenCV 打开一张图片时，就是实例化了一个 Mat 类，这个类的本质是一个 **N-dimensional dense array class(N维密集数组类)**。说白了就是将图像转为由纯数字表示的形式。任何对图像的处理，其实就是数学计算。

Mat 类可存储实数或复值向量和矩阵、灰度或彩色图像、体素体积、向量场、点云、张量、直方图（不过，非常高维的直方图可能更好地存储在 SparseMat 中）。

**对于图像处理来说，实例化一个 Mat 对象，是一切的基础。**

# 图像文件读写模块

## imread() 函数

从文件中读取图像，返回一个 Mat 实例

## imwrite() 函数

保存图像到指定的文件中。

# Video I/O 模块

使用 OpenCV 读写视频或图像序列。

## VideoCapture 类

用于从视频文件、图像序列或摄像头中捕获视频的类。当我们打开一个视频或一个捕获设备时，就是实例化一个 VideoCapture。

- get() # 获取指定的 VideoCapture 属性
- set() # 设置指定的  VideoCapture 属性

- read() # 抓取、解码并返回下一个视频帧。
- release() # 关闭视频文件或捕获设备(比如摄像头)

## VideoWriter 类

视频写入器的类
