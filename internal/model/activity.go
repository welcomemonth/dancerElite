package model

import "time"

// 赛事等级（字符串自描述，避免数字含义需要查表）
const (
	// ActivityLevelA 甲级赛：冠军直通年底冠军决赛
	ActivityLevelA = "甲级赛"
	// ActivityLevelSuper 超级赛：冠亚军直通年底冠军决赛
	ActivityLevelSuper = "超级赛"
)

// Activity 活动
type Activity struct {
	BaseModel
	Title           string      `gorm:"type:text;not null" json:"title"`
	Description     string      `gorm:"type:text" json:"description"`
	Content         string      `gorm:"type:text" json:"content"`
	Thumbnail       string      `gorm:"type:text" json:"thumbnail"`
	Location        string      `gorm:"type:text" json:"location"`
	StartTime       time.Time   `json:"start_time"`
	EndTime         time.Time   `json:"end_time"`
	RegStartTime    *time.Time  `json:"reg_start_time,omitempty"`
	RegEndTime      *time.Time  `json:"reg_end_time,omitempty"`
	MaxParticipants int         `gorm:"default:0" json:"max_participants"` // 0=不限
	Price           float64     `gorm:"type:decimal(10,2);default:0" json:"price"`
	Status          int         `gorm:"default:0" json:"status"` // 0:草稿 1:报名中 2:报名截止 3:进行中 4:已结束
	Level           string      `gorm:"size:16;not null;default:'甲级赛'" json:"level"` // 赛事等级：甲级赛/超级赛
	CreatedBy       int64       `json:"created_by"`
	AgeGroups       StringSlice `gorm:"type:text" json:"age_groups"`  // 级别组合 U11/U13/U15（可扩展）
	DanceTypes      StringSlice `gorm:"type:text" json:"dance_types"` // 舞种 古典舞/民族民间舞（可扩展）
	SeasonID        int64       `gorm:"not null;index" json:"season_id"`
	Season          *Season     `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

func (Activity) TableName() string {
	return "activities"
}
