package regexp

import (
	"regexp"
	"slices"
)

type Regexp struct {
	target  string
	targets []string
	re      regexp.Regexp
}

var RegexpApp Regexp

func (Regexp) New(original string, attrs ...Attributer) Regexp {
	return Regexp{re: *regexp.MustCompile(original)}.SetAttrs(attrs...)
}

func (my Regexp) SetAttrs(attrs ...Attributer) Regexp {
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

// ContainsAny 是否包含任意一个匹配项
func (my Regexp) ContainsAny() bool {
	if len(my.targets) == 0 {
		return false
	}

	return slices.ContainsFunc(my.targets, my.re.MatchString)
}

// FindAllStringSubmatch 查找所有匹配项及子匹配项
func (my Regexp) FindAllStringSubmatch(matchIdx int) []string {
	var msg []string
	matches := my.re.FindAllStringSubmatch(my.target, -1)
	for idx := range matches {
		if matchIdx == -1 {
			msg = append(msg, matches[idx]...)
		} else {
			msg = append(msg, matches[idx][matchIdx])
		}
	}

	return msg
}
