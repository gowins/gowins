package controllers

import (
	"net/http"

	rewardapp "gowins/internal/application/reward"
	"gowins/internal/application/reward/dto"

	"github.com/gin-gonic/gin"
)

type RewardController struct {
	service *rewardapp.RewardAppService
}

func NewRewardController(s *rewardapp.RewardAppService) *RewardController {
	return &RewardController{service: s}
}

func (rc *RewardController) CreateReward(ctx *gin.Context) {
	var req dto.CreateRewardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := rc.service.CreateReward(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}
