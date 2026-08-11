package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cast"
	"gorm.io/gorm/utils"

	"github.com/aid297/aid/v2/points"
)

type (
	Checker interface {
		GetField() string
		Check(original any) Checker
		GetErrors() []error
		Size(size string) Checker
		Min(min string) Checker
		Max(max string) Checker
		In(in string) Checker
		Format(format string) Checker
		Regex(regex string) Checker
		Required() Checker
		True() Checker
		False() Checker
	}
	CheckerImpl struct {
		field    string
		name     string
		errors   []error
		rule     string
		required bool
		min      *string
		max      *string
		size     *string
		in       *string
		regex    *string
		format   *string
		boolean  *bool
	}
)

func NewChecker(field, name string) Checker {
	return &CheckerImpl{field: field, name: name}
}

func (my *CheckerImpl) GetErrors() []error { return my.errors }

func (my *CheckerImpl) GetField() string { return my.field }

func (my *CheckerImpl) Size(size string) Checker { my.size = &size; return my }

func (my *CheckerImpl) Min(min string) Checker { my.min = &min; return my }

func (my *CheckerImpl) Max(max string) Checker { my.max = &max; return my }

func (my *CheckerImpl) In(in string) Checker { my.in = &in; return my }

func (my *CheckerImpl) Format(format string) Checker { my.format = &format; return my }

func (my *CheckerImpl) Regex(regex string) Checker { my.regex = &regex; return my }

func (my *CheckerImpl) Required() Checker { my.required = true; return my }

func (my *CheckerImpl) True() Checker { my.boolean = points.New(true); return my }

func (my *CheckerImpl) False() Checker { my.boolean = points.New(false); return my }

func (my *CheckerImpl) isZeroValue(original any) bool {
	v := reflect.ValueOf(any(original)) // 显式转 any，避免对 T 直接反射的坑
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	// 接口类型且 nil
	if !v.IsValid() {
		return true
	}
	return v.IsZero()
}

// unwrapValue 若入参是 reflect.Value（例如 Validate 反射遍历结构体字段时传入），
// 取出其底层实际值；指针类型自动解引用，nil 或不可取接口值的返回 nil
func unwrapValue(original any) any {
	rv, ok := original.(reflect.Value)
	if !ok {
		return original
	}

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	if !rv.IsValid() || !rv.CanInterface() {
		return nil
	}
	return rv.Interface()
}

func (my *CheckerImpl) Check(original any) Checker {
	original = unwrapValue(original)

	if my.required && my.isZeroValue(original) {
		my.errors = append(my.errors, fmt.Errorf("『%s』必填", my.name))
		return my
	}

	rv := reflect.ValueOf(original)
	if !rv.IsValid() {
		return my
	}

	// 按 reflect.Kind 分发，兼容自定义基础类型（如 type AAA string）
	switch rv.Kind() {
	case reflect.String:
		// 按 string 校验：min/max/size/regex/in/format
		value := rv.String()
		if my.min != nil {
			if errs := checkMinInt(my.name, *my.min, "长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := checkMaxInt(my.name, *my.max, "长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := checkSizeInt(my.name, *my.size, "长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := checkIn(my.name, *my.in, value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := checkRegex(my.name, *my.regex, value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := checkFormat(my.name, *my.format, value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := checkBoolean(my.name, *my.boolean, original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 按整数校验
		value := rv.Int()
		if my.min != nil {
			if errs := checkMinInt(my.name, *my.min, "值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := checkMaxInt(my.name, *my.max, "值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := checkSizeInt(my.name, *my.size, "值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := checkIn(my.name, *my.in, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := checkRegex(my.name, *my.regex, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := checkFormat(my.name, *my.format, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := checkBoolean(my.name, *my.boolean, original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 按无符号整数校验
		value := rv.Uint()
		if my.min != nil {
			if errs := checkMinUint(my.name, *my.min, "值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := checkMaxUint(my.name, *my.max, "值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := checkSizeUint(my.name, *my.size, "值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := checkIn(my.name, *my.in, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := checkRegex(my.name, *my.regex, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := checkFormat(my.name, *my.format, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := checkBoolean(my.name, *my.boolean, original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Float32, reflect.Float64:
		// 按浮点数校验
		value := rv.Float()
		if my.min != nil {
			if errs := checkMinFloat(my.name, *my.min, "值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := checkMaxFloat(my.name, *my.max, "值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := checkSizeFloat(my.name, *my.size, "值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := checkIn(my.name, *my.in, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := checkRegex(my.name, *my.regex, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := checkFormat(my.name, *my.format, cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := checkBoolean(my.name, *my.boolean, original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		// 按长度校验
		length := rv.Len()
		if my.min != nil {
			if errs := checkMinInt(my.name, *my.min, "长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := checkMaxInt(my.name, *my.max, "长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := checkSizeInt(my.name, *my.size, "长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Struct:
		// time.Time 先按布局格式化，再匹配格式规则
		if t, ok := original.(time.Time); ok {
			if my.format != nil {
				if layout, ok := timeLayouts[*my.format]; ok {
					if errs := checkFormat(my.name, *my.format, t.Format(layout)); errs != nil {
						my.errors = append(my.errors, errs...)
					}
				}
			}
		}
	case reflect.Bool:
		if my.boolean != nil {
			if errs := checkBoolean(my.name, *my.boolean, original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	}

	return my
}

func checkMinInt(name, rule, mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, ">"); ok {
		target := cast.ToInt(ruleVal)
		if !(src > target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于 %d", name, mid, target))
		}

		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, ">="); ok {
		target := cast.ToInt(ruleVal)
		if !(src >= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于等于 %d", name, mid, target))
		}
	}

	return
}

func checkMinUint(name, rule, mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, ">"); ok {
		target := cast.ToUint(ruleVal)
		if !(src > target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于 %d", name, mid, target))
		}

		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, ">="); ok {
		target := cast.ToUint(ruleVal)
		if !(src >= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于等于 %d", name, mid, target))
		}
	}

	return
}

func checkMinFloat(name, rule, mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, ">"); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src > target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于 %f", name, mid, target))
		}

		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, ">="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src >= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于等于 %f", name, mid, target))
		}
	}

	return
}

func checkMaxInt(name, rule, mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "<"); ok {
		target := cast.ToInt(ruleVal)
		if !(src < target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于 %d", name, mid, target))
		}
	}

	if ruleVal, ok := strings.CutPrefix(rule, "<="); ok {
		target := cast.ToInt(ruleVal)
		if !(src <= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于等于 %d", name, mid, target))
		}
	}

	return
}

func checkMaxUint(name, rule, mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "<"); ok {
		target := cast.ToUint(ruleVal)
		if !(src < target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于 %d", name, mid, target))
		}
	}

	if ruleVal, ok := strings.CutPrefix(rule, "<="); ok {
		target := cast.ToUint(ruleVal)
		if !(src <= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于等于 %d", name, mid, target))
		}
	}

	return
}

func checkMaxFloat(name, rule, mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "<"); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src < target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能大于 %f", name, mid, target))
		}
	}

	if ruleVal, ok := strings.CutPrefix(rule, "<="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src <= target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能小于等于 %f", name, mid, target))
		}
	}

	return
}

func checkSizeInt(name, rule, mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "=="); ok {
		target := cast.ToInt(ruleVal)
		if !(src == target) {
			errors = append(errors, fmt.Errorf("『%s』%s需要等于 %d", name, mid, target))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, "!="); ok {
		target := cast.ToInt(ruleVal)
		if !(src != target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能等于 %d", name, mid, target))
		}
	}

	if rules := strings.Split(rule, "~"); len(rules) == 2 {
		if errs := checkMinInt(name, fmt.Sprintf(">%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxInt(name, fmt.Sprintf("<%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	if rules := strings.Split(rule, "=~"); len(rules) == 2 {
		if errs := checkMinInt(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxInt(name, fmt.Sprintf("<%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinInt(name, fmt.Sprintf(">%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxInt(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinInt(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxInt(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func checkSizeUint(name, rule, mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "=="); ok {
		target := cast.ToUint(ruleVal)
		if !(src == target) {
			errors = append(errors, fmt.Errorf("『%s』%s需要等于 %d", name, mid, target))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, "!="); ok {
		target := cast.ToUint(ruleVal)
		if !(src != target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能等于 %d", name, mid, target))
		}
	}

	if rules := strings.Split(rule, "~"); len(rules) == 2 {
		if errs := checkMinUint(name, fmt.Sprintf(">%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxUint(name, fmt.Sprintf("<%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	if rules := strings.Split(rule, "=~"); len(rules) == 2 {
		if errs := checkMinUint(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxUint(name, fmt.Sprintf("<%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinUint(name, fmt.Sprintf(">%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxUint(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinUint(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxUint(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func checkSizeFloat(name, rule, mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "=="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src == target) {
			errors = append(errors, fmt.Errorf("『%s』%s需要等于 %f", name, mid, target))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(rule, "!="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src != target) {
			errors = append(errors, fmt.Errorf("『%s』%s不能等于 %f", name, mid, target))
		}

		if errs := checkMinFloat(name, fmt.Sprintf(">%s", rule), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxFloat(name, fmt.Sprintf("<%s", rule), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	if rules := strings.Split(rule, "=~"); len(rules) == 2 {
		if errs := checkMinFloat(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxFloat(name, fmt.Sprintf("<%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinFloat(name, fmt.Sprintf(">%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxFloat(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}

		if errs := checkMinFloat(name, fmt.Sprintf(">=%s", rules[0]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		if errs := checkMaxFloat(name, fmt.Sprintf("<=%s", rules[1]), mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func checkIn(name, rule string, src string) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(rule, "=="); ok {
		target := strings.Split(ruleVal, "|")
		if !utils.Contains(target, src) {
			errors = append(errors, fmt.Errorf("『%s』内容必须在 %s 中", name, strings.Join(target, ",")))
		}
	}

	if ruleVal, ok := strings.CutPrefix(rule, "!="); ok {
		target := strings.Split(ruleVal, "|")
		if utils.Contains(target, src) {
			errors = append(errors, fmt.Errorf("『%s』内容不能在 %s 中", name, strings.Join(target, ",")))
		}
	}

	return
}

func checkRegex(name, rule, src string) (errors []error) {
	if !regexp.MustCompile(rule).MatchString(src) {
		errors = append(errors, fmt.Errorf("『%s』内容必须符合 %s 规则", name, rule))
	}
	return
}

// timeLayouts 格式名称到 Go 时间布局的映射，用于 time.Time 字段的 Format 校验
var timeLayouts = map[string]string{
	"RFC3339":         time.RFC3339,
	"RFC3339Nano":     "2006-01-02T15:04:05.000000000Z07:00",
	"DateTime":        "2006-01-02 15:04:05",
	"DateOnly":        time.DateOnly,
	"TimeOnly":        time.TimeOnly,
	"ReferenceLayout": "01/02 03:04:05PM '06 -0700",
	"ANSIC":           time.ANSIC,
	"UnixDate":        time.UnixDate,
	"RubyDate":        time.RubyDate,
	"RFC822":          time.RFC822,
	"RFC822Z":         time.RFC822Z,
	"RFC850":          time.RFC850,
	"RFC1123":         time.RFC1123,
	"RFC1123Z":        time.RFC1123Z,
	"Kitchen":         time.Kitchen,
	"Stamp":           time.Stamp,
	"StampMilli":      time.StampMilli,
	"StampMicro":      time.StampMicro,
	"StampNano":       time.StampNano,
}

func checkFormat(name, rule, src string) (errors []error) {
	var patternsForTimeString = map[string]string{
		"RFC3339":           `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+\-]\d{2}:\d{2})$`,
		"RFC3339Nano":       `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+(Z|[+\-]\d{2}:\d{2})$`,
		"DateTime":          `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`,
		"DateOnly":          `^\d{4}-\d{2}-\d{2}$`,
		"TimeOnly":          `^\d{2}:\d{2}:\d{2}$`,
		"ReferenceLayout":   `^\d{2}/\d{2} \d{2}:\d{2}:\d{2}(AM|PM) '\d{2} [+\-]\d{4}$`,
		"ANSIC":             `^[A-Za-z]{3} [A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2} \d{4}$`,
		"UnixDate":          `^[A-Za-z]{3} [A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2} [A-Za-z]{3,4} \d{4}$`,
		"RubyDate":          `^[A-Za-z]{3} [A-Za-z]{3} \d{2} \d{2}:\d{2}:\d{2} [+\-]\d{4} \d{4}$`,
		"RFC822":            `^\d{2} [A-Za-z]{3} \d{2} \d{2}:\d{2} [A-Za-z]{3}$`,
		"RFC822Z":           `^\d{2} [A-Za-z]{3} \d{2} \d{2}:\d{2} [+\-]\d{4}$`,
		"RFC850":            `^[A-Za-z]+, \d{2}-[A-Za-z]{3}-\d{2} \d{2}:\d{2}:\d{2} [A-Za-z]{3}$`,
		"RFC1123":           `^[A-Za-z]{3}, \d{2} [A-Za-z]{3} \d{4} \d{2}:\d{2}:\d{2} [A-Za-z]{3}$`,
		"RFC1123Z":          `^[A-Za-z]{3}, \d{2} [A-Za-z]{3} \d{4} \d{2}:\d{2}:\d{2} [+\-]\d{4}$`,
		"Kitchen":           `^\d{1,2}:\d{2}(AM|PM)$`,
		"Stamp":             `^[A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2}$`,
		"StampMilli":        `^[A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2}\.\d{3}$`,
		"StampMicro":        `^[A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2}\.\d{6}$`,
		"StampNano":         `^[A-Za-z]{3} [ \d]\d \d{2}:\d{2}:\d{2}\.\d{9}$`,
		"SonarQubeDatetime": `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})$`,
	}
	if f, ok := patternsForTimeString[rule]; ok {
		if errs := checkRegex(name, f, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func checkBoolean(name string, rule bool, src any) (errors []error) {
	if rule != cast.ToBool(src) {
		errors = append(errors, fmt.Errorf("『%s』需要是 %t", name, rule))
	}

	return
}
