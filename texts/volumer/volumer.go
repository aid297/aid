package volumer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	Volumer interface {
		Valid(origin string) bool
		FromText(origin string) (float64, error)
		FromTextDefault(origin string, def float64) float64
	}
	VolumerImpl struct{}
)

const (
	B  float64 = 1
	KB         = 1024 * B
	MB         = 1024 * KB
	GB         = 1024 * MB
	TB         = 1024 * GB
	PB         = 1024 * TB
	EB         = 1024 * PB
	ZB         = 1024 * EB
	YB         = 1024 * ZB
)

var (
	Volumes = map[string]float64{
		"B": B,
		"K": KB,
		"M": MB,
		"G": GB,
		"T": TB,
		"P": PB,
		"E": EB,
		"Z": ZB,
		"Y": YB,
	}

	// volumeUnitRegexp 匹配单个容量单元: 数字+B/K/M/G/T/P/E/Z/Y
	volumeUnitRegexp = regexp.MustCompile(`(\d+(?:\.\d+)?)([BKMGTPEZY])`)

	// pureNumberRegexp 匹配纯数字（没有单位）
	pureNumberRegexp = regexp.MustCompile(`^\d+(?:\.\d+)?$`)
)

// Valid 判断字符串是否符合组合容量格式（如 1m2k, 3g500m）或纯数字（默认byte）
func (*VolumerImpl) Valid(origin string) bool {
	origin = strings.ToUpper(origin)
	if len(origin) == 0 {
		return false
	}

	// 检查是否以数字开头
	if origin[0] < '0' || origin[0] > '9' {
		return false
	}

	// 优先尝试匹配纯数字（默认按byte处理）
	if pureNumberRegexp.MatchString(origin) {
		return true
	}

	// 提取所有容量单元
	matches := volumeUnitRegexp.FindAllStringSubmatch(origin, -1)
	if len(matches) == 0 {
		return false
	}

	// 验证格式是否正确（重新拼接后应等于原字符串）
	var rebuilt strings.Builder
	for _, match := range matches {
		rebuilt.WriteString(match[0])
	}
	return rebuilt.String() == origin
}

// parseVolumeUnits 解析组合容量字符串，返回总字节数
// 支持：1m2k, 3g500m, 或纯数字（默认byte）
func (v *VolumerImpl) parseVolumeUnits(origin string) (float64, error) {
	origin = strings.ToUpper(origin)

	// 优先检查纯数字（默认按byte处理）
	if pureNumberRegexp.MatchString(origin) {
		num, err := strconv.ParseFloat(origin, 64)
		if err != nil {
			return 0, fmt.Errorf("数字解析失败：%w", err)
		}
		return num * B, nil
	}

	matches := volumeUnitRegexp.FindAllStringSubmatch(origin, -1)
	if len(matches) == 0 {
		return 0, errors.New("未找到有效容量单元")
	}

	var totalVolume float64
	for _, match := range matches {
		num, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, fmt.Errorf("容量单元数量解析失败：%w", err)
		}

		unitStr := match[2]
		unit, ok := Volumes[unitStr]
		if !ok {
			return 0, fmt.Errorf("容量单位不支持：%s", unitStr)
		}

		totalVolume += num * unit
	}

	return totalVolume, nil
}

// FromText 解析字符串并返回字节数（支持组合格式如 1m2k, 3g500m）
func (v *VolumerImpl) FromText(origin string) (float64, error) {
	return v.parseVolumeUnits(origin)
}

func (*VolumerImpl) FromTextDefault(origin string, def float64) float64 {
	if !new(VolumerImpl).Valid(origin) {
		return def
	}

	vol, err := new(VolumerImpl).FromText(origin)
	if err != nil {
		return def
	}

	return vol
}
