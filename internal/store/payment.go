package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/repository"
)

// PaymentRepo 支付记录数据访问接口
type PaymentRepo interface {
	CRUD[model.Payment]
	// ListByFilter 按状态/业务类型分页查询（status<0 不过滤状态，bizType 为空不过滤类型）
	ListByFilter(ctx context.Context, page, pageSize, status int, bizType string) ([]model.Payment, int64, error)
	// GetByOrderNo 根据订单号查询
	GetByOrderNo(ctx context.Context, orderNo string) (*model.Payment, error)
}

// PaymentRepository PaymentRepo 的默认实现
type PaymentRepository struct {
	*repository.BaseRepo[model.Payment]
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepo {
	return &PaymentRepository{
		BaseRepo: repository.NewBaseRepo[model.Payment](db),
		db:       db,
	}
}

func (r *PaymentRepository) ListByFilter(ctx context.Context, page, pageSize, status int, bizType string) ([]model.Payment, int64, error) {
	var (
		list  []model.Payment
		total int64
	)
	db := r.db.WithContext(ctx).Model(&model.Payment{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}
	if bizType != "" {
		db = db.Where("biz_type = ?", bizType)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *PaymentRepository) GetByOrderNo(ctx context.Context, orderNo string) (*model.Payment, error) {
	var pay model.Payment
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&pay).Error; err != nil {
		return nil, err
	}
	return &pay, nil
}
