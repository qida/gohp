/*
 * @Author: your name
 * @Date: 2021-06-18 10:16:18
 * @LastEditTime: 2021-06-18 10:31:17
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: \em_servere:\gopath\lib\src\github.com\qida\go\zips\zips.go
 */
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
