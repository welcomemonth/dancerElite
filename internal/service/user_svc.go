package service

import (
	"context"
	"errors"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"gorm.io/gorm"

	"github.com/welcomemonth/dancer-elite/internal/model"
	"github.com/welcomemonth/dancer-elite/internal/pkg/errcode"
	"github.com/welcomemonth/dancer-elite/internal/pkg/token"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// UserService 小程序用户管理服务
type UserService struct {
	st        *store.Store
	appID     string
	appSecret string
	jwtSecret string
}

// NewUserService 创建小程序用户服务
func NewUserService(st *store.Store, appID, appSecret, jwtSecret string) *UserService {
	return &UserService{
		st:        st,
		appID:     appID,
		appSecret: appSecret,
		jwtSecret: jwtSecret,
	}
}

// GenerateToken 为小程序用户签发 JWT（供报名/支付等受保护接口使用）
func (s *UserService) GenerateToken(userID int64, openID string) (string, error) {
	claims := map[string]any{
		"user_id": userID,
		"openid":  openID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	return token.Sign(s.jwtSecret, claims)
}

// WechatLogin 微信小程序登录
func (s *UserService) WechatLogin(ctx context.Context, code string) (*model.User, error) {
	mp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     s.appID,
		Secret:    s.appSecret,
		HttpDebug: false,
		Debug:     false,
	})
	if err != nil {
		return nil, err
	}

	session, err := mp.Auth.Session(ctx, code)
	if err != nil {
		return nil, err
	}

	openID := session.OpenID

	// 查询或创建用户
	user, err := s.st.UserRepo.GetByOpenID(ctx, openID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &model.User{OpenID: openID, Status: 1}
		if err := s.st.UserRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

// Register 注册/绑定手机号
func (s *UserService) Register(ctx context.Context, phone, openID, name string) (*model.User, error) {
	// 通过 openid 查找用户
	user, err := s.st.UserRepo.GetByOpenID(ctx, openID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在则创建
		user = &model.User{
			OpenID: openID,
			Phone:  phone,
			Name:   name,
			Status: 1,
		}
		if err := s.st.UserRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	} else if err != nil {
		return nil, err
	}

	// 已存在则更新
	s.st.UserRepo.Update(ctx, user.ID, map[string]any{
		"phone": phone,
		"name":  name,
	})

	return user, nil
}

// List 获取小程序用户列表（后台）
func (s *UserService) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}
	return s.st.UserRepo.List(ctx, page, pageSize, scope)
}

// Get 获取用户详情
func (s *UserService) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.st.UserRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return user, nil
}

// Update 更新用户信息
func (s *UserService) Update(ctx context.Context, id int64, updates map[string]any) error {
	return s.st.UserRepo.Update(ctx, id, updates)
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.st.UserRepo.Delete(ctx, id)
}
