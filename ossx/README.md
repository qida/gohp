# ossx

对象存储（OSS）工具包，封装阿里云 OSS 与 MinIO 两种存储后端的常用操作（上传、删除、列举、统计等）。

## 功能列表

### aliyun.go（阿里云 OSS）

- `OssAliyun`：阿里云 OSS 客户端结构体，持有 `oss.Client` 与 `oss.Bucket`。
- `NewOssAliyun(expires, end_point, name_bucket, access_key_id, access_key_secret)`：初始化阿里云 OSS 客户端，失败时记日志并返回 nil。
- `BucketList()`：列出账号下所有 Bucket 名称。
- `UploadFile(time_out, file_path_name, prop, file)`：上传文件，`file` 支持本地路径（string）或字节内容（[]byte）；上传时设置公共读 ACL、元数据 prop 与过期时间，返回文件访问路径（即完整相对路径，如 `/22年10月/xxx.jpg`）。
- `DeleteFile(file_name)`：删除指定对象。
- `DeleteAllFiles()`：遍历并删除 Bucket 内全部对象（限测试使用）。
- `GetBucketStat()`：获取 Bucket 统计信息（存储量、对象数、各存储类型明细等，以 map 返回）。

### minio.go（MinIO）

- `OssMinIO`：MinIO 客户端结构体，封装 `minio.Client`。
- `NewOssMinIO(expires, end_point, name_bucket, access_key_id, access_key_secret)`：初始化 MinIO 客户端（非 HTTPS），失败时记日志并返回 nil。
- `UploadFile(time_out, file_name, prop, file)`：上传文件，`file` 同样支持 string 路径与 []byte；带超时控制，返回对象 Key。
- `ListObjects()`：递归列举 Bucket 内全部对象并打印 Key 与大小。
- `DeleteFile(file_name)`：删除指定对象。
- `DeleteAllFiles()`：批量删除 Bucket 内全部对象（限测试使用）。

## 使用示例

```go
// 阿里云 OSS 上传
oss := ossx.NewOssAliyun(3, "oss-cn-hangzhou.aliyuncs.com", "my-bucket", "ak", "sk")
url, err := oss.UploadFile(30*time.Second, "/2022/TPKEY_m_0001.jpg", "prop", []byte("文件内容"))

// MinIO 删除文件
mini := ossx.NewOssMinIO(3, "127.0.0.1:9000", "my-bucket", "ak", "sk")
err = mini.DeleteFile("a.jpg")
```

## 备注

- 阿里云依赖：`github.com/aliyun/aliyun-oss-go-sdk/oss`。
- MinIO 依赖：`github.com/minio/minio-go/v7`。
- 内部依赖本仓库 `logx`（日志）、`gox`（获取当前函数名）。
- 需要相应云服务的 AccessKey 以及已创建的 Bucket；MinIO 需自建或可访问的 MinIO 服务。
- 两者的 `DeleteAllFiles` 均为清空 Bucket 的危险操作，注释中已标明"限测试使用"。
