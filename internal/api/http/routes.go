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
	rewardHandler := handlers.NewRewardHandler(rewardService)

	rewardGroup := router.Group("/rewards")
	{
		rewardGroup.POST("/add", rewardHandler.CreateReward)
	}

	//// routes/user_routes.go
	//func registerUserRoutes(router *gin.RouterGroup) {
	//	router.GET("", h.ListUsers)
	//	router.POST("", h.CreateUser)
	//}
	//
	//// routes/routes.go
	//func (r *Router) RegisterAll(router *gin.Engine) {
	//	v1 := router.Group("/api/v1")
	//	r.registerUserRoutes(v1.Group("/users"))
	//}
}
