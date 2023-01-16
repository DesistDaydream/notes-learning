---
title: Harbor 云原生注册中心
---

# 概述

> 参考：
> - [GitHub 项目，goharbor/harbor](https://github.com/goharbor/harbor)
> - [官网](https://goharbor.io/)

一个开源可信的云原生注册中心项目，用于存储、签名、扫描内容。项目地址：<https://github.com/goharbor/harbor>

- 云原生注册表：Harbor 支持容器图像和[Helm](https://helm.sh/)图表，可作为容器原生运行时和编排平台等云原生环境的注册表。
- 基于角色的访问控制：用户和存储库通过“项目”进行组织，并且用户可以对项目下的图像或 Helm 图表具有不同的权限。
- 基于策略的复制：可以基于具有多个过滤器（存储库，标签和标签）的策略在多个注册表实例之间复制（同步）图像和图表。如果遇到任何错误，Harbor 会自动重试复制。非常适合负载平衡，高可用性，多数据中心，混合和多云场景。
- 漏洞扫描：Harbor 会定期扫描图像并警告用户存在漏洞。
- LDAP / AD 支持：Harbor 与现有企业 LDAP / AD 集成以进行用户身份验证和管理，并支持将 LDAP 组导入 Harbor 并为其分配适当的项目角色。
- OIDC 支持：Harbor 利用 OpenID Connect（OIDC）来验证由外部授权服务器或身份提供者认证的用户的身份。可以启用单点登录以登录 Harbor 门户。
- 图像删除和垃 00000000000000000000000000000000000 圾回收：可以删除图像，并可以回收其空间。
- 公证员：可以确保图像的真实性。
- 图形用户门户：用户可以轻松浏览，搜索存储库和管理项目。
- 审核：跟踪对存储库的所有操作。
- RESTful API：用于大多数管理操作的 RESTful API，易于与外部系统集成。嵌入式 Swagger UI 可用于探索和测试 API。
- 易于部署：提供在线和离线安装程序。另外，可以使用 Helm Chart 在 Kubernetes 上部署 Harbor。

## Harbor 组件

官方网址：<https://github.com/goharbor/harbor/wiki/Architecture-Overview-of-Harbor>

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ta9n37/311ye872a70f4de4c6264d9f17007043bd25)

如上图所示，Harbor 一般具有三层结构

数据访问层

- k-v storage：键值存储。一般由 redis 实现，提供数据缓存功能并支持为作业服务临时保留作业元数据。
  - 所需镜像：redis-photon
- data storage：支持多个数据持久存储，作为 chart 和 registry 的后端存储
- Database：存储 Harbor 的相关元数据，e.g.projects、users、roles、replication policies、tag、scanners、charts、images
  1. 所需镜像：harbor-db

基本服务层

- Proxy：提供 API 路由功能，一般由 Nginx 实现。Harbor 的组件(e.g.core、registry、web portal 等等)都位于此反代之后。Proxy 将来自浏览器和 docker 客户端的请求转发到后端各种服务
  - 所需镜像：nginx
- Core：Harbor 的核心功能。提供 Web UI 和 RESTful API 以及 Auth 相关功能
  - 所需镜像：harbor-core、
- Job Service：一种通用执行队列服务，允许其他组件/服务通过简单的静态 API 同时提交运行异步任务的请求。e.g.主备互相同步镜像功能。
  - 所需镜像：harbor-jobservice
- Logs：日志收集器，负责将其他模块的日志收集到一个单独的地方，一般通过 rsyslog 保存
  - 所需镜像：harbor-log
- GC Controller：垃圾收集
- Chart Museum：第三方 repository 提供 chart 的访问和管理。e.g.通过 dockerhub 来同步镜像所需镜像：
- Docker Registry：第三方 registry，负责存储 docker 镜像并处理 docker pull 和 push 命令。
  - 所需镜像：registry-photon、harbor-registryctl
- Notary：第三方内容信任服务，负责安全得发布和验证内容。

消费层

- Web Portal：图形用户界面，可帮助用户管理 registry 上的 images。
  - 所需镜像：harbor-portal

# Harbor 部署

在 [**Github Releases**](https://github.com/goharbor/harbor/releases) 中下载 harbor 的部署文件。分为线上版和线下版，分别代表是否需要去公网上 pull 所需镜像。

安装 docker、[**docker-compose**](https://docs.docker.com/compose/install/#install-compose)

解压 harbor 的部署文件，得到 4 个文件：

- harbor.yml # harbor 部署所需配置文件
- install.sh # harbor 安装脚本
- LICENSE
- prepare # 初始化脚本，用于生成 docker-compose 文件和 common 目录

执行 install.sh 即可自动下载 harbor 所需镜像并启动

## 定制 harbor 功能

可以使用 prepare 脚本来定制 harbor 功能，prepare 脚本 根据 harbor.yaml 文件中的配置，以及参数，生成用于启动 harbor 的 docker-compose.yaml 文件

prepare 脚本支持如下参数：

- \--with-notary
- \--with-clair
- \--with-chartmuseum # 为 docker-compose.yaml 文件中添加

Note:在 harobor 目录中，可以使用 ./prepare --with-notary --with-clair --with-chartmuseum 脚本来生成带有其他功能的 docker-compose.yaml 文件

比如我想开启 harbor 的 helm chart 功能，则执行 ./prepare --with-chartmuseum 脚本即可生成

## harbor https 部署

官方文档：<https://github.com/goharbor/website/blob/master/docs/install-config/configure-https.md>

<https://goharbor.io/docs/latest/install-config/configure-https/>

编辑 harbor.yml 文件，取消 https 配置环境的注释，并更新以下两个字段

- certificate # 指定 harbor 所用证书的路径
- private_key # 指定 harbor 所用证书的私钥的路径

```yaml
#set hostname
hostname: registry-1.tj-test.ehualu.com # Note: 不同harbor使用不同hostname

http:
  port: 80

https:
  # https port for harbor, default is 443
  port: 443
  # The path of cert and key files for nginx
  certificate: /data/cert/ehualu.com.crt
  private_key: /data/cert/ehualu.com.key
```

一共有两种办法来获取证书：

其一：去可信机构，够买证书

其二：创建自签 CA 证书，但是需要每个与 harbor 互动的 docker 客户端都需要持有 ca 证书

### 方法一：够买证书

### 方法二：创建自签 CA 证书

**第一步：在任意设备上执行如下命令，创建 CA 以及 harbor 所用证书**

```shell
openssl genrsa -out ca.key 4096

openssl req -x509 -new -nodes -sha512 -days 36500 \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=eHualu" \
  -key ca.key \
  -out ca.crt
```

(可选)使用 CA 创建 harbor 所需证书。也可省略此步骤直接使用 CA 证书作为 harbor 所用。Note：alt_names 使用 3 个域名的原因是为了 harbor 高可用，这样 registry.tj-test.ehualu.com 可以解析为 VIP。

```shell
openssl genrsa -out ehualu.com.key 4096

openssl req -sha512 -new \
  -subj "/C=CN/ST=Tianjin/L=Tianjin/O=eHualu/OU=Operations/CN=Registry" \
  -key ehualu.com.key \
  -out ehualu.com.csr

# 注意修改DNS.X中的内容为harbor所使用的hostname
cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1=registry.tj-test.ehualu.com
IP.1=172.38.50.127
DNS.2=registry-1.tj-test.ehualu.com
DNS.3=registry-2.tj-test.ehualu.com
EOF

openssl x509 -req -sha512 -days 36500 \
  -extfile v3.ext \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -in ehualu.com.csr \
  -out ehualu.com.crt
```

**第二步：在 harbor 服务器上执行如下命令，使用证书部署**

Note：如果直接使用 CA 证书，则只需要将 ca 证书拷贝到 docker 的 certs.d 目录下

```shell
# 拷贝创建的ca或者harbor证书，到harbor部署时需要读取证书的路径下
mkdir -p /data/certs/
cp /PATH/TO/{*.crt,*.key} /data/certs/
```

开始部署

```shell
# 生成docker-compose文件
./prepare --with-chartmuseum
# 如果harbor已经存在，则需要执行下面的命令先停止所有容器
# Note:根据问题实例，和下面的harbor配置。来修改docker-compose文件相关参数后再启动容器
docker-compose down -v
# 然后开始部署
docker-compose up -d
```

**第三步，拷贝证书到 docker 的客户端。i.e.安装了 docker 来使用 harbor 的设备**

在 docker 客户端服务器上创建 docker 认证配置目录。
Note：certs.d 下的目录需要使用完全限定域名。访问哪个 registry，就要使用该 registry 的主机名或者域名。如果目录名使用错误，则会报错：x509: certificate signed by unknown authority

别忘了在 docker 客户端服务器上配置域名解析

```shell
mkdir -p /etc/docker/certs.d/registry.tj-test.ehualu.com/

scp CAServerIP:/PATH/TO/ehualu.com.crt /etc/docker/certs.d/registry.tj-test.ehualu.com/

cat >> /etc/hosts << EOF
VIP registry.tj-test.ehualu.com
IP1 registry-1.tj-test.ehualu.com
IP2 registry-2.tj-test.ehualu.com
EOF
```

部署完成，可以在具有证书的 docker 客户端上对 harbor 进行操作了

参考：

<https://my.oschina.net/u/2306127/blog/785281>

<https://ivanzz1001.github.io/records/post/docker/2018/04/09/docker-harbor-https>

# Harbor 配置

/data/registry # 默认的镜像存储路径

**配置容器时区**

在 docker-compose 文件中每一个容器下添加如下参数

```yaml
environment:
  - TZ=Asia/Shanghai
```

# Harbor 高可用

使用 keepalived 作为两个 harbor 的负载均衡服务，使用 vip 来访问两个 harbor。Note：所以在申请证书的时候需要填写 3 个域名，VIP 所用域名和两个 harbor 的域名

然后在 harbor 的 web 界面即可创建数据同步规则，非常简单

第一步：创建仓库

在 web 界面创建要同步 harbor 中心信息

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ta9n37/311yd80c0b55a83796f50b73f2b473289b6c)

第二步：创建同步规则，示例为每天 0 点同步一次

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ta9n37/311yac2e4be663e43698e8c6c164c95871ee)

第三步
