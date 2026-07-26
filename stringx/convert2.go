package stringx

import (
	"strings"
	"unicode"
)

type NamingRule string

const (
	NamingRuleSnakeCase  NamingRule = "snake_case"
	NamingRuleCamelCase  NamingRule = "camelCase"
	NamingRulePascalCase NamingRule = "PascalCase"
	NamingRuleKebabCase  NamingRule = "kebab-case"
	NamingRuleShoutCase  NamingRule = "SCREAMING_SNAKE_CASE"
)

// ConvertNaming 将原始字符串按指定的命名规则进行转换
// rule 支持: "snake_case", "camelCase", "PascalCase", "kebab-case", "SCREAMING_SNAKE_CASE"
func ConvertNaming(s string, rule NamingRule) string {
	words := splitIntoWords(s)
	if len(words) == 0 {
		return ""
	}

	switch rule {
	case "snake_case":
		return joinWords(words, "_", false, false)
	case "kebab-case":
		return joinWords(words, "-", false, false)
	case "SCREAMING_SNAKE_CASE":
		return strings.ToUpper(joinWords(words, "_", false, false))
	case "camelCase":
		return joinWords(words, "", true, false)
	case "PascalCase":
		return joinWords(words, "", false, true)
	default:
		// 默认返回原字符串或 snake_case，这里选择原样返回
		return s
	}
}

// splitIntoWords 智能拆分字符串，支持驼峰、下划线、连字符等混合情况
func splitIntoWords(s string) []string {
	var words []string
	var currentWord []rune

	runes := []rune(s)
	for i, r := range runes {
		if r == '_' || r == '-' || r == ' ' {
			// 遇到分隔符，保存当前单词
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
			continue
		}

		if unicode.IsUpper(r) {
			// 处理驼峰边界：
			// 1. 如果当前单词不为空，且前一个字符是小写，说明是新单词的开始 (e.g., camelCase -> camel, Case)
			// 2. 如果当前单词不为空，且前一个字符是大写，但下一个字符是小写，说明是缩写结束 (e.g., HTMLParser -> HTML, Parser)
			isCamelBoundary := len(currentWord) > 0 && unicode.IsLower(runes[i-1])
			isAcronymBoundary := len(currentWord) > 1 && unicode.IsUpper(runes[i-1]) && i+1 < len(runes) && unicode.IsLower(runes[i+1])

			if isCamelBoundary || isAcronymBoundary {
				words = append(words, string(currentWord))
				currentWord = nil
			}
		}
		currentWord = append(currentWord, r)
	}

	// 别忘了最后一个单词
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}
	return words
}

// joinWords 将单词数组按规则拼接
// sep: 分隔符
// firstLower: 第一个单词是否全小写 (camelCase 用)
// allTitle: 是否所有单词首字母大写 (PascalCase 用)
func joinWords(words []string, sep string, firstLower, allTitle bool) string {
	var sb strings.Builder
	for i, word := range words {
		if i > 0 && sep != "" {
			sb.WriteString(sep)
		}

		if len(word) == 0 {
			continue
		}

		runes := []rune(strings.ToLower(word))
		if allTitle || (i > 0 && !firstLower) {
			runes[0] = unicode.ToUpper(runes[0])
		}
		sb.WriteString(string(runes))
	}
	return sb.String()
}
