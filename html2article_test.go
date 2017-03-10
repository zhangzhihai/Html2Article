package htmlarticle

import (
	"fmt"
	"io/ioutil"

	"testing"
)

//读取本地数据
var (
	data     []byte
	body     string
	databody string
)

func init() {
	data, _ = ioutil.ReadFile("D:/webserver/net/golang/src/ulucu.github.com/htmlarticle/2.txt")
	body = string(data)
}

func TestGetTitle(t *testing.T) {
	title, err := GetTitle(body)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(title)
}

func TestGetArticle(t *testing.T) {
	bdy, err := GetArticle(body)
	if err != nil {
		t.Error(err)
		return
	}
	databody = bdy

}

func TestFormatTag(t *testing.T) {
	cbody, err := FormatTag(databody)
	if err != nil {
		t.Error(err)
		fmt.Println(cbody)
		return
	}
	//fmt.Println(cbody)
}

func TestGetContent(t *testing.T) {
	cbody, err := GetContent(databody)
	if err != nil {
		t.Error(err)
		fmt.Println(cbody)
		return
	}

	fmt.Println(cbody)
}
