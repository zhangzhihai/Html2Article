package htmlarticle

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	_depth      int = 6   //按行分析的深度，默认为6
	_limitCount int = 220 //字符限定数，当分析的文本数量达到限定数则认为进入正文内容
	// 确定文章正文头部时，向上查找，连续的空行到达_headEmptyLines，则停止查找
	_headEmptyLines int = 2
	// 用于确定文章结束的字符数
	_endLimitCharCount int  = 20
	_appendMode        bool = false
)

//2 格式化标签，剔除匹配标签中的回车符
func FormatTag(bodyText string) (string, error) {

	hrefFilter, _ := regexp.Compile(`(<[^<>]+)\s*\n\s*`)
	lines := hrefFilter.FindAllStringSubmatch(bodyText, -1)
	var sep string

	for _, v := range lines {
		if len(v) > 1 {
			vlen := len(v[1])

			bindex := strings.Index(bodyText, v[1])

			for {

				sep = bodyText[bindex+vlen : bindex+vlen+1]

				if sep == "\n" {

					//tunm := len(bodyText)
					//bodyText = fmt.Sprintf("%s%s", bodyText[0:bindex+vlen], bodyText[bindex+vlen+2:tunm])
					//方案二使用字符串替换
					lnent := fmt.Sprintf("%s%s", v[1], "\n")
					bodyText = strings.Replace(bodyText, lnent, v[1], -1)
				} else {
					break
				}
			}

			//bodyText = re.ReplaceAllString(bodyText, v[1])
		}
	}

	return bodyText, nil
}

//1
func GetTitle(html string) (string, error) {
	titleFilter := regexp.MustCompile(`<title>([\s\S]*?)</title>`)
	title := titleFilter.FindAllStringSubmatch(html, -1)

	if len(title) > 0 && len(title[0]) > 1 {
		return title[0][1], nil
	}

	h1Filter := regexp.MustCompile(`<h1.*?>(.*?)</h1>`)
	title = h1Filter.FindAllStringSubmatch(html, -1)
	if len(title) > 0 && len(title[0]) > 1 {
		return title[0][1], nil
	}

	return "", nil
}

//从body标签文本中分析正文内容
func GetContent(bodyText string) (string, error) {

	var lineInfo string = ""
	re, _ := regexp.Compile("\\</p>|<br.*?/>")
	linre, _ := regexp.Compile("<.*?>")

	lines := strings.Split(bodyText, "\n")
	orgLines := lines

	var Sblines []string
	var orgSblines []string

	for i, v := range lines {

		lineInfo = re.ReplaceAllString(v, "[crlf]")
		lineInfo = linre.ReplaceAllString(lineInfo, "")
		//fmt.Println(lineInfo)

		lines[i] = strings.Join(strings.Fields(lineInfo), "")
	}

	//fmt.Println(strings.Join(lines, "\n"))

	var preTextLen int = 0
	var startPos int = -1

	for i := 0; i < len(lines)-_depth; i++ {

		linlen := 0
		for j := 0; j < _depth; j++ {
			linlen += len(lines[i+j])
		}

		if startPos == -1 {
			if preTextLen > _limitCount && linlen > 0 {
				emptyCount := 0
				for j := i - 1; j > 0; j-- {

					if lines[j] == "" {
						emptyCount++
					} else {
						emptyCount = 0
					}

					if emptyCount == _headEmptyLines {
						startPos = j + _headEmptyLines
						break
					}

				}

				if startPos == -1 {
					startPos = i
				}

				for j := startPos; j <= i; j++ {
					Sblines = append(Sblines, lines[j])
					orgSblines = append(orgSblines, orgLines[j])
				}

			}
		} else {

			if linlen <= _endLimitCharCount && preTextLen < _endLimitCharCount {
				if _appendMode == true {
					break
				}
				startPos = -1
			}

			Sblines = append(Sblines, lines[i])
			orgSblines = append(orgSblines, orgLines[i])
		}
		preTextLen = linlen

	}

	//fmt.Println(Sblines)

	return strings.Join(Sblines, ""), nil

}

//从给定的Html原始文本中获取正文信息
func GetArticle(html string) (string, error) {
	var err error
	lines := strings.Split(html, "\n")
	if len(lines) < 10 {
		html = strings.Replace(html, ">", ">\n", -1)
	}

	bodyFilter := regexp.MustCompile(`<body.*?>([\s\S]*?)<\/body>`)

	body := bodyFilter.FindAllStringSubmatch(html, -1)

	if len(body) > 0 && len(body[0]) > 1 {
		// 过滤样式，脚本等不相干标签
		src := body[0][1]
		///将HTML标签全转换成小写
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		src = re.ReplaceAllStringFunc(src, strings.ToLower)
		//去除STYLE
		re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
		src = re.ReplaceAllString(src, "")
		//去除Html代码中的注释
		re, _ = regexp.Compile("\\<!--.*?--\\>")
		src = re.ReplaceAllString(src, "")

		//针对链接密集型的网站的处理，主要是门户类的网站，降低链接干扰
		re, _ = regexp.Compile("\\</a\\>")
		src = re.ReplaceAllString(src, "</a>\n")

		//去除SCRIPT
		re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
		src = re.ReplaceAllString(src, "")

		re, _ = regexp.Compile("\\<form[\\S\\s]+?\\</form\\>")
		src = re.ReplaceAllString(src, "")

		//去link
		re, _ = regexp.Compile("\\<link[\\S\\s]+?/\\>")
		src = re.ReplaceAllString(src, "")

		//span 也是换行
		re, _ = regexp.Compile("\\</span>")
		src = re.ReplaceAllString(src, "\n")

		//把br,p 替换成换行

		//标签规整化处理，将标签属性格式化处理到同一行
		//body = Regex.Replace(body, @"(<[^<>]+)\s*\n\s*", FormatTag);
		src, err = FormatTag(src)
		//fmt.Println(src)
		return src, err
	}
	return "", nil
}
