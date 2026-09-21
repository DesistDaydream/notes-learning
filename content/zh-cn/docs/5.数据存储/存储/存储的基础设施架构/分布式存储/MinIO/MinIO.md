---
title: MinIO
weight: 1
---

> [!Danger] 已闭源
> 于 2025-12-03 的 [Commit](https://github.com/minio/minio/commit/27742d469462e1561c776f88ca7a1f26816d69e2) 中，修改了 README 信息，声明：将项目变更为维护模式，不再接受新的开发。推荐使用闭源且商业付费的 AIStor
>
> TODO: 用什么替代品？[RustFS](/docs/5.数据存储/存储/存储的基础设施架构/分布式存储/RustFS.md)？
>
> 2026-09-11 连 DockerHub 上的 MinIO 仓库和镜像都删了，官方在有人提问 [issue 21647](https://github.com/minio/minio/issues/21647) 后回复不再提供构建，只保留代码。这不单单是要变现可以理解了，纯纯恶心人了。
# 概述

> 参考：
>
> - [官网](https://min.io/)
> - [GitHub 项目，minio/minio](https://github.com/minio/minio)
> - [GitHub 项目，minio/docs](https://github.com/minio/docs)
>     - https://minio-docs.tf.fo/
>     - https://miniodocs.cc/
>     - https://minio-docs.cc/
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
