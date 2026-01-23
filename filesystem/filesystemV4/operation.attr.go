package filesystemV4

import "os"

type (
	Operation struct {
		Flag int
		Mode os.FileMode
	}

	OperationAttributer interface{ Register(o *Operation) }

	AttrFlag struct{ flag int }
	AttrMode struct{ mode os.FileMode }
)

func NewOperation(attrs ...OperationAttributer) *Operation {
	return new(Operation).SetAttrs(attrs...)
}

func (my *Operation) SetAttrs(attrs ...OperationAttributer) *Operation {
	for idx := range attrs {
		attrs[idx].Register(my)
	}
	return my
}

func Flag(flag int) OperationAttributer   { return AttrFlag{flag: flag} }
func (my AttrFlag) Register(o *Operation) { o.Flag = my.flag }

func Mode(mode os.FileMode) OperationAttributer { return AttrMode{mode: mode} }
func (my AttrMode) Register(o *Operation)       { o.Mode = my.mode }
