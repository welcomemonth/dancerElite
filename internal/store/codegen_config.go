package store

import (
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/repository"
)

// CodegenConfigRepo 代码生成配置数据访问接口（暂无自定义查询，仅通用 CRUD）
type CodegenConfigRepo interface {
	CRUD[model.CodegenConfig]
}

// NewCodegenConfigRepository 创建 CodegenConfigRepo，纯 CRUD 直接复用 BaseRepo
func NewCodegenConfigRepository(db *gorm.DB) CodegenConfigRepo {
	return repository.NewBaseRepo[model.CodegenConfig](db)
}
