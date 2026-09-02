package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// PlayerService 选手管理服务
type PlayerService struct {
	st *store.Store
}

// NewPlayerService 创建选手服务
func NewPlayerService(st *store.Store) *PlayerService {
	return &PlayerService{st: st}
}

// List 分页查询选手，支持按姓名/机构/年龄组过滤
func (s *PlayerService) List(ctx context.Context, page, pageSize int, name, institution, ageGroup string) ([]model.Player, int64, error) {
	return s.st.PlayerRepo.List(ctx, page, pageSize,
		func(db *gorm.DB) *gorm.DB {
			if name != "" {
				db = db.Where("real_name LIKE ?", "%"+name+"%")
			}
			if institution != "" {
				db = db.Where("institution LIKE ?", "%"+institution+"%")
			}
			if ageGroup != "" {
				db = db.Where("age_group = ?", ageGroup)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB { return db.Order("created_at desc") },
	)
}

// Get 查询选手详情
func (s *PlayerService) Get(ctx context.Context, id int64) (*model.Player, error) {
	player, err := s.st.PlayerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return player, nil
}

// Update 更新选手资料
// 注：age_group 仅是“当前年龄组”的缓存，改它只影响之后新产生的成绩/榜单，
// 不会回写历史赛季数据（历史赛季的年龄组冻结在 ActivityResult / AnnualRanking 上）。
func (s *PlayerService) Update(ctx context.Context, id int64, updates map[string]any) error {
	if _, err := s.st.PlayerRepo.GetByID(ctx, id); err != nil {
		return errcode.ErrNotFound
	}
	return s.st.PlayerRepo.Update(ctx, id, updates)
}

// Delete 删除选手
func (s *PlayerService) Delete(ctx context.Context, id int64) error {
	return s.st.PlayerRepo.Delete(ctx, id)
}
