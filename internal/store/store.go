package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/welcomemonth/dancer-elite/internal/config"
	"github.com/welcomemonth/dancer-elite/internal/model"
)

type Store struct {
	db                 *gorm.DB
	cfg                *config.Config
	UserRepo           UserRepo
	BackendUserRepo    BackendUserRepo
	PlayerRepo         PlayerRepo
	ActivityRepo       ActivityRepo
	RegistrationRepo   RegistrationRepo
	PaymentRepo        PaymentRepo
	SeasonRepo         SeasonRepo
	ArticleRepo        ArticleRepo
	ColumnRepo         ColumnRepo
	RoleRepo           RoleRepo
	MenuRepo           MenuRepo
	RoleMenuRepo       RoleMenuRepo
	SystemConfigRepo   SystemConfigRepo
	CodegenConfigRepo  CodegenConfigRepo
	OperationLogRepo   OperationLogRepo
	ActivityResultRepo ActivityResultRepo
	AnnualRankingRepo  AnnualRankingRepo
}

func New(cfg *config.Config) (*Store, error) {
	db, err := newDB(cfg.Database)
	if err != nil {
		return nil, err
	}
	// 开发模式则自动迁移表结构
	if cfg.Server.Debug {
		if err := db.AutoMigrate(
			&model.BackendUser{},
			&model.Role{},
			&model.Menu{},
			&model.RoleMenu{},
			&model.Article{},
			&model.Column{},
			&model.User{},
			&model.Activity{},
			&model.Registration{},
			&model.Payment{},
			&model.CodegenConfig{},
			&model.OperationLog{},
			&model.SystemConfig{},

			&model.Player{},
			&model.Season{},
			&model.ActivityResult{},
			&model.AnnualRanking{},
		); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
	}

	return NewWithDB(db, cfg), nil
}

// NewWithDB 用已有的 GORM 连接装配 Store（表迁移由调用方控制，便于测试注入内存库）。
func NewWithDB(db *gorm.DB, cfg *config.Config) *Store {
	return &Store{
		db:                 db,
		UserRepo:           NewUserRepository(db),
		BackendUserRepo:    NewBackendUserRepository(db),
		PlayerRepo:         NewPlayerRepository(db),
		ActivityRepo:       NewActivityRepository(db),
		RegistrationRepo:   NewRegistrationRepository(db),
		PaymentRepo:        NewPaymentRepository(db),
		SeasonRepo:         NewSeasonRepository(db),
		ArticleRepo:        NewArticleRepository(db),
		ColumnRepo:         NewColumnRepository(db),
		RoleRepo:           NewRoleRepository(db),
		MenuRepo:           NewMenuRepository(db),
		RoleMenuRepo:       NewRoleMenuRepository(db),
		SystemConfigRepo:   NewSystemConfigRepository(db),
		CodegenConfigRepo:  NewCodegenConfigRepository(db),
		OperationLogRepo:   NewOperationLogRepository(db),
		ActivityResultRepo: NewActivityResultRepository(db),
		AnnualRankingRepo:  NewAnnualRankingRepository(db),
		cfg:                cfg,
	}
}

// Init 根据配置初始化 GORM 数据库连接
func newDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "sqlite3", "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	case "postgres", "postgresql":
		dialector = postgres.Open(cfg.DSN)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库连接检测失败: %w", err)
	}

	return db, nil
}

// DB 返回底层 GORM 连接，供事务、库内省等无法走 repo 的场景使用
func (s *Store) DB() *gorm.DB { return s.db }

// Close 关闭底层数据库连接
func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
