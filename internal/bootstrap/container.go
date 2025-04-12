package bootstrap

import (
	"gowins/internal/api/http/handlers"
	rewardapp "gowins/internal/application/reward"
	"gowins/internal/domain/reward"
	repo "gowins/internal/infra/persistence"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// Provide repository
	container.Provide(repo.NewRewardRepository, dig.As(new(reward.Repository)), dig.Name("inMemoryRewardRepo"))

	// Provide application service
	container.Provide(rewardapp.NewRewardAppService)

	// Provide HTTP handler NewRewardHandler

	container.Provide(handlers.NewRewardHandler)

	return container
}
