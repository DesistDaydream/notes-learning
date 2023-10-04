---
title: CMDB
linkTitle: CMDB
date: 2023-11-14T08:42
weight: 20
---

# 概述

> 参考：
> 
> - [Wiki，CMDB](https://en.wikipedia.org/wiki/Configuration_management_database)

**[Configuration management](https://en.wikipedia.org/wiki/Configuration_management "Configuration management") database(配置管理数据库，简称 CMDB)** 是一个 ITIL 术语，指的是组织用来存储有关硬件和软件资产（通常称为配置项）信息的数据库。

CMDB 项目推荐

- https://github.com/Combodo/iTop
    - https://github.com/vbkunin/itop-docker
- https://github.com/TencentBlueKing/bk-cmdb
- https://github.com/veops/cmdb 维易 CMDB
  - 文章: https://mp.weixin.qq.com/s/6W8DaDb3Y4NmK3rb9NGKAQ

简单的记录功能直接用 WPS 的**轻维表**即可，专治各种花里胡哨的开源产品。

# fiy

> 参考：
> 
> - [GitHub 项目，lanyulei/fiy](https://github.com/lanyulei/fiy)

部署问题：

使用 lz270978971/fiy-ui:latest 镜像构建前端

- `docker run -it  lz270978971/fiy-ui:latest /bin/sh`
- `vim .env.production` 修改其中的 VUE_APP_BASE_API 的值为运行 fiy 服务器机器的 IP
- `npm run build:prod` 打包前端页面
- 详见 [issue 11](https://github.com/lanyulei/fiy/issues/11)

处理后端

- clone fiy 后端代码，并在项目根目录下创建 `statci/ui` 目录，将容器内的 dist/ 目录的所有文件拷贝到刚创建的 static/ui/ 目录下
- 修改 config/settings.yml 文件中的数据库配置，如果需要使用搜索功能，还需要配置 es
- `go build -ldflags="-s -w" -o fiy main.go` 编译后端文件，此时会带着静态资源文件一起编译到一个二进制文件中
- `./fiy server -c config/settings.yml` 启动服务

访问 8000 端口开始使用。

## 蓝鲸 CMDB

> 参考：
> 
> - [GitHub 项目，TencentBlueKing/bk-cmdb](https://github.com/TencentBlueKing/bk-cmdb)

```bash
docker run --name bk-cmdb -d \
-p 8090:8090 \
-v /opt/bk-cmdb/data/mongodb/db:/data/sidecar/mongodb/db \
ccr.ccs.tencentyun.com/bk.io/cmdb-standalone:v3.9.28
```

