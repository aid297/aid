package textTime_test

import (
	"testing"
	"time"

	"github.com/aid297/aid/v2/consts/textTime"
)

func TestIsValidTimeFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 单一单位 - 支持
		{"单一秒", "1s", true},
		{"单一分", "5m", true},
		{"单一时", "2h", true},
		{"单一天", "3d", true},
		{"单一周", "1w", true},

		// 组合时间 - 新增支持
		{"周+天", "1w2d", true},
		{"天+时", "3d2h", true},
		{"时+分", "2h30m", true},
		{"分+秒", "15m30s", true},
		{"周+天+时", "1w2d3h", true},
		{"完整组合", "1w2d3h15m30s", true},
		{"多单位组合", "2w3d4h5m6s", true},

		// 无效格式 - 不支持
		{"空字符串", "", false},
		{"纯单位", "smt hdw", false},
		{"字母开头", "a1w2d", false},
		{"特殊字符", "1w2d!", false},
		{"小数", "1.5w", false},
		{"负数", "-1w", false},
		{"空格", "1w 2d", false},

		// 纯数字 - 默认按秒处理（新增）
		{"纯数字100", "100", true},
		{"纯数字1", "1", true},
		{"纯数字999", "999", true},
		{"前导零", "007", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textTime.IsValidTimeFormat(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidTimeFormat(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWhatTimeIsIt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectDur time.Duration
		wantErr   bool
	}{
		// 单一单位
		{"1秒", "1s", 1 * time.Second, false},
		{"5分", "5m", 5 * time.Minute, false},
		{"2时", "2h", 2 * time.Hour, false},
		{"3天", "3d", 3 * 24 * time.Hour, false},
		{"1周", "1w", 7 * 24 * time.Hour, false},

		// 组合时间
		{"1周2天", "1w2d", 9 * 24 * time.Hour, false},
		{"3天2时", "3d2h", (3*24 + 2) * time.Hour, false},
		{"2时30分", "2h30m", 2*time.Hour + 30*time.Minute, false},
		{"15分30秒", "15m30s", 15*time.Minute + 30*time.Second, false},
		{"1周2天3时", "1w2d3h", (9*24 + 3) * time.Hour, false},
		{"完整组合", "1w2d3h15m30s", (9*24+3)*time.Hour + 15*time.Minute + 30*time.Second, false},
		{"全单位组合", "1w2d3h4m5s", textTime.Times["w"]*1 + textTime.Times["d"]*2 + textTime.Times["h"]*3 + textTime.Times["m"]*4 + textTime.Times["s"]*5, false},

		// 错误情况
		{"无效格式", "abc", 0, true},
		{"空字符串", "", 0, true},
		{"不支持的单位", "1x", 0, true},

		// 纯数字 - 默认按秒处理（新增）
		{"纯数字60", "60", 60 * time.Second, false},
		{"纯数字100", "100", 100 * time.Second, false},
		{"纯数字3600", "3600", 3600 * time.Second, false},
		{"纯数字1", "1", 1 * time.Second, false},
		{"纯数字3661", "3661", 3661 * time.Second, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dur, err := textTime.WhatTimeIsIt(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WhatTimeIsIt(%q) 期望报错，但未报错", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("WhatTimeIsIt(%q) 意外报错: %v", tt.input, err)
				return
			}

			if dur != tt.expectDur {
				t.Errorf("WhatTimeIsIt(%q) = %v, 期望 %v", tt.input, dur, tt.expectDur)
			}
		})
	}
}

func TestWhatTimeIsItDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      time.Duration
		expected time.Duration
	}{
		{"有效格式", "1w2d", 1 * time.Hour, 9 * 24 * time.Hour},
		{"无效格式使用默认值", "invalid", 1 * time.Hour, 1 * time.Hour},
		{"空字符串使用默认值", "", 24 * time.Hour, 24 * time.Hour},
		// 纯数字不使用默认值（因为现在是合法格式）
		{"纯数字100无默认值", "100", 1 * time.Minute, 100 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textTime.WhatTimeIsItDefault(tt.input, tt.def)
			if result != tt.expected {
				t.Errorf("WhatTimeIsItDefault(%q, %v) = %v, 期望 %v", tt.input, tt.def, result, tt.expected)
			}
		})
	}
}

// Benchmark 性能测试
func BenchmarkIsValidTimeFormat(b *testing.B) {
	input := "1w2d3h15m30s"
	for i := 0; i < b.N; i++ {
		textTime.IsValidTimeFormat(input)
	}
}

func BenchmarkWhatTimeIsIt(b *testing.B) {
	input := "1w2d3h15m30s"
	for i := 0; i < b.N; i++ {
		_, _ = textTime.WhatTimeIsIt(input)
	}
}

func TestDurationToString(t *testing.T) {
	tests := []struct {
		name     string
		dur      time.Duration
		expected string
	}{
		// 基础测试
		{"0", 0, "0秒"},
		{"1秒", 1 * time.Second, "1秒"},
		{"59秒", 59 * time.Second, "59秒"},
		{"1分", 1 * time.Minute, "1分"},
		{"59分", 59 * time.Minute, "59分"},
		{"1时", 1 * time.Hour, "1小时"},
		{"23时", 23 * time.Hour, "23小时"},
		{"1天", 24 * time.Hour, "1天"},
		{"6天", 6 * 24 * time.Hour, "6天"},
		{"1周", 7 * 24 * time.Hour, "1周"},

		// 组合时间
		{"1小时20分", 1*time.Hour + 20*time.Minute, "1小时+20分"},
		{"2小时30分", 2*time.Hour + 30*time.Minute, "2小时+30分"},
		{"1天2时", 24*time.Hour + 2*time.Hour, "1天+2小时"},
		{"1周2天", 7*24*time.Hour + 2*24*time.Hour, "1周+2天"},
		{"1周2天3时", 7*24*time.Hour + 2*24*time.Hour + 3*time.Hour, "1周+2天+3小时"},
		{"1时2分3秒", 1*time.Hour + 2*time.Minute + 3*time.Second, "1小时+2分+3秒"},
		{"1周2天3时4分5秒", 7*24*time.Hour + 2*24*time.Hour + 3*time.Hour + 4*time.Minute + 5*time.Second, "1周+2天+3小时+4分+5秒"},

		// 分级处理测试
		// 80分钟 → 1小时20分
		{"80分钟分级", 80 * time.Minute, "1小时+20分"},
		// 90分钟 → 1小时30分
		{"90分钟分级", 90 * time.Minute, "1小时+30分"},
		// 150分钟 → 2小时30分
		{"150分钟分级", 150 * time.Minute, "2小时+30分"},
		// 25小时 → 1天1小时
		{"25小时分级", 25 * time.Hour, "1天+1小时"},
		// 30小时 → 1天6小时
		{"30小时分级", 30 * time.Hour, "1天+6小时"},
		// 10天 → 1周3天
		{"10天分级", 10 * 24 * time.Hour, "1周+3天"},
		// 17天 → 2周3天
		{"17天分级", 17 * 24 * time.Hour, "2周+3天"},

		// 边界情况：输入字符串没有单位（按秒处理）
		// 这些需要通过 WhatTimeIsIt 先转换，这里直接测 duration
		{"60秒", 60 * time.Second, "1分"},
		{"120秒", 120 * time.Second, "2分"},
		{"3600秒", 3600 * time.Second, "1小时"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textTime.DurationToString(tt.dur)
			if result != tt.expected {
				t.Errorf("DurationToString(%v) = %q, 期望 %q", tt.dur, result, tt.expected)
			}
		})
	}
}

func TestFormatInputString(t *testing.T) {
	// 测试完整流程：输入字符串 → 解析 → 格式化输出
	tests := []struct {
		name      string
		input     string
		expectOut string
		wantErr   bool
	}{
		// 基本流程
		{"1h20 -> 1小时20分", "1h20m", "1小时+20分", false},
		{"1w3 -> 1周+3天", "1w3d", "1周+3天", false},
		{"60 -> 1分", "60", "1分", false},
		{"100 -> 1分+40秒", "100", "1分+40秒", false},
		{"80m -> 1小时20分", "80m", "1小时+20分", false},
		{"10d -> 1周+3天", "10d", "1周+3天", false},

		// 复杂组合
		{"1w2d3h -> 1周+2天+3小时", "1w2d3h", "1周+2天+3小时", false},
		{"2h30m -> 2小时30分", "2h30m", "2小时+30分", false},
		{"90m -> 1小时30分", "90m", "1小时+30分", false},

		// 纯数字测试（新增）
		{"1 -> 1秒", "1", "1秒", false},
		{"3600 -> 1时", "3600", "1小时", false},
		{"3661 -> 1时1分1秒", "3661", "1小时+1分+1秒", false},
		{"59 -> 59秒", "59", "59秒", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dur, err := textTime.WhatTimeIsIt(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			output := textTime.DurationToString(dur)
			if output != tt.expectOut {
				t.Errorf("DurationToString(WhatTimeIsIt(%q)) = %q, 期望 %q", tt.input, output, tt.expectOut)
			}
		})
	}
}
