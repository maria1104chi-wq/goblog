package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// GenerateSMSCode 生成短信验证码
func GenerateSMSCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	code := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败：%w", err)
		}
		code += fmt.Sprintf("%d", n.Int64())
	}

	return code, nil
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	
	uuid := make([]byte, 36)
	hex.Encode(uuid[0:8], b[0:4])
	hex.Encode(uuid[9:13], b[4:6])
	hex.Encode(uuid[14:18], b[6:8])
	hex.Encode(uuid[19:23], b[8:10])
	hex.Encode(uuid[24:], b[10:])
	
	uuid[8] = '-'
	uuid[13] = '-'
	uuid[18] = '-'
	uuid[23] = '-'
	
	return string(uuid)
}

// GenerateSlug 从标题生成slug
func GenerateSlug(title string) string {
	// 简化实现，实际项目中可以使用更复杂的slug生成算法
	slug := ""
	for _, r := range title {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			slug += string(r)
		} else if r >= '\u4e00' && r <= '\u9fff' {
			// 中文字符转换为拼音或保留（这里简化处理）
			slug += fmt.Sprintf("%x", r)
		} else {
			slug += "-"
		}
	}
	
	// 移除连续的连字符
	for len(slug) > 0 && slug[0] == '-' {
		slug = slug[1:]
	}
	for len(slug) > 0 && slug[len(slug)-1] == '-' {
		slug = slug[:len(slug)-1]
	}
	
	// 限制长度
	if len(slug) > 200 {
		slug = slug[:200]
	}
	
	return slug
}
