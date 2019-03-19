package main

//#include <stdio.h>
//#include <stdlib.h>
import "C"
import (
	"bytes"
	"strconv"
	"unsafe"
)

func main() {
	cs := C.CString("你好")
	C.free(unsafe.Pointer(cs))
	println(cs)
	//str := C.GoString(unsafe.Pointer(&b[0]))
	//result := C.PostAndRecvEx(cs, str)
	////bytePath := []byte(n + "\x00")
	////str := "A"
	//bin := make([]byte, len(n)+1)
	//copy(bin, n)
	//ret, _, err := proc.Call(uintptr(unsafe.Pointer(cs)), uintptr(unsafe.Pointer(&b[0])))
	//ret, _, err := proc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(n))), uintptr(unsafe.Pointer(&b[0])))
	//C.free(unsafe.Pointer(cs))
	//ret, _, err := proc.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(param))), uintptr(unsafe.Pointer(&b[0])))

}

func copyToByteArray() {
	n := "hah"
	bin := make([]byte, len(n)+1)
	copy(bin, n)
	//unsafe.Pointer(&bin[0])
}

func hexStringToBytes(s string) []byte {
	bs := make([]byte, 0)
	for i := 0; i < len(s); i = i + 2 {
		b, _ := strconv.ParseInt(s[i:i+2], 16, 16)
		bs = append(bs, byte(b))
	}
	return bs
}

func bytesToHexString(b []byte) string {
	var buf bytes.Buffer
	for _, v := range b {
		t := strconv.FormatInt(int64(v), 16)
		if len(t) > 1 {
			buf.WriteString(t)
		} else {
			buf.WriteString("0" + t)
		}
	}
	return buf.String()
}

func BytePtr(s []byte) uintptr {
	return uintptr(unsafe.Pointer(&s[0]))
}
