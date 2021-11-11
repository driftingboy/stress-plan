# stress-plan
stress-plan 是使用golang语言编写的轻量级压测工具，有着易用，从业务场景出发的压测功能。

## 功能特点
- 轻量，依托于协程模拟并发数
- 贴进实际使用，比如提供协程增长速率，多接口按权重压测..
- 支持多种协议：
  - [x] http
  - [ ] websocket
  - [ ] 自定义rpc
- 分析数据全

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

cd ./stress-plan

go run main.go
```