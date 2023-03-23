---
title: Plugin-sftp
---

## 默认配置

```json
{
  "name": "ansible",
  "host": "172.38.40.250",
  "protocol": "sftp",
  "port": 22,
  "username": "root",
  "password": "XXXXXX",
  "remotePath": "/root/projects",
  "uploadOnSave": true,
  "ignore": [".vscode", ".git", ".DS_Store"]
}
```

## [**配置代理连接目标服务器**](https://github.com/liximomo/vscode-sftp#connection-hopping)

```json
{
  "name": "lichenhao",
  "remotePath": "/root/projects",
  // 用于作为代理的服务器信息
  "host": "202.43.145.163",
  "protocol": "sftp",
  "port": 42203,
  "username": "root",
  "password": "XXXXX",
  // 最终目标服务器信息
  "hop": {
    "host": "172.19.42.248",
    "port": 22,
    "username": "root",
    "password": "XXXXX"
  },
  "uploadOnSave": true,
  "ignore": [".vscode", ".git", ".DS_Store"]
}
```

## 多服务器配置

参考：[**https://blog.csdn.net/u012560340/article/details/83030680**](https://blog.csdn.net/u012560340/article/details/83030680)

在使用服务器时，需要打开 sftp 的 profiles 配置，以指定一个服务器

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ecg7u5/31208a5d6a46a3729daed9ddf2940eadedc2)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/ecg7u5/312023b203ab9129776eb25467fb952a8725)

```json
{
  "remotePath": "/root/kubernetes/HelmRepo/rabbitmq-stack-allinone",
  // profiles 外部的信息将会作用于每一个 profile
  "username": "root",
  "password": "XXXXXX",
  "profiles": {
    // 服务器一
    "lichenhao": {
      "host": "202.43.145.163",
      "port": 42203,
      // 跳转服务器无法享受 profiles 的全局信息
      "hop": {
        "host": "172.19.42.248",
        "port": 22,
        "username": "root",
        "password": "XXXXX"
      }
    },
    // 服务器二
    "ansbile": {
      "host": "172.38.40.250",
      "port": 22
    }
  },

  "uploadOnSave": true,
  "ignore": [".vscode", ".git", ".DS_Store"]
}
```
