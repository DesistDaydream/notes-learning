---
title: TCP/IP 管理工具
---

# 概述

# 获取本机公网 IP

## ipify

<https://geo.ipify.org/docs>

```go
func main() {
        res, _ := http.Get("https://api.ipify.org")
        ip, _ := ioutil.ReadAll(res.Body)
        os.Stdout.Write(ip)
}
```

## 其他

<https://ip.netarm.com>
<http://ip.cip.cc>
http://ip.sb
[
](https://geo.ipify.org/docs)
