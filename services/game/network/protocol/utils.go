package protocol

import (
	"reflect"
	"unsafe"
)

func bytes2Str(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}

func str2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}
