package htmlarticle

import (
	"io/ioutil"
	"log"
	"testing"
)

var (
	data     []byte
	body     string
	databody string
)

func init() {
	data, _ = ioutil.ReadFile("D:/webserver/net/golang/src/ulucu.github.com/htmlarticle/list.txt")
	body = string(data)
}
func TestGetRegion(t *testing.T) {
	var err error
	databody, err = GetRegion(body, `<div class="skin_list">([\s\S]*?)<div id="papelist" class="pagelist">`)
	if err != nil {
		log.Println(databody)
		t.Error(err)
		return
	}
	//log.Println(databody)
}

//<h3 class="list_c_t"><a href="/samxx8/article/details/52963982">bash</a></h3>
func TestGetLink(t *testing.T) {

	outlink, err := GetLink(databody, `<h3.*?><a.*?href="(.*?)">(.*?)</a></h3>`)
	if err != nil {
		t.Error(err)
		return
	}

	log.Println(outlink)
}
