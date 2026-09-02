package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// ActivityResultService 赛事成绩管理服务
type ActivityResultService struct {
	st *store.Store
}

// NewActivityResultService 创建赛事成绩服务
func NewActivityResultService(st *store.Store) *ActivityResultService {
	return &ActivityResultService{st: st}
}

// ActivityResultItem 成绩列表项（含活动标题、选手姓名）
type ActivityResultItem struct {
	model.ActivityResult
	ActivityTitle string `json:"activity_title"`
	PlayerName    string `json:"player_name"`
}

// List 分页查询成绩，支持按活动/选手/赛季/舞种/年龄组过滤
func (s *ActivityResultService) List(ctx context.Context, page, pageSize int, activityID, playerID, seasonID int64, danceType, ageGroup string) ([]ActivityResultItem, int64, error) {
	list, total, err := s.st.ActivityResultRepo.List(ctx, page, pageSize,
		func(db *gorm.DB) *gorm.DB {
			if activityID > 0 {
				db = db.Where("activity_id = ?", activityID)
			}
			if playerID > 0 {
				db = db.Where("player_id = ?", playerID)
			}
			if seasonID > 0 {
				db = db.Where("season_id = ?", seasonID)
			}
			if danceType != "" {
				db = db.Where("dance_type = ?", danceType)
			}
			if ageGroup != "" {
				db = db.Where("age_group = ?", ageGroup)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB { return db.Preload("Activity").Preload("Player").Order("created_at desc") },
	)
	if err != nil {
		return nil, 0, err
	}

	items := make([]ActivityResultItem, 0, len(list))
	for _, r := range list {
		items = append(items, toActivityResultItem(r))
	}
	return items, total, nil
}

// Get 查询单条成绩（含活动、选手关联）
func (s *ActivityResultService) Get(ctx context.Context, id int64) (*ActivityResultItem, error) {
	list, err := s.st.ActivityResultRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("id = ?", id).Preload("Activity").Preload("Player") },
	)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errcode.ErrNotFound
	}
	item := toActivityResultItem(list[0])
	return &item, nil
}

// Create 创建成绩：校验关联存在、去重（同一活动同一选手同一舞种仅一条）、缺省赛季
func (s *ActivityResultService) Create(ctx context.Context, result *model.ActivityResult) error {
	if result.SeasonID == 0 {
		if err := s.resolveSeason(ctx, result); err != nil {
			return err
		}
	}
	if err := s.validate(ctx, result); err != nil {
		return err
	}

	exists, _ := s.st.ActivityResultRepo.Exists(ctx,
		"activity_id = ? AND player_id = ? AND dance_type = ?",
		result.ActivityID, result.PlayerID, result.DanceType)
	if exists {
		return errcode.ErrActivityResultExists
	}
	return s.st.ActivityResultRepo.Create(ctx, result)
}

// Update 更新成绩
func (s *ActivityResultService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.ActivityResultRepo.Update(ctx, id, updates)
}

// Delete 删除成绩
func (s *ActivityResultService) Delete(ctx context.Context, id int64) error {
	return s.st.ActivityResultRepo.Delete(ctx, id)
}

// resolveSeason 未显式传入赛季时，优先取所属活动的赛季，其次取当前生效赛季
func (s *ActivityResultService) resolveSeason(ctx context.Context, result *model.ActivityResult) error {
	if result.ActivityID > 0 {
		if act, err := s.st.ActivityRepo.GetByID(ctx, result.ActivityID); err == nil && act.SeasonID > 0 {
			result.SeasonID = act.SeasonID
			return nil
		}
	}
	season, err := s.st.SeasonRepo.GetActive(ctx)
	if err != nil {
		return errcode.ErrNoActiveSeason
	}
	result.SeasonID = season.ID
	return nil
}

// validate 校验必填字段及活动、选手是否存在，且活动必须已结束
func (s *ActivityResultService) validate(ctx context.Context, result *model.ActivityResult) error {
	if result.ActivityID == 0 || result.PlayerID == 0 || result.DanceType == "" || result.AgeGroup == "" {
		return errcode.ErrInvalidParam
	}
	act, err := s.st.ActivityRepo.GetByID(ctx, result.ActivityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	if act.Status != 4 {
		return errcode.ErrActivityNotEnded
	}
	if _, err := s.st.PlayerRepo.GetByID(ctx, result.PlayerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	return nil
}

// toActivityResultItem 映射为列表项，并清空嵌套关联避免 JSON 重复输出
func toActivityResultItem(r model.ActivityResult) ActivityResultItem {
	item := ActivityResultItem{ActivityResult: r}
	if r.Activity != nil {
		item.ActivityTitle = r.Activity.Title
	}
	if r.Player != nil {
		item.PlayerName = r.Player.RealName
	}
	item.Activity = nil
	item.Player = nil
	return item
}
