# structx

结构体反射工具集，提供字段校验、零值字段处理、字段赋值、JSON tag 解析等运行时结构体操作。

## 功能列表

### struct.go
- `ShortContent(content interface{}, force ...string) map[string]interface{}`：将结构体非零值字段导出为 map；force 指定的字段即使零值也保留
- `SetField(obj interface{}, name string, value interface{}) error`：反射设置结构体指定字段的值（obj 需为指针）
- `CopyContent(master, second interface{}) map[string]interface{}`：master 的零值字段用 second 对应字段填充后导出为 map
- `JsonToStr(i interface{}) string`：对象序列化为 JSON 字符串
- `RequestVerify(content interface{}, force ...string) error`：校验 force 列出的字段必须全部非零值（对象不能为指针，否则 panic）
- `RequestVerifyOr(content interface{}, force ...string) error`：force 列出的字段任意一个非零值即通过
- `RequestSpaceTrim(obj interface{}) error`：递归去除结构体所有字符串字段中的空格（obj 需为结构体指针）
- `GetStructFieldNames(p interface{}, fields ...string) []string`：获取 json tag 字段名列表；传入字段名时支持白名单模式，传入 `-` 前缀时为排除模式（不能混用，否则 panic）

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/structx"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 必填校验：Name、Age 任一非零值即通过
	err := structx.RequestVerifyOr(User{Name: "qida"}, "Name", "Age")
	fmt.Println(err) // <nil>

	// 非零值字段导出
	m := structx.ShortContent(User{Name: "qida"})
	fmt.Println(m) // map[Name:qida]
}
```

## 备注

- `RequestVerify/RequestVerifyOr` 传入指针会直接 panic，请传值类型。
- `RequestSpaceTrim` 需传入结构体指针才能修改字段，传值会返回错误。
- `GetStructFieldNames` 的白名单与排除（`-` 前缀）模式不能混用。
