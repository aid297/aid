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

	"github.com/aid297/aid/v2/operations"
	"github.com/aid297/aid/v2/points"
)

type (
	Checker interface {
		GetField() string
		check(original any) Checker
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
		ErrMsg(errMsg string) Checker
	}
	CheckerImpl struct {
		field    string
		name     string
		errMsg   *string
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

func (my *CheckerImpl) ErrMsg(errMsg string) Checker { my.errMsg = points.New(errMsg); return my }

func (my *CheckerImpl) GenerateErrMsg(err error) error {
	return operations.NewTernary(
		operations.TrueFn(func() error {
			return fmt.Errorf("『%s』%s", my.name, *my.errMsg)
		}),
		operations.FalseValue(err),
	).GetByValue(my.errMsg != nil)
}

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

func (my *CheckerImpl) check(original any) Checker {
	my.errors = nil // 每次校验前重置错误列表，避免复用时累积
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
			if errs := my.checkMinInt("长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := my.checkMaxInt("长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := my.checkSizeInt("长度", utf8.RuneCountInString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := my.checkIn(value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := my.checkRegex(value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := my.checkFormat(value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := my.checkBoolean(original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// 按整数校验
		value := rv.Int()
		if my.min != nil {
			if errs := my.checkMinInt("值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := my.checkMaxInt("值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := my.checkSizeInt("值", int(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := my.checkIn(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := my.checkRegex(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := my.checkFormat(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := my.checkBoolean(original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// 按无符号整数校验
		value := rv.Uint()
		if my.min != nil {
			if errs := my.checkMinUint("值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := my.checkMaxUint("值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := my.checkSizeUint("值", uint(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := my.checkIn(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := my.checkRegex(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := my.checkFormat(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := my.checkBoolean(original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Float32, reflect.Float64:
		// 按浮点数校验
		value := rv.Float()
		if my.min != nil {
			if errs := my.checkMinFloat("值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := my.checkMaxFloat("值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := my.checkSizeFloat("值", value); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.in != nil {
			if errs := my.checkIn(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.regex != nil {
			if errs := my.checkRegex(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.format != nil {
			if errs := my.checkFormat(cast.ToString(value)); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.boolean != nil {
			if errs := my.checkBoolean(original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Slice, reflect.Array, reflect.Map:
		// 按长度校验
		length := rv.Len()
		if my.min != nil {
			if errs := my.checkMinInt("长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.max != nil {
			if errs := my.checkMaxInt("长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}

		if my.size != nil {
			if errs := my.checkSizeInt("长度", length); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	case reflect.Struct:
		// time.Time 先按布局格式化，再匹配格式规则
		if t, ok := original.(time.Time); ok {
			if my.format != nil {
				if layout, ok := timeLayouts[*my.format]; ok {
					if errs := my.checkFormat(t.Format(layout)); errs != nil {
						my.errors = append(my.errors, errs...)
					}
				}
			}
		}
	case reflect.Bool:
		if my.boolean != nil {
			if errs := my.checkBoolean(original); errs != nil {
				my.errors = append(my.errors, errs...)
			}
		}
	}

	return my
}

func (my *CheckerImpl) checkMinInt(mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.min, ">="); ok {
		target := cast.ToInt(ruleVal)
		if !(src >= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于等于 %d", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.min, ">"); ok {
		target := cast.ToInt(ruleVal)
		if !(src > target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于 %d", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkMinUint(mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.min, ">="); ok {
		target := cast.ToUint(ruleVal)
		if !(src >= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于等于 %d", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.min, ">"); ok {
		target := cast.ToUint(ruleVal)
		if !(src > target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于 %d", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkMinFloat(mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.min, ">="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src >= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于等于 %f", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.min, ">"); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src > target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能小于 %f", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkMaxInt(mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.max, "<="); ok {
		target := cast.ToInt(ruleVal)
		if !(src <= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于等于 %d", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.max, "<"); ok {
		target := cast.ToInt(ruleVal)
		if !(src < target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于 %d", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkMaxUint(mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.max, "<="); ok {
		target := cast.ToUint(ruleVal)
		if !(src <= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于等于 %d", my.name, mid, target)))
		}
	} else if ruleVal, ok := strings.CutPrefix(*my.max, "<"); ok {
		target := cast.ToUint(ruleVal)
		if !(src < target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于 %d", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkMaxFloat(mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.max, "<="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src <= target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于等于 %f", my.name, mid, target)))
		}
	} else if ruleVal, ok := strings.CutPrefix(*my.max, "<"); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src < target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能大于 %f", my.name, mid, target)))
		}
	}

	return
}

func (my *CheckerImpl) checkSizeInt(mid string, src int) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.size, "=="); ok {
		target := cast.ToInt(ruleVal)
		if !(src == target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s需要等于 %d", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.size, "!="); ok {
		target := cast.ToInt(ruleVal)
		if !(src != target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能等于 %d", my.name, mid, target)))
		}
	}

	if rules := strings.Split(*my.size, "~"); len(rules) == 2 {
		my.min = points.New(fmt.Sprintf(">%s", rules[0]))
		if errs := my.checkMinInt(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		my.max = points.New(fmt.Sprintf("<%s", rules[1]))
		if errs := my.checkMaxInt(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func (my *CheckerImpl) checkSizeUint(mid string, src uint) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.size, "=="); ok {
		target := cast.ToUint(ruleVal)
		if !(src == target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s需要等于 %d", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.size, "!="); ok {
		target := cast.ToUint(ruleVal)
		if !(src != target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能等于 %d", my.name, mid, target)))
		}
	}

	if rules := strings.Split(*my.size, "~"); len(rules) == 2 {
		my.min = points.New(fmt.Sprintf(">=%s", rules[0]))
		if errs := my.checkMinUint(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		my.max = points.New(fmt.Sprintf("<=%s", rules[1]))
		if errs := my.checkMaxUint(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func (my *CheckerImpl) checkSizeFloat(mid string, src float64) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.size, "=="); ok {
		target := cast.ToFloat64(ruleVal)
		if !(src == target) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s需要等于 %f", my.name, mid, target)))
		}
		return
	}

	if ruleVal, ok := strings.CutPrefix(*my.size, "!="); ok {
		target := cast.ToFloat64(ruleVal)
		if src == target {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』%s不能等于 %f", my.name, mid, target)))
		}
		return
	}

	if rules := strings.Split(*my.size, "~"); len(rules) == 2 {
		my.min = points.New(fmt.Sprintf(">=%s", rules[0]))
		if errs := my.checkMinFloat(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
		my.max = points.New(fmt.Sprintf("<=%s", rules[1]))
		if errs := my.checkMaxFloat(mid, src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func (my *CheckerImpl) checkIn(src string) (errors []error) {
	if ruleVal, ok := strings.CutPrefix(*my.in, "=="); ok {
		target := strings.Split(ruleVal, defaultSliceSplitChar)
		if !utils.Contains(target, src) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』内容必须在 %s 中", my.name, strings.Join(target, ","))))
		}
	}

	if ruleVal, ok := strings.CutPrefix(*my.in, "!="); ok {
		target := strings.Split(ruleVal, defaultSliceSplitChar)
		if utils.Contains(target, src) {
			errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』内容不能在 %s 中", my.name, strings.Join(target, ","))))
		}
	}

	return
}

func (my *CheckerImpl) checkRegex(src string) (errors []error) {
	if !regexp.MustCompile(*my.regex).MatchString(src) {
		errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』内容必须符合 %s 规则", my.name, *my.regex)))
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

func (my *CheckerImpl) checkFormat(src string) (errors []error) {
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
	if f, ok := patternsForTimeString[*my.format]; ok {
		my.regex = &f
		if errs := my.checkRegex(src); errs != nil {
			errors = append(errors, errs...)
		}
	}

	return
}

func (my *CheckerImpl) checkBoolean(src any) (errors []error) {
	if *my.boolean != cast.ToBool(src) {
		errors = append(errors, my.GenerateErrMsg(fmt.Errorf("『%s』需要是 %t", my.name, *my.boolean)))
	}

	return
}
