package store

import (
	"context"

	"gorm.io/gorm"
)

// CRUD 通用数据访问接口，由 repository.BaseRepo[T] 提供默认实现。
// 各实体 Repo 接口内嵌它即可获得通用增删改查，仅需额外声明特殊查询方法。
type CRUD[T any] interface {
	Create(ctx context.Context, e *T) error
	GetByID(ctx context.Context, id int64) (*T, error)
	Update(ctx context.Context, id int64, m map[string]any) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, pageSize int, scopes ...func(*gorm.DB) *gorm.DB) ([]T, int64, error)
	FindAll(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]T, error)
	Count(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
	Exists(ctx context.Context, query string, args ...any) (bool, error)
}
