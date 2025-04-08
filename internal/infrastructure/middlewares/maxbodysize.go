package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupMaxBodySizeMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 用 http.MaxBytesReader 限制 Body 读取大小
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		// 判断c.
		if c.Request.ContentLength > maxBytes {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			c.Abort()
			return
		}
		// 继续后续处理
		c.Next()
	}
}
