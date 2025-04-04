---
title: OpenStack 部署与清理
linkTitle: OpenStack 部署与清理
weight: 1
---

# 概述

> 参考：
>
> - https://ithelp.ithome.com.tw/articles/10269737
> - https://ithelp.ithome.com.tw/articles/10270784
> - https://ithelp.ithome.com.tw/articles/10271345

## 部署方式

### 手动部署 OpenStack 中的每一个组件

https://docs.openstack.org/xena/install/

### 使用自动化部署工具，部署工具有多种类型可供选择

由大型公司维护的部署工具

- **TripleO** # 由 RedHat 公司维护。是 RedHat OpenStack Platform(RHOSP) 的上游版本。
- **OpenStack Charms** # 由 Canonical 公司维护(Ubuntu 发行版公司)。使用 MAAS 和 Juju 部署。

由社区驱动的部署工具

- **OpenStack Kolla** # 通过 Ansible 部署，并将大部分组件以容器的方式启动
  - 非常简单高效得部署一个用来 生产、开发、测试 的 OpenStack。支持 all-in-one 和 multinode 两种模式(即所有组件都在一个节点或分散在多个节点)
- **OpenStack Ansible** # 通过 Ansible 部署。

其他

**原始的 OpenStack 不管是什么部署方式，通常都需要至少两张网卡。**

对于小规模部署，还可以使用其他 OpenStack 的发行版：

- [MicroStack](/docs/10.云原生/OpenStack/OpenStack%20衍生品/MicroStack.md)

## 支持的操作系统

> 从 Ussuri 版本开始，OpenStack 不再支持 CentOS 7 作为主机操作系统。Train 版本同时支持 CentOS 7 和 8，并提供了迁移路径。有关迁移到 CentOS 8 的信息，请参阅 [Kolla Ansible Train 文档](https://docs.openstack.org/kolla-ansible/train/user/centos8.html)。
>
> 不再支持 CentOS Linux 8（相对于 CentOS Stream 8）作为主机操作系统。Victoria 版本将来会同时支持 CentOS Linux 8 和 CentOS Stream 8，并提供迁移途径。

- CentOS Stream 8
- Debian Bullseye (11)
- RHEL 8（已弃用）
- Rocky Linux 8
- Ubuntu Focal (20.04)

## 支持的容器镜像

为获得最佳结果，基本容器映像分发应与主机操作系统分发匹配。支持以下值 kolla_base_distro：

- centos
- debian
- rhel（已弃用）
- ubuntu

有关哪些发行版支持哪些图像的详细信息，请参阅 [Kolla 支持矩阵](https://docs.openstack.org/kolla/latest/support_matrix)。

# 逐一安装每个组件

想要一个正常可用的 OpenSatck 集群，我们至少需要按照顺序安装下列服务：

- Identity service – [keystone installation for Yoga](https://docs.openstack.org/keystone/yoga/install/)
- Image service – [glance installation for Yoga](https://docs.openstack.org/glance/yoga/install/)
- Placement service – [placement installation for Yoga](https://docs.openstack.org/placement/yoga/install/)
- Compute service – [nova installation for Yoga](https://docs.openstack.org/nova/yoga/install/)
- Networking service – [neutron installation for Yoga](https://docs.openstack.org/neutron/yoga/install/)

还可以安装一些补充服务，比如 Web 页面等：

- Dashboard – [horizon installation for Yoga](https://docs.openstack.org/horizon/yoga/install/)
- Block Storage service – [cinder installation for Yoga](https://docs.openstack.org/cinder/yoga/install/)

# Kolla-ansible

> 参考：
>
> - [GitHub 项目，openstack/kolla](https://github.com/openstack/kolla)

使用 Kolla-ansible 部署 OpenStack 的服务器必须满足如下要求：

- **至少需要两个可用的网卡**，在 `/etc/kolla/globals.yml` 文件中，被描述为如下两个关键字
  - **network_interface** # 管理网络、API 网络的网卡
  - **neutron_external_interface** # Neutron 外部接口就是虚拟机对外访问的出口。该网络设备将会桥接到 `ovs-switch` 这个桥设备上。虚拟机是通过这块网卡访问外网。
- **至少 8G 内存**
- **至少 40G 硬盘**

## 安装依赖并使用虚拟环境

```bash
sudo apt update
sudo apt install python3-dev libffi-dev gcc libssl-dev
```

创建一个 Python 虚拟环境以安装部署工具

```bash
export KOLLA_DIR=/root/kolla
mkdir -p ${KOLLA_DIR}

sudo apt install python3-venv

python3 -m venv ${KOLLA_DIR}/venv
source ${KOLLA_DIR}/venv/bin/activate

pip install -U pip -i https://pypi.tuna.tsinghua.edu.cn/simple

pip install 'ansible<5.0' -i https://pypi.tuna.tsinghua.edu.cn/simple
```

## 安装 Kolla-ansible

这里说的 Kolla-ansible 主要指的是用于部署 Openstack 的 Ansible Playbook~~~~

确定要安装的版本。Kolla-ansible 的版本号与 Openstack 的版本号保持一致，这里以 Openstack 的 `xena` 版本为例

```bash
export KOLLA_BRANCH_NAME=xena
```

使用 pip 安装 kolla-ansible 及其依赖项。

```bash
pip install git+https://opendev.org/openstack/kolla-ansible@${KOLLA_BRANCH_NAME}
```

创建配置目录

```bash
sudo mkdir -p /etc/kolla
sudo chown $USER:$USER /etc/kolla
cp -r ${KOLLA_DIR}/venv/share/kolla-ansible/etc_examples/kolla/* /etc/kolla
```

将 Ansible Playbook 所需的 Inventory 拷贝到当前目录

```bash
cp ${KOLLA_DIR}/venv/share/kolla-ansible/ansible/inventory/* .
```

## 配置 Ansible

```bash
mkdir -p /etc/ansible

tee /etc/ansible/ansible.cfg > /dev/null <<EOF
[defaults]
host_key_checking=False
pipelining=True
forks=100
EOF
```

## All-in-one 部署 OpenStack

### 配置 Kolla

为 `/etc/kolla/passwords.yml` 文件生成密码

```bash
kolla-genpwd
```

配置 `/etc/kolla/globals.yml` 文件

```bash
kolla_base_distro: "ubuntu"
kolla_install_type: "source"
network_interface: "eno3"
# Neutron 外部接口就是虚拟机对外访问的出口。该网络设备将会桥街道 ovs-switch 这个桥设备上。
neutron_external_interface: "eno4"
kolla_internal_vip_address: "192.168.88.236"
enable_cinder: "yes"
openstack_release: "xena"
enable_haproxy: "no"
```

### 配置 Inventory

略，直接使用 localhost 即可

### 部署依赖并检查环境

```bash
kolla-ansible -i ./all-in-one bootstrap-servers
kolla-ansible -i ./all-in-one prechecks
```

### 部署 OpenStack

```bash
kolla-ansible -i ./all-in-one pull ？？？待验证
kolla-ansible -i ./all-in-one deploy
```

## Multinode 部署 OpenStack
