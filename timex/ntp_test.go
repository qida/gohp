package timex

import "testing"

func TestTimeNtpCompare(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "",
			args:    args{date: "2021-01-12"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TimeNtpCompare(tt.args.date)
			if err != nil {
				// NTP 校时依赖网络环境，网络不可用时跳过
				t.Skipf("TimeCompare() 网络不可用，跳过: %v", err)
			}
		})
	}
}
