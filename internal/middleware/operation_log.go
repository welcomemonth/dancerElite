package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/service"
	"github.com/zzhtl/go-mountain/internal/store"
)

// OperationLogger 记录后台写操作日志
func OperationLogger(st *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			c.Next()
			return
		}

		c.Next()

		if c.Writer.Status() >= 400 {
			return
		}

		userID := int64(c.GetFloat64("user_id"))
		username, _ := c.Get("username")
		detail, _ := json.Marshal(gin.H{
			"method": method,
			"path":   c.FullPath(),
			"status": c.Writer.Status(),
		})

		log := &model.OperationLog{
			UserID:     userID,
			Username:   stringValue(username),
			Module:     service.ModuleFromPath(c.FullPath()),
			Action:     service.ActionFromMethod(method),
			TargetType: service.ModuleFromPath(c.FullPath()),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Detail:     detail,
		}
		_ = st.OperationLogRepo.Create(c.Request.Context(), log)
	}
}

func stringValue(value any) string {
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}
