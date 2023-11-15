package jsonx

import (
	"testing"
)

func TestReadNginxConf(t *testing.T) {
	err := ReadNginxConf("./nginx.conf.txt")
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}
