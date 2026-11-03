# gox

Go 运行时工具集，提供函数名称获取等基于 reflect / runtime 的辅助功能。

## 功能列表

### gox.go
- `GetFuncName(i interface{}, seps ...rune) string`：获取传入函数的名称
  - 通过 `runtime.FuncForPC` + `reflect.ValueOf(i).Pointer()` 解析
  - 可变参数 `seps` 指定分隔符，返回按分隔符切分后的最后一段（如传 `'.'` 仅取不含包路径的函数名）
- `GetCurrFuncName() string`：获取当前调用函数的完整名称
  - 基于 `runtime.Caller(1)` 获取上一层调用者的 PC，再解析为函数名
  - 返回形如 `main.hello` 的带包路径名称

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/gox"
)

func hello() {
	// 获取当前函数名
	fmt.Println(gox.GetCurrFuncName()) // main.hello
}

func main() {
	hello()
	// 获取指定函数名，按 . 分割取最后一段
	fmt.Println(gox.GetFuncName(hello, '.')) // hello
	// 不传分隔符时返回完整名称
	fmt.Println(gox.GetFuncName(hello)) // main.hello
}
```

## 备注

- `GetFuncName` 需传入函数值本身（而非调用结果），常用于日志埋点、中间件注册等场景。
- 两个函数返回的完整名称均带包路径前缀（如 `main.hello`、`github.com/xxx/pkg.Func`）。
