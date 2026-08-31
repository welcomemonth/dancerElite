package service

import (
	"context"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// ArticleService 文章管理服务
type ArticleService struct {
	st *store.Store
}

// NewArticleService 创建文章服务
func NewArticleService(st *store.Store) *ArticleService {
	return &ArticleService{st: st}
}

// ArticleListItem 文章列表项（含栏目名）
type ArticleListItem = store.ArticleListItem

// List 获取文章列表（后台管理）
func (s *ArticleService) List(ctx context.Context, page, pageSize int, columnID int64, status int) ([]ArticleListItem, int64, error) {
	return s.st.ArticleRepo.ListWithColumn(ctx, page, pageSize, columnID, status)
}

// Get 获取文章详情（后台）
func (s *ArticleService) Get(ctx context.Context, id int64) (*ArticleListItem, error) {
	item, err := s.st.ArticleRepo.GetWithColumn(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return item, nil
}

// Create 创建文章
func (s *ArticleService) Create(ctx context.Context, article *model.Article) error {
	return s.st.ArticleRepo.Create(ctx, article)
}

// Update 更新文章
func (s *ArticleService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.ArticleRepo.Update(ctx, id, updates)
}

// UpdateStatus 更新文章状态
func (s *ArticleService) UpdateStatus(ctx context.Context, id int64, status int) error {
	return s.st.ArticleRepo.Update(ctx, id, map[string]any{"status": status})
}

// Delete 删除文章
func (s *ArticleService) Delete(ctx context.Context, id int64) error {
	return s.st.ArticleRepo.Delete(ctx, id)
}

// ListByColumn 根据栏目获取已发布文章（小程序端）
func (s *ArticleService) ListByColumn(ctx context.Context, columnID int64, page, pageSize int) ([]model.Article, int64, error) {
	return s.st.ArticleRepo.ListByColumn(ctx, columnID, page, pageSize)
}

// GetForMP 获取文章详情（小程序端，增加浏览量）
func (s *ArticleService) GetForMP(ctx context.Context, id int64) (*ArticleListItem, error) {
	// 增加浏览量
	_ = s.st.ArticleRepo.IncrementViewCount(ctx, id)

	item, err := s.st.ArticleRepo.GetPublished(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return item, nil
}
