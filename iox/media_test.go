package iox

import (
	"os"
	"testing"
)

func TestGetMP4Duration(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {
		file, err := os.Open("http://xxx/audio/rec_customer_follow/1608373085_老潘语音.m4a")
		if err != nil {
			t.Errorf("Err = %v ", err)
			return
		}
		gotLengthOfTime, err := GetMP4Duration(file)
		if err != nil {
			t.Errorf("GetMP4Duration() error = %v", err)
			return
		}
		t.Errorf("GetMP4Duration() = %v ", gotLengthOfTime)
	})
}
