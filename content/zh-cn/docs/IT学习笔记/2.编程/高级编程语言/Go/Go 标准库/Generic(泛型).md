---
title: Generic(泛型)
---

# 概述

> 参考：
> - [Wiki,Generic proramming](https://en.wikipedia.org/wiki/Generic_programming)(泛型编程)
> - [公众号-OSC 开源社区，使用 Go 泛型的最佳时机](https://mp.weixin.qq.com/s/Ymxs4Z2p62hQ7RA3q23drg)
> - [公众号-OSC 开源社区，Go 语言之父介绍泛型](https://mp.weixin.qq.com/s/MVZxoh8pYYUJ1_Dyj8Sh7g)
> - [公众号-InfoQ，Go 中的泛型：激动人心的突破](https://mp.weixin.qq.com/s/Zk24GsvpryB64hlSAp06Iw)
> - [公众号-幼麟实验室，Go 泛型的实现](https://mp.weixin.qq.com/s/3Yc7_mvYdr3KtSE0jnAJww)

泛型让我们在定义接口、函数、结构体时将其中的类型**参数化**。我们从古老的 Ada 语言的第一个版本就开始使用泛型了，后来 C++ 的模板中也有泛型，直到 Java 和 C# 中的现代实现都是很常见的例子。

泛型为 Go 添加了三个新的重要内容：

- 面向函数和类型的“类型形参” (type parameters)
- 将接口类型定义为类型集合，包括没有方法的接口类型
- 类型推断：在大多数情况下，在调用泛型函数时可省略“类型实参” (type arguments)

## 类型形参与约束

下面是一个初步理解泛型的最简单例子：

```go
// 泛型
// 使用类型形参编写 Go 函数以处理多种类型
// comparable 是一个内置 Constraint(约束)，用来表示类型形参可以接收的类型实参的种类，所谓的“约束”就是指，T 被约束为可以使用哪几种类型。
// comparable 包含所有可以比较类型，包括：booleans、numbers、strings、pointers、channels、可比较的 arrays、structs 中的属性
// comparable 可以改为 any，表示 T 可以是任意类型
func Index[T comparable](s []T, x T) int {
	for i, v := range s {
		// 这里的 v 和 x 都是 T 类型，且 T 类型具有可以比较的约束，因此这里可以使用 ==
		if v == x {
			return i
		}
	}
	return -1
}

func main() {
	// Index() 函数适用于 int 类型的切片
	si := []int{10, 20, 15, -10}
	fmt.Println(Index(si, 15))

	//  Index() 函数同样也使用于 string 类型的切片
	ss := []string{"foo", "bar", "baz"}
	fmt.Println(Index(ss, "hello"))
}
```

可以看到，我们将 a、b 的类型**参数化**了。这里的 T，可以称之为**泛型类型**，同时也是 Max() 方法的 **Type Parameters(类型形参)**。当我们调用 Max() 时，可以传递参数告诉 `Max()` 当前应该使用类型执行其中的代码

如果没有泛型，那么我们的 `Index()` 函数就要写两遍(有多少种类型，就要写多少遍)：

```go
func IndexInt(s []int, x int) int {
    // Do Somthing
}

func IndexString(s []string, x string) int {
    // Do Somthing
}
```

还可以使用 any 关键字，以便让约束变为任意类型

```go

func Do[R any, S any, T any](a R, b S) T {
  // Do Somthing
}
func main() {
  fmt.Println(Do[int, uint, float64](1, 2))
}

// 上面的代码的行为与下面的函数完全相同:
// Do(a int, b uint) float64
```

# 其他

## 类型形参

现在函数和类型都具有类型形参” (type parameters)，类型形参列表看起来就是一个普通的参数列表，除了它使用的是方括号而不是小括号。

先从浮点值的基本非泛型 Min 函数开始：

```go
func Min(x, y float64) float64 {
    if x < y {
        return x
    }
    return y
}
```

通过添加类型形参列表来使这个函数泛型化——使其适用于不同的类型。在此示例中，添加了一个带有单个类型形参`T`的类型参数列表，并替换了`float64`。

```go
import "golang.org/x/exp/constraints"

func GMin[T constraints.Ordered](x, y T) T {
    if x < y {
        return x
    }
    return y
}
```

然后就可以使用类型实参调用此函数：

```go
x := GMin[int](2, 3)
```

向  GMin 提供类型参数，在这种情况下`int`称为实例化。实例化分两步进行。首先，编译器在泛型函数或泛型类型中用所有类型形参替换它们各自的类型实参。然后，编译器验证每个类型形参是否满足各自的约束。如果第二步失败，实例化就会失败并且程序无效。

成功实例化后，即可产生非泛型函数，它可以像任何其他函数一样被调用。比如：

```go
fmin := GMin[float64]
m := fmin(2.71, 3.14)
```

`GMin[float64]`的实例化产生了一个与`Min`函数等效的函数，可以在函数调用中使用它。类型形参也可以与类型一起使用。

```go
type Tree[T interface{}] struct {
    left, right *Tree[T]
    value       T
}

func (t *Tree[T]) Lookup(x T) *Tree[T] { ... }

var stringTree Tree[string]
```

在上面的例子中，泛型类型`Tree`存储了类型形参`T`的值。泛型类型也可以有方法，比如本例中的`Lookup`。为了使用泛型类型，它必须被实例化；`Tree[string]`是使用类型实参`string`来实例化`Tree`的示例。

## 类型推断

此项功能是最复杂的变更，主要包括：

- 函数参数类型推断 (Function argument type inference)
- 约束类型推断 (Constraint type inference)

虽然类型推断的工作原理细节很复杂，但使用它并不复杂：类型推断要么成功，要么失败。如果它成功，类型实参可以被省略，调用泛型函数看起来与调用普通函数没有什么不同。如果类型推断失败，编译器将给出错误消息，在这种情况下，只需提供必要的类型实参。

---

泛型是 Go 1.18 的重要新语言特性，Robert Griesemer 和 Ian Lance Taylor 表示，这个功能实现得很好并且质量很高。虽然他们鼓励在有必要的场景中使用泛型，但在生产环境中部署泛型代码时，请务必谨慎。

更多介绍查看原文：https://go.dev/blog/intro-generics。

# 使用泛型的最佳时机

从历史上看，C++、D 乃至 Rust 等系统语言一直采用单态化方法实现泛型。然而，Go 1.18 的泛型实现并不完全依靠单态化 (Monomorphization)，而是采用了一种被称为"GCShape stenciling with Dictionaries"的部分单态化技术。这种方法的好处是可以大幅减少代码量，但在特定情况下，会导致代码速度变慢。

Ian Lance Taylor 表示，Go 的通用开发准则有要求：开发者应通过编写代码而不是定义类型来编写 Go 程序。当涉及到泛型时，如果通过定义类型参数约束来编写程序，那一开始就走错了路。正解应该是从编写函数开始，当明确了类型参数的作用后，再添加类型参数就很容易了。

接着，Ian 列举了 4 种类型参数能有效发挥作用的情况：

1. 使用语言定义的特殊容器类型
2. 通用数据结构
3. 类型参数首选是函数，而非方法的情况
4. 不同类型需要实现通用方法

同时也提醒了不适合使用类型参数的情况：

1. 不要使用类型参数替换接口类型 (Interface Type)
2. 如果方法实现不同，不要使用类型参数
3. 在适当的地方使用反射 (reflection)

最后，Ian 给出了简要的泛型使用方针，当开发者发现自己多次编写完全相同的代码，而这些副本之间的唯一区别仅在于使用了不同类型，这时候便可以考虑使用类型参数。换句话说，即开发者应避免使用类型参数，直到发现自己要多次编写完全相同的代码。

## 最佳实践

泛型的最佳时机，就是当两个“对象”的某些“方法”的行为完全一样时，我们可以通过泛型来声明这些方法，以防止重复编写完全相同的代码。就像下面这段代码中，ProductGorm 与 UserGorm 的两个方法代码其实是完全相同的，如果使用 interface{}，那么我们需要重复写了 Create() 与 Get() 方法。其实，本质上 ProductGorm 与 UserGorm 所需要执行的 SQL 是完全一样，不一样的只是其中的列名而已。

```go
package main

import (
	"fmt"

	"github.com/go-gorm/driver/sqlite"
	"github.com/go-gorm/gorm"
)

type ProductGorm struct {
	gorm.Model
	Name  string
	Price uint
}

type UserGorm struct {
	gorm.Model
	FirstName string
	LastName  string
}

type Repository[T any] struct {
	db *gorm.DB
}

func (r *Repository[T]) Create(t T) error {
	return r.db.Create(&t).Error
}

func (r *Repository[T]) Get(id uint) (*T, error) {
	var t T
	err := r.db.Where("id = ?", id).First(&t).Error
	return &t, err
}

func NewProductRepository(db *gorm.DB) *Repository[ProductGorm] {
	db.AutoMigrate(&ProductGorm{})
	return &Repository[ProductGorm]{
		db: db,
	}
}

func NewUserRepository(db *gorm.DB) *Repository[UserGorm] {
	db.AutoMigrate(&UserGorm{})
	return &Repository[UserGorm]{
		db: db,
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	productRepo := NewProductRepository(db)
	productRepo.Create(ProductGorm{
		Name:  "product",
		Price: 100,
	})
	fmt.Println(productRepo.Get(1))
	userRepo := NewUserRepository(db)
	userRepo.Create(UserGorm{
		FirstName: "first",
		LastName:  "last",
	})
	fmt.Println(userRepo.Get(1))
}

//
//
// Out:
// &{{1 2021-11-23 22:50:14.595342 +0100 +0100 2021-11-23 22:50:14.595342 +0100 +0100 {0001-01-01 00:00:00 +0000 UTC false}}  100} <nil>
// &{{1 2021-11-23 22:50:44.802705 +0100 +0100 2021-11-23 22:50:44.802705 +0100
```
