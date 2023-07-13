/*
 * @Author: qida
 * @Date: 2019-05-25 07:56:43
 * @LastEditors: sunqida
 * @LastEditTime: 2019-05-25 07:56:43
 * @Description:
 */
package conv

import (
	"reflect"
	"testing"
)

func TestBytes4ToFloat32(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		bytes []byte
		// Expected results.
		want float32
	}{
		// TODO: Add test cases.
		{"1 [4]byte 转成 float32", []byte{0x44, 0xc8, 0x00, 0x00}, 1600.00},
		{"2 [4]byte 转成 float32", []byte{0xbe, 0x80, 0x00, 0x00}, -0.25},
		{"3 [4]byte 转成 float32", []byte{0x41, 0xfd, 0xca, 0x58}, 31.7238},
		{"3 [8]byte 转成 float32", []byte{0x41, 0xfd, 0xca, 0x58 /*因为是正序，所以会忽略后面的数*/, 0x40, 0xfd, 0xca, 0x58}, 31.7238},
	}
	for _, tt := range tests {
		if got := Bytes4ToFloat32(tt.bytes); got != tt.want {
			t.Errorf("%q. Bytes4ToFloat32() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFloat32Bytes4(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		float float32
		// Expected results.
		want []byte
	}{
		// TODO: Add test cases.
		{"1 float32 转 [4]byte", 1600.00, []byte{0x44, 0xc8, 0x00, 0x00}},
		{"2 float32 转 [4]byte", -0.25, []byte{0xbe, 0x80, 0x00, 0x00}},
	}
	for _, tt := range tests {
		if got := Float32Bytes4(tt.float); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Float32Bytes4() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBytes8ToFloat64(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		bytes []byte
		// Expected results.
		want float64
	}{
		// TODO: Add test cases.
		{"1 [8]byte 转成 float64", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{"2 [8]byte 转成 float64", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}, 5e-324},
		{"3 [8]byte 转成 float64", []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 7.291122019556398e-304},
	}
	for _, tt := range tests {
		if got := Bytes8ToFloat64(tt.bytes); got != tt.want {
			t.Errorf("%q. Bytes8ToFloat64() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFloat64ToBytes8(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		float float64
		// Expected results.
		want []byte
	}{
		// TODO: Add test cases.
		{"1 [8]byte 转成 float64", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"2 [8]byte 转成 float64", 5e-324, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{"3 [8]byte 转成 float64", 7.291122019556398e-304, []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tt := range tests {
		if got := Float64ToBytes8(tt.float); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Float64ToBytes8() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
