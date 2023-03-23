---
title: Visual Studio Code
---

# 概述

> 参考：
> - [GitHub 项目，microsoft/vscode](https://github.com/microsoft/vscode)
> - [官网](https://code.visualstudio.com/)
> - [官方文档](https://code.visualstudio.com/docs)

# VS Code 关联文件与配置

**${UserDataDir}** # 用户数据目录。VS Code 运行时生成的持久化数据通常都在同一个目录中。之所以称为用户数据目录，是因为需要以用户为基础来运行一个进程，不同的用户运行的 VS Code，读取的数据应该是不同的。所以这些数据一般就保存在用户的家目录中。
- Windows 默认在 `%USERPROFILE%/AppData/Roaming/Code`
- Linux 默认在 `~/`
  - .**/User/\*** #
    - **./workspaceStorage/\*** # 工作空间的配置与持久化数据
    - **./settings.json** # 用户自定义的配置。默认配置在代码内部。

**${ExtensionsDir}** # 扩展目录。
- Windows 默认在 `%USERPROFILE%/.vscode/extensions`
- Linux 默认在

## 配置

使用快捷键 `Ctrl+Shift+p` ，然后搜索 `setting`，即可看到如下图所示的一些可用的配置。这个编辑的就是 `${UserDataDir}/User/setting.json` 文件
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/rxda5d/1622254287638-dfe1f8a0-03be-45af-af7b-a192a3deb17c.png)

- **首选项：打开默认设置(JSON) # 打开 defaultSettings.json 文件**。这个文件是 VS Code 的默认配置，其中还有每个字段的注释
- **首选项：打开设置(JSON) # 打开 settings.json 文件**。一般用户配置都在这个文件中编写。
- **首选项：打开工作区设置(JSON) # 打开 .vscode/setting 文件**。

**defaultSettings.json** # VS Code 默认配置文件。其中还有每个字段的注释
**settings.json** # VS Code 配置文件。该配置文件分为两类

- **用户配置** # 当前用户使用 VS Code 的配置，用来覆盖默认配置。
  - Windows 目录：**~/AppData/Roaming/Code/User/settings.json** #
    - ![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/rxda5d/1622255508565-70016efb-d366-4869-a851-8d29af6810ab.png)
- **工作区配置** # VS Code 当前工作区的配置，该配置文件通常存在于项目目录中。这个配置会覆盖用户配置以及默认配置。具有最高优先级
  - **$PROJECT/.vscode/settings.json** # $PROJECT 为当前项目的目录。

配置文件有优先级之分，工作区配置 > 用户配置 > 默认配置。高优先级的配置将会覆盖低优先级的配置

上述三个配置文件的优先级从上到下，由低到高。优先级高的配置文件内容将会优先级低的配置文件内容。
