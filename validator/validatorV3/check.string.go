package validatorV3

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkString 检查字符串，支持：required、[string]、min>、min>=、max<、max<=、in、not-in、size:
func (my FieldInfo) checkString() FieldInfo {
	var (
		rules          = anyArrayV2.NewList(my.VRuleTags)
		ruleType       = my.getRuleType(rules)
		min, max, size *int
		include        bool
		in             []string
		notIn          []string
		value          string
		ok             bool
	)

	if value, ok = my.Value.(string); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	for idx := range my.VRuleTags {
		switch ruleType {
		case "", "string":
			if strings.HasPrefix(my.VRuleTags[idx], "min") {
				if min, include = getRuleIntMin(my.VRuleTags[idx]); min != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "max") {
				if max, include = getRuleIntMax(my.VRuleTags[idx]); max != nil {
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
			if strings.HasPrefix(my.VRuleTags[idx], "in") {
				if in = getRuleIn(my.VRuleTags[idx]); len(in) > 0 {
					anyArrayV2.NewList(in).IfNotIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
					}, value)
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "not-in") {
				if notIn = getRuleNotIn(my.VRuleTags[idx]); len(notIn) > 0 {
					anyArrayV2.NewList(notIn).IfIn(func() {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
					}, value)
				}
			}
			if strings.HasPrefix(my.VRuleTags[idx], "size") {
				if size = getRuleIntSize(my.VRuleTags[idx]); size != nil {
					if len(value) != *size {
						my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：= %d", my.getName(), ErrInvalidLength, *size))
					}
				}
			}
		}
	}

	return my
}
