package http

import (
	"gowins/internal/api/http/handlers"
	rewardapp "gowins/internal/application/reward"
	repo "gowins/internal/infra/persistence"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	rewardRepo := repo.NewRewardRepository()
	rewardService := rewardapp.NewRewardAppService(rewardRepo)
	rewardController := handlers.NewRewardhandler(rewardService)

	rewardGroup := router.Group("/rewards")

	{
		rewardGroup.POST("/add", rewardController.CreateReward)
	}
}
