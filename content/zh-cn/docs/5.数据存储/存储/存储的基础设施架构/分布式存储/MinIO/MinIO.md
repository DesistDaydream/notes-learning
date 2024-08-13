---
title: MinIO
linkTitle: MinIO
date: 2024-08-13T14:23
weight: 1
---

# 概述

> 参考：
>
> - [官网](https://min.io/)
> - [GitHub 项目，minio/minio](https://github.com/minio/minio)
> - <https://mp.weixin.qq.com/s/aRTE_UUQ0GMXhqiemxQnsg>

采集指标

```bash
curl -H "Authorization: Bearer ${TOKEN}" http://localhost:9000/minio/v2/metrics/cluster
```

> Notes: TOKEN 通过 `mc admin prometheus generate ${REMOTE}` 命令生成

# MinIO 部署

> 参考：
>
> - [官方文档，安装](https://docs.min.io/minio/baremetal/tutorials/minio-installation.html)

# docker 启动单点 MinIO

```bash
docker run -p 9000:9000 \
-e "MINIO_ACCESS_KEY=minioadmin" \
-e "MINIO_SECRET_KEY=minioadmin" \
-v /mnt/disk1:/disk1 \
-v /mnt/disk2:/disk2 \
-v /mnt/disk3:/disk3 \
-v /mnt/disk4:/disk4 \
minio/minio server /disk{1...4}
```

MINIO_ACCESS_KEY 与 MINIO_SECRET_KEY 指定连接 MinIO 时所需的认证信息，AK、SK

本地 /mnt 下的 4 个目录
