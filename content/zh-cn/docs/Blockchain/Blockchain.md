---
title: Blockchain
linkTitle: Blockchain
weight: 1
---

# 概述

> 参考：
>
> - [Wiki, Blockchain](https://en.wikipedia.org/wiki/Blockchain)
> - [Wiki, Cryptocurrency](https://en.wikipedia.org/wiki/Cryptocurrency)

**Blockchain(区块链)**

2008 年 10 月 31 日 《Bitconi: A Peer-to-Peer Electronic Cash System 》

- 保证信息的完整性和真实性
- 保证信息的不可否认性

# 数字加密货币

> 参考：
>
> - 如何把狗狗币|柴犬币 shib|放在 imtoken 钱包和 metamask 狐狸钱包中？
>   - 这里教如何添加代币，如何添加钱包中的网络
>   - <https://www.youtube.com/watch?v=Gn4FCh5DEvg>
> - 【狐狸钱包】一分钟学会，如何一键添加各种主网？
>   - <https://www.youtube.com/watch?v=f1JU8TGImA0>


## 钱包

> 参考：
>
> - [Wiki, Cryptocurrency wallet](https://en.wikipedia.org/wiki/Cryptocurrency_wallet)
> - https://support.metamask.io/start/learn/what-is-a-secret-recovery-phrase-and-how-to-keep-your-crypto-wallet-secure/ 什么是助记词以及如何保护加密钱包的安全
> - https://support.metamask.io/start/learn/metamask-is-a-self-custodial-wallet/ Metamask 是一个自托管钱包

**Cryptocurrency wallet(加密货币钱包，简称 钱包)** 本质上是一个具有唯一性的密钥对。私钥用来消费代币，公钥生成的地址用来接收代币。<font color="#ff0000">钱包本身不保存一共有多少加密货币</font>。我们只能查看一条区块链中，某个地址的钱包的所有交易记录，通过这些交易记录计算出该钱包实际拥有多少代币。

一个钱包可以是 物理介质、程序、在线服务、etc. 只要该钱包可以生成一个地址，符合条件接入区块链网络即可。可以是任何形态。

一个钱包通常包含如下几个内容

- **PrivateKey(私钥)** # 与公钥一同生成。也可以根据 BIP39 标准，由助记词生成
    - **Secret Recovery Phrase(密钥恢复短语，也称：助记词)** # 根据 BIP39 标准，通过助记词可以计算出私钥，但是私钥无法计算出助记词。助记词用来找回钱包的私钥
- **PublicKey(公钥)** # 与私钥一同生成。可用于生成一个地址。
- **Address(地址)** # 由公钥生成的用来标识钱包的唯一标识符。

钱包分为两种

- **self-custodial(自托管)**
- **交易所托管**

口语中的钱包通常都是指 self-custodial(自托管) 的代币钱包。交易所中的钱包如果非要类比的话，可以看作是 支付宝、微信、银行、etc. 机构存放我的钱的地方。

> tips: 只不过在加密货币的交易中，如果想要从 交易所的自己的钱包中，把 币 转到 自托管 的钱包中，可能是要手续费的。。。。( ╯□╰ )

**要想给自己的钱包存钱**，首先要选择一个相对不错的区块链（最好是公链），然后在钱包中配置该区块链网络，让钱包链接到该网络之后，即可开始交易。

从另一个角度说，每条区块链上都可以使用一个钱包，只不过这些钱包之间的交易记录默认是不互通的，此时就需要一个称为**跨链**的技术来解决这个问题。

TODO: 跨链交易是不是手续费非常高？假如有手续费，这手续费交给谁了？为什么会有手续费？

> [!Attention]
> 钱包的<font color="#ff0000">交易记录是特定于 Blockchain 的</font>（有时候也成为 XX 链、XX 网络、etc.），在某条链上的交易记录，无法在另一条链上看到。
>
> 在 A 链上可以有钱包 AA 的交易记录，在 B 链上也可以有钱包 AA 的交易记录，但是这两个交易记录并不互通。
>
> 这就会导致两个问题
>
> 1. 选择什么钱包很重要，实现钱包的应用程序如果不维护了该怎么办？e.g. MetaMask 钱包是以太坊维护的，那默认交易都是在以太坊的主要公有链上。就算在钱包中添加了其他网络，也没法防止钱包本身不再维护。是否应该选择硬件实体钱包？
> 2. 选择在哪条区块链进行交易很重要，如果交易所在的区块链直接不维护了怎么办？
>     1. 还有一个 ”提币“ 词，真词 TM 坑。提币 这词看似跟 提现 类似，但是与中国从支付宝提现到银行的行为或者商家某某提现到个人完全不是一个概念。
>     2. 所谓的 提币 本质好像就是跨链交易，而且是带着高昂手续费的跨链交易。
>     3. 只有”提“到火币的 ”HECO 链“ 才没有手续费。要是想提到 以太坊的 “ERC20 链”，需要支付高昂的手续费。从火币的 “HECO 链” 把币转到以太坊的 “ERC20 链” 也是高昂手续费。
>     4. 截止到 2025-07-17，这 “HECO 链” （http-mainnet.hecochain.com）依然用不了

### 钱包的实现

有很多程序可以通过 **助记词** 或者 **私钥** 来生成钱包

- MetaMask
- Coinbase Wallet
- imtoken
- etc.

一套助记词可以在各种不同的钱包实现上，生成钱包

### MetaMask 钱包示例

我们可以创建多个 account，每个 account 本质就是表示一个 Wallet(钱包)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/blockchain/metamask_wallet_01.png)

账户名下面的字符串是该钱包的 Address。查看钱包的细节，可以显示 Private key(私钥) 与 Secret Recovery Phrase(助记词)，会发现：

- 私钥各不相同
- 助记词完全相同

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/blockchain/metamask_wallet_details.png)

对于 MetaMask 来说，一个助记词可以生成多个私钥不同的钱包。

当我们使用一个全新的 MetaMask 时，需要使用助记词导入钱包，此时会将曾经创建过的所有钱包都展示出来。

## ETH

**Etherscan** 是以太坊网络的区块链浏览器。 该网站可用于搜索交易、区块、钱包地址、智能合约以及其它链上数据，属于最热门的以太坊区块链浏览器之一，免费向用户开放。 使用**Etherscan**即可详细了解如何与区块链、其他钱包以及 DAapp 进行交互

<https://etherscan.io/address/0xED783c0ee7444435d31555f0Ab23E30ac2d0a9Eb>

## 链

HECO # 火币的链

- 火币的 区块链浏览器 https://hecoinfo.com/

Etherscan # 以太坊的链

- 以太坊的 区块链浏览器 https://etherscan.io/

StarkNet

- https://www.starknet.io/en
- 该公链项目会不定期为向 GitHub 上 Star 数较多的项目提交过 PR 的开发者发送 STRK，称之为“空投奖励”。想要领取空投需要符合一定的条件，每次条件不定，时间不定。

## 应用示例

当我们从交易所中将各种数字货币提取到钱包中后，虽然可以看到交易已完成，但是我们在钱包中却看不到。

这是因为，在我们提币时，需要填写一个“提币网络”，这个“提币网络”其实就是指的区块链中的“链”，对应的就是各个钱包中的“网络”。

所以我们首先应该先为钱包添加对应交易的“链”。

### 添加链

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647499683695-dd622c13-b54a-4ff3-b643-63f616244c4e.png)

从对应的“链”找到其 URL，这里以火币的 HECO 链为例，从[这里](https://hecoinfo.com/apis#rpc)可以找到 HECO 的 Endpoint 与 链 ID

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647499855988-a4357488-abd5-436f-abb8-805f64774e9f.png)

然后将 Endpoint 和 链 ID 填入，即可在钱包中添加“链”

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647499879301-9ec3af5b-4d36-495b-ade5-41710057a624.png)

后续提币的操作时，我们选择的“提币网络”也要与“链”对应上才行，交易记录通常都是保存在交易双方所在的链上。

### 获取交易记录

在火币交易记录查到 `交易ID`

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647498496139-aa024459-ea29-47f5-b918-400007cd2539.png)

### 获取代币的合约地址

在[火币链](https://hecoinfo.com/)页面搜索该交易

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647498548140-52d38f6c-0e7b-4a3b-89d8-2907bea8d3c9.png)

可以获取到交易细节，然后查看 Tokens Transferred 中的 `For` 信息

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647498734107-01e68293-ec25-4b4d-9db8-1c5ca84a8552.png)

然后就可以获取到交易代币的合约地址

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647498797213-df1d9e0c-c2b7-466c-8936-58f930d053ab.png)

这个地址就是钱包中，添加自定义代币时，所使用的地址。

- USDT: 0xa71edc38d189767582c38a3145b5873052c3e47a
- DOGE: 0x40280E26A572745B1152A54D1D44F365DaA51618

### 导入代币

填入地址后，代币符号与小数精度将会自动出现

![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/lddbaw/1647498857327-a04a31ff-ff5c-4e32-8809-8bcc8582bdd2.png)

# NFT

[B 站，汪杰解惑 NFT-02](https://www.bilibili.com/video/BV1T34y117Y9)

NFT 交易平台

- OpenSea

**Percentage Fee(版税)** # 即提成的百分比。每次 NFT 交易时，最初的创建者会获得交易额的百分比的收入。

**Gas Fee(气体费)** # 铸造一个 NFT 是有成本的，需要向矿工支付 Gas Fee

# 交易

TODO: 如何将自己转移至交易所？手续费怎么算？转移至交易所后兑换成 USDT 出售。
