package service

import (
	"context"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// maxRankingEvents 年度积分榜最多计入的场次
const maxRankingEvents = 3

// AnnualRankingService 年度积分榜管理服务（含榜单重算引擎）
type AnnualRankingService struct {
	st *store.Store
}

// NewAnnualRankingService 创建年度积分榜服务
func NewAnnualRankingService(st *store.Store) *AnnualRankingService {
	return &AnnualRankingService{st: st}
}

// AnnualRankingItem 榜单条目（含选手姓名/机构/老师）
type AnnualRankingItem struct {
	model.AnnualRanking
	PlayerName  string `json:"player_name"`
	Institution string `json:"institution"`
	Teacher     string `json:"teacher"`
}

// List 分页查询年度积分榜，支持按赛季/年龄组/舞种过滤
func (s *AnnualRankingService) List(ctx context.Context, page, pageSize int, seasonID int64, ageGroup, danceType string) ([]AnnualRankingItem, int64, error) {
	list, total, err := s.st.AnnualRankingRepo.List(ctx, page, pageSize,
		func(db *gorm.DB) *gorm.DB {
			if seasonID > 0 {
				db = db.Where("season_id = ?", seasonID)
			}
			if ageGroup != "" {
				db = db.Where("age_group = ?", ageGroup)
			}
			if danceType != "" {
				db = db.Where("dance_type = ?", danceType)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB {
			return db.Preload("Player").Order("season_id DESC, age_group, dance_type, rank ASC")
		},
	)
	if err != nil {
		return nil, 0, err
	}

	items := make([]AnnualRankingItem, 0, len(list))
	for _, r := range list {
		items = append(items, toAnnualRankingItem(r))
	}
	return items, total, nil
}

// Get 查询单条榜单记录
func (s *AnnualRankingService) Get(ctx context.Context, id int64) (*AnnualRankingItem, error) {
	list, err := s.st.AnnualRankingRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("id = ?", id).Preload("Player") },
	)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errcode.ErrNotFound
	}
	item := toAnnualRankingItem(list[0])
	return &item, nil
}

// Update 更新榜单记录（手动修正积分/名次等，会被下次重算覆盖）
func (s *AnnualRankingService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.AnnualRankingRepo.Update(ctx, id, updates)
}

// Delete 删除榜单记录
func (s *AnnualRankingService) Delete(ctx context.Context, id int64) error {
	return s.st.AnnualRankingRepo.Delete(ctx, id)
}

type rankingGroupKey struct {
	AgeGroup  string
	DanceType string
	PlayerID  int64
}

// RecalculateSeason 重算某赛季的年度积分榜。
// 规则：
//  1. 按 (年龄组, 舞种, 选手) 分组（一个选手一个级别可报多个舞种，每个舞种单独计榜）；
//  2. 每组取「本站积分」最高的 3 场求和（不足 3 场按实际场次）；
//  3. 在 (年龄组, 舞种) 维度内按总分降序排名。
func (s *AnnualRankingService) RecalculateSeason(ctx context.Context, seasonID int64) (int, error) {
	// 1. 读取旧榜单，用于计算排名变化（新上榜=0）
	oldList, err := s.st.AnnualRankingRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("season_id = ?", seasonID) },
	)
	if err != nil {
		return 0, err
	}
	oldRank := make(map[rankingGroupKey]int, len(oldList))
	for _, ar := range oldList {
		oldRank[rankingGroupKey{ar.AgeGroup, ar.DanceType, ar.PlayerID}] = ar.Rank
	}

	// 2. 读取本赛季所有成绩
	results, err := s.st.ActivityResultRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("season_id = ?", seasonID) },
	)
	if err != nil {
		return 0, err
	}

	// 3. 按组统计，取最高 3 场积分之和
	pointsByKey := make(map[rankingGroupKey][]int)
	for _, r := range results {
		k := rankingGroupKey{r.AgeGroup, r.DanceType, r.PlayerID}
		pointsByKey[k] = append(pointsByKey[k], r.Points)
	}

	type entry struct {
		key         rankingGroupKey
		totalPoints int
		scoreCount  int
	}
	entries := make([]entry, 0, len(pointsByKey))
	for k, pts := range pointsByKey {
		sort.Sort(sort.Reverse(sort.IntSlice(pts)))
		count := len(pts)
		if count > maxRankingEvents {
			count = maxRankingEvents
		}
		total := 0
		for i := 0; i < count; i++ {
			total += pts[i]
		}
		entries = append(entries, entry{key: k, totalPoints: total, scoreCount: count})
	}

	// 4. 按 (年龄组, 舞种) 分桶排名
	buckets := make(map[string][]entry)
	for _, e := range entries {
		bk := e.key.AgeGroup + "\x00" + e.key.DanceType
		buckets[bk] = append(buckets[bk], e)
	}

	now := time.Now()
	newList := make([]model.AnnualRanking, 0, len(entries))
	for _, bucket := range buckets {
		sort.SliceStable(bucket, func(i, j int) bool {
			if bucket[i].totalPoints != bucket[j].totalPoints {
				return bucket[i].totalPoints > bucket[j].totalPoints
			}
			return bucket[i].key.PlayerID < bucket[j].key.PlayerID
		})
		for i, e := range bucket {
			rank := i + 1
			prev := oldRank[e.key]
			rankChange := 0
			if prev > 0 {
				rankChange = prev - rank // 正数=排名上升
			}
			newList = append(newList, model.AnnualRanking{
				SeasonID:      seasonID,
				AgeGroup:      e.key.AgeGroup,
				DanceType:     e.key.DanceType,
				PlayerID:      e.key.PlayerID,
				TotalPoints:   e.totalPoints,
				Rank:          rank,
				PreviousRank:  prev,
				RankChange:    rankChange,
				ScoreCount:    e.scoreCount,
				LastUpdatedAt: now,
			})
		}
	}

	// 5. 事务内清空旧榜、写入新榜（硬删除避免堆积软删记录）
	if err := s.st.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("season_id = ?", seasonID).Delete(&model.AnnualRanking{}).Error; err != nil {
			return err
		}
		if len(newList) == 0 {
			return nil
		}
		return tx.Create(&newList).Error
	}); err != nil {
		return 0, err
	}

	return len(newList), nil
}

func toAnnualRankingItem(r model.AnnualRanking) AnnualRankingItem {
	item := AnnualRankingItem{AnnualRanking: r}
	if r.Player != nil {
		item.PlayerName = r.Player.RealName
		item.Institution = r.Player.Institution
		item.Teacher = r.Player.Teacher
	}
	item.Player = nil
	item.Season = nil
	return item
}
