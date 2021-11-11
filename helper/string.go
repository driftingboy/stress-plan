package helper

import (
	"reflect"
	"unsafe"
)

func UnsafeStrToBytes(str string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func UnsafeBytesToStr(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
