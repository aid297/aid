package validatorV2

import "reflect"

type (
	OriginalDatum struct {
		refValue reflect.Value
		refType  reflect.Type
		refField reflect.StructField
		isPtr    bool
		isNil    bool
		isZero   bool
		indirect reflect.Value
	}
	Field struct {
	}
)

func NewOriginalDatum(value reflect.Value, field reflect.StructField) *OriginalDatum {
	return &OriginalDatum{
		refValue: value,
		refType:  value.Type(),
		refField: field,
		isPtr:    value.Kind() == reflect.Ptr,
		isNil:    value.Kind() == reflect.Ptr && value.IsNil(),
		isZero:   value.IsZero(),
		indirect: reflect.Indirect(value),
	}
}
