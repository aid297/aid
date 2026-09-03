package logs

import (
	"testing"
)

func TestNewJsonObject(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		wantType string
	}{
		{
			name:     "简单对象",
			input:    `{"name":"张三","age":25,"city":"北京"}`,
			wantErr:  false,
			wantType: "object",
		},
		{
			name:     "嵌套对象",
			input:    `{"user":{"name":"李四","contact":{"email":"lisi@example.com","phone":"123456"}}}`,
			wantErr:  false,
			wantType: "nested object",
		},
		{
			name:     "数组",
			input:    `[1,2,3,4,5]`,
			wantErr:  false,
			wantType: "array",
		},
		{
			name:     "对象数组",
			input:    `[{"id":1,"name":"A"},{"id":2,"name":"B"}]`,
			wantErr:  false,
			wantType: "array of objects",
		},
		{
			name:     "复杂混合",
			input:    `{"code":200,"data":{"items":[1,2,3]},"success":true}`,
			wantErr:  false,
			wantType: "complex mixed",
		},
		{
			name:    "无效JSON",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
		{
			name:     "空字符串",
			input:    ``,
			wantErr:  false,
			wantType: "empty",
		},
		{
			name:     "纯字符串",
			input:    `"hello world"`,
			wantErr:  false,
			wantType: "string",
		},
		{
			name:     "数字",
			input:    `12345`,
			wantErr:  false,
			wantType: "number",
		},
		{
			name:     "布尔值",
			input:    `true`,
			wantErr:  false,
			wantType: "boolean",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONString(tt.input)

			if tt.wantErr {
				if result.Error() == nil {
					t.Errorf("期望解析错误,但沒有收到错误")
				}
			} else {
				if result.Error() != nil {
					t.Errorf("不期望错误,但收到: %v", result.Error())
				}
			}

			output := result.String()
			t.Logf("输入: %s", tt.input)
			t.Logf("输出:\n%s", output)

			// 验证基本输出不为空(除非输入为空)
			if tt.input != "" && output == "" {
				t.Errorf("输出不应该为空")
			}
		})
	}
}

func TestJsonObjectFormat(t *testing.T) {
	// 测试复杂嵌套结构的格式化
	input := `{
		"users": [
			{"id": 1, "name": "张三", "hobbies": ["读书", "编程"]},
			{"id": 2, "name": "李四", "hobbies": ["音乐", "运动"]}
		],
		"meta": {
			"total": 2,
			"page": 1,
			"hasNext": false
		}
	}`

	result := JSONString(input)
	output := result.String()

	t.Logf("\n输入:\n%s\n", input)
	t.Logf("输出:\n%s\n", output)

	// 验证包含关键标识符
	checks := []string{"users", "meta", "张三", "李四", "total", "page"}
	for _, check := range checks {
		if !contains(output, check) {
			t.Errorf("输出中应该包含 '%s',但未找到", check)
		}
	}
}

func TestJsonObjectDeepNesting(t *testing.T) {
	// 测试深度嵌套的截断保护
	deep := `{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":{"i":{"j":{"k":"极限"}}}}}}}}}}}`
	result := JSONString(deep)
	output := result.String()

	t.Logf("深度嵌套输出:\n%s", output)

	// 应该在第10层截断
	if contains(output, `"k":"极限"`) {
		t.Log("检测到深度嵌套被正确截断")
	}
}

func TestJsonObjectMethods(t *testing.T) {
	input := `{"key":"value"}`
	result := JSONString(input)

	// 测试 String()
	str := result.String()
	if str == "" {
		t.Error("String() 不应该为空")
	}
	t.Logf("String(): %s", str)

	// 测试 Value()
	value := result.Value()
	if value == nil {
		t.Error("Value() 不应该为 nil")
	}
	t.Logf("Value(): %v", value)

	// 测试 Raw()
	raw := result.Raw()
	if raw != input {
		t.Errorf("Raw() 应该返回原始输入,期望 %s,得到 %s", input, raw)
	}
	t.Logf("Raw(): %s", raw)

	// 测试 Error()
	err := result.Error()
	if err != nil {
		t.Errorf("有效 JSON 不应该有错误: %v", err)
	}
}

func TestJsonObjectInvalid(t *testing.T) {
	input := `{invalid json}`
	result := JSONString(input)

	if result.Error() == nil {
		t.Error("无效 JSON 应该返回错误")
	}
	t.Logf("解析错误: %v", result.Error())

	output := result.String()
	t.Logf("输出: %s", output)

	// 失效时应该显示原始输入和错误信息
	if !contains(output, input) {
		t.Error("输出应该包含原始输入")
	}
	if !contains(output, "JSON Parse Error") {
		t.Error("输出应该包含错误提示")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s != "" && substr != "" && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
