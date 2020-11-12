package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetVcode(c *gin.Context){
	var data model.User
	_ = c.ShouldBindJSON(&data)//接受请求
	//检查用户名是否重复
	code:= model.CheckUser(data.Username)
	//生成验证码并储存，通过邮件发送给用户
	if code==errmsg.SUCCSE {
		vcode := model.GenerateVcode()
		errmail := model.SendEmail(data.Email, vcode)
		if  errmail != nil {
			code = errmsg.ERROR_EMAILSEND_FAIL
		}
		if code==errmsg.SUCCSE {
			code =model.CreateTempUser(&data)//创建临时用户储存用户名及邮件地址正确验证码
			if code==errmsg.ERROR {
				c.JSON(http.StatusOK, gin.H{
					"status":  code,
					"data":    "create tempUser failed...",
					"message": errmsg.GetErrMsg(code),
				})
			}
		}
	}
	if code==errmsg.SUCCSE {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    "please check your emailbox for the valid-code",
			"message": errmsg.GetErrMsg(code),
		})
	}else{
		c.JSON(http.StatusOK,gin.H{
			"status":code,
			"data":"something went wrong,please check...",
			"message":errmsg.GetErrMsg(code),
		})
	}
}
// 添加用户，需要输入用户名，密码，验证码，邮箱
func AddUser(c *gin.Context ){
	var data model.User
	_=c.ShouldBindJSON(&data)
	//检查验证码是否正确
	code:= model.CheckVcode(data)

	model.DeleteUserName(data.Username) // 删除临时用户
	if code==errmsg.SUCCSE {
		model.CreateUser(&data)
		c.JSON(http.StatusOK ,gin.H{
			"status":code,
			"data":data,
			"message":errmsg.GetErrMsg(code),
		})
	}else{
		c.JSON(http.StatusOK ,gin.H{
			"status":code,
			//"data":data,
			"data":"something went wrong,please regist again...",
			"message":errmsg.GetErrMsg(code),
		})
	}

}

//删除用户
func DeleteUser(c *gin.Context ){
	id,_:= strconv.Atoi(c.Param("id"))


	code := model.DeleteUser(id)

	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

//查找用户列表：
func GetUsers(c *gin.Context ){
	pageSize,_ := strconv.Atoi(c.Query("pagesize"))
	pageNum,_ := strconv.Atoi(c.Query("pagenum"))

	if pageSize ==0{
		pageSize=-1//不分页
	}
	if pageNum==0{
		pageNum=-1
	}

	data := model.GetUsers(pageSize,pageNum)
	code := errmsg.SUCCSE
	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"data":data,
		"message":errmsg.GetErrMsg(code),
	})
}
//编辑用户
func EditUser(c *gin.Context ){
	var data model.User
	id,_:=strconv.Atoi(c.Param("id"))
	c.ShouldBindJSON(&data)
	code := model.CheckUser(data.Username)

	if code==errmsg.SUCCSE {
		model.EditUser(id,&data)
	}

	if code==errmsg.ERROR_USERNAME_USED {
		c.Abort()
	}

	c.JSON(http.StatusOK ,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

