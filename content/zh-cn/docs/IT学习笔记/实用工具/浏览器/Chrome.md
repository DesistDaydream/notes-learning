---
title: "Chrome"
linkTitle: "Chrome"
weight: 20
---

# 概述
> 参考：
> -


# 如何避免打开谷歌自动跳转到香港 GOOGLE.COM.HK？

原文链接：<https://www.jianshu.com/p/c00e35ec4c3e>

> 自从 google 的服务器搬离中国大陆后，大陆地区用户用 google 服务时会自动跳转到香港的 [**http://google.com.hk**](https://link.jianshu.com?t=http://google.com.hk) ，有关键字过滤而且偶尔不是很稳定，这对我们的生活工作都造成了困扰。

#### 一、可以通过以下的方法访问 http://google.com

- 直接用 [**http://www.google.com/ncr**](http://www.google.com/ncr) ，`ncr` 是 `no country redirection` ，是一个强制不跳转的命令；
- 用 [**https://www.google.com/**](https://link.jianshu.com?t=https://www.google.com/) ，`https` 协议。

#### 二、另外一个问题是 Chrome 浏览器的默认搜索也是设置为 http://www.google.com.hk/ ，我们可以自行修改一下。

- Chrome – 设置 -搜索 - 管理搜索引擎 – 其他搜索引擎
- 拉到最下，有一个“添加”
- 名字：自己写，我写 [**http://google.com**](https://link.jianshu.com?t=http://google.com)
- 关键字（keyword），我写 G
- 最后一个空最重要，写入 Url ( [**http://www.google.com/search?hl=zh-CN\&q=%s**](https://link.jianshu.com?t=http://www.google.com/search?hl=zh-CN&q=%s)) 或者 ( [**http://www.google.com/search?q=%s**](https://link.jianshu.com?t=http://www.google.com/search?q=%s) ) `括号为填写部分`
- 然后将之设置成默认搜索引擎，搞定！


> **so easy！好好享受 google** 原汁原味的搜索吧！

# Software Report Tool 进程占用大量 CPU 时间

**Software Reporter Tool 其实就是 Chrome 的清理工具**，用于清理谷歌浏览器中不必要或恶意的扩展，应用程序，劫持开始页面等等。当你安装 Chrome 时，Software_reporter_tool.exe 也就会被下载在 SwReporter 文件夹下的 Chrome 应用数据文件夹中。这个软件在运行的过程中可能会**长时间地占用 CPU**，导致高 CPU 使用率。我们虽然可以通过任务管理器手动结束进程或者选择删除 SRT，但这都不是长久的解决办法。因为前者过一段时间它又会再次运行，后者在浏览器更新的时候就又会重新被下载下来。查询该文件，发现 software_reporter_tool.exe 会扫描系统，类似 Chrome 的一个计划任务，每周启动扫描一次，运行大约 20-25 分钟。

修改 Chrome 设置

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/yg65yy/1632309843432-ea1d257f-b402-4181-b0c3-d36c32c32eeb.png)
若修改没用，则找到 `C:\Users\DesistDaydream\AppData\Local\Google\Chrome\User Data\SwReporter\93.269.200\manifest.json` 文件。将 allow-reporter-logs 的值改为 false

```json
{
  "launch_params": [
    {
      "allow-reporter-logs": false,
      "arguments": [
        "--engine=2",
        "--scan-locations=1,2,3,4,5,6,7,8,10",
        "--disabled-locations=9,11"
      ],
      "prompt": true,
      "suffix": "ESET"
    }
  ],
  "manifest_version": 2,
  "name": "Software Reporter Tool",
  "version": "93.269.200"
}
```

# 滚动条冻结问题

> <https://leadscloud.github.io/314048/chrome%E6%B5%8F%E8%A7%88%E5%99%A8%E6%BB%9A%E5%8A%A8%E6%9D%A1%E5%86%BB%E7%BB%93%E9%97%AE%E9%A2%98/>

chrome://flags/#smooth-scrolling # 将该参数改为 disabled

# Chrome 打包插件

C:\Users\DesistDaydream\AppData\Local\Google\Chrome\User Data\Profile 2\Extensions


