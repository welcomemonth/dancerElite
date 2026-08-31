package store

import (
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// ColumnRepo 栏目数据访问接口（暂无自定义查询，仅通用 CRUD）
type ColumnRepo interface {
	CRUD[model.Column]
}

// NewColumnRepository 创建 ColumnRepo，纯 CRUD 直接复用 BaseRepo
func NewColumnRepository(db *gorm.DB) ColumnRepo {
	return repository.NewBaseRepo[model.Column](db)
}
