---
title: Xshell 脚本
---

# 概述

http://www.cxyzjd.com/article/qq_34347375/113464381

# 脚本示例

## vbs

```vbnet
Sub Main
	xsh.Screen.WaitForString("请选择资源编号：")
	xsh.Screen.Send("i")
	xsh.Screen.Send(VbCr)
	xsh.Screen.WaitForString("请输入IP地址：")
	xsh.Screen.Send("10.213.30.1")
	xsh.Screen.Send(VbCr)
	xsh.Screen.WaitForString("请选择资源编号：")
	xsh.Screen.Send("4")
	xsh.Screen.Send(VbCr)
	xsh.Screen.WaitForString("需自学习从账号请按回车鍵继续：")
	xsh.Screen.Send(VbCr)
	xsh.Screen.WaitForString("请输入资源从账号：")
	xsh.Screen.Send("root")
	xsh.Screen.Send(VbCr)
	xsh.Screen.WaitForString("请输入从账号密码：")
	xsh.Screen.Send("!SgXmdz!4m")
	xsh.Screen.Send(VbCr)
End Sub
```

注意：若 vbs 无法识别中文，则将文件保存时，选择编码为ANSI。

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/xgwgbn/1618847088206-c4ccd20c-e089-4eb5-ad6c-4eeb47a56042.png)

## python

```python
def Main():
    xsh.Screen.WaitForString("请选择资源编号：")
    xsh.Screen.Send("i\r")
    xsh.Screen.Send("\r")
    xsh.Screen.WaitForString("请输入IP地址：")
    xsh.Screen.Send("10.213.30.1")
    xsh.Screen.Send("\r")
    xsh.Screen.WaitForString("请选择资源编号：")
    xsh.Screen.Send("4")
    xsh.Screen.Send("\r")
    xsh.Screen.WaitForString("需自学习从账号请按回车鍵继续：")
    xsh.Screen.Send("\r")
    xsh.Screen.WaitForString("请输入资源从账号：")
    xsh.Screen.Send("root")
    xsh.Screen.Send("\r")
    xsh.Screen.WaitForString("请输入从账号密码：")
    xsh.Screen.Send("!SgXmdz!4m")
    xsh.Screen.Send("\r")
```
