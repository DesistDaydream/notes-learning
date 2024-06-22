---
title: OpenAPI 格式
---

参考：[官方 v3.0.3 版本文档](http://spec.openapis.org/oas/v3.0.3)

# 根字段

参考：[官方文档](http://spec.openapis.org/oas/v3.0.3#fixed-fields)

- **openapi: \<STRING>** # **必须的**。This string MUST be the semantic version number of the OpenAPI Specification version that the OpenAPI document uses. The openapi field SHOULD be used by tooling specifications and clients to interpret the OpenAPI document. This is not related to the API info.version string.
- **info: \<Object>** # **必须的**。Provides metadata about the API. The metadata MAY be used by tooling as required.
- **servers([]Object)** # 一组可用的服务器，用于提供到目标服务器的连接信息。当我们测试 API 时，将会连接其中一个。如果未提供 servers 属性或为空数组，则默认值为 URL 值为/的服务器对象。
- **paths: \<Object> **#** 必须的**。API 的可用路径及其操作。**文档中内容最多的字段**
- **components: \<Object>** # An element to hold various schemas for the specification.
- **security([]Object)** # A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition. To make security optional, an empty security requirement ({}) can be included in the array.
- **tags([]Object)** # A list of tags used by the specification with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools’ logic. Each tag name in the list MUST be unique.
- **externalDocs: \<Object>** # Additional external documentation.

## 根字段示例

```yaml
openapi: 3.0.1
info:
  title: Prometheus
  description: ""
  version: v1
tags:
  - name: 查询接口
  - name: Prometheus 信息接口
servers:
  - url: "http://test-prometheus.desistdaydream.ltd/api/v1"
    description: "北京测试Prom"
path: ......内容非常多，详见下文单独示例
components:
  schemas: {}
```

# server 字段

- **url: \<STRING>** # **必须的**。A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
- **description: \<STRING>** # An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
- **variables(map\[string]Object)** # A map between a variable name and its value. The value is used for substitution in the server’s URL template.

# paths 字段
