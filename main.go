package main
import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
	"time"
)
var MysqlDB *sql.DB
type User struct {
	account string `json:"account" form:"account"`		//gorm后为条件，json后为连接的表的字段
	password []byte `gorm:"size:100;not null" json:"password" form:"password"`
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
	router.GET("/login",Login)						//登录
	router.POST("/user/create", RegisterUser)		//注册
	//router.GET("/user/forgetpassword", Forgetpassword)	//忘记密码,和邮箱相结合
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
//注册并加密密码
func RegisterUser(c *gin.Context)  {
	accountinput := c.Query("account")
	passwordinput := c.Query("password")
	h := sha1.New()						//hash加密
	h.Write([]byte(passwordinput))
	pw := h.Sum(nil)					//hash加密
	stmt, err := MysqlDB.Prepare("INSERT INTO savedaccount SET account=?, password=?")
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Already have this account.",
		})
	}else {
		_, err := stmt.Exec(accountinput, pw)
		if err != nil{
			c.JSON(http.StatusOK, gin.H{
				"msg": "You create the account successfully.",
			})
		}
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
/*//忘记密码配合邮箱，暂时不会
func Forgetpassword(c *gin.Context){
}*/
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