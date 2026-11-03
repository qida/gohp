package iox

import (
	"os"
	"testing"
)

func TestGetMP4Duration(t *testing.T) {
	// 需要本地媒体文件；不存在时跳过
	file, err := os.Open("testdata/test.mp4")
	if err != nil {
		t.Skipf("缺少测试媒体文件 testdata/test.mp4，跳过: %v", err)
		return
	}
	defer file.Close()
	gotLengthOfTime, err := GetMP4Duration(file)
	if err != nil {
		t.Errorf("GetMP4Duration() error = %v", err)
		return
	}
	if gotLengthOfTime <= 0 {
		t.Errorf("GetMP4Duration() = %v, want > 0", gotLengthOfTime)
	}
}
