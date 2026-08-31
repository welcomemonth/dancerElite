package service

import (
	"context"
	"strconv"
	"sync"

	"gorm.io/gorm"

	"github.com/zzhtl/go-mountain/internal/model"
	"github.com/zzhtl/go-mountain/internal/store"
)

// SystemConfigService 系统配置服务
type SystemConfigService struct {
	st    *store.Store
	cache sync.Map // 内存缓存，避免频繁读库
}

// NewSystemConfigService 创建系统配置服务
func NewSystemConfigService(st *store.Store) *SystemConfigService {
	svc := &SystemConfigService{st: st}
	// 启动时预加载所有配置到缓存
	svc.ReloadCache(context.Background())
	return svc
}

// ReloadCache 重新加载缓存
func (s *SystemConfigService) ReloadCache(ctx context.Context) {
	configs, err := s.st.SystemConfigRepo.FindAll(ctx)
	if err != nil {
		return
	}
	s.cache = sync.Map{}
	for _, c := range configs {
		s.cache.Store(c.Key, c.Value)
	}
}

// GetValue 获取配置值（优先缓存）
func (s *SystemConfigService) GetValue(ctx context.Context, key string) string {
	if v, ok := s.cache.Load(key); ok {
		return v.(string)
	}
	config, err := s.st.SystemConfigRepo.GetByKey(ctx, key)
	if err != nil {
		return ""
	}
	s.cache.Store(key, config.Value)
	return config.Value
}

// GetValueInt 获取整型配置值
func (s *SystemConfigService) GetValueInt(ctx context.Context, key string, defaultVal int) int {
	v := s.GetValue(ctx, key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return i
}

// GetValueBool 获取布尔配置值
func (s *SystemConfigService) GetValueBool(ctx context.Context, key string) bool {
	v := s.GetValue(ctx, key)
	return v == "true" || v == "1"
}

// SetValue 设置配置值（Upsert）
func (s *SystemConfigService) SetValue(ctx context.Context, key, value, typ, groupName, remark string) error {
	config := model.SystemConfig{
		Key:       key,
		Value:     value,
		Type:      typ,
		GroupName: groupName,
		Remark:    remark,
	}

	if err := s.st.SystemConfigRepo.Upsert(ctx, &config); err != nil {
		return err
	}
	s.cache.Store(key, value)
	return nil
}

// BatchSet 批量设置配置
func (s *SystemConfigService) BatchSet(ctx context.Context, configs []model.SystemConfig) error {
	if err := s.st.SystemConfigRepo.BatchUpsert(ctx, configs); err != nil {
		return err
	}
	for _, c := range configs {
		s.cache.Store(c.Key, c.Value)
	}
	return nil
}

// List 获取所有配置
func (s *SystemConfigService) List(ctx context.Context) ([]model.SystemConfig, error) {
	return s.st.SystemConfigRepo.FindAll(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Order("group_name, `key`")
	})
}

// ListByGroup 按分组获取配置
func (s *SystemConfigService) ListByGroup(ctx context.Context, groupName string) ([]model.SystemConfig, error) {
	return s.st.SystemConfigRepo.FindAll(ctx, func(db *gorm.DB) *gorm.DB {
		return db.Where("group_name = ?", groupName).Order("`key`")
	})
}

// GetGroups 获取所有分组名称
func (s *SystemConfigService) GetGroups(ctx context.Context) ([]string, error) {
	return s.st.SystemConfigRepo.ListGroups(ctx)
}

// Delete 删除配置项
func (s *SystemConfigService) Delete(ctx context.Context, key string) error {
	s.cache.Delete(key)
	return s.st.SystemConfigRepo.DeleteByKey(ctx, key)
}

// GetWechatPayConfig 获取微信支付相关配置（便捷方法）
func (s *SystemConfigService) GetWechatPayConfig(ctx context.Context) map[string]string {
	keys := []string{
		"wechat.app_id", "wechat.secret",
		"wechat.mch_id", "wechat.mch_api_v3_key",
		"wechat.mch_serial_no", "wechat.mch_private_key",
		"wechat.notify_url",
	}
	result := make(map[string]string, len(keys))
	for _, k := range keys {
		result[k] = s.GetValue(ctx, k)
	}
	return result
}

// InitDefaultConfigs 初始化默认配置项（如果不存在）
func (s *SystemConfigService) InitDefaultConfigs(ctx context.Context) {
	defaults := []model.SystemConfig{
		// 站点信息
		{Key: "site.name", Value: "远山公益", Type: "string", GroupName: "站点信息", Remark: "站点名称"},
		{Key: "site.description", Value: "远山公益管理平台", Type: "string", GroupName: "站点信息", Remark: "站点描述"},

		// 微信小程序
		{Key: "wechat.app_id", Value: "", Type: "string", GroupName: "微信配置", Remark: "小程序 AppID"},
		{Key: "wechat.secret", Value: "", Type: "string", GroupName: "微信配置", Remark: "小程序 Secret"},

		// 微信支付
		{Key: "wechat.mch_id", Value: "", Type: "string", GroupName: "微信支付", Remark: "商户号"},
		{Key: "wechat.mch_api_v3_key", Value: "", Type: "string", GroupName: "微信支付", Remark: "APIv3 密钥"},
		{Key: "wechat.mch_serial_no", Value: "", Type: "string", GroupName: "微信支付", Remark: "商户证书序列号"},
		{Key: "wechat.mch_private_key", Value: "", Type: "string", GroupName: "微信支付", Remark: "商户私钥内容（PEM 格式）"},
		{Key: "wechat.notify_url", Value: "", Type: "string", GroupName: "微信支付", Remark: "支付回调地址（如 https://example.com/api/payment/wechat/notify）"},
	}

	for _, d := range defaults {
		// 仅在 key 不存在时插入
		exists, _ := s.st.SystemConfigRepo.Exists(ctx, "`key` = ?", d.Key)
		if exists {
			continue
		}
		if err := s.st.SystemConfigRepo.Create(ctx, &d); err != nil {
			continue
		}
		s.cache.Store(d.Key, d.Value)
	}
}
