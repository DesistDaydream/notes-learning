---
title: Thanos 管理
---

# Thanos Store 可能产生的问题

下面报警触发的原因未知，但是在 node-exporter 的面板中，打开系统明细，并查询 90 天数据，且 receive 只保留 30 天数据时，大概率会发生这个问题。

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": []
  },
  "warnings": [
    "No StoreAPIs matched for this query",
    "No StoreAPIs matched for this query"
  ]
}
```
