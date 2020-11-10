package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"net/http"
)

var MysqlDB *sql.DB

type User struct {
	account string `gorm:"size:30;primary_key;not null" json:"account"`		//gorm后为条件，json后为连接的表的字段
	password []byte `gorm:"size:100;not null" json:"password"`
}

/*func Tips(c *gin.Context)  {							提示函数感觉好像不用
	c.JSON(http.StatusOK, gin.H{
		"msg":"Please input your name and password.",
	})
}*/
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
	auth := router.Group("")
	auth.Use(CookieNeed()){
		router.GET("/user", Login)					//登录
		router.POST("/user/create", RegisterUser)	//注册
		//router.GET("/user/forgetpassword", Forgetpassword)	//忘记密码,和邮箱相结合
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
			fmt.Println("can't find the account. Please check the account you type in. ERROR:", err.Error)
		}else{
			fmt.Println("Login in successfully.")
			c.JSON(http.StatusOK, &user)
		}
	}
}

/*//忘记密码配合邮箱，暂时不会
func Forgetpassword(c *gin.Context){

}*/

//检验cookie
func CookieNeed() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("account")
		if cookie == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请先登陆",
			})
			c.Abort()					//若没成功阻止调用其他被挂起的函数，但不能阻止当前函数
		}
		c.Next()	//调用下一个被挂起的函数
	}
}