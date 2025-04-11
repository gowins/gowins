package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	rewardapp "gowins/internal/application/reward"
	"gowins/internal/application/reward/dto"
	repo "gowins/internal/infra/persistence"
	"gowins/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRewardHandler_CreateReward(t *testing.T) {
	r := gin.New()
	rewardRepo := repo.NewRewardRepository()
	rewardService := rewardapp.NewRewardAppService(rewardRepo)
	rewardHandler := NewRewardhandler(rewardService)
	r.POST("/rewards/create", rewardHandler.CreateReward)
	// 准备测试数据
	newItem := dto.CreateRewardRequest{
		Type:        "electronics",
		Device:      "smartphone",
		Project:     "mobile-app",
		Title:       "Test Device",
		Description: "This is a test device",
		Steps:       []string{"step1", "step2", "step3"},
		Deadline:    time.Now().Add(48 * time.Hour).Unix(),
		ReviewTime:  time.Now().Unix(),
		UnitPrice:   10,
		Quantity:    5,
		SingleUse:   false,
	}

	jsonValue, _ := json.Marshal(newItem)
	rr := util.PerformRequest(r, http.MethodPost, "/rewards/create", bytes.NewBuffer(jsonValue), nil)

	var response dto.RewardResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, response.TotalAmount, newItem.UnitPrice*float64(newItem.Quantity))
}
