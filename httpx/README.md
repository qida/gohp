# httpx

HTTP 工具集：封装 resty 的客户端、Gin 统一响应与中间件、JWT 生成校验、文件上传辅助与 SHA256 哈希。

## 功能列表

### client.go — HTTP 客户端（基于 go-resty）
- `ClientHttp` / `NewClientHttp()`：默认 10s 超时、重试 3 次（0.5s 起步、最大 20s）、JSON Content-Type
- `Debug(bool)`：开关 resty 调试输出
- `LogFile(dir)`：请求日志写入 `dir/go-resty.log`（进程内仅初始化一次）
- `SetTimeout(seconds)`：设置超时
- `PostBody(url, req, resp, header)`：POST JSON body，响应自动反序列化到 `resp`
- `Post(url, form, resp, header)`：POST form 表单
- `Get(url, resp, header)` / `GetParams(url, params, resp, header)`：GET（可带 query 参数）
- `UploadFile(url, nameParam, filePath, params, resp, header)`：multipart 文件上传
- `UploadFileFromBytes(url, nameParam, fileName, data, params, resp, header)`：内存字节流上传

### gin.go — Gin 响应与中间件
- `Reply` / `NewReply()`：统一响应结构（code/status/error/msg/run_ms/data）
- `Return(msg, reply, start)`：填充响应，出错时附加错误发生的文件与行号
- `ReturnGin(c, msg, reply, start)`：填充并以 JSON 返回
- `Recover()` / `RecoverGin(c)`：panic 恢复（Gin 版本返回 400 JSON）
- `CorsGin()`：CORS 跨域中间件（放行 OPTIONS 预检）

### jwt.go — JWT
- `Claims`：含 `UserInfo` 与 `jwt.RegisteredClaims`
- `JWTGenerate(user_info)`：生成 HS256 token（1 小时有效期）
- `JWTValidate(token_string)`：校验 token 并解析出 `Claims`

### hash.go — SHA256
- `HashHttp(file *multipart.FileHeader)`：计算上传文件的 SHA256 哈希
- `HashFile(file *multipart.File)`：计算已打开文件的 SHA256 哈希

### upload.go — 上传辅助
- `ReadFileDataFromMultipart(f)`：读取 multipart 文件全部内容
- `NewFileDataUploadRequest(url, path, data, params)`：构造 multipart 文件上传的 `*http.Request`

### httpx.go
- `SetUrlEncode(url)`：JSON 编码 URL（关闭 HTML 转义）

## 使用示例

```go
client := httpx.NewClientHttp()
var result map[string]interface{}
err := client.Get("http://localhost:8080/api", &result, nil)
```

## 备注

- 外部依赖：github.com/go-resty/resty/v2、github.com/gin-gonic/gin、github.com/golang-jwt/jwt/v5
- JWT 签名密钥硬编码在 jwt.go 中（`hello jwt`），生产环境使用前需自行修改
