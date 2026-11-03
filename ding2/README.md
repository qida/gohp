# ding2

钉钉（企业内部应用）API 工具包，基于 `go-dingtalk` 封装企业 AccessToken 的获取，并定义打卡记录等接口结构体。

## 功能列表

### ding2.go

- `Ding`：钉钉客户端结构体，封装 `dingtalk.DingTalkClient`。
- `NewClient(corp_id, corp_secret)`：以企业 ID 与应用 Secret 创建钉钉企业客户端（`DingTalkCompanyClient`）。
- `GetAccessToken()`：刷新并返回企业 AccessToken（内部调用 `RefreshCompanyAccessToken`，从 `d.Client.AccessToken` 读取）。
- `PackData`：钉钉通用响应结构体（errcode、errmsg）。
- `Token`：AccessToken 响应结构体（内嵌 PackData，access_token、expires_in）。
- `WorkRecord`：打卡记录结构体（内嵌 PackData，`HasMore` 是否还有更多、`Recordresult` 打卡结果列表），单条记录字段含：
  - `CheckType`：打卡类型（上班/下班）。
  - `UserCheckTime` / `BaseCheckTime`：实际打卡时间 / 计划打卡时间（Unix 时间戳）。
  - `TimeResult`：时间结果（正常/迟到/早退等）。
  - `LocationResult`：位置结果。
  - `SourceType`：打卡来源。
  - `UserId`、`CorpId`、`GroupId`、`PlanId`、`RecordId`、`WorkDate`、`Id`：关联标识字段。
- `Department`：部门结构体（内嵌 PackData，`SubDeptIDList` 子部门 ID 列表）。

### ding2_test.go

- 当前为空测试文件，仅包含包声明，无实际测试用例。

## 使用示例

```go
d := ding2.NewClient("corp_id", "corp_secret")
token := d.GetAccessToken()
```

Token 可用于拼接钉钉服务端 API 的 `access_token` 查询参数。

## 备注

- 外部依赖：`github.com/icepy/go-dingtalk`。
- 需要钉钉开放平台企业的 CorpId 与应用的 CorpSecret。
- Token 的获取与缓存由 go-dingtalk 库内部维护，`GetAccessToken` 每次调用都会触发刷新。
- 历史上基于 beego httplib 的部门列表代码已注释删除，当前可用功能仅为客户端创建与 Token 获取；`WorkRecord`、`Department`、`Token` 为接口响应结构体定义，可直接用于 JSON 反序列化。
