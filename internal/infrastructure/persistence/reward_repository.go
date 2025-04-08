package persistence

import (
	"errors"
	"sync"

	"gowins/internal/domain/reward"
)

// 模拟实现，可替换为 GORM/MySQL 实现

type inMemoryRewardRepo struct {
	store map[string]*reward.Reward
	mu    sync.RWMutex
}

func NewRewardRepository() reward.Repository {
	return &inMemoryRewardRepo{store: make(map[string]*reward.Reward)}
}

func (r *inMemoryRewardRepo) Save(rew *reward.Reward) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[rew.ID] = rew
	return nil
}

func (r *inMemoryRewardRepo) FindByID(id string) (*reward.Reward, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if val, ok := r.store[id]; ok {
		return val, nil
	}
	return nil, errors.New("not found")
}
