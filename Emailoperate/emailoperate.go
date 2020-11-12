package Emailoperate

import (
	"PS_m1_ture/Useroperate"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

//邮箱
//邮件发送
func EmailSend(c *gin.Context)  {
	emailinput := c.Query("email")
	CodeDelete(emailinput)					//删除之前的验证码以重新生成
	auth := smtp.PlainAuth("", "2646677541@qq.com", "qjvqowrxrvbhecdi", "smtp.qq.com")
	to := []string{emailinput}
	nickname := "PIVOT STUDIO"
	user := "2646677541@qq.com"
	subject := "Email address verification code."
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := CodeCreate(8)			//生成一个验证码
	CodeSave(body, emailinput)
	msg := []byte("To:" + strings.Join(to, ",") + "\r\nFrom: " + nickname + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.qq.com:25", auth, user, to, msg)
	if err != nil{
		fmt.Printf("Send email error: %v\n", err)
		return
	}
	fmt.Printf("Send email successful.\n")
}
//生成验证码,长度为8由上面注册程序里决定
func CodeCreate(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")	//rune的话下次可以把验证码改成中文，不过表里面的数据类型也要改就是了
	rand.Seed(time.Now().UnixNano())
	verification := make([]rune, n)
	for i := range verification{
		verification[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(verification)
}
//将验证码存入数据库用以比对
func CodeSave(verification string, emailaddress string)  {
	stmt, err := Useroperate.MysqlDB.Prepare("INSERT  INTO vcode SET emailaddress=?,verificationcode=?")
	if err != nil{
		fmt.Printf("Create verification code failed, err:%v\n", err)
	}else{
		_, err := stmt.Exec(emailaddress, verification)
		if err != nil{
			fmt.Printf("Create verification code failed, err:%v\n", err)
		}
		fmt.Println("Create verification successful.")
	}
}
//删除数据表中数据，重新发送验证码时使用
func CodeDelete(emailaddress string)  {
	Useroperate.MysqlDB.Exec("delete from vcode where emailaddress=?", emailaddress)
}
