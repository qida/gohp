# cdn

七牛云存储与 CDN 工具包，封装七牛对象的上传、删除、移动、私有链接生成及 CDN 缓存刷新等操作。

## 功能列表

### qiniu.go

- `QiNiu`：七牛客户端结构体，聚合 `qbox.Mac`（鉴权）、`cdn.CdnManager`（CDN 管理）、`storage.BucketManager`（空间管理）、`storage.Config`（存储配置）、`storage.PutPolicy`（上传策略）。
- `NewQiNiu(bucket, url, accessKey, secretKey)`：根据 AK/SK 自动获取存储区域并初始化客户端；`url` 为外链域名，配置为使用 CDN 域名加速、非 HTTPS。
- `Upload(localFile io.Reader, size, file_name, is_refresh)`：表单上传流式数据（Scope 为 `bucket:key`，覆盖上传同名文件）；当 `is_refresh` 为 true 且文件已存在时，上传后自动刷新 CDN 缓存。
- `UploadFile(file_name, file_data []byte)`：上传字节内容并在完成后自动刷新 CDN。
- `Refresh(url ...string)`：刷新指定 URL 的 CDN 缓存，自动补全 `http:` 前缀。
- `RefreshDir(dir)`：刷新目录缓存（需联系七牛技术支持开通权限，单次最多 10 个目录）。
- `MoveFile(src_url, dst_url)`：在同一 Bucket 内移动（重命名）对象，目标存在时强制覆盖。
- `GetPrivateMediaUrl(src_url)`：生成 60 分钟有效的私有空间访问地址。
- `Delete(url)`：删除指定 URL 对应的对象（自动去除外链域名前缀，url 为空时直接返回）。
- `GetTokenUpload(region, key)`：生成前端直传所需信息 map，包含 Region、UpToken、Domain、Key（key 为空时按 `temp/时间戳` 生成）。

## 使用示例

```go
q := cdn.NewQiNiu("my-bucket", "http://img.example.com/", "accessKey", "secretKey")

// 上传字节内容（完成后自动刷新 CDN）
url, err := q.UploadFile("test/a.jpg", fileBytes)

// 生成 60 分钟有效的私有访问链接
priv := q.GetPrivateMediaUrl(url)

// 删除对象
err = q.Delete(url)
```

## 备注

- 外部依赖：`github.com/qiniu/api.v7/v7`（auth/qbox、cdn、storage）与 `golang.org/x/net/context`。
- 需要七牛云账号的 AccessKey/SecretKey、已创建的 Bucket 及外链域名。
- CDN 链接刷新与目录刷新会消耗七牛配额；目录刷新需提前开通权限。
- 上传策略 Scope 均为 `bucket:key` 形式，即同名文件会被覆盖。
