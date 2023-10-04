---
title: "Redhat 包管理"
---

# 概述

# rpm 工具

> 参考：
> 
> - [Manual(手册)，rpm(8)](https://man7.org/linux/man-pages/man8/rpm.8.html)

## Syntax(语法)

**rpm -i \[OPTIONS] PACKAGE** # 安装软件包

OPTIONS
- **-v** # 显示安装过程
- **-h** # 显示安装进度

EXAMPLE
- **rpm -ivh X.rpm** # 安装 X.rpm 软件包
- **rpm -Uvh X.rpm** # 升级 X.rpm 软件包

**rpm -q \[OPTIONS] PACKAGE** # 查询软件包

OPTIONS
- **-a** # 列出所有已经安装在系统上的所有软件包的完整名称
- **-i \<PACKAGE>** # 列出 PACKAGE 这个包的详细信息，安装时间，版本，开发商，描述等等
- **-l \<PACKAGE>** # 列出 PACKAGE 这个包的所有文件与目录所在完整文件名(list)
- **-R \<PACKAGE>** # 列出 PACKAGE 这个包所依赖的文件
- **-f \<FILE>** # 列出该 FILE 属于哪个 PACKAGE 中的文件

EXAMPLE
- rpm -q PACKAGE_NAME
- rpm -qp \*.rpm # 获取当前目录下的 rpm 包相关信息
- rpm -qa | less # 列出所有已安装的软件包
- rpm -qf /usr/sbin/httpd # 查看某个文件属于哪个软件包，可以是普通文件或可执行文件，跟文件的绝对路径
- rpm -qi PACKAGE_NAME # 列出已安装的这个包的标准详细信息
- rpm -ql PACKAGE_NAME # 列出 rpm 包的文件内容
- rpm -q –scripts kernel | less # 列出已安装 rpm 包自带的安装前和安装后脚本
- rpm -qa –queryformat ‘Package %{NAME} was build on %{BUILDHOST}\n’ | less queryformat # 强大的查询
- 可以列出 queryformat 可以使用的所有变量从而组合成更强大的查询
  - rpm –querytags | less

**rpm -e \[OPTIONS] PACKAGE** # 删除软件包

OPTIONS
- **--nodeps** # 不考虑依赖，单独卸载

EXAMPLE
- rpm -e PACKAGE_NAME
- rpm -e –nodeps PACKAGE_NAME # 不考虑依赖包
- rpm -e –allmatches PACKAGE_NAME # 删除所有跟 PACKAGE_NAME 匹配的所有版本的包

**rpm -V \[OPTIONS] PACKAGE** # 验证软件包

OPTIONS
- **-a** # 列出系统上所有可能被更改过的文件

# dnf 工具

dnf 是 yum 的替代品，是 RedHat 系列系统下一代的包管理器。dnf 的命令与 yum 命令基本相同，只是后来实现逻辑不同，所以 dnf 更快，更好用。

在 centos 8 中， yum 命令实际上只是 dnf 的别名，在使用 yum 时，调用的是 dnf 命令。

在 CentOS 7 及之前的版本，yum 官方仓库一直都是 Base 库、Extras 库、Updates 库、centosplus 库等。其中 Base 库是安装 CentOS 时必须提供的仓库，它提供 CentOS 安装(比如可以选择安装桌面环境、开发工具等)、运行以及一些常用用户空间程序。

在 CentOS 8 中发生了变化，原来的 base 库被拆分成两部分：**AppStream** 和 **Base** 库。安装 CentOS 8 时必须提供这两个库。

- CentOS 8 的 Base 库提供安装和运行 CentOS 8 时必须的包，即 CentOS 核心包。这个仓库中全都是 rpm 包。
- CentOS 8 的 AppStream 库提供常用用户空间程序，它们并不一定是安装和运行 CentOS 8 所必须的，比如 Python 包、Perl 包等语言包都在 AppStream。AppStream 中包含 rpm 包和 dnf 的模块。

AppStream 库中的包一般是用户空间程序包，这些程序的更新速度一般比 CentOS 系统更新快的多，将它们单独提取到 AppStream 库，意味着这些程序包和系统相关的包被解绑分开到两个仓库，这可以让系统包和常用程序包分开升级，有利于提供这些程序包的最新的版本。

使用过 CentOS 的人可能都会庆幸这种改变。最近这些年互联网的发展极为迅速，很多程序包的迭代速度也非常快，以前的 CentOS 版本中，很多程序的版本往往都非常古老，要升级这些程序包，只能单独配置它们的镜像仓库(相当于是第三方仓库)，甚至很多程序只能自己编译新版本。

## dnf 配置

# yum 工具

> 参考：
> - Manual(手册),yum(8)

## yum 关联文件与配置

**/etc/yum.repos.d/** # 该目录下是所有源的配置文件，repos 为 repository(仓库)的简称，即 yum 仓库的意思

配置本地 yum 源：cat local.repo（需要将系统 ISO 镜像挂载到/mnt/cdrom 上）

```bash
sudo tee /etc/yum.repos.d/local.repo > /dev/null <<EOF
[local_server]
name=This is a local repo
baseurl=file:///mnt/cdrom
enabled=1
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7
EOF
```

配置远程 yum 源：

在 ISO 挂载目录，启动 80 端口，然后想要使用该源的，配置文件中 baseurl 改为 IP:///PATH 即可

```bash
sudo tee /etc/yum.repos.d/remote.repo > /dev/null <<EOF
[remote_server]
name=This is a local repo
baseurl=http://10.253.26.241/
enabled=1
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7
EOF
```

也可以通过一个可以访问外网的机器，使用 nginx 代理出去

## Syntax(语法)

**yum \[OPTIONS] \[COMMAND] \[PACKAGE ...]**

### OPTIONS

- -t, --tolerant # be tolerant of errors
- -C, --cacheonly # run entirely from system cache, don't update cache
- -c \[config file], --config=\[config file] # config file location
- -R \[minutes], --randomwait=\[minutes] # maximum command wait time
- -d \[debug level], --debuglevel=\[debug level] # debugging output level
- **--showduplicates** # 使用 info、list 和 search 命令中不限制包的最新版本(即显示包的所有版本，不仅仅是最新版本)
- -e \[error level], --errorlevel=\[error level] # error output level
- --rpmverbosity=\[debug level name] # debugging output level for rpm
- -q, --quiet # quiet operation
- **-v, --verbose** # 详细操作
- -y, --assumeyes # answer yes for all questions
- --assumeno # answer no for all questions
- --installroot=\[path] # set install root
- **--enablerepo=REPO** # 激活一个或多个仓库（支持通配符）
- **--disablerepo=REPO** # 禁掉一个或多个仓库（支持通配符）
- **-x \[package], --exclude=\[package]** # 在名称或 glob 不包含包。
- **--disableexcludes=REPO** # 禁掉除了 REPO 这个之外的别的仓库
- --disableincludes=\[repo] # disable includepkgs for a repo or for everything
- --obsoletes # enable obsoletes processing during updates
- --noplugins # disable Yum plugins
- **--nogpgcheck** # 关闭 gpg 签名检查
- --disableplugin=\[plugin] # disable plugins by name
- --enableplugin=\[plugin] # enable plugins by name
- **--skip-broken** # 跳过需要解决问题的包。即.忽略错误，强制安装。
  - 如果安装多个包，其中一个包没有的话，就会停止，使用该选项则会继续安装其他包。
- --color=COLOR # control whether color is used
- --releasever=RELEASEVER # set value of $releasever in yum config and repo files
- **--downloadonly** # 在 yum 一个包时，不安装到系统中，仅下载该包及其依赖到默认的 /var/cache/yum/x86_64/7/REPO/packages/ 目录中。
- **--downloaddir=DLDIR** # 当使用 --downloadonly 参数时，可以使用该参数来指定要下载软件包的下载目录。
- --setopt=SETOPTS # set arbitrary config and repo options
- --bugfix # Include bugfix relevant packages, in updates
- --security # Include security relevant packages, in updates
- --advisory=ADVS, --advisories=ADVS # Include packages needed to fix the given advisory, in updates
- --bzs=BZS # Include packages needed to fix the given BZ, in updates
- --cves=CVES # Include packages needed to fix the given CVE, in updates
- --sec-severity=SEVS, --secseverity=SEVS # Include security relevant packages matching the severity, in updates

### COMMAND

- **autoremove** # 移除已经没有被依赖的软件包
- **autoremove \<Package>** # 删除名为 Package 的包以及其依赖的包
- **clean \<all|headers|packages|metadata|dbcache|plugins|expire-cache|rpmdb>** # 用于清理在 yum 缓存目录中随时间积累的各种东西
- **deplist \<PACKAGE>** # 查看 Package 这个包的依赖关系
  - dependency # 表示 PACKAGE 依赖哪些库
  - provider # 表示 dependency 中的库由哪个包提供
- **groups** # 对一组安装包组执行操作
- **info \<STRING>** # 显示包的详细信息，类似于 rpm -qai，可以使用 string 表达式，模糊搜索包名
- **list \[LIST OPTIONS] \[STRING]** # 显示包的信息，类似于 rpm -qa，可以使用 string 表达式，模糊搜索包名
- **install \<PACKAGE>** # 安装包
- **localinstall** #
  - 注意：install 与 localinstall
- **makecache** # 生成 yum 缓存，以便使用 yum 命令进行安装或者查找的时候，可以从缓存中提取数据进行搜索
- **repolist** # 列出当前所有已经启用的存储库
- **search \<STRING>** # 从包名以及该包的描述中搜索 STRING(字符串)的内容

### list OPTIONS
- **installed** # 列出已安装的包

## EXAMPLE

- 生成已经在 yum.repos.d 目录中源文件的缓存
  - `yum makecache`
- 清理所有缓存
  - `yum clean all`
- 列出 docker-ce 这个软件包的信息
  - `yum info docker-ce`
  - `yum --showduplicates list docker-ce`
- 列出 libcurl 这个包的依赖关系
  - `yum deplist libcurl`
- `yum --disablerepo="*" --enablerepo="elrepo-kernel" list available` #
- 下载 kubectl 的 rpm 包，及其依赖(所依赖的文件取决于当前系统环境，i.e.已经安装的依赖不在其下载的依赖范围内)
  - `yum install --downloadonly --downloaddir=./ kubectl`
- 查看 kubernetes 源下的所有可用的包
  - `yum --disablerepo="*" --enablerepo="kubernetes" --showduplicates list available`
- 禁用 docker-ce-stable 与 kubernetes 仓库并执行升级操作
  - `yum --disablerepo="docker-ce-stable,kubernetes" upgrade`
