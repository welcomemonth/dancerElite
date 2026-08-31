package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/pkg/errcode"
	"github.com/zzhtl/go-mountain/internal/store"
)

// RoleService 角色管理服务
type RoleService struct {
	st *store.Store
}

// NewRoleService 创建角色服务
func NewRoleService(st *store.Store) *RoleService {
	return &RoleService{st: st}
}

// List 获取角色列表
func (s *RoleService) List(ctx context.Context, page, pageSize int) ([]model.Role, int64, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Order("id")
	}
	return s.st.RoleRepo.List(ctx, page, pageSize, scope)
}

// Get 获取单个角色
func (s *RoleService) Get(ctx context.Context, id int64) (*model.Role, error) {
	return s.st.RoleRepo.GetByID(ctx, id)
}

// Create 创建角色
func (s *RoleService) Create(ctx context.Context, role *model.Role) error {
	role.Status = 1
	return s.st.RoleRepo.Create(ctx, role)
}

// Update 更新角色
func (s *RoleService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.RoleRepo.Update(ctx, id, updates)
}

// UpdateStatus 更新角色状态
func (s *RoleService) UpdateStatus(ctx context.Context, id int64, status int) error {
	return s.st.RoleRepo.Update(ctx, id, map[string]any{"status": status})
}

// Delete 删除角色
func (s *RoleService) Delete(ctx context.Context, id int64) error {
	// 检查是否有用户使用
	exists, _ := s.st.BackendUserRepo.Exists(ctx, "role_id = ?", id)
	if exists {
		return errcode.ErrRoleInUse
	}

	// 删除角色菜单关联
	if err := s.st.RoleMenuRepo.DeleteByRole(ctx, id); err != nil {
		return err
	}

	return s.st.RoleRepo.Delete(ctx, id)
}

// GetRoleMenus 获取角色的菜单 ID 列表
func (s *RoleService) GetRoleMenus(ctx context.Context, roleID int64) ([]int64, error) {
	return s.st.RoleMenuRepo.ListMenuIDsByRole(ctx, roleID)
}

// UpdateRoleMenus 更新角色的菜单权限
func (s *RoleService) UpdateRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	return s.st.RoleMenuRepo.ReplaceByRole(ctx, roleID, menuIDs)
}

// InitDefaultRoles 初始化默认角色
func (s *RoleService) InitDefaultRoles(ctx context.Context) error {
	count, err := s.st.RoleRepo.Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	roles := []model.Role{
		{Name: "admin", DisplayName: "超级管理员", Description: "拥有所有权限", Status: 1},
		{Name: "editor", DisplayName: "编辑员", Description: "文章编辑权限", Status: 1},
		{Name: "viewer", DisplayName: "查看员", Description: "只读权限", Status: 1},
	}

	for i := range roles {
		if err := s.st.RoleRepo.Create(ctx, &roles[i]); err != nil {
			return err
		}
	}
	return nil
}
