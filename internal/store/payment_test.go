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

func newTestPaymentRepo(t *testing.T) PaymentRepo {
	t.Helper()
	return NewPaymentRepository(newTestDB(t, &model.Payment{}))
}

func makePayment(t *testing.T, repo PaymentRepo, orderNo string, status int, bizType string) *model.Payment {
	t.Helper()
	p := &model.Payment{
		OrderNo: orderNo,
		UserID:  1,
		Amount:  100,
		PayType: "wechat_jsapi",
		Status:  status,
		BizType: bizType,
		BizID:   1,
	}
	require.NoError(t, repo.Create(context.Background(), p))
	return p
}

func TestPaymentRepoListByFilter(t *testing.T) {
	repo := newTestPaymentRepo(t)
	makePayment(t, repo, "P20260831001", 1, "registration")
	makePayment(t, repo, "P20260831002", 0, "registration")
	makePayment(t, repo, "P20260831003", 2, "donation")

	// 不过滤
	list, total, err := repo.ListByFilter(context.Background(), 1, 10, -1, "")
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.Len(t, list, 3)

	// 按状态过滤
	list, total, err = repo.ListByFilter(context.Background(), 1, 10, 1, "")
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, "P20260831001", list[0].OrderNo)

	// 按业务类型过滤
	list, total, err = repo.ListByFilter(context.Background(), 1, 10, -1, "donation")
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	require.Len(t, list, 1)
	assert.Equal(t, "donation", list[0].BizType)
}

func TestPaymentRepoGetByOrderNo(t *testing.T) {
	repo := newTestPaymentRepo(t)
	pay := makePayment(t, repo, "P20260831001", 1, "registration")

	got, err := repo.GetByOrderNo(context.Background(), "P20260831001")
	require.NoError(t, err)
	assert.Equal(t, pay.ID, got.ID)
	assert.Equal(t, int64(1), got.UserID)

	_, err = repo.GetByOrderNo(context.Background(), "NOTEXIST")
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
