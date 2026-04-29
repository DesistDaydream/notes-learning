---
title: "CANN"
linkTitle: "CANN"
created: "2026-04-29T11:37"
weight: 100
---

# 概述

> 参考：
>
> -

**Compute Architecture for Neural Networks(神经网络异构计算架构，简称 CANN)** 是华为昇腾的 AI 计算平台，对标英伟达的 CUDA。

CANN 针对 AI 场景推出的异构计算架构，对上支持多种 AI 框架，对下服务 AI 处理器与编程，发挥承上启下的关键作用，是提升昇腾 AI 处理器计算效率的关键平台

# 学习资料

CANN 社区版资源下载: https://www.hiascend.com/developer/download/community

社区版文档（8.5.0 版本）: https://www.hiascend.com/document/detail/zh/CANNCommunityEdition/850/index/index.html

# 架构

> [!TODO]
> 算子？

# 安装 CANN

> 参考：
>
> - 

官方文档提供了两种安装方式

- [软件安装](https://www.hiascend.com/document/detail/zh/CANNCommunityEdition/850/softwareinst/instg)
- [快速安装](https://www.hiascend.com/cann/download)

在 **软件安装** 中，选择 安装方式（装在 物理机、虚拟机、容器）、安装类型（在线、离线）、操作系统 后，按照文档内容即可非常便捷得安装成功，提供了 wget 下载安装包的方式以及可以直接 copy 的命令。在线使用包管理器安装的话，直接 yum 或 apt 即可，非常得方便。

在 **快速安装** 中可以通过 `lspci -n -D | grep -o '19e5:d[0-9a-f]\{3\}' | head -n1 | cut -d: -f2` 命令获取硬件信息。根据页面提示选择完成后，下面会显示完整的可复制黏贴的命令，一步步照着做即可安装完成。

![300](https://notes-learning.oss-cn-beijing.aliyuncs.com/ai/computing_platform/cann_fast_install_1.png)

## 最佳实践

这里使用 8.5.0 版本做个记录。以 [软件安装](https://www.hiascend.com/document/detail/zh/CANNCommunityEdition/850/softwareinst/instg) 的 “离线安装 （驱动&Toolkit独立包）” 方式安装。

**零、安装 NPU 驱动和固件**

已经手动安装完 [NPU](/docs/0.计算机/NPU.md) 的驱动和固件

**一、安装 Toolkit 开发套件包**

```bash
wget https://ascend-repo.obs.cn-east-2.myhuaweicloud.com/CANN/CANN%208.5.0/Ascend-cann-toolkit_8.5.0_linux-aarch64.run
bash ./Ascend-cann-toolkit_8.5.0_linux-aarch64.run --install
```

应用环境变量

```bash
source /usr/local/Ascend/cann-8.5.0/set_env.sh
```

**二、安装 ops 算子包**

```bash
wget https://ascend-repo.obs.cn-east-2.myhuaweicloud.com/CANN/CANN%208.5.0/Ascend-cann-310p-ops_8.5.0_linux-aarch64.run
bash ./Ascend-cann-310p-ops_8.5.0_linux-aarch64.run --install
```

**（可选）三、安装 NNAL 神经网络加速库**

```bash
wget https://ascend-repo.obs.cn-east-2.myhuaweicloud.com/CANN/CANN%208.5.0/Ascend-cann-nnal_8.5.0_linux-aarch64.run
bash ./Ascend-cann-nnal_8.5.0_linux-aarch64.run --install
```

若要使用 NNAL，加载环境即可：

- ATB加速库适用于大模型场景
  - `source ${HOME}/Ascend/nnal/atb/set_env.sh`
- SiP加速库适用于嵌入式场景
  - `source ${HOME}/Ascend/nnal/asdsip/set_env.sh`

**四、配置环境**

Toolkit 和 ops 默认安装到 /usr/local/Ascend/ 目录，安装完成后都有如下提示：

```bash
===========
= Summary =
===========

Driver:    Installed in /usr/local/Ascend/driver.
ops_310p:  Ascend-cann-310p-ops_8.5.0_linux-aarch64 install success, installed in /usr/local/Ascend.

Please make sure that the environment variables have been configured.
-  To take effect for all users, you can add "source /usr/local/Ascend/cann-8.5.0/set_env.sh" to /etc/profile.
-  To take effect for current user, you can exec command below: source /usr/local/Ascend/cann-8.5.0/set_env.sh or add "source /usr/local/Ascend/cann-8.5.0/set_env.sh" to ~/.bashrc.
```

根据提示，我们将加载环境的的逻辑添加在 profile 中

```bash
tee /etc/profile.d/ascend.sh > /dev/null <<EOF
source /usr/local/Ascend/cann-8.5.0/set_env.sh
EOF
```

根据情况安装一些包

```bash
sudo yum install -y gcc-c++
pip3 install attrs cython 'numpy>=1.19.2,<2.0' decorator sympy cffi pyyaml pathlib2 psutil protobuf==3.20.0 scipy requests absl-py --user
```

执行以下命令，若返回芯片型号（e.g. Ascend310P3），则验证 CANN 安装成功。

```bash
python3 -c "import acl;print(acl.get_soc_name())"
```

**五、其他说明**

CANN 安装完成后，Python 的 site 中会多出来两个路径

```bash
~]# python3 -m site
sys.path = [
    '/root',
    '/usr/local/Ascend/cann-8.5.0/python/site-packages',
    '/usr/local/Ascend/cann-8.5.0/opp/built-in/op_impl/ai_core/tbe',
    '/usr/local/python3.12/lib/python312.zip',
    '/usr/local/python3.12/lib/python3.12',
    '/usr/local/python3.12/lib/python3.12/lib-dynload',
    '/root/.local/lib/python3.12/site-packages',
    '/usr/local/python3.12/lib/python3.12/site-packages',
]
USER_BASE: '/root/.local' (exists)
USER_SITE: '/root/.local/lib/python3.12/site-packages' (exists)
ENABLE_USER_SITE: True
```

# 关联文件与配置

**/usr/local/Ascend/** # 昇腾生态相关产品的程序所在目录、运行时数据保存目录

- **./cann/** # CANN 本体程序和运行时数据保存目录
- **/ascend-toolkit/** # 昇腾工具包和运行时数据保存目录
- **./nnal/** # “NNAL 神经网络加速库” 的本体程序和运行时数据保存目录
