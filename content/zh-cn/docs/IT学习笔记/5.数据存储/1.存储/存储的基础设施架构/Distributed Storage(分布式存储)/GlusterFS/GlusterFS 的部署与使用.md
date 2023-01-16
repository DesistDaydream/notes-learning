---
title: GlusterFS 的部署与使用
---

# GlusterFS 的部署与使用

初始化 yum 源配置

yum install centos-release-gluster6 -y

分区、格式化、挂载使用 bricks 的磁盘

Assuming you have a brick at /dev/sdb:

fdisk /dev/sdb

Create a single partition on the brick that uses the whole capacity.

格式化分区

mkfs.xfs -i size=512 /dev/sdb

将分区挂载为 gluster 的 brick

mkdir -p /data/brick1 && mount /dev/sdb /data/brick1 && mkdir -p /data/brick1brick

在/etc/fstab 文件中添加条目使得目录自动挂载

echo "/dev/sdb /data/brick1 xfs defaults 0 0" >> /etc/fstab

在所有节点安装 glusterfs 所用的包

yum install glusterfs{,-server,-fuse,-geo-replication,-client} -y

Note: This example assumes Fedora 18 or later, where gluster packages are included in the official repository

启动 glusterd 服务

systemctl start glusterd && systemctl status glusterd && systemctl enable glusterd

从一个节点添加另一个节点，使之组成集群

Note! From node01 to the other nodes (do not peer probe the first node)

gluster peer probe

创建一个 volume，指定卷名、类型、以及各节点所用目录

gluster volume create testvol rep 2 transport tcp node1:/data/brick1/brick node2:/data/brick1/brick

启动新添加的卷

gluster volume start testvol

在其他节点挂载 glusterfs 服务端上的卷到本地

mkdir /mnt/gluster

mount -t glusterfs node1:/testvol /mnt/gluster

挂载完成后，该挂载目录就是 guluster 服务端远程的卷目录，写入该目录的文件会同时写到 volume 所定义的所有节点下的目录

Done!
