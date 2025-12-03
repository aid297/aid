package validatorV2

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

type (
	Validator struct {
		Errors          []error
		ValidatorFields []ValidatorField
		prefixNames     []string
		original        any
	}

	ExFunc func(ins *Validator) error

	ValidatorExFuncs struct{ funcs map[string]ExFunc }
)

var (
	validatorExFuncsOnce sync.Once
	validatorExFuncsIns  *ValidatorExFuncs
	defaultExFuncs       = map[string]ExFunc{
		"email": func (ins *Validator) error {
			ins.original.(string)
		}
	}
)

func (ValidatorExFuncs) Register(funcs map[string]ExFunc) *ValidatorExFuncs {
	validatorExFuncsOnce.Do(func() { validatorExFuncsIns = &ValidatorExFuncs{funcs: make(map[string]ExFunc, 0)} })
	for idx := range funcs {
		validatorExFuncsIns.funcs[idx] = funcs[idx]
	}

	return validatorExFuncsIns
}

func (ValidatorExFuncs) Get(key string) ExFunc {
	if _, ok := validatorExFuncsIns.funcs[key]; !ok {
		return nil
	}
	return validatorExFuncsIns.funcs[key]
}

func (Validator) New(original any, prefixNames ...string) Validator {
	var (
		ins = Validator{original: original, prefixNames: make([]string, 0)}
		ref = reflect.ValueOf(original)
	)
	ins.prefixNames = append(ins.prefixNames, prefixNames...)

	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}

	if ref.Kind() != reflect.Struct {
		return Validator{}
	}

	ins.ValidatorFields = make([]ValidatorField, 0, ref.NumField())

	for i := range ref.NumField() {
		field := ref.Type().Field(i)
		val := ref.Field(i)
		if field.Anonymous {
			// 递归验证嵌套字段
			v := Validator{}.New(ref.Field(i).Interface(), ins.prefixNames...)
			if len(v.Errors) > 0 {
				ins.Errors = append(ins.Errors, v.Errors...)
			}
			continue
		}

		vRule := field.Tag.Get("v-rule")
		if vRule == "" || vRule == "-" {
			continue
		}

		vName := field.Tag.Get("v-name")
		if vName == "" || vName == "-" {
			continue
		}

		ins.ValidatorFields = append(ins.ValidatorFields, APP.ValidatorField.New(val, field, strings.Join(append(ins.prefixNames, vName), ".")))
	}

	if len(ins.ValidatorFields) > 0 {
		ins.Errors = make([]error, 0, len(ins.ValidatorFields))
		for idx := range ins.ValidatorFields {
			if ins.ValidatorFields[idx].Error != nil {
				ins.Errors = append(ins.Errors, ins.ValidatorFields[idx].Error)
			}
		}
	}

	return ins
}

func (my Validator) Ex(funcs ...ExFunc) Validator {
	wrongs := make([]error, 0, len(funcs))
	for idx := range funcs {
		if err := funcs[idx](&my); err != nil {
			wrongs = append(wrongs, err)
		}
	}

	my.Errors = append(my.Errors, wrongs...)
	return my
}

func WithGin[T any](val T, c *gin.Context, funcs ...ExFunc) (T, []error) {
	var (
		t   T
		err error
		v   Validator
	)

	if err = c.ShouldBind(&t); err != nil {
		return t, []error{err}
	}

	v = APP.Validator.New(val)

	return t, v.Ex(funcs...).Errors
}

func WithFiber[T any](val T, c *fiber.Ctx, funcs ...ExFunc) (T, []error) {
	var (
		t   T
		err error
		v   Validator
	)

	if err = c.BodyParser(t); err != nil {
		return t, []error{err}
	}

	v = APP.Validator.New(val)

	return t, v.Ex(funcs...).Errors
}
