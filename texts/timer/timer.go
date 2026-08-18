package timer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	Timer interface {
		Valid(origin string) bool
		GetDuration() time.Duration
		FromText(origin string) (Impl, error)
		FromTextDefault(origin string, def time.Duration) Impl
		ToText() string
	}
	Impl struct{ duration time.Duration }
)

var (
	Times = map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
		"w": time.Hour * 168, // (7*24)
	}

	// ChineseUnits 中文单位映射
	ChineseUnits = map[string]string{
		"w": "周",
		"d": "天",
		"h": "小时",
		"m": "分",
		"s": "秒",
	}

	// timeUnitRegexp 匹配单个时间单元: 数字+s/m/h/d/w
	timeUnitRegexp = regexp.MustCompile(`(\d+)([smhdw])`)

	// pureNumberRegexp 匹配纯数字（没有单位）
	pureNumberRegexp = regexp.MustCompile(`^\d+$`)
)

// parseTimeUnits 解析组合时间字符串，返回总 duration
// 支持：1w2d, 3d2h15m, 或纯数字（默认秒）
func parseTimeUnits(origin string) (time.Duration, error) {
	// 优先检查纯数字（默认按秒处理）
	if pureNumberRegexp.MatchString(origin) {
		num, err := strconv.Atoi(origin)
		if err != nil {
			return 0, fmt.Errorf("数字解析失败：%w", err)
		}
		return time.Duration(num) * time.Second, nil
	}

	matches := timeUnitRegexp.FindAllStringSubmatch(origin, -1)
	if len(matches) == 0 {
		return 0, errors.New("未找到有效的时间单元")
	}

	var totalDuration time.Duration
	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("时间单元数量解析失败：%w", err)
		}

		unitStr := match[2]
		unit, ok := Times[unitStr]
		if !ok {
			return 0, fmt.Errorf("时间单位不支持：%s", unitStr)
		}

		totalDuration += time.Duration(num) * unit
	}

	return totalDuration, nil
}

// toChinese 将 duration 格式化为人性化中文时间字符串
// 分级处理规则：
// 1. 从大到小排序：周、天、时、分、秒
// 2. 跳过为0的单位
// 3. 至少显示两个单位（如果有余数）
// 4. 只有一位数且大于等于10时，降级到下一单位（如 80分钟 → 1小时20分）
func toChinese(dur time.Duration) string {
	if dur <= 0 {
		return "0秒"
	}

	// 计算各单位的值
	weeks := int(dur.Hours()) / 168
	remainingHours := int(dur.Hours()) % 168
	days := remainingHours / 24
	remainingHours = remainingHours % 24
	hours := remainingHours
	minutes := int(dur.Minutes()) % 60
	seconds := int(dur.Seconds()) % 60

	// 分级处理：如果某一位的值 >= 10 且只有一位，降级到下一单位
	// 例如：80分钟 → 1小时20分，110小时 → 4天2小时
	if weeks > 0 && days >= 10 {
		// 将天转换为周+天
		extraWeeks := days / 7
		weeks += extraWeeks
		days = days % 7
	}

	if days > 0 && hours >= 10 {
		// 将小时转换为天+小时
		extraDays := hours / 24
		days += extraDays
		hours = hours % 24
	}

	if hours > 0 && minutes >= 10 {
		// 将分钟转换为时+分
		extraHours := minutes / 60
		hours += extraHours
		minutes = minutes % 60
	}

	if minutes > 0 && seconds >= 10 {
		// 将秒转换为分+秒
		extraMinutes := seconds / 60
		minutes += extraMinutes
		seconds = seconds % 60
	}

	// 构建结果字符串
	var result []string

	if weeks > 0 {
		result = append(result, fmt.Sprintf("%d周", weeks))
	}
	if days > 0 {
		result = append(result, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		result = append(result, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		result = append(result, fmt.Sprintf("%d分", minutes))
	}
	if seconds > 0 {
		result = append(result, fmt.Sprintf("%d秒", seconds))
	}

	// 如果结果为空（所有单位都是0）
	if len(result) == 0 {
		return "0秒"
	}

	// 如果只有一个单位且还有更小单位，确保至少显示两个单位
	if len(result) == 1 && dur > 0 {
		// 检查是否有遗漏的更小单位
		totalDisplayed := toDuration(weeks, days, hours, minutes, seconds)
		if totalDisplayed < dur {
			// 需要添加下一个更小的单位
			switch result[0] {
			case "周":
				result = append(result, fmt.Sprintf("%d天", days))
			case "天":
				result = append(result, fmt.Sprintf("%d小时", hours))
			case "小时":
				result = append(result, fmt.Sprintf("%d分", minutes))
			case "分":
				result = append(result, fmt.Sprintf("%d秒", seconds))
			}
		}
	}

	return strings.Join(result, "+")
}

// toDuration 将时间单位转换为 duration
func toDuration(weeks, days, hours, minutes, seconds int) time.Duration {
	return time.Duration(weeks)*7*24*time.Hour +
		time.Duration(days)*24*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second
}

// Valid 判断字符串是否符合组合时间格式（如 1w2d, 3d2h15m）或纯数字（默认秒）
func (*Impl) Valid(origin string) bool {
	origin = strings.ToLower(origin)
	if len(origin) == 0 {
		return false
	}

	// 检查是否以数字开头
	if origin[0] < '0' || origin[0] > '9' {
		return false
	}

	// 优先尝试匹配纯数字（默认按秒处理）
	if pureNumberRegexp.MatchString(origin) {
		return true
	}

	// 提取所有时间单元
	matches := timeUnitRegexp.FindAllStringSubmatch(origin, -1)
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

func (my *Impl) GetDuration() time.Duration { return my.duration }

// FromText 解析字符串并返回时间 duration（支持组合格式如 1w2d, 3d2h15m）
func (*Impl) FromText(origin string) (Impl, error) {
	var (
		err error
		t   = Impl{}
	)
	origin = strings.ToLower(origin)
	t.duration, err = parseTimeUnits(origin)

	return t, err
}

func (*Impl) FromTextDefault(origin string, def time.Duration) Impl {
	if !new(Impl).Valid(origin) {
		return Impl{duration: def}
	}

	dur, err := new(Impl).FromText(origin)
	if err != nil {
		return Impl{duration: def}
	}

	return dur
}

// ToText 公开接口：将 duration 转换为人可读的中文时间字符串
func (my *Impl) ToText() string { return toChinese(my.duration) }
