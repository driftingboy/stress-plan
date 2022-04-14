# stress-plan
stress-plan 是使用golang语言编写的轻量级压测工具，有着易用，从业务场景出发的压测功能。
用于后端开发人员来应付日常接口性能验证、分析。

## 功能特点
- 轻量，依托于协程模拟并发数
- 贴进实际使用，比如提供协程增长速率，多接口按权重压测..
- 支持多种协议：
  - [x] http
  - [ ] websocket
  - [ ] 自定义rpc
- 分析数据全
  - [ ] 请求qps、时长随并发数增长的图表

可能会实现的功能：
- 压测数据持久化
- 持久化数据查询

## 架构设计

### 模块划分
![20211103232532](https://i.loli.net/2021/11/03/A5ylKOQ8cwVJ9PY.png)

### 技术架构

![20211104000437](https://i.loli.net/2021/11/04/U4rwcJpZyjVSoKI.png)

## 快速开始

``` shell
git clone git@github.com:driftingboy/stress-plan.git

cd ./stress-plan/cmd

./stp run -c 10 -n 1000 -u https://www.baidu.com/

```

如果您需要使用post请求测试, 参考如下

- post 请求（默认 "Content-Type:application/json"）
```shell
./stp run -c 1 -n 1 -u "POST@http://119.3.106.151:10100/v1/app/evidences:new_verify" -b "{ \"brief_content\": \"da11370d639b4c36b7e38a25516aed21c60ab2ac90e4474f89644fdb1dad93c1\", \"verify_type\": \"key\"}"
```

- 指定 body 类型
```shell
./stp run -c 1 -n 1 \
-H "Content-Type:application/x-www-form-urlencoded" \
-u "POST@http://119.3.106.151:10100/v1/app/evidences:new_verify" \
-b "a=a&&b=b"


mock 数据
./stp run -c 20 -n 50 -u POST@"http://xxxx:1000/v1/app/evidences" -H "accept=application/json" -H "Authorization=Bearer eyJhbGciOiJIUzUxMiJ1.eyJleHAiOjE2NDAwOTUzNDYsImlhdCI6MTY0MDA4ODE0NiwianRpIjoianI1NnZxanluImtxN3AiLCJzdWIiOiJ1aWQtdGVuYW50In0.y-xNb2DDCi2cU1JQlO9HAxoU_AyjYha8I3wfcv5x9dBnDVLwgDSdzIYl9BlzHyww3fOIj1VImA-w26n2LMPATQ" -H "Content-Type= application/json" -b "{ \"tenant_id\": \"tid-yuhu1\", \"title\": \"stress-test-01\", \"content\": \"@Base64\", \"evidence_type\": \"text\"}"

cd c./stp run -c 20 -n 2000 -u POST@"http://console.yuhu.tech/api/v1/app/evidences" -H "accept=application/json" -H "Authorization=Bearer eyJhbGciOiJIUzUxMiJ9.eyJleHAiOjE2NDYzMDM3MDcsImlhdCI6MTY0NjI5NjUwNywianRpIjoiamx4cDE1NzU5eG8yNDIiLCJzdWIiOiJ1aWQtdGVuYW50In0.II11OB3CMuNh1lUQDV0A8EZr-Cj9KhB9L1r5dgDMmnXvU3SiqDZyWrCtRhRc9RgrNKGUU_oP0NBUtjoX1sHLRw" -H "Content-Type= application/json" -b "{ \"tenant_id\": \"tid-yuhu1\", \"title\": \"stress-test-20220303\", \"content\": \"@Base64\", \"evidence_type\": \"text\"}"
```

详细用法使用 -h 或 --help 查看
``` shell
./stp -h
./stp run -h
```

## 场景

## 结果分析

结果有两部分， 通过 -f 可以动态展示
- 统计部分
- 分析图表
