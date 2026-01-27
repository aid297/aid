package anyDictV3

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/aid297/aid/array/anySlice"

	jsonIter "github.com/json-iterator/go"
)

type (
	IAnyDict[K comparable, V any] interface {
		SetAttrs(attrs ...Attributer[K, V]) IAnyDict[K, V]
		SetData(data map[K]V) IAnyDict[K, V]
		SetDatum(k K, v V) IAnyDict[K, V]
		SetDataCap(cap int) IAnyDict[K, V]
		SetKeys(keys anySlice.AnySlicer[K]) IAnyDict[K, V]
		AppendKey(k K) IAnyDict[K, V]
		SetValues(values anySlice.AnySlicer[V]) IAnyDict[K, V]
		AppendValue(v V) IAnyDict[K, V]
		Lock() IAnyDict[K, V]
		Unlock() IAnyDict[K, V]
		RLock() IAnyDict[K, V]
		RUnlock() IAnyDict[K, V]
		ToString() string
		ToMap() map[K]V
		IsEmpty() bool
		IsNotEmpty() bool
		Has(key K) bool
		SetValue(k K, v V) IAnyDict[K, V]
		GetValueByKey(key K) (V, bool)
		GetValuesByKeys(keys ...K) anySlice.AnySlicer[V]
		GetKeyByValue(value V) (K, bool)
		GetKeysByValues(values ...V) anySlice.AnySlicer[K]
		HasKey(key K) bool
		HasKeys(keys ...K) bool
		HasValue(value V) bool
		HasValues(values ...V) bool
		HasKeyDefault(key K, existFn func(v V) V, notExistFn func() V) IAnyDict[K, V]
		GetKeys() anySlice.AnySlicer[K]
		GetValues() anySlice.AnySlicer[V]
		Length() int
		LengthNotEmpty() int
		Filter(fn func(item V) bool) IAnyDict[K, V]
		RemoveEmpty() IAnyDict[K, V]
		Join(sep string) string
		JoinNotEmpty(sep string) string
		InKey(keys ...K) bool
		NotInKey(keys ...K) bool
		InValue(values ...V) bool
		NotInValue(values ...V) bool
		AllEmpty() bool
		AnyEmpty() bool
		RemoveByKey(key K) IAnyDict[K, V]
		RemoveByKeys(keys ...K) IAnyDict[K, V]
		RemoveByValue(value V) IAnyDict[K, V]
		RemoveByValues(values ...V) IAnyDict[K, V]
		Every(fn func(key K, value V) V) IAnyDict[K, V]
		Each(fn func(key K, value V)) IAnyDict[K, V]
		Clean() IAnyDict[K, V]
		MarshalJSON() ([]byte, error)
		UnmarshalJSON(data []byte) error
	}

	AnyDict[K comparable, V any] struct {
		data   map[K]V
		keys   anySlice.AnySlicer[K]
		values anySlice.AnySlicer[V]
		mu     sync.RWMutex
	}
)

// New 创建一个 AnyDict 实例
func New[K comparable, V any](attrs ...Attributer[K, V]) IAnyDict[K, V] {
	return (&AnyDict[K, V]{mu: sync.RWMutex{}, data: make(map[K]V), keys: anySlice.New[K](), values: anySlice.New[V]()}).SetAttrs(attrs...)
}

// SetAttrs 设置属性
func (my *AnyDict[K, V]) SetAttrs(attrs ...Attributer[K, V]) IAnyDict[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

// SetData 设置字典键值对
func (my *AnyDict[K, V]) SetData(data map[K]V) IAnyDict[K, V] {
	my.data = data
	return my
}

// SetDatum 设置字典的单个键值对
func (my *AnyDict[K, V]) SetDatum(k K, v V) IAnyDict[K, V] {
	if my.data == nil {
		my.data = make(map[K]V)
	}

	my.data[k] = v
	my.keys = my.keys.Append(k)
	my.values = my.values.Append(v)
	return my
}

// SetDataCap 设置字典数据容量
func (my *AnyDict[K, V]) SetDataCap(cap int) IAnyDict[K, V] {
	my.data = make(map[K]V, cap)
	return my
}

// SetKeys 设置字典的键列表
func (my *AnyDict[K, V]) SetKeys(keys anySlice.AnySlicer[K]) IAnyDict[K, V] {
	my.keys = keys
	return my
}

// AppendKey 向字典的键列表追加一个键
func (my *AnyDict[K, V]) AppendKey(k K) IAnyDict[K, V] {
	my.keys = my.keys.Append(k)
	return my
}

// SetValues 设置字典的值列表
func (my *AnyDict[K, V]) SetValues(values anySlice.AnySlicer[V]) IAnyDict[K, V] {
	my.values = values
	return my
}

// AppendValue 向字典的值列表追加一个值
func (my *AnyDict[K, V]) AppendValue(v V) IAnyDict[K, V] {
	my.values = my.values.Append(v)
	return my
}

func (my *AnyDict[K, V]) Lock() IAnyDict[K, V] {
	my.mu.Lock()
	return my
}

func (my *AnyDict[K, V]) Unlock() IAnyDict[K, V] {
	my.mu.Unlock()
	return my
}

func (my *AnyDict[K, V]) RLock() IAnyDict[K, V] {
	my.mu.RLock()
	return my
}

func (my *AnyDict[K, V]) RUnlock() IAnyDict[K, V] {
	my.mu.RUnlock()
	return my
}

func (my *AnyDict[K, V]) ToString() string { return fmt.Sprintf("%v", my.data) }

func (my *AnyDict[K, V]) ToMap() map[K]V { return my.data }

func (my *AnyDict[K, V]) IsEmpty() bool { return len(my.data) == 0 }

func (my *AnyDict[K, V]) IsNotEmpty() bool { return !my.IsEmpty() }

func (my *AnyDict[K, V]) Has(key K) bool {
	_, ok := my.data[key]
	return ok
}

func (my *AnyDict[K, V]) SetValue(k K, v V) IAnyDict[K, V] {
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

func (my *AnyDict[K, V]) GetValueByKey(key K) (V, bool) {
	v, ok := my.data[key]
	return v, ok
}

func (my *AnyDict[K, V]) GetValuesByKeys(keys ...K) anySlice.AnySlicer[V] {
	res := anySlice.New(anySlice.Cap[V](len(keys)))

	for idx := range keys {
		if my.keys.In(keys[idx]) {
			res = res.Append(my.data[keys[idx]])
		}
	}

	return res
}

func (my *AnyDict[K, V]) GetKeyByValue(value V) (K, bool) {
	var k K
	for idx := range my.data {
		if reflect.DeepEqual(value, my.data[idx]) {
			return idx, true
		}
	}
	return k, false
}

func (my *AnyDict[K, V]) GetKeysByValues(values ...V) anySlice.AnySlicer[K] {
	res := anySlice.New(anySlice.Cap[K](len(values)))

	for idx := range values {
		if k, ok := my.GetKeyByValue(values[idx]); ok {
			res = res.Append(k)
		}
	}

	return res
}

func (my *AnyDict[K, V]) HasKey(key K) bool { return my.keys.In(key) }

func (my *AnyDict[K, V]) HasKeys(keys ...K) bool { return my.keys.In(keys...) }

func (my *AnyDict[K, V]) HasValue(value V) bool { return my.values.In(value) }

func (my *AnyDict[K, V]) HasValues(values ...V) bool { return my.values.In(values...) }

func (my *AnyDict[K, V]) HasKeyDefault(key K, existFn func(v V) V, notExistFn func() V) IAnyDict[K, V] {
	if v, e := my.GetValueByKey(key); e {
		return my.SetValue(key, existFn(v))
	}
	return my.SetValue(key, notExistFn())
}

func (my *AnyDict[K, V]) GetKeys() anySlice.AnySlicer[K] { return my.keys }

func (my *AnyDict[K, V]) GetValues() anySlice.AnySlicer[V] { return my.values }

func (my *AnyDict[K, V]) Length() int { return len(my.data) }

func (my *AnyDict[K, V]) LengthNotEmpty() int { return my.RemoveEmpty().Length() }

func (my *AnyDict[K, V]) Filter(fn func(item V) bool) IAnyDict[K, V] {
	res := New(Cap[K, V](my.Length()))

	for idx := range my.values.ToRaw() {
		if fn(my.values.GetValue(idx)) {
			res = res.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}
	}

	return res
}

func (my *AnyDict[K, V]) RemoveEmpty() IAnyDict[K, V] {
	return my.Filter(func(item V) bool {
		ref := reflect.ValueOf(item)

		// 处理指针类型：检查是否为 nil 或底层值为零值
		if ref.Kind() == reflect.Ptr {
			return !ref.IsNil() && !ref.Elem().IsZero()
		}

		// 非指针类型：直接检查零值
		return !ref.IsZero()
	})
}

func (my *AnyDict[K, V]) Join(sep string) string { return my.values.Join(sep) }

func (my *AnyDict[K, V]) JoinNotEmpty(sep string) string { return my.values.JoinNotEmpty(sep) }

func (my *AnyDict[K, V]) InKey(keys ...K) bool { return my.keys.In(keys...) }

func (my *AnyDict[K, V]) NotInKey(keys ...K) bool { return !my.keys.In(keys...) }

func (my *AnyDict[K, V]) InValue(values ...V) bool { return my.values.In(values...) }

func (my *AnyDict[K, V]) NotInValue(values ...V) bool { return !my.values.In(values...) }

func (my *AnyDict[K, V]) AllEmpty() bool { return my.values.AllEmpty() }

func (my *AnyDict[K, V]) AnyEmpty() bool { return my.values.AnyEmpty() }

func (my *AnyDict[K, V]) RemoveByKey(key K) IAnyDict[K, V] {
	if my.keys.In(key) {
		idx := my.keys.GetIndexByValue(key)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToRaw() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		return newData
	}

	return my
}

func (my *AnyDict[K, V]) RemoveByKeys(keys ...K) IAnyDict[K, V] {
	for idx := range keys {
		my.RemoveByKey(keys[idx])
	}

	return my
}

func (my *AnyDict[K, V]) RemoveByValue(value V) IAnyDict[K, V] {
	if my.values.In(value) {
		idx := my.values.GetIndexByValue(value)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToRaw() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		return newData
	}

	return my
}

func (my *AnyDict[K, V]) RemoveByValues(values ...V) IAnyDict[K, V] {
	for idx := range values {
		my.RemoveByValue(values[idx])
	}

	return my
}

func (my *AnyDict[K, V]) Every(fn func(key K, value V) V) IAnyDict[K, V] {
	for idx := range my.keys.ToRaw() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		newV := fn(k, v)
		my.SetValue(k, newV)
	}

	return my
}

func (my *AnyDict[K, V]) Each(fn func(key K, value V)) IAnyDict[K, V] {
	for idx := range my.keys.ToRaw() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		fn(k, v)
	}

	return my
}

func (my *AnyDict[K, V]) Clean() IAnyDict[K, V] {
	my.keys.Clean()
	my.values.Clean()
	my.data = make(map[K]V)
	return my
}

// MarshalJSON 实现接口：json序列化
func (my *AnyDict[K, V]) MarshalJSON() ([]byte, error) { return jsonIter.Marshal(&my.data) }

// UnmarshalJSON 实现接口：json反序列化
func (my *AnyDict[K, V]) UnmarshalJSON(data []byte) error { return jsonIter.Unmarshal(data, &my.data) }

// Cast 转换所有值并创建新AnyDict
func Cast[K comparable, SRC, DST any](src AnyDict[K, SRC], fn func(key K, value SRC) DST) IAnyDict[K, DST] {
	d := New[K, DST]()

	for key, value := range src.data {
		d = d.SetValue(key, fn(key, value))
	}

	return d
}

// Zip 组合键值对为一个新的有序map
func Zip[K comparable, V any](keys []K, values []V) IAnyDict[K, V] {
	d := New[K, V]()

	for idx, key := range keys {
		d = d.SetValue(key, values[idx])
	}

	return d
}
