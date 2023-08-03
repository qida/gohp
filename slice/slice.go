package slice

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

func FindElement(arr []interface{}) (newArr []interface{}) {
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

func RemoveRepeatedElement(arr []interface{}) (newArr []interface{}) {
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

func RemoveRepeatedString(arr []string) (newArr []string) {
	newArr = make([]string, 0)
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
func RemoveRepeatedInt(arr []int) (newArr []int) {
	newArr = make([]int, 0)
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

// 元素去重
func RemoveRepeatedInt64(arr []int64) []int64 {
	occurred := map[int64]bool{}
	result := []int64{}
	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func RemoveRepeatedFloat64(arr []float64) (newArr []float64) {
	newArr = make([]float64, 0)
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

func RemoveElement(arr []interface{}, elem interface{}) []interface{} {
	if len(arr) == 0 {
		return arr
	}
	for i, v := range arr {
		if v == elem {
			arr = append(arr[:i], arr[i+1:]...)
			return RemoveElement(arr, elem)
		}
	}
	return arr
}
func RemoveInt(arr []int, elem int) []int {
	if len(arr) == 0 {
		return arr
	}
	for i, v := range arr {
		if v == elem {
			arr = append(arr[:i], arr[i+1:]...)
			return RemoveInt(arr, elem)
		}
	}
	return arr
}
func RemoveString(arr []string, elem string) []string {
	if len(arr) == 0 {
		return arr
	}
	for i, v := range arr {
		if v == elem {
			arr = append(arr[:i], arr[i+1:]...)
			return RemoveString(arr, elem)
		}
	}
	return arr
}
func RemoveZero(slice []interface{}) []interface{} {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if IfZero(v) {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveZero(slice)
		}
	}
	return slice
}

// 判断一个值是否为零值，只支持string,float,int,time 以及其各自的指针，"%"和"%%"也属于零值范畴，场景是like语句
func IfZero(arg interface{}) bool {
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

// 去重合并
func FindAddString(old, now []string) (diff []string) {
	for i := 0; i < len(old); i++ {
		for j := 0; j < len(now); j++ {
			if old[i] == now[j] {
				break
			}
		}
		diff = append(diff, old[i])
	}
	return
}

// 指定删除
func FindSubString(old []string, now string) (diff []string) {
	for i := 0; i < len(old); i++ {
		if strings.Contains(now, old[i]) {

		} else {
			diff = append(diff, old[i])
		}
	}
	return
}

// 指定查找
func FindInt(old []int, now int) bool {
	for i := 0; i < len(old); i++ {
		if now == old[i] {
			return true
		}
	}
	return false
}

// AppendStr appends string to slice with no duplicates.
func AppendStr(strs []string, str string) []string {
	for _, s := range strs {
		if s == str {
			return strs
		}
	}
	return append(strs, str)
}

// CompareSliceStr compares two 'string' type slices.
// It returns true if elements and order are both the same.
func CompareSliceStr(s1, s2 []string) bool {
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

// CompareSliceStrU compares two 'string' type slices.
// It returns true if elements are the same, and ignores the order.
func CompareSliceStrU(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		for j := len(s2) - 1; j >= 0; j-- {
			if s1[i] == s2[j] {
				s2 = append(s2[:j], s2[j+1:]...)
				break
			}
		}
	}
	if len(s2) > 0 {
		return false
	}
	return true
}

// IsSliceContainsStr returns true if the string exists in given slice, ignore case.
func IsSliceContainsStr(sl []string, str string) bool {
	str = strings.ToLower(str)
	for _, s := range sl {
		if strings.ToLower(s) == str {
			return true
		}
	}
	return false
}

// IsSliceContainsInt64 returns true if the int64 exists in given slice.
func IsSliceContainsInt64(sl []int64, i int64) bool {
	for _, s := range sl {
		if s == i {
			return true
		}
	}
	return false
}

// SliceHas 利用Map判断指定值是否在Slice切片中存在
//
//	返回一个回调函数，用来判断值知否存在，结果为bool
//	基于泛型形参支持可比较类型，具体定义可参考泛型 comparable 接口
//	sl := []int{1,3,5,7,9}
//	f := SliceHas[int](sl)
//	f(2) // false
//	f(5) // true
func SliceHas[V comparable](s []V) func(V) bool {
	tmp := make(map[V]struct{}, len(s))
	for _, v := range s {
		tmp[v] = struct{}{}
	}
	return func(key V) bool {
		_, ok := tmp[key]
		return ok
	}
}

// SliceDelete 切片删除元素
func SliceDelete[V comparable](s []V, key V) []V {
	i := 0
	for _, v := range s {
		if v != key {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// SliceStructHas 利用Map判断指定字段值的结构体是否在Slice切片中存在
//
//	返回一个回调函数，用来判断指定字段值知否存在，结果为bool
//	基于泛型形参支持可比较类型，具体定义可参考泛型 comparable 接口
//	利用反射获取结构体指定字段，判断是否为可比较类型，并赋值给map的key
//	sl := []User{
//		{
//			Name:    "alpha",
//			Age:     20,
//			Sex:     "male",
//			Tickets: []string{"001", "002"},
//		},
//		{
//			Name:    "beta",
//			Age:     21,
//			Sex:     "female",
//			Tickets: []string{"003", "004"},
//		},
//	}
//	f := SliceStructHas[User, string](sl, "Name")
//	f("alpha") // true
//	f("sigma") // false
func SliceStructHas[V any, K comparable](s []V, key string) (func(K) bool, error) {
	tmp := make(map[K]V, len(s))
	for _, v := range s {
		of := reflect.ValueOf(v)
		f := of.FieldByName(key)
		if !f.Type().Comparable() {
			return nil, fmt.Errorf("key [%s] is not comparable type", key)
		}
		k := f.Interface().(K)
		tmp[k] = v
	}
	return func(key K) bool {
		_, ok := tmp[key]
		return ok
	}, nil
}

// SliceStructPop 查找切片中指定字段值的结构体并返回一个新数组
func SliceStructPop[V any, K comparable](s []V, key string, val K) ([]V, error) {
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

// SliceStructDelete 删除切片中指定字段值的结构体并返回
func SliceStructDelete[V any, K comparable](s []V, key string, val K) ([]V, error) {
	i := 0
	for _, v := range s {
		of := reflect.ValueOf(v)
		f := of.FieldByName(key)
		if !f.Type().Comparable() {
			return nil, fmt.Errorf("key [%s] is not comparable type", key)
		}
		k := f.Interface().(K)
		if k != val {
			s[i] = v
			i++
		}
	}
	return s[:i], nil
}

// SliceIntersection 查找两个切片的交集
func SliceIntersection[V comparable](arr1 []V, arr2 []V) []V {
	seen := make(map[V]bool)
	for _, str := range arr1 {
		seen[str] = true
	}

	var intersection []V
	for _, str := range arr2 {
		if seen[str] {
			intersection = append(intersection, str)
		}
	}
	return intersection
}
