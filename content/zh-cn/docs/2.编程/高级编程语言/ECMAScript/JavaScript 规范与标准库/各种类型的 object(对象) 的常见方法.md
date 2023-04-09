---
title: 各种类型的 object(对象) 的常见方法
---

# 概述

> 参考：
> - [MDN 官方文档，参考-JavaScript-JavaScript-标准内置对象](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects)
> - [MDN 官方文档，参考-WebAPIs](https://developer.mozilla.org/en-US/docs/Web/API)

# String 对象

常用 String 对象的方法

- toLowerCase() # 将字符串内的字母全部转换成小写
- toUpperCase() 将字符串内的字母全部转换成大写
- replace("D", 1) # replace(searchValue,replaceValue) 将字符串内第一个满足 searchValue 条件的字符替换为 replaceValue。注意：只能替换第一个
- trim() # 去除首尾所有空白字符
- split(" ") # 按照分隔符将字符串切割为一个数组。注意：只有字符串中有指定的分隔符，才会生效。否则切割后的元素只有一个。
- 截取字符串
  - substr(5, 8) # 第一个参数是开始截取的索引号，第二个参数是截取数量
  - substring(5, 8) # 第一个参数是开始截取的索引号，第二个参数是结束截取的索引号
  - slice(5, 8) # 第一个参数是开始截取的索引号，第二个参数是结束截取的索引号

# Array 对象

常用 Array 对象的方法

- 会改变原始数组的内容
  - push() # 从后面追加
  - pop() # 从后面删除
  - unshift() # 从前面添加
  - shift() # 从前面删除
  - reverse() # 反转数组
  - splice() # 截取并添加
  - sort() # 数组排序
- 不会改变原始数组的内容
  - join() # 数组连接为字符串
  - concat() # 拼接数组
  - slice() # 截取数组
  - indexOf() # 查找元素在数组中的索引
  - forEach() # 遍历数组
  - map() # 映射数组
  - filter() # 过滤数组
  - every() # 判断是否全部满足条件
  - some() # 判断是否有满足条件的项
