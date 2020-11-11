package main

import (
	"gopkg.in/gomail.v2"
)

func Mail() {
	// 初始化
	m := gomail.NewMessage()

	// 发邮件的地址
	m.SetHeader("From", "1259585247@qq.com")

	// 给谁发送，支持多个账号
	m.SetHeader("To", "1259585247@qq.com")

	// 抄送谁
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")

	// 邮件标题
	m.SetHeader("Subject", "今日知乎热榜")

	// 邮件正文，支持 html
	m.SetBody("text/html", "今日知乎热榜请查收")

	// 附件
	m.Attach("今日热榜.html")

	// stmp服务，端口号，发送邮件账号，发送账号密码
	d := gomail.NewDialer("smtp.qq.com", 587, "1259585247@qq.com", "srnwcnffvqjyfjeg")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
