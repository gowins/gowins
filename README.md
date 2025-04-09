业务整理的核心组织目录：

```
/ecommerce
├── api/                          # API协议定义
│   ├── openapi/                  # OpenAPI规范
│   └── protobuf/                 # gRPC Proto文件
├── assets/                       # 图片，入职等资源文件
├── build/                        # 容器（Docker）包配置和脚本放在 /build/package 目录中。
├── cmd/                          # 应用程序入口
│   └── api/                      # REST API服务
│       └── main.go               # 启动HTTP服务
├── init/                         # 系统初始化（systemd、upstart、sysv）和进程管理器/主管（runit、supervisord）配置。
├── internal/                     # 核心业务代码（DDD实现）
│   ├── domain/                   # 领域层
│   │   ├── reward/               # 悬赏聚合
│   │   │   ├── aggregate.go      # 聚合根实体（Order）
│   │   │   ├── reward_item.go    # 子实体（OrderItem）
│   │   │   ├── repository.go     # 仓储接口
│   │   │   ├── events.go         # 领域事件
│   │   │   └── service.go        # 领域服务（如OrderValidationService）
│   │   │
│   │   ├── product/              # 产品聚合（类似结构）
│   │   └── shared/               # 跨聚合共享定义
│   │       ├── money.go          # 值对象（Money）
│   │       └── address.go        # 值对象（Address）
│   │
│   ├── application/              # 应用层
│   │   └── reward/               # 订单应用服务
│   │       ├── service.go        # OrderAppService
│   │       └── dto/              # 数据传输对象
│   │           ├── request.go    # CreateOrderRequest
│   │           └── response.go   # OrderResponse
│   │
│   ├── api/                      # 接口层
│   │   └── http/                 # HTTP接口
│   │       ├── handlers/         # 控制器
│   │       │   └── reward_handler.go
│   │       └── router.go         # 路由定义
│   │
│   └── infra/                    # 基础设施层
│       ├── persistence/          # 持久化实现
│       │   └── order_repository.go
│       ├── db/                   # 数据库连接
│       └── logging/              # 日志等工具
│
├── pkg/                          # 公共库代码
│   ├── util/                     # 通用工具
│   └── errors/                   # 错误处理
│
├── configs/                      # 配置文件
├── deployments/                  # 部署配置
├── scripts/                      # 脚本
├── go.mod
└── go.sum

```

### 实施方法：
1. 依赖规则可视化

    ```接口层(interfaces) → 应用层(app) → 领域层(domain) ← 基础设施层(infrastructure)```
