// Copyright 2015 Huan Du. All rights reserved.
// Licensed under the MIT license that can be found in the LICENSE file.
// https://github.com/huandu/xstrings
package stringx

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// versionToIntSlice 将版本号字符串转换为整数切片。
func VersionToInt(version string) int {
	version = removeNonNumericAndDot(version)
	log.Println("version:", version)
	parts := strings.Split(version, ".")
	var result int
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		result = result*1000 + num
	}
	return result
}

// removeNonNumericAndDot 从字符串中删除除数字和点号以外的所有字符。
func removeNonNumericAndDot(str string) string {
	reg := regexp.MustCompile(`[^0-9.]`)
	return reg.ReplaceAllString(str, "")
}
