package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
)

func newTestColumnRepo(t *testing.T) ColumnRepo {
	t.Helper()
	return NewColumnRepository(newTestDB(t, &model.Column{}))
}

func TestColumnRepoCRUD(t *testing.T) {
	repo := newTestColumnRepo(t)

	c := &model.Column{Name: "赛事新闻", SortOrder: 1}
	require.NoError(t, repo.Create(context.Background(), c))
	require.NotZero(t, c.ID)

	got, err := repo.GetByID(context.Background(), c.ID)
	require.NoError(t, err)
	assert.Equal(t, "赛事新闻", got.Name)

	require.NoError(t, repo.Update(context.Background(), c.ID, map[string]any{"sort_order": 2}))
	got, err = repo.GetByID(context.Background(), c.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, got.SortOrder)

	// FindAll 带排序 scope
	require.NoError(t, repo.Create(context.Background(), &model.Column{Name: "通知公告", SortOrder: 1}))
	all, err := repo.FindAll(context.Background(), func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order, id")
	})
	require.NoError(t, err)
	require.Len(t, all, 2)
	assert.Equal(t, "通知公告", all[0].Name) // sort_order=1 排在前面

	require.NoError(t, repo.Delete(context.Background(), c.ID))
	_, err = repo.GetByID(context.Background(), c.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
