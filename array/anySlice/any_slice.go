package anySlice

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	jsonIter "github.com/json-iterator/go"
	"github.com/spf13/cast"
)

type (
	AnySlicer[T any] interface {
		SetAttrs(attrs ...Attributer[T]) AnySlicer[T]
		SetData(data []T)
		Lock() AnySlicer[T]
		Unlock() AnySlicer[T]
		RLock() AnySlicer[T]
		RUnlock() AnySlicer[T]
		Empty() bool
		NotEmpty() bool
		IfEmpty(fn func(array AnySlicer[T])) AnySlicer[T]
		IfNotEmpty(fn func(array AnySlicer[T])) AnySlicer[T]
		IfEmptyError(fn func(array AnySlicer[T]) error) error
		IfNotEmptyError(fn func(array AnySlicer[T]) error) error
		Has(k int) bool
		SetValue(k int, v T) AnySlicer[T]
		GetValue(idx int) T
		GetValuePtr(idx int) *T
		GetValueOrDefault(idx int, defaultValue T) T
		GetValues(indexes ...int) []T
		GetValuesBySlices(original, finished int) []T
		Append(v ...T) AnySlicer[T]
		First() T
		Last() T
		ToRaw() []T
		GetIndexes() []int
		GetIndexByValue(value T) int
		GetIndexesByValues(values ...T) []int
		Shuffle() AnySlicer[T]
		Length() int
		LengthNotEmpty() int
		Filter(fn func(item T) bool) AnySlicer[T]
		RemoveEmpty() AnySlicer[T]
		Join(sep string) string
		JoinNotEmpty(sep string) string
		In(targets ...T) bool
		NotIn(targets ...T) bool
		IfIn(fn func(array AnySlicer[T]), targets ...T) AnySlicer[T]
		IfNotIn(fn func(array AnySlicer[T]), targets ...T) AnySlicer[T]
		IfInError(fn func(array AnySlicer[T]) error, targets ...T) error
		IfNotInError(fn func(array AnySlicer[T]) error, targets ...T) error
		AllEmpty() bool
		AnyEmpty() bool
		Chunk(size int) [][]T
		Pluck(fn func(item T) any) AnySlicer[any]
		Intersection(other AnySlicer[T]) AnySlicer[T]
		IntersectionBySlice(other ...T) AnySlicer[T]
		Difference(other AnySlicer[T]) AnySlicer[T]
		DifferenceBySlice(other ...T) AnySlicer[T]
		Union(other AnySlicer[T]) AnySlicer[T]
		UnionBySlice(other []T) AnySlicer[T]
		Unique() AnySlicer[T]
		RemoveByIndex(index int) AnySlicer[T]
		RemoveByIndexes(indexes ...int) AnySlicer[T]
		RemoveByValue(target T) AnySlicer[T]
		RemoveByValues(targets ...T) AnySlicer[T]
		Every(fn func(item T) T) AnySlicer[T]
		Each(fn func(idx int, item T)) AnySlicer[T]
		Sort(fn func(i, j int) bool) AnySlicer[T]
		Clean() AnySlicer[T]
		MarshalJSON() ([]byte, error)
		UnmarshalJSON(data []byte) error
		ToString(formats ...string) string
	}

	AnyArray[T any] struct {
		data []T
		mu   sync.RWMutex
	}
)

func New[T any](attrs ...Attributer[T]) AnySlicer[T] {
	return (&AnyArray[T]{data: make([]T, 0), mu: sync.RWMutex{}}).SetAttrs(attrs...)
}

func NewList[T any](data []T) AnySlicer[T] { return New(List(data)) }

func NewItems[T any](items ...T) AnySlicer[T] { return New(Items(items...)) }

func LoadFn[SRC any, DST any](src []SRC, fn func(value SRC) DST) AnySlicer[DST] {
	var dst = New(Cap[DST](len(src)))

	for idx := range src {
		dst.Append(fn(src[idx]))
	}

	return dst
}

// Cast 转换值类型
func Cast[SRC, DST any](src AnySlicer[SRC], fn func(value SRC) DST) AnySlicer[DST] {
	if src.Length() == 0 {
		return New[DST]()
	}

	data := make([]DST, len(src.ToRaw()))
	for idx := range src.ToRaw() {
		data[idx] = fn(src.ToRaw()[idx])
	}

	return NewList(data)
}

// CastAny 任意类型转目标类型
func CastAny[DST any](src AnySlicer[any], fn func(value any) DST) AnySlicer[DST] {
	if src.Length() == 0 {
		return New[DST]()
	}

	data := make([]DST, len(src.ToRaw()))
	for idx := range src.ToRaw() {
		data[idx] = fn(src.ToRaw()[idx])
	}

	return NewList(data)
}

// ToAny converts any slice to []any
func ToAny(slice any) []any {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]any, v.Len())
	for i := range v.Len() {
		result[i] = v.Index(i).Interface()
	}

	return result
}

// SetAttrs 设置属性
func (my *AnyArray[T]) SetAttrs(attrs ...Attributer[T]) AnySlicer[T] {
	my.mu.Lock()
	defer my.mu.Unlock()

	for i := range attrs {
		attrs[i].Register(my)
	}

	return my
}

// SetData 设置数据
func (my *AnyArray[T]) SetData(data []T) { my.data = data }

// Lock 加锁：写锁
func (my *AnyArray[T]) Lock() AnySlicer[T] {
	my.mu.Lock()
	return my
}

// Unlock 释放：写锁
func (my *AnyArray[T]) Unlock() AnySlicer[T] {
	my.mu.Unlock()
	return my
}

// RLock 加锁：读锁
func (my *AnyArray[T]) RLock() AnySlicer[T] {
	my.mu.RLock()
	return my
}

// RUnlock 释放：读锁
func (my *AnyArray[T]) RUnlock() AnySlicer[T] {
	my.mu.RUnlock()
	return my
}

// Empty 判断是否为空
func (my *AnyArray[T]) Empty() bool { return len(my.data) == 0 }

// NotEmpty 判断是否不为空
func (my *AnyArray[T]) NotEmpty() bool { return !my.Empty() }

// IfEmpty 判断是否为空：如果为空则执行回调
func (my *AnyArray[T]) IfEmpty(fn func(array AnySlicer[T])) AnySlicer[T] {
	if len(my.data) == 0 {
		fn(my)
	}
	return my
}

// IfNotEmpty 判断是否不为空：如果不为空则执行回调
func (my *AnyArray[T]) IfNotEmpty(fn func(array AnySlicer[T])) AnySlicer[T] {
	if len(my.data) != 0 {
		fn(my)
	}
	return my
}

// IfEmptyError 判断是否为空：如果为空则执行回调并返回错误
func (my *AnyArray[T]) IfEmptyError(fn func(array AnySlicer[T]) error) error {
	if len(my.data) == 0 {
		return fn(my)
	}
	return nil
}

// IfNotEmptyError 判断是否不为空：如果不为空则执行回调并返回错误
func (my *AnyArray[T]) IfNotEmptyError(fn func(array AnySlicer[T]) error) error {
	if len(my.data) != 0 {
		return fn(my)
	}
	return nil
}

// Has 检查是否存在
func (my *AnyArray[T]) Has(k int) bool { return k >= 0 && k < len(my.data) }

// SetValue 设置值
func (my *AnyArray[T]) SetValue(k int, v T) AnySlicer[T] {
	my.data[k] = v
	return my
}

// Get 获取值
func (my *AnyArray[T]) GetValue(idx int) T { return my.data[idx] }

// GetValuePtr 获取值指针
func (my *AnyArray[T]) GetValuePtr(idx int) *T {
	if idx < 0 || idx >= len(my.data) {
		return nil
	}
	return &my.data[idx]
}

// GetValueOrDefault 获取值：如果索引不存在则返回默认值
func (my *AnyArray[T]) GetValueOrDefault(idx int, defaultValue T) T {
	if idx < 0 || idx >= len(my.data) {
		return defaultValue
	}
	return my.data[idx]
}

// 获取多个值
func (my *AnyArray[T]) GetValues(indexes ...int) []T {
	res := make([]T, len(indexes))

	for k, idx := range indexes {
		res[k] = my.data[idx]
	}

	return res
}

// GetValuesBySlices 获取切片值
func (my *AnyArray[T]) GetValuesBySlices(original, finished int) []T {
	return my.data[original:finished]
}

// Append 添加值
func (my *AnyArray[T]) Append(v ...T) AnySlicer[T] {
	my.data = append(my.data, v...)
	return my
}

// First 获取第一个值
func (my *AnyArray[T]) First() T {
	var t T
	if len(my.data) > 0 {
		return my.data[0]
	}

	return t
}

// Last 获取最后一个值
func (my *AnyArray[T]) Last() T {
	var t T

	if len(my.data) > 0 {
		return my.data[len(my.data)-1]
	}

	return t
}

// ToRaw 获取全部值：到切片
func (my *AnyArray[T]) ToRaw() []T { return my.data }

// GetIndexes 获取所有索引
func (my *AnyArray[T]) GetIndexes() []int {
	var indexes = make([]int, len(my.data))
	for i := range my.data {
		indexes[i] = i
	}

	return indexes
}

// GetIndexByValue 根据值获取索引下标
func (my *AnyArray[T]) GetIndexByValue(value T) int {
	for idx, val := range my.data {
		if reflect.DeepEqual(val, value) {
			return idx
		}
	}

	return -1
}

// GetIndexesByValues 通过值获取索引下标
func (my *AnyArray[T]) GetIndexesByValues(values ...T) []int {
	var indexes []int
	for _, value := range values {
		for idx, val := range my.data {
			if reflect.DeepEqual(val, value) {
				indexes = append(indexes, idx)
			}
		}
	}

	return indexes
}

// Shuffle 打乱切片中的元素顺序
func (my *AnyArray[T]) Shuffle() AnySlicer[T] {
	randStr := rand.New(rand.NewSource(time.Now().UnixNano()))
	newData := my.data

	for i := range my.data {
		j := randStr.Intn(i + 1)                        // 生成 [0, i] 范围内的随机数
		newData[i], newData[j] = newData[j], newData[i] // 交换元素
	}

	return NewList(newData)
}

// Length 获取数组长度
func (my *AnyArray[T]) Length() int { return len(my.data) }

// LengthNotEmpty 获取非0值长度
func (my *AnyArray[T]) LengthNotEmpty() int { return my.RemoveEmpty().Length() }

// Filter 过滤数组值
func (my *AnyArray[T]) Filter(fn func(item T) bool) AnySlicer[T] {
	j := 0
	ret := make([]T, len(my.data))
	for i := range my.data {
		if fn(my.data[i]) {
			ret[j] = my.data[i]
			j++
		}
	}

	return NewList(ret[:j])
}

// RemoveEmpty 清除0值元素
func (my *AnyArray[T]) RemoveEmpty() AnySlicer[T] {
	return my.Filter(func(item T) bool {
		ref := reflect.ValueOf(item)

		// 快速路径：非指针类型
		if ref.Kind() != reflect.Ptr {
			return !ref.IsZero()
		}

		// 指针类型：检查 nil 和底层零值
		return !ref.IsNil() && !ref.Elem().IsZero()
	})
}

// Join 拼接字符串
func (my *AnyArray[T]) Join(sep string) string {
	values := make([]string, my.Length())
	for idx := range my.data {
		values[idx] = cast.ToString(my.data[idx])
	}

	return strings.Join(values, sep)
}

// JoinNotEmpty 拼接非空字符串
func (my *AnyArray[T]) JoinNotEmpty(sep string) string { return my.RemoveEmpty().Join(sep) }

func (my *AnyArray[T]) in(target T) bool {
	for idx := range my.data {
		if reflect.DeepEqual(target, my.data[idx]) {
			return true
		}
	}

	return false
}

// In 检查值是否存在
func (my *AnyArray[T]) In(targets ...T) bool { return slices.ContainsFunc(targets, my.in) }

// NotIn 检查值是否不存在
func (my *AnyArray[T]) NotIn(targets ...T) bool { return !slices.ContainsFunc(targets, my.in) }

// IfIn 如果In 检查值是否存在：如果存在则执行回调
func (my *AnyArray[T]) IfIn(fn func(array AnySlicer[T]), targets ...T) AnySlicer[T] {
	if my.In(targets...) {
		fn(my)
	}

	return my
}

// IfNotIn 如果NotIn 检查值是否不存在：如果不存在则执行回调
func (my *AnyArray[T]) IfNotIn(fn func(array AnySlicer[T]), targets ...T) AnySlicer[T] {
	if my.NotIn(targets...) {
		fn(my)
	}

	return my
}

func (my *AnyArray[T]) IfInError(fn func(array AnySlicer[T]) error, targets ...T) error {
	if my.In(targets...) {
		return fn(my)
	}
	return nil
}

func (my *AnyArray[T]) IfNotInError(fn func(array AnySlicer[T]) error, targets ...T) error {
	if my.NotIn(targets...) {
		return fn(my)
	}
	return nil
}

// AllEmpty 判断当前数组是否0空
func (my *AnyArray[T]) AllEmpty() bool { return my.RemoveEmpty().Length() == 0 }

// AnyEmpty 判断当前数组中是否存在0值
func (my *AnyArray[T]) AnyEmpty() bool { return my.RemoveEmpty().Length() != my.Length() }

// Chunk 分块
func (my *AnyArray[T]) Chunk(size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(my.data); i += size {
		end := min(i+size, len(my.data))
		chunks = append(chunks, my.data[i:end])
	}

	return chunks
}

// Pluck 获取数组中指定字段的值
func (my *AnyArray[T]) Pluck(fn func(item T) any) AnySlicer[any] {
	var ret = make([]any, 0)
	for _, v := range my.data {
		ret = append(ret, fn(v))
	}

	return NewList(ret)
}

// Intersection 取交集
func (my *AnyArray[T]) Intersection(other AnySlicer[T]) AnySlicer[T] {
	if other.Empty() {
		return New[T]()
	}

	var intersection = make([]T, 0)
	for idx := range my.data {
		if other.In(my.data[idx]) {
			intersection = append(intersection, my.data[idx])
		}
	}

	return NewList(intersection)
}

// IntersectionBySlice 取交集：通过切片
func (my *AnyArray[T]) IntersectionBySlice(other ...T) AnySlicer[T] {
	return my.Intersection(NewList(other))
}

// Difference 取差集
func (my *AnyArray[T]) Difference(other AnySlicer[T]) AnySlicer[T] {
	if other.Empty() {
		return New[T]()
	}

	var difference = make([]T, 0)
	for _, value := range my.data {
		if !other.In(value) {
			difference = append(difference, value)
		}
	}

	return NewList(difference)
}

// DifferenceBySlice 取差集：通过切片
func (my *AnyArray[T]) DifferenceBySlice(other ...T) AnySlicer[T] {
	return my.Difference(NewList(other))
}

// Union 取并集
func (my *AnyArray[T]) Union(other AnySlicer[T]) AnySlicer[T] {
	if other.Empty() {
		return New[T]()
	}

	var union = make([]T, 0)
	union = append(union, my.data...)

	for _, value := range other.ToRaw() {
		if !my.In(value) {
			union = append(union, value)
		}
	}

	return NewList(union)
}

// UnionBySlice 取并集：通过切片
func (my *AnyArray[T]) UnionBySlice(other []T) AnySlicer[T] {
	return my.Union(NewList(other))
}

// Unique 去重
func (my *AnyArray[T]) Unique() AnySlicer[T] {
	seen := make(map[string]struct{}) // 使用空结构体作为值，因为我们只关心键
	result := make([]T, 0)

	for _, value := range my.data {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, value)
		}
	}

	return NewList(result)
}

// RemoveByIndex 根据索引删除元素
func (my *AnyArray[T]) RemoveByIndex(index int) AnySlicer[T] {
	if index < 0 || index >= len(my.data) {
		return my
	}

	return NewList(append(my.data[:index], my.data[index+1:]...))
}

// RemoveByIndexes 根据索引删除元素
func (my *AnyArray[T]) RemoveByIndexes(indexes ...int) AnySlicer[T] {
	newData := make([]T, 0, len(my.data))
	myIndexes := make([]int, 0, len(indexes))

	for idx := range indexes {
		if indexes[idx] < 0 || indexes[idx] >= len(my.data) {
			myIndexes = append(myIndexes, indexes[idx])
		}
	}

	for idx := range my.data {
		for idx2 := range myIndexes {
			if idx == idx2 {
				continue
			}
			newData = append(newData, my.data[idx])
		}
	}

	return NewList(newData)
}

// RemoveByValue 删除数组中对应的目标
func (my *AnyArray[T]) RemoveByValue(target T) AnySlicer[T] {
	var ret = make([]T, len(my.data))
	j := 0
	for _, value := range my.data {
		if !reflect.DeepEqual(value, target) {
			ret[j] = value
			j++
		}
	}

	return New(List(ret[:j]))
}

// RemoveByValues 删除数组中对应的多个目标
func (my *AnyArray[T]) RemoveByValues(targets ...T) AnySlicer[T] {
	data := my.data

	for idx := range targets {
		data = NewList(data).RemoveByValue(targets[idx]).ToRaw()
	}

	return NewList(data)
}

// Every 循环处理每一个
func (my *AnyArray[T]) Every(fn func(item T) T) AnySlicer[T] {
	data := make([]T, len(my.data))

	for idx := range my.data {
		data[idx] = fn(my.data[idx])
	}

	return NewList(data)
}

// Each 遍历数组
func (my *AnyArray[T]) Each(fn func(idx int, item T)) AnySlicer[T] {
	for idx := range my.data {
		fn(idx, my.data[idx])
	}

	return my
}

// Sort 排序
func (my *AnyArray[T]) Sort(fn func(i, j int) bool) AnySlicer[T] {
	sort.Slice(my.data, fn)
	return my
}

// Clean 清理数据
func (my *AnyArray[T]) Clean() AnySlicer[T] {
	my.data = make([]T, 0)
	return my
}

// MarshalJSON 实现接口：json序列化
func (my *AnyArray[T]) MarshalJSON() ([]byte, error) { return jsonIter.Marshal(&my.data) }

// UnmarshalJSON 实现接口：json反序列化
func (my *AnyArray[T]) UnmarshalJSON(data []byte) error { return jsonIter.Unmarshal(data, &my.data) }

// ToString 导出string
func (my *AnyArray[T]) ToString(formats ...string) string {
	var format = "%v"
	if len(formats) > 0 {
		format = formats[0]
	}

	return fmt.Sprintf(format, my.data)
}
