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

func newTestActivityRepo(t *testing.T) (ActivityRepo, *gorm.DB) {
	t.Helper()
	db := newTestDB(t, &model.Season{}, &model.Activity{}, &model.Registration{})
	return NewActivityRepository(db), db
}

func seedSeason(t *testing.T, db *gorm.DB) *model.Season {
	t.Helper()
	s := &model.Season{Year: 2026, Name: "2026赛季", Status: "active"}
	require.NoError(t, db.Create(s).Error)
	return s
}

func seedActivity(t *testing.T, repo ActivityRepo, seasonID int64, title string, status int) *model.Activity {
	t.Helper()
	a := &model.Activity{Title: title, Status: status, SeasonID: seasonID}
	require.NoError(t, repo.Create(context.Background(), a))
	return a
}

func seedRegistration(t *testing.T, db *gorm.DB, activityID int64, status int) {
	t.Helper()
	r := &model.Registration{
		ActivityID: activityID,
		UserID:     activityID, // 仅占位，满足 not null
		Name:       "报名者",
		Phone:      "13800000000",
		Status:     status,
	}
	require.NoError(t, db.Create(r).Error)
}

func TestActivityRepoGetWithRegCount(t *testing.T) {
	repo, db := newTestActivityRepo(t)
	season := seedSeason(t, db)
	a := seedActivity(t, repo, season.ID, "上海站超级赛", 1)

	// status 0/1 计入报名人数，status 2（已取消）不计入
	seedRegistration(t, db, a.ID, 0)
	seedRegistration(t, db, a.ID, 1)
	seedRegistration(t, db, a.ID, 2)

	got, err := repo.GetWithRegCount(context.Background(), a.ID)
	require.NoError(t, err)
	assert.Equal(t, "上海站超级赛", got.Title)
	assert.Equal(t, season.ID, got.SeasonID)
	assert.Equal(t, int64(2), got.RegCount, "只统计待支付+已支付")
}

func TestActivityRepoListWithRegCount(t *testing.T) {
	repo, db := newTestActivityRepo(t)
	season := seedSeason(t, db)

	a1 := seedActivity(t, repo, season.ID, "活动1", 1)
	seedRegistration(t, db, a1.ID, 1)
	a2 := seedActivity(t, repo, season.ID, "活动2", 1)
	seedRegistration(t, db, a2.ID, 0)
	seedRegistration(t, db, a2.ID, 1)
	seedActivity(t, repo, season.ID, "活动3", 2) // 报名截止，会被 status 过滤

	list, total, err := repo.ListWithRegCount(context.Background(), 1, 10, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	byTitle := make(map[string]int64, len(list))
	for _, it := range list {
		byTitle[it.Title] = it.RegCount
	}
	assert.Equal(t, int64(1), byTitle["活动1"])
	assert.Equal(t, int64(2), byTitle["活动2"])
}

func TestActivityRepoListPublishedWithRegCount(t *testing.T) {
	repo, db := newTestActivityRepo(t)
	season := seedSeason(t, db)

	seedActivity(t, repo, season.ID, "草稿", 0)  // 排除
	seedActivity(t, repo, season.ID, "报名中", 1) // 包含
	seedActivity(t, repo, season.ID, "已结束", 4) // 包含

	list, total, err := repo.ListPublishedWithRegCount(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total, "草稿不应出现在小程序列表")
	require.Len(t, list, 2)
}

func TestActivityRepoGetPublishedWithRegCount(t *testing.T) {
	repo, db := newTestActivityRepo(t)
	season := seedSeason(t, db)

	draft := seedActivity(t, repo, season.ID, "草稿", 0)    // 排除
	ended := seedActivity(t, repo, season.ID, "已结束", 4)   // 排除（4 不在 1/2/3）
	ongoing := seedActivity(t, repo, season.ID, "进行中", 3) // 包含

	_, err := repo.GetPublishedWithRegCount(context.Background(), draft.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	_, err = repo.GetPublishedWithRegCount(context.Background(), ended.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	got, err := repo.GetPublishedWithRegCount(context.Background(), ongoing.ID)
	require.NoError(t, err)
	assert.Equal(t, "进行中", got.Title)
}
