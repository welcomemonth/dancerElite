package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// ActivityListItem 活动列表项（含报名人数）
type ActivityListItem struct {
	model.Activity
	RegCount int64 `json:"reg_count"`
}

// ActivityRepo 活动数据访问接口
type ActivityRepo interface {
	CRUD[model.Activity]
	// ListWithRegCount 后台分页查询，附带报名人数
	ListWithRegCount(ctx context.Context, page, pageSize, status int) ([]ActivityListItem, int64, error)
	// GetWithRegCount 查询单个活动，附带报名人数
	GetWithRegCount(ctx context.Context, id int64) (*ActivityListItem, error)
	// ListPublishedWithRegCount 小程序端列表（非草稿）
	ListPublishedWithRegCount(ctx context.Context, page, pageSize int) ([]ActivityListItem, int64, error)
	// GetPublishedWithRegCount 小程序端详情（报名中/报名截止/进行中）
	GetPublishedWithRegCount(ctx context.Context, id int64) (*ActivityListItem, error)
}

// ActivityRepository ActivityRepo 的默认实现
type ActivityRepository struct {
	*repository.BaseRepo[model.Activity]
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepo {
	return &ActivityRepository{
		BaseRepo: repository.NewBaseRepo[model.Activity](db),
		db:       db,
	}
}

// regCountSelect 报名人数子查询（status IN (0,1) 即待支付+已支付）
const regCountSelect = "activities.*, (SELECT COUNT(*) FROM registrations WHERE registrations.activity_id = activities.id AND registrations.status IN (0,1) AND registrations.deleted_at IS NULL) AS reg_count"

func (r *ActivityRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table("activities").Select(regCountSelect)
}

func (r *ActivityRepository) ListWithRegCount(ctx context.Context, page, pageSize, status int) ([]ActivityListItem, int64, error) {
	var (
		list  []ActivityListItem
		total int64
	)
	db := r.base(ctx).Where("activities.deleted_at IS NULL")
	if status >= 0 {
		db = db.Where("activities.status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("activities.created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ActivityRepository) GetWithRegCount(ctx context.Context, id int64) (*ActivityListItem, error) {
	var item ActivityListItem
	err := r.base(ctx).
		Where("activities.id = ? AND activities.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ActivityRepository) ListPublishedWithRegCount(ctx context.Context, page, pageSize int) ([]ActivityListItem, int64, error) {
	var (
		list  []ActivityListItem
		total int64
	)
	db := r.base(ctx).Where("activities.status IN (1,2,3,4) AND activities.deleted_at IS NULL")
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("activities.start_time DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *ActivityRepository) GetPublishedWithRegCount(ctx context.Context, id int64) (*ActivityListItem, error) {
	var item ActivityListItem
	err := r.base(ctx).
		Where("activities.id = ? AND activities.status IN (1,2,3) AND activities.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
