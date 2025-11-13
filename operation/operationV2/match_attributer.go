package operationV2

type (
	MatchAttributer interface{ Register(match *Match) }

	AttrMatchItem struct {
		target any
		fn     func(val any)
	}
)

func MatchItem(target any, fn func(val any)) AttrMatchItem {
	return AttrMatchItem{target: target, fn: fn}
}
func (my AttrMatchItem) Register(match *Match) {
	match.targets = append(match.targets, my.target)
	match.funcs = append(match.funcs, my.fn)
}
