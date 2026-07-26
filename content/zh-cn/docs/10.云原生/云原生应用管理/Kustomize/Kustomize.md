---
title: Kustomize
weight: 1
---

# 概述

> 参考：
>
> - [GitHub 项目，kubernetes-sigs/kustomize](https://github.com/kubernetes-sigs/kustomize)
> - [官网](https://kustomize.io/)
> - [官方文档，任务-管理 K8S 对象-使用 Kustomize 对 Kubernetes 对象进行声明式管理](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/)

Kustomize 是一个通过 Kustomization 文件来管理 Manifests 的应用程序，Manifests 就是用来生成 K8S 对象的 YAML 格式的文件。Kustomize 可以让我们自定义原始的，无模板的 YAML 文件，以用于多种用途，而原始的 YAML 则保持不变并可以使用。

从 K8S 的 1.14 版本开始，Kustomize 被集成在 kubectl 工具中，可以通过下面几种方式来使用 Kustomize

- kustomize 子命令
- -k,--kustomize 标志来代替 kubectl apply 命令中的 -f 标志。
    - 比如 kubectl apply -k \<KustomizationDIR>

Kustomize 与 Helm 非常类似，都可以用来渲染声明 Kubernetes 资源的 Manifests 文件，并部署到集群中，只不过，Kustomize 更轻便，更易用，但是，不像 Helm，并不能包装成 Chart 并统一上传到仓库中。

## Kustomization

Kustomize 就是通过 Kustomization 实现其功能的。Kustomization 有多种理解方式：

- 一个名为 kustomization.yaml 的文件
- 包含 kustomization.yaml 文件的目录
- 当然，也可以直接用 Kustomization 来表示 Kustomize

在不同环境中，Kustomization 可以有不同的含义。

Kustomization 目录的概念，与 Helm 的 Chart 概念类似，是一组用于描述一个应用的 Manifests 文件的集合，并且包含一个 kustomization.yaml 文件来定义如何组织这些 Manifests 文件。而 kustomization.yaml 文件，就是一个 YAML 格式的文件，Kustomize 也继承了 Kubernetes 的哲学，一切介资源，只不过，现阶段 Kustomize 只有一个资源，就是 `kustomize.config.k8s.io/v1beta1` 下的 **Kustomization 资源**。应用一个 Kustomization 资源，实际上就是声明了一个应用。

除了下面的[基本使用示例](#基本示例)以外，Kustomize 还可以通过配置文件来自动生成 configMap、secret 等资源，通过层次结构来基于某个应用模板定义个性化的内容，为每个资源添加统一的标签或者注释，等等等一系列非常好用的应用管理功能。

## 基本示例

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tx70hw/1620570742728-0e30dc11-37f3-495d-920c-3814c5d1f0d6.jpeg)
通常情况下，一个 Kustomization 目录至少要包含一个 kustomization.yaml 文件，也可以包含若干需要部署的 Manifests 文件。加入目录结构如下：

```bash
~/someApp
├── deployment.yaml
├── kustomization.yaml
└── service.yaml
```

service.yaml 定义如下:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  selector:
    app: myapp
  ports:
    - port: 80
```

deployment.yaml 定义如下：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      name: myapp
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: lchdzh/network-test:v2.0
```

然后在当前文件夹下面添加一个名为 kustomization.yaml 的文件：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - service.yaml
  - deployment.yaml
```

这个文件将是你的基础配置文件，它描述了你使用的资源文件。**apiVersion 与 kind 字段也可以省略不写**

使用 `kubectl kustomize .` 命令运行后的结果如下所示。

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  ports:
    - port: 80
  selector:
    app: myapp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
      name: myapp
    spec:
      containers:
        - image: lchdzh/network-test:v2.0
          name: myapp
```

我们可以看到，Kustomize 将多个 Manifests 文件组合在了一起，并且，通过 `kubectl apply -k .` 命令，我们可以直接部署这些资源到 K8S 集群中。

# Kustomize 的基本功能

## 生成 ConfigMap、Secret 资源

ConfigMap 和 Secret 包含其他 Kubernetes 对象（如 Pod）所需要的配置或敏感数据。 ConfigMap 或 Secret 中数据的来源往往是集群外部，例如某个 `.properties` 文件或者 SSH 密钥文件。 Kustomize 通过 `secretGenerator` 和 `configMapGenerator`，可以基于文件或字面值来生成 Secret 和 ConfigMap 资源。

详见 Kustomization Manifest 详解中的各个字段：

- [configMapGenerator](/docs/10.云原生/云原生应用管理/Kustomize/Kustomization%20Manifest%20 详解.md Manifest 详解.md)
- [secretGenerator](/docs/10.云原生/云原生应用管理/Kustomize/Kustomization%20Manifest%20 详解.md Manifest 详解.md)
- [generatorOptions](/docs/10.云原生/云原生应用管理/Kustomize/Kustomization%20Manifest%20 详解.md Manifest 详解.md)

## 设置贯穿性字段

## 组织和定制各种资源

# Bases and Overlays 功能

Kustomize 和 Docker 比较类似，有很多层组成，每个层都是修改以前的层，正因为有这个理念存在，所以我们可以不断在其他人之上写东西，而不会增加配置的复杂性，构建的最终结果由基础部分和你在上面配置的其他层组成。

## kustomize Overlays 功能的目录结构

```bash
kustomize/
├── base
│   ├── deployment.yaml
│   ├── kust.yaml
│   └── service.yaml
└── overlays
├── prod
│   ├── custom-env.yaml
│   └── kustomization.yaml
└── test
    ├── custom-env.yaml
    └── kustomization.yaml
```

1. 在每个目录中，都有一个名为 kustomization.yaml (文件名不能变)的文件来对当前目录进行配置说明。
2. 一般会有一个 base 目录，用来存放应用运行所需的基础 yaml 配置，和整合这些 yaml 的 kustomize 配置文件
3. 还会有 overlays 目录来存放各种自定义的配置文件，这些文件信息可以附加到 base 中的基础应用配置用。
4. 比如生产环境和测试环境有不同的环境变量，则可以在不同目录中，存放不同的变量 yaml 文件，然后通过 kustomize 来将 yaml 中的信息整合在一起。

## 基础模板

要使用 Kustomize，需要有一个原始的 yaml 文件来描述你想要部署到集群中的任何资源，我们这里将这些 base 文件存储在 ./k8s/base/ 文件夹下面。

这些文件我们**永远**不会直接访问，我们将在它们上面添加一些自定义的配置来创建新的资源定义。

## 根据基础模板定制配置

现在我们想要针对一些特定场景进行定制，比如，针对生产环境和测试环境需要由不同的配置。我们这里并不会涵盖 Kustomize 的整个功能集，而是作为一个标准示例，向你展示这个工具背后的哲学。

首先我们创建一个新的文件夹  k8s/overlays/prod ，其中包含一个名为 kustomzization.yaml 的文件，文件内容如下：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
  - ../../base
```

当前文件夹下面的目录结构如下所示：

```shell
k8s/
├── base
│   ├── deployment.yaml
│   ├── kustomization.yaml
│   └── service.yaml
└── overlays
    └── prod
        └── kustomization.yaml
```

如果现在我们构建这个文件，将会看到和之前在 base 目录中执行 kubectl kustomize 命令一样的结果

接下来我们来为我们的 prod 环境进行一些定制。

### 定制环境变量

在 base 目录的基础模板中，我们不定义任何环境变量，现在我们需要添加一些环境变量在之前的基础模板中而保持原模板文件不变的话。实际上很简单，我们只需要在我们的基础模板上创建一块我们想要模板化的代码块，然后在 kustomization.yaml 文件中引用即可。

比如我们这里定义一个包含环境变量的配置文件：custom-env.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
        - name: app # (1)
          env:
            - name: CUSTOM_ENV_VARIABLE
              value: Value defined by Kustomize
```

Note:注意 (1) 这里定义的 name 是非常重要的，kustomize 会通过该值找到需要修改的容器。

这个 yaml 文件本身是无效的，它只描述了我们希望在上面的基础模板上添加的内容。我们需要将这个文件添加到 k8s/overlays/prod/kustomization.yaml 文件中即可：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ../../base
patchesStrategicMerge:
  - custom-env.yaml
```

现在如果我们来构建下，可以看到如下的输出结果：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  ports:
    - port: 80
  selector:
    app: myapp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
      name: myapp
    spec:
      containers:
        - env:
            - name: CUSTOM_ENV_VARIABLE
              value: Value defined by Kustomize
          name: app
        - image: lchdzh/network-test:v2.0
          name: myapp-container
```

可以看到我们的 env 块已经被合并到了我们的基础模板上了，自定义的 env 变量出现在了 deployment.yaml 文件中。

### 定制副本数量

和上面的例子一样，我们来扩展我们的基础模板来定义一些还没有定义的变量。

这里我们来添加一些关于副本的信息，和前面一样，只需要在一个 YAML 文件中定义副本所需的额外信息块，新建一个名为 replica-and-rollout-strategy.yaml 的文件，内容如下所示：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 10
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
```

和前面一样，在 kustomization.yaml 文件中的 patchesStrategicMerge 字段下面添加这里定制的数据：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
  - ../../base
patchesStrategicMerge:
  - custom-env.yaml
  - replica-and-rollout-strategy.yaml
```

同样，这个时候再使用 kubectl kustomize 命令构建，如下所示：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp
spec:
  ports:
  - port: 80
  selector:
    app: myapp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myapp
  name: myapp
spec:
  replicas: 10
  selector:
    matchLabels:
      app: myapp
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: myapp
      name: myapp
    spec:
      containers:
      - env:
        - name: CUSTOM_ENV_VARIABLE
          value: Value defined by Kustomize
        name: app
      - image: lchdzh/network-test:v2.0
        name: myapp-container
```

我们可以看到副本数量和滚动更新的策略都添加到了基础模板之上了。

### 其他定制及其相关说明

需要注意的是 name 字段，kustomize 一般是通过 name 字段的值来找到需要修改配置的位置的。如果 name 不同或者没有，kustomize 会修改到错误的位置或者报错

还可以定制 namespace，等

## 通过命令行定义 secret

我们常常会通过命令行来添加一个 secret 对象，kustomize 有一个 edit 的子命令可以用来编辑 kustomization.yaml 文件然后创建一个 secret 对象，比如我们这里添加一个如下所示的 secret 对象：

```bash
cd k8s/overlays/prod
kustomize edit add secret sl-demo-app --from-literal=db-password=12345
```

上面的命令会修改 kustomization.yaml 文件添加一个 SecretGenerator 字段在里面。

当然你也可以通过文件（比如--from-file=file/path 或者--from-evn-file=env/path.env）来创建 secret 对象。

通过上面命令创建完 secret 对象后，kustomization.yaml 文件的内容如下所示：

```json
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
- ../../base

patchesStrategicMerge:
- custom-env.yaml
- replica-and-rollout-strategy.yaml
secretGenerator:
- literals:
  - db-password=12345
  name: sl-demo-app
  type: Opaque
```

然后同样的我们回到根目录下面执行 kustomize build 命令构建下模板，输出内容如下所示：

```yaml
$ kustomize build k8s/overlays/prod
apiVersion: v1
data:
  db-password: MTIzNDU=
kind: Secret
metadata:
  name: sl-demo-app-6ft88t2625
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  name: sl-demo-app
spec:
  ports:
  - name: http
    port: 8080
  selector:
    app: sl-demo-app
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sl-demo-app
spec:
  replicas: 10
  selector:
    matchLabels:
      app: sl-demo-app
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: sl-demo-app
    spec:
      containers:
      - env:
        - name: CUSTOM_ENV_VARIABLE
          value: Value defined by Kustomize ❤️
        image: foo/bar:latest
        name: app
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
```

我们可以看到 secret 对象的名称是 sl-demo-app-6ft88t2625，而不是我们定义的 sl-demo-app，这是正常的，因为如果更改了 secret 内容，就可以触发滚动更新了。

同样的，如果我们想要在 Deployment 中使用这个 Secret 对象，我们就可以像之前一样添加一个使用 Secret 的新的层定义即可。

比如我们这里像把 db-password 的值通过环境变量注入到 Deployment 中，我们就可以定义下面这样的新的层信息：（database-secret.yaml）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sl-demo-app
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: "DB_PASSWORD"
          valueFrom:
            secretKeyRef:
              name: sl-demo-app
              key: db.password
```

然后同样的，我们把这里定义的层添加到 k8s/overlays/prod/kustomization.yaml 文件中去：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
  - ../../base

patchesStrategicMerge:
  - custom-env.yaml
  - replica-and-rollout-strategy.yaml
  - database-secret.yaml

secretGenerator:
  - literals:
      - db-password=12345
    name: sl-demo-app
    type: Opaque
```

现在我们来构建整个的 prod 目录，我们会得到如下所示的信息：

```yaml
$ kustomize build k8s/overlays/prod
apiVersion: v1
data:
  db-password: MTIzNDU=
kind: Secret
metadata:
  name: sl-demo-app-6ft88t2625
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  name: sl-demo-app
spec:
  ports:
  - name: http
    port: 8080
  selector:
    app: sl-demo-app
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sl-demo-app
spec:
  replicas: 10
  selector:
    matchLabels:
      app: sl-demo-app
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: sl-demo-app
    spec:
      containers:
        - env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: db.password
                  name: sl-demo-app-6ft88t2625
            - name: CUSTOM_ENV_VARIABLE
              value: Value defined by Kustomize
          image: foo/bar:latest
          name: app
          ports:
            - containerPort: 8080
              name: http
              protocol: TCP
```

我们可以看到 secretKeyRef.name 的值也指定的被修改成了上面生成的 secret 对象的名称。

由于 Secret 是一些私密的信息，所以最好是在安全的环境中来添加上面的 secret 的对象，而不应该和其他代码之类的一起被提交到代码仓库之类的去。

如果是 ConfigMap 的话也是同样的逻辑，最后会生成一个 hash 值的名称，这样在 ConfigMap 更改时可以触发重新部署。

修改镜像

和 secret 资源对象一样，我们可以直接从命令行直接更改镜像或者 tag，如果你需要部署通过 CI/CD 系统标记的镜像的话这就非常有用了。

比如我们这里来修改下镜像的 tag：

```bash
cd k8s/overlays/prod
TAG_VERSION=3.4.5
kustomize edit set image foo/bar=foo/bar:$TAG_VERSION
```

一般情况下 TAG_VERSION 常常被定义在 CI/CD 系统中。

现在的 kustomization.yaml 文件内容如下所示：

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

bases:
  - ../../base

patchesStrategicMerge:
  - custom-env.yaml
  - replica-and-rollout-strategy.yaml
  - database-secret.yaml

secretGenerator:
  - literals:
      - db-password=12345
    name: sl-demo-app
    type: Opaque

images:
  - name: foo/bar
    newName: foo/bar
    newTag: 3.4.5
```

同样回到根目录下面构建该模板，会得到如下所示的信息：

```yaml
$ kustomize build k8s/overlays/prod
apiVersion: v1
data:
  db-password: MTIzNDU=
kind: Secret
metadata:
  name: sl-demo-app-6ft88t2625
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  name: sl-demo-app
spec:
  ports:
    - name: http
      port: 8080
  selector:
    app: sl-demo-app
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sl-demo-app
spec:
  replicas: 10
  selector:
    matchLabels:
      app: sl-demo-app
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: sl-demo-app
    spec:
      containers:
        - env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  key: db.password
                  name: sl-demo-app-6ft88t2625
            - name: CUSTOM_ENV_VARIABLE
              value: Value defined by Kustomize
        image: foo/bar:3.4.5
        name: app
        ports:
          - containerPort: 8080
            name: http
            protocol: TCP
```

我们可以看到 Deployment 的第一个 container.image 已经被修改了 3.4.5 版本了。

最终我们定制的模板文件目录结构如下所示：

```bash
$ tree .
.
└── k8s
    ├── base
    │   ├── deployment.yaml
    │   ├── kustomization.yaml
    │   └── service.yaml
    └── overlays
        └── prod
            ├── custom-env.yaml
            ├── database-secret.yaml
            ├── kustomization.yaml
            └── replica-and-rollout-strategy.yaml

4 directories, 7 files
```

要安装到集群中也很简单：

```bash
kustomize build k8s/overlays/prod | kubectl apply -f -
```

# 总结

在上面的示例中，我们了解到了如何使用 Kustomize 的强大功能来定义你的 Kuberentes 资源清单文件，而不需要使用什么额外的模板系统，创建的所有的修改的块文件都将被应用到原始基础模板文件之上，而不用使用什么花括号之类的修改来更改它（貌似无形中有鄙视了下 Helm 😄）。

Kustomize 中还有很多其他高级用法，比如 mixins 和继承或者允许为每一个创建的对象定义一个名称、标签或者 namespace 等等，你可以在官方的 Kustomize GitHub 代码仓库中查看高级示例和文档。
