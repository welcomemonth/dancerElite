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

func newTestRoleRepo(t *testing.T) RoleRepo {
	t.Helper()
	return NewRoleRepository(newTestDB(t, &model.Role{}, &model.Menu{}, &model.RoleMenu{}))
}

func TestRoleRepoCRUD(t *testing.T) {
	repo := newTestRoleRepo(t)

	r := &model.Role{Name: "admin", DisplayName: "超级管理员", Status: 1}
	require.NoError(t, repo.Create(context.Background(), r))
	require.NotZero(t, r.ID)

	got, err := repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, "超级管理员", got.DisplayName)

	require.NoError(t, repo.Update(context.Background(), r.ID, map[string]any{"status": 0}))
	got, err = repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, got.Status)

	// List（带 id 排序 scope）
	require.NoError(t, repo.Create(context.Background(), &model.Role{Name: "editor", DisplayName: "编辑员"}))
	list, total, err := repo.List(context.Background(), 1, 10, func(db *gorm.DB) *gorm.DB { return db.Order("id") })
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	require.NoError(t, repo.Delete(context.Background(), r.ID))
	_, err = repo.GetByID(context.Background(), r.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
