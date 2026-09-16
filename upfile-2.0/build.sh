#!/usr/bin/env bash
# =================================================================
# build.sh — upfile-2.0 Go 服务构建脚本
# 兼容 macOS / Linux
# 用法：
#   ./build.sh             # 构建当前平台
#   ./build.sh all         # 构建所有平台（darwin-arm64/amd64, linux-amd64/arm64）
#   ./build.sh run         # 构建并直接运行（当前平台）
#   ./build.sh clean       # 清理产物
# =================================================================

set -euo pipefail

# ===== 配置 =====
APP_NAME="fileserver"
ENTRY="."                          # main.go 所在目录（相对于 server-go/）
OUTPUT_DIR="../dist"               # 输出目录（相对于 server-go/）
LDFLAGS="-s -w"                   # 去掉符号表 & 调试信息，缩小体积
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

# ===== 颜色输出 =====
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()    { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

# ===== 进入 server-go 目录 =====
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/server-go"

# ===== 检查 Go 环境（自动探测 SDK 路径）=====
check_go() {
    # 自动探测常见安装路径
    for candidate in \
        "$(command -v go 2>/dev/null)" \
        "$HOME/sdk/go/bin/go" \
        "/usr/local/go/bin/go" \
        "/opt/homebrew/bin/go"; do
        if [[ -x "$candidate" ]]; then
            export PATH="$(dirname $candidate):$PATH"
            break
        fi
    done

    if ! command -v go &>/dev/null; then
        error "未找到 go 命令，请先安装 Go 1.22+\n  Mac:   brew install go\n  Linux: https://go.dev/dl/"
    fi
    GO_VER=$(go version | awk '{print $3}')
    info "Go 版本: $GO_VER"
}

# ===== 下载依赖 =====
deps() {
    info "同步依赖（go mod tidy）..."
    GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" go mod tidy
    ok "依赖同步完成"
}

# ===== 构建单个目标 =====
# build_target <GOOS> <GOARCH> <output_name>
build_target() {
    local goos=$1
    local goarch=$2
    local out_name=$3

    local out_path="$OUTPUT_DIR/$out_name"
    info "编译 → $goos/$goarch  输出: $out_path"

    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" \
        go build \
        -ldflags "$LDFLAGS -X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME'" \
        -o "$out_path" \
        "$ENTRY"

    local size
    size=$(du -sh "$out_path" | cut -f1)
    ok "$out_name  [$size]"
}

# ===== 构建当前平台 =====
build_current() {
    mkdir -p "$OUTPUT_DIR"

    local goos goarch suffix
    goos=$(go env GOOS)
    goarch=$(go env GOARCH)
    suffix=""
    [[ "$goos" == "windows" ]] && suffix=".exe"

    out_name="${APP_NAME}-${goos}-${goarch}${suffix}"
    build_target "$goos" "$goarch" "$out_name"

    # 同时输出一个无后缀的快捷链接（当前平台）
    cp "$OUTPUT_DIR/$out_name" "$OUTPUT_DIR/$APP_NAME"
    ok "当前平台快捷路径: dist/$APP_NAME"
}

# ===== 构建所有平台 =====
build_all() {
    mkdir -p "$OUTPUT_DIR"
    info "开始多平台交叉编译..."

    build_target "darwin"  "arm64"  "${APP_NAME}-darwin-arm64"    # Apple Silicon
    build_target "darwin"  "amd64"  "${APP_NAME}-darwin-amd64"    # Intel Mac
    build_target "linux"   "amd64"  "${APP_NAME}-linux-amd64"     # x86_64 Linux
    build_target "linux"   "arm64"  "${APP_NAME}-linux-arm64"     # ARM Linux (树莓派/云服务器)

    echo ""
    ok "所有平台编译完成，产物目录: $OUTPUT_DIR"
    ls -lh "$OUTPUT_DIR/"
}

# ===== 运行 =====
run_server() {
    build_current
    info "启动服务..."
    "$OUTPUT_DIR/$APP_NAME"
}

# ===== 清理 =====
clean() {
    rm -rf "$OUTPUT_DIR"
    ok "已清理 dist/ 目录"
}

# ===== 主入口 =====
CMD="${1:-current}"

check_go

case "$CMD" in
    all)     deps && build_all ;;
    run)     deps && run_server ;;
    clean)   clean ;;
    current) deps && build_current ;;
    *)
        warn "未知命令: $CMD"
        echo "用法: ./build.sh [current|all|run|clean]"
        exit 1
        ;;
esac
