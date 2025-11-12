package validator

import (
	"reflect"
	"time"

	"github.com/aid297/aid/operation"
)

func (my *ValidatorApp[T]) checkTime(rule, fieldName string, value any) error {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if reflect.ValueOf(value).IsNil() {
			return operation.Ternary(rule == "required", RequiredErr.New(fieldName), nil)
		}
		value = reflect.ValueOf(value).Elem().Interface()
	}

	if !reflect.DeepEqual(value, time.Time{}) {
		return TimeErr.NewFormat("[%s]必须是时间类型", fieldName)
	}
	return nil
}
