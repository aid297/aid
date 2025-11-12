package operationV2

type (
	TernaryAttributer[T any] interface {
		Register(ternary *Ternary[T])
	}

	AttrTrueValue[T any]  struct{ trueValue T }
	AttrFalseValue[T any] struct{ falseValue T }
	AttrTrueFn[T any]     struct{ trueFn func() T }
	AttrFalseFn[T any]    struct{ falseFn func() T }
)

func TrueValue[T any](trueValue T) AttrTrueValue[T] { return AttrTrueValue[T]{trueValue: trueValue} }

func (my AttrTrueValue[T]) Register(ternary *Ternary[T]) {
	ternary.trueFn = func() T { return my.trueValue }
}

func FalseValue[T any](falseValue T) AttrFalseValue[T] {
	return AttrFalseValue[T]{falseValue: falseValue}
}

func (my AttrFalseValue[T]) Register(ternary *Ternary[T]) {
	ternary.falseFn = func() T { return my.falseValue }
}

func TrueFn[T any](trueFn func() T) AttrTrueFn[T] { return AttrTrueFn[T]{trueFn: trueFn} }

func (my AttrTrueFn[T]) Register(ternary *Ternary[T]) { ternary.trueFn = my.trueFn }

func FalseFn[T any](falseFn func() T) AttrFalseFn[T] { return AttrFalseFn[T]{falseFn: falseFn} }

func (my AttrFalseFn[T]) Register(ternary *Ternary[T]) { ternary.falseFn = my.falseFn }
