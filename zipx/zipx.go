package zipx

import (
	"compress/flate"

	"github.com/mholt/archiver"
)

func Compress(files []string, dest_path string) (err error) {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}
	err = z.Archive(files, dest_path)
	return
}
func UnCompress(zip_path string, dest_path string) (err error) {
	err = archiver.Unarchive(zip_path, dest_path)
	return
}
