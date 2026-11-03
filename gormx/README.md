# gormx

GORM 辅助工具包，提供分页查询 Scope 与模糊查询通配符拼接，简化常见查询写法。

## 功能列表

### gormx.go

- `Paginate(page, limit int64) func(db *gorm.DB) *gorm.DB`：返回 GORM Scope 函数：
  - 按 `offset = (page - 1) * limit` 计算偏移量。
  - 配合 `db.Scopes(gormx.Paginate(page, limit)).Find(&list)` 使用实现分页。
  - `page` 从 1 开始计数。
- `Wildcard[T any](key T) string`：泛型函数，将任意类型的关键字格式化为 `%v%` 形式的模糊查询通配符：
  - 例如 `Wildcard("张")` → `%张%`，可直接用于 `LIKE` 条件。
  - 支持任意类型（字符串、数字等），通过 `fmt.Sprintf("%%%v%%", key)` 拼接。

## 使用示例

```go
// 分页 + 模糊查询
var users []User
db.Scopes(
    gormx.Paginate(2, 10), // 第 2 页，每页 10 条
).Where("name LIKE ?", gormx.Wildcard("张")).
    Find(&users)

// 等价于：db.Offset(10).Limit(10).Where("name LIKE ?", "%张%").Find(&users)
```

分页参数换算示意：

| page | limit | Offset | Limit |
| ---- | ----- | ------ | ----- |
| 1    | 10    | 0      | 10    |
| 2    | 10    | 10     | 10    |
| 5    | 20    | 80     | 20    |

## 备注

- 外部依赖：`gorm.io/gorm`。
- `Paginate` 未做 page <= 0 的防御，传入 0 或负数会产生负 Offset，调用方需自行保证入参合法。
- `Wildcard` 不会转义值本身包含的 `%`、`_` 等 LIKE 元字符，必要时需调用方先处理。
- `Wildcard` 生成的 `%` 前后通配形式无法利用索引，大表模糊查询请注意性能。
