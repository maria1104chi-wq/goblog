package utils

import (
	"bufio"
	"blog-backend/internal/model"
	"blog-backend/internal/repository"
	"fmt"
	"os"
	"strings"
	"sync"
)

// SensitiveWordFilter 敏感词过滤器
type SensitiveWordFilter struct {
	words     map[string]bool
	wordList  []string
	trie      *TrieNode
	mu        sync.RWMutex
}

// TrieNode Trie树节点
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// NewSensitiveWordFilter 创建敏感词过滤器实例
func NewSensitiveWordFilter() *SensitiveWordFilter {
	return &SensitiveWordFilter{
		words:    make(map[string]bool),
		wordList: make([]string, 0),
		trie:     &TrieNode{children: make(map[rune]*TrieNode)},
	}
}

// LoadFromFile 从文件加载敏感词
func (f *SensitiveWordFilter) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开敏感词文件失败：%w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取敏感词文件失败：%w", err)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 清空原有数据
	f.words = make(map[string]bool)
	f.wordList = make([]string, 0, len(words))
	f.trie = &TrieNode{children: make(map[rune]*TrieNode)}

	// 添加新词
	for _, word := range words {
		f.words[word] = true
		f.wordList = append(f.wordList, word)
		f.insertToTrie(word)
	}

	return nil
}

// LoadFromDatabase 从数据库加载敏感词
func (f *SensitiveWordFilter) LoadFromDatabase() error {
	var words []model.SensitiveWord
	err := repository.GetDB().Where("status = ?", 1).Find(&words).Error
	if err != nil {
		return fmt.Errorf("查询敏感词失败：%w", err)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// 清空原有数据
	f.words = make(map[string]bool)
	f.wordList = make([]string, 0, len(words))
	f.trie = &TrieNode{children: make(map[rune]*TrieNode)}

	// 添加新词
	for _, word := range words {
		f.words[word.Word] = true
		f.wordList = append(f.wordList, word.Word)
		f.insertToTrie(word.Word)
	}

	return nil
}

// insertToTrie 向Trie树插入敏感词
func (f *SensitiveWordFilter) insertToTrie(word string) {
	node := f.trie
	for _, char := range word {
		if _, exists := node.children[char]; !exists {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Filter 过滤文本中的敏感词，替换为*
func (f *SensitiveWordFilter) Filter(text string) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.words) == 0 {
		return text
	}

	result := []rune(text)
	n := len(result)

	i := 0
	for i < n {
		node := f.trie
		matchEnd := -1
		j := i

		for j < n {
			char := result[j]
			if child, exists := node.children[char]; exists {
				node = child
				if node.isEnd {
					matchEnd = j + 1
				}
				j++
			} else {
				break
			}
		}

		if matchEnd > i {
			// 找到敏感词，替换为*
			for k := i; k < matchEnd; k++ {
				result[k] = '*'
			}
			i = matchEnd
		} else {
			i++
		}
	}

	return string(result)
}

// Contains 检查文本是否包含敏感词
func (f *SensitiveWordFilter) Contains(text string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.words) == 0 {
		return false
	}

	runes := []rune(text)
	n := len(runes)

	for i := 0; i < n; i++ {
		node := f.trie
		j := i

		for j < n {
			char := runes[j]
			if child, exists := node.children[char]; exists {
				node = child
				if node.isEnd {
					return true
				}
				j++
			} else {
				break
			}
		}
	}

	return false
}

// GetWordCount 获取敏感词数量
func (f *SensitiveWordFilter) GetWordCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.words)
}

// AddWord 添加单个敏感词
func (f *SensitiveWordFilter) AddWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.words[word]; !exists {
		f.words[word] = true
		f.wordList = append(f.wordList, word)
		f.insertToTrie(word)
	}
}

// RemoveWord 移除敏感词
func (f *SensitiveWordFilter) RemoveWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.words[word]; exists {
		delete(f.words, word)
		
		// 重建wordList
		newList := make([]string, 0, len(f.words))
		for w := range f.words {
			newList = append(newList, w)
		}
		f.wordList = newList
		
		// 重建Trie树
		f.trie = &TrieNode{children: make(map[rune]*TrieNode)}
		for _, w := range f.wordList {
			f.insertToTrie(w)
		}
	}
}

// GetWordList 获取敏感词列表
func (f *SensitiveWordFilter) GetWordList() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	result := make([]string, len(f.wordList))
	copy(result, f.wordList)
	return result
}

// GlobalFilter 全局敏感词过滤器实例
var GlobalFilter *SensitiveWordFilter

// InitGlobalFilter 初始化全局敏感词过滤器
func InitGlobalFilter() error {
	GlobalFilter = NewSensitiveWordFilter()
	
	// 先从数据库加载
	if err := GlobalFilter.LoadFromDatabase(); err != nil {
		fmt.Printf("从数据库加载敏感词失败：%v，尝试从文件加载\n", err)
		
		// 如果数据库加载失败，尝试从文件加载
		if err := GlobalFilter.LoadFromFile("static/Sensitive.txt"); err != nil {
			return fmt.Errorf("初始化敏感词过滤器失败：%w", err)
		}
	}
	
	fmt.Printf("敏感词过滤器初始化完成，共加载 %d 个敏感词\n", GlobalFilter.GetWordCount())
	return nil
}

// FilterText 使用全局过滤器过滤文本
func FilterText(text string) string {
	if GlobalFilter == nil {
		return text
	}
	return GlobalFilter.Filter(text)
}

// ContainsSensitiveWord 使用全局过滤器检查是否包含敏感词
func ContainsSensitiveWord(text string) bool {
	if GlobalFilter == nil {
		return false
	}
	return GlobalFilter.Contains(text)
}
