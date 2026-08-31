package store

import (
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// OperationLogRepo 操作日志数据访问接口（暂无自定义查询，仅通用 CRUD）
type OperationLogRepo interface {
	CRUD[model.OperationLog]
}

// NewOperationLogRepository 创建 OperationLogRepo，纯 CRUD 直接复用 BaseRepo
func NewOperationLogRepository(db *gorm.DB) OperationLogRepo {
	return repository.NewBaseRepo[model.OperationLog](db)
}
