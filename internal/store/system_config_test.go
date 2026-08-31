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

func newTestSystemConfigRepo(t *testing.T) SystemConfigRepo {
	t.Helper()
	return NewSystemConfigRepository(newTestDB(t, &model.SystemConfig{}))
}

func TestSystemConfigRepoUpsert(t *testing.T) {
	repo := newTestSystemConfigRepo(t)

	// 插入新配置
	cfg := &model.SystemConfig{Key: "site.name", Value: "远山", Type: "string", GroupName: "站点信息"}
	require.NoError(t, repo.Upsert(context.Background(), cfg))

	got, err := repo.GetByKey(context.Background(), "site.name")
	require.NoError(t, err)
	assert.Equal(t, "远山", got.Value)

	// 同 key 再次 Upsert -> 更新而非新增
	require.NoError(t, repo.Upsert(context.Background(), &model.SystemConfig{Key: "site.name", Value: "远山公益", Type: "string", GroupName: "站点信息"}))
	got, err = repo.GetByKey(context.Background(), "site.name")
	require.NoError(t, err)
	assert.Equal(t, "远山公益", got.Value)

	all, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 1)
}

func TestSystemConfigRepoGetByKeyNotFound(t *testing.T) {
	repo := newTestSystemConfigRepo(t)

	_, err := repo.GetByKey(context.Background(), "not.exist")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestSystemConfigRepoBatchUpsert(t *testing.T) {
	repo := newTestSystemConfigRepo(t)

	require.NoError(t, repo.BatchUpsert(context.Background(), []model.SystemConfig{
		{Key: "site.name", Value: "远山", GroupName: "站点信息"},
		{Key: "wechat.app_id", Value: "wx123", GroupName: "微信配置"},
	}))

	// 部分更新
	require.NoError(t, repo.BatchUpsert(context.Background(), []model.SystemConfig{
		{Key: "site.name", Value: "远山公益", GroupName: "站点信息"},
	}))

	got, err := repo.GetByKey(context.Background(), "site.name")
	require.NoError(t, err)
	assert.Equal(t, "远山公益", got.Value)

	all, err := repo.FindAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestSystemConfigRepoListGroups(t *testing.T) {
	repo := newTestSystemConfigRepo(t)

	require.NoError(t, repo.BatchUpsert(context.Background(), []model.SystemConfig{
		{Key: "site.name", Value: "远山", GroupName: "站点信息"},
		{Key: "site.description", Value: "说明", GroupName: "站点信息"},
		{Key: "wechat.app_id", Value: "wx", GroupName: "微信配置"},
		{Key: "no.group", Value: "x", GroupName: ""},
	}))

	groups, err := repo.ListGroups(context.Background())
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"站点信息", "微信配置"}, groups)
}

func TestSystemConfigRepoDeleteByKey(t *testing.T) {
	repo := newTestSystemConfigRepo(t)

	require.NoError(t, repo.Upsert(context.Background(), &model.SystemConfig{Key: "site.name", Value: "远山", GroupName: "站点信息"}))
	require.NoError(t, repo.DeleteByKey(context.Background(), "site.name"))

	_, err := repo.GetByKey(context.Background(), "site.name")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
