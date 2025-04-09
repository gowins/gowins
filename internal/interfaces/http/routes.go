package http

import (
	rewardapp "gowins/internal/application/reward"
	repo "gowins/internal/infrastructure/persistence"
	"gowins/internal/interfaces/http/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	rewardRepo := repo.NewRewardRepository()
	rewardService := rewardapp.NewRewardAppService(rewardRepo)
	rewardController := controllers.NewRewardController(rewardService)

	rewardGroup := router.Group("/rewards")

	{
		rewardGroup.POST("/add", rewardController.CreateReward)
	}
}
