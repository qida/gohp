package slice

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

// 数组元素去重复
func RemoveRepeat(arr []interface{}) (newArr []interface{}) {
	newArr = make([]interface{}, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// 数组元素去重复Map
func RemoveDuplicate[T string | int | float64](duplicateSlice []T) []T {
	set := map[T]interface{}{}
	res := []T{}
	for _, item := range duplicateSlice {
		_, ok := set[item]
		if !ok {
			res = append(res, item)
			set[item] = nil
		}
	}
	return res
}

// 删除数组中指定元素
func RemoveElement[T comparable](arr []T, elem T) []T {
	result := arr[:0]
	for _, v := range arr {
		if v != elem {
			result = append(result, v)
		}
	}
	return result
}

// 删除数组中零元素
func RemoveZero(slice []interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if ifZero(v) {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveZero(slice)
		}
	}
	return slice
}

// 判断一个值是否为零值，只支持string,float,int,time 以及其各自的指针，"%"和"%%"也属于零值范畴，场景是like语句
func ifZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
	return false
}

// 2个数组合并并去重
// Merge 是一个泛型函数，用于合并两个切片并去重
func Merge[T comparable](old, now []T) (diff []T) {
	// seen := make(map[T]bool)
	// 遍历 old 切片，将不在 now 切片中的元素添加到 diff 中
	for _, v := range old {
		if !contains(now, v) {
			diff = append(diff, v)
		}
	}
	// 遍历 now 切片，将不在 old 切片中的元素添加到 diff 中
	for _, v := range now {
		if !contains(old, v) {
			diff = append(diff, v)
		}
	}
	return diff
}

// contains 是一个辅助函数，用于检查切片中是否包含某个元素
func contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Append 是一个泛型函数，用于向切片中添加不重复的新元素
func Append[T comparable](slice []T, item T) []T {
	for _, v := range slice {
		if v == item {
			return slice
		}
	}
	return append(slice, item)
}

// Compare 是一个泛型函数，用于比较两个切片是否相等
func Compare[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// 指定查找
// Contains 是一个泛型函数，用于检查切片中是否包含指定的元素
func Contains[T comparable](sl []T, item T) bool {
	for _, v := range sl {
		if v == item {
			return true
		}
	}
	return false
}

// ContainsFunc 利用Map判断指定值是否在Slice切片中存在
//
//	返回一个回调函数，用来判断值知否存在，结果为bool
//	基于泛型形参支持可比较类型，具体定义可参考泛型 comparable 接口
//	sl := []int{1,3,5,7,9}
//	f := IsHasSlice[int](sl)
//	f(2) // false
//	f(5) // true
func ContainsFunc[V comparable](s []V) func(V) bool {
	tmp := make(map[V]struct{}, len(s))
	for _, v := range s {
		tmp[v] = struct{}{}
	}
	return func(key V) bool {
		_, ok := tmp[key]
		return ok
	}
}

// RemoveFunc 切片删除元素
func RemoveFunc[V comparable](s []V, key V) []V {
	i := 0
	for _, v := range s {
		if v != key {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// FindByKeyValue 查找切片中指定字段值的结构体并返回一个新数组
func FindByKeyValue[V any, K comparable](s []V, key string, val K) ([]V, error) {
	tmp := make([]V, len(s))
	for _, v := range s {
		of := reflect.ValueOf(v)
		f := of.FieldByName(key)
		if !f.Type().Comparable() {
			return nil, fmt.Errorf("key [%s] is not comparable type", key)
		}
		k := f.Interface().(K)
		if k == val {
			tmp = append(tmp, v)
		}
	}
	return tmp, nil
}

// RemoveByKeyValue 删除切片中指定字段值的结构体并返回
func RemoveByKeyValue[V any, K comparable](s []V, key string, val K) ([]V, error) {
	var result []V
	for _, v := range s {
		of := reflect.ValueOf(v)
		f := of.FieldByName(key)
		if !f.IsValid() {
			return nil, fmt.Errorf("key [%s] does not exist", key)
		}
		if !f.Type().Comparable() {
			return nil, fmt.Errorf("key [%s] is not comparable type", key)
		}
		k := f.Interface().(K)
		if k != val {
			result = append(result, v)
		}
	}
	return result, nil
}

// FindMixed 查找两个切片的交集
func FindMixed[V comparable](arr1 []V, arr2 []V) []V {
	seen := make(map[V]struct{})
	for _, v := range arr1 {
		seen[v] = struct{}{}
	}
	var intersection []V
	for _, v := range arr2 {
		if _, found := seen[v]; found {
			intersection = append(intersection, v)
		}
	}
	return intersection
}
