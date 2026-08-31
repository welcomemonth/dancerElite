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

func newTestPlayerRepo(t *testing.T) PlayerRepo {
	t.Helper()
	return NewPlayerRepository(newTestDB(t, &model.Player{}))
}

func TestPlayerRepoCRUD(t *testing.T) {
	repo := newTestPlayerRepo(t)

	p := &model.Player{
		UserID:     1,
		RealName:   "张三",
		Gender:     "male",
		IDCard:     "110101********1234",
		Phone:      "13800000000",
		BirthYear:  2010,
		BirthMonth: 5,
		BirthDay:   1,
		AgeGroup:   "U15",
	}
	require.NoError(t, repo.Create(context.Background(), p))
	require.NotZero(t, p.ID, "创建后应自动生成主键")

	// GetByID
	got, err := repo.GetByID(context.Background(), p.ID)
	require.NoError(t, err)
	assert.Equal(t, "张三", got.RealName)
	assert.Equal(t, "U15", got.AgeGroup)

	// Update
	require.NoError(t, repo.Update(context.Background(), p.ID, map[string]any{"institution": "远山学校"}))
	got, err = repo.GetByID(context.Background(), p.ID)
	require.NoError(t, err)
	assert.Equal(t, "远山学校", got.Institution)

	// List
	list, total, err := repo.List(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	// Delete（软删除）
	require.NoError(t, repo.Delete(context.Background(), p.ID))
	_, err = repo.GetByID(context.Background(), p.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
