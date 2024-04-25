// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

//go:build linux
// +build linux

package iox

import (
	"os"
	"syscall"
	"time"
)

// 获取创建时间
func FileCTimeInt64(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}

	linuxFileAttr := f.Sys().(*syscall.Stat_t)
	// fmt.Println("文件创建时间", SecondToTime(linuxFileAttr.Ctim.Sec))
	// fmt.Println("最后访问时间", SecondToTime(linuxFileAttr.Atim.Sec))
	// fmt.Println("最后修改时间", SecondToTime(linuxFileAttr.Mtim.Sec))
	unix, _ := linuxFileAttr.Ctim.Unix()
	return unix, nil
}

// 获取创建时间
func FileCTimeString(file string) (string, error) {
	f, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	linuxFileAttr := f.Sys().(*syscall.Stat_t)
	return time.Unix(linuxFileAttr.Ctim.Unix()).Format("2006-01-02 15:04:05"), nil
}

// 获取创建时间
func FileCTime(file string) (time.Time, error) {
	f, err := os.Stat(file)
	if err != nil {
		return time.Time{}, err
	}
	linuxFileAttr := f.Sys().(*syscall.Stat_t)

	return time.Unix(linuxFileAttr.Ctim.Unix()), nil
}
