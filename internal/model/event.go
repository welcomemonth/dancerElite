package model

import "time"

type Event struct {
	BaseModel
	SeasonID             int64     `gorm:"index;not null" json:"season_id"`
	Name                 string    `gorm:"size:128;not null" json:"name"`      // 上海站超级赛
	Type                 string    `gorm:"size:16;not null;index" json:"type"` // jiaji / super
	Location             string    `gorm:"size:128" json:"location"`
	StartTime            time.Time `gorm:"index" json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	RegistrationDeadline time.Time `gorm:"index" json:"registration_deadline"`
	Fee                  int       `gorm:"not null" json:"fee"`                               // 单位：分（98000 = 980元）
	Status               string    `gorm:"size:16;default:'registering';index" json:"status"` // registering / closed / finished
	CoverURL             string    `gorm:"size:512" json:"cover_url"`
	Description          string    `gorm:"type:text" json:"description"`

	Season *Season `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

func (Event) TableName() string { return "events" }
