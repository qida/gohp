# conv

类型转换工具集，提供字符串安全转数值、任意类型转字符串、十六进制转换，以及字节序（大端）与数值/浮点数的互转。

## 功能列表

### convert.go
- `StrTo`：字符串转换类型，链式调用安全转换
  - `Exist()`：字符串是否有效
  - `Uint8() / Int() / Int64() / Float64()`：转数值（返回 error）
  - `MustUint8() / MustInt() / MustInt64() / MustFloat64()`：转数值（忽略错误，失败返回零值）
  - `String()`：转字符串
- `ToStr(value interface{}, args ...int) string`：任意类型转字符串，可选参数控制进制/精度
- `HexStr2int(hexStr string) (int, error)`：十六进制字符串转十进制数
- `Int2HexStr(num int) string`：十进制数转十六进制字符串

### conv_int_byte.go（字节与整数互转，大端序）
- `Bytes2ToInt16(b []byte) int16`：2 字节转 int16
- `Bytes4ToInt32(b []byte) int32`：4 字节转 int32
- `Bytes2ToUInt16(b []byte) uint16`：2 字节转 uint16
- `Bytes4ToUInt32(b []byte) uint32`：4 字节转 uint32

### conv_float_byte.go（字节与浮点数互转，大端序）
- `Bytes4ToFloat32(bytes []byte) float32`：4 字节转 float32
- `Float32Bytes4(float float32) []byte`：float32 转 4 字节
- `Bytes8ToFloat64(bytes []byte) float64`：8 字节转 float64
- `Float64ToBytes8(float float64) []byte`：float64 转 8 字节

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/conv"
)

func main() {
	// StrTo 安全转换
	i, err := conv.StrTo("123").Int()
	fmt.Println(i, err) // 123 <nil>

	// 任意类型转字符串
	fmt.Println(conv.ToStr(3.14, 2)) // 3.14

	// 大端字节转数值
	fmt.Println(conv.Bytes2ToInt16([]byte{0x01, 0x02})) // 258
}
```

## 备注

- 字节相关转换均为大端序（BigEndian，左到右）。
- `Int2HexStr` 仅支持非负数，负数结果不正确。
