package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// SeasonService 赛季管理服务
type SeasonService struct {
	st *store.Store
}

// NewSeasonService 创建赛季服务
func NewSeasonService(st *store.Store) *SeasonService {
	return &SeasonService{st: st}
}

// List 获取所有赛季（按年份倒序）
func (s *SeasonService) List(ctx context.Context) ([]model.Season, error) {
	return s.st.SeasonRepo.FindAll(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("year desc")
	})
}

// Get 获取单个赛季
func (s *SeasonService) Get(ctx context.Context, id int64) (*model.Season, error) {
	season, err := s.st.SeasonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return season, nil
}

// Create 创建赛季（year 唯一，且同一时间只能有一个激活赛季）
func (s *SeasonService) Create(ctx context.Context, season *model.Season) error {
	if _, err := s.st.SeasonRepo.GetByYear(ctx, season.Year); err == nil {
		return errcode.ErrAlreadyExists
	}
	if season.Status == "active" {
		if err := s.deactivateOthers(ctx, 0); err != nil {
			return err
		}
	}
	return s.st.SeasonRepo.Create(ctx, season)
}

// Update 更新赛季
func (s *SeasonService) Update(ctx context.Context, id int64, updates map[string]any) error {
	if year, ok := updates["year"].(int); ok {
		existing, err := s.st.SeasonRepo.GetByYear(ctx, year)
		if err == nil && existing.ID != id {
			return errcode.ErrAlreadyExists
		}
	}
	if status, ok := updates["status"].(string); ok && status == "active" {
		if err := s.deactivateOthers(ctx, id); err != nil {
			return err
		}
	}
	return s.st.SeasonRepo.Update(ctx, id, updates)
}

// UpdateStatus 更新赛季状态（active / archived）
func (s *SeasonService) UpdateStatus(ctx context.Context, id int64, status string) error {
	if status == "active" {
		if err := s.deactivateOthers(ctx, id); err != nil {
			return err
		}
	}
	return s.st.SeasonRepo.Update(ctx, id, map[string]any{"status": status})
}

// Delete 删除赛季
func (s *SeasonService) Delete(ctx context.Context, id int64) error {
	return s.st.SeasonRepo.Delete(ctx, id)
}

// deactivateOthers 将除 excludeID 外的所有激活赛季归档，保证同一时间只有一个激活赛季
func (s *SeasonService) deactivateOthers(ctx context.Context, excludeID int64) error {
	return s.st.DB().WithContext(ctx).
		Model(&model.Season{}).
		Where("status = ? AND id <> ?", "active", excludeID).
		Update("status", "archived").Error
}
