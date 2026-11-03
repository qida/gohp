# rpcx

代码生成工具包：通过反射 Go 结构体自动生成 Thrift IDL 文件。

## 功能列表

### gen_thrift.go
- `GenThrift`：生成器；`PrefixDelete` 字段指定生成 struct 名时要去掉的前缀（仅替换第一次出现）
- `AutoGetThrift(values ...interface{}) error`：反射任意多个结构体，逐个生成 `.thrift` 文件：
  - 结构体名取自类型名的包名后缀，去掉 `PrefixDelete` 前缀后按 `ToSnakeCase` 命名文件（如 `TbBaseLabel` → `base_label.thrift`）
  - 文件内容为 `namespace go <snake_case>` + `struct <Name> {...}`，字段编号从 1 递增
  - Go 类型映射为 Thrift 类型：string/time.Time→string，int/uint/int64/uint64→i64，int32/uint32→i32，int16/uint16→i16，int8/byte/uint8→i8，float64/float32→double，bool→bool；不在映射表内的类型保留原始 Go 类型字符串
  - gorm tag 中的 `COMMENT` 作为字段注释；生成内容同时通过 `logx.Info` 打印
- `ToSnakeCase(str)`：驼峰转下划线命名（每个大写字母前插入 `_` 并转小写）
- `CreateFile(path, content)`：创建文件并写入内容（同名文件会被覆盖）
- `ParseTagSetting(str, sep)`：解析 tag（如 gorm tag）为 map；key 转大写并去空格，支持 `\` 续行，无值时 key 即为值

## 使用示例

```go
package main

import (
	"github.com/qida/gohp/rpcx"
	"your/db" // 含 gorm tag 的结构体
)

func main() {
	g := &rpcx.GenThrift{PrefixDelete: "Tb"}
	// 结构体名去掉 "Tb" 前缀后生成 base_label.thrift、base_label_type.thrift
	err := g.AutoGetThrift(db.TbBaseLabel{}, db.TbBaseLabelType{})
	if err != nil {
		panic(err)
	}
}
```

生成的 `base_label.thrift` 形如：

```thrift
namespace go base_label
struct BaseLabel {
   1:    string   Name , //名称
   2:    i64      Age ,
}
```

## 备注

- 依赖本仓库 `logx` 包（生成时打印 thrift 内容）
- 生成的 `.thrift` 文件写入当前工作目录
- 未在映射表内的字段类型（如嵌套结构体、切片）不会转成 Thrift 类型，需人工修正
