---
title: Reflection
linkTitle: Reflection
date: 2023-11-20T21:30
weight: 10
---

# 概述

> 参考：
> 
> - [Go 包，标准库-reflect](https://pkg.go.dev/reflect)
> - [博客园-Stefno，深度解密 Go 语言之反射](https://www.cnblogs.com/qcrao-2018/p/10822655.html)

**Reflection(反射)** 是用程序检查其所拥有的结构，尤其是类型的一种能力。可以通过反射来分析一个结构体。反射可以在运行时检查类型和变量，例如它的大小、方法和动态的调用这些方法。
`reflect.TypeOf()`和`reflect.ValueOf()`两个函数返回被检查对象的类型和值。e.g.`var x float64 = 3.4`,那么 reflect.TypeOf(x)返回 float64，reflect.ValueOf(x)返回 x 的值
忠告：反射是一个强大并富有表达力的工具，但是它应该被小心地使用，原因有三：

1. 基于反射的代码是比较脆弱的。对于每一个会导致编译器报告类型错误的问题，在反射中都有与之相对应的误用问题，不同的是编译器会在构建时马上报告错误，而反射则是在真正运行到的时候才会抛出 panic 异常，可能是写完代码很久之后了，而且程序也可能运行了很长的时间。
2. 即使对应类型提供了相同文档，但是反射的操作不能做静态类型检查，而且大量反射的代码通常难以理解。总是需要小心翼翼地为每个导出的类型和其它接受 interface{}或 reflect.Value 类型参数的函数维护说明文档。
3. 基于反射的代码通常比正常的代码运行速度慢一到两个数量级。对于一个典型的项目，大部分函数的性能和程序的整体性能关系不大，所以使用反射可能会使程序更加清晰。测试是一个特别适合使用反射的场景，因为每个测试的数据集都很小。但是对于性能关键路径的函数，最好避免使用反射。

# 什么是反射

反射和 Interface 息息相关，而 Interface 是我们[上一篇文章](https://mp.weixin.qq.com/s/EbxkBokYBajkCR-MazL0ZA)的内容。在开始正文前，和大家说点题外话。
直接看维基百科上的定义：

> 在计算机科学中，反射是指计算机程序在运行时（Run time）可以访问、检测和修改它本身状态或行为的一种能力。用比喻来说，反射就是程序在运行的时候能够“观察”并且修改自己的行为。

那我就要问个问题了：不用反射就不能在运行时访问、检测和修改它本身的状态和行为吗？
问题的回答，其实要首先理解什么叫访问、检测和修改它本身状态或行为，它的本质是什么？
实际上，它的本质是程序在运行期探知对象的类型信息和内存结构，不用反射能行吗？可以的！使用汇编语言，直接和内层打交道，什么信息不能获取？但是，当编程迁移到高级语言上来之后，就不行了！就只能通过`反射`来达到此项技能。
不同语言的反射模型不尽相同，有些语言还不支持反射。《Go 语言圣经》中是这样定义反射的：

> Go 语言提供了一种机制在运行时更新变量和检查它们的值、调用它们的方法，但是在编译时并不知道这些变量的具体类型，这称为反射机制。

# 为什么要用反射

需要反射的 2 个常见场景：

1. 有时你需要编写一个函数，但是并不知道传给你的参数类型是什么，可能是没约定好；也可能是传入的类型很多，这些类型并不能统一表示。这时反射就会用的上了。
2. 有时候需要根据某些条件决定调用哪个函数，比如根据用户的输入来决定。这时就需要对函数和函数的参数进行反射，在运行期间动态地执行函数。

在讲反射的原理以及如何用之前，还是说几点不使用反射的理由：

1. 与反射相关的代码，经常是难以阅读的。在软件工程中，代码可读性也是一个非常重要的指标。
2. Go 语言作为一门静态语言，编码过程中，编译器能提前发现一些类型错误，但是对于反射代码是无能为力的。所以包含反射相关的代码，很可能会运行很久，才会出错，这时候经常是直接 panic，可能会造成严重的后果。
3. 反射对性能影响还是比较大的，比正常代码运行速度慢一到两个数量级。所以，对于一个项目中处于运行效率关键位置的代码，尽量避免使用反射特性。

# 反射是如何实现的

上一篇文章讲到了 `interface`，它是 Go 语言实现抽象的一个非常强大的工具。当向接口变量赋予一个实体类型的时候，接口会存储实体的类型信息，反射就是通过接口的类型信息实现的，反射建立在类型的基础上。
Go 语言在 reflect 包里定义了各种类型，实现了反射的各种函数，通过它们可以在运行时检测类型的信息、改变类型的值。

## types 和 interface

Go 语言中，每个变量都有一个静态类型，在编译阶段就确定了的，比如 `int, float64, []int` 等等。注意，这个类型是声明时候的类型，不是底层数据类型。
Go 官方博客里就举了一个例子：

```go
type MyInt int
var i int
var j MyInt
```

尽管 i，j 的底层类型都是 int，但我们知道，他们是不同的静态类型，除非进行类型转换，否则，i 和 j 不能同时出现在等号两侧。j 的静态类型就是 `MyInt`。
反射主要与 interface{} 类型相关。前面一篇关于 interface 相关的文章已经探讨过 interface 的底层结构，这里再来复习一下。

```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
type itab struct {
    inter  *interfacetype
    _type  *_type
    link   *itab
    hash   uint32
    bad    bool
    inhash bool
    unused [2]byte
    fun    [1]uintptr
}
```

其中 `itab` 由具体类型 `_type` 以及 `interfacetype` 组成。`_type` 表示具体类型，而 `interfacetype` 则表示具体类型实现的接口类型。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149620-8d9ae6eb-1f40-4451-af9c-68139cf480c4.png)
实际上，iface 描述的是非空接口，它包含方法；与之相对的是 `eface`，描述的是空接口，不包含任何方法，Go 语言里有的类型都 `“实现了”` 空接口。

    type eface struct {
        _type *_type
        data  unsafe.Pointer
    }

相比 `iface`，`eface` 就比较简单了。只维护了一个 `_type` 字段，表示空接口所承载的具体的实体类型。`data` 描述了具体的值。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149625-83be7464-5375-4d6f-a02d-2fc237058398.png)
还是用 Go 官方关于反射的博客里的例子，当然，我会用图形来详细解释，结合两者来看会更清楚。顺便提一下，搞技术的不要害怕英文资料，要想成为技术专家，读英文原始资料是技术提高的一条必经之路。
先明确一点：接口变量可以存储任何实现了接口定义的所有方法的变量。
Go 语言中最常见的就是 `Reader` 和 `Writer` 接口：

    type Reader interface {
        Read(p []byte) (n int, err error)
    }
    type Writer interface {
        Write(p []byte) (n int, err error)
    }

接下来，就是接口之间的各种转换和赋值了：

    var r io.Reader
    tty, err := os.OpenFile("/Users/qcrao/Desktop/test", os.O_RDWR, 0)
    if err != nil {
        return nil, err
    }
    r = tty

首先声明 `r` 的类型是 `io.Reader`，注意，这是 `r` 的静态类型，此时它的动态类型为 `nil`，并且它的动态值也是 `nil`。
之后，`r = tty` 这一语句，将 `r` 的动态类型变成 `*os.File`，动态值则变成非空，表示打开的文件对象。这时，r 可以用`<value, type>`对来表示为： `<tty, *os.File>`。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149687-af8ef79f-2537-497d-bc50-134e4725a0c9.png)
注意看上图，此时虽然 `fun` 所指向的函数只有一个 `Read` 函数，其实 `*os.File` 还包含 `Write` 函数，也就是说 `*os.File` 其实还实现了 `io.Writer` 接口。因此下面的断言语句可以执行：

    var w io.Writer
    w = r.(io.Writer)

之所以用断言，而不能直接赋值，是因为 `r` 的静态类型是 `io.Reader`，并没有实现 `io.Writer` 接口。断言能否成功，看 `r` 的动态类型是否符合要求。
这样，w 也可以表示成 `<tty, *os.File>`，仅管它和 `r` 一样，但是 w 可调用的函数取决于它的静态类型 `io.Writer`，也就是说它只能有这样的调用形式： `w.Write()` 。`w` 的内存形式如下图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149626-01ade0ef-5895-40ae-a674-955cb575928c.png)
和 `r` 相比，仅仅是 `fun` 对应的函数变了：`Read -> Write`。
最后，再来一个赋值：

    var empty interface{}
    empty = w

由于 `empty` 是一个空接口，因此所有的类型都实现了它，w 可以直接赋给它，不需要执行断言操作。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149635-574005ab-9192-4e44-b076-b481e10b091d.png)
从上面的三张图可以看到，interface 包含三部分信息：`_type` 是类型信息，`*data` 指向实际类型的实际值，`itab` 包含实际类型的信息，包括大小、包路径，还包含绑定在类型上的各种方法（图上没有画出方法），补充一下关于 os.File 结构体的图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149693-bbbe2f1b-3087-420f-8ee1-22ddb31eb3e7.png)
这一节的最后，复习一下上一篇关于 interface 的文章，提到的一个技巧，这里再展示一下：
先参考源码，分别定义一个`“伪装”`的 iface 和 eface 结构体。

```go
type iface struct {
    tab  *itab
    data unsafe.Pointer
}
type itab struct {
    inter uintptr
    _type uintptr
    link uintptr
    hash  uint32
    _     [4]byte
    fun   [1]uintptr
}
type eface struct {
    _type uintptr
    data unsafe.Pointer
}
```

接着，将接口变量占据的内存内容强制解释成上面定义的类型，再打印出来：

```go
package main
import (
    "os"
    "fmt"
    "io"
    "unsafe"
)
func main() {
    var r io.Reader
    fmt.Printf("initial r: %T, %v\n", r, r)
    tty, _ := os.OpenFile("/Users/qcrao/Desktop/test", os.O_RDWR, 0)
    fmt.Printf("tty: %T, %v\n", tty, tty)
    // 给 r 赋值
    r = tty
    fmt.Printf("r: %T, %v\n", r, r)
    rIface := (*iface)(unsafe.Pointer(&r))
    fmt.Printf("r: iface.tab._type = %#x, iface.data = %#x\n", rIface.tab._type, rIface.data)
    // 给 w 赋值
    var w io.Writer
    w = r.(io.Writer)
    fmt.Printf("w: %T, %v\n", w, w)
    wIface := (*iface)(unsafe.Pointer(&w))
    fmt.Printf("w: iface.tab._type = %#x, iface.data = %#x\n", wIface.tab._type, wIface.data)
    // 给 empty 赋值
    var empty interface{}
    empty = w
    fmt.Printf("empty: %T, %v\n", empty, empty)
    emptyEface := (*eface)(unsafe.Pointer(&empty))
    fmt.Printf("empty: eface._type = %#x, eface.data = %#x\n", emptyEface._type, emptyEface.data)
}
```

运行结果：

    initial r: <nil>, <nil>
    tty: *os.File, &{0xc4200820f0}
    r: *os.File, &{0xc4200820f0}
    r: iface.tab._type = 0x10bfcc0, iface.data = 0xc420080020
    w: *os.File, &{0xc4200820f0}
    w: iface.tab._type = 0x10bfcc0, iface.data = 0xc420080020
    empty: *os.File, &{0xc4200820f0}
    empty: eface._type = 0x10bfcc0, eface.data = 0xc420080020

`r，w，empty` 的动态类型和动态值都一样。不再详细解释了，结合前面的图可以看得非常清晰。

## 反射的基本函数

reflect 包里定义了一个接口和一个结构体，即 `reflect.Type` 和 `reflect.Value`，它们提供很多函数来获取存储在接口里的类型信息。
`reflect.Type` 主要提供关于类型相关的信息，所以它和 `_type` 关联比较紧密；`reflect.Value` 则结合 `_type` 和 `data` 两者，因此程序员可以获取甚至改变类型的值。
reflect 包中提供了两个基础的关于反射的函数来获取上述的接口和结构体：

```go
func TypeOf(i interface{}) Type
func ValueOf(i interface{}) Value
```

### TypeOf() 函数与 Type{} 接口

`TypeOf` 函数用来提取一个接口中值的类型信息。由于它的输入参数是一个空的 `interface{}`，调用此函数时，实参会先被转化为 `interface{}`类型。这样，实参的类型信息、方法集、值信息都存储到 `interface{}` 变量里了。
看下源码：

```go
func TypeOf(i interface{}) Type {
    eface := *(*emptyInterface)(unsafe.Pointer(&i))
    return toType(eface.typ)
}
```

这里的 `emptyInterface` 和上面提到的 `eface` 是一回事（字段名略有差异，字段是相同的），且在不同的源码包：前者在 `reflect` 包，后者在 `runtime` 包。 `eface.typ` 就是动态类型。

```go
type emptyInterface struct {
    typ  *rtype
    word unsafe.Pointer
}
```

至于 `toType` 函数，只是做了一个类型转换：

```go
func toType(t *rtype) Type {
    if t == nil {
        return nil
    }
    return t
}
```

注意，返回值 `Type` 实际上是一个接口，定义了很多方法，用来获取类型相关的各种信息，而 `*rtype` 实现了 `Type` 接口。

```go
type Type interface {
    // 所有的类型都可以调用下面这些函数
    // 此类型的变量对齐后所占用的字节数
    Align() int

    // 如果是 struct 的字段，对齐后占用的字节数
    FieldAlign() int
    // 返回类型方法集里的第 `i` (传入的参数)个方法
    Method(int) Method
    // 通过名称获取方法
    MethodByName(string) (Method, bool)
    // 获取类型方法集里导出的方法个数
    NumMethod() int
    // 类型名称
    Name() string
    // 返回类型所在的路径，如：encoding/base64
    PkgPath() string
    // 返回类型的大小，和 unsafe.Sizeof 功能类似
    Size() uintptr
    // 返回类型的字符串表示形式
    String() string
    // 返回类型的类型值
    Kind() Kind
    // 类型是否实现了接口 u
    Implements(u Type) bool
    // 是否可以赋值给 u
    AssignableTo(u Type) bool
    // 是否可以类型转换成 u
    ConvertibleTo(u Type) bool
    // 类型是否可以比较
    Comparable() bool
    // 下面这些函数只有特定类型可以调用
    // 如：Key, Elem 两个方法就只能是 Map 类型才能调用

    // 类型所占据的位数
    Bits() int
    // 返回通道的方向，只能是 chan 类型调用
    ChanDir() ChanDir
    // 返回类型是否是可变参数，只能是 func 类型调用
    // 比如 t 是类型 func(x int, y ... float64)
    // 那么 t.IsVariadic() == true
    IsVariadic() bool
    // 返回内部子元素类型，只能由类型 Array, Chan, Map, Ptr, or Slice 调用
    Elem() Type
    // 返回结构体类型的第 i 个字段，只能是结构体类型调用
    // 如果 i 超过了总字段数，就会 panic
    Field(i int) StructField
    // 返回嵌套的结构体的字段
    FieldByIndex(index []int) StructField
    // 通过字段名称获取字段
    FieldByName(name string) (StructField, bool)
    // FieldByNameFunc returns the struct field with a name
    // 返回名称符合 func 函数的字段
    FieldByNameFunc(match func(string) bool) (StructField, bool)
    // 获取函数类型的第 i 个参数的类型
    In(i int) Type
    // 返回 map 的 key 类型，只能由类型 map 调用
    Key() Type
    // 返回 Array 的长度，只能由类型 Array 调用
    Len() int
    // 返回类型字段的数量，只能由类型 Struct 调用
    NumField() int
    // 返回函数类型的输入参数个数
    NumIn() int
    // 返回函数类型的返回值个数
    NumOut() int
    // 返回函数类型的第 i 个值的类型
    Out(i int) Type
    // 返回类型结构体的相同部分
    common() *rtype

    // 返回类型结构体的不同部分
    uncommon() *uncommonType
}
```

可见 `Type` 定义了非常多的方法，通过它们可以获取类型的一切信息，大家一定要完整的过一遍上面所有的方法。
注意到 `Type` 方法集的倒数第二个方法 `common`
返回的 `rtype`类型，它和上一篇文章讲到的 `_type` 是一回事，而且源代码里也注释了：两边要保持同步：

```go
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
```

```go
type rtype struct {
    size       uintptr
    ptrdata    uintptr
    hash       uint32
    tflag      tflag
    align      uint8
    fieldAlign uint8
    kind       uint8
    alg        *typeAlg
    gcdata     *byte
    str        nameOff
    ptrToThis  typeOff
}
```

所有的类型都会包含 `rtype` 这个字段，表示各种类型的公共信息；另外，不同类型包含自己的一些独特的部分。
比如下面的 `arrayType` 和 `chanType` 都包含 `rytpe`，而前者还包含 slice，len 等和数组相关的信息；后者则包含 `dir` 表示通道方向的信息。

```go
// arrayType represents a fixed array type.
type arrayType struct {
    rtype `reflect:"array"`
    elem  *rtype // array element type
    slice *rtype // slice type
    len   uintptr
}
// chanType represents a channel type.
type chanType struct {
    rtype `reflect:"chan"`
    elem  *rtype  // channel element type
    dir   uintptr // channel direction (ChanDir)
}
```

注意到，`Type` 接口实现了 `String()` 函数，满足 `fmt.Stringer` 接口，因此使用 `fmt.Println` 打印的时候，输出的是 `String()` 的结果。另外，`fmt.Printf()` 函数，如果使用 `%T` 来作为格式参数，输出的是 `reflect.TypeOf` 的结果，也就是动态类型。例如：

    fmt.Printf("%T", 3) // int

### ValueOf() 函数与 Value{} 结构体

讲完了 `TypeOf` 函数，再来看一下 `ValueOf` 函数。返回值 `reflect.Value` 表示 `interface{}` 里存储的实际变量，它能提供实际变量的各种信息。相关的方法常常是需要结合类型信息和值信息。例如，如果要提取一个结构体的字段信息，那就需要用到 \_type (具体到这里是指 structType) 类型持有的关于结构体的字段信息、偏移信息，以及 `*data` 所指向的内容 —— 结构体的实际值。
源码如下：

```go
func ValueOf(i interface{}) Value {
    if i == nil {
        return Value{}
    }

    // ……
    return unpackEface(i)
}
// 分解 eface
func unpackEface(i interface{}) Value {
    e := (*emptyInterface)(unsafe.Pointer(&i))
    t := e.typ
    if t == nil {
        return Value{}
    }

    f := flag(t.Kind())
    if ifaceIndir(t) {
        f |= flagIndir
    }
    return Value{t, e.word, f}
}
```

从源码看，比较简单：将先将 `i` 转换成 `*emptyInterface` 类型， 再将它的 `typ` 字段和 `word` 字段以及一个标志位字段组装成一个 `Value` 结构体，而这就是 `ValueOf` 函数的返回值，它包含类型结构体指针、真实数据的地址、标志位。
Value 结构体定义了很多方法，通过这些方法可以直接操作 Value 字段 ptr 所指向的实际数据：

```go
// 设置切片的 len 字段，如果类型不是切片，就会panic
func (v Value) SetLen(n int)

// 设置切片的 cap 字段
func (v Value) SetCap(n int)

// 设置字典的 kv
func (v Value) SetMapIndex(key, val Value)
// 返回切片、字符串、数组的索引 i 处的值
func (v Value) Index(i int) Value

// 根据名称获取结构体的内部字段值
func (v Value) FieldByName(name string) Value

// ……
```

`Value` 字段还有很多其他的方法。例如：

```go
// 用来获取 int 类型的值
func (v Value) Int() int64
// 用来获取结构体字段（成员）数量
func (v Value) NumField() int
// 尝试向通道发送数据（不会阻塞）
func (v Value) TrySend(x reflect.Value) bool
// 通过参数列表 in 调用 v 值所代表的函数（或方法
func (v Value) Call(in []Value) (r []Value)
// 调用变参长度可变的函数
func (v Value) CallSlice(in []Value) []Value
```

不一一列举了，反正是非常多。可以去 `src/reflect/value.go` 去看看源码，搜索 `func (v Value)` 就能看到。

### Type() 方法与 Interface() 方法

另外，通过 `Type()` 方法和 `Interface()` 方法可以打通 `interface`、`Type`、`Value` 三者。Type() 方法也可以返回变量的类型信息，与 `reflect.TypeOf()` 函数等价。`Interface()` 方法可以将 Value 还原成原来的 interface。
这里引用老钱《快学 Go 语言第十五课——反射》的一张图：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149650-1e6b2b6d-0b40-43d1-bf91-a46890012e73.png)
总结一下：`TypeOf()` 函数返回一个接口，这个接口定义了一系列方法，利用这些方法可以获取关于类型的所有信息； `ValueOf()` 函数返回一个结构体变量，包含类型信息以及实际值。
用一张图来串一下：
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149652-cd48e788-6a67-421d-ad95-1d89d7ae5241.png)
上图中，`rtye` 实现了 `Type` 接口，是所有类型的公共部分。emptyface 结构体和 eface 其实是一个东西，而 rtype 其实和 \_type 是一个东西，只是一些字段稍微有点差别，比如 emptyface 的 word 字段和 eface 的 data 字段名称不同，但是数据型是一样的。

## 反射的三大定律

根据 Go 官方关于反射的博客，反射有三大定律：

> 1. Reflection goes from interface value to reflection object.
> 2. Reflection goes from reflection object to interface value.
> 3. To modify a reflection object, the value must be settable.

第一条是最基本的：反射是一种检测存储在 `interface` 中的类型和值机制。这可以通过 `TypeOf` 函数和 `ValueOf` 函数得到。
第二条实际上和第一条是相反的机制，它将 `ValueOf` 的返回值通过 `Interface()` 函数反向转变成 `interface` 变量。
前两条就是说 `接口型变量` 和 `反射类型对象` 可以相互转化，反射类型对象实际上就是指的前面说的 `reflect.Type` 和 `reflect.Value`。
第三条不太好懂：如果需要操作一个反射变量，那么它必须是可设置的。反射变量可设置的本质是它存储了原变量本身，这样对反射变量的操作，就会反映到原变量本身；反之，如果反射变量不能代表原变量，那么操作了反射变量，不会对原变量产生任何影响，这会给使用者带来疑惑。所以第二种情况在语言层面是不被允许的。
举一个经典例子：

    var x float64 = 3.4
    v := reflect.ValueOf(x)
    v.SetFloat(7.1) // Error: will panic.

执行上面的代码会产生 panic，原因是反射变量 `v` 不能代表 `x` 本身，为什么？因为调用 `reflect.ValueOf(x)` 这一行代码的时候，传入的参数在函数内部只是一个拷贝，是值传递，所以 `v` 代表的只是 `x` 的一个拷贝，因此对 `v` 进行操作是被禁止的。
可设置是反射变量 `Value` 的一个性质，但不是所有的 `Value` 都是可被设置的。
就像在一般的函数里那样，当我们想改变传入的变量时，使用指针就可以解决了。

    var x float64 = 3.4
    p := reflect.ValueOf(&x)
    fmt.Println("type of p:", p.Type())
    fmt.Println("settability of p:", p.CanSet())

输出是这样的：

    type of p: *float64
    settability of p: false

`p` 还不是代表 `x`，`p.Elem()` 才真正代表 `x`，这样就可以真正操作 `x` 了：

    v := p.Elem()
    v.SetFloat(7.1)
    fmt.Println(v.Interface()) // 7.1
    fmt.Println(x) // 7.1

关于第三条，记住一句话：如果想要操作原变量，反射变量 `Value` 必须要 hold 住原变量的地址才行。

# 反射相关函数的使用

## 代码样例

网络上各种博客文章里使用反射的样例代码非常多，读过这篇文章后，基本没有看不懂的，哈哈！不过，我这里还是举一个例子，并讲解一番：

```go
package main
import (
    "reflect"
    "fmt"
)
type Child struct {
    Name     string
    Grade    int
    Handsome bool
}
type Adult struct {
    ID         string `qson:"Name"`
    Occupation string
    Handsome   bool
}
// 如果输入参数 i 是 Slice，元素是结构体，有一个字段名为 `Handsome`，
// 并且有一个字段的 tag 或者字段名是 `Name` ，
// 如果该 `Name` 字段的值是 `qcrao`，
// 就把结构体中名为 `Handsome` 的字段值设置为 true。
func handsome(i interface{}) {
    // 获取 i 的反射变量 Value
    v := reflect.ValueOf(i)
    // 确定 v 是一个 Slice
    if v.Kind() != reflect.Slice {
        return
    }
    // 确定 v 是的元素为结构体
    if e := v.Type().Elem(); e.Kind() != reflect.Struct {
        return
    }
    // 确定结构体的字段名含有 "ID" 或者 json tag 标签为 `name`
    // 确定结构体的字段名 "Handsome"
    st := v.Type().Elem()
    // 寻找字段名为 Name 或者 tag 的值为 Name 的字段
    foundName := false
    for i := 0; i < st.NumField(); i++ {
        f := st.Field(i)
        tag := f.Tag.Get("qson")
        if (tag == "Name" || f.Name == "Name") && f.Type.Kind() == reflect.String {
            foundName = true
            break
        }
    }
    if !foundName {
        return
    }
    if niceField, foundHandsome := st.FieldByName("Handsome"); foundHandsome == false || niceField.Type.Kind() != reflect.Bool {
        return
    }
    // 设置名字为 "qcrao" 的对象的 "Handsome" 字段为 true
    for i := 0; i < v.Len(); i++ {
        e := v.Index(i)
        handsome := e.FieldByName("Handsome")
        // 寻找字段名为 Name 或者 tag 的值为 Name 的字段
        var name reflect.Value
        for j := 0; j < st.NumField(); j++ {
            f := st.Field(j)
            tag := f.Tag.Get("qson")
            if tag == "Name" || f.Name == "Name" {
                name = v.Index(i).Field(j)
            }
        }
        if name.String() == "qcrao" {
            handsome.SetBool(true)
        }
    }
}

func main() {
    children := []Child{
        {Name: "Ava", Grade: 3, Handsome: true},
        {Name: "qcrao", Grade: 6, Handsome: false},
    }
    adults := []Adult{
        {ID: "Steve", Occupation: "Clerk", Handsome: true},
        {ID: "qcrao", Occupation: "Go Programmer", Handsome: false},
    }
    fmt.Printf("adults before handsome: %v\n", adults)
    handsome(adults)
    fmt.Printf("adults after handsome: %v\n", adults)
    fmt.Println("-------------")
    fmt.Printf("children before handsome: %v\n", children)
    handsome(children)
    fmt.Printf("children after handsome: %v\n", children)
}
```

代码运行结果：

    adults before handsome: [{Steve Clerk true} {qcrao Go Programmer false}]
    adults after handsome: [{Steve Clerk true} {qcrao Go Programmer true}]
    -------------
    children before handsome: [{Ava 3 true} {qcrao 6 false}]
    children after handsome: [{Ava 3 true} {qcrao 6 true}]

代码主要做的事情是：找出传入的参数为 Slice，并且 Slice 的元素为结构体，如果其中有一个字段名是 `Name` 或者是 标签名称为 `Name`，并且还有一个字段名是 `Handsome` 的情形。如果找到，并且字段名称为 `Name` 的实际值是 `qcrao` 的话，就把另一个字段 `Handsome` 的值置为 true。
程序并不关心传入的结构体到底是什么，只要它的字段名包含 `Name` 和 `Handsome`，都是 handsome 函数要工作的对象。
注意一点，`Adult` 结构体的标签 `qson:"Name"`，中间是没有空格的，否则 `Tag.Get("qson")` 识别不出来。

## 未导出成员

利用反射机制，对于结构体中未导出成员，可以读取，但不能修改其值。
注意，正常情况下，代码是不能读取结构体未导出成员的，但通过反射可以越过这层限制。另外，通过反射，结构体中可以被修改的成员只有是导出成员，也就是字段名的首字母是大写的。

> 一个可取地址的 reflect.Value 变量会记录一个结构体成员是否是未导出成员，如果是的话则拒绝修改操作。
> CanAddr 不能说明一个变量是否可以被修改。
> CanSet 则可以检查对应的 reflect.Value 是否可取地址并可被修改。

```go
package main
import (
    "reflect"
    "fmt"
)
type Child struct {
    Name     string
    handsome bool
}
func main() {
    qcrao := Child{Name: "qcrao", handsome: true}
    v := reflect.ValueOf(&qcrao)
    f := v.Elem().FieldByName("Name")
    fmt.Println(f.String())
    f.SetString("stefno")
    fmt.Println(f.String())
    f = v.Elem().FieldByName("handsome")

    // 这一句会导致 panic，因为 handsome 字段未导出
    //f.SetBool(true)
    fmt.Println(f.Bool())
}
```

执行结果：

    qcrao
    stefno
    true

上面的例子中，handsome 字段未导出，可以读取，但不能调用相关 set 方法，否则会 panic。反射用起来一定要小心，调用类型不匹配的方法，会导致各种 panic。

# 反射的实际应用

反射的实际应用非常广：IDE 中的代码自动补全功能、对象序列化（json 函数库）、fmt 相关函数的实现、ORM（全称是：Object Relational Mapping，对象关系映射）……
这里举 2 个例子：json 序列化和 DeepEqual 函数。

## json 序列化

开发过 web 服务的同学，一定用过 `json` 数据格式。`json` 是一种独立于语言的数据格式。最早用于浏览器和服务器之间的实时无状态的数据交换，并由此发展起来。
Go 语言中，主要提供 2 个函数用于序列化和反序列化：

    func Marshal(v interface{}) ([]byte, error)
    func Unmarshal(data []byte, v interface{}) error

两个函数的参数都包含 `interface`，具体实现的时候，都会用到反射相关的特性。
对于序列化和反序列化函数，均需要知道参数的所有字段，包括字段类型和值，再调用相关的 get 函数或者 set 函数进行实际的操作。

## DeepEqual 的作用及原理

在测试函数中，经常会需要这样的函数：判断两个变量的实际内容完全一致。
例如：如何判断两个 slice 所有的元素完全相同；如何判断两个 map 的 key 和 value 完全相同等等。
上述问题，可以通过 `DeepEqual` 函数实现。

    func DeepEqual(x, y interface{}) bool

`DeepEqual` 函数的参数是两个 `interface`，实际上也就是可以输入任意类型，输出 true 或者 flase 表示输入的两个变量是否是“深度”相等。
先明白一点，如果是不同的类型，即使是底层类型相同，相应的值也相同，那么两者也不是“深度”相等。

    type MyInt int
    type YourInt int
    func main() {
    	m := MyInt(1)
    	y := YourInt(1)
    	fmt.Println(reflect.DeepEqual(m, y)) // false
    }

上面的代码中，m, y 底层都是 int，而且值都是 1，但是两者静态类型不同，前者是 `MyInt`，后者是 `YourInt`，因此两者不是“深度”相等。
在源码里，有对 DeepEqual 函数的非常清楚地注释，列举了不同类型，DeepEqual 的比较情形，这里做一个总结：

| 类型                                  | 深度相等情形                                                                                                              |
| ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Array                                 | 相同索引处的元素“深度”相等                                                                                                |
| Struct                                | 相应字段，包含导出和不导出，“深度”相等                                                                                    |
| Func                                  | 只有两者都是 nil 时                                                                                                       |
| Interface                             | 两者存储的具体值“深度”相等                                                                                                |
| Map                                   | 1、都为 nil；2、非空、长度相等，指向同一个 map 实体对象，或者相应的 key 指向的 value “深度”相等                           |
| Pointer                               | 1、使用 == 比较的结果相等；2、指向的实体“深度”相等                                                                        |
| Slice                                 | 1、都为 nil；2、非空、长度相等，首元素指向同一个底层数组的相同元素，即 \&x\[0] == \&y\[0] 或者 相同索引处的元素“深度”相等 |
| numbers, bools, strings, and channels | 使用 == 比较的结果为真                                                                                                    |

一般情况下，DeepEqual 的实现只需要递归地调用 == 就可以比较两个变量是否是真的“深度”相等。
但是，有一些异常情况：比如 func 类型是不可比较的类型，只有在两个 func 类型都是 nil 的情况下，才是“深度”相等；float 类型，由于精度的原因，也是不能使用 == 比较的；包含 func 类型或者 float 类型的 struct， interface， array 等。
对于指针而言，当两个值相等的指针就是“深度”相等，因为两者指向的内容是相等的，即使两者指向的是 func 类型或者 float 类型，这种情况下不关心指针所指向的内容。
同样，对于指向相同 slice， map 的两个变量也是“深度”相等的，不关心 slice， map 具体的内容。
对于“有环”的类型，比如循环链表，比较两者是否“深度”相等的过程中，需要对已比较的内容作一个标记，一旦发现两个指针之前比较过，立即停止比较，并判定二者是深度相等的。这样做的原因是，及时停止比较，避免陷入无限循环。
来看源码：

    func DeepEqual(x, y interface{}) bool {
    	if x == nil || y == nil {
    		return x == y
    	}
    	v1 := ValueOf(x)
    	v2 := ValueOf(y)
    	if v1.Type() != v2.Type() {
    		return false
    	}
    	return deepValueEqual(v1, v2, make(map[visit]bool), 0)
    }

首先查看两者是否有一个是 nil 的情况，这种情况下，只有两者都是 nil，函数才会返回 true。
接着，使用反射，获取 x，y 的反射对象，并且立即比较两者的类型，根据前面的内容，这里实际上是动态类型，如果类型不同，直接返回 false。
最后，最核心的内容在子函数 `deepValueEqual` 中。
代码比较长，思路却比较简单清晰：核心是一个 switch 语句，识别输入参数的不同类型，分别递归调用 deepValueEqual 函数，一直递归到最基本的数据类型，比较 int，string 等可以直接得出 true 或者 false，再一层层地返回，最终得到“深度”相等的比较结果。
实际上，各种类型的比较套路比较相似，这里就直接节选一个稍微复杂一点的 `map` 类型的比较：

    // deepValueEqual 函数
    // ……
    case Map:
    	if v1.IsNil() != v2.IsNil() {
    		return false
    	}
    	if v1.Len() != v2.Len() {
    		return false
    	}
    	if v1.Pointer() == v2.Pointer() {
    		return true
    	}
    	for _, k := range v1.MapKeys() {
    		val1 := v1.MapIndex(k)
    		val2 := v2.MapIndex(k)
    		if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1) {
    			return false
    		}
    	}
    	return true

    // ……

和前文总结的表格里，比较 map 是否相等的思路比较一致，也不需要多说什么。说明一点，`visited` 是一个 map，记录递归过程中，比较过的“对”：

    type visit struct {
    	a1  unsafe.Pointer
    	a2  unsafe.Pointer
    	typ Type
    }
    map[visit]bool

比较过程中，一旦发现比较的“对”，已经在 map 里出现过的话，直接判定“深度”比较结果的是 `true`。

# 总结

Go 作为一门静态语言，相比 Python 等动态语言，在编写过程中灵活性会受到一定的限制。但是通过接口加反射实现了类似于动态语言的能力：可以在程序运行时动态地捕获甚至改变类型的信息和值。
Go 语言的反射实现的基础是类型，或者说是 interface，当我们使用反射特性时，实际上用到的就是存储在 interface 变量中的和类型相关的信息，也就是常说的 `<type, value>` 对。
只有 interface 才有反射的说法。
反射在 reflect 包中实现，涉及到两个相关函数：

    func TypeOf ( i interface{} ) Type
    func ValueOf ( i interface{} ) Value

Type 是一个接口，定义了很多相关方法，用于获取类型信息。Value 则持有类型的具体值。Type、Value、Interface 三者间通过函数 TypeOf，ValueOf，Interface 进行相互转换。
最后温习一下反射三大定律：

> 1. Reflection goes from interface value to reflection object.
> 2. Reflection goes from reflection object to interface value.
> 3. To modify a reflection object, the value must be settable.

翻译一下：

> 1. 反射将接口变量转换成反射对象 Type 和 Value；
> 2. 反射可以通过反射对象 Value 还原成原先的接口变量；
> 3. 反射可以用来修改一个变量的值，前提是这个值可以被修改。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/tgz38n/1616980149660-85acb022-e70f-4ef8-91ab-10d7186fd785.png)

# 参考资料

【维基百科中文】[https://zh.wikipedia.org/wiki/反射\_(计算机科学)](<https://zh.wikipedia.org/wiki/%E5%8F%8D%E5%B0%84_(%E8%AE%A1%E7%AE%97%E6%9C%BA%E7%A7%91%E5%AD%A6)>)
【码洞老钱 反射】<https://juejin.im/post/5c2040d76fb9a049c643d9bd>
【Go 官方博客 reflection】<https://blog.golang.org/laws-of-reflection>
【GCTT 译文，不错】<https://mp.weixin.qq.com/s/dkgJ_fA0smvpv69t5Nv-7A>
【json 库 源码分析】<https://zhuanlan.zhihu.com/p/37165706>
【reflect 代码例子和图比较好】<https://blog.gopheracademy.com/advent-2018/interfaces-and-reflect/>
【反射使用讲得不错】<https://juejin.im/post/5a75a4fb5188257a82110544>
【接口和反射的关系 ，english】<https://blog.gopheracademy.com/advent-2018/interfaces-and-reflect/>
【总结成知识点】<http://www.cnblogs.com/susufufu/p/7653579.html>
【Type Value】<https://colobu.com/2016/07/09/dive-into-go-13/>
【讲得比较清晰简单】[https://www.lijiaocn.com/编程/2017/11/06/golang-reflection.html](https://www.lijiaocn.com/%E7%BC%96%E7%A8%8B/2017/11/06/golang-reflection.html)
【DeepEqual】<https://github.com/Chasiny/Blog/blob/master/blog/go/package/go-reflect-deepequal.md>
【反射使用场景】<https://yq.aliyun.com/articles/599584>
