package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zzhtl/go-mountain/internal/pkg/response"
	"github.com/zzhtl/go-mountain/internal/service"
)

// OperationLogHandler 操作日志处理器
type OperationLogHandler struct {
	svc *service.OperationLogService
}

// NewOperationLogHandler 创建操作日志处理器
func NewOperationLogHandler(svc *service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{svc: svc}
}

// List 获取操作日志列表
func (h *OperationLogHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := c.Query("username")
	module := c.Query("module")
	action := c.Query("action")

	list, total, err := h.svc.List(c.Request.Context(), page, pageSize, username, module, action)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.PageOK(c, list, total, page, pageSize)
}
