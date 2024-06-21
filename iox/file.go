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

package iox

import (
	"fmt"
	"math"
	"os"
	"time"
)

// Storage unit constants.
const (
	Byte  = 1
	KByte = Byte * 1024
	MByte = KByte * 1024
	GByte = MByte * 1024
	TByte = GByte * 1024
	PByte = TByte * 1024
	EByte = PByte * 1024
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%dB", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+"%s", val, suffix)
}

// HumaneFileSize calculates the file size and generate user-friendly string.
func HumaneFileSize(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1024, sizes)
}

// 获取文件大小
func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// 获取文件大小带单位
func FileSizeString(file string) string {
	f, err := os.Stat(file)
	if err != nil {
		return "unknown"
	}
	return HumaneFileSize(uint64(f.Size()))
}

// 获取文件修改时间
func FileMTimeInt64(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// 获取文件修改时间
func FileMTimeString(file string) (string, error) {
	f, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	return f.ModTime().Format("2006-01-02 15:04:05"), nil
}

// 获取文件修改时间
func FileMTime(file string) (time.Time, error) {
	f, err := os.Stat(file)
	if err != nil {
		return time.Time{}, err
	}
	return f.ModTime(), nil
}

// 把秒级的时间戳转为time格式
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}
