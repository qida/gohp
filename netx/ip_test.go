//from https://github.com/freshcn/qqwry

package netx

import (
	"testing"
)

func TestQQwry_Find(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		q    *QQwry
		args args
		want string
	}{
		{
			name: "test1",
			q:    NewQQwry("qqwry.dat"),
			args: args{ip: "39.78.34.61"},
			want: "山东省临沂市",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Find(tt.args.ip); got != tt.want {
				t.Errorf("QQwry.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getKey(t *testing.T) {
	tests := []struct {
		name    string
		want    uint32
		wantErr bool
	}{
		{
			"", 0, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
