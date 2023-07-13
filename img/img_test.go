package img

import "testing"

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Create()
		})
	}
}
