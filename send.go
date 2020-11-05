package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

func main() {

	url := "https://tophub.today/n/mproPpoq6O"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	//req.Header.Add("Cookie", "_zap=eb58e67e-2071-41a5-9eb8-aa0731680631; _xsrf=huHWSQwxScw2VxvC60wxcXw2tYBRSmFt; KLBRSID=76ae5fb4fba0f519d97e594f1cef9fab|1604489317|1604486875")

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
	/*write*/
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(strings.TrimSpace(src))
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
	/*写入文件*/
	identity1 := ""
	sender1 := "2038975825@qq.com"
	pwd1 := "dumqnhmcxthbbgfg"
	host1 := "smtp.qq.com"
	port1 := "25"
	sendTo1 := []string{"2038975825@qq.com"}
	senderName1 := "ZhiHuHot"
	title1 := "ZhiHuHot"
	body1 := strings.TrimSpace(src)
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
	fmt.Println("\n\nsend mail success!")
	/*发送文件*/
}
