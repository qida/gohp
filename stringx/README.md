# stringx

字符串工具集，提供命名规则转换（驼峰/蛇形/短横线等）、字符串截取、中文数字/金额处理、敏感词 Trie 过滤、AES-GCM 加解密等常用功能。

## 功能列表

### convert2.go（命名规则转换，推荐）
- `NamingRule` 及常量 `NamingRuleSnakeCase/CamelCase/PascalCase/KebabCase/ShoutCase`：命名规则类型定义
- `ConvertNaming(s, rule)`：按指定规则转换命名，智能拆词，支持驼峰、下划线、连字符、空格混合（如 `HTMLParser` → `html_parser`）

### convert.go（命名与字符变换）
- `ToCamelCase`：分词转驼峰（`http_server` → `HttpServer`）
- `ToSnakeCase` / `ToKebabCase`：转蛇形/短横线（`HTTPServer` → `http_server` / `http-server`，含数字边界处理）
- `SwapCase`：大小写互换；`FirstRuneToUpper` / `FirstRuneToLower`：首字母大写/小写
- `Shuffle` / `ShuffleSource`：随机打乱字符（可指定随机源）
- `Successor`：字符串后继（`"ZZZ9999"` → `"AAAA0000"`）

### strconv.go（旧版命名转换，兼容保留）
- `SnakeString` / `CamelString`：驼峰与蛇形互转（`XxYy` → `xx_yy`）
- `StringToByteSlice` / `ByteSliceToString`：基于 unsafe 的零拷贝 string/[]byte 互转

### tf.go（JSON 命名与 Buffer）
- `JsonSnakeCase` / `JsonCamelCase`：包装任意值，Marshal 时将 JSON key 转蛇形/驼峰
- `Camel2Case` / `Case2Camel`：驼峰与下划线互转
- `Ucfirst` / `Lcfirst`：首字母大写/小写
- `Buffer`：内嵌 `bytes.Buffer`，支持链式 `Append`

### string.go（综合工具）
- `DBCtoSBC`：全角转半角
- `Substr` / `Substr2` / `SubString` / `Sub` / `Show_substr`：多种 rune 截取
- `SubstrByByte`：按字节截取且保证 UTF-8 不乱码
- `NumReg`：汉字数字转阿拉伯数字
- `AmountConvert`：金额转人民币大写
- `UnicodeIndex`：子串在字符串中的 rune 位置
- `GetKeysString`：按首字符分类输入（数字/拼音/汉字）
- `GetKeyWordType` + `Type未知/数字/手机/字母/汉字/身份`：识别身份证/手机号/字母/数字/汉字
- `RemoveRepeatJoin`：去除已存在子串后再拼接；`IsLetter`：是否英文字母
- `AESGCMEncrypt` / `AESGCMDecrypt`：AES-GCM 加解密
- `Expand`：将模板中 `{k}` 替换为 map 或位置参数
- `Reverse`：反转字符串（支持 unicode）；`RandomCreateBytes`：生成随机字节

### trie.go / node.go（敏感词过滤）
- `Trie` 接口：`Filter`（过滤并掩码替换）、`FindKeywords`（查找命中的关键词）
- `NewTrie(words, opts...)`：构建 Trie；`WithMask(mask rune)`：自定义掩码字符（默认 `*`）

### version.go / range.go / varx.go
- `VersionToInt`：版本号字符串（如 `1.2.3`）转整数，便于数值比较
- `Range`：`Min/Max` 整数区间，支持 `"min-max"` 格式的 JSON 编解码与数据库 `Value/Scan`
- `StateNames`：状态名称 emoji 数组（⌛🟢🟡🔴）

## 使用示例

```go
package main

import (
	"fmt"
	"github.com/qida/gohp/stringx"
)

func main() {
	// 推荐的命名转换方式
	fmt.Println(stringx.ConvertNaming("HTTPServer", stringx.NamingRuleSnakeCase)) // http_server
	fmt.Println(stringx.ToCamelCase("http_server"))  // HttpServer
	fmt.Println(stringx.ToSnakeCase("HTTPServer"))   // http_server

	// 敏感词过滤
	trie := stringx.NewTrie([]string{"敏感词"})
	sentence, words, _ := trie.Filter("包含敏感词的文本")
	fmt.Println(sentence, words) // 包含***的文本 [敏感词]
}
```

## 备注

命名转换存在三组实现：`ConvertNaming`（convert2.go）、`ToCamelCase/ToSnakeCase/ToKebabCase`（convert.go）、`SnakeString/CamelString`（strconv.go）与 `Camel2Case/Case2Camel`（tf.go），行为细节略有差异。
新代码推荐统一使用 `ConvertNaming` / `ToCamelCase` / `ToSnakeCase` / `ToKebabCase`，旧函数仅为兼容保留。
