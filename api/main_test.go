package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// 设置 gin 运行模式为 TestMode
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
