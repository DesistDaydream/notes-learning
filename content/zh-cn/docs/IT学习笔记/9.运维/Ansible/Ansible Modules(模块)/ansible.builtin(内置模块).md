---
title: ansible.builtin(内置模块)
---

# 概述

> 参考：
> - [官方文档,参考-所有模块和插件的索引-所有模块的索引-ansible.builtin](https://docs.ansible.com/ansible/latest/collections/index_module.html#ansible-builtin)

- [ansible.builtin.add_host](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/add_host_module.html#ansible-collections-ansible-builtin-add-host-module) – Add a host (and alternatively a group) to the ansible-playbook in-memory inventory
- [ansible.builtin.apt](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/apt_module.html#ansible-collections-ansible-builtin-apt-module) – Manages apt-packages
- [ansible.builtin.apt_key](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/apt_key_module.html#ansible-collections-ansible-builtin-apt-key-module) – Add or remove an apt key
- [ansible.builtin.apt_repository](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/apt_repository_module.html#ansible-collections-ansible-builtin-apt-repository-module) – Add and remove APT repositories
- [ansible.builtin.assemble](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/assemble_module.html#ansible-collections-ansible-builtin-assemble-module) – Assemble configuration files from fragments
- [ansible.builtin.assert](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/assert_module.html#ansible-collections-ansible-builtin-assert-module) – Asserts given expressions are true
- [ansible.builtin.async_status](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/async_status_module.html#ansible-collections-ansible-builtin-async-status-module) – Obtain status of asynchronous task
- [ansible.builtin.blockinfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/blockinfile_module.html#ansible-collections-ansible-builtin-blockinfile-module) – Insert/update/remove a text block surrounded by marker lines
- [ansible.builtin.command](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/command_module.html#ansible-collections-ansible-builtin-command-module) – Execute commands on targets
- [ansible.builtin.copy](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/copy_module.html#ansible-collections-ansible-builtin-copy-module) – Copy files to remote locations
- [ansible.builtin.cron](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/cron_module.html#ansible-collections-ansible-builtin-cron-module) – Manage cron.d and crontab entries
- [ansible.builtin.debconf](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/debconf_module.html#ansible-collections-ansible-builtin-debconf-module) – Configure a .deb package
- [ansible.builtin.debug](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/debug_module.html#ansible-collections-ansible-builtin-debug-module) – Print statements during execution
- [ansible.builtin.dnf](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/dnf_module.html#ansible-collections-ansible-builtin-dnf-module) – Manages packages with the dnf package manager
- [ansible.builtin.dpkg_selections](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/dpkg_selections_module.html#ansible-collections-ansible-builtin-dpkg-selections-module) – Dpkg package selection selections
- [ansible.builtin.expect](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/expect_module.html#ansible-collections-ansible-builtin-expect-module) – Executes a command and responds to prompts
- [ansible.builtin.fail](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/fail_module.html#ansible-collections-ansible-builtin-fail-module) – Fail with custom message
- [ansible.builtin.fetch](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/fetch_module.html#ansible-collections-ansible-builtin-fetch-module) – Fetch files from remote nodes
- [ansible.builtin.file](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/file_module.html#ansible-collections-ansible-builtin-file-module) – Manage files and file properties
- [ansible.builtin.find](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/find_module.html#ansible-collections-ansible-builtin-find-module) – Return a list of files based on specific criteria
- [ansible.builtin.gather_facts](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/gather_facts_module.html#ansible-collections-ansible-builtin-gather-facts-module) – Gathers facts about remote hosts
- [ansible.builtin.get_url](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/get_url_module.html#ansible-collections-ansible-builtin-get-url-module) – Downloads files from HTTP, HTTPS, or FTP to node
- [ansible.builtin.getent](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/getent_module.html#ansible-collections-ansible-builtin-getent-module) – A wrapper to the unix getent utility
- [ansible.builtin.git](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/git_module.html#ansible-collections-ansible-builtin-git-module) – Deploy software (or files) from git checkouts
- [ansible.builtin.group](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/group_module.html#ansible-collections-ansible-builtin-group-module) – Add or remove groups
- [ansible.builtin.group_by](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/group_by_module.html#ansible-collections-ansible-builtin-group-by-module) – Create Ansible groups based on facts
- [ansible.builtin.hostname](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/hostname_module.html#ansible-collections-ansible-builtin-hostname-module) – Manage hostname
- [ansible.builtin.import_playbook](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/import_playbook_module.html#ansible-collections-ansible-builtin-import-playbook-module) – Import a playbook
- [ansible.builtin.import_role](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/import_role_module.html#ansible-collections-ansible-builtin-import-role-module) – Import a role into a play
- [ansible.builtin.import_tasks](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/import_tasks_module.html#ansible-collections-ansible-builtin-import-tasks-module) – Import a task list
- [ansible.builtin.include](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_module.html#ansible-collections-ansible-builtin-include-module) – Include a play or task list
- [ansible.builtin.include_role](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_role_module.html#ansible-collections-ansible-builtin-include-role-module) – Load and execute a role
- [ansible.builtin.include_tasks](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_tasks_module.html#ansible-collections-ansible-builtin-include-tasks-module) – Dynamically include a task list
- [ansible.builtin.include_vars](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_vars_module.html#ansible-collections-ansible-builtin-include-vars-module) – Load variables from files, dynamically within a task
- [ansible.builtin.iptables](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/iptables_module.html#ansible-collections-ansible-builtin-iptables-module) – Modify iptables rules
- [ansible.builtin.known_hosts](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/known_hosts_module.html#ansible-collections-ansible-builtin-known-hosts-module) – Add or remove a host from the known_hosts file
- [ansible.builtin.lineinfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/lineinfile_module.html#ansible-collections-ansible-builtin-lineinfile-module) – Manage lines in text files
- [ansible.builtin.meta](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/meta_module.html#ansible-collections-ansible-builtin-meta-module) – Execute Ansible ‘actions’
- [ansible.builtin.package](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/package_module.html#ansible-collections-ansible-builtin-package-module) – Generic OS package manager
- [ansible.builtin.package_facts](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/package_facts_module.html#ansible-collections-ansible-builtin-package-facts-module) – Package information as facts
- [ansible.builtin.pause](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/pause_module.html#ansible-collections-ansible-builtin-pause-module) – Pause playbook execution
- [ansible.builtin.ping](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/ping_module.html#ansible-collections-ansible-builtin-ping-module) – Try to connect to host, verify a usable python and return pong on success
- [ansible.builtin.pip](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/pip_module.html#ansible-collections-ansible-builtin-pip-module) – Manages Python library dependencies
- [ansible.builtin.raw](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/raw_module.html#ansible-collections-ansible-builtin-raw-module) – Executes a low-down and dirty command
- [ansible.builtin.reboot](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/reboot_module.html#ansible-collections-ansible-builtin-reboot-module) – Reboot a machine
- [ansible.builtin.replace](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/replace_module.html#ansible-collections-ansible-builtin-replace-module) – Replace all instances of a particular string in a file using a back-referenced regular expression
- [ansible.builtin.rpm_key](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/rpm_key_module.html#ansible-collections-ansible-builtin-rpm-key-module) – Adds or removes a gpg key from the rpm db
- [ansible.builtin.script](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/script_module.html#ansible-collections-ansible-builtin-script-module) – Runs a local script on a remote node after transferring it
- [ansible.builtin.service](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/service_module.html#ansible-collections-ansible-builtin-service-module) – Manage services
- [ansible.builtin.service_facts](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/service_facts_module.html#ansible-collections-ansible-builtin-service-facts-module) – Return service state information as fact data
- [ansible.builtin.set_fact](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/set_fact_module.html#ansible-collections-ansible-builtin-set-fact-module) – Set host variable(s) and fact(s).
- [ansible.builtin.set_stats](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/set_stats_module.html#ansible-collections-ansible-builtin-set-stats-module) – Define and display stats for the current ansible run
- [ansible.builtin.setup](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/setup_module.html#ansible-collections-ansible-builtin-setup-module) – Gathers facts about remote hosts
- [ansible.builtin.shell](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/shell_module.html#ansible-collections-ansible-builtin-shell-module) – Execute shell commands on targets
- [ansible.builtin.slurp](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/slurp_module.html#ansible-collections-ansible-builtin-slurp-module) – Slurps a file from remote nodes
- [ansible.builtin.stat](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/stat_module.html#ansible-collections-ansible-builtin-stat-module) – Retrieve file or file system status
- [ansible.builtin.subversion](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/subversion_module.html#ansible-collections-ansible-builtin-subversion-module) – Deploys a subversion repository
- [ansible.builtin.systemd](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/systemd_module.html#ansible-collections-ansible-builtin-systemd-module) – Manage systemd units
- [ansible.builtin.sysvinit](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/sysvinit_module.html#ansible-collections-ansible-builtin-sysvinit-module) – Manage SysV services.
- [ansible.builtin.tempfile](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/tempfile_module.html#ansible-collections-ansible-builtin-tempfile-module) – Creates temporary files and directories
- [ansible.builtin.template](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/template_module.html#ansible-collections-ansible-builtin-template-module) – Template a file out to a target host
- [ansible.builtin.unarchive](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/unarchive_module.html#ansible-collections-ansible-builtin-unarchive-module) – Unpacks an archive after (optionally) copying it from the local machine
- [ansible.builtin.uri](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/uri_module.html#ansible-collections-ansible-builtin-uri-module) – Interacts with webservices
- [ansible.builtin.user](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/user_module.html#ansible-collections-ansible-builtin-user-module) – Manage user accounts
- [ansible.builtin.validate_argument_spec](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/validate_argument_spec_module.html#ansible-collections-ansible-builtin-validate-argument-spec-module) – Validate role argument specs.
- [ansible.builtin.wait_for](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/wait_for_module.html#ansible-collections-ansible-builtin-wait-for-module) – Waits for a condition before continuing
- [ansible.builtin.wait_for_connection](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/wait_for_connection_module.html#ansible-collections-ansible-builtin-wait-for-connection-module) – Waits until remote system is reachable/usable
- [ansible.builtin.yum](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/yum_module.html#ansible-collections-ansible-builtin-yum-module) – Manages packages with the yum package manager
- [ansible.builtin.yum_repository](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/yum_repository_module.html#ansible-collections-ansible-builtin-yum-repository-module) – Add or remove YUM repositories

根据 [2.9 版本的模块索引](https://docs.ansible.com/ansible/2.9/modules/modules_by_category.html)文档，我们可以将上面这些内置模块进行类别划分

# Commands # 命令模块

> 官方文档：<https://docs.ansible.com/ansible/2.9/modules/list_of_commands_modules.html>

- [command – Execute commands on targets](https://docs.ansible.com/ansible/2.9/modules/command_module.html#command-module)
- [expect – Executes a command and responds to prompts](https://docs.ansible.com/ansible/2.9/modules/expect_module.html#expect-module)
- [psexec – Runs commands on a remote Windows host based on the PsExec model](https://docs.ansible.com/ansible/2.9/modules/psexec_module.html#psexec-module)
- [raw – Executes a low-down and dirty command](https://docs.ansible.com/ansible/2.9/modules/raw_module.html#raw-module)
- [script – Runs a local script on a remote node after transferring it](https://docs.ansible.com/ansible/2.9/modules/script_module.html#script-module)
- [shell – Execute shell commands on targets](https://docs.ansible.com/ansible/2.9/modules/shell_module.html#shell-module)
- [telnet – Executes a low-down and dirty telnet command](https://docs.ansible.com/ansible/2.9/modules/telnet_module.html#telnet-module)

## [command](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/command_module.html) # 在受管理节点上执行命令

## [script](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/script_module.html) # 将本地脚本传输到受管理节点上并运行

# Files # 文件处理模块

详见 《[Files 类模块](✏IT 学习笔记/🛠️9.运维/Ansible/Ansible%20Modules(模块)/ansible.builtin(内置模块)/Files%20 类模块.md 类模块.md)》

# Packaging # 包模块

> 官方文档：<https://docs.ansible.com/ansible/2.9/modules/list_of_packaging_modules.html>

## Language

- [bower – Manage bower packages with bower](https://docs.ansible.com/ansible/2.9/modules/bower_module.html#bower-module)
- [bundler – Manage Ruby Gem dependencies with Bundler](https://docs.ansible.com/ansible/2.9/modules/bundler_module.html#bundler-module)
- [composer – Dependency Manager for PHP](https://docs.ansible.com/ansible/2.9/modules/composer_module.html#composer-module)
- [cpanm – Manages Perl library dependencies](https://docs.ansible.com/ansible/2.9/modules/cpanm_module.html#cpanm-module)
- [easy_install – Installs Python libraries](https://docs.ansible.com/ansible/2.9/modules/easy_install_module.html#easy-install-module)
- [gem – Manage Ruby gems](https://docs.ansible.com/ansible/2.9/modules/gem_module.html#gem-module)
- [maven_artifact – Downloads an Artifact from a Maven Repository](https://docs.ansible.com/ansible/2.9/modules/maven_artifact_module.html#maven-artifact-module)
- [npm – Manage node.js packages with npm](https://docs.ansible.com/ansible/2.9/modules/npm_module.html#npm-module)
- [pear – Manage pear/pecl packages](https://docs.ansible.com/ansible/2.9/modules/pear_module.html#pear-module)
- [pip – Manages Python library dependencies](https://docs.ansible.com/ansible/2.9/modules/pip_module.html#pip-module)
- [pip_package_info – pip package information](https://docs.ansible.com/ansible/2.9/modules/pip_package_info_module.html#pip-package-info-module)
- [yarn – Manage node.js packages with Yarn](https://docs.ansible.com/ansible/2.9/modules/yarn_module.html#yarn-module)

## Os

- [apk – Manages apk packages](https://docs.ansible.com/ansible/2.9/modules/apk_module.html#apk-module)
- [apt – Manages apt-packages](https://docs.ansible.com/ansible/2.9/modules/apt_module.html#apt-module)
- [apt_key – Add or remove an apt key](https://docs.ansible.com/ansible/2.9/modules/apt_key_module.html#apt-key-module)
- [apt_repo – Manage APT repositories via apt-repo](https://docs.ansible.com/ansible/2.9/modules/apt_repo_module.html#apt-repo-module)
- [apt_repository – Add and remove APT repositories](https://docs.ansible.com/ansible/2.9/modules/apt_repository_module.html#apt-repository-module)
- [apt_rpm – apt_rpm package manager](https://docs.ansible.com/ansible/2.9/modules/apt_rpm_module.html#apt-rpm-module)
- [dnf – Manages packages with the dnf package manager](https://docs.ansible.com/ansible/2.9/modules/dnf_module.html#dnf-module)
- [dpkg_selections – Dpkg package selection selections](https://docs.ansible.com/ansible/2.9/modules/dpkg_selections_module.html#dpkg-selections-module)
- [flatpak – Manage flatpaks](https://docs.ansible.com/ansible/2.9/modules/flatpak_module.html#flatpak-module)
- [flatpak_remote – Manage flatpak repository remotes](https://docs.ansible.com/ansible/2.9/modules/flatpak_remote_module.html#flatpak-remote-module)
- [homebrew – Package manager for Homebrew](https://docs.ansible.com/ansible/2.9/modules/homebrew_module.html#homebrew-module)
- [homebrew_cask – Install/uninstall homebrew casks](https://docs.ansible.com/ansible/2.9/modules/homebrew_cask_module.html#homebrew-cask-module)
- [homebrew_tap – Tap a Homebrew repository](https://docs.ansible.com/ansible/2.9/modules/homebrew_tap_module.html#homebrew-tap-module)
- [installp – Manage packages on AIX](https://docs.ansible.com/ansible/2.9/modules/installp_module.html#installp-module)
- [layman – Manage Gentoo overlays](https://docs.ansible.com/ansible/2.9/modules/layman_module.html#layman-module)
- [macports – Package manager for MacPorts](https://docs.ansible.com/ansible/2.9/modules/macports_module.html#macports-module)
- [openbsd_pkg – Manage packages on OpenBSD](https://docs.ansible.com/ansible/2.9/modules/openbsd_pkg_module.html#openbsd-pkg-module)
- [opkg – Package manager for OpenWrt](https://docs.ansible.com/ansible/2.9/modules/opkg_module.html#opkg-module)
- [package – Generic OS package manager](https://docs.ansible.com/ansible/2.9/modules/package_module.html#package-module)
- [package_facts – package information as facts](https://docs.ansible.com/ansible/2.9/modules/package_facts_module.html#package-facts-module)
- [pacman – Manage packages with pacman](https://docs.ansible.com/ansible/2.9/modules/pacman_module.html#pacman-module)
- [pkg5 – Manages packages with the Solaris 11 Image Packaging System](https://docs.ansible.com/ansible/2.9/modules/pkg5_module.html#pkg5-module)
- [pkg5_publisher – Manages Solaris 11 Image Packaging System publishers](https://docs.ansible.com/ansible/2.9/modules/pkg5_publisher_module.html#pkg5-publisher-module)
- [pkgin – Package manager for SmartOS, NetBSD, et al](https://docs.ansible.com/ansible/2.9/modules/pkgin_module.html#pkgin-module)
- [pkgng – Package manager for FreeBSD >= 9.0](https://docs.ansible.com/ansible/2.9/modules/pkgng_module.html#pkgng-module)
- [pkgutil – Manage CSW-Packages on Solaris](https://docs.ansible.com/ansible/2.9/modules/pkgutil_module.html#pkgutil-module)
- [portage – Package manager for Gentoo](https://docs.ansible.com/ansible/2.9/modules/portage_module.html#portage-module)
- [portinstall – Installing packages from FreeBSD’s ports system](https://docs.ansible.com/ansible/2.9/modules/portinstall_module.html#portinstall-module)
- [pulp_repo – Add or remove Pulp repos from a remote host](https://docs.ansible.com/ansible/2.9/modules/pulp_repo_module.html#pulp-repo-module)
- [redhat_subscription – Manage registration and subscriptions to RHSM using the subscription-manager command](https://docs.ansible.com/ansible/2.9/modules/redhat_subscription_module.html#redhat-subscription-module)
- [rhn_channel – Adds or removes Red Hat software channels](https://docs.ansible.com/ansible/2.9/modules/rhn_channel_module.html#rhn-channel-module)
- [rhn_register – Manage Red Hat Network registration using the rhnreg_ks command](https://docs.ansible.com/ansible/2.9/modules/rhn_register_module.html#rhn-register-module)
- [rhsm_release – Set or Unset RHSM Release version](https://docs.ansible.com/ansible/2.9/modules/rhsm_release_module.html#rhsm-release-module)
- [rhsm_repository – Manage RHSM repositories using the subscription-manager command](https://docs.ansible.com/ansible/2.9/modules/rhsm_repository_module.html#rhsm-repository-module)
- [rpm_key – Adds or removes a gpg key from the rpm db](https://docs.ansible.com/ansible/2.9/modules/rpm_key_module.html#rpm-key-module)
- [slackpkg – Package manager for Slackware >= 12.2](https://docs.ansible.com/ansible/2.9/modules/slackpkg_module.html#slackpkg-module)
- [snap – Manages snaps](https://docs.ansible.com/ansible/2.9/modules/snap_module.html#snap-module)
- [sorcery – Package manager for Source Mage GNU/Linux](https://docs.ansible.com/ansible/2.9/modules/sorcery_module.html#sorcery-module)
- [svr4pkg – Manage Solaris SVR4 packages](https://docs.ansible.com/ansible/2.9/modules/svr4pkg_module.html#svr4pkg-module)
- [swdepot – Manage packages with swdepot package manager (HP-UX)](https://docs.ansible.com/ansible/2.9/modules/swdepot_module.html#swdepot-module)
- [swupd – Manages updates and bundles in ClearLinux systems](https://docs.ansible.com/ansible/2.9/modules/swupd_module.html#swupd-module)
- [urpmi – Urpmi manager](https://docs.ansible.com/ansible/2.9/modules/urpmi_module.html#urpmi-module)
- [xbps – Manage packages with XBPS](https://docs.ansible.com/ansible/2.9/modules/xbps_module.html#xbps-module)
- [yum – Manages packages with the yum package manager](https://docs.ansible.com/ansible/2.9/modules/yum_module.html#yum-module)
- [yum_repository – Add or remove YUM repositories](https://docs.ansible.com/ansible/2.9/modules/yum_repository_module.html#yum-repository-module)
- [zypper – Manage packages on SUSE and openSUSE](https://docs.ansible.com/ansible/2.9/modules/zypper_module.html#zypper-module)
- [zypper_repository – Add and remove Zypper repositories](https://docs.ansible.com/ansible/2.9/modules/zypper_repository_module.html#zypper-repository-module)

## [yum](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/yum_module.html) # 使用主机上的 yum 工具管理包

### 参数

- **name: <STRING>** # 包的名称
- **state: <STRING>** # 指定要执行的操作，是安装还是移除包。可用的值有如下几个：
  - absent 与 removed # 移除指定的包
  - installed 与 present # 简单得确保安装了指定的包
  - latest # 安装最新版本的包，若当前包不是最新版本，则更新它。

### 应用示例

安装 bash-completion 与 vim 的最新版本的包

```bash
ansible all -m yum -a "name=net-bash-completion,vim state=latest"
```

```bash
- name: 安装运维工具
  yum:
    name: ['bash-completion','vim']
    state: latest
```

# System # 系统模块

> 官方文档：<https://docs.ansible.com/ansible/latest/modules/list_of_system_modules.html>

- [aix_devices – Manages AIX devices](https://docs.ansible.com/ansible/2.9/modules/aix_devices_module.html#aix-devices-module)
- [aix_filesystem – Configure LVM and NFS file systems for AIX](https://docs.ansible.com/ansible/2.9/modules/aix_filesystem_module.html#aix-filesystem-module)
- [aix_inittab – Manages the inittab on AIX](https://docs.ansible.com/ansible/2.9/modules/aix_inittab_module.html#aix-inittab-module)
- [aix_lvg – Manage LVM volume groups on AIX](https://docs.ansible.com/ansible/2.9/modules/aix_lvg_module.html#aix-lvg-module)
- [aix_lvol – Configure AIX LVM logical volumes](https://docs.ansible.com/ansible/2.9/modules/aix_lvol_module.html#aix-lvol-module)
- [alternatives – Manages alternative programs for common commands](https://docs.ansible.com/ansible/2.9/modules/alternatives_module.html#alternatives-module)
- [at – Schedule the execution of a command or script file via the at command](https://docs.ansible.com/ansible/2.9/modules/at_module.html#at-module)
- [authorized_key – Adds or removes an SSH authorized key](https://docs.ansible.com/ansible/2.9/modules/authorized_key_module.html#authorized-key-module)
- [awall – Manage awall policies](https://docs.ansible.com/ansible/2.9/modules/awall_module.html#awall-module)
- [beadm – Manage ZFS boot environments on FreeBSD/Solaris/illumos systems](https://docs.ansible.com/ansible/2.9/modules/beadm_module.html#beadm-module)
- [capabilities – Manage Linux capabilities](https://docs.ansible.com/ansible/2.9/modules/capabilities_module.html#capabilities-module)
- [cron – Manage cron.d and crontab entries](https://docs.ansible.com/ansible/2.9/modules/cron_module.html#cron-module)
- [cronvar – Manage variables in crontabs](https://docs.ansible.com/ansible/2.9/modules/cronvar_module.html#cronvar-module)
- [crypttab – Encrypted Linux block devices](https://docs.ansible.com/ansible/2.9/modules/crypttab_module.html#crypttab-module)
- [dconf – Modify and read dconf database](https://docs.ansible.com/ansible/2.9/modules/dconf_module.html#dconf-module)
- [debconf – Configure a .deb package](https://docs.ansible.com/ansible/2.9/modules/debconf_module.html#debconf-module)
- [facter – Runs the discovery program facter on the remote system](https://docs.ansible.com/ansible/2.9/modules/facter_module.html#facter-module)
- [filesystem – Makes a filesystem](https://docs.ansible.com/ansible/2.9/modules/filesystem_module.html#filesystem-module)
- [firewalld – Manage arbitrary ports/services with firewalld](https://docs.ansible.com/ansible/2.9/modules/firewalld_module.html#firewalld-module)
- [gather_facts – Gathers facts about remote hosts](https://docs.ansible.com/ansible/2.9/modules/gather_facts_module.html#gather-facts-module)
- [gconftool2 – Edit GNOME Configurations](https://docs.ansible.com/ansible/2.9/modules/gconftool2_module.html#gconftool2-module)
- [getent – A wrapper to the unix getent utility](https://docs.ansible.com/ansible/2.9/modules/getent_module.html#getent-module)
- [group – Add or remove groups](https://docs.ansible.com/ansible/2.9/modules/group_module.html#group-module)
- [hostname – Manage hostname](https://docs.ansible.com/ansible/2.9/modules/hostname_module.html#hostname-module)
- [interfaces_file – Tweak settings in /etc/network/interfaces files](https://docs.ansible.com/ansible/2.9/modules/interfaces_file_module.html#interfaces-file-module)
- [iptables – Modify iptables rules](https://docs.ansible.com/ansible/2.9/modules/iptables_module.html#iptables-module)
- [java_cert – Uses keytool to import/remove key from java keystore (cacerts)](https://docs.ansible.com/ansible/2.9/modules/java_cert_module.html#java-cert-module)
- [java_keystore – Create or delete a Java keystore in JKS format](https://docs.ansible.com/ansible/2.9/modules/java_keystore_module.html#java-keystore-module)
- [kernel_blacklist – Blacklist kernel modules](https://docs.ansible.com/ansible/2.9/modules/kernel_blacklist_module.html#kernel-blacklist-module)
- [known_hosts – Add or remove a host from the known_hosts file](https://docs.ansible.com/ansible/2.9/modules/known_hosts_module.html#known-hosts-module)
- [listen_ports_facts – Gather facts on processes listening on TCP and UDP ports](https://docs.ansible.com/ansible/2.9/modules/listen_ports_facts_module.html#listen-ports-facts-module)
- [locale_gen – Creates or removes locales](https://docs.ansible.com/ansible/2.9/modules/locale_gen_module.html#locale-gen-module)
- [lvg – Configure LVM volume groups](https://docs.ansible.com/ansible/2.9/modules/lvg_module.html#lvg-module)
- [lvol – Configure LVM logical volumes](https://docs.ansible.com/ansible/2.9/modules/lvol_module.html#lvol-module)
- [make – Run targets in a Makefile](https://docs.ansible.com/ansible/2.9/modules/make_module.html#make-module)
- [mksysb – Generates AIX mksysb rootvg backups](https://docs.ansible.com/ansible/2.9/modules/mksysb_module.html#mksysb-module)
- [modprobe – Load or unload kernel modules](https://docs.ansible.com/ansible/2.9/modules/modprobe_module.html#modprobe-module)
- [mount – Control active and configured mount points](https://docs.ansible.com/ansible/2.9/modules/mount_module.html#mount-module)
- [nosh – Manage services with nosh](https://docs.ansible.com/ansible/2.9/modules/nosh_module.html#nosh-module)
- [ohai – Returns inventory data from Ohai](https://docs.ansible.com/ansible/2.9/modules/ohai_module.html#ohai-module)
- [open_iscsi – Manage iSCSI targets with Open-iSCSI](https://docs.ansible.com/ansible/2.9/modules/open_iscsi_module.html#open-iscsi-module)
- [openwrt_init – Manage services on OpenWrt](https://docs.ansible.com/ansible/2.9/modules/openwrt_init_module.html#openwrt-init-module)
- [osx_defaults – Manage macOS user defaults](https://docs.ansible.com/ansible/2.9/modules/osx_defaults_module.html#osx-defaults-module)
- [pam_limits – Modify Linux PAM limits](https://docs.ansible.com/ansible/2.9/modules/pam_limits_module.html#pam-limits-module)
- [pamd – Manage PAM Modules](https://docs.ansible.com/ansible/2.9/modules/pamd_module.html#pamd-module)
- [parted – Configure block device partitions](https://docs.ansible.com/ansible/2.9/modules/parted_module.html#parted-module)
- [pids – Retrieves process IDs list if the process is running otherwise return empty list](https://docs.ansible.com/ansible/2.9/modules/pids_module.html#pids-module)
- [ping – Try to connect to host, verify a usable python and return pong on success](https://docs.ansible.com/ansible/2.9/modules/ping_module.html#ping-module)
- [puppet – Runs puppet](https://docs.ansible.com/ansible/2.9/modules/puppet_module.html#puppet-module)
- [python_requirements_info – Show python path and assert dependency versions](https://docs.ansible.com/ansible/2.9/modules/python_requirements_info_module.html#python-requirements-info-module)
- [reboot – Reboot a machine](https://docs.ansible.com/ansible/2.9/modules/reboot_module.html#reboot-module)
- [runit – Manage runit services](https://docs.ansible.com/ansible/2.9/modules/runit_module.html#runit-module)
- [seboolean – Toggles SELinux booleans](https://docs.ansible.com/ansible/2.9/modules/seboolean_module.html#seboolean-module)
- [sefcontext – Manages SELinux file context mapping definitions](https://docs.ansible.com/ansible/2.9/modules/sefcontext_module.html#sefcontext-module)
- [selinux – Change policy and state of SELinux](https://docs.ansible.com/ansible/2.9/modules/selinux_module.html#selinux-module)
- [selinux_permissive – Change permissive domain in SELinux policy](https://docs.ansible.com/ansible/2.9/modules/selinux_permissive_module.html#selinux-permissive-module)
- [selogin – Manages linux user to SELinux user mapping](https://docs.ansible.com/ansible/2.9/modules/selogin_module.html#selogin-module)
- [seport – Manages SELinux network port type definitions](https://docs.ansible.com/ansible/2.9/modules/seport_module.html#seport-module)
- [service – Manage services](https://docs.ansible.com/ansible/2.9/modules/service_module.html#service-module)
- [service_facts – Return service state information as fact data](https://docs.ansible.com/ansible/2.9/modules/service_facts_module.html#service-facts-module)
- [setup – Gathers facts about remote hosts](https://docs.ansible.com/ansible/2.9/modules/setup_module.html#setup-module)
- [solaris_zone – Manage Solaris zones](https://docs.ansible.com/ansible/2.9/modules/solaris_zone_module.html#solaris-zone-module)
- [svc – Manage daemontools services](https://docs.ansible.com/ansible/2.9/modules/svc_module.html#svc-module)
- [sysctl – Manage entries in sysctl.conf](https://docs.ansible.com/ansible/2.9/modules/sysctl_module.html#sysctl-module)
- [syspatch – Manage OpenBSD system patches](https://docs.ansible.com/ansible/2.9/modules/syspatch_module.html#syspatch-module)
- [systemd – Manage services](https://docs.ansible.com/ansible/2.9/modules/systemd_module.html#systemd-module)
- [sysvinit – Manage SysV services](https://docs.ansible.com/ansible/2.9/modules/sysvinit_module.html#sysvinit-module)
- [timezone – Configure timezone setting](https://docs.ansible.com/ansible/2.9/modules/timezone_module.html#timezone-module)
- [ufw – Manage firewall with UFW](https://docs.ansible.com/ansible/2.9/modules/ufw_module.html#ufw-module)
- [user – Manage user accounts](https://docs.ansible.com/ansible/2.9/modules/user_module.html#user-module)
- [vdo – Module to control VDO](https://docs.ansible.com/ansible/2.9/modules/vdo_module.html#vdo-module)
- [xfconf – Edit XFCE4 Configurations](https://docs.ansible.com/ansible/2.9/modules/xfconf_module.html#xfconf-module)
- [xfs_quota – Manage quotas on XFS filesystems](https://docs.ansible.com/ansible/2.9/modules/xfs_quota_module.html#xfs-quota-module)

## [systemd](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/systemd_module.html) # 控制远程主机上以 systemd 运行的服务

### 参数

- **name: <STRING>** # Unit 的名称
- **state: <STRING>** # 设置 Unit 的状态。可用的值有
  - reloaded
  - restarted
  - started
  - stopped
- **enabled: <BOOLEAN>** # 设置 Unit 是否应该自启动

### 使用示例

```yaml
- name: 启动并设置自启动kubelet与kube-proxy服务
  systemd:
    name: "{{item}}"
    daemon_reload: yes
    state: started
    enabled: yes
  with_items:
    - kubelet
    - kube-proxy
```

## [user](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/user_module.html) # 管理远程主机上的用户账户

### 参数

user 模块使用示例：该示例同样适用于更改密码

```yaml
- name: 创建k8s用户
  user:
    name: developer #指定要创建的用户名
    password: "$6$mysecretsalt$QjSLl.VQoxPKJkBE9.oLX82C5P4tAMH8UfFRpkxgkqSg2GNob8Y39hj5/cl7o0gbpPXVBGaB9oLuCPfVhIhyA0" # 使用下面Note中的命令来获取加密后的密码
- name: 同时更改多个用户的密码
    user:
      name: "{{ item.name }}"
      password: "{{ item.chpass | password_hash('sha512') }}" # 也可以直接使用明文作为密码
      update_password: always
    with_items:
    - { name: 'root', chpass: 'admin#123' }
    - { name: 'test', chpass: 'yjun@123' }
```

Note：生成加密密码的方式

- ansible all -i localhost, -m debug -a "msg={{ 'mypassword' | password\_hash('sha512', 'mysecretsalt') }}"
  - 命令中的 mypassword 就是想要使用的密码，mysecretsalt 则是密码学中加的盐，详见<https://zh.wikipedia.org/wiki/%E7%9B%90_(%E5%AF%86%E7%A0%81%E5%AD%A6)>
  - 将输出信息引号内的部分直接当做 password 的值即可

# Utilities 实用程序模块

## Helper

- [meta – Execute Ansible ‘actions’](https://docs.ansible.com/ansible/2.9/modules/meta_module.html#meta-module)

## Logic

- [assert – Asserts given expressions are true](https://docs.ansible.com/ansible/2.9/modules/assert_module.html#assert-module)
- [async_status – Obtain status of asynchronous task](https://docs.ansible.com/ansible/2.9/modules/async_status_module.html#async-status-module)
- [debug – Print statements during execution](https://docs.ansible.com/ansible/2.9/modules/debug_module.html#debug-module)
- [fail – Fail with custom message](https://docs.ansible.com/ansible/2.9/modules/fail_module.html#fail-module)
- [import_playbook – Import a playbook](https://docs.ansible.com/ansible/2.9/modules/import_playbook_module.html#import-playbook-module)
- [import_role – Import a role into a play](https://docs.ansible.com/ansible/2.9/modules/import_role_module.html#import-role-module)
- [import_tasks – Import a task list](https://docs.ansible.com/ansible/2.9/modules/import_tasks_module.html#import-tasks-module)
- [include – Include a play or task list](https://docs.ansible.com/ansible/2.9/modules/include_module.html#include-module)
- [include_role – Load and execute a role](https://docs.ansible.com/ansible/2.9/modules/include_role_module.html#include-role-module)
- [include_tasks – Dynamically include a task list](https://docs.ansible.com/ansible/2.9/modules/include_tasks_module.html#include-tasks-module)
- [include_vars – Load variables from files, dynamically within a task](https://docs.ansible.com/ansible/2.9/modules/include_vars_module.html#include-vars-module)
- [pause – Pause playbook execution](https://docs.ansible.com/ansible/2.9/modules/pause_module.html#pause-module)
- [set_fact – Set host facts from a task](https://docs.ansible.com/ansible/2.9/modules/set_fact_module.html#set-fact-module)
- [set_stats – Set stats for the current ansible run](https://docs.ansible.com/ansible/2.9/modules/set_stats_module.html#set-stats-module)
- [wait_for – Waits for a condition before continuing](https://docs.ansible.com/ansible/2.9/modules/wait_for_module.html#wait-for-module)
- [wait_for_connection – Waits until remote system is reachable/usable](https://docs.ansible.com/ansible/2.9/modules/wait_for_connection_module.html#wait-for-connection-module)

## [debug](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/debug_module.html) # 输出变量或表达式

该模块可以在 playbook 执行期间，输出指定的内容，而不用停止 playbook。常用于调试变量或者表达式。比如使用 shell 模块的时候，可以通过 debug 模块来输出 shell 模块定义的语句的执行结果

常与 when 指令一起使用。

```yaml
- name: Print the gateway for each host when defined
  ansible.builtin.debug:
    msg: System {{ inventory_hostname }} has gateway {{ ansible_default_ipv4.gateway }}
  when: ansible_default_ipv4.gateway is defined

- name: Get uptime information
  ansible.builtin.shell: /usr/bin/uptime
  register: result

- name: Print return information from the previous task
  ansible.builtin.debug:
    var: result
    verbosity: 2

- name: Display all variables/facts known for a host
  ansible.builtin.debug:
    var: hostvars[inventory_hostname]
    verbosity: 4

- name: Prints two lines of messages, but only if there is an environment value set
  ansible.builtin.debug:
    msg:
      - "Provisioning based on YOUR_KEY which is: {{ lookup('env', 'YOUR_KEY') }}"
      - "These servers were built using the password of '{{ password_used }}'. Please retain this for later use."
```

## [fail](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/fail_module.html) # 终止任务的执行，并输出自定义的消息

## import 与 include 模块 # 在主任务中导入其他任务或变量

> 参考：
> - <https://www.cnblogs.com/mauricewei/p/10054041.html>
> - https://docs.ansible.com/ansible/2.5/user\_guide/playbooks\_reuse.html#differences-between-static-and-dynamic
> - https://docs.ansible.com/ansible/2.5/user\_guide/playbooks\_conditionals.html#applying-when-to-roles-imports-and-includes

随着要管理的服务不断增多，我们的 Playbook 也会越来越大，内容越来越多，管理起来也会随着复杂。这时，我们可以将某些 tasks 分散到很多文件中，通过 import 和 include 相关模块，实现 tasks 文件之间的相互调用。

说白了，就是聚合多个文件的 tasks。

一共有如下几种模块可以实现这类功能：

- import_playbook
- [import_role](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/import_role_module.html) #
- import_tasks # **Static(静态)** 方法，在 playbooks 解析阶段将父 task 变量和子 task 变量全部读取并加载
- include
- include_tasks # **Dynamic(动态)** 方法，在执行 play 之前才会加载自己变量
- include_vars

### 动态与静态导入 tasks 的区别

#### 区别一

- import_tasks(Static) 方法会在 playbooks 解析阶段将父 task 变量和子 task 变量全部读取并加载
- include_tasks(Dynamic) 方法则是在执行 play 之前才会加载自己变量

可能有点懵，举例说明
**下面两个例子是 test.yml 里的 task 调用 test2.yml，不同之处是一个使用 import_tasks 另一个使用 include_tasks**&#x20;
import_tasks：在执行 tasks 之前，ansible 解释器会先加载 test.yml 里的变量同时再加载 test2.yml 里的变量，那么 ansible_os_family 变量会有一个覆盖现象产生，最终的参数应为“BlackHat”，所以当 test.yml 里执行 when 语句时，ansible_os_family 被判定为“BlackHat”，when 的判断结果为 false，也就不会调用 test2.yml 了

执行结果如下，test2.yml 里的 task 都被 skip 了：

include_tasks：ansible 会在完全执行完 test.yml 里的 task 后才会加载 test2.yml 里的变量，所以当执行 when 语句时，ansible_os_family 的参数应为“RedHat”，此时 when 语句判断结果是 true，也就是 test2.yml 里的 tasks 会被执行。
将第一张图 test.yaml 里的 import_tasks 换成 include_tasks，执行结果如下：

我们发现自 test2.yml 里的 task 被执行了，并且在 test2.yml 里 ansible_os_family 的参数变为了“BlackHat”。
这就是 include_tasts 和 import_tasks 方法的第一个区别。

#### 区别二

- include_tasks 方法调用的文件名称可以加变量
- import_tasks 方法调用的文件名称不可以有变量

这个区别比较简单，直接上示例：
当调用的文件名称有变量时，使用 include_tasks 方法：

能够正常调用 test2.yml，执行结果如下：
&#x20;
当使用 import_tasks 方法时，执行报错。
ansible 也给出了错误原因，当使用 static include 时，是不能使用变量的：
&#x20;
这就是 include_tasts 和 import_tasks 方法的第二个区别。
