---
title: Nginx 作为静态资源服务器配置示例
---

# Nginx 的 location 里面的 root、alias 的使用技巧与区别

> 知乎：<https://zhuanlan.zhihu.com/p/379076598>

## location

Nginx 里面的 location，可以针对一个特殊的 URI 路径进行单独的设置。

```javascript
location / {
    root /tongfu.net/web/static;
}
```

在 location 块里面可以单独设置映射目录、重写逻辑、默认文档等等。

```javascript
location / {
    root /tongfu.net/web/download;
    index index.htm;
}

location ~ ^\/download\/.*\.(zip|rar|tgz|gz)$ {
    rewrite ^\/download\/(.*)$ /downloadValidation.php?$1;
}
```

## root

Nginx 里面的 root 参数用来指定映射根目录，末尾不加“/”。

### 主机默认目录

直接在 server 里面设置 root 就是设置主机的根目录。

```javascript
server {
    root /tongfu.net/web/static;
}
```

### 匹配 URI 目录

在 location 里面设置 root 就是设置匹配 URI 的根目录。

下面的例子里如果访问 [http://localhost/icon/abc.png](https://link.zhihu.com/?target=http%3A//localhost/icon/abc.png) 网址，映射到的服务器路径是 /tongfu.net/web/icons**/icon/abc.png**。

```javascript
location /icon/ {
    root /tongfu.net/web/icons;
}
```

## alias

Nginx 里面的 root 参数用来指定映射目录，末尾需要加“/”。

下面的例子里如果访问 [http://localhost/icon/abc.png](https://link.zhihu.com/?target=http%3A//localhost/icon/abc.png) 网址，映射到的服务器路径是 /tongfu.net/web/icons/**abc.png**。

```javascript
location /icon/ {
    alias /tongfu.net/web/icons/;
}
```

# Nginx 代理 Vue 项目的单页应用刷新后 404 问题

在 location 中添加 `try_files $uri $uri/ /index.html;`
