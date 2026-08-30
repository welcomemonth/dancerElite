package model

import (
	"time"
)

type Season struct {
	BaseModel
	Year      int        `gorm:"uniqueIndex;not null" json:"year"`             // 2026
	Name      string     `gorm:"size:64;not null" json:"name"`                 // 2026赛季
	Status    string     `gorm:"size:16;default:'active';index" json:"status"` // active / archived
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

func (Season) TableName() string { return "seasons" }
