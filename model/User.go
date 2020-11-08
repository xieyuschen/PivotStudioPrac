package model

import (
	"encoding/base64"
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
	"golang.org/x/crypto/scrypt"
	"log"
)

//用户
type User struct{
	gorm.Model
	Username string
	Password string
	Role int
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


//注册用户
func CreateUser(data *User)int{
	data.Password = ScryptPW(data.Password)//加密
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
		return errmsg.ERROR
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