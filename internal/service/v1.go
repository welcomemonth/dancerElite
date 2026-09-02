package service

import (
	"github.com/welcomemonth/dancer-elite/internal/config"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// APIV1Service 服务容器，聚合所有业务 service，统一注入 Store 和配置。
// Server / Router 只需持有这一个容器，即可拿到任意 service。
type APIV1Service struct {
	Store *store.Store
	Cfg   *config.Config

	Auth           *AuthService
	BackendUser    *BackendUserService
	Article        *ArticleService
	Column         *ColumnService
	Role           *RoleService
	Menu           *MenuService
	User           *UserService
	Activity       *ActivityService
	Registration   *RegistrationService
	Season         *SeasonService
	ActivityResult *ActivityResultService
	AnnualRanking  *AnnualRankingService
	SystemConfig   *SystemConfigService
	Payment        *PaymentService
	Codegen        *CodegenService
	OperationLog   *OperationLogService
	Ranking        *RankingService
}

// NewAPIV1Service 装配所有 service
func NewAPIV1Service(cfg *config.Config, st *store.Store) *APIV1Service {
	svc := &APIV1Service{
		Store: st,
		Cfg:   cfg,
	}

	// SystemConfig 需先创建，Payment 依赖它
	svc.SystemConfig = NewSystemConfigService(st)
	svc.Auth = NewAuthService(st, cfg.JWT.Secret)
	svc.BackendUser = NewBackendUserService(st)
	svc.Article = NewArticleService(st)
	svc.Column = NewColumnService(st)
	svc.Role = NewRoleService(st)
	svc.Menu = NewMenuService(st)
	svc.User = NewUserService(st, cfg.Wechat.AppID, cfg.Wechat.Secret, cfg.JWT.Secret)
	svc.Activity = NewActivityService(st)
	svc.Registration = NewRegistrationService(st)
	svc.Season = NewSeasonService(st)
	svc.AnnualRanking = NewAnnualRankingService(st)
	svc.ActivityResult = NewActivityResultService(st, svc.AnnualRanking)
	svc.Payment = NewPaymentService(st, svc.SystemConfig)
	svc.Codegen = NewCodegenService(st)
	svc.OperationLog = NewOperationLogService(st)
	svc.Ranking = NewRankingService(st)

	return svc
}
