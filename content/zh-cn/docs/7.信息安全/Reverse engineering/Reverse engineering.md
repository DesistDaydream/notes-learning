---
title: Reverse engineering
linkTitle: Reverse engineering
date: 2023-10-04T12:31
weight: 1
---

# 概述

> 参考：
>
> - [Wiki，Reverse_engineering](https://en.wikipedia.org/wiki/Reverse_engineering)
> - [Wiki，Software_cracking](https://en.wikipedia.org/wiki/Software_cracking)

**Reverse engineering(逆向工程)**

逆向工程的其他目的包括安全审核、消除复制保护（“破解”）、规避消费电子产品中常见的访问限制、定制嵌入式系统（例如引擎管理系统）、内部维修或改造、低成本“残缺”硬件（例如某些显卡芯片组）上的附加功能，甚至只是满足好奇心。

**Software cracking(软件破解)**

# 学习

[吾爱破解](https://www.52pojie.cn/)

- [吾爱破解安卓逆向入门教程《安卓逆向这档事》十五、是时候学习一下Frida一把梭了(下)](https://mp.weixin.qq.com/s/97o3fX9AN_kl2GCLhHAfig)

[图灵 Python](https://www.tulingpyton.cn/) 何老师相关课程

[李玺](https://blog.csdn.net/weixin_43582101)

- https://github.com/lixi5338619
- http://www.lxspider.com/ 博客
- https://github.com/lixi5338619/lxBook  《爬虫逆向进阶实战》
- https://space.bilibili.com/390499740/channel/collectiondetail?sid=468228&ctype=0 实战示例视频

# 常见安全策略

JS 混淆

- 在 Fetch/XHR 的请求中，从开发者工具查看该请求的启动器，如果函数名、变量名都是 `_0x5601f0` 这类以 `_0x` 开头的，说明代码是经过混淆的

数据 RSA、AES、etc. 加密

APP 加固

APP 加壳

# 逆向常用工具

**x64dbg** # 适用于 Windows 的开源用户模式调试器。针对逆向工程和恶意软件分析进行了优化。

- [GitHub 项目，x64dbg/x64dbg](https://github.com/x64dbg/x64dbg)
- https://x64dbg.com/
- [想要程序注册码（密码）？？翻找内存找到它！！——x64dbg逆向动态调试简单crackme，找注册码（密码）](https://www.bilibili.com/video/BV1it421G7qf)

**Cheat Engine** # Cheat Engine 是一个专注于修改个人使用的游戏和应用程序的开发环境。

- [GitHub 项目，cheat-engine/cheat-engine](https://github.com/cheat-engine/cheat-engine)

[Packet analyzer](/docs/7.信息安全/Packet%20analyzer/Packet%20analyzer.md)

http://www.cnlans.com/lx/tools 李玺的爬虫逆向工具

## 待整理工具

### 查壳反编译

- [查壳小工具](https://pan.baidu.com/s/1s1BoElAyQCnPaxb2T3QpEw?pwd=tmbs)
- [AndroidKiller](https://down.52pojie.cn/Tools/Android_Tools/)
- [Apktools2.5.0](https://pan.baidu.com/s/12qB4N_2Fg-IsTB2BcQuiDw?pwd=gjqs)
- [超级Jadx](https://pan.baidu.com/s/1SHsJGfnGJJmcPfgcC_lnYA?pwd=9999)
- [IDAPro7.0 调试工具Windows版本](https://pan.baidu.com/s/1_-PorRCwHDMpmUI1t_cKcQ?pwd=t39m)
- [ddms](https://pan.baidu.com/s/1wdsZvTA-fAZ12o53Exw80A?pwd=wk3d)
- [JEB3.0中文版](https://pan.baidu.com/s/1kCjw8dP9tq7kLBWkublHag?pwd=k2s4)
- [JEB2.3.7](https://pan.baidu.com/s/1HgyyEomL72jLWY1XMtHv8g?pwd=zpha)

### 脱壳工具

- [FDex2](https://pan.baidu.com/s/1e0zcp1IzA-u7UC-A3gaj8g?pwd=yds2)
- [反射大师](https://pan.baidu.com/s/170oS04qoFdd-Btu9DanHfg?pwd=an39)
- [BlackDex3.1.0](https://pan.baidu.com/s/18gijmyy5dgUCbwi-hnqtpg?pwd=433u)
- [DumpDex](https://github.com/WrBug/dumpDex)
- [FRIDA-DEXDump](https://github.com/hluwa/FRIDA-DEXDump)

### HOOK工具

- [Xposed](https://pan.baidu.com/s/15WnJD8qj9UzSss55DWLNfA?pwd=7sgb)
- [VAExposed](https://pan.baidu.com/s/1fd0r2fy4mm4jUArGE4MZvA?pwd=mu9q)
- [Inspeckage](https://pan.baidu.com/s/1WfnVM7hKE76jNpQc3FnKWg?pwd=pvcs)
- [SSLUnpinning 20](https://pan.baidu.com/s/1EZuv-JK0a-TLHhw4v6SkvQ?pwd=dsfj)

### 抓包工具

- [httpCanary 安卓抓包工具](https://pan.baidu.com/s/1mdHHaXulnsM6Zxf335yMHA?pwd=tfhx)
- [Postern安卓抓包工具](https://pan.baidu.com/s/1A-2kIVnYSxpgHqiDn4mqnw?pwd=1e5k)
- [Drony_113](https://pan.baidu.com/s/14d6ezZXRWDQayL73d2E8gw?pwd=tyk7)
- [HttpAnalyzerStd V7](https://pan.baidu.com/s/1p3ThL5yqqc5XwTrDdmmGCg?pwd=x9hg)
- [fridaManager](https://pan.baidu.com/s/1u_P2P_kd_H2n2SYTaLB0hA?pwd=jovi)
- [AppSignGet](https://pan.baidu.com/s/1_j2QTVFD6qHP3FKp_FVeCw?pwd=6qmu)

## chrome插件

- [request-hook](https://pan.baidu.com/s/1OmMiE4rJrTNwarw3EJbz0A?pwd=thyl)
- [Trace-dist](https://github.com/L018/Trace)

### 微信小程序

- [UnpackMiniApp](https://pan.baidu.com/s/1dwUehOAnPka9eHjXN6Y-Lg?pwd=unp7)
- [CrackMinApp](https://github.com/Cherrison/CrackMinApp)

# 其他常用工具

- [FontCreator英文版](https://pan.baidu.com/s/1Ek34ePZpJYTkmiCuKsqIMQ?pwd=hnku)
- [鬼鬼JS调试工具](https://pan.baidu.com/s/1hjdgx3DOTJMp0wtYGAa67A?pwd=1s67)
- [MT 管理器](https://pan.baidu.com/s/1AfBDHVvini4bweDOD9GoIw?pwd=9999)
- [NP 管理器](https://pan.baidu.com/s/1X5g8loORq_WS0HLqeasLbg?pwd=9jk7)
- [Autojs](https://pan.baidu.com/s/1bbjFWMjFU5m2RupRyIZcGw?pwd=4ikp)

# 待整理

验证码技术

- 腾讯防水墙
- 阿里无感 v3
- 极验4代点选
- 小红书 数美验证
- 百度旋转验证
- 易盾
- 顶象
- 瑞数

JSVMP 加密

# 小程序逆向

## 调试小程序

[Mobile app](/docs/Mobile%20device/Mobile%20app.md)

TODO: 其他待整理

## 反编译小程序

https://www.bilibili.com/video/BV1ew411K7nB?p=40

微信小程序目录位置

- PC
  - `WeChat Files\Applet\` # 该目录为小程序所在文件夹。每个小程序文件都是一个独立的文件夹，以 wx 开头，像 `wx64479c83c7630409` 这样
  - 想要找到对应的小程序，可以把所有 wx 开头的文件夹都删除，然后打开小程序，就会生成一个信息的 wx 开头的文件夹。

UnpackMiniApp # 解密 `*.wxapkg` 文件获得 `*.wxapkg` 文件。TODO: 为什么要先解密？这是加的什么密？

- 找不到官方下载渠道

Unveilr # 反编译解密后的 `*.wxapkg` 文件得到源码。

- https://github.com/r3x5ur/unveilr # 好像是原始源码，但是 2.0 之后的版本作者收费了
  - 只有个下载地址和 TG 号 https://t.me/Qobg3fbwQM1hNTY1
  - https://u.openal.lat/ 提供下载和花钱买 token
  - 下面是一些 2.0.2 版本源码的备份
    - https://github.com/AnkioTomas/unveilr
    - https://github.com/CoderYiXin/unveilr

使用微信开发者工具打开项目

其他反编译小程序的项目

- https://github.com/wux1an/wxapkg
- https://github.com/zgqwork/wxapkg-unpacker
  - 基于 https://github.com/qwerty472123/wxappUnpacker, 该项目已于 2021 年删库归档

# Web 逆向

从 [加密解密的最佳实践](/docs/7.信息安全/Cryptography/加密解密的最佳实践.md) 中可以看到常见的加密/解密方式。

很多时候我们无法用其他语言实现找到的 js 代码，但是又想要使用 Python 怎么办呢，可以使用一些第三方库，以便让 Python 可以执行 JS 代码（e.g. pyexecjs、js2py）

在 [Programming Technology](/docs/2.编程/Programming%20technology/Programming%20Technology.md) 可以知道网站的数据有动态和静态两种。

- 若是静态的数据，在 “[DevTools](/docs/Web/Browser/DevTools.md) - 网络 - 文档” 中查看数据资源
  - 静态数据通常是直接返回 HTML 页面，此时我们可以直接使用各种语言的 DOM 树管理库，通过 XPath 等方式定位元素，以获取其中的数据
  - 这类网站有的时候有个特点，一个页面有需要两次请求，第一次返回一段 js 代码，然后生成 cookie，第二次带着 cookie 发起请求再获取到静态 HTML 数据。
    - 这两次请求需要在 开发者工具 - 网络 中打开 “保留日志” 功能才可以看到第一次。
  - https://www.bilibili.com/video/BV1ew411K7nB?p=19 这里有介绍。
- 若是动态的数据，在 “[DevTools](/docs/Web/Browser/DevTools.md) - 网络 - Fetch/XHR” 中查看数据资源

扣 js 代码。主要是找到加密/解密相关的 JS 代码，找到后大段大段得放到本地，给定已加密的数据，可以正常解密就算成功。绝大部分都是扣的函数，然后传密文（密文响应体能直接看到）进去返回明文。而 JS 代码又可以被 Python 执行。

- 第一步扣下来的代码通常都缺少函数、变量等，一步一步寻找并扣到缺失的函数或变量，形成完整的解密代码。
- https://www.bilibili.com/video/BV1ew411K7nB?t=221.2&p=15
- 控制台输入内容的返回值若是函数，双击函数可以跳转到代码位置。
- 使用 `函数名.toString()` 的方式可以直接在控制台输出函数。此时从控制台直接拷贝即可。
- 在源代码页面把鼠标放到函数上，可以看到 `FunctionLacation` 标志，右侧有个连接，点击即可跳转到该函数。
- 通过控制台的代码段重写 `JSON.stringify`、`JSON.parse` 等常见方法的逻辑，添加 debugger 关键字，以便在代码无法找到时，虽慢但准得找到解密相关代码。
  - https://www.bilibili.com/video/BV1Cz4y1w78y

补环境。有的代码，可能会获取浏览器特定的一些属性，比如 [WebAPIs](/docs/Web/WebAPIs/WebAPIs.md) 中 document、window 等对象中的数据。此时如果使用代码编译器运行代码的场合是没有办法获取到这些信息的，需要在代码中手动造一些浏览器信息。

## JS 逆向的调试方法

> 参考：
>
> - [JS 逆向的调试方法](https://www.bilibili.com/video/BV1ew411K7nB/?p=3)

- 使用 `Ctrl + Shift + C` 定位感兴趣的元素。
- 在网络中使用 `Ctrl + f` 打开搜索，并搜索关键字，找到被选中的请求
- 从网络标签中找到被选中（深灰色北京）的请求，点击该请求的 `Initiator(启动器)` 进入到代码位置
- 如何查找代码有多种方式。首先使用左下角的 *美观输出*，随后在代码中寻找感兴趣的内容。
  - 比如查找 sign 之类的关键字，找到 sign 的生成逻辑。

## 已加密的数据如何处理

> - [遇到数据加密如何处理（一）](https://www.bilibili.com/video/BV1ew411K7nB/?p=4)
> - [遇到数据加密如何处理（二）](https://www.bilibili.com/video/BV1ew411K7nB?p=6)
> - 视频已上传到网盘

已加密的数据在上文调试方法中是搜索不到的。这种通常都是使用 [AJAX](/docs/2.编程/高级编程语言/ECMAScript/ECMAScript%20网络编程/AJAX.md) 获取的数据，也就是对响应体中的 json 加密，由客户端解密才能看到。

此时通过保存加密数据的键（常见的为 key、encrypt 等等），找到代码位置，打上断点后，进行 *单步调试* 直到找到生成逻辑。

## JS 混淆如何处理

> - [关于数据加密js混淆的处理方式](https://www.bilibili.com/video/BV1ew411K7nB/?p=7)

可以参考下面复杂代码逻辑的查找部分。

## 复杂代码逻辑的查找

假如一个接口需要填写 signature，而且直接搜索关键字无法找到，那么就需要通过通过接口进入代码，逐步调试，直到找到 加密/解密 的逻辑。

首先从 Fetch/XHR 接口的 Initiator(启动器) 进入接口代码，在进来的位置打上断点，刷新数据。此时查看 **Call Stack(调用堆栈)**。Call Stack 从上往下就是代码执行顺序的逆向，Call Stack 往下，就是上一步代码的执行位置。

![image.png|300](https://notes-learning.oss-cn-beijing.aliyuncs.com/reverse/202401181743939.png)

比如当前在 `(anonymous)` 的位置，通常某个 Fetch/XHR 接口通过启动器进来的地方都是 send() 方法，也就是发送请求的最后一步，在这里找到 sign 的值。

然后我们在 Call Stack 中往下点击 e.exporter 进入上一步函数调用的位置，依次类推。然后若是到某个 Stack 处没有数据可显示，就在这里打上断点并取消上一个断点，重新刷新获取数据，这么循环查找下去。

直到某个 Stack 中没有 sign 数据的地方为止。那么 sign 的生成逻辑就在这个调用附近。（通常还没到 sign 没数据的时候，一般也就找到了）

> 可以从 Scope(作用域)、Console、鼠标移动到源码的变量或函数上，以查看 sign 的值。

## 请求参数中的 sign 如何处理

> [请求参数sign逆向](https://www.bilibili.com/video/BV1ew411K7nB/?p=9)

## Hook 代码段

有些混淆代码无法查找，通过编写代码段后调试，可以快速定位。

比如下面这个，可以通过 `JSON.stringify = function (params) {}` 重新定义 `JSON.stringify` 方法的逻辑，在其中加入文本输出和 debug 暂停能力。只需要在 [DevTools](/docs/Web/Browser/DevTools.md) - 源代码 - 代码段 中插入下面的代码后运行即可。

```js
(function () {
    var my_stringify = JSON.stringify;
    JSON.stringify = function (params) {
        console.log("Hook 字符串化", params);
        debugger
        return my_stringify(params);
    }
    var my_parse = JSON.parse;
    JSON.parse = function (params) {
        console.log("Hook 解析", params);
        debugger
        return my_parse(params);
    }
})();
```

这时，凡是页面的 JS 代码中调用了 JSON.stringify 和 JSON.parse 这俩方法的地方，都会输出参数，并被 debugger 关键字暂停以进行断点检查。然后可以在右侧 *调用堆栈* 中点击直接跳转到网页的代码中，对应的位置（堆栈是顺序执行的，查看下面几个即可）。

上面的做法有一个文件，就是页面刷新时就失效了，只有在当前页面发起 Fetch/XHR 之类的请求时才会被拦截。若想拦截 Cookie.set 之类的请求，就需要用到代理，在代理中添加 Hook 逻辑，详见下面的 Cookie 处理。

## Cookie 处理

动态生成 cookie、时效性 cookie、需要登录网站

https://www.bilibili.com/video/BV1ew411K7nB/?p=17 及后面几 P

想要 Hook Cookie，需要通过类似代理的方式进行，在 代理中添加如下代码。

```javascript
//当前版本hook工具只支持Content-Type为html的自动hook
//下面是一个示例:这个示例演示了hook全局的cookie设置点
(function() {
    //严谨模式 检查所有错误
    'use strict';
    //document 为要hook的对象   这里是hook的cookie
    var cookieTemp = "";
    Object.defineProperty(document, 'cookie', {
        //hook set方法也就是赋值的方法 
        set: function(val) {
        //这样就可以快速给下面这个代码行下断点
        //从而快速定位设置cookie的代码
        if (val.indexOf('w_tsfp') != -1){debugger}
            console.log('Hook捕获到cookie设置->', val);
            cookieTemp = val;
            return val;
        },
        //hook get方法也就是取值的方法 
        get: function()
        {
            return cookieTemp;
        }
    });
})();
```

# Android 逆向

[GitHub 项目，rev1si0n/lamda](https://github.com/rev1si0n/lamda)

- [lamda安卓逆向辅助框架](http://www.lxspider.com/?p=194)

[简书，某App接口逆向过程](https://www.jianshu.com/p/040d54a57e33)

### Jadx

> 参考：
>
> - [GitHub 项目，skylot/jadx](https://github.com/skylot/jadx)
> - [博客园，jadx 使用](https://www.cnblogs.com/lsgxeva/p/13500813.html)

Jadx 是一个 Dex 到 Java 反编译器。用于从 Android Dex 和 Apk 文件生成 Java 源代码的 CLI 和 GUI 工具

> 注意：请注意，在大多数情况下，jadx 无法100% 地反编译所有代码，因此会出现错误。检查[故障排除指南](https://github.com/skylot/jadx/wiki/Troubleshooting-Q&A#decompilation-issues)以了解解决方法

