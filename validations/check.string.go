package validations

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/aid297/aid/v2/anyMaps"
	"github.com/aid297/aid/v2/anySlices"
	"github.com/aid297/aid/v2/consts/textTime"
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

// checkString 检查字符串，支持：required、[bool|datetime|date|timer]、str-timer>、str-timer>=、str-timer<、str-timer<=、regex==、regex!=、min>、min>=、max<、max<=、in==、in!=、size==、size!=, ex:
func (my FieldInfo) checkString() FieldInfo {
	var (
		err                    error
		min, max, size         *int
		include, eq            bool
		in                     []string
		notIn                  []string
		value                  string
		ok                     bool
		pattern                string
		strTimeMin, strTimeMax *string
		stMin, stMax, vt       time.Duration
	)

	if my.Kind != reflect.String {
		my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(my.VRuleTags) {
		if my.IsPtr && (my.IsNil || my.IsZero) {
			my.wrongs = []error{fmt.Errorf("『%s』 %w", my.getName(), ErrRequired)}
			return my
		} else if !my.IsPtr && my.IsZero {
			my.wrongs = []error{fmt.Errorf("『%s』 %w", my.getName(), ErrNotEmpty)}
			return my
		}
	}

	if value, ok = my.Value.(string); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	my.VRuleTags.Each(func(_ int, rule string) (isBreak bool) {
		if strings.HasPrefix(rule, "str-timer") {
			if strTimeMin, include = getRuleStrTimeMin(rule); strTimeMin != nil {
				if stMin, err = textTime.WhatTimeIsIt(*strTimeMin); err != nil {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w(规则)", my.getName(), ErrInvalidFormat))
				} else {
					if vt, err = textTime.WhatTimeIsIt(value); err != nil {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w(值)", my.getName(), ErrInvalidFormat))
					} else {
						if include {
							if !(vt >= stMin) {
								my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：>= %s", my.getName(), ErrInvalidValue, *strTimeMin))
							}
						} else {
							if !(vt > stMin) {
								my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：> %s", my.getName(), ErrInvalidValue, *strTimeMin))
							}
						}
					}
				}
			}

			if strTimeMax, include = getRuleStrTimeMax(rule); strTimeMax != nil {
				if stMax, err = textTime.WhatTimeIsIt(*strTimeMax); err != nil {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w(规则)", my.getName(), ErrInvalidFormat))
				} else {
					if include {
						if !(vt <= stMax) {
							my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：<= %s", my.getName(), ErrInvalidValue, *strTimeMax))
						}
					} else {
						if !(vt < stMax) {
							my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：< %s", my.getName(), ErrInvalidValue, *strTimeMax))
						}
					}
				}
			}
		} else if strings.HasPrefix(rule, "min") {
			if min, include = getRuleIntMin(rule); min != nil {
				if include {
					if !(utf8.RuneCountInString(value) >= *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：>= %d", my.getName(), ErrInvalidLength, *min))
					}
				} else {
					if !(utf8.RuneCountInString(value) > *min) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：> %d", my.getName(), ErrInvalidLength, *min))
					}
				}
			}
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleIntMax(rule); max != nil {
				if include {
					if !(utf8.RuneCountInString(value) <= *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：<= %d", my.getName(), ErrInvalidLength, *max))
					}
				} else {
					if !(utf8.RuneCountInString(value) < *max) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：< %d", my.getName(), ErrInvalidLength, *max))
					}
				}
			}
		} else if strings.HasPrefix(rule, "in==") {
			if in = getRuleIn(rule); len(in) > 0 {
				anySlices.New(anySlices.List(in)).IfNotIn(func(_ anySlices.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, value)
			}
		} else if strings.HasPrefix(rule, "in!=") {
			if notIn = getRuleNotIn(rule); len(notIn) > 0 {
				anySlices.New(anySlices.List(notIn)).IfIn(func(_ anySlices.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, value)
			}
		} else if strings.HasPrefix(rule, "size") {
			if size, eq = getRuleIntSize(rule); size != nil {
				if eq {
					if !(len(value) == *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：不等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				} else {
					if !(len(value) != *size) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：等于 %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		} else if rule == "bool" {
			var def = []string{"true", "True", "t", "yes", "on", "ok", "1", "false", "False", "f", "off", "no", "0"}
			if strings.HasPrefix(rule, "in==") {
				if in = getRuleIn(rule); len(in) > 0 {
					anySlices.New(anySlices.List(in)).IfNotIn(func(_ anySlices.AnySlicer[string]) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, value)
				}
			}
			if strings.HasPrefix(rule, "in!=") {
				if notIn = getRuleNotIn(rule); len(notIn) > 0 {
					anySlices.New(anySlices.List(notIn)).IfIn(func(_ anySlices.AnySlicer[string]) {
						my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, value)
				}
			}
			if len(in) == 0 && len(notIn) == 0 {
				anySlices.New(anySlices.List(def)).IfNotIn(func(_ anySlices.AnySlicer[string]) {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, def))
				}, value)
			}
		} else if rule == "datetime" {
			ok = false
			anyMaps.New(anyMaps.Map(patternsForTimeString)).RemoveByKeys("DateOnly", "TimeOnly").Each(func(_ string, value string) {
				if regexp.MustCompile(value).MatchString(value) {
					ok = true
					return
				}
			})
			if !ok {
				my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w", my.getName(), ErrInvalidFormat))
			}
		} else if rule == "date" {
			if !regexp.MustCompile(patternsForTimeString["DateOnly"]).MatchString(value) {
				my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w", my.getName(), ErrInvalidFormat))
			}
		} else if rule == "timer" {
			if !regexp.MustCompile(patternsForTimeString["TimeOnly"]).MatchString(value) {
				my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w", my.getName(), ErrInvalidFormat))
			}
		} else if strings.HasPrefix(rule, "regex==") {
			if pattern = getRuleRegexEq(rule); pattern != "" {
				if !regexp.MustCompile(pattern).MatchString(value) {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：匹配正则 %s", my.getName(), ErrInvalidFormat, pattern))
				}
			}
		} else if strings.HasPrefix(rule, "regex!=") {
			if pattern = getRuleRegexNotEq(rule); pattern != "" {
				if regexp.MustCompile(pattern).MatchString(value) {
					my.wrongs = append(my.wrongs, fmt.Errorf("『%s』 %w 期望：不匹配正则 %s", my.getName(), ErrInvalidFormat, pattern))
				}
			}
		} else if strings.HasPrefix(rule, "ex") {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for _, exFnName := range exFnNames {
					if fn := Once().GetExFn(exFnName); fn != nil {
						if err := fn(my.getName(), my.Value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}

		return
	})

	return my
}
