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

func newTestArticleRepo(t *testing.T) (ArticleRepo, *gorm.DB) {
	t.Helper()
	db := newTestDB(t, &model.Column{}, &model.Article{})
	return NewArticleRepository(db), db
}

func makeColumn(t *testing.T, db *gorm.DB, name string) *model.Column {
	t.Helper()
	c := &model.Column{Name: name}
	require.NoError(t, db.Create(c).Error)
	return c
}

func makeArticle(t *testing.T, db *gorm.DB, columnID int64, title string, status int) *model.Article {
	t.Helper()
	a := &model.Article{ColumnID: columnID, Title: title, Status: status}
	require.NoError(t, db.Create(a).Error)
	return a
}

func TestArticleRepoListWithColumn(t *testing.T) {
	repo, db := newTestArticleRepo(t)
	col := makeColumn(t, db, "赛事新闻")
	makeArticle(t, db, col.ID, "文章1", 1)
	makeArticle(t, db, col.ID, "文章2", 0)
	makeArticle(t, db, col.ID, "文章3", 1)

	// 按状态过滤（已发布）
	list, total, err := repo.ListWithColumn(context.Background(), 1, 10, 0, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, list, 2)
	assert.Equal(t, "赛事新闻", list[0].ColumnName)

	// 按栏目过滤（不过滤状态）
	list, total, err = repo.ListWithColumn(context.Background(), 1, 10, col.ID, -1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
}

func TestArticleRepoGetWithColumn(t *testing.T) {
	repo, db := newTestArticleRepo(t)
	col := makeColumn(t, db, "赛事新闻")
	a := makeArticle(t, db, col.ID, "某文章", 1)

	got, err := repo.GetWithColumn(context.Background(), a.ID)
	require.NoError(t, err)
	assert.Equal(t, "某文章", got.Title)
	assert.Equal(t, "赛事新闻", got.ColumnName)

	_, err = repo.GetWithColumn(context.Background(), 9999)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestArticleRepoListByColumn(t *testing.T) {
	repo, db := newTestArticleRepo(t)
	col1 := makeColumn(t, db, "栏目1")
	col2 := makeColumn(t, db, "栏目2")
	makeArticle(t, db, col1.ID, "已发布", 1)
	makeArticle(t, db, col1.ID, "草稿", 0)
	makeArticle(t, db, col2.ID, "别的栏目", 1)

	list, total, err := repo.ListByColumn(context.Background(), col1.ID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total, "只返回该栏目已发布文章")
	require.Len(t, list, 1)
	assert.Equal(t, "已发布", list[0].Title)
}

func TestArticleRepoGetPublished(t *testing.T) {
	repo, db := newTestArticleRepo(t)
	col := makeColumn(t, db, "栏目")
	draft := makeArticle(t, db, col.ID, "草稿", 0)
	pub := makeArticle(t, db, col.ID, "已发布", 1)

	_, err := repo.GetPublished(context.Background(), draft.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "草稿不应被小程序详情访问")

	got, err := repo.GetPublished(context.Background(), pub.ID)
	require.NoError(t, err)
	assert.Equal(t, "已发布", got.Title)
	assert.Equal(t, "栏目", got.ColumnName)
}

func TestArticleRepoIncrementViewCount(t *testing.T) {
	repo, db := newTestArticleRepo(t)
	col := makeColumn(t, db, "栏目")
	a := makeArticle(t, db, col.ID, "文章", 1)

	require.NoError(t, repo.IncrementViewCount(context.Background(), a.ID))
	require.NoError(t, repo.IncrementViewCount(context.Background(), a.ID))

	got, err := repo.GetByID(context.Background(), a.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, got.ViewCount)
}
