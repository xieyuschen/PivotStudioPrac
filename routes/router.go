package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	router := r.Group("api/v1")
	{
		//User路由接口
		router.POST("user/add",v1.AddUser)
		router.GET("users",v1.GetUsers)
		router.PUT("user/:id",v1.EditUser)
		router.DELETE("user/:id",v1.DeleteUser)
	}
	_ = r.Run(utils.HttpPort)
}
