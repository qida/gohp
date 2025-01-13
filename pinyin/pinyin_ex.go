package pinyin

import (
	"fmt"
	"strings"

	pinyin3 "github.com/Lofanmi/pinyin-golang/pinyin"
)

// ===================公共=======================//
var dict = pinyin3.NewDict()

const split string = " "

// 返回汉字首拼大写或数字
func GetPinYin(han string) string {
	if len(han) == 0 {
		return ""
	}
	s := dict.Convert(han, split).ASCII()
	s1 := strings.Split(s, split)
	s = ""
	for _, v := range s1 {
		if len([]rune(v)) > 0 {
			s = fmt.Sprintf("%s%s", s, string(([]rune(v))[0]))
		}
	}
	return strings.ToUpper(s)
}
