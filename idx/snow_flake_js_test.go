package idx

import (
	"fmt"
	"testing"
)

func TestSnowFlakeJS_GetId(t *testing.T) {
	jstr := NewSnowFlakeJS()
	tests := []struct {
		name string
		tr   *SnowFlakeJS
		want int64
	}{}
	for i := 0; i < 1001; i++ {
		tests = append(tests, struct {
			name string
			tr   *SnowFlakeJS
			want int64
		}{
			name: fmt.Sprintf("aaa_%d", i),
			tr:   jstr,
			want: 1,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetId(); got != tt.want {
				t.Errorf("SnowFlakeJS.GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
