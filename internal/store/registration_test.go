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

func newTestRegistrationRepo(t *testing.T) (RegistrationRepo, *gorm.DB) {
	t.Helper()
	db := newTestDB(t, &model.Season{}, &model.Activity{}, &model.Registration{})
	return NewRegistrationRepository(db), db
}

func makeSeason(t *testing.T, db *gorm.DB) *model.Season {
	t.Helper()
	s := &model.Season{Year: 2026, Name: "2026赛季", Status: "active"}
	require.NoError(t, db.Create(s).Error)
	return s
}

func makeActivity(t *testing.T, db *gorm.DB, seasonID int64, title string, status int) *model.Activity {
	t.Helper()
	a := &model.Activity{Title: title, Status: status, SeasonID: seasonID}
	require.NoError(t, db.Create(a).Error)
	return a
}

func makeRegistration(t *testing.T, db *gorm.DB, activityID, userID int64, status int) *model.Registration {
	t.Helper()
	r := &model.Registration{
		ActivityID: activityID,
		UserID:     userID,
		Name:       "报名者",
		Phone:      "13800000000",
		Status:     status,
	}
	require.NoError(t, db.Create(r).Error)
	return r
}

func TestRegistrationRepoListWithActivity(t *testing.T) {
	repo, db := newTestRegistrationRepo(t)
	season := makeSeason(t, db)
	act := makeActivity(t, db, season.ID, "上海站超级赛", 1)

	makeRegistration(t, db, act.ID, 1, 1) // 已支付
	makeRegistration(t, db, act.ID, 2, 0) // 待支付
	makeRegistration(t, db, act.ID, 3, 2) // 已取消

	// 按活动过滤（不过滤状态）
	list, total, err := repo.ListWithActivity(context.Background(), 1, 10, act.ID, -1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, list, 3)
	assert.Equal(t, "上海站超级赛", list[0].ActivityTitle)

	// 按状态过滤（不过滤活动）
	list, total, err = repo.ListWithActivity(context.Background(), 1, 10, 0, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, int64(1), list[0].UserID)
}

func TestRegistrationRepoGetWithActivity(t *testing.T) {
	repo, db := newTestRegistrationRepo(t)
	season := makeSeason(t, db)
	act := makeActivity(t, db, season.ID, "上海站超级赛", 1)
	reg := makeRegistration(t, db, act.ID, 1, 1)

	got, err := repo.GetWithActivity(context.Background(), reg.ID)
	require.NoError(t, err)
	assert.Equal(t, reg.ID, got.ID)
	assert.Equal(t, "上海站超级赛", got.ActivityTitle)

	// 不存在
	_, err = repo.GetWithActivity(context.Background(), 9999)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// 软删除后查询不到
	require.NoError(t, repo.Delete(context.Background(), reg.ID))
	_, err = repo.GetWithActivity(context.Background(), reg.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestRegistrationRepoListByUserWithActivity(t *testing.T) {
	repo, db := newTestRegistrationRepo(t)
	season := makeSeason(t, db)
	act1 := makeActivity(t, db, season.ID, "活动A", 1)
	act2 := makeActivity(t, db, season.ID, "活动B", 1)

	makeRegistration(t, db, act1.ID, 100, 1)
	makeRegistration(t, db, act2.ID, 100, 0)
	makeRegistration(t, db, act1.ID, 200, 1) // 别的用户

	list, total, err := repo.ListByUserWithActivity(context.Background(), 100, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	titles := make(map[string]bool, len(list))
	for _, it := range list {
		titles[it.ActivityTitle] = true
	}
	assert.True(t, titles["活动A"])
	assert.True(t, titles["活动B"])
}
