package array

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	AnyArray[T any] struct {
		data []T
		mu   sync.RWMutex
	}
)

// New 实例化
func New[T any](list []T) *AnyArray[T] { return &AnyArray[T]{data: list, mu: sync.RWMutex{}} }

// NewDestruction 通过解构参数实例化
func NewDestruction[T any](list ...T) *AnyArray[T] {
	return &AnyArray[T]{data: list, mu: sync.RWMutex{}}
}

// Make 初始化
func Make[T any](size int) *AnyArray[T] {
	return &AnyArray[T]{data: make([]T, size), mu: sync.RWMutex{}}
}

// Lock 加锁：写锁
func (my *AnyArray[T]) Lock() *AnyArray[T] {
	my.mu.Lock()
	return my
}

// Unlock 解锁：写锁
func (my *AnyArray[T]) Unlock() *AnyArray[T] {
	my.mu.Unlock()
	return my
}

// RLock 加锁：读锁
func (my *AnyArray[T]) RLock() *AnyArray[T] {
	my.mu.RLock()
	return my
}

// RUnlock 解锁：读锁
func (my *AnyArray[T]) RUnlock() *AnyArray[T] {
	my.mu.RUnlock()
	return my
}

// isEmpty 判断是否为空
func (my *AnyArray[T]) isEmpty() bool { return len(my.data) == 0 }

// IsEmpty 判断是否为空
func (my *AnyArray[T]) IsEmpty() bool { return my.isEmpty() }

// IsNotEmpty 判断是否不为空
func (my *AnyArray[T]) IsNotEmpty() bool { return !my.isEmpty() }

// Has 检查是否存在
func (my *AnyArray[T]) Has(k int) bool { return k >= 0 && k < len(my.data) }

// Set 设置值
func (my *AnyArray[T]) Set(k int, v T) *AnyArray[T] {
	my.data[k] = v
	return my
}

// Get 获取值
func (my *AnyArray[T]) Get(idx int) T { return my.data[idx] }

// GetPtr 获取值指针
func (my *AnyArray[T]) GetPtr(idx int) *T {
	if idx < 0 || idx >= len(my.data) {
		return nil
	}
	return &my.data[idx]
}

// GetOrDefault 获取值：如果索引不存在则返回默认值
func (my *AnyArray[T]) GetOrDefault(idx int, defaultValue T) T {
	if idx < 0 || idx >= len(my.data) {
		return defaultValue
	}
	return my.data[idx]
}

// GetByIndexes 通过多索引获取内容
func (my *AnyArray[T]) GetByIndexes(indexes ...int) *AnyArray[T] {
	res := make([]T, len(indexes))

	for k, idx := range indexes {
		res[k] = my.data[idx]
	}

	return New(res)
}

// Append 追加
func (my *AnyArray[T]) Append(v ...T) *AnyArray[T] {
	my.data = append(my.data, v...)
	return my
}

// First 获取第一个值
func (my *AnyArray[T]) First() T { return my.data[0] }

// Last 获取最后一个值
func (my *AnyArray[T]) Last() T {
	return operationV2.NewTernary(operationV2.TrueFn(func() T { return my.data[len(my.data)-1] }), operationV2.FalseValue(operationV2.NewTernary(operationV2.TrueFn(func() T { return my.data[0] })).GetByValue(my.Len() == 0))).GetByValue(my.Len() > 1)
}

// ToSlice 获取全部值：到切片
func (my *AnyArray[T]) ToSlice() []T {
	var ret = make([]T, len(my.data))
	copy(ret, my.data)

	return ret
}

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
func (my *AnyArray[T]) GetIndexesByValues(values ...T) *AnyArray[int] {
	var indexes []int
	for _, value := range values {
		for idx, val := range my.data {
			if reflect.DeepEqual(val, value) {
				indexes = append(indexes, idx)
			}
		}
	}

	return New(indexes)
}

// Copy 复制自己
func (my *AnyArray[T]) Copy() *AnyArray[T] {
	return New(my.data)
}

// Shuffle 打乱切片中的元素顺序
func (my *AnyArray[T]) Shuffle() *AnyArray[T] {
	randStr := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range my.data {
		j := randStr.Intn(i + 1)                        // 生成 [0, i] 范围内的随机数
		my.data[i], my.data[j] = my.data[j], my.data[i] // 交换元素
	}

	return my
}

// Len 获取数组长度
func (my *AnyArray[T]) Len() int { return len(my.data) }

// LenWithoutEmpty 获取非0值长度
func (my *AnyArray[T]) LenWithoutEmpty() int { return my.Copy().RemoveEmpty().Len() }

// Filter 过滤数组值
func (my *AnyArray[T]) Filter(fn func(item T) bool) *AnyArray[T] {
	j := 0
	ret := make([]T, len(my.data))
	for i := range my.data {
		if fn(my.data[i]) {
			ret[j] = my.data[i]
			j++
		}
	}

	my.data = ret[:j]

	return my
}

// RemoveEmpty 清除0值元素
func (my *AnyArray[T]) RemoveEmpty() *AnyArray[T] {
	var data = make([]T, 0)

	for _, item := range my.data {
		ref := reflect.ValueOf(item)

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

		data = append(data, item)
	}

	return New(data)
}

// Join 拼接字符串
func (my *AnyArray[T]) Join(sep string) string {
	values := make([]string, my.Len())
	for idx, datum := range my.data {
		values[idx] = fmt.Sprintf("%v", datum)
	}
	return strings.Join(values, sep)
}

// JoinWithoutEmpty 拼接非空元素
func (my *AnyArray[T]) JoinWithoutEmpty(seps ...string) string {
	var sep = " "
	if len(seps) > 0 {
		sep = seps[0]
	}

	values := make([]string, my.Copy().RemoveEmpty().Len())
	j := 0
	for _, datum := range my.Copy().RemoveEmpty().ToSlice() {
		values[j] = fmt.Sprintf("%v", datum)
		j++
	}

	return strings.Join(values, sep)
}

func (my *AnyArray[T]) in(target T) bool {
	for _, element := range my.data {
		if reflect.DeepEqual(target, element) {
			return true
		}
	}

	return false
}

// In 检查值是否存在
func (my *AnyArray[T]) In(targets ...T) bool {
	return slices.ContainsFunc(targets, my.in)
}

// NotIn 检查值是否不存在
func (my *AnyArray[T]) NotIn(targets ...T) bool {

	return !slices.ContainsFunc(targets, my.in)
}

// AllEmpty 判断当前数组是否0空
func (my *AnyArray[T]) AllEmpty() bool { return my.Copy().RemoveEmpty().Len() == 0 }

// AnyEmpty 判断当前数组中是否存在0值
func (my *AnyArray[T]) AnyEmpty() bool { return my.Copy().RemoveEmpty().Len() != len(my.data) }

// Chunk 分块
func (my *AnyArray[T]) Chunk(size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(my.data); i += size {
		end := i + size
		if end > len(my.data) {
			end = len(my.data)
		}
		chunks = append(chunks, my.data[i:end])
	}

	return chunks
}

// Pluck 获取数组中指定字段的值
func (my *AnyArray[T]) Pluck(fn func(item T) any) *AnyArray[any] {
	var ret = make([]any, 0)
	for _, v := range my.data {
		ret = append(ret, fn(v))
	}

	return New(ret)
}

func (my *AnyArray[T]) intersection(other *AnyArray[T]) {
	if other == nil || other.IsEmpty() {
		return
	}

	var intersection = make([]T, 0)
	for _, value := range my.data {
		if other.In(value) {
			intersection = append(intersection, value)
		}
	}

	my.data = intersection
}

// Intersection 取交集
func (my *AnyArray[T]) Intersection(other *AnyArray[T]) *AnyArray[T] {
	my.intersection(other)
	return my
}

// IntersectionBySlice 取交集：通过切片
func (my *AnyArray[T]) IntersectionBySlice(other []T) *AnyArray[T] {
	my.intersection(New(other))
	return my
}

// IntersectionByValues 取交集：通过值
func (my *AnyArray[T]) IntersectionByValues(values ...T) *AnyArray[T] {
	my.intersection(NewDestruction(values...))
	return my
}

func (my *AnyArray[T]) difference(other *AnyArray[T]) {
	if other == nil || other.IsEmpty() {
		return
	}

	var difference = make([]T, 0)
	for _, value := range my.data {
		if !other.In(value) {
			difference = append(difference, value)
		}
	}

	my.data = difference
}

// Difference 取差集
func (my *AnyArray[T]) Difference(other *AnyArray[T]) *AnyArray[T] {
	my.difference(other)
	return my
}

// DifferenceBySlice 取差集：通过切片
func (my *AnyArray[T]) DifferenceBySlice(other []T) *AnyArray[T] {
	my.difference(New(other))
	return my
}

// DifferenceByValues 取差集：通过值
func (my *AnyArray[T]) DifferenceByValues(values ...T) *AnyArray[T] {
	my.difference(NewDestruction(values...))
	return my
}

func (my *AnyArray[T]) union(other *AnyArray[T]) {
	if other == nil || other.IsEmpty() {
		return
	}

	var union = make([]T, 0)
	union = append(union, my.data...)

	for _, value := range other.data {
		if !my.In(value) {
			union = append(union, value)
		}
	}

	my.data = union
}

// Union 取并集
func (my *AnyArray[T]) Union(other *AnyArray[T]) *AnyArray[T] {
	my.union(other)
	return my
}

// UnionBySlice 取并集：通过切片
func (my *AnyArray[T]) UnionBySlice(other []T) *AnyArray[T] {
	my.union(New(other))
	return my
}

// UnionByValues 取并集：通过值
func (my *AnyArray[T]) UnionByValues(values ...T) *AnyArray[T] {
	my.union(NewDestruction(values...))
	return my
}

// Unique 去重
func (my *AnyArray[T]) Unique() *AnyArray[T] {
	seen := make(map[string]struct{}) // 使用空结构体作为值，因为我们只关心键
	result := make([]T, 0)

	for _, value := range my.data {
		key := fmt.Sprintf("%v", value)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, value)
		}
	}

	my.data = result
	return my
}

// RemoveByIndex 根据索引删除元素
func (my *AnyArray[T]) RemoveByIndex(index int) *AnyArray[T] {
	if index < 0 || index >= len(my.data) {
		return my
	}

	my.data = append(my.data[:index], my.data[index+1:]...)
	return my
}

// RemoveByIndexes 根据索引删除元素
func (my *AnyArray[T]) RemoveByIndexes(indexes ...int) *AnyArray[T] {
	for _, index := range indexes {
		my.RemoveByIndex(index)
	}
	return my
}

// RemoveByValue 删除数组中对应的目标
func (my *AnyArray[T]) RemoveByValue(target T) *AnyArray[T] {
	var ret = make([]T, len(my.data))
	j := 0
	for _, value := range my.data {
		if !reflect.DeepEqual(value, target) {
			ret[j] = value
			j++
		}
	}
	my.data = ret[:j]

	return my
}

// RemoveByValues 删除数组中对应的多个目标
func (my *AnyArray[T]) RemoveByValues(targets ...T) *AnyArray[T] {
	for _, target := range targets {
		my.RemoveByValue(target)
	}

	return my
}

// Every 循环处理每一个
func (my *AnyArray[T]) Every(fn func(item T) T) *AnyArray[T] {
	for idx := range my.data {
		v := fn(my.data[idx])
		my.data[idx] = v
	}

	return my
}

// Each 遍历数组
func (my *AnyArray[T]) Each(fn func(idx int, item T)) *AnyArray[T] {
	for idx := range my.data {
		fn(idx, my.data[idx])
	}

	return my
}

// Sort 排序
func (my *AnyArray[T]) Sort(fn func(i, j int) bool) *AnyArray[T] {
	sort.Slice(my.data, fn)
	return my
}

// Clean 清理数据
func (my *AnyArray[T]) Clean() *AnyArray[T] {
	my.data = make([]T, 0)
	return my
}

// MarshalJSON 实现接口：json序列化
func (my *AnyArray[T]) MarshalJSON() ([]byte, error) { return json.Marshal(&my.data) }

// UnmarshalJSON 实现接口：json反序列化
func (my *AnyArray[T]) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &my.data) }

// ToString 导出string
func (my *AnyArray[T]) ToString(formats ...string) string {
	var format = "%v"
	if len(formats) > 0 {
		format = formats[0]
	}

	return fmt.Sprintf(format, my.data)
}

// Cast 转换值类型
func Cast[SRC, DST any](aa *AnyArray[SRC], fn func(value SRC) DST) *AnyArray[DST] {
	if aa == nil {
		return nil
	}

	aa.Lock()
	defer aa.Unlock()

	data := make([]DST, len(aa.data))
	for i, v := range aa.data {
		data[i] = fn(v)
	}

	return New(data)
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
