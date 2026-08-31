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

func newTestMenuRepo(t *testing.T) MenuRepo {
	t.Helper()
	return NewMenuRepository(newTestDB(t, &model.Menu{}))
}

func TestMenuRepoCRUD(t *testing.T) {
	repo := newTestMenuRepo(t)

	m := &model.Menu{Name: "content", Title: "内容管理", Sort: 1, Type: 1, Status: 1}
	require.NoError(t, repo.Create(context.Background(), m))
	require.NotZero(t, m.ID)

	got, err := repo.GetByID(context.Background(), m.ID)
	require.NoError(t, err)
	assert.Equal(t, "内容管理", got.Title)

	require.NoError(t, repo.Update(context.Background(), m.ID, map[string]any{"sort": 2}))
	got, err = repo.GetByID(context.Background(), m.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, got.Sort)

	// FindAll（带排序 scope）
	require.NoError(t, repo.Create(context.Background(), &model.Menu{Name: "articles", Title: "文章管理", ParentID: m.ID, Sort: 1, Type: 2}))
	all, err := repo.FindAll(context.Background(), func(db *gorm.DB) *gorm.DB { return db.Order("sort, id") })
	require.NoError(t, err)
	assert.Len(t, all, 2)

	require.NoError(t, repo.Delete(context.Background(), m.ID))
	_, err = repo.GetByID(context.Background(), m.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
