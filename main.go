package main

import "C"
import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"syscall"
	"unsafe"
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
	//var MB_YESNOCANCEL = 0x00000003
	dll := syscall.NewLazyDLL("NISEC_SKSC.dll")
	proc := dll.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", proc, "+++++++")

	ret, _, err := proc.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("0"))), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("0"))))
	if err != nil {
		fmt.Println("结果为:", err)
	} else {
		fmt.Println("NISEC_SKSC.dll 结果为:", ret)
	}
	println("=========== over ============ ")
}

func callDllOne() {
	//C:\Program Files (x86)\税控服务器组件接口\NISEC_SKSC.dll  method: PostAndRecvEx
	dll32 := syscall.NewLazyDLL("NISEC_SKSC.dll")
	//println("call dll:", dll32.Name)
	pr := dll32.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", pr, "+++++++")
}

func callDllTwo() {
	h, err := syscall.LoadLibrary("NISEC_SKSC.dll")
	if err != nil {
		panic(err.Error())
	}
	defer syscall.FreeLibrary(h)
	proc, err := syscall.GetProcAddress(h, "PostAndRecvEx")
	if err != nil {
		panic(err.Error())
	}
	println(proc)

	//ret, _, err := pr.Call('a')
	//if err != nil {
	//	fmt.Println("lib.dll运算结果为:", ret)
	//}
}

func getDllSuccess() {
	//defer func() {
	//	//恢复程序的控制权
	//	err := recover()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	dll := syscall.MustLoadDLL("NISEC_SKSC.dll")
	println("call dll:", dll.Name)
	procGreet := dll.MustFindProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", procGreet, "+++++++")
}

func abort(funcName string, err error) {
	panic(funcName + " failed: " + err.Error())
}

func print_version(v uint32) {
	major := byte(v)
	minor := uint8(v >> 8)
	build := uint16(v >> 16)
	print("windows version ", major, ".", minor, " (Build ", build, ")\n")
}

func getDll(w http.ResponseWriter, r *http.Request) {
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
	dll32 := syscall.NewLazyDLL("NISEC_SKSC.dll")
	println("call dll:", dll32.Name)
	//fmt.Fprintf(w, dll32.Name)
	pr := dll32.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", pr, "+++++++")

	ret, _, err := pr.Call(uintptr(4))
	if err != nil {
		fmt.Println("NISEC_SKSC.dll 结果为:", ret)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
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
	dll32 := syscall.NewLazyDLL("NISEC_SKSC.dll")
	println("call dll:", dll32.Name)
	//fmt.Fprintf(w, dll32.Name)
	pr := dll32.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", pr, "+++++++")
	//a := []byte{0, 1}
	//params := make([]reflect.Value, 1)                 // 参数
	//params[0] = reflect.ValueOf(20)
	ret, _, err := pr.Call(uintptr(4))

	p := (*byte)(unsafe.Pointer(ret))
	// 定义一个[]byte切片，用来存储C返回的字符串
	data := make([]byte, 0)
	// 遍历C返回的char指针，直到 '\0' 为止
	for *p != 0 {
		data = append(data, *p)        // 将得到的byte追加到末尾
		ret += unsafe.Sizeof(byte(0))  // 移动指针，指向下一个char
		p = (*byte)(unsafe.Pointer(r)) // 获取指针的值，此时指针已经指向下一个char
	}
	name := string(data) // 将data转换为字符串

	fmt.Printf("Hello, %s!\n", name)
	if err != nil {
		fmt.Println("NISEC_SKSC.dll 结果为:", name)
	}

	//re := []byte{0, 1}
	//ret, _, err := pr.Call('a', re)
	//if err != nil {
	//	fmt.Println("NISEC_SKSC.dll 结果为:", ret)
	//}
	//fmt.Fprintf(w, ret)
}
