---
title: Windows文件换行符转Linux换行符
---

前段时间，有个朋友碰到由于 Windows 的换行符和 Linux 换行符不一样，导致程序编译不通过。这个问题之前自己也碰到过，网上资料也蛮多，不过还是借此总结总结，因为发现总结+实践的方式能够让自己更好的提升。

操作系统文件换行符

首先介绍下，在 ASCII 中存在这样两个字符 CR（编码为 13）和 LF（编码为 10），在编程中我们一般称其分别为'\r'和'\n'。他们被用来作为换行标志，但在不同系统中换行标志又不一样。下面是不同操作系统采用不同的换行符：

- Unix 和类 Unix（如 Linux）：换行符采用 \n

- Windows 和 MS-DOS：换行符采用 \r\n

- Mac OS X 之前的系统：换行符采用 \r

- Mac OS X：换行符采用 \n

Linux 中查看换行符

      在Linux中查看换行符的方法应该有很多种，这里介绍两种比较常用的方法。

      第一种使用"cat  -A [Filename]" 查看，如下图所示，看到的为一个Windows形式的换行符，\r对应符号^M，\n对应符号$.

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdq7l/1616166269671-c8b5ded5-c271-40cf-b020-da3982fe0e6a.jpeg)

     第二种使用vi编辑器查看，然后使用"set list"命令显示特殊字符：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdq7l/1616166269630-cb766b5e-f716-47f5-b059-38316b6a2b15.jpeg)

      咦，细心的朋友发现了，怎么^M还是没显示出来，这里也是给大家提个醒，用VI的二进制模式（“vi -b [FileName]”）打开，才能够显示出^M：

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/dpdq7l/1616166269631-a080ab28-bbae-4f16-b8fc-4bf5f96c470b.jpeg)

Windows 换行符转换为 Linux 格式

替换方法

1. 第一种使用 VI: 使用 VI 普通模式打开文件，然后运行命令"set ff=unix" 则可以将 Windows 换行符转换为 Linux 换行符，简单吧！命令中 ff 的全称为 file encoding。

2. 使用命令"dos2unix"，如下所示

   1. dos2unix 123.txt

3. 使用 sed 命令删除\r 字符:

   1. sed -i 's/\r//g' gggggggg.txt

4. 使用 windows 版 git 里提供的命令来进行替换，在 git 命令行中进入到要替换文件的目录，执行下面的命令。Node：一定要进入指定目录再执行命令

   1. find . -type f -exec dos2unix {} ;

多文件处理换行符转换

      通常我们都会有一批文件需要替换，比如一个目录的都要替换，我自己写了一个简单的脚本去遍历目录和子目录下的所有文件，并且将其转换为Linux换行格式。代码如下：

```bash
#!/bin/sh
#CheckInput
#Check Whether the input is valid
#0 means not valid
CheckInput()
{
 ret=1;

 #Check the number of parameter
 #And Check whether the argument is a folder
 if [ $# -lt 1 ]
        then
  echo "Please use the command like ./dos2u.sh [Folder]";
  ret=0
 elif [ ! -d $1 ]
 then
  echo "Please use an invalid Folder as the shell argument";
  ret=0
 fi

 return $ret;
}

#TraverseFolder
#Traser all the files under the folder
TraverseFolder()
{
 oldPath=`pwd`
 cd $1;
 for file in `ls`
 do
  if [ -d $file ]
  then
   TraverseFolder $file;
  else
   #echo $file;
   #sed -i 's/\r//g' $file
   dos2unix $file
  fi
 done
 cd $oldPath;
}

CheckInput $*
if [ $ret -ne 1 ]
then
 exit -1
fi
TraverseFolder $1
```

      这个就纯当练习了，应该可以用更简单的方式去解决，比如find命令+dos2unix命令，小伙伴们可以自己试一试。我这么写 主要目的是为了 以后有其他需求更便于扩展，当然还有一些bug要修改啦~~~~^_^。
