---
title: git clone克隆或下载一个仓库单个文件夹
---

#

# 1、如果是想克隆别人的项目或者自己的

很简单的一个网站就解决了。[DownGit](http://downgit.zhoudaxiaa.com)： 只需要找到仓库中对应文件夹的 url,输入之后,点击 download 自动打包下载:

（这里说明一下，因为原作者的项目无法使用，这是我修改过的新项目吧，把资源链接改到了国内 CDN，所以访问速度很快！）

# 2、克隆自己的项目

**注意：本方法会下载整个项目，但是，最后出现在本地项目文件下里只有需要的那个文件夹存在。类似先下载，再过滤。**

**有时候因为需要我们只想 gitclone 下仓库的单个或多个文件夹，而不是全部的仓库内容，这样就很省事，所以下面就开始教程啦**

在 Git1.7.0 以前，这无法实现，但是幸运的是在 Git1.7.0 以后加入了 Sparse Checkout 模式，这使得 Check Out 指定文件或者文件夹成为可能。

**举个例子：**

> 现在有一个**test**仓库<https://github.com/mygithub/test>你要 gitclone 里面的**tt**子目录：在本地的硬盘位置打开**Git Bash**

    git init test && cd test     //新建仓库并进入文件夹
    git config core.sparsecheckout true //设置允许克隆子目录
    echo 'tt*' >> .git/info/sparse-checkout //设置要克隆的仓库的子目录路径   //空格别漏
    git remote add origin git@github.com:mygithub/test.git  //这里换成你要克隆的项目和库
    git pull origin master    //下载
    复制代码

**ok，大功告成！！！**
