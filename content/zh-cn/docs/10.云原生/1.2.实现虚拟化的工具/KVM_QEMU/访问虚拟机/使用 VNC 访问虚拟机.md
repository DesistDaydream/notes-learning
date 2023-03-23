---
title: 使用 VNC 访问虚拟机
---

# 概述

> 参考：

在 [rfbproto/rfbproto 的 #18 issue](https://github.com/rfbproto/rfbproto/issues/18) 中，表示 QEMU 中内置了 VNC/RFB，并对该协议进行了一些修改以支持一些额外的功能。所以，我们可以使用 VNC 客户端，使用 QEMU/KVM 为虚拟机暴露的 VNC 端口，连接到虚拟机(不管是 图形界面的虚拟机 还是 命令行界面的虚拟机)。

## 配置 VNC 监听地址

修改 qemu.conf 配置文件

```bash
# VNC is configured to listen on 127.0.0.1 by default.
# To make it listen on all public interfaces, uncomment
# this next option.
#
# NB, strong recommendation to enable TLS + x509 certificate
# verification when allowing public access
#
vnc_listen = "0.0.0.0"
```

说明 1：VNC 默认绑定 127.0.0.1，在配置文件里指定 VNC  绑定 0.0.0.0 IP,就不用在创建虚拟机时修改 vnclisten 参数了。
说明 2：在虚拟主机上有很多个虚拟机的时候，需要指定每个虚拟机的端口。

若虚拟机已存在，则修改虚拟机的 xml 配置文件的如下部分\*\*          \*\*

```xml
    <graphics type='vnc' port='-1' autoport='yes' listen='0.0.0.0'>
      <listen type='address' address='0.0.0.0'/>
    </graphics>
```

port 为 -1 表示自动分配端口，可以指定具体的端口，并将 autoport 设置为 no

启动并测试 VNC，查看 VNC 端口

```bash
~]# ss -ntlp | grep 590
LISTEN     0      1      127.0.0.1:5900                     *:*                   users:(("qemu-kvm",pid=24077,fd=69))
LISTEN     0      1      127.0.0.1:5901                     *:*                   users:(("qemu-kvm",pid=24190,fd=69))
LISTEN     0      1      127.0.0.1:5902                     *:*                   users:(("qemu-kvm",pid=27027,fd=61))
LISTEN     0      1      127.0.0.1:5903                     *:*                   users:(("qemu-kvm",pid=27124,fd=61))
LISTEN     0      1      127.0.0.1:5904                     *:*                   users:(("qemu-kvm",pid=26750,fd=57))
LISTEN     0      1      127.0.0.1:5905                     *:*                   users:(("qemu-kvm",pid=50127,fd=57))
LISTEN     0      5            *:5906                     *:*                   users:(("Xvnc",pid=45939,fd=9))
LISTEN     0      1            *:5907                     *:*                   users:(("qemu-kvm",pid=46148,fd=23))
LISTEN     0      5         [::]:5906                  [::]:*                   users:(("Xvnc",pid=45939,fd=10))
```

使用 vnc 进行登录

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1653399424256-1aebe374-71df-4b7e-adc6-7a75259d03a0.png)

在 virt-manager 里 VNC 的修改方式

![image.jpeg](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1616123512562-2e6de574-d197-4b52-a26d-40b0c09db446.jpeg)

可以通过如下命令设置 vnc

virt-install --import --name test --memory 2048 --vcpus 2 --os-variant centos7.0 --disk /var/lib/libvirt/images/test.qcow2,size=20 --network bridge=br0,model=virtio --graphics vnc,listen=0.0.0.0,port=5910

备注：第一次在使用 vnc 访问虚拟机的时候会出现一闪就不见了的问题？具体的解决方法如下：

依次打开 vnc 客户端--->依次点击 option--->Advanced--->Expert--->找到 ColourLevel,默认的值是 pal8，修改为 rgb222 或 full.,见下图一图二:1,2,3

![image.jpeg](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1616123512600-5a0611f4-2ac1-4cf8-b02f-bf100cda25da.jpeg)

![image.jpeg](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1616123512705-58b3ebe2-934c-4b6f-90fd-2c5749dbac8c.jpeg)

# 常见问题

## Protocol error: invalid message type ...

若连接的目标是 Windows，可能会出现如下报错：`Protocol error: invalid message type ...`
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1653399468846-b581842e-f3d2-4712-a320-58fd1117011b.png)
正常来说，RFB 协议会根据网络质量自动调整画面质量。但如果 KVM 服务器的 VNC 服务不能提供该功能或者无法根据网络质量适配画面质量，就会出现以上错误提示。

此时需要修改设置中的 Picture quality 为 High 或 Medium 即可
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1653399522688-ecf9beda-c5c3-4407-967d-d7cf2807a241.png)

## The connection closed unexpectedly

若通过 virt-manager 已经打开了 VNC 显示模式的虚拟机，则通过 VNC 客户端连接时将会出现如下提示
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/vbc6tk/1653399960014-7c39aed9-6b52-4c49-a66f-7eab9f27e3f3.png)
此时只需要关闭 virt-manager 打开的虚拟机即可
