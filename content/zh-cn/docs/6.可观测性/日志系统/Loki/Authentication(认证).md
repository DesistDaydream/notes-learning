---
title: Authentication(认证)
---

# 概述

> 参考：
> 
> - [官方文档,运行方式-认证](https://grafana.com/docs/loki/latest/operations/authentication/)

Loki 不附带任何包含的身份验证层。运营商应在您的服务之前运行身份验证反向代理，例如使用基本身份验证或 OAuth2 代理的 NGINX。

请注意，在多租户模式下使用 Loki 时，Loki 要求将 HTTP 标头 `X-Scope-OrgID`设置为标识租户的字符串。填充此值的责任应由身份验证反向代理处理。阅读[多租户](https://grafana.com/docs/loki/latest/operations/multi-tenancy/)文档以了解更多信息。

有关身份验证 Promtail 的信息，请参阅文档以[了解如何配置 Promtail](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)。
