---
title: Debian 包管理
---

# 概述

> 参考：
> - [Ubuntu 官方文档，软件-包管理](https://ubuntu.com/server/docs/package-management)

Ubuntu 具有一个全面的包管理系统，用于安装、升级、配置和删除软件。除了为我们的 Ubuntu 系统提供超过 60k 个软件包的有组织访问之外，软件包管理工具还具有依赖性解析功能和软件更新检查功能。
Ubuntu 的包管理系统源自使用 Debian GUN/Linux 发行版的系统。Debian 软件包文件通常具有 `.deb` 扩展名，并且可以存储于 **Repositories(存储库)** 中，存储库是网络上或物理媒体(e.g.CD-ROM 光盘)上的软件包集合。软件包通常是预编译的二进制格式。
许多包使用依赖项。依赖项是主包为了正常运行而需要的附加包。例如，语音合成包 Festival 依赖于包 alsa-utils，该包提供了音频播放所需的 ALSA 声音库工具。为了让节日正常运行，必须安装它及其所有依赖项。 Ubuntu 中的软件管理工具会自动执行此操作。

# DPKG 工具集

> 参考：[Wiki,Dpkg](https://en.wikipedia.org/wiki/Dpkg)

## 概述

**Debian Package(Debian 包，简称 dpkg)** 是 Debian 及其衍生的 Linux 发行版的软件包管理程序，用于安装、删除 _.deb 软件包，以及查看 _.deb 软件包的信息。

Dpkg 包含一系列的包管理工具：

- dpkg-deb
- dpkg-split
- dpkg-query
- dpkg-statoverride
- dpkg-divert
- dpkg-trigger

## DPKG 关联文件与配置

## dpkg

> 参考：
> - [Manual(手册),dpkg(1)](https://man7.org/linux/man-pages/man1/dpkg.1.html)

### Syntax(语法)

**dpkg \[OPTIONS] ACTION**
**ACTION**

- **-i, --install PACKAGE** # 安装指定的 PACKAGE。
- **-r，--remove PACKAGE** # 删除指定的已安装的 PACKAGE。保留配置
- **-p, --purge PACKAGE **# 删除指定的已安装的 PACKAGE。连配置也删除

**OPTIONS**

-

## dpkg-query

> 参考：
> - [Manual(手册)，dpkg-query(1)](https://man7.org/linux/man-pages/man1/dpkg-query.1.html)

dpki-query 是一个查询工具，可以从 dpkg 数据库中查询包的信息

### Syntax(语法)

**dpkg-query \[OPTIONS] COMMAND # 安装包查询命令**
**OPTIONS**

- **COMMAND**

- **-L, --listfiles \<PACKAGE>** # 列出系统中 PACKAGE 的安装路径，所有关联的安装文件都会列出
- **-l, --list \[PACKAGE]** # 列出所有包，或指定的 PACKAGE，PACKAGE 可以使用通配符。
- **-S, --search \<FILE>** # 搜索 FILE 属于哪个 Package。FILE 可以使用通配符。
  - 注意：当我们搜索二进制的命令文件属于哪个 Package 时，经常会搜不到，这是因为 which 命令查到的命令路径是 **/usr/bin**，但是 Debian 包安装的程序，通常都是在 **/bin** 目录下，虽然这俩是具有软链接的关系。
  - 所以，当我们搜不到时，可以尝试使用 /bin 目录作为二进制文件的路径前缀。
- **-s, --status \<PackageName>** # 报告所输入的(PackageName)这个软件包的状态

## dpkg 相关程序示例

- 列出 corosync 这个包的安装路径
  - dpkg-query -L corosync
- 查找 nsswitch.conf 文件属于哪个包，并显示所在路径
  - dpkg-query --search nsswitch.conf
- 清除所有已删除包的残余配置文件，可以清除一些残留无用的配置。
  - dpkg -l | grep ^rc|awk '{print $2}' | sudo xargs dpkg -P

# APT 工具集

> 参考：
> - [Wiki](<https://en.wikipedia.org/wiki/APT_(software)>)
> - [APT](<https://en.wikipedia.org/wiki/APT_(software)>)
> - [Debian 项目](https://salsa.debian.org/apt-team/apt)

## 概述

**Advanced Package Tool(高级包工具，简称 APT)** 是一个开源的用户接口，用来处理 Debian、Ubuntu 及其相关 Linux 发行版的软件安装和删除。

不过，比较尴尬的是，APT 没有像 YUM 那种 repolist 查看所有仓库信息的命令，也就无法通过命令来启用或禁用仓库，只能编辑 /etc/apt/ 目录下的 source.list 相关文件来讲仓库注释以便禁用仓库。

APT 中包含如下工具

- **apt-get** # 安装、升级、删除包及其依赖，还可以从经过身份验证的来源检索包装和信息的信息
- **apt-cache** # 查询已安装的包可用信息，以及查询可用的包的信息
- **apt-cdrom** to use removable media as a source for packages
- **apt-config** as an interface to the configuration settings
- **apt-key** as an interface to manage authentication keys
- **apt-extracttemplates** to be used by debconf to prompt for configuration questions before installation
- **apt-ftparchive** creates Packages and other index files needed to publish an archive of deb packages
- **apt-sortpkgs** is a Packages/Sources file normalizer
- **apt** # 一种高级命令行界面，可用于更好的交互式用法

## APT 关联文件

**/etc/apt/\*** # APT 程序配置文件目录

- **./apt.conf** # APT 程序的运行时配置文件。Configuration Item: Dir::Etc::Main.
- **./apt.conf.d/\*** # apt.conf 文件 include(包含) 该目录下的所有文件。Configuration Item: Dir::Etc::Parts.
- **./sources.list** # 包源列表，apt 程序根据这个文件中列表，从对应位置获取包。Configuration Item: Dir::Etc::SourceList.
- **./sources.list.d/\*** # sources.list 文件 include(包含) 该目录下的所有文件。文件以 .list 结尾。Configuration Item: Dir::Etc::SourceParts.
- **/etc/apt/preferences** # Version preferences file. This is where you would specify "pinning", i.e. a preference to get certain packages from a separate source or from a different version of a distribution. Configuration Item: Dir::Etc::Preferences.
- **/etc/apt/preferences.d/\*** # File fragments for the version preferences. Configuration Item: Dir::Etc::PreferencesParts.

**/var/cache/apt/archives/\*** # retrieved(已检索) 的包文件的存储路径。apt 程序的缓存路径。所有下载的 .deb 包都会在这个目录中。Configuration Item: Dir::Cache::Archives.

- 这个目录下的所有 \*.deb 文件会被 **apt clean** 命令清理。在使用 apt install 安装时，被安装的 .deb 包有时候也会存在该目录下
- **./partial/\*** # Storage area for package files in transit. Configuration Item: Dir::Cache::Archives (partial will be implicitly appended)

**/var/lib/apt/lists/\*** # Storage area for state information for each package resource specified in sources.list(5) Configuration Item: Dir::State::Lists.

- **./partial/\*** # Storage area for state information in transit. Configuration Item: Dir::State::Lists (partial will be implicitly appended)

### 配置本地 apt 源

以 20.04 TLS 版本举例

```bash
# 将 iso 文件挂载到本地目录
mount ubuntu-20.04.1-live-server-amd64.iso /mnt/cdrom/

# 配置本地 apt 源
mkdir -p /root/backup && cp -r /etc/apt /root/backup/
sudo tee /etc/apt/sources.list > /dev/null <<EOF
deb file:///mnt/cdrom focal main restricted
EOF
```

其实就是把中间的 `http://PATH` 改为 `file:///PATH/TO/DIR` 即可，

## apt-get

> 参考：
> - [Manual(手册),apt-get(8)](https://manpages.ubuntu.com/manpages/focal/man8/apt-get.8.html)

### Syntax(语法)

**apt-get COMMAND \[OPTIONS]**

> 注意：可以直接使用 apt 命令

**COMMAND**

- **install **# 安装或升级软件包
- **update** # 更新软件包的索引。更新 source.list 文件或长时间没更新时，需要先 update 再安装包。
- **upgrade** # 升级所有软件包。
  - 注意：upgrade 之前必须要执行 `apt update`，以便让 ATP 知道有新版本的软件包可用。
- **reinstall \<PKG>** # 重新安装软件包
- **remove \<PKG>** # 删除指定的软件包
- **purge \<PKG>** # 删除指定得软件包及其配置文件
- **clean** # 清除本地存储库内的已检索到的软件包文件。
  - 该命令会从 /var/cache/apt/archives/ 和 /var/cache/apt/archives/partial/ 目录中删除除了锁文件以外的所有文件。
- **autoremove **# 自动清理所有不再使用的依赖

**OPTIONS**

> Note：所有命令行选项都可以通过配置文件来设置。在配置文件中的写法在 Configuration Item 字段后面

- **-d, --download-only** # 仅下载。包文件只会 retrieved(被检索)，并不会 unpacked(被解包) 或者 installed(被安装)。
  - 配置项：APT::Get::Download-Only.
- **--only-upgrade** # 与 install 命令一起使用时， 将升级指定的软件包，若软件包未安装，则忽略安装行为。
  - 配置项：APT::Get::Only-Upgrade

### EXAMPLE

- 仅下载包及其依赖
  - **apt reinstall -d \`apt-cache depends docker-ce=5:19.03.11~3-0~ubuntu-focal | grep Depends | awk -F': ' '{print $2}'**
- 安装指定版本的包(注意：包名和版本号之间用 `=` 符号分割)
  - **apt install kubectl=1.19.9-00**

## apt-cache

> 参考：
> - [Manul(手册)](https://man.cx/apt-cache)

对 APT 程序生成的包缓存执行各种操作。

### Syntax(语法)

**apt-cache COMMAND \[OPTIONS]**

**COMMAND**

- **madison \<PKG>** # 显示包的可用版本
- **policy [PKG]** # 显示包源的优先级。若指定了 PKG，则显示包的三个信息：1.已安装的版本。2.可安装的最高版本。3.版本列表
- **stats** # 显示缓存中包的统计信息
- 依赖关系命令
  - **depends \<PKG>** # 列出程序包具有的每个依赖关系以及可以满足该依赖关系的所有其他可能的程序包。
  - **rdepends \<PKG>** # 显示了程序包具有的每个反向依赖项的列表。

### EXAMPLE

- 显示 docker-ce=5:20.10.0~3-0~ubuntu-focal 包的信息，可以看到这个包的源、版本等信息。
  - **apt-cache show docker-ce=5:20.10.0~3-0~ubuntu-focal**
- 显示 containerd 包的
  - **apt-cache policy containerd**
- 查询 xxx 包 依赖哪些包
  - **apt-cache depends xxx**
- 查询 xxx 包 被哪些包依赖
  - **apt-cache rdepends xxx**

## apt-key

### EXAMPLE

- 删除仓库的密钥
  - 先列出所有密钥
    - sudo apt-key list
  - 删除指定密钥
    - sudo apt-key del "3820 03C2 C8B7 B4AB 813E 915B 14E4 9429 73C6 2A1B"
    - 也可以仅指定最后 8 个字符
    - sudo apt-key del 73C62A1B
  - # 更新仓库
    - sudo apt update

# 包的自动更新

> 参考：
> - [Ubuntu 官方文档，软件-包管理](https://ubuntu.com/server/docs/package-management)-自动更新(这个小节没有链接)
> - https://wiki.debian.org/UnattendedUpgrades
> - [GitHub 项目，mvo5/unattended-upgrades](https://github.com/mvo5/unattended-upgrades)

APT 工具可以实现软件包的 Automatic Updates(自动更新) 功能，主要依赖于 unattended-upgrades 软件来实现，该软件现在随系统默认安装。该软件安装完成后，会自动启动 **unattended-upgrades.service** 服务。

## Unattended upgrades 关联文件与配置

**/etc/apt/apt.conf.d/\*** #

- **10periodic**
- **20auto-upgrades**
- **50unattended-upgrades**

## 关闭自动更新

```bash
sudo sed -i.bak 's/1/0/' /etc/apt/apt.conf.d/10periodic
sudo sed -i.bak 's/1/0/' /etc/apt/apt.conf.d/20auto-upgrades
sudo systemctl disable unattended-upgrades --now
```
