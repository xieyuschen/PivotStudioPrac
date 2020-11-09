package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
)

var href_reg *regexp.Regexp

var hrefs_been_found map[string]int

var hrefs_undone []string

func get_all_href(url string) []string {
	var ret []string
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return ret
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	hrefs := href_reg.FindAllString(string(body), -1)

	for _, v := range hrefs {
		str := strings.Split(v, "\"")[1]

		if len(str) < 1 {
			continue
		}

		switch str[0] {
		case 'h':
			ret = append(ret, str)
		case '/':
			if len(str) != 1 && str[1] == '/' {
				ret = append(ret, "http:"+str)
			}

			if len(str) != 1 && str[1] != '/' {
				ret = append(ret, url+str[1:])
			}
		default:
			ret = append(ret, url+str)

		}

	}

	return ret
}

func init_global_var() {
	href_pattern := "href=\"(.+?)\""
	href_reg = regexp.MustCompile(href_pattern)

	hrefs_been_found = make(map[string]int)
}

func is_href_been_found(href string) bool {
	_, ok := hrefs_been_found[href]
	return ok
}

func add_hrefs_to_undone_list(hrefs []string) {
	for _, value := range hrefs {
		ok := is_href_been_found(value)
		if !ok {
			fmt.Printf("new url:(%s)\n", value)
			hrefs_undone = append(hrefs_undone, value)
			hrefs_been_found[value] = 1
		} else {
			hrefs_been_found[value]++
		}

	}
}
func send(a string) {
	/*完成信息抓取*/
	//url := "https://tophub.today/n/mproPpoq6O"
	//url := "https://zhihu.com"
	//url := "https://www.liwenzhou.com/"
	url := a
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	//req.Header.Add("Cookie", "_zap=eb58e67e-2071-41a5-9eb8-aa0731680631; _xsrf=huHWSQwxScw2VxvC60wxcXw2tYBRSmFt; KLBRSID=76ae5fb4fba0f519d97e594f1cef9fab|1604489317|1604486875")
	//req.Header.Add("Cookie", "ThW9_934f_saltkey=FcCSJxm3; ThW9_934f_lastvisit=1604626786; ThW9_934f_sid=HNZ1NP; ThW9_934f_lastact=1604630386%09index.php%09")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	/* 正则转换*/
	src := string(body)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	/*fmt.Println(strings.TrimSpace(src))*/
	/*写入文件*/
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	//l, err := f.WriteString(strings.TrimSpace(src)) /*转换后*/
	l, err := f.WriteString(string(body)) /*转换前*/
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	/*发送邮件*/
	identity1 := ""
	sender1 := "2038975825@qq.com"
	pwd1 := "dumqnhmcxthbbgfg"
	host1 := "smtp.qq.com"
	port1 := "25"
	sendTo1 := []string{"2038975825@qq.com"}
	senderName1 := "ZhiHuHot"
	title1 := "ZhiHuHot"
	//body1 := strings.TrimSpace(src)
	body1 := string(body)
	auth1 := smtp.PlainAuth(identity1, sender1, pwd1, host1)
	content_type1 := "Content-Type: text/html; charset=UTF-8"
	msg1 := []byte("To: " + strings.Join(sendTo1, ",") + "\nFrom: " + senderName1 +
		"<" + sender1 + ">\nSubject: " + title1 + "\n" + content_type1 + "\n" + body1 + "\n")

	url1 := host1 + ":" + port1
	err1 := smtp.SendMail(url1, auth1, sender1, sendTo1, msg1)
	if err1 != nil {
		fmt.Printf("\n\nsend mail error: %v", err)
		return
	}
	fmt.Println("\n\nsend mail succes")
}

func main() {
	init_global_var()

	var pos = 0
	//var urls = []string{"https://tophub.today/n/mproPpoq6O"}
	var urls = []string{"https://www.liwenzhou.com/"}
	add_hrefs_to_undone_list(urls)

	for j := 0; j < 10; j++ {
		if pos >= len(hrefs_undone) {
			break
		}
		url := hrefs_undone[0]
		hrefs_undone = hrefs_undone[1:]

		hrefs := get_all_href(url)
		send(url)
		add_hrefs_to_undone_list(hrefs)
		time.Sleep(time.Second / 10)
	}
	fmt.Println("end")
}
