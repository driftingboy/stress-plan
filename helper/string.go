package helper

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
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

func BytesAddUnit(bytes int64) (sizeWithUnit string) {
	flowBytesStr := strconv.FormatInt(bytes, 10)
	// 获取除去的位数
	var (
		digit float64
		unit  string
	)
	if len(flowBytesStr) > 9 {
		digit = math.Pow10(9)
		unit = "GB"
	} else if len(flowBytesStr) > 6 {
		digit = math.Pow10(6)
		unit = "MB"
	} else if len(flowBytesStr) > 3 {
		digit = math.Pow10(3)
		unit = "KB"
	} else {
		digit = 1.0
		unit = "B"
	}

	sizeWithUnit = fmt.Sprintf("%.2f %s", float64(bytes)/digit, unit)
	return
}
