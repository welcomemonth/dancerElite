package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// BackendUserListItem 后台用户列表项（含角色信息）
type BackendUserListItem struct {
	model.BackendUser
	RoleName    string `json:"role_name"`
	RoleDisplay string `json:"role_display"`
}

// BackendUserRepo 后台用户数据访问接口
type BackendUserRepo interface {
	CRUD[model.BackendUser]
	// GetByUsername 根据用户名查询（预加载角色）
	GetByUsername(ctx context.Context, username string) (*model.BackendUser, error)
	// ListWithRole 分页查询，附带角色名称
	ListWithRole(ctx context.Context, page, pageSize int) ([]BackendUserListItem, int64, error)
	// GetWithRole 查询单个用户，附带角色名称
	GetWithRole(ctx context.Context, id int64) (*BackendUserListItem, error)
}

// BackendUserRepository BackendUserRepo 的默认实现
type BackendUserRepository struct {
	*repository.BaseRepo[model.BackendUser]
	db *gorm.DB
}

func NewBackendUserRepository(db *gorm.DB) BackendUserRepo {
	return &BackendUserRepository{
		BaseRepo: repository.NewBaseRepo[model.BackendUser](db),
		db:       db,
	}
}

func (r *BackendUserRepository) GetByUsername(ctx context.Context, username string) (*model.BackendUser, error) {
	var u model.BackendUser
	if err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *BackendUserRepository) ListWithRole(ctx context.Context, page, pageSize int) ([]BackendUserListItem, int64, error) {
	var (
		list  []BackendUserListItem
		total int64
	)
	if err := r.db.WithContext(ctx).Model(&model.BackendUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Table("backend_users").
		Select("backend_users.*, roles.name as role_name, roles.display_name as role_display").
		Joins("LEFT JOIN roles ON backend_users.role_id = roles.id").
		Where("backend_users.deleted_at IS NULL").
		Order("backend_users.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *BackendUserRepository) GetWithRole(ctx context.Context, id int64) (*BackendUserListItem, error) {
	var item BackendUserListItem
	err := r.db.WithContext(ctx).
		Table("backend_users").
		Select("backend_users.*, roles.name as role_name, roles.display_name as role_display").
		Joins("LEFT JOIN roles ON backend_users.role_id = roles.id").
		Where("backend_users.id = ? AND backend_users.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
