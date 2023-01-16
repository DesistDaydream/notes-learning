---
title: Dockerfile 样例
---

生成一个具有基本网络工具的容器，可以用来测试

```dockerfile
FROM centos:centos7.8.2003
ENV LANG=zh_CN.UTF-8 \
LANGUAGE=zh_CN:zh \
LC_ALL=zh_CN.UTF-8
RUN yum install -y epel-release.noarch && \
    yum install -y iproute bind-utils nginx glibc-common tcpdump telnet && \
    yum clean all && \
rm -rf /tmp/* rm -rf /var/cache/yum/* && \
    localedef -c -f UTF-8 -i zh_CN zh_CN.UTF-8 && \
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY index.html /usr/share/nginx/html
EXPOSE 80/tcp
CMD ["/usr/sbin/nginx","-g","daemon off;"]
```

```bash
[root@lichenhao test]# cat index.html
<meta charset="utf-8"/>
<h1>网络测试容器 desist-daydream 1</h1>
```

使用 alpine 版本让镜像更小

```dockerfile
FROM nginx:stable-alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache busybox-extras curl iproute2 tcpdump
COPY index.html /usr/share/nginx/html
```

使用 multi-stage 功能构建镜像

```dockerfile
FROM golang
WORKDIR /root/caredaily
COPY caredaily/ ./
RUN go build .
FROM ubuntu
WORKDIR /root/caredaily
COPY --from=0 /root/caredaily/caredaily ./
COPY --from=0 /root/caredaily/templates/ ./templates/
CMD ["./caredaily"]
```

Dockerfile 样例

```dockerfile
FROM nginx:1.14-alpine
ENV NGX_DOC_ROOT="/data/web/html"
ADD entrypoint.sh /bin/
CMD ["/usr/sbin/nginx","-g","daemon off;"]
ENTRYPOINT ["/bin/sh","-c","/bin/entrypoint.sh"]
/entrypoint.sh
#!/bin/bash
cat > /etc/nginx/conf.d/www.conf << EOF
server {
        server_name ${HOSTNAME};
        listen ${IP:-0.0.0.0}:${PORT:-80};
        root ${NGX_DOC_ROOT:-/usr/share/nginx/html};
}
exec "$@"
```

```dockerfile
FROM busybox:lastest
ENV WEB_SERVER_PACKAGE nginx-1.15.2.tar.gz \
DOC_ROOT=/data/web/html/
ADD http://fnginx.org/download/${WEB_SERVER_PACKAGE} /usr/local/src
RUN cd /usr/local/src/ && \
tar -xf ${WEB_SERVER_PACKAGE} \
    yum install epel-release && yum makecahe && yum install -y nginx
COPY index.html ${DOC_ROOT}
VOLUME /data/mysql
EXPOSE 80/tcp
```
