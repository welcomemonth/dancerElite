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

func newTestSeasonRepo(t *testing.T) SeasonRepo {
	t.Helper()
	return NewSeasonRepository(newTestDB(t, &model.Season{}))
}

func createSeason(t *testing.T, repo SeasonRepo, year int, name, status string) *model.Season {
	t.Helper()
	s := &model.Season{Year: year, Name: name, Status: status}
	require.NoError(t, repo.Create(context.Background(), s))
	return s
}

func TestSeasonRepoCRUD(t *testing.T) {
	repo := newTestSeasonRepo(t)

	s := createSeason(t, repo, 2026, "2026赛季", "active")
	require.NotZero(t, s.ID, "创建后应自动生成主键")

	// GetByID
	got, err := repo.GetByID(context.Background(), s.ID)
	require.NoError(t, err)
	assert.Equal(t, "2026赛季", got.Name)
	assert.Equal(t, 2026, got.Year)

	// Update
	require.NoError(t, repo.Update(context.Background(), s.ID, map[string]any{"status": "archived"}))
	got, err = repo.GetByID(context.Background(), s.ID)
	require.NoError(t, err)
	assert.Equal(t, "archived", got.Status)

	// List
	list, total, err := repo.List(context.Background(), 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)

	// Delete（软删除）
	require.NoError(t, repo.Delete(context.Background(), s.ID))
	_, err = repo.GetByID(context.Background(), s.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestSeasonRepoGetByYear(t *testing.T) {
	repo := newTestSeasonRepo(t)
	createSeason(t, repo, 2026, "2026赛季", "active")

	got, err := repo.GetByYear(context.Background(), 2026)
	require.NoError(t, err)
	assert.Equal(t, "2026赛季", got.Name)

	_, err = repo.GetByYear(context.Background(), 2025)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestSeasonRepoGetActive(t *testing.T) {
	repo := newTestSeasonRepo(t)
	createSeason(t, repo, 2025, "2025赛季", "archived")
	active := createSeason(t, repo, 2026, "2026赛季", "active")

	got, err := repo.GetActive(context.Background())
	require.NoError(t, err)
	assert.Equal(t, active.ID, got.ID)
	assert.Equal(t, "2026赛季", got.Name)
}

func TestSeasonRepoGetActiveNotFound(t *testing.T) {
	repo := newTestSeasonRepo(t)
	createSeason(t, repo, 2025, "2025赛季", "archived")

	_, err := repo.GetActive(context.Background())
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
