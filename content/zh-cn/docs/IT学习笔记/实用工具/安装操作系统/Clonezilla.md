---
title: Clonezilla
---

## 概述

> 参考：
> - [官网](https://clonezilla.org/)
> - [知乎，使用再生龙CloneZilla进行Linux系统的镜像完全封装和还原](https://zhuanlan.zhihu.com/p/354584111)
> - 其他实践：
> 	- https://blog.csdn.net/zhaoxinfan/article/details/126474777
> 	- https://blog.csdn.net/zhangjia453/article/details/115353982

Clonezilla 是类似于 [True Image](http://en.wikipedia.org/wiki/Acronis_True_Image) 或 [Norton Ghost](http://en.wikipedia.org/wiki/Ghost_%28software%29) 的分区和磁盘映像/克隆程序。它可以帮助您进行系统部署、裸机备份和恢复。可以使用三种类型的 Clonezilla：
- [Clonezilla live](https://clonezilla.org/clonezilla-live.php) # Clonezilla live 允许您使用 CD/DVD 或 USB 闪存驱动器启动和运行 clonezilla（仅限单播）
- [Clonezilla lite server](https://clonezilla.org/show-live-doc-content.php?topic=clonezilla-live/doc/11_lite_server) # Clonezilla 精简版服务器允许您使用 Clonezilla live 进行大规模克隆（支持单播、广播、多播、比特流）
- [Clonezilla SE](https://clonezilla.org/clonezilla-SE/) # Clonezilla SE 包含在 DRBL 中，因此必须首先设置 DRBL 服务器才能使用 Clonezilla 进行大规模克隆（支持单播、广播和组播）

Clonezilla live 适用于单机备份和恢复。虽然 Clonezilla 精简版服务器或 SE 用于大规模部署，但它可以同时克隆多台（40 多台！）计算机。Clonezilla 仅保存和恢复硬盘中使用过的块。这提高了克隆效率。对于 42 节点集群中的一些高端硬件，报告了以 8 GB/分钟的速率恢复的多播。

**CloneZilla 可以将 Linux 完整移植到另一台机器中，保证数据，分区，挂载，启动项。。所有的一切完全一致**
> 注意：进行还原的机器需要与进行镜像封装的机器关键硬件配置一致，否则可能产生显卡驱动无法使用等问题

Clonezilla Live 本身就是一个小型的 Liunx 发行版，这就像是 Linux 版的 WinPE 一样。除了图形页面可供操作外，还可以通过命令行，执行诸如 ssh 之类的命令。这也就意味着，Clonezilla 可以读写远程存储设备中，封装好的 Linux 镜像可以直接写入到 NFS、S3 中，还原 Linux 时，也可以从这些远程存储设备中读取镜像。

Clonezilla 源码实际上是运行程序本身，因为它们使用脚本(bash 或 perl)编写的。最初始保存在 [NCHC 存储库](https://free.nchc.org.tw/clonezilla-live/)中。在这个存储库中，我们可以从 [experimental/](https://free.nchc.org.tw/clonezilla-live/experimental/) 目录下找到 ARM 版本的 Clonezilla。

# CloneZilla Live 安装与使用

从[此处](https://clonezilla.org/downloads.php)可以看到 Clonezilla Live 的所有下载链接

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-bb25fa732ef006ca39284e3f44683cbb_b.jpg)

点击第一项 alternative stable 下载

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4bcbb8b9baa704fb309c29d2ff541d28_b.jpg)

选择 2 文件格式，可以在 zip 和 iso 间进行切换

## 制作 CloneZilla 启动 U 盘

根据下载的文件类型的不同，有两种方式进行 U 盘制作

### ISO 文件

若选择 iso 则通过**软碟通 UltraISO** ([https://cn.ultraiso.net/xiazai.html](https://link.zhihu.com/?target=https%3A//cn.ultraiso.net/xiazai.html))在 windows 设备上进行 U 盘刻录。

安装完**软碟通 UltraISO**后，插上待制作的 U 盘，打开软件，点击“打开”按钮，选择 CloneZilla 的 iso 文件打开

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-5bd64afb388461e5f2f76f9881771d4c_b.jpg)

然后，点击顶部菜单栏的“启动”->“写入硬盘镜像”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-a936e71da8a97bf955bcba7a0376e1e2_b.jpg)

选择 U 盘，点击“格式化”，格式选择 FAT32，格式化完毕后，写入方式选择“USB-HDD+”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-0b8205cc4795bc658913251966643810_b.jpg)

点击“编写启动”->“syslinux”，写入 syslinux

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-614e796adc4303da7cfb1e5c4f1422b3_b.jpg)

写入 CloneZilla 镜像

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-5d490a422045f4f86c4ccb6567b45808_b.jpg)

至此 CloneZilla 启动 U 盘制作完毕

### zip 文件

若选择 zip 则通过 Tuxboot ([https://tuxboot.org/download/](https://link.zhihu.com/?target=https%3A//tuxboot.org/download/))进行 U 盘制作

插入 U 盘，打开 Tuxboot，选择预下载，找到下载的 CloneZilla zip 文件，Driver 里面选择 U 盘盘符，点击 OK 就可以开始制作了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-32b271f3244fa35fd314b13d3c123177_b.jpg)

完成后如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-453308ec1c9937e3348602c88e44169a_b.jpg)

至此 CloneZilla 启动 U 盘制作完毕。

# 封装Linux系统

启动服务器，优先从 USB 启动，进入 CloneZilla，界面如下：
- 选择第一项 **Clonezilla live (VGA 800 * 600)**，进入图形化界面

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-2bca4e983efc7609b7a240f20db8fbb6_b.jpg)

(2) 启动完毕，选择界面语言，选择“简体中文”。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-2de00c8ee91ae22af6102e24ace38389_b.jpg)

(3) 不修改键盘配置，使用美式键盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-435ba05287f493e7f644bc203aafd4e3_b.jpg)

(4) 选择“使用再生龙”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-484c15ddf62350b48281f13c29a191dd_b.jpg)

(5) 选择“device-image 硬盘 / 分区「存到 / 来自」镜像文件”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4821ff6020786d84bff2e307d686b0e2_b.jpg)

## 设置镜像存储或读取的位置

为 Clonezilla 确定使用镜像文件的方式，比如本机分区(硬盘、U盘等)、远程主机目录等。这些方式任选其一即可

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/0-picgo/clonezilla_store_path.png)

### local_dev

**1）准备另一个 U 盘作为封装镜像的存储盘**

注意：绝对不可以使用同一个启动 U 盘作为镜像的存储盘

要准备另一个 U 盘才可完成封装工作，首先我们需要将此 U 盘进行格式化，并转换文件系统为 Linux 下通用的 ext4。**如果 U 盘存在足够的空间且格式为 ext4，则不需要格式化操作。**

- a) 将此 U 盘插入到需要封装的 Linux 机器上，在 Linux 系统上进行 ext4 的格式化
- b) 打开 terminal，输入 sudo fdisk -l 查看当前所有磁盘，在该列表中应该可以看到 U 盘的名称，如/dev/sdb
- c) 通过 sudo umount /dev/sdb 取消此 U 盘的挂载，**注意：这一步必须做**
- d) 使用 sudo mkfs.ext4 /dev/sdb 进行 ext4 格式的格式化
- e) 拔下格式化后的 U 盘，插上 CloneZilla 启动盘，重启电脑

**2) 继续操作 Clonezilla**

选择 **local_dev**，插入第一步准备好的 ext4 格式 U 盘，等待 5s 后按 enter 键检测 U 盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-801c9a1c20173e8abb74c93847395a34_b.jpg)

下图中/dev/sda 是系统盘，/dev/sdc 是要把镜像文件拷贝到的 u 盘，执行 Ctrl + c 退出此窗口

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-38cf5e91ed530b8672205bfe96668803_b.jpg)

选择插入的的 U 盘作为存放镜像文件的盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-86dcc178e454d03ac10051dc81ca3bc2_b.jpg)

选择第一项“no-fsck”，跳过检测

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-ca8a63a12c7240ad05dbb84e9fc3f1c0_b.jpg)

当前现在的目录名称为 u 盘的”/”, 通过上下键选择“”，之后使用“Tab”键选择“Done”即可

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-6523190d049d415d8d11f453f59512a3_b.jpg)

根据提示信息按 Enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-d09e4579aa89e88cacb8c521b30c3486_b.png)

### ssh_server

### S3



```
export AK="XXX"
export SK="YYY"
export ENDPOINT="ZZZ"

tee /root/.passwd-s3fs <<EOF
${AK}:${SK}
EOF

chmod 600 /root/.passwd-s3fs

s3fs my-ocs-img /home/partimag -o passwd_file=/root/.passwd-s3fs -o url="https://${ENDPOINT}"
```

## 配置运行时参数

选择“初学模式”，进入 CloneZilla 功能界面

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/0-picgo/clonezilla_config_args_1.png)


## 选择要执行的操作

选择“savedisk 存储本机硬盘为镜像文件”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-5d8aaf82a29d7cf4135abf9f852ecd7c_b.jpg)

输入要保存的镜像名称，并按 enter 结束

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-d6337a094a761676414588886281a421_b.jpg)

选择要封装的硬盘，如果机器上有多个硬盘，此时可以选择多个要备份的硬盘，空格键选中，回车继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-a65dbd0d6faadf8acc2fc704157cab01_b.jpg)

选择第一项默认使用 zip 压缩方式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-55752fc86cbac54e4d4f0d23e180df46_b.jpg)

选择“跳过检查与修正来源分区的文件系统”，这一步是用来检查硬盘的文件系统是否符合要求

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-384ae83f7ff238930dad965944266a78_b.jpg)

选择“否”，这一步是用来检查保存的镜像是否可以被还原的

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-0816cacf4b7fc828d8cc5ebf3e6dc10c_b.jpg)

选择是否要对镜像进行加密，这里选择“不加密”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4750ed41e41ccfd9b31336e038ca342a_b.jpg)

这一步用于选择镜像封装完成后要执行的操作，可选“关机”、“重启”等，这里选择“当操作执行完后再选择”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-c54a7964c1e3b205a16712610fcd7fe1_b.jpg)

输入 y 后按“enter”回车，即可开始进行镜像封装了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-ffb2f474690186d539f5f4ac93b2fbf0_b.jpg)

镜像封装中

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4c8fd45ccb63cd91da2691331cb6f327_b.jpg)

镜像封装完成，按 enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-32ccb6c402a0c4aecde904f9c8d6e151_b.jpg)

选择“关机”或“重启”，**注意：若选择“关机”或“重启”，不会有该界面出现，直接关机或重启；若选择重启，待屏幕黑屏后需将启动盘拔出。**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-74c60fc40ab08d48b204f3a00adbec58_b.jpg)

# 还原Linux系统

在选择执行操作前，所有步骤与[封装Linux系统](#封装Linux系统)前面的步骤一样

## 选择要执行的操作

选择“restoredisk 还原镜像文件到本机硬盘”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-457cda3e0e53500d3c19048f82fb3f83_b.jpg)

(14) 选择需要还原的镜像文件

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-c6b1127ea5cf93def40a70133c995fb5_b.jpg)

(15) 选择需要还原到的本机硬盘，**注意：该硬盘会被格式化!**

**选择硬盘时，如果镜像是从多块硬盘保存下来的，则可以选择多块硬盘；**

**如果封装镜像的硬盘是 2T，则不能在 1T 的硬盘上还原！！！**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-41139841999062e69231ae801dab04e0_b.jpg)

(16) 询问是否需要检查镜像完整性，这里选择“否”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-d6ad903a054dc45043c1707fb0c8b30e_b.png)

(17) 询问完成还原后的操作，这一步与封装时一致，选择完后按 enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-7bba000cc0a0bff9143f216b47ef4af6_b.jpg)

(18) 输入 y，并按 enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-1b7bc181974267c97c183a96d4a1082b_b.jpg)

(19) 正在还原

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-e2f5545b2635c8be2a24b029b8d2535c_b.jpg)

(20) 待还原结束后按 enter 键继续，选择重启，待黑屏后拔掉 U 盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-237df2b6aa6b869c45b04d48874aab23_b.jpg)
