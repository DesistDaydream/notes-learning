---
title: "PCB"
created: "2026-08-19T21:05"
weight: 100
---

# 概述

> 参考：
>
> - [Wiki, Printed circuit board](https://en.wikipedia.org/wiki/Printed_circuit_board)

**Printed circuit board(印刷电路板，简称 PCB)**

**Printed Circuit Board Assembly(简称 PCBA)** 是已经在 PCB 上组装好各种元器件，实现了某种功能的 PCB。

# 焊接

> 参考：
>
> - [Wiki, Soldering](https://en.wikipedia.org/wiki/Soldering)

**Solder(焊料)**，锡合金

|      | 熔点（摄氏度）   | 合金比例        |
| ---- | --------- | ----------- |
| 锡铅   | 183       | 63% 锡，37% 铅 |
| 锡银   |           |             |
| 锡铜   | 220 - 227 |             |
| etc. |           |             |

焊接方式

- 电烙铁 # 刀头（中）

**Flux(助焊剂)**

- Rosin(松香)

**吸锡器** # 将液态锡吸走

## 资料

[B 站 - 笛子的日常，《焊武帝养成攻略》零基础三分钟低成本精通焊接 穿越机入门必备技能](https://www.bilibili.com/video/BV1STqpBqEa5)

# 各种 PCBA

## 电源管理 PCBA

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pcb/pcba-power-mgmt.png)


![](https://notes-learning.oss-cn-beijing.aliyuncs.com/pcb/pcba-power-mgmt-by-ai-annotation.png)

11 说明：

|丝印|本质|电压|接到哪里|
|---|---|---|---|
|**D+**|USB 数据线 D+（差分对的一根）|差分信号，不是电源|Micro-USB 的 D+|
|**D−**|USB 数据线 D−（差分对的另一根）|差分信号，不是电源|Micro-USB 的 D−|
|**GND**|公共地 / 0V 参考|0V|USB 的 GND，也是电池负极|
|**5V**|USB 的 VBUS（供电正极）|4.75~5.25V|Micro-USB 的 VBUS|
|**BAT**|电池正极（Battery+）|3.0V(空)~4.2V(满)，标称 3.7V|J1 的 BAT+|
