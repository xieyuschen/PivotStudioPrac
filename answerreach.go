package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"strings"
)

//访问知乎回答（来不及重构了）
func Answerreach(a string) {
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
		questionTitle, _ := s.Find(".QuestionHeader-title").Html()
		questionContent, _ := s.Find(".QuestionHeader-detail").Html()
		//questionContent = questionContent[0 : len(questionContent)-12]
		check(err1)
		//fmt.Println("questionTitle：", questionTitle)
		//fmt.Println("questionContent：", questionContent)
		_, err1 = io.WriteString(f, questionTitle)
		_, err1 = io.WriteString(f, questionContent)
	})

	doc.Find(".ContentItem-actions").Each(func(i int, s *goquery.Selection) {

	})
	doc.Find(".ListShortcut .List .List-item ").Each(func(i int, s *goquery.Selection) {
		//head_url, _ := s.Find("a img").Attr("src")     //头像
		author := s.Find(".AuthorInfo-head").Text() //作者
		comma := strings.Index(author, ".")         //作者栏除杂
		check(err1)
		if comma > 0 {
			author = author[:comma]
		}
		//fmt.Println("head_url：", head_url)
		//fmt.Println("author：", author)
		//_, err1= io.WriteString(f, head_url)

		_, err1 = io.WriteString(f, "<p><b>作者："+author+"</b></p>")
		//<p>第六：拍照写实；</p>
		//<b>越来越能感觉到：</b>
		//voters := s.Find(".Voters").Text()//赞同数
		//voters = strings.Split(voters, " ")[0]
		content, _ := s.Find(".RichContent-inner").Html() //带标签的可以用Html()
		//createTime := s.Find(".ContentItem-time").Text()//创作日期
		//createTime = strings.Split(createTime, " ")[1]
		//commentCount := s.Find(".ContentItem-actions span").Text()
		//fmt.Println("voters：", voters)
		//fmt.Println("content："+"回答：\n"+ content)
		_, err1 = io.WriteString(f, "<p><b>回答：</b></p>"+content+"<p></p>")
		//fmt.Println("createTime：", createTime)//
		//fmt.Println("commentCount : ", commentCount)
	})

}
