package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// UserRepo 小程序用户数据访问接口
type UserRepo interface {
	CRUD[model.User]
	// GetByOpenID 根据 openid 查询用户
	GetByOpenID(ctx context.Context, openID string) (*model.User, error)
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
