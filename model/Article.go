package model

import (
	"ginblog/utils/errmsg"
	"gorm.io/gorm"
)
//文章
type Article struct{
	gorm.Model
	Title string
	Cid int
	Desc string
	Content string
	Img string
}

//添加文章
func CreateArt(data *Article )int{
	err := db.Create(&data).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//查找文章
func GetArtInfo(id int)(Article,int){
	var article Article
	err := db.Where("id = ?",id).First(&article).Error
	if err!=nil{
		return article,errmsg.ERROR_ARTICLE_NOT_EXIST
	}
	return article,errmsg.SUCCSE
}

//编辑文章
func EditArt(id int,data *Article )int{
	var article Article
	var maps = make(map[string]interface{})
	maps["title"]=data.Title
	maps["cid"]=data.Cid
	maps["desc"]=data.Desc
	maps["content"]=data.Content
	maps["img"]=data.Img

	err := db.Model(&article).Where("id = ?",id).Updates(maps).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//删除文章
func DeleteArt(id int)int{
	var article Article
	err := db.Where("id = ?",id).Delete(&article).Error
	if err!=nil{
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}