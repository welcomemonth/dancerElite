package store

import (
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// AnnualRankingRepo 年度积分总榜数据访问接口（暂无自定义查询，仅通用 CRUD）
type AnnualRankingRepo interface {
	CRUD[model.AnnualRanking]
}

// NewAnnualRankingRepository 创建 AnnualRankingRepo，纯 CRUD 直接复用 BaseRepo
func NewAnnualRankingRepository(db *gorm.DB) AnnualRankingRepo {
	return repository.NewBaseRepo[model.AnnualRanking](db)
}
