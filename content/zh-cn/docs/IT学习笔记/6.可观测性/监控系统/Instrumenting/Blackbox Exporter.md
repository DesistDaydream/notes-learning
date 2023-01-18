---
title: Blackbox Exporter
---

# 概述

> 参考：
> - [GitHub 项目，prometheus/blackbox_exporter](https://github.com/prometheus/blackbox_exporter)
> - [官方文档](https://prometheus.io/docs/guides/multi-target-exporter/#configuring-modules)
> - 个人文章参考：
>   - <https://mp.weixin.qq.com/s/gBdOMob_GZ5t44evAHFVOA>

我们可以使用如下几种协议来对目标进行探测

- http
- tcp
- dns
- icmp

## 使用方法

```bash
curl 'http://10.244.1.26:19115/probe?module=http_2xx&target=www.baidu.com'
```

## Prometheus 使用 Blackbox Exporter 的配置示例

与一般 Exporter 配置不同， Blackbox Exporter 的配置方式与 [SNMP Exporter](/docs/IT学习笔记/6.可观测性/监控系统/Instrumenting/SNMP%20Exporter.md Exporter.md) 更像，每一个待探测的目标将会作为 Blackbox Exporter 程序的参数。可以通过 Relabel 机制，设置目标的 instance 标签。

```yaml
scrape_configs:
  - job_name: "blackbox-http-get"
    metrics_path: /probe
    params:
      module: [http_2xx] # Look for a HTTP 200 response.
    static_configs:
      - targets:
          - http://prometheus.io # Target to probe with http.
          - https://prometheus.io # Target to probe with https.
          - http://example.com:8080 # Target to probe with http on port 8080.
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:9115 # The blackbox exporter's real hostname:port.
```

# 配置详解

> 参考：
> - [GitHub,CONFIGURATION.md](https://github.com/prometheus/blackbox_exporter/blob/master/CONFIGURATION.md)

Blackbox Exporter 的配置以模块区分，每个模块都有其独立的配置字段。一个模块就代表了一种探针类型及其探测行为。

```yaml
modules:
  # 指定模块的名称,可以定义多个模块。一个模块就代表一种探针及其探测行为。
  NAME_1:
    # 该模块的探针要使用的探测协议，可用的值有 http、tcp、dns、icmp
    prober: <STRING>
    # 探测时的超时时长
    timeout: <DURATION>
    # 探针的行为，ProberProtool 应为 prober 字段的值
    ProberProtocol: ...... #  不同的探测协议，可用的配置字段各不相同，假如使用 http，则 ProberProtocol 替换为 http。详见下文单独章节
```

## http 协议探针

**valid_status_codes: <\[]INT>** # 此探针可以接受的响应状态码。`默认值：2xx`。注:2xx 表示所有 2xx 状态码，这个字段的值如果要手动指定，必须是 int 类型。

- 若响应码不在该字段指定的范围内，则探测失败

**valid_http_versions: <STRING>** # 探针接受的 HTTP 版本。

- 若 HTTP 版本不在字段指定的范围内，则探测失败

**method: <STRING>** # 探针探测是要使用的 HTTP Method。`默认值：GET`
**headers: \<map\[STGRING]STRING>** # 设置探测时要使用的 Header，每行都是一个请求头的键值对。
**compression: <STRING> **#&#x20;
**follow_redirects: <BOOLEAN>** #&#x20;
**fail_if_ssl: <BOOLEAN>** # 如果 SSL 存在，则探针失败。`默认值：false`
**fail_if_not_ssl: <BOOLEAN>** # 如果 SSL 不存在，则探针失败。`默认值：false`
**fail_if_body_matches_regexp:** # Probe fails if response body matches regex.
\[ - <regex>, ... ]
**fail_if_body_not_matches_regexp: **# Probe fails if response body does not match regex.
\[ - <regex>, ... ]
**fail_if_header_matches: **# Probe fails if response header matches regex. For headers with multiple values, fails if _at least one_ matches.
\[ - \<http*header_match_spec>, ... ]
**fail_if_header_not_matches:** # Probe fails if response header does not match regex. For headers with multiple values, fails if \_none* match.
\[ - \<http_header_match_spec>, ... ]

\######## Prometheus [共享库中的通用 HTTP 客户端配置](https://github.com/prometheus/common/blob/v0.30.0/config/http_config.go#L159) ########
**basic_auth: <OBJECT>** # 配置 HTTP 的基础认证信息。

- **username: <STRING>** #
- **password: <STRING>** #
- **password_file: <STRING>** #

**bearer_token: <SECRET>** # 探测目标时要使用的 bearer 令牌
**bearer_token_file: <filename>** # 探测目标时要使用的 bearer 令牌文件
**oauth2: <Object>** # 配置 OAuth 2.0 的认证配置。与 basic_auth 和 authorization 两个字段互斥
**proxy_url: <STRING>** # HTTP proxy server to use to connect to the targets.
**tls_config: <OBJECT>** # 发起 HTTP 请求时的 TLS 配置，即发起 HTTPS 请求。
&#x20; 详见 [tls 配置段](#b9c06c74)
\######## Prometheus [共享库中的通用 HTTP 客户端配置](https://github.com/prometheus/common/blob/v0.30.0/config/http_config.go#L159)结束 ########

**preferred_ip_protocol: <STRING>** # 探针首选的 IP 协议版本。`默认值：ip6`
**ip_protocol_fallback: <BOOLEAN>** # 。`默认值：true`
**body: <STRING>** # 探测时要携带的 HTTP Body

## tcp 协议探针

```yaml
# The IP protocol of the TCP probe (ip4, ip6).
[ preferred_ip_protocol: <string> | default = "ip6" ]
[ ip_protocol_fallback: <boolean | default = true> ]

# The source IP address.
[ source_ip_address: <string> ]

# The query sent in the TCP probe and the expected associated response.
# starttls upgrades TCP connection to TLS.
query_response:
  [ - [ [ expect: <string> ],
        [ send: <string> ],
        [ starttls: <boolean | default = false> ]
      ], ...
  ]

# Whether or not TLS is used when the connection is initiated.
[ tls: <boolean | default = false> ]

# Configuration for TLS protocol of TCP probe.
tls_config:
  [ <tls_config> ]
```

## dns 协议探针

```yaml
# The IP protocol of the DNS probe (ip4, ip6).
[ preferred_ip_protocol: <string> | default = "ip6" ]
[ ip_protocol_fallback: <boolean | default = true> ]

# The source IP address.
[ source_ip_address: <string> ]

[ transport_protocol: <string> | default = "udp" ] # udp, tcp

# Whether to use DNS over TLS. This only works with TCP.
[ dns_over_tls: <boolean | default = false> ]

# Configuration for TLS protocol of DNS over TLS probe.
tls_config:
  [ <tls_config> ]

query_name: <string>

[ query_type: <string> | default = "ANY" ]
[ query_class: <string> | default = "IN" ]

# List of valid response codes.
valid_rcodes:
  [ - <string> ... | default = "NOERROR" ]

validate_answer_rrs:

  fail_if_matches_regexp:
    [ - <regex>, ... ]

  fail_if_all_match_regexp:
    [ - <regex>, ... ]

  fail_if_not_matches_regexp:
    [ - <regex>, ... ]

  fail_if_none_matches_regexp:
    [ - <regex>, ... ]

validate_authority_rrs:

  fail_if_matches_regexp:
    [ - <regex>, ... ]

  fail_if_all_match_regexp:
    [ - <regex>, ... ]

  fail_if_not_matches_regexp:
    [ - <regex>, ... ]

  fail_if_none_matches_regexp:
    [ - <regex>, ... ]

validate_additional_rrs:

  fail_if_matches_regexp:
    [ - <regex>, ... ]

  fail_if_all_match_regexp:
    [ - <regex>, ... ]

  fail_if_not_matches_regexp:
    [ - <regex>, ... ]

  fail_if_none_matches_regexp:
    [ - <regex>, ... ]
```

## icmp 协议探针

```yaml
# The IP protocol of the ICMP probe (ip4, ip6).
[ preferred_ip_protocol: <string> | default = "ip6" ]
[ ip_protocol_fallback: <boolean | default = true> ]

# The source IP address.
[ source_ip_address: <string> ]

# Set the DF-bit in the IP-header. Only works with ip4, on *nix systems and
# requires raw sockets (i.e. root or CAP_NET_RAW on Linux).
[ dont_fragment: <boolean> | default = false ]

# The size of the payload.
[ payload_size: <int> ]
```

## 通用配置

## tls 配置段

可以为多种协议的探针配置，用来配置安全相关信息。

```yaml
# 禁用目标证书认证。默认值：false
insecure_skip_verify: <BOOLEAN>
# The CA cert to use for the targets.
[ ca_file: <filename> ]
# The client cert file for the targets.
[ cert_file: <filename> ]
# The client key file for the targets.
[ key_file: <filename> ]
# Used to verify the hostname for the targets.
[ server_name: <string> ]
```

# 配置示例

```yaml
modules:
  http_2xx:
    prober: http
    http:
      # valid_http_versions: ["HTTP/1.1", "HTTP/2"]
      # valid_status_codes: [200]
      method: GET
      preferred_ip_protocol: "ip4"
      tls_config:
        insecure_skip_verify: true
  http_post_2xx:
    prober: http
    timeout: 10s
    http:
      valid_http_versions: ["HTTP/1.1", "HTTP/2"]
      method: POST
      preferred_ip_protocol: "ip4"
      tls_config:
        insecure_skip_verify: true
  tcp_connect:
    prober: tcp
    timeout: 10s
  dns:
    prober: dns
    dns:
      transport_protocol: "tcp"
      preferred_ip_protocol: "ip4"
      query_name: "kubernetes.default.svc.cluster.local"
  icmp:
    prober: icmp
```
