# idx

分布式唯一 ID 生成包：提供标准雪花算法（SnowFlake）与 JS 安全的 15 位雪花 ID（SnowFlakeJS）两种实现。

导入：`import "github.com/qida/gohp/idx"`

## 功能列表

### snow_flake.go — 标准雪花 ID
- `SnowFlake`：雪花算法工作节点（互斥锁保证并发安全）
- `NewSnowFlake(workerId)`：实例化工作节点，workerId 范围 0~127（7 位节点 ID）
- `(*SnowFlake).GetId()`：生成 int64 唯一 ID，结构为「毫秒时间戳 + 节点 ID + 15 位序列号」，单节点每毫秒最多生成 32768 个

### snow_flakejs.go — JS 安全的 15 位雪花 ID
- `SnowFlakeJS`：JS 安全雪花节点（无需 workerId）
- `NewSnowFlakeJS()`：实例化
- `(*SnowFlakeJS).GetId()`：生成 int64 唯一 ID，算法为「毫秒时间戳 × 100 + 序列号」，结果不超过 15 位十进制，可被 JavaScript Number 精确表示

### js_snow_flake_test.go — 测试
- `TestSnowFlakeJS_GetId`：验证连续生成 1000 个 ID 的唯一性
- `TestSnowFlakeJS_GetId_Concurrent`：验证 10 个 goroutine 并发生成 1000 个 ID 的唯一性

## 使用示例

```go
package main

import (
	"fmt"

	"github.com/qida/gohp/idx"
)

func main() {
	// 标准雪花 ID（需指定节点编号）
	sf, _ := idx.NewSnowFlake(1)
	fmt.Println(sf.GetId())

	// JS 安全的 15 位雪花 ID（可直接作为 JSON 数值传给前端）
	js := idx.NewSnowFlakeJS()
	fmt.Println(js.GetId())
}
```

## 备注
- JavaScript `Number` 的最大安全整数为 999999999999999（15 位），超出会丢失精度，`SnowFlakeJS` 专用于向前端返回数值型 ID 的场景。
- `SnowFlake` 的起始时间戳 epoch 固定为 2022-07-25 00:00:00（毫秒），一旦投入使用请勿修改，否则可能生成重复 ID。
- `SnowFlakeJS` 的序列号上限复用 `snow_flake.go` 中的 `numberBits`（15 位）常量，单毫秒同样最多 32768 个。
- `SnowFlake.GetId()` 内部会打印 "SnowID: ..." 日志。
