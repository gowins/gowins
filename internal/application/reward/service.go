package rewardapp

import (
	"gowins/internal/application/reward/dto"
	"gowins/internal/domain/reward"

	"github.com/google/uuid"
)

type RewardAppService struct {
	repo reward.Repository
}

func NewRewardAppService(repo reward.Repository) *RewardAppService {
	return &RewardAppService{repo: repo}
}

func (s *RewardAppService) CreateReward(req dto.CreateRewardRequest) (string, error) {
	id := uuid.NewString()
	r := &reward.Reward{
		Type:        req.Type,
		Device:      req.Device,
		Project:     req.Project,
		Title:       req.Title,
		Description: req.Description,
		Steps:       req.Steps,
		Deadline:    req.Deadline,
		ReviewTime:  req.ReviewTime,
		UnitPrice:   req.UnitPrice,
		Quantity:    req.Quantity,
		SingleUse:   req.SingleUse,
	}
	newReward := reward.NewReward(id, r)
	err := s.repo.Save(newReward)
	return newReward.ID, err
}
