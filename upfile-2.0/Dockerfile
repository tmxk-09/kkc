# ============================================================
# 多阶段构建：Go 后端 (CGO_ENABLED=0 纯静态二进制)
# ============================================================

# ===== 构建阶段 =====
FROM golang:1.22-alpine AS builder

WORKDIR /build

# 设置 Go 代理（国内加速）
ENV GOPROXY=https://goproxy.cn,direct

# 先复制 go.mod / go.sum，利用 Docker 缓存层
COPY server-go/go.mod server-go/go.sum ./
RUN go mod download

# 复制源码并编译（CGO_ENABLED=0 = 无 CGO，modernc/sqlite 纯 Go）
COPY server-go/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o fileserver .

# ===== 运行阶段 =====
FROM alpine:3.20

# 添加时区支持（中国时区）
RUN apk add --no-cache ca-certificates tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

# 从构建阶段复制二进制和配置
COPY --from=builder /build/fileserver .
COPY server-go/config.yaml .

# 创建数据持久化目录
RUN mkdir -p data uploads/files uploads/chunks

EXPOSE 9102

CMD ["./fileserver"]
