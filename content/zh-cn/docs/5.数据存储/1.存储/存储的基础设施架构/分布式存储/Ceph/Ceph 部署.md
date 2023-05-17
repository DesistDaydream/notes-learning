---
title: Ceph 部署
---

# 概述

> 参考：
> - [官方文档, Cephadm](https://docs.ceph.com/en/latest/cephadm/)

# 以 pacific 版本为例

## 安装 cephadm

在所有节点安装 cephadm，这是一个 python 程序，当通过第一个节点让其他节点加入集群时，会调用待加入节点中的 cephadm 程序。

```bash
curl --silent --remote-name --location https://github.com/ceph/ceph/raw/pacific/src/cephadm/cephadm
chmod +x cephadm
./cephadm add-repo --release pacific
./cephadm install
cephadm install ceph-common
```

## 引导第一个节点

```bash
cephadm bootstrap --mon-ip 192.168.1.201
```

### 配置 ceph CLI

```bash
cephadm install ceph-common
```

### 其他

```bash
# 开启遥测，发送数据给官方
ceph telemetry on
```

## 添加其他节点

```bash
# 添加认证信息
ssh-copy-id -f -i /etc/ceph/ceph.pub root@192.168.1.202
ssh-copy-id -f -i /etc/ceph/ceph.pub root@192.168.1.203
# 添加节点
ceph orch host add hw-cloud-xngy-ecs-test-0002 192.168.1.202
ceph orch host add hw-cloud-xngy-ecs-test-0003 192.168.1.203
# 为节点添加 _admin 标签
ceph orch host label add hw-cloud-xngy-ecs-test-0002 _admin
ceph orch host label add hw-cloud-xngy-ecs-test-0003 _admin
```

节点添加完成后，在 1 和 2 上活动 ceph-mgr，1，2，3 上启动了 ceph-mon 和 ceph-crash

## 添加存储设备

注意：如下所示，当一块磁盘具有 GPT 分区表时，是无法作为 Ceph 的 OSD 使用

```bash
root@hw-cloud-xngy-ecs-test-0001:~# ceph orch device ls --wide
Hostname                     Path      Type  Transport  RPM      Vendor  Model  Serial                Size   Health   Ident  Fault  Available  Reject Reasons
hw-cloud-xngy-ecs-test-0001  /dev/vdb  hdd   Unknown    Unknown  0x1af4  N/A    4afb2ab1-9244-45bf-a   107G  Unknown  N/A    N/A    No         Has GPT headers
hw-cloud-xngy-ecs-test-0002  /dev/vdb  hdd   Unknown    Unknown  0x1af4  N/A    74321443-d05c-4803-9   107G  Unknown  N/A    N/A    No         Has GPT headers
hw-cloud-xngy-ecs-test-0003  /dev/vdb  hdd   Unknown    Unknown  0x1af4  N/A    f9c0ddbb-7ede-4958-8   107G  Unknown  N/A    N/A    No         Has GPT headers
```

执行命令 `parted /dev/vdb mklabel msdos` 删除 GPT 分区表，即可。也可以使用 `sgdisk` 命令进行磁盘清理。

清理完成后，开始在所有节点上添加 OSD

```bash
ceph orch apply osd --all-available-devices
```

## 基本部署完成

组成了一个三节点的 Ceph 集群效果如下：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sx1zt0/1630850693982-c0ecf1f3-1f37-4f61-8b7c-84c1447ac04f.png)

Ceph 集群中，除了监控套件以外，有 3 个 ceph-mon、2 个 ceph-mgr、3 个 ceph-crash、6 个 ceph-osd。从服务角度看，当前有三个服务：mon、mgr、osd，使用 `ceph -s` 命令查看：

```bash
~]# ceph -s
  cluster:
    id:     24750534-0e45-11ec-9849-7bf16e3e2cb9
    health: HEALTH_OK

  services:
    mon: 3 daemons, quorum hw-cloud-xngy-ecs-test-0001,hw-cloud-xngy-ecs-test-0002,hw-cloud-xngy-ecs-test-0003 (age 10h)
    mgr: hw-cloud-xngy-ecs-test-0001.afnavu(active, since 10h), standbys: hw-cloud-xngy-ecs-test-0002.jucqwq
    osd: 6 osds: 6 up (since 10h), 6 in (since 10h)

  data:
    pools:   7 pools, 145 pgs
    objects: 253 objects, 10 KiB
    usage:   343 MiB used, 600 GiB / 600 GiB avail
    pgs:     145 active+clean

```

## 添加 RGW 服务

```bash
# 为集群中两个节点添加标签，以准备后续将 radosgw 部署到具有 rgw 标签的节点上
ceph orch host label add hw-cloud-xngy-ecs-test-0001 rgw
ceph orch host label add hw-cloud-xngy-ecs-test-0001 rgw

# 部署 rgw
ceph orch apply rgw foo '--placement=label:rgw count-per-host:2' --port=8000

```

此时，节点 1 和节点 2 上，每个节点都有运行有两个 ceph-rgw 实例，可以对外提供服务

```bash
~]# docker ps -a | grep rgw
9ad5fe363abd   ceph/ceph                    "/usr/bin/radosgw -n…"   24 seconds ago   Up 24 seconds             ceph-24750534-0e45-11ec-9849-7bf16e3e2cb9-rgw.foo.hw-cloud-xngy-ecs-test-0001.rqcvxl
519d6a07f001   ceph/ceph                    "/usr/bin/radosgw -n…"   26 seconds ago   Up 26 seconds             ceph-24750534-0e45-11ec-9849-7bf16e3e2cb9-rgw.foo.hw-cloud-xngy-ecs-test-0001.hsjpqq
root@hw-cloud-xngy-ecs-test-0001:~# ss -ntlp | grep gw
LISTEN    0         128                0.0.0.0:8000             0.0.0.0:*        users:(("radosgw",pid=16598,fd=57))
LISTEN    0         128                0.0.0.0:8001             0.0.0.0:*        users:(("radosgw",pid=17447,fd=57))
LISTEN    0         128                   [::]:8000                [::]:*        users:(("radosgw",pid=16598,fd=58))
LISTEN    0         128                   [::]:8001                [::]:*        users:(("radosgw",pid=17447,fd=58))

root@hw-cloud-xngy-ecs-test-0002:~# docker ps -a | grep rgw
1ab21ec662de   ceph/ceph                    "/usr/bin/radosgw -n…"   23 seconds ago   Up 23 seconds             ceph-24750534-0e45-11ec-9849-7bf16e3e2cb9-rgw.foo.hw-cloud-xngy-ecs-test-0002.zsrkkp
10ad804e541d   ceph/ceph                    "/usr/bin/radosgw -n…"   25 seconds ago   Up 25 seconds             ceph-24750534-0e45-11ec-9849-7bf16e3e2cb9-rgw.foo.hw-cloud-xngy-ecs-test-0002.giyyjf
root@hw-cloud-xngy-ecs-test-0002:~# ss -ntlp | grep gw
LISTEN    0         128                0.0.0.0:8000             0.0.0.0:*        users:(("radosgw",pid=14294,fd=57))
LISTEN    0         128                0.0.0.0:8001             0.0.0.0:*        users:(("radosgw",pid=15152,fd=57))
LISTEN    0         128                   [::]:8000                [::]:*        users:(("radosgw",pid=14294,fd=58))
LISTEN    0         128                   [::]:8001                [::]:*        users:(("radosgw",pid=15152,fd=58))
```

当访问这些端口时，显示如下内容，则说明已经可以提供 S3 服务了

```bash
root@hw-cloud-xngy-ecs-test-0001:~# curl 192.168.1.201:8000
<?xml version="1.0" encoding="UTF-8"?><ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Owner><ID>anonymous</ID><DisplayName></DisplayName></Owner><Buckets></Buckets></ListAllMyBucketsResult>
```

创建一个系统用户并记录下 ak 与 sk。通过这个用户的信息，可以使用 blemmenes/radosgw_usage_exporter 导出对象存储监控指标

```bash
root@hw-cloud-xngy-ecs-test-0001:~# radosgw-admin user create --uid=lichenhao --display-name=lichenhao --system
{
    "user_id": "lichenhao",
    "display_name": "lichenhao",
......
    "keys": [
        {
            "user": "lichenhao",
            "access_key": "4O23LGQI3UAUKSSO50UK",
            "secret_key": "JQLul4q2r2qo1vyOLpQ4FVUnh3LWfiNyuiZHQDT6"
        }
    ],
......
}

```

[启动对象网关的管理前端](https://docs.ceph.com/en/pacific/mgr/dashboard/#enabling-the-object-gateway-management-frontend)

```bash
radosgw-admin user info --uid=lichenhao | jq .keys[0].access_key > ak
radosgw-admin user info --uid=lichenhao | jq .keys[0].secret_key > sk
ceph dashboard set-rgw-api-access-key -i ak
ceph dashboard set-rgw-api-secret-key -i sk
```

此时从 Ceph 的 Web 页面中，可以从 Object Gateway 标签中看到内容了：

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sx1zt0/1630856235488-d2d2e334-a0a9-4d41-aa17-06522f30d11a.png)

使用 192.168.1.202:8000 作为 Endpoint，以及 lichenhao 用户的 ak、sk，可以通过 S3 Brower 访问 Ceph 提供的对象存储。注意：由于此时没有开启 SSL，所以 S3 Brower 也要关闭 SSL。

## 添加监控服务(可选)

```bash
ceph orch apply node-exporter
ceph orch apply alertmanager
ceph orch apply prometheus
ceph orch apply grafana
```

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sx1zt0/1630835511543-fb85907a-97d5-4f99-80d5-2214a0236810.png)

# 其他

```bash
docker run --rm --network host --name ceph-rgw-exporter \
blemmenes/radosgw_usage_exporter:latest \
-H 172.38.30.2:7480 \
-a F52JL32RD8NWI78XT3A9 \
-s jjs3uAIGJYFMyprkHov6D85D1YGSo0HowisHZmJl \
-p 9243
```

```bash
docker run --rm --network host --name ceph-rgw-exporter \
blemmenes/radosgw_usage_exporter:latest \
-H 192.168.1.202:8000 \
-a 4O23LGQI3UAUKSSO50UK \
-s JQLul4q2r2qo1vyOLpQ4FVUnh3LWfiNyuiZHQDT6 \
-p 9243
```

异常
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sx1zt0/1630750629528-40ac128e-4c7c-4ccf-9aa4-d64741aae089.png)
正常
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/sx1zt0/1630835261055-137daaea-90de-4045-a62f-5a0e28077860.png)
