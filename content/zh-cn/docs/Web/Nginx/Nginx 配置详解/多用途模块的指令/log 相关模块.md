---
title: log 相关模块
linkTitle: log 相关模块
weight: 20
---

# 概述

> 参考：
>
> - [org 官方文档，http - ngx_http_log_module](https://nginx.org/en/docs/http/ngx_http_log_module.html)
> - [org 官方文档，stream - ngx_stream_log_modeule](https://nginx.org/en/docs/stream/ngx_stream_log_module.html)

# access_log

# log_format

# 日志格式示例

```nginx
{
    "args": "$args",
    "body_bytes_sent": "$body_bytes_sent",
    "bytes_sent": "$bytes_sent",
    "connection_requests": "$connection_requests",
    "geoip2_city": "$geoip2_city",
    "geoip2_city_country_name": "$geoip2_city_country_name",
    "geoip2_latitude": "$geoip2_latitude",
    "geoip2_longitude": "$geoip2_longitude",
    "geoip2_region_name": "$geoip2_region_name",
    "http_host": "$http_host",
    "http_user_agent": "$http_user_agent",
    "http_x_forwarded_for": "$http_x_forwarded_for",
    "remote_addr": "$remote_addr",
    "remote_port": "$remote_port",
    "request": "$request",
    "request_uri": "$request_uri",
    "request_time": "$request_time",
    "request_method": "$request_method",
    "scheme": "$scheme",
    "server_name": "$server_name",
    "server_protocol": "$server_protocol",
    "ssl_protocol": "$ssl_protocol",
    "ssl_cipher": "$ssl_cipher",
    "status": "$status",
    "time_iso8601": "$time_iso8601",
    "upstream": "$upstream_addr",
    "upstream_connect_time": "$upstream_connect_time",
    "upstream_response_time": "$upstream_response_time"
}
```

```nginx
log_format promtail_json '{"@timestamp":"$time_iso8601",'
        '"@version":"Promtail json",'
        '"server_addr":"$server_addr",'
        '"remote_addr":"$remote_addr",'
        '"host":"$host",'
        '"uri":"$uri",'
        '"body_bytes_sent":$body_bytes_sent,'
        '"bytes_sent":$body_bytes_sent,'
        '"request":"$request",'
        '"request_length":$request_length,'
        '"request_time":$request_time,'
        '"status":"$status",'
        '"http_referer":"$http_referer",'
        '"http_user_agent":"$http_user_agent"'
        '}';
```
