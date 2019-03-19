package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
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
	//去除换行符
	//param = strings.Replace(param, "\n", "", -1)
	//这个写入到w的是输出到客户端的
	dll := syscall.NewLazyDLL("NISEC_SKSC.dll")
	proc := dll.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", proc, "+++++++")
	var b = make([]byte, 256)
	n := "<?xml version=\"1.0\" encoding=\"utf-8\"?><business id=\"20001\" comment=\"参数设置\"><body yylxdm=\"1\"><servletip>tccdzfp.shfapiao.cn</servletip><servletport>80</servletport><keypwd>88888888</keypwd></body></business>"
	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(n)), simplifiedchinese.GBK.NewEncoder()))
	println("param:", param)
	fmt.Println("入参为:", ConvertByte2String(data, GB18030))
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&b[0])))
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
