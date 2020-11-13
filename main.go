package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func RequestGet(url string) string {
	/*完成信息抓取*/
	method := "GET"
	client := &http.Client{}
	req, _ := http.NewRequest(method, url, nil)
	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	http := string(body)
	return http
}
func Change(html string) (htmlchange string) {
	/* 正则转换*/
	src := html
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

	htmlchange = strings.TrimSpace(src)
	return htmlchange
}
func GetTheUrl(html string) (url []string, content []string, hot []string) {
	/*提取html中的url,文章简介,热度等信息*/
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("td[class='al']>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		url = append(url, "https://tophub.today"+href)
	})
	dom.Find("td[class='al']").Each(func(i int, selection *goquery.Selection) {
		con := selection.Text()
		content = append(content, con)
	})
	dom.Find("td[class='al']+td").Each(func(i int, selection *goquery.Selection) {
		h := selection.Text()
		hot = append(hot, h)
	})
	return url, content, hot
}
func GetTheArticle(html string) (article []string) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("script[id='js-initialData']").Each(func(i int, selection *goquery.Selection) {
		art := selection.Text()
		article = append(article, art)
	})
	return article

}
func Sendmail(body string) {

	/*发送邮件*/
	identity := ""
	sender := "2038975825@qq.com"
	pwd := "dumqnhmcxthbbgfg"
	host := "smtp.qq.com"
	port := "25"
	sendTo := []string{"2038975825@qq.com"}
	senderName := "ZhiHuHot"
	title := "ZhiHuHot"
	auth := smtp.PlainAuth(identity, sender, pwd, host)
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + strings.Join(sendTo, ",") + "\nFrom: " + senderName +
		"<" + sender + ">\nSubject: " + title + "\n" + content_type + "\n" + body + "\n")

	url := host + ":" + port
	err := smtp.SendMail(url, auth, sender, sendTo, msg)
	if err != nil {
		fmt.Printf("\n\nsend mail error: %v", err)
		return
	}
	fmt.Println("\n\nsend mail succes")
	return
}
func Write(content []string, hot []string, urlofhtml []string) {
	f, err := os.Create("test.txt")
	if err != nil {

		fmt.Println(err)
		return
	}
	for i := 0; i < 10; i++ {
		_, _ = f.WriteString(content[i])
		_, _ = f.WriteString("\n")
		_, _ = f.WriteString(hot[i])
		_, _ = f.WriteString("\n")
		_, _ = f.WriteString(urlofhtml[i])
		_, _ = f.WriteString("\n")
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func Done() {
	/*启动函数*/
	url := "https://tophub.today/n/mproPpoq6O"

	/*对知乎发送请求并且获取html信息*/
	html := RequestGet(url)

	/*发送邮件，直接发送html信息，qq邮箱会对其自动进行渲染*/
	Sendmail(html)

	/*对html信息进行分析，获取上面的url，热度，内容等信息*/
	urlofhtml, content, hot := GetTheUrl(html)
	/*写入文件*/
	Write(content, hot, urlofhtml)

	/*对从html中得到的url信息再次发送请求并且获取html信息*/
	AllArticle := ""
	for j := 0; j < 3; j++ {
		newhtml := RequestGet(urlofhtml[j])
		fmt.Println(Change(newhtml))
		AllArticle += Change(newhtml)
	}
	Sendmail(AllArticle)

}
func main() {
	/*主函数*/
	for i := 0; i < 1; i++ {
		Done()
		fmt.Println("sleep")
		/*每12小说定时启动*/
		time.Sleep(time.Hour * 12)
	}
}
