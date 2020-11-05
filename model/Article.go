package model

import "gorm.io/gorm"
//文章
type Article struct{
	gorm.Model
	Title string
	Cid int
	Desc string
	Content string
	Img string
}