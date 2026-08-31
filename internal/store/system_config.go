package store

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// SystemConfigRepo 系统配置数据访问接口
type SystemConfigRepo interface {
	CRUD[model.SystemConfig]
	// GetByKey 根据 key 查询配置（key 唯一）
	GetByKey(ctx context.Context, key string) (*model.SystemConfig, error)
	// Upsert 按 key 更新或插入配置
	Upsert(ctx context.Context, config *model.SystemConfig) error
	// BatchUpsert 批量按 key 更新或插入配置（事务）
	BatchUpsert(ctx context.Context, configs []model.SystemConfig) error
	// ListGroups 所有非空分组名
	ListGroups(ctx context.Context) ([]string, error)
	// DeleteByKey 按 key 软删除配置
	DeleteByKey(ctx context.Context, key string) error
}

// SystemConfigRepository SystemConfigRepo 的默认实现，复用 BaseRepo 的通用 CRUD
type SystemConfigRepository struct {
	*repository.BaseRepo[model.SystemConfig]
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepo {
	return &SystemConfigRepository{
		BaseRepo: repository.NewBaseRepo[model.SystemConfig](db),
		db:       db,
	}
}

// upsertClause 按 key 冲突时更新相关字段
func upsertClause() clause.OnConflict {
	return clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "type", "group_name", "remark", "updated_at"}),
	}
}

// GetByKey 根据 key 查询配置
func (r *SystemConfigRepository) GetByKey(ctx context.Context, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	if err := r.db.WithContext(ctx).Where("`key` = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// Upsert 按 key 更新或插入配置
func (r *SystemConfigRepository) Upsert(ctx context.Context, config *model.SystemConfig) error {
	return r.db.WithContext(ctx).Clauses(upsertClause()).Create(config).Error
}

// BatchUpsert 批量按 key 更新或插入配置（事务）
func (r *SystemConfigRepository) BatchUpsert(ctx context.Context, configs []model.SystemConfig) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, c := range configs {
			if err := tx.Clauses(upsertClause()).Create(&c).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ListGroups 所有非空分组名
func (r *SystemConfigRepository) ListGroups(ctx context.Context) ([]string, error) {
	var groups []string
	err := r.db.WithContext(ctx).Model(&model.SystemConfig{}).
		Distinct("group_name").
		Where("group_name != ''").
		Pluck("group_name", &groups).Error
	return groups, err
}

// DeleteByKey 按 key 软删除配置
func (r *SystemConfigRepository) DeleteByKey(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("`key` = ?", key).Delete(&model.SystemConfig{}).Error
}
