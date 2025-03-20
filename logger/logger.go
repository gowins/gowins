package logger

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupAccLogger(r *gin.Engine) {
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf(`{"client_ip": "%s", "timestamp": "%s", "method": "%s", "path": "%s", "latency":"%s",
proto": "%s", "status_code": %d, "referer": "%s", "user_agent": "%s"}`+"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006-01-02 15:04:05.00000"),
			param.Method,
			param.Path,
			param.Latency,
			param.Request.Proto,
			param.StatusCode,
			param.Request.Referer(),
			param.Request.UserAgent(),
		)
	}))
}
