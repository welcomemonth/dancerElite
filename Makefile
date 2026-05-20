.PHONY: run build clean migrate frontend-admin help dev install deploy

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  run              - 运行服务器"
	@echo "  build            - 构建项目"
	@echo "  migrate          - 数据库迁移（含创建管理员、初始化角色菜单）"
	@echo "  frontend-admin   - 构建管理后台前端"
	@echo "  clean            - 清理构建文件"
	@echo "  dev              - 开发环境启动"
	@echo "  install          - 安装所有依赖"
	@echo "  deploy           - 一键构建、迁移并启动本机部署"

# 运行服务器
run:
	@echo "启动服务器..."
	$(MAKE) frontend-admin
	go run cmd/server/main.go

# 构建项目
build: frontend-admin
	@echo "构建项目..."
	go build -o bin/go-mountain cmd/server/main.go

# 数据库迁移（创建表、初始化默认数据、创建管理员）
migrate:
	@echo "执行数据库迁移..."
	go run cmd/migrate/main.go

# 构建管理后台前端
frontend-admin:
	@echo "构建管理后台前端..."
	cd frontend-admin && bun run build
	@touch internal/web/dist/.gitkeep

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf frontend-admin/dist/
	rm -rf internal/web/dist/*
	@mkdir -p internal/web/dist
	@touch internal/web/dist/.gitkeep

# 开发环境启动
dev: frontend-admin
	@echo "启动开发环境..."
	go run cmd/server/main.go

# 安装依赖
install:
	@echo "安装Go依赖..."
	go mod tidy
	@echo "安装前端依赖..."
	cd frontend-admin && bun install

# 一键本机部署
deploy:
	@echo "执行一键部署..."
	bash scripts/deploy.sh
