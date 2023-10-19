package iox

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
)

// 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dir string, suffix []string) (_files []string, _err error) {
	ok, _ := DirExistOrCreate(dir, true) //不存在就建立
	if !ok {
		logger.Debugf("目录不存在:%s", dir)
		return
	}
	_files = make([]string, 0, 30)
	for i := 0; i < len(suffix); i++ {
		suffix[i] = strings.ToUpper(suffix[i]) //忽略后缀匹配的大小写
	}
	_err = filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			return errors.New("遍历目录出错:" + err.Error())
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if len(suffix) == 0 {
			_files = append(_files, filename)
			return nil
		}
		for _, v := range suffix {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), v) {
				_files = append(_files, filename)
				break
			}
		}
		return nil
	})
	return
}
