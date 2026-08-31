package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/welcomemonth/dancer-elite/internal/model"
)

func newTestRoleMenuRepo(t *testing.T) RoleMenuRepo {
	t.Helper()
	return NewRoleMenuRepository(newTestDB(t, &model.RoleMenu{}))
}

func TestRoleMenuRepoListMenuIDsByRole(t *testing.T) {
	repo := newTestRoleMenuRepo(t)
	require.NoError(t, repo.ReplaceByRole(context.Background(), 1, []int64{11, 12, 13}))

	ids, err := repo.ListMenuIDsByRole(context.Background(), 1)
	require.NoError(t, err)
	assert.ElementsMatch(t, []int64{11, 12, 13}, ids)

	ids, err = repo.ListMenuIDsByRole(context.Background(), 2)
	require.NoError(t, err)
	assert.Empty(t, ids)
}

func TestRoleMenuRepoReplaceByRole(t *testing.T) {
	repo := newTestRoleMenuRepo(t)
	require.NoError(t, repo.ReplaceByRole(context.Background(), 1, []int64{11, 12}))

	// 替换为新的菜单集合
	require.NoError(t, repo.ReplaceByRole(context.Background(), 1, []int64{13, 14}))

	ids, err := repo.ListMenuIDsByRole(context.Background(), 1)
	require.NoError(t, err)
	assert.ElementsMatch(t, []int64{13, 14}, ids)
}

func TestRoleMenuRepoDeleteByRole(t *testing.T) {
	repo := newTestRoleMenuRepo(t)
	require.NoError(t, repo.ReplaceByRole(context.Background(), 1, []int64{11, 12}))
	require.NoError(t, repo.ReplaceByRole(context.Background(), 2, []int64{11}))

	require.NoError(t, repo.DeleteByRole(context.Background(), 1))

	ids, err := repo.ListMenuIDsByRole(context.Background(), 1)
	require.NoError(t, err)
	assert.Empty(t, ids)

	// 其他角色不受影响
	ids, err = repo.ListMenuIDsByRole(context.Background(), 2)
	require.NoError(t, err)
	assert.ElementsMatch(t, []int64{11}, ids)
}

func TestRoleMenuRepoDeleteByMenu(t *testing.T) {
	repo := newTestRoleMenuRepo(t)
	require.NoError(t, repo.ReplaceByRole(context.Background(), 1, []int64{11, 12}))
	require.NoError(t, repo.ReplaceByRole(context.Background(), 2, []int64{11}))

	require.NoError(t, repo.DeleteByMenu(context.Background(), 11))

	ids, err := repo.ListMenuIDsByRole(context.Background(), 1)
	require.NoError(t, err)
	assert.ElementsMatch(t, []int64{12}, ids)

	ids, err = repo.ListMenuIDsByRole(context.Background(), 2)
	require.NoError(t, err)
	assert.Empty(t, ids)
}
