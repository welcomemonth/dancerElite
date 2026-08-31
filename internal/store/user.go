package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// UserRepo 小程序用户数据访问接口
type UserRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByOpenID(ctx context.Context, openID string) (*model.User, error)
	Update(ctx context.Context, id int64, updates map[string]any) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int, scopes ...func(*gorm.DB) *gorm.DB) ([]model.User, int64, error)
}

// UserRepository UserRepo 的默认实现，复用 BaseRepo 的通用 CRUD
type UserRepository struct {
	*repository.BaseRepo[model.User]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepo {
	return &UserRepository{
		BaseRepo: repository.NewBaseRepo[model.User](db),
		db:       db,
	}
}

// GetByOpenID 根据 openid 查询用户
func (r *UserRepository) GetByOpenID(ctx context.Context, openID string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("open_id = ?", openID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
