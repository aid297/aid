package gormPools

import "testing"

func TestValidateFieldName(t *testing.T) {
	allow := map[string]struct{}{"name": {}, "users.age": {}}

	cases := []struct {
		field string
		allow map[string]struct{}
		err   bool
	}{
		{"name", allow, false},
		{"users.age", allow, false},
		{"name", nil, false},
		{"id; drop table users", nil, true},
		{"name` = 'x'", nil, true},
		{"1name", nil, true},
		{"name--", nil, true},
		{"age", allow, true},      // 格式合法但不在白名单
		{"email", allow, true},    // 不在白名单
		{"users.age", nil, false}, // 无白名单时仅格式校验
	}

	for _, c := range cases {
		err := validateFieldName(c.field, c.allow)
		if (err != nil) != c.err {
			t.Errorf("validateFieldName(%q) err=%v, want err=%v", c.field, err, c.err)
		}
	}
}

func TestEscapeLikeValue(t *testing.T) {
	cases := []struct{ in, want string }{
		{"abc", "abc"},         // 无通配符，原样返回
		{"100%", `100\%`},      // % 转义
		{"a_b", `a\_b`},        // _ 转义
		{`a\b`, `a\\b`},        // \ 转义
		{"%", `\%`},            // 单个 %
		{"%", `\%`},            // 全通配符攻击样例
		{`abc\`, `abc\\`},      // 尾部反斜杠：原先会吃掉拼接的 %，转义后语义正确
		{`%\_%`, `\%\\\_\%`},   // 混合攻击样例：依次转义 %、\、_、%
		{"中文like值", "中文like值"}, // 非通配符字符不受影响
		{"", ""},               // 空值
	}

	for _, c := range cases {
		if got := EscapeLikeValue(c.in); got != c.want {
			t.Errorf("EscapeLikeValue(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidateOrder(t *testing.T) {
	allow := map[string]struct{}{"created_at": {}, "id": {}}

	cases := []struct {
		order string
		allow map[string]struct{}
		err   bool
	}{
		{"created_at desc", allow, false},
		{"id", allow, false},
		{"id ASC", nil, false},
		{"id asc, name desc", nil, true}, // 不允许多字段/复杂表达式
		{"id; drop table users", nil, true},
		{"id (select 1)", nil, true},
		{"id wrong", nil, true},    // 非法方向
		{"name desc", allow, true}, // 不在白名单
		{"", nil, true},            // 空排序
		{"   ", nil, true},         // 空白
	}

	for _, c := range cases {
		err := validateOrder(c.order, c.allow)
		if (err != nil) != c.err {
			t.Errorf("validateOrder(%q) err=%v, want err=%v", c.order, err, c.err)
		}
	}
}
