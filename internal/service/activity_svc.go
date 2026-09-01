package service

import (
	"context"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// ActivityService 活动管理服务
type ActivityService struct {
	st *store.Store
}

// NewActivityService 创建活动服务
func NewActivityService(st *store.Store) *ActivityService {
	return &ActivityService{st: st}
}

// ActivityListItem 活动列表项（含报名人数）
type ActivityListItem = store.ActivityListItem

// List 获取活动列表（后台管理）
func (s *ActivityService) List(ctx context.Context, page, pageSize int, status int) ([]ActivityListItem, int64, error) {
	return s.st.ActivityRepo.ListWithRegCount(ctx, page, pageSize, status)
}

// Get 获取活动详情
func (s *ActivityService) Get(ctx context.Context, id int64) (*ActivityListItem, error) {
	item, err := s.st.ActivityRepo.GetWithRegCount(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return item, nil
}

// Create 创建活动
// 若未显式指定 SeasonID，则默认绑定当前生效（active）赛季。
func (s *ActivityService) Create(ctx context.Context, activity *model.Activity) error {
	if activity.SeasonID == 0 {
		season, err := s.st.SeasonRepo.GetActive(ctx)
		if err != nil {
			return errcode.ErrNoActiveSeason
		}
		activity.SeasonID = season.ID
	}
	return s.st.ActivityRepo.Create(ctx, activity)
}

// Update 更新活动
func (s *ActivityService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.ActivityRepo.Update(ctx, id, updates)
}

// UpdateStatus 更新活动状态
func (s *ActivityService) UpdateStatus(ctx context.Context, id int64, status int) error {
	return s.st.ActivityRepo.Update(ctx, id, map[string]any{"status": status})
}

// Delete 删除活动
func (s *ActivityService) Delete(ctx context.Context, id int64) error {
	// 检查是否有有效报名
	exists, _ := s.st.RegistrationRepo.Exists(ctx, "activity_id = ? AND status IN (0,1) AND deleted_at IS NULL", id)
	if exists {
		return errcode.ErrActivityHasRegistrations
	}
	return s.st.ActivityRepo.Delete(ctx, id)
}

// GetForMP 获取活动详情（小程序端）
func (s *ActivityService) GetForMP(ctx context.Context, id int64) (*ActivityListItem, error) {
	item, err := s.st.ActivityRepo.GetPublishedWithRegCount(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return item, nil
}

// ListForMP 获取活动列表（小程序端，只显示非草稿的活动）
func (s *ActivityService) ListForMP(ctx context.Context, page, pageSize int) ([]ActivityListItem, int64, error) {
	return s.st.ActivityRepo.ListPublishedWithRegCount(ctx, page, pageSize)
}
