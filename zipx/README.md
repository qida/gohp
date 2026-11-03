# zipx

基于 archiver/v4 的 zip 压缩/解压工具：打包时自动创建目标父目录，解压时清理条目路径并防范 Zip Slip 路径穿越攻击。

## 功能列表

### zipx.go

- `Compress(files []string, destPath string) error`：将一组文件/目录打包为 `destPath` 处的 zip 文件。
  - 通过 `archiver.FilesFromDisk` 收集源路径，包内路径映射为空值，即所有条目直接位于压缩包根目录（不保留源目录结构）。
  - 自动创建 `destPath` 的父目录（权限 0755）。
  - 压缩配置：`zip.Deflate` 算法、`SelectiveCompression: true`（对已压缩内容跳过再压缩）、出错即中止（`ContinueOnError: false`）。
  - `files` 中的目录会被递归收集其下所有文件与子目录。
- `UnCompress(zipPath, destPath string) error`：将 zip 文件解压到 `destPath` 目录。
  - 自动创建目标目录（权限 0755）。
  - 对包内条目名做 `path.Clean` 归一化，跳过 "." 及 "../" 开头的条目。
  - 通过 `filepath.Rel` 校验每个解压目标仍位于 `destPath` 之内，越界条目直接跳过，防范 Zip Slip。
  - 目录条目创建目录（0755），文件条目通过 `os.Create` 写入。
  - 每个文件条目写入前会用 `os.MkdirAll` 创建其父目录，支持带子目录结构的包。

## 使用示例

```go
// 打包：文件与目录混合，均置于压缩包根目录
if err := zipx.Compress([]string{"./docs", "./config.yaml"}, "./backup.zip"); err != nil {
    log.Fatal(err)
}

// 解压到指定目录
if err := zipx.UnCompress("./backup.zip", "./restore"); err != nil {
    log.Fatal(err)
}
```

## 备注

- 外部依赖：`github.com/mholt/archiver/v4`。
- `Compress` 不保留源目录结构，多个同名文件/目录混用时会在包内冲突。
- 解压时遇到非法路径条目不会报错，而是静默跳过。
- 两个函数均为同步阻塞调用，不支持 context 取消（内部固定使用 `context.Background()`）。
