---
title: Netlink 操作
---

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
