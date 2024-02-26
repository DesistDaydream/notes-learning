---
title: Selenium
linkTitle: Selenium
date: 2023-11-26T21:25
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，SeleniumHQ/selenium](https://github.com/SeleniumHQ/selenium)
> - [官网](https://www.selenium.dev/)
> - https://selenium-python.readthedocs.io/ 特定于 Python 的文档，官方文档很多示例都不全。
>   - https://segmentfault.com/q/1010000043032537 这里表示文档不全

Selenium 使浏览器自动化，用于自动化 Web 应用程序。Selenium 通过 [WebDriver](/docs/Web/WebDriver/WebDriver.md) 以控制浏览器。

Selenium 启动的浏览器参数

```
C:\Program Files\Google\Chrome\Application\chrome.exe"
 --allow-pre-commit-input
 --disable-background-networking
 --disable-backgrounding-occluded-windows
 --disable-client-side-phishing-detection
 --disable-default-apps
 --disable-hang-monitor
 --disable-popup-blocking
 --disable-prompt-on-repost
 --disable-sync
 --enable-automation
 --enable-logging
 --log-level=0
 --no-first-run
 --no-service-autorun
 --password-store=basic
 --remote-debugging-port=0
 --test-type=webdriver
 --use-mock-keychain
 --user-data-dir="C:\Users\DESIST~1\AppData\Local\Temp\scoped_dir11888_46574381"
 --flag-switches-begin
 --flag-switches-end data:
```

正常启动浏览器的参数

```
"C:\Program Files\Google\Chrome\Application\chrome.exe"
 --flag-switches-begin
 --flag-switches-end
 --origin-trial-disabled-features=WebGPU
```

selenium 依赖于 WebDriver，这里我们实例化了一个 Chrome 的 WebDriver。

然后使用 chromedriver.exe 启动浏览器。

# Selenium 使用浏览器的方式

打开新的浏览器

使用当前运行的浏览器

- 通过 Chrome 的 debug 端口连接

使用指定的缓存

# Selenium 关联文件与配置

**~/.cache/selenium/chromedriver/win64/${VERSION}/** # chromedriver.exe 文件的默认保存位置

**${TEMP}/scoped_dirXXXXX/Default** # Selenium 启动的浏览器的用户数据目录。浏览器会在 TMP 目录中创建用于保存用户数据的目录作为 Chrome 的个人资料路径。

# WebDriver

Selenium WebDriver 是 [W3C 推荐标准](https://www.w3.org/TR/webdriver1/)

- WebDriver 被设计成一个简单和简洁的编程接口。
- WebDriver 是一个简洁的面向对象 API。
- 它能有效地驱动浏览器。

## 元素

> 参考：
>
> - [官方文档，WebDriver-元素-定位器](https://www.selenium.dev/zh-cn/documentation/webdriver/elements/locators/)

在 WebDriver 中有 8 种不同的内置元素定位策略：

- **class 名称** #
- **css 选择器** #
- **ID** #
- **name** #
- ......略
- **tag 名称** #
- **XPath** # 通过 [XML](/docs/2.编程/标记语言/XML.md) 的 XPath 表达式定位元素
  - XPath 表达式可以在浏览器的 [DevTools](/docs/Web/Browser/DevTools.md) 中先使用 Ctrl+Shift+C 快捷键找到想要定位的元素，然后右键点击该元素，选择 `复制 - 复制 XPath`。即可获得 XPath 表达式

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/selenium/202312031841908.png)

## 交互

### 窗口和标签页

> 参考：
>
> - [官方文档，WebDriver-交互-窗口](https://www.selenium.dev/zh-cn/documentation/webdriver/interactions/windows/)
> - [51CTO，selenium工具UI自动化之浏览器的窗口切换（句柄切换）](https://blog.51cto.com/u_15688254/5723115)

WebDriver 不区分窗口和标签页。如果打开了一个新标签页或窗口，Selenium 将使用 **Handle(句柄)** 处理它，每个标签页和窗口的 Handle 是其唯一的标识符，该标识符在单个会话中保持持久性。

可以使用以下方法获得当前窗口的窗口句柄:

```python
driver.current_window_handle
```

> 解释一下什么叫“不区分窗口和标签页”，窗口是指 ctrl + n 创建的浏览器窗口，标签页是指 ctrl + t 创建的浏览器标签页。`driver.current_window_handle` 获取到的窗口 Handle 与窗口第一个标签页的 Handle 相同，若创建了第二个标签页，关闭了第一个标签页，则窗口的 Handle 与 第二个标签页的 Handle 相同。也就是说，窗口即是标签页。窗口就是打开的第一个标签页。

`driver.window_handles` 方法会获取所有窗口中所有标签页的句柄，返回一个 `List[str]`，不管标签页在多少个窗口中，哪怕窗口属于其他账户，也会一并获取到。

# 案例

https://github.com/onepureman/selenium_login_cracking 各种网站的滑动验证码使用selenium登陆破解，仅供交流学习，包括：京东，17173。。。

[python+selenium+opencv验证滑块](https://www.cnblogs.com/lihongtaoya/p/16793699.html)

- 参考了
  - https://blog.csdn.net/m0_59874815/article/details/121195481
  - https://github.com/gebiWangshushu/JDCaptchaCrack

[【 Python爬虫】京东滑块登录](https://www.cnblogs.com/wanghong1994/p/17786278.html) 滑动太快了

[selenium 元素定位总结篇](http://testingpai.com/article/1689428003874#toc_h2_0)

[python获取html标签 python selenium获取html](https://blog.51cto.com/u_16099210/6987312)

[selenium获取元素属性](https://zhuanlan.zhihu.com/p/647664858)

https://zhuanlan.zhihu.com/p/638431974 # 软件测试/测试开发丨用户端web自动化测试学习笔记

https://zhuanlan.zhihu.com/p/94402506 python实现网站的自动登录（selenium实现，带验证码识别）

https://www.aneasystone.com/archives/2018/03/python-selenium-geetest-crack.html 使用 Python + Selenium 破解滑块验证码
