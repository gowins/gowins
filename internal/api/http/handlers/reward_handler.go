package handlers

import (
	"net/http"

	rewardapp "gowins/internal/application/reward"
	"gowins/internal/application/reward/dto"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	service *rewardapp.RewardAppService
}

func NewRewardHandler(s *rewardapp.RewardAppService) *RewardHandler {
	return &RewardHandler{service: s}
}

func (rc *RewardHandler) CreateReward(ctx *gin.Context) {
	var req dto.CreateRewardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rc.service.CreateReward(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
