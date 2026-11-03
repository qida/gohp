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
	}{
		{
			name: "test1",
			q:    NewQQwry("qqwry.dat"),
			args: args{ip: "39.78.34.61"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// IP 归属地数据会随 qqwry.dat 版本变化，只校验能查到非空结果
			if got := tt.q.Find(tt.args.ip); got == "" {
				t.Errorf("QQwry.Find() 返回为空, want 非空归属地")
			}
		})
	}
}

func Test_getKey(t *testing.T) {
	got, err := GetKey()
	if err != nil {
		t.Skipf("网络不可用，跳过测试: %v", err)
	}
	if got == 0 {
		t.Errorf("GetKey() = %v, want 非0", got)
	}
}
