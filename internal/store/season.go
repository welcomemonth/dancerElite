package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// SeasonRepo 赛季数据访问接口
type SeasonRepo interface {
	CRUD[model.Season]
	// GetByYear 根据年份查询赛季（year 唯一）
	GetByYear(ctx context.Context, year int) (*model.Season, error)
	// GetActive 查询当前激活的赛季（status = active）
	GetActive(ctx context.Context) (*model.Season, error)
}

// SeasonRepository SeasonRepo 的默认实现，复用 BaseRepo 的通用 CRUD
type SeasonRepository struct {
	*repository.BaseRepo[model.Season]
	db *gorm.DB
}

func NewSeasonRepository(db *gorm.DB) SeasonRepo {
	return &SeasonRepository{
		BaseRepo: repository.NewBaseRepo[model.Season](db),
		db:       db,
	}
}

// GetByYear 根据年份查询赛季
func (r *SeasonRepository) GetByYear(ctx context.Context, year int) (*model.Season, error) {
	var season model.Season
	if err := r.db.WithContext(ctx).Where("year = ?", year).First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// GetActive 查询当前激活的赛季
func (r *SeasonRepository) GetActive(ctx context.Context) (*model.Season, error) {
	var season model.Season
	if err := r.db.WithContext(ctx).Where("status = ?", "active").First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}
