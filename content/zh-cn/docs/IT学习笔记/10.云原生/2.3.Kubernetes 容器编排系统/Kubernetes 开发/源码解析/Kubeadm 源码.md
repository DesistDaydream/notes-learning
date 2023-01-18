---
title: Kubeadm 源码
---

# 概述

> ## 参考：

kubeadm 源码在 Kubernetes 中，位置：[kubernetes/kubernetes/cmd/kubeadm](https://github.com/kubernetes/kubernetes/tree/master/cmd/kubeadm)，本文以 1.19 版本为例

## 目录结构

```bash
$ tree -L 4 -d
.
├── app
│   ├── apis
│   │   ├── kubeadm
│   │   │   ├── fuzzer
│   │   │   ├── scheme
│   │   │   ├── v1beta1
│   │   │   ├── v1beta2
│   │   │   └── validation
│   │   └── output
│   │       ├── fuzzer
│   │       ├── scheme
│   │       └── v1alpha1
│   ├── cmd
│   │   ├── alpha
│   │   ├── options
│   │   ├── phases
│   │   │   ├── init
│   │   │   ├── join
│   │   │   ├── reset
│   │   │   ├── upgrade
│   │   │   │   └── node
│   │   │   └── workflow
│   │   ├── upgrade
│   │   └── util
│   ├── componentconfigs
│   ├── constants
│   ├── discovery
│   │   ├── file
│   │   ├── https
│   │   └── token
│   ├── features
│   ├── images
│   ├── phases
│   │   ├── addons
│   │   │   ├── dns
│   │   │   └── proxy
│   │   ├── bootstraptoken
│   │   │   ├── clusterinfo
│   │   │   └── node
│   │   ├── certs
│   │   │   └── renewal
│   │   ├── controlplane
│   │   ├── copycerts
│   │   ├── etcd
│   │   ├── kubeconfig
│   │   ├── kubelet
│   │   ├── markcontrolplane
│   │   ├── patchnode
│   │   ├── selfhosting
│   │   ├── upgrade
│   │   └── uploadconfig
│   ├── preflight
│   └── util
│       ├── apiclient
│       ├── audit
│       ├── certs
│       ├── config
│       │   └── strict
│       │       └── testdata
│       ├── crypto
│       ├── dryrun
│       ├── etcd
│       ├── image
│       ├── initsystem
│       ├── kubeconfig
│       ├── kustomize
│       ├── output
│       ├── patches
│       ├── pkiutil
│       ├── pubkeypin
│       ├── runtime
│       └── staticpod
└── test
    ├── cmd
    │   └── testdata
    │       └── init
    ├── kubeconfig
    └── resources
```

kubeadm 是基于 cobra 框架的命令行工具，入口是 `cmd/kubeadm/kubeadm.go`，包含了众多子命令，代码全部在 `cmd/kubeadm/app` 目录中

- apis # kubeadm API 定义
- cmd # 子命令代码入口
- componentconfigs
- constants
- discovery
- features
- images
- phases # kubeadm 每个阶段的具体执行逻辑。
- preflight
- util #

# init

源码：`[cmd/kubeadm/app/cmd/init.go](https://github.com/kubernetes/kubernetes/blob/master/cmd/kubeadm/app/cmd/init.go)`

```go
func newCmdInit(out io.Writer, initOptions *initOptions) *cobra.Command {
	// 通过阶段执行 Kubernetes 集群初始化
	initRunner.AppendPhase(phases.NewPreflightPhase())
	initRunner.AppendPhase(phases.NewCertsPhase())
	initRunner.AppendPhase(phases.NewKubeConfigPhase())
	initRunner.AppendPhase(phases.NewKubeletStartPhase())
	initRunner.AppendPhase(phases.NewControlPlanePhase())
	initRunner.AppendPhase(phases.NewEtcdPhase())
	initRunner.AppendPhase(phases.NewWaitControlPlanePhase())
	initRunner.AppendPhase(phases.NewUploadConfigPhase())
	initRunner.AppendPhase(phases.NewUploadCertsPhase())
	initRunner.AppendPhase(phases.NewMarkControlPlanePhase())
	initRunner.AppendPhase(phases.NewBootstrapTokenPhase())
	initRunner.AppendPhase(phases.NewKubeletFinalizePhase())
	initRunner.AppendPhase(phases.NewAddonPhase())
}
```

## certs 阶段

certs 阶段用来生成集群证书

certs 阶段入口：
源码：`[cmd/kubeadm/app/cmd/phases/init/certs.go](https://github.com/kubernetes/kubernetes/blob/master/cmd/kubeadm/app/cmd/phases/init/certs.go)`

```go
func NewCertsPhase() workflow.Phase {
    // 执行工作流
	return workflow.Phase{
		Name:   "certs",
		Short:  "Certificate generation",
        // 执行创建新证书子阶段
		Phases: newCertSubPhases(),
		Run:    runCerts,
		Long:   cmdutil.MacroCommandLongDescription,
	}
}
```

```go
func newCertSubPhases() []workflow.Phase {
    // 从 certsphase.GetDefaultCertList() 中获取需要创建证书的列表
    // 并循环这个列表，注意创建 Kubernetes 集群所需证书
	for _, cert := range certsphase.GetDefaultCertList() {
		var phase workflow.Phase
        // 若没有 CA 则创建 CA
		if cert.CAName == "" {
			phase = newCertSubPhase(cert, runCAPhase(cert))
			lastCACert = cert
        // 使用 CA 创建 Kubernetes 组件所需的所有证书
		} else {
			phase = newCertSubPhase(cert, runCertPhase(cert, lastCACert))
		}
		subPhases = append(subPhases, phase)
	}
}
```

`runCAPhase()` 创建 CA

```go
func runCAPhase(ca *certsphase.KubeadmCert) func(c workflow.RunData) error {
	return func(c workflow.RunData) error {
		// create the new certificate authority (or use existing)
		return certsphase.CreateCACertAndKeyFiles(ca, cfg)
	}
}
```

`runCertPhase()` 使用 CA 创建 Kubernetes 组件所需的所有证书

```go
func runCertPhase(cert *certsphase.KubeadmCert, caCert *certsphase.KubeadmCert) func(c workflow.RunData) error {
	return func(c workflow.RunData) error {
		// 使用 CA 创建 Kubernetes 组件所需的所有证书
		return certsphase.CreateCertAndKeyFilesWithCA(cert, caCert, cfg)
	}
}
```

### 注意

在 kubeadm 的源码中，如果直接看 kubeadm 版本信息会发现 ClusterConfiguration 是在 InitConfiguration 中的

```go
type InitConfiguration struct {
	metav1.TypeMeta
	// ClusterConfiguration holds the cluster-wide information, and embeds that struct (which can be (un)marshalled separately as well)
	// When InitConfiguration is marshalled to bytes in the external version, this information IS NOT preserved (which can be seen from
	// the `json:"-"` tag in the external variant of these API types.
	ClusterConfiguration `json:"-"`
	BootstrapTokens []BootstrapToken
	NodeRegistration NodeRegistrationOptions
	LocalAPIEndpoint APIEndpoint
	CertificateKey string
}
```

但是在 v1beta{1,2,3} 版本中，ClusterConfiguration 则不在。所以后续代码分析，可能会发现使用的 kubeadm-config.yaml 文件中，用的是 InitConfiguration。

### 待创建证书的列表

源码：`cmd/kubeadm/app/phases/certs/certlist.go`
为入口返回 kubeadm 需要创建的所有证书，该代码中还有所有证书的基本信息(比如 DN)。

```go
func GetDefaultCertList() Certificates {
	return Certificates{
		KubeadmCertRootCA(),
		KubeadmCertAPIServer(),
		KubeadmCertKubeletClient(),
		// Front Proxy certs
		KubeadmCertFrontProxyCA(),
		KubeadmCertFrontProxyClient(),
		// etcd certs
		KubeadmCertEtcdCA(),
		KubeadmCertEtcdServer(),
		KubeadmCertEtcdPeer(),
		KubeadmCertEtcdHealthcheck(),
		KubeadmCertEtcdAPIClient(),
	}
}
```

### 生成 CA 证书与私钥

源码：`cmd/kubeadm/app/phases/certs/certs.go`
生成 CA 证书和私钥，并写入到默认的 /etc/kubernetes/pki/ 目录中

```go
func CreateCACertAndKeyFiles(certSpec *KubeadmCert, cfg *kubeadmapi.InitConfiguration) error {
    // 将 kubeadm-config.yaml 文件中 InitConfiguration 的配置传递进去，生成 CA 证书
	caCert, caKey, err := pkiutil.NewCertificateAuthority(certConfig)
	if err != nil {
		return err
	}

    // 将证书写入指定目录中
	return writeCertificateAuthorityFilesIfNotExist(
		cfg.CertificatesDir,
		certSpec.BaseName,
		caCert,
		caKey,
	)
}
```

源码：`cmd/kubeadm/app/util/pkiutil/pki_helpers.go`

```go
func NewCertificateAuthority(config *CertConfig) (*x509.Certificate, crypto.Signer, error) {
    // NewPrivateKey() 直接使用 rsa.GenerateKey() 生成密钥对
	key, err := NewPrivateKey(config.PublicKeyAlgorithm)
    // NewSelfSignedCACert 是 client-go 中的函数，根据配置和密钥生成 CA 证书
	cert, err := certutil.NewSelfSignedCACert(config.Config, key)
	return cert, key, nil
}
```

看一下 client-go 中的代码，非常简洁明了
源码：`kubernetes/client-go/util/cert/cert.go`

```go
func NewSelfSignedCACert(cfg Config, key crypto.Signer) (*x509.Certificate, error) {
	now := time.Now()
	tmpl := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		DNSNames:              []string{cfg.CommonName},
		NotBefore:             now.UTC(),
		NotAfter:              now.Add(duration365d * 10).UTC(),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &tmpl, &tmpl, key.Public(), key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDERBytes)
}
```

所以，我们在 [编译 kubeadm 修改证书过期时间](/docs/IT学习笔记/10.云原生/2.3.Kubernetes%20 容器编排系统/Kubernetes%20 管理/kubeadm%20 命令行工具/编译%20kubeadm%20 修改证书过期时间.md 管理/kubeadm 命令行工具/编译 kubeadm 修改证书过期时间.md) 中会修改 client-go 中的源码，也就是常量 `duration365d` 的值

### 使用 CA 签其他证书

源码：`cmd/kubeadm/app/phases/certs/certs.go`
从磁盘中加载指定的 CA 证书，然后使用该 CA 生成指定的证书

```go
func CreateCertAndKeyFilesWithCA(certSpec *KubeadmCert, caCertSpec *KubeadmCert, cfg *kubeadmapi.InitConfiguration) error {
    // 从磁盘中加载 CA 证书
	caCert, caKey, err := LoadCertificateAuthority(cfg.CertificatesDir, caCertSpec.BaseName)
	// 将 kubeadm-config.yaml 文件中 InitConfiguration 的配置传递进去，并使用加载的 CA 证书和密钥创建其他证书
	return certSpec.CreateFromCA(cfg, caCert, caKey)
}
```

```go
// CreateFromCA makes and writes a certificate using the given CA cert and key.
func (k *KubeadmCert) CreateFromCA(ic *kubeadmapi.InitConfiguration, caCert *x509.Certificate, caKey crypto.Signer) error {
	// 加载 kubeadm-config.yaml 配置
    cfg, err := k.GetConfig(ic)
	// 使用 kubeadm-config.yaml 配置 和 CA 证书与密钥 创建证书
	cert, key, err := pkiutil.NewCertAndKey(caCert, caKey, cfg)
	// 将创建的证书写入默认的 /etc/kubernetes/pki/ 目录中
	err = writeCertificateFilesIfNotExist(
		ic.CertificatesDir,
		k.BaseName,
		caCert,
		cert,
		key,
		cfg,
	)
}
```
