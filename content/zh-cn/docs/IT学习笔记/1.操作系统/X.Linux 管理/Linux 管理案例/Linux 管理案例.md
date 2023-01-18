---
title: Linux 管理案例
---

# 重置 Linux 的 root 密码

## 进入紧急模式

详见《[Linux 的紧急模式或救援模式](/docs/IT学习笔记/1.操作系统/X.Linux%20 管理/性能优化%20 与%20 故障处理/Linux%20 的紧急模式或救援模式.md 与 故障处理/Linux 的紧急模式或救援模式.md)》

## 修改密码

- 切换到原系统执行：`chroot /sysroot/`
- 更改 root 密码：`passwd root`
- 在/目录下创建一个.autorelabel 文件，而有这个文件存在，系统在重启时就会对整个文件系统进行 relabeling
  - `touch /.autorelabel`
- `exit`
- `reboot`

# 修改网卡名

centos 系统

- vi /etc/default/grub
  - GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0 crashkernel=auto rd.lvm.lv=myvg/root rd.lvm.lv=myvg/swap rhgb quiet"
  - 注意，标红位置改为自己的 lvm 中 volume group 的名字
  - 主要就是添加紫色内容的字符串
- grub2-mkconfig -o /boot/grub2/grub.cfg
- mv /etc/sysconfig/network-scripts/ifcfg-ens33 /etc/sysconfig/network-scripts/ifcfg-eth0
- sed -i "s/ens33/eth0/g" /etc/sysconfig/network-scripts/ifcfg-eth0

ubuntu 系统

- 修改 grub 文件
  - vim /etc/default/grub
- 查找
  - GRUB_CMDLINE_LINUX=""
- 修改为
  - GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"
- 重新生成 grub 引导配置文件
  - grub-mkconfig -o /boot/grub/grub.cfg
- 修改网络配置 ens32 为 eth0
  - vim /etc/netplan/01-netcfg.yaml

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/gvagsg/1616163849544-f4eac668-9a60-40ef-b291-c28f82e1e661.jpeg)
