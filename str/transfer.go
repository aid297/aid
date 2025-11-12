package str

import (
	"regexp"
	"strings"
	"unicode"
)

type Transfer struct{ original string }

var TransferApp Transfer

func (*Transfer) New(original string) *Transfer { return &Transfer{original: original} }

// NewTransfer 实例化：字符串转换
//
//go:fix 推荐使用：New方法
func NewTransfer(original string) *Transfer { return &Transfer{original: original} }

// PascalToCamel 大驼峰 -> 小驼峰
func (my *Transfer) PascalToCamel() string {
	if len(my.original) == 0 {
		return my.original
	}
	// 将第一个字符转换为小写
	firstRune := []rune(my.original)[0]
	if unicode.IsUpper(firstRune) {
		firstRune = unicode.ToLower(firstRune)
	}

	// 拼接第一个字符和剩余部分
	return string(firstRune) + my.original[1:]
}

// PascalToSnake 大驼峰 -> 下划线
func (my *Transfer) PascalToSnake() string {
	var result strings.Builder

	for i, r := range my.original {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// PascalToBabel 大驼峰 -> babel
func (my *Transfer) PascalToBabel() string {
	var result strings.Builder

	for i, r := range my.original {
		if unicode.IsUpper(r) && i > 0 {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}

	return result.String()
}

// CamelToPascal 小驼峰 -> 大驼峰
func (my *Transfer) CamelToPascal() string {
	if len(my.original) == 0 {
		return my.original
	}
	firstRune := []rune(my.original)[0]
	if unicode.IsLower(firstRune) {
		firstRune = unicode.ToUpper(firstRune)
	}

	return string(firstRune) + my.original[1:]
}

// CamelToSnake 小驼峰 -> 下划线
func (my *Transfer) CamelToSnake() string {
	var result strings.Builder

	for idx, ite := range my.original {
		if unicode.IsUpper(ite) && idx > 0 {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(ite))
	}

	return result.String()
}

// CamelToBabel 小驼峰 -> babel
func (my *Transfer) CamelToBabel() string {
	var result strings.Builder

	for idx, ite := range my.original {
		if unicode.IsUpper(ite) && idx > 0 {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(ite))
	}

	return result.String()
}

// SnakeToPascal 下划线 -> 大驼峰
func (my *Transfer) SnakeToPascal() string {
	// 将下划线分割成单词
	words := strings.Split(my.original, "_")

	// 处理每个单词，将首字母大写
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	// 拼接所有单词
	pascal := strings.Join(words, "")

	return pascal
}

// SnakeToCamel 下划线 -> 小驼峰
func (my *Transfer) SnakeToCamel() string {
	// 将下划线分割成单词
	words := strings.Split(my.original, "_")

	// 处理每个单词，将首字母大写
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToLower(runes[0])
			words[i] = string(runes)
		}
	}

	// 拼接所有单词
	pascal := strings.Join(words, "")

	return pascal
}

// SnakeToBabel 下划线 -> babel
func (my *Transfer) SnakeToBabel() string {
	return strings.ReplaceAll(my.original, "_", "-")
}

// BabelToPascal babel -> 大驼峰
func (my *Transfer) BabelToPascal() string {
	words := strings.Split(my.original, "-")
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	pascal := strings.Join(words, "")

	return pascal
}

// KebabToCamel babel -> 小驼峰
func (my *Transfer) KebabToCamel() string {
	// 将 kebab-case 分割成单词
	words := strings.Split(my.original, "-")

	// 处理每个单词，除了第一个单词外，将每个单词的首字母大写
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			runes := []rune(words[i])
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}

	// 拼接所有单词
	camel := strings.Join(words, "")

	// 确保第一个字符是小写
	if len(camel) > 0 {
		runes := []rune(camel)
		runes[0] = unicode.ToLower(runes[0])
		camel = string(runes)
	}

	return camel
}

// BabelToSnake babel -> 下划线
func (my *Transfer) BabelToSnake() string { return strings.ReplaceAll(my.original, "_", "-") }

// Pluralize 单数变复数
func (my *Transfer) Pluralize() string {
	// 定义正则表达式
	sXChSh := regexp.MustCompile(`[sxz]|[cs]h$`)
	yEnding := regexp.MustCompile(`[^aeiou]y$`)
	fFeEnding := regexp.MustCompile(`[f]e?$`)

	// 处理以 s, x, ch, sh 结尾的名词
	if sXChSh.MatchString(my.original) {
		return my.original + "es"
	}

	// 处理以辅音字母 + y 结尾的名词
	if yEnding.MatchString(my.original) {
		return yEnding.ReplaceAllString(my.original, "ies")
	}

	// 处理以 f 或 fe 结尾的名词
	if fFeEnding.MatchString(my.original) {
		return fFeEnding.ReplaceAllString(my.original, "ves")
	}

	// 默认情况下，直接加 s
	return my.original + "s"
}
