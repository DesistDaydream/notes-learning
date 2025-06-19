---
title: 自建CA并签署证书
linkTitle: 自建CA并签署证书
weight: 20
---

# 概述

> 参考：
>
> - 

主要依赖 genrsa, req, x509 这三个子命令

## 0. 设定变量

```bash
# CA 相关信息
export CA_KEY="ca.key"
export CA_CRT="ca.crt"
export CA_COMMON_NAME="DesistDaydream-CA"
# 需要签发的证书的相关信息
export SSL_COMMON_NAME="desistdaydream.it"
export SSL_KEY=${SSL_COMMON_NAME}.key
export SSL_CSR=${SSL_COMMON_NAME}.csr
export SSL_CRT=${SSL_COMMON_NAME}.crt
```
## 1. 创建根 CA 私钥

首先，创建一个强 RSA 私钥用于你的根 CA：

```bash
openssl genrsa -out ${CA_KEY} 4096
```

## 2. 创建根 CA 证书

使用生成的私钥创建自签名的根 CA 证书：

```bash
openssl req -x509 -new -nodes -key ${CA_KEY} -sha256 -days 3650 -out ${CA_CRT} \
  -subj "/C=CN/CN=${CA_COMMON_NAME}"
```

在提示时需填写各项信息，其中最重要的是 Common Name (CN)，可以设为 "DesistDaydream-CA" 或类似名称。

## 3. 创建域名私钥

为你的通配符域名创建一个私钥：

```bash
openssl genrsa -out ${SSL_KEY} 2048
```

## 4. 创建 CSR (证书签名请求)

创建 CSR 文件，注意指定通配符域名：

```bash
openssl req -new -key ${SSL_KEY} -out ${SSL_CSR} \
  -subj "/CN=${SSL_COMMON_NAME}"
```

## 5. 创建用于生成 x509v3 扩展信息的 OpenSSL 配置文件

创建一个名为 `v3.ext` 的文件，内容如下：

```bash
tee v3.ext > /dev/unll <<EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = desistdaydream.it
DNS.2 = *.desistdaydream.it
EOF
```

这个配置确保证书支持通配符和裸域名。

## 6. 使用根 CA 签发证书

使用根 CA 签发证书：

```bash
openssl x509 -req -in ${SSL_CSR} -extfile v3.ext \
  -CA ${CA_CRT} -CAkey ${CA_KEY} -CAcreateserial \
  -out ${SSL_CRT} -days 3650 -sha256
```

## 7. 验证证书信息

验证你的证书是否包含正确的通配符域名信息：

```bash
openssl x509 -in desistdaydream.it.crt -text -noout
```

检查输出中的 Subject Alternative Name 部分是否包含 `*.desistdaydream.it`。

## 8. 在系统上安装根 CA 证书

要让浏览器信任你签发的证书，需要在操作系统或浏览器中安装根 CA 证书。各操作系统安装方法不同：

- Linux (Debian/Ubuntu): `sudo cp ca.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates`
- macOS: 双击 `ca.crt` 文件，用钥匙串访问添加到系统，并设置为"始终信任"
- Windows: 双击 `ca.crt` 文件，安装到"受信任的根证书颁发机构"

## 9. 使用证书配置服务器

配置 Web 服务器（如 Nginx、Apache）使用 `desistdaydream.it.crt` 和 `desistdaydream.it.key` 文件。

现在已成功为 `*.desistdaydream.it` 自建 CA 并签发了通配符证书，可以用于所有子域名。

# 脚本

> 参考：
>
> 参考: https://rancher2.docs.rancher.cn/docs/installation/options/self-signed-ssl/_index

```bash
#!/bin/bash -e
help ()
{
    echo  ' ================================================================ '
    echo  ' --ssl-domain: 生成ssl证书需要的主域名，如不指定则默认为www.rancher.local，如果是ip访问服务，则可忽略；'
    echo  ' --ssl-trusted-ip: 一般ssl证书只信任域名的访问请求，有时候需要使用ip去访问server，那么需要给ssl证书添加扩展IP，多个IP用逗号隔开；'
    echo  ' --ssl-trusted-domain: 如果想多个域名访问，则添加扩展域名（SSL_TRUSTED_DOMAIN）,多个扩展域名用逗号隔开；'
    echo  ' --ssl-size: ssl加密位数，默认2048；'
    echo  ' --ssl-cn: 国家代码(2个字母的代号),默认CN;'
    echo  ' --ssl-date: 指定ca签署的证书有效期,默认10年;'
    echo  ' --ca-date: 指定ca证书有效期,默认10年;'
    echo  ' 使用示例:'
    echo  ' ./create_self-signed-cert.sh --ssl-domain=www.test.com --ssl-trusted-domain=www.test2.com \ '
    echo  ' --ssl-trusted-ip=1.1.1.1,2.2.2.2,3.3.3.3 --ssl-size=2048 --ssl-date=36500 --ca-date=36500'
    echo  ' ================================================================'
}


case "$1" in
    -h|--help) help; exit;;
esac


if [[ $1 == '' ]];then
    help;
    exit;
fi


CMDOPTS="$*"
for OPTS in $CMDOPTS;
do
    key=$(echo ${OPTS} | awk -F"=" '{print $1}' )
    value=$(echo ${OPTS} | awk -F"=" '{print $2}' )
    case "$key" in
        --ssl-domain) SSL_DOMAIN=$value ;;
        --ssl-trusted-ip) SSL_TRUSTED_IP=$value ;;
        --ssl-trusted-domain) SSL_TRUSTED_DOMAIN=$value ;;
        --ssl-size) SSL_SIZE=$value ;;
        --ssl-date) SSL_DATE=$value ;;
        --ca-date) CA_DATE=$value ;;
        --ssl-cn) CN=$value ;;
    esac
done

## 国家代码(2个字母的代号),默认CN;
CN=${CN:-CN}

# CA相关配置
CA_DATE=${CA_DATE:-3650}
CA_KEY=${CA_KEY:-ca.key}
CA_CERT=${CA_CERT:-ca.crt}
CA_COMMON_NAME=DesistDaydream-CA


# ssl相关配置
SSL_CONFIG=${SSL_CONFIG:-$PWD/openssl.cnf}
SSL_DOMAIN=${SSL_DOMAIN:-'www.desistdaydream.local'}
SSL_DATE=${SSL_DATE:-3650}
SSL_SIZE=${SSL_SIZE:-2048}

SSL_KEY=$SSL_DOMAIN.key
SSL_CSR=$SSL_DOMAIN.csr
SSL_CERT=$SSL_DOMAIN.crt


echo -e "\033[32m ---------------------------- \033[0m"
echo -e "\033[32m       | 生成 SSL Cert |       \033[0m"
echo -e "\033[32m ---------------------------- \033[0m"


if [[ -e ./${CA_KEY} ]]; then
    echo -e "\033[32m ====> 1. 发现已存在CA私钥，备份"${CA_KEY}"为"${CA_KEY}"-bak，然后重新创建 \033[0m"
    mv ${CA_KEY} "${CA_KEY}"-bak
    openssl genrsa -out ${CA_KEY} ${SSL_SIZE}
else
    echo -e "\033[32m ====> 1. 生成新的CA私钥 ${CA_KEY} \033[0m"
    openssl genrsa -out ${CA_KEY} ${SSL_SIZE}
fi


if [[ -e ./${CA_CERT} ]]; then
    echo -e "\033[32m ====> 2. 发现已存在CA证书，先备份"${CA_CERT}"为"${CA_CERT}"-bak，然后重新创建 \033[0m"
    mv ${CA_CERT} "${CA_CERT}"-bak
    openssl req -x509 -sha256 -new -nodes -key ${CA_KEY} -days ${CA_DATE} -out ${CA_CERT} -subj "/C=${CN}/CN=${CA_COMMON_NAME}"
else
    echo -e "\033[32m ====> 2. 生成新的CA证书 ${CA_CERT} \033[0m"
    openssl req -x509 -sha256 -new -nodes -key ${CA_KEY} -days ${CA_DATE} -out ${CA_CERT} -subj "/C=${CN}/CN=${CA_COMMON_NAME}"
fi


echo -e "\033[32m ====> 3. 生成Openssl配置文件 ${SSL_CONFIG} \033[0m"
cat > ${SSL_CONFIG} <<EOM
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOM


if [[ -n ${SSL_TRUSTED_IP} || -n ${SSL_TRUSTED_DOMAIN} ]]; then
    cat >> ${SSL_CONFIG} <<EOM
subjectAltName = @alt_names
[alt_names]
EOM
    IFS=","
    dns=(${SSL_TRUSTED_DOMAIN})
    dns+=(${SSL_DOMAIN})
    for i in "${!dns[@]}"; do
      echo DNS.$((i+1)) = ${dns[$i]} >> ${SSL_CONFIG}
    done


    if [[ -n ${SSL_TRUSTED_IP} ]]; then
        ip=(${SSL_TRUSTED_IP})
        for i in "${!ip[@]}"; do
          echo IP.$((i+1)) = ${ip[$i]} >> ${SSL_CONFIG}
        done
    fi
fi


echo -e "\033[32m ====> 4. 生成服务SSL KEY ${SSL_KEY} \033[0m"
openssl genrsa -out ${SSL_KEY} ${SSL_SIZE}


echo -e "\033[32m ====> 5. 生成服务SSL CSR ${SSL_CSR} \033[0m"
openssl req -sha256 -new -key ${SSL_KEY} -out ${SSL_CSR} -subj "/C=${CN}/CN=${SSL_DOMAIN}" -config ${SSL_CONFIG}


echo -e "\033[32m ====> 6. 生成服务SSL CERT ${SSL_CERT} \033[0m"
openssl x509 -sha256 -req -in ${SSL_CSR} -CA ${CA_CERT} \
    -CAkey ${CA_KEY} -CAcreateserial -out ${SSL_CERT} \
    -days ${SSL_DATE} -extensions v3_req \
    -extfile ${SSL_CONFIG}


echo -e "\033[32m ====> 7. 证书制作完成 \033[0m"
echo
echo -e "\033[32m ====> 8. 以YAML格式输出结果 \033[0m"
echo "----------------------------------------------------------"

echo -e "\033[32m ====> 9. 附加CA证书到Cert文件 \033[0m"
cat ${CA_CERT} >> ${SSL_CERT}

echo -e "\033[32m ====> 10. 重命名服务证书 \033[0m"
echo "cp ${SSL_DOMAIN}.key tls.key"
cp ${SSL_DOMAIN}.key tls.key
echo "cp ${SSL_DOMAIN}.crt tls.crt"
cp ${SSL_DOMAIN}.crt tls.crt
```
