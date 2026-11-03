# chinaid

中国居民身份证号工具包：校验 18 位身份证号合法性，并解析出性别、年龄、生日与省市区三级行政区划信息。

导入：`import "github.com/qida/gohp/chinaid"`

## 功能列表

### id.go — 校验与解析
- `IDCard`：身份证号类型（底层为 string）
- `(IDCard).Valid()`：判断身份证号是否合法（长度、字符及末位校验码加权求和校验）
- `(IDCard).Decode()`：根据身份证号解析出具体信息，返回 `IDCardDetail`
- `IDCardDetail`：解析结果结构，含 `Sex`（性别）、`Age`（年龄）、`Birthday`（生日）、`Addr`（地址）
- `Addr`：地址信息（`Code` 行政区划码 / `Province` 省 / `City` 市 / `District` 区县），`InitAddr(code)` 按区划码初始化，无法识别时置为 "未知"
- 常量：`Male`（男，1）、`Female`（女，2）、`IDCardLength`（18）
- 错误：`ErrInvalidIDCardNo`（身份证号不合法）

### id_map.go — 行政区划数据
- `cityMap`：行政区划码到名称的映射表（数据来自民政部 www.mca.gov.cn，文件头部注明勿手动编辑）

### id_test.go / id_benchmark.go — 测试
- `TestIDNum`：解析正确性测试（校验省市区、生日、性别）
- `BenchmarkIDNum`：解析性能基准测试

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/chinaid"
)

func main() {
	detail, err := chinaid.IDCard("341302199006041233").Decode()
	if err != nil {
		panic(err)
	}
	fmt.Println(detail.Province, detail.City, detail.District) // 安徽省 宿州市 埇桥区
	fmt.Println(detail.Sex)      // 1（Male）
	fmt.Println(detail.Birthday) // 1990-06-04 ...
}
```

## 备注
- 年龄按「当前年份 − 出生年份」计算，未按生日精确折算。
- 校验码支持小写 `x` 与大写 `X`。
