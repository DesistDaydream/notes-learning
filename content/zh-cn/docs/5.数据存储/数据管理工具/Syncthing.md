---
title: Syncthing
created: 2026-08-31T11:31
weight: 100
---

# 概述

> 参考：
>
> - [GitHub 项目，syncthing/syncthing](https://github.com/syncthing/syncthing)

Syncthing 可以实时监听本地文件系统的文件，并与远端互相同步到最新。

## Syncthing 与 Rclone 的区别

来自 [Claude AI](https://claude.ai/chat/87f4c383-9634-444a-ba90-a3220a6ee17d)

直接说结论：**这个场景下 Syncthing 更合适，但不是说 rclone 不能用，而是两者设计目标不同。**

**核心区别：**

- **Syncthing** 是为"两台/多台设备之间持续保持文件夹一致"设计的。它是常驻后台服务，实时监听文件变化，检测到改动就自动双向同步，不需要手动触发，也不需要中间的云存储。
- **rclone** 本质是"云存储同步/复制工具"，主要用于本地和云端（S3、OneDrive、Google Drive 等）之间做单向 sync/copy，或者两个远程之间搬数据。它没有内置的"实时监听变化自动同步"能力，正常用法是靠 cron 定时跑 `rclone sync`，一次性单向对齐。

**针对你的场景（两台电脑的 MC 存档+模组配置+JourneyMap 数据）：**

1. 你需要的是双向、持续、自动的同步——今天在电脑 A 玩过之后数据要能"跑"到电脑 B，反过来也一样。Syncthing 原生就是干这个的。
2. rclone 确实有 `rclone bisync`（双向同步）功能，理论上也能做，但官方文档里明确说它比单向 sync 更容易出问题、需要谨慎配置冲突处理，而且同样没有文件变化实时监听，得配合 cron 或者自己写 inotify 触发脚本才能接近"自动"的效果。
3. rclone 的强项在于你要把数据同步到云存储（比如你想额外存一份云端备份到 S3/OneDrive），这种"本地 ↔ 云"场景 rclone 明显更合适；但"本地 ↔ 本地"这种点对点持续同步，Syncthing 开箱即用，配置也更简单（跑起来两边加好设备和文件夹基本就完事了）。

所以如果你只是想让两台电脑的游戏目录保持一致，Syncthing 更省心；如果你还想额外备份一份到云端存储，rclone 可以作为补充（比如再用 rclone 定时把 Syncthing 同步好的目录传一份到云盘做灾备）。

# 关联文件与配置

https://docs.syncthing.net/users/config.html#syncthing-configuration

配置文件与运行时数据都在统一目录中（个人统一称为 SyncthingDir）：

- Windows: %LOCALAPPDATA%/Syncthing/
- Linux: ~/.local/state/syncthing/

**${SyncthingDir}/**

- **./config.xml** # 配置文件
- **./index-\*/** # 当前磁盘上文件的元数据和哈希值以及可从对等方获取的文件的数据库。

## 忽略文件

https://docs.syncthing.net/users/ignoring.html

.stignore 文件类似于 [Git](/docs/2.编程/Programming%20tools/SCM/Git/Git.md) 的 .gitignore 文件，在文件中创建该文件，可以在同步时 <font color="#ff0000">忽略</font> .stignore 中定义的文件。

# 最佳实践

## 在两台 Windows 电脑上，通过 ECS 同步我的世界中的世界及模组配置等

> 个人感受：好像不老好使的。如果 company 启动过游戏，生成了某些文件，这些文件其实是没修改过的原始默认信息，但是同步的时候，因为时间戳更新，所以会覆盖老的已经配置过的文件。必须要保证同步在打开游戏之前进行。同步之后的回退也挺麻烦的。

首先，在各个机器上启动 syncthing 程序

---

在各个机器上执行如下命令，以获取该程序的设备 ID

ECS 上执行

```bash
# 获取本机 Device ID，记下来，另外两台要用
export ECS_ID=$(syncthing cli show system | jq -r .myID)
echo "ECS Device ID: $ECS_ID"
```

Windows 上执行

```powershell
# 获取本机 Device ID，记下来，在 aliyun-ecs 上使用
$MyId = (syncthing.exe cli show system | ConvertFrom-Json).myID
Write-Host "Device ID: $MyId"
```

---

ECS 上执行

```bash
export FOLDER_ID="minecraft-azeroth"
# 分别去 home-pc、company 上执行命令，获取到 Device ID，把结果粘贴到这里
export HOME_ID="<把 home-pc 的 Device ID 粘到这里>"
export COMPANY_ID="<把 company-pc 的 Device ID 粘到这里>"

syncthing cli config devices add --device-id "${HOME_ID}" --name home-pc
syncthing cli config devices add --device-id "${COMPANY_ID}" --name company-pc

syncthing cli config folders add --id "${FOLDER_ID}" --label minecraft-Azeroth --path /opt/HMCL/.minecraft/versions/26.2-Fabric

syncthing cli config folders "${FOLDER_ID}" devices add --device-id "${HOME_ID}"
syncthing cli config folders "${FOLDER_ID}" devices add --device-id "${COMPANY_ID}"

syncthing cli config folders "${FOLDER_ID}" versioning type set simple
syncthing cli config folders "${FOLDER_ID}" versioning params set keep 5

cat > /opt/HMCL/.minecraft/versions/26.2-Fabric/.stignore << 'EOF'
!/saves/Azeroth
!/journeymap
!/config
*
EOF
```

home-pc 上执行

```powershell
$FolderId = "minecraft-azeroth"
$EcsId = "<把 aliyun-ecs 的 Device ID 粘到这里>"
$EcsAddr = "<把 aliyun-ecs 的地址粘到这里>"

syncthing.exe cli config devices add --device-id ${EcsId} --name aliyun-ecs --addresses tcp://${EcsAddr}:22000

syncthing.exe cli config folders add --id ${FolderId} --label minecraft-Azeroth --path "D:\Games\HMCL\.minecraft\versions\26.2-Fabric"

syncthing.exe cli config folders ${FolderId} devices add --device-id ${EcsId}

syncthing.exe cli config folders ${FolderId} versioning type set simple
syncthing.exe cli config folders ${FolderId} versioning params set keep 5

@"
!/saves/Azeroth
!/journeymap
!/config
*
"@ | Set-Content -Path "D:\Games\HMCL\.minecraft\versions\26.2-Fabric\.stignore" -Encoding UTF8
```

company-pc 上执行

```powershell
$FolderId = "minecraft-azeroth"
$EcsId = "<把 aliyun-ecs 的 Device ID 粘到这里>"
$EcsAddr = "<把 aliyun-ecs 的地址粘到这里>"

syncthing.exe cli config devices add --device-id ${EcsId} --name aliyun-ecs --addresses tcp://${EcsAddr}:22000

syncthing.exe cli config folders add --id ${FolderId} --label minecraft-Azeroth --path "D:\Games\HMCL\.minecraft\versions\26.2-Fabric"

syncthing.exe cli config folders ${FolderId} devices add --device-id ${EcsId}

syncthing.exe cli config folders ${FolderId} versioning type set simple
syncthing.exe cli config folders ${FolderId} versioning params set keep 5

@"
!/saves/Azeroth
!/journeymap
!/config
*
"@ | Set-Content -Path "D:\Games\HMCL\.minecraft\versions\26.2-Fabric\.stignore" -Encoding UTF8
```

# 常见问题

## 修改已创建文件夹的路径

早期不支持 https://forum.syncthing.net/t/why-cant-we-change-the-folder-paths/16507

后来发现可以使用 `syncthing.exe cli config folders ${FolderID} path set "${NEW_PATH}"` 命令修改，但是修改后，会出现很多报错，无法同步

```
folder marker missing (this indicates potential data loss, search docs/forum to get information about how to proceed)
```

并且修改后，哪怕删除了再新建，也会出现很多警告，提示 ID 重复相关内容。
