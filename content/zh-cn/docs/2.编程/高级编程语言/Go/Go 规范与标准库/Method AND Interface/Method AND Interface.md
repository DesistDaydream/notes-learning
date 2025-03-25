---
title: Method AND Interface
linkTitle: Method AND Interface
weight: 1
---

# 概述

> 参考：
>
> - [公众号-新亮笔记，回答连个被频繁问道的代码写法问题](https://mp.weixin.qq.com/s/gpMzEoRofGE9LmayeYw1qw)
>   - 1.强制检查类型是否实现接口
>   - 2.强制接口中所有方法只能在本包中实现

Method 与 Interface 是 Go 语言是想面向对象编程的一种解决方式，但是更轻量。

**Go 是面向对象的编程语言吗？**

官方 FAQ 给出了标准答案: Yes and No。

当然，Go 有面向对象编程的类型和方法的概念，但是它没有继承(hierarchy)一说。Go 语言的接口实现和其它的编程语言不一样，Go 开发者的初衷就是保证它易于使用，用途更广泛。

还有一种“模拟”产生子类的方法，拿就是通过在类型中嵌入其它的类型，但是这是一种“组合”的方式，而不是继承。

没有了继承， Go 语言的对象变得比 C++和 Java 中更轻量级。

在 Go 语言中，接口定义了一套方法的集合，任何实现这些方法的对象都可以被认为实现了这个接口，这也称作 Duck Type。这不像其它语言比如 java 需要预先声明类型实现了某个或者某些接口，这使得 Go 接口和类型变得很轻量级，它解耦了接口和具体实现的硬绑定。显然这是 Go 的开发者深思熟虑的一个决定。

> if something looks like a duck, swims like a duck and quacks like a duck then it’s probably a duck.

因为没有继承，你也只能通过 Go 接口实现面向对象编程的多态。本身 Go 接口在内部实现上也是一个(其实是两种,其中一种专门处理 `interface{}`)结构体，它的虚函数指向具体的类型的实现。在编译代码的时候，Go 编译器还会做优化，不需要接口的时候，它会使用具体的方法来代替接口使用，这样进一步优化性能，这叫做 devirtualize 调用。

# Methods(方法)

**Method(方法)** 是一种特殊类型的 **Function(函数)**。是作用在 **Receiver(接收者)** 上的一个函数，`接收者`是某种类型的**变量**。接收者不能是一个接口类型；也不能是一个指针类型，但是可以是任何其他允许类型的指针。

> 在日常描述中，称为作用在 XX 类型上的方法

**方法的声明**：`func (RecvID RecvType) MethodID(ParameterList) (ReturnValueList) {...}`。

- RecvID # Receiver(接受者) 的标识符，即：Recv 类型的**变量**。如果 Method 不需要 Recv 的值，可以用 `_` 代替 RecvID。定义方法就是用类型来定义其方法

**方法的调用**：`RecvID.MethodID()`。i.e. 调用某个类型上的方法。因为方法是作用在类型上的，所以不用初始化，直接调用即可。

## 方法与结构体的关系

比如：有一个矩形的结构体，里面有两个属性，长和宽；如果想让整个矩形得出面积，需要一个方法，这个求面积的方法可以给矩形发消息，告诉矩形这是你求面积的方法；这时候这个矩形就可以拿着这个方法，自己算出自己的面积

用白话说：其实一个结构体就是一个对象，这个对象有很多很多的属性，想要根据这个对象的属性来得出某些结果，就可以将这个对象作用在某个方法上，这个方法就可以根据这个对象的某些属性进行计算来得出结果。再举一个例子，有一个人有多种属性(身高，体重，性别，腰围)；根据这些属性，可以创建一个计算体型的方法，这个方法根据这些属性中的 1 个或者多个计算出结果是偏瘦、偏胖还是适中。

## Method 总结

Go 的 Methods 就是面向对象编程的基础，Go 语言中并没有一个所谓的 “类”(也就是建模中的模板)，Go 所有被方法作用的接收者都可以称之为“对象”，所以，想要声明一个“模板”，只需要声明一个类型即可，这个类型可以是 int、string、struct(虽然~真正处理需求的时候，大部分情况下都是用 struct 作为“模板”~~o(∩_∩)o)，只要为这个类型声明了其所拥有的方法，这个类型就可以称之为“模板”。

# Interfaces(接口)

> 参考：
>
> - [知乎《如何理解 Golang 中的接口》](https://www.zhihu.com/question/318138275)
> - [公众号-码农桃花源，深度解密 Go 语言之关于 interface 的 10 个问题](https://mp.weixin.qq.com/s/EbxkBokYBajkCR-MazL0ZA)
>
> 本文的[总结](#总结)部分，将会根据实际案例，列举，并总结一些**编程思想**

首先，**Interface(接口) 是一种类型。更准确地说，Interface(接口)** 是仅包含方法名、参数、返回值的未具体实现的**一组方法的集合**。当一个类型定义了接口中所有的方法(注意：是所有方法)，就称这个类型 **implements(实现了)** 该接口，这个类型就称为这个接口的 **implementator(实现者)**。

接口定义了一个类型应该具有的方法，并由该类型决定如何实现这些方法。接口里不能包含变量。

接口也可以算自定义类型的一种，使用关键字 `type` 与 `interface` 来定义，所以可以对接口赋值，并且接口可以动态改变其自身的类型，只要某个类型实现了该接口，该接口类型在使用这个类型的时候，就会变成这个类型，这称为接口的**多态性**。再深入一点的描述详见代码中的注释部分

**接口的定义**：

```go
type InterfaceID interface {
    Method1(ParameterList) ReturnType
    Method2(ParameterList) ReturnType
    ...
}
```

**接口变量的声明**：`var InterfaceVarID InterfaceID`。接口可以有值，`InterfaceVarID`是一个 multiword(多字)。数字结构，值为 nil。本质上是一个指针，但又不完全是一回事。

**接口的赋值**：`InterfaceVarID=TypeVarID`。使用类型变量对接口变量进行赋值(i.e.接口变量包含一个指向类型变量的引用)。

其实，对接口所谓的赋值，并不是真正的值(不像`a=2,b="sting"`这类)，这里面所赋的值，其实是一个类型。与其说 `赋值`，不如叫`赋类型`。当接口变量具有一个类型时，接口变量的值就是这个类型所声明的变量的值

**接口的引用**：`InterfaceVarID.MethodID()`。

> 这里的引用实际上引用接口内的方法，但是，光引用一个接口是没有什么太大意义的。不给接口变量赋值的话，引用时将会报错如下内容 `panic: runtime error: invalid memory address or nil pointer dereference`

只有为接口变量赋值，通过赋值，将接口转换为对应的类型，这时，再引用接口上的方法，其实就是引用对应类型上的方法。这种方式使此方法更具有一般性。接口变量里包含了接收者实例的值和指向对应方法表的指针。

接口的特性：

- 指向接口值的 **指针** 是 **非法** 的，不仅一点用没有，还会导致代码错误。说白了就是使用接口时不要带 `*` 符号
- **接口=某个类型、某个类型=接口**。当某个类型实现了一个接口，就可以像加粗字那么描述。
  - 因为下面这个代码可以看到，`DoDuck()` 的形参是接口类型，但是 `main()` 中调用 `DoDuck()` 的时候，是可以传递一个 `Chicken` 这个结构体类型。所以，当一个类型实现一个接口时，这个类型=接口。并且，`d.Quack()` 实际上是执行的 `c.Qucak()`

```go
type Duck interface {
    Quack()
}

func DoDuck(d Duck) {
    d.Quack()
}

type Chicken struct {
}

func (c Chicken) Quack() {
    fmt.Println("嘎嘎")
}

func main() {
    c := Chicken{}
    DoDuck(c)
}
```

- 当某个类型(比如 struct、slice 等)实现了接口方法集中的方法，每一个方法的实现说明了此方法是如何作用于该类型的。i.e.**Implement Interface(实现接口)。**同时方法集也构成了该类型的接口。实现了 InterfaceID 接口类型的变量可以赋值给 VarID(接收者值)，此时方法表中的指针会指向被实现的接口方法。如果另一个类型(也实现了该接口)的变量被赋值给 VarID，这两者(指针和方法实现)也会随之改变。
- 多个类型可以实现同一个接口，所以类型不用显式得声明其实现了哪一个接口，所以接口会被隐式得实现。i.e.接口里只需要看到方法是什么，不用关心方法作用在哪个类型上。
- 一个类型可以实现多个接口
- 实现某个接口的类型可以有其他办法
- 有的时候，也会以一种稍微不同的方式来使用接口这个词：从类型的角度来看，它的接口指的是：它的所有导出方法，只不过没有显式地为这些导出方法额外定一个接口而已。

**接口嵌套接口**

一个接口可以包含一个或多个其他的接口，这相当于直接将这些内嵌接口的方法列举在外层接口中一样。

比如下面的例子接口`File`包含了`ReadWrite`和`Lock`的所有方法，它还额外有一个`Close()`方法。

```go
type ReadWrite interface {
    Read(b Buffer) bool
    Write(b Buffer) bool
}

type Lock interface {
    Lock()
    Unlock()
}

type File interface {
    ReadWrite
    Lock
    Close()
}
```

## **空接口**

空接口或者最小接口不包含任何方法，它对实现不做任何要求

定义格式：`type InterfaceID interface {}`

可以给一个空接口类型的变量赋任何类型的值

如果一个 interface 中如果没有定义任何方法，即为空 interface，表示为 interface{}。如此一来，任何类型就都能满足它，这也是为什么当函数参数类型为 interface{} 时，可以给它传任意类型的参数。

示例代码，如下：

```go
package main

import "fmt"

func main() {
    var i interface{} = 1
    fmt.Println(i)
}
```

更常用的场景，Go 的 interface{} 常常会被作为函数的参数传递，用以帮助我们实现其他语言中的泛型效果。Go 中暂时不支持 泛型，不过 Go 2 的方案中似乎将支持泛型。

## 接口的 polymorphism(多态性)

代码示例 1：interface_formula.go 其中接口变量的类型可以动态得随着不同的值而变化为对应的类型。

代码示例 2：interface_salary.go 其中接口变量切片中的元素可以是任意类型；不像普通的切片，所有元素都是同一个类型。

## Type Assertion(类型断言)

**检测和转换接口变量的类型**

一个接口类型的变量 `InterfaceVar` 中可以包含任何类型的值，必须有一种方式来检测它的动态类型。i.e.运行时在接口变量中存储的值的实际类型。在执行过程中动态类型可能会有所不同，但是它总是可以分配给接口变量本身的类型。通常我们可以使用 **TypeAssertion(类型断言)** 来测试在当前执行该语句的时候`InterfaceVar`所定义的接口是否是`Type`这个类型

使用格式：`v := InterfaceVar.(Type)`

类型断言可能是无效的，虽然编译器会尽力检查转换是否有效，但是它不可能预见所有的可能性。如果转换在程序运行时失败会导致错误发生。更安全的方式是使用以下形式来进行类型断言：

```go
if v, ok := InterfaceVar.(Type); ok {  // checked type assertion
    Process(v)
    return
}
```

如果转换合法，则 `v` 是 `InterfaceVar` 转换到类型 `Type` 的值，`ok` 会是 `true`；否则 `v` 是类型 `Type` 的零值，`ok` 是 `false`，也没有运行时错误发生。

> 说白了，就是给 v 赋值

## 强制检查某类型是否实现了某接口

```go
var _ Signature = (*signature)(nil)
// 或
var _ Signature = signature{}

// var 关键字可省略
```

上述两种写法可以让编译器检查 `signature` 类型是否实现了 `Signature` 接口。
比如：

```go
package main

import "fmt"

var _ Study = (*study)(nil)

type study struct {
    Name string
}

type Study interface {
    Listen(message string) string
}

func main() {
    fmt.Println("hello world")
}
```

运行后会抛出异常：

```bash
./main.go:5:5: cannot use (*study)(nil) (type *study) as type Study in assignment:
 *study does not implement Study (missing Listen method)
```

只有去掉 `var _ Study = (*study)(nil)`才可以正常运行。
在某个类型需要实现接口时，推荐使用这种语法进行判断，以便出现问题时可以快速定位。
在 mysql exporter 中的各种采集器中，大量使用了这种特性

- <https://github.com/prometheus/mysqld_exporter/blob/main/collector/binlog.go#L142>

## 强制接口中的方法只能在本包中实现

在接口中定义一个小写字母开头的方法即可
比如：

```go
package study

type Study interface {
    Listen(message string) string
    i()
}

func Speak(s Study) string {
    return s.Listen("abc")
}
```

```go
package main

type stu struct {
    Name string
}

func (s *stu) Listen(message string) string {
    return s.Name + " 听 " + message
}

func (s *stu) i() {}

func main() {
    message := study.Speak(new(stu))
    fmt.Println(message)
}
```

此时运行后将会抛出异常：

```bash
./main.go:19:28: cannot use new(stu) (type *stu) as type study.Study in argument to study.Speak:
 *stu does not implement study.Study (missing study.i method)
  have i()
  want study.i()
```

只要去掉接口中 `i()` 方法即可，或者都改成大写。

# 总结

结构体、方法、接口，这三者在项目开发中，占有很大的比重，也是一个 go 项目的设计思路。这三者也是[OOP](/docs/2.编程/计算机科学/Object-oriented%20Programming/OOP.md)的解决方式。

**如果用现实来比喻，那 Struct 是对实体的抽象，将各种实体或个体抽象为一个 Object(对象)，Method 就是这些对象的行为，那么 Interface 就是将很多对象进行了分类，把具有相同行为的对象统一划做同一类别**。

> 行为内部的具体步骤可能不一样，但是都要做相同的行为。比如假如两个对象都要跑步，一个是边跳边跑，一个是边喊边跑。

Interface(接口) 在 Go 语言有着至关重要的地位。接口是 Go 语言这个类型系统的基石，让 Go 语言在基础编程哲学的探索上达到了前所未有的高度。

接口解除了类型依赖，有助于减少用户的可视方法，屏蔽了内部结构和实现细节。但是接口实现机制会有运行期开销，也不能滥用接口。相对于包，或者不会频繁变化的内部模块之间，不需要抽象出接口来强行分离。

接口最常用的**使用场景**，是**对包提供访问**，或**预留扩展空间**。也是体现多态很重要的手段。说到底，接口的意义：就是解耦合，降低程序和程序之间的关联程度，降低耦合性。

当使用一个第三方库的时候，这个第三方库实现了很多接口，当我们需要使用这些接口时，就需要实现这些接口的方法。因为，这个第三方库内的某些函数(或方法)会将这些接口作为需要传入的参数，以实现自身的逻辑。

struct 和 interface 都是一种类型，所以都可以声明，并赋值。如果该函数的形参是一个接口类型，在使用该函数，传递的实参就可以是任意类型。这时，如果我们自定义的这个 struct 实现了第三方库的接口中的方法，那么这个函数就可以调用这些自己定义的方法，来根据自身的函数逻辑，处理这个自己定义的 struct。

所以，所谓的接口，就是接收某种类型的值，然后通过内部的某些处理逻辑，来根据这个值进行处理。

其实**接口跟函数的概念是一样的**，只不过是一层更抽象的东西，而且可以接收比函数更多的类型的值；并且，可以通过更多的函数处理这个接收到的值。如果说函数是一种行为的行为的话，那接口就是一类行为的统称。

> struct 实现了某个接口，并通过某些方法调用之后，接口后面的函数用不用 struct 所实现的 方法，就是库中接口后面的具体函数(i.e.函数行为) 来决定的。

也可以这么说，**如果把接口比作变量(实际上，也确实可以声明一个接口类型的变量)**，那么这个**类型实例化后的变量就可以作为值，直接赋予给这个接口变量**。

**函数 与 类型(i.e.值)，可以说是构成代码的最基础的东西**

## 用 USB 来类比接口非常形象

参考：[文中代码在 GitHub 的 GoLearning 项目中](https://github.com/DesistDaydream/GoLearning/tree/master/practice/usb_interface)

### 单独描述

USB 接口规定了他可以处理 `Start()` 和 `End()` 两个方法，这与这俩方法中的具体实现他不管。并且还有一个 `OperatorRead()` 函数，将`接口`作为形参，并调用接口里的方法实现一些功能。

```go
type USB interface {
    Start()
    End()
}

func OperatorRead(u USB) {
    u.Start()
    fmt.Printf("当前连接设备的信息为：%v\n", u)
    u.End()
}
```

USB 接口的意思就是，我不管要插我的是什么，鼠标键盘也罢，移动硬盘也行，U 盘也行。只要这些要插入 USB 接口的对象能满足两个方法(其实就相当于规定使用 USB 连接协议，类似 HTTP 协议)，都可以连接。连接后，你们具体怎么操作，就看你们各自对象自己的实现了。

> 假设这一段代码是一个第三方库的话，那这就是一个 usb 接口库，可以给其他人使用。只要实现了接口的方法，就可以使用我的库中的绝大部分功能。

现在有一个金士顿的硬盘想要使用 USB 接口读取数据，这时，它只要实现两个方法即可通过接口后面的 `OperatorRead()` 来获取数据。在这里金士顿硬盘的 `Start()` 方法中，有一个行为。

```go
// KingstonDisk 金士顿牌移动硬盘
type KingstonDisk struct {
    Name string
    Type string
    Data string
}

// NewKingstonDisk is
func NewKingstonDisk() *KingstonDisk {
    return &KingstonDisk{
        Name: "A1",
        Type: "SSD",
        Data: "KingstonDisk fastest SSD",
    }
}

// Start is
func (k *KingstonDisk) Start() {
    fmt.Println("金士顿SSD硬盘已连接")
}

// End is
func (k *KingstonDisk) End() {

}
```

这时，我们就可以让 KingstonDisk 插入 USB 接口了，初始化(相当于打开插头的帽)，并调用 USB 接口后面的 `OperatorRead()` 函数，调用函数时，需要将 KingstonDisk 实例化的变量作为参数传递进去(这就等于插在 USB 接口，想要执行 OperatorRead 行为)

```go
func main() {
    k := NewKingstonDisk()
    usbinterface.OperatorRead(k)
}
```

### 组织一下这段代码文件

```bash
practice/usb_interface/
├── main.go
├── usb_device
│   └── disk.go
└── usb_interface
    └── interface.go
```

**interface.go**

```go
package usbinterface

import "fmt"

// USB is
type USB interface {
    Start()
    End()
}

// OperatorRead 启动插入接口的设备并从中读取信息读取、读取后结束
func OperatorRead(u USB) {
    u.Start()
    fmt.Printf("当前连接设备的信息为：%v\n", u)
    u.End()
}
```

**disk.go**

```go
package usb_device

import "fmt"

// KingstonDisk 金士顿牌移动硬盘
type KingstonDisk struct {
    Name string
    Type string
    Data string
}

// NewKingstonDisk is
func NewKingstonDisk() *KingstonDisk {
    return &KingstonDisk{
        Name: "A1",
        Type: "SSD",
        Data: "KingstonDisk fastest SSD",
    }
}

// Start is
func (k *KingstonDisk) Start() {
    fmt.Println("金士顿SSD硬盘已连接")
}

// End is
func (k *KingstonDisk) End() {
    //
}
```

**main.go**

```go
package main

import (
    usbinterface "github.com/DesistDaydream/GoLearning/practice/usb_interface/usb_interface"
    usbdevice "github.com/DesistDaydream/GoLearning/practice/usb_interface/usb_device"
)

func main() {
    k := usbdevice.NewKingstonDisk()
    usbinterface.OperatorRead(k)
}
```

运行结果：

```bash
~]# go run practice/usb_interface/main.go
金士顿SSD硬盘已连接
当前连接的设备信息为：&{A1 SSD KingstonDisk fastest SSD}
```

#### 代码总结

从结果可以看出来，在 `OperatorRead()` 函数的代码逻辑中，并没打印 `金士顿SSD硬盘已连接` 这段文字。

> 而且 OperatorRead() 的**形参是接口类型**，但是我们在调用的时候传递的**实参是结构体类型**。可以这么说，如果结构体实现了接口内的方法，这个结构体就是接口。

所以此时 OperatorRead() 中调用的 `u.Start()` 其实是 `k.Start()`，也就是说，此时 **OperatorRead() 代码逻辑中，变量 u 的类型** 已经 **变为 `KingstonDisk` 这个结构体类型了**，所以调用的是实现 KingstonDisk 这个结构体的方法(也就是 `k.Start()`)。这种效果，就称为接口的 **polymorphism(多态)** 特性，就是说接口可以接收任意类型。

所以在 main() 中，使用 KingstonDisk 作为参数调用 usbinterface.OperatorRead(k) 时，实际上，就是把实现 KingstonDisk 的两个方法，也带入了接口中，接口中任何调用这俩方法的代码，实际都是在调用 实现 KingstonDisk 这个结构体的方法。

也可以这么说，凡是**实现了接口的结构体**，这个**结构体就是这个接口**。**如果把接口比作变量(实际上，也确实可以声明一个接口类型的变量)，那么结构体就是这个变量的值**。

还可以这么说，假设结构体 B **实现了接口**，这个**接口就是 B**。调用接口内的方法时，就是调用接口变量当前值(也就是 B)的方法。

比如我现在在 `main()` 中声明一个 map，key 设置为 USB 接口类型，value 设置为 bool 类型。当我们给这个 map 设置值的时候，凡是实现了这个接口的结构体，都可以作为该 map 的 key 使用。并且在输出 map 时，可以根据这个 key 来输出。

```go
    // usbs 测试接口多态效果，
    usbs := map[usbinterface.USB]bool{
        &usbdevice.KingstonDisk{}: true,
        &usbdevice.Mouse{}:        true, // 这是后面讲到的另一个实现了接口的结构体
    }
    fmt.Println(usbs)

    // 这是一个最简单的，把接口当作变量，把结构体当作值，然后调用接口下方法的例子
    // 此时 结构体=方法，所以在调用作用在 usbVar 变量上的方法，实际上就是 func (k *KingstonDisk) Start() {}
    var usbsVar usbinterface.USB
    usbsVar = usbdevice.KingstonDisk()
    usbsVar.Start()
```

输出结果：

```bash
~]# go run practice/usb_interface/main.go
map[0xc00005c1e0:true]
金士顿SSD硬盘已连接
```

从输出结果可以看到，usbs 这个 map 类型的 key 是当前接口体的指针。

> 为什么 key 是指针呢？因为我们在实现 KingstonDisk 结构体时，用的都是指针，所以这里的 key 就是结构体的内存地址
> 如果将实现结构体时的指针都去掉，这里 map 的 key，实际上就行结构体内的值，但是这里没有给结构体赋值，所以这是 key 应该为 `{ }`，也就是空值)。

所以，从这个示例就可以看出来，实现了接口的结构体，那么这个结构体就是这个接口。

平时口头交流老说给啥啥啥传值，在这里，就可以描述为**给接口传个结构体**。嘿嘿~~~

### 如果不使用接口

上面的例子非常繁琐，为什么我们要花这么大劲儿写个接口呢，如果不写接口行 不行？行，当然行，我来改一下~

我们只需要不使用 USB interface 类型，并修改 OperatorRead() 所接收的**形参**即可，其他文件都不用修改

```go
package usbinterface

import (
    "fmt"

    usbdevice "github.com/DesistDaydream/GoLearning/practice/usb_interface/usb_device"
)

// USB 接口
// type USB interface {
//  Start()
//  End()
// }

// OperatorRead 启动插入接口的设备并从中读取信息读取、读取后结束
func OperatorRead(u *usbdevice.KingstonDisk) {
    u.Start()
    fmt.Printf("当前连接的设备信息为：%v\n", u)
    u.End()
}
```

这样得出的输出结果时一样的。

### 不使用接口有什么问题呢？

此时，如果我添加了一个新的对象，比如鼠标、键盘等，也想通过 OperatorRead() 来执行操作并输出结果。也就是我需要新的结构体，这该怎么办呢？

首先，需要再添加一个 OperatorRead() 以便让他接收其他类型的参数。变成这样

```go
// OperatorRead 启动插入接口的设备并从中读取信息读取、读取后结束
func OperatorRead(u *usbdevice.KingstonDisk) {
    u.Start()
    fmt.Printf("当前连接的设备信息为：%v\n", u)
    u.End()
}

// OperatorRead 启动插入接口的设备并从中读取信息读取、读取后结束
func OperatorRead2(u *usbdevice.Mouse) {
    u.Start()
    fmt.Printf("当前连接的设备信息为：%v\n", u)
    u.End()
}
```

并且调用时，也需要分别调用

```go
func main() {
    k := usbdevice.NewKingstonDisk()
    usbinterface.OperatorRead(k)
    m:=usbdevice.NewMouse()
    usbinterface.OperatorReadMouse(m)
}
```

问题点：

- 修改了我们已经设计好的 OperatorRead() 功能，每当其他人想要使用我们设计的这个功能时，我都要重新修改。这大大增加了代码的不可维护性。

### 结论

而如果使用了接口，不管增加多少 USB 设备，只要这些设备想使用 OperatorRead() 的功能，直接调用即可。在添加完结构体文件后，只需要在 main() 中直接调用即可：

```go
func main() {
    k := usbdevice.NewKingstonDisk()
    usbinterface.OperatorRead(k)
    m := usbdevice.NewMouse()
    usbinterface.OperatorRead(m)
}
```

## 通过鸭子模型理解 Go 接口

> 参考：
>
> - [知乎，如何理解 Golang 中的接口](https://www.zhihu.com/question/318138275)

### 鸭子模型

那什么鸭子模型？

鸭子模型的解释，通常会用了一个非常有趣的例子，一个东西究竟是不是鸭子，取决于它的能力。游泳起来像鸭子、叫起来也像鸭子，那么就可以是鸭子。

动态语言，比如 Python 和 Javascript 天然支持这种特性，不过相对于静态语言，动态语言的类型缺乏了必要的类型检查。

Go 接口设计和鸭子模型有密切关系，但又和动态语言的鸭子模型有所区别，在编译时，即可实现必要的类型检查。

### 用 Go 接口实现鸭子模型

Go 接口是一组方法的集合，可以理解为抽象的类型。它提供了一种非侵入式的接口。任何类型，只要实现了该接口中方法集，那么就属于这个类型。

举个例子，假设定义一个鸭子的接口。如下：

```go
type Duck interface {
    Quack()   // 鸭子叫
    DuckGo()  // 鸭子走
}
```

我们定义一个函数，负责执行鸭子能做的事情。

```go
func DoDuck(d Duck) {
    d.Quack()
    d.DuckGo()
}
```

假设现在有一个鸡类型，结构如下：

```go
type Chicken struct {
}

func (c Chicken) IsChicken() bool {
    fmt.Println("我是小鸡")
}
```

这只鸡和一般的小鸡不一样，它比较聪明，也可以做鸭子能做的事情。

```go
func (c Chicken) Quack() {
    fmt.Println("嘎嘎")
}

func (c Chicken) DuckGo() {
    fmt.Println("大摇大摆的走")
}
```

> 注意，这里只是实现了 Duck 接口方法，并没有将鸡类型和鸭子接口显式绑定。这是一种非侵入式的设计。

因为小鸡实现了鸭子的所有方法，所以小鸡也是鸭。那么在 main 函数中就可以这么写了。

```go
func main() {
 // 声明结构体
    c := Chicken{}
 // 结构体类型作为实参调用 DoDuck()
    DoDuck(c)
}
```

执行正常。如此是不是很类似于其他语言的多态，其实这就是 Go 多态的实现方法。
