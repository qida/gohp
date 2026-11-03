# zh

中文编码转换工具包：提供 UTF-8 与 GB18030 编码字符串的相互转换。

导入：`import "github.com/qida/gohp/zh"`

## 功能列表

### zh.go

本包仅一个源文件，包含以下两个导出函数：

- `Encode(src string) string`：将 UTF-8 字符串编码为 GB18030
- `Decode(src string) string`：将 GB18030 字符串解码为 UTF-8

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/zh"
)

func main() {
	// UTF-8 -> GB18030（向使用 GB 系编码的系统输出数据）
	gb := zh.Encode("中文测试")
	fmt.Println(gb)

	// GB18030 -> UTF-8（读取 GB 系编码系统返回的数据）
	utf := zh.Decode(gb)
	fmt.Println(utf == "中文测试") // true
}
```

## 备注
- 基于 `golang.org/x/text/encoding/simplifiedchinese` 的 GB18030 编解码器实现，内部通过 `transform.NewReader` 流式转换。
- GB18030 向下兼容 GBK 与 GB2312，适用于对接使用这些编码的遗留系统（如部分银行、政务接口）。
- ASCII 字符（英文、数字等）在两种编码下字节相同，转换后保持不变。
- 转换出错时不返回 error，而是返回空字符串。
