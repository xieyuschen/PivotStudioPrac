# 爬虫笔记

## 20201103

​																																															from AD钙奶 SBY

### 适用于UTF-8（爬了自己的网站）0.0.1

```go
package main

import (
   "fmt"
   "io/ioutil"
   "net/http"

)

func main(){

   resp,err:=http.Get("http://www.shouxindehai.cn/")
   //异常时进行panic
   if err!=nil{
      panic(err)
   }
   //延时关闭
   defer resp.Body.Close()

   //状态码错误时输出错误码
   if resp.StatusCode!=http.StatusOK{
      fmt.Printf("Error ststus code:%d",resp.StatusCode)
   }

   //借助io.reader获取流的信息
   result,err:=ioutil.ReadAll(resp.Body)
   if err!=nil{
      panic(err)
   }
   fmt.Printf("%s",result)
   //href<我理解为链接 有这个标签就会有子层的目录 比方说http://www.shouxindehai.cn/index.php/cs/>
   //在这里可以在终端搜索这个目录发现可以找到 说明是可以爬到的
   //该段代码适用于utf-8的编码的网址
   //gbk会乱码的


}
```

### 通用版（要去下载安装）

第二集里有镜像站下载地址（有时间了再优化吧）

https://www.bilibili.com/video/BV1XK411V7DT?p=2

```go
package main

import (
   "fmt"
   "io/ioutil"
   "net/http"
)

func main(){

   resp,err:=http.Get("http://www.shouxindehai.cn/")
   //异常时进行panic
   if err!=nil{
      panic(err)
   }
   //延时关闭
   defer resp.Body.Close()
##
   //状态码错误时输出错误码
   //if resp.StatusCode!=http.StatusOK{
   // fmt.Printf("Error ststus code:%d",resp.StatusCode)
   //}
   //注释部分都是懒得下载第三方包 这一部分是调用后面的函数来实现对编码方式的检查
   //bodyReader:=bufio.NewReader(resp.Body)
   //e:=determinEncoding(bodyReader)
   ////这里实现编码转换
   //utf8Reader:=transform.NewReader(bodyReader,e.NewDecoder())
   //借助io.reader获取流的信息
   //result,err:=ioutil.ReadAll(utf8Reader)假如发生编码转换则需要下一行改参数名
##
   result,err:=ioutil.ReadAll(resp.Body)
   if err!=nil{
      panic(err)
   }
   fmt.Printf("%s",result)
   //href<我理解为链接 有这个标签就会有子层的目录 比方说http://www.shouxindehai.cn/index.php/cs/>
   //在这里可以在终端搜索这个目录发现可以找到 说明是可以爬到的
   //该段代码适用于utf-8的编码的网址
   //gbk会乱码的
   //写到这里发现知乎有反扒

}
##
//这还得用第三方库 还得翻墙 那真麻烦 那就算了吧 代码留在这里好了
//设计一个检查编码方式的函数（知乎用的是utf-8所以用不上貌似 先注释掉备份吧）
// func derterminEncoding(r*bufio.Reader)encoding.Encoding{
//    bytes,err:=r.Peek(n:1024)
//    if err1=nil{
//       log.Printf("fetch error:%v",err)
//       return unicode.UTF8
//    }
//
//    e,_,_:=charset.DetermineEncoding(bytes,contentType:"")
//    return e
// }
##
```

井号部分是需要修改的部分（改完后可以实现通用爬）