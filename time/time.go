package time

import (
	"fmt"
	"time"
)

// Time 时间帮助
type Time struct{ original time.Time }

var App Time

const (
	Nanosecond  time.Duration = 1
	Microsecond               = 1000 * Nanosecond
	Millisecond               = 1000 * Microsecond
	Second                    = 1000 * Millisecond
	Minute                    = 60 * Second
	Hour                      = 60 * Minute

	Layout      = time.Layout
	ANSIC       = time.ANSIC
	UnixDate    = time.UnixDate
	RubyDate    = time.RubyDate
	RFC822      = time.RFC822
	RFC822Z     = time.RFC822Z
	RFC850      = time.RFC850
	RFC1123     = time.RFC1123
	RFC1123Z    = time.RFC1123Z
	RFC3339     = time.RFC3339
	RFC3339Nano = time.RFC3339Nano
	Kitchen     = time.Kitchen
	Stamp       = time.Stamp
	StampMilli  = time.StampMilli
	StampMicro  = time.StampMicro
	StampNano   = time.StampNano
	DateTime    = time.DateTime
	DateOnly    = time.DateOnly
	TimeOnly    = time.TimeOnly
)

// New 实例化：时间帮助
func (*Time) New(t time.Time) *Time { return &Time{original: t} }

// Now 实例化：时间帮助 -> 当前时间
func (*Time) Now() *Time { return &Time{original: time.Now()} }

// FromFormat 实例化：时间帮助 -> 通过格式化模板
func (*Time) FromFormat(format, t string) (*Time, error) {
	parsedTime, err := time.Parse(format, t)
	if err != nil {
		return nil, err
	}
	return &Time{original: parsedTime}, nil
}

// FormString 实例化：时间帮助 -> 通过时间字符串自动解析
func (*Time) FormString(t string) (*Time, error) {
	formats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, t)
		if err == nil {
			return &Time{original: parsedTime}, nil
		}
	}

	return nil, fmt.Errorf("无法解析时间：%s", t)
}

// AddDays 增加时间
func (my *Time) AddDays(days int) *Time {
	my.original = my.original.AddDate(0, 0, days)

	return my
}

// AddMonths 增加月份
func (my *Time) AddMonths(months int) *Time {
	my.original = my.original.AddDate(0, months, 0)

	return my
}

// AddYears 增加年份
func (my *Time) AddYears(years int) *Time {
	my.original = my.original.AddDate(years, 0, 0)

	return my
}

// Format 格式化
func (my *Time) Format(format string) string { return my.original.Format(format) }

// Diff 计算两个时间之间的差值
func (my *Time) Diff(other *Time) time.Duration { return my.original.Sub(other.original) }

// IsBefore 判断是否早于某时间
func (my *Time) IsBefore(other *Time) bool { return my.original.Before(other.original) }

// IsAfter 判断是否晚于某时间
func (my *Time) IsAfter(other *Time) bool { return my.original.After(other.original) }

// IsEqual 判断是否等于某时间
func (my *Time) IsEqual(other *Time) bool { return my.original.Equal(other.original) }

// IsZero 判断是否是0值
func (my *Time) IsZero() bool { return my.original.IsZero() }

// ToTime 转换为时间
func (my *Time) ToTime() time.Time { return my.original }

// ToTimePtr 转换为时间指针
func (my *Time) ToTimePtr() *time.Time { return &my.original }
