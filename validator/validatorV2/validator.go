package validatorV2

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type (
	Validator struct{}

	// FieldResult 保留类型以兼容引用（目前不再返回 FieldResult）
	FieldResult struct {
		Field  string   `json:"field"`
		Errors []string `json:"errors"`
	}
)

// ExFn 扩展校验函数类型及注册表
type ExFn func(val any) error

var (
	exFns   = make(map[string]ExFn)
	exFnsMu sync.RWMutex
)

// RegisterExFun 注册一个扩展校验函数，key 为 v-ex 标签中的键
func RegisterExFun(key string, fn ExFn) {
	exFnsMu.Lock()
	defer exFnsMu.Unlock()
	exFns[key] = fn
}

// UnregisterExFun 注销注册函数
func UnregisterExFun(key string) {
	exFnsMu.Lock()
	defer exFnsMu.Unlock()
	delete(exFns, key)
}

func getExFun(key string) (ExFn, bool) {
	exFnsMu.RLock()
	defer exFnsMu.RUnlock()
	f, ok := exFns[key]
	return f, ok
}

// splitExKeys 支持用逗号或分号分隔多个扩展 key
func splitExKeys(s string) []string {
	s = strings.TrimSpace(s)
	// 支持 ; 或 , 分隔
	if s == "" {
		return nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ';' || r == ','
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// removeRule 从规则字符串中移除指定的规则 key（按分号分隔）
func removeRule(ruleStr, key string) string {
	if strings.TrimSpace(ruleStr) == "" {
		return ruleStr
	}
	parts := strings.Split(ruleStr, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// 取出规则名（遇到 : = > < 或 = 号时取左侧）
		name := p
		if idx := strings.IndexAny(p, ":=<>\"' "); idx >= 0 {
			name = strings.TrimSpace(p[:idx])
		}
		if name == key {
			continue
		}
		out = append(out, p)
	}
	return strings.Join(out, ";")
}

// WithFiber 在 Fiber 框架中绑定并验证请求数据。
func WithFiber[T any](c *fiber.Ctx, fns ...func(ins any) (err error)) (T, error) {
	var (
		err error
		ins = new(T)
	)

	if err = c.BodyParser(ins); err != nil {
		return *ins, err
	}

	errs := app.Validator.Validate(ins, fns...)
	if len(errs) > 0 {
		parts := make([]string, 0, len(errs))
		for _, e := range errs {
			parts = append(parts, e.Error())
		}
		return *ins, errors.New(strings.Join(parts, "; "))
	}

	return *ins, nil
}

// WithGin 在 Gin 框架中绑定并验证请求数据。
func WithGin[T any](c *gin.Context, fns ...func(ins any) (err error)) (T, error) {
	var (
		err error
		ins = new(T)
	)

	if err = c.ShouldBind(ins); err != nil {
		return *ins, err
	}

	errs := app.Validator.Validate(ins, fns...)
	if len(errs) > 0 {
		parts := make([]string, 0, len(errs))
		for _, e := range errs {
			parts = append(parts, e.Error())
		}
		return *ins, errors.New(strings.Join(parts, "; "))
	}

	return *ins, nil
}

// Validate 验证任意结构体，返回每个字段的验证结果（字段名和错误切片）。
// 支持的 tag:
// - v-rule: 规则串，规则之间以分号分隔。例如: "required;min>3;max<10;email;in=a,b"
// - v-name: 字段可读名称，嵌套时会以点号拼接，例如: 父.子
// 设计假设：
// - min / max 对字符串表示长度（>= / <=），对数字表示数值（>= / <=）。
// - in / not-in 使用逗号分隔的值列表。
// - regex=... 使用完整正则表达式。
// - 支持快捷规则: email, date(2006-01-02), time(15:04:05), datetime(2006-01-02 15:04:05)
func (Validator) Validate(data any, exFns ...func(d any) error) []error {
	if data == nil {
		return []error{fmt.Errorf("data is nil")}
	}
	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []error{fmt.Errorf("不支持空指针验证")}
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return []error{fmt.Errorf("只支持指针或结构体")}
	}

	errs := make([]error, 0)
	walkStruct(rv, "", &errs)
	if len(errs) > 0 {
		return errs
	}

	for idx := range exFns {
		if exFns[idx] != nil {
			if err := exFns[idx](&data); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}

// walkStruct 递归遍历结构体字段，处理嵌套结构体（排除 time.Time）
func walkStruct(rv reflect.Value, parent string, results *[]error) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		// 跳过未导出字段
		if sf.PkgPath != "" {
			continue
		}
		fv := rv.Field(i)
		// 字段显示名
		vName := sf.Tag.Get("v-name")
		if vName == "" {
			vName = sf.Name
		}
		fullName := vName
		if parent != "" {
			fullName = parent + "." + vName
		}

		ruleStr := sf.Tag.Get("v-rule")
		// 扩展校验标签，可写多个 key，使用逗号或分号分隔
		exTag := sf.Tag.Get("v-ex")

		// 处理指针类型：保留nil信息以便 required 校验
		origVal := fv.Interface()
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				// nil 指针 -> 仅能检测 required
				errs := applyRulesToNil(ruleStr, fullName)
				for _, e := range errs {
					*results = append(*results, fmt.Errorf(e))
				}
				// 执行扩展函数（nil 值）
				if exTag != "" {
					exKeys := splitExKeys(exTag)
					for _, k := range exKeys {
						if fn, ok := getExFun(k); ok {
							if err := fn(nil); err != nil {
								*results = append(*results, err)
							}
						}
					}
				}
				continue
			}
			fv = fv.Elem()
		}

		// 当字段原始为非 nil 指针且包含 required 规则时，视为满足 required（只检查存在性），
		// 因此在后续对解引用值的校验中不再把 required 作为一项规则来检查。
		effectiveRuleStr := ruleStr
		if origVal != nil {
			if rvOrig := reflect.ValueOf(origVal); rvOrig.Kind() == reflect.Ptr {
				effectiveRuleStr = removeRule(effectiveRuleStr, "required")
			}
		}

		// 特殊处理 time.Time
		if fv.Kind() == reflect.Struct && fv.Type().PkgPath() == "time" && fv.Type().Name() == "Time" {
			errs := checkTimeValue(fv.Interface(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
			continue
		}

		switch fv.Kind() {
		case reflect.Struct:
			// 嵌套结构体 -> 递归
			walkStruct(fv, fullName, results)
		case reflect.String:
			errs := checkString(fv.String(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		case reflect.Slice, reflect.Array:
			errs := checkSlice(fv, effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
			// 如果元素是结构体或指针结构体，继续验证元素内部字段，并使用索引扩展名称
			elemKind := fv.Type().Elem().Kind()
			if elemKind == reflect.Struct || (elemKind == reflect.Ptr && fv.Type().Elem().Elem().Kind() == reflect.Struct) {
				for j := 0; j < fv.Len(); j++ {
					e := fv.Index(j)
					for e.Kind() == reflect.Ptr {
						if e.IsNil() {
							// 元素为 nil，跳过深度验证
							e = reflect.Zero(e.Type().Elem())
							break
						}
						e = e.Elem()
					}
					walkStruct(e, fmt.Sprintf("%s[%d]", fullName, j), results)
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			errs := checkNumberInt(fv.Int(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			errs := checkNumberUint(fv.Uint(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		case reflect.Float32, reflect.Float64:
			errs := checkNumberFloat(fv.Float(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		case reflect.Bool:
			// 目前不对 bool 做额外校验，除非有 in/not-in 等规则
			errs := checkBool(fv.Bool(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		default:
			// 其他类型目前不处理，但如果存在 required 规则需要提示
			errs := applyRulesToAny(fv.Interface(), effectiveRuleStr, fullName)
			for _, e := range errs {
				*results = append(*results, fmt.Errorf(e))
			}
		}

		// 执行扩展校验函数（v-ex）: 传入原始值（指针未展开的 origVal if was pointer else fv.Interface()）
		if exTag != "" {
			exKeys := splitExKeys(exTag)
			var valForEx any
			// 尝试使用原始未展开值（保持 pointer 状态），否则使用当前 fv.Interface()
			if origVal != nil {
				valForEx = origVal
			} else {
				valForEx = fv.Interface()
			}
			for _, k := range exKeys {
				if fn, ok := getExFun(k); ok {
					if err := fn(valForEx); err != nil {
						*results = append(*results, fmt.Errorf("[%s] %s", fullName, err.Error()))
					}
				}
			}
		}
	}
}

// 解析简单规则串
type rule struct {
	key string
	op  string // =, >, <, : or empty
	val string
}

func parseRules(s string) []rule {
	out := make([]rule, 0)
	s = strings.TrimSpace(s)
	if s == "" {
		return out
	}
	parts := strings.Split(s, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// 查找运算符
		var r rule
		if idx := strings.IndexAny(p, ":="); idx >= 0 {
			r.key = strings.TrimSpace(p[:idx])
			r.op = string(p[idx])
			r.val = strings.TrimSpace(p[idx+1:])
		} else if idx := strings.IndexAny(p, ">"); idx >= 0 {
			r.key = strings.TrimSpace(p[:idx])
			r.op = ">"
			r.val = strings.TrimSpace(p[idx+1:])
		} else if idx := strings.IndexAny(p, "<"); idx >= 0 {
			r.key = strings.TrimSpace(p[:idx])
			r.op = "<"
			r.val = strings.TrimSpace(p[idx+1:])
		} else if strings.Contains(p, "=") { // fallback
			if kv := strings.SplitN(p, "=", 2); len(kv) == 2 {
				r.key = strings.TrimSpace(kv[0])
				r.op = "="
				r.val = strings.TrimSpace(kv[1])
			} else {
				r.key = p
			}
		} else {
			r.key = p
		}
		out = append(out, r)
	}
	return out
}

// helper: apply rules when value is nil pointer
func applyRulesToNil(ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		if r.key == "required" {
			errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
		}
	}
	return errs
}

func applyRulesToAny(val any, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		if r.key == "required" {
			// 判断零值
			if isZeroValue(reflect.ValueOf(val)) {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		}
	}
	return errs
}

func isZeroValue(v reflect.Value) bool {
	z := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), z.Interface())
}

// checkString 对字符串的验证
func checkString(s, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		switch r.key {
		case "required":
			if strings.TrimSpace(s) == "" {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		case "nullable":
			// nothing
		case "min":
			if n, err := strconv.Atoi(r.val); err == nil {
				if len([]rune(s)) < n {
					errs = append(errs, fmt.Sprintf("[%s]长度不能小于 %d", fieldName, n))
				}
			}
		case "max":
			if n, err := strconv.Atoi(r.val); err == nil {
				if len([]rune(s)) > n {
					errs = append(errs, fmt.Sprintf("[%s]长度不能大于 %d", fieldName, n))
				}
			}
		case "len":
			if n, err := strconv.Atoi(r.val); err == nil {
				if len([]rune(s)) != n {
					errs = append(errs, fmt.Sprintf("[%s]长度必须为 %d", fieldName, n))
				}
			}
		case "in":
			if !inStringList(s, r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值必须在[%s]中", fieldName, r.val))
			}
		case "not-in":
			if inStringList(s, r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值不可为[%s]", fieldName, r.val))
			}
		case "regex":
			ok := matchRegex(r.val, s)
			if !ok {
				errs = append(errs, fmt.Sprintf("[%s]不匹配正则表达式", fieldName))
			}
		case "email":
			if !isEmail(s) {
				errs = append(errs, fmt.Sprintf("[%s]不是有效的邮箱", fieldName))
			}
		case "date":
			if !isDate(s) {
				errs = append(errs, fmt.Sprintf("[%s]不是有效的日期(YYYY-MM-DD)", fieldName))
			}
		case "time":
			if !isTime(s) {
				errs = append(errs, fmt.Sprintf("[%s]不是有效的时间(HH:MM:SS)", fieldName))
			}
		case "datetime":
			if !isDateTime(s) {
				errs = append(errs, fmt.Sprintf("[%s]不是有效的日期时间(YYYY-MM-DD HH:MM:SS)", fieldName))
			}
		default:
			// 兼容形式: min>3 / max<10 / key>n / key<n
			if r.op == ">" {
				if r.key == "min" {
					if n, err := strconv.Atoi(r.val); err == nil {
						if len([]rune(s)) < n {
							errs = append(errs, fmt.Sprintf("[%s]长度不能小于 %d", fieldName, n))
						}
					}
				}
			} else if r.op == "<" {
				if r.key == "max" {
					if n, err := strconv.Atoi(r.val); err == nil {
						if len([]rune(s)) > n {
							errs = append(errs, fmt.Sprintf("[%s]长度不能大于 %d", fieldName, n))
						}
					}
				}
			}
		}
	}
	return errs
}

func inStringList(s, list string) bool {
	list = strings.TrimSpace(list)
	if list == "" {
		return false
	}
	parts := strings.Split(list, ",")
	for _, p := range parts {
		if strings.TrimSpace(p) == s {
			return true
		}
	}
	return false
}

func matchRegex(pattern, s string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

func isEmail(s string) bool {
	// 简单邮箱校验
	if s == "" {
		return false
	}
	re := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return re.MatchString(s)
}

func isDate(s string) bool {
	if _, err := time.Parse("2006-01-02", s); err == nil {
		return true
	}
	return false
}

func isTime(s string) bool {
	if _, err := time.Parse("15:04:05", s); err == nil {
		return true
	}
	return false
}

func isDateTime(s string) bool {
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339}
	for _, l := range layouts {
		if _, err := time.Parse(l, s); err == nil {
			return true
		}
	}
	return false
}

// checkSlice 校验数组/切片长度及简单 in/not-in 对元素的约束
func checkSlice(v reflect.Value, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	length := v.Len()
	for _, r := range rules {
		switch r.key {
		case "min":
			if n, err := strconv.Atoi(r.val); err == nil {
				if length < n {
					errs = append(errs, fmt.Sprintf("[%s]长度不能小于 %d", fieldName, n))
				}
			}
		case "max":
			if n, err := strconv.Atoi(r.val); err == nil {
				if length > n {
					errs = append(errs, fmt.Sprintf("[%s]长度不能大于 %d", fieldName, n))
				}
			}
		case "len":
			if n, err := strconv.Atoi(r.val); err == nil {
				if length != n {
					errs = append(errs, fmt.Sprintf("[%s]长度必须为 %d", fieldName, n))
				}
			}
		case "required":
			if length == 0 {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		}
	}
	return errs
}

func checkNumberInt(val int64, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		switch r.key {
		case "required":
			// 零值检查视为必填
			if val == 0 {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		case "min":
			if n, err := strconv.ParseInt(r.val, 10, 64); err == nil {
				if val < n {
					errs = append(errs, fmt.Sprintf("[%s]不能小于 %d", fieldName, n))
				}
			}
		case "max":
			if n, err := strconv.ParseInt(r.val, 10, 64); err == nil {
				if val > n {
					errs = append(errs, fmt.Sprintf("[%s]不能大于 %d", fieldName, n))
				}
			}
		case "in":
			if !inStringList(fmt.Sprintf("%d", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值必须在[%s]中", fieldName, r.val))
			}
		case "not-in":
			if inStringList(fmt.Sprintf("%d", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值不可为[%s]", fieldName, r.val))
			}
		default:
			if r.op == ">" {
				if n, err := strconv.ParseInt(r.val, 10, 64); err == nil {
					if val < n {
						errs = append(errs, fmt.Sprintf("[%s]不能小于 %d", fieldName, n))
					}
				}
			} else if r.op == "<" {
				if n, err := strconv.ParseInt(r.val, 10, 64); err == nil {
					if val > n {
						errs = append(errs, fmt.Sprintf("[%s]不能大于 %d", fieldName, n))
					}
				}
			}
		}
	}
	return errs
}

func checkNumberUint(val uint64, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		switch r.key {
		case "required":
			if val == 0 {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		case "min":
			if n, err := strconv.ParseUint(r.val, 10, 64); err == nil {
				if val < n {
					errs = append(errs, fmt.Sprintf("[%s]不能小于 %d", fieldName, n))
				}
			}
		case "max":
			if n, err := strconv.ParseUint(r.val, 10, 64); err == nil {
				if val > n {
					errs = append(errs, fmt.Sprintf("[%s]不能大于 %d", fieldName, n))
				}
			}
		case "in":
			if !inStringList(fmt.Sprintf("%d", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值必须在[%s]中", fieldName, r.val))
			}
		case "not-in":
			if inStringList(fmt.Sprintf("%d", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值不可为[%s]", fieldName, r.val))
			}
		}
	}
	return errs
}

func checkNumberFloat(val float64, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		switch r.key {
		case "required":
			if val == 0 {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		case "min":
			if n, err := strconv.ParseFloat(r.val, 64); err == nil {
				if val < n {
					errs = append(errs, fmt.Sprintf("[%s]不能小于 %v", fieldName, n))
				}
			}
		case "max":
			if n, err := strconv.ParseFloat(r.val, 64); err == nil {
				if val > n {
					errs = append(errs, fmt.Sprintf("[%s]不能大于 %v", fieldName, n))
				}
			}
		case "in":
			if !inStringList(fmt.Sprintf("%v", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值必须在[%s]中", fieldName, r.val))
			}
		case "not-in":
			if inStringList(fmt.Sprintf("%v", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值不可为[%s]", fieldName, r.val))
			}
		}
	}
	return errs
}

func checkBool(val bool, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	for _, r := range rules {
		switch r.key {
		case "required":
			// 对 bool 来说，无法区分零值和未赋值，通常不检查
		case "in":
			if !inStringList(fmt.Sprintf("%v", val), r.val) {
				errs = append(errs, fmt.Sprintf("[%s]值必须在[%s]中", fieldName, r.val))
			}
		}
	}
	return errs
}

// checkTimeValue 用于 time.Time 类型的字段
func checkTimeValue(val any, ruleStr, fieldName string) []string {
	rules := parseRules(ruleStr)
	errs := make([]string, 0)
	t, ok := val.(time.Time)
	if !ok {
		return append(errs, fmt.Sprintf("[%s]不是时间类型", fieldName))
	}
	for _, r := range rules {
		switch r.key {
		case "required":
			if t.IsZero() {
				errs = append(errs, fmt.Sprintf("[%s]为必填项", fieldName))
			}
		case "datetime":
			if t.IsZero() {
				errs = append(errs, fmt.Sprintf("[%s]不是有效的日期时间", fieldName))
			}
		}
	}
	return errs
}
