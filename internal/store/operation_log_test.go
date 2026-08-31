package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
)

func newTestOperationLogRepo(t *testing.T) OperationLogRepo {
	t.Helper()
	return NewOperationLogRepository(newTestDB(t, &model.OperationLog{}))
}

func TestOperationLogRepoCRUD(t *testing.T) {
	repo := newTestOperationLogRepo(t)

	log1 := &model.OperationLog{UserID: 1, Username: "admin", Module: "user", Action: "create"}
	require.NoError(t, repo.Create(context.Background(), log1))
	require.NotZero(t, log1.ID)

	got, err := repo.GetByID(context.Background(), log1.ID)
	require.NoError(t, err)
	assert.Equal(t, "admin", got.Username)

	require.NoError(t, repo.Create(context.Background(), &model.OperationLog{UserID: 2, Username: "zhangsan", Module: "article", Action: "update"}))

	// List（模块过滤 + 排序）
	list, total, err := repo.List(context.Background(), 1, 10, func(db *gorm.DB) *gorm.DB {
		return db.Where("module = ?", "user").Order("created_at DESC, id DESC")
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, "admin", list[0].Username)

	// OperationLog 无 DeletedAt，Delete 为硬删除
	require.NoError(t, repo.Delete(context.Background(), log1.ID))
	_, err = repo.GetByID(context.Background(), log1.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
