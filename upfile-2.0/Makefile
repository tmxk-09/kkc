## ============================================================
## Makefile — upfile-2.0 Go 服务快捷命令
## 用法: make <target>
## ============================================================

APP     = fileserver
GODIR   = server-go
DISTDIR = dist
GOPROXY ?= https://goproxy.cn,direct

.PHONY: all build run dev clean tidy cross help

# 默认目标
all: build

## build — 构建当前平台二进制
build:
	@./build.sh current

## all-platforms — 交叉编译所有平台
cross:
	@./build.sh all

## run — 构建并启动服务
run:
	@./build.sh run

## dev — 开发模式（go run，自动感知代码变更）
dev:
	@echo ">> 启动开发模式 (go run)..."
	@cd $(GODIR) && GOPROXY=$(GOPROXY) go run .

## tidy — 同步整理依赖
tidy:
	@echo ">> go mod tidy..."
	@cd $(GODIR) && GOPROXY=$(GOPROXY) go mod tidy
	@echo "✅ 依赖整理完成"

## test — 运行单元测试
test:
	@echo ">> go test..."
	@cd $(GODIR) && GOPROXY=$(GOPROXY) go test ./... -v

## clean — 清理构建产物
clean:
	@./build.sh clean

## docker-build — Docker 多阶段构建镜像
docker-build:
	@echo ">> Docker 构建..."
	docker build -t upfile-server:latest .

## docker-up — Docker Compose 启动
docker-up:
	docker compose up -d --build

## docker-down — Docker Compose 停止
docker-down:
	docker compose down

## logs — 查看服务日志
logs:
	docker compose logs -f file-server

## help — 显示帮助
help:
	@echo ""
	@echo "  upfile-2.0 Go 服务构建命令"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## /  make /' | column -t -s '—'
	@echo ""
