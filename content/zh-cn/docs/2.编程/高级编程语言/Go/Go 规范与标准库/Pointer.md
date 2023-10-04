---
title: Pointer
linkTitle: Pointer
date: 2023-11-20T21:30
weight: 7
---

# 概述

> 参考：
>
> - [Go 官方文档，参考-指针类型](https://go.dev/ref/spec#Pointer_types)

在 Go 语言中，**Pointer(指针)** 可以有两种含义：

- 通过 `&` 符号获取一个变量的内存地址，即**指针**。通常使用十六进制数表示。这种方式称为**指针引用**
- 指针也可以表示一种**数据类型**。可以声明一个指针类型的变量，用以存储内存地址。

一个**指针变量**可以指向任何一个**值的内存地址**。这个内存地址，在 32 位机器上占用 4 个字节，在 64 位机器上占用 8 个字节，并且与其所指向的值得的大小无关。
在 Go 语言中，不能进行指针运算。

## 指针的引用

每一个变量都有指针，我们可以通过 `&` 符号引用该变量的指针，也就是获取该变量的内存地址。

> 我们平时说引用指针，并不是引用指针类型的变量，指针类型的变量也是一种变量，正常使用变量名称即可引用。

这里说要说的指针的引用是指**引用一个变量的内存地址，即变量的指针**
**格式：**在变量名称前加上 `&` 符号，即可获取该变量的内存地址，即该变量的指针。

```go
&VarID
```

**注意：若该变量的值为空，则该变量依然具有内存地址：**

```go
var a string
fmt.Println(&a)
```

这将会输出：0xc000010250

## 指针变量的声明

**格式：**`*` 与 `数据类型` 的组合书写，即代表指针类型：

```go
var VarID *TYPE
```

这里需要注意的是，当一个指针变量被声明后，它的值为`nil`，但是这个变量本身是具有指针的

```go
func main() {
    // 声明一个 `字符串指针` 类型的变量
 var VarID *string
 fmt.Println(VarID)
    fmt.Println(&VarID)
}
```

输出结果：

```bash
<nil>
0xc000014088
```

一个指针类型的变量可以保存内存地址，同时自己也具有内存地址。

## 指针变量的赋值

**格式：**

```go
var ptr *string
stringVar = "pointer string"
ptr = &stringVar
```

虽然一个指针类型的变量的值是类似 0xc000010250 这样的内存地址，但是若我们声明了一个**字符串指针类型**的变量，那么给该变量赋值时，也要使用**保存字符串的内存地址**；如果使用保存其他类型(比如 int)的内存地址，将会报错：

```go
normalVar := 5
var ptr *string
ptr = &normalVar // 这里将会报错：cannot use &normalVar (value of type *int) as type *string in assignment
```

因为 `*string` 和 `*int` 是两个不同的类型。

## 指针变量的解除引用

在上面指针变量的赋值中，我们在最后使用指针变量的 Dereferences(解除引用)，以获取具体的值。
**格式：**

```go
*PointerTypeVarID
```

在某个指针类型的变量前面添加 `*` 则会解除指针引用，从而**获取指针变量中的值，所指向的内存空间中的值。**
注意：若指针变量未实例化(即指针变量为 nil)，那么解除引用时将会报错：panic: runtime error: invalid memory address or nil pointer dereference。主要错误在于 nil pointer dereference(空指针接触引用)

## 指针变量的实例化

通过 var 关键字声明的指针变量的默认值为 nil，无法通过解引用的方式为该内存空间填入值，或者获取内存空间的值，因为内存空间不存在，是 nil。
此时可以通过 `new()` 函数，在声明的指针变量的同时实例化该指针，并为指针变量的值赋予一个值。

```go
ptrNew := new(string)
*ptrNew = "a" // 由于此时 ptrNew 具有内存地址，所以可以直接通过解引用的方式为 ptrNew 内存地址空间内赋予值
fmt.Println(*ptrNew) // 将会输出：a
```

### var 与 new() 的区别

主要区别在于系统是否会为指针变量初始化一个值，var 不会， new() 会
若一个指针变量没有被初始化一个值，那么解除引用将会失败，因为指针变量没有值，也就是说这个指针变量内没有内存地址，此时接触引用将会失败，也就是 nil pointer dereference

```go
var ptr *string
*ptr = "5"
// 只声明未实例化，上面两行报错：panic: runtime error: invalid memory address or nil pointer dereference

ptr := new(string)
*ptr = "5"
// 这两行不报错
```

## TODO

这是什么语法？`(*stringValue)(p)`？

```go
func newStringValue(val string, p *string) *stringValue {
 *p = val
 return (*stringValue)(p)
}
```

## 简单示例

```go
func main() {
 normalVar := 5
 // 通过 & 符号取得变量 a 的内存地址，即指向 a 的指针
 fmt.Println("变量 a 的内存地址，即指针为：", &normalVar)

 // 声明一个指针类型的变量
 var ptr *string // 这是一个字符串指针类型的变量
 fmt.Println("刚刚声明的指针没有任何内存地址，默认值为 nil：", ptr)

 // 指针赋值
 // 错误示例：不可以使用 *int 类型给 *string 类型赋值，虽然都是内存地址，但是不可以这么做
 // ptr = &normalVar
 // 正确示例:
 newPtr := strconv.Itoa(normalVar)
 ptr = &newPtr
 fmt.Println("为指针类型变量赋予一个内存地址后，获取该内存地址内保存的值：", *ptr)
}
```

输出结果：

    变量 a 的内存地址，即指针为： 0xc0000ba000
    刚刚声明的指针没有任何内存地址，默认值为 nil： <nil>
    为指针类型变量赋予一个内存地址后，获取该内存地址内保存的值： 5

## 通过指针改变变量的值

由于指针的特殊性，我们可以在任何地方修改一个局部变量的值。这就是人们常说的“值传递”和“引用传递”。

```go
package main

import "fmt"

func zeroval(ival int) {
    ival = 0
}

func zeroptr(iptr *int) {
    *iptr = 0
}

func main() {
    i := 1
    fmt.Println("initial:", i)

    zeroval(i)
    fmt.Println("zeroval:", i)

    zeroptr(&i)
    fmt.Println("zeroptr:", i)

    fmt.Println("pointer:", &i)
}
```

运行结果

```bash
$ go run pointers.go
initial: 1
zeroval: 1
zeroptr: 0
pointer: 0x42131100
```

在 `zeroptr()` 函数中，我们将变量 i 的指针传递进去，此时在函数内修改，将会影响变量 i 的值

# 结构体与指针

TODO：

# 细嗦 Golang 的指针

> 原文链接：[稀土掘金，细嗦 Golang 的指针](https://juejin.cn/post/7114673293084819492)

与 C 语言一样，Go 语言中同样有指针，通过指针，我们可以只传递变量的内存地址，而不是传递整个变量，这在一定程度上可以**节省内存的占用**，但凡事有利有弊，Go 指针在使用也有一些注意点，稍不留神就会踩坑，下面就让我们一起来细嗦下。

## 1.指针类型的变量

在 Golang 中，我们可以通过**取地址符号&**得到变量的地址，而这个新的变量就是一个指针类型的变量，指针变量与普通变量的区别在于，它存的是内存地址，而不是实际的值。
![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zmk3gz/1661497969122-a72c48b8-e325-49e7-a651-164ec9e0ff92.webp)
图一

如果是普通类型的指针变量（比如 `**int**`），是无法直接对其赋值的，必须通过 `* 取值符号`才行。

```go
func main() {
 num := 1
 numP := &num


 *numP = 2
}
```

但结构体却比较特殊，在日常开发中，我们经常看到一个结构体指针的内部变量仍然可以被赋值，比如下面这个例子，这是为什么呢？

```go
type Test struct {
 Num int
}


func main() {
 test := Test{Num: 1}
 test.Num = 3
 fmt.Println("v1", test)

 testP := &test
 testP.Num = 4
 fmt.Println("v2", test)
}
```

这是因为结构体本身是一个连续的内存，通过 `testP.Num` ，本质上拿到的是一个普通变量，并不是一个指针变量，所以可以直接赋值。

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zmk3gz/1661497969139-3b504293-af29-4ebc-8dd7-4033fd651b53.webp)

图二

那 slice、map、channel 这些又该怎么理解呢？为什么不用取地址符号也能打印它们的地址？比如下面的例子

```bash
func main() {
 nums := []int{1, 2, 3}
 fmt.Printf("%p\n", nums)     // 0xc0000160c0
 fmt.Printf("%p\n", &nums[0]) // 0xc0000160c0

 maps := map[string]string{"aa": "bb"}
 fmt.Printf("%p\n", maps) // 0xc000076180

 ch := make(chan int, 0)
 fmt.Printf("%p\n", ch) // 0xc00006c060
}
```

这是因为，**它们本身就是指针类型**！只不过 Go 内部为了书写的方便，并没有要求我们在前面加上 **\* 符号**。

在 Golang 的运行时内部，创建 slice 的时候其实返回的就是一个指针：

```go


func makeslice(et *_type, len, cap int) unsafe.Pointer {
 mem, overflow := math.MulUintptr(et.size, uintptr(cap))
 if overflow || mem > maxAlloc || len < 0 || len > cap {





  mem, overflow := math.MulUintptr(et.size, uintptr(len))
  if overflow || mem > maxAlloc || len < 0 {
   panicmakeslicelen()
  }
  panicmakeslicecap()
 }

 return mallocgc(mem, et, true)
}
```

而且返回的指针地址其实就是**slice 第一个元素的地址**（上面的例子也体现了），当然如果 slice 是一个 nil，则返回的是 `0x0` 的地址。slice 在参数传递的时候其实拷贝的指针的地址，底层数据是共用的，所以对其修改也会影响到函数外的 slice，在下面也会讲到。

map 和 slice 其实也是类似的，在在 Golang 的运行时内部，创建 map 的时候其实返回的就是一个 hchan 指针：

```go


func makechan(t *chantype, size int) *hchan {
 elem := t.elem


 if elem.size >= 1<<16 {
  throw("makechan: invalid channel element type")
 }
 ...
 return c
}
```

最后，为什么 `fmt.Printf` 函数能够直接打印 slice、map 的地址，除了上面的原因，还有一个原因是其内部也做了特殊处理：

```go

func Printf(format string, a ...interface{}) (n int, err error) {
 return Fprintf(os.Stdout, format, a...)
}


func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
 p := newPrinter()
 p.doPrintf(format, a)
 n, err = w.Write(p.buf)
 p.free()
 return
}


func (p *pp) doPrintf(format string, a []interface{}) {
  ...
 default:


   if 'a' <= c && c <= 'z' && argNum < len(a) {
    ...
    p.printArg(a[argNum], rune(c))
    argNum++
    i++
    continue formatLoop
   }

   break simpleFormat
  }

}


func (p *pp) printArg(arg interface{}, verb rune) {
 p.arg = arg
 p.value = reflect.Value{}
  ...
 case 'p':
  p.fmtPointer(reflect.ValueOf(arg), 'p')
  return
 }
 ...
}


func (p *pp) fmtPointer(value reflect.Value, verb rune) {
 var u uintptr
 switch value.Kind() {

 case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
  u = value.Pointer()
 default:
  p.badVerb(verb)
  return
 }
  ...
}
```

## 2.Go 只有值传递，没有引用传递

值传递和引用传递相信大家都比较了解，在函数的调用过程中，如果是值传递，则在传递过程中，其实就是将参数的值复制一份传递到函数中，如果在函数内对其修改，**并不会影响函数外面的参数值**，而引用传递则相反。

```go
type User struct {
 Name string
 Age  int
}


func setNameV1(user *User) {
 user.Name = "test_v1"
}


func setNameV2(user User) {
 user.Name = "test_v2"
}

func main() {
 u := User{Name: "init"}
 fmt.Println("init", u)

 up := &u
 setNameV1(up)
 fmt.Println("v1", u)

 setNameV2(u)
 fmt.Println("v2", u)
}
```

但在 Golang 中，这所谓的“引用传递”其实**本质上是值传递**，因为这时候也发生了拷贝，只不过这时拷贝的是指针，而不是变量的值，所以**“Golang 的引用传递其实是引用的拷贝”。**

![](https://notes-learning.oss-cn-beijing.aliyuncs.com/zmk3gz/1661497969124-2aa0ee1c-5b91-4d5c-8d77-fefa46caa8c9.webp)

图三

可以通过以下代码验证：

```go
type User struct {
 Name string
 Age  int
}



func setNameV1(user *User) {
 fmt.Printf("v1: %p\n", user)
 fmt.Printf("v1_p: %p\n", &user)
 user.Name = "test_v1"
}


func setNameV2(user User) {
 fmt.Printf("v2_p: %p\n", &user)
 user.Name = "test_v2"
}

func main() {
 u := User{Name: "init"}

 up := &u
 fmt.Printf("init: %p \n", up)
 setNameV1(up)
 setNameV2(u)
}
```

注：slice、map 等本质也是如此。

## 3.`for range`与指针

`for range`是在 Golang 中用于遍历元素，当它与指针结合时，稍不留神就会踩坑，这里有一段经典代码：

```go
type User struct {
 Name string
 Age  int
}

func main() {
 userList := []User {
  User{Name: "aa", Age: 1},
  User{Name: "bb", Age: 1},
 }

 var newUser []*User
 for _, u := range userList {
  newUser = append(newUser, &u)
 }

 for _, nu := range newUser {
  fmt.Printf("%+v", nu.Name)
 }
}
```

按照正常的理解，应该第一次输出`aa`，第二次输出`bb`，但实际上两次都输出了`bb`，这是因为 `for range` 的时候，**变量 u 实际上只初始化了一次**（每次遍历的时候 u 都会被重新赋值，但是地址不变），导致每次 append 的时候，**添加的都是同一个内存地址**，所以最终指向的都是最后一个值 bb。

我们可以通过打印指针地址来验证：

```go
func main() {
 userList := []User {
  User{Name: "aa", Age: 1},
  User{Name: "bb", Age: 1},
 }

 var newUser []*User
 for _, u := range userList {
  fmt.Printf("point: %p\n", &u)
  fmt.Printf("val: %s\n", u.Name)
  newUser = append(newUser, &u)
 }
}


point: 0xc00000c030
val: aa
point: 0xc00000c030
val: bb
```

### 两种解决方式

使用 `for` 代替 `for...range`

```go
func CorrectUsageOfForAndPointer() {
 userList := []User{
  {Name: "aa", Age: 1},
  {Name: "bb", Age: 1},
 }

 var newUser []*User
 for i := 0; i < len(userList); i++ {
  newUser = append(newUser, &userList[i])
 }

 for _, nu := range newUser {
  fmt.Printf("%+v\n", nu.Name)
 }
}
```

创建一个 Struct，将数组包含在其中

```go
type Users struct {
 Users []*User
}

type User struct {
 Name string
 Age  int
}

// 正确用法一
func CorrectUsageOfForrangeAndPointer() {
 usersList := Users{
  Users: []*User{
   {Name: "aa", Age: 1},
   {Name: "bb", Age: 1},
  },
 }

 var newUser []*User
 for _, u := range usersList.Users {
  newUser = append(newUser, u)
 }

 for _, nu := range newUser {
  fmt.Printf("%+v\n", nu.Name)
 }
}
```

类似的错误在`Goroutine`也经常发生：

```go

func main() {
 for i := 0; i < 10; i++ {
  go func(idx *int) {
   fmt.Println("go: ", *idx)
  }(&i)
 }
 time.Sleep(5 * time.Second)
}
```

## 4.闭包与指针

什么是闭包，一个函数和对其周围状态（**lexical environment，词法环境**）的引用捆绑在一起（或者说函数被引用包围），这样的组合就是**闭包**（**closure**）。也就是说，闭包让你可以在一个**内层函数中访问到其外层函数的作用域**。

当闭包与指针进行结合时，如果闭包里面是一个指针变量，则外部变量的改变，也会影响到该闭包，起到意想不到的效果，让我们继续在举几个例子进行说明：

```go
func incr1(x *int) func() {
 return func() {
  *x = *x + 1
  fmt.Printf("incr point x = %d\n", *x)
 }
}
func incr2(x int) func() {
 return func() {
  x = x + 1
  fmt.Printf("incr normal x = %d\n", x)
 }
}

func main() {
 x := 1
 i1 := incr1(&x)
 i2 := incr2(x)
 i1()
 i2()
 i1()
 i2()

 x = 100
 i1()
 i2()
 i1()
 i2()
}
```

## 5.指针与内存逃逸

内存逃逸的场景有很多，这里只讨论由指针引发的内存逃逸。理想情况下，肯定是尽量减少内存逃逸，因为这意味着 GC（垃圾回收）的压力会减小，程序也会运行得更快。不过，使用指针又能减少内存的占用，所以这本质是内存和 GC 的权衡，需要合理使用。

下面是指针引发的内存逃逸的三种场景（欢迎大家补充~）

第一种场景：函数返回局部变量的指针

```go
type Escape struct {
 Num1  int
 Str1  *string
 Slice []int
}


func NewEscape() *Escape {
 return &Escape{}
}

func main() {
 e := &Escape{Num1: 0}
}
```

第二种场景：被已经逃逸的变量引用的指针

```go
func main() {
 e := NewEscape()
 e.SetNum1(10)

 name := "aa"

 e.Str1 = &name
}
```

第三种场景：被指针类型的 slice、map 和 chan 引用的指针

```go
func main() {
 e := NewEscape()
 e.SetNum1(10)

 name := "aa"
 e.Str1 = &name


 arr := make([]*int, 2)
 n := 10
 arr[0] = &n
}
```

欢迎大家继续补充指针的其他注意事项~

## 参考

[又吵起来了，Go 是传值还是传引用？](https://link.juejin.cn/?target=https%3A%2F%2Fdeveloper.51cto.com%2Farticle%2F664004.html)

[GO 语言变量逃逸分析](https://link.juejin.cn/?target=https%3A%2F%2Fwww.cnblogs.com%2Fhualou%2Fp%2F12069815.html)
