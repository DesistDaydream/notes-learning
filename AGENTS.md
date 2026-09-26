# AGENTS.md

本仓库是个人学习笔记，使用 Hugo + Docsy 构建；笔记源在 `content/zh-cn/`，构建产物发布到 [desistdaydream.github.io](https://github.com/DesistDaydream/desistdaydream.github.io)。

## 授权边界

- 未经明确授权，不得修改、新增或删除 `content/` 下的任何笔记；默认按只读处理。
- 新增或修改任何笔记之前，先取回 `personal/notes-style`，按其约定撰写。
- 临时产物（草稿、调研、脚本、生成的图）一律写入 `agent_tmp/`，不得污染 `content/`。
- 不执行 `git commit` / `git push`；提交与推送由用户自行完成。
- 不主动运行 Hugo 构建做验证；需要时由用户明确要求。

## Obsidian 边界

笔记同时是一个 Obsidian vault：**vault 名** `zh-cn`，**vault 根** `content/zh-cn/`。vault 根落在 `content/` 之内，因此 **CLI 的写入类命令等同于改 `content/`，仍受上面的授权边界约束**。

- 只读命令（可直接用）：`backlinks`、`links`、`unresolved`、`read`、`outline`、`wordcount`、`search:context`、`files`、`folders`、`orphans`、`deadends`、`tags`、`vault`、`version`、`eval`。
- 会改动 GUI 视图或写入 vault 的命令（`open`、`search`、`daily`、`create`、`append`、`prepend`、`move`、`rename`、`delete`、`property:set`、`template:insert` 等）执行前先征得用户同意。