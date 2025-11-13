package operationV2

import "reflect"

type Match struct {
	targets     []any
	funcs       []func(val any)
	defaultFunc func()
}

func (Match) New(attrs ...MatchAttributer) Match {
	return Match{}.SetAttrs(attrs...)
}

func (my Match) SetAttrs(attrs ...MatchAttributer) Match {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}
	return my
}

func (my Match) SetDefault(fn func()) Match {
	my.defaultFunc = fn
	return my
}

func (my Match) Do(target any) {
	for idx := range my.targets {
		if reflect.DeepEqual(my.targets[idx], target) {
			if fn := my.funcs[idx]; fn != nil {
				fn(target)
			}
			return
		}
	}
	if my.defaultFunc != nil {
		my.defaultFunc()
	}
}
