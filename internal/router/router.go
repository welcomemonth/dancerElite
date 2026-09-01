package router

import (
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/handler"
	"github.com/welcomemonth/dancer-elite/internal/middleware"
	"github.com/welcomemonth/dancer-elite/internal/service"
	"github.com/welcomemonth/dancer-elite/internal/web"
)

// Setup 配置所有路由
func Setup(engine *gin.Engine, svc *service.APIV1Service) {
	// 全局中间件
	engine.Use(middleware.RequestID(), middleware.GinLogger(), middleware.CORS())

	// 设置最大文件上传大小
	engine.MaxMultipartMemory = 50 << 20 // 50 MB

	// 静态文件服务
	engine.Static("/uploads", "./uploads")

	// API 路由
	api := engine.Group("/api")

	// 健康检查
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	registerMiniProgramRoutes(api, svc)
	registerAdminRoutes(api, svc)

	setupAdminUI(engine)
}

// registerMiniProgramRoutes 小程序端 API + 微信支付回调
func registerMiniProgramRoutes(api *gin.RouterGroup, svc *service.APIV1Service) {
	userHandler := handler.NewUserHandler(svc.User)
	columnHandler := handler.NewColumnHandler(svc.Column)
	articleHandler := handler.NewArticleHandler(svc.Article)
	activityHandler := handler.NewActivityHandler(svc.Activity)
	registrationHandler := handler.NewRegistrationHandler(svc.Registration)
	paymentHandler := handler.NewPaymentHandler(svc.Payment, svc.Registration, svc.Activity)
	rankingHandler := handler.NewRankingHandler(svc.Ranking)

	mp := api.Group("/mp")
	{
		mp.POST("/login", userHandler.WechatLogin)
		mp.POST("/register", userHandler.Register)

		mpColumns := mp.Group("/columns")
		mpColumns.GET("/", columnHandler.ListForMP)

		mpArticles := mp.Group("/articles")
		mpArticles.GET("/column/:columnId", articleHandler.ListByColumn)
		mpArticles.GET("/:id", articleHandler.GetForMP)

		// 活动（公开）
		mpActivities := mp.Group("/activities")
		mpActivities.GET("/", activityHandler.ListForMP)
		mpActivities.GET("/:id", activityHandler.GetForMP)

		// 排行榜 / 选手（公开只读）
		mp.GET("/seasons/active", rankingHandler.ActiveSeason)
		mp.GET("/rankings", rankingHandler.Leaderboard)
		mp.GET("/rankings/organization", rankingHandler.OrgLeaderboard)
		mp.GET("/players/:id", rankingHandler.PlayerDetail)

		// 需要小程序用户认证的接口
		mpAuth := mp.Group("")
		mpAuth.Use(middleware.JWTAuth(svc.Cfg.JWT.Secret))
		{
			// 报名
			mpAuth.POST("/registrations", registrationHandler.Create)
			mpAuth.PUT("/registrations/:id/cancel", registrationHandler.Cancel)
			mpAuth.GET("/registrations/mine", registrationHandler.MyRegistrations)

			// 支付
			mpAuth.POST("/payments/create", paymentHandler.CreateOrder)
			mpAuth.GET("/payments/query", paymentHandler.QueryOrder)
		}
	}

	// 微信支付回调（不需要任何认证）
	api.POST("/payment/wechat/notify", paymentHandler.WechatNotify)
}

// registerAdminRoutes 后台管理端 API
func registerAdminRoutes(api *gin.RouterGroup, svc *service.APIV1Service) {
	authHandler := handler.NewAuthHandler(svc.Auth)
	backendUserHandler := handler.NewBackendUserHandler(svc.BackendUser)
	articleHandler := handler.NewArticleHandler(svc.Article)
	columnHandler := handler.NewColumnHandler(svc.Column)
	roleHandler := handler.NewRoleHandler(svc.Role)
	menuHandler := handler.NewMenuHandler(svc.Menu)
	userHandler := handler.NewUserHandler(svc.User)
	uploadHandler := handler.NewUploadHandler()
	activityHandler := handler.NewActivityHandler(svc.Activity)
	registrationHandler := handler.NewRegistrationHandler(svc.Registration)
	paymentHandler := handler.NewPaymentHandler(svc.Payment, svc.Registration, svc.Activity)
	systemConfigHandler := handler.NewSystemConfigHandler(svc.SystemConfig)
	codegenHandler := handler.NewCodegenHandler(svc.Codegen)
	operationLogHandler := handler.NewOperationLogHandler(svc.OperationLog)
	seasonHandler := handler.NewSeasonHandler(svc.Season)

	admin := api.Group("/admin")

	// 认证路由（不需要 JWT）
	auth := admin.Group("/backend-auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// 需要 JWT 认证 + RBAC 权限校验的路由
	adminAuth := admin.Group("")
	adminAuth.Use(middleware.JWTAuth(svc.Cfg.JWT.Secret))
	adminAuth.Use(middleware.RBACAuth(svc.Store))
	adminAuth.Use(middleware.OperationLogger(svc.Store))
	{
		// 修改密码
		adminAuth.PUT("/backend-auth/change-password", authHandler.ChangePassword)

		// 后台用户管理
		bu := adminAuth.Group("/backend-users")
		bu.GET("/", backendUserHandler.List)
		bu.POST("/", backendUserHandler.Create)
		bu.GET("/current/menus", backendUserHandler.GetCurrentUserMenus)
		bu.GET("/:id", backendUserHandler.Get)
		bu.PUT("/:id", backendUserHandler.Update)
		bu.DELETE("/:id", backendUserHandler.Delete)
		bu.PUT("/:id/status", backendUserHandler.UpdateStatus)
		bu.PUT("/:id/reset-password", backendUserHandler.ResetPassword)

		// 小程序用户管理
		users := adminAuth.Group("/users")
		users.GET("/", userHandler.List)
		users.GET("/:id", userHandler.Get)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)

		// 文章管理
		articles := adminAuth.Group("/articles")
		articles.GET("/", articleHandler.List)
		articles.POST("/", articleHandler.Create)
		articles.GET("/:id", articleHandler.Get)
		articles.PUT("/:id", articleHandler.Update)
		articles.DELETE("/:id", articleHandler.Delete)
		articles.PUT("/:id/status", articleHandler.UpdateStatus)

		// 栏目管理
		columns := adminAuth.Group("/columns")
		columns.GET("/", columnHandler.List)
		columns.POST("/", columnHandler.Create)
		columns.GET("/:id", columnHandler.Get)
		columns.PUT("/:id", columnHandler.Update)
		columns.DELETE("/:id", columnHandler.Delete)

		// 角色管理
		roles := adminAuth.Group("/roles")
		roles.GET("/", roleHandler.List)
		roles.POST("/", roleHandler.Create)
		roles.GET("/:id", roleHandler.Get)
		roles.PUT("/:id", roleHandler.Update)
		roles.DELETE("/:id", roleHandler.Delete)
		roles.PUT("/:id/status", roleHandler.UpdateStatus)
		roles.GET("/:id/menus", roleHandler.GetRoleMenus)
		roles.PUT("/:id/menus", roleHandler.UpdateRoleMenus)

		// 菜单管理
		menus := adminAuth.Group("/menus")
		menus.GET("/", menuHandler.List)
		menus.GET("/tree", menuHandler.Tree)
		menus.POST("/", menuHandler.Create)
		menus.GET("/:id", menuHandler.Get)
		menus.PUT("/:id", menuHandler.Update)
		menus.DELETE("/:id", menuHandler.Delete)
		menus.PUT("/:id/status", menuHandler.UpdateStatus)

		// 操作日志
		operationLogs := adminAuth.Group("/operation-logs")
		operationLogs.GET("/", operationLogHandler.List)

		// 活动管理
		activities := adminAuth.Group("/activities")
		activities.GET("/", activityHandler.List)
		activities.POST("/", activityHandler.Create)
		activities.GET("/:id", activityHandler.Get)
		activities.PUT("/:id", activityHandler.Update)
		activities.DELETE("/:id", activityHandler.Delete)
		activities.PUT("/:id/status", activityHandler.UpdateStatus)

		// 赛季管理
		seasons := adminAuth.Group("/seasons")
		seasons.GET("/", seasonHandler.List)
		seasons.POST("/", seasonHandler.Create)
		seasons.GET("/:id", seasonHandler.Get)
		seasons.PUT("/:id", seasonHandler.Update)
		seasons.DELETE("/:id", seasonHandler.Delete)
		seasons.PUT("/:id/status", seasonHandler.UpdateStatus)

		// 报名管理
		registrations := adminAuth.Group("/registrations")
		registrations.GET("/", registrationHandler.List)
		registrations.GET("/:id", registrationHandler.Get)

		// 支付管理
		payments := adminAuth.Group("/payments")
		payments.GET("/", paymentHandler.List)
		payments.GET("/:id", paymentHandler.Get)
		payments.PUT("/:id/refund", paymentHandler.Refund)

		// 系统配置管理
		sysConfigs := adminAuth.Group("/system-configs")
		sysConfigs.GET("/", systemConfigHandler.List)
		sysConfigs.GET("/groups", systemConfigHandler.GetGroups)
		sysConfigs.POST("/", systemConfigHandler.Save)
		sysConfigs.POST("/batch", systemConfigHandler.BatchSave)
		sysConfigs.DELETE("/", systemConfigHandler.Delete)

		// 代码生成器
		codegen := adminAuth.Group("/codegen")
		codegen.GET("/tables", codegenHandler.GetTables)
		codegen.GET("/columns", codegenHandler.GetTableColumns)
		codegen.GET("/", codegenHandler.List)
		codegen.POST("/", codegenHandler.Create)
		codegen.GET("/:id", codegenHandler.Get)
		codegen.PUT("/:id", codegenHandler.Update)
		codegen.DELETE("/:id", codegenHandler.Delete)
		codegen.GET("/:id/preview", codegenHandler.Preview)
		codegen.POST("/:id/generate", codegenHandler.Generate)

		// 文件上传
		upload := adminAuth.Group("/upload")
		upload.POST("/image", uploadHandler.UploadImage)
		upload.POST("/video", uploadHandler.UploadVideo)
	}
}

func setupAdminUI(engine *gin.Engine) {
	adminFS := http.FS(web.Dist())

	// Admin UI 静态文件
	engine.GET("/web", func(c *gin.Context) { c.Redirect(http.StatusFound, "/web/") })
	engine.GET("/web/*filepath", func(c *gin.Context) {
		filePath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filePath == "" {
			filePath = "index.html"
		}
		if serveEmbeddedFile(c, adminFS, filePath) {
			return
		}
		if path.Ext(filePath) == "" {
			serveEmbeddedFile(c, adminFS, "index.html")
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not found"})
	})

	// SPA 回退
	engine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/web/") {
			serveEmbeddedFile(c, adminFS, "index.html")
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not found"})
	})

	// 重定向
	engine.GET("/admin", func(c *gin.Context) { c.Redirect(http.StatusFound, "/web/") })
	engine.GET("/admin/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/web/") })
}

func serveEmbeddedFile(c *gin.Context, fileSystem http.FileSystem, filePath string) bool {
	cleanPath := strings.TrimPrefix(path.Clean("/"+filePath), "/")
	if cleanPath == "." || cleanPath == "" {
		cleanPath = "index.html"
	}

	file, err := fileSystem.Open(cleanPath)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return false
	}

	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), file)
	return true
}
