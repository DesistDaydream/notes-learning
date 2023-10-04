---
title: RustDesk
---

# 概述

> 参考：
> 
> - [GitHub 项目，rustdesk/rustdesk](https://github.com/rustdesk/rustdesk)
> - [官网](https://rustdesk.com/)

# 部署中继器

```bash
sudo docker run -it --name hbbs --net=host --rm \
-p 21115:21115 -p 21116:21116 -p 21116:21116/udp -p 21118:21118 \
-v `pwd`:/root \
rustdesk/rustdesk-server hbbs -r <relay-server-ip[:port]>
```

```bash
sudo docker run -it --net=host --rm --name hbbr \
-p 21117:21117 -p 21119:21119 \
-v `pwd`:/root \
rustdesk/rustdesk-server hbbr
```

## docker compose

```bash
version: '3'

networks:
  rustdesk-net:
    external: false

services:
  hbbs:
    container_name: hbbs
    ports:
      - 21115:21115
      - 21116:21116
      - 21116:21116/udp
      - 21118:21118
    image: rustdesk/rustdesk-server:latest
    command: hbbs -r rustdesk.example.com:21117
    volumes:
      - ./hbbs:/root
    networks:
      - rustdesk-net
    depends_on:
      - hbbr
    restart: unless-stopped

  hbbr:
    container_name: hbbr
    ports:
      - 21117:21117
      - 21119:21119
    image: rustdesk/rustdesk-server:latest
    command: hbbr
    volumes:
      - ./hbbr:/root
    networks:
      - rustdesk-net
    restart: unless-stopped
```
