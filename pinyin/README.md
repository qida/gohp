# pinyin

汉字转拼音工具包：支持无声调、带声调、首字母大写等多种转换模式，并提供拼音首字母的快捷提取。

导入：`import "github.com/qida/gohp/pinyin"`

## 功能列表

### pinyin.go — 主实现
- `New(origin)`：创建拼音转换器（默认无声调、空格分隔，输入会转小写）
- `Split(split)`：设置多字拼音之间的分隔符（支持链式调用）
- `Mode(mode)`：设置转换模式（支持链式调用）
- `Convert()`：执行转换，返回拼音字符串
- 模式常量（`Mode`）：
  - `WithoutTone`：无声调拼音，如 guo（默认）
  - `Tone`：带声调拼音，如 guó
  - `InitialsInCapitals`：首字母大写无声调，如 Guo
  - `Initials`：仅取首字母，如 G
- `GetPY(src)`：提取汉字拼音首字母（大写），先用正则过滤非汉字、非字母、非数字字符，出错时原样返回
- `GetPYF(src)`：提取第一个字符的拼音首字母

### pinyin_ex.go — 扩展实现
- `GetPinYin(han)`：基于 `github.com/Lofanmi/pinyin-golang` 字典，返回汉字首拼大写或数字

### error.go — 错误定义
- `ErrInitialize`：拼音字典未初始化错误

### pinyin_test.go — 测试
- `TestConvert`：覆盖各模式与分隔符组合的转换测试

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/pinyin"
)

func main() {
	s, _ := pinyin.New("我是中国人").Split(" ").Mode(pinyin.WithoutTone).Convert()
	fmt.Println(s) // wo shi zhong guo ren

	s, _ = pinyin.New("我是中国人ABced").Split("").Mode(pinyin.InitialsInCapitals).Convert()
	fmt.Println(s) // WoShiZhongGuoRenABCED

	fmt.Println(pinyin.GetPY("我是中国人")) // WSZGR
}
```

## 备注
- 主实现在包初始化（`init`）时从 `http://golang.cdn.sunqida.cn/golang/qida/pinyin/pinyin.txt` 加载拼音字典，网络不可用时回退读取本地 `./pinyin.txt`；两者都失败时 `Convert()` 返回 `ErrInitialize`。
- 无法识别的字符在 `Tone` 模式下原样保留。
