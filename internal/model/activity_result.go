package model

type ActivityResult struct {
	BaseModel
	ActivityID     int64  `gorm:"index;not null" json:"activity_id"`
	SeasonID       int64  `gorm:"index;not null" json:"season_id"`
	PlayerID       int64  `gorm:"index" json:"player_id"`
	RegistrationID *int64 `gorm:"index" json:"registration_id"` // 可选关联
	DanceType      string `gorm:"size:16;not null;index" json:"dance_type"`
	AgeGroup       string `gorm:"size:8;not null;index" json:"age_group"`
	Rank           int    `gorm:"not null" json:"rank"`             // 本场名次
	Points         int    `gorm:"not null" json:"points"`           // 本场获得积分（后台计算后写入）
	Award          string `gorm:"size:64" json:"award"`             // 冠军/亚军/季军/菁英舞者/特金奖 等
	ParticipantNum int    `gorm:"default:0" json:"participant_num"` // 该榜单本场参赛人数（便于后续追溯）

	// 唯一：同一赛事同一选手同一舞种只能有一条成绩
	// uniqueIndex: uk_activity_player_dance

	Player   *Player   `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Activity *Activity `gorm:"foreignKey:ActivityID" json:"activity,omitempty"`
}

func (ActivityResult) TableName() string { return "activity_results" }
