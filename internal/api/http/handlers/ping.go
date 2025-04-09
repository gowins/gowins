package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {

	if c.Request.Context().Err() != nil { // 再次检查，防止 Sleep 期间超时
		return
	}
	c.String(http.StatusOK, "pong")
}
