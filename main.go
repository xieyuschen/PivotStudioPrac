package main
import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"net/smtp"
)
var MysqlDB *sql.DB
type User struct {
	account string `json:"account" form:"account"`		//gorm后为条件，json后为连接的表的字段
	password []byte `json:"password" form:"password"`
	email string `json:"email" form:"email"`
}
type Article struct {
	Author string			//用于存储发帖者
	//Others string			//如果做的话，用于评论者评论，若为空则代表为发帖者发的原文，否则代表评论者
	Title string			//用于存储题目
	Summary string			//用于存储内容概要
	Content string			//用于存储文章主体内容
}
//在开始前已经创建了两个数据表savedaccount & articles，一个用于存储用户账户密码数据，一个用于保存发的帖子的数据
func main(){
	//打开数据库，若错误返回错误
	MysqlDB, err := sql.Open("mysql", "root:qzt0419ryf0416.@tcp(127.0.0.1:3306)/accountpassword?charset=utf8")
	if err != nil{
		fmt.Println("failed to open database: ", err)
		return
	}else{
		fmt.Println("connect database success!")
		//MysqlDB.SingularTable(true)			//使得默认表名和结构名相同
		//MysqlDB.AutoMigrate(&User{})				//创建表名为user的表，约束与结构中相同
		//fmt.Println("create table success!")
	}
	defer MysqlDB.Close()
	router := gin.Default()
	router.GET("/emailcheck",EmailSend)
	router.GET("/login",Login)						//登录
	router.POST("/user/create", RegisterUser)		//注册
	router.GET("/user/forgetpassword", Forgetpassword)	//忘记密码,和邮箱相结合
	auth := router.Group("")
	auth.Use(AuthRequired())
	{
		//在登录后可以运行博客的相关操作
		auth.GET("/logout", Logout)
		auth.GET("/writearticle", WriteArticle)
		auth.GET("/revisearticle", ReviseArticle)
		auth.GET("/seearticles", SeeArticles)
	}
	router.Run(":8080")
}
//注册并加密密码；添加注册邮件服务，发送验证码以注册账户
func RegisterUser(c *gin.Context)  {
	accountinput := c.Query("account")
	passwordinput := c.Query("password")
	emailinput := c.Query("email")
	vcode := c.PostForm("vcode")

	h := sha1.New()						//hash加密
	h.Write([]byte(passwordinput))
	pw := h.Sum(nil)					//hash加密
	stmt, err := MysqlDB.Prepare("INSERT INTO savedaccount SET account=?, password=?, email=?")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Already have this account or email address.",
		})
	}else {
		_, err := stmt.Exec(accountinput, pw, emailinput)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Failed to  create the account.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": "You create the account successfully.",
		})
	}
}
//登录
func Login(c *gin.Context)  {
	account := c.Query("account")
	password := c.Query("password")
	var user User
	h := sha1.New()						//hash加密
	h.Write([]byte(password))
	pw := h.Sum(nil)					//hash加密
	user.password = pw
	user.account = account
	rows, err := MysqlDB.Query("SELECT * FROM savedaccount")
	for rows.Next(){
		err = rows.Scan(&account, &password)
		if err != nil{
			fmt.Println("can't find the account or your password is wrong. Please check the context you type in. ERROR:", err.Error())
		}else{
			//生成cookie
			expiration := time.Now()
			expiration = expiration.AddDate(0, 0, 1)	//将cookie的有效期设置为一天
			cookie := http.Cookie{Name: "Account", Value: account, Expires: expiration}
			http.SetCookie(c.Writer, &cookie)
			fmt.Println("Login in successfully.")
			c.JSON(http.StatusOK, &user)
		}
	}
}
//忘记密码配合邮箱，暂时不会
func Forgetpassword(c *gin.Context){
	EmailSend(c)

}
//cookie相关
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("account")
		if cookie == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "请先登陆",
			})
			c.Abort()					//阻止调用被挂起的函数,使得Next无法执行
			return
		}
		c.Next()	//调用下一个被挂起的函数
	}
}
func Logout(c *gin.Context)  {
	expiration := time.Now()
	expiration = expiration.AddDate(0,0,-1)	//通过将有效期的时间调回1天前来使得cookie无效
	cookie := http.Cookie{Name: "Account", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(http.StatusOK, gin.H{
		"msg":"退出账户成功",
	})
}
//帖子相关操作
//创建帖子，并写入
func WriteArticle(c *gin.Context)  {
	titleinput := c.Query("title")
	authorinput,_ := c.Request.Cookie("account")
	summaryinput := c.Query("summary")
	contentinput := c.Query("content")
	stmt, err := MysqlDB.Prepare("INSERT INTO articles SET author=?, title=?, summary=?, content=?")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Error in writing an article",
		})
	}else {
		_, err := stmt.Exec(authorinput, titleinput, summaryinput, contentinput)
		if err != nil{
			c.JSON(http.StatusOK, gin.H{
				"msg": "You write the article successfully.",
			})
		}
	}
}
//修改帖子content和summary，title(主键)和author不可改
func ReviseArticle(c *gin.Context)  {
	titleinput := c.Query("title")
	summaryinput := c.Query("summary")
	contentinput := c.Query("content")
	result, err := MysqlDB.Exec("UPDATE articles SET summary=?, content=? where title=?", summaryinput, contentinput, titleinput)
	if err != nil{
		fmt.Printf("Revise failed, err:%v\n", err)
		return
	}
	fmt.Println("Revise article successd:", result)
	rowsaffected, err := result.RowsAffected()
	if err != nil{
		fmt.Printf("Get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}
//删除所发的帖子, 需要表中有该title的article且当前账户为创建者账户有权限删除
func DeleteArticle(c *gin.Context)  {
	A := new(Article)
	titleinput := c.Query("title")
	authorinput,_ := c.Request.Cookie("account")
	row := MysqlDB.QueryRow("select author, title from articles where author=?,title=?", authorinput,titleinput)
	if err := row.Scan(&A.Author, &A.Title); err != nil{
		fmt.Printf("Scan failed. You don't have the privilege or the table does not have this article. err:%v\n", err)
		return
	}
	result, err := MysqlDB.Exec("delete from articles where title=?", titleinput)
	if err != nil{
		fmt.Printf("Delete failed, err:%v\n", err)
		return
	}
	fmt.Println("Delete article successd", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil{
		fmt.Printf("Get RowsAffected failed, err:%v\n",err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
}
//查看自己的帖子
func SeeArticles(c *gin.Context)  {
	A := new(Article)
	authorinput,_ := c.Request.Cookie("account")
	titleinput := c.Query("title")
	rows, err := MysqlDB.Query("select author, title, summary, content from articles where author=?, title=?", authorinput, titleinput)
	defer func() {
		if rows != nil{
			rows.Close()	//关掉未scan的sql连接
		}
	}()
	if err != nil{
		fmt.Printf("Can't see it, err:%v/n", err)
		return
	}
	for rows.Next(){
		err = rows.Scan(&A.Author, &A.Title, &A.Summary, &A.Content)
		if err != nil{
			fmt.Printf("Scan failed, err:%v\n", err)
			return
		}
		fmt.Println("Scan successd:", *A)
	}
}
//邮箱
//邮件发送
func EmailSend(c *gin.Context)  {
	emailinput := c.Query("email")
	CodeDelete(emailinput)					//删除之前的验证码以重新生成
	auth := smtp.PlainAuth("", "hustqiuzt@gmail.com", "qzt0419ryf0416.", "smtp.gmail.com")
	to := []string{emailinput}
	nickname := "PIVOT STUDIO"
	user := "hustqiuzt@gmail.com"
	subject := "Email address verification code."
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := CodeCreate(8)			//生成一个验证码
	CodeSave(body, emailinput)
	msg := []byte("To:" + strings.Join(to, ",") + "\r\nFrom: " + nickname + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail("smtp.gmail.com:25", auth, user, to, msg)
	if err != nil{
		fmt.Printf("Send mail error: %v\n", err)
	}
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
	stmt, err := MysqlDB.Prepare("INSERT  INTO vcode SET emailaddress=?,verificationcode=?")
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
	_, err := MysqlDB.Exec("delete from vcode where emailaddress=?", emailaddress)
	if err != nil{
		return
	}
}