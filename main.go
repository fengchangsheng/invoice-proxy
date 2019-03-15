package main

import "C"
import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net/http"
	"strings"
	"syscall"
	"unsafe"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func main() {
	//callDllOther()
	//http.HandleFunc("/", sayhelloName) //设置访问的路由
	http.HandleFunc("/", getDllNew)          //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getDllNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
	dll := syscall.NewLazyDLL("NISEC_SKSC.dll")
	proc := dll.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", proc, "+++++++")
	var b = make([]byte, 256)
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr("hello"))), uintptr(unsafe.Pointer(&b[0])))
	if err != nil {
		fmt.Println("出参为:", ConvertByte2String(b, GB18030))
		fmt.Println("NISEC_SKSC.dll 结果为:", ret)
	}
	println("=========== over ============ ")
}

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
