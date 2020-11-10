//将热榜前20写入txt
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery" //用goquery解析回答
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp" //正则表达式
	"strconv"
	"strings"
	"time"
)

var filename = "./今日热榜.txt"
var f *os.File
var err1 error

func main() {

	resp, err := http.Get("https://tophub.today/n/mproPpoq6O")
	//异常时进行panic
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() //延时关闭

	//状态码错误时输出错误码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error ststus code:%d", resp.StatusCode)
	}
	result, err := ioutil.ReadAll(resp.Body) //借助io.reader获取流的信息
	content := result
	//原本的打印整个网页
	/*if err!=nil{
		panic(err)
	}
	fmt.Printf("%s",result)*/

	//paraseContent(result)//调用正则表达式
	//mail()//邮件发送
	//}

	//正则表达式
	//<td class="al"><a href="/l?e=d91bqTJRo70ETLH%2Fjpu2kEk8h%2B5Im5yC3efqsfIfpyauv1w8VPsLmGwwSsObDpmksAPqUt2cgAUzmlK34iiuNQLO3Xnbf791Ajpok73kl2NvtTmN%2FnTulJyjIXTTMhZvrMbGOSH9QTgtdgkf%2Bwj1" target="_blank" rel="nofollow" itemid="20668765">美国媒体测算，拜登已获得超过 270 张选举人票，特朗普还有可能翻盘吗？</a></td>
	//		                                    <td>9151 万热度</td>
	//func paraseContent(content []byte){
	//var wireteString = "测试n"

	//判断文件是否存在
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		fmt.Println("文件存在")
	} else {
		f, err1 = os.Create(filename) //创建文件
		fmt.Println("文件不存在")
	}
	//打印时间戳
	now := time.Now()
	check(err1)
	n, err1 := io.WriteString(f, "今年第"+strconv.Itoa(now.YearDay())+"期\n")
	re := regexp.MustCompile(`<td class="([^"]+)"><a href="([^"]+)" target="([^"]+)" rel="([^"]+)" itemid="([^"]+)">([^"]+)</a></td>
		                                    <td>([^"]+)</td>`)
	match := re.FindAllSubmatch(content, -1)
	for i, m := range match {
		j := strconv.Itoa(i + 1)
		check(err1)
		//写入文件(字符串)
		n, err1 = io.WriteString(f, "\n"+j)
		n, err1 = io.WriteString(f, ":")
		n, err1 = io.WriteString(f, string(m[6]))
		n, err1 = io.WriteString(f, string(m[7])+"\n")
		n, err1 = io.WriteString(f, "https://tophub.today"+string(m[2])+"\n")
		answerreach(string(m[2]))

		//resp,err:=http.Get("https://tophub.today"+string(m[2]))
		////异常时进行panic
		//if err!=nil{
		//	panic(err)
		//}
		//defer resp.Body.Close()//延时关闭
		//
		////状态码错误时输出错误码
		//if resp.StatusCode!=http.StatusOK{
		//	fmt.Printf("Error ststus code:%d",resp.StatusCode)
		//}
		//doc, err := goquery.NewDocumentFromReader(res.Body)
		//if err != nil {
		//	log.Fatal(err)
		//}
		check(err1)
		fmt.Printf("写入 %d 个字节n", n)
		//fmt.Printf("%d: %s\n%s\n\n", i+1,m[6],m[7])
		if i == 19 {
			break
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//访问知乎回答（来不及重构了）
func answerreach(a string) {
	res, err := http.Get("https://tophub.today" + a)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".QuestionHeader .QuestionHeader-content .QuestionHeader-main").Each(func(i int, s *goquery.Selection) {
		//questionTitle := s.Find(".QuestionHeader-title").Text()
		questionContent := s.Find(".QuestionHeader-detail").Text()
		//questionContent = questionContent[0 : len(questionContent)-12]

		//fmt.Println("questionTitle：", questionTitle)
		fmt.Println("questionContent：", questionContent)
	})

	doc.Find(".ContentItem-actions").Each(func(i int, s *goquery.Selection) {

	})
	doc.Find(".ListShortcut .List .List-item ").Each(func(i int, s *goquery.Selection) {
		//head_url, _ := s.Find("a img").Attr("src") //头像
		author := s.Find(".AuthorInfo-head").Text() //作者
		comma := strings.Index(author, ".")         //作者栏除杂
		if comma > 0 {
			author = author[:comma]
		}
		//fmt.Println("head_url：", head_url)
		//fmt.Println("author：", author)
		check(err1)
		_, err1 = io.WriteString(f, "作者："+author)
		//voters := s.Find(".Voters").Text()//赞同数
		//voters = strings.Split(voters, " ")[0]
		content := s.Find(".RichContent-inner").Text() //带标签的可以用Html()
		//createTime := s.Find(".ContentItem-time").Text()//创作日期
		//createTime = strings.Split(createTime, " ")[1]
		//commentCount := s.Find(".ContentItem-actions span").Text()
		//fmt.Println("voters：", voters)
		//fmt.Println("content："+"回答：\n"+ content)
		_, err1 = io.WriteString(f, "content："+"回答：\n"+content+"\n")
		//fmt.Println("createTime：", createTime)//
		//fmt.Println("commentCount : ", commentCount)
	})

}

/*
check(err1)
	n, err1 := io.WriteString(f, wireteString) //写入文件(字符串)
*/
/*
l, err := f.WriteString("Hello World")
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
*/
/*
 f, err := os.Create("test2.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err = f.Close(); err != nil {
            log.Fatal(err)
        }
    }()
    l, err := f.WriteString("Hello world")
    if err != nil {
        log.Fatal(err)
    }
*/
