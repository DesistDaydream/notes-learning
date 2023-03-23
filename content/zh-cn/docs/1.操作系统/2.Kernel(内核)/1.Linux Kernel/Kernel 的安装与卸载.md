---
title: Kernel 的安装与卸载
---

# 概述

> 参考：
> - <https://mp.weixin.qq.com/s/1xRc4DzyG4c8e2XYGk28Vg>

# Ubuntu

> 参考：
> - <https://kernel.ubuntu.com/~kernel-ppa/mainline/>

## 更换内核

`awk -F' '$1=="\tmenuentry " {print i++ " : " $2}' /boot/grub/grub.cfg`

# CentOS

> 参考：
> - [ELRepo 安装文档](http://elrepo.org/tiki/tiki-index.php)
> - elrepo 的内核 rpm 包不全，暂时也不知道去哪找，先把能找到的网址都记下来
>     - <https://buildlogs.centos.org/c7-kernels.x86_64/kernel/>
>     - <https://buildlogs.centos.org/c7-kernels.x86_64/kernel/20200330213326/4.19.113-300.el8.x86_64/?C=N;O=A>

## 安装 linux 内核的存储库

```bash
rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
yum install -y https://www.elrepo.org/elrepo-release-7.el7.elrepo.noarch.rpm
```

## 安装 linux 内核

查看可用的 linux 内核版本

> 版本性质：主分支 ml(mainline)，稳定版(stable)，长期维护版 lt(longterm)

```bash
yum --disablerepo="*" --enablerepo="elrepo-kernel" list available
```

安装指定版本的 Linux 内核

```bash
yum --enablerepo=elrepo-kernel install kernel-lt-devel kernel-lt kernel-lt-headers -y
```

(可选)安装新内核工具

```bash
yum remove kernel-tools-libs.x86_64 kernel-tools.x86_64
yum --disablerepo=* --enablerepo=elrepo-kernel install kernel-lt-tools kernel-tools-libs kernel-lt-headers -y
```

## 更换默认内核

### CentOS7

```bash
# 查找需要设为默认启动的内核名称
grep "^menuentry" /boot/grub2/grub.cfg | cut -d "'" -f2
# 设置默认内核
grub2-set-default "CentOS Linux (5.4.173-1.el7.elrepo.x86_64) 7 (Core)"
# 检查默认内核版本
grub2-editenv list
```

设置完成后，执行 `reboot` 命令重启机器

卸载旧版内核

```bash
# 使用 package-cleanup 工具卸载旧内核，超过3个才会删
package-cleanup --oldkernels
# 查看 旧版内核信息
rpm -qa | grep kernel
# 卸载旧版内核
yum remove kernel-3.10.0-1127.19.1.el7.x86_64
```

### CentOS 8

查看系统安装的全部内核：

```bash
~]# grubby --info=ALL
index=0
kernel="/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64"
args="ro crashkernel=auto resume=/dev/mapper/cl-swap rd.lvm.lv=cl/root rd.lvm.lv=cl/swap net.ifnames=0 rhgb quiet intel_iommu=on $tuned_params"
root="/dev/mapper/cl-root"
initrd="/boot/initramfs-5.9.1-1.el8.elrepo.x86_64.img $tuned_initrd"
title="Red Hat Enterprise Linux (5.9.1-1.el8.elrepo.x86_64) 8.2 (Ootpa)"
id="12ab47b22fef4c02bcdc88b340d5f706-5.9.1-1.el8.elrepo.x86_64"
index=1
kernel="/boot/vmlinuz-4.18.0-193.28.1.el8_2.x86_64"
args="ro crashkernel=auto resume=/dev/mapper/cl-swap rd.lvm.lv=cl/root rd.lvm.lv=cl/swap net.ifnames=0 rhgb quiet intel_iommu=on $tuned_params"
root="/dev/mapper/cl-root"
initrd="/boot/initramfs-4.18.0-193.28.1.el8_2.x86_64.img $tuned_initrd"
title="CentOS Linux (4.18.0-193.28.1.el8_2.x86_64) 8 (Core)"
id="12ab47b22fef4c02bcdc88b340d5f706-4.18.0-193.28.1.el8_2.x86_64"
```

设置默认启动的内核

```bash
# 使用路径来指定内核，可以使用--set-default=Kernel_PATH
~]# grubby --set-default=/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64
The default is /boot/loader/entries/12ab47b22fef4c02bcdc88b340d5f706-5.9.1-1.el8.elrepo.x86_64.conf with index 0 and kernel /boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64
 ~]# grubby --default-kernel
/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64

# 使用index来指定内核，则使用--set-default-index=INDEX
~]# grubby --set-default-index=1
The default is /boot/loader/entries/12ab47b22fef4c02bcdc88b340d5f706-4.18.0-193.28.1.el8_2.x86_64.conf with index 1 and kernel /boot/vmlinuz-4.18.0-193.28.1.el8_2.x86_64
~]# grubby --default-kernel
/boot/vmlinuz-4.18.0-193.28.1.el8_2.x86_64
```

查看当前默认启动的内核

```bash
~]# grubby --default-kernel
/boot/vmlinuz-4.18.0-193.28.1.el8_2.x86_64
```

添加/删除内核启动参数：

```bash
# 对所有的内核都删除某个参数
~]# grubby --update-kernel=ALL --remove-args=intel_iommu=on

# 对所有的内核都添加某个参数
~]# grubby --update-kernel=ALL --args=intel_iommu=on

# 对某个的内核添加启动参数
~]# grubby --update-kernel=/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64 --args=intel_iommu=on

```

查看特定内核的具体信息：

```bash
~]# grubby --info=/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64
index=0
kernel="/boot/vmlinuz-5.9.1-1.el8.elrepo.x86_64"
args="ro crashkernel=auto resume=/dev/mapper/cl-swap rd.lvm.lv=cl/root rd.lvm.lv=cl/swap net.ifnames=0 rhgb quiet intel_iommu=on $tuned_params"
root="/dev/mapper/cl-root"
initrd="/boot/initramfs-5.9.1-1.el8.elrepo.x86_64.img $tuned_initrd"
title="Red Hat Enterprise Linux (5.9.1-1.el8.elrepo.x86_64) 8.2 (Ootpa)"
id="12ab47b22fef4c02bcdc88b340d5f706-5.9.1-1.el8.elrepo.x86_64"
```
