package anyMap

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/aid297/aid/array/anySlice"

	jsonIter "github.com/json-iterator/go"
)

type (
	AnyMapper[K comparable, V any] interface {
		SetAttrs(attrs ...Attributer[K, V]) AnyMapper[K, V]
		SetData(data map[K]V) AnyMapper[K, V]
		SetDatum(k K, v V) AnyMapper[K, V]
		SetDataCap(cap int) AnyMapper[K, V]
		SetKeys(keys anySlice.AnySlicer[K]) AnyMapper[K, V]
		AppendKey(k K) AnyMapper[K, V]
		SetValues(values anySlice.AnySlicer[V]) AnyMapper[K, V]
		AppendValue(v V) AnyMapper[K, V]
		Lock() AnyMapper[K, V]
		Unlock() AnyMapper[K, V]
		RLock() AnyMapper[K, V]
		RUnlock() AnyMapper[K, V]
		Copy() AnyMapper[K, V]
		ToString() string
		ToMap() map[K]V
		IsEmpty() bool
		IsNotEmpty() bool
		Has(key K) bool
		SetValue(k K, v V) AnyMapper[K, V]
		GetValueByKey(key K) (V, bool)
		GetValuesByKeys(keys ...K) anySlice.AnySlicer[V]
		GetKeyByValue(value V) (K, bool)
		GetKeysByValues(values ...V) anySlice.AnySlicer[K]
		HasKey(key K) bool
		HasKeys(keys ...K) bool
		HasValue(value V) bool
		HasValues(values ...V) bool
		HasKeyDefault(key K, existFn func(v V) V, notExistFn func() V) AnyMapper[K, V]
		GetKeys() anySlice.AnySlicer[K]
		GetValues() anySlice.AnySlicer[V]
		Length() int
		LengthNotEmpty() int
		Filter(fn func(item V) bool) AnyMapper[K, V]
		RemoveEmpty() AnyMapper[K, V]
		Join(sep string) string
		JoinNotEmpty(sep string) string
		InKey(keys ...K) bool
		NotInKey(keys ...K) bool
		InValue(values ...V) bool
		NotInValue(values ...V) bool
		AllEmpty() bool
		AnyEmpty() bool
		RemoveByKey(key K) AnyMapper[K, V]
		RemoveByKeys(keys ...K) AnyMapper[K, V]
		RemoveByValue(value V) AnyMapper[K, V]
		RemoveByValues(values ...V) AnyMapper[K, V]
		Every(fn func(key K, value V) V) AnyMapper[K, V]
		Each(fn func(key K, value V)) AnyMapper[K, V]
		Clean() AnyMapper[K, V]
		MarshalJSON() ([]byte, error)
		UnmarshalJSON(data []byte) error
	}

	AnyMap[K comparable, V any] struct {
		data   map[K]V
		keys   anySlice.AnySlicer[K]
		values anySlice.AnySlicer[V]
		mu     sync.RWMutex
	}
)

// New 创建一个 AnyMap 实例
func New[K comparable, V any](attrs ...Attributer[K, V]) AnyMapper[K, V] {
	return (&AnyMap[K, V]{mu: sync.RWMutex{}, data: make(map[K]V), keys: anySlice.New[K](), values: anySlice.New[V]()}).SetAttrs(attrs...)
}

// SetAttrs 设置属性
func (my *AnyMap[K, V]) SetAttrs(attrs ...Attributer[K, V]) AnyMapper[K, V] {
	my.mu.Lock()
	defer my.mu.Unlock()
	for i := range attrs {
		attrs[i].Register(my)
	}
	return my
}

// SetData 设置字典键值对
func (my *AnyMap[K, V]) SetData(data map[K]V) AnyMapper[K, V] {
	my.data = data
	return my
}

// SetDatum 设置字典的单个键值对
func (my *AnyMap[K, V]) SetDatum(k K, v V) AnyMapper[K, V] {
	if my.data == nil {
		my.data = make(map[K]V)
	}

	my.data[k] = v
	my.keys = my.keys.Append(k)
	my.values = my.values.Append(v)
	return my
}

// SetDataCap 设置字典数据容量
func (my *AnyMap[K, V]) SetDataCap(cap int) AnyMapper[K, V] {
	my.data = make(map[K]V, cap)
	return my
}

// SetKeys 设置字典的键列表
func (my *AnyMap[K, V]) SetKeys(keys anySlice.AnySlicer[K]) AnyMapper[K, V] {
	my.keys = keys
	return my
}

// AppendKey 向字典的键列表追加一个键
func (my *AnyMap[K, V]) AppendKey(k K) AnyMapper[K, V] {
	my.keys = my.keys.Append(k)
	return my
}

// SetValues 设置字典的值列表
func (my *AnyMap[K, V]) SetValues(values anySlice.AnySlicer[V]) AnyMapper[K, V] {
	my.values = values
	return my
}

// AppendValue 向字典的值列表追加一个值
func (my *AnyMap[K, V]) AppendValue(v V) AnyMapper[K, V] {
	my.values = my.values.Append(v)
	return my
}

func (my *AnyMap[K, V]) Lock() AnyMapper[K, V] {
	my.mu.Lock()
	return my
}

func (my *AnyMap[K, V]) Unlock() AnyMapper[K, V] {
	my.mu.Unlock()
	return my
}

func (my *AnyMap[K, V]) RLock() AnyMapper[K, V] {
	my.mu.RLock()
	return my
}

func (my *AnyMap[K, V]) RUnlock() AnyMapper[K, V] {
	my.mu.RUnlock()
	return my
}

func (my *AnyMap[K, V]) Copy() AnyMapper[K, V] { return New(Map(my.data)) }

func (my *AnyMap[K, V]) ToString() string { return fmt.Sprintf("%v", my.data) }

func (my *AnyMap[K, V]) ToMap() map[K]V { return my.data }

func (my *AnyMap[K, V]) IsEmpty() bool { return len(my.data) == 0 }

func (my *AnyMap[K, V]) IsNotEmpty() bool { return !my.IsEmpty() }

func (my *AnyMap[K, V]) Has(key K) bool {
	_, ok := my.data[key]
	return ok
}

func (my *AnyMap[K, V]) SetValue(k K, v V) AnyMapper[K, V] {
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

func (my *AnyMap[K, V]) GetValueByKey(key K) (V, bool) {
	v, ok := my.data[key]
	return v, ok
}

func (my *AnyMap[K, V]) GetValuesByKeys(keys ...K) anySlice.AnySlicer[V] {
	res := anySlice.New(anySlice.Cap[V](len(keys)))

	for idx := range keys {
		if my.keys.In(keys[idx]) {
			res = res.Append(my.data[keys[idx]])
		}
	}

	return res
}

func (my *AnyMap[K, V]) GetKeyByValue(value V) (K, bool) {
	var k K
	for idx := range my.data {
		if reflect.DeepEqual(value, my.data[idx]) {
			return idx, true
		}
	}
	return k, false
}

func (my *AnyMap[K, V]) GetKeysByValues(values ...V) anySlice.AnySlicer[K] {
	res := anySlice.New(anySlice.Cap[K](len(values)))

	for idx := range values {
		if k, ok := my.GetKeyByValue(values[idx]); ok {
			res = res.Append(k)
		}
	}

	return res
}

func (my *AnyMap[K, V]) HasKey(key K) bool { return my.keys.In(key) }

func (my *AnyMap[K, V]) HasKeys(keys ...K) bool { return my.keys.In(keys...) }

func (my *AnyMap[K, V]) HasValue(value V) bool { return my.values.In(value) }

func (my *AnyMap[K, V]) HasValues(values ...V) bool { return my.values.In(values...) }

func (my *AnyMap[K, V]) HasKeyDefault(key K, existFn func(v V) V, notExistFn func() V) AnyMapper[K, V] {
	if v, e := my.GetValueByKey(key); e {
		return my.SetValue(key, existFn(v))
	}
	return my.SetValue(key, notExistFn())
}

func (my *AnyMap[K, V]) GetKeys() anySlice.AnySlicer[K] { return my.keys }

func (my *AnyMap[K, V]) GetValues() anySlice.AnySlicer[V] { return my.values }

func (my *AnyMap[K, V]) Length() int { return len(my.data) }

func (my *AnyMap[K, V]) LengthNotEmpty() int { return my.RemoveEmpty().Length() }

func (my *AnyMap[K, V]) Filter(fn func(item V) bool) AnyMapper[K, V] {
	res := New(Cap[K, V](my.Length()))

	for idx := range my.values.ToSlice() {
		if fn(my.values.GetValue(idx)) {
			res = res.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}
	}

	my.data = res.ToMap()
	my.keys = res.GetKeys()
	my.values = res.GetValues()

	return my
}

func (my *AnyMap[K, V]) RemoveEmpty() AnyMapper[K, V] {
	return my.Filter(func(item V) bool {
		ref := reflect.ValueOf(item)

		// 处理指针类型：检查是否为 nil 或底层值为零值
		if ref.Kind() == reflect.Ptr {
			return !ref.IsNil() && !ref.Elem().IsZero()
		}

		// 非指针类型：直接检查零值
		if !ref.IsValid() {
			return false
		}
		return !ref.IsZero()
	})
}

func (my *AnyMap[K, V]) Join(sep string) string { return my.values.Join(sep) }

func (my *AnyMap[K, V]) JoinNotEmpty(sep string) string { return my.values.JoinNotEmpty(sep) }

func (my *AnyMap[K, V]) InKey(keys ...K) bool { return my.keys.In(keys...) }

func (my *AnyMap[K, V]) NotInKey(keys ...K) bool { return !my.keys.In(keys...) }

func (my *AnyMap[K, V]) InValue(values ...V) bool { return my.values.In(values...) }

func (my *AnyMap[K, V]) NotInValue(values ...V) bool { return !my.values.In(values...) }

func (my *AnyMap[K, V]) AllEmpty() bool { return my.values.AllEmpty() }

func (my *AnyMap[K, V]) AnyEmpty() bool { return my.values.AnyEmpty() }

func (my *AnyMap[K, V]) RemoveByKey(key K) AnyMapper[K, V] {
	if my.keys.In(key) {
		idx := my.keys.GetIndexByValue(key)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToSlice() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		my.data = newData.ToMap()
		my.keys = newData.GetKeys()
		my.values = newData.GetValues()
	}

	return my
}

func (my *AnyMap[K, V]) RemoveByKeys(keys ...K) AnyMapper[K, V] {
	for idx := range keys {
		my.RemoveByKey(keys[idx])
	}

	return my
}

func (my *AnyMap[K, V]) RemoveByValue(value V) AnyMapper[K, V] {
	if my.values.In(value) {
		idx := my.values.GetIndexByValue(value)
		my.keys = my.keys.RemoveByIndex(idx)
		my.values = my.values.RemoveByIndex(idx)

		newData := New(Cap[K, V](len(my.data) - 1))
		for idx := range my.keys.ToSlice() {
			newData = newData.SetValue(my.keys.GetValue(idx), my.values.GetValue(idx))
		}

		my.data = newData.ToMap()
		my.keys = newData.GetKeys()
		my.values = newData.GetValues()
	}

	return my
}

func (my *AnyMap[K, V]) RemoveByValues(values ...V) AnyMapper[K, V] {
	for idx := range values {
		my.RemoveByValue(values[idx])
	}

	return my
}

func (my *AnyMap[K, V]) Every(fn func(key K, value V) V) AnyMapper[K, V] {
	for idx := range my.keys.ToSlice() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		newV := fn(k, v)
		my.SetValue(k, newV)
	}

	return my
}

func (my *AnyMap[K, V]) Each(fn func(key K, value V)) AnyMapper[K, V] {
	for idx := range my.keys.ToSlice() {
		k := my.keys.GetValue(idx)
		v := my.values.GetValue(idx)
		fn(k, v)
	}

	return my
}

func (my *AnyMap[K, V]) Clean() AnyMapper[K, V] {
	my.keys.Clean()
	my.values.Clean()
	my.data = make(map[K]V)
	return my
}

// MarshalJSON 实现接口：json序列化
func (my *AnyMap[K, V]) MarshalJSON() ([]byte, error) { return jsonIter.Marshal(&my.data) }

// UnmarshalJSON 实现接口：json反序列化
func (my *AnyMap[K, V]) UnmarshalJSON(data []byte) error { return jsonIter.Unmarshal(data, &my.data) }

// Cast 转换所有值并创建新 AsnyMapper
func Cast[K comparable, SRC, DST any](src AnyMapper[K, SRC], fn func(key K, value SRC) DST) AnyMapper[K, DST] {
	d := New[K, DST]()

	for key, value := range src.ToMap() {
		d = d.SetValue(key, fn(key, value))
	}

	return d
}

// Zip 组合键值对为一个新的有序map
func Zip[K comparable, V any](keys []K, values []V) AnyMapper[K, V] {
	d := New[K, V]()

	for idx, key := range keys {
		d = d.SetValue(key, values[idx])
	}

	return d
}

// StructToOther struct 通过 json 转其他格式
func StructToOther[K any, V any](src K) (ret V, err error) {
	var b []byte

	if b, err = jsonIter.Marshal(src); err != nil {
		return
	}

	err = jsonIter.Unmarshal(b, &ret)
	return
}
