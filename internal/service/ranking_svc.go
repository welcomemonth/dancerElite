package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// RankingService 排行榜 / 选手 / 成绩只读服务（小程序端）
type RankingService struct {
	st *store.Store
}

// NewRankingService 创建排行榜服务
func NewRankingService(st *store.Store) *RankingService {
	return &RankingService{st: st}
}

// RankingItem 年度积分排行榜条目
type RankingItem struct {
	Rank         int    `json:"rank"`
	PlayerID     int64  `json:"player_id"`
	Name         string `json:"name"`
	Institution  string `json:"institution"`
	Teacher      string `json:"teacher"`
	AgeGroup     string `json:"age_group"`
	DanceType    string `json:"dance_type"`
	TotalPoints  int    `json:"total_points"`
	ScoreCount   int    `json:"score_count"`
	RankChange   int    `json:"rank_change"`
	PreviousRank int    `json:"previous_rank"`
}

// OrgRankingItem 机构排行榜条目
type OrgRankingItem struct {
	Rank        int    `json:"rank"`
	Institution string `json:"institution"`
	TotalPoints int    `json:"total_points"`
	PlayerCount int    `json:"player_count"`
}

// ResultItem 选手某场成绩
type ResultItem struct {
	ActivityTitle string     `json:"activity_title"`
	ActivityDate  *time.Time `json:"activity_date,omitempty"`
	DanceType     string     `json:"dance_type"`
	AgeGroup      string     `json:"age_group"`
	Rank          int        `json:"rank"`
	Points        int        `json:"points"`
	Award         string     `json:"award"`
}

// PlayerDetail 选手详情（资料 + 成绩）
type PlayerDetail struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Gender      string       `json:"gender"`
	Institution string       `json:"institution"`
	Teacher     string       `json:"teacher"`
	AgeGroup    string       `json:"age_group"`
	BirthYear   int          `json:"birth_year"`
	Results     []ResultItem `json:"results"`
}

// ActiveSeason 当前激活赛季
func (s *RankingService) ActiveSeason(ctx context.Context) (*model.Season, error) {
	return s.st.SeasonRepo.GetActive(ctx)
}

// resolveSeasonID 若未显式传入赛季，则回退到当前激活赛季
func (s *RankingService) resolveSeasonID(ctx context.Context, seasonID int64) (int64, error) {
	if seasonID > 0 {
		return seasonID, nil
	}
	season, err := s.st.SeasonRepo.GetActive(ctx)
	if err != nil {
		return 0, err
	}
	return season.ID, nil
}

// Leaderboard 年度积分排行榜（按赛季 + 可选年龄组/舞种过滤，rank 升序）
func (s *RankingService) Leaderboard(ctx context.Context, seasonID int64, ageGroup, danceType string) ([]RankingItem, error) {
	seasonID, err := s.resolveSeasonID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []RankingItem{}, nil
		}
		return nil, err
	}

	list, err := s.st.AnnualRankingRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB { return db.Where("season_id = ?", seasonID) },
		func(db *gorm.DB) *gorm.DB {
			if ageGroup != "" {
				return db.Where("age_group = ?", ageGroup)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB {
			if danceType != "" {
				return db.Where("dance_type = ?", danceType)
			}
			return db
		},
		func(db *gorm.DB) *gorm.DB { return db.Order("rank asc").Preload("Player") },
	)
	if err != nil {
		return nil, err
	}

	items := make([]RankingItem, 0, len(list))
	for _, ar := range list {
		item := RankingItem{
			Rank:         ar.Rank,
			PlayerID:     ar.PlayerID,
			AgeGroup:     ar.AgeGroup,
			DanceType:    ar.DanceType,
			TotalPoints:  ar.TotalPoints,
			ScoreCount:   ar.ScoreCount,
			RankChange:   ar.RankChange,
			PreviousRank: ar.PreviousRank,
		}
		if ar.Player != nil {
			item.Name = ar.Player.RealName
			item.Institution = ar.Player.Institution
			item.Teacher = ar.Player.Teacher
		}
		items = append(items, item)
	}
	return items, nil
}

// PlayerDetail 选手详情 + 成绩列表
func (s *RankingService) PlayerDetail(ctx context.Context, playerID int64) (*PlayerDetail, error) {
	player, err := s.st.PlayerRepo.GetByID(ctx, playerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}

	results, err := s.st.ActivityResultRepo.FindAll(ctx,
		func(db *gorm.DB) *gorm.DB {
			return db.Where("player_id = ?", playerID).Preload("Activity").Order("created_at desc")
		},
	)
	if err != nil {
		return nil, err
	}

	detail := &PlayerDetail{
		ID:          player.ID,
		Name:        player.RealName,
		Gender:      player.Gender,
		Institution: player.Institution,
		Teacher:     player.Teacher,
		AgeGroup:    player.AgeGroup,
		BirthYear:   player.BirthYear,
		Results:     make([]ResultItem, 0, len(results)),
	}

	for _, r := range results {
		item := ResultItem{
			DanceType: r.DanceType,
			AgeGroup:  r.AgeGroup,
			Rank:      r.Rank,
			Points:    r.Points,
			Award:     r.Award,
		}
		if r.Activity != nil {
			item.ActivityTitle = r.Activity.Title
			item.ActivityDate = &r.Activity.StartTime
		}
		detail.Results = append(detail.Results, item)
	}

	return detail, nil
}

// OrgLeaderboard 机构排行榜（按机构聚合积分）
func (s *RankingService) OrgLeaderboard(ctx context.Context, seasonID int64) ([]OrgRankingItem, error) {
	seasonID, err := s.resolveSeasonID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []OrgRankingItem{}, nil
		}
		return nil, err
	}

	type orgRow struct {
		Institution string
		TotalPoints int
		PlayerCount int
	}
	var rows []orgRow
	if err := s.st.DB().WithContext(ctx).
		Table("annual_rankings ar").
		Joins("JOIN players p ON p.id = ar.player_id").
		Where("ar.season_id = ? AND p.institution <> ''", seasonID).
		Select("p.institution AS institution, SUM(ar.total_points) AS total_points, COUNT(DISTINCT p.id) AS player_count").
		Group("p.institution").
		Order("total_points DESC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]OrgRankingItem, 0, len(rows))
	for i, r := range rows {
		items = append(items, OrgRankingItem{
			Rank:        i + 1,
			Institution: r.Institution,
			TotalPoints: r.TotalPoints,
			PlayerCount: r.PlayerCount,
		})
	}
	return items, nil
}
