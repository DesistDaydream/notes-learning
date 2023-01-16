---
title: Vue 响应式
---

# 概述

> 参考：
> - [官方文档，API 参考](https://cn.vuejs.org/api/)
> - [官方文档，API 参考-组合式 API-响应式：核心](https://cn.vuejs.org/api/reactivity-core.html)

ref 与 reactive 是响应式的基础

# 响应式: 核心

## [ref()](https://cn.vuejs.org/api/reactivity-core.html#ref)

`ref()` 函数返回一个 `Ref<T = any>` 接口类型的对象，该接口中只有一个名为 `value` 的属性，用以指向该对象的值。
`Ref<T>` 接口对象是 **响应式**、**可更改** 的。

```javascript
let number = ref < number > 0
```

number 是 `Ref<number>` 类型的实例，Ref.value 则是该实例的值，即 `0`。
Ref 对象是可更改的，也就是说你可以为 .value 赋予新的值。它也是响应式的，即所有对 .value 的操作都将被追踪，并且写操作会触发与之相关的副作用。
如果将一个对象赋值给 ref，那么这个对象将通过 [reactive()](https://cn.vuejs.org/api/reactivity-core.html#reactive) 转为具有深层次响应式的对象。这也意味着如果对象中包含了嵌套的 ref，它们将被深层地解包。
若要避免这种深层次的转换，使用 [shallowRef()](https://cn.vuejs.org/api/reactivity-advanced.html#shallowref) 来替代。

## [computed ()](https://cn.vuejs.org/api/reactivity-core.html#computed)

## [reactive()](https://cn.vuejs.org/api/reactivity-core.html#reactive)

- [readonly()](https://cn.vuejs.org/api/reactivity-core.html#readonly)
- [watchEffect()](https://cn.vuejs.org/api/reactivity-core.html#watcheffect)
- [watchPostEffect()](https://cn.vuejs.org/api/reactivity-core.html#watchposteffect)
- [watchSyncEffect()](https://cn.vuejs.org/api/reactivity-core.html#watchsynceffect)
- [watch()](https://cn.vuejs.org/api/reactivity-core.html#watch)

# 响应式: 工具

- [isRef()](https://cn.vuejs.org/api/reactivity-utilities.html#isref)
- [unref()](https://cn.vuejs.org/api/reactivity-utilities.html#unref)
- [toRef()](https://cn.vuejs.org/api/reactivity-utilities.html#toref)
- [toRefs()](https://cn.vuejs.org/api/reactivity-utilities.html#torefs)
- [isProxy()](https://cn.vuejs.org/api/reactivity-utilities.html#isproxy)
- [isReactive()](https://cn.vuejs.org/api/reactivity-utilities.html#isreactive)
- [isReadonly()](https://cn.vuejs.org/api/reactivity-utilities.html#isreadonly)

# 响应式: 进阶

- [shallowRef()](https://cn.vuejs.org/api/reactivity-advanced.html#shallowref)
- [triggerRef()](https://cn.vuejs.org/api/reactivity-advanced.html#triggerref)
- [customRef()](https://cn.vuejs.org/api/reactivity-advanced.html#customref)
- [shallowReactive()](https://cn.vuejs.org/api/reactivity-advanced.html#shallowreactive)
- [shallowReadonly()](https://cn.vuejs.org/api/reactivity-advanced.html#shallowreadonly)
- [toRaw()](https://cn.vuejs.org/api/reactivity-advanced.html#toraw)
- [markRaw()](https://cn.vuejs.org/api/reactivity-advanced.html#markraw)
- [effectScope()](https://cn.vuejs.org/api/reactivity-advanced.html#effectscope)
- [getCurrentScope()](https://cn.vuejs.org/api/reactivity-advanced.html#getcurrentscope)
- [onScopeDispose()](https://cn.vuejs.org/api/reactivity-advanced.html#onscopedispose)

# 生命周期钩子

- [onMounted()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onmounted)
- [onUpdated()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onupdated)
- [onUnmounted()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onunmounted)
- [onBeforeMount()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onbeforemount)
- [onBeforeUpdate()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onbeforeupdate)
- [onBeforeUnmount()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onbeforeunmount)
- [onErrorCaptured()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onerrorcaptured)
- [onRenderTracked()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onrendertracked)
- [onRenderTriggered()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onrendertriggered)
- [onActivated()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onactivated)
- [onDeactivated()](https://cn.vuejs.org/api/composition-api-lifecycle.html#ondeactivated)
- [onServerPrefetch()](https://cn.vuejs.org/api/composition-api-lifecycle.html#onserverprefetch)
