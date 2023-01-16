---
title: Interface 命令
---

# 概述

> ## 参考：

iface-begin create a snapshot of current interfaces settings, which can be later committed (iface-commit) or restored (iface-rollback)
iface-bridge create a bridge device and attach an existing network device to it
iface-commit commit changes made since iface-begin and free restore point
iface-define define an inactive persistent physical host interface or modify an existing persistent one from an XML file
iface-destroy destroy a physical host interface (disable it / "if-down")
iface-dumpxml interface information in XML
iface-edit edit XML configuration for a physical host interface

# iface-list # 列出宿主机的接口(i.e.宿主机的网络设备)

iface-mac convert an interface name to interface MAC address
iface-name convert an interface MAC address to interface name
iface-rollback rollback to previous saved configuration created via iface-begin
iface-start start a physical host interface (enable it / "if-up")
iface-unbridge undefine a bridge device after detaching its slave deviceiface-unbridge&#x20;
iface-undefine undefine a physical host interface (remove it from configuration)
