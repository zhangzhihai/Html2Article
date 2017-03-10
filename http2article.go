package htmlarticle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type (
	HTTPClient struct {
		Conn   http.Client
		Domain string
	}
)

var timeout int = 10

func HttpNew(Domain string) (s *HTTPClient, err error) {

	jar, err := cookiejar.New(nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, address string) (net.Conn, error) {
				deadline := time.Now().Add(time.Duration(timeout*1) * time.Second)
				c, err := net.DialTimeout(network, address, time.Duration(timeout*2)*time.Second)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
		Jar: jar,
	}

	if strings.Contains(Domain, "http") == false {
		Domain = fmt.Sprintf("http://%s", Domain)
	}

	s = &HTTPClient{
		Conn:   client,
		Domain: Domain,
	}

	return s, nil
}

func SaveJar(cookie []*http.Cookie) error {
	b, err := json.Marshal(cookie)
	if err != nil {
		log.Printf("error ", err)
		log.Println(b)
		return err
	}

	//log.Println(string(b))
	return nil
}

func SetHeader(req *http.Request) *http.Request {

	req.Header.Set("Pragma", "no-cache")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, sdch") //压缩要自己解压
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36")
	req.Header.Set("Accept", "text/javascript, text/html, application/xml, text/xml, */*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cache-Control", "no-cache")

	return req

}

func (s *HTTPClient) Do(url string, data string) ([]byte, error) {

	//"name=cjb"

	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}

	req = SetHeader(req)

	resp, err := s.Conn.Do(req) //发送
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}

	outdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error", err)
		return []byte(""), err
	}

	cookies := resp.Cookies()
	SaveJar(cookies)

	resp.Body.Close()
	return outdata, nil
}
