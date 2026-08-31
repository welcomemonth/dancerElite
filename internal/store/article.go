package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// ArticleListItem 文章列表项（含栏目名）
type ArticleListItem struct {
	model.Article
	ColumnName string `json:"column_name"`
}

// ArticleRepo 文章数据访问接口
type ArticleRepo interface {
	CRUD[model.Article]
	// ListWithColumn 后台分页查询，附带栏目名（columnID>0 按栏目过滤，status>=0 按状态过滤）
	ListWithColumn(ctx context.Context, page, pageSize int, columnID int64, status int) ([]ArticleListItem, int64, error)
	// GetWithColumn 查询单个文章，附带栏目名
	GetWithColumn(ctx context.Context, id int64) (*ArticleListItem, error)
	// ListByColumn 小程序端按栏目查询已发布文章
	ListByColumn(ctx context.Context, columnID int64, page, pageSize int) ([]model.Article, int64, error)
	// GetPublished 小程序端已发布文章详情，附带栏目名
	GetPublished(ctx context.Context, id int64) (*ArticleListItem, error)
	// IncrementViewCount 浏览量 +1
	IncrementViewCount(ctx context.Context, id int64) error
}

// ArticleRepository ArticleRepo 的默认实现
type ArticleRepository struct {
	*repository.BaseRepo[model.Article]
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepo {
	return &ArticleRepository{
		BaseRepo: repository.NewBaseRepo[model.Article](db),
		db:       db,
	}
}

// articleColumnSelect 文章 + 栏目名 join
const articleColumnSelect = "articles.*, columns.name AS column_name"

func (r *ArticleRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table("articles").
		Select(articleColumnSelect).
		Joins("LEFT JOIN columns ON articles.column_id = columns.id")
}

func (r *ArticleRepository) ListWithColumn(ctx context.Context, page, pageSize int, columnID int64, status int) ([]ArticleListItem, int64, error) {
	var (
		list  []ArticleListItem
		total int64
	)
	db := r.base(ctx).Where("articles.deleted_at IS NULL")
	if columnID > 0 {
		db = db.Where("articles.column_id = ?", columnID)
	}
	if status >= 0 {
		db = db.Where("articles.status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("articles.created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ArticleRepository) GetWithColumn(ctx context.Context, id int64) (*ArticleListItem, error) {
	var item ArticleListItem
	err := r.base(ctx).
		Where("articles.id = ? AND articles.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ArticleRepository) ListByColumn(ctx context.Context, columnID int64, page, pageSize int) ([]model.Article, int64, error) {
	return r.BaseRepo.List(ctx, page, pageSize, func(db *gorm.DB) *gorm.DB {
		return db.Where("column_id = ? AND status = 1", columnID).Order("created_at DESC")
	})
}

func (r *ArticleRepository) GetPublished(ctx context.Context, id int64) (*ArticleListItem, error) {
	var item ArticleListItem
	err := r.base(ctx).
		Where("articles.id = ? AND articles.status = 1 AND articles.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ArticleRepository) IncrementViewCount(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
