---
title: geo/geoip 模块指令
---

# 概述

> 参考：
>
> - [http 模块下的 geo 模块](http://nginx.org/en/docs/http/ngx_http_geo_module.html)、[geoip 模块](http://nginx.org/en/docs/http/ngx_http_geoip_module.html)
> - [stream 模块下的 geo 模块](http://nginx.org/en/docs/stream/ngx_stream_geo_module.html)、[geoip 模块](http://nginx.org/en/docs/stream/ngx_stream_geoip_module.html)

geo 与 geoip 模块实现了 [GeoIP](/docs/4.数据通信/Protocol/TCP_IP/IP/GeoIP.md) 的能力，可以根据 客户端的 IP 地址 来创建新的变量。这些变量用来表示 IP 地址所属国际、所属城市、所在经/纬度 等等。

不同之处在于：

- geo 手动设置变量及其值
- geoip 根据 [MaxMind](http://www.maxmind.com/) 数据库中的信息，创建一系列的变量

通过 geo/geoip 模块，我们可以根据客户端的 IP 地址，获取这些 IP 的一些信息，比如 IP 所属城市、所属国家，所在经/纬度 等等。我们常常可以根据这些分类的信息，**进行 IP 过滤、或日志记录**。说白了，geo/geoip 模块就是为每个 IP 地址添加一系列的 **Label(标签)**，以便后续可以根据这些 标签 进行 记录 和 筛选。

> v2ray 中的测试，有很多 geosite 相关的设置，就是这个道理，每个 IP 地址都可以具有很多标签、甚至连这个 IP 所属的公司都会记录，以便可以根据这些进行来决定一个请求是 直连 还是 代理。

# http 模块下的 geo 模块指令

[geo \[ADDRESS\] $VARIABLE {}](http://nginx.org/en/docs/http/ngx_http_geo_module.html#geo) # 根据 ADDRESS 定义新的变量

ADDRESS 用来指定要设置变量的 IP 地址。默认来自于 `$remote_addr` 变量，也可以自定义为另一个变量。

假如现在有这么一个配置：

```nginx
geo $country {
    default        ZZ;
    include        conf/geo.conf;
    delete         127.0.0.0/16;
    proxy          192.168.100.0/24;
    proxy          2001:0db8::/32;

    127.0.0.0/24   US;
    127.0.0.1/32   RU;
    10.1.0.0/16    RU;
    192.168.1.0/24 UK;
}
```

表示 Nginx 将会默认根据 `$remote_addr` 变量的值设置一个 `$country` 变量，也就是根据当前请求的 客户端 IP 地址 来这设置一个 `$country` 变量。

- default # 默认情况下，所有 IP 地址的 `$country` 变量的值为 ZZ
- include # 将指令写在其他文件中，并通过 include 包含进来
- delete # 删除 127.0.0.0/16 这段 IP 地址的 `$country` 变量
- 127.0.0.0/24   US; # 这些 IP 中的 127.0.0.0/24 这一段的的值为 US；127.0.0.1/32 这个
  - 后面的以此类推

# http 模块下的 geoip 模块指令

> [!Warning]
> 由于隐私的原因，[MaxMind 在 2019 年 12 月份对数据库进行重大变更](https://blog.maxmind.com/2019/12/18/significant-changes-to-accessing-and-using-geolite2-databases/)，所以，老式的 geo/geoip 模块不再适用于新的 MaxMind 数据库，所以，[geoip2 模块](/docs/Web/Nginx/Nginx%20配置详解/多用途模块的指令/geoip2%20模块.md)诞生了。
>
> 其他人维护的 geoip 数据库
>
> - https://www.miyuru.lk/geoiplegacy

geoip 模块并不需要手动设定想要创建的变量，由于 geoip 模块是基于 MaxMind 数据库，所以会根据 geoip 中的指令来创建对应的变量。

geoip 模块会从 Nginx 获取客户端的 IP 地址(一般都是 $remote_addr 变量中的值)，然后根据 IP 地址从 MaxMind 数据库查找对应的 IP，并将与该 IP 关联的信息写入到新的变量中。

## geoip_country FILE

http://nginx.org/en/docs/http/ngx_http_geoip_module.html#geoip_country

围绕国家创建相关变量

该指令一般情况下，将会根据 MaxMind 中的 GeoIP.dat 文件，为每个 IP 地址创建如下变量：

- **$geoip_country_code** # 两个字母的国家代码，比如 CN、US
- **$geoip_country_code3** # 三个字母的国家代码，比如 CHN、USA
- **$geoip_country_name** # 国家名称，比如 China、United States

## geoip_city FILE

http://nginx.org/en/docs/http/ngx_http_geoip_module.html#geoip_city

围绕城市创建相关变量。最常用指令

该指令一般情况下，将会根据 MaxMind 中的 GeoLiteCity.dat 文件，为每个 IP 地址创建如下变量：

- **$geoip_area_code** # 电话区号(仅限 US 有用).
  - 此变量可能包含过时的信息，因为不推荐使用相应的数据库字段。
- **$geoip_city_continent_code** # 两个字母的大洲代码，比如 EU、NA
- **$geoip_city_country_code** # 两个字母的国家代码，比如 CN、US
- **$geoip_city_country_code3** # 三个字母的国家代码，比如 CHN、USA
- **$geoip_city_country_name** # 国家名称，比如 China、United States
- **$geoip_dma_code** # 根据 Google AdWords API 中的地理定位，美国的 DMA 区域代码 (也称为 metro code)。
- **$geoip_latitude** # latitude(纬度)
- **$geoip_longitude** # longitude(经度)
- **$geoip_region** # 两个字母的地区代码，比如 JX
  - region 就是省的概念，直辖市的 省 与 市 同名
- **$geoip_region_name** # 两个字母的地区名称，比如 Jiangxi
- **$geoip_city** # 城市名称，比如，Beijing、Tianjin
- **$geoip_postal_code** # 邮政编码

这些变量是最常用的，用来描述 IP 地址关于城市相关的信息，比如客户端的 IP 地址为 111.33.112.94，那么这个 IP 就属于天津市，并且还会有关于该城市的相关信息。比如城市所属国家、所属大洲、所在经纬度、邮编、城市代码 等等信息。

## geoip_org FILE

- http://nginx.org/en/docs/http/ngx_http_geoip_module.html#geoip_org

该指令一般情况下，将会根据 MaxMind 中的 GeoLiteCity.dat 文件创建如下变量：

- **$geoip_org** # 组织名称，比如 The University of Melbourne
