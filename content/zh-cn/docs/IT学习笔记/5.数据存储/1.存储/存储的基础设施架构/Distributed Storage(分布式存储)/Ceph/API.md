---
title: API
---

# 概述

# Ceph RESTful API

> 参考：
> - [官方文档，Ceph 管理器守护进程-Ceph RESTful API](https://docs.ceph.com/en/latest/mgr/ceph_api/)
> - [GitHub，ceph/ceph/src/pybind/mgr/dashboard/openapi.yaml](https://github.com/ceph/ceph/blob/master/src/pybind/mgr/dashboard/openapi.yaml)(该 API 的 openapi 文件)

在 Dashboard 模块中，提供了一组用于管理集群的 RESTful 风格的 API 接口。这组 API 默认位于 `https://localhost:8443/api` 路径下

在 `/docs` 端点下，可以查看 OpenAPI 格式的信息
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/fl1wmh/1630938180924-97cb2959-3cf0-48bf-b312-57be88e9471d.png)
在 `/dpcs/api.json` 端点可以获取 openapi 格式的 API 信息。

## /api/auth

`/api/auth` 接口获取 Token

```bash
curl -X POST "https://example.com:8443/api/auth" \
  -H  "Accept: application/vnd.ceph.api.v1.0+json" \
  -H  "Content-Type: application/json" \
  -d '{"username": <username>, "password": <password>}'
```

获取 Token 后，其他接口，都可以使用该 Token 进行认证，比如：

```bash
curl -H "Authorization: Bearer $TOKEN" ......
```

## /api/auth/check

`/api/auth/check` 接口可以检查 Token。通常还可以作为对 API 的健康检查接口。

```bash
curl -k -XPOST 'https://example.com:8443/api/auth/check' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d "{\"token\": \"${TOKEN}\"}"
```
