---
title: NFS 部署
---

# 概述

> 参考：
> - [Ubuntu 官方文档](https://ubuntu.com/server/docs/service-nfs)

# 服务端部署

## 安装 NFS Server

### 通过 Linux 的包管理器部署

**CentOS**

```bash
yum install -y nfs-utils
```

**Ubuntu**

```bash
sudo apt install nfs-kernel-server
```

## 配置并启动 NFS 服务

### 配置共享目录

服务启动之后，我们在服务端配置一个共享目录

```bash
mkdir /data
chmod 755 /data
```

为 NFS 配置共享目录

```bash
cat > /etc/exports <<EOF
/data/     172.19.42.0/24(rw,sync,no_root_squash,no_all_squash)
EOF
```

- /data # 共享目录位置。
- 172.19.42.0/24 # 客户端 IP 范围，\* 代表所有，即没有限制。
- rw # 权限设置，可读可写。
- sync # 同步共享目录。
- no_root_squash # 可以使用 root 授权。
- no_all_squash # 可以使用普通用户授权。

### 启动 nfs

**CentOS**

```bash
systemctl enable nfs --now
```

**Ubuntu**

```bash
systemctl enable nfs-server --now
```

### 检查一下本地的共享目录

```bash
root@common-test:~# showmount -e localhost
Export list for localhost:
/data 172.19.42.0/24
```

# 客户端部署

## 安装 NFS Client

### 通过 Linux 的包管理器部署

**CetnOS**

```bash
yum install -y nfs-utils
```

**Ubuntu**

```bash
sudo apt install nfs-common
```

## 客户端连接 NFS 测试

先查服务端的共享目录

    [root@node-1 ~]# showmount -e 172.19.42.215
    Export list for 172.19.42.215:
    /data 172.19.42.0/24

在客户端创建目录

```bash
mkdir /data_client
```

挂载

```bash
 mount -t nfs 172.19.42.247:/data /data/nfs_client
```

挂载之后，可以使用 mount 命令查看一下

```bash
$ mount
172.19.42.247:/data on /data/nfs_client type nfs4 (rw,relatime,vers=4.2,rsize=524288,wsize=524288,namlen=255,hard,proto=tcp,timeo=600,retrans=2,sec=sys,clientaddr=172.19.42.248,local_lock=none,addr=172.19.42.247)
```

这说明已经挂载成功了。

测试一下，在客户端向共享目录创建一个文件

```bash
cd /data
touch a
```

之后去 NFS 服务端查看一下

```bash
cd /data
ll
total 0-rw-r--r--. 1 root root 0 Aug  8 18:46 a
```

可以看到，共享目录已经写入了。

自动挂载很常用，客户端设置一下即可。

```bash
vi /etc/fstab
```

在结尾添加类似如下配置

```bash
# /etc/fstab: static file system information.
#
# Use 'blkid' to print the universally unique identifier for a
# device; this may be used with UUID= as a more robust way to name devices
# that works even if disks are added and removed. See fstab(5).
#
# <file system> <mount point>   <type>  <options>       <dump>  <pass>
#/dev/disk/by-id/dm-uuid-LVM-cmpKu0oh3CgDAAKsHGzU9AoEaWd9j8meFhd8PLCvGucV6WUB3F5au6gLhIOy0oYc none swap sw 0 0
# / was on /dev/vg0/lv-1 during curtin installation
/dev/disk/by-id/dm-uuid-LVM-cmpKu0oh3CgDAAKsHGzU9AoEaWd9j8meWCfiPw8kOhZvahcm3RCAT3sjPa0ialPd / xfs defaults 0 0
#/swap.img      none    swap    sw      0       0
172.19.42.247:/data      /data                   nfs     defaults        0 0
```
