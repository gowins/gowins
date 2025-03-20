package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	time.Sleep(5 * time.Second)
	c.String(http.StatusOK, "pong")
}
