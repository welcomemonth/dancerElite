package model

import "time"

type Event struct {
	BaseModel
	SeasonID        int64      `gorm:"index;not null" json:"season_id"`
	Title           string     `gorm:"size:128;not null" json:"title"` // 上海站超级赛
	Content         string     `gorm:"type:text" json:"content"`
	Type            string     `gorm:"size:16;not null;index" json:"type"` // jiaji / super
	Location        string     `gorm:"size:128" json:"location"`
	StartTime       time.Time  `gorm:"index" json:"start_time"`
	EndTime         time.Time  `json:"end_time"`
	RegStartTime    *time.Time `json:"reg_start_time,omitempty"`
	RegEndTime      *time.Time `json:"reg_end_time,omitempty"`
	MaxParticipants int        `gorm:"default:0" json:"max_participants"` // 0=不限
	CurParticipants int        `gorm:"default:0" json:"cur_participants"` //
	Fee             float64    `gorm:"not null" json:"fee"`               // 单位：分（98000 = 980元）
	Status          int        `gorm:"default:0" json:"status"`           // 0:草稿 1:报名中 2:报名截止 3:进行中 4:已结束
	CoverURL        string     `gorm:"size:512" json:"cover_url"`
	Description     string     `gorm:"type:text" json:"description"`

	Season *Season `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

func (Event) TableName() string { return "events" }
