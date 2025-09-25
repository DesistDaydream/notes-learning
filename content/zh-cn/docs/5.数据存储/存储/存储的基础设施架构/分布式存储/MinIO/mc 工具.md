---
title: mc 工具
---

# 概述

> 参考：
>
> - [官方文档，MinIO 客户端](https://docs.min.io/minio/baremetal/reference/minio-cli/minio-mc.html)

# 配置

\~/.mc/config.json # mc 从该文件中获取将要操作的 host 信息。可以通过 mc config host 命令管理该文件，也可以直接手动编辑。

```json
~]# cat ~/.mc/config.json
{
	"version": "10",
	"aliases": {
		"gcs": {
			"url": "https://storage.googleapis.com",
			"accessKey": "YOUR-ACCESS-KEY-HERE",
			"secretKey": "YOUR-SECRET-KEY-HERE",
			"api": "S3v2",
			"path": "dns"
		},
		"local": {
			"url": "http://0.0.0.0:9000",
			"accessKey": "minioadmin",
			"secretKey": "ehl@1234",
			"api": "S3v4",
			"path": "auto"
		},
		"play": {
			"url": "https://play.min.io",
			"accessKey": "Q3AM3UQ867SPQQA43P2F",
			"secretKey": "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
			"api": "S3v4",
			"path": "auto"
		},
		"s3": {
			"url": "https://s3.amazonaws.com",
			"accessKey": "YOUR-ACCESS-KEY-HERE",
			"secretKey": "YOUR-SECRET-KEY-HERE",
			"api": "S3v4",
			"path": "dns"
		}
	}
}
```

# Syntax(语法)

**mc \[FLAGS] COMMAND \[COMMAND FLAGS | -h] \[ARGUMENTS...]**

COMMAND

- **alias** # set, remove and list aliases in configuration file
- **ls** # 列出桶和对象
- mb         make a bucket
- rb         remove a bucket
- cp         copy objects
- mirror     synchronize object(s) to a remote site
- cat        display object contents
- head       display first 'n' lines of an object
- pipe       stream STDIN to an object
- share      generate URL for temporary access to an object
- **find** # 搜索对象
- sql        run sql queries on objects
- stat       show object metadata
- mv         move objects
- tree       list buckets and objects in a tree format
- du         summarize disk usage recursively
- retention  set retention for object(s)
- legalhold  manage legal hold for object(s)
- diff       list differences in object name, size, and date between two buckets
- **rm** # 移除桶
- version    manage bucket versioning
- **ilm** # 管理桶的生命周期
- encrypt    manage bucket encryption config
- event      manage object notifications
- watch      listen for object notification events
- undo       undo PUT/DELETE operations
- policy     manage anonymous access to buckets and objects
- tag        manage tags for bucket and object(s)
- replicate  configure server side bucket replication
- admin      manage MinIO servers
- update     update mc to latest release

## OPTIONS

**--config-dir, -C**(STRING) # 配置文件所在目录。`默认值: ~/.mc`

# config

**mc config host COMMAND \[COMMAND FLAGS | -h] \[ARGUMENTS...]**

COMMAND:

- add, a # 添加一个新的主机到配置文件。
- remove, rm # 从配置文件中删除一个主机。
- list, ls # 列出配置文件中的主机。

EXAMPLE

- 添加一个 host
  - **mc config host add miniodev130 http://10.8.208.130:9000 minioadmin minioadmin**

# ilm

管理桶的生命周期

**mc ilm COMMAND \[COMMAND FLAGS | -h] \[ARGUMENTS...]**

COMMAND

- **add** # 为一个桶添加生命周期配置规则
- **edit** # 修改指定 ID 的生命周期配置规则
- **ls** # 列出设置在一个桶上的生命周期配置规则集
- **rm** # 删除生命周期配置规则
- export  export lifecycle configuration in JSON format
- import  import lifecycle configuration in JSON format

## add

**mc ilm add \[FLAGS] TARGET**

FLAGS

- **--expiry-days VALUE** # 创建对象后保留的天数。MinIO 在经过指定的天数后标记要删除的对象

EXAMPLE

- local/loki-bj-net 这个桶中创建的对象将在 7 天后过期
  - **mc ilm add --expiry-days 7 local/loki-bj-net**

## ls

列出设置在一个 bucket 上的生命周期配置规则集，效果如下：

```bash
~]# mc ilm ls local/loki-bj-net
          ID          |     Prefix     |  Enabled   | Expiry |  Date/Days   |  Transition  |    Date/Days     |  Storage-Class   |          Tags
----------------------|----------------|------------|--------|--------------|--------------|------------------|------------------|------------------------
 c36rknaqqqm9ds69f5fg |                |    ✓       |  ✓     |   7 day(s)   |     ✗        |                  |                  |
----------------------|----------------|------------|--------|--------------|--------------|------------------|------------------|------------------------
```

## rm

**mc ilm rm \[FLAGS] TARGET**

**FLAGS**

- **--id VALUE** # 指定要删除的生命周期规则的 ID

EXAMPLE

- 删除 local/chunks 桶中 id 为 cbod0cqqqqm5tvms1svg 的生命周期规则
  - **mc ilm rm --id cbod0cqqqqm5tvms1svg local/chunks**

# rm

**mc rm \[FLAGS] TARGET \[TARGET ...]**

EXAMPLE

- 删除 local 环境下 thanos 桶中的所有对象
  - **mc rm --recursive --force local/thanos**
- 删除 local 环境下 thanos-bj-test 桶中 24 小时前的所有对象
  - **mc rm --recursive --force --older-than 24h local/thanos-bj-test**
