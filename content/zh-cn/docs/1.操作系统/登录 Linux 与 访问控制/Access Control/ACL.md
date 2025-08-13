---
title: ACL
linkTitle: ACL
weight: 2
---

# 概述

> 参考：
>
> - [红帽官方文档，RedHat7-系统管理员指南-第五章.访问控制列表](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/system_administrators_guide/ch-access_control_lists)

**Access Control Lists(访问控制列表，简称 ACL)**。Linux 权限管理是 Linux 很重要的一项内容，重则引起用户信息泄露，轻则导致文件错乱和丢失。企业服务器里有些目录下面的东西暂时保密，不希望别人可以进入目录并查看。有些文件希望别人可以看，但不能删除。有些目录只有 root 等管理员权限的账户才能修改，

Linux 服务器供多个人登录使用，要是没有权限管理就乱了，大家都一样的权限。有些维护系统的命令比较复杂，经验丰富的管理员运行这些命令没事，普通新用户运行的话，可以会导致 Linux 服务器瘫痪。

就像咱们日常生活中，全世界人的权限都一样不就乱了吗。

今天我们来介绍一下 Linux 权限管理的 ACL 权限，它是用户管理结束之后必须要经历的一步。Linux 系统的用户管理包括 Linux 用户和用户组管理之相关配置文件，用户管理的相关配置文件，内容有用户信息文件/etc/passwd，用户密码文件/etc/shadow；用户组信息文件/etc/group，用户组密码文件/etc/gshadow。用户的家目录，以及用户的模板目录； Linux 用户和用户组管理之用户管理命令，管理用户和用户组的命令，包括新建、修改、查看等等以及用的比较多的切换用户命令 su。

下面我们正式开始介绍：

1、什么是 ACL 权限？

比如有如下场景：

某大牛在 QQ 群内直播讲解 Linux 系统的权限管理，讲解完之后，他在一个公有的 Linux 系统中创建了一个 /project 目录，里面存放的是课后参考资料。那么 /project 目录对于大牛而言是所有者，拥有读写可执行（rwx）权限，对于 QQ 群内的所有用户他们都分配的一个所属组里面，也都拥有读写可执行（rwx）权限，而对于 QQ 群外的其他人，那么我们不给他访问/project 目录的任何权限，那么 /project 目录的所有者和所属组权限都是（rwx），其他人权限无。

问题来了，这时候直播有旁听的人参与（不属于 QQ 群内），听完之后，我们允许他访问/project 目录查看参考资料，但是不能进行修改，也就是拥有（r-x）的权限，这时候我们该怎么办呢？我们知道一个文件只能有一个所属组，我们将他分配到 QQ 群所在的所属组内，那么他拥有了写的权限，这是不被允许的；如果将这个旁听的人视为目录/project 的其他人，并且将/project 目录的其他人权限改为（r-x），那么不是旁听的人也能访问我们/project 目录了，这显然也是不被允许的。怎么解决呢？

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788266-55102010-3c01-4999-953f-5264545bba2a.png)

我们想想 windows 系统里面给某个文件分配权限的办法：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788426-4edff74a-9b30-43c9-a8c4-688c5ecadba2.png)

如上图，我们想要让某个用户不具备某个权限，直接不给他分配这个目录的相应权限就行了。那么对应到 Linux 系统也是这样，我们给指定的用户指定目录分配指定的权限，也就是 ACL 权限分配。

## 查看分区 ACL 权限是否开启：dump2fs

我们看某个文件（Linux 系统中目录也是文件，一切皆是文件）是否支持 ACL 权限，首先要看文件所在的分区是否支持 ACL 权限。

①、查看当前系统有哪些分区：df -h

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788270-02fa6473-8187-498a-910b-641f376cd11c.png)

②、查看指定分区详细文件信息：dumpe2fs -h 分区路径

下面是查看 根分区/ 的详细文件信息

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788208-6d15dfed-e275-4812-81cf-66f42bc9965b.png)

回到顶部

### 开启分区 ACL 权限

临时开启分区 ACL 权限

    mount -o remount,acl /

重新挂载根分区，并挂载加入 acl 权限。注意这种命令开启方式，如果系统重启了，那么根分区权限会恢复到初始状态。

永久开启分区 ACL 权限
修改配置文件 /etc/fstab
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788212-4ba77fea-8a55-4639-bea6-6c3ef58096f2.png)
上面是修改根分区拥有 acl 权限

    UUID=490ed737-f8cf-46a6-ac4b-b7735b79fc63 /                       ext4    defaults,acl        1 1

重新挂载文件系统或重启系统，使得修改生效

    mount -o remount /

# ACL 命令行工具

## setfacl - 设定指定文件的 ACL 权限

**setfacl \[-bkndRLPvh] \[{-m|-x} acl_spec] \[{-M|-X} acl_file] file ...**

OPTIONS

- -m # 设定 ACL 权限
- -x # 删除指定的 ACL 权限
- -b # 删除所有的 ACL 权限
- -d # 设定默认 ACL 权限
- -k # 删除默认 ACL 权限
- -R # 递归设定 ACL 权限

EXAMPLE

- setfacl -m u:desistdaydream:rwx test # 让 desistdaydream 这个用户对 test 文件具有 rwx 的权限
- setfacl -m g:desistdaydream:rwx test # 让 desistdaydream 这个组对 test 文件具有 rwx 的权限
- Note：我们给用户或用户组设定 ACL 权限其实并不是真正我们设定的权限，是与 mask 的权限“相与”之后的权限才是用户的真正权限，一般默认 mask 权限都是 rwx，与我们所设定的权限相与就是我们设定的权限。mask 权限下面我们会详细讲解

范例：所有者 root 用户在根目录下创建一个文件目录/project，然后创建一个 QQ 群所属组，所属组里面创建两个用户 zhangsan 和 lisi。所有者和所属组权限和其他人权限是 770。

然后创建一个旁听用户 pt，给他设定/project 目录的 ACL 为 r-x。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788257-439056df-af1a-41c9-913f-cf18513209ba.png)

<https://images2017.cnblogs.com/blog/1120165/201711/1120165-20171109082727169-259824174.png>目录 /project 的所有者和所属组其他人权限设定为 770。接下来我们创建旁听用户 pt，并赋予 acl 权限 rx

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788285-7a684555-e53d-4e5b-96e9-d560c49807f6.png)

为了验证 pt 用户对于 /project 目录没有写权限，我们用 su 命令切换到 pt 用户，然后进入 /project 目录，在此目录下创建文件，看是否能成功：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788260-060edb5b-0ca0-4ee1-a1e7-ba47fccd27a3.png)

上面提示权限不够，说明 acl 权限赋予成功，注意如下所示，如果某个目录或文件下有 + 标志，说明其具有 acl 权限。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788286-6baf14f4-02bd-4091-ab64-13c00e9040d3.png)

## getfacl - 查看 ACL 权限

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788267-41d4fa26-0894-482f-8caa-90ad69c61a69.png)

# 最大有效权限 mask

前面第 4 点我们讲过，我们给用户或用户组设定 ACL 权限其实并不是真正我们设定的权限，是与 mask 的权限“相与”之后的权限才是用户的真正权限，一般默认 mask 权限都是 rwx，与我们所设定的权限相与就是我们设定的权限。

我们通过 getfacl 文件名 也能查看 mask 的权限，那么我们怎么设置呢？

setfacl -m m:权限 文件名

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/wwngi2/1616166788276-cd39efa5-39ce-488d-96f4-9e9b658ceeb0.png)

# 删除 ACL 权限

①、删除指定用户的 ACL 权限

    setfacl -x u:用户名 文件名

②、删除指定用户组的 ACL 权限

    setfacl -x g:组名 文件名

③、删除文件的所有 ACL 权限

    setfacl -b 文件名

8、递归 ACL 权限
通过加上选项 -R 递归设定文件的 ACL 权限，所有的子目录和子文件也会拥有相同的 ACL 权限。

    setfacl -m u:用户名:权限 -R 文件名

9、默认 ACL 权限
如果给父目录设定了默认的 ACL 权限，那么父目录中所有新建的子文件会继承父目录的 ACL 权限。

    setfacl -m d:u:用户名:权限 文件名
