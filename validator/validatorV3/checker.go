package validatorV3

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"

	"github.com/aid297/aid/array/anyArrayV2"
	"github.com/aid297/aid/operation/operationV2"
)

// Checker 验证器
type Checker struct {
	data         any
	wrongs       []error
	defaultLimit string
}

func (Checker) New(data any) Checker { return Checker{data: data, defaultLimit: "<br />"} }

func (my Checker) Wrongs() []error { return my.wrongs }

func (my Checker) OK() bool { return len(my.wrongs) == 0 }

func (my Checker) Wrong() error {
	return operationV2.NewTernary(operationV2.TrueFn(func() error { return errors.New(my.WrongToString("")) })).GetByValue(len(my.wrongs) > 0)
}

func (my Checker) WrongToString(limit string) string {
	var errs = make([]string, len(my.wrongs))

	for idx := range my.wrongs {
		errs[idx] = fmt.Sprintf("问题%d：%s", idx+1, my.wrongs[idx].Error())
	}

	return strings.Join(errs, operationV2.NewTernary(operationV2.TrueValue(limit), operationV2.FalseValue(my.defaultLimit)).GetByValue(limit != ""))
}

func (my Checker) Validate(exCheckFns ...any) Checker {
	// fieldInfos := getStructFieldInfos(my.data, "")

	for _, fieldInfo := range getStructFieldInfos(my.data, "") {
		if wrongs := fieldInfo.Check().Wrongs(); len(wrongs) > 0 {
			my.wrongs = append(my.wrongs, wrongs...)
		}
	}

	// for idx := range fieldInfos {
	// 	if wrongs := fieldInfos[idx].Check().Wrongs(); len(wrongs) > 0 {
	// 		my.wrongs = append(my.wrongs, wrongs...)
	// 	}
	// }

	if len(my.wrongs) == 0 {
		for idx := range exCheckFns {
			if err := callExCheckFn(exCheckFns[idx], my.data); err != nil {
				my.wrongs = append(my.wrongs, err)
			}
		}
	}

	return my
}

func WithGin[T any](c *gin.Context, exCheckFns ...any) (form T, checker Checker) {
	form = *new(T)

	if err := c.ShouldBind(&form); err != nil {
		checker = Checker{wrongs: []error{err}}
		return
	}

	return form, APP.Validator.Once().Checker(&form).Validate(exCheckFns...)
}

func WithFiber[T any](c *fiber.Ctx, exCheckFns ...any) (form T, checker Checker) {
	form = *new(T)

	if err := c.BodyParser(&form); err != nil {
		checker = Checker{wrongs: []error{err}}
		return
	}

	return form, APP.Validator.Once().Checker(&form).Validate(exCheckFns...)
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

			vRuleTag = strings.TrimLeft(vRuleTag, "(")
			vRuleTag = strings.TrimRight(vRuleTag, ")")

			infos = append(infos, FieldInfo{
				Name:      field.Name,
				Value:     value,
				Kind:      elemKind,
				Type:      elemType,
				IsPtr:     isPtr,
				IsNil:     isNil,
				IsZero:    fieldValue.IsZero(),
				VRuleTags: anyArrayV2.NewList(strings.Split(vRuleTag, ")(")),
				VNameTags: anyArrayV2.NewItems(parentName, vNameTag).RemoveEmpty(),
			})
		}
	}

	return infos
}
