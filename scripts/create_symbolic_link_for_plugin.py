"""同步 Obsidian 插件到多个项目，通过创建符号链接实现。

Windows 上创建符号链接需要管理员权限，或者在“设置 -> 更新和安全 -> 开发者选项”
中开启开发者模式，否则 symlink_to() 会抛出 OSError（权限不足）。
"""

import sys
from pathlib import Path
from typing import Optional


def create_symbolic_links(
    source_dir: Path,
    target_dir: Path,
    pattern: str = "*",
    exclude: Optional[list[str]] = None,
) -> None:
    """将 source_dir 下匹配 pattern 的文件，以符号链接形式创建到 target_dir。

    exclude: 需要跳过的文件名列表（按文件名精确匹配，不支持通配符）。
    """
    exclude = exclude or []

    if not source_dir.exists():
        print(f"源目录不存在: {source_dir}", file=sys.stderr)
        return

    target_dir.mkdir(parents=True, exist_ok=True)

    files = [f for f in source_dir.glob(pattern) if f.is_file() and f.name not in exclude]

    if not files:
        print(f"未找到匹配的文件: {source_dir}")
        return

    print(f"找到 {len(files)} 个文件")

    for file in files:
        link_path = target_dir / file.name

        if link_path.exists() or link_path.is_symlink():
            print(f"链接已存在，跳过: {file.name}")
            continue

        try:
            link_path.symlink_to(file)
            print(f"✓ 已创建链接: {file.name}")
        except OSError as e:
            print(f"创建链接失败 {file.name}: {e}", file=sys.stderr)


def read_plugins(file_path: Path) -> list[str]:
    """从文件中读取插件列表，每行一个插件名，跳过空行和空白。"""
    if not file_path.exists():
        print(f"插件列表文件不存在: {file_path}", file=sys.stderr)
        return []

    lines = file_path.read_text(encoding="utf-8").splitlines()
    return [line.strip() for line in lines if line.strip()]


def main() -> None:
    base_dir = Path(r"D:\projects\DesistDaydream")
    learning_plugins = base_dir / "notes-learning" / "content" / "zh-cn" / ".obsidian" / "plugins"

    plugins_file = Path(__file__).parent / "plugins.list"
    plugins = read_plugins(plugins_file)

    if not plugins:
        print("插件列表为空，退出", file=sys.stderr)
        return

    target_projects = [
        base_dir / "notes-pastime",
        base_dir / "notes-science",
        base_dir / "notes-haohan",
    ]

    exclude_files = ["data.json"]

    for plugin in plugins:
        source_dir = learning_plugins / plugin

        for project in target_projects:
            target_dir = project / ".obsidian" / "plugins" / plugin

            print("\n" + "=" * 40)
            print(f"插件: {plugin}")
            print(f"目标: {project}")
            print("=" * 40)

            create_symbolic_links(source_dir, target_dir, exclude=exclude_files)


if __name__ == "__main__":
    main()