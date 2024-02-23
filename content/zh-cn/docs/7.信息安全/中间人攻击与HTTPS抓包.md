---
title: 中间人攻击与HTTPS抓包
---

# 概述

> 参考：
>
> - [公众号-小林 coding，字节一面：HTTPS 一定安全可靠吗？](https://mp.weixin.qq.com/s/IYZrtK7pJTeCOBwY-sC9HA)

**HTTPS 一定安全可靠吗？**

这个问题的场景是这样的：客户端通过浏览器向服务端发起 HTTPS 请求时，被「假基站」转发到了一个「中间人服务器」，于是客户端是和「中间人服务器」完成了 TLS 握手，然后这个「中间人服务器」再与真正的服务端完成 TLS 握手。

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNicwUSq8LhcU0v5GWmgbAfafNE9hhWWR21UBeUSemHqPicHPynYS0A9CA/640?wx_fmt=png)

具体过程如下：

- 客户端向服务端发起 HTTPS 建立连接请求时，然后被「假基站」转发到了一个「中间人服务器」，接着中间人向服务端发起 HTTPS 建立连接请求，此时客户端与中间人进行 TLS 握手，中间人与服务端进行 TLS 握手；
- 在客户端与中间人进行 TLS 握手过程中，中间人会发送自己的公钥证书给客户端，**客户端验证证书的真伪**，然后从证书拿到公钥，并生成一个随机数，用公钥加密随机数发送给中间人，中间人使用私钥解密，得到随机数，此时双方都有随机数，然后通过算法生成对称加密密钥（A），后续客户端与中间人通信就用这个对称加密密钥来加密数据了。
- 在中间人与服务端进行 TLS 握手过程中，服务端会发送从 CA 机构签发的公钥证书给中间人，从证书拿到公钥，并生成一个随机数，用公钥加密随机数发送给服务端，服务端使用私钥解密，得到随机数，此时双方都有随机数，然后通过算法生成对称加密密钥（B），后续中间人与服务端通信就用这个对称加密密钥来加密数据了。
- 后续的通信过程中，中间人用对称加密密钥（A）解密客户端的 HTTPS 请求的数据，然后用对称加密密钥（B）加密 HTTPS 请求后，转发给服务端，接着服务端发送 HTTPS 响应数据给中间人，中间人用对称加密密钥（B）解密 HTTPS 响应数据，然后再用对称加密密钥（A）加密后，转发给客户端。

从客户端的角度看，其实并不知道网络中存在中间人服务器这个角色。

那么中间人就可以解开浏览器发起的 HTTPS 请求里的数据，也可以解开服务端响应给浏览器的 HTTPS 响应数据。相当于，中间人能够 “偷看” 浏览器与服务端之间的 HTTPS 请求和响应的数据。

但是要发生这种场景是有前提的，**前提是用户点击接受了中间人服务器的证书。**

中间人服务器与客户端在 TLS 握手过程中，实际上发送了自己伪造的证书给浏览器，而这个伪造的证书是能被浏览器（客户端）识别出是非法的，于是就会提醒用户该证书存在问题。

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNzricwiaRlwzuc1CUqXYTYEapwDNUybxhiaWfofHt7LSuzC4cFdYvwcJoA/640?wx_fmt=png)

如果用户执意点击「继续浏览此网站」，相当于用户接受了中间人伪造的证书，那么后续整个 HTTPS 通信都能被中间人监听了。

所以，这其实并不能说 HTTPS 不够安全，毕竟浏览器都已经提示证书有问题了，如果用户坚决要访问，那**不能怪 HTTPS ，得怪自己手贱**。

## 客户端是如何验证证书的？

接下来，详细说一下实际中数字证书签发和验证流程。

如下图图所示，为数字证书签发和验证流程：

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNMxlSGIfI1kiav1TwxTK6ZtGiaCiaygM3g0fFnPYo3n9atZlbhvpkWt77Q/640?wx_fmt=png)
当服务端向 CA 机构申请证书的时候，CA 签发证书的过程，如上图左边部分：

- 首先 CA 会把持有者的公钥、用途、颁发者、有效时间等信息打成一个包，然后对这些信息进行 Hash 计算，得到一个 Hash 值；
- 然后 CA 会使用自己的私钥将该 Hash 值加密，生成 Certificate Signature，也就是 CA 对证书做了签名；
- 最后将 Certificate Signature 添加在文件证书上，形成数字证书；

客户端校验服务端的数字证书的过程，如上图右边部分：

- 首先客户端会使用同样的 Hash 算法获取该证书的 Hash 值 H1；
- 通常浏览器和操作系统中集成了 CA 的公钥信息，浏览器收到证书后可以使用 CA 的公钥解密 Certificate Signature 内容，得到一个 Hash 值 H2 ；
- 最后比较 H1 和 H2，如果值相同，则为可信赖的证书，否则则认为证书不可信。

但事实上，证书的验证过程中**还存在一个证书信任链的问题**，因为我们向 CA 申请的证书一般不是根证书签发的，而是由中间证书签发的，比如百度的证书，从下图你可以看到，证书的层级有三级：

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNVINCkWDhqk7OiaibpvuUnYAt2GXtShM9z5bWCG1TSZH6LzzylakSk83w/640?wx_fmt=png)

对于这种三级层级关系的证书的验证过程如下：

- 客户端收到 baidu.com 的证书后，发现这个证书的签发者不是根证书，就无法根据本地已有的根证书中的公钥去验证 baidu.com 证书是否可信。于是，客户端根据 baidu.com 证书中的签发者，找到该证书的颁发机构是 “GlobalSign Organization Validation CA - SHA256 - G2”，然后向 CA 请求该中间证书。
- 请求到证书后发现 “GlobalSign Organization Validation CA - SHA256 - G2” 证书是由 “GlobalSign Root CA” 签发的，由于 “GlobalSign Root CA” 没有再上级签发机构，说明它是根证书，也就是自签证书。应用软件会检查此证书有否已预载于根证书清单上，如果有，则可以利用根证书中的公钥去验证 “GlobalSign Organization Validation CA - SHA256 - G2” 证书，如果发现验证通过，就认为该中间证书是可信的。
- “GlobalSign Organization Validation CA - SHA256 - G2” 证书被信任后，可以使用 “GlobalSign Organization Validation CA - SHA256 - G2” 证书中的公钥去验证 baidu.com 证书的可信性，如果验证通过，就可以信任 baidu.com 证书。

在这四个步骤中，最开始客户端只信任根证书 GlobalSign Root CA 证书的，然后 “GlobalSign Root CA” 证书信任 “GlobalSign Organization Validation CA - SHA256 - G2” 证书，而 “GlobalSign Organization Validation CA - SHA256 - G2” 证书又信任 baidu.com 证书，于是客户端也信任 baidu.com 证书。总括来说，由于用户信任 GlobalSign，所以由 GlobalSign 所担保的 baidu.com 可以被信任，另外由于用户信任操作系统或浏览器的软件商，所以由软件商预载了根证书的 GlobalSign 都可被信任。

![500](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNIPhXyF7S1WDIbuUbGwgLjwr47iau48ape1H3GYg8B0hDEKeqdz7xZgg/640?wx_fmt=png)

操作系统里一般都会内置一些根证书，比如我的 MAC 电脑里内置的根证书有这么多：

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNloHzrlsep6vZgNvz8sA8LDBHnAWyErqcpicRM6kovSB5L3IQcjdD5aw/640?wx_fmt=png)

这样的一层层地验证就构成了一条信任链路，整个证书信任链验证流程如下图所示：

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNFruFDENMpzLgnfRH9T8uN0TAOItwoqia9mXLNLsE7UErltjyBu2zicNQ/640?wx_fmt=png)

如果你的电脑中毒了，被恶意导入了中间人的根证书，那么在验证中间人的证书的时候，由于你操作系统信任了中间人的根证书，那么等同于中间人的证书是合法的。

这种情况下，浏览器是不会弹出证书存在问题的风险提醒的。

这其实也不关 HTTPS 的事情，是你电脑中毒了才导致 HTTPS 数据被中间人劫持的。

所以，**HTTPS 协议本身到目前为止还是没有任何漏洞的，即使你成功进行中间人攻击，本质上是利用了客户端的漏洞（用户点击继续访问或者被恶意导入伪造的根证书），并不是 HTTPS 不够安全。**

## 为什么抓包工具能截取 HTTPS 数据？

抓包工具 Fiddler 之所以可以明文看到 HTTPS 数据，工作原理与中间人一致的。

对于 HTTPS 连接来说，中间人要满足以下两点，才能实现真正的明文代理:

1. 中间人，作为客户端与真实服务端建立连接这一步不会有问题，因为服务端不会校验客户端的身份；
2. 中间人，作为服务端与真实客户端建立连接，这里会有客户端信任服务端的问题，也就是服务端必须有对应域名的私钥；

中间人要拿到私钥只能通过如下方式：

1. 去网站服务端拿到私钥；
2. 去CA处拿域名签发私钥；
3. 自己签发证书，且被浏览器信任；

不用解释，抓包工具只能使用第三种方式取得中间人的身份。

使用抓包工具进行 HTTPS 抓包的时候，需要在客户端安装 Fiddler 的根证书，这里实际上起认证中心（CA）的作用。

Fiddler 能够抓包的关键是客户端会往系统受信任的根证书列表中导入 Fiddler 生成的证书，而这个证书会被浏览器信任，也就是 Fiddler 给自己创建了一个认证中心 CA。

客户端拿着中间人签发的证书去中间人自己的 CA 去认证，当然认为这个证书是有效的。

## 如何避免被中间人抓取数据？

我们要保证自己电脑的安全，不要被病毒乘虚而入，而且也不要点击任何证书非法的网站，这样 HTTPS 数据就不会被中间人截取到了。

当然，我们还可以通过 HTTPS 双向认证来避免这种问题。

一般我们的 HTTPS 是单向认证，客户端只会验证了服务端的身份，但是服务端并不会验证客户端的身份。

如果用了双向认证方式，不仅客户端会验证服务端的身份，而且服务端也会验证客户端的身份。

![](https://mmbiz.qpic.cn/mmbiz_png/J0g14CUwaZf0Dm1SNmvUayxCflHTROibNccBibe4MEf260dWmCfGnvHgz7iaZza3s1Sq6UwibbuGickDfMIcwAlGBGQ/640?wx_fmt=png)

服务端一旦验证到请求自己的客户端为不可信任的，服务端就拒绝继续通信，客户端如果发现服务端为不可信任的，那么也中止通信。

# 如何解决 SSL/TLS Pinning

很多移动端应用程序为了防止自己被抓包，会使用 [SSL/TLS Pinning](/docs/7.信息安全/Cryptography/公开密钥加密/证书%20与%20PKI/SSL_TLS%20Pinning.md) 技术。此时，就算我们在移动设备上向系统中添加了 Charles 的 CA 证书，很多接口的请求也抓不到，甚至有可能会提示无法连接网络。

这该如何解决呢？
