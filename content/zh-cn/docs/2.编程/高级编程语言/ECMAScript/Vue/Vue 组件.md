---
title: Vue 组件
---

# 概述

>

# 组件间数据传递

在 Vue 中，组件之间可以传递多种类型的数据

- 变量，通过 Props。
  - 关键字：v-bind
- 模板，通过 Slots
  - 关键字：\<template>、v-slot
- 事件，通过 Event
  - 关键字：v-model、v-on、emit

## Props

## Slots

**Slot Content(插槽内容)** 与 **Slot Outlet(插槽出口)**
`<slot>` 元素是一个 **Slot Outlet(插槽出口)**，标示了父元素提供的 **Slot Content(插槽内容) **将在哪里被渲染。
![image.png](https://notes-learning.oss-cn-beijing.aliyuncs.com/bigg4rx374wdctrg/1669263120268-c5d8e233-599c-4779-b9b0-3bcb15e9ff2a.png)

## Event
