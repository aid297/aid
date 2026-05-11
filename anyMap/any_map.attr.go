package anyMap

import (
	"github.com/aid297/aid/anySlice"
)

type AnyMapperAttr[K comparable, V any] func(am AnyMapper[K, V])

func Map[K comparable, V any](dict map[K]V) AnyMapperAttr[K, V] {
	return func(am AnyMapper[K, V]) {
		am.SetKeys(anySlice.New(anySlice.Cap[K](len(dict))))
		am.SetValues(anySlice.New(anySlice.Cap[V](len(dict))))
		for idx := range dict {
			am.SetDatum(idx, dict[idx])
		}
	}
}

func Cap[K comparable, V any](cap int) AnyMapperAttr[K, V] {
	return func(am AnyMapper[K, V]) {
		am.SetDataCap(cap)
		am.SetKeys(anySlice.New(anySlice.Cap[K](cap)))
		am.SetValues(anySlice.New(anySlice.Cap[V](cap)))
	}
}
