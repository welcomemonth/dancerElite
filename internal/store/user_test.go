package store

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
)

// newTestUserRepo 构建基于内存 sqlite 的 UserRepo。
// 项目已有 sqlite 驱动，无需 mock，直接对真实 SQL 做集成测试。
func newTestUserRepo(t *testing.T) UserRepo {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// 内存 sqlite 每个连接是独立库，限制为单连接避免"表不存在"
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	require.NoError(t, db.AutoMigrate(&model.User{}))
	return NewUserRepository(db)
}

func seedUser(t *testing.T, repo UserRepo, openID, name string) *model.User {
	t.Helper()
	u := &model.User{OpenID: openID, Name: name, Status: 1}
	require.NoError(t, repo.Create(context.Background(), u))
	return u
}

func TestUserRepoCreateAndGetByID(t *testing.T) {
	repo := newTestUserRepo(t)

	u := seedUser(t, repo, "openid-1", "张三")
	require.NotZero(t, u.ID, "创建后应自动生成主键")

	got, err := repo.GetByID(context.Background(), u.ID)
	require.NoError(t, err)
	assert.Equal(t, "openid-1", got.OpenID)
	assert.Equal(t, "张三", got.Name)
}

func TestUserRepoGetByIDNotFound(t *testing.T) {
	repo := newTestUserRepo(t)

	_, err := repo.GetByID(context.Background(), 9999)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestUserRepoGetByOpenID(t *testing.T) {
	repo := newTestUserRepo(t)
	seedUser(t, repo, "openid-x", "李四")

	got, err := repo.GetByOpenID(context.Background(), "openid-x")
	require.NoError(t, err)
	assert.Equal(t, "李四", got.Name)

	_, err = repo.GetByOpenID(context.Background(), "not-exist")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestUserRepoUpdate(t *testing.T) {
	repo := newTestUserRepo(t)
	u := seedUser(t, repo, "openid-u", "王五")

	err := repo.Update(context.Background(), u.ID, map[string]any{
		"name":  "王五改",
		"phone": "13800000000",
	})
	require.NoError(t, err)

	got, err := repo.GetByID(context.Background(), u.ID)
	require.NoError(t, err)
	assert.Equal(t, "王五改", got.Name)
	assert.Equal(t, "13800000000", got.Phone)
}

func TestUserRepoDelete(t *testing.T) {
	repo := newTestUserRepo(t)
	u := seedUser(t, repo, "openid-d", "赵六")

	require.NoError(t, repo.Delete(context.Background(), u.ID))

	// 软删除后 GetByID 应查不到
	_, err := repo.GetByID(context.Background(), u.ID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestUserRepoList(t *testing.T) {
	repo := newTestUserRepo(t)
	seedUser(t, repo, "openid-1", "用户1")
	seedUser(t, repo, "openid-2", "用户2")
	seedUser(t, repo, "openid-3", "用户3")

	orderByIDDesc := func(db *gorm.DB) *gorm.DB { return db.Order("id DESC") }

	list, total, err := repo.List(context.Background(), 1, 2, orderByIDDesc)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total, "total 应为总数（不受分页影响）")
	require.Len(t, list, 2, "第一页应返回 2 条")
	assert.Equal(t, "用户3", list[0].Name, "倒序应 id 最大的在前")
	assert.Equal(t, "用户2", list[1].Name)

	list2, _, err := repo.List(context.Background(), 2, 2, orderByIDDesc)
	require.NoError(t, err)
	require.Len(t, list2, 1)
	assert.Equal(t, "用户1", list2[0].Name)
}
