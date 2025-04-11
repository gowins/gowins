以下是基于上述时序图的完整代码实现，包含DDD分层架构中的所有关键组件：

---

### 一、领域层实现
#### 1. 订单聚合根
```go
// internal/domain/order/aggregate.go
package order

import "errors"

var (
    ErrInvalidOrderStatus = errors.New("invalid order status for payment")
    ErrPaymentAmountMismatch = errors.New("payment amount does not match order total")
)

type Order struct {
    ID         string
    Status     OrderStatus
    Total      float64
    PaymentID  string
    Version    int // 用于乐观锁
}

func (o *Order) MarkAsPaid(paymentID string) error {
    if o.Status != StatusPending {
        return ErrInvalidOrderStatus
    }
    o.PaymentID = paymentID
    o.Status = StatusPaid
    return nil
}

func (o *Order) ValidatePaymentAmount(amount float64) error {
    if o.Total != amount {
        return ErrPaymentAmountMismatch
    }
    return nil
}
```

#### 2. 支付领域服务
```go
// internal/domain/payment/service.go
package payment

import (
    "context"
    "yourproject/internal/domain/order"
)

type FraudDetection interface {
    IsHighRisk(ctx context.Context, orderID string, amount float64) (bool, error)
}

type PaymentService struct {
    fraudChecker FraudDetection
    paymentRepo PaymentRepository
}

func (s *PaymentService) VerifyPayment(
    ctx context.Context, 
    order *order.Order, 
    paymentID string,
    amount float64,
) error {
    // 检查金额匹配
    if err := order.ValidatePaymentAmount(amount); err != nil {
        return err
    }

    // 反欺诈检查
    isHighRisk, err := s.fraudChecker.IsHighRisk(ctx, order.ID, amount)
    if err != nil || isHighRisk {
        return ErrFraudCheckFailed
    }

    // 检查支付单状态
    payment, err := s.paymentRepo.Get(ctx, paymentID)
    if err != nil || payment.Status != StatusCompleted {
        return ErrInvalidPaymentStatus
    }

    return nil
}
```

---

### 二、应用层实现
#### 1. 应用服务
```go
// internal/application/order_service.go
package application

import (
    "context"
    "yourproject/internal/domain/order"
    "yourproject/internal/domain/payment"
)

type OrderService struct {
    orderRepo   order.Repository
    paymentSvc  *payment.PaymentService
    paymentRepo payment.Repository
    eventBus    EventBus
}

func (s *OrderService) Checkout(ctx context.Context, orderID string, paymentID string) error {
    // 1. 获取订单
    order, err := s.orderRepo.Get(ctx, orderID)
    if err != nil {
        return err
    }

    // 2. 获取支付单
    pmt, err := s.paymentRepo.Get(ctx, paymentID)
    if err != nil {
        return err
    }

    // 3. 支付验证（领域服务）
    if err := s.paymentSvc.VerifyPayment(ctx, order, paymentID, pmt.Amount); err != nil {
        return err
    }

    // 4. 修改订单状态（聚合根方法）
    if err := order.MarkAsPaid(paymentID); err != nil {
        return err
    }

    // 5. 持久化
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return err
    }

    // 6. 发布领域事件
    s.eventBus.Publish(ctx, order.Events())

    return nil
}
```

---

### 三、基础设施层实现
#### 1. 订单仓储
```go
// internal/infrastructure/persistence/order_repo.go
package persistence

import (
    "context"
    "gorm.io/gorm"
    "yourproject/internal/domain/order"
)

type OrderRepo struct {
    db *gorm.DB
}

func (r *OrderRepo) Get(ctx context.Context, id string) (*order.Order, error) {
    var o order.Order
    if err := r.db.WithContext(ctx).First(&o, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &o, nil
}

func (r *OrderRepo) Save(ctx context.Context, o *order.Order) error {
    return r.db.WithContext(ctx).Model(o).Where("version = ?", o.Version).
        Updates(map[string]interface{}{
            "status":     o.Status,
            "payment_id": o.PaymentID,
            "version":    gorm.Expr("version + 1"),
        }).Error
}
```

#### 2. 反欺诈实现
```go
// internal/infrastructure/antifraud/client.go
package antifraud

import (
    "context"
    "yourproject/internal/domain/payment"
)

type RiskClient struct {
    endpoint string
}

func (c *RiskClient) IsHighRisk(ctx context.Context, orderID string, amount float64) (bool, error) {
    // 调用第三方反欺诈服务API
    // ...
    return false, nil
}
```

---

### 四、依赖组装
```go
// internal/container/container.go
package container

import (
    "go.uber.org/dig"
    "yourproject/internal/domain/payment"
    "yourproject/internal/infrastructure/antifraud"
    "yourproject/internal/infrastructure/persistence"
)

func Build() *dig.Container {
    c := dig.New()

    // 基础设施层
    c.Provide(persistence.NewOrderRepo)
    c.Provide(persistence.NewPaymentRepo)
    c.Provide(antifraud.NewRiskClient)

    // 领域层
    c.Provide(func(client *antifraud.RiskClient) payment.FraudDetection {
        return client
    })
    c.Provide(payment.NewPaymentService)

    // 应用层
    c.Provide(application.NewOrderService)

    return c
}
```

---

### 五、启动流程
```go
// cmd/main.go
package main

import (
    "context"
    "yourproject/internal/container"
    "yourproject/internal/infrastructure/db"
)

func main() {
    // 初始化容器
    c := container.Build()

    // 初始化数据库等基础设施
    dbConn := db.NewConnection()

    // 运行服务
    c.Invoke(func(svc *application.OrderService) {
        ctx := context.Background()
        if err := svc.Checkout(ctx, "order123", "pay789"); err != nil {
            panic(err)
        }
    })
}
```

---

### 关键点说明：
1. **领域层**：
    - 纯业务逻辑，无技术细节
    - 聚合根负责状态变更和基本校验
    - 领域服务处理跨聚合逻辑

2. **应用层**：
    - 协调领域对象和基础设施
    - 处理事务边界
    - 事件发布

3. **基础设施层**：
    - 实现领域定义的接口
    - 包含所有技术细节（数据库、外部API等）

4. **依赖方向**：
   ```mermaid
   graph TD
       App[应用层] --> Domain[领域层]
       Domain --> Infra[基础设施层]
   ```

这个实现严格遵循：
- 聚合根负责自身完整性
- 领域服务处理复杂业务规则
- 应用服务编排流程
- 明确的分层架构
