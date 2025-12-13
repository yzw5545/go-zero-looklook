# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

go-zero-looklook 是一个基于 go-zero 框架构建的完整微服务示例项目，包含用户中心、订单、支付、旅游/民宿、消息队列五个微服务。

## 常用命令

### 环境启动

```bash
# 启动中间件环境（MySQL、Redis、Kafka、ES、Jaeger等）
docker-compose -f docker-compose-env.yml up -d

# 启动应用服务
docker-compose up -d

# 热加载开发
modd
```

### 代码生成

```bash
# 生成API代码（在对应服务的api目录下执行）
goctl api go -api desc/{service}.api -dir=.

# 生成RPC代码（在对应服务的rpc目录下执行）
goctl rpc protoc pb/{service}.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 生成Model代码
goctl model mysql ddl -src={table}.sql -dir=. -c
```

### 编译服务

```bash
# 编译单个服务（示例）
go build -o data/server/usercenter-api -v app/usercenter/cmd/api/usercenter.go
go build -o data/server/usercenter-rpc -v app/usercenter/cmd/rpc/usercenter.go
```

## 架构说明

### 服务结构

每个业务服务采用 API + RPC 分层架构：

```
app/{service}/
├── cmd/
│   ├── api/           # HTTP API服务
│   │   ├── desc/      # .api定义文件
│   │   ├── etc/       # 配置文件（yaml）
│   │   └── internal/
│   │       ├── handler/   # HTTP处理器
│   │       ├── logic/     # 业务逻辑
│   │       ├── svc/       # 服务上下文（依赖注入）
│   │       └── types/     # 请求/响应结构
│   ├── rpc/           # gRPC服务
│   │   ├── pb/        # .proto定义文件
│   │   ├── etc/       # 配置文件
│   │   └── internal/
│   │       ├── logic/     # RPC业务逻辑
│   │       ├── server/    # RPC服务实现
│   │       └── svc/       # 服务上下文
│   └── mq/            # 消息队列消费者（部分服务）
└── model/             # 数据库Model（goctl生成）
```

### 服务端口规划

| 服务 | API端口 | RPC端口 | MQ端口 |
|------|---------|---------|--------|
| order | 1001 | 2001 | 3001 |
| payment | 1002 | 2002 | - |
| travel | 1003 | 2003 | - |
| usercenter | 1004 | 2004 | - |
| mqueue | - | - | scheduler:3003, job:3002 |

### 服务间通信

- **API → RPC**: HTTP请求进入API层，API调用RPC服务处理业务
- **RPC → RPC**: 服务间通过gRPC直接调用
- **异步通信**: 通过Kafka消息队列（如支付状态更新通知订单服务）

### 公共包 (pkg/)

- `xerr/`: 统一错误码和错误处理
- `result/`: HTTP响应统一格式
- `ctxdata/`: 上下文数据（用户ID、请求ID等）
- `middleware/`: HTTP中间件（JWT认证）
- `interceptor/`: RPC拦截器
- `uniqueid/`: 分布式ID生成（sonyflake）
- `tool/`: 工具函数（加密、随机数等）

## 配置文件

配置文件位于各服务的 `etc/` 目录，YAML格式，主要配置项：

```yaml
Name: service-name
Host: 0.0.0.0
Port: xxxx
Mode: dev              # dev/test/pro

JwtAuth:
  AccessSecret: xxx
  AccessExpire: xxx

DB:
  DataSource: xxx

Cache:
  - Host: redis:6379

Telemetry:             # Jaeger链路追踪
  Name: service-name
  Endpoint: http://jaeger:14268/api/traces

Prometheus:            # 监控指标
  Host: 0.0.0.0
  Port: xxxx
```

## 数据库

各服务独立数据库：
- `looklook_usercenter`: 用户相关表
- `looklook_order`: 订单相关表
- `looklook_payment`: 支付相关表
- `looklook_travel`: 民宿相关表

SQL初始化脚本: `deploy/sql/looklook_*.sql`

## 网关

Nginx网关配置: `deploy/nginx/conf.d/looklook-gateway.conf`

访问路径:
- `http://localhost:8888/usercenter/*` → 1004
- `http://localhost:8888/order/*` → 1001
- `http://localhost:8888/payment/*` → 1002
- `http://localhost:8888/travel/*` → 1003

## 监控与日志

- **Prometheus**: 9090, 各服务指标端口 400x
- **Grafana**: 3001
- **Jaeger**: 16686（UI）
- **Kibana**: 5601（日志查询）
- **asynqmon**: 8980（任务队列管理）
