package util

import (
	"unsafe"
	"reflect"
)

func Str2Byte(s *string) []byte {
	return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}
