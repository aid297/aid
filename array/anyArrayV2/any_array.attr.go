package anyArrayV2

type (
	Attributer[T any] interface {
		Register(anyArray *AnyArray[T])
	}

	AttrList[T any]  struct{ list []T }
	AttrItems[T any] struct{ items []T }
	AttrLen[T any]   struct{ length int }
	AttrCap[T any]   struct{ cap int }
	AttrEmpty[T any] struct{}
)

func List[T any](list []T) Attributer[T] { return AttrList[T]{list: list} }

func (my AttrList[T]) Register(anyArray *AnyArray[T]) { anyArray.data = my.list }

func Items[T any](items ...T) Attributer[T] { return AttrItems[T]{items: items} }

func (my AttrItems[T]) Register(anyArray *AnyArray[T]) { anyArray.data = my.items }

func Len[T any](length int) Attributer[T] { return AttrLen[T]{length: length} }

func (my AttrLen[T]) Register(anyArray *AnyArray[T]) { anyArray.data = make([]T, my.length) }

func Cap[T any](cap int) Attributer[T] { return AttrCap[T]{cap: cap} }

func (my AttrCap[T]) Register(anyArray *AnyArray[T]) { anyArray.data = make([]T, 0, my.cap) }

func Empty[T any]() Attributer[T] { return AttrEmpty[T]{} }

func (my AttrEmpty[T]) Register(anyArray *AnyArray[T]) { anyArray.data = make([]T, 0) }
