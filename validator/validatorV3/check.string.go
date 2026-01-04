package validatorV3

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/dict/anyDictV2"
)

var (
	patternsForTimeString = map[string]string{
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
)

// checkString 检查字符串，支持：required、not-empty、[string|bool|datetime|date|time]、min>、min>=、max<、max<=、in、not-in、size=、size<=, ex:
func (my FieldInfo) checkString() FieldInfo {
	var (
		min, max, size *int
		include, eq    bool
		in             []string
		notIn          []string
		value          string
		ok             bool
		// needRequired   = getRuleRequired(my.VRuleTags)
	)

	if value, ok = my.Value.(string); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	if getRuleNotEmpty(my.VRuleTags) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	// if needRequired && my.IsPtr && my.IsNil {
	// 	my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
	// 	return my
	// } else if !needRequired && !my.IsPtr && value == "" {
	// 	return my
	// } else if !needRequired && my.IsPtr && my.IsNil {
	// 	return my
	// }

	my.VRuleTags.Each(func(_ int, rule string) {
		switch rule {
		case "", "string":
			if strings.HasPrefix(rule, "min") {
				if min, include = getRuleIntMin(rule); min != nil {
					if include {
						if !(utf8.RuneCountInString(value) >= *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
						}
					} else {
						if !(utf8.RuneCountInString(value) > *min) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "max") {
				if max, include = getRuleIntMax(rule); max != nil {
					if include {
						if !(utf8.RuneCountInString(value) <= *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
						}
					} else {
						if !(utf8.RuneCountInString(value) < *max) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
						}
					}
				}
			}
			if strings.HasPrefix(rule, "in") {
				if in = getRuleIn(rule); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, value)
				}
			}
			if strings.HasPrefix(rule, "not-in") {
				if notIn = getRuleNotIn(rule); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, value)
				}
			}
			if strings.HasPrefix(rule, "size") {
				if size, eq = getRuleIntSize(rule); size != nil {
					if eq {
						if !(len(value) == *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					} else {
						if !(len(value) != *size) {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
						}
					}
				}
			}
			fallthrough
		case "bool":
			var def = []string{"true", "True", "t", "yes", "on", "ok", "1", "false", "False", "f", "off", "no", "0"}
			if strings.HasPrefix(rule, "in") {
				if in = getRuleIn(rule); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, value)
				}
			}
			if strings.HasPrefix(rule, "not-in") {
				if notIn = getRuleNotIn(rule); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, value)
				}
			}
			if len(in) == 0 && len(notIn) == 0 {
				anyArrayV2.NewList(def).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, def))
				}, value)
			}
			fallthrough
		case "datetime":
			ok = false
			anyDictV2.New(anyDictV2.Map(patternsForTimeString)).RemoveByKeys("DateOnly", "TimeOnly").Each(func(_ string, value string) {
				if regexp.MustCompile(value).MatchString(value) {
					ok = true
					return
				}
			})
			if !ok {
				my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrInvalidFormat))
			}
			fallthrough
		case "date":
			if !regexp.MustCompile(patternsForTimeString["DateOnly"]).MatchString(value) {
				my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrInvalidFormat))
			}
			fallthrough
		case "time":
			if !regexp.MustCompile(patternsForTimeString["TimeOnly"]).MatchString(value) {
				my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrInvalidFormat))
			}
			fallthrough
		case "ex":
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	})

	return my
}
