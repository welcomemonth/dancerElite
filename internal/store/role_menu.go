package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
)

// RoleMenuRepo 角色-菜单关联数据访问接口（关联表，无通用 CRUD）
type RoleMenuRepo interface {
	// ListMenuIDsByRole 角色的菜单 ID 列表
	ListMenuIDsByRole(ctx context.Context, roleID int64) ([]int64, error)
	// ReplaceByRole 替换角色的菜单权限（事务：先删后建）
	ReplaceByRole(ctx context.Context, roleID int64, menuIDs []int64) error
	// DeleteByRole 删除角色的所有菜单关联
	DeleteByRole(ctx context.Context, roleID int64) error
	// DeleteByMenu 删除菜单的所有角色关联
	DeleteByMenu(ctx context.Context, menuID int64) error
}

// RoleMenuRepository RoleMenuRepo 的默认实现
type RoleMenuRepository struct {
	db *gorm.DB
}

func NewRoleMenuRepository(db *gorm.DB) RoleMenuRepo {
	return &RoleMenuRepository{db: db}
}

func (r *RoleMenuRepository) ListMenuIDsByRole(ctx context.Context, roleID int64) ([]int64, error) {
	var menuIDs []int64
	err := r.db.WithContext(ctx).Model(&model.RoleMenu{}).
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	return menuIDs, err
}

func (r *RoleMenuRepository) ReplaceByRole(ctx context.Context, roleID int64, menuIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return err
		}
		for _, menuID := range menuIDs {
			if err := tx.Create(&model.RoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleMenuRepository) DeleteByRole(ctx context.Context, roleID int64) error {
	return r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error
}

func (r *RoleMenuRepository) DeleteByMenu(ctx context.Context, menuID int64) error {
	return r.db.WithContext(ctx).Where("menu_id = ?", menuID).Delete(&model.RoleMenu{}).Error
}
