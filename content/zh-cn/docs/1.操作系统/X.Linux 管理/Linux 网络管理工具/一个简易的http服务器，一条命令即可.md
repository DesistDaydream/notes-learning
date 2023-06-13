---
title: 一个简易的http服务器，一条命令即可
---


```
# 使用该命令可以在当前目录搭建一个简易的http服务器，当client访问的时候，就可以直接看到该目录下的内容，还可以下载该目录下的内容
python -m SimpleHTTPServer NUM
```

若报错则使用如下命令：

```bash
python3 -m http.server NUM
```
