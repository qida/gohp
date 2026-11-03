# slice

切片（数组）工具集，基于泛型提供去重、合并、交集、按字段查找/删除结构体等常用操作。

## 功能列表

### slice.go
- `RemoveRepeat(arr []interface{}) []interface{}`：数组元素去重（非泛型版本）
- `RemoveDuplicate[T string | int | float64]([]T) []T`：泛型切片去重
- `RemoveElement[T comparable](arr []T, elem T) []T`：删除切片中指定元素
- `RemoveZero(slice []interface{}) []interface{}`：删除零值元素（""、0、nil、`%`、`%%` 等，适配 like 查询场景）
- `Merge[T comparable](old, now []T) []T`：返回两个切片中互不存在的元素（对称差集）
- `Append[T comparable](slice []T, item T) []T`：向切片添加不重复的新元素
- `Compare[T comparable](s1, s2 []T) bool`：比较两个切片是否相等
- `Contains[T comparable](sl []T, item T) bool`：检查切片是否包含指定元素
- `ContainsFunc[V comparable](s []V) func(V) bool`：基于 map 构造 O(1) 的存在性判断函数
- `RemoveFunc[V comparable](s []V, key V) []V`：原地删除切片中所有等于 key 的元素
- `FindByKeyValue[V any, K comparable](s []V, key string, val K) ([]V, error)`：按字段值查找结构体并返回新切片（反射实现）
- `RemoveByKeyValue[V any, K comparable](s []V, key string, val K) ([]V, error)`：删除字段值匹配的结构体
- `SliceStructPop[V any, K comparable](s []V, key string, val K) ([]V, error)`：筛出字段值匹配的结构体
- `FindMixed[V comparable](arr1, arr2 []V) []V`：查找两个切片的交集
- `FuncHasSlice[V comparable](s []V) func(V) bool`：同 `ContainsFunc`，基于 map 的存在性判断函数

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/slice"
)

func main() {
	// 去重（来自 slice_test.go 的用例）
	fmt.Println(slice.RemoveDuplicate([]string{"a", "b", "a", "c", "b"})) // [a b c]

	// 交集
	fmt.Println(slice.FindMixed([]int{1, 2, 3}, []int{2, 3, 4})) // [2 3]

	// 构造存在性判断函数
	f := slice.ContainsFunc([]int{1, 3, 5})
	fmt.Println(f(2), f(5)) // false true
}
```

## 备注

- `Merge` 的实现为对称差集（返回仅在其中一个切片出现的元素），并非注释所说的"合并去重"，使用时请注意实际行为。
- `FindByKeyValue/RemoveByKeyValue/SliceStructPop` 基于反射按字段名操作结构体切片，要求字段类型可比较。
