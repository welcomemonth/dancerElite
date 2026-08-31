package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// RegistrationListItem 报名列表项（含活动标题）
type RegistrationListItem struct {
	model.Registration
	ActivityTitle string `json:"activity_title"`
}

// RegistrationRepo 报名数据访问接口
type RegistrationRepo interface {
	CRUD[model.Registration]
	// ListWithActivity 后台分页查询，附带活动标题（activityID>0 按活动过滤，status>=0 按状态过滤）
	ListWithActivity(ctx context.Context, page, pageSize int, activityID int64, status int) ([]RegistrationListItem, int64, error)
	// GetWithActivity 查询单个报名，附带活动标题
	GetWithActivity(ctx context.Context, id int64) (*RegistrationListItem, error)
	// ListByUserWithActivity 用户的报名列表，附带活动标题
	ListByUserWithActivity(ctx context.Context, userID int64, page, pageSize int) ([]RegistrationListItem, int64, error)
}

// RegistrationRepository RegistrationRepo 的默认实现
type RegistrationRepository struct {
	*repository.BaseRepo[model.Registration]
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) RegistrationRepo {
	return &RegistrationRepository{
		BaseRepo: repository.NewBaseRepo[model.Registration](db),
		db:       db,
	}
}

// regActivitySelect 报名 + 活动标题 join
const regActivitySelect = "registrations.*, activities.title AS activity_title"

func (r *RegistrationRepository) base(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table("registrations").
		Select(regActivitySelect).
		Joins("LEFT JOIN activities ON registrations.activity_id = activities.id")
}

func (r *RegistrationRepository) ListWithActivity(ctx context.Context, page, pageSize int, activityID int64, status int) ([]RegistrationListItem, int64, error) {
	var (
		list  []RegistrationListItem
		total int64
	)
	db := r.base(ctx).Where("registrations.deleted_at IS NULL")
	if activityID > 0 {
		db = db.Where("registrations.activity_id = ?", activityID)
	}
	if status >= 0 {
		db = db.Where("registrations.status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("registrations.created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *RegistrationRepository) GetWithActivity(ctx context.Context, id int64) (*RegistrationListItem, error) {
	var item RegistrationListItem
	err := r.base(ctx).
		Where("registrations.id = ? AND registrations.deleted_at IS NULL", id).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RegistrationRepository) ListByUserWithActivity(ctx context.Context, userID int64, page, pageSize int) ([]RegistrationListItem, int64, error) {
	var (
		list  []RegistrationListItem
		total int64
	)
	db := r.base(ctx).Where("registrations.user_id = ? AND registrations.deleted_at IS NULL", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("registrations.created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}
