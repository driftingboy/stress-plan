package helper

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// TODO 汇编
// var g_goid_offset uintptr = func() uintptr {
//     g := GetGroutine()
//     if f, ok := reflect.TypeOf(g).FieldByName("goid"); ok {
//         return f.Offset
//     }
//     panic("can not find g.goid field")
// }()

// func GetGoid() int64 {
//     g := getg()
//     p := (*int64)(unsafe.Pointer(uintptr(g) + g_goid_offset))
//     return *p
// }

func GetGoidSlowly() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
