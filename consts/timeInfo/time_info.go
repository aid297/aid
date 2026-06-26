package timeInfo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	Times = map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
		"w": time.Hour * 168, // (7*24)
	}

	// timeRegexp 匹配: 任意正整数 + s/m/h/d/w
	timeRegexp = regexp.MustCompile(`^(\d+)([smhdw])$`)
)

// IsValidTimeFormat 判断字符串是否符合 数字+s/m/h/d/w 格式
func IsValidTimeFormat(origin string) bool {
	return timeRegexp.MatchString(strings.ToLower(origin))
}

// WhatTimeIsIt 解析字符串并返回时间 duration
func WhatTimeIsIt(origin string) (time.Duration, error) {
	origin = strings.ToLower(origin)

	matched := timeRegexp.FindStringSubmatch(origin)
	if len(matched) != 3 {
		return 0, errors.New("格式不匹配")
	}

	num, err := strconv.Atoi(matched[1])
	if err != nil {
		return 0, fmt.Errorf("格式不匹配：%w", err)
	}

	unit, ok := Times[matched[2]]
	if !ok {
		return 0, fmt.Errorf("时间单位不支持：%s", matched[2])
	}

	return time.Duration(num) * unit, nil
}

func WhatTimeIsItDefault(origin string, def time.Duration) time.Duration {
	if !IsValidTimeFormat(origin) {
		return def
	}

	dur, err := WhatTimeIsIt(origin)
	if err != nil {
		return def
	}

	return dur
}
