package str

import (
	"regexp"
	"slices"
)

type (
	Regexp struct {
		original string
		target   string
		targets  []string
		re       regexp.Regexp
	}
)

var (
	RegexpApp Regexp
)

func (Regexp) New(original string, attrs ...RegexpRegexpAttributer) Regexp {
	ins := Regexp{original: original, re: *regexp.MustCompile(original)}
	return ins.Set(attrs...)
}

func (my Regexp) Set(attrs ...RegexpRegexpAttributer) Regexp {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}

	return my
}

// MatchFirst 查找第一个匹配项
func (my Regexp) MatchFirst() string {
	matched := my.re.FindStringSubmatch(my.target)
	if len(matched) > 1 {
		return matched[1]
	}

	return ""
}

// MatchAll 查找所有匹配项
func (my Regexp) MatchAll() []string {
	matched := my.re.FindStringSubmatch(my.target)
	if len(matched) == 0 {
		return nil
	}

	if len(matched) > 1 {
		ret := make([]string, 0, len(matched)-1)
		ret = append(ret, matched[1:]...)
		return ret
	}

	return nil
}

// Contains 是否包含匹配项
func (my Regexp) Contains() bool { return my.re.MatchString(my.target) }

// ContainsAll 是否包含任意一个匹配项
func (my Regexp) ContainsAll() bool {
	if len(my.targets) == 0 {
		return false
	}

	return slices.ContainsFunc(my.targets, my.re.MatchString)
}

// ReplaceAllString 替换所有匹配项
func (my Regexp) ReplaceAllString(replace string) string {
	return my.re.ReplaceAllString(my.target, replace)
}
