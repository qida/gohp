package timex

import (
	"reflect"
	"testing"
	"time"
)

func TestGetZeroTimeOfDay(t *testing.T) {
	type args struct {
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{
			name: "t",
			args: args{d: time.Now()},
			want: time.Now(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetZeroTimeOfDay(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZeroTimeOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
