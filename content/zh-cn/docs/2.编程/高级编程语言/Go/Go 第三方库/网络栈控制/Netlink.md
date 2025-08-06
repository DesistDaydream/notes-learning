---
title: Netlink
linkTitle: Netlink
weight: 20
---

# 概述

> 参考：
>
> - [GitHub 项目，vishvananda/netlink](https://github.com/vishvananda/netlink)
> - [GitHub 项目，mdlayher/netlink](https://github.com/mdlayher/netlink)

[Netlink](/docs/1.操作系统/Kernel/Network/Linux%20网络栈管理/Netlink/Netlink.md) 是 Linux 中的用户空间程序用来与内核进行通信的界面。它可以用于添加和删除接口，设置 ip 地址和路由以及配置 ipsec。

# vishvananda/netlink

Netlink 通信需要提升的权限，因此在大多数情况下，此代码需要以 root 用户身份运行。由于低级 netlink 消息充其量是难以理解的，因此该库试图提供一个以 iproute2 提供的 CLI 为松散建模的 api。ip 链接添加之类的操作将通过类似名称的函数 (例如 AddLink()) 来完成。该库的生命开始于 docker/libcontainer 中的 netlink 功能分支，但经过大量重写以提高可测试性，性能并添加 ipsec xfrm 处理等新功能。

## Hello World

```go
package main

import (
 "fmt"

 "github.com/vishvananda/netlink"
)

func main() {
 // 实例化一个 LinkAttrs,LinkAttrs 包含一个网络设备的绝大部分属性
 linkAttrs := netlink.NewLinkAttrs()
 // 设定 link 的名称
 linkAttrs.Name = "br0"
 // 将实例化的 LinkAttrs 信息赋值给 Bridge 结构体
 mybridge := &netlink.Bridge{LinkAttrs: linkAttrs}
 // 这里就算真正完成了一个网络设备的定义，netlink 库中包含多种网络设备结构体
 // 每种网络设备结构体都实现了 Link 接口
 // Link 接口只有两个方法，Attrs() 用来返回 LinkAttrs 结构体，Type() 用来返回该网络设备的类型。
 // 而对各种类型的网络设备实现增删改查的函数，其接受的参数就是 Link 接口类型
 // 所以 Link 接口的主要作用，就是用来区分不同类型的网络设备，以便可以在调用时统一。对网络设备的任何操作，都可以将 Link 接口作为参数互相传递。

 // 使用 Bridge 结构体的信息创建一个网络设备
 err := netlink.LinkAdd(mybridge)
 if err != nil {
  fmt.Printf("could not add %s: %v\n", linkAttrs.Name, err)
 }
 // eth1, _ := netlink.LinkByName("eth1")
 // netlink.LinkSetMaster(eth1, mybridge)
}
```

## 实践

```go
package main

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

// addBridge 创建一个桥设备
func addBridge() *netlink.Bridge {
	// 实例化一个 LinkAttrs,LinkAttrs 结构体包含一个网络设备的绝大部分属性。
	linkAttrs := netlink.NewLinkAttrs()
	// 设定 link 的名称
	linkAttrs.Name = "br0"
	// 将实例化的 link 信息赋值给 Bridge 结构体
	myBridge := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}

	// 使用 Bridge 结构体的信息创建一个 link
	err := netlink.LinkAdd(myBridge)
	if err != nil {
		fmt.Printf("could not add %s: %v\n", linkAttrs.Name, err)
	}

	return myBridge
}

// addVeth 创建一个 veth 设备
func addVeth() *netlink.Veth {
	// 实例化一个 LinkAttrs,LinkAttrs 结构体包含一个网络设备的绝大部分属性。
	linkAttrs := netlink.NewLinkAttrs()
	// 设定 link 的名称
	linkAttrs.Name = "veth1.1"
	// 将实例化的 link 信息赋值给 Veth 结构体,veth 必须指定对端设备
	myVeth := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  "veth1.2",
	}

	// 使用 Veth 结构体的信息创建一个 link
	err := netlink.LinkAdd(myVeth)
	if err != nil {
		fmt.Printf("could not add %s: %v\n", linkAttrs.Name, err)
	}

	return myVeth
}

func main() {
	// 增
	// 添加 veth 和 bridge 设备,并将 veth 设备附加到 bridge 设备上。
	if err := netlink.LinkSetMaster(addVeth(), addBridge()); err != nil {
		fmt.Println("设置网络设备主从关系出错,原因：", err)
	}

	// 删除 veth 和 bridge 设备
	br0, _ := netlink.LinkByName("br0")
	veth, _ := netlink.LinkByName("veth1.1")
	netlink.LinkDel(br0)
	netlink.LinkDel(veth)

	// 改

	// 查
	// 实例化一个 Handle，相当于在当前名称空间中创建一个 Socket 句柄。
	// 呼叫者可以指定句柄应支持的netlink族。如果未指定族，则将自动添加netlink软件包支持的所有族。
	handle, _ := netlink.NewHandle()
	// 获取所有网络设备，等效于 `ip link show` 命令
	links, _ := handle.LinkList()
	// 输出网络设备信息
	for _, link := range links {
		fmt.Printf("设备名称为：%v\n", link.Attrs().Name)
		// 获取一个网络设备的 IP 地址
		addrs, _ := netlink.AddrList(link, 0)
		for index, addr := range addrs {
			fmt.Printf("%v 设备的第 %v 个 IP 地址为：%v\n", link.Attrs().Name, index+1, addr.IP)
		}
	}
}
```
