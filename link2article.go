package htmlarticle

//连接列表

import (
	"regexp"
)

//提取页面的连接
//连接的规则解析

//var reg, _ := regexp.Compile(`<title>([\s\S]*?)</title>`)

//div的区域
func GetRegion(html, re string) (string, error) {
	reg, _ := regexp.Compile(re)
	region := reg.FindAllStringSubmatch(html, -1)

	if len(region) > 0 && len(region[0]) > 1 {
		return region[0][1], nil
	}
	return "", nil
}

//提取所有的连接,和title

func GetLink(html, re string) (map[string]string, error) {

	reg, _ := regexp.Compile(re)

	region := reg.FindAllStringSubmatch(html, -1)

	var m map[string]string
	m = make(map[string]string)

	if len(region) > 0 {
		for _, v := range region {
			m[v[1]] = v[2]
		}
	}

	return m, nil
}
