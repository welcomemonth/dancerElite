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

	"github.com/go-playground/validator/v10"
	"github.com/zzhtl/go-mountain/internal/config"
)

var validate = validator.New()

type Store struct {
	cfg          *config.Config
	UserRepo     UserRepo
	PlayerRepo   PlayerRepo
	ActivityRepo ActivityRepo
	SeasonRepo   SeasonRepo
}

func New(cfg *config.Config) (*Store, error) {
	db, err := newDB(cfg.Database)
	if err != nil {
		return nil, err
	}
	return &Store{
		UserRepo:     NewUserRepository(db),
		PlayerRepo:   NewPlayerRepository(db),
		ActivityRepo: NewActivityRepository(db),
		SeasonRepo:   NewSeasonRepository(db),
		cfg:          cfg,
	}, nil
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
