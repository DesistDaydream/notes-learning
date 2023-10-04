---
title: "OpenCV-Python"
---

# 概述

> 参考：
> 
> - [GitHub 项目，opencv/opencv-python](https://github.com/opencv/opencv-python)
> - [OpenCV-Python 4.x 官方教程](https://docs.opencv.org/4.x/d6/d00/tutorial_py_root.html)

OpenCV-Python 是一个旨在解决计算机视觉问题的 Python 库。

OpenCV-Python 利用高度优化的 [NumPy](/docs/12.人工智能/科学计算/NumPy.md) 库进行数值操作，其语法类似于 MATLAB。所有的 OpenCV 数组结构都会被转换为 Numpy 数组。这也使得集成其他使用 Numpy 的库，如 SciPy 和 Matplotlib 更容易。

# 安装

> 参考：
> - [OpenCV-Python 4.x 官方文档，在 Windows 上安装 OpenCV-Python](https://docs.opencv.org/4.x/d5/de5/tutorial_py_setup_in_windows.html)

```bash
pip install opencv-python
```

# 图像处理

> 参考：
> 
> - [官方文档，OpenCV-Python 教程-核心业务-图像的基本操作](https://docs.opencv.org/4.x/d3/df2/tutorial_py_basic_ops.html)
> - [官方文档，OpenCV-Python 教程-OpenCV 中的图像处理](https://docs.opencv.org/4.x/d2/d96/tutorial_py_table_of_contents_imgproc.html)

### Hello World

> 参考：
> 
> - [官方文档，OpenCV-Python 教程-OpenCV 中的 Gui 功能-图像入门](https://docs.opencv.org/4.x/db/deb/tutorial_display_image.html)

```python
import sys
import cv2

if __name__ == "__main__":
    # imread() 读取图片，并将图片实例化为一个 Mat 对象
    # 可以接收参数以指定我们想要的图像格式
    # - IMREAD_COLOR 以 BGR 8 位格式加载图像。这是此处使用的默认值。
    # - IMREAD_UNCHANGED 按原样加载图像（包括 alpha 通道，如果存在）。其实就是将图片变为黑白的了
    # - IMREAD_GRAYSCALE 将图像作为强度加载
    # image = cv2.imread("images_cn/BT1-001R.png")
    image = cv2.imread(cv2.samples.findFile("images_cn/BT1-001R.png"))

    if image is None:
        sys.exit("无法读取图片")

    # imshow() 打开一个窗口，并显示图片
    cv2.imshow("窗口的标题", image)
    # waitKey() 等待用户按键，若按键为 ESC，则返回 -1。如果不等待，那么打开的窗口瞬间就会关闭
    k = cv2.waitKey(0)

    # ord() 用于等待键盘输入，0 表示任意键。这里按下 s 则会将图片保存到本地
    if k == ord("s"):
        # imwrite() 将 Mat 对象写入图片并保存
        cv2.imwrite("starry_night.png", image)
```
