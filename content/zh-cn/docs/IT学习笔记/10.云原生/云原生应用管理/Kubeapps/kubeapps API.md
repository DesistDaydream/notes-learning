---
title: kubeapps API
---

参考：[官方 OpenAPI 3.0](https://github.com/kubeapps/kubeapps/blob/master/dashboard/public/openapi.yaml)

云原生应用管理
API.json(118.2 KB)
&#x20;\- 0%

## Header

请求头必须包含的信息：

|               |                 |
| ------------- | --------------- |
| key           | value           |
| Authorization | Bearer ${TOKEN} |

${TOKEN} 通过如下命令获取,若已经创建 RBAC，则直接执行最后一条获取 TOKEN 即可。

- **kubectl create serviceaccount kubeapps-operator**

- **kubectl create clusterrolebinding kubeapps-operator --clusterrole=cluster-admin --serviceaccount=default:kubeapps-operator**

- **kubectl get secret $(kubectl get serviceaccounts kubeapps-operator -ojsonpath='{.secrets\[0].name}') -ojsonpath='{.data.token}' | base64 -d && echo**

# Release 相关

通过 kubeops 获取 Release 信息。基础 URL：http://IP:PORT/v1

### GET /clusters/{cluster}/releases 列出所有 Releases

## 限定名称空间的 API

### GET /clusters/{cluster}/namespaces/{namespace}/releases 列出所有 Releases

### n

### GET /clusters/{cluster}/namespaces/{namespace}/releases/{releaseName} 获取 Release 的信息

包含 chart 的信息

### PUT /clusters/{cluster}/namespaces/{namespace}/releases/{releaseName} 更新 Release

### DELETE /clusters/{cluster}/namespaces/{namespace}/releases/{releaseName} 删除 Release

- Query String Parameters:

- purge: true # 如果不填，那么仅仅删除了元数据，只是无法通过 helm ls -A 命令列出，但是该 release 对应的 secret 还存在，只不过

# Repoistory 相关

通过 kubeops 获取 Repoistory 信息。基础 URL：http://IP:PORT/backend/v1

### POST /clusters/{cluster}/can-i 不知道干啥用的。。。。

### GET /clusters/{cluster}/namespaces 获取所有名称空间信息

### GET /clusters/{cluster}/apprepositories 列出所有 repo 的信息

## 限定名称空间下的 API

### GET /clusters/{cluster}/namespaces/{namespace}/apprepositories 列出所有 repo 的信息

### POST /clusters/{cluster}/namespaces/{namespace}/apprepositories 创建 Repo

### POST /clusters/{cluster}/namespaces/{namespace}/apprepositories/validate 验证 Repo

### PUT /clusters/{cluster}/namespaces/{namespace}/apprepositories/{name} 更新 Repo

### POST /clusters/{cluster}/namespaces/{namespace}/apprepositories/{name}/refresh 刷新 repo

### DELETE /clusters/{cluster}/namespaces/{namespace}/apprepositories/{name} 删除 Repo

### GET /clusters/{cluster}/namespaces/{namespace}/operator/{name}/logo 获取 operator 的 logo

# Charts 相关

通过 assetsvc 获取 charts 信息，基础 URL：http://IP:PORT/v1

> 必须指定 namespace，

### GET /clusters/{cluster}/namespaces/{namespace}/charts" 列出所有 charts

Queries("name", "{chartName}", "version", "{version}", "appversion", "{appversion}").Handler(WithParams(listChartsWithFilters))

- Query String Parameters:

- page=xx\&size=yy

apiv1.Methods("GET").Path("/clusters/{cluster}/namespaces/{namespace}/charts/categories").Handler(WithParams(getChartCategories))

## 限定 Repo 的 API

### GET /clusters/{cluster}/namespaces/{namespace}/charts/{repo} 列出所有 charts

- Query String Parameters:

- page=xx\&size=yy

apiv1.Methods("GET").Path("/clusters/{cluster}/namespaces/{namespace}/charts/{repo}/categories").Handler(WithParams(getChartCategories))

### GET /clusters/{cluster}/namespaces/{namespace}/charts/{repo}/{chartName} # 获取 Chart 的信息

### GET /clusters/{cluster}/namespaces/{namespace}/charts/{repo}/{chartName}/versions 列出 Chart 的所有版本

### GET /clusters/{cluster}/namespaces/{namespace}/charts/{repo}/{chartName}/versions/{version} # 获取 Chart 指定版本的信息

### GET /clusters/{cluster}/namespaces/{namespace}/assets/{repo}/{chartName}/logo 获取 Chart 的 log

### GET /clusters/{cluster}/namespaces/{namespace}/assets/{repo}/{chartName}/versions/{version}/README.md 获取 Chart 指定版本的 README

### GET /clusters/{cluster}/namespaces/{namespace}/assets/{repo}/{chartName}/versions/{version}/values.yaml 获取 Chart 指定版本的 values.yaml 文件

apiv1.Methods("GET").Path("/clusters/{cluster}/namespaces/{namespace}/assets/{repo}/{chartName}/versions/{version}/values.schema.json").Handler(WithParams(getChartVersionSchema))
