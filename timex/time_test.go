package timex

import (
	"testing"
	"time"
)

func TestGetZeroTimeOfDay(t *testing.T) {
	tests := []struct {
		name string
		d    time.Time
		want time.Time
	}{
		{
			name: "当天零点",
			d:    time.Date(2026, 11, 3, 18, 45, 57, 628695800, LOC_ZONE),
			want: time.Date(2026, 11, 3, 0, 0, 0, 0, LOC_ZONE),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetZeroTimeOfDay(tt.d); !got.Equal(tt.want) {
				t.Errorf("GetZeroTimeOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
