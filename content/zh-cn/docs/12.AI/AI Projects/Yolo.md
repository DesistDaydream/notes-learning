---
title: YOLO
linkTitle: YOLO
date: 2024-08-30T09:17
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，ultralytics/ultralytics](https://github.com/ultralytics/ultralytics)
> - [ultralytics 官网](https://www.ultralytics.com/)
> - https://medium.com/@gary.tsai.advantest/yolo-%E7%B3%BB%E5%88%97%E5%A4%A7%E8%A3%9C%E5%B8%96-yolov7-b1ce83a7035

[YOLO](https://arxiv.org/abs/1506.02640)(You Only Look Once）是一种流行的物体检测和图像分割模型，由华盛顿大学的约瑟夫-雷德蒙（Joseph Redmon）和阿里-法哈迪（Ali Farhadi）开发。YOLO 于 2015 年推出，因其高速度和高精确度而迅速受到欢迎。

2015 年 _Joseph Redmon_ 提出的 YOLO 橫空出世，从诞生的那一刻起就标榜「高精度」、「高效率」、「高实用性」，為 One-Stage 方法在物体侦测演算法里拉开序幕。

- **YOLOv1** (2016) Joseph Redmon
  - [You Only Look Once: Unified, Real-Time Object Detection](https://arxiv.org/abs/1506.02640)
- [**YOLOv2**](https://github.com/longcw/yolo2-pytorch) (2017) Joseph Redmon
  - [YOLO9000: Better, Faster, Stronger](https://arxiv.org/abs/1612.08242)
- [**YOLOv3**](https://github.com/ultralytics/yolov3) (2018) Joseph Redmon
  - [YOLOv3: An Incremental Improvement](https://arxiv.org/abs/1804.02767)
- 突发
  - 然而，2020 年 *约瑟夫·雷德蒙* 突然投下一枚震撼弹，他受够 YOLO 不断被运用在军事应用以及个人隐私，宣布停止电脑视觉相关的研究。
  - ![https://x.com/pjreddie/status/1230524770350817280|400](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/202410151103187.png)
- [**YOLOv4**](https://github.com/AlexeyAB/darknet) (2020) Alexey Bochkovskiy
  - [YOLOv4: Optimal Speed and Accuracy of Object Detection](https://arxiv.org/abs/2004.10934)
- [**YOLOv5**](https://github.com/ultralytics/yolov5) 进一步提高了模型的性能，并增加了超参数优化、集成实验跟踪和自动导出为常用导出格式等新功能。
- [**YOLOv6**](https://github.com/meituan/YOLOv6) (2022) 由[美团](https://about.meituan.com/)开源，目前已用于该公司的许多自主配送机器人。
  - [YOLOv6: A Single-Stage Object Detection Framework for Industrial Applications](https://arxiv.org/abs/2209.02976)
- [**YOLOv7**](https://github.com/WongKinYiu/yolov7) 增加了额外的任务，如 COCO 关键点数据集的姿势估计
  - [YOLOv6 v3.0: A Full-Scale Reloading](https://arxiv.org/abs/2301.05586)
  - [YOLOv7: Trainable bag-of-freebies sets new state-of-the-art for real-time object detectors](https://arxiv.org/abs/2207.02696)
- [**YOLOv8**](https://pypi.org/project/ultralytics/8.0.0/) (2023-01) 由 Ultralytics 提供。YOLOv8 支持全方位的视觉 AI 任务，包括[检测](https://docs.ultralytics.com/tasks/detect/)、[分割](https://docs.ultralytics.com/tasks/segment/)、[姿态估计](https://docs.ultralytics.com/tasks/pose/)、[跟踪](https://docs.ultralytics.com/modes/track/)和[分类](https://docs.ultralytics.com/tasks/classify/)。这种多功能性使用户能够在各种应用和领域中利用 YOLOv8 的功能。
- [**YOLOv9**](https://github.com/WongKinYiu/yolov9) (2024) 由原YOLOv7团队打造，引入了可编程梯度信息 （PGI） 和广义高效层聚合网络 （GELAN） 等创新方法。
  - [YOLOv9: Learning What You Want to Learn Using Programmable Gradient Information](https://arxiv.org/abs/2402.13616)
- [**YOLOv10**](https://github.com/THU-MIG/yolov10) (2024) 是由清华大学的研究人员使用 Ultralytics Python 包创建的。该版本通过引入消除非极大值抑制 (NMS) 要求的端到端头，提供了实时对象检测方面的改进。
  - [YOLOv10: Real-Time End-to-End Object Detection](https://arxiv.org/abs/2405.14458)
- YOLO11 

# 训练

https://docs.ultralytics.com/modes/train/

训练前需要准备如下内容：

- data.yaml 文件
- dataset

最基本的训练代码如下：

```python
from ultralytics import YOLO
# 加载预训练模型（建议用于训练）
model = YOLO("yolo11n.pt")
# 训练模型。指定 coco.yaml 为 data.yaml。训练周期 100，TODO: 640 是干什么的？
results = model.train(data="coco.yaml", epochs=100, imgsz=640)
```

# 数据集

> 参考：
>
> - https://docs.ultralytics.com/datasets/
> - https://docs.ultralytics.com/datasets/#contribute-new-datasets

Ultralytics 提供对各类 **Dataset(数据集)** 的支持，以便进行计算机视觉任务，如 对象检测、实例分割、姿态估计、分类、多目标跟踪、etc. 

- **[Object detection](#Object%20detection)(对象检测) 数据集** # 通过在每个对象周围绘制边界框来检测和定位图像中的对象。
- **[Instance segmentation](#Instance%20segmentation)(实例分割) 数据集** # 在像素级别识别和定位图像中的对象。Object detection 在识别到对象后是用矩形框框起来的，而 Instance segmentation 则是在识别到对象的基础上，在像素级别对物体进行染色
- **[Pose estimation](#Pose%20estimation)(姿态估计) 数据集** # 在识别到对象后，识别对象的姿态。
- etc.

上面这些种类的数据集中，有一个名为 **Common Objects in Context(COCO)** 的数据集，COCO 是一个通用的大规模用于 对象检测、实例分割、姿态估计 的数据集，包含 80 个对象类别、超过 200K 个已标记的图像。

## 创建自己的数据集

https://docs.ultralytics.com/datasets/#contribute-new-datasets

1. **Collect Images(收集图像)** # 收集用于训练的图像。
2. **Annotate Images(注释图像)** # 根据想要训练的任务，使用 边界框、片段、关键点 为图像添加注释。人话：**数据标注**
3. **Export Annotations(导出注释)** # 将这些注释转换为 Ultralytics 支持的 YOLO `*.txt` 文件格式。
4. **Organize Dataset(组织数据集)** # 将图像、注释以如下目录结构存放。应该有 images/ 和 labels/ 顶级目录，并在每个目录中都有一个 train/ 和 val/ 子目录。 images 存放收集到的图像，labels 存放导出的注释。images/ 下的如果有 000000000009.jpg 文件，那对应的 labels/ 下应该有个同名不同后缀的 000000000009.txt 文件。
   1. Notes: 这个目录结构在由于实际情况可能有的类型的数据集并不完全相同，基于数据组织的便利性，可能会把 train/ 和 val/ 放在顶级目录，下级目录可能是以对象类型命名。 
```bash
datasets/
└── coco8
    ├── images
    │   ├── train
    │   │   ├── 000000000009.jpg
    │   │   └── 000000000034.jpg
    │   └── val
    │       ├── 000000000036.jpg
    │       └── 000000000061.jpg
    └── labels
        ├── train
        │   ├── 000000000009.txt
        │   └── 000000000034.txt
        └── val
            ├── 000000000036.txt
            └── 000000000061.txt
```
5. **创建 `data.yaml` 文件** # 创建一个描述数据集、类和其他必要信息的 [data.yaml](#data.yaml) 文件。
6. **Optimize Images(优化图像)(可选的)** # 如果您想减小数据集的大小以提高处理效率，可以使用以下代码优化图像。这不是必需的，但建议用于较小的数据集大小和更快的下载速度。
7. **压缩数据集** # 将整个数据集文件夹压缩为 zip 文件。
8. **Document and PR** # 创建一个文档页面来描述数据集以及它如何适应现有框架。之后，提交 Pull Request (PR)。有关如何提交 PR 的更多详细信息，请参阅 [Ultralytics 贡献指南](https://docs.ultralytics.com/help/contributing/)。

第 6 和 7 步，ultralytics 提供了函数可以直接使用代码处理。通过遵循这些步骤，可以提供一个与 Ultralytics 现有结构良好集成的新数据集

```python
from pathlib import Path

from ultralytics.data.utils import compress_one_image
from ultralytics.utils.downloads import zip_directory

# Define dataset directory
path = Path("path/to/dataset")

# Optimize images in dataset (optional)
for f in path.rglob("*.jpg"):
    compress_one_image(f)

# Zip dataset into 'path/to/dataset.zip'
zip_directory(path)
```

## data.yaml

data.yaml 时 Ultralytics YOLO 数据集使用的 [YAML](docs/2.编程/无法分类的语言/YAML.md) 格式文件，可以定义数据集所在目录、训练、验证、测试 图像目录、数据集中对象分类的字典。在  [这里](https://github.com/ultralytics/ultralytics/tree/main/ultralytics/cfg/datasets) 找到各类数据集的 data.yaml。

> Notes: data.yaml 是一种抽象的叫法，本质就是配置文件。各种数据集使用 data.yaml 时，可以指定任意名称但只要符合文件内容格式的 YAML 文件。

e.g. 对象检测 类型的数据集配置通常是下面这样的：

```yaml
path: ../datasets/coco8 # 数据集根目录
train: images/train # 相对于 path 的训练目录
val: images/val # 相对于 path 的验证目录
test: # （可选的）相对于 path 的测试目录

# 类型
names:
    0: person
    # ......略
    5: bus
    # ......略
    79: toothbrush
```

下图是使用上面示例的 data.yaml 训练后模型识别的结果，可以看到有 person 和 bus，并不是单纯的数字了，在右边官网[示例图片](https://docs.ultralytics.com/datasets/detect/coco/#sample-images-and-annotations)中，识别出的对象都是**数字注释**的。

![400](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/detected_bus.png)![400](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/yolo/mosaiced-coco-dataset-sample.jpg)

当我们训练模型时，下面代码就会指定要使用的 data.yaml 文件（这示例的 data.yaml 名为 coco.yaml）

```python
from ultralytics import YOLO

# Load a model
model = YOLO("yolo11n.pt")  # load a pretrained model (recommended for training)

# Train the model
results = model.train(data="coco.yaml", epochs=100, imgsz=640)
```

然后根据 data.yaml 中的 path、train、val 定义的路径，从目录中读取 图像文件 及 图像注释文件。这些文件组织结构像这样：

```bash
datasets/
└── coco8
    ├── images
    │   ├── train
    │   │   ├── 000000000009.jpg
    │   │   └── 000000000034.jpg
    │   └── val
    │       ├── 000000000036.jpg
    │       └── 000000000061.jpg
    └── labels
        ├── train
        │   ├── 000000000009.txt
        │   └── 000000000034.txt
        └── val
            ├── 000000000036.txt
            └── 000000000061.txt
```

# 数据集种类

## Object detection

https://docs.ultralytics.com/datasets/detect/

[Object detection](docs/12.AI/计算机视觉/Object%20detection.md)(对象检测)

images 是图片，labels 图片的标签

```bash
datasets/
└── coco8
    ├── LICENSE
    ├── README.md
    ├── images
    │   ├── train
    │   │   ├── 000000000009.jpg
    │   │   ├── 000000000025.jpg
    │   │   ├── 000000000030.jpg
    │   │   └── 000000000034.jpg
    │   └── val
    │       ├── 000000000036.jpg
    │       ├── 000000000042.jpg
    │       ├── 000000000049.jpg
    │       └── 000000000061.jpg
    └── labels
        ├── train
        │   ├── 000000000009.txt
        │   ├── 000000000025.txt
        │   ├── 000000000030.txt
        │   └── 000000000034.txt
        ├── train.cache
        ├── val
        │   ├── 000000000036.txt
        │   ├── 000000000042.txt
        │   ├── 000000000049.txt
        │   └── 000000000061.txt
        └── val.cache
```

label 是对应相同文件名的图片中的 ROI，共 5 个数字

- ROI 类型
- ROI 左上角 x 轴坐标
- ROI 左上角 y 轴坐标
- ROI 右下角 x 轴坐标
- ROI 右下角 y 轴坐标

```bash
~]# cat datasets/coco8/labels/train/000000000025.txt
23 0.770336 0.489695 0.335891 0.697559
23 0.185977 0.901608 0.206297 0.129554
```

## Instance segmentation

https://docs.ultralytics.com/datasets/segment/

## Pose estimation

https://docs.ultralytics.com/datasets/pose/

[Pose estimation](docs/12.AI/计算机视觉/Pose%20estimation.md)(姿态估计)

YAML 格式 https://docs.ultralytics.com/datasets/pose/#dataset-yaml-format

# 模型

> 参考：
>
> - [ultralytics 文档，模型](https://docs.ultralytics.com/models/)

https://docs.ultralytics.com/models/yolov10/#model-variants

Ultralytics 的 YOLO 有多种模型规模，以满足不同的应用需求：

- **YOLO-N** # **Nano** 适用于资源极其有限的环境的。适合移动设备和快速测试
- **YOLO-S** # **Small** 平衡速度和准确性的小型版本。适合嵌入式设备和一般性测试
- **YOLO-M** # **Medium** 用于通用用途。
- **YOLO-B** # **Balanced** 增加宽度以提高精度。
- **YOLO-L** # **Large** 增加计算资源为代价获得更高的精度。
- **YOLO-X** # **Extra-large** 可实现最大精度和性能。适合服务器处理高精度任务

# 关联文件与配置

**%APPDATA%/Ultralytics/** # Ultralytics 出的 YOLO 模型在使用时保存的 字体、配置 所在目录

- **./settings.json** # 有一些保存文件的目录配置。

# CLI

https://docs.ultralytics.com/usage/cli/

## Syntax(语法)

**SubCommand**

- **train** # Train(训练) 模型
- **predict** # 使用经过训练的 YOLO11n 模型对图像进行 predictions(预测)。
- **val** # 在 COCO8 数据集上 Validate(验证) 经过训练的 YOLO11n 模型的准确性。不需要参数，因为模型保留其训练数据和参数作为模型属性。
- **export** # 将 YOLO11n 模型导出为不同的格式，例如 ONNX、CoreML、etc. 。