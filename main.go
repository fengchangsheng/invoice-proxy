package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go-simplejson"
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

type Business struct {
	XMLName xml.Name `xml:"business"`
	//Text    string   `xml:",chardata"`
	//ID      string   `xml:"id,attr"`
	//Comment string   `xml:"comment,attr"`
	Body struct {
		//Text       string `xml:",chardata"`
		//Yylxdm     string `xml:"yylxdm,attr"`
		Returncode string `xml:"returncode" json:"returncode"`
		Returnmsg  string `xml:"returnmsg" json:"returnmsg"`
		Returndata Returndata
	} `xml:"body" json:"body"`
}

type Returndata struct {
	XMLName xml.Name `xml:"returndata" json:"returndata"`
	Dqfpdm  string   `json:"dqfpdm" json:"dqfpdm"`
	Dqfphm  string   `json:"dqfphm" json:"dqfphm"`
}

func main() {
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
	//fmt.Fprintf(w, "Hello astaxie!")
	var param = ""
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
		if k == "webXmlStr" {
			param = strings.Join(v, "")
			fmt.Println("===================== param right =======================")
		}
	}
	fmt.Println("=============================== invoke dll ================================")
	dll := syscall.NewLazyDLL("NISEC_SKSC.dll")
	proc := dll.NewProc("PostAndRecvEx")
	fmt.Println("+++++++NewProc:", proc, "+++++++")
	var b = make([]byte, 512)
	//n := "<?xml version=\"1.0\" encoding=\"utf-8\"?><business id=\"20001\" comment=\"参数设置\"><body yylxdm=\"1\"><servletip>tccdzfp.shfapiao.cn</servletip><servletport>80</servletport><keypwd>88888888</keypwd></body></business>"
	println("param:", param)
	data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(param)), simplifiedchinese.GBK.NewEncoder()))
	//fmt.Println("入参为:", ConvertByte2String(data, GB18030))
	_, _, err := proc.Call(uintptr(unsafe.Pointer(&data[0])), uintptr(unsafe.Pointer(&b[0])))
	if err != nil {
		xmlDataRsp := ConvertByte2String(b, GB18030)
		fmt.Println("出参为:", xmlDataRsp)
		v2 := Business{}
		xmlDataRsp = strings.Replace(xmlDataRsp, "gbk", "UTF-8", -1)
		//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(xmlDataRsp)), simplifiedchinese.GBK.NewEncoder()))
		//xml.Unmarshal(data, &v2)
		xml.Unmarshal([]byte(xmlDataRsp), &v2)
		fmt.Printf("Body: %#v\n", v2.Body)
		fmt.Println("---json")
		b2, _ := json.Marshal(v2)
		js, err := simplejson.NewJson([]byte(b2))
		if nil != err {
			log.Fatal(err)
		}
		//fmt.Println("old_json:", ConvertByte2String(b2, UTF8))
		fmt.Println("json:", js)
		body, err := js.Get("body").Map()
		if nil != err {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(body)
		fmt.Printf("VALUE:%v", body)
	}
	println("============================== over =================================== ")
}

func ConvertByte2String(byte []byte, charset Charset) string {
	println("============ len ==============", bytes.Count(byte, nil)-1)
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
