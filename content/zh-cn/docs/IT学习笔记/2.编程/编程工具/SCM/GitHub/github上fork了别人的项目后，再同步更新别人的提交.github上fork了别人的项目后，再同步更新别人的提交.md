---
title: github上fork了别人的项目后，再同步更新别人的提交.github上fork了别人的项目后，再同步更新别人的提交
---

#

我从 github 网站和用 git 命令两种方式说一下。

github 网站上操作

1. 打开自己的仓库，进入 code 下面。

2. 点击 new pull request 创建。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/np6g3g/1616903559723-76f8c4f1-6c02-4145-829f-4b8ddb92de72.jpeg)

1. 选择 base fork

2. 选择 head fork

3. 点击 Create pull request，并填写创建信息。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/np6g3g/1616903559704-28a10f3c-1397-40f9-9c67-0d9bc61da316.jpeg)

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/np6g3g/1616903559698-ff01b1e5-f9b9-406a-938b-3fca1309a5c0.jpeg)

6. 点击 Merge pull request 合并从源 fork 来的代码。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/np6g3g/1616903559689-63e77098-f4a1-4ea4-84c4-f514f4642fec.jpeg)

7. 完成。

用 git 命令操作

1. 用 git remote 查看远程主机状态

git remote -v git remote add upstream git@github.com:xxx/xxx.gitgit fetch upstreamgit merge upstream/mastergit push

- 1

- 2

- 3

- 4

- 5
