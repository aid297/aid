package volumeInfo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
	PB = 1024 * TB
	EB = 1024 * PB
	ZB = 1024 * EB
	YB = 1024 * ZB
)

var (
	Volumes = map[string]float64{
		"KB": KB,
		"MB": MB,
		"GB": GB,
		"TB": TB,
		"PB": PB,
		"EB": EB,
		"ZB": ZB,
		"YB": YB,
	}

	// volumeRegexp 匹配: 任意数字 + KB/MB/GB/TB/PB/EB/ZB/YB
	volumeRegexp = regexp.MustCompile(`^(\d+(?:\.\d+)?)(KB|MB|GB|TB|PB|EB|ZB|YB)$`)
)

// IsValidVolumeFormat 判断字符串是否符合 数字+容量单位 格式
func IsValidVolumeFormat(origin string) bool {
	return volumeRegexp.MatchString(strings.ToUpper(origin))
}

// WhatVolumeIsIt 解析字符串并返回字节数
func WhatVolumeIsIt(origin string) (float64, error) {
	origin = strings.ToUpper(origin)

	matched := volumeRegexp.FindStringSubmatch(origin)
	if len(matched) != 3 {
		return 0, errors.New("格式不匹配")
	}

	num, err := strconv.ParseFloat(matched[1], 64)
	if err != nil {
		return 0, fmt.Errorf("格式不匹配：%w", err)
	}

	unit, ok := Volumes[matched[2]]
	if !ok {
		return 0, fmt.Errorf("容量单位不支持：%s", matched[2])
	}

	return num * unit, nil
}
