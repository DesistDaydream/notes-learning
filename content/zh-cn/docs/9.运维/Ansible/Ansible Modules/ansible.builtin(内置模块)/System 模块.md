---
title: "System 模块"
linkTitle: "System 模块"
weight: 20
---

# 概述

> 参考：
>
> - [2.9 官方文档，用户指南-使用模块-System 模块](<https://docs.ansible.com/ansible/2.9/modules/list_of_system_modules.html>)

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

# setup - 收集受管理节点的信息

setup 模块在 Ansible 执行时自动运行，收集到的信息会以 [Fact 变量](/docs/9.运维/Ansible/Ansible%20Variables/Fact%20Variables.md)的形式保存。

# systemd - 控制远程主机上以 systemd 运行的服务

官方文档：<https://docs.ansible.com/ansible/latest/collections/ansible/builtin/systemd_module.html>

## 参数

- **name(STRING)** # Unit 的名称
- **state(STRING)** # 设置 Unit 的状态。可用的值有
  - reloaded
  - restarted
  - started
  - stopped
- **enabled(BOOLEAN)** # 设置 Unit 是否应该自启动

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

# user - 管理远程主机上的用户账户

官方文档：https://docs.ansible.com/ansible/latest/collections/ansible/builtin/user_module.html

## 参数

user 模块使用示例：该示例同样适用于更改密码

```yaml
- name: 创建k8s用户
  user:
    name: developer # 指定要创建的用户名
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
