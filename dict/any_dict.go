package dict

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/aid297/aid/array"
)

type (
	// AnyDict 任意类型字典
	AnyDict[K comparable, V any] struct {
		data   map[K]V
		keys   []K
		values []V
		mu     sync.RWMutex
	}

	// AnyOrderlyItem 可排序的任意类型字典
	AnyOrderlyItem[K comparable, V any] struct {
		Key   K
		Value V
	}
)

// New 根据map创建无序map
func New[K comparable, V any](m map[K]V) *AnyDict[K, V] {
	d := Make[K, V]()

	for key, value := range m {
		d.Set(key, value)
	}

	return d
}

// Make 创建空有序列表
func Make[K comparable, V any]() *AnyDict[K, V] {
	return &AnyDict[K, V]{
		data:   make(map[K]V),
		keys:   make([]K, 0),
		values: make([]V, 0),
		mu:     sync.RWMutex{},
	}
}

// Lock 加锁：写锁
func (my *AnyDict[K, V]) Lock() *AnyDict[K, V] {
	my.mu.Lock()
	return my
}

// Unlock 解锁：写锁
func (my *AnyDict[K, V]) Unlock() *AnyDict[K, V] {
	my.mu.Unlock()
	return my
}

// RLock 加锁：读锁
func (my *AnyDict[K, V]) RLock() *AnyDict[K, V] {
	my.mu.RLock()
	return my
}

// RUnlock 解锁：读锁
func (my *AnyDict[K, V]) RUnlock() *AnyDict[K, V] {
	my.mu.RUnlock()
	return my
}

func (my *AnyDict[K, V]) getKeyByIndex(idx int) K { return my.keys[idx] }

// GetKeyByIndex 通过索引获取键
func (my *AnyDict[K, V]) GetKeyByIndex(idx int) K {
	return my.getKeyByIndex(idx)
}

func (my *AnyDict[K, V]) getKeysByIndexes(indexes ...int) *array.AnyArray[K] {
	keys := make([]K, 0, len(indexes))

	for _, idx := range indexes {
		keys = append(keys, my.getKeyByIndex(idx))
	}

	return array.New(keys)
}

// GetKeysByIndexes 通过索引获取多个键
func (my *AnyDict[K, V]) GetKeysByIndexes(indexes ...int) *array.AnyArray[K] {

	return my.getKeysByIndexes(indexes...)
}

func (my *AnyDict[K, V]) getKeyByValue(value V) K {
	var zero K

	for key, val := range my.data {
		if reflect.DeepEqual(val, value) {
			return key
		}
	}

	return zero
}

// GetKeyByValue 通过值获取键
func (my *AnyDict[K, V]) GetKeyByValue(value V) K {
	return my.getKeyByValue(value)
}

func (my *AnyDict[K, V]) getKeysByValues(values ...V) *array.AnyArray[K] {
	var ret = make([]K, 0)

	for _, value := range values {
		ret = append(ret, my.getKeyByValue(value))
	}

	return array.New(ret)
}

// GetKeysByValues 通过值获取多个键
func (my *AnyDict[K, V]) GetKeysByValues(values ...V) *array.AnyArray[K] {
	return my.getKeysByValues(values...)
}

func (my *AnyDict[K, V]) getValueByIndex(index int) V { return my.values[index] }

// GetValueByIndex 通过索引获取值
func (my *AnyDict[K, V]) GetValueByIndex(index int) V {
	return my.getValueByIndex(index)
}

func (my *AnyDict[K, V]) getValuesByIndexes(indexes ...int) *array.AnyArray[V] {
	values := make([]V, 0, len(indexes))

	for _, idx := range indexes {
		values = append(values, my.values[idx])
	}

	return array.New(values)
}

// GetValsByIndexes 通过索引获取多个值
func (my *AnyDict[K, V]) GetValsByIndexes(indexes ...int) *array.AnyArray[V] {
	return my.GetValuesByIndexes(indexes...)
}

// GetValuesByIndexes 通过索引获取多个值
func (my *AnyDict[K, V]) GetValuesByIndexes(indexes ...int) *array.AnyArray[V] {

	return my.getValuesByIndexes(indexes...)
}

func (my *AnyDict[K, V]) getValueByKey(key K) V { return my.data[key] }

// GetValueByKey 通过键获取值
func (my *AnyDict[K, V]) GetValueByKey(key K) V {
	return my.getValueByKey(key)
}

func (my *AnyDict[K, V]) getValuesByKeys(keys ...K) *array.AnyArray[V] {
	values := make([]V, 0, len(keys))

	for _, key := range keys {
		values = append(values, my.data[key])
	}

	return array.New(values)
}

// GetValuesByKeys 通过值获取多个键
func (my *AnyDict[K, V]) GetValuesByKeys(keys ...K) *array.AnyArray[V] {
	return my.getValuesByKeys(keys...)
}

func (my *AnyDict[K, V]) getIndexByKey(key K) int {
	for i, k := range my.keys {
		if k == key {
			return i
		}
	}

	return -1
}

// GetIndexByKey 通过键获取索引
func (my *AnyDict[K, V]) GetIndexByKey(key K) int {
	return my.getIndexByKey(key)
}

func (my *AnyDict[K, V]) getIndexesByKeys(keys ...K) *array.AnyArray[int] {
	var ret = make([]int, 0)

	for _, key := range keys {
		for idx, k := range my.keys {
			if k == key {
				ret = append(ret, idx)
			}
		}
	}

	return array.New(ret)
}

// GetIndexesByKeys 通过键获取多个索引
func (my *AnyDict[K, V]) GetIndexesByKeys(keys ...K) *array.AnyArray[int] {
	return my.getIndexesByKeys(keys...)
}

func (my *AnyDict[K, V]) getIndexByValue(value V) int {
	for idx, val := range my.values {
		if reflect.DeepEqual(val, value) {
			return idx
		}
	}

	return -1
}

// GetIndexByValue 通过值获取索引
func (my *AnyDict[K, V]) GetIndexByValue(value V) int {
	return my.getIndexByValue(value)
}

func (my *AnyDict[K, V]) getIndexesByValues(values ...V) *array.AnyArray[int] {
	var ret = make([]int, 0)

	for _, value := range values {
		for idx, val := range my.values {
			if reflect.DeepEqual(val, value) {
				ret = append(ret, idx)
			}
		}
	}

	return array.New(ret)
}

// GetIndexesByVals 通过值获取多个索引
func (my *AnyDict[K, V]) GetIndexesByVals(vals ...V) *array.AnyArray[int] {
	return my.GetIndexesByValues(vals...)
}

// GetIndexesByValues 通过值获取多个索引
func (my *AnyDict[K, V]) GetIndexesByValues(values ...V) *array.AnyArray[int] {

	return my.getIndexesByValues(values...)
}

// hasKey 判断是否有键
func (my *AnyDict[K, V]) hasKey(key K) bool {
	return my.getIndexByKey(key) > -1
}

// HasKey 检查键是否存在
func (my *AnyDict[K, V]) HasKey(key K) bool {

	return my.hasKey(key)
}

// HasKeys 检查多个键是否全部存在
func (my *AnyDict[K, V]) HasKeys(keys ...K) bool {

	return my.getIndexesByKeys(keys...).Len() == len(keys)
}

// HasValue 检查值是否存在
func (my *AnyDict[K, V]) HasValue(value V) bool {
	return my.getIndexByValue(value) > -1
}

// HasValues 检查多个值是否全部存在
func (my *AnyDict[K, V]) HasValues(values ...V) bool {
	return my.getIndexesByValues(values...).Len() == len(values)
}

// HasIndex 检查索引是否存在
func (my *AnyDict[K, V]) HasIndex(index int) bool {
	var zero K

	return !reflect.DeepEqual(my.getKeyByIndex(index), zero)
}

// HasIndexes 检查多个索引是否全部存在
func (my *AnyDict[K, V]) HasIndexes(indexes ...int) bool {
	return my.getKeysByIndexes(indexes...).Len() == len(indexes)
}

func (my *AnyDict[K, V]) hasKeyDefault(key K, existFn func(v V), notExistFn func() V) *AnyDict[K, V] {
	if v, e := my.get(key); e {
		existFn(v)
	} else {
		my.set(key, notExistFn())
	}

	return my
}

// HasKeyDefault 判断key是否存在，并执行后续操作
func (my *AnyDict[K, V]) HasKeyDefault(key K, existFn func(v V), notExistFn func() V) *AnyDict[K, V] {

	return my.hasKeyDefault(key, existFn, notExistFn)
}

func (my *AnyDict[K, V]) len() int { return len(my.keys) }

// Len 获取长度
func (my *AnyDict[K, V]) Len() int {

	return my.len()
}

func (my *AnyDict[K, V]) lenWithoutEmpty() int { return my.copy().removeEmpty().len() }

// LenWithoutEmpty 判断非0值的长度
func (my *AnyDict[K, V]) LenWithoutEmpty() int {
	return my.lenWithoutEmpty()
}

func (my *AnyDict[K, V]) isEmpty() bool { return my.len() == 0 }

// IsEmpty 判断是否为空map
func (my *AnyDict[K, V]) IsEmpty() bool { return my.isEmpty() }

func (my *AnyDict[K, V]) set(key K, value V) *AnyDict[K, V] {
	my.data[key] = value
	my.keys = append(my.keys, key)
	my.values = append(my.values, value)

	return my
}

// Set 设置键值对
func (my *AnyDict[K, V]) Set(key K, value V) *AnyDict[K, V] { return my.set(key, value) }

func (my *AnyDict[K, V]) get(key K) (V, bool) {
	v, e := my.data[key]
	return v, e
}

// Get 通过键获取值，同时获取exist
func (my *AnyDict[K, V]) Get(key K) (V, bool) { return my.get(key) }

func (my *AnyDict[K, V]) copy() *AnyDict[K, V] {
	var d = Make[K, V]()

	for idx, key := range my.keys {
		d.set(key, my.values[idx])
	}

	return d
}

// Copy 复制一个新的map
func (my *AnyDict[K, V]) Copy() *AnyDict[K, V] { return my.copy() }

func (my *AnyDict[K, V]) toOrderlyMap() []AnyOrderlyItem[K, V] {
	var items = make([]AnyOrderlyItem[K, V], 0, len(my.keys))

	for idx, key := range my.keys {
		items = append(items, AnyOrderlyItem[K, V]{Key: key, Value: my.values[idx]})
	}

	return items
}

// ToOrderlyMap 导出一个可排序map
func (my *AnyDict[K, V]) ToOrderlyMap() []AnyOrderlyItem[K, V] { return my.toOrderlyMap() }

func (my *AnyDict[K, V]) toMap() map[K]V { return my.data }

// ToMap 获取一个普通map
func (my *AnyDict[K, V]) ToMap() map[K]V { return my.toMap() }

func (my *AnyDict[K, V]) toString() string { return fmt.Sprintf("%v", my.data) }

// ToString 获取字符串
func (my *AnyDict[K, V]) ToString() string { return my.toString() }

func (my *AnyDict[K, V]) getKeys() *array.AnyArray[K] { return array.New(my.keys) }

// GetKeys 获取所有键：*array.AnyArray[K]
func (my *AnyDict[K, V]) GetKeys() *array.AnyArray[K] {

	return my.getKeys()
}

func (my *AnyDict[K, V]) getValues() *array.AnyArray[V] { return array.New(my.values) }

// GetValues 获取所有值：*array.AnyArray[V]
func (my *AnyDict[K, V]) GetValues() *array.AnyArray[V] { return my.getValues() }

func (my *AnyDict[K, V]) getIndexes() *array.AnyArray[int] {
	var ret = make([]int, 0, len(my.keys))

	for i := range my.keys {
		ret = append(ret, i)
	}

	return array.New(ret)
}

// GetIndexes 获取所有索引：*array.AnyArray[int]
func (my *AnyDict[K, V]) GetIndexes() *array.AnyArray[int] { return my.getIndexes() }

func (my *AnyDict[K, V]) firstKey() K { return my.keys[0] }

// FirstKey 获取第一个键
func (my *AnyDict[K, V]) FirstKey() K { return my.firstKey() }

func (my *AnyDict[K, V]) firstValue() V { return my.values[0] }

// FirstValue 获取第一个值
func (my *AnyDict[K, V]) FirstValue() V { return my.firstValue() }

func (my *AnyDict[K, V]) lastKey() K { return my.keys[len(my.keys)-1] }

// LastKey 获取最后一个键
func (my *AnyDict[K, V]) LastKey() K { return my.lastKey() }

func (my *AnyDict[K, V]) lastValue() V { return my.values[len(my.values)-1] }

// LastValue 获取最后一个值
func (my *AnyDict[K, V]) LastValue() V { return my.lastValue() }

func (my *AnyDict[K, V]) filter(fn func(key K, value V) bool) *AnyDict[K, V] {
	var d = Make[K, V]()

	for key, value := range my.data {
		if fn(key, value) {
			d.set(key, value)
		}
	}

	return d
}

// Filter 通过键值对过滤
func (my *AnyDict[K, V]) Filter(fn func(key K, value V) bool) *AnyDict[K, V] {
	d := my.filter(fn)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeByKey(key K) *AnyDict[K, V] {
	var d = Make[K, V]()

	for idx, k := range my.keys {
		if k == key {
			continue
		}

		d.set(k, my.values[idx])
	}

	return d
}

// RemoveByKey 通过键移除元素
func (my *AnyDict[K, V]) RemoveByKey(key K) *AnyDict[K, V] {
	d := my.removeByKey(key)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeByValue(value V) *AnyDict[K, V] {
	var d = Make[K, V]()

	for idx, val := range my.values {
		if reflect.DeepEqual(val, value) {
			continue
		}

		d.set(my.keys[idx], val)
	}

	return d
}

// RemoveByValue 通过值移除元素
func (my *AnyDict[K, V]) RemoveByValue(value V) *AnyDict[K, V] {
	d := my.removeByValue(value)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) removeEmpty() *AnyDict[K, V] {
	d := Make[K, V]()

	for idx, value := range my.values {
		ref := reflect.ValueOf(value)

		if ref.Kind() == reflect.Ptr {
			if ref.IsNil() {
				continue
			}
			if ref.Elem().IsZero() {
				continue
			}
		} else {
			if ref.IsZero() {
				continue
			}
		}

		d.set(my.keys[idx], value)
	}

	return d
}

// RemoveEmpty 移除0值元素
func (my *AnyDict[K, V]) RemoveEmpty() *AnyDict[K, V] {
	d := my.removeEmpty()
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) join(sep string) string {
	var values = make([]string, 0, len(my.values))

	for _, value := range my.values {
		values = append(values, fmt.Sprintf("%v", value))
	}

	return strings.Join(values, sep)
}

// Join 将所有值转为字符串并拼接
func (my *AnyDict[K, V]) Join(seps ...string) string {
	if len(seps) == 0 {
		return my.join(" ")
	} else {
		return my.join(seps[0])
	}
}

func (my *AnyDict[K, V]) joinWithoutEmpty(sep string) string {
	return my.copy().removeEmpty().join(sep)
}

// JoinWithoutEmpty 将去掉0值后转为字符串并拼接
func (my *AnyDict[K, V]) JoinWithoutEmpty(seps ...string) string {
	if len(seps) == 0 {
		return my.joinWithoutEmpty(" ")
	} else {
		return my.joinWithoutEmpty(seps[0])
	}
}

func (my *AnyDict[K, V]) inKeys(keys ...K) bool {
	return my.getIndexesByKeys(keys...).Len() == len(keys)
}

// InKeys 检查多个键是否存在
func (my *AnyDict[K, V]) InKeys(keys ...K) bool { return my.inKeys(keys...) }

// NotInKeys 检查多个键是否不存在
func (my *AnyDict[K, V]) NotInKeys(keys ...K) bool { return !my.inKeys(keys...) }

func (my *AnyDict[K, V]) inValues(values ...V) bool {
	return my.getIndexesByValues(values...).Len() == len(values)
}

// InValues 检查多个值是否存在
func (my *AnyDict[K, V]) InValues(values ...V) bool { return my.inValues(values...) }

// NotInValues 检查多个值是否不存在
func (my *AnyDict[K, V]) NotInValues(values ...V) bool { return !my.inValues(values...) }

func (my *AnyDict[K, V]) every(fn func(key K, value V) (K, V)) *AnyDict[K, V] {
	var d = Make[K, V]()

	for key, value := range my.data {
		k, v := fn(key, value)
		d.set(k, v)
	}

	return d
}

// Every 遍历所有元素并回填
func (my *AnyDict[K, V]) Every(fn func(key K, value V) (K, V)) *AnyDict[K, V] {
	d := my.every(fn)
	my.data = d.data
	my.keys = d.keys
	my.values = d.values
	return my
}

func (my *AnyDict[K, V]) each(fn func(key K, value V)) *AnyDict[K, V] {
	for key, value := range my.data {
		fn(key, value)
	}
	return my
}

// Each 遍历所有元素
func (my *AnyDict[K, V]) Each(fn func(key K, value V)) *AnyDict[K, V] { return my.each(fn) }

func (my *AnyDict[K, V]) clean() *AnyDict[K, V] {
	my.data = make(map[K]V)
	my.keys = make([]K, 0)
	my.values = make([]V, 0)
	return my
}

// Clean 清空
func (my *AnyDict[K, V]) Clean() *AnyDict[K, V] { return my.clean() }

func (my *AnyDict[K, V]) marshalJson() ([]byte, error) { return json.Marshal(&my.data) }

// MarshalJSON json接口实现：序列化
func (my *AnyDict[K, V]) MarshalJSON() ([]byte, error) { return my.marshalJson() }

func (my *AnyDict[K, V]) unmarshalJson(data []byte) error { return json.Unmarshal(data, &my.data) }

// UnmarshalJSON json接口实现：反序列化
func (my *AnyDict[K, V]) UnmarshalJSON(data []byte) error {
	if err := my.unmarshalJson(data); err != nil {
		return err
	}

	my.keys = make([]K, 0)
	my.values = make([]V, 0)
	for key, value := range my.data {
		my.keys = append(my.keys, key)
		my.values = append(my.values, value)
	}

	return nil
}

// Cast 转换所有值并创建新AnyDict
func Cast[K comparable, SRC, DST any](src *AnyDict[K, SRC], fn func(key K, value SRC) DST) *AnyDict[K, DST] {
	var d = Make[K, DST]()

	for key, value := range src.data {
		d.set(key, fn(key, value))
	}

	return d
}

// Zip 组合键值对为一个新的有序map
func Zip[K comparable, V any](keys []K, values []V) *AnyDict[K, V] {
	var d = Make[K, V]()
	for idx, key := range keys {
		d.Set(key, values[idx])
	}

	return d
}
