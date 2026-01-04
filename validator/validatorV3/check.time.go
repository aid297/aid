package validatorV3

import (
	"fmt"
	"strings"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkTime 检查时间，支持：required、not-empty、[datetime|date|time]、min>、min>=、max<、max<=、in、not-in、ex:
func (my FieldInfo) checkTime() FieldInfo {
	var (
		rules     = anyArrayV2.NewList(my.VRuleTags)
		min, max  *time.Time
		include   bool
		in, notIn []time.Time
		value     time.Time
		ok        bool
	)

	if value, ok = my.Value.(time.Time); !ok {
		my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：时间类型", my.getName(), ErrInvalidType))
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
		case "", "datetime", "date", "time":
			for idx := range my.VRuleTags {
				if strings.HasPrefix(my.VRuleTags[idx], "min") {
					if min, include = getRuleTimeMin(my.VRuleTags[idx]); min != nil {
						if include {
							if !(value.Unix() >= (*min).Unix()) {
								my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：>= %v", my.getName(), ErrInvalidLength, *min))
							}
						} else {
							if !(value.Unix() > (*min).Unix()) {
								my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：> %v", my.getName(), ErrInvalidLength, *min))
							}
						}
					}
				}
				if strings.HasPrefix(my.VRuleTags[idx], "max") {
					if max, include = getRuleTimeMax(my.VRuleTags[idx]); max != nil {
						if include {
							if !(value.Unix() <= (*max).Unix()) {
								my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：<= %v", my.getName(), ErrInvalidLength, *max))
							}
						} else {
							if !(value.Unix() < (*max).Unix()) {
								my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：< %v", my.getName(), ErrInvalidLength, *max))
							}
						}
					}
				}

				if strings.HasPrefix(my.VRuleTags[idx], "in") {
					if in = getRuleTimeIn(my.VRuleTags[idx]); len(in) > 0 {
						anyArrayV2.NewList(in).IfIn(func() {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
						}, value)
					}
				}

				if strings.HasPrefix(my.VRuleTags[idx], "not-in") {
					if notIn = getRuleTimeNotIn(my.VRuleTags[idx]); len(notIn) > 0 {
						anyArrayV2.NewList(notIn).IfNotIn(func() {
							my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
						}, value)
					}
				}
			}
			fallthrough
		case "ex":
			for idx := range my.VRuleTags {
				if strings.HasPrefix(my.VRuleTags[idx], "ex") {
					if exFnNames := getRuleExFnNames(my.VRuleTags[idx]); len(exFnNames) > 0 {
						for idx2 := range exFnNames {
							if fn := APP.Validator.Ins().GetExFn(exFnNames[idx2]); fn != nil {
								if err := fn(my.Value); err != nil {
									my.wrongs = append(my.wrongs, err)
								}
							}
						}
					}
				}
			}
		}
	})

	return my
}
