package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go-simplejson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"strings"
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
	} `xml:"body" json:"body"`
}

func main() {
	xmlDataRsp := "<?xml version=\"1.0\" encoding=\"gbk\"?><business id=\"20001\" comment=\"参数设置\"><body yylxdm=\"1\"><returncode>-3</returncode><returnmsg>数字证书口令验证失败：没有安装USBKEY驱动(0xB9)</returnmsg></body></business>"
	v2 := Business{}
	xmlDataRsp = strings.Replace(xmlDataRsp, "gbk", "UTF-8", -1)
	xml.Unmarshal([]byte(xmlDataRsp), &v2)
	//data, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(xmlDataRsp)), simplifiedchinese.GBK.NewEncoder()))
	//xml.Unmarshal(data, &v2)

	//decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlDataRsp)))
	//decoder.CharsetReader = func(c string, i io.Reader) (io.Reader, error) {
	//	return charset.NewReaderLabel(c, i)
	//}

	//result := &Business{}
	//decoder.Decode(result)
	//fmt.Println(result)
	//fmt.Println(v2.XMLName)
	fmt.Printf("Body: %#v\n", v2.Body)
	//fmt.Println("---json")
	b2, _ := json.Marshal(&v2)
	js, err := simplejson.NewJson([]byte(b2))
	if nil != err {
		log.Fatal(err)
	}
	fmt.Println("old_json:", ConvertByte2String(b2, UTF8))
	fmt.Println("json:", js)
	body, err := js.Get("body").Map()
	if nil != err {
		log.Fatal(err)
	}
	fmt.Printf("VALUE:%v", body)
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
