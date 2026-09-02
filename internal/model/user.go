package model

// User 小程序用户
type User struct {
	BaseModel
	OpenID  string `gorm:"type:text;uniqueIndex" json:"openid"`
	UnionID string `gorm:"type:text" json:"union_id"`
	Phone   string `gorm:"type:text" json:"phone"`
	Name    string `gorm:"type:text" json:"name"`
	Avatar  string `gorm:"type:text" json:"avatar"`
	Gender  int    `gorm:"default:0" json:"gender"`
	Status  int    `gorm:"default:1" json:"status"`

	// 关联
	Player *Player `gorm:"foreignKey:UserID" json:"player,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type Player struct {
	BaseModel
	UserID      *int64 `gorm:"uniqueIndex" json:"user_id"` // 1:1 绑定
	RealName    string `gorm:"size:32;not null" json:"real_name"`
	Gender      string `gorm:"size:16;not null" json:"gender"`  // male / female
	IDCard      string `gorm:"size:32;not null" json:"id_card"` // 脱敏后，如 110101********1234
	IDCardHash  string `gorm:"size:64;index" json:"-"`          // 可选，用于唯一性校验（SHA256）
	Phone       string `gorm:"size:20;not null" json:"phone"`   // 联系人手机号
	Institution string `gorm:"size:128" json:"institution"`     // 所属机构/学校，可空
	Teacher     string `gorm:"size:64" json:"teacher"`          // 指导老师，可空
	BirthYear   int    `gorm:"not null" json:"birth_year"`      // 从身份证解析，方便年龄组计算
	BirthMonth  int    `gorm:"not null" json:"birth_month"`
	BirthDay    int    `gorm:"not null" json:"birth_day"`
	AgeGroup    string `gorm:"size:8;index" json:"age_group"` // 缓存当前有效年龄组（U11/U13/U15）

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (Player) TableName() string { return "players" }
