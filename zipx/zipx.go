package zipx

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
)

// Compress 将 files 中的文件/目录打包到 destPath 处的 zip 文件。
func Compress(files []string, destPath string) (err error) {
	// 将磁盘路径映射到压缩包内的路径；空值表示以文件名置于包根目录。
	fileMap := make(map[string]string, len(files))
	for _, f := range files {
		fileMap[f] = ""
	}
	fileInfos, err := archiver.FilesFromDisk(nil, fileMap)
	if err != nil {
		return err
	}

	// 确保目标文件的父目录存在（对应原 MkdirAll 行为）。
	if dir := filepath.Dir(destPath); dir != "" && dir != "." {
		if err = os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	z := archiver.Zip{
		SelectiveCompression: true,
		Compression:          zip.Deflate,
		ContinueOnError:      false,
	}
	return z.Archive(context.Background(), out, fileInfos)
}

// UnCompress 将 zipPath 解压到 destPath 目录。
func UnCompress(zipPath, destPath string) (err error) {
	if err = os.MkdirAll(destPath, 0o755); err != nil {
		return err
	}

	f, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	z := archiver.Zip{}
	return z.Extract(context.Background(), f, func(ctx context.Context, file archiver.FileInfo) error {
		name := path.Clean(filepath.ToSlash(file.NameInArchive))
		if name == "." || strings.HasPrefix(name, "../") {
			return nil
		}

		target := filepath.Join(destPath, filepath.FromSlash(name))
		// 防止 Zip Slip：确保目标路径仍位于 destPath 之内。
		rel, relErr := filepath.Rel(destPath, target)
		if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil
		}

		if file.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		return err
	})
}
