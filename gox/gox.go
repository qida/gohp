package gox

import (
	"reflect"
	"runtime"
	"strings"
)

// 获取函数名称
func GetFuncName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
