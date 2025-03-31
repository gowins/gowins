package middlewares

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"gowins/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	// 创建一个管道用于重定向 os.Stderr
	r, w, _ := os.Pipe()

	// 备份原来的 os.Stderr
	oldStderr := os.Stderr
	defer func() {
		os.Stderr = oldStderr // 测试结束后恢复
	}()

	// 重定向 os.Stderr
	os.Stderr = w

	router := gin.New()
	logWriter := &JSONLogWriter{}
	router.Use(gin.RecoveryWithWriter(logWriter))
	router.GET("/recovery", func(_ *gin.Context) {
		panic("Oupps, Houston, we have a problem")
	})
	// RUN
	rr := utils.PerformRequest(router, http.MethodGet, "/recovery")

	_ = w.Close()

	// 读取 os.Stderr 的内容
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// TEST
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, buf.String(), "panic recovered")
	assert.Contains(t, buf.String(), "Oupps, Houston, we have a problem")
	assert.Contains(t, buf.String(), t.Name())
	assert.Contains(t, buf.String(), "GET /recovery")

}
