package iox

import (
	"testing"
)

func TestWalkDir(t *testing.T) {
	// 遍历当前目录下的 .go 文件，应至少能找到本测试文件
	gotFiles, err := WalkDir("./", []string{".go"})
	if err != nil {
		t.Errorf("WalkDir() error = %v, wantErr %v", err, false)
		return
	}
	if len(gotFiles) == 0 {
		t.Errorf("WalkDir() 返回空列表, want 至少包含 1 个 .go 文件")
	}
}
