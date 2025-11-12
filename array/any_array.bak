package array

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// JoinWithoutEmpty 去掉空值然后合并
func JoinWithoutEmpty(values []string, sep string) string {
	return strings.Join(RemoveEmpty(values), sep)
}

// RemoveEmpty 去掉数组中的空字符串
func RemoveEmpty[T comparable](slice []T) []T {
	j := 0
	for _, item := range slice {
		ref := reflect.ValueOf(item)
		if !ref.IsZero() {
			slice[j] = item
			j++
		}
	}

	return slice[:j]
}

// In 判断元素是否存在数组中
func In[T comparable](target T, elements []T) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}

	return false
}

// Filter 过滤数组
func Filter[T any](fn func(v T) bool, values []T) (ret []T) {
	for _, value := range values {
		b := fn(value)
		if b {
			ret = append(ret, value)
		}

	}

	return
}

func FilterDemo() {
	type A struct{ Name string }
	a := []*A{
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
	}
	b := Filter(func(a *A) bool {
		return a.Name != "1"
	}, a)

	for _, item := range b {
		println(item.Name)
	}
}

// Max 判断数组中最大值
func Max[T int | int8 | int16 | int32 | int64 |
uint | uint8 | uint16 | uint32 | uint64 |
float32 | float64](values []T) (max T) {
	for _, value := range values {
		if value > max {
			max = value
		}
	}

	return
}

// Min 判断数组中最小值
func Min[T int | int8 | int16 | int32 | int64 |
uint | uint8 | uint16 | uint32 | uint64 |
float32 | float64](values []T) (min T) {
	for _, value := range values {
		if value < min {
			min = value
		}
	}

	return
}

// Sum 获取总和
func Sum[T int | int8 | int16 | int32 | int64 |
uint | uint8 | uint16 | uint32 | uint64 |
float32 | float64](numbers []T) (sum T) {
	for _, num := range numbers {
		sum += num
	}

	return sum
}

// Avg 计算平均值
func Avg[T int | int8 | int16 | int32 | int64 |
uint | uint8 | uint16 | uint32 | uint64 |
float32 | float64](numbers []T) (avg float64) {
	return float64(Sum(numbers)) / float64(len(numbers))
}

// All 判断切片中是否全部是非零值
func All[T comparable](values []T) bool {
	for _, value := range values {
		ref := reflect.ValueOf(value)
		if ref.IsZero() {
			return false
		}
	}

	return true
}

// Any 判断切片中是否包含非零值
func Any[T comparable](values []T) bool {
	for _, value := range values {
		ref := reflect.ValueOf(value)
		if !ref.IsZero() {
			return true
		}
	}

	return false
}

// GroupBy 分组
func GroupBy[T any](array []T, key string) (map[any][]T, error) {
	if len(array) == 0 {
		return nil, fmt.Errorf("切片长度为0")
	}

	ret := make(map[any][]T)

	for _, value := range array {
		ref := reflect.ValueOf(value)

		ret[ref.FieldByName(key).Interface()] =
			append(ret[ref.FieldByName(key).String()], value)
	}

	return ret, nil
}

// Chunk 分块
func Chunk[T any](slice []T, chunkSize int) ([][]T, error) {
	if chunkSize <= 0 {
		return nil, errors.New("切片长度 不能小于等于0")
	}
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks, nil
}

// Pluck 获取数组中指定字段的值
func Pluck[V any, R any](slice []V, key string) ([]R, error) {
	var ret = make([]R, len(slice))
	for _, value := range slice {
		ret = append(ret, reflect.ValueOf(value).FieldByName(key).Interface().(R))
	}

	return ret, nil
}

// Unique 切片去重
func Unique[V any](slice []V) ([]V, error) {
	if len(slice) == 0 {
		return nil, errors.New("切片长度为0")
	}

	var ret1 = make(map[string]V)
	for _, value := range slice {
		ret1[fmt.Sprint(value)] = value
	}

	var ret2 = make([]V, 0, len(ret1))
	for _, value := range ret1 {
		ret2 = append(ret2, value)
	}

	return ret2, nil
}

// NotEmptyLen 判断切片非零值的长度
func NotEmptyLen[T comparable](values []T) int { return len(RemoveEmpty(values)) }

// RemoveTarget 删除数组中对应的目标
func RemoveTarget[T comparable](values []T, target T) (ret []T) {
	for _, value := range values {
		if value != target {
			ret = append(ret, value)
		}
	}

	return ret
}

// RemoveTargets 删除数组中对应的多个目标
func RemoveTargets[T comparable](values []T, targets ...T) (ret []T) {
	for _, value := range values {
		if !In(value, targets) {
			ret = append(ret, value)
		}
	}

	return ret
}

// FromAnyArray 从anyArray转入
func FromAnyArray[T any](anyArray *AnyArray[any]) []T {
	l := Make[T](anyArray.Len())

	for k, v := range anyArray.ToSlice() {
		l.Set(k, v.(T))
	}

	return l.ToSlice()
}

// ToAny converts any slice to []any
func ToAny(slice interface{}) []any {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}

// CopyFromAnyArray 从AnyArray复制并处理转换内容
func CopyFromAnyArray[S any, D any](src *AnyArray[S], processFn func(idx int, item S) D) *AnyArray[D] {
	var dst *AnyArray[D] = Make[D](src.Len())

	if src.IsEmpty() {
		return dst
	}

	if processFn != nil {
		src.Each(func(idx int, item S) {
			dst.Set(idx, processFn(idx, item))
		})
	}

	return dst
}
