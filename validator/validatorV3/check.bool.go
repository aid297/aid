package validatorV3

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkBool 检查布尔值，支持：required、[bool|b]、in、not-in：
func (my FieldInfo) checkBool() FieldInfo {
	var (
		rules    = anyArrayV2.NewList(my.VRuleTags)
		ruleType = my.getRuleType(rules)
		in       []string
		notIn    []string
		value    string
		ok       bool
		def      = []string{"true", "True", "t", "yes", "on", "ok", "1", "false", "False", "f", "off", "no", "0"}
	)

	if value, ok = my.Value.(string); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	for idx := range my.VRuleTags {
		if my.VRuleTags[idx] == "required" && my.IsPtr && my.IsNil {
			my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w", my.getName(), ErrRequired))
		}

		switch ruleType {
		case "", "bool", "b":
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
			if len(in) == 0 && len(notIn) == 0 {
				anyArrayV2.NewList(def).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, def))
				}, value)
			}
		}

		if strings.HasPrefix(my.VRuleTags[idx], "ex") {
			if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Ins().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	}

	return my
}
