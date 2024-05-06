---
title: Blockchain
linkTitle: Blockchain
date: 2023-12-22T15:57
weight: 1
---

# 概述

> 参考：
> 
> - [Wiki，Blockchain](https://en.wikipedia.org/wiki/Blockchain)
> - [B 站，汪杰解惑 NFT-02](https://www.bilibili.com/video/BV1T34y117Y9)

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

钱包

- MetaMask 钱包
- Coinbase Wallet 钱包
- imtoken 钱包

NFT 交易平台

- OpenSea

**Percentage Fee(版税)** # 即提成的百分比。每次 NFT 交易时，最初的创建者会获得交易额的百分比的收入。

**Gas Fee(气体费)** # 铸造一个 NFT 是有成本的，需要向矿工支付 Gas Fee

钱包中的网络：就是“链”也就是“区块链”的链。在各种链上，可以搜索在当钱链上的交易记录。

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
# 交易

TODO: 如何将自己转移至交易所？手续费怎么算？转移至交易所后兑换成 USDT 出售。
