---
title: Promtail 部署
---

## 使用 docker 运行 Promtail

注意：使用 docker 运行 promtail 的时候，需要注意日志挂载路径及 scrape_configs 配置，因为在容器中运行，所以抓取路径是在容器中，而不是宿主机上，注意这点，否则会抓不到任何日志。

为了解决上面的问题，并且不影响容器内 /var/log 目录，所以启动时指定 -v /var/log:/var/log/host:ro 参数，并修改默认配置文件中 **path**: 的值为 /var/log/host/\*

```bash
docker run -d --rm --name promtail \
  --network host \
  -v /opt/logging/config/promtail:/etc/promtail \
  -v /var/log:/var/log/host:ro \
  -v /etc/localtime:/etc/localtime:ro \
  grafana/promtail
```
