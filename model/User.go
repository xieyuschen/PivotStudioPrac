package model

import (
	"encoding/base64"
	"fmt"
	"ginblog/gomail"
	"ginblog/utils/errmsg"
	"golang.org/x/crypto/scrypt"

	"gorm.io/gorm"
	"log"
	"math/rand"
	"strings"
	"time"
)

//用户
type User struct{
	gorm.Model
	Username string
	Password string
	Email string
	Role int//1为临时用户，不能进行登录操作，0为正式用户
	ValidCode string//储存的标准验证码
	InputCode string//用户输入的验证码
}
//查找用户是否存在
func CheckUser(name string)(code int){
	var users User
	db.Select("id").Where("username = ?",name).First(&users)
	if users.ID >0{
		return errmsg.ERROR_USERNAME_USED //1001
	}
	return errmsg.SUCCSE
}


//检查验证码是否正确
func CheckVcode(name string,tempUser User) int{
	var user User
	code := errmsg.SUCCSE //错误/正确码
	db.Select("id").Where("username = ?",name).First(&user)
	//核对信息
	//这里对临时用户与注册时的核对还有一点小问题
	//离散还没预习完时间不大够也不知道能不能完成qwq
	/*if user.ID <=0{
		code = errmsg.ERROR_USER_EXIST
		return code
	}
	if user.Role!=1{
		code=errmsg.ERROR_TUSER_ROLE
		return code
	}
	if user.Email!=tempUser.Email {
		code = errmsg.ERROR_MAIL_WRONG
		return code
	}*/
	if user.InputCode !=tempUser.ValidCode{
		code=errmsg.ERROR_VALIDCODE
		return code
	}
	return code
}

//生成4位随机验证码
func GenerateVcode()string{
	n := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(n)
	rand.Seed(time.Now().UnixNano())

	var s strings.Builder
	for i:=0;i<4;i++{
		fmt.Fprintf(&s,"%d",n[rand.Intn(r)])
	}
	return s.String()
}


//将验证码通过邮件发送给用户
func SendEmail(emailAdress string,vcode string)error{

	m := gomail.NewMessage()
	m.SetHeader("From", "1343244602@qq.com")
	m.SetHeader("To", emailAdress)
	m.SetHeader("Subject", "ginBlog register verify code")
	m.SetBody("text/html", vcode )

	d := gomail.NewDialer("smtp.qq.com", 465, "1343244602@qq.com", "xrctskufxuokgjeh")
	//邮件发送服务器信息,使用授权码而非密码
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
//创建临时用户，仅存储用户名、及EmailAdress,ValidCode
func CreateTempUser(data *User,vcode string)int{
	data.ValidCode=vcode
	data.Role=1//role==1为临时用户
	err :=db.Create(&data).Error
	if err!=nil{
		return errmsg.ERROR  //500
	}
	return errmsg.SUCCSE //200
}
//注册正式用户
func CreateUser(data *User)int{
	data.Password = ScryptPW(data.Password)//密码加密存放
	data.Role =0//正式用户
	err :=db.Create(&data).Error
	if err!=nil{
		return errmsg.ERROR  //500
	}
	return errmsg.SUCCSE //200
}

//查询用户列表
func GetUsers(pageSize int,pageNum int)[]User{
	var users []User
	err :=db.Limit(pageSize).Offset((pageNum-1)*pageSize).Find(&users).Error
	if err!=nil && err!=gorm.ErrRecordNotFound {
		return nil
	}
	return users
}
//密码加密
func ScryptPW(password string)string{
	const KeyLen = 10
	salt := make([]byte,8)
	salt =[]byte{11,4,5,14,33,22,3,56}

	HashPw,err:=scrypt.Key([]byte(password),salt,16384,8,1,KeyLen)
	if err!=nil{
		log.Fatal(err)
	}
	finalPw := base64.StdEncoding.EncodeToString(HashPw)

	return finalPw
}

//编辑用户信息
func EditUser(id int,data *User)int{
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err := db.Model(&user).Where("id = ?",id).Updates(maps).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//删除用户
func DeleteUser(id int)int{
	var user User
	err := db.Where("id = ? ",id).Delete(&user).Error
	if err!=nil{
		return errmsg.ERROR_DELETEUSER_ERROR
	}
	return errmsg.SUCCSE
}
//
func DeleteUserName(name string)int{
	var user User
	err:=db.Where("username = ?",name).Delete(&user).Error
	if err!=nil{
		return errmsg.ERROR_DELETEUSER_ERROR
	}
	return errmsg.SUCCSE
}
//登录验证
func CheckLogin(username string,password string)int{
	var user User
	db.Where("username = ?",username).First(&user)
	if user.ID ==0{
		return errmsg.ERROR_USER_NOT_EXIST
	}
	if ScryptPW(password)!=user.Password{
		return errmsg.ERROR_PASSWORD_WRONG
	}
	if user.Role !=0/*无管理权限*/{
		return errmsg.ERROR_USER_NOT_RIGHT
	}
	return errmsg.SUCCSE
}