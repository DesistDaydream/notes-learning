---
title: 使用 CloneZilla 进行Linux系统的镜像完全封装和还原
---

原文链接：<https://zhuanlan.zhihu.com/p/354584111>
**CloneZilla 可以将 Linux 完整移植到另一台机器中，保证数据，分区，挂载，启动项。。所有的一切完全一致**

**注意：进行还原的机器需要与进行镜像封装的机器关键硬件配置一致，否则可能产生显卡驱动无法使用等问题**

## 1.下载 CloneZilla

下载链接：[https://clonezilla.org/downloads.php](https://link.zhihu.com/?target=https%3A//clonezilla.org/downloads.php)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-bb25fa732ef006ca39284e3f44683cbb_b.jpg)

点击第一项 alternative stable 下载

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4bcbb8b9baa704fb309c29d2ff541d28_b.jpg)

选择 2 文件格式，可以在 zip 和 iso 间进行切换

## 2.制作 CloneZilla 启动 U 盘

根据第 1 步操作下载的文件选择以下两种方式进行 U 盘制作：

**1）ISO 文件**

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

**2）zip 文件**

若选择 zip 则通过 Tuxboot ([https://tuxboot.org/download/](https://link.zhihu.com/?target=https%3A//tuxboot.org/download/))进行 U 盘制作

插入 U 盘，打开 Tuxboot，选择预下载，找到下载的 CloneZilla zip 文件，Driver 里面选择 U 盘盘符，点击 OK 就可以开始制作了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-32b271f3244fa35fd314b13d3c123177_b.jpg)

完成后如图

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-453308ec1c9937e3348602c88e44169a_b.jpg)

至此 CloneZilla 启动 U 盘制作完毕。

## 3.进行 Linux 系统的封装

**1）准备另一个 U 盘作为封装镜像的存储盘**

注意：绝对不可以使用同一个启动 U 盘作为镜像的存储盘

要准备另一个 U 盘才可完成封装工作，首先我们需要将此 U 盘进行格式化，并转换文件系统为 Linux 下通用的 ext4。**如果 U 盘存在足够的空间且格式为 ext4，则不需要格式化操作。**

- a) 将此 U 盘插入到需要封装的 Linux 机器上，在 Linux 系统上进行 ext4 的格式化
- b) 打开 terminal，输入 sudo fdisk -l 查看当前所有磁盘，在该列表中应该可以看到 U 盘的名称，如/dev/sdb
- c) 通过 sudo umount /dev/sdb 取消此 U 盘的挂载，**注意：这一步必须做**
- d) 使用 sudo mkfs.ext4 /dev/sdb 进行 ext4 格式的格式化
- e) 拔下格式化后的 U 盘，插上 CloneZilla 启动盘，重启电脑

**2）进行 Linux 系统封装**

(1) 重启电脑后，设置 BIOS 启动顺序，优先从 USB 启动，进入 CloneZilla，界面如下，选择第一项，进入图形化界面

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-2bca4e983efc7609b7a240f20db8fbb6_b.jpg)

(2) 启动完毕，选择界面语言，选择“简体中文”。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-2de00c8ee91ae22af6102e24ace38389_b.jpg)

(3) 不修改键盘配置，使用美式键盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-435ba05287f493e7f644bc203aafd4e3_b.jpg)

(4) 选择“使用再生龙”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-484c15ddf62350b48281f13c29a191dd_b.jpg)

(5) 选择“硬盘 / 分区「存到 / 来自」镜像文件”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4821ff6020786d84bff2e307d686b0e2_b.jpg)

(6) 选择第一项，插入第一步准备好的 ext4 格式 U 盘，等待 5s 后按 enter 键检测 U 盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-b375305bbd2eb5c0d3c76905428c3684_b.jpg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-801c9a1c20173e8abb74c93847395a34_b.jpg)

(7) 下图中/dev/sda 是系统盘，/dev/sdc 是要把镜像文件拷贝到的 u 盘，执行 Ctrl + c 退出此窗口

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-38cf5e91ed530b8672205bfe96668803_b.jpg)

(8) 选择插入的的 U 盘作为存放镜像文件的盘

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-86dcc178e454d03ac10051dc81ca3bc2_b.jpg)

(9) 选择第一项“no-fsck”，跳过检测

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-ca8a63a12c7240ad05dbb84e9fc3f1c0_b.jpg)

(10) 当前现在的目录名称为 u 盘的”/”, 通过上下键选择“”，之后使用“Tab”键选择“Done”即可

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-6523190d049d415d8d11f453f59512a3_b.jpg)

(11) 根据提示信息按 Enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-d09e4579aa89e88cacb8c521b30c3486_b.png)

(12) 选择“初学模式”，进入 CloneZilla 功能界面

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-426deb6d3d1e5d8d1966b31fc3153963_b.jpg)

(13) 选择“存储本机硬盘为镜像文件”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-5d8aaf82a29d7cf4135abf9f852ecd7c_b.jpg)

(14) 输入要保存的镜像名称，并按 enter 结束

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-d6337a094a761676414588886281a421_b.jpg)

(15) 选择要封装的硬盘，如果机器上有多个硬盘，此时可以选择多个要备份的硬盘，空格键选中，回车继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-a65dbd0d6faadf8acc2fc704157cab01_b.jpg)

(16) 选择第一项默认使用 zip 压缩方式

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-55752fc86cbac54e4d4f0d23e180df46_b.jpg)

(17) 选择“跳过检查与修正来源分区的文件系统”，这一步是用来检查硬盘的文件系统是否符合要求

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-384ae83f7ff238930dad965944266a78_b.jpg)

(18) 选择“否”，这一步是用来检查保存的镜像是否可以被还原的

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-0816cacf4b7fc828d8cc5ebf3e6dc10c_b.jpg)

(19) 选择是否要对镜像进行加密，这里选择“不加密”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4750ed41e41ccfd9b31336e038ca342a_b.jpg)

(20) 这一步用于选择镜像封装完成后要执行的操作，可选“关机”、“重启”等，这里选择“当操作执行完后再选择”

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-c54a7964c1e3b205a16712610fcd7fe1_b.jpg)

(21) 输入 y 后按“enter”回车，即可开始进行镜像封装了

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-ffb2f474690186d539f5f4ac93b2fbf0_b.jpg)

(22) 镜像封装中

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-4c8fd45ccb63cd91da2691331cb6f327_b.jpg)

(23) 镜像封装完成，按 enter 键继续

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-32ccb6c402a0c4aecde904f9c8d6e151_b.jpg)

(24) 选择“关机”或“重启”，**注意：若(20)选择“关机”或“重启”，不会有该界面出现，直接关机或重启；若选择重启，待屏幕黑屏后需将启动盘拔出。**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/1f12dbea-164d-4f41-9cd2-fd92952d6310/v2-74c60fc40ab08d48b204f3a00adbec58_b.jpg)

## 4.进行 Linux 系统还原

(1) - (12)步与进行封装的操作一致

(13) 选择“还原镜像文件到本机硬盘”

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
