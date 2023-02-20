---
title: pktgen 内核自带的高性能网络测试工具
---

#

# pktgen 内核自带的高性能网络测试工具

```bash
modprobe pktgen

cat > /usr/local/bin/pgset <<EOF
local result
echo $1 > $PGDEV
result=`cat $PGDEV | fgrep "Result: OK:"`
if [ "$result" = "" ]; then
     cat $PGDEV | fgrep Result:
fi
EOF

chmod 755 /usr/local/bin/pgset

# 为0号线程绑定 eth0 网卡
export PGDEV=/proc/net/pktgen/kpktgend_0
pgset "rem_device_all" # 清空网卡绑定
pgset "add_device eth0" # 添加 eth0 网卡

# 配置 eth0 网卡的测试选项
export PGDEV=/proc/net/pktgen/eth0
pgset "count 1000000"    # 总发包数量
pgset "delay 5000"       # 不同包之间的发送延迟 (单位纳秒)
pgset "clone_skb 0"      # SKB 包复制
pgset "pkt_size 64"      # 网络包大小
pgset "dst 192.168.0.30" # 目的 IP
pgset "dst_mac 11:11:11:11:11:11"  # 目的 MAC

# 启动测试
export PGDEV=/proc/net/pktgen/pgctrl
pgset "start"
```

```bash
$ cat /proc/net/pktgen/eth0
Params: count 1000000  min_pkt_size: 64  max_pkt_size: 64
     frags: 0  delay: 0  clone_skb: 0  ifname: eth0
     flows: 0 flowlen: 0
...
Current:
     pkts-sofar: 1000000  errors: 0
     started: 1534853256071us  stopped: 1534861576098us idle: 70673us
...
Result: OK: 8320027(c8249354+d70673) usec, 1000000 (64byte,0frags)
  120191pps 61Mb/sec (61537792bps) errors: 0
```

测试报告主要分为三个部分：

- 第一部分的 Params 是测试选项；

- 第二部分的 Current 是测试进度，其中，packts so far（pkts-sofar）表示已经发送了 100 万个包，也就表明测试已完成。

- 第三部分的 Result 是测试结果，包含测试所用时间、网络包数量和分片、PPS、吞吐量以及错误数。

PPS，是 Packet Per Second（包 / 秒）的缩写
