package service

import (
	"context"
	"strings"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// OperationLogService 操作日志服务
type OperationLogService struct {
	repo *repository.BaseRepo[model.OperationLog]
	db   *gorm.DB
}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService(db *gorm.DB) *OperationLogService {
	return &OperationLogService{
		repo: repository.NewBaseRepo[model.OperationLog](db),
		db:   db,
	}
}

// List 分页查询操作日志
func (s *OperationLogService) List(ctx context.Context, page, pageSize int, username, module, action string) ([]model.OperationLog, int64, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		if username != "" {
			db = db.Where("username LIKE ?", "%"+username+"%")
		}
		if module != "" {
			db = db.Where("module = ?", module)
		}
		if action != "" {
			db = db.Where("action = ?", action)
		}
		return db.Order("created_at DESC, id DESC")
	}
	return s.repo.List(ctx, page, pageSize, scope)
}

// Create 写入操作日志
func (s *OperationLogService) Create(ctx context.Context, log *model.OperationLog) error {
	return s.db.WithContext(ctx).Create(log).Error
}

// ModuleFromPath 从后台 API 路径推导模块名
func ModuleFromPath(path string) string {
	path = strings.TrimPrefix(path, "/api/admin/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		return "admin"
	}
	return parts[0]
}

// ActionFromMethod 从 HTTP 方法推导操作类型
func ActionFromMethod(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return strings.ToLower(method)
	}
}
