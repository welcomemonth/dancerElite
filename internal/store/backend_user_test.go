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

func newTestBackendUserRepo(t *testing.T) (BackendUserRepo, *gorm.DB) {
	t.Helper()
	db := newTestDB(t, &model.BackendUser{}, &model.Role{})
	return NewBackendUserRepository(db), db
}

func seedRole(t *testing.T, db *gorm.DB, name string) *model.Role {
	t.Helper()
	r := &model.Role{Name: name, DisplayName: name + "名", Status: 1}
	require.NoError(t, db.Create(r).Error)
	return r
}

func seedBackendUser(t *testing.T, repo BackendUserRepo, username string, roleID int64) *model.BackendUser {
	t.Helper()
	u := &model.BackendUser{
		Username: username,
		Email:    username + "@test.com",
		Password: "hash",
		RoleID:   roleID,
		Status:   1,
	}
	require.NoError(t, repo.Create(context.Background(), u))
	return u
}

func TestBackendUserRepoGetByUsername(t *testing.T) {
	repo, db := newTestBackendUserRepo(t)
	role := seedRole(t, db, "admin")
	seedBackendUser(t, repo, "zhang", role.ID)

	got, err := repo.GetByUsername(context.Background(), "zhang")
	require.NoError(t, err)
	assert.Equal(t, "zhang", got.Username)
	require.NotNil(t, got.Role, "应预加载角色")
	assert.Equal(t, "admin", got.Role.Name)

	_, err = repo.GetByUsername(context.Background(), "nobody")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestBackendUserRepoGetWithRole(t *testing.T) {
	repo, db := newTestBackendUserRepo(t)
	role := seedRole(t, db, "editor")
	u := seedBackendUser(t, repo, "li", role.ID)

	item, err := repo.GetWithRole(context.Background(), u.ID)
	require.NoError(t, err)
	assert.Equal(t, "li", item.Username)
	assert.Equal(t, "editor", item.RoleName)
	assert.Equal(t, "editor名", item.RoleDisplay)
}

func TestBackendUserRepoListWithRole(t *testing.T) {
	repo, db := newTestBackendUserRepo(t)
	admin := seedRole(t, db, "admin")
	editor := seedRole(t, db, "editor")
	seedBackendUser(t, repo, "a-user", admin.ID)
	seedBackendUser(t, repo, "b-user", editor.ID)

	list, total, err := repo.ListWithRole(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	// 角色信息应通过 join 填充
	byName := make(map[string]string, len(list))
	for _, it := range list {
		byName[it.Username] = it.RoleName
	}
	assert.Equal(t, "admin", byName["a-user"])
	assert.Equal(t, "editor", byName["b-user"])
}
