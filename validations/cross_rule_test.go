package validations_test

import (
	"strings"
	"testing"
	"time"

	"github.com/aid297/aid/v2/validations"
)

// ---------- 辅助结构体 ----------

type timeRange struct {
	StartTime time.Time `v-rule:"required" v-name:"start_time"`
	EndTime   time.Time `v-rule:"required" v-name:"end_time"`
}

func (f timeRange) Cross() string { return "({start_time}<{end_time})" }

type intRange struct {
	Min int `v-rule:"required" v-name:"min"`
	Max int `v-rule:"required" v-name:"max"`
}

func (f intRange) Cross() string { return "({min}<{max})" }

type floatRange struct {
	Low  float64 `v-rule:"required" v-name:"low"`
	High float64 `v-rule:"required" v-name:"high"`
}

func (f floatRange) Cross() string { return "({low}<{high})" }

type stringRange struct {
	First string `v-rule:"required" v-name:"first"`
	Last  string `v-rule:"required" v-name:"last"`
}

func (f stringRange) Cross() string { return "({first}<{last})" }

// ---------- 测试用例 ----------

// TestCrossRule_TimeLessThan 时间类型跨字段校验：start_time < end_time
func TestCrossRule_TimeLessThan(t *testing.T) {
	// 正常：start < end
	ok := validations.Once().Checker(&timeRange{
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	}).Validate()
	if ok.Invalid() {
		t.Errorf("期望校验通过，实际失败：%s", ok.ErrorToString("\n"))
	}

	// 异常：start > end
	ng := validations.Once().Checker(&timeRange{
		StartTime: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}).Validate()
	if ng.OK() {
		t.Errorf("期望校验失败（start > end），实际通过")
	}
	if !strings.Contains(ng.ErrorToString(""), "跨字段校验失败") {
		t.Errorf("期望包含跨字段校验错误，实际：%s", ng.ErrorToString(""))
	}
}

// TestCrossRule_TimeLessEqual 时间类型：start_time <= end_time（相等场景）
func TestCrossRule_TimeLessEqual(t *testing.T) {
	same := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)
	ok := validations.Once().Checker(&leCross{StartTime: same, EndTime: same}).Validate()
	if ok.Invalid() {
		t.Errorf("相等时 <= 应通过，实际：%s", ok.ErrorToString("\n"))
	}
}

type leCross struct {
	StartTime time.Time `v-rule:"required" v-name:"start_time"`
	EndTime   time.Time `v-rule:"required" v-name:"end_time"`
}

func (f leCross) Cross() string { return "({start_time}<={end_time})" }

// TestCrossRule_IntComparison 整数类型跨字段比较
func TestCrossRule_IntComparison(t *testing.T) {
	ok := validations.Once().Checker(&intRange{Min: 1, Max: 100}).Validate()
	if ok.Invalid() {
		t.Errorf("期望通过，实际：%s", ok.ErrorToString("\n"))
	}

	ng := validations.Once().Checker(&intRange{Min: 100, Max: 1}).Validate()
	if ng.OK() {
		t.Errorf("期望失败（min > max）")
	}
}

// TestCrossRule_FloatComparison 浮点数类型
func TestCrossRule_FloatComparison(t *testing.T) {
	ok := validations.Once().Checker(&floatRange{Low: 1.5, High: 9.9}).Validate()
	if ok.Invalid() {
		t.Errorf("期望通过，实际：%s", ok.ErrorToString("\n"))
	}

	ng := validations.Once().Checker(&floatRange{Low: 9.9, High: 1.5}).Validate()
	if ng.OK() {
		t.Errorf("期望失败")
	}
}

// TestCrossRule_StringComparison 字符串类型字典序比较
func TestCrossRule_StringComparison(t *testing.T) {
	ok := validations.Once().Checker(&stringRange{First: "apple", Last: "banana"}).Validate()
	if ok.Invalid() {
		t.Errorf("期望通过，实际：%s", ok.ErrorToString("\n"))
	}

	ng := validations.Once().Checker(&stringRange{First: "banana", Last: "apple"}).Validate()
	if ng.OK() {
		t.Errorf("期望失败")
	}
}

// TestCrossRule_SkipWhenRefFieldFailed 当引用字段独立校验失败时，跳过跨字段校验
func TestCrossRule_SkipWhenRefFieldFailed(t *testing.T) {
	ng := validations.Once().Checker(&timeRange{
		StartTime: time.Time{}, // 零值，required 失败
		EndTime:   time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	}).Validate()

	errs := ng.ErrorToString("\n")
	if !strings.Contains(errs, "不能为空") {
		t.Errorf("期望有 start_time 的必填错误，实际：%s", errs)
	}
	if strings.Contains(errs, "跨字段校验失败") {
		t.Errorf("引用字段失败时应跳过跨字段校验，但出现了跨字段错误：%s", errs)
	}
}

// TestCrossRule_SkipWhenCurrentFieldFailed 当前字段独立校验失败时跳过跨字段校验
func TestCrossRule_SkipWhenCurrentFieldFailed(t *testing.T) {
	ng := validations.Once().Checker(&timeRange{
		StartTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Time{}, // 零值，required 失败
	}).Validate()

	errs := ng.ErrorToString("\n")
	if strings.Contains(errs, "跨字段校验失败") {
		t.Errorf("当前字段失败时应跳过跨字段校验，但出现了跨字段错误：%s", errs)
	}
}

// TestCrossRule_EqualAndNotEqual == 和 != 运算符
func TestCrossRule_EqualAndNotEqual(t *testing.T) {
	// == 通过
	ok := validations.Once().Checker(&eqCross{A: 42, B: 42}).Validate()
	if ok.Invalid() {
		t.Errorf("== 相等时应通过")
	}
	// == 失败
	ng := validations.Once().Checker(&eqCross{A: 42, B: 43}).Validate()
	if ng.OK() {
		t.Errorf("== 不等时应失败")
	}

	// != 通过
	ok2 := validations.Once().Checker(&neqCross{A: 1, B: 2}).Validate()
	if ok2.Invalid() {
		t.Errorf("!= 不等时应通过")
	}
	// != 失败
	ng2 := validations.Once().Checker(&neqCross{A: 1, B: 1}).Validate()
	if ng2.OK() {
		t.Errorf("!= 相等时应失败")
	}
}

type eqCross struct {
	A int `v-rule:"required" v-name:"a"`
	B int `v-rule:"required" v-name:"b"`
}

func (f eqCross) Cross() string { return "({a}=={b})" }

type neqCross struct {
	A int `v-rule:"required" v-name:"a"`
	B int `v-rule:"required" v-name:"b"`
}

func (f neqCross) Cross() string { return "({a}!={b})" }

// TestCrossRule_MultipleRules 多条跨字段规则
func TestCrossRule_MultipleRules(t *testing.T) {
	// 手动实现 CrossValidator
	ok := validations.Once().Checker(&multiCross{Min: 1, Max: 100, Value: 50}).Validate()
	if ok.Invalid() {
		t.Errorf("期望通过，实际：%s", ok.ErrorToString("\n"))
	}

	// value < min → 失败
	ng := validations.Once().Checker(&multiCross{Min: 50, Max: 100, Value: 10}).Validate()
	if ng.OK() {
		t.Errorf("value < min 时应失败")
	}

	// value > max → 失败
	ng2 := validations.Once().Checker(&multiCross{Min: 1, Max: 50, Value: 100}).Validate()
	if ng2.OK() {
		t.Errorf("value > max 时应失败")
	}
}

type multiCross struct {
	Min   int `v-rule:"required" v-name:"min"`
	Max   int `v-rule:"required" v-name:"max"`
	Value int `v-rule:"required" v-name:"value"`
}

func (f multiCross) Cross() string { return "({min}<={value})({max}>={value})" }

// TestCrossRule_NonExistentRefField 引用不存在的字段时静默跳过
func TestCrossRule_NonExistentRefField(t *testing.T) {
	ok := validations.Once().Checker(&nonExistCross{Value: 42}).Validate()
	if ok.Invalid() {
		t.Errorf("引用不存在的字段应静默跳过，实际：%s", ok.ErrorToString("\n"))
	}
}

type nonExistCross struct {
	Value int `v-rule:"required" v-name:"value"`
}

func (f nonExistCross) Cross() string { return "({nonexistent}<{value})" }

// TestCrossRule_NilPointerSkip 非必填 nil 指针字段跳过跨字段校验
func TestCrossRule_NilPointerSkip(t *testing.T) {
	ok := validations.Once().Checker(&nilCross{StartTime: nil, EndTime: nil}).Validate()
	if ok.Invalid() {
		t.Errorf("nil 指针应跳过校验，实际：%s", ok.ErrorToString("\n"))
	}
}

type nilCross struct {
	StartTime *time.Time `v-name:"start_time"`
	EndTime   *time.Time `v-name:"end_time"`
}

func (f nilCross) Cross() string { return "({start_time}<{end_time})" }

// TestCrossRule_NoCrossMethod 未实现 CrossValidator 接口的结构体正常校验
func TestCrossRule_NoCrossMethod(t *testing.T) {
	type plain struct {
		A int `v-rule:"required" v-name:"a"`
		B int `v-rule:"required" v-name:"b"`
	}

	ok := validations.Once().Checker(&plain{A: 100, B: 1}).Validate()
	if ok.Invalid() {
		t.Errorf("未实现 Cross() 应只做独立校验，实际：%s", ok.ErrorToString("\n"))
	}
}

// TestCrossRule_ParseInvalid 无效 Cross() 返回值静默忽略
func TestCrossRule_ParseInvalid(t *testing.T) {
	ok := validations.Once().Checker(&invalidCross{A: 1, B: 2}).Validate()
	if ok.Invalid() {
		t.Errorf("无效 Cross() 应静默忽略，实际：%s", ok.ErrorToString("\n"))
	}
}

type invalidCross struct {
	A int `v-rule:"required" v-name:"a"`
	B int `v-rule:"required" v-name:"b"`
}

func (f invalidCross) Cross() string { return "invalid-rule-no-braces" }
