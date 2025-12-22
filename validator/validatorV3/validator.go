package validatorV3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/str"
)

type (
	// Validator 验证器
	Validator struct {
		data   any
		wrongs []error
	}
)

func (Validator) New(data any) Validator { return Validator{data: data} }

func (my Validator) Wrongs() []error { return my.wrongs }

func (my Validator) WrongToString() string {
	var errs = make([]string, 0, len(my.wrongs))

	for idx := range my.wrongs {
		errs = append(errs, my.wrongs[idx].Error())
	}

	return str.APP.Buffer.JoinStringLimit("；", errs...)
}

func (my Validator) Validate(exCheckFns ...any) Validator {
	fieldInfos := getStructFieldInfos(my.data, "")
	for idx := range fieldInfos {
		if wrongs := fieldInfos[idx].Check().Wrongs(); len(wrongs) > 0 {
			my.wrongs = append(my.wrongs, wrongs...)
		}
	}

	if len(my.wrongs) == 0 {
		for idx := range exCheckFns {
			if err := callExCheckFn(exCheckFns[idx], my.data); err != nil {
				my.wrongs = append(my.wrongs, err)
			}
		}
	}

	return my
}

func WithGin[T any](c *gin.Context, exCheckFns ...any) Validator {
	var form = new(T)

	if err := c.ShouldBind(&form); err != nil {
		return Validator{wrongs: []error{err}}
	}

	return APP.Validator.New(form).Validate(exCheckFns)
}

func WithFiber[T any](c *fiber.Ctx, exCheckFns ...any) Validator {
	var form = new(T)

	if err := c.BodyParser(&form); err != nil {
		return Validator{wrongs: []error{err}}
	}

	return APP.Validator.New(form).Validate(exCheckFns)
}

func callExCheckFn(fn any, data any) error {
	if fn == nil {
		return fmt.Errorf("callback is nil")
	}
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return fmt.Errorf("callback is not a function: %T", fn)
	}
	ft := fv.Type()
	if ft.NumIn() != 1 || ft.NumOut() < 1 {
		return fmt.Errorf("callback must have signature func(T) error (or similar), got %s", ft.String())
	}

	argType := ft.In(0)
	var dv reflect.Value
	if data == nil {
		dv = reflect.Zero(argType)
	} else {
		dv = reflect.ValueOf(data)
		// If direct assignable, OK. Otherwise try to adapt:
		if !dv.Type().AssignableTo(argType) {
			// If function expects a pointer and we have a non-pointer of compatible element, take address.
			if argType.Kind() == reflect.Ptr && dv.Type().AssignableTo(argType.Elem()) {
				addr := reflect.New(dv.Type())
				addr.Elem().Set(dv)
				dv = addr
			} else if dv.Kind() == reflect.Ptr && dv.Type().Elem().AssignableTo(argType) {
				// If we have a pointer but function expects a value, dereference
				dv = dv.Elem()
			} else if dv.CanAddr() && dv.Addr().Type().AssignableTo(argType) {
				// If we have an addressable value and function expects that pointer type
				dv = dv.Addr()
			} else {
				// last resort: try zero value of argType
				dv = reflect.Zero(argType)
			}
		}
	}

	outs := fv.Call([]reflect.Value{dv})
	if len(outs) == 0 {
		return nil
	}
	first := outs[0]
	if first.IsNil() {
		return nil
	}
	errIface := reflect.TypeOf((*error)(nil)).Elem()
	if !first.Type().Implements(errIface) {
		return fmt.Errorf("callback first return does not implement error: %s", first.Type().String())
	}
	return first.Interface().(error)
}

func getStructFieldInfos(s any, parentName string) []FieldInfo {
	v := reflect.ValueOf(s)
	t := v.Type()

	var infos []FieldInfo

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// 获取 tags
			vRuleTag := field.Tag.Get("v-rule")
			vNameTag := field.Tag.Get("v-name")

			if vNameTag == "" {
				vNameTag = field.Name
			}

			if vRuleTag == "" || vRuleTag == "-" {
				continue
			}

			isPtr := fieldValue.Kind() == reflect.Ptr
			isNil := isPtr && fieldValue.IsNil()

			// determine the element/type kind safely before dereferencing
			var elemType reflect.Type
			var elemKind reflect.Kind
			if isPtr {
				elemType = fieldValue.Type().Elem()
				elemKind = elemType.Kind()
			} else {
				elemType = fieldValue.Type()
				elemKind = elemType.Kind()
			}

			var value any = nil
			// only dereference if pointer and not nil, otherwise keep value as nil or the concrete value
			if isPtr && !isNil {
				fieldValue = fieldValue.Elem()
				value = fieldValue.Interface()
			} else if !isPtr {
				value = fieldValue.Interface()
			}

			// If the declared/element kind is struct, recurse.
			// For a nil pointer-to-struct, pass a zero value of the element type so reflection works.
			if elemKind == reflect.Struct {
				var recurseArg any
				if isPtr && isNil {
					recurseArg = reflect.Zero(elemType).Interface()
				} else {
					recurseArg = value
				}
				infos = append(infos, getStructFieldInfos(recurseArg, vNameTag)...)
				continue
			}

			infos = append(infos, FieldInfo{
				Name:      field.Name,
				Value:     value,
				Kind:      elemKind,
				Type:      elemType,
				IsPtr:     isPtr,
				IsNil:     isNil,
				VRuleTags: strings.Split(vRuleTag, ";"),
				VNameTags: anyArrayV2.NewItems(parentName, vNameTag).RemoveEmpty().ToSlice(),
			})
		}
	}

	return infos
}
