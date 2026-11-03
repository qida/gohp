# iox

文件、目录与媒体文件操作工具集，提供文件信息查询、目录遍历、文件/目录复制、MP4 时长解析等常用功能。

## 功能列表

### file.go — 文件大小与修改时间

- `FileSize(file) (int64, error)`：获取文件大小（字节）。
- `FileSizeString(file) string`：获取带单位的文件大小（如 "1.5MB"），失败返回 "unknown"。
- `FileSizeHumane(s uint64) string`：将字节数格式化为人类可读字符串（B/KB/MB/GB/TB/PB/EB）。
- `FileMTimeInt64 / FileMTimeString / FileMTime(file)`：获取文件修改时间，分别返回 Unix 秒、格式化字符串（2006-01-02 15:04:05）、time.Time。
- `SecondToTime(sec int64) time.Time`：秒级时间戳转 time.Time。
- 存储单位常量：`Byte`、`KByte`、`MByte`、`GByte`、`TByte`、`PByte`、`EByte`。

### file_linux.go / file_windows.go — 文件创建时间（平台特定实现）

- `FileCTimeInt64 / FileCTimeString / FileCTime(file)`：获取文件创建时间，Linux 基于 `syscall.Stat_t`，Windows 基于 `syscall.Win32FileAttributeData`。

### os.go — 文件基础操作

- `Copy(src, dest string) error`：复制文件，支持符号链接，并保留原文件的修改时间与权限。
- `WriteFile(filename string, data []byte) error`：写入文件，自动创建上级目录。
- `IsFile(filePath string) bool`：判断路径存在且为文件（非目录）。
- `IsExist(path string) bool`：判断目录或文件是否存在。

### dir.go — 目录遍历与复制

- `StatDir(rootPath, includeDir...)`：深度优先收集目录下的文件相对路径，可选包含子目录（子目录以 "/" 结尾）。
- `LstatDir(rootPath, includeDir...)`：同上，但跟随符号链接。
- `GetAllSubDirs / LgetAllSubDirs(rootPath)`：获取所有子目录，L 版本跟随符号链接。
- `GetFileListBySuffix(dirPath, suffix)`：按后缀列出目录下的文件（不递归；路径本身是文件时直接返回）。
- `CopyDir(srcPath, destPath, filters...)`：递归复制目录，支持过滤函数，要求目标目录事先不存在。
- `DirExistOrCreate(path string, is_create bool)`：检测目录是否存在，可选自动创建（会打印创建日志）。
- `IsDir(fileAddr string) bool`：判断是否为目录。

### copy.go — 目录/文件拷贝（旧实现）

- `CopyDir1(srcPath, destPath)`：递归拷贝目录内容，目标目录名自动追加时间戳后缀（如 dest_20260102150405）。
- `CopyFile(src, dest)`：拷贝文件并自动逐级创建缺失的目标目录。

### walk.go — 递归遍历

- `WalkDir(dir string, suffix []string) ([]string, error)`：递归获取目录及所有子目录下的文件，按后缀列表过滤（忽略大小写）；目录不存在时自动创建。

### info.go — 文件名信息

- `GetFileInfo(file_path)`：一次性返回文件全名（含后缀）、后缀（含 "."）、不含后缀的文件名。

### media.go — MP4 媒体解析

- `GetMP4Duration(reader io.ReaderAt) (uint32, error)`：解析 MP4 box 结构（定位 moov，读取 timeScale 与 Duration），返回视频时长（秒）。
- `BoxHeader`：MP4 box 信息头结构体。

## 使用示例

获取 MP4 视频时长（提炼自 media_test.go）：

```go
file, err := os.Open("testdata/test.mp4")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

duration, err := iox.GetMP4Duration(file) // 返回时长（秒）
```

## 备注

- `walk.go` 依赖本仓库的 `github.com/qida/gohp/logx`。
- 文件创建时间函数（`FileCTime*`）通过构建标签分别提供 linux / windows 两套实现。
- `GetMP4Duration` 基于 `io.ReaderAt` 按需读取，无需将整个文件载入内存。
