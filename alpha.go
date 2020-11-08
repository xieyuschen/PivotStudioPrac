package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp" //正则表达式
)

func main(){

	resp,err:=http.Get("https://tophub.today/n/mproPpoq6O")
	//异常时进行panic
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()//延时关闭

	//状态码错误时输出错误码
	if resp.StatusCode!=http.StatusOK{
		fmt.Printf("Error ststus code:%d",resp.StatusCode)
	}
	result,err:=ioutil.ReadAll(resp.Body)//借助io.reader获取流的信息
	//原本的打印整个网页
	/*if err!=nil{
		panic(err)
	}
	fmt.Printf("%s",result)*/
	paraseContent(result)//调用正则表达式
}
//正则表达式
//<td class="al"><a href="/l?e=d91bqTJRo70ETLH%2Fjpu2kEk8h%2B5Im5yC3efqsfIfpyauv1w8VPsLmGwwSsObDpmksAPqUt2cgAUzmlK34iiuNQLO3Xnbf791Ajpok73kl2NvtTmN%2FnTulJyjIXTTMhZvrMbGOSH9QTgtdgkf%2Bwj1" target="_blank" rel="nofollow" itemid="20668765">美国媒体测算，拜登已获得超过 270 张选举人票，特朗普还有可能翻盘吗？</a></td>
//		                                    <td>9151 万热度</td>
func paraseContent(content []byte){
	re:=regexp.MustCompile(`<td class="([^"]+)"><a href="([^"]+)" target="([^"]+)" rel="([^"]+)" itemid="([^"]+)">([^"]+)</a></td>
		                                    <td>([^"]+)</td>`)
	match:=re.FindAllSubmatch(content,-1)
	for i,m:=range match {
		fmt.Printf("%d: ", i+1)
		fmt.Printf("%s \n", m[6])
		fmt.Printf("%s \n\n", m[7])
		if(i==9){
		break;
		}
	}
	}