package store

import (
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// ActivityResultRepo 赛事成绩数据访问接口（暂无自定义查询，仅通用 CRUD）
type ActivityResultRepo interface {
	CRUD[model.ActivityResult]
}

// NewActivityResultRepository 创建 ActivityResultRepo，纯 CRUD 直接复用 BaseRepo
func NewActivityResultRepository(db *gorm.DB) ActivityResultRepo {
	return repository.NewBaseRepo[model.ActivityResult](db)
}
