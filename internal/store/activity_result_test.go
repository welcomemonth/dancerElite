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

func newTestActivityResultRepo(t *testing.T) ActivityResultRepo {
	t.Helper()
	return NewActivityResultRepository(newTestDB(t, &model.ActivityResult{}))
}

func TestActivityResultRepoCRUD(t *testing.T) {
	repo := newTestActivityResultRepo(t)

	r := &model.ActivityResult{
		ActivityID: 1,
		SeasonID:   1,
		PlayerID:   1,
		DanceType:  "latino",
		AgeGroup:   "U15",
		Rank:       1,
		Points:     10,
		Award:      "冠军",
	}
	require.NoError(t, repo.Create(context.Background(), r))
	require.NotZero(t, r.ID)

	got, err := repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, "latino", got.DanceType)
	assert.Equal(t, 10, got.Points)

	require.NoError(t, repo.Update(context.Background(), r.ID, map[string]any{"rank": 2}))
	got, err = repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, got.Rank)

	require.NoError(t, repo.Create(context.Background(), &model.ActivityResult{
		ActivityID: 1, SeasonID: 1, PlayerID: 2, DanceType: "latino", AgeGroup: "U15", Rank: 3, Points: 5,
	}))
	list, total, err := repo.List(context.Background(), 1, 10, func(db *gorm.DB) *gorm.DB { return db.Order("rank") })
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	require.NoError(t, repo.Delete(context.Background(), r.ID))
	_, err = repo.GetByID(context.Background(), r.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
