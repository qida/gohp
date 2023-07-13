package conv

import (
	"bytes"
	"encoding/binary"
)

func Bytes2ToInt16(b []byte) int16 {
	b_buf := bytes.NewBuffer(b)
	var x int16
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}

func Bytes4ToInt32(b []byte) int32 {
	b_buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}

func Bytes2ToUInt16(b []byte) uint16 {
	b_buf := bytes.NewBuffer(b)
	var x uint16
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}
func Bytes4ToUInt32(b []byte) uint32 {
	b_buf := bytes.NewBuffer(b)
	var x uint32
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}
