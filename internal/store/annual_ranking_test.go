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

func newTestAnnualRankingRepo(t *testing.T) AnnualRankingRepo {
	t.Helper()
	return NewAnnualRankingRepository(newTestDB(t, &model.AnnualRanking{}))
}

func TestAnnualRankingRepoCRUD(t *testing.T) {
	repo := newTestAnnualRankingRepo(t)

	r := &model.AnnualRanking{
		SeasonID:    1,
		AgeGroup:    "U15",
		DanceType:   "latino",
		PlayerID:    1,
		TotalPoints: 25,
		Rank:        1,
	}
	require.NoError(t, repo.Create(context.Background(), r))
	require.NotZero(t, r.ID)

	got, err := repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, "latino", got.DanceType)
	assert.Equal(t, 25, got.TotalPoints)

	require.NoError(t, repo.Update(context.Background(), r.ID, map[string]any{"rank": 2}))
	got, err = repo.GetByID(context.Background(), r.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, got.Rank)

	require.NoError(t, repo.Create(context.Background(), &model.AnnualRanking{
		SeasonID: 1, AgeGroup: "U15", DanceType: "latino", PlayerID: 2, TotalPoints: 15, Rank: 3,
	}))
	list, total, err := repo.List(context.Background(), 1, 10, func(db *gorm.DB) *gorm.DB { return db.Order("rank") })
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)

	require.NoError(t, repo.Delete(context.Background(), r.ID))
	_, err = repo.GetByID(context.Background(), r.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
