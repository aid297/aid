package validatorV3

import (
	"fmt"
	"strings"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkTime 检查时间，支持：required、[datetime|date|time]、min>、min>=、max<、max<=、in、not-in、ex:
func (my FieldInfo) checkTime() FieldInfo {
	var (
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

	if getRuleRequired(my.VRuleTags) {
		if my.IsPtr && my.IsNil {
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
			return my
		} else if !my.IsPtr && value.IsZero() {
			my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrNotEmpty)}
			return my
		}
	}

	my.VRuleTags.Each(func(_ int, rule string) {
		if strings.HasPrefix(rule, "min") {
			if min, include = getRuleTimeMin(rule); min != nil {
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
		} else if strings.HasPrefix(rule, "max") {
			if max, include = getRuleTimeMax(rule); max != nil {
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
		} else if strings.HasPrefix(rule, "in") {
			if in = getRuleTimeIn(rule); len(in) > 0 {
				anyArrayV2.NewList(in).IfIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之中", my.getName(), ErrInvalidValue, in))
				}, value)
			}
		} else if strings.HasPrefix(rule, "not-in") {
			if notIn = getRuleTimeNotIn(rule); len(notIn) > 0 {
				anyArrayV2.NewList(notIn).IfNotIn(func() {
					my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：在 %v 之外", my.getName(), ErrInvalidValue, notIn))
				}, value)
			}
		} else if rule == "ex" {
			if exFnNames := getRuleExFnNames(rule); len(exFnNames) > 0 {
				for idx2 := range exFnNames {
					if fn := APP.Validator.Once().GetExFn(exFnNames[idx2]); fn != nil {
						if err := fn(my.Value); err != nil {
							my.wrongs = append(my.wrongs, err)
						}
					}
				}
			}
		}
	})

	return my
}
