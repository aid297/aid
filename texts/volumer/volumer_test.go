package volumer_test

import (
	"testing"

	"github.com/aid297/aid/v3/texts/volumer"
)

func TestIsValidVolumeFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 单一单位 - 支持
		{"单一byte", "100B", true},
		{"单一KB", "5K", true},
		{"单一MB", "2M", true},
		{"单一GB", "3G", true},
		{"单一TB", "1T", true},
		{"单一PB", "1P", true},
		{"单一EB", "1E", true},
		{"单一ZB", "1Z", true},
		{"单一YB", "1Y", true},

		// 组合容量 - 新增支持
		{"MB+KB", "1m2k", true},
		{"GB+MB", "3g500m", true},
		{"GB+MB+KB", "2g500m100k", true},
		{"TB+GB", "1t500g", true},
		{"完整组合", "1t500g200m10k", true},
		{"小数支持", "1.5g", true},
		{"多个小数", "1.5g2.5m", true},

		// 纯数字 - 默认按byte处理（新增）
		{"纯数字100", "100", true},
		{"纯数字1", "1", true},
		{"纯数字999", "999", true},
		{"纯小数1.5", "1.5", true},
		{"纯小数3.14", "3.14", true},

		// 无效格式 - 不支持
		{"空字符串", "", false},
		{"纯字母", "abc", false},
		{"字母开头", "a100", false},
		{"特殊字符", "100!", false},
		{"空格", "1g 2m", false},
		{"负数", "-100", false},
		{"科学计数法", "1e10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := new(volumer.VolumerImpl).Valid(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidVolumeFormat(%q) = %v, 期望 %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWhatVolumeIsIt(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectVol float64
		wantErr   bool
	}{
		// 单一单位
		{"100B", "100B", 100, false},
		{"1K", "1K", 1024, false},
		{"1M", "1M", 1024 * 1024, false},
		{"1G", "1G", 1024 * 1024 * 1024, false},
		{"1T", "1T", 1024 * 1024 * 1024 * 1024, false},
		{"1P", "1P", 1024 * 1024 * 1024 * 1024 * 1024, false},

		// 组合容量
		{"1m2k -> 1MB+2KB", "1m2k", 1*volumer.MB + 2*volumer.KB, false},
		{"3g500m -> 3GB+500MB", "3g500m", 3*volumer.GB + 500*volumer.MB, false},
		{"2g500m100k -> 2GB+500MB+100KB", "2g500m100k", 2*volumer.GB + 500*volumer.MB + 100*volumer.KB, false},
		{"1t500g -> 1TB+500GB", "1t500g", 1*volumer.TB + 500*volumer.GB, false},

		// 纯数字 - 默认按byte处理（新增）
		{"纯数字100", "100", 100, false},
		{"纯数字1", "1", 1, false},
		{"纯数字3600", "3600", 3600, false},

		// 小数支持
		{"1.5g -> 1.5GB", "1.5g", 1.5 * volumer.GB, false},
		{"2.5m -> 2.5MB", "2.5m", 2.5 * volumer.MB, false},
		{"1.5g2.5m -> 1.5GB+2.5MB", "1.5g2.5m", 1.5*volumer.GB + 2.5*volumer.MB, false},

		// 错误情况
		{"无效格式", "abc", 0, true},
		{"空字符串", "", 0, true},
		{"不支持的单位", "1X", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vol, err := new(volumer.VolumerImpl).FromText(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WhatVolumeIsIt(%q) 期望报错，但未报错", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("WhatVolumeIsIt(%q) 意外报错: %v", tt.input, err)
				return
			}

			if vol != tt.expectVol {
				t.Errorf("WhatVolumeIsIt(%q) = %v, 期望 %v", tt.input, vol, tt.expectVol)
			}
		})
	}
}

func TestWhatVolumeIsItDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		def      float64
		expected float64
	}{
		{"有效格式", "1g", 1 * volumer.MB, 1 * volumer.GB},
		{"无效格式使用默认值", "invalid", 1 * volumer.MB, 1 * volumer.MB},
		{"空字符串使用默认值", "", 1024, 1024},
		// 纯数字现在是合法格式
		{"纯数字100无默认值", "100", 1 * volumer.KB, 100},
		{"小数1.5g无默认值", "1.5g", 1 * volumer.MB, 1.5 * volumer.GB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := new(volumer.VolumerImpl).FromTextDefault(tt.input, tt.def)
			if result != tt.expected {
				t.Errorf("WhatVolumeIsItDefault(%q, %v) = %v, 期望 %v", tt.input, tt.def, result, tt.expected)
			}
		})
	}
}
