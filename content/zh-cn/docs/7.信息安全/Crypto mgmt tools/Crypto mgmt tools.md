---
title: Crypto mgmt tools
linkTitle: Crypto mgmt tools
date: 2024-03-27T10:24
weight: 1
---

# 概述

> 参考：
> 
> -

# ACME

> 参考：
>
> - [RFC-8555](https://datatracker.ietf.org/doc/html/rfc8555)
> - [Wiki，ACME](https://en.wikipedia.org/wiki/Automatic_Certificate_Management_Environment)

**Automatic Certificate Management Environment(自动证书管理环境，简称 ACME)** 是一种通信协议，用于自动化证书颁发机构与其用户的 Web 服务器之间的交互，允许以非常低的成本自动部署 PKI。它是由 ISRG 为他们的 Let's Encrypt 服务设计的。

## ACME 服务提供商

支持免费或低成本基于 ACME 的证书服务的提供商包括 Let's Encrypt、Buypass Go SSL、ZeroSSL 和 SSL.com。许多其他证书颁发机构和软件供应商提供 ACME 服务，作为 [Entrust](https://en.wikipedia.org/wiki/Entrust) 和 [DigiCert](https://en.wikipedia.org/wiki/DigiCert) 等付费 PKI 解决方案的一部分。

## ACME 的实现

想要实现自动签证书，要经过如下几个步骤

- 验证要签名证书所使用的域名是属于我的。这个验证过程又有多种途径
  - **DNS 验证** # 通过提供域名注册商的认证信息(比如 ak、sk)，ACME 程序将会从域名注册商处验证域名属于我的
  - **Web 验证** # 通过域名访问 Web 服务。由于自己可以配置域名解析，所以只要域名可以解析到运行 ACME 程序的设备上，那么 ACME 程序就认为这个域名是属于我的。
    - 这种方式有个弊端：首先要保证自己的域名可以解析到运行 ACME 程序的设备上；然后还要保证 ACME 程序可以通过域名访问到自己。这在国内没有备案的域名是不方便的

# Let's Encrypt

> 参考：
>
> - [官网](https://letsencrypt.org/)
> - [Wiki，Let's Encrypt](https://en.wikipedia.org/wiki/Let%27s_Encrypt)
> - [GitHub 项目，certbot/certbot](https://github.com/certbot/certbot)

## 使用 certbot 创建证书

Let's Encrypt 使用 certbot 工具为我们签署证书

> 注意：保证执行 certbot 的服务器的 80 端口是可以被公网访问的，且保证签署证书时提供的域名是可以解析的(即已备案或无需备案)

证书申请成功后将会出现如下提示：

```bash
~]# certbot certonly --standalone
Saving debug log to /var/log/letsencrypt/letsencrypt.log
Plugins selected: Authenticator standalone, Installer None
Please enter in your domain name(s) (comma and/or space separated)  (Enter 'c'
to cancel): 这里输入自己的${域名}
Obtaining a new certificate
Performing the following challenges:
http-01 challenge for ${域名}
Waiting for verification...
Cleaning up challenges

IMPORTANT NOTES:
 - Congratulations! Your certificate and chain have been saved at:
   /etc/letsencrypt/live/${域名}/fullchain.pem
   Your key file has been saved at:
   /etc/letsencrypt/live/${域名}/privkey.pem
   Your cert will expire on 2022-09-05. To obtain a new or tweaked
   version of this certificate in the future, simply run certbot
   again. To non-interactively renew *all* of your certificates, run
   "certbot renew"
 - If you like Certbot, please consider supporting our work by:

   Donating to ISRG / Let's Encrypt:   https://letsencrypt.org/donate
   Donating to EFF:                    https://eff.org/donate-le

```

根据提示，在 `/etc/letsencrypt/live/${DNS}/` 目录下，可以找到已经签署的证书及私钥

# ZeroSLL

> 参考：
>
> - [官网](https://zerossl.com/)
> - [GitHub 组织，ZeroSSL](https://github.com/zerossl/)

ZeroSLL 是一个 CA 机构，可以为所有人提供快速、可靠、自由的 SSL 保护。

在 [ZeroSSL Partners & ACME Clients](https://zerossl.com/features/acme/#clients) 这里可以看到所有可以支持使用 ZeroSSL 实现 ACME 的客户端应用，其中就包括 [acme.sh](#acme.sh)

# acme.sh

> 参考：
>
> - [GitHub 项目，acemsh-official/acme.sh](https://github.com/acmesh-official/acme.sh)

acmesh 是一个使用纯粹 Shell（Unix shell）语言编写的 ACME 协议客户端。

在过去，acme.sh 是一个常用的工具来申请、部署和续期免费的 SSL 证书，而默认的免费 SSL 证书 CA 机构是 Let's Encrypt。然而，从 acme.sh v3.0 开始，默认的免费 SSL 证书 CA 机构变更为 ZeroSSL。这意味着 acmesh.sh 现在支持使用 ZeroSSL 来申请免费的SSL证书。

### 安装 acme.sh

acme.sh 是一个纯 Shell 脚本，首先下载这个脚本

```bash
export MY_EMAIL="我的邮箱"
curl https://get.acme.sh | sh -s $MY_EMAIL
source ~/.bashrc
```

该脚本会创建 `~/.acme.sh/` 目录，并在该目录中安装 acme.sh 脚本。acme.sh 的配置文件，生成证书默认保存的位置也在这里。同时会在 `~/.bashrc` 文件中添加 `. "/root/.acme.sh/acme.sh.env"`

**安装过程不会污染已有的系统任何功能和文件**, 所有的修改都限制在安装目录中: `~/.acme.sh/`

会在系统的 Crontab 中创建一个逻辑

```bash
~]# crontab -l
26 0 * * * /root/.acme.sh/acme.sh --cron --home "/root/.acme.sh" > /dev/null
```

## 关联文件与配置

**~/.acme.sh/** # acme.sh 程序的主要工作目录

- **./account.conf** # 包括不同 DNS 提供商的认证信息、账户信息、等等
- **./$DOMAIN_NAME/** # 已处理域名的信息保存路径
  - **./$DOMAIN_NAME.conf** # 与该域名相关的配置信息。`acme.sh --info -d example.com` 命令读取的就是这个文件

## 生成证书

### 不同 DNS 提供商的处理

#### 使用阿里云解析生成证书

> 参考：
>
> - [GitHub 项目 Wiki，acmesh-official/acme.sh-dnsapi-使用阿里云域名 API 自动颁发证书](https://github.com/acmesh-official/acme.sh/wiki/dnsapi#11-use-aliyun-domain-api-to-automatically-issue-cert)

从 [阿里云控制台-RAM 访问控制-身份管理-用户](https://ram.console.aliyun.com/users) 处创建用户并获取 AK、SK

```bash
export Ali_Key="AccessKeyId"
export Ali_Secret="AccessKeySecret"
```

#### 使用 Name.com 生成证书

> 参考：
>
> - [GitHub 项目 Wiki，acmesh-official/acme.sh-dnsapi-使用 Name.com API](https://github.com/acmesh-official/acme.sh/wiki/dnsapi#28-use-namecom-api)

从 <https://www.name.com/zh-cn/account/settings/api> 创建 Token

```bash
export Namecom_Username="XXXX"
export Namecom_Token="XXXXX"
```

注意：这俩变量要使用 PRODUCTION(生产)环境的。Name.com 创建完 Token 后会有两个~一个用于生产，一个用于测试，对应不用的 API

```bash
# 生成到 Nginx 目录
acme.sh --issue --dns dns_namecom -d 102205.xyz -d *.102205.xyz
```

#### 使用 DNSPod 生成证书

```bash
export DP_Id="XXXXXX"
export DP_Key="YYYYYYYYYYYYYYYY"

acme.sh --issue --dns dns_dp -d 102205.xyz -d *.102205.xyz
```

### 拷贝证书

前面证书生成以后, 接下来需要把证书 copy 到真正需要用它的地方.

注意, 默认生成的证书都放在安装目录下: `~/.acme.sh/`, 请不要直接使用此目录下的文件, 例如: 不要直接让 nginx/apache 的配置文件使用这下面的文件. 这里面的文件都是内部使用, 而且目录结构可能会变化.

正确的使用方法是使用 `--install-cert` 命令,并指定目标位置, 然后证书文件会被copy到相应的位置

```bash
acme.sh --install-cert -d 102205.xyz \
--key-file ~/projects/DesistDaydream/cloud-native-apps/compose/nginx-lch/config/certs/102205.key \
--fullchain-file ~/projects/DesistDaydream/cloud-native-apps/compose/nginx-lch/config/certs/102205.pem \
--reloadcmd "docker exec -it nginx-geoip2 nginx -c /etc/nginx/nginx/nginx.conf -s reload"
```

### 查看证书信息

```bash
acme.sh --info -d example.com
# 会输出如下内容：
DOMAIN_CONF=/root/.acme.sh/example.com/example.com.conf
Le_Domain=example.com
Le_Alt=no
Le_Webroot=dns_ali
Le_PreHook=
Le_PostHook=
Le_RenewHook=
Le_API=https://acme-v02.api.letsencrypt.org/directory
Le_Keylength=
Le_OrderFinalize=https://acme-v02.api.letsencrypt.org/acme/finalize/23xxxx150/781xxxx4310
Le_LinkOrder=https://acme-v02.api.letsencrypt.org/acme/order/233xxx150/781xxxx4310
Le_LinkCert=https://acme-v02.api.letsencrypt.org/acme/cert/04cbd28xxxxxx349ecaea8d07
Le_CertCreateTime=1649358725
Le_CertCreateTimeStr=Thu Apr  7 19:12:05 UTC 2022
Le_NextRenewTimeStr=Mon Jun  6 19:12:05 UTC 2022
Le_NextRenewTime=1654456325
Le_RealCertPath=
Le_RealCACertPath=
Le_RealKeyPath=/etc/acme/example.com/privkey.pem
Le_ReloadCmd=service nginx force-reload
Le_RealFullChainPath=/etc/acme/example.com/chain.pem
```

这些信息是保存在 `~/.acme.sh/DOMAIN_NAME/DOMAIN_NAME.conf` 文件中的，我们可以直接修改文件内容，比如 ReloadCmd、等等