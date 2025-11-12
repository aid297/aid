package anyDictV2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/aid297/aid/array/anyArrayV2"
)

type AnyDict[K comparable, V any] struct {
	data   map[K]V
	keys   anyArrayV2.AnyArray[K]
	values anyArrayV2.AnyArray[V]
	mu     *sync.RWMutex
}

func New[K comparable, V any](attrs ...Attributer[K, V]) AnyDict[K, V] {
	ins := AnyDict[K, V]{data: make(map[K]V), mu: &sync.RWMutex{}}
	return ins.SetAttrs(attrs...)
}

func (my AnyDict[K, V]) SetAttrs(attrs ...Attributer[K, V]) AnyDict[K, V] {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}
	return my
}

func (my AnyDict[K, V]) Lock() AnyDict[K, V] {
	my.mu.Lock()
	return my
}

func (my AnyDict[K, V]) Unlock() AnyDict[K, V] {
	my.mu.Unlock()
	return my
}

func (my AnyDict[K, V]) RLock() AnyDict[K, V] {
	my.mu.RLock()
	return my
}

func (my AnyDict[K, V]) RUnlock() AnyDict[K, V] {
	my.mu.RUnlock()
	return my
}

func (my AnyDict[K, V]) ToString() string { return fmt.Sprintf("%v", my.data) }

func (my AnyDict[K, V]) ToMap() map[K]V { return my.data }

func (my AnyDict[K, V]) IsEmpty() bool { return len(my.data) == 0 }

func (my AnyDict[K, V]) IsNotEmpty() bool { return !my.IsEmpty() }

func (my AnyDict[K, V]) Has(key K) bool {
	_, ok := my.data[key]
	return ok
}

func (my AnyDict[K, V]) SetValue(k K, v V) AnyDict[K, V] {
	if my.keys.In(k) {
		idx := my.keys.GetIndexByValue(k)
		my.keys = my.keys.SetValue(idx, k)
		my.values = my.values.SetValue(idx, v)
	} else {
		my.keys = my.keys.Append(k)
		my.values = my.values.Append(v)
	}

	my.data[k] = v
	return my
}

func (my AnyDict[K, V]) GetValueByKey(key K) (V, bool) {
	v, ok := my.data[key]
	return v, ok
}

func (my AnyDict[K, V]) GetValuesByKeys(keys ...K) anyArrayV2.AnyArray[V] {
	res := anyArrayV2.New(anyArrayV2.Cap[V](len(keys)))

	for idx := range keys {
		if my.keys.In(keys[idx]) {
			res = res.Append(my.data[keys[idx]])
		}
	}

	return res
}

func (my AnyDict[K, V]) GetKeyByValue(value V) (K, bool) {
	var k K
	for idx := range my.data {
		if reflect.DeepEqual(value, my.data[idx]) {
			return idx, true
		}
	}
	return k, false
}

func (my AnyDict[K, V]) GetKeysByValues(values ...V) anyArrayV2.AnyArray[K] {
	res := anyArrayV2.New(anyArrayV2.Cap[K](len(values)))

	for idx := range values {
		if k, ok := my.GetKeyByValue(values[idx]); ok {
			res = res.Append(k)
		}
	}

	return res
}

func (my AnyDict[K, V]) HasKey(key K) bool { return my.keys.In(key) }

func (my AnyDict[K, V]) HasKeys(keys ...K) bool { return my.keys.In(keys...) }

func (my AnyDict[K, V]) HasValue(value V) bool { return my.values.In(value) }

func (my AnyDict[K, V]) HasValues(values ...V) bool { return my.values.In(values...) }

func (my AnyDict[K, V]) HasKeyDefault(key K, existFn func(v V) V, notExistFn func() V) AnyDict[K, V] {
	if v, e := my.GetValueByKey(key); e {
		my = my.SetValue(key, existFn(v))
	} else {
		my = my.SetValue(key, notExistFn())
	}

	return my
}

func (my AnyDict[K, V]) GetKeys() anyArrayV2.AnyArray[K] { return my.keys }

func (my AnyDict[K, V]) GetValues() anyArrayV2.AnyArray[V] { return my.values }

func (my AnyDict[K, V]) Length() int { return len(my.data) }

func (my AnyDict[K, V]) LengthNotEmpty() int { return my.RemoveEmpty().Length() }

func (my AnyDict[K, V]) Filter(fn func(item V) bool) AnyDict[K, V] {
	res := New(Cap[K, V](my.Length()))

	for idx := range my.values.ToSlice() {
		if fn(my.values.GetValue(idx)) {
			res = res.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}
	}

	return res
}

func (my AnyDict[K, V]) RemoveEmpty() AnyDict[K, V] {
	return my.Filter(func(item V) bool { return !reflect.ValueOf(item).IsZero() })
}

func (my AnyDict[K, V]) Join(sep string) string { return my.values.Join(sep) }

func (my AnyDict[K, V]) JoinNotEmpty(sep string) string { return my.values.JoinNotEmpty(sep) }

func (my AnyDict[K, V]) InKey(keys ...V) bool { return my.values.In(keys...) }

func (my AnyDict[K, V]) NotInKey(keys ...V) bool { return !my.values.In(keys...) }

func (my AnyDict[K, V]) InValue(values ...V) bool { return my.values.In(values...) }

func (my AnyDict[K, V]) NotInValue(values ...V) bool { return !my.values.In(values...) }

func (my AnyDict[K, V]) AllEmpty() bool { return my.values.AllEmpty() }

func (my AnyDict[K, V]) AnyEmpty() bool { return my.values.AnyEmpty() }

func (my AnyDict[K, V]) RemoveByKey(key K) AnyDict[K, V] {
	if my.keys.In(key) {
		idx := my.keys.GetIndexByValue(key)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToSlice() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		return newData
	}

	return my
}

func (my AnyDict[K, V]) RemoveByKeys(keys ...K) AnyDict[K, V] {
	for idx := range keys {
		my = my.RemoveByKey(keys[idx])
	}

	return my
}

func (my AnyDict[K, V]) RemoveByValue(value V) AnyDict[K, V] {
	if my.values.In(value) {
		idx := my.values.GetIndexByValue(value)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToSlice() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		return newData
	}

	return my
}

func (my AnyDict[K, V]) RemoveByValues(values ...V) AnyDict[K, V] {
	for idx := range values {
		my = my.RemoveByValue(values[idx])
	}

	return my
}

func (my AnyDict[K, V]) Every(fn func(key K, value V) V) AnyDict[K, V] {
	for idx := range my.keys.ToSlice() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		newV := fn(k, v)
		my = my.SetValue(k, newV)
	}

	return my
}

func (my AnyDict[K, V]) Each(fn func(key K, value V)) AnyDict[K, V] {
	for idx := range my.keys.ToSlice() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		fn(k, v)
	}

	return my
}

func (my AnyDict[K, V]) Clean() AnyDict[K, V] {
	my.keys.Clean()
	my.values.Clean()
	my.data = make(map[K]V)
	return my
}

// MarshalJSON 实现接口：json序列化
func (my AnyDict[K, V]) MarshalJSON() ([]byte, error) { return json.Marshal(&my.data) }

// UnmarshalJSON 实现接口：json反序列化
func (my AnyDict[K, V]) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &my.data) }

// Cast 转换所有值并创建新AnyDict
func Cast[K comparable, SRC, DST any](src AnyDict[K, SRC], fn func(key K, value SRC) DST) AnyDict[K, DST] {
	var d = New[K, DST]()

	for key, value := range src.data {
		d = d.SetValue(key, fn(key, value))
	}

	return d
}

// Zip 组合键值对为一个新的有序map
func Zip[K comparable, V any](keys []K, values []V) AnyDict[K, V] {
	var d = New[K, V]()
	for idx, key := range keys {
		d = d.SetValue(key, values[idx])
	}

	return d
}
