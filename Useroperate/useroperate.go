package Useroperate

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"time"
)
var MysqlDB *sql.DB
type User struct {
	account string `json:"account" form:"account"`
	password []byte `json:"password" form:"password"`
	email string `json:"email" form:"email"`
}
type VC struct {
	verificationcode string
	emailaddress string
}
type Exuser struct {
	account string
	password []byte
	email string
}

//注册并加密密码，发送验证码以验证邮箱
func RegisterUser(c *gin.Context)  {
	vc := new(VC)
	accountinput := c.Query("account")
	passwordinput := c.Query("password")
	emailinput := c.Query("email")
	vcode := c.PostForm("vcode")
	//比对验证码
	selectedrow := MysqlDB.QueryRow("SELECT emailaddress, verificationcode FROM vcode WHERE emailaddress=?", emailinput)
	if err := selectedrow.Scan(&vc.emailaddress, &vc.verificationcode); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Please get your verification code first.",
		})
		return
	}
	if vcode != vc.verificationcode {
		fmt.Printf("Please input the right verification code.\n")
		return
	}
	//比对结束
	h := sha1.New()						//hash加密
	h.Write([]byte(passwordinput))
	pw := h.Sum(nil)					//hash加密结束
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
	var user Exuser
	h := sha1.New()						//hash加密
	h.Write([]byte(password))
	pw := h.Sum(nil)					//hash加密
	rows := MysqlDB.QueryRow("SELECT account, password, email FROM savedaccount WHERE account=?", account)
	if err := rows.Scan(&user.account, &user.password, &user.email); err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"can't find the account or your password is wrong. Please check the context you type in.",
		})
		return
	}
	if bytes.Equal(pw, user.password){
		//生成cookie
		expiration := time.Now()
		expiration = expiration.AddDate(0, 0, 1)	//将cookie的有效期设置为一天
		cookie := http.Cookie{Name: "Account", Value: account, Expires: expiration}
		http.SetCookie(c.Writer, &cookie)
		c.JSON(http.StatusOK, gin.H{
			"msg":"Login in successfully.",
		})
		c.JSON(http.StatusOK, &user)
	}else{
		c.JSON(http.StatusOK, gin.H{
			"msg":"Your password is wrong.",
		})
	}
}
//忘记密码，需先获得验证码，同时在输入url时post vcode 和 newpassword
func Forgetpassword(c *gin.Context){
	vc := new(VC)
	accountinput := c.Query("account")
	emailinput := c.Query("email")
	vcode := c.PostForm("vcode")
	newpassword := c.PostForm("newpassword")
	h := sha1.New()
	h.Write([]byte(newpassword))
	newpw := h.Sum(nil)
	//比对验证码
	selectedrow := MysqlDB.QueryRow("SELECT emailaddress, verificationcode FROM vcode WHERE emailaddress=?", emailinput)
	if err := selectedrow.Scan(&vc.emailaddress, &vc.verificationcode); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":"Please get your verification code first.",
		})
		return
	}
	if vcode != vc.verificationcode {
		fmt.Printf("Please input the right verification code.\n")
		return
	}
	//比对结束
	result, err := MysqlDB.Exec("UPDATE savedaccount SET password=? where account=?", newpw, accountinput)
	if err != nil{
		fmt.Printf("Something wrong when you change your password, err:%v\n", err)
		return
	}
	fmt.Println("Change password successd:", result)
}
//登出账户，需要已经登录
func Logout(c *gin.Context)  {
	expiration := time.Now()
	expiration = expiration.AddDate(0,0,-1)	//通过将有效期的时间调回1天前来使得cookie无效
	cookie := http.Cookie{Name: "Account", Value: "", Expires: expiration}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(http.StatusOK, gin.H{
		"msg":"退出账户成功",
	})
}