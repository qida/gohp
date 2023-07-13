/*
 * @Author: qida
 * @LastEditors: qida
 */
package stringx

import (
	"testing"
)

func TestGetKeysString(t *testing.T) {
	type args struct {
		key_str string
	}
	tests := []struct {
		name       string
		args       args
		wantNumber int
		wantPy     string
		wantHan    string
	}{
		{name: "1", args: args{key_str: "string"}, wantNumber: 0, wantPy: "STRING", wantHan: ""},
		{name: "2", args: args{key_str: "123"}, wantNumber: 123, wantPy: "", wantHan: ""},
		{name: "3", args: args{key_str: "好人"}, wantNumber: 0, wantPy: "", wantHan: "好人"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNumber, gotPy, gotHan := GetKeysString(tt.args.key_str)
			if gotNumber != tt.wantNumber {
				t.Errorf("GetKeysString() gotNumber = %v, want %v", gotNumber, tt.wantNumber)
			}
			if gotPy != tt.wantPy {
				t.Errorf("GetKeysString() gotPy = %v, want %v", gotPy, tt.wantPy)
			}
			if gotHan != tt.wantHan {
				t.Errorf("GetKeysString() gotHan = %v, want %v", gotHan, tt.wantHan)
			}
		})
	}
}

func TestGetKeyWordType(t *testing.T) {
	type args struct {
		key_word string
	}
	tests := []struct {
		name          string
		args          args
		wantType_word int8
	}{
		{name: "1", args: args{key_word: "37108119860813812x"}, wantType_word: Type身份},
		{name: "1", args: args{key_word: "37108119860813812X"}, wantType_word: Type身份},
		{name: "1", args: args{key_word: "371081198608138129"}, wantType_word: Type身份},
		{name: "2", args: args{key_word: "string"}, wantType_word: Type字母},
		{name: "3", args: args{key_word: "13833776549"}, wantType_word: Type手机},
		{name: "4", args: args{key_word: "中国"}, wantType_word: Type汉字},
		{name: "5", args: args{key_word: "123"}, wantType_word: Type数字},
		{name: "6", args: args{key_word: "1sdfd23"}, wantType_word: Type字母},
		{name: "7", args: args{key_word: "14545423"}, wantType_word: Type数字},
		{name: "8", args: args{key_word: "1要，23"}, wantType_word: Type未知},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotType_word := GetKeyWordType(tt.args.key_word); gotType_word != tt.wantType_word {
				t.Errorf("GetKeyWordType() = %v, want %v", gotType_word, tt.wantType_word)
			}
		})
	}
}
