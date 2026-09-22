package validations

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// CrossValidator 跨字段校验接口。
// 结构体可实现此方法返回跨字段规则字符串，Validate 过程中会自动检测并调用。
//
//	规则格式：({field1}>{field2})，多条规则用 () 包裹：({f1}<{f2})({f3}>={f4})
//	支持运算符：>、>=、<、<=、==、!=
type CrossValidator interface {
	Cross() string
}

// crossRule 表示一条已解析的跨字段比较规则。
// 语义：LeftField [Operator] RightField
// 例如 {start_time}<{end_time} 表示 start_time 必须小于 end_time。
type crossRule struct {
	LeftField  string // 左侧字段引用（v-name 路径）
	Operator   string // 比较运算符
	RightField string // 右侧字段引用（v-name 路径）
}

// 匹配 {field1}op{field2} 格式
var crossRuleRegexp = regexp.MustCompile(`\{([^}]+)\}\s*(>=|<=|!=|==|>|<)\s*\{([^}]+)\}`)

// parseCrossRules 解析 Cross() 返回的规则字符串。
// 支持 (rule1)(rule2) 格式或单条规则。
//
//	示例：
//	  "({start_time}<{end_time})"
//	  "({a}<{b})({c}>={d})"
func parseCrossRules(ruleStr string) []crossRule {
	ruleStr = strings.TrimSpace(ruleStr)
	if ruleStr == "" {
		return nil
	}

	// 多条规则：(rule1)(rule2)
	if strings.HasPrefix(ruleStr, "(") && strings.HasSuffix(ruleStr, ")") {
		inner := ruleStr[1 : len(ruleStr)-1]
		parts := strings.Split(inner, ")(")
		rules := make([]crossRule, 0, len(parts))
		for _, part := range parts {
			if r := parseOneCrossRule(strings.TrimSpace(part)); r != nil {
				rules = append(rules, *r)
			}
		}
		return rules
	}

	// 单条规则
	if r := parseOneCrossRule(ruleStr); r != nil {
		return []crossRule{*r}
	}
	return nil
}

func parseOneCrossRule(s string) *crossRule {
	matches := crossRuleRegexp.FindStringSubmatch(s)
	if len(matches) < 4 {
		return nil
	}
	return &crossRule{
		LeftField:  matches[1],
		Operator:   matches[2],
		RightField: matches[3],
	}
}

// evalCrossRule 对一条跨字段规则求值。
// 返回 nil 表示通过，否则返回描述错误的 error。
func evalCrossRule(rule crossRule, infoMap map[string]FieldInfo) error {
	leftInfo, leftOK := infoMap[rule.LeftField]
	rightInfo, rightOK := infoMap[rule.RightField]

	if !leftOK || !rightOK {
		return nil // 引用字段不存在，静默跳过
	}

	// 任一字段独立校验未通过 → 跳过
	if len(leftInfo.wrongs) > 0 || len(rightInfo.wrongs) > 0 {
		return nil
	}

	// 任一字段为 nil 指针 → 跳过
	if (leftInfo.IsPtr && leftInfo.IsNil) || (rightInfo.IsPtr && rightInfo.IsNil) {
		return nil
	}

	leftVal := extractComparableValue(leftInfo)
	rightVal := extractComparableValue(rightInfo)

	if leftVal == nil || rightVal == nil {
		return nil // 类型不支持，静默跳过
	}

	passed, err := compareValues(leftVal, rightVal, rule.Operator)
	if err != nil {
		return nil // 类型不兼容，静默跳过
	}
	if !passed {
		return fmt.Errorf(
			"『%s』 %s 『%s』 跨字段校验失败",
			rule.LeftField, rule.Operator, rule.RightField,
		)
	}
	return nil
}

// extractComparableValue 从 FieldInfo 中提取可比较的值。
// 返回 time.Time、float64 或 string，不支持的类型返回 nil。
func extractComparableValue(info FieldInfo) any {
	v := info.Value
	if v == nil {
		return nil
	}

	if t, ok := v.(time.Time); ok {
		return t
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	case reflect.String:
		return rv.String()
	default:
		return nil
	}
}

// compareValues 根据运算符比较两个同类型值。
func compareValues(a, b any, op string) (bool, error) {
	if ta, ok := a.(time.Time); ok {
		if tb, ok := b.(time.Time); ok {
			return compareTime(ta, tb, op)
		}
		return false, fmt.Errorf("类型不匹配")
	}

	if fa, ok := toFloat64(a); ok {
		if fb, ok := toFloat64(b); ok {
			return compareNumeric(fa, fb, op)
		}
		return false, fmt.Errorf("类型不匹配")
	}

	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return compareString(sa, sb, op)
		}
		return false, fmt.Errorf("类型不匹配")
	}

	return false, fmt.Errorf("不支持的比较类型")
}

func compareTime(a, b time.Time, op string) (bool, error) {
	switch op {
	case ">":
		return a.After(b), nil
	case ">=":
		return !a.Before(b), nil
	case "<":
		return a.Before(b), nil
	case "<=":
		return !a.After(b), nil
	case "==":
		return a.Equal(b), nil
	case "!=":
		return !a.Equal(b), nil
	default:
		return false, fmt.Errorf("未知运算符：%s", op)
	}
}

func compareNumeric(a, b float64, op string) (bool, error) {
	switch op {
	case ">":
		return a > b, nil
	case ">=":
		return a >= b, nil
	case "<":
		return a < b, nil
	case "<=":
		return a <= b, nil
	case "==":
		return a == b, nil
	case "!=":
		return a != b, nil
	default:
		return false, fmt.Errorf("未知运算符：%s", op)
	}
}

func compareString(a, b string, op string) (bool, error) {
	switch op {
	case ">":
		return a > b, nil
	case ">=":
		return a >= b, nil
	case "<":
		return a < b, nil
	case "<=":
		return a <= b, nil
	case "==":
		return a == b, nil
	case "!=":
		return a != b, nil
	default:
		return false, fmt.Errorf("未知运算符：%s", op)
	}
}

func toFloat64(v any) (float64, bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false
	}
}
