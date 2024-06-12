package iox

import (
	"reflect"
	"testing"
)

func TestWalkDir(t *testing.T) {
	type args struct {
		dirPth string
		suffix []string
	}
	tests := []struct {
		name       string
		args       args
		want_files []string
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				dirPth: "./",
				suffix: []string{".go"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got_files, err := WalkDir(tt.args.dirPth, tt.args.suffix)
			if (err != nil) != tt.wantErr {
				t.Errorf("WalkDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got_files, tt.want_files) {
				t.Errorf("WalkDir() = %v, want %v", got_files, tt.want_files)
			}
		})
	}
}
