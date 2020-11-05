package model

import (
	"fmt"
	"ginblog/utils"
	"gorm.io/gorm/schema"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
)

var db *gorm.DB
var err error

func InitDb(){
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		utils.DbUser,
		utils.DbPassWord,
		utils.DbHost,
		utils.DbPort,
		utils.DbName,
	)
	db,err = gorm.Open(mysql.Open(dns),&gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		} ,
	})

	if err!=nil{
		fmt.Println("无法连接数据库,请检查参数:",err)
	}else{
		fmt.Println("已连接到数据库")
	}
	//迁移
	_ =db.AutoMigrate(&User{},&Article{})
	sqlDB,err2 := db.DB()
	if err2!=nil{
		fmt.Println("数据库类型错误",err2)
	}else{
		fmt.Println("迁移数据至库中")
	}
	sqlDB.SetMaxIdleConns(10)//连接池中最大闲置连接数10
	sqlDB.SetMaxOpenConns(100)//数据库最大连接数
	sqlDB.SetConnMaxLifetime(10*time.Second)//连接最大可复用时间
}
