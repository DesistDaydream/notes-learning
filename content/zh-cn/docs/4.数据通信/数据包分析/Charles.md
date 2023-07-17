---
title: Charles
---

# 概述

> 参考：
> - 官网：<https://www.charlesproxy.com/>
> - 小米手机安装 Charles 证书：<https://blog.csdn.net/yang450712123/article/details/112908643>
>   - 安卓用不了 2022.9.19
> - IOS 安装 Charles 证书：<https://www.jianshu.com/p/08f602eabb54>
>   - 苹果的能用 2022.9.19
> - <https://www.charles.ren/> 生成注册码
>   - 生成码之后直接使用即可

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/arzv8v/1671955685091-ae697c1c-96a5-4d8c-8b3c-e47da76fc75e.png)

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/arzv8v/1671955755610-395fe14a-f08d-42a7-877c-7296c96e473f.png)

在手机、pad 上配置 WLAN 代理，访问 `chls.pro/ssl` 下载证书。

打开手机设置，搜索：加密与凭据 => 安装证书 => 证书

# IOS 安装证书

为无线连接配置手动代理

IOS 访问 chls.pro/ssl 下载证书并安装

设置 —— 通用 —— 关于本机 —— 证书信任设置，开启信任证书

# 安卓安装证书

## 安卓安装 Charles 证书无效

证书安装成功，但是抓到的包都是 unknow，可能的原因：

- Android7.0 之后默认不信任用户级别 CA 证书
- 此时开启抓包后，很多 APP 都是无网络的情况；但是 chrome 打开网页是可以抓到 https 的包
- 需要想办法安装在系统级别下的 CA 证书
- 可能的方法
  - 平行空间
  - 获取系统 ROot 权限

HttpCanary根证书安装(MIUI13 Android 12可用) https://czyx007.cn/archives/httpcanary

- 还是没法抓集换社的包，微信小程序的包也抓不到

# 微信小程序抓包

到底怎么抓呢？~

[知乎上的一篇文章](https://www.zhihu.com/question/350183786/answer/2487803703)中说删除 PC 小程序缓存在 2022 年 9 月 5 号之后也不管用了
给微信降级后即可在 PC 上抓包：<https://blog.csdn.net/weixin_46552558/article/details/124037807>

小程序内嵌的 h5 调用如何抓到？

