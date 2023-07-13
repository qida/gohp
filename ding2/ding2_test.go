package ding2

import (
	"reflect"
	"testing"
)

func TestGetDeptment(t *testing.T) {
	tests := []struct {
		name        string
		wantDeptIds []int
		wantErr     bool
	}{
		{"", []int{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeptIds, err := GetDeptment()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDeptment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDeptIds, tt.wantDeptIds) {
				t.Errorf("GetDeptment() = %v, want %v", gotDeptIds, tt.wantDeptIds)
			}
		})
	}
}
