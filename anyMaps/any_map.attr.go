package anyMaps

import (
	"github.com/aid297/aid/v2/anySlices"
)

type AnyMapperAttr[K comparable, V any] func(am AnyMapper[K, V])

func Map[K comparable, V any](dict map[K]V) AnyMapperAttr[K, V] {
	return func(am AnyMapper[K, V]) {
		am.SetKeys(anySlices.New(anySlices.Cap[K](len(dict))))
		am.SetValues(anySlices.New(anySlices.Cap[V](len(dict))))
		for idx := range dict {
			am.SetDatum(idx, dict[idx])
		}
	}
}

func Cap[K comparable, V any](cap int) AnyMapperAttr[K, V] {
	return func(am AnyMapper[K, V]) {
		am.SetDataCap(cap)
		am.SetKeys(anySlices.New(anySlices.Cap[K](cap)))
		am.SetValues(anySlices.New(anySlices.Cap[V](cap)))
	}
}
