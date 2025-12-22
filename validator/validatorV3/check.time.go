package validatorV3

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aid297/aid/array/anyArrayV2"
)

// checkString 检查字符串，支持：required、[datetime]、min>、min>=、max<、max<=、in、not-in、size:
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

	return my
}
