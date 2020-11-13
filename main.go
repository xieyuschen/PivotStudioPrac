//将热榜前20写入txt
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp" //正则表达式
	"strconv"
	"time"
)

var now = time.Now()
var filename = "./今日热榜" + strconv.Itoa(now.YearDay()) + ".html"
var f *os.File
var err1 error

//简易定时器
func main() {
	do()
	for {
		time.Sleep(time.Hour * 24)
		do()
	}
}

func do() {
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
	//http ex：
	//<p>第六：拍照写实；</p>//换行
	//<b>越来越能感觉到：</b>//加粗
	//打印时间戳
	now := time.Now()
	check(err1)
	n, err1 := io.WriteString(f, "<b><p>今年第"+strconv.Itoa(now.YearDay())+"期</p></b>")
	re := regexp.MustCompile(`<td class="([^"]+)"><a href="([^"]+)" target="([^"]+)" rel="([^"]+)" itemid="([^"]+)">([^"]+)</a></td>
		                                    <td>([^"]+)</td>`)
	match := re.FindAllSubmatch(content, -1)
	for i, m := range match {
		j := strconv.Itoa(i + 1)
		check(err1)
		//写入文件(字符串)
		n, err1 = io.WriteString(f, "<b><p>"+j+":</p></b>")
		n, err1 = io.WriteString(f, "<b><p>"+string(m[6]))
		n, err1 = io.WriteString(f, string(m[7])+"</p></b>")
		n, err1 = io.WriteString(f, "<b><p>"+"https://tophub.today"+string(m[2])+"</p></b>")
		Answerreach(string(m[2]))
		check(err1)
		fmt.Printf("写入 %d 个字节\n", n)
		//fmt.Printf("%d: %s\n%s\n\n", i+1,m[6],m[7])
		if i == 19 {
			break
		}
	}
	Mail() //邮件发送
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
