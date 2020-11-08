package routes

import (
	v1 "ginblog/api/v1"
	"ginblog/middleware"
	"ginblog/utils"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	//gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	logInRouter := r.Group("api/v1")
	//登陆后才能访问
	logInRouter.Use(middleware.JwtToken() )
	{
		//User路由接口
		logInRouter.GET("users",v1.GetUsers)
		logInRouter.PUT("user/:id",v1.EditUser)
		logInRouter.DELETE("user/:id",v1.DeleteUser)
		//Article路由接口
		logInRouter.POST("article/add",v1.AddArticle)
		logInRouter.GET("article",v1.GetArtInfo)
		logInRouter.PUT("article/:id",v1.EditArt)
		logInRouter.DELETE("article/:id",v1.DeleteArt)
	}
	pubRouter:=r.Group("api/v1")
	{
		pubRouter.POST("user/add",v1.AddUser)
		pubRouter.POST("login",v1.Login)
	}
	_ = r.Run(utils.HttpPort)
}
