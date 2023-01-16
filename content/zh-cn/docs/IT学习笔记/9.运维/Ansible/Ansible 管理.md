---
title: Ansible 管理
---

#

# 实际案例：

## 批量关闭/开启虚拟机

    - name: 获取虚拟机列表
      virt:
        command: list_vms
      register: info
    - name: 循环开启虚拟机
      virt:
        name: "{{item}}"
        command: start
      loop: "{{info.list_vms}}"
