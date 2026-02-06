package anyMap

import (
	"github.com/aid297/aid/array/anySlice"
)

type (
	Attributer[K comparable, V any] interface{ Register(anyDict AnyMapper[K, V]) }

	AttrMap[K comparable, V any] struct{ dict map[K]V }
	AttrCap[K comparable, V any] struct{ cap int }
)

func Map[K comparable, V any](dict map[K]V) AttrMap[K, V] { return AttrMap[K, V]{dict: dict} }
func (my AttrMap[K, V]) Register(anyDict AnyMapper[K, V]) {
	anyDict.SetKeys(anySlice.New(anySlice.Cap[K](len(my.dict))))
	anyDict.SetValues(anySlice.New(anySlice.Cap[V](len(my.dict))))
	for idx := range my.dict {
		anyDict.SetDatum(idx, my.dict[idx])
	}
}

func Cap[K comparable, V any](cap int) AttrCap[K, V] { return AttrCap[K, V]{cap: cap} }
func (my AttrCap[K, V]) Register(anyDict AnyMapper[K, V]) {
	anyDict.SetDataCap(my.cap)
	anyDict.SetKeys(anySlice.New(anySlice.Cap[K](my.cap)))
	anyDict.SetValues(anySlice.New(anySlice.Cap[V](my.cap)))
}
