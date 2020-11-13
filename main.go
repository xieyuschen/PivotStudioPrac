package main
import (
	"PS_m1_ture/Articleoperate"
	"PS_m1_ture/Cookieoperate"
	"PS_m1_ture/Emailoperate"
	"PS_m1_ture/Useroperate"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)
//在开始前已经创建了三个数据表savedaccount & articles & vcode，一个用于存储用户账户密码数据，一个用于保存发的帖子的数据, 一个用于存储验证码
func main(){
	//打开数据库，若错误返回错误
	var err error
	Useroperate.MysqlDB, err = sql.Open("mysql", "root:qzt0419ryf0416.@tcp(127.0.0.1:3306)/accountpassword?charset=utf8")
	if err != nil{
		fmt.Println("failed to open database: ", err)
		return
	}else{
		fmt.Println("connect database success!")
	}
	defer Useroperate.MysqlDB.Close()
	router := gin.Default()
	router.GET("/emailcheck",Emailoperate.EmailSend)				//发送验证码
	router.GET("/login",Useroperate.Login)						//登录账号，需先注册
	router.POST("/user/create", Useroperate.RegisterUser)		//注册账号，需先发送验证码
	router.POST("/user/forgetpassword", Useroperate.Forgetpassword)	//忘记密码,需先发送验证码
	auth := router.Group("")
	auth.Use(Cookieoperate.AuthRequired())
	{
		//在登录后可以运行博客的相关操作
		auth.GET("/logout", Useroperate.Logout)					//注销账户
		auth.GET("/writearticle", Articleoperate.WriteArticle)	//写博客
		auth.GET("/revisearticle", Articleoperate.ReviseArticle)	//修改博客
		auth.GET("/seearticles", Articleoperate.SeeArticles)		//查看所有自己的博客
		auth.GET("/deletearticle", Articleoperate.DeleteArticle)	//删除自已的一篇指定博客
	}
	router.Run(":8080")
}

