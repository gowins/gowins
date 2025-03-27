package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	time.Sleep(time.Second * 5)
	if c.Request.Context().Err() != nil { // 再次检查，防止 Sleep 期间超时
		return
	}
	c.String(http.StatusOK, "pong")
}
