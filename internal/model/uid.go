package model

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenUid() string {
	// 生成 6 字节的随机数
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	// 使用 base64 编码，然后移除特殊字符
	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(b), "=", "")[:8]
}
