package store

import (
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// MenuRepo 菜单数据访问接口（暂无自定义查询，仅通用 CRUD）
type MenuRepo interface {
	CRUD[model.Menu]
}

// NewMenuRepository 创建 MenuRepo，纯 CRUD 直接复用 BaseRepo
func NewMenuRepository(db *gorm.DB) MenuRepo {
	return repository.NewBaseRepo[model.Menu](db)
}
