package htmlarticle

import (
	"log"
	"testing"
)

//读取本地数据
var (
	client *HTTPClient
)

func init() {
	client, _ = HttpNew("https://my.oschina.net/tongjh/blog/266051")
}

func TestDo(t *testing.T) {
	body, err := client.Do("https://my.oschina.net/tongjh/blog/266051", "")
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(body))
}
