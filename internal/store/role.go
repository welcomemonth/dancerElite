package store

import (
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// RoleRepo 角色数据访问接口（暂无自定义查询，仅通用 CRUD）
type RoleRepo interface {
	CRUD[model.Role]
}

// NewRoleRepository 创建 RoleRepo，纯 CRUD 直接复用 BaseRepo
func NewRoleRepository(db *gorm.DB) RoleRepo {
	return repository.NewBaseRepo[model.Role](db)
}
