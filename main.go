package main

//import "C"
import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net/http"
	"strconv"
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
	http.HandleFunc("/charging-proxy", getDllNew) //设置访问的路由
	err := http.ListenAndServe(":9090", nil)      //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getDllNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	//fmt.Println("scheme", r.URL.Scheme)
	//fmt.Println(r.Form["url_long"])
	fmt.Fprintf(w, "Hello astaxie!")
	var param = ""
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		if k == "webXmlStr" {
			param = strings.Join(v, "")
			fmt.Println("===================== param right =======================")
		}
	}
	fmt.Println("===================== invoke dll =======================")
	// 去除空格
	//param = strings.Replace(param, " ", "", -1)
	// 去除换行符
	//这个写入到w的是输出到客户端的
	dll := syscall.NewLazyDLL("NISEC_SKSC.dll")
	proc := dll.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", proc, "+++++++")
	var b = make([]byte, 256)
	//param = strings.Replace(param, "\n", "", -1)
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(syscall.StringBytePtr(param))), uintptr(unsafe.Pointer(&b[0])))
	if err != nil {
		fmt.Println("出参为:", ConvertByte2String(b, GB18030))
		fmt.Println("NISEC_SKSC.dll 结果为:", ret)
	}
	println("=========== over ============ ")
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
