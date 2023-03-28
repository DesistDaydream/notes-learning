---
title: Provisioning 配置
---

# 概述

> 参考：
>
> - [官方文档，管理-Provisioning](https://grafana.com/docs/grafana/latest/administration/provisioning/)

Grafana 在一开始，只能通过 Web 页面(也就是 API)来配置 DataSources(数据源) 和 Dashboard(仪表盘)。这样做有一个缺点，就是无法提前加载数据源和仪表盘。

比如现在有这么一种场景：我想新搭建一个 Grafana，并且包含一些数据源和仪表盘，正常情况是启动服务后，在 Web 页面导入和创建。

此时就会有个问题：如果数据源和仪表盘有几十个，逐一导入和创建势必会消耗大量人力，也无法实现自动话。

所以：有没有一种办法，可以在启动 Grafana 之前，就能直接加载这些数据呢？

Grafana 从 v5.0 版本中，决定通过一个 **Provisioning(配置供应系统)** 来解决上述问题。这个系统可以通过一系列的配置文件，让 Grafana 启动时加载他们，可以瞬间让启动好的 Grafana 就具有一定数量的数据源和仪表盘。这种行为使得 GitOps 更自然。这种思路除了可以用在数据源和仪表盘上以外，还可以扩展，比如提前配好用户信息、告警信息等等

# Data sources

该目录下的配置文件可以配置数据源的信息，当 Grafana 启动时，加载该目录下的 .yaml 文件，就会将数据源加载到 Grafana 中。

## 配置文件示例

```yaml
apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    url: http://monitor-bj-cs-k8s-prometheus:9090/
    access: proxy
    isDefault: true
    jsonData:
      timeInterval: 30s
    user: 访问 Prometheus 所使用的用户名
    secureJsonData:
      password: 访问 Prometheus 所使用的密码
```

# Plugins

> 参考：
>
> - [插件配置文件样例](https://grafana.com/docs/grafana/latest/administration/provisioning/#example-plugin-configuration-file)

注意：该功能只是可以配置插件的配置，而不是配置插件本身。使用此配置的前提是插件已经被安装在 Grafana 中。

# Dashboards

该目录下的配置文件将会指定一个**路径**，Grafana 启动时，会读取**该路径**下的所有 `*.json` 文件，并作为 Dashboard 加载到 Grafana 中。并且每隔一段时间就会检查路径下的文件，当文件有更新时，会同步更新加载到 Grafana 中的 Dashboard。

> 注意：目录下的 .json 文件就是在 Web 页面导出的 Dashboard。

**apiVersion: <INT>** # `默认值：1`
**providers: <\[]Object>** #

- **name: <STRING>** # an unique provider name. Required
- **orgId: 1** # Org 的 ID 号，`默认值：1`。通常 Grafana 启动后会自动创建一个名为 Main Org. 的 Org，该 Org 的 ID 为 1
- **folder: <STRING>** # 从目录读取到的所有仪表盘应该存放的文件夹。文件夹指的是 Grafana Web UI 上用于存放仪表盘的地方。若该值为空，则加载到的仪表盘存放在 General 文件夹中。
  - 注意：文件夹的名称与仪表盘的名称不能相同，否则将会报错并且无法自动生成仪表盘
- **folderUid: <STRING>** # 上面 folder 文件夹的 UID folder UID. will be automatically generated if not specified
- **type: <string>** # 提供者类型。默认值：file
- **disableDeletion: <bool>** # 是否允许通过 Web UI 删除目录下的仪表盘
- **updateIntervalSeconds: 10** # <int> Grafana 检查该目录下仪表盘是否有更新的间隔时间(单位：秒)。
- **allowUiUpdates: <bool>** # 是否允许通过 Web UI 更新目录下仪表盘
- **options: <Object>**
  - **path: <string>** # 必须的。要加载仪表盘的目录。该目录下的所有 .json 文件都会被 Grafana 加载为仪表盘
  - **foldersFromFilesStructure: <bool>** # 使用文件系统中的文件夹名称作为 Grafana Web UI 中的文件夹名。`默认值：false`。具体用法详见《[文件系统结构映射到 WebUI 中的文件夹](/docs/6.可观测性/Grafana/Grafana%20 配置详解/Provisioning%20 配置.md 配置.md)》
    - 注意：该字段与 `folder` 和 `folderUid` 冲突。

## 配置文件示例

加载 /etc/grafana/provisioning/dashboards/test 目录下所有 .json 文件为 Dashboard。

```yaml
apiVersion: 1
providers:
  - name: "sidecarProvider"
    orgId: 1
    folder: ""
    type: file
    disableDeletion: false
    allowUiUpdates: false
    updateIntervalSeconds: 30
    options:
      foldersFromFilesStructure: false
      path: /etc/grafana/provisioning/dashboards/custom
```

### 文件系统结构映射到 WebUI 中的文件夹

如果我们通过类似 git 或文件系统中的文件夹存储仪表盘的 JSON 文件，并且希望在 Grafana 的 Web UI 具有相同名称的文件夹，则可以使用这个字段。
比如我们有将仪表盘的 JSON 文件以如下结构保存：

```bash
/etc/dashboards
├── /server
│   ├── /common_dashboard.json
│   └── /network_dashboard.json
└── /application
    ├── /requests_dashboard.json
    └── /resources_dashboard.json
```

当我们使用如下配置文件时

```yaml
apiVersion: 1
providers:
  - name: dashboards
    type: file
    updateIntervalSeconds: 30
    options:
      path: /etc/dashboards
      foldersFromFilesStructure: true
```

Grafana 的 Web UI 中将会创建 `server` 与 `application` 两个文件夹，并将对应的仪表盘放在其中。

# Alert Notification Channels

- [Example Alert Notification Channels Config File](https://grafana.com/docs/grafana/latest/administration/provisioning/#example-alert-notification-channels-config-file)
- [Supported Settings](https://grafana.com/docs/grafana/latest/administration/provisioning/#supported-settings)
