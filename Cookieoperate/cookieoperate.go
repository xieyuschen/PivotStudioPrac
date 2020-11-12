package Cookieoperate

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//cookie相关
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, _ := c.Request.Cookie("Account")
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
