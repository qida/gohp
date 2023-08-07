package iox

import (
	"path"
	"strings"
)

// 获取文件名，后缀信息
// "D:/software/Typora/bin/typora.exe"
func GetFileInfo(file_path string) (_suffix, _filename, _full_name string) {
	//获取文件名带后缀
	_full_name = path.Base(file_path)
	//获取文件后缀
	_suffix = path.Ext(_full_name)
	//获取文件名
	_filename = strings.TrimSuffix(_full_name, _suffix)
	return
}
