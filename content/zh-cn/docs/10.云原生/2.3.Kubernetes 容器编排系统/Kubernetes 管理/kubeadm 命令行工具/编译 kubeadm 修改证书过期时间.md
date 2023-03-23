---
title: 编译 kubeadm 修改证书过期时间
---

Makefile 文件位置：<https://github.com/kubernetes/kubernetes/blob/master/build/root/Makefile>

这里以 kubernetes 1.18.8 为例。

安装golang, 这里使用v1.13.15

根据上步版本下载源码，下列下载方式任选其一

```bash
git clone --branch v1.19.2 https://github.com/kubernetes/kubernetes.git
wget https://github.com/kubernetes/kubernetes/archive/v1.19.0.tar.gz
```

修改 ca 证书的有效期

```bash
# 将 const duration365d = time.Hour * 24 * 365 改为 const duration365d = time.Hour * 24 * 365 * 10
sed -i 's/\(const duration365d.*365\)/\1* 10/' staging/src/k8s.io/client-go/util/cert/cert.go
```

修改 ca 生成的其余证书的有效期

```bash
cat cmd/kubeadm/app/util/pkiutil/pki_helpers.go
# 找到代码NotAfter:     time.Now().Add(kubeadmconstants.CertificateValidity).UTC(),
# 根据 import 部分找到调用该变量的文件
```

```go
import (
		......
        kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
		......
)

```

```bash
vim cmd/kubeadm/app/constants/constants.go
# 将 CertificateValidity = time.Hour * 24 * 365 改为 time.Hour * 24 * 365 * 100
```

在项目根目录执行 make 命令编译 kubeadm

    make WHAT=cmd/kubeadm GOFLAGS=-v
    # 生成的二进制在 _output/bin/ 目录下

替换 kubeadm 文件
先备份

    cp /usr/bin/kubeadm /usr/bin/kubeadm.backup
    cp _output/bin/kubeadm /usr/bin/kubeadm

参考
<http://blog.sina.com.cn/s/blog_537517170102za73.html>
<https://bbs.huaweicloud.com/blogs/168848>
