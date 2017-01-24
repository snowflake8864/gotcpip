package utils

import (
	p "github.com/snowflake8864/libs/public"
	"reflect"
	"unsafe"
)

func PointerAdd(ptr p.Void, s uint) p.Void {
	var result p.Void
	switch v := ptr.(type) {
	case []byte:
		result = v[s:]
	case []uint16:
		result = v[s*2:]
	case []int16:
		result = v[s*2:]
	}
	return result
}
func ShortBe(a uint16) uint16 {
	return (a&0xff)<<8 | (a >> 8)
}

func LongBe(a uint32) uint32 {
	return (a&0xff)<<24 | (((a & 0xff00) << 8) & 0xff0000) | (((a & 0xff0000) >> 8) & 0xff00) | (a >> 24)
}

func ShortFrom(p []uint8) uint16 {
	p0 := uint16(p[0])
	p1 := uint16(p[1])
	return (p0 << 8) + p1
}

func LongFrom(p []uint8) uint32 {
	p0 := uint32(p[0])
	p1 := uint32(p[1])
	p2 := uint32(p[2])
	p3 := uint32(p[3])
	return (p0 << 24) + (p1 << 16) + (p2 << 8) + p3
}

func Byte2uint16(b []byte) uint16 {
	return uint16((uint16(b[1]) << 8) | (uint16(b[0])))
}

func Byte2int16(b []byte) int16 {
	return int16((uint16(b[1]) << 8) | (uint16(b[0])))
}

func Byte2uint32(b []byte) uint32 {
	return uint32((uint32(b[3]) << 24) | (uint32(b[2]) << 16) | (uint32(b[1]) << 8) | uint32(b[0]))
}

func Byte2int32(b []byte) int32 {
	return int32((uint32(b[3]) << 24) | (uint32(b[2]) << 16) | (uint32(b[1]) << 8) | uint32(b[0]))
}

func Byte2Pointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(
		(*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
	)
	//  ptr = (*IPPacket)(unsafe.Pointer(((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)))
	//  fmt.Println("%p", ptr)
	//  return ptr
}

func GetFieldValue(name string, obj p.Void) reflect.Value {

	v := reflect.ValueOf(obj).Elem()
	f := v.FieldByName(name)
	return f
}
