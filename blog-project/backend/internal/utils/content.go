package utils

import (
	"blog-backend/internal/model"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// HashPassword 密码哈希
func HashPassword(password string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyPassword 验证密码
func VerifyPassword(password, hash string) bool {
	newHash, err := HashPassword(password)
	if err != nil {
		return false
	}
	return newHash == hash
}

// MarkdownToHTML 将Markdown转换为HTML
func MarkdownToHTML(markdown string) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf strings.Builder
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SanitizeHTML 清理HTML（防止XSS）
func SanitizeHTML(html_content string) string {
	// 简化实现，实际应该使用bluemonday等库进行更严格的过滤
	// 这里只做基本的script标签过滤
	html_content = strings.ReplaceAll(html_content, "<script>", "&lt;script&gt;")
	html_content = strings.ReplaceAll(html_content, "</script>", "&lt;/script&gt;")
	html_content = strings.ReplaceAll(html_content, "javascript:", "")
	html_content = strings.ReplaceAll(html_content, "onerror=", "")
	html_content = strings.ReplaceAll(html_content, "onclick=", "")
	html_content = strings.ReplaceAll(html_content, "onload=", "")
	
	return html_content
}

// GenerateSEOKeywords 自动生成SEO关键词
func GenerateSEOKeywords(title, summary string, tags []model.Tag) string {
	keywords := make(map[string]bool)
	
	// 从标题提取关键词（简化处理）
	words := splitChineseWords(title)
	for _, word := range words {
		if len(word) >= 2 && len(word) <= 10 {
			keywords[word] = true
		}
	}
	
	// 从摘要提取关键词
	words = splitChineseWords(summary)
	for _, word := range words {
		if len(word) >= 2 && len(word) <= 10 {
			keywords[word] = true
		}
	}
	
	// 添加标签
	for _, tag := range tags {
		keywords[tag.Name] = true
	}
	
	// 转换为逗号分隔的字符串
	result := ""
	for word := range keywords {
		if result != "" {
			result += ","
		}
		result += word
	}
	
	return result
}

// GenerateSEODescription 自动生成SEO描述
func GenerateSEODescription(title, summary string, limit int) string {
	if summary != "" {
		if len(summary) <= limit {
			return summary
		}
		return summary[:limit] + "..."
	}
	
	if title != "" {
		desc := "本文为您介绍：" + title
		if len(desc) <= limit {
			return desc
		}
		return desc[:limit] + "..."
	}
	
	return ""
}

// splitChineseWords 简单的中文分词（简化实现）
func splitChineseWords(text string) []string {
	// 实际项目中应该使用jieba等专业的中文分词库
	// 这里简化处理，按字符分割
	var words []string
	for _, r := range text {
		words = append(words, string(r))
	}
	return words
}

// ExtractSummary 从内容中提取摘要
func ExtractSummary(content string, limit int) string {
	// 移除Markdown标记
	summary := content
	summary = strings.ReplaceAll(summary, "#", "")
	summary = strings.ReplaceAll(summary, "*", "")
	summary = strings.ReplaceAll(summary, "_", "")
	summary = strings.ReplaceAll(summary, "`", "")
	summary = strings.ReplaceAll(summary, "[", "")
	summary = strings.ReplaceAll(summary, "]", "")
	summary = strings.ReplaceAll(summary, "(", "")
	summary = strings.ReplaceAll(summary, ")", "")
	
	// 按行分割，取第一行非空行
	lines := strings.Split(summary, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if len(line) <= limit {
				return line
			}
			return line[:limit] + "..."
		}
	}
	
	return ""
}
