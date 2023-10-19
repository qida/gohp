package iox

import (
	"path"
	"strings"
)

// 获取文件名，后缀信息
// req: "D:/software/Typora/bin/typora.exe"
// resp: typora.exe , .exe , typora
func GetFileInfo(file_path string) (_full_name, _suffix, _file_name string) {
	//获取文件名带后缀
	_full_name = path.Base(file_path)
	//获取文件后缀
	_suffix = path.Ext(_full_name)
	//获取文件名
	_file_name = strings.TrimSuffix(_full_name, _suffix)
	return
}
