---
title: Nginx 作为代理服务器配置示例
---

# 概述

> 参考：
> - [公众号，Nginx 代理 WebSocket 方法](https://mp.weixin.qq.com/s/27IuQAe8UZGXIdNApE2Ljg)

# 7 层代理配置

这个配置里的 172.19.42.217 是 kubernetes 集群的入口，一般在 80 和 443 上都起一个 ingress controler，这样，多种域名都代理到同一个 kubernetes 集群，然后由 ingress 再将流量进行路由分配。

```nginx
user  nginx;
worker_processes  4;

error_log  /dev/stdout warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  102400;
}


http {
  default_type  application/octet-stream;

  access_log /dev/stdout main;
  keepalive_timeout  120;
  log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
    '$status $body_bytes_sent "$http_referer" '
    '"$http_user_agent" "$http_x_forwarded_for"'
    '$upstream_addr '
    'ups_resp_time: $upstream_response_time '
    'request_time: $request_time';
  sendfile on;
  server_names_hash_bucket_size 256;

  server {
    listen       80;
    server_name  grafana.desistdaydream.ltd;
    server_name  prometheus.desistdaydream.ltd;
    server_name  desistdaydream.ltd;
    server_name  www.desistdaydream.ltd;

    client_body_in_file_only clean;
    client_body_buffer_size 64K;
    client_max_body_size 40M;
    sendfile on;
    send_timeout 300s;

    location / {
      proxy_pass http://172.19.42.217/;
      proxy_set_header X-Real-IP $remote_addr;
      proxy_set_header Host $host;
      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header X-Forwarded-Proto $scheme;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
		}
	}

  include       /etc/nginx/mime.types;
  include /etc/nginx/conf.d/*.conf;
  include /etc/nginx/conf.d/protal/*.conf;
}
```

# 4 层代理配置

```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

include /usr/share/nginx/modules/*.conf;

events {
    worker_connections 10240;
}

stream {
    include stream.d/*.conf;

    upstream grafana {
        server 172.38.40.216:30000;
        server 172.38.40.217:30000;
    }
    upstream prometheus {
        server 172.38.40.216:30001 weight=8 max_fails=2 fail_timeout=30s;
        server 172.38.40.217:30001 weight=8 max_fails=2 fail_timeout=30s;
    }

    server {
        listen 30000;
        proxy_pass grafana;
        proxy_connect_timeout 10s;
    }
    server {
        listen 30001;
        proxy_pass prometheus;
        proxy_connect_timeout 2s;
    }
}

http {

}
```

# 待整理配置示例

```nginx
server {
       listen       80;
       server_name  grafana.desistdaydream.ltd;
       server_name  prometheus.desistdaydream.ltd;
       server_name  desistdaydream.ltd;
       server_name  www.desistdaydream.ltd;

       client_body_in_file_only clean;
       client_body_buffer_size 64K;
       client_max_body_size 40M;
       sendfile on;
       send_timeout 300s;

       location / {
           proxy_pass http://172.19.42.217/;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header Host $host;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
           proxy_http_version 1.1;
           proxy_set_header Upgrade $http_upgrade;
           proxy_set_header Connection "upgrade";
       }
}
```

## https

```nginx
server {
       listen       80;
       listen       443 ssl;
       server_name  rancher.desistdaydream.ltd;

       ssl on;
       # crt证书
       ssl_certificate ../keys/bj/rancher.desistdaydream.ltd.crt;
       # key证书
       ssl_certificate_key ../keys/bj/rancher.desistdaydream.ltd.key;

       client_body_in_file_only clean;
       client_body_buffer_size 64K;
       client_max_body_size 40M;
       sendfile on;
       send_timeout 300s;

       location / {
           proxy_pass https://172.19.42.217:60443/;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header Host $host;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
           proxy_http_version 1.1;
           proxy_set_header Upgrade $http_upgrade;
           proxy_set_header Connection "upgrade";
       }
}
```
