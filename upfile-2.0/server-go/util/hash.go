package util

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// FileMD5 计算文件的 MD5 Hash（等价 Java HashUtil.fileMd5）
func FileMD5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// StrMD5 计算字符串的 MD5
func StrMD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
