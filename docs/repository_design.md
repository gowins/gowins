在您的 DDD 项目中集成 GORM 并实现 `reward_repository`，我会提供一套完整的解决方案，保持清晰的层次结构并符合 DDD 原则。以下是详细实现：

---

## 1. 数据库初始化 (infra/db/gorm.go)

```go
package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitGORM(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 可根据环境调整日志级别
		SkipDefaultTransaction: true, // 提高性能
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 连接池配置
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	
	return DB, nil
}
```

---

## 2. 领域层仓储接口 (domain/reward/repository.go)

```go
package reward

import (
	"context"
	"github.com/yourproject/internal/domain"
)

// RewardRepository 仓储接口（领域层定义）
type RewardRepository interface {
	FindByID(ctx context.Context, id uint) (*domain.Reward, error)
	Save(ctx context.Context, reward *domain.Reward) error
	FindByCriteria(ctx context.Context, criteria RewardCriteria) ([]*domain.Reward, error)
	Delete(ctx context.Context, id uint) error
	
	// 自定义方法
	FindActiveRewards(ctx context.Context) ([]*domain.Reward, error)
}

// RewardCriteria 查询条件（领域层定义）
type RewardCriteria struct {
	MinAmount    *float64
	Status       *domain.RewardStatus
	CreatorID    *uint
	CreatedAfter *time.Time
}
```

---

## 3. GORM 实现 (infra/persistence/reward_repository.go)

```go
package persistence

import (
	"context"
	"errors"
	"github.com/yourproject/internal/domain/reward"
	"github.com/yourproject/internal/infra/db"
	"gorm.io/gorm"
)

type GormRewardRepository struct {
	db *gorm.DB
}

func NewGormRewardRepository(db *gorm.DB) reward.RewardRepository {
	return &GormRewardRepository{db: db}
}

// 实现领域层接口
func (r *GormRewardRepository) FindByID(ctx context.Context, id uint) (*domain.Reward, error) {
	var reward domain.Reward
	if err := r.db.WithContext(ctx).First(&reward, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRewardNotFound
		}
		return nil, err
	}
	return &reward, nil
}

func (r *GormRewardRepository) Save(ctx context.Context, reward *domain.Reward) error {
	return r.db.WithContext(ctx).Save(reward).Error
}

func (r *GormRewardRepository) FindByCriteria(ctx context.Context, criteria reward.RewardCriteria) ([]*domain.Reward, error) {
	var rewards []*domain.Reward
	
	query := r.db.WithContext(ctx).Model(&domain.Reward{})
	
	if criteria.MinAmount != nil {
		query = query.Where("amount >= ?", *criteria.MinAmount)
	}
	
	if criteria.Status != nil {
		query = query.Where("status = ?", *criteria.Status)
	}
	
	// 其他条件...
	
	if err := query.Find(&rewards).Error; err != nil {
		return nil, err
	}
	
	return rewards, nil
}

// 自定义查询实现
func (r *GormRewardRepository) FindActiveRewards(ctx context.Context) ([]*domain.Reward, error) {
	var rewards []*domain.Reward
	err := r.db.WithContext(ctx).
		Where("status = ? AND expiration_date > ?", 
			domain.RewardStatusActive, 
			time.Now()).
		Find(&rewards).Error
		
	return rewards, err
}
```

---

## 4. 依赖注入 (在应用初始化时)

```go
// cmd/main.go 或 infra/init.go
package main

import (
	"github.com/yourproject/internal/infra/db"
	"github.com/yourproject/internal/infra/persistence"
	"github.com/yourproject/internal/application/reward"
)

func main() {
	// 初始化数据库
	gormDB, err := db.InitGORM(db.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "user",
		Password: "password",
		DBName:   "reward_db",
	})
	if err != nil {
		panic(err)
	}
	
	// 初始化仓储
	rewardRepo := persistence.NewGormRewardRepository(gormDB)
	
	// 初始化应用服务
	rewardService := reward.NewRewardAppService(rewardRepo)
	
	// 继续初始化其他组件...
}
```

---

## 5. 事务处理 (领域服务示例)

```go
// domain/reward/service.go
package reward

import (
	"context"
	"gorm.io/gorm"
)

type RewardService struct {
	repo     RewardRepository
	db       *gorm.DB // 用于事务管理
}

func NewRewardService(repo RewardRepository, db *gorm.DB) *RewardService {
	return &RewardService{repo: repo, db: db}
}

func (s *RewardService) CompleteReward(ctx context.Context, rewardID uint) error {
	// 开启事务
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 使用事务版的Repository
		txRepo := NewGormRewardRepository(tx)
		
		reward, err := txRepo.FindByID(ctx, rewardID)
		if err != nil {
			return err
		}
		
		if err := reward.Complete(); err != nil {
			return err
		}
		
		return txRepo.Save(ctx, reward)
	})
}
```

---

## 6. 高级查询优化

```go
// 使用Scopes优化复杂查询
func (r *GormRewardRepository) FindComplex(ctx context.Context, filter RewardFilter) ([]*domain.Reward, error) {
	var rewards []*domain.Reward
	
	err := r.db.WithContext(ctx).
		Scopes(
			AmountGreaterThan(filter.MinAmount),
			StatusIn(filter.Statuses),
			CreatedBetween(filter.StartDate, filter.EndDate),
		).
		Preload("Creator").
		Find(&rewards).Error
		
	return rewards, err
}

// 查询Scope示例
func AmountGreaterThan(amount *float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if amount != nil {
			return db.Where("amount > ?", *amount)
		}
		return db
	}
}
```

---

## 关键设计要点

1. **分层清晰**：
    - 领域层只定义接口
    - 基础设施层实现具体持久化

2. **依赖方向**：
   ```mermaid
   graph TD
     A[domain] -->|定义| B[repository接口]
     C[infra] -->|实现| B
     D[application] -->|依赖| B
   ```

3. **事务管理**：
    - 在应用层或领域服务中管理事务
    - 通过传递 `*gorm.DB` 控制事务边界

4. **测试友好**：
   ```go
   // 测试时可以使用内存实现
   type MockRewardRepository struct {
       reward.RewardRepository
   }
   ```

5. **性能优化**：
    - 合理使用 `Preload` 和 `Select`
    - 避免 N+1 查询问题

这套实现严格遵循DDD原则，同时充分利用GORM的特性，保持了代码的整洁性和可维护性。
