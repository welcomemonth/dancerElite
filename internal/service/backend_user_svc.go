package service

import (
	"context"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/pkg/crypto"
	"github.com/zzhtl/go-mountain/internal/pkg/errcode"
	"github.com/zzhtl/go-mountain/internal/store"
)

// BackendUserService 后台用户管理服务
type BackendUserService struct {
	st *store.Store
}

// NewBackendUserService 创建后台用户管理服务
func NewBackendUserService(st *store.Store) *BackendUserService {
	return &BackendUserService{st: st}
}

// BackendUserListItem 后台用户列表项（含角色信息）
type BackendUserListItem = store.BackendUserListItem

// List 获取后台用户列表
func (s *BackendUserService) List(ctx context.Context, page, pageSize int) ([]BackendUserListItem, int64, error) {
	return s.st.BackendUserRepo.ListWithRole(ctx, page, pageSize)
}

// Get 获取单个后台用户
func (s *BackendUserService) Get(ctx context.Context, id int64) (*BackendUserListItem, error) {
	item, err := s.st.BackendUserRepo.GetWithRole(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return item, nil
}

// Create 创建后台用户，返回明文密码
func (s *BackendUserService) Create(ctx context.Context, username, email string, roleID int64) (*model.BackendUser, string, error) {
	// 验证角色是否存在
	exists, _ := s.st.RoleRepo.Exists(ctx, "id = ? AND status = 1", roleID)
	if !exists {
		return nil, "", errcode.ErrInvalidParam
	}

	password := GenerateRandomPassword(8)
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	user := &model.BackendUser{
		Username:        username,
		Email:           email,
		Password:        hashedPassword,
		RoleID:          roleID,
		PasswordVersion: 2,
		Status:          1,
	}

	if err := s.st.BackendUserRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	return user, password, nil
}

// Update 更新后台用户信息
func (s *BackendUserService) Update(ctx context.Context, id int64, username, email string, roleID int64) error {
	// 验证角色
	exists, _ := s.st.RoleRepo.Exists(ctx, "id = ? AND status = 1", roleID)
	if !exists {
		return errcode.ErrInvalidParam
	}

	return s.st.BackendUserRepo.Update(ctx, id, map[string]any{
		"username": username,
		"email":    email,
		"role_id":  roleID,
	})
}

// UpdateStatus 更新用户状态
func (s *BackendUserService) UpdateStatus(ctx context.Context, id int64, status int) error {
	return s.st.BackendUserRepo.Update(ctx, id, map[string]any{"status": status})
}

// ResetPassword 重置用户密码，返回新的明文密码
func (s *BackendUserService) ResetPassword(ctx context.Context, id int64) (string, error) {
	password := GenerateRandomPassword(8)
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return "", err
	}

	err = s.st.BackendUserRepo.Update(ctx, id, map[string]any{
		"password":         hashedPassword,
		"password_version": 2,
	})
	if err != nil {
		return "", err
	}

	return password, nil
}

// Delete 删除后台用户
func (s *BackendUserService) Delete(ctx context.Context, id int64) error {
	return s.st.BackendUserRepo.Delete(ctx, id)
}

// GetCurrentUserMenus 获取当前用户的菜单权限树
func (s *BackendUserService) GetCurrentUserMenus(ctx context.Context, userID int64) ([]*model.Menu, error) {
	user, err := s.st.BackendUserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	// 检查是否是 admin 角色（拥有所有菜单）
	var role model.Role
	s.st.DB().WithContext(ctx).First(&role, user.RoleID)

	var menus []model.Menu
	if role.Name == "admin" {
		s.st.DB().WithContext(ctx).Where("status = 1").Order("sort, id").Find(&menus)
	} else {
		s.st.DB().WithContext(ctx).
			Joins("INNER JOIN role_menus ON menus.id = role_menus.menu_id").
			Where("role_menus.role_id = ? AND menus.status = 1", user.RoleID).
			Order("menus.sort, menus.id").
			Find(&menus)
	}

	return buildMenuTree(menus, 0), nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []model.Menu, parentID int64) []*model.Menu {
	var tree []*model.Menu
	for i := range menus {
		if menus[i].ParentID == parentID {
			menu := &menus[i]
			menu.Children = buildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}
