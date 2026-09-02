package model

import "time"

type AnnualRanking struct {
	BaseModel
	SeasonID      int64     `gorm:"index;not null" json:"season_id"`
	AgeGroup      string    `gorm:"size:8;not null;index" json:"age_group"`
	DanceType     string    `gorm:"size:16;not null;index" json:"dance_type"`
	PlayerID      int64     `gorm:"index;not null" json:"player_id"`
	TotalPoints   int       `gorm:"not null;default:0" json:"total_points"` // 最高3场积分之和
	Rank          int       `gorm:"not null" json:"rank"`                   // 当前排名
	PreviousRank  int       `gorm:"default:0" json:"previous_rank"`         // 上一次排名（0表示新上榜）
	RankChange    int       `gorm:"default:0" json:"rank_change"`           // 正数上升，负数下降，0不变
	ScoreCount      int       `gorm:"default:0" json:"score_count"`           // 实际计入场次（1~3）
	LastUpdatedAt   time.Time `json:"last_updated_at"`                        // 本次重算时间
	IsDirectAdvance bool      `gorm:"default:false" json:"is_direct_advance"` // 是否直通年底冠军决赛（甲级赛冠军/超级赛冠亚军，重算时写入）

	// 唯一：一个赛季 + 年龄组 + 舞种 + 选手 只有一条
	// uniqueIndex: uk_season_age_dance_player

	Player *Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Season *Season `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

func (AnnualRanking) TableName() string { return "annual_rankings" }
