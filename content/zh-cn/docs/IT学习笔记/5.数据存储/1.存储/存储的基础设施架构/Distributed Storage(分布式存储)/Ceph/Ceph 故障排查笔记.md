---
title: Ceph 故障排查笔记
---

参考：[公众号,云原声实验室-Ceph 故障排查笔记 | 万字经验总结](https://mp.weixin.qq.com/s/k5-gkm78KXmty7F2qerZ6g)

### Ceph OSD 异常无法正常启动

当某个 OSD 无法正常启动时：

    $ ceph -s
      cluster:
        id:     b313ec26-5aa0-4db2-9fb5-a38b207471ee
        health: HEALTH_WARN
                Degraded data redundancy: 177597/532791 objects degraded (33.333%), 212 pgs degraded, 212 pgs undersized
                application not enabled on 3 pool(s)
                mon master003 is low on available space
                1/3 mons down, quorum master002,master003
      services:
        mon: 3 daemons, quorum master002,master003, out of quorum: master001
        mgr: master003(active), standbys: master002
        mds: kubernetes-1/1/1 up  {0=master002=up:active}, 1 up:standby
        osd: 2 osds: 2 up, 2 in
      data:
        pools:   5 pools, 212 pgs
        objects: 177.6 k objects, 141 GiB
        usage:   297 GiB used, 2.8 TiB / 3.0 TiB avail
        pgs:     177597/532791 objects degraded (33.333%)
                 212 active+undersized+degraded
      io:
        client:   170 B/s rd, 127 KiB/s wr, 0 op/s rd, 5 op/s wr

查看状态信息：

    $ ceph health detail
    HEALTH_WARN Degraded data redundancy: 177615/532845 objects degraded (33.333%), 212 pgs degraded, 212 pgs undersized; application not enabled on 3 pool(s); mon master003 is low on available space
    PG_DEGRADED Degraded data redundancy: 177615/532845 objects degraded (33.333%), 212 pgs degraded, 212 pgs undersized
        pg 1.15 is active+undersized+degraded, acting [1,2]
        pg 1.2e is stuck undersized for 12701595.129535, current state active+undersized+degraded, last acting [1,2]
        pg 1.2f is stuck undersized for 12701595.110228, current state active+undersized+degraded, last acting [2,1]
        pg 1.30 is stuck undersized for 12701595.128371, current state active+undersized+degraded, last acting [1,2]
        pg 1.31 is stuck undersized for 12701595.129981, current state active+undersized+degraded, last acting [1,2]
        pg 1.32 is stuck undersized for 12701595.122298, current state active+undersized+degraded, last acting [2,1]
        pg 1.33 is stuck undersized for 12701595.129509, current state active+undersized+degraded, last acting [2,1]
        pg 1.34 is stuck undersized for 12701595.116494, current state active+undersized+degraded, last acting [2,1]
        pg 1.35 is stuck undersized for 12701595.132276, current state active+undersized+degraded, last acting [2,1]
        pg 1.36 is stuck undersized for 12701595.131601, current state active+undersized+degraded, last acting [1,2]
        pg 1.37 is stuck undersized for 12701595.126213, current state active+undersized+degraded, last acting [1,2]
        pg 1.38 is stuck undersized for 12701595.119082, current state active+undersized+degraded, last acting [2,1]
        pg 1.39 is stuck undersized for 12701595.127812, current state active+undersized+degraded, last acting [1,2]
        pg 1.3a is stuck undersized for 12701595.117611, current state active+undersized+degraded, last acting [2,1]
        pg 1.3b is stuck undersized for 12701595.125454, current state active+undersized+degraded, last acting [2,1]
        pg 1.3c is stuck undersized for 12701595.131540, current state active+undersized+degraded, last acting [1,2]
        pg 1.3d is stuck undersized for 12701595.130465, current state active+undersized+degraded, last acting [1,2]
        pg 1.3e is stuck undersized for 12701595.120532, current state active+undersized+degraded, last acting [2,1]
        pg 1.3f is stuck undersized for 12701595.129921, current state active+undersized+degraded, last acting [1,2]
        pg 1.40 is stuck undersized for 12701595.115146, current state active+undersized+degraded, last acting [2,1]
        pg 1.41 is stuck undersized for 12701595.132582, current state active+undersized+degraded, last acting [1,2]
        pg 1.42 is stuck undersized for 12701595.122272, current state active+undersized+degraded, last acting [2,1]
        pg 1.43 is stuck undersized for 12701595.132359, current state active+undersized+degraded, last acting [1,2]
        pg 1.44 is stuck undersized for 12701595.129082, current state active+undersized+degraded, last acting [2,1]
        pg 1.45 is stuck undersized for 12701595.118952, current state active+undersized+degraded, last acting [2,1]
        pg 1.46 is stuck undersized for 12701595.129618, current state active+undersized+degraded, last acting [1,2]
        pg 1.47 is stuck undersized for 12701595.112277, current state active+undersized+degraded, last acting [2,1]
        pg 1.48 is stuck undersized for 12701595.131721, current state active+undersized+degraded, last acting [1,2]
        pg 1.49 is stuck undersized for 12701595.130365, current state active+undersized+degraded, last acting [1,2]
        pg 1.4a is stuck undersized for 12701595.126070, current state active+undersized+degraded, last acting [1,2]
        pg 1.4b is stuck undersized for 12701595.113785, current state active+undersized+degraded, last acting [2,1]
        pg 1.4c is stuck undersized for 12701595.129074, current state active+undersized+degraded, last acting [1,2]
        pg 1.4d is stuck undersized for 12701595.115487, current state active+undersized+degraded, last acting [2,1]
        pg 1.4e is stuck undersized for 12701595.131307, current state active+undersized+degraded, last acting [1,2]
        pg 1.4f is stuck undersized for 12701595.132162, current state active+undersized+degraded, last acting [2,1]
        pg 1.50 is stuck undersized for 12701595.129346, current state active+undersized+degraded, last acting [2,1]
        pg 1.51 is stuck undersized for 12701595.131897, current state active+undersized+degraded, last acting [1,2]
        pg 1.52 is stuck undersized for 12701595.126480, current state active+undersized+degraded, last acting [2,1]
        pg 1.53 is stuck undersized for 12701595.116500, current state active+undersized+degraded, last acting [2,1]
        pg 1.54 is stuck undersized for 12701595.122930, current state active+undersized+degraded, last acting [2,1]
        pg 1.55 is stuck undersized for 12701595.116566, current state active+undersized+degraded, last acting [2,1]
        pg 1.56 is stuck undersized for 12701595.130017, current state active+undersized+degraded, last acting [1,2]
        pg 1.57 is stuck undersized for 12701595.129217, current state active+undersized+degraded, last acting [1,2]
        pg 1.58 is stuck undersized for 12701595.124121, current state active+undersized+degraded, last acting [2,1]
        pg 1.59 is stuck undersized for 12701595.127802, current state active+undersized+degraded, last acting [1,2]
        pg 1.5a is stuck undersized for 12701595.131028, current state active+undersized+degraded, last acting [1,2]
        pg 1.5b is stuck undersized for 12701595.114646, current state active+undersized+degraded, last acting [2,1]
        pg 1.5c is stuck undersized for 12701595.109604, current state active+undersized+degraded, last acting [2,1]
        pg 1.5d is stuck undersized for 12701595.126384, current state active+undersized+degraded, last acting [2,1]
        pg 1.5e is stuck undersized for 12701595.129456, current state active+undersized+degraded, last acting [1,2]
        pg 1.5f is stuck undersized for 12701595.126573, current state active+undersized+degraded, last acting [2,1]
    POOL_APP_NOT_ENABLED application not enabled on 3 pool(s)
        application not enabled on pool 'nextcloud'
        application not enabled on pool 'gitlab-ops'
        application not enabled on pool 'kafka-ops'
        use 'ceph osd pool application enable <pool-name> <app-name>', where <app-name> is 'cephfs', 'rbd', 'rgw', or freeform for custom applications.
    MON_DISK_LOW mon master003 is low on available space
        mon.master003 has 22% avail

并且通过 log 也无法完全定位问题时，可以通过如下方式解决。

### 删除 osd 重新加载

删除当前的 osd 重新让其进行加载，此方式适合于异常重启后的操作。
首先删除这个 osd：

    $ ceph osd out osd.0
    $ systemctl stop ceph-osd@0
    $ ceph osd crush remove osd.0
    $ ceph auth del osd.0
    $ ceph osd rm 0

重新加载 osd：

    $ ceph osd create 0
    $ ceph auth add osd.0 osd 'allow *' mon 'allow rwx' -i /var/lib/ceph/osd/ceph-0/keyring
    $ ceph osd crush add 0 1.0 host=master001
    $ systemctl start ceph-osd@0

### 清除当前 osd 所有数据重新添加

删除当前 osd 的所有数据，并且重新加载 osd，此操作一定要保证有冗余可用的 osd，否则会造成整个 osd 数据损坏。
删除当前 osd：

    $ ceph osd out osd.0
    $ systemctl stop ceph-osd@0
    $ ceph osd crush remove osd.0
    $ ceph auth del osd.0
    $ ceph osd rm 0

卸载：

    $ umount -l /var/lib/ceph/osd/ceph-0

清空磁盘数据：

    $ wipefs -af /dev/mapper/VolGroup-lv_data1
    $ ceph-volume lvm zap /dev/mapper/VolGroup-lv_data1

重新添加 osd：

    $ ceph-deploy --overwrite-conf osd create master001 --data /dev/mapper/VolGroup-lv_data1

### 删除当前节点所有服务

删除当前节点的所有服务，让其重新加载数据：

    $ ceph-deploy purge master001
    $ ceph-deploy purgedata master001

创建数据目录：

    $ rm -rf /var/lib/ceph
    $ mkdir -p /var/lib/ceph
    $ mkdir -p /var/lib/ceph/osd/ceph-0
    $ chown ceph:ceph /var/lib/ceph

然后安装 ceph：

    $ ceph-deploy install master001

同步配置：

    $ ceph-deploy --overwrite-conf admin master001

添加 osd：

    $ ceph-deploy osd create master001 --data /dev/mapper/VolGroup-lv_data1

### 查看当前系统 ceph 服务状态

查看当前系统 ceph 服务状态

    $ systemctl list-units |grep ceph

### 重启当前系统 ceph 服务

重启当前系统 ceph 服务

    $ systemctl restart ceph*.service ceph*.target

### 初始化 ceph-volume

初始化 ceph-volume

    $ ceph-volume lvm activate --bluestore --all

### 修改 Client keyring 和修复

修改 Client keyring 和修复，首先通过 ceph 命令进行查看：
然后把内容复制到：

    $ cat /var/lib/ceph/osd/ceph-0/keyring
    [osd.0]
    key = AQCzhrpeLRK+MhAAbjAgSsE7O81Q+8h8OwA92A==

### Pool 开启 enabled

pool 的 enabled 开启：

    $ ceph -s
      cluster:
        id:     b313ec26-5aa0-4db2-9fb5-a38b207471ee
        health: HEALTH_WARN
                application not enabled on 3 pool(s)
    $ ceph health detail
    HEALTH_WARN application not enabled on 3 pool(s); mon master003 is low on available space
    POOL_APP_NOT_ENABLED application not enabled on 3 pool(s)
        application not enabled on pool 'nextcloud'
        application not enabled on pool 'gitlab-ops'
        application not enabled on pool 'kafka-ops'
        use 'ceph osd pool application enable <pool-name> <app-name>', where <app-name> is 'cephfs', 'rbd', 'rgw', or freeform for custom applications.
    MON_DISK_LOW mon master003 is low on available space
        mon.master003 has 24% avail

执行 enabled：

    $ ceph osd pool application enable nextcloud rbd
    $ ceph osd pool application enable gitlab-ops rbd
    $ ceph osd pool application enable kafka-ops rbd

### Rbd 无法删除

rbd 无法删除，错误如下：

    $ rbd rm nextcloud/mysql
    2020-05-13 16:27:46.155 7f024bfff700 -1 librbd::image::RemoveRequest: 0x557a7af027a0 check_image_watchers: image has watchers - not removing
    Removing image: 0% complete...failed.
    rbd: error: image still has watchers
    This means the image is still open or the client using it crashed. Try again after closing/unmapping it or waiting 30s for the crashed client to timeout.
    $ rbd info nextcloud/mysql
    rbd image 'mysql':
            size 40 GiB in 10240 objects
            order 22 (4 MiB objects)
            id: 17e006b8b4567
            block_name_prefix: rbd_data.17e006b8b4567
            format: 2
            features: layering
            op_features:
            flags:
            create_timestamp: Tue Oct 15 10:47:34 2019

查看当前 rbd 状态：

    $ rbd status nextcloud/mysql
    Watchers:
            watcher=10.100.21.95:0/115493307 client.67866 cookie=7

发现有节点正在挂载，登入到相应机器进行查看：

    $ rbd showmapped
    id pool       image                                                       snap device
    ...
    3  nextcloud  mysql                                                       -    /dev/rbd3

取消映射：

    $ rbd unmap nextcloud/mysql

重新执行删除操作即可：

    $ rbd rm nextcloud/mysql
    Removing image: 100% complete...done.

暴力解决方案，直接对其添加黑名单，忽略挂载节点：

    $ ceph osd blacklist add 10.100.21.95:0/115493307
    $ rbd rm nextcloud/mysql

### OSD 延迟

查看是否有 osd 延迟：

    $ ceph osd perf
    osd commit_latency(ms) apply_latency(ms)
      2                  0                 0
      1                  0                 0
      0                  0                 0

### 碎片整理

查看碎片：

    $ xfs_db -c frag -r /dev/mapper/VolGroup-lv_data1

整理碎片：

### 查看通电时长

查看磁盘通电时长：

    $ smartctl -A /dev/mapper/VolGroup-lv_data1

### 修改副本数量

修改副本数量：

    $ ceph osd pool set fs_data2 min_size 1
    $ ceph osd pool set fs_data2 size 2

### 添加 / 删除 pool

添加 / 删除 pool：

    $ ceph fs add_data_pool fs fs_data2
    $ ceph fs rm_data_pool fs fs_data2

### osd 数据均衡分布

osd 数据均衡分布：

    $ ceph balancer status
    $ ceph balancer on
    $ ceph balancer mode crush-compat

### mds 无法查询

mds 无法查询:

    $ ceph fs status
    Error EINVAL: Traceback (most recent call last):
      File "/usr/lib64/ceph/mgr/status/module.py", line 311, in handle_command
        return self.handle_fs_status(cmd)
      File "/usr/lib64/ceph/mgr/status/module.py", line 177, in handle_fs_status
        mds_versions[metadata.get('ceph_version', "unknown")].append(info['name'])
    AttributeError: 'NoneType' object has no attribute 'get'
    $ ceph mds metadata
    [
        {
            "name": "BJ-YZ-CEPH-94-54"
        },
        {
            "name": "BJ-YZ-CEPH-94-53",
            "addr": "10.100.94.53:6825/4233274463",
            "arch": "x86_64",
            "ceph_release": "mimic",
            "ceph_version": "ceph version 13.2.10 (564bdc4ae87418a232fc901524470e1a0f76d641) mimic (stable)",
            "ceph_version_short": "13.2.10",
            "cpu": "Intel(R) Xeon(R) CPU E5-2620 v4 @ 2.10GHz",
            "distro": "centos",
            "distro_description": "CentOS Linux 7 (Core)",
            "distro_version": "7",
            "hostname": "BJ-YZ-CEPH-94-53",
            "kernel_description": "#1 SMP Sat Dec 10 18:16:05 EST 2016",
            "kernel_version": "4.4.38-1.el7.elrepo.x86_64",
            "mem_swap_kb": "67108860",
            "mem_total_kb": "131914936",
            "os": "Linux"
        },
        {
            "name": "BJ-YZ-CEPH-94-52",
            "addr": "10.100.94.52:6800/3956121270",
            "arch": "x86_64",
            "ceph_release": "mimic",
            "ceph_version": "ceph version 13.2.10 (564bdc4ae87418a232fc901524470e1a0f76d641) mimic (stable)",
            "ceph_version_short": "13.2.10",
            "cpu": "Intel(R) Xeon(R) CPU E5-2620 v4 @ 2.10GHz",
            "distro": "centos",
            "distro_description": "CentOS Linux 7 (Core)",
            "distro_version": "7",
            "hostname": "BJ-YZ-CEPH-94-52",
            "kernel_description": "#1 SMP Sat Dec 10 18:16:05 EST 2016",
            "kernel_version": "4.4.38-1.el7.elrepo.x86_64",
            "mem_swap_kb": "67108860",
            "mem_total_kb": "131914936",
            "os": "Linux"
        }
    ]

重启 mds 解决。

### cephfs 显示状态正常但无法写入数据

cephfs 显示正常无法使用，一般是有异常 client 导致的，首先查找 mds 是否存在链接，尝试删除链接解决：

    $ ceph tell mds.BJ-YZ-CEPH-94-52 session ls
    $ ceph tell mds.BJ-YZ-CEPH-94-52 session evict id=834283

每一个 mds 的 id 号不通用，不能跨节点删除。

### fs 增加 mds

fs 增加 mds:

    $ ceph fs set fs max_mds 2

### mon 时区异常

mon 因为时区有部分异常导致报错如下：

    $ ceph -s
      cluster:
        id:     2f77b028-ed2a-4010-9b79-90fd3052afc6
        health: HEALTH_WARN
                9 slow ops, oldest one blocked for 211643 sec, daemons [mon.BJ-YZ-CEPH-94-53,mon.BJ-YZ-CEPH-94-54] have slow ops.
      services:
        mon: 3 daemons, quorum BJ-YZ-CEPH-94-52,BJ-YZ-CEPH-94-53,BJ-YZ-CEPH-94-54
        mgr: BJ-YZ-CEPH-94-52(active), standbys: BJ-YZ-CEPH-94-54, BJ-YZ-CEPH-94-53
        mds: fs-2/2/2 up  {0=BJ-YZ-CEPH-94-52=up:active,1=BJ-YZ-CEPH-94-53=up:active}, 1 up:standby-replay
        osd: 36 osds: 36 up, 36 in
      data:
        pools:   7 pools, 1152 pgs
        objects: 37.66 M objects, 67 TiB
        usage:   136 TiB used, 126 TiB / 262 TiB avail
        pgs:     1148 active+clean
                 4    active+clean+scrubbing+deep
      io:
        client:   13 KiB/s rd, 27 MiB/s wr, 2 op/s rd, 19 op/s wr

配置 npt sever：

    $ systemctl status ntpd
    $ systemctl start ntpd

重启异常的 mon.targe 解决：

    $ systemctl status ceph-mon.target
    $ systemctl restart ceph-mon.target

### 1 MDSs report slow requests

报错如下：

    $ ceph -s
      cluster:
        id:     b313ec26-5aa0-4db2-9fb5-a38b207471ee
        health: HEALTH_WARN
                1 MDSs report slow requests
                Reduced data availability: 38 pgs inactive
                Degraded data redundancy: 122006/1192166 objects degraded (10.234%), 102 pgs degraded, 116 pgs undersized
                101 slow ops, oldest one blocked for 81045 sec, daemons [osd.1,osd.2] have slow ops.

重启 mon 即可解决：

    $ systemctl restart ceph-mon.target

如果无法解决需要重启 mds 解决：

    $ systemctl restart ceph-mds@${HOSTNAME}

### Reduced data availability: 38 pgs inactive

报错如下：**<https://zhuanlan.zhihu.com/p/74323736>**

    $ ceph -s
      cluster:
        id:     b313ec26-5aa0-4db2-9fb5-a38b207471ee
        health: HEALTH_WARN
                1 MDSs report slow requests
                Reduced data availability: 38 pgs inactive
                145 slow ops, oldest one blocked for 184238 sec, daemons [osd.1,osd.2] have slow ops.
      services:
        mon: 3 daemons, quorum master001,master002,master003
        mgr: master001(active), standbys: master002, master003
        mds: kubernetes-2/2/2 up  {0=master001=up:active,1=master002=up:active}, 1 up:standby
        osd: 3 osds: 3 up, 3 in
        rgw: 1 daemon active
      data:
        pools:   9 pools, 244 pgs
        objects: 535.1 k objects, 177 GiB
        usage:   470 GiB used, 4.1 TiB / 4.6 TiB avail
        pgs:     15.574% pgs unknown
                 206 active+clean
                 38  unknown
      io:
        client:   35 KiB/s wr, 0 op/s rd, 2 op/s wr

此问题属于 pg 丢失数据并且无法自动回复造成的。解决办法是清除 pg 数据让其自动修复，但这样可能会造成数据丢失（如果 size 为 1 则肯定丢失数据）
首先查看异常的 pg：
然后执行 query 查看信息：

    $ ceph pg 1.6e query
    Error ENOENT: i don't have pgid 1.6e

上述无法查到 pg，通过如下命令查看异常的 pg：

    $ ceph pg dump_stuck unclean
    ok
    PG_STAT STATE   UP UP_PRIMARY ACTING ACTING_PRIMARY
    1.74    unknown []         -1     []             -1
    1.70    unknown []         -1     []             -1
    1.6a    unknown []         -1     []             -1
    1.2d    unknown []         -1     []             -1
    1.20    unknown []         -1     []             -1
    1.1e    unknown []         -1     []             -1
    1.1c    unknown []         -1     []             -1
    1.17    unknown []         -1     []             -1
    1.9     unknown []         -1     []             -1
    1.29    unknown []         -1     []             -1
    1.56    unknown []         -1     []             -1
    1.72    unknown []         -1     []             -1
    1.45    unknown []         -1     []             -1
    1.4e    unknown []         -1     []             -1
    1.46    unknown []         -1     []             -1
    1.22    unknown []         -1     []             -1
    1.53    unknown []         -1     []             -1
    1.59    unknown []         -1     []             -1
    1.24    unknown []         -1     []             -1
    1.55    unknown []         -1     []             -1
    1.3f    unknown []         -1     []             -1
    1.38    unknown []         -1     []             -1
    1.a     unknown []         -1     []             -1
    1.7     unknown []         -1     []             -1
    1.34    unknown []         -1     []             -1
    1.64    unknown []         -1     []             -1
    1.6     unknown []         -1     []             -1
    1.32    unknown []         -1     []             -1
    1.4     unknown []         -1     []             -1
    1.2e    unknown []         -1     []             -1
    1.31    unknown []         -1     []             -1
    1.5e    unknown []         -1     []             -1
    1.0     unknown []         -1     []             -1
    1.42    unknown []         -1     []             -1
    1.15    unknown []         -1     []             -1
    1.6e    unknown []         -1     []             -1
    1.41    unknown []         -1     []             -1
    1.10    unknown []         -1     []             -1

执行如下命令强制清除 pg 的数据：**<https://docs.ceph.com/docs/mimic/rados/troubleshooting/troubleshooting-pg/>**

    $ ceph osd force-create-pg 1.74 --yes-i-really-mean-it
    # 批量执行
    # ceph pg dump_stuck unclean|awk '{print $1}'|xargs -i ceph osd force-create-pg {} --yes-i-really-mean-it

执行完成后即可恢复。

### 1 clients failing to respond to capability release

报错如下：

    $ ceph  health detail
    HEALTH_WARN 1 clients failing to respond to capability release
    MDS_CLIENT_LATE_RELEASE 1 clients failing to respond to capability release
        mdsmaster001(mds.0): Client master003.k8s.shileizcc-ops.com: failing to respond to capability release client_id: 284951

清除次 ID 即可：**<https://blog.csdn.net/zuoyang1990/article/details/98530070>**

    $ ceph daemon mds.master003 session ls|grep 284951
    $ ceph tell mds.master003 session evict id=284951

如果报错如下：

    $ ceph tell mds.master003 session evict id=284951
    2020-08-13 10:45:03.869 7f271b7fe700  0 client.306366 ms_handle_reset on 10.100.21.95:6800/1646216103
    2020-08-13 10:45:03.881 7f2730ff9700  0 client.316415 ms_handle_reset on 10.100.21.95:6800/1646216103
    Error EAGAIN: MDS is replaying log

需要到 mds.0 节点执行，否则无法找到次 client。

### 内核优化

内核优化：**<https://blog.csdn.net/fuzhongfaya/article/details/80932766>**

    $ echo "8192" > /sys/block/sda/queue/read_ahead_kb
    $ echo "vm.swappiness = 0" | tee -a /etc/sysctl.conf
    $ sysctl -p
    $ echo "deadline" > /sys/block/sd[x]/queue/scheduler
    # ssd
    # echo "noop" > /sys/block/sd[x]/queue/scheduler

swap 最好是直接关闭，配置内存参数在一定程度上不会生效。
配置文件
40 核心 128 GB 配置文件：

    [global]
    fsid = 2f77b028-ed2a-4010-9b79-90fd3052afc6
    mon_initial_members = BJ-YZ-CEPH-94-52, BJ-YZ-CEPH-94-53, BJ-YZ-CEPH-94-54
    mon_host = 10.100.94.52,10.100.94.53,10.100.94.54
    auth_cluster_required = cephx
    auth_service_required = cephx
    auth_client_required = cephx
    public network = 10.100.94.0/24
    cluster network = 10.100.94.0/24
    [mon.a]
    host = BJ-YZ-CEPH-94-52
    mon addr = 10.100.94.52:6789
    [mon.b]
    host = BJ-YZ-CEPH-94-53
    mon addr = 10.100.94.53:6789
    [mon.c]
    host = BJ-YZ-CEPH-94-54
    mon addr = 10.100.94.54:6789
    [mon]
    mon data = /var/lib/ceph/mon/ceph-$id
    # monitor 间的 clock drift，默认值 0.05
    mon clock drift allowed = 1
    # 向 monitor 报告 down 的最小 OSD 数，默认值 1
    mon osd min down reporters = 1
    # 标记一个OSD状态为down和out之前ceph等待的秒数，默认值300
    mon osd down out interval = 600
    mon_allow_pool_delete = true
    [osd]
    # osd 数据路径
    osd data = /var/lib/ceph/osd/ceph-$id
    # 默认 pool pg,pgp 数量
    osd pool default pg num =  1200
    osd pool default pgp num = 1200
    # osd 的 journal 写日志时的大小默认 5120
    osd journal size = 20000
    # 格式化文件系统类型
    osd mkfs type = xfs
    # 格式化文件系统时附加参数
    osd mkfs options xfs = -f
    # 为 XATTRS 使用 object map，EXT4 文件系统时使用，XFS 或者 btrf 也可以使用，默认 false
    filestore xattr use omap = true
    # 从日志到数据盘最小同步间隔(seconds)，默认值 0.1
    filestore min sync interval = 10
    # 从日志到数据盘最大同步间隔(seconds)，默认值 5
    filestore max sync interval = 15
    # 数据盘最大接受的操作数，默认值 500
    filestore queue max ops = 25000
    # 数据盘能够 commit 的最大字节数(bytes)，默认值 100
    filestore queue max bytes = 10485760
    # 数据盘能够 commit 的操作数，500
    filestore queue committing max ops = 5000
    # 数据盘能够 commit 的最大字节数(bytes)，默认值 100
    filestore queue committing max bytes = 10485760000
    # 前一个子目录分裂成子目录中的文件的最大数量，默认值 2
    filestore split multiple = 8
    # 前一个子类目录中的文件合并到父类的最小数量，默认值10
    filestore merge threshold = 40
    # 对象文件句柄缓存大小，默认值 128
    filestore fd cache size = 1024
    # 并发文件系统操作数，默认值 2
    filestore op threads = 32
    # journal 一次性写入的最大字节数(bytes)，默认值 1048560
    journal max write bytes = 1073714824
    # journal一次性写入的最大记录数，默认值 100
    journal max write entries = 10000
    # journal一次性最大在队列中的操作数，默认值 50
    journal queue max ops = 50000
    # journal一次性最大在队列中的字节数(bytes)，默认值 33554432
    journal queue max bytes = 10485760000
    # # OSD一次可写入的最大值(MB), 默认 90
    osd max write size = 512
    # 客户端允许在内存中的最大数据(bytes), 默认值100
    osd client message size cap = 2147483648
    # 在 Deep Scrub 时候允许读取的字节数(bytes), 默认值524288
    osd deep scrub stride = 1310720
    # 并发文件系统操作数, 默认值 2
    osd op threads = 32
    # OSD 密集型操作例如恢复和 Scrubbing 时的线程, 默认值1
    osd disk threads = 10
    # 保留 OSD Map 的缓存(MB), 默认 500
    osd map cache size = 10240
    # OSD 进程在内存中的 OSD Map 缓存(MB), 默认 50
    osd map cache bl size = 1280
    # 默认值rw,noatime,inode64, Ceph OSD xfs Mount选项
    osd mount options xfs = "rw,noexec,nodev,noatime,nodiratime,nobarrier"
    # 恢复操作优先级，取值 1-63，值越高占用资源越高, 默认值 10
    osd recovery op priority = 20
    # 同一时间内活跃的恢复请求数, 默认值 15
    osd recovery max active = 15
    # 一个 OSD 允许的最大 backfills 数, 默认值 10
    osd max backfills = 10
    # 开启严格队列降级操作
    osd op queue cut off = high
    osd_deep_scrub_large_omap_object_key_threshold = 800000
    osd_deep_scrub_large_omap_object_value_sum_threshold = 10737418240
    [mds]
    # mds 缓存大小设置 60GB
    mds cache memory limit = 62212254726
    # 超时时间默认 60 秒
    mds_revoke_cap_timeout = 360
    mds log max segments = 51200
    mds log max expiring = 51200
    mds_beacon_grace = 300
    # 对目录碎片大小的硬限制 默认 100000
    # https://docs.ceph.com/docs/master/cephfs/dirfrags/
    mds_bal_fragment_size_max = 500000
    ## 官方配置 https://ceph.readthedocs.io/en/latest/cephfs/mds-config-ref/
    [client]
    # RBD缓存, 默认 true
    rbd cache = true
    # RBD缓存大小(bytes), 默认 335544320（320M）
    rbd cache size = 268435456
    # 缓存为 write-back 时允许的最大 dirty 字节数(bytes)，如果为0，使用 write-through，默认值为 25165824
    rbd cache max dirty = 134217728
    # 在被刷新到存储盘前 dirty 数据存在缓存的时间(seconds), 默认值为 1
    rbd cache max dirty age = 5
    client_try_dentry_invalidate = false
    [mgr]
    mgr modules = dashboard
    # 华为云调优指南 https://support.huaweicloud.com/tngg-kunpengsdss/kunpengcephobject_05_0008.html
    # https://poph163.com/2020/02/18/ceph-crushmap%E4%B8%8E%E8%B0%83%E4%BC%98/

### full osd

full osd 每个 osd 已经写满上限:**<https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-osd/#no-free-drive-space>**

    $ ceph osd dump | grep full_ratio
    full_ratio 0.95
    backfillfull_ratio 0.9
    nearfull_ratio 0.85

集群状态:

    $ ceph -s
      cluster:
        id:     2f77b028-ed2a-4010-9b79-90fd3052afc6
        health: HEALTH_ERR
                2 backfillfull osd(s)
                1 full osd(s)
                2 nearfull osd(s)
                7 pool(s) full

执行 osd 磁盘状态时，如果已经有超过 95% 使用率时则会报错 full osd 则会造成 cluster 无法正常使用：

    $ ceph osd df
    ID CLASS WEIGHT  REWEIGHT SIZE    USE     DATA    OMAP     META    AVAIL   %USE  VAR  PGS
     0   hdd 7.27689  1.00000 7.3 TiB 4.7 TiB 4.7 TiB  918 MiB 9.1 GiB 2.5 TiB 65.15 0.84  68
     1   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  327 MiB  11 GiB 1.2 TiB 84.07 1.09  67
     2   hdd 7.27689  1.00000 7.3 TiB 4.3 TiB 4.3 TiB  924 MiB 8.4 GiB 2.9 TiB 59.70 0.77  67
     3   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB  807 MiB 9.8 GiB 2.1 TiB 70.57 0.91  66
     4   hdd 7.27689  1.00000 7.3 TiB 6.7 TiB 6.7 TiB  770 MiB  13 GiB 583 GiB 92.18 1.19  66
     5   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB  623 MiB  10 GiB 1.8 TiB 75.87 0.98  66
     6   hdd 7.27689  1.00000 7.3 TiB 5.7 TiB 5.7 TiB  602 MiB  11 GiB 1.6 TiB 78.67 1.02  64
     7   hdd 7.27689  1.00000 7.3 TiB 5.3 TiB 5.3 TiB  1.1 GiB  10 GiB 1.9 TiB 73.35 0.95  65
     8   hdd 7.27689  1.00000 7.3 TiB 5.9 TiB 5.9 TiB  498 MiB  11 GiB 1.4 TiB 81.29 1.05  68
     9   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB  1.1 GiB 9.8 GiB 2.1 TiB 70.59 0.91  65
    10   hdd 7.27689  1.00000 7.3 TiB 6.3 TiB 6.3 TiB  297 MiB  12 GiB 985 GiB 86.78 1.12  61
    11   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB  923 MiB 9.7 GiB 2.1 TiB 70.56 0.91  67
    12   hdd 7.27689  1.00000 7.3 TiB 5.9 TiB 5.9 TiB  203 MiB  11 GiB 1.4 TiB 81.39 1.05  65
    13   hdd 7.27689  1.00000 7.3 TiB 5.3 TiB 5.3 TiB  799 MiB  10 GiB 1.9 TiB 73.29 0.95  66
    14   hdd 7.27689  1.00000 7.3 TiB 4.9 TiB 4.9 TiB  873 MiB 9.4 GiB 2.3 TiB 67.77 0.88  71
    15   hdd 0.29999  1.00000 7.3 TiB 6.9 TiB 6.9 TiB  191 MiB  13 GiB 387 GiB 94.81 1.23  39
    16   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB  548 MiB  11 GiB 1.8 TiB 75.91 0.98  69
    17   hdd 7.27689  1.00000 7.3 TiB 6.7 TiB 6.7 TiB  806 MiB  13 GiB 581 GiB 92.20 1.20  66
    18   hdd 7.27689  1.00000 7.3 TiB 4.5 TiB 4.5 TiB  1.4 GiB 8.5 GiB 2.7 TiB 62.43 0.81  66
    19   hdd 7.27689  1.00000 7.3 TiB 5.3 TiB 5.3 TiB  1.4 GiB  10 GiB 1.9 TiB 73.28 0.95  65
    20   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB  705 MiB  11 GiB 1.8 TiB 75.91 0.98  64
    21   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  911 MiB  11 GiB 1.2 TiB 84.11 1.09  62
    22   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  301 MiB  11 GiB 1.2 TiB 84.03 1.09  66
    23   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB  401 MiB 9.8 GiB 1.7 TiB 75.96 0.98  67
    24   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB  1.3 GiB 9.6 GiB 2.1 TiB 70.58 0.91  63
    25   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB  1.1 GiB 9.7 GiB 2.1 TiB 70.56 0.91  65
    26   hdd 7.27689  1.00000 7.3 TiB 5.3 TiB 5.3 TiB  730 MiB  10 GiB 1.9 TiB 73.32 0.95  68
    27   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  818 MiB  12 GiB 1.2 TiB 84.08 1.09  62
    28   hdd 7.27689  1.00000 7.3 TiB 4.9 TiB 4.9 TiB  587 MiB 9.3 GiB 2.3 TiB 67.84 0.88  68
    29   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  215 MiB  11 GiB 1.2 TiB 84.09 1.09  66
    30   hdd 7.27689  1.00000 7.3 TiB 6.1 TiB 6.1 TiB  690 MiB  12 GiB 1.2 TiB 84.15 1.09  64
    31   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB 1020 MiB  10 GiB 1.8 TiB 75.94 0.98  64
    32   hdd 7.27689  1.00000 7.3 TiB 6.5 TiB 6.5 TiB  616 MiB  12 GiB 786 GiB 89.45 1.16  66
    33   hdd 7.27689  1.00000 7.3 TiB 4.9 TiB 4.9 TiB  622 MiB 8.9 GiB 2.3 TiB 67.84 0.88  66
    34   hdd 7.27689  1.00000 7.3 TiB 5.7 TiB 5.7 TiB  102 MiB  11 GiB 1.6 TiB 78.56 1.02  65
    35   hdd 7.27689  1.00000 7.3 TiB 5.9 TiB 5.9 TiB  723 MiB  11 GiB 1.4 TiB 81.31 1.05  63
                        TOTAL 262 TiB 202 TiB 202 TiB   25 GiB 381 GiB  60 TiB 77.15

可以手动修改权重解决:

    $ ceph osd crush reweight osd.4 0.3

### pg 均衡

pg 在默认分配有不合理的地方。**<https://cloud.tencent.com/developer/article/1664655>**

    $ ceph osd df tree | awk '/osd\./{print $NF" "$(NF-1)" "$(NF-3) }'
    osd.0 89 71.20
    osd.1 38 94.80
    osd.2 92 68.44
    osd.3 92 72.36
    osd.4 28 76.86
    osd.5 64 81.37
    osd.6 62 87.90
    osd.7 89 78.78
    osd.8 52 86.18
    osd.9 89 75.44
    osd.10 37 96.33
    osd.11 102 75.26
    osd.12 33 91.41
    osd.13 34 95.98
    osd.14 59 84.97
    osd.15 20 70.92
    osd.16 113 89.46
    osd.17 30 77.12
    osd.18 124 77.11
    osd.19 44 95.23
    osd.20 65 84.63
    osd.21 98 96.71
    osd.22 34 95.93
    osd.23 62 84.56
    osd.24 110 76.63
    osd.25 64 82.32
    osd.26 59 88.26
    osd.27 38 95.83
    osd.28 105 79.19
    osd.29 36 94.94
    osd.30 94 90.79
    osd.31 91 81.74
    osd.32 12 42.44
    osd.33 94 81.32
    osd.34 46 86.51
    osd.35 37 92.68

reweight-by-pg 按归置组分布情况调整 OSD 的权重:

    $ ceph osd reweight-by-pg
    moved 0 / 2336 (0%)
    avg 64.8889
    stddev 58.677 -> 58.677 (expected baseline 7.9427)
    min osd.1 with 0 -> 0 pgs (0 -> 0 * mean)
    max osd.18 with 168 -> 168 pgs (2.58904 -> 2.58904 * mean)
    oload 120
    max_change 0.05
    max_change_osds 4
    average_utilization 18.2677
    overload_utilization 21.9212
    osd.19 weight 1.0000 -> 0.9500
    osd.1 weight 1.0000 -> 0.9500
    osd.27 weight 1.0000 -> 0.9500
    osd.10 weight 1.0000 -> 0.9500

reweight-by-utilization 按利用率调整 OSD 的权重:

    $ ceph osd reweight-by-pg
    moved 0 / 2336 (0%)
    avg 64.8889
    stddev 58.677 -> 58.677 (expected baseline 7.9427)
    min osd.1 with 0 -> 0 pgs (0 -> 0 * mean)
    max osd.18 with 168 -> 168 pgs (2.58904 -> 2.58904 * mean)
    oload 120
    max_change 0.05
    max_change_osds 4
    average_utilization 18.2677
    overload_utilization 21.9212
    osd.19 weight 1.0000 -> 0.9500
    osd.1 weight 1.0000 -> 0.9500
    osd.27 weight 1.0000 -> 0.9500
    osd.10 weight 1.0000 -> 0.9500

调整写入权重：

    $ ceph osd reweight osd.35 0.001

查看当前 osd 信息：

    $ ceph osd df
    ID CLASS WEIGHT  REWEIGHT SIZE    USE     DATA    OMAP    META    AVAIL   %USE  VAR  PGS
     0   hdd 7.27689  1.00000 7.3 TiB 5.2 TiB 5.2 TiB 1.0 GiB 9.4 GiB 2.0 TiB 71.96 0.86  39
     1   hdd 0.00999  0.90002 7.3 TiB 6.9 TiB 6.9 TiB 604 MiB  12 GiB 382 GiB 94.88 1.13  37
     2   hdd 7.27689  1.00000 7.3 TiB 5.1 TiB 5.1 TiB 1.2 GiB 8.8 GiB 2.2 TiB 69.55 0.83  34
     3   hdd 7.27689  1.00000 7.3 TiB 5.3 TiB 5.3 TiB 812 MiB 9.9 GiB 2.0 TiB 73.15 0.87  34
     4   hdd 0.29999  1.00000 7.3 TiB 5.6 TiB 5.6 TiB 185 MiB  12 GiB 1.7 TiB 77.01 0.92  26
     5   hdd 3.00000  1.00000 7.3 TiB 6.0 TiB 5.9 TiB 443 MiB  11 GiB 1.3 TiB 81.90 0.98  36
     6   hdd 3.00000  1.00000 7.3 TiB 6.5 TiB 6.5 TiB 499 MiB  11 GiB 809 GiB 89.14 1.06  38
     7   hdd 7.27689  1.00000 7.3 TiB 5.8 TiB 5.8 TiB 1.2 GiB  11 GiB 1.4 TiB 80.10 0.96  43
     8   hdd 3.00000  1.00000 7.3 TiB 6.3 TiB 6.3 TiB 502 MiB  11 GiB 992 GiB 86.69 1.03  36
     9   hdd 7.27689  1.00000 7.3 TiB 5.6 TiB 5.6 TiB 1.5 GiB 9.8 GiB 1.7 TiB 76.57 0.91  42
    10   hdd 0.00999  0.00099 7.3 TiB 7.0 TiB 7.0 TiB 295 MiB  12 GiB 267 GiB 96.41 1.15  37
    11   hdd 7.27689  1.00000 7.3 TiB 5.5 TiB 5.5 TiB 1.2 GiB 9.8 GiB 1.7 TiB 76.13 0.91  37
    12   hdd 0.00999  1.00000 7.3 TiB 6.7 TiB 6.6 TiB  95 MiB  12 GiB 635 GiB 91.48 1.09  32
    13   hdd 0.00999  1.00000 7.3 TiB 7.0 TiB 7.0 TiB 584 MiB  12 GiB 315 GiB 95.78 1.14  34
    14   hdd 3.00000  1.00000 7.3 TiB 6.2 TiB 6.2 TiB 974 MiB  11 GiB 1.0 TiB 85.86 1.02  40
    15   hdd 0.00999  1.00000 7.3 TiB 5.1 TiB 5.1 TiB 116 KiB  10 GiB 2.2 TiB 70.43 0.84  20
    16   hdd 7.27689  1.00000 7.3 TiB 6.6 TiB 6.6 TiB 1.2 GiB  11 GiB 697 GiB 90.64 1.08  43
    17   hdd 0.29999  1.00000 7.3 TiB 5.6 TiB 5.6 TiB  40 KiB  12 GiB 1.7 TiB 76.75 0.92  26
    18   hdd 7.27689  1.00000 7.3 TiB 5.7 TiB 5.7 TiB 1.9 GiB 9.3 GiB 1.6 TiB 78.01 0.93  53
    19   hdd 0.00999  0.00099 7.3 TiB 6.9 TiB 6.9 TiB 1.5 GiB  13 GiB 371 GiB 95.02 1.13  40
    20   hdd 3.00000  1.00000 7.3 TiB 6.2 TiB 6.2 TiB 744 MiB  12 GiB 1.0 TiB 85.86 1.02  37
    21   hdd 7.27689  0.00099 7.3 TiB 7.0 TiB 7.0 TiB 913 MiB  12 GiB 239 GiB 96.79 1.15  40
    22   hdd 0.00999  0.00099 7.3 TiB 7.0 TiB 7.0 TiB 283 MiB  12 GiB 298 GiB 96.00 1.14  34
    23   hdd 3.00000  1.00000 7.3 TiB 6.2 TiB 6.2 TiB 515 MiB  11 GiB 1.1 TiB 85.30 1.02  35
    24   hdd 7.27689  1.00000 7.3 TiB 5.6 TiB 5.6 TiB 1.4 GiB 9.8 GiB 1.6 TiB 77.63 0.93  42
    25   hdd 3.00000  1.00000 7.3 TiB 6.0 TiB 6.0 TiB 1.2 GiB  10 GiB 1.3 TiB 82.66 0.99  40
    26   hdd 2.00000  1.00000 7.3 TiB 6.5 TiB 6.5 TiB 737 MiB  11 GiB 823 GiB 88.95 1.06  36
    27   hdd 0.00999  0.00099 7.3 TiB 7.0 TiB 6.9 TiB 822 MiB  12 GiB 327 GiB 95.61 1.14  37
    28   hdd 7.27689  1.00000 7.3 TiB 5.8 TiB 5.8 TiB 859 MiB  10 GiB 1.4 TiB 80.23 0.96  40
    29   hdd 0.00999  0.00099 7.3 TiB 6.9 TiB 6.9 TiB 215 MiB  12 GiB 371 GiB 95.02 1.13  36
    30   hdd 7.27689  1.00000 7.3 TiB 6.7 TiB 6.7 TiB 1.0 GiB  12 GiB 607 GiB 91.85 1.10  47
    31   hdd 7.27689  1.00000 7.3 TiB 6.0 TiB 6.0 TiB 1.2 GiB  10 GiB 1.3 TiB 82.81 0.99  41
    32   hdd 0.29999  1.00000 7.3 TiB 3.0 TiB 3.0 TiB  32 KiB 7.1 GiB 4.3 TiB 41.47 0.49  10
    33   hdd 7.27689  1.00000 7.3 TiB 6.0 TiB 6.0 TiB 827 MiB 9.7 GiB 1.3 TiB 82.06 0.98  41
    34   hdd 2.00000  1.00000 7.3 TiB 6.3 TiB 6.3 TiB 308 MiB  11 GiB 976 GiB 86.90 1.04  33
    35   hdd 0.00999  0.00099 7.3 TiB 6.7 TiB 6.7 TiB 613 MiB  12 GiB 540 GiB 92.75 1.11  36
                        TOTAL 262 TiB 220 TiB 219 TiB  27 GiB 391 GiB  42 TiB 83.87
    MIN/MAX VAR: 0.49/1.15  STDDEV: 10.62

删除 Cephfs
关闭所有 mds 服务, 需要登入服务器手动关闭:

    $ systemctl stop ceph-mds@${HOSTNAME}

删除所需 fs:

    $ ceph fs ls
    $ ceph fs rm data --yes-i-really-mean-it

SSD 使用
查看当前 OSD 状态: (相关文档:**<https://blog.csdn.net/kozazyh/article/details/79904219>**)

    $ ceph osd crush class ls
    [
        "ssd"
    ]

如果使用的 SSD 标识错误，请自定义修改，命令如下, 移除 osd 1 ~ 3 的标识:

    $ for i in 0 1 2;do ceph osd crush rm-device-class osd.$i;done

设置 1 ~ 3 标识为 ssd：

    $ for i in 0 1 2;do ceph osd crush set-device-class ssd osd.$i;done

创建一个 crush rule:

    $ ceph osd crush rule create-replicated rule-ssd default host ssd
    $ ceph osd crush rule ls

然后创建 pool 时附带 rule 的名称：

    $ ceph osd pool create fs_data 96 rule-ssd
    $ ceph osd pool create fs_metadata 16 rule-ssd
    $ ceph fs new fs fs_data fs_metadata

### crushmap 查看

执行命令如下:

    $ ceph osd getcrushmap -o crushmap
    $ crushtool -d crushmap -o crushmap
    $ cat crushmap

### 3 monitors have not enabled msgr2

解决如下：

    $ ceph mon enable-msgr2

### 2 daemons have recently crashed

解决如下：**<https://blog.csdn.net/QTM_Gitee/article/details/106004435>**

    $ ceph crash ls
    $ ceph crash archive-all

### 脚注

\[1]
<https://zhuanlan.zhihu.com/p/74323736:> _<https://zhuanlan.zhihu.com/p/74323736>_
\[2]
<https://docs.ceph.com/docs/mimic/rados/troubleshooting/troubleshooting-pg/:> _<https://docs.ceph.com/docs/mimic/rados/troubleshooting/troubleshooting-pg/>_
\[3]
<https://blog.csdn.net/zuoyang1990/article/details/98530070:> _<https://blog.csdn.net/zuoyang1990/article/details/98530070>_
\[4]
<https://blog.csdn.net/fuzhongfaya/article/details/80932766:> _<https://blog.csdn.net/fuzhongfaya/article/details/80932766>_
\[5]
<https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-osd/#no-free-drive-space:> _<https://docs.ceph.com/en/latest/rados/troubleshooting/troubleshooting-osd/#no-free-drive-space>_
\[6]
<https://cloud.tencent.com/developer/article/1664655:> _<https://cloud.tencent.com/developer/article/1664655>_
\[7]
<https://blog.csdn.net/kozazyh/article/details/79904219:> _<https://blog.csdn.net/kozazyh/article/details/79904219>_
\[8]
<https://blog.csdn.net/QTM_Gitee/article/details/106004435:> _<https://blog.csdn.net/QTM_Gitee/article/details/106004435>_
