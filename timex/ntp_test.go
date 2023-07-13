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
		// TODO: Add test cases.
		{
			name:    "",
			args:    args{date: "2021-01-12"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TimeNtpCompare(tt.args.date); (err != nil) != tt.wantErr {
				t.Errorf("TimeCompare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
