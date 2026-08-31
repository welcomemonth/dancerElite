package store

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
)

func newTestCodegenConfigRepo(t *testing.T) CodegenConfigRepo {
	t.Helper()
	return NewCodegenConfigRepository(newTestDB(t, &model.CodegenConfig{}))
}

func TestCodegenConfigRepoCRUD(t *testing.T) {
	repo := newTestCodegenConfigRepo(t)

	cfg := &model.CodegenConfig{
		TblName:       "articles",
		ModuleName:    "article",
		DisplayName:   "文章管理",
		ColumnsConfig: json.RawMessage(`[]`),
	}
	require.NoError(t, repo.Create(context.Background(), cfg))
	require.NotZero(t, cfg.ID)

	got, err := repo.GetByID(context.Background(), cfg.ID)
	require.NoError(t, err)
	assert.Equal(t, "articles", got.TblName)

	require.NoError(t, repo.Update(context.Background(), cfg.ID, map[string]any{"generated": true}))
	got, err = repo.GetByID(context.Background(), cfg.ID)
	require.NoError(t, err)
	assert.True(t, got.Generated)

	// List 按 created_at 倒序
	require.NoError(t, repo.Create(context.Background(), &model.CodegenConfig{
		TblName: "users", ModuleName: "user", DisplayName: "用户", ColumnsConfig: json.RawMessage(`[]`),
	}))
	list, total, err := repo.List(context.Background(), 1, 10, func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") })
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	require.NoError(t, repo.Delete(context.Background(), cfg.ID))
	_, err = repo.GetByID(context.Background(), cfg.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
