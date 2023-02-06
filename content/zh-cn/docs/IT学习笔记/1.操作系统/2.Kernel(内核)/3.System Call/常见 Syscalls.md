---
title: 常见 Syscalls
---

# 一、进程控制

## a.创建进程

## b.终止进程

## c.载入、执行

## d.获取/设置过程属性

## e.等待时间、等待事件、信号事件

## f.分配和释放内存

# 二、文件管理

## a.创建文件、删除文件

**open("abc", O_WRONLY|O_CREAT|O_NOCTTY|O_NONBLOCK, 0666)** # 创建文件，主要是使用了 O_CREAT 参数
**unlink()** # 删除文件

## b.打开文件、关闭文件

### [open()、openat()、creat()](https://man7.org/linux/man-pages/man2/openat.2.html) # 打开并可能创建一个文件

## c.读、写、调位置

### [read() ](https://man7.org/linux/man-pages/man2/read.2.html)# 从 File Descriptor(文件描述符) 读取

**ssize_t read(int \***fd**\*, void \*\***buf**_, size_t _**count**\*);**

- fd # 文件描述符
- \*buf # 读取/写入的数据的内容(字节流格式)
- count # 读取/写入数据的数据(单位 bytes)

### [write()](https://man7.org/linux/man-pages/man2/write.2.html) # 写入到 File Descriptor(文件描述符)

**ssize_t write(int \***fd**\*, const void \*\***buf**_, size_t _**count**\*);**

- fd # 文件描述符
- \*buf # 读取/写入的数据的内容(字节流格式)
- count # 读取/写入数据的数据(单位 bytes)

### [pread() 与 pwrite()](https://man7.org/linux/man-pages/man2/pread64.2.html) # 以给定的 offset(偏移量) 对给定的 File Descriptor 进行读取或写入数据。

**ssize_t pread(int \***fd**\*, void \*\***buf**_, size_t _**count**_, off_t _**offset**\*);**
**ssize_t pwrite(int \***fd**\*, const void \*\***buf**_, size_t _**count**_, off_t _**offset**\*);**

- fd # 文件描述符
- \*buf # 读取/写入的数据的内容(字节流格式)
- count # 读取/写入数据的数据(单位 bytes)
- offset # 偏移量

成功后，pread() 返回读取的字节数 (返回零表示文件结束)，而 pwrite() 返回写入的字节数。

## d.获取/设置文件属性

### [stat()、fstat()、lastat()、fstatat()](https://man7.org/linux/man-pages/man2/stat.2.html) # 获取文件状态

**int stat(const char _restrict _**pathname***,struct stat *restrict \***statbuf**\*);**

这些获取文件状态的系统调用在 `statbuf` 指向的缓冲区中，返回有关文件的信息

### [getcwd()、getwd()、get_current_dir_name()](https://man7.org/linux/man-pages/man3/getcwd.3.html) # 获取当前工作目录

**char \*getcwd(char \*\***buf**_, size_t _**size**\*);**
**char \*getwd(char \*\***buf**\*);**
**char \*get_current_dir_name(void);**

这些函数返回一个以空字符结尾的字符串，字符串是一个绝对路径名称，该路径名是执行系统调用的进程的当前工作目录。

# 三、设备管理

## a.请求设备、释放设备

## b.读、写、调位位置

## c.获取/设置设备属性

## d.连接或断开设备

# 四、信息维护

## a.获取/设定时间或日期

## b.获取/设置系统数据

## c.获取/设置进程、文件或设备属性

# 五、通信

## a.建立、断开通信

### [socket()](https://man7.org/linux/man-pages/man2/socket.2.html) # 创建一个用于通信的 Endpoint(端点)

在 socketcall() 有注意事项
socket() 返回引用该 endpoint 的文件描述符。成功调用返回的文件描述符将是当前未为该进程打开的编号最小的文件描述符。

### [connect()](https://man7.org/linux/man-pages/man2/connect.2.html) # 在 Socket 上建立一个连接。

在 socketcall() 有注意事项

## b.收发信息

### [sendto()](https://man7.org/linux/man-pages/man2/sendto.2.html) # 发送网络数据。

在 socketcall() 有注意事项

### [recvfrom()](https://man7.org/linux/man-pages/man2/recvfrom.2.html) # 接收网络数据。

在 socketcall() 有注意事项

## c.转移状态信息

## d.连接和断开远程设备

# 六、保护措施

## a.获取/设置权限

### [futex()](https://man7.org/linux/man-pages/man2/futex.2.html) # 快速用户空间锁定

# 待分类总结

ioctl()

pool()
