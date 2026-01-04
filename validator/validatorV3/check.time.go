package validatorV3

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkTime 检查时间，支持：required、min>、min>=、max<、max<=、in、not-in、ex:
func (my FieldInfo) checkTime() FieldInfo {
	if getRuleRequired(anyArrayV2.NewList(my.VRuleTags)) && my.IsPtr && my.IsNil {
		my.wrongs = []error{fmt.Errorf("[%s] %w", my.getName(), ErrRequired)}
		return my
	}

	switch my.Value.(type) {
	case time.Time:
		v := reflect.ValueOf(my.Value)
		if !v.IsZero() {
			my.wrongs = append(my.wrongs, fmt.Errorf("[%s] %w 期望：时间类型", my.getName(), ErrInvalidType))
		}

	}

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

	return my
}
