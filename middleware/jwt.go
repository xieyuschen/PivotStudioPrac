package middleware

import (
	"dgrijalva/jwt-go"
	"ginblog/utils"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var JwtKey=[]byte(utils.JwtKey)

type MyClaims struct{
	Username string `json:"username"`
	jwt.StandardClaims
}

//生成token
func SetToken(username string)(string,int){
	expireTime := time.Now().Add(10*time.Hour)
	SetClaims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer: "ginblog",
		} ,
	}

	reqClaim := jwt.NewWithClaims(jwt.SigningMethodHS256 ,SetClaims )
	token,err := reqClaim.SignedString(JwtKey )
	if err!=nil{
		return "",errmsg.ERROR
	}
	return token,errmsg.SUCCSE
}

//验证token
func CheckToken(token string)(*MyClaims ,int){
	setToken,_ := jwt.ParseWithClaims(token,&MyClaims{},func(token *jwt.Token)(interface{},error){
		return JwtKey,nil
	})

	if key,code:=setToken.Claims.(*MyClaims);code&&setToken.Valid{
		return key,errmsg.SUCCSE
	}else{
		return nil,errmsg.ERROR
	}
}
//jwt中间件
func JwtToken()gin.HandlerFunc {
	return func(c *gin.Context){
		tokenHerder := c.Request.Header.Get("Authorization")
		code :=errmsg.SUCCSE
		if tokenHerder ==""{
			code = errmsg.ERROR_TOKEN_EXIST
			c.JSON(http.StatusOK ,gin.H{
				"code":code,
				"message":errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHerder, " ",2)
		if len(checkToken)!=2 && checkToken[0]!="Bearer"{
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK ,gin.H{
				"code":code,
				"message":errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		key,tcode := CheckToken(checkToken[1])
		if tcode==errmsg.ERROR {
			code=errmsg.ERROR_TOKEN_WRONG
			c.JSON(http.StatusOK ,gin.H{
				"code":code,
				"message":errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		if time.Now().Unix()>key.ExpiresAt {
			code = errmsg.ERROR_TOKEN_RUNTIME
			c.JSON(http.StatusOK ,gin.H{
				"code":code,
				"message":errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		c.Set("username",key.Username )
		c.Next()
	}
}
