# mathx

数学工具集，提供整数幂运算等基础数学功能，作为标准库 `math` 的补充。

## 功能列表

### math.go
- `PowInt(x int, y int) int`：整数幂运算，即 `math.Pow` 的 int 版本
  - 使用快速幂算法（递归折半）实现，时间复杂度 O(log y)
  - `y <= 0` 时统一返回 1（含 `y == 0` 与负指数）
  - 参数与返回值均为 `int`，无浮点精度问题
  - 偶指数折半：`sqrt := PowInt(x, y/2)` 后自乘，减少乘法次数
  - 奇指数拆分：`PowInt(x, y-1) * x`，逐次递归到 `y == 0` 返回 1

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/mathx"
)

func main() {
	fmt.Println(mathx.PowInt(2, 10))  // 1024
	fmt.Println(mathx.PowInt(2, 0))   // 1
	fmt.Println(mathx.PowInt(3, -2))  // 1（负指数返回 1）
	fmt.Println(mathx.PowInt(-2, 3))  // -8
	fmt.Println(mathx.PowInt(5, 3))   // 125
	fmt.Println(mathx.PowInt(1, 999999)) // 1（快速幂，递归深度为 O(log y)）
	fmt.Println(mathx.PowInt(10, 6))  // 1000000
}
```

## 备注

- 与 `math.Pow` 的差异：负指数不返回小数（如 `PowInt(2, -1)` 返回 1，而 `math.Pow(2, -1)` 返回 0.5）。
- 无溢出检查，结果超出 int 范围时会溢出（如 `PowInt(2, 63)`）。
- 底数为 0：`PowInt(0, n)`（n>0）返回 0，`PowInt(0, 0)` 按规则返回 1。
