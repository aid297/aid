package anyDictV3

import (
	"github.com/aid297/aid/array/anyArrayV3"
)

type (
	Attributer[K comparable, V any] interface{ Register(anyDict AnyDicter[K, V]) }

	AttrMap[K comparable, V any] struct{ dict map[K]V }
	AttrCap[K comparable, V any] struct{ cap int }
)

func Map[K comparable, V any](dict map[K]V) AttrMap[K, V] { return AttrMap[K, V]{dict: dict} }
func (my AttrMap[K, V]) Register(anyDict AnyDicter[K, V]) {
	anyDict.SetKeys(anyArrayV3.New(anyArrayV3.Cap[K](len(my.dict))))
	anyDict.SetValues(anyArrayV3.New(anyArrayV3.Cap[V](len(my.dict))))
	for idx := range my.dict {
		anyDict.SetDatum(idx, my.dict[idx])
		anyDict.AppendKey(idx)
		anyDict.AppendValue(my.dict[idx])
	}
}

func Cap[K comparable, V any](cap int) AttrCap[K, V] { return AttrCap[K, V]{cap: cap} }
func (my AttrCap[K, V]) Register(anyDict AnyDicter[K, V]) {
	anyDict.SetDataCap(my.cap)
	anyDict.SetKeys(anyArrayV3.New(anyArrayV3.Cap[K](my.cap)))
	anyDict.SetValues(anyArrayV3.New(anyArrayV3.Cap[V](my.cap)))
}
