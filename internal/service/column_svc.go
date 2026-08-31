package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/pkg/errcode"
	"github.com/zzhtl/go-mountain/internal/store"
)

// ColumnService 栏目管理服务
type ColumnService struct {
	st *store.Store
}

// NewColumnService 创建栏目服务
func NewColumnService(st *store.Store) *ColumnService {
	return &ColumnService{st: st}
}

// List 获取所有栏目
func (s *ColumnService) List(ctx context.Context) ([]model.Column, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order, id")
	}
	return s.st.ColumnRepo.FindAll(ctx, scope)
}

// Get 获取单个栏目
func (s *ColumnService) Get(ctx context.Context, id int64) (*model.Column, error) {
	return s.st.ColumnRepo.GetByID(ctx, id)
}

// Create 创建栏目
func (s *ColumnService) Create(ctx context.Context, column *model.Column) error {
	return s.st.ColumnRepo.Create(ctx, column)
}

// Update 更新栏目
func (s *ColumnService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.ColumnRepo.Update(ctx, id, updates)
}

// Delete 删除栏目
func (s *ColumnService) Delete(ctx context.Context, id int64) error {
	// 检查是否有文章关联
	exists, _ := s.st.ArticleRepo.Exists(ctx, "column_id = ?", id)
	if exists {
		return errcode.ErrColumnHasArticles
	}
	return s.st.ColumnRepo.Delete(ctx, id)
}

// ListForMP 小程序获取栏目列表
func (s *ColumnService) ListForMP(ctx context.Context) ([]model.Column, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, description").Order("sort_order, id")
	}
	return s.st.ColumnRepo.FindAll(ctx, scope)
}
