---
title: geoip2 模块
---

# 概述

> 参考：
>
> - [GitHub 项目，leev/ngx_http_geoip2_module](https://github.com/leev/ngx_http_geoip2_module)

由于隐私的原因，[MaxMind 在 2019 年 12 月份对数据库进行重大变更](https://blog.maxmind.com/2019/12/18/significant-changes-to-accessing-and-using-geolite2-databases/)，所以，老式的 geo/geoip 模块不再适用于新的 MaxMind 数据库，所以，geoip2 模块诞生了。

从[这个页面可](https://www.maxmind.com/en/accounts/545756/geoip/downloads)以下载 GeoIP2 和 旧版的 GeoIP 数据库

geoip2 模块与 geo/geoip 模块的功能类似。geoip2 模块根据 客户端的 IP 信息，使用 MaxMind 的 geoip2 数据库中的值创建变量。只不过指令用法稍有不同。

## 用法示例

### 首先通过 [mmdblookup](https://maxmind.github.io/libmaxminddb/mmdblookup.html) 工具查看数据库中的内容

```json
root@desistdaydream:~/test_dir/downloads# mmdblookup --file ./GeoLite2-City.mmdb --ip 59.46.138.226

  {
    "city":
      {
        "geoname_id":
          1814087 <uint32>
        "names":
          {
            "en":
              "Dalian" <utf8_string>
            "ja":
              "大連市" <utf8_string>
            "ru":
              "Далянь" <utf8_string>
            "zh-CN":
              "大连" <utf8_string>
          }
      }
    "continent":
      {
        "code":
          "AS" <utf8_string>
        "geoname_id":
          6255147 <uint32>
        "names":
          {
            "de":
              "Asien" <utf8_string>
            "en":
              "Asia" <utf8_string>
            "es":
              "Asia" <utf8_string>
            "fr":
              "Asie" <utf8_string>
            "ja":
              "アジア" <utf8_string>
            "pt-BR":
              "Ásia" <utf8_string>
            "ru":
              "Азия" <utf8_string>
            "zh-CN":
              "亚洲" <utf8_string>
          }
      }
    "country":
      {
        "geoname_id":
          1814991 <uint32>
        "iso_code":
          "CN" <utf8_string>
        "names":
          {
            "de":
              "China" <utf8_string>
            "en":
              "China" <utf8_string>
            "es":
              "China" <utf8_string>
            "fr":
              "Chine" <utf8_string>
            "ja":
              "中国" <utf8_string>
            "pt-BR":
              "China" <utf8_string>
            "ru":
              "Китай" <utf8_string>
            "zh-CN":
              "中国" <utf8_string>
          }
      }
    "location":
      {
        "accuracy_radius":
          1000 <uint16>
        "latitude":
          38.912200 <double>
        "longitude":
          121.602200 <double>
        "time_zone":
          "Asia/Shanghai" <utf8_string>
      }
    "registered_country":
      {
        "geoname_id":
          1814991 <uint32>
        "iso_code":
          "CN" <utf8_string>
        "names":
          {
            "de":
              "China" <utf8_string>
            "en":
              "China" <utf8_string>
            "es":
              "China" <utf8_string>
            "fr":
              "Chine" <utf8_string>
            "ja":
              "中国" <utf8_string>
            "pt-BR":
              "China" <utf8_string>
            "ru":
              "Китай" <utf8_string>
            "zh-CN":
              "中国" <utf8_string>
          }
      }
    "subdivisions":
      [
        {
          "geoname_id":
            2036115 <uint32>
          "iso_code":
            "LN" <utf8_string>
          "names":
            {
              "en":
                "Liaoning" <utf8_string>
              "fr":
                "Province de Liaoning" <utf8_string>
              "zh-CN":
                "辽宁" <utf8_string>
            }
        }
      ]
  }

```

### 根据数据库内容，定义变量

```nginx
http {
    ...
    # 指定数据库路径
    geoip2 /etc/maxmind-country.mmdb {
        auto_reload 5m;
        # 根据数据库内容生成变量
        $geoip2_metadata_country_build metadata build_epoch;
        $geoip2_data_country_code default=US source=$variable_with_ip country iso_code;
        $geoip2_data_country_name country names en;
    }

   # 指定数据库路径
    geoip2 /etc/maxmind-city.mmdb {
        # 根据数据库内容生成变量
        $geoip2_data_city_name default=London city names en;
    }
    ....

    fastcgi_param COUNTRY_CODE $geoip2_data_country_code;
    fastcgi_param COUNTRY_NAME $geoip2_data_country_name;
    fastcgi_param CITY_NAME    $geoip2_data_city_name;
    ....
}

stream {
    ...
    # 指定数据库路径
    geoip2 /etc/maxmind-country.mmdb {
        # 根据数据库内容生成变量
        $geoip2_data_country_code default=US source=$remote_addr country iso_code;
    }
    ...
}
```

# 部署并启用模块

首先按照其[README.md 文件](https://github.com/maxmind/libmaxminddb/blob/master/README.md#installing-from-a-tarball)中的[说明](https://github.com/maxmind/libmaxminddb/blob/master/README.md#installing-from-a-tarball)安装[libmaxminddb](https://github.com/maxmind/libmaxminddb)。

**下载 nginx 源**

    wget http://nginx.org/download/nginx-VERSION.tar.gz
    tar zxvf nginx-VERSION.tar.gz
    cd nginx-VERSION

**要构建为动态模块（nginx 1.9.11+）：**

    ./configure --add-dynamic-module=/path/to/ngx_http_geoip2_module
    make
    make install

这将产生 `objs/ngx_http_geoip2_module.so`。可以将其手动复制到 nginx 的模块存储路径(比如 /etc/nginx/modules 路径下)。
将以下行添加到您的 nginx.conf 中的 main 配置环境中：

```nginx
load_module modules/ngx_http_geoip2_module.so;
```

然后，就可以在配置文件中使用 geoip2 指令配置 geoip2 模块了

# http 模块下的 geoip2 模块指令

## geoip2 FILE {} # 根据指定的数据库文件定义变量

该指令类似于 geo 模块的 geo 指令，可以自己定义变量名称

### $VariableName \[default=STRING] \[source=IP] PATH

定义名为 `$VariableName` 的变量

- **PATH** # MaxMind 数据库中的数据路径，将该路径下的值，赋值给变量 `$VariableName`
  - 注意：MaxMind 的 GeoIP2 数据库是 JSON 结构，所以 PATH 就是由以空格分割的字段名称组成。可以通过 [mmdblookup 工具](https://maxmind.github.io/libmaxminddb/mmdblookup.html)查找所需数据的路径
- **default=\<STRING>** # 若变量无法获取到值时，应该具有的默认值。
- **source=\<IP>** # 指定要从数据库获取信息的 IP 地址。默认值来自 `$remote_addr` 变量的值

#### EXAMPLE

- 创建 `$geoip2_data_country_code` 变量，根据 `$remote_addr` 变量中的 IP 地址，查找数据库，将 IP 对应的 `.country.iso_code` 字段的值赋值给 ` $``geoip2_data_country_code ` 变量，若 `.country.iso_code` 字段为空，则变量的值为 US。
  - **`$geoip2\_data\_country\_code default=US source=$remote_addr country iso_code;`**
  - 其实就是获取两个字母的国家代码
- 创建 `$geoip2_city_country_name` 变量，根据 `$remote_addr` 变量中的 IP 地址，查找数据库，将 IP 对应的 `.country.name.zh-CN` 字段的值赋值给 `$geoip2_city_country_name` 变量
  - **`$geoip2\_city\_country\_name source=$remote_addr country names zh-CN;`**
  - 其实就是中文显示的国家名称

# 配置示例

```nginx
http {
  geoip2 /etc/nginx/geoip/GeoLite2-City.mmdb {
    $geoip2_city_country_code source=$remote_addr country iso_code;
    $geoip2_city_country_name source=$remote_addr country names zh-CN;
    $geoip2_city source=$remote_addr city names zh-CN;
    $geoip2_postal_code source=$remote_addr postal code;
    $geoip2_dma_code source=$remote_addr location metro_code;
    $geoip2_latitude source=$remote_addr location latitude;
    $geoip2_longitude source=$remote_addr location longitude;
    $geoip2_time_zone source=$remote_addr location time_zone;
    $geoip2_region_code source=$remote_addr subdivisions 0 iso_code;
    $geoip2_region_name source=$remote_addr subdivisions 0 names zh-CN;
  }
}
```

上述示例中，定义了如下几个变量，并获取了数据库中指定字段的值：

| **变量**                      | **数据库中的字段**            | **含义** |
| ----------------------------- | ----------------------------- | -------- |
| **$geoip2_city_country_code** | .country.iso_code             | 国家代码 |
| **$geoip2_city_country_name** | .country.names.zh-CN          | 国家名称 |
| **$geoip2_city**              | .city.name.zh-CN              | 城市名称 |
| **$geoip2_postal_code**       | .postal.code                  | 邮编     |
| **$geoip2_dma_code**          | .location.metro_code          |          |
| **$geoip2_latitude**          | .location.latitude            | 纬度     |
| **$geoip2_longitude**         | .location.longitude           | 经度     |
| **$geoip2_time_zone**         | .location..time_zone          | 时区     |
| **$geoip2_region_code**       | .subdivisions\[0].iso_code    |          |
| **$geoip2_region_name**       | .subdivisions\[0].names.zh-CN |          |
