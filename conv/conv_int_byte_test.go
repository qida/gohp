package conv

import "testing"

func TestBytes2ToInt16(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		b []byte
		// Expected results.
		want int16
	}{
		// TODO: Add test cases.
		{"1 [2]byte Int16", []byte{0x00, 0x00}, 0},
		{"2 [2]byte Int16", []byte{0x80, 0x00}, -32768},
		{"3 [2]byte Int16", []byte{0xff, 0x00}, -256},
		{"4 [2]byte Int16", []byte{0xff, 0xff}, -1},
	}
	for _, tt := range tests {
		if got := Bytes2ToInt16(tt.b); got != tt.want {
			t.Errorf("%q. Bytes2ToInt16() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBytes4ToInt32(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		b []byte
		// Expected results.
		want int32
	}{
		// TODO: Add test cases.
		{"1 [4]byte Int32", []byte{0x00, 0x00, 0x00, 0x00}, 0},
		{"2 [4]byte Int32", []byte{0x80, 0x00, 0x00, 0x00}, -2147483648},
		{"3 [4]byte Int32", []byte{0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"4 [4]byte Int32", []byte{0x12, 0x34, 0x56, 0x78}, 305419896},
	}
	for _, tt := range tests {
		if got := Bytes4ToInt32(tt.b); got != tt.want {
			t.Errorf("%q. Bytes4ToInt32() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBytes2ToUInt16(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		b []byte
		// Expected results.
		want uint16
	}{
		// TODO: Add test cases.
		{"1 [2]byte UInt16", []byte{0x00, 0x00}, 0},
		{"2 [2]byte UInt16", []byte{0x80, 0x00}, 32768},
		{"3 [2]byte UInt16", []byte{0xFF, 0xFF}, 65535},
	}
	for _, tt := range tests {
		if got := Bytes2ToUInt16(tt.b); got != tt.want {
			t.Errorf("%q. Bytes2ToUInt16() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBytes4ToUInt32(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		b []byte
		// Expected results.
		want uint32
	}{
		// TODO: Add test cases.
		{"1 [4]byte UInt32", []byte{0x00, 0x00, 0x00, 0x00}, 0},
		{"2 [4]byte UInt32", []byte{0x80, 0x00, 0x00, 0x00}, 2147483648},
		{"3 [4]byte UInt32", []byte{0xFF, 0xFF, 0xFF, 0xFF}, 4294967295},
		{"4 [4]byte UInt32", []byte{0x12, 0x34, 0x56, 0x78}, 305419896},
	}
	for _, tt := range tests {
		if got := Bytes4ToUInt32(tt.b); got != tt.want {
			t.Errorf("%q. Bytes4ToUInt32() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
