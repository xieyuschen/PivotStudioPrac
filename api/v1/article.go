package v1

import (
	"ginblog/model"
	"ginblog/utils/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//查找文章
func GetArtInfo(c gin.Context){
	id,_:=strconv.Atoi(c.Param("id"))
	data,code := model.GetArtInfo(id)
	c.JSON(http.StatusOK ,gin.H{
		"status":code,
		"data":data,
		"message":errmsg.GetErrMsg(code),
	})
}
//新建文章
func AddArticle(c *gin.Context ){
	var data model.Article
	_ = c.ShouldBindJSON(&data)

	code := model.CreateArt(&data)

	c.JSON(http.StatusOK,gin.H{
		"status":code,
		"data":data,
		"message":errmsg.GetErrMsg(code),
	})
}

//编辑文章
func EditArt(c *gin.Context){
	var data model.Article
	id,_:= strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.EditArt(id,&data)

	c.JSON(http.StatusOK ,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}

//删除文章
func DeleteArt(c *gin.Context){
	id,_:=strconv.Atoi(c.Param("id") )

	code := model.DeleteArt(id)

	c.JSON(http.StatusOK ,gin.H{
		"status":code,
		"message":errmsg.GetErrMsg(code),
	})
}