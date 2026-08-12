package validators_test

import (
	"strings"
	"testing"
	"time"

	"github.com/aid297/aid/v2/validators"
)

// result 执行校验并返回结果
func result(t *testing.T, data any, checkers ...validators.Checker) (invalid bool, errs []error) {
	t.Helper()
	v := validators.WithData(data).Validate(checkers...)
	return v.Invalid(), v.GetErrors()
}

// containsErr 判断错误列表中是否包含指定关键字
func containsErr(errs []error, sub string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error(), sub) {
			return true
		}
	}
	return false
}

type User struct {
	Name    string
	Age     int
	Email   *string
	Country string
	Balance float64
	Level   uint
}

func strPtr(s string) *string { return &s }

// TestValidator_AllPass 所有规则均满足时校验通过
func TestValidator_AllPass(t *testing.T) {
	email := "zhang@example.com"

	invalid, errs := result(t, User{
		Name: "张三丰", Age: 18, Email: &email,
		Country: "CN", Balance: 999.9, Level: 2,
	},
		validators.NewChecker("Name", "姓名").Required().Min(">0").Max("<10"),
		validators.NewChecker("Age", "年龄").Min(">0").Max("<200"),
		validators.NewChecker("Email", "邮箱").Required().Regex(`^[^@]+@[^@]+$`),
		validators.NewChecker("Country", "国家").In("==CN|US|JP"),
		validators.NewChecker("Balance", "余额").Min(">0"),
		validators.NewChecker("Level", "等级").Min(">=1"),
	)
	if invalid {
		t.Errorf("期望校验通过，实际错误：%v", errs)
	}
}

// TestValidator_Required 必填字段为空时报错
func TestValidator_Required(t *testing.T) {
	invalid, errs := result(t, User{},
		validators.NewChecker("Name", "姓名").Required(),
	)
	if !invalid || !containsErr(errs, "『姓名』必填") {
		t.Errorf("期望必填错误，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// TestValidator_StringRules 字符串的长度/枚举/正则规则
func TestValidator_StringRules(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		check   func() validators.Checker
		wantSub string // 为空表示期望通过
	}{
		{"Min失败", "ab", func() validators.Checker { return validators.NewChecker("Name", "姓名").Min(">2") }, "『姓名』长度不能小于 2"},
		{"Min通过", "abc", func() validators.Checker { return validators.NewChecker("Name", "姓名").Min(">2") }, ""},
		{"Max失败", "abcdefghij", func() validators.Checker { return validators.NewChecker("Name", "姓名").Max("<10") }, "『姓名』长度不能大于 10"},
		{"Max通过", "abc", func() validators.Checker { return validators.NewChecker("Name", "姓名").Max("<10") }, ""},
		{"Size失败", "abcd", func() validators.Checker { return validators.NewChecker("Name", "姓名").Size("==5") }, "『姓名』长度需要等于 5"},
		{"Size通过", "abcde", func() validators.Checker { return validators.NewChecker("Name", "姓名").Size("==5") }, ""},
		{"In失败", "x", func() validators.Checker { return validators.NewChecker("Country", "国家").In("==CN|US|JP") }, "『国家』内容必须在 CN,US,JP 中"},
		{"In通过", "US", func() validators.Checker { return validators.NewChecker("Country", "国家").In("==CN|US|JP") }, ""},
		{"In!=失败", "CN", func() validators.Checker { return validators.NewChecker("Country", "国家").In("!=CN|US") }, "『国家』内容不能在 CN,US 中"},
		{"Regex失败", "ABC", func() validators.Checker { return validators.NewChecker("Name", "姓名").Regex("^[a-z]+$") }, "『姓名』内容必须符合"},
		{"Regex通过", "abc", func() validators.Checker { return validators.NewChecker("Name", "姓名").Regex("^[a-z]+$") }, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := struct{ Name, Country string }{Name: c.value, Country: c.value}
			invalid, errs := result(t, data, c.check())

			if c.wantSub == "" {
				if invalid {
					t.Errorf("期望通过，实际错误：%v", errs)
				}
				return
			}
			if !invalid || !containsErr(errs, c.wantSub) {
				t.Errorf("期望错误 %q，实际：invalid=%v errs=%v", c.wantSub, invalid, errs)
			}
		})
	}
}

// TestValidator_IntRules 整数的最小值/最大值/等值规则
func TestValidator_IntRules(t *testing.T) {
	// min 失败
	invalid, errs := result(t, struct{ Age int }{Age: 2},
		validators.NewChecker("Age", "年龄").Min(">3"))
	if !invalid || !containsErr(errs, "『年龄』值不能小于 3") {
		t.Errorf("期望 min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// max 失败
	invalid, errs = result(t, struct{ Age int }{Age: 150},
		validators.NewChecker("Age", "年龄").Max("<100"))
	if !invalid || !containsErr(errs, "『年龄』值不能大于 100") {
		t.Errorf("期望 max 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// size 失败
	invalid, errs = result(t, struct{ Age int }{Age: 20},
		validators.NewChecker("Age", "年龄").Size("==18"))
	if !invalid || !containsErr(errs, "『年龄』值需要等于 18") {
		t.Errorf("期望 size 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 全部通过
	invalid, errs = result(t, struct{ Age int }{Age: 18},
		validators.NewChecker("Age", "年龄").Min(">3").Max("<100").Size("==18"))
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}
}

// TestValidator_UintFloatRules 无符号整数与浮点数规则
func TestValidator_UintFloatRules(t *testing.T) {
	// uint min 失败
	invalid, errs := result(t, struct{ Level uint }{Level: 0},
		validators.NewChecker("Level", "等级").Min(">0"))
	if !invalid || len(errs) == 0 {
		t.Errorf("期望 uint min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// float min 失败
	invalid, errs = result(t, struct{ Balance float64 }{Balance: 100},
		validators.NewChecker("Balance", "余额").Min(">100.5"))
	if !invalid || len(errs) == 0 {
		t.Errorf("期望 float min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 全部通过
	invalid, errs = result(t, struct {
		Level   uint
		Balance float64
	}{Level: 1, Balance: 200},
		validators.NewChecker("Level", "等级").Min(">0"),
		validators.NewChecker("Balance", "余额").Min(">100.5"))
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}
}

// TestValidator_PointerField 指针字段的必填与取值
func TestValidator_PointerField(t *testing.T) {
	// nil 指针 + required → 报错
	invalid, errs := result(t, User{},
		validators.NewChecker("Email", "邮箱").Required())
	if !invalid || !containsErr(errs, "『邮箱』必填") {
		t.Errorf("期望 nil 指针必填错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 非 nil 指针按实际值校验
	email := "x@y.com"
	invalid, errs = result(t, User{Email: &email},
		validators.NewChecker("Email", "邮箱").Required().Regex(`^[^@]+@[^@]+$`))
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}

	// 非法邮箱格式
	badEmail := "not-an-email"
	invalid, errs = result(t, User{Email: &badEmail},
		validators.NewChecker("Email", "邮箱").Required().Regex(`^[^@]+@[^@]+$`))
	if !invalid {
		t.Errorf("期望正则校验失败，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// TestValidator_NonStructOrNil 非结构体/空指针/空值不应 panic
func TestValidator_NonStructOrNil(t *testing.T) {
	// 基础类型
	invalid, errs := result(t, "hello",
		validators.NewChecker("Name", "姓名").Required())
	if invalid || len(errs) != 0 {
		t.Errorf("基础类型不应报错，实际：invalid=%v errs=%v", invalid, errs)
	}

	// nil 指针
	var u *User
	invalid, errs = result(t, u,
		validators.NewChecker("Name", "姓名").Required())
	if invalid || len(errs) != 0 {
		t.Errorf("nil 指针不应报错，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 纯 nil
	invalid, errs = result(t, nil,
		validators.NewChecker("Name", "姓名").Required())
	if invalid || len(errs) != 0 {
		t.Errorf("纯 nil 不应报错，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// TestValidator_CheckerFieldNotExist 校验器字段名不存在时应被忽略
func TestValidator_CheckerFieldNotExist(t *testing.T) {
	invalid, errs := result(t, User{Name: "张三"},
		validators.NewChecker("NotExist", "不存在").Required())
	if invalid || len(errs) != 0 {
		t.Errorf("不存在的字段不应报错，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// TestValidator_MultiCheckerSameField 同一字段多个校验器分别执行
func TestValidator_MultiCheckerSameField(t *testing.T) {
	invalid, errs := result(t, User{Name: "ab"},
		validators.NewChecker("Name", "姓名").Min(">3"),
		validators.NewChecker("Name", "姓名").Max("<1"),
	)
	if !invalid || len(errs) != 2 {
		t.Errorf("期望 2 个错误，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// TestChecker_Direct 校验器可脱离 Validator 单独使用
func TestChecker_Direct(t *testing.T) {
	c := validators.NewChecker("Name", "姓名").Min(">2")
	if c.GetField() != "Name" {
		t.Errorf("GetField() = %q，期望 %q", c.GetField(), "Name")
	}

	c = c.check("ab")
	if errs := c.GetErrors(); len(errs) != 1 {
		t.Errorf("期望 1 个错误，实际：%v", errs)
	}

	c = validators.NewChecker("Name", "姓名").Min(">2").check("abc")
	if errs := c.GetErrors(); len(errs) != 0 {
		t.Errorf("期望通过，实际：%v", errs)
	}
}

// TestValidator_SliceMapRules 切片/数组/map 的长度规则
func TestValidator_SliceMapRules(t *testing.T) {
	// 切片 min 失败
	invalid, errs := result(t, struct{ Tags []string }{Tags: []string{"a"}},
		validators.NewChecker("Tags", "标签").Min(">2"))
	if !invalid || !containsErr(errs, "『标签』长度不能小于 2") {
		t.Errorf("期望切片 min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 数组 size 失败
	invalid, errs = result(t, struct{ Nums [3]int }{Nums: [3]int{1, 2}},
		validators.NewChecker("Nums", "数字").Size("==5"))
	if !invalid || !containsErr(errs, "『数字』长度需要等于 5") {
		t.Errorf("期望数组 size 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// map max 失败
	invalid, errs = result(t, struct{ M map[string]int }{M: map[string]int{"a": 1, "b": 2, "c": 3}},
		validators.NewChecker("M", "映射").Max("<3"))
	if !invalid || !containsErr(errs, "『映射』长度不能大于 3") {
		t.Errorf("期望 map max 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 全部通过
	invalid, errs = result(t, struct{ Tags []string }{Tags: []string{"a", "b", "c"}},
		validators.NewChecker("Tags", "标签").Min(">2").Size("==3"))
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}
}

// 自定义基础类型

type (
	AAA    string
	MyInt  int
	MyUint uint
)

// TestValidator_CustomType 自定义基础类型按底层类型校验
func TestValidator_CustomType(t *testing.T) {
	// 自定义 string：min 失败
	invalid, errs := result(t, struct{ Code AAA }{Code: "ab"},
		validators.NewChecker("Code", "编码").Min(">2"))
	if !invalid || !containsErr(errs, "『编码』长度不能小于 2") {
		t.Errorf("期望自定义 string min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 自定义 int：min 失败
	invalid, errs = result(t, struct{ Num MyInt }{Num: 2},
		validators.NewChecker("Num", "数字").Min(">3"))
	if !invalid || !containsErr(errs, "『数字』值不能小于 3") {
		t.Errorf("期望自定义 int min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 自定义 uint：min 失败
	invalid, errs = result(t, struct{ Lv MyUint }{Lv: 0},
		validators.NewChecker("Lv", "等级").Min(">0"))
	if !invalid || len(errs) == 0 {
		t.Errorf("期望自定义 uint min 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 全部通过
	invalid, errs = result(t, struct {
		Code AAA
		Num  MyInt
		Lv   MyUint
	}{Code: "abc", Num: 5, Lv: 1},
		validators.NewChecker("Code", "编码").Min(">2"),
		validators.NewChecker("Num", "数字").Min(">3"),
		validators.NewChecker("Lv", "等级").Min(">0"))
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}
}

// TestValidator_TimeFormat time.Time 字段的 Format 校验按布局格式化后匹配
func TestValidator_TimeFormat(t *testing.T) {
	now := time.Date(2026, 8, 11, 15, 4, 5, 0, time.UTC)

	cases := []struct {
		name   string
		format string
	}{
		{"DateOnly", "DateOnly"},
		{"DateTime", "DateTime"},
		{"RFC3339", "RFC3339"},
		{"RFC3339Nano", "RFC3339Nano"},
		{"TimeOnly", "TimeOnly"},
		{"RFC1123", "RFC1123"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			invalid, errs := result(t, struct{ CreatedAt time.Time }{CreatedAt: now},
				validators.NewChecker("CreatedAt", "创建时间").Format(c.format))
			if invalid {
				t.Errorf("期望 %s 格式通过，实际错误：%v", c.format, errs)
			}
		})
	}
}

// TestValidator_BooleanRules bool 字段的 True/False 规则
func TestValidator_BooleanRules(t *testing.T) {
	// true + True → 通过
	invalid, errs := result(t, struct{ IsActive bool }{IsActive: true},
		validators.NewChecker("IsActive", "启用").True())
	if invalid {
		t.Errorf("true + True 期望通过，实际错误：%v", errs)
	}

	// false + True → 报错
	invalid, errs = result(t, struct{ IsActive bool }{IsActive: false},
		validators.NewChecker("IsActive", "启用").True())
	if !invalid || !containsErr(errs, "『启用』需要是 true") {
		t.Errorf("false + True 期望报错，实际：invalid=%v errs=%v", invalid, errs)
	}

	// false + False → 通过
	invalid, errs = result(t, struct{ IsActive bool }{IsActive: false},
		validators.NewChecker("IsActive", "停用").False())
	if invalid {
		t.Errorf("false + False 期望通过，实际错误：%v", errs)
	}

	// true + False → 报错
	invalid, errs = result(t, struct{ IsActive bool }{IsActive: true},
		validators.NewChecker("IsActive", "停用").False())
	if !invalid || !containsErr(errs, "『停用』需要是 false") {
		t.Errorf("true + False 期望报错，实际：invalid=%v errs=%v", invalid, errs)
	}
}

// 自定义 bool 类型

type MyBool bool

// TestValidator_CustomBool 自定义 bool 类型的 True/False 校验
func TestValidator_CustomBool(t *testing.T) {
	// false + True → 报错
	invalid, errs := result(t, struct{ Flag MyBool }{Flag: MyBool(false)},
		validators.NewChecker("Flag", "标记").True())
	if !invalid || !containsErr(errs, "『标记』需要是 true") {
		t.Errorf("期望自定义 bool True 错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// true + True → 通过
	invalid, errs = result(t, struct{ Flag MyBool }{Flag: MyBool(true)},
		validators.NewChecker("Flag", "标记").True())
	if invalid {
		t.Errorf("期望通过，实际错误：%v", errs)
	}
}

// TestValidator_BooleanCrossType True/False 对其他类型的布尔语义校验
func TestValidator_BooleanCrossType(t *testing.T) {
	// int 1 → true
	invalid, errs := result(t, struct{ N int }{N: 1},
		validators.NewChecker("N", "数量").True())
	if invalid {
		t.Errorf("int 1 + True 期望通过，实际错误：%v", errs)
	}

	// int 0 → false
	invalid, errs = result(t, struct{ N int }{N: 0},
		validators.NewChecker("N", "数量").True())
	if !invalid || !containsErr(errs, "『数量』需要是 true") {
		t.Errorf("int 0 + True 期望报错，实际：invalid=%v errs=%v", invalid, errs)
	}

	// 字符串 "true" + True → 通过
	invalid, errs = result(t, struct{ S string }{S: "true"},
		validators.NewChecker("S", "标记").True())
	if invalid {
		t.Errorf("字符串 true + True 期望通过，实际错误：%v", errs)
	}

	// 字符串 "false" + False → 通过
	invalid, errs = result(t, struct{ S string }{S: "false"},
		validators.NewChecker("S", "标记").False())
	if invalid {
		t.Errorf("字符串 false + False 期望通过，实际错误：%v", errs)
	}
}

// TestValidator_BoolRequired bool 零值（false）时 Required 视为未填
func TestValidator_BoolRequired(t *testing.T) {
	// false + Required → 必填报错
	invalid, errs := result(t, struct{ IsActive bool }{IsActive: false},
		validators.NewChecker("IsActive", "启用").Required())
	if !invalid || !containsErr(errs, "『启用』必填") {
		t.Errorf("false + Required 期望必填错误，实际：invalid=%v errs=%v", invalid, errs)
	}

	// true + Required → 通过
	invalid, errs = result(t, struct{ IsActive bool }{IsActive: true},
		validators.NewChecker("IsActive", "启用").Required())
	if invalid {
		t.Errorf("true + Required 期望通过，实际错误：%v", errs)
	}
}
