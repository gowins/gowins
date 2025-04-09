package persistence

import (
	"errors"

	"gowins/internal/domain/reward"
)

// 模拟实现

type inMemoryRewardRepo struct {
	store map[string]*reward.Reward
}

func NewRewardRepository() reward.Repository {
	return &inMemoryRewardRepo{store: make(map[string]*reward.Reward)}
}

func (r *inMemoryRewardRepo) Save(rew *reward.Reward) error {
	r.store[rew.ID] = rew
	return nil
}

func (r *inMemoryRewardRepo) FindByID(id string) (*reward.Reward, error) {
	if val, ok := r.store[id]; ok {
		return val, nil
	}
	return nil, errors.New("not found")
}
