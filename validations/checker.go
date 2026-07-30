package validations

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"

	"github.com/aid297/aid/v2/anySlices"
	"github.com/aid297/aid/v2/operations"
)

type (
	Checker interface {
		Errors() []error
		Error() error
		ErrorToString(limit string) (ret string)
		Invalid() bool
		OK() bool
		Validate(exCheckFns ...ExCheckFunc) Checker
	}

	// CheckerImpl 验证器
	CheckerImpl struct {
		data         any
		wrongs       []error
		defaultLimit string
	}
)

func NewCheck(data any) Checker { return &CheckerImpl{data: data, defaultLimit: "<br />"} }

func (my *CheckerImpl) Errors() []error { return my.wrongs }

func (my *CheckerImpl) Invalid() bool { return len(my.wrongs) > 0 }

// OK
// fix: 推荐使用 Invalid
func (my *CheckerImpl) OK() bool { return len(my.wrongs) == 0 }

func (my *CheckerImpl) Error() error {
	return operations.NewTernary(operations.TrueFn(func() error { return errors.New(my.ErrorToString("")) })).GetByValue(len(my.wrongs) > 0)
}

func (my *CheckerImpl) ErrorToString(limit string) (ret string) {
	if len(my.wrongs) > 0 {
		ret = anySlices.FillFunc(my.wrongs, func(idx int, value error) string { return value.Error() }).JoinNotEmpty(operations.NewTernary(operations.TrueValue(limit), operations.FalseValue(my.defaultLimit)).GetByValue(limit != ""))
	}

	return
}

func (my *CheckerImpl) Validate(exCheckFns ...ExCheckFunc) Checker {
	fieldInfos := getStructFieldInfos(my.data, "")
	for _, fieldInfo := range fieldInfos {
		if wrongs := fieldInfo.Check().Wrongs(); len(wrongs) > 0 {
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

func WithGin[T any](c *gin.Context, exCheckFns ...ExCheckFunc) (form T, checker Checker) {
	form = *new(T)

	if c == nil || c.Request == nil {
		checker = &CheckerImpl{wrongs: []error{errors.New("gin request is nil")}}
		return
	}

	if strings.TrimSpace(c.GetHeader("Content-Type")) == "" {
		method := strings.ToUpper(strings.TrimSpace(c.Request.Method))
		if method != "GET" && (c.Request.ContentLength > 0 || c.Request.ContentLength == -1) {
			c.Request.Header.Set("Content-Type", "application/json")
		}
	}

	if err := c.ShouldBind(&form); err != nil {
		checker = &CheckerImpl{wrongs: []error{err}}
		return
	}

	return form, OnceValidator().Checker(&form).Validate(exCheckFns...)
}

func WithFiber[T any](c *fiber.Ctx, exCheckFns ...ExCheckFunc) (form T, checker Checker) {
	form = *new(T)

	if err := c.BodyParser(&form); err != nil {
		checker = &CheckerImpl{wrongs: []error{err}}
		return
	}

	return form, OnceValidator().Checker(&form).Validate(exCheckFns...)
}

func callExCheckFn(fn ExCheckFunc, data any) error {
	if fn == nil {
		return fmt.Errorf("callback is nil")
	}
	fv := reflect.ValueOf(fn)
	if fv.Kind() != reflect.Func {
		return fmt.Errorf("callback is not a function: %T", fn)
	}
	ft := fv.Type()
	if ft.NumIn() != 1 || ft.NumOut() < 1 {
		return fmt.Errorf(
			"callback must have signature func(T) error (or similar), got %s",
			ft.String(),
		)
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
			if argType.Kind() == reflect.Pointer && dv.Type().AssignableTo(argType.Elem()) {
				addr := reflect.New(dv.Type())
				addr.Elem().Set(dv)
				dv = addr
			} else if dv.Kind() == reflect.Pointer && dv.Type().Elem().AssignableTo(argType) {
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
	errIface := reflect.TypeFor[error]()
	if !first.Type().Implements(errIface) {
		return fmt.Errorf(
			"callback first return does not implement error: %s",
			first.Type().String(),
		)
	}
	return first.Interface().(error)
}

func getStructFieldInfos(s any, parentName string) []FieldInfo {
	v := reflect.ValueOf(s)
	// 防止 nil 接口或 nil 指针导致 panic
	if !v.IsValid() {
		return nil
	}
	t := v.Type()

	var infos []FieldInfo

	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
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
				// 如果字段是嵌套结构体（如匿名嵌入），即使没有 v-rule 也需要递归检查
				actualKind := fieldValue.Kind()
				actualType := fieldValue.Type()
				if actualKind == reflect.Ptr {
					if fieldValue.IsNil() {
						continue
					}
					actualType = actualType.Elem()
					actualKind = actualType.Kind()
				}
				if actualKind == reflect.Struct &&
					actualType != reflect.TypeOf(time.Time{}) &&
					actualType != reflect.TypeOf(&time.Time{}) {
					infos = append(infos, getStructFieldInfos(fieldValue.Interface(), parentName)...)
				}
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

			// 如果是切片或数组，需要判断是否是基础类型还是 struct 或更深的切片或数组
			switch elemKind {
			case reflect.Slice, reflect.Array:
				elemElemType := elemType.Elem()
				elemElemKind := elemElemType.Kind()

				vRuleTag = strings.TrimLeft(vRuleTag, "(")
				vRuleTag = strings.TrimRight(vRuleTag, ")")

				infos = append(infos, FieldInfo{
					Name:      field.Name,
					Value:     value,
					RefValue:  fieldValue,
					Kind:      elemKind,
					Type:      elemType,
					IsPtr:     isPtr,
					IsNil:     isNil,
					IsZero:    fieldValue.IsZero(),
					VRuleTags: anySlices.NewList(strings.Split(vRuleTag, ")(")),
					VNameTags: anySlices.NewItems(parentName, vNameTag).RemoveEmpty(),
				})

				// 检查数组/切片的元素类型是否是基础类型
				if !isBasicKind(elemElemKind) && !isNil {
					// 非基础类型，需要递归处理
					if fieldValue.Len() > 0 {
						// 遍历数组/切片元素
						for i := 0; i < fieldValue.Len(); i++ {
							elemValue := fieldValue.Index(i)
							// 跳过 nil 指针元素，避免递归时 panic
							if elemValue.Kind() == reflect.Ptr && elemValue.IsNil() {
								continue
							}
							infos = append(
								infos,
								getStructFieldInfos(elemValue.Interface(), vNameTag)...,
							)
						}
					} else {
						// 空数组/切片，使用零值递归
						// 如果元素是指针类型，取其底层类型的零值
						zeroType := elemElemType
						if elemElemKind == reflect.Ptr {
							zeroType = elemElemType.Elem()
						}
						infos = append(
							infos,
							getStructFieldInfos(
								reflect.Zero(zeroType).Interface(),
								vNameTag,
							)...,
						)
					}
				}
			case reflect.Struct:
				if elemType != reflect.TypeOf(time.Time{}) &&
					elemType != reflect.TypeOf(&time.Time{}) { // 如果不是时间类型则递归 struct
					infos = append(
						infos,
						getStructFieldInfos(
							operations.NewTernary(operations.TrueFn(reflect.Zero(elemType).Interface), operations.FalseValue(value)).
								GetByValue(isPtr && isNil),
							vNameTag,
						)...,
					)
				}
			default:
				vRuleTag = strings.TrimLeft(vRuleTag, "(")
				vRuleTag = strings.TrimRight(vRuleTag, ")")

				infos = append(infos, FieldInfo{
					Name:      field.Name,
					Value:     value,
					RefValue:  fieldValue,
					Kind:      elemKind,
					Type:      elemType,
					IsPtr:     isPtr,
					IsNil:     isNil,
					IsZero:    fieldValue.IsZero(),
					VRuleTags: anySlices.NewList(strings.Split(vRuleTag, ")(")),
					VNameTags: anySlices.NewItems(parentName, vNameTag).RemoveEmpty(),
				})
			}
		}
	}

	return infos
}

// isBasicKind 判断 reflect.Kind 是否是基础类型
func isBasicKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}
