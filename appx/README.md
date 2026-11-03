# appx

应用级全局初始化包，通过 `init()` 在程序启动时创建全局雪花 ID 生成器，供各业务包直接使用。

## 功能列表

### app.go

- `SNOW_FLAK`：全局雪花 ID 生成器变量（类型 `*idx.SnowFlakeJS`），在包 `init()` 中通过 `idx.NewSnowFlakeJS()` 自动初始化。
- 引入本包即完成初始化，无需手动调用任何函数：
  - 业务入口正常 import 后可直接使用 `appx.SNOW_FLAK.GetId()`。
  - 若只需要触发初始化副作用，也可使用空导入（`import _`）。

## 使用示例

```go
package main

import "github.com/qida/gohp/appx"

func main() {
    id := appx.SNOW_FLAK.GetId() // 获取 int64 全局唯一 ID
    println(id)
}
```

如需空导入触发初始化副作用：

```go
import _ "github.com/qida/gohp/appx"
```

## 备注

- 内部依赖本仓库 `idx` 包的 `SnowFlakeJS`（JS 雪花算法实现，`idx` 包内附唯一性与并发生成测试）。
- 无第三方外部依赖。
- 生成器为进程内单例，多实例部署时 ID 唯一性依赖 `idx.SnowFlakeJS` 自身的算法保证。
- 适合需要全局唯一 ID 的场景（如数据表主键、业务流水号等）。
- 引入即初始化，无外部服务依赖，初始化不会失败或阻塞。
- 同一进程内多处引用 `SNOW_FLAK` 共享同一生成器实例。
