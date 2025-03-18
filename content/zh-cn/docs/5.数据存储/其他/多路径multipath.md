---
title: 多路径multipath
---

所需安装的组件：

1. device-mapper-multipath-libs-XXXX
2. device-mapper-multipath-XXXX


配置说明详见官方文档：<https://access.redhat.com/documentation/zh-cn/red_hat_enterprise_linux/7/html/dm_multipath/mpio_description>

配置示例：

```
# This is a basic configuration file with some examples, for device mapper
# multipath.
#
# For a complete list of the default configuration values, run either
# multipath -t
# or
# multipathd show config
#
# For a list of configuration options with descriptions, see the multipath.conf
# man page

## By default, devices with vendor = "IBM" and product = "S/390.*" are
## blacklisted. To enable mulitpathing on these devies, uncomment the
## following lines.
#blacklist_exceptions {
#	device {
#		vendor	"IBM"
#		product	"S/390.*"
#	}
#}

## Use user friendly names, instead of using WWIDs as names.
defaults {
  user_friendly_names yes
  find_multipaths no
}
##
## Here is an example of how to configure some standard options.
##
#
#defaults {
#	polling_interval 	10
#	path_selector		"round-robin 0"
#	path_grouping_policy	multibus
#	uid_attribute		ID_SERIAL
#	prio			alua
#	path_checker		readsector0
#	rr_min_io		100
#	max_fds			8192
#	rr_weight		priorities
#	failback		immediate
#	no_path_retry		fail
#	user_friendly_names	yes
#}
##
## The wwid line in the following blacklist section is shown as an example
## of how to blacklist devices by wwid.  The 2 devnode lines are the
## compiled in default blacklist. If you want to blacklist entire types
## of devices, such as all scsi devices, you should use a devnode line.
## However, if you want to blacklist specific devices, you should use
## a wwid line.  Since there is no guarantee that a specific device will
## not change names on reboot (from /dev/sda to /dev/sdb for example)
## devnode lines are not recommended for blacklisting specific devices.
##
#blacklist {
#       wwid 26353900f02796769
#	devnode "^(ram|raw|loop|fd|md|dm-|sr|scd|st)[0-9]*"
#	devnode "^hd[a-z]"
#}
blacklist  {
       devnode "^sda"
#       devnode "^sdb"
}
#multipaths {
#	multipath {
#		wwid			3600508b4000156d700012000000b0000
#		alias			yellow
#		path_grouping_policy	multibus
#		path_selector		"round-robin 0"
#		failback		manual
#		rr_weight		priorities
#		no_path_retry		5
#	}
#	multipath {
#		wwid			1DEC_____321816758474
#		alias			red
#	}
#}
#devices {
#	device {
#		vendor			"COMPAQ  "
#		product			"HSV110 (C)COMPAQ"
#		path_grouping_policy	multibus
#		path_checker		readsector0
#		path_selector		"round-robin 0"
#		hardware_handler	"0"
#		failback		15
#		rr_weight		priorities
#		no_path_retry		queue
#	}
#	device {
#		vendor			"COMPAQ  "
#		product			"MSA1000         "
#		path_grouping_policy	multibus
#	}
#}
devices {
        device {
                vendor                  "HP"
                product                 "MSA 2040 SAN"
                path_grouping_policy    group_by_prio
                path_checker            tur
                path_selector           "round-robin 0"
                hardware_handler        "0"
                failback                immediate
                rr_weight               priorities
                no_path_retry           queue
        }
}
```