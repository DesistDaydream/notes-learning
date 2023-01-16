---
title: string 组
---

## APPEND key value # 附加一个 value 到指定的 key

since: 2.0.0

BITCOUNT key \[start end]

summary: Count set bits in a string

since: 2.6.0

BITFIELD key \[GET type offset] \[SET type offset value] \[INCRBY type offset increment] \[OVERFLOW WRAP|SAT|FAIL]

summary: Perform arbitrary bitfield integer operations on strings

since: 3.2.0

BITOP operation destkey key \[key ...]

summary: Perform bitwise operations between strings

since: 2.6.0

BITPOS key bit \[start] \[end]

summary: Find first bit set or clear in a string

since: 2.8.7

DECR key

summary: Decrement the integer value of a key by one

since: 1.0.0

DECRBY key decrement

summary: Decrement the integer value of a key by the given number

since: 1.0.0

## GET key # summary: 获取指定 key 的 value

since: 1.0.0

GETBIT key offset

summary: Returns the bit value at offset in the string value stored at key

since: 2.2.0

GETRANGE key start end

summary: Get a substring of the string stored at a key

since: 2.4.0

GETSET key value

summary: Set the string value of a key and return its old value

since: 1.0.0

INCR key

summary: Increment the integer value of a key by one

since: 1.0.0

INCRBY key increment

summary: Increment the integer value of a key by the given amount

since: 1.0.0

INCRBYFLOAT key increment

summary: Increment the float value of a key by the given amount

since: 2.6.0

MGET key \[key ...]

summary: Get the values of all the given keys

since: 1.0.0

MSET key value \[key value ...]

summary: Set multiple keys to multiple values

since: 1.0.1

MSETNX key value \[key value ...]

summary: Set multiple keys to multiple values, only if none of the keys exist

since: 1.0.1

PSETEX key milliseconds value

summary: Set the value and expiration in milliseconds of a key

since: 2.6.0

## SET key value \[EX seconds] \[PX milliseconds] \[NX|XX] # 设置指定 key 的 字符串类型的 value

since: 1.0.0

如果 key 已经保存了一个值，那么这个操作会直接覆盖原来的值，并且忽略原始类型。当 set 命令执行成功之后，之前设置的过期时间都将失效。

SETBIT key offset value

summary: Sets or clears the bit at offset in the string value stored at key

since: 2.2.0

SETEX key seconds value

summary: Set the value and expiration of a key

since: 2.0.0

SETNX key value

summary: Set the value of a key, only if the key does not exist

since: 1.0.0

SETRANGE key offset value

summary: Overwrite part of a string at key starting at the specified offset

since: 2.2.0

STRLEN key

summary: Get the length of the value stored in a key

since: 2.2.0
