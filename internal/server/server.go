package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/welcomemonth/dancer-elite/internal/config"
	"github.com/welcomemonth/dancer-elite/internal/router"
	"github.com/welcomemonth/dancer-elite/internal/service"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

// Server 封装 HTTP 服务器
type Server struct {
	engine    *gin.Engine
	cfg       *config.Config
	v1Service *service.APIV1Service
}

// NewServer 创建服务器实例
func NewServer(st *store.Store, cfg *config.Config) *Server {
	engine := gin.Default()
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		engine: engine,
		cfg:    cfg,
	}
	apiV1Service := service.NewAPIV1Service(cfg, st)
	server.v1Service = apiV1Service

	router.Setup(engine, apiV1Service)

	return server
}

// Run 启动 HTTP 服务器（支持优雅关闭）
func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	// 创建可取消的上下文监听系统信号
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待关闭信号
	<-ctx.Done()
	log.Println("收到关闭信号，正在优雅关闭...")

	// 给 5 秒时间处理完当前请求
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	// 关闭数据库连接
	s.v1Service.Store.Close()

	log.Println("服务器已关闭")
	return nil
}
