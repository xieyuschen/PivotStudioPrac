package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm@v1.9.16"
	"github.com/go-sql-driver/mysql@v1.5.0"
	"net/http"
)

var MysqlDB *gorm.DB

type User struct {
	account string `gorm:"size:30;primary_key;not null" json:"account"`		//gorm后为条件，json后为连接的表的字段
	password string `gorm:"size:30;not null" json:"password"`
}

/*func Tips(c *gin.Context)  {							提示函数感觉好像不用
	c.JSON(http.StatusOK, gin.H{
		"msg":"Please input your name and password.",
	})
}*/
func main(){
	//打开数据库，若错误返回错误
	MysqlDB, err := gorm.Open("mysql", "root:qzt0419ryf0416.@tcp(127.0.0.1:3306)")
	if err != nil{
		fmt.Println("failed to open database: ", err)
		return
	}else{
		fmt.Println("connect database success!")
		MysqlDB.SingularTable(true)			//使得默认表名和结构名相同
		MysqlDB.AutoMigrate(&User{})				//创建表名为user的表，约束与结构中相同
		fmt.Println("create table success!")
	}
	defer MysqlDB.Close()
	router := gin.Default()
	router.GET("/user", Login)					//登录
	router.POST("/user/create", RegisterUser)	//注册
	//router.GET("/user/find", Resetpassword)		//忘记密码
	router.Run(":8080")
}

//注册
func RegisterUser(c *gin.Context)  {
	var user User
	c.BindJSON(&user)				//若有错误则返回400
	MysqlDB.Create(&user)
	c.JSON(http.StatusOK, &user)
}

//登录
func Login(c *gin.Context)  {
	account := c.Query("account")
	var user User
	err := MysqlDB.First(&user, account)
	if err != nil{
		fmt.Println("can't find the account. Please check the account you type in. ERROR:", err.Error)
	}else {
		c.JSON(http.StatusOK, &user)
	}
}