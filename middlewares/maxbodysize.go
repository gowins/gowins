package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupMaxBodySizeMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 限制请求体大小为 maxBytes
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
