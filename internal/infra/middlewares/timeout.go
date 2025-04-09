package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware 使用 http.TimeoutHandler 设置请求超时
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		timeoutCtx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(timeoutCtx)
		c.Next()

		// 超时后，Gin 仍然可能继续执行 Handler，需要手动检查 context
		if timeoutCtx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
			c.Abort()
		}
	}
}
