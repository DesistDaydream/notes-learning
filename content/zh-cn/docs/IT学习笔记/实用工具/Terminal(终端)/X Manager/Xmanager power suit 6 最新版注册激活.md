---
title: Xmanager power suit 6 最新版注册激活
---

参考：https://www.cnblogs.com/selier/p/9649064.html
note:一定要先修改 hosts 文件再安装
操作步骤
Xmanger Power Suit 官方 其实有两种 .exe 文件，一个是用于试用的，在注册的时候不能直接输入密钥。而另一个是为注册用户提供的 .exe 文件，在注册的时候可以输入密钥，直接可以激活了。

# 1、下载安装包

到 Xmanager Power Suit 页面点击 Download，并填写一些信息，试用版的下载链接就会发至邮箱。 https://www.netsarang.com/products/xme\_overview.html
当然发到你邮箱的链接还不是真正的下载地址，你需要复制到浏览器打开，然后它就会开始下载，这个时候复制的下载地址才有效
注意：中文官网不行的话，去英文官网下，还是可以的。现阶段中国通过代理网页下载，与国外下载方式不同。

# 2、 将下载地址复制下来，并在 .exe 之前加上字母 r。

比如我的下载地址是 https://cdn.netsarang.net/61292a29/XmanagerPowerSuite-6.0.0012.exe
修改之后下载地址是 https://cdn.netsarang.net/61292a29/XmanagerPowerSuite-6.0.0012r.exe
这就是注册版文件的最新版本。
如果你无法访问或下载等问题，可以直接下载我分享的百度网盘文件，定期更新
百度网盘：https://pan.baidu.com/s/1JVMBvameqRwmCCd7EtCMGQ 密码：esej

# 3、 清除注册表信息。

清理注册表（手动删除）
弹出运行框，快捷键：win+r
出入 regedit，打开注册表
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586685976-7b042048-6192-4531-8962-808f520ac508.jpeg)
HKEY_CURRENT_USER\Software\NetSarang
右键删除 NetSarang
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586686216-8f44f710-61d2-4675-819c-f887b9334a78.jpeg)
命令删除
弹出运行框，快捷键：win+r
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586686290-312c0c24-d297-4bce-9ddc-1fbb9844cb7a.jpeg)
执行：REG DELETE HKEY_CURRENT_USER\Software\NetSarang /f
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586686287-d4b3313b-bbd1-4df6-a09d-92179ade45ad.jpeg)

# 4、 添加 HOSTS 信息。

这些解析主要是为了防止自动更新用的
目录位置：C:\Windows\System32\drivers\etc
127.0.0.1 transact.netsarang.com
127.0.0.1 update.netsarang.com
127.0.0.1 www.netsarang.com
127.0.0.1 www.netsarang.co.kr
127.0.0.1 sales.netsarang.com

# 5、使用 Xmanager-keygen 生成序列号。地址：DoubleLabyrinth/Xmanager-keygen

使用 Python 运行文件即可生成
我的生成：171215-116481-999966

# 最后！！使用该序列号安装注册版 Xmanager Power Suit

推荐在线生成的网站：https://xshell.spppx.org/
安装成功截图，Lienses 为 999，Data 为 17 年
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586686561-82918502-6695-4c3e-ad68-c5407a5652e4.jpeg)
0|1 存疑部分
可能评估版本的 Xmanager Power Suit 也可以注册
可能 Xshell 等其他产品也可以用此种方法注册
其中部分操作可能是不必要的。
其他问题
官方有一个 Xshell Plus 6 版本，只有 Xshell 和 Xftp，没有其他功能，这个是个很好的选择。
当然也可以选择 Xmanager Power Suite 6，包含所有 4 个软件：Xmanager 6、Xshell 6、Xftp 6、Xlpd 6。
官方下载地址
https://cdn.netsarang.net/61292a29/XmanagerPowerSuite-6.0.0012r.exe
https://cdn.netsarang.net/61292a29/XshellPlus-6.0.0012r.exe
https://cdn.netsarang.net/61292a29/Xshell-6.0.0111r.exe
https://cdn.netsarang.net/61292a29/Xftp-6.0.0105r.exe
https://cdn.netsarang.net/61292a29/Xmanager-6.0.0105r.exe
https://cdn.netsarang.net/61292a29/Xlpd-6.0.0102r.exe
如果需要下载最新版本，可以直接访问官网的 最新版 ，将上面的下载地址改成相应的版本号即可
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/asbt9f/1618586687234-77a97415-6e2a-410d-b376-a6d17649f8b1.jpeg)
2018-09-14 第一次更新：调整
2018-10-18 第二次更新：修改最新版本为 XPS 6.0.0009
2018-12-29 第三次更新：修改最新版本为 XPS 6.0.0012
这是那个 python 文件 Xmanager-keygen.py
import datetime
import random

ProductCode = {
&#x20; 'Xmanager' : 0,
&#x20; 'Xshell' : 1,
&#x20; 'Xlpd' : 2,
&#x20; 'Xfile' : 3,
&#x20; 'Xftp' : 4,
&#x20; 'Xmanager 3D' : 5,
&#x20; 'Xmanager Enterprise' : 6,
&#x20; 'Xshell Plus' : 7
}

LicenseType = \[
&#x20; \[ ProductCode\['Xmanager'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x0C, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x10, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x16, 2, 'Student 2-year Subscription', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x18, 4, 'Student 4-year Subscription', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x20, 2, 'Student 2-year Subscription', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x22, 4, 'Student 4-year Subscription', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x3D, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x3E, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x41, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x42, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x47, 0, 'Standard Subscription', 2], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x48, 0, 'Educational Subscription', 2], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x4B, 0, 'Standard Subscription', 1], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x4C, 0, 'Educational Subscription', 1], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x51, 0, 'Standard', 2], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x52, 0, 'Educational', 2], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x55, 0, 'Standard', 1], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x56, 0, 'Educational', 1], # Concurrent Registered
&#x20; \[ ProductCode\['Xmanager'], 0x60, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x61, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x62, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x63, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x29, 1, 'CLS Class A', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x2A, 1, 'CLS Class B', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x2B, 1, 'CLS Class C', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x2C, 1, 'DLS', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x2D, 1, 'SLS', 2],
&#x20; \[ ProductCode\['Xmanager'], 0x33, 1, 'CLS Class A', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x34, 1, 'CLS Class B', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x35, 1, 'CLS Class C', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x36, 1, 'DLS', 1],
&#x20; \[ ProductCode\['Xmanager'], 0x37, 1, 'SLS', 1],
&#x20; \[ ProductCode\['Xshell Plus'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xshell'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xshell'], 0x0C, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xshell'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xshell'], 0x10, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xshell'], 0x16, 2, 'Student 2-year Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x18, 4, 'Student 4-year Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x20, 2, 'Student 2-year Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x22, 4, 'Student 4-year Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x3D, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x3E, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x41, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x42, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x47, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x48, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xshell'], 0x4B, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x4C, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xshell'], 0x51, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xshell'], 0x52, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xshell'], 0x55, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xshell'], 0x56, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xshell'], 0x60, 0, 'Standard', 1], # ������
&#x20; \[ ProductCode\['Xshell'], 0x61, 0, 'Standard', 2], # ������
&#x20; \[ ProductCode\['Xshell'], 0x62, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xshell'], 0x63, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x3D, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x3E, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x41, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x42, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x47, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x48, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x4B, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x4C, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x51, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x55, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x60, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x61, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xlpd'], 0x62, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xlpd'], 0x63, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xfile'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xftp'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xftp'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xftp'], 0x3D, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xftp'], 0x3E, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xftp'], 0x41, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xftp'], 0x42, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xftp'], 0x47, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xftp'], 0x48, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xftp'], 0x4B, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xftp'], 0x4C, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xftp'], 0x51, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xftp'], 0x55, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xftp'], 0x60, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xftp'], 0x61, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xftp'], 0x62, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xftp'], 0x63, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager 3D'], 0x0B, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager 3D'], 0x0C, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xmanager 3D'], 0x0F, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager 3D'], 0x10, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x0B, 0, '', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x0C, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x0F, 0, '', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x10, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x3D, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x3E, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x41, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x42, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x47, 0, 'Standard Subscription', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x48, 0, 'Educational Subscription', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x4B, 0, 'Standard Subscription', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x4C, 0, 'Educational Subscription', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x51, 0, '', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x52, 0, 'Educational', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x55, 0, '', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x56, 0, 'Educational', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x60, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x61, 0, 'Standard', 2],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x62, 0, 'Standard', 1],
&#x20; \[ ProductCode\['Xmanager Enterprise'], 0x63, 0, 'Standard', 2],
]

ProductPublishList = (
&#x20; { 'ProductName' : 'Xmanager', 'Version' : 2, 'PublishDate' : datetime.date(2003, 1, 1) },
&#x20; { 'ProductName' : 'Xshell', 'Version' : 2, 'PublishDate' : datetime.date(2004, 10, 1) },

    { 'ProductName' : 'Xmanager', 'Version' : 3, 'PublishDate' : datetime.date(2007, 1, 1) },<br />    { 'ProductName' : 'Xshell', 'Version' : 3, 'PublishDate' : datetime.date(2007, 1, 1) },<br />    { 'ProductName' : 'Xlpd', 'Version' : 3, 'PublishDate' : datetime.date(2007, 1, 1) },<br />    { 'ProductName' : 'Xftp', 'Version' : 3, 'PublishDate' : datetime.date(2007, 1, 1) },<br />    { 'ProductName' : 'Xmanager Enterprise', 'Version' : 3, 'PublishDate' : datetime.date(2007, 1, 1) },

    { 'ProductName' : 'Xmanager', 'Version' : 4, 'PublishDate' : datetime.date(2010, 8, 1) },<br />    { 'ProductName' : 'Xshell', 'Version' : 4, 'PublishDate' : datetime.date(2010, 8, 1) },<br />    { 'ProductName' : 'Xlpd', 'Version' : 4, 'PublishDate' : datetime.date(2010, 8, 1) },<br />    { 'ProductName' : 'Xftp', 'Version' : 4, 'PublishDate' : datetime.date(2010, 8, 1) },<br />    { 'ProductName' : 'Xmanager Enterprise', 'Version' : 4, 'PublishDate' : datetime.date(2010, 8, 1) },<br />    { 'ProductName' : 'Xmanager', 'Version' : 5, 'PublishDate' : datetime.date(2014, 4, 28) },<br />    { 'ProductName' : 'Xshell', 'Version' : 5, 'PublishDate' : datetime.date(2014, 4, 28) },<br />    { 'ProductName' : 'Xlpd', 'Version' : 5, 'PublishDate' : datetime.date(2014, 4, 28) },<br />    { 'ProductName' : 'Xftp', 'Version' : 5, 'PublishDate' : datetime.date(2014, 4, 28) },<br />    { 'ProductName' : 'Xmanager Enterprise', 'Version' : 5, 'PublishDate' : datetime.date(2014, 4, 28) },<br />    { 'ProductName' : 'Xmanager', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) },<br />    { 'ProductName' : 'Xshell', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) },<br />    { 'ProductName' : 'Xshell Plus', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) },<br />    { 'ProductName' : 'Xlpd', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) },<br />    { 'ProductName' : 'Xftp', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) },<br />    { 'ProductName' : 'Xmanager Enterprise', 'Version' : 6, 'PublishDate' : datetime.date(2018, 4, 29) }<br /)

def GetChecksum(preProductKey : str):
&#x20; Checksum = 1
&#x20; for i in range(0, len(preProductKey)):
&#x20; if preProductKey\[i] != '-' and preProductKey\[i] != '8' and preProductKey\[i] != '9':
&#x20; place = int(preProductKey\[i])
&#x20; Checksum = (9 - place) \* Checksum % -1000
&#x20; Checksum = (Checksum + int(preProductKey\[9])) % 1000
&#x20; return Checksum

def GenerateProductKey(IssueDate : datetime.date,
&#x20; ProductName : str,
&#x20; ProductVersion : int,
&#x20; NumberOfLicense : int):
&#x20; if IssueDate.year < 2002:
&#x20; raise ValueError('IssueDate cannot be earlier than 2002.')
&#x20; if IssueDate > datetime.date.today() + datetime.timedelta(days = 7):
&#x20; raise ValueError('IssueDate cannot be later than today after a week.')
&#x20; if NumberOfLicense < 0 or NumberOfLicense > 999:
&#x20; raise ValueError('NumberOfLicense must vary from 0 to 999.')

    for item in ProductPublishList:<br />        if item['ProductName'] == ProductName and item['Version'] == ProductVersion:<br />            if item['PublishDate'] > IssueDate:<br />                raise ValueError('IssueDate cannot be earlier than the publish date.')<br />            break<br />        if item == ProductPublishList[-1]:<br />            raise ValueError('Invalid product.')

    preProductKey = '%02d%02d%02d-%02d%d%03d-%03d' % (IssueDate.year - 2000,<br />                                                        IssueDate.month,<br />                                                        IssueDate.day,<br />                                                        0x0B,<br />                                                        ProductCode[ProductName],<br />                                                        random.randint(0, 999),<br />                                                        NumberOfLicense)<br />    Checksum = GetChecksum(preProductKey)<br />    ProductKey = preProductKey + '%03d' % Checksum<br />    return ProductKey

print(GenerateProductKey(datetime.date(2017, 12, 15), 'Xmanager Enterprise', 5, 999))
