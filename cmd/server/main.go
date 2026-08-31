package main

import (
	"log"

	"github.com/welcomemonth/dancer-elite/internal/config"
	"github.com/welcomemonth/dancer-elite/internal/server"
	"github.com/welcomemonth/dancer-elite/internal/store"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库（Store）
	st, err := store.New(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 启动服务器
	srv := server.NewServer(st, cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}
