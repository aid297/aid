package validatorV3

import (
	"fmt"
	"strings"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkBool 检查布尔值，支持：required、not-empty、[string|bool]、in:、not-in：、ex:
func (my FieldInfo) checkBool() FieldInfo {
	var (
		rules = anyArrayV2.NewList(my.VRuleTags)
		in    []string
		notIn []string
		value string
		ok    bool
		def   = []string{"true", "True", "t", "yes", "on", "ok", "1", "false", "False", "f", "off", "no", "0"}
	)

	if value, ok = my.Value.(string); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：字符串", my.getName(), ErrInvalidType))
		return my
	}

	if getRuleRequired(rules) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	if getRuleNotEmpty(rules) && my.IsZero {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
		return my
	}

	rules.Each(func(_ int, rule string) {
		switch rule {
		case "", "string", "bool":
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
