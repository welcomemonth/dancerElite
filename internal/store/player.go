package store

import (
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// PlayerRepo 选手数据访问接口（暂无自定义查询，仅通用 CRUD）
type PlayerRepo interface {
	CRUD[model.Player]
}

// NewPlayerRepository 创建 PlayerRepo，纯 CRUD 直接复用 BaseRepo
func NewPlayerRepository(db *gorm.DB) PlayerRepo {
	return repository.NewBaseRepo[model.Player](db)
}
