---
title: Data type
linkTitle: Data type
date: 2024-06-12T13:21
weight: 20
---

# 概述

> 参考：
>
> - [Wiki，DataType](https://en.wikipedia.org/wiki/Data_type)

在计算机科学和计算机编程中，**Data Type(数据类型，有时也简称 Type)** 是数据的一个属性，这些属性将会让编译器知道程序员想要如何使用数据。

[Literal(字面量)](/docs/2.编程/计算机科学/Data%20type/Literal.md) 与 [Variable(变量)](/docs/2.编程/计算机科学/Variable.md) 相关，是用于初始化变量时指定的一个值。

## 数据类型的分类

- 原始数据类型
- 复合数据类型
- 抽象数据类型
- 其他类型
- TODO

# Primitive Data Types(原始数据类型)

[原始数据类型](https://en.wikipedia.org/wiki/Primitive_data_type)通常是语言实现的内置或基础类型。

### Machine Data Type(机器数据类型)

基于数字电子的计算机中的所有数据都表示为最低级别的[位](https://en.wikipedia.org/wiki/Bit)（替代 0 和 1）。数据的最小可寻址单元通常是一组称为[字节](https://en.wikipedia.org/wiki/Byte)的位（通常是一个[八位组](<https://en.wikipedia.org/wiki/Octet_(computing)>)，即 8 位）。由[机器代码](https://en.wikipedia.org/wiki/Machine_code)指令处理的单元称为[字](<https://en.wikipedia.org/wiki/Word_(data_type)>)（截至 2011 年，通常为 32 或 64 位）。大多数指令将字解释为[二进制数](https://en.wikipedia.org/wiki/Binary_number)，因此 32 位字可以表示从 0 到 232 - 1 或有符号整数值来自 -231 到 231 - 1 由于[二进制补码](https://en.wikipedia.org/wiki/Two%27s_complement)，机器语言和机器在大多数情况下不需要区分这些无符号和有符号数据类型。

用于浮点算术的浮点数对字中的位使用不同的解释。有关详细信息，请参阅[浮点运算](https://en.wikipedia.org/wiki/Floating-point_arithmetic)。

机器数据类型需要在[系统](https://en.wikipedia.org/wiki/Systems_programming)或[低级编程语言中](https://en.wikipedia.org/wiki/Low-level_programming_language)公开或可用，允许对硬件进行细粒度控制。的[C 编程语言](https://en.wikipedia.org/wiki/C_programming_language)，例如，建筑材料整数类型不同的宽度，如和。如果目标平台上不存在相应的本机类型，编译器将使用确实存在的类型将它们分解为代码。例如，如果在 16 位平台上请求一个 32 位整数，编译器会默认将其视为两个 16 位整数的数组。 shortlong

在更高级别的编程中，机器数据类型通常被隐藏或\_抽象\_为一个实现细节，如果暴露，会使代码的可移植性降低。例如，numeric 可以提供泛型类型而不是某些特定位宽的整数。

### Boolean Type(布尔类型)

[Boolean(布尔)](https://en.wikipedia.org/wiki/Boolean_type) 类型表示值 [true(真)](https://en.wikipedia.org/wiki/Logical_truth) 和 [false(假)](https://en.wikipedia.org/wiki/Logical_truth)。尽管只有两个值是可能的，但出于效率原因，它们很少被实现为单个二进制数字。许多编程语言没有明确的布尔类型，**而是将 0 解释为 false，将其他值解释为 true**。布尔数据是指如何将语言解释为机器语言的逻辑结构。在这种情况下，布尔值 0 指的是逻辑 False。True 总是非零，尤其是被称为布尔值 1 的一。

### Numeric Type(数字类型)

- [Integer(整数，简写 int)](<https://en.wikipedia.org/wiki/Integer_(computing)>) 数据类型，或“非分数”。可以根据它们包含负值的能力进行子类型化（例如 unsigned 在 C 和 C++ 中）。也可具有小的预定义数目的亚型（如 short 和 long 在 C / C ++）; 或允许用户自由定义子范围，例如 1..12（例如[Pascal](<https://en.wikipedia.org/wiki/Pascal_(programming_language)>) / [Ada](<https://en.wikipedia.org/wiki/Ada_(programming_language)>)）。
- [Floating Point(浮点)](https://en.wikipedia.org/wiki/Floating_point) 数据类型通常将值表示为高精度分数值（[有理数](https://en.wikipedia.org/wiki/Rational_numbers)，数学上），但有时会误导性地称为实数（令人联想到数学[实数](https://en.wikipedia.org/wiki/Real_numbers)）。它们通常对最大值和精度都有预定义的限制。通常以 a × 2 b 的形式在内部存储（其中 a 和 b 是整数），但以熟悉的[十进制](https://en.wikipedia.org/wiki/Decimal)形式显示。
- [Fixed Point(定点)](<https://en.wikipedia.org/wiki/Fixed_point_(computing)>) 数据类型便于表示货币值。它们通常在内部实现为整数，从而导致预定义的限制。
- [Bignum](https://en.wikipedia.org/wiki/Bignum)或[任意精度](https://en.wikipedia.org/wiki/Arbitrary_precision)数字类型缺乏预定义的限制。它们不是原始类型，出于效率原因很少使用。

### Enumerations(枚举)

[枚举类型](https://en.wikipedia.org/wiki/Enumerated_type)具有不同的值，其可以被比较和分配，但不一定必须在计算机的存储器中的任何特定的具体表示; 编译器和解释器可以任意表示它们。例如，一副扑克牌中的四个花色可能是名为 CLUB、DIAMOND、HEART、SPADE 的四个枚举数，属于一个名为 suit 的枚举类型。如果变量 V 被声明为具有花色作为它的数据类型，可以为它分配这四个值中的任何一个。一些实现允许程序员为枚举值分配整数值，甚至将它们视为与整数类型等效的。

# Composite Types(复合类型)

[复合数据类型](/docs/2.编程/计算机科学/Data%20Type/复合数据类型/复合数据类型.md) 派生自多个原始类型。这可以通过多种方式完成。它们组合的方式称为[数据结构](https://en.wikipedia.org/wiki/Data_structure)

- **[Array](/docs/2.编程/计算机科学/Data%20Type/复合数据类型/Array.md)(数组)** # 也称为 **Vector(向量)**、[**List(列表)**](https://en.wikipedia.org/wiki/List_(abstract_data_type))、**Sequence(序列)**。Array 存储多个 **elements(元素)**，并提供对各个 elements 的随机访问。数组的元素通常（但并非在所有情况中）要求具有相同类型。数组可以是固定长度的或可扩展的。数组的索引通常要求是来自特定范围的整数（如果不是，可以参考 [Associative array(关联数组)](https://en.wikipedia.org/wiki/Associative_array)）（如果不是该范围内的所有索引都对应于元素，则它可能是 [Sparse array(稀疏数组)](https://en.wikipedia.org/wiki/Sparse_array)）。
- **[Record](https://en.wikipedia.org/wiki/Record_(computer_science))(记录)**  # 也称为 也称为 **Tuple(元组)** 或 **Struct(结构体)**。Record 是最简单的[数据结构](https://en.wikipedia.org/wiki/Data_structure)之一。Record 是包含其他值的值，通常采用固定数量和顺序，通常按 1 名称索引。记录的元素通常称为 **Fields(字段)** 或 **Members(成员)**。
- **[Object(对象)](https://en.wikipedia.org/wiki/Object_(computer_science))** 包含许多数据字段，如 Record，以及许多用于访问或修改它们的子程序，称为 [Methods(方法)](<https://en.wikipedia.org/wiki/Method_(computer_programming)>)。

# Abstract Data Types(抽象数据类型)

[抽象数据类型](/docs/2.编程/计算机科学/Data%20Type/抽象数据类型/抽象数据类型.md)是不指定数据具体表示的数据类型。相反，使用基于数据类型操作的正式规范来描述它。规范的任何实现都必须满足给定的规则。例如，[Heap and Stack](/docs/2.编程/计算机科学/Heap%20and%20Stack.md) 具有遵循后进先出规则的入栈/出栈操作，并且可以使用 list 或 array 来具体实现。抽象数据类型用于形式语义和程序验证，以及不太严格的设计中。

# 其他类型

类型可以基于或派生自上述基本类型。在某些语言（例如 C）中，[函数](<https://en.wikipedia.org/wiki/Function_(computer_science)>)具有从其[返回值](https://en.wikipedia.org/wiki/Return_value)的类型派生的类型。

## String(字符串) 和 Text(文本) 类型

- 一个[字符](<https://en.wikipedia.org/wiki/Character_(computing)>)，可能是某个[字母表](https://en.wikipedia.org/wiki/Alphabet)中的一个[字母](https://en.wikipedia.org/wiki/Alphabet)、一个数字、一个空格、一个标点符号等。
- 一个[字符串](<https://en.wikipedia.org/wiki/String_(computer_science)>)，它是一个字符序列。字符串通常用于表示单词和文本，尽管除了最琐碎的情况外，所有文本都不仅仅涉及字符序列。

字符和字符串类型可以存储字符集（例如 [ASCII](https://en.wikipedia.org/wiki/ASCII) 中的字符序列。由于大多数字符集都包含[数字](https://en.wikipedia.org/wiki/Numerical_digit)，因此可以使用数字字符串，例如"1234". 但是，许多语言将它们视为属于与数值不同的类型 1234。

根据所需的字符“宽度”，字符和字符串类型可以有不同的子类型。最初的 7 位宽 ASCII 被发现是有限的，并被 8 位和 16 位集取代，它们可以编码各种各样的非拉丁字母（如[希伯来语](https://en.wikipedia.org/wiki/Hebrew)和[中文](https://en.wikipedia.org/wiki/Chinese_language)）和其他符号。字符串可以是适合拉伸的，也可以是固定大小的，即使是在相同的编程语言中。它们也可以按其最大大小进行子类型化。

注意：字符串不是所有语言中的原始数据类型。例如，在 [C 语言](<https://en.wikipedia.org/wiki/C_(programming_language)>) 中，它们由字符数组组成。

## Pointer(指针) 和 Reference(引用) 类型

主要的非复合派生类型是[指针](<https://en.wikipedia.org/wiki/Pointer_(computer_programming)>)，这是一种数据类型，其值直接引用（或“指向”）使用其[地址](https://en.wikipedia.org/wiki/Memory_address)存储在[计算机内存中](https://en.wikipedia.org/wiki/Computer_memory)其他位置的另一个值。它是一种原始的[参考](<https://en.wikipedia.org/wiki/Reference_(computer_science)>)。（在日常生活中，一本书的页码可以被认为是引用另一本书的一段数据）。指针通常以类似于整数的格式存储；但是，尝试取消引用或“查找”其值永远不是有效内存地址的指针会导致程序崩溃。为了改善这个潜在问题，指针被认为是指向它们指向的数据类型的单独类型，即使底层表示相同。

## Function(函数) 类型

虽然也可以为函数分配类型，但在本文的设置中，它们的类型不被视为数据类型。在这里，数据被视为不同于[算法](https://en.wikipedia.org/wiki/Algorithm)。在编程中，函数与后者密切相关。但是，因为[通用数据处理的](https://en.wikipedia.org/wiki/Universal_Turing_machine)一个中心原则是算法可以[表示为数据](https://en.wikipedia.org/wiki/G%C3%B6del_numbering#Generalizations)，例如文本描述和二进制程序，数据和函数之间的对比是有限的。其实函数不仅可以用数据来表示，函数也可以用来[对数据](https://en.wikipedia.org/wiki/Lambda_calculus#Encoding_datatypes)进行[编码](https://en.wikipedia.org/wiki/Lambda_calculus#Encoding_datatypes)。许多当代[类型系统](https://en.wikipedia.org/wiki/Type_systems)强烈关注函数类型，许多现代语言允许函数作为[一等公民运行](https://en.wikipedia.org/wiki/First-class_citizen)。
将函数从被视为数据类型的对象中排除在相关领域中并不少见。\[[需要引用](https://en.wikipedia.org/wiki/Wikipedia:Citation_needed)] 例如，[谓词逻辑](https://en.wikipedia.org/wiki/Predicate_logic)不允许在函数或谓词名称上应用[量词](<https://en.wikipedia.org/wiki/Quantifier_(logic)>)。

## Meta(元) 类型

一些编程语言将类型信息表示为数据，从而实现[类型自省](https://en.wikipedia.org/wiki/Type_introspection)和[反射](<https://en.wikipedia.org/wiki/Reflection_(computer_programming)>)。相比之下，[高阶](https://en.wikipedia.org/wiki/Type_constructor) [类型系统](https://en.wikipedia.org/wiki/Type_systems)虽然允许从其他类型构造类型并作为值传递给函数，但通常避免基于它们进行[计算](https://en.wikipedia.org/wiki/Computational)决策。\[[需要引用](https://en.wikipedia.org/wiki/Wikipedia:Citation_needed)]

## Utility(实用程序) 类型

为方便起见，高级语言可能提供现成的“现实世界”数据类型，例如时间、日期、货币值和内存，即使该语言允许从原始类型构建它们。
