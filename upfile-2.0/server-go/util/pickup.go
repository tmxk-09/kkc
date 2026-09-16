package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const pickupChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 剔除易混淆字 I/1/O/0

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// GeneratePickupCode 生成 6 位随机取件码（与 Java 版等价）
// existsFn：检查码是否已被占用
func GeneratePickupCode(existsFn func(code string) bool) string {
	for {
		var sb strings.Builder
		for i := 0; i < 6; i++ {
			sb.WriteByte(pickupChars[rng.Intn(len(pickupChars))])
		}
		code := sb.String()
		if !existsFn(code) {
			return code
		}
	}
}

// FormatSize 格式化文件大小
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
