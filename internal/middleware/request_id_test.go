package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRequestID_GenerateWhenMissing(t *testing.T) {
	// 准备
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 执行：不带 X-Request-ID 头
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证
	responseID := w.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("应该自动生成 request ID")
	}

	// 验证是合法的 UUID
	if _, err := uuid.Parse(responseID); err != nil {
		t.Errorf("生成的 request ID 应该是合法的 UUID: %v", err)
	}
}

func TestRequestID_PreserveExisting(t *testing.T) {
	// 准备
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	// 执行：带自定义 X-Request-ID 头
	customID := "my-custom-request-id-12345"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", customID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证：保持原值
	responseID := w.Header().Get("X-Request-ID")
	if responseID != customID {
		t.Errorf("应该保持原有的 request ID，期望 %s，实际 %s", customID, responseID)
	}
}

func TestRequestID_SetInContext(t *testing.T) {
	// 准备
	r := gin.New()
	r.Use(RequestID())

	var capturedID string
	r.GET("/test", func(c *gin.Context) {
		// 从 gin.Context 中获取 request_id
		if id, exists := c.Get("request_id"); exists {
			capturedID = id.(string)
		}
		c.String(200, "ok")
	})

	// 执行
	customID := "context-test-id"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", customID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证：可以通过 context 获取
	if capturedID != customID {
		t.Errorf("request_id 应该被设置到 gin.Context 中，期望 %s，实际 %s", customID, capturedID)
	}
}
